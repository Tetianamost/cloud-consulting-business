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

