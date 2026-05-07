package storage

import (
	"database/sql"
)

// PostgresStore implements PostgreSQL storage
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// GetDB returns the underlying database connection (for testing purposes)
func (s *PostgresStore) GetDB() *sql.DB {
	return s.db
}
