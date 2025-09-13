// Package shared provides test configuration management
package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// TestEnvironment represents different test environments
type TestEnvironment string

const (
	TestEnvUnit        TestEnvironment = "unit"
	TestEnvIntegration TestEnvironment = "integration"
	TestEnvPerformance TestEnvironment = "performance"
	TestEnvE2E         TestEnvironment = "e2e"
	TestEnvLocal       TestEnvironment = "local"
	TestEnvCI          TestEnvironment = "ci"
)

// DatabaseConfig holds database configuration for tests
type DatabaseConfig struct {
	Driver          string        `json:"driver"`
	URL             string        `json:"url"`
	MaxConnections  int           `json:"max_connections"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	MigrationsPath  string        `json:"migrations_path"`
	SeedDataPath    string        `json:"seed_data_path"`
}

// RedisConfig holds Redis configuration for tests
type RedisConfig struct {
	URL         string        `json:"url"`
	Password    string        `json:"password"`
	DB          int           `json:"db"`
	MaxRetries  int           `json:"max_retries"`
	DialTimeout time.Duration `json:"dial_timeout"`
	ReadTimeout time.Duration `json:"read_timeout"`
}

// AWSConfig holds AWS service configuration for tests
type AWSConfig struct {
	Region          string `json:"region"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token"`
	UseMocks        bool   `json:"use_mocks"`
}

