package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nullexp/limiter-x/internal/domain/model"
	"github.com/nullexp/limiter-x/internal/port/driven/db"
	"github.com/nullexp/limiter-x/internal/port/driven/db/repository"
)

type UserRepositoryFactory struct{}

func NewUserRepositoryFactory() *UserRepositoryFactory {
	return &UserRepositoryFactory{}
}

func (f *UserRepositoryFactory) New(handler db.DbHandler) repository.UserRepository {
	return NewUserRepository(handler)
}

type UserRepository struct {
	handler db.DbHandler
}

func NewUserRepository(handler db.DbHandler) *UserRepository {
	return &UserRepository{handler: handler}
}

func (ur UserRepository) CreateUser(ctx context.Context, user model.User) (string, error) {
	query := `
        INSERT INTO users (username, password, role_id, is_admin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id
    `
	var id string
	err := ur.handler.QueryRowContext(ctx, query, user.Username, user.Password, user.RoleId, user.IsAdmin).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (ur UserRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	query := `
	SELECT id, username, password, role_id, is_admin, created_at, updated_at
	FROM users
	WHERE id = $1
`
	var user model.User
	err := ur.handler.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Password, &user.RoleId, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil when no rows are found
		}
		return nil, err
	}
	return &user, nil
}

func (ur UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	query := `
	SELECT id, username, password, role_id, is_admin, created_at, updated_at
	FROM users
`
	rows, err := ur.handler.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []model.User{}, nil // Return nil when no rows are found
		}
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.RoleId, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur UserRepository) UpdateUser(ctx context.Context, user model.User) error {
	query := `
	UPDATE users
	SET username = $1, password = $2, role_id = $3, is_admin = $4, updated_at = NOW()
	WHERE id = $5
`
	_, err := ur.handler.ExecContext(ctx, query, user.Username, user.Password, user.RoleId, user.IsAdmin, user.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // Return nil when no rows are found
		}
		return err
	}
	return nil
}

func (ur UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
	DELETE FROM users
	WHERE id = $1
`
	_, err := ur.handler.ExecContext(ctx, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // Return nil when no rows are found
		}
		return err
	}
	return nil
}

func (ur UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
	SELECT id, username, password, role_id, is_admin, created_at, updated_at
	FROM users
	WHERE username = $1 
`
	var user model.User
	err := ur.handler.QueryRowContext(ctx, query, username).Scan(&user.Id, &user.Username, &user.Password, &user.RoleId, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil when no rows are found
		}
		return &model.User{}, err
	}
	return &user, nil
}

func (ur UserRepository) GetUsersWithPagination(ctx context.Context, offset, limit int) ([]model.User, error) {
	query := `
        SELECT id, username, password, role_id, is_admin, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        OFFSET $1 LIMIT $2
    `
	rows, err := ur.handler.QueryContext(ctx, query, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil when no rows are found
		}
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.RoleId, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur UserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := ur.handler.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
