package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
	"github.com/nullexp/limiter-x/internal/port/driven/db/repository"
)

type UserRateLimitRepositoryFactory struct{}

func NewUserRateLimitRepositoryFactory() *UserRateLimitRepositoryFactory {
	return &UserRateLimitRepositoryFactory{}
}

func (f *UserRateLimitRepositoryFactory) New(handler db.DbHandler) repository.UserRateLimitRepository {
	return NewUserRateLimitRepository(handler)
}

type UserRateLimitRepository struct {
	handler db.DbHandler
}

func NewUserRateLimitRepository(handler db.DbHandler) *UserRateLimitRepository {
	return &UserRateLimitRepository{handler: handler}
}

// CreateRateLimit inserts a new user rate limit record with an auto-generated ID
func (ur *UserRateLimitRepository) CreateRateLimit(ctx context.Context, rateLimit model.UserRateLimit) (string, error) {
	query := `
        INSERT INTO user_rate_limits (user_id, request_count, rate_limit, timestamp)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var id string
	err := ur.handler.QueryRowContext(ctx, query, rateLimit.UserId, rateLimit.RequestCount, rateLimit.RateLimit, rateLimit.Timestamp).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetRateLimitByUserId retrieves a user's rate limit by user ID
func (ur *UserRateLimitRepository) GetRateLimitByUserId(ctx context.Context, userId string) (*model.UserRateLimit, error) {
	query := `
        SELECT id, user_id, request_count, rate_limit, timestamp
        FROM user_rate_limits
        WHERE user_id = $1
    `

	var rateLimit model.UserRateLimit
	err := ur.handler.QueryRowContext(ctx, query, userId).Scan(&rateLimit.Id, &rateLimit.UserId, &rateLimit.RequestCount, &rateLimit.RateLimit, &rateLimit.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil if no records found
		}
		return nil, err
	}

	return &rateLimit, nil
}

// UpdateRateLimit updates an existing user's rate limit
func (ur *UserRateLimitRepository) UpdateRateLimit(ctx context.Context, rateLimit model.UserRateLimit) error {
	query := `
        UPDATE user_rate_limits
        SET request_count = $1, rate_limit = $2, timestamp = $3
        WHERE id = $4
    `
	_, err := ur.handler.ExecContext(ctx, query, rateLimit.RequestCount, rateLimit.RateLimit, rateLimit.Timestamp, rateLimit.Id)
	return err
}

// DeleteRateLimit removes a user's rate limit record
func (ur *UserRateLimitRepository) DeleteRateLimit(ctx context.Context, userId string) error {
	query := `
        DELETE FROM user_rate_limits
        WHERE user_id = $1
    `
	_, err := ur.handler.ExecContext(ctx, query, userId)
	return err
}

// UpdateUserRateLimit updates a user's customizable rate limit
func (ur *UserRateLimitRepository) UpdateUserRateLimit(ctx context.Context, userId string, newRateLimit int) error {
	query := `
        UPDATE user_rate_limits
        SET rate_limit = $1
        WHERE user_id = $2
    `
	_, err := ur.handler.ExecContext(ctx, query, newRateLimit, userId)
	return err
}
