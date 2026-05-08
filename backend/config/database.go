package config

import (
	"database/sql"
	"log"
	"nowyouseeme/internal/database"
	"os"
)

// DBConfig holds database connection parameters
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDBConfig loads database config from environment
// psql -h localhost -p 5432 -U liuzhenhua -d nowyouseeme
func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "liuzhenhua"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "nowyouseeme"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// ConnectDB establishes PostgreSQL connection using the internal database client
func ConnectDB(config *DBConfig) (*sql.DB, error) {
	// Convert to internal database config format
	dbConfig := database.PostgresConfig{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		DBName:   config.DBName,
		SSLMode:  config.SSLMode,
	}

	// Use the new PostgreSQL client
	db, err := database.NewPostgresClient(dbConfig)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to PostgreSQL database: %s", config.DBName)
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