// BedrockConfig holds Bedrock-specific test configuration
type BedrockConfig struct {
	BaseURL     string        `json:"base_url"`
	Model       string        `json:"model"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
	Timeout     time.Duration `json:"timeout"`
	UseMock     bool          `json:"use_mock"`
}

// SESConfig holds SES-specific test configuration
type SESConfig struct {
	SenderEmail    string        `json:"sender_email"`
	ReplyToEmail   string        `json:"reply_to_email"`
	Timeout        time.Duration `json:"timeout"`
	UseMock        bool          `json:"use_mock"`
	VerifiedEmails []string      `json:"verified_emails"`
}

// LoggingConfig holds logging configuration for tests
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	Structured bool   `json:"structured"`
}

// PerformanceConfig holds performance testing configuration
type PerformanceConfig struct {
	MaxResponseTime    time.Duration `json:"max_response_time"`
	MaxMemoryUsage     int64         `json:"max_memory_usage"`
	MaxCPUUsage        float64       `json:"max_cpu_usage"`
	ConcurrentRequests int           `json:"concurrent_requests"`
	TestDuration       time.Duration `json:"test_duration"`
	WarmupDuration     time.Duration `json:"warmup_duration"`
}

// QualityConfig holds quality assurance configuration
type QualityConfig struct {
	MinAccuracy       float64 `json:"min_accuracy"`
	MinCompleteness   float64 `json:"min_completeness"`
	MinRelevance      float64 `json:"min_relevance"`
	MinActionability  float64 `json:"min_actionability"`
	MinTechnicalDepth float64 `json:"min_technical_depth"`
	MinBusinessValue  float64 `json:"min_business_value"`
}

// TestConfig represents the complete test configuration
type TestConfig struct {
	Environment TestEnvironment   `json:"environment"`
	Database    DatabaseConfig    `json:"database"`
	Redis       RedisConfig       `json:"redis"`
	AWS         AWSConfig         `json:"aws"`
	Bedrock     BedrockConfig     `json:"bedrock"`
	SES         SESConfig         `json:"ses"`
	Logging     LoggingConfig     `json:"logging"`
	Performance PerformanceConfig `json:"performance"`
	Quality     QualityConfig     `json:"quality"`
	Timeout     time.Duration     `json:"timeout"`
	Parallel    bool              `json:"parallel"`
	Verbose     bool              `json:"verbose"`
	SkipSlow    bool              `json:"skip_slow"`
	CleanupDB   bool              `json:"cleanup_db"`
	SeedData    bool              `json:"seed_data"`
	CustomVars  map[string]string `json:"custom_vars"`
}

// ConfigManager manages test configurations for different environments
type ConfigManager struct {
	configs map[TestEnvironment]*TestConfig
	logger  *logrus.Logger
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(logger *logrus.Logger) *ConfigManager {
	return &ConfigManager{
		configs: make(map[TestEnvironment]*TestConfig),
		logger:  logger,
	}
}

// LoadConfig loads configuration for a specific environment
func (cm *ConfigManager) LoadConfig(env TestEnvironment) (*TestConfig, error) {
	// Check if already loaded
	if config, exists := cm.configs[env]; exists {
		return config, nil
	}

	// Load from file first
	config, err := cm.loadFromFile(env)
	if err != nil {
		// If file doesn't exist, create default config
		config = cm.createDefaultConfig(env)
	}

	// Override with environment variables
	cm.overrideWithEnvVars(config)

	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration for environment %s: %w", env, err)
	}

	// Cache the configuration
	cm.configs[env] = config

	cm.logger.WithFields(logrus.Fields{
		"environment":  env,
		"database_url": config.Database.URL,
		"redis_url":    config.Redis.URL,
		"use_mocks":    config.Bedrock.UseMock,
	}).Info("Test configuration loaded")

	return config, nil
}

// loadFromFile loads configuration from a JSON file
func (cm *ConfigManager) loadFromFile(env TestEnvironment) (*TestConfig, error) {
	configPath := cm.getConfigPath(env)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	config.Environment = env
	return &config, nil
}

// getConfigPath returns the configuration file path for an environment
func (cm *ConfigManager) getConfigPath(env TestEnvironment) string {
	configDir := os.Getenv("TEST_CONFIG_DIR")
	if configDir == "" {
		configDir = "backend/testing/config"
	}

	return filepath.Join(configDir, fmt.Sprintf("%s.json", env))
}

// createDefaultConfig creates a default configuration for an environment
func (cm *ConfigManager) createDefaultConfig(env TestEnvironment) *TestConfig {
	config := &TestConfig{
		Environment: env,
		Database: DatabaseConfig{
			Driver:          "postgres",
			URL:             "postgres://test:test@localhost/test_db?sslmode=disable",
			MaxConnections:  10,
			ConnMaxLifetime: 30 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
			MigrationsPath:  "backend/scripts",
			SeedDataPath:    "backend/testing/fixtures",
		},
		Redis: RedisConfig{
			URL:         "redis://localhost:6379/1",
			DB:          1,
			MaxRetries:  3,
			DialTimeout: 5 * time.Second,
			ReadTimeout: 3 * time.Second,
		},
		AWS: AWSConfig{
			Region:   "us-east-1",
			UseMocks: true,
		},
		Bedrock: BedrockConfig{
			BaseURL:     "http://localhost:8080",
			Model:       "anthropic.claude-3-sonnet-20240229-v1:0",
			MaxTokens:   1000,
			Temperature: 0.7,
			Timeout:     30 * time.Second,
			UseMock:     true,
		},
		SES: SESConfig{
			SenderEmail:    "test@example.com",
			ReplyToEmail:   "test@example.com",
			Timeout:        10 * time.Second,
			UseMock:        true,
			VerifiedEmails: []string{"test@example.com"},
		},
		Logging: LoggingConfig{
			Level:      "error",
			Format:     "json",
			Output:     "stdout",
			Structured: true,
		},
		Performance: PerformanceConfig{
			MaxResponseTime:    5 * time.Second,
			MaxMemoryUsage:     100 * 1024 * 1024, // 100MB
			MaxCPUUsage:        80.0,
			ConcurrentRequests: 10,
			TestDuration:       30 * time.Second,
			WarmupDuration:     5 * time.Second,
		},
		Quality: QualityConfig{
			MinAccuracy:       0.80,
			MinCompleteness:   0.75,
			MinRelevance:      0.80,
			MinActionability:  0.75,
			MinTechnicalDepth: 0.70,
			MinBusinessValue:  0.75,
		},
		Timeout:    30 * time.Second,
		Parallel:   false,
		Verbose:    false,
		SkipSlow:   false,
		CleanupDB:  true,
		SeedData:   false,
		CustomVars: make(map[string]string),
	}

	// Environment-specific overrides
	switch env {
	case TestEnvUnit:
		config.Database.URL = "sqlite://memory"
		config.Parallel = true
		config.SkipSlow = true
		config.CleanupDB = false

	case TestEnvIntegration:
		config.SeedData = true
		config.Bedrock.UseMock = false
		config.SES.UseMock = false
		config.AWS.UseMocks = false

	case TestEnvPerformance:
		config.Performance.ConcurrentRequests = 50
		config.Performance.TestDuration = 5 * time.Minute
		config.Logging.Level = "warn"

	case TestEnvE2E:
		config.SeedData = true
		config.Bedrock.UseMock = false
		config.SES.UseMock = false
		config.AWS.UseMocks = false
		config.Timeout = 2 * time.Minute

	case TestEnvCI:
		config.Parallel = true
		config.SkipSlow = true
		config.Logging.Level = "error"
		config.Performance.MaxResponseTime = 10 * time.Second
	}

	return config
}

// overrideWithEnvVars overrides configuration with environment variables
func (cm *ConfigManager) overrideWithEnvVars(config *TestConfig) {
	// Database overrides
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		config.Database.URL = url
	}
	if driver := os.Getenv("TEST_DATABASE_DRIVER"); driver != "" {
		config.Database.Driver = driver
	}

	// Redis overrides
	if url := os.Getenv("TEST_REDIS_URL"); url != "" {
		config.Redis.URL = url
	}

	// AWS overrides
	if region := os.Getenv("AWS_REGION"); region != "" {
		config.AWS.Region = region
	}
	if keyID := os.Getenv("AWS_ACCESS_KEY_ID"); keyID != "" {
		config.AWS.AccessKeyID = keyID
	}
	if secret := os.Getenv("AWS_SECRET_ACCESS_KEY"); secret != "" {
		config.AWS.SecretAccessKey = secret
	}
	if token := os.Getenv("AWS_SESSION_TOKEN"); token != "" {
		config.AWS.SessionToken = token
	}

	// Bedrock overrides
	if baseURL := os.Getenv("BEDROCK_BASE_URL"); baseURL != "" {
		config.Bedrock.BaseURL = baseURL
	}
	if model := os.Getenv("BEDROCK_MODEL"); model != "" {
		config.Bedrock.Model = model
	}
	if useMock := os.Getenv("BEDROCK_USE_MOCK"); useMock != "" {
		config.Bedrock.UseMock = strings.ToLower(useMock) == "true"
	}

	// SES overrides
	if sender := os.Getenv("SES_SENDER_EMAIL"); sender != "" {
		config.SES.SenderEmail = sender
	}
	if replyTo := os.Getenv("SES_REPLY_TO_EMAIL"); replyTo != "" {
		config.SES.ReplyToEmail = replyTo
	}
	if useMock := os.Getenv("SES_USE_MOCK"); useMock != "" {
		config.SES.UseMock = strings.ToLower(useMock) == "true"
	}

	// Logging overrides
	if level := os.Getenv("TEST_LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}

	// General test overrides
	if timeout := os.Getenv("TEST_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			config.Timeout = d
		}
	}
	if parallel := os.Getenv("TEST_PARALLEL"); parallel != "" {
		config.Parallel = strings.ToLower(parallel) == "true"
	}
	if verbose := os.Getenv("TEST_VERBOSE"); verbose != "" {
		config.Verbose = strings.ToLower(verbose) == "true"
	}
	if skipSlow := os.Getenv("TEST_SKIP_SLOW"); skipSlow != "" {
		config.SkipSlow = strings.ToLower(skipSlow) == "true"
	}
	if cleanupDB := os.Getenv("TEST_CLEANUP_DB"); cleanupDB != "" {
		config.CleanupDB = strings.ToLower(cleanupDB) == "true"
	}
	if seedData := os.Getenv("TEST_SEED_DATA"); seedData != "" {
		config.SeedData = strings.ToLower(seedData) == "true"
	}
}

// validateConfig validates the configuration
func (cm *ConfigManager) validateConfig(config *TestConfig) error {
	// Validate database configuration
	if config.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if config.Database.MaxConnections <= 0 {
		return fmt.Errorf("database max connections must be positive")
	}

	// Validate Redis configuration
	if config.Redis.URL == "" {
		return fmt.Errorf("Redis URL is required")
	}

	// Validate AWS configuration
	if config.AWS.Region == "" {
		return fmt.Errorf("AWS region is required")
	}

	// Validate Bedrock configuration
	if config.Bedrock.MaxTokens <= 0 {
		return fmt.Errorf("Bedrock max tokens must be positive")
	}
	if config.Bedrock.Temperature < 0 || config.Bedrock.Temperature > 1 {
		return fmt.Errorf("Bedrock temperature must be between 0 and 1")
	}

	// Validate SES configuration
	if config.SES.SenderEmail == "" {
		return fmt.Errorf("SES sender email is required")
	}

	// Validate performance configuration
	if config.Performance.MaxResponseTime <= 0 {
		return fmt.Errorf("max response time must be positive")
	}
	if config.Performance.ConcurrentRequests <= 0 {
		return fmt.Errorf("concurrent requests must be positive")
	}

	// Validate quality configuration
	if config.Quality.MinAccuracy < 0 || config.Quality.MinAccuracy > 1 {
		return fmt.Errorf("min accuracy must be between 0 and 1")
	}

	return nil
}

// SaveConfig saves configuration to a file
func (cm *ConfigManager) SaveConfig(env TestEnvironment, config *TestConfig) error {
	configPath := cm.getConfigPath(env)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	cm.logger.WithField("config_path", configPath).Info("Configuration saved")
	return nil
}

// GetConfig returns the cached configuration for an environment
func (cm *ConfigManager) GetConfig(env TestEnvironment) *TestConfig {
	return cm.configs[env]
}

// ListEnvironments returns all available test environments
func (cm *ConfigManager) ListEnvironments() []TestEnvironment {
	return []TestEnvironment{
		TestEnvUnit,
		TestEnvIntegration,
		TestEnvPerformance,
		TestEnvE2E,
		TestEnvLocal,
		TestEnvCI,
	}
}

// Global configuration manager instance
var globalConfigManager *ConfigManager

// InitGlobalConfigManager initializes the global configuration manager
func InitGlobalConfigManager(logger *logrus.Logger) {
	globalConfigManager = NewConfigManager(logger)
}

// GetGlobalConfig returns the global configuration for an environment
func GetGlobalConfig(env TestEnvironment) (*TestConfig, error) {
	if globalConfigManager == nil {
		return nil, fmt.Errorf("global config manager not initialized")
	}
	return globalConfigManager.LoadConfig(env)
}

// Helper functions for common configuration tasks

// GetDatabaseURL returns the database URL for an environment
func GetDatabaseURL(env TestEnvironment) (string, error) {
	config, err := GetGlobalConfig(env)
	if err != nil {
		return "", err
	}
	return config.Database.URL, nil
}

// GetRedisURL returns the Redis URL for an environment
func GetRedisURL(env TestEnvironment) (string, error) {
	config, err := GetGlobalConfig(env)
	if err != nil {
		return "", err
	}
	return config.Redis.URL, nil
}

// ShouldUseMocks returns whether to use mocks for an environment
func ShouldUseMocks(env TestEnvironment) (bool, error) {
	config, err := GetGlobalConfig(env)
	if err != nil {
		return true, err
	}
	return config.Bedrock.UseMock && config.SES.UseMock && config.AWS.UseMocks, nil
}

// GetTestTimeout returns the test timeout for an environment
func GetTestTimeout(env TestEnvironment) (time.Duration, error) {
	config, err := GetGlobalConfig(env)
	if err != nil {
		return 30 * time.Second, err
	}
	return config.Timeout, nil
}
