package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
)

// DatabaseConnection manages a database connection with configuration from config.DatabaseConfig
type DatabaseConnection struct {
	db     *sql.DB
	config *config.DatabaseConfig
	logger *logrus.Logger
}

// NewDatabaseConnection creates a new database connection using the provided configuration
func NewDatabaseConnection(cfg *config.DatabaseConfig, logger *logrus.Logger) (*DatabaseConnection, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config cannot be nil")
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Open database connection
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	if cfg.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	} else {
		db.SetMaxOpenConns(25) // Default
	}

	if cfg.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	} else {
		db.SetMaxIdleConns(5) // Default
	}

	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	} else {
		db.SetConnMaxLifetime(30 * time.Minute) // Default
	}

	conn := &DatabaseConnection{
		db:     db,
		config: cfg,
		logger: logger,
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"max_open_conns":    cfg.MaxOpenConnections,
		"max_idle_conns":    cfg.MaxIdleConnections,
		"conn_max_lifetime": cfg.ConnMaxLifetime,
		"email_events":      cfg.EnableEmailEvents,
	}).Info("Database connection initialized successfully")

	return conn, nil
}

// GetDB returns the underlying sql.DB instance
func (c *DatabaseConnection) GetDB() *sql.DB {
	return c.db
}

// Close closes the database connection
func (c *DatabaseConnection) Close() error {
	if c.db != nil {
		c.logger.Info("Closing database connection")
		return c.db.Close()
	}
	return nil
}

// Ping tests the database connection
func (c *DatabaseConnection) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// IsHealthy checks if the database connection is healthy
func (c *DatabaseConnection) IsHealthy(ctx context.Context) bool {
	if err := c.Ping(ctx); err != nil {
		c.logger.WithError(err).Error("Database health check failed")
		return false
	}
	return true
}

// GetStats returns database connection pool statistics
func (c *DatabaseConnection) GetStats() sql.DBStats {
	return c.db.Stats()
}

// RunMigration runs the email events migration if email events are enabled
func (c *DatabaseConnection) RunMigration(ctx context.Context, migrationSQL string) error {
	if !c.config.EnableEmailEvents {
		c.logger.Info("Email events disabled, skipping migration")
		return nil
	}

	c.logger.Info("Running email events migration")

	_, err := c.db.ExecContext(ctx, migrationSQL)
	if err != nil {
		c.logger.WithError(err).Error("Failed to run email events migration")
		return fmt.Errorf("failed to run migration: %w", err)
	}

	c.logger.Info("Email events migration completed successfully")
	return nil
}
