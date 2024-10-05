package model

import (
	"context"
	"time"

	validator "github.com/go-playground/validator/v10"
)

type UserReadable struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	RoleId    string    `json:"role_id"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,gte=1"`
	Password string `json:"password" validate:"required,gte=1"`
	RoleId   string `json:"role_id" validate:"required,uuid"`
}

func (dto CreateUserRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type CreateUserResponse struct {
	Id string `json:"id"`
}
type GetUserByIdRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto GetUserByIdRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetUserByIdResponse struct {
	User UserReadable `json:"user"`
}
type GetAllUsersResponse struct {
	Users []UserReadable `json:"users"`
}
type UpdateUserRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required,gte=1"`
	RoleId   string `json:"role_id" validate:"required,uuid"`
}

func (dto UpdateUserRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type DeleteUserRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto DeleteUserRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetUserByUsernameAndPasswordRequest struct {
	Username string `json:"username" validate:"required,gte=1"`
	Password string `json:"password" validate:"required,gte=1"`
}

func (dto GetUserByUsernameAndPasswordRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetUsersWithPaginationRequest struct {
	Limit  int `json:"limit" validate:"gte=0"`
	Offset int `json:"offset" validate:"gte=0"`
}

func (dto GetUsersWithPaginationRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetUserByUsernameAndPasswordResponse struct {
	User UserReadable `json:"user"`
}

type GetUsersWithPaginationResponse struct {
	Users      []UserReadable `json:"users"`
	TotalCount int64          `json:"totalCount" validate:"gte=0"`
}
