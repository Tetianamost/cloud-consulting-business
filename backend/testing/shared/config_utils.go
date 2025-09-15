// Package shared provides configuration utilities for testing
package shared

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ConfigUtils provides utility functions for test configuration
type ConfigUtils struct {
	configManager *ConfigManager
	logger        *logrus.Logger
}

// NewConfigUtils creates a new configuration utilities instance
func NewConfigUtils(logger *logrus.Logger) *ConfigUtils {
	return &ConfigUtils{
		configManager: NewConfigManager(logger),
		logger:        logger,
	}
}

// DetectEnvironment detects the test environment from various sources
func (cu *ConfigUtils) DetectEnvironment() TestEnvironment {
	// Check explicit environment variable
	if env := os.Getenv("TEST_ENVIRONMENT"); env != "" {
		switch strings.ToLower(env) {
		case "unit":
			return TestEnvUnit
		case "integration":
			return TestEnvIntegration
		case "performance":
			return TestEnvPerformance
		case "e2e":
			return TestEnvE2E
		case "local":
			return TestEnvLocal
		case "ci":
			return TestEnvCI
		}
	}

	// Check CI environment variables
	if cu.isCI() {
		return TestEnvCI
	}

	// Check Go test flags
	if cu.isShortTest() {
		return TestEnvUnit
	}

	// Check if running in performance mode
	if cu.isPerformanceTest() {
		return TestEnvPerformance
	}

	// Default to local environment
	return TestEnvLocal
}

// isCI checks if running in a CI environment
func (cu *ConfigUtils) isCI() bool {
	ciVars := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"GITHUB_ACTIONS",
		"GITLAB_CI",
		"JENKINS_URL",
		"TRAVIS",
		"CIRCLECI",
	}

	for _, ciVar := range ciVars {
		if os.Getenv(ciVar) != "" {
			return true
		}
	}

	return false
}

// isShortTest checks if running with -short flag
func (cu *ConfigUtils) isShortTest() bool {
	// This would be set by the test runner
	return os.Getenv("TEST_SHORT") == "true"
}

// isPerformanceTest checks if running performance tests
func (cu *ConfigUtils) isPerformanceTest() bool {
	return os.Getenv("TEST_PERFORMANCE") == "true"
}

// SetupTestEnvironment sets up the test environment with appropriate configuration
func (cu *ConfigUtils) SetupTestEnvironment() (*TestConfig, error) {
	env := cu.detectEnvironmentFromContext()

	config, err := cu.configManager.LoadConfig(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for environment %s: %w", env, err)
	}

	// Initialize global config manager if not already done
	if globalConfigManager == nil {
		InitGlobalConfigManager(cu.logger)
	}

	cu.logger.WithFields(logrus.Fields{
		"environment": env,
		"parallel":    config.Parallel,
		"timeout":     config.Timeout,
		"use_mocks":   config.Bedrock.UseMock,
	}).Info("Test environment configured")

	return config, nil
}

// detectEnvironmentFromContext detects environment from current context
func (cu *ConfigUtils) detectEnvironmentFromContext() TestEnvironment {
	// Check command line arguments or environment
	env := cu.DetectEnvironment()

	// Override based on specific conditions
	if cu.shouldUseIntegrationEnv() {
		return TestEnvIntegration
	}

	if cu.shouldUseE2EEnv() {
		return TestEnvE2E
	}

	return env
}

// shouldUseIntegrationEnv checks if integration environment should be used
func (cu *ConfigUtils) shouldUseIntegrationEnv() bool {
	// Check if database is available
	if dbURL := os.Getenv("TEST_DATABASE_URL"); dbURL != "" && !strings.Contains(dbURL, "memory") {
		return true
	}

	// Check if external services are available
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		return true
	}

	return false
}

// shouldUseE2EEnv checks if E2E environment should be used
func (cu *ConfigUtils) shouldUseE2EEnv() bool {
	// Check for E2E specific environment variables
	return os.Getenv("TEST_E2E") == "true" || os.Getenv("E2E_TESTS") == "true"
}

// CreateTestConfigFromEnv creates a test configuration from environment variables
func (cu *ConfigUtils) CreateTestConfigFromEnv(baseEnv TestEnvironment) (*TestConfig, error) {
	// Load base configuration
	config, err := cu.configManager.LoadConfig(baseEnv)
	if err != nil {
		return nil, err
	}

	// Apply environment-specific overrides
	cu.applyEnvironmentOverrides(config)

	return config, nil
}

// applyEnvironmentOverrides applies environment-specific configuration overrides
func (cu *ConfigUtils) applyEnvironmentOverrides(config *TestConfig) {
	// Database overrides
	if url := os.Getenv("DATABASE_URL"); url != "" {
		config.Database.URL = url
	}

	// Redis overrides
	if url := os.Getenv("REDIS_URL"); url != "" {
		config.Redis.URL = url
	}

	// Timeout overrides
	if timeout := os.Getenv("TEST_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.Timeout = d
		}
	}

	// Parallel execution override
	if parallel := os.Getenv("TEST_PARALLEL"); parallel != "" {
		config.Parallel = strings.ToLower(parallel) == "true"
	}

	// Verbose output override
	if verbose := os.Getenv("TEST_VERBOSE"); verbose != "" {
		config.Verbose = strings.ToLower(verbose) == "true"
	}
}

