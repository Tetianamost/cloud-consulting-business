package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
)

// DatabasePool manages database connections with connection pooling
type DatabasePool struct {
	db     *sql.DB
	config *DatabaseConfig
	logger *logrus.Logger
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewDatabasePool creates a new database pool with optimized connection settings
func NewDatabasePool(config *DatabaseConfig, logger *logrus.Logger) (*DatabasePool, error) {
	if config == nil {
		return nil, fmt.Errorf("database config cannot be nil")
	}

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	pool := &DatabasePool{
		db:     db,
		config: config,
		logger: logger,
	}

	pool.configureConnectionPool()

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"host":              config.Host,
		"port":              config.Port,
		"database":          config.Database,
		"max_open_conns":    config.MaxOpenConns,
		"max_idle_conns":    config.MaxIdleConns,
		"conn_max_lifetime": config.ConnMaxLifetime,
	}).Info("Database connection pool initialized successfully")

	return pool, nil
}

// configureConnectionPool sets up optimal connection pool settings for chat workloads
func (p *DatabasePool) configureConnectionPool() {
	// Set maximum number of open connections (increased for chat workloads)
	if p.config.MaxOpenConns > 0 {
		p.db.SetMaxOpenConns(p.config.MaxOpenConns)
	} else {
		p.db.SetMaxOpenConns(50) // Increased default: 50 connections for chat
	}

	// Set maximum number of idle connections (increased for quick response)
	if p.config.MaxIdleConns > 0 {
		p.db.SetMaxIdleConns(p.config.MaxIdleConns)
	} else {
		p.db.SetMaxIdleConns(10) // Increased default: 10 idle connections
	}

	// Set maximum lifetime of connections (increased for stability)
	if p.config.ConnMaxLifetime > 0 {
		p.db.SetConnMaxLifetime(p.config.ConnMaxLifetime)
	} else {
		p.db.SetConnMaxLifetime(10 * time.Minute) // Increased default: 10 minutes
	}

	// Set maximum idle time for connections (increased for chat sessions)
	if p.config.ConnMaxIdleTime > 0 {
		p.db.SetConnMaxIdleTime(p.config.ConnMaxIdleTime)
	} else {
		p.db.SetConnMaxIdleTime(2 * time.Minute) // Increased default: 2 minutes
	}
}

// GetDB returns the underlying sql.DB instance
func (p *DatabasePool) GetDB() *sql.DB {
	return p.db
}

// Close closes the database connection pool
func (p *DatabasePool) Close() error {
	if p.db != nil {
		p.logger.Info("Closing database connection pool")
		return p.db.Close()
	}
	return nil
}

// Ping tests the database connection
func (p *DatabasePool) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// GetStats returns database connection pool statistics
func (p *DatabasePool) GetStats() sql.DBStats {
	return p.db.Stats()
}

// IsHealthy checks if the database connection is healthy
func (p *DatabasePool) IsHealthy(ctx context.Context) bool {
	if err := p.Ping(ctx); err != nil {
		p.logger.WithError(err).Error("Database health check failed")
		return false
	}
	return true
}

// BeginTx starts a new database transaction
func (p *DatabasePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, opts)
}

// WithTx executes a function within a database transaction
func (p *DatabasePool) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := p.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			p.logger.WithError(rollbackErr).Error("Failed to rollback transaction")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// LogStats logs current database connection pool statistics
func (p *DatabasePool) LogStats() {
	stats := p.GetStats()
	p.logger.WithFields(logrus.Fields{
		"open_connections":     stats.OpenConnections,
		"in_use_connections":   stats.InUse,
		"idle_connections":     stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}).Debug("Database connection pool statistics")
}

// DefaultDatabaseConfig returns a default database configuration
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		Database:        "consulting",
		Username:        "postgres",
		Password:        "password",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

// LoadDatabaseConfigFromEnv loads database configuration from environment variables
func LoadDatabaseConfigFromEnv() *DatabaseConfig {
	config := DefaultDatabaseConfig()

	if host := getEnv("DB_HOST", ""); host != "" {
		config.Host = host
	}
	if port := getEnvAsInt("DB_PORT", 0); port > 0 {
		config.Port = port
	}
	if database := getEnv("DB_NAME", ""); database != "" {
		config.Database = database
	}
	if username := getEnv("DB_USER", ""); username != "" {
		config.Username = username
	}
	if password := getEnv("DB_PASSWORD", ""); password != "" {
		config.Password = password
	}
	if sslMode := getEnv("DB_SSL_MODE", ""); sslMode != "" {
		config.SSLMode = sslMode
	}
	if maxOpen := getEnvAsInt("DB_MAX_OPEN_CONNS", 0); maxOpen > 0 {
		config.MaxOpenConns = maxOpen
	}
	if maxIdle := getEnvAsInt("DB_MAX_IDLE_CONNS", 0); maxIdle > 0 {
		config.MaxIdleConns = maxIdle
	}
	if lifetime := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 0); lifetime > 0 {
		config.ConnMaxLifetime = lifetime
	}
	if idleTime := getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 0); idleTime > 0 {
		config.ConnMaxIdleTime = idleTime
	}

	return config
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
