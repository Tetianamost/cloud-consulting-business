package interfaces

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Database defines the interface for database operations
type Database interface {
	// Connection management
	Connect(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error
	IsHealthy(ctx context.Context) bool
	
	// Transaction management
	BeginTx(ctx context.Context) (Transaction, error)
	WithTx(ctx context.Context, fn func(tx Transaction) error) error
	
	// Migration management
	Migrate(ctx context.Context) error
	Rollback(ctx context.Context, version string) error
	GetMigrationVersion(ctx context.Context) (string, error)
	
	// Connection pool management
	GetStats() DatabaseStats
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
	
	// Raw access (use with caution)
	GetDB() *gorm.DB
	GetSQLDB() *sql.DB
}

// Transaction defines the interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
	GetDB() *gorm.DB
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Database        string        `json:"database"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	Timezone        string        `json:"timezone"`
}

// DatabaseStats represents database connection statistics
type DatabaseStats struct {
	OpenConnections     int           `json:"open_connections"`
	InUseConnections    int           `json:"in_use_connections"`
	IdleConnections     int           `json:"idle_connections"`
	WaitCount           int64         `json:"wait_count"`
	WaitDuration        time.Duration `json:"wait_duration"`
	MaxIdleClosed       int64         `json:"max_idle_closed"`
	MaxIdleTimeClosed   int64         `json:"max_idle_time_closed"`
	MaxLifetimeClosed   int64         `json:"max_lifetime_closed"`
}

// MigrationInfo represents information about a database migration
type MigrationInfo struct {
	Version     string    `json:"version"`
	Name        string    `json:"name"`
	AppliedAt   time.Time `json:"applied_at"`
	ExecutionTime int64   `json:"execution_time_ms"`
}

// QueryResult represents the result of a database query
type QueryResult struct {
	RowsAffected int64         `json:"rows_affected"`
	LastInsertID int64         `json:"last_insert_id,omitempty"`
	ExecutionTime time.Duration `json:"execution_time"`
	Error        error         `json:"error,omitempty"`
}

// DatabaseHealthStatus represents the health status of the database
type DatabaseHealthStatus struct {
	Status        HealthStatusType `json:"status"`
	ResponseTime  time.Duration    `json:"response_time"`
	ConnectionCount int            `json:"connection_count"`
	Error         string           `json:"error,omitempty"`
	LastChecked   time.Time        `json:"last_checked"`
}