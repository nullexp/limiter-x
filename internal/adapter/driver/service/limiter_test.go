package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nullexp/limiter-x/internal/adapter/driven/cache"
	"github.com/nullexp/limiter-x/internal/adapter/driven/db"
	"github.com/nullexp/limiter-x/internal/adapter/driven/db/repository"
	"github.com/nullexp/limiter-x/internal/domain/model"
	portRepo "github.com/nullexp/limiter-x/internal/port/driven/db/repository"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitService_RateLimit(t *testing.T) {
	cache := cache.NewMemoryClient(time.Hour, time.Hour)
	repoFactory := repository.NewUserRateLimitRepositoryFactoryMock()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	window := time.Second * 10 // Set your window time here

	service := NewRateLimitService(repoFactory, cache, transactionFactory, window)

	tests := []struct {
		name       string
		userId     string
		limit      int
		setupRepo  func(userId string) portRepo.UserRateLimitRepository
		setupCache func(userId string)
		clean      func()
		expect     bool
	}{
		{
			name:   "First request, should succeed",
			userId: uuid.New().String(),
			limit:  5,
			setupRepo: func(userId string) portRepo.UserRateLimitRepository {
				transaction := transactionFactory.NewTransaction()
				handler, err := transaction.Begin(context.Background())
				assert.Nil(t, err)
				repo := repoFactory.New(handler)
				repo.CreateRateLimit(context.Background(), model.UserRateLimit{
					UserId:       userId, // Use userId from the test case
					RequestCount: 0,
					RateLimit:    100,  // This is the DB rate limit
					Timestamp:    time.Now(),
				})
				err = transaction.Commit(context.Background())
				assert.Nil(t, err)
				return repo
			},
			setupCache: func(userId string) {
				cache.Connect()
			},
			clean: func() {
				cache.Disconnect()
			},
			expect: true,
		},
		{
			name:   "Rate limit from DB when limit is 0",
			userId: uuid.New().String(),
			limit:  0, // This should cause the service to use the DB-stored rate limit
			setupRepo: func(userId string) portRepo.UserRateLimitRepository {
				transaction := transactionFactory.NewTransaction()
				handler, err := transaction.Begin(context.Background())
				assert.Nil(t, err)
				repo := repoFactory.New(handler)
				repo.CreateRateLimit(context.Background(), model.UserRateLimit{
					UserId:       userId,
					RequestCount: 0,
					RateLimit:    10, // DB-stored rate limit
					Timestamp:    time.Now(),
				})
				err = transaction.Commit(context.Background())
				assert.Nil(t, err)
				return repo
			},
			setupCache: func(userId string) {
				cache.Connect()
			},
			clean: func() {
				cache.Disconnect()
			},
			expect: true,
		},
		{
			name:   "Rate limit from cache when limit is 0",
			userId: uuid.New().String(),
			limit:  0, // Should use the rate limit from the cache
			setupRepo: func(userId string) portRepo.UserRateLimitRepository {
				transaction := transactionFactory.NewTransaction()
				handler, err := transaction.Begin(context.Background())
				assert.Nil(t, err)
				repo := repoFactory.New(handler)
				repo.CreateRateLimit(context.Background(), model.UserRateLimit{
					UserId:       userId,
					RequestCount: 0,
					RateLimit:    10, // DB-stored rate limit
					Timestamp:    time.Now(),
				})
				err = transaction.Commit(context.Background())
				assert.Nil(t, err)
				return repo
			},
			setupCache: func(userId string) {
				cache.Connect()
				// Set a rate limit of 8 in cache
				err := cache.Set(context.Background(), userId, []byte(`{"UserId":"`+userId+`","RequestCount":0,"RateLimit":8,"Timestamp":"`+time.Now().Format(time.RFC3339)+`"}`), time.Hour)
				assert.Nil(t, err)
			},
			clean: func() {
				cache.Disconnect()
			},
			expect: true,
		},
		{
			name:   "Override DB limit with passed limit",
			userId: uuid.New().String(),
			limit:  15, // Should override the DB-stored rate limit
			setupRepo: func(userId string) portRepo.UserRateLimitRepository {
				transaction := transactionFactory.NewTransaction()
				handler, err := transaction.Begin(context.Background())
				assert.Nil(t, err)
				repo := repoFactory.New(handler)
				repo.CreateRateLimit(context.Background(), model.UserRateLimit{
					UserId:       userId,
					RequestCount: 0,
					RateLimit:    10, // DB-stored rate limit
					Timestamp:    time.Now(),
				})
				err = transaction.Commit(context.Background())
				assert.Nil(t, err)
				return repo
			},
			setupCache: func(userId string) {
				cache.Connect()
			},
			clean: func() {
				cache.Disconnect()
			},
			expect: true,
		},
		{
			name:   "Rate limit reached, should fail",
			userId: uuid.New().String(),
			limit:  2, // Passed limit should be used
			setupRepo: func(userId string) portRepo.UserRateLimitRepository {
				transaction := transactionFactory.NewTransaction()
				handler, err := transaction.Begin(context.Background())
				assert.Nil(t, err)

				repo := repoFactory.New(handler)
				repo.CreateRateLimit(context.Background(), model.UserRateLimit{
					UserId:       userId,
					RequestCount: 2,
					RateLimit:    5,
					Timestamp:    time.Now(),
				})
				err = transaction.Commit(context.Background())
				assert.Nil(t, err)
				return repo
			},
			setupCache: func(userId string) {
				cache.Connect()
				err := cache.Set(context.Background(), userId, []byte(`{"UserId":"`+userId+`","RequestCount":2,"Timestamp":"`+time.Now().Format(time.RFC3339)+`"}`), time.Hour)
				assert.Nil(t, err)
			},
			clean: func() {
				cache.Disconnect()
			},
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup repository and cache for each test

			test.setupCache(test.userId)

			// Call the rate limit service method
			allowed, err := service.RateLimit(context.Background(), test.userId, test.limit)
			assert.Nil(t, err)

			// Check the expected outcome
			assert.Equal(t, test.expect, allowed)

			test.clean()
		})
	}
}

