package repository

import (
	"context"

	"github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
)

type UserRepository interface {
	CreateUser(context.Context, model.User) (string, error)
	GetUserById(context.Context, string) (*model.User, error)
	GetAllUsers(context.Context) ([]model.User, error)
	UpdateUser(context.Context, model.User) error
	DeleteUser(context.Context, string) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUsersWithPagination(ctx context.Context, offset, limit int) ([]model.User, error)
	Count(ctx context.Context) (int64, error)
}

type UserRepositoryFactory interface {
	New(db.DbHandler) UserRepository
}
