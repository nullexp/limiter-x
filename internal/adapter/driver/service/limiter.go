package service

import (
	"context"
	"encoding/json"
	"time"

	domainModel "github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
	"github.com/nullexp/limiter-x/internal/port/driven/db/repository"
	"github.com/nullexp/limiter-x/internal/port/driver/service"
	"github.com/pkg/errors"
)

type RateLimitService struct {
	repoFactory          repository.UserRateLimitRepositoryFactory
	cache                driven.Cache
	dbTransactionFactory db.DbTransactionFactory
	window               time.Duration
}

// NewRateLimitService creates a new instance of RateLimitService.
func NewRateLimitService(repo repository.UserRateLimitRepositoryFactory, cache driven.Cache, dbTransactionFactory db.DbTransactionFactory, window time.Duration) *RateLimitService {
	return &RateLimitService{
		repoFactory:          repo,
		cache:                cache,
		dbTransactionFactory: dbTransactionFactory,
		window:               window,
	}
}

// RateLimit checks if the request is allowed for the user within the defined rate limit.
func (rls *RateLimitService) RateLimit(ctx context.Context, userId string, limit int) (bool, error) {
	// Start the transaction
	tx := rls.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	allowed, err := rls.rateLimit(ctx, transaction, userId, limit)
	if err != nil {
		return false, err
	}

	// Commit the transaction if everything went well
	err = tx.Commit(ctx)
	return allowed, err
}

// rateLimit handles the rate limiting logic within a transaction.
func (rls *RateLimitService) rateLimit(ctx context.Context, tx db.DbHandler, userId string, limit int) (bool, error) {
	repo := rls.repoFactory.New(tx)

	// Try to fetch from cache first
	cachedRateLimitData, err := rls.cache.Fetch(ctx, userId)
	if err == nil && cachedRateLimitData != nil {
		var cachedRateLimit domainModel.UserRateLimit
		if err := json.Unmarshal(cachedRateLimitData, &cachedRateLimit); err != nil {
			return false, errors.Wrap(err, "failed to unmarshal cached rate limit")
		}

		// Determine which limit to use (parameter or database value)
		effectiveLimit := limit
		if limit == 0 {
			effectiveLimit = cachedRateLimit.RateLimit
		}

		// Cache hit, validate against effective limit
		if cachedRateLimit.RequestCount >= effectiveLimit {
			// Deny request if the request count has reached or exceeded the limit
			return false, nil
		}

		// Increment the request count and update cache
		cachedRateLimit.RequestCount++
		if err := rls.updateCache(ctx, userId, &cachedRateLimit); err != nil {
			return false, err
		}
		return true, nil
	}

	// Cache miss or request count exceeded, fallback to repository
	rateLimit, err := repo.GetRateLimitByUserId(ctx, userId)
	if err != nil {
		return false, errors.Wrap(err, "failed to get user rate limit from repo")
	}

	now := time.Now()
	if rateLimit == nil {
		// User has no existing rate limit, create new one with default or parameter limit
		effectiveLimit := limit
		if limit == 0 {
			effectiveLimit = 100 // Default rate limit if no record and limit is 0
		}
		rateLimit = &domainModel.UserRateLimit{
			UserId:       userId,
			RequestCount: 1,
			RateLimit:    effectiveLimit, // Store the effective limit in the database
			Timestamp:    now,
		}
		_, err = repo.CreateRateLimit(ctx, *rateLimit)
		if err != nil {
			return false, errors.Wrap(err, "failed to create user rate limit in repo")
		}

		// Store new rate limit in cache
		if err := rls.setCache(ctx, userId, rateLimit); err != nil {
			return false, err
		}
		return true, nil
	}

	// Determine which limit to use (parameter or database value)
	effectiveLimit := limit
	if limit == 0 {
		effectiveLimit = rateLimit.RateLimit
	}

	// Check if request window has passed
	if now.Sub(rateLimit.Timestamp) > rls.window {
		// Reset rate limit if outside the window
		rateLimit.RequestCount = 1
		rateLimit.Timestamp = now
		if err := repo.UpdateRateLimit(ctx, *rateLimit); err != nil {
			return false, errors.Wrap(err, "failed to update user rate limit")
		}
		// Update the cache with the reset data
		if err := rls.setCache(ctx, userId, rateLimit); err != nil {
			return false, err
		}
		return true, nil
	}

	// If within window, check the request count against effective limit
	if rateLimit.RequestCount >= effectiveLimit {
		// Deny the request if the count is equal to or exceeds the limit
		return false, nil
	}

	// Increment the request count and update repository
	rateLimit.RequestCount++
	if err := repo.UpdateRateLimit(ctx, *rateLimit); err != nil {
		return false, errors.Wrap(err, "failed to update user rate limit in repo")
	}
	// Update the cache with the incremented count
	if err := rls.updateCache(ctx, userId, rateLimit); err != nil {
		return false, err
	}
	return true, nil
}