func setupBenchmark(b *testing.B) (*RateLimitService, func(userId string) error, func(userId string) error, func() error) {
	cache := cache.NewMemoryClient(time.Hour, time.Hour)
	repoFactory := repository.NewUserRateLimitRepositoryFactoryMock()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	window := time.Second * 10 // Set your window time here

	service := NewRateLimitService(repoFactory, cache, transactionFactory, window)

	// Setup repository for each user
	setupRepo := func(userId string) error {
		transaction := transactionFactory.NewTransaction()
		handler, err := transaction.Begin(context.Background())
		if err != nil {
			return err
		}
		repo := repoFactory.New(handler)
		_, err = repo.CreateRateLimit(context.Background(), model.UserRateLimit{
			UserId:       userId,
			RequestCount: 0,
			Timestamp:    time.Now(),
		})
		if err != nil {
			return err
		}
		err = transaction.Commit(context.Background())
		return err
	}

	// Setup cache connection and any initial data
	setupCache := func(userId string) error {
		err := cache.Connect()
		if err != nil {
			return err
		}
		return nil
	}

	// Clean up by disconnecting from the cache
	clean := func() error {
		err := cache.Disconnect()
		return err
	}

	return service, setupRepo, setupCache, clean
}

func BenchmarkRateLimit_CacheHit(b *testing.B) {
	service, setupRepo, setupCache, clean := setupBenchmark(b)
	defer clean()

	userId := uuid.New().String()
	limit := 5

	// Set up repository and cache
	setupRepo(userId)
	setupCache(userId)

	// Simulate a cache hit scenario by setting data in the cache
	err := service.cache.Set(context.Background(), userId, []byte(`{"UserId":"`+userId+`","RequestCount":0,"Timestamp":"`+time.Now().Format(time.RFC3339)+`"}`), time.Hour)
	if err != nil {
		b.Fatalf("failed to set cache: %v", err)
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.RateLimit(context.Background(), userId, limit)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkRateLimit_CacheMiss(b *testing.B) {
	service, setupRepo, setupCache, clean := setupBenchmark(b)
	defer clean()

	userId := uuid.New().String()
	limit := 5

	// Set up repository but simulate a cache miss by not setting cache data
	setupRepo(userId)
	setupCache(userId)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.RateLimit(context.Background(), userId, limit)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkRateLimit_LimitReached(b *testing.B) {
	service, setupRepo, setupCache, clean := setupBenchmark(b)
	defer clean()

	userId := uuid.New().String()
	limit := 2

	// Set up repository and simulate reaching the limit
	setupRepo(userId)
	setupCache(userId)

	err := service.cache.Set(context.Background(), userId, []byte(`{"UserId":"`+userId+`","RequestCount":2,"Timestamp":"`+time.Now().Format(time.RFC3339)+`"}`), time.Hour)
	if err != nil {
		b.Fatalf("failed to set cache: %v", err)
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.RateLimit(context.Background(), userId, limit)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
