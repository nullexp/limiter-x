package service

import (
	"context"

	domainError "github.com/nullexp/limiter-x/internal/domain/error"
	domainModel "github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
	"github.com/nullexp/limiter-x/internal/port/driven/db/repository"
	"github.com/nullexp/limiter-x/internal/port/driver/model"
	"github.com/pkg/errors"
)

type UserService struct {
	userRepositoryFactory repository.UserRepositoryFactory
	passwordService       driven.PasswordService
	dbTransactionFactory  db.DbTransactionFactory
}

func NewUserService(userRepositoryFactory repository.UserRepositoryFactory, passwordService driven.PasswordService, dbTransactionFactory db.DbTransactionFactory) *UserService {
	return &UserService{userRepositoryFactory: userRepositoryFactory, passwordService: passwordService, dbTransactionFactory: dbTransactionFactory}
}

func (us UserService) CreateUser(ctx context.Context, request model.CreateUserRequest) (out *model.CreateUserResponse, err error) {
	if err = request.Validate(ctx); err != nil {
		return
	}

	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	out, err = us.createUser(ctx, transaction, request)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) createUser(ctx context.Context, tx db.DbHandler, request model.CreateUserRequest) (*model.CreateUserResponse, error) {
	ps, err := us.passwordService.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	userRepo := us.userRepositoryFactory.New(tx)

	id, err := userRepo.CreateUser(ctx, domainModel.User{
		Username: request.Username,
		RoleId:   request.RoleId,
		IsAdmin:  false,
		Password: ps,
	})
	if err != nil {
		return nil, err
	}

	return &model.CreateUserResponse{Id: id}, nil
}

func castUserToReadable(user *domainModel.User) model.UserReadable {
	return model.UserReadable{
		Id:        user.Id,
		Username:  user.Username,
		RoleId:    user.RoleId,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (us UserService) GetUserById(ctx context.Context, request model.GetUserByIdRequest) (out *model.GetUserByIdResponse, err error) {
	if err = request.Validate(ctx); err != nil {
		return
	}

	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	out, err = us.getUserById(ctx, transaction, request)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) getUserById(ctx context.Context, tx db.DbHandler, request model.GetUserByIdRequest) (*model.GetUserByIdResponse, error) {
	userRepo := us.userRepositoryFactory.New(tx)
	user, err := userRepo.GetUserById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domainError.ErrUserNotFound
	}

	return &model.GetUserByIdResponse{
		User: castUserToReadable(user),
	}, nil
}

func (us UserService) GetAllUsers(ctx context.Context) (out *model.GetAllUsersResponse, err error) {
	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	out, err = us.getAllUsers(ctx, transaction)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) getAllUsers(ctx context.Context, tx db.DbHandler) (*model.GetAllUsersResponse, error) {
	repo := us.userRepositoryFactory.New(tx)
	users, err := repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	out := []model.UserReadable{}

	for _, user := range users {
		out = append(out, castUserToReadable(&user))
	}
	return &model.GetAllUsersResponse{
		Users: out,
	}, nil
}

func (us UserService) UpdateUser(ctx context.Context, request model.UpdateUserRequest) (err error) {
	if err = request.Validate(ctx); err != nil {
		return err
	}

	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.RollbackUnlessCommitted(ctx)

	err = us.updateUser(ctx, transaction, request)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) updateUser(ctx context.Context, tx db.DbHandler, request model.UpdateUserRequest) error {
	repo := us.userRepositoryFactory.New(tx)
	user, err := repo.GetUserById(ctx, request.Id)
	if err != nil {
		return err
	}

	if user == nil {
		return domainError.ErrUserNotFound
	}

	user.RoleId = request.RoleId

	ps, err := us.passwordService.HashPassword(request.Password)
	if err != nil {
		return err
	}

	user.Password = ps

	return repo.UpdateUser(ctx, *user)
}

func (us UserService) DeleteUser(ctx context.Context, request model.DeleteUserRequest) (err error) {
	if err = request.Validate(ctx); err != nil {
		return
	}

	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.RollbackUnlessCommitted(ctx)

	err = us.deleteUser(ctx, transaction, request)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) deleteUser(ctx context.Context, tx db.DbHandler, request model.DeleteUserRequest) error {
	repo := us.userRepositoryFactory.New(tx)
	user, err := repo.GetUserById(ctx, request.Id)
	if err != nil {
		return err
	}

	if user == nil {
		return domainError.ErrUserNotFound
	}

	if user.IsAdmin {
		return domainError.ErrAdminCantBeRemoved
	}

	return repo.DeleteUser(ctx, request.Id)
}

func (us UserService) GetUserByUsernameAndPassword(ctx context.Context, request model.GetUserByUsernameAndPasswordRequest) (out *model.GetUserByUsernameAndPasswordResponse, err error) {
	if err = request.Validate(ctx); err != nil {
		return
	}

	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.RollbackUnlessCommitted(ctx)

	out, err = us.getUserByUsernameAndPassword(ctx, transaction, request)
	if err != nil {
		txErr := tx.Rollback(ctx)
		if err != nil {
			err = errors.Wrap(err, txErr.Error())
			return
		}
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) getUserByUsernameAndPassword(ctx context.Context, tx db.DbHandler, request model.GetUserByUsernameAndPasswordRequest) (*model.GetUserByUsernameAndPasswordResponse, error) {
	repo := us.userRepositoryFactory.New(tx)

	user, err := repo.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, domainError.ErrUserNotFound
	}

	err = us.passwordService.ComparePassword(user.Password, request.Password)
	if err != nil {
		// For security, we do not provide more info
		return nil, domainError.ErrUserNotFound
	}

	return &model.GetUserByUsernameAndPasswordResponse{User: castUserToReadable(user)}, nil
}

func (us UserService) GetUsersWithPagination(ctx context.Context, request model.GetUsersWithPaginationRequest) (out *model.GetUsersWithPaginationResponse, err error) {
	if err = request.Validate(ctx); err != nil {
		return nil, err
	}

	tx := us.dbTransactionFactory.NewTransaction()
	transaction, err := tx.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.RollbackUnlessCommitted(ctx)

	out, err = us.getUsersWithPagination(ctx, transaction, request)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	return
}

func (us UserService) getUsersWithPagination(ctx context.Context, tx db.DbHandler, request model.GetUsersWithPaginationRequest) (*model.GetUsersWithPaginationResponse, error) {
	repo := us.userRepositoryFactory.New(tx)
	users, err := repo.GetUsersWithPagination(ctx, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}

	count, err := repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	out := []model.UserReadable{}
	for _, user := range users {
		out = append(out, castUserToReadable(&user))
	}
	return &model.GetUsersWithPaginationResponse{
		Users:      out,
		TotalCount: count,
	}, nil
}
