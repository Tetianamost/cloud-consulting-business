package config

import (
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
}

// BedrockConfig holds Amazon Bedrock configuration
type BedrockConfig struct {
	APIKey    string
	Region    string
	ModelID   string
	BaseURL   string
	Timeout   int
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

// Load loads basic configuration from environment variables without validation
func Load() (*Config, error) {
	// Load .env file if it exists (ignore errors for minimal setup)
	_ = godotenv.Load()
	
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		LogLevel:           getEnvAsInt("LOG_LEVEL", 4), // Info level
		GinMode:            getEnv("GIN_MODE", "debug"),
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		JWTSecret:          getEnv("JWT_SECRET", "cloud-consulting-demo-secret"),
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

