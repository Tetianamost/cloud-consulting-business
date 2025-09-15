package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds basic configuration for the application
type Config struct {
	Port               string
	LogLevel           int
	GinMode            string
	CORSAllowedOrigins []string
	JWTSecret          string
	Bedrock            BedrockConfig
	SES                SESConfig
	Chat               ChatConfig
	Database           DatabaseConfig
}

// ChatConfig holds chat system configuration including feature flags
type ChatConfig struct {
	Mode                    string // "websocket", "polling", or "auto"
	EnableWebSocketFallback bool   // Enable automatic fallback from WebSocket to polling
	WebSocketTimeout        int    // WebSocket connection timeout in seconds
	PollingInterval         int    // Default polling interval in milliseconds
	MaxReconnectAttempts    int    // Maximum reconnection attempts before fallback
	FallbackDelay           int    // Delay before fallback in milliseconds
}

// BedrockConfig holds Amazon Bedrock configuration
type BedrockConfig struct {
	APIKey  string
	Region  string
	ModelID string
	BaseURL string
	Timeout int
}

// SESConfig holds AWS SES configuration
type SESConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	SenderEmail     string
	ReplyToEmail    string
	Timeout         int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL                string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    int // in minutes
	EnableEmailEvents  bool
}

// Load loads basic configuration from environment variables without validation
func Load() (*Config, error) {
	// Load .env file if it exists (ignore errors for minimal setup)
	_ = godotenv.Load()

	cfg := &Config{
		// Prefer BACKEND_PORT to avoid conflicts with nginx listening on :80
		// Fallback to PORT if BACKEND_PORT is not set, and default to 8061
		Port:     getEnv("BACKEND_PORT", getEnv("PORT", "8061")),
		LogLevel: getEnvAsInt("LOG_LEVEL", 4), // Info level
		GinMode:  getEnv("GIN_MODE", "debug"),
		// Added http://localhost:3007 for development only
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{
			"http://localhost:3000",
			"http://localhost:3007", // for development only
		}),
		JWTSecret: getEnv("JWT_SECRET", "cloud-consulting-demo-secret"),
		Bedrock: BedrockConfig{
			APIKey:  getEnv("AWS_BEARER_TOKEN_BEDROCK", ""),
			Region:  getEnv("BEDROCK_REGION", "us-east-1"),
			ModelID: getEnv("BEDROCK_MODEL_ID", "amazon.nova-lite-v1:0"),
			BaseURL: getEnv("BEDROCK_BASE_URL", "https://bedrock-runtime.us-east-1.amazonaws.com"),
			Timeout: getEnvAsInt("BEDROCK_TIMEOUT_SECONDS", 30),
		},
		SES: SESConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_SES_REGION", "us-east-1"),
			SenderEmail:     getEnv("SES_SENDER_EMAIL", ""),
			ReplyToEmail:    getEnv("SES_REPLY_TO_EMAIL", ""),
			Timeout:         getEnvAsInt("SES_TIMEOUT_SECONDS", 30),
		},
		Chat: ChatConfig{
			Mode:                    getEnv("CHAT_MODE", "auto"), // "websocket", "polling", or "auto"
			EnableWebSocketFallback: getEnvAsBool("CHAT_ENABLE_WEBSOCKET_FALLBACK", true),
			WebSocketTimeout:        getEnvAsInt("CHAT_WEBSOCKET_TIMEOUT", 10),
			PollingInterval:         getEnvAsInt("CHAT_POLLING_INTERVAL", 3000),
			MaxReconnectAttempts:    getEnvAsInt("CHAT_MAX_RECONNECT_ATTEMPTS", 3),
			FallbackDelay:           getEnvAsInt("CHAT_FALLBACK_DELAY", 5000),
		},
		Database: DatabaseConfig{
			URL:                getEnv("DATABASE_URL", ""),
			MaxOpenConnections: getEnvAsInt("DB_MAX_OPEN_CONNECTIONS", 25),
			MaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
			ConnMaxLifetime:    getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30),
			EnableEmailEvents:  getEnvAsBool("ENABLE_EMAIL_EVENTS", true), // Default to true for production readiness
		},
	}

	log.Println("[CONFIG DEBUG] about to print CORSAllowedOrigins")
	fmt.Printf("[CONFIG DEBUG] CORSAllowedOrigins loaded: %v\n", cfg.CORSAllowedOrigins)
	log.Println("[CONFIG DEBUG] finished printing CORSAllowedOrigins")

	// Validate email event tracking configuration
	if err := cfg.ValidateEmailEventTracking(); err != nil {
		log.Printf("[CONFIG WARN] Email event tracking validation failed: %v", err)
	}

	return cfg, nil
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// ValidateEmailEventTracking validates the configuration for email event tracking
func (c *Config) ValidateEmailEventTracking() error {
	if !c.Database.EnableEmailEvents {
		return nil // Email events are disabled, no validation needed
	}

	// Check if database URL is provided when email events are enabled
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required when email event tracking is enabled (ENABLE_EMAIL_EVENTS=true)")
	}

	// Check if SES configuration is complete for email event tracking
	if c.SES.SenderEmail == "" {
		return fmt.Errorf("SES sender email is required for email event tracking")
	}

	if c.SES.AccessKeyID == "" {
		return fmt.Errorf("AWS access key ID is required for email event tracking")
	}

	if c.SES.SecretAccessKey == "" {
		return fmt.Errorf("AWS secret access key is required for email event tracking")
	}

	return nil
}

// IsEmailEventTrackingEnabled returns true if email event tracking is properly configured
func (c *Config) IsEmailEventTrackingEnabled() bool {
	return c.Database.EnableEmailEvents &&
		c.Database.URL != "" &&
		c.SES.SenderEmail != "" &&
		c.SES.AccessKeyID != "" &&
		c.SES.SecretAccessKey != ""
}
