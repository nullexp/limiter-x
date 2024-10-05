package db

import (
	"context"

	"github.com/nullexp/limiter-x/internal/port/driven/db"
)

// PostgresTransactionMock is a mock implementation of the DbTransaction interface
type PostgresTransactionMock struct{}

func (m PostgresTransactionMock) Begin(ctx context.Context) (db.DbHandler, error) {
	// Simulate returning a handler (could be a mock handler)
	return &MockDbHandler{}, nil
}

func (m PostgresTransactionMock) Commit(ctx context.Context) error {
	return nil
}

func (m PostgresTransactionMock) Rollback(ctx context.Context) error {
	return nil
}

func (m PostgresTransactionMock) RollbackUnlessCommitted(ctx context.Context) {
	// DO nothing
}

// PostgresTransactionFactoryMock is a mock implementation of the DbTransactionFactory interface
type PostgresTransactionFactoryMock struct{}

func (m PostgresTransactionFactoryMock) NewTransaction() db.DbTransaction {
	return PostgresTransactionMock{}
}

func NewPostgresTransactionFactoryMock() PostgresTransactionFactoryMock {
	return PostgresTransactionFactoryMock{}
}
