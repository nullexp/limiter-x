package grpc

import (
	"context"

	userv1 "github.com/nullexp/limiter-x/internal/adapter/driver/grpc/proto/user/v1"
	"github.com/nullexp/limiter-x/internal/port/driver/model"
	driverService "github.com/nullexp/limiter-x/internal/port/driver/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userv1.UnimplementedUserServiceServer
	service driverService.UserService
}

func NewUserService(us driverService.UserService) *UserService {
	return &UserService{service: us}
}

func (us UserService) CreateUser(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	rs, err := us.service.CreateUser(ctx, model.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
		RoleId:   request.RoleId,
	})
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{Id: rs.Id}, nil
}

func castReadableToGrpcUser(readable model.UserReadable) *userv1.User {
	return &userv1.User{
		Id:        readable.Id,
		Username:  readable.Username,
		RoleId:    readable.RoleId,
		IsAdmin:   readable.IsAdmin,
		CreatedAt: readable.CreatedAt.String(),
		UpdatedAt: readable.UpdatedAt.String(),
	}
}

func (us UserService) GetUserById(ctx context.Context, request *userv1.GetUserByIdRequest) (*userv1.GetUserByIdResponse, error) {
	rs, err := us.service.GetUserById(ctx, model.GetUserByIdRequest{Id: request.Id})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetUserById failed : %v", err)
	}

	return &userv1.GetUserByIdResponse{User: castReadableToGrpcUser(rs.User)}, nil
}

func (us *UserService) GetAllUsers(ctx context.Context, request *userv1.GetAllUsersRequest) (*userv1.GetAllUsersResponse, error) {
	rs, err := us.service.GetAllUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetAllUsers failed: %v", err)
	}

	var users []*userv1.User
	for _, user := range rs.Users {
		userv1User := castReadableToGrpcUser(user)
		users = append(users, userv1User)
	}

	return &userv1.GetAllUsersResponse{Users: users}, nil
}

func (us *UserService) UpdateUser(ctx context.Context, request *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	err := us.service.UpdateUser(ctx, model.UpdateUserRequest{
		Id:       request.Id,
		Password: request.Password,
		RoleId:   request.RoleId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UpdateUser failed: %v", err)
	}

	return &userv1.UpdateUserResponse{}, nil
}

func (us *UserService) DeleteUser(ctx context.Context, request *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	err := us.service.DeleteUser(ctx, model.DeleteUserRequest{Id: request.Id})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteUser failed: %v", err)
	}

	return &userv1.DeleteUserResponse{}, nil
}

func (us *UserService) GetUserByUsernameAndPassword(ctx context.Context, request *userv1.GetUserByUsernameAndPasswordRequest) (*userv1.GetUserByUsernameAndPasswordResponse, error) {
	rs, err := us.service.GetUserByUsernameAndPassword(ctx, model.GetUserByUsernameAndPasswordRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetUserByUsernameAndPassword failed: %v", err)
	}

	user := castReadableToGrpcUser(rs.User)

	return &userv1.GetUserByUsernameAndPasswordResponse{User: user}, nil
}

func (us *UserService) GetUsersWithPagination(ctx context.Context, request *userv1.GetUsersWithPaginationRequest) (*userv1.GetUsersWithPaginationResponse, error) {
	rs, err := us.service.GetUsersWithPagination(ctx, model.GetUsersWithPaginationRequest{
		Limit:  int(request.Limit),
		Offset: int(request.Offset),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "GetUsersWithPagination failed: %v", err)
	}

	var users []*userv1.User
	for _, user := range rs.Users {
		userv1User := castReadableToGrpcUser(user)
		users = append(users, userv1User)
	}

	return &userv1.GetUsersWithPaginationResponse{Users: users, TotalCount: rs.TotalCount}, nil
}
