package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nullexp/limiter-x/internal/adapter/driven/db"
	"github.com/nullexp/limiter-x/internal/adapter/driven/db/repository"
	"github.com/nullexp/limiter-x/internal/domain/model"
	portModels "github.com/nullexp/limiter-x/internal/port/driver/model"
	"github.com/stretchr/testify/assert"
)

type MockPasswordService struct{}

func (m MockPasswordService) HashPassword(password string) (string, error) {
	return "hashed" + password, nil
}

func (m MockPasswordService) ComparePassword(hashedPassword, textPassword string) error {
	if hashedPassword == "hashed"+textPassword {
		return nil
	}
	return errors.New("password does not match")
}

func NewMockPasswordService() *MockPasswordService {
	return &MockPasswordService{}
}

func TestUserServiceCreateUser(t *testing.T) {
	uid := uuid.New().String()
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})

	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()
	request := portModels.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		RoleId:   uid,
	}

	response, err := us.CreateUser(ctx, request)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Id)

	createdUser, err := repo.GetUserById(ctx, response.Id)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", createdUser.Username)
	assert.Equal(t, "hashedpassword", createdUser.Password)
	assert.Equal(t, uid, createdUser.RoleId)
}

func TestUserServiceGetUserById(t *testing.T) {
	uid := uuid.New().String()
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})

	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()
	mockUser := model.User{
		Username:  "testuser",
		Password:  "hashedpassword",
		RoleId:    uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := repo.CreateUser(ctx, mockUser)
	assert.NoError(t, err)

	request := portModels.GetUserByIdRequest{Id: id}
	response, err := us.GetUserById(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.User.Username)
}

func TestUserServiceGetAllUsers(t *testing.T) {
	uid := uuid.New().String()
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})
	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()

	mockUsers := []model.User{
		{
			Username:  "testuser1",
			Password:  "hashedpassword",
			RoleId:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "testuser2",
			Password:  "hashedpassword",
			RoleId:    uid,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range mockUsers {
		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
	}

	response, err := us.GetAllUsers(ctx)
	assert.NoError(t, err)
	assert.Len(t, response.Users, 2)
}

func TestUserServiceUpdateUser(t *testing.T) {
	uid := uuid.New().String()
	uid2 := uuid.New().String()
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})

	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()
	mockUser := model.User{
		Username:  "testuser",
		Password:  "hashedpassword",
		RoleId:    uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := repo.CreateUser(ctx, mockUser)
	assert.NoError(t, err)

	request := portModels.UpdateUserRequest{
		Id:       id,
		Password: "newpassword",
		RoleId:   uid2,
	}

	err = us.UpdateUser(ctx, request)
	assert.NoError(t, err)

	updatedUser, err := repo.GetUserById(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, uid2, updatedUser.RoleId)
	assert.Equal(t, "hashednewpassword", updatedUser.Password)
}

func TestUserServiceDeleteUser(t *testing.T) {
	uid := uuid.New().String()
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})
	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()
	mockUser := model.User{
		Username:  "testuser",
		Password:  "hashedpassword",
		RoleId:    uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := repo.CreateUser(ctx, mockUser)
	assert.NoError(t, err)

	request := portModels.DeleteUserRequest{Id: id}
	err = us.DeleteUser(ctx, request)
	assert.NoError(t, err)

	deletedUser, err := repo.GetUserById(ctx, id)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserServiceGetUserByUsernameAndPassword(t *testing.T) {
	uid := uuid.New().String()
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})

	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()
	mockUser := model.User{
		Username:  "testuser",
		Password:  "hashedpassword", // This should match the logic in MockPasswordService
		RoleId:    uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := repo.CreateUser(ctx, mockUser)
	assert.NoError(t, err)

	request := portModels.GetUserByUsernameAndPasswordRequest{
		Username: "testuser",
		Password: "password", // This should match the logic in MockPasswordService
	}

	response, err := us.GetUserByUsernameAndPassword(ctx, request)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.User.Username)
}

func TestUserServiceGetUsersWithPagination(t *testing.T) {
	mockUserRepositoryFactory := repository.NewUserRepositoryFactoryMock()
	mockPasswordService := NewMockPasswordService()
	transactionFactory := db.NewPostgresTransactionFactoryMock()
	repo := mockUserRepositoryFactory.New(&db.MockDbHandler{})
	us := NewUserService(mockUserRepositoryFactory, mockPasswordService, transactionFactory)

	ctx := context.Background()

	mockUsers := []model.User{
		{
			Username:  "testuser1",
			Password:  "hashedpassword",
			RoleId:    "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "testuser2",
			Password:  "hashedpassword",
			RoleId:    "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range mockUsers {
		_, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
	}

	request := portModels.GetUsersWithPaginationRequest{
		Offset: 0,
		Limit:  2,
	}

	response, err := us.GetUsersWithPagination(ctx, request)
	assert.NoError(t, err)
	assert.Len(t, response.Users, 2)
}
