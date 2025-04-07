// Package database handles database connections and schema management
package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

// DB is the global database connection pool
var DB *sql.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DefaultConfig returns the default database configuration
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "081358226221Ra",
		DBName:   "medicine_reminder",
		SSLMode:  "disable",
	}
}

// BuildConnectionString creates a PostgreSQL connection string from config
func BuildConnectionString(config Config) string {
	password := url.QueryEscape(config.Password)
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, password, config.Host, config.Port, config.DBName, config.SSLMode)
}

// InitDB initializes the database connection and creates required tables
func InitDB() {
	config := DefaultConfig()
	connStr := BuildConnectionString(config)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Create medicines table
	err = createMedicinesTable()
	if err != nil {
		log.Fatalf("Error creating medicines table: %v", err)
	}

	log.Println("Database connection established successfully")
}

// createMedicinesTable creates the medicines table if it doesn't exist
func createMedicinesTable() error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS medicines (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			dosage VARCHAR(255) NOT NULL,
			frequency VARCHAR(255) NOT NULL,
			time_of_day VARCHAR(255) NOT NULL,
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := DB.Exec(createTableQuery)
	return err
}
