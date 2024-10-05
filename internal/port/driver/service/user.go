package driver

import (
	"context"

	"github.com/nullexp/limiter-x/internal/port/driver/model"
)

type UserService interface {
	CreateUser(context.Context, model.CreateUserRequest) (*model.CreateUserResponse, error)
	GetUserById(context.Context, model.GetUserByIdRequest) (*model.GetUserByIdResponse, error)
	GetAllUsers(context.Context) (*model.GetAllUsersResponse, error)
	UpdateUser(context.Context, model.UpdateUserRequest) error
	DeleteUser(context.Context, model.DeleteUserRequest) error
	GetUserByUsernameAndPassword(context.Context, model.GetUserByUsernameAndPasswordRequest) (*model.GetUserByUsernameAndPasswordResponse, error)
	GetUsersWithPagination(context.Context, model.GetUsersWithPaginationRequest) (*model.GetUsersWithPaginationResponse, error)
}
