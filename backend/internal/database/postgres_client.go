package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PostgresConfig holds PostgreSQL connection configuration
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresClient creates a new PostgreSQL connection with connection pooling
func NewPostgresClient(cfg PostgresConfig) (*sql.DB, error) {
	// Build connection string
	connStr := buildConnString(cfg)

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection with ping
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database at %s:%s: %w", cfg.Host, cfg.Port, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)           // Maximum open connections
	db.SetMaxIdleConns(5)            // Maximum idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
	db.SetConnMaxIdleTime(1 * time.Minute) // Idle connection timeout

	return db, nil
}

// buildConnString constructs PostgreSQL connection string
func buildConnString(cfg PostgresConfig) string {
	if cfg.Password != "" {
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SSLMode,
	)
}
