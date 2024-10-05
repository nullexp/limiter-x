package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
	"github.com/nullexp/limiter-x/internal/port/driven/db/repository"
)

type UserRateLimitRepositoryFactoryMock struct{}

func NewUserRateLimitRepositoryFactoryMock() *UserRateLimitRepositoryFactoryMock {
	return &UserRateLimitRepositoryFactoryMock{}
}

func (f *UserRateLimitRepositoryFactoryMock) New(handler db.DbHandler) repository.UserRateLimitRepository {
	return NewMockUserRateLimitRepository()
}

type MockUserRateLimitRepository struct {
	rateLimits map[string]model.UserRateLimit // Simulated in-memory database
}

func NewMockUserRateLimitRepository() *MockUserRateLimitRepository {
	return &MockUserRateLimitRepository{
		rateLimits: make(map[string]model.UserRateLimit),
	}
}

func (m *MockUserRateLimitRepository) CreateRateLimit(ctx context.Context, rateLimit model.UserRateLimit) (string, error) {
	id := uuid.New().String() // Generate UUID
	rateLimit.Id = id
	m.rateLimits[id] = rateLimit
	return id, nil
}

func (m *MockUserRateLimitRepository) GetRateLimitByUserId(ctx context.Context, userId string) (*model.UserRateLimit, error) {
	for _, rateLimit := range m.rateLimits {
		if rateLimit.UserId == userId {
			return &rateLimit, nil
		}
	}
	return nil, nil
}

func (m *MockUserRateLimitRepository) UpdateRateLimit(ctx context.Context, rateLimit model.UserRateLimit) error {
	if _, ok := m.rateLimits[rateLimit.Id]; !ok {
		return nil
	}
	m.rateLimits[rateLimit.Id] = rateLimit
	return nil
}

func (m *MockUserRateLimitRepository) DeleteRateLimit(ctx context.Context, userId string) error {
	for id, rateLimit := range m.rateLimits {
		if rateLimit.UserId == userId {
			delete(m.rateLimits, id)
			return nil
		}
	}
	return nil
}