// setCache stores the user rate limit in the cache.
func (rls *RateLimitService) setCache(ctx context.Context, userId string, rateLimit *domainModel.UserRateLimit) error {
	data, err := json.Marshal(rateLimit)
	if err != nil {
		return errors.Wrap(err, "failed to marshal rate limit")
	}
	return rls.cache.Set(ctx, userId, data, rls.window)
}

// updateCache updates the user rate limit in the cache.
func (rls *RateLimitService) updateCache(ctx context.Context, userId string, rateLimit *domainModel.UserRateLimit) error {
	data, err := json.Marshal(rateLimit)
	if err != nil {
		return errors.Wrap(err, "failed to marshal rate limit")
	}
	return rls.cache.Set(ctx, userId, data, rls.window)
}

func (rls *RateLimitService) GetUserRateLimit(ctx context.Context, userId string) (*service.RateLimitModel, error) {
	// Try to fetch from cache first
	cachedRateLimitData, err := rls.cache.Fetch(ctx, userId)
	if err == nil && cachedRateLimitData != nil {
		var cachedRateLimit domainModel.UserRateLimit
		if err := json.Unmarshal(cachedRateLimitData, &cachedRateLimit); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cached rate limit")
		}

		// Return the cached rate limit as a RateLimitModel
		return &service.RateLimitModel{
			UserId:    cachedRateLimit.UserId,
			Limit:     cachedRateLimit.RateLimit,    // Limit in this case refers to the request count in the cache
			Remaining: cachedRateLimit.RequestCount, // Example: returning the same count for simplicity
			Window:    rls.window.String(),
		}, nil
	}

	// Cache miss, fallback to repository
	tx := rls.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repo := rls.repoFactory.New(transaction)
	rateLimit, err := repo.GetRateLimitByUserId(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rate limit from repository")
	}
	if rateLimit == nil {
		return nil, errors.New("rate limit not found")
	}

	// Return the rate limit fetched from the repository
	return &service.RateLimitModel{
		UserId:    rateLimit.UserId,
		Limit:     rateLimit.RateLimit,
		Remaining: rateLimit.RequestCount, // You can adjust how you calculate the remaining requests
		Window:    rls.window.String(),
	}, nil
}

// UpdateUserRateLimit updates the rate limit for a specific user
func (rls *RateLimitService) UpdateUserRateLimit(ctx context.Context, userId string, newLimit int) error {
	// Begin a transaction
	tx := rls.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repo := rls.repoFactory.New(transaction)

	// Fetch the current rate limit from the repository
	rateLimit, err := repo.GetRateLimitByUserId(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "failed to get user rate limit from repository")
	}

	// If no rate limit exists for the user, create a new one
	if rateLimit == nil {
		rateLimit = &domainModel.UserRateLimit{
			UserId:       userId,
			RequestCount: newLimit, // Set the new limit
			Timestamp:    time.Now(),
		}
		_, err = repo.CreateRateLimit(ctx, *rateLimit)
		if err != nil {
			return errors.Wrap(err, "failed to create user rate limit in repository")
		}
	} else {
		// Update the rate limit in the repository
		rateLimit.RequestCount = newLimit
		err = repo.UpdateRateLimit(ctx, *rateLimit)
		if err != nil {
			return errors.Wrap(err, "failed to update user rate limit in repository")
		}
	}

	// Commit the transaction if successful
	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	// Update the cache with the new limit
	if err := rls.setCache(ctx, userId, rateLimit); err != nil {
		return errors.Wrap(err, "failed to update cache with new rate limit")
	}

	return nil
}