// ValidateTestEnvironment validates that the test environment is properly configured
func (cu *ConfigUtils) ValidateTestEnvironment(config *TestConfig) error {
	// Validate database connectivity if not using mocks
	if !strings.Contains(config.Database.URL, "memory") {
		if err := cu.validateDatabaseConnection(config); err != nil {
			return fmt.Errorf("database validation failed: %w", err)
		}
	}

	// Validate Redis connectivity if configured
	if config.Redis.URL != "" {
		if err := cu.validateRedisConnection(config); err != nil {
			cu.logger.WithError(err).Warn("Redis validation failed, continuing with degraded functionality")
		}
	}

	// Validate AWS credentials if not using mocks
	if !config.AWS.UseMocks {
		if err := cu.validateAWSCredentials(config); err != nil {
			return fmt.Errorf("AWS validation failed: %w", err)
		}
	}

	return nil
}

// validateDatabaseConnection validates database connectivity
func (cu *ConfigUtils) validateDatabaseConnection(config *TestConfig) error {
	// This would implement actual database connection testing
	// For now, just check if URL is provided
	if config.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	cu.logger.WithField("database_url", config.Database.URL).Debug("Database connection validated")
	return nil
}

// validateRedisConnection validates Redis connectivity
func (cu *ConfigUtils) validateRedisConnection(config *TestConfig) error {
	// This would implement actual Redis connection testing
	// For now, just check if URL is provided
	if config.Redis.URL == "" {
		return fmt.Errorf("Redis URL is required")
	}

	cu.logger.WithField("redis_url", config.Redis.URL).Debug("Redis connection validated")
	return nil
}

// validateAWSCredentials validates AWS credentials
func (cu *ConfigUtils) validateAWSCredentials(config *TestConfig) error {
	if config.AWS.AccessKeyID == "" || config.AWS.SecretAccessKey == "" {
		return fmt.Errorf("AWS credentials are required when not using mocks")
	}

	cu.logger.WithField("aws_region", config.AWS.Region).Debug("AWS credentials validated")
	return nil
}

// GetConfigForTest returns the appropriate configuration for a specific test
func (cu *ConfigUtils) GetConfigForTest(testName string) (*TestConfig, error) {
	// Determine environment based on test name patterns
	env := cu.determineEnvironmentFromTestName(testName)

	return cu.configManager.LoadConfig(env)
}

// determineEnvironmentFromTestName determines the environment based on test name
func (cu *ConfigUtils) determineEnvironmentFromTestName(testName string) TestEnvironment {
	testName = strings.ToLower(testName)

	if strings.Contains(testName, "unit") || strings.Contains(testName, "mock") {
		return TestEnvUnit
	}

	if strings.Contains(testName, "integration") || strings.Contains(testName, "api") {
		return TestEnvIntegration
	}

	if strings.Contains(testName, "performance") || strings.Contains(testName, "load") || strings.Contains(testName, "benchmark") {
		return TestEnvPerformance
	}

	if strings.Contains(testName, "e2e") || strings.Contains(testName, "end") {
		return TestEnvE2E
	}

	// Default to unit tests
	return TestEnvUnit
}

// SetupTestLogger creates a logger configured for the test environment
func (cu *ConfigUtils) SetupTestLogger(config *TestConfig) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logger.SetLevel(level)

	// Set formatter
	if config.Logging.Structured {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	switch config.Logging.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		logger.SetOutput(os.Stdout)
	}

	return logger
}

// PrintConfigSummary prints a summary of the current configuration
func (cu *ConfigUtils) PrintConfigSummary(config *TestConfig) {
	cu.logger.WithFields(logrus.Fields{
		"environment":     config.Environment,
		"database_driver": config.Database.Driver,
		"use_mocks":       config.Bedrock.UseMock && config.SES.UseMock,
		"parallel":        config.Parallel,
		"timeout":         config.Timeout,
		"cleanup_db":      config.CleanupDB,
		"seed_data":       config.SeedData,
		"skip_slow":       config.SkipSlow,
		"log_level":       config.Logging.Level,
	}).Info("Test configuration summary")
}

// Global configuration utilities instance
var globalConfigUtils *ConfigUtils

// InitGlobalConfigUtils initializes the global configuration utilities
func InitGlobalConfigUtils(logger *logrus.Logger) {
	globalConfigUtils = NewConfigUtils(logger)
}

// GetGlobalConfigUtils returns the global configuration utilities instance
func GetGlobalConfigUtils() *ConfigUtils {
	return globalConfigUtils
}

// Helper functions for common configuration tasks

// MustGetConfig returns configuration or panics if not available
func MustGetConfig(env TestEnvironment) *TestConfig {
	config, err := GetGlobalConfig(env)
	if err != nil {
		panic(fmt.Sprintf("failed to get config for environment %s: %v", env, err))
	}
	return config
}

// IsTestEnvironment checks if running in a specific test environment
func IsTestEnvironment(env TestEnvironment) bool {
	current := os.Getenv("TEST_ENVIRONMENT")
	return strings.ToLower(current) == string(env)
}

// ShouldSkipSlow returns whether slow tests should be skipped
func ShouldSkipSlow() bool {
	if globalConfigUtils != nil {
		env := globalConfigUtils.DetectEnvironment()
		if config, err := GetGlobalConfig(env); err == nil {
			return config.SkipSlow
		}
	}

	// Fallback to environment variable
	return os.Getenv("TEST_SKIP_SLOW") == "true"
}

// ShouldRunInParallel returns whether tests should run in parallel
func ShouldRunInParallel() bool {
	if globalConfigUtils != nil {
		env := globalConfigUtils.DetectEnvironment()
		if config, err := GetGlobalConfig(env); err == nil {
			return config.Parallel
		}
	}

	// Fallback to environment variable
	return os.Getenv("TEST_PARALLEL") == "true"
}
