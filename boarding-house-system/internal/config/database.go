package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB               *sql.DB
	ErrDBUnavailable = errors.New("database unavailable")
)

// InitDB initializes the database connection with retry logic
func InitDB(cfg *Config) error {
	const maxRetries = 3
	const retryDelay = 5 * time.Second

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to open database: %v", attempt, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		// Configure connection pool
		DB.SetConnMaxLifetime(5 * time.Minute)
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(25)
		DB.SetConnMaxIdleTime(2 * time.Minute)

		// Verify connection
		if err = DB.Ping(); err != nil {
			log.Printf("Attempt %d/%d: Database ping failed: %v", attempt, maxRetries, err)
			_ = DB.Close()
			time.Sleep(retryDelay)
			continue
		}
		// Verify database selection
		var dbName string
		if err = DB.QueryRow("SELECT DATABASE()").Scan(&dbName); err != nil {
			return fmt.Errorf("failed to verify database: %w", err)
		}
		log.Printf("Using database: %s", dbName)
		log.Println("Successfully connected to database")
		return nil
	}

	return fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
}

// CloseDB safely closes the database connection
func CloseDB() {
	if DB == nil {
		return
	}

	// Attempt graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close all idle connections first
	DB.SetMaxIdleConns(0)
	DB.SetConnMaxLifetime(1 * time.Second)

	err := DB.Close()
	switch {
	case err == nil:
		log.Println("Database connection closed gracefully")
	case errors.Is(err, sql.ErrConnDone):
		log.Println("Database connection already closed")
	default:
		log.Printf("Error closing database: %v", err)
	}
}

// HealthCheck verifies database availability
func HealthCheck() error {
	if DB == nil {
		return ErrDBUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return DB.PingContext(ctx)
}

// WithTransaction executes a function within a database transaction
func WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
