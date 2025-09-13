// Package shared provides common test utilities and helper functions
// that can be used across different test categories (integration, performance, email, etc.)
package shared

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// SimpleTestConfig holds basic configuration for test utilities (legacy)
type SimpleTestConfig struct {
	DatabaseURL string
	RedisURL    string
	LogLevel    string
	Timeout     time.Duration
}

// LoadSimpleTestConfig loads basic test configuration from environment variables
func LoadSimpleTestConfig() *SimpleTestConfig {
	return &SimpleTestConfig{
		DatabaseURL: getEnvOrDefault("TEST_DATABASE_URL", "postgres://test:test@localhost/test_db?sslmode=disable"),
		RedisURL:    getEnvOrDefault("TEST_REDIS_URL", "redis://localhost:6379/1"),
		LogLevel:    getEnvOrDefault("TEST_LOG_LEVEL", "error"),
		Timeout:     30 * time.Second,
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupSimpleTestLogger creates a logger configured for testing
func SetupSimpleTestLogger(level string) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}
	logger.SetLevel(logLevel)

	// Use JSON formatter for structured logging in tests
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	return logger
}

// CreateTestContext creates a context with timeout for tests
func CreateTestContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T, reason string) {
	if testing.Short() {
		t.Skipf("Skipping test in short mode: %s", reason)
	}
}

// RequireNoError is a helper that fails the test if error is not nil
func RequireNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
}

// AssertEventuallyTrue waits for a condition to become true within timeout
func AssertEventuallyTrue(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if condition() {
				return
			}
		case <-timeoutCh:
			t.Fatalf("Condition never became true within %v: %s", timeout, message)
		}
	}
}

// CleanupFunc represents a cleanup function that should be called after test
type CleanupFunc func()

// TestDatabase provides utilities for database testing
type TestDatabase struct {
	DB     *sql.DB
	config *SimpleTestConfig
	logger *logrus.Logger
}

// NewTestDatabase creates a new test database instance
func NewTestDatabase(config *SimpleTestConfig, logger *logrus.Logger) (*TestDatabase, error) {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	return &TestDatabase{
		DB:     db,
		config: config,
		logger: logger,
	}, nil
}

// Close closes the test database connection
func (td *TestDatabase) Close() error {
	if td.DB != nil {
		return td.DB.Close()
	}
	return nil
}

// TruncateAllTables truncates all tables in the test database
func (td *TestDatabase) TruncateAllTables(t *testing.T) {
	tables := []string{
		"chat_messages",
		"chat_sessions",
		"inquiries",
		"reports",
		"email_events",
	}

	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Log warning but don't fail test - table might not exist
			td.logger.WithError(err).Warnf("Failed to truncate table %s", table)
		}
	}
}

// SimpleMockHTTPServer provides utilities for creating mock HTTP servers in tests
type SimpleMockHTTPServer struct {
	// This will be expanded in future tasks when we move HTTP-related test files
}

// TestMetrics provides utilities for collecting test metrics
type TestMetrics struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// NewTestMetrics creates a new test metrics instance
func NewTestMetrics() *TestMetrics {
	return &TestMetrics{
		StartTime: time.Now(),
	}
}

// Finish marks the test as finished and calculates duration
func (tm *TestMetrics) Finish() {
	tm.EndTime = time.Now()
	tm.Duration = tm.EndTime.Sub(tm.StartTime)
}

// LogMetrics logs the test metrics
func (tm *TestMetrics) LogMetrics(logger *logrus.Logger, testName string) {
	logger.WithFields(logrus.Fields{
		"test_name":  testName,
		"start_time": tm.StartTime,
		"end_time":   tm.EndTime,
		"duration":   tm.Duration,
	}).Info("Test metrics")
}

// TestSuite provides a structured way to run test suites with setup/teardown
type TestSuite struct {
	Name     string
	Config   *SimpleTestConfig
	Logger   *logrus.Logger
	Database *TestDatabase
	Metrics  *TestMetrics
}

// NewTestSuite creates a new test suite with common setup
func NewTestSuite(name string) (*TestSuite, error) {
	config := LoadSimpleTestConfig()
	logger := SetupSimpleTestLogger(config.LogLevel)

	database, err := NewTestDatabase(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup test database: %w", err)
	}

	return &TestSuite{
		Name:     name,
		Config:   config,
		Logger:   logger,
		Database: database,
		Metrics:  NewTestMetrics(),
	}, nil
}

// Cleanup performs cleanup operations for the test suite
func (ts *TestSuite) Cleanup() {
	if ts.Metrics != nil {
		ts.Metrics.Finish()
		ts.Metrics.LogMetrics(ts.Logger, ts.Name)
	}

	if ts.Database != nil {
		ts.Database.Close()
	}
}

// SetupTest performs common test setup operations
func (ts *TestSuite) SetupTest(t *testing.T) {
	t.Helper()

	// Clean database tables
	if ts.Database != nil {
		ts.Database.TruncateAllTables(t)
	}

	ts.Logger.WithField("test_name", t.Name()).Info("Starting test")
}

// TeardownTest performs common test teardown operations
func (ts *TestSuite) TeardownTest(t *testing.T) {
	t.Helper()

	ts.Logger.WithField("test_name", t.Name()).Info("Finished test")
}
