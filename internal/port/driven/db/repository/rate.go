package repository

import (
	"context"

	"github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
)

type UserRateLimitRepository interface {
	CreateRateLimit(context.Context, model.UserRateLimit) (string, error)
	GetRateLimitByUserId(context.Context, string) (*model.UserRateLimit, error)
	UpdateRateLimit(context.Context, model.UserRateLimit) error
	DeleteRateLimit(context.Context, string) error
}

type UserRateLimitRepositoryFactory interface {
	New(db.DbHandler) UserRateLimitRepository
}
