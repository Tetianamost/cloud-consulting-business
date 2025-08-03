package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
)

// RedisCache provides caching functionality for chat data
type RedisCache struct {
	client *redis.Client
	logger *logrus.Logger
	config *RedisCacheConfig
}

// RedisCacheConfig holds Redis cache configuration
type RedisCacheConfig struct {
	Host               string
	Port               int
	Password           string
	Database           int
	MaxRetries         int
	PoolSize           int
	MinIdleConns       int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	SessionTTL         time.Duration
	MessageTTL         time.Duration
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(config *RedisCacheConfig, logger *logrus.Logger) (*RedisCache, error) {
	if config == nil {
		return nil, fmt.Errorf("redis config cannot be nil")
	}

	// Create Redis client options
	opts := &redis.Options{
		Addr:               fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:           config.Password,
		DB:                 config.Database,
		MaxRetries:         config.MaxRetries,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	cache := &RedisCache{
		client: client,
		logger: logger,
		config: config,
	}

	logger.WithFields(logrus.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"database": config.Database,
	}).Info("Redis cache initialized successfully")

	return cache, nil
}

// Close closes the Redis connection
func (c *RedisCache) Close() error {
	if c.client != nil {
		c.logger.Info("Closing Redis connection")
		return c.client.Close()
	}
	return nil
}

// Ping tests the Redis connection
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// IsHealthy checks if Redis is healthy
func (c *RedisCache) IsHealthy(ctx context.Context) bool {
	if err := c.Ping(ctx); err != nil {
		c.logger.WithError(err).Error("Redis health check failed")
		return false
	}
	return true
}

// Session caching methods

// SetSession caches a chat session
func (c *RedisCache) SetSession(ctx context.Context, session *domain.ChatSession) error {
	key := c.sessionKey(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.config.SessionTTL).Err(); err != nil {
		c.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to cache session")
		return fmt.Errorf("failed to cache session: %w", err)
	}

	c.logger.WithField("session_id", session.ID).Debug("Session cached successfully")
	return nil
}

// GetSession retrieves a cached chat session
func (c *RedisCache) GetSession(ctx context.Context, sessionID string) (*domain.ChatSession, error) {
	key := c.sessionKey(sessionID)

	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to get cached session")
		return nil, fmt.Errorf("failed to get cached session: %w", err)
	}

	var session domain.ChatSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to unmarshal cached session")
		return nil, fmt.Errorf("failed to unmarshal cached session: %w", err)
	}

	c.logger.WithField("session_id", sessionID).Debug("Session retrieved from cache")
	return &session, nil
}

// DeleteSession removes a session from cache
func (c *RedisCache) DeleteSession(ctx context.Context, sessionID string) error {
	key := c.sessionKey(sessionID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to delete cached session")
		return fmt.Errorf("failed to delete cached session: %w", err)
	}

	c.logger.WithField("session_id", sessionID).Debug("Session deleted from cache")
	return nil
}

// SetUserSessions caches user's session list
func (c *RedisCache) SetUserSessions(ctx context.Context, userID string, sessions []*domain.ChatSession) error {
	key := c.userSessionsKey(userID)

	data, err := json.Marshal(sessions)
	if err != nil {
		return fmt.Errorf("failed to marshal user sessions: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.config.SessionTTL).Err(); err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("Failed to cache user sessions")
		return fmt.Errorf("failed to cache user sessions: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"user_id":       userID,
		"session_count": len(sessions),
	}).Debug("User sessions cached successfully")
	return nil
}

// GetUserSessions retrieves cached user sessions
func (c *RedisCache) GetUserSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	key := c.userSessionsKey(userID)

	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.WithError(err).WithField("user_id", userID).Error("Failed to get cached user sessions")
		return nil, fmt.Errorf("failed to get cached user sessions: %w", err)
	}

	var sessions []*domain.ChatSession
	if err := json.Unmarshal([]byte(data), &sessions); err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("Failed to unmarshal cached user sessions")
		return nil, fmt.Errorf("failed to unmarshal cached user sessions: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"user_id":       userID,
		"session_count": len(sessions),
	}).Debug("User sessions retrieved from cache")
	return sessions, nil
}

// InvalidateUserSessions removes user sessions from cache
func (c *RedisCache) InvalidateUserSessions(ctx context.Context, userID string) error {
	key := c.userSessionsKey(userID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("Failed to invalidate user sessions cache")
		return fmt.Errorf("failed to invalidate user sessions cache: %w", err)
	}

	c.logger.WithField("user_id", userID).Debug("User sessions cache invalidated")
	return nil
}

// Message caching methods

// SetSessionMessages caches recent messages for a session
func (c *RedisCache) SetSessionMessages(ctx context.Context, sessionID string, messages []*domain.ChatMessage) error {
	key := c.sessionMessagesKey(sessionID)

	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal session messages: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.config.MessageTTL).Err(); err != nil {
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to cache session messages")
		return fmt.Errorf("failed to cache session messages: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"message_count": len(messages),
	}).Debug("Session messages cached successfully")
	return nil
}

// GetSessionMessages retrieves cached messages for a session
func (c *RedisCache) GetSessionMessages(ctx context.Context, sessionID string) ([]*domain.ChatMessage, error) {
	key := c.sessionMessagesKey(sessionID)

	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to get cached session messages")
		return nil, fmt.Errorf("failed to get cached session messages: %w", err)
	}

	var messages []*domain.ChatMessage
	if err := json.Unmarshal([]byte(data), &messages); err != nil {
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to unmarshal cached session messages")
		return nil, fmt.Errorf("failed to unmarshal cached session messages: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"message_count": len(messages),
	}).Debug("Session messages retrieved from cache")
	return messages, nil
}

// InvalidateSessionMessages removes session messages from cache
func (c *RedisCache) InvalidateSessionMessages(ctx context.Context, sessionID string) error {
	key := c.sessionMessagesKey(sessionID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to invalidate session messages cache")
		return fmt.Errorf("failed to invalidate session messages cache: %w", err)
	}

	c.logger.WithField("session_id", sessionID).Debug("Session messages cache invalidated")
	return nil
}

// Cache key generation methods

func (c *RedisCache) sessionKey(sessionID string) string {
	return fmt.Sprintf("chat:session:%s", sessionID)
}

func (c *RedisCache) userSessionsKey(userID string) string {
	return fmt.Sprintf("chat:user_sessions:%s", userID)
}

func (c *RedisCache) sessionMessagesKey(sessionID string) string {
	return fmt.Sprintf("chat:session_messages:%s", sessionID)
}

// Cache statistics and monitoring

// GetStats returns Redis cache statistics
func (c *RedisCache) GetStats(ctx context.Context) (*CacheStats, error) {
	info, err := c.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis stats: %w", err)
	}

	poolStats := c.client.PoolStats()

	stats := &CacheStats{
		RedisInfo:   info,
		PoolStats:   poolStats,
		IsHealthy:   c.IsHealthy(ctx),
		LastChecked: time.Now(),
	}

	return stats, nil
}

// CacheStats represents cache statistics
type CacheStats struct {
	RedisInfo   string           `json:"redis_info"`
	PoolStats   *redis.PoolStats `json:"pool_stats"`
	IsHealthy   bool             `json:"is_healthy"`
	LastChecked time.Time        `json:"last_checked"`
}

// DefaultRedisCacheConfig returns a default Redis cache configuration
func DefaultRedisCacheConfig() *RedisCacheConfig {
	return &RedisCacheConfig{
		Host:               "localhost",
		Port:               6379,
		Password:           "",
		Database:           0,
		MaxRetries:         3,
		PoolSize:           10,
		MinIdleConns:       2,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolTimeout:        4 * time.Second,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
		SessionTTL:         30 * time.Minute,
		MessageTTL:         15 * time.Minute,
	}
}

// LoadRedisCacheConfigFromEnv loads Redis cache configuration from environment variables
func LoadRedisCacheConfigFromEnv() *RedisCacheConfig {
	config := DefaultRedisCacheConfig()

	if host := getEnv("REDIS_HOST", ""); host != "" {
		config.Host = host
	}
	if port := getEnvAsInt("REDIS_PORT", 0); port > 0 {
		config.Port = port
	}
	if password := getEnv("REDIS_PASSWORD", ""); password != "" {
		config.Password = password
	}
	if database := getEnvAsInt("REDIS_DATABASE", -1); database >= 0 {
		config.Database = database
	}
	if maxRetries := getEnvAsInt("REDIS_MAX_RETRIES", 0); maxRetries > 0 {
		config.MaxRetries = maxRetries
	}
	if poolSize := getEnvAsInt("REDIS_POOL_SIZE", 0); poolSize > 0 {
		config.PoolSize = poolSize
	}
	if minIdle := getEnvAsInt("REDIS_MIN_IDLE_CONNS", 0); minIdle > 0 {
		config.MinIdleConns = minIdle
	}
	if dialTimeout := getEnvAsDuration("REDIS_DIAL_TIMEOUT", 0); dialTimeout > 0 {
		config.DialTimeout = dialTimeout
	}
	if readTimeout := getEnvAsDuration("REDIS_READ_TIMEOUT", 0); readTimeout > 0 {
		config.ReadTimeout = readTimeout
	}
	if writeTimeout := getEnvAsDuration("REDIS_WRITE_TIMEOUT", 0); writeTimeout > 0 {
		config.WriteTimeout = writeTimeout
	}
	if poolTimeout := getEnvAsDuration("REDIS_POOL_TIMEOUT", 0); poolTimeout > 0 {
		config.PoolTimeout = poolTimeout
	}
	if idleTimeout := getEnvAsDuration("REDIS_IDLE_TIMEOUT", 0); idleTimeout > 0 {
		config.IdleTimeout = idleTimeout
	}
	if sessionTTL := getEnvAsDuration("REDIS_SESSION_TTL", 0); sessionTTL > 0 {
		config.SessionTTL = sessionTTL
	}
	if messageTTL := getEnvAsDuration("REDIS_MESSAGE_TTL", 0); messageTTL > 0 {
		config.MessageTTL = messageTTL
	}

	return config
}

// Generic cache methods for simple key-value operations

// SetGeneric sets a generic key-value pair with TTL
func (c *RedisCache) SetGeneric(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to set generic cache value")
		return fmt.Errorf("failed to set generic cache value: %w", err)
	}
	return nil
}

// GetGeneric gets a generic value by key
func (c *RedisCache) GetGeneric(ctx context.Context, key string) (string, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key not found")
		}
		c.logger.WithError(err).WithField("key", key).Error("Failed to get generic cache value")
		return "", fmt.Errorf("failed to get generic cache value: %w", err)
	}
	return data, nil
}

// DeleteGeneric deletes a generic key
func (c *RedisCache) DeleteGeneric(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to delete generic cache key")
		return fmt.Errorf("failed to delete generic cache key: %w", err)
	}
	return nil
}

// Set sets a key-value pair with TTL (for quality assurance service compatibility)
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttlSeconds int) error {
	ttl := time.Duration(ttlSeconds) * time.Second
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to set cache value")
		return fmt.Errorf("failed to set cache value: %w", err)
	}
	return nil
}

// Get gets a value by key (for quality assurance service compatibility)
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key not found")
		}
		c.logger.WithError(err).WithField("key", key).Error("Failed to get cache value")
		return nil, fmt.Errorf("failed to get cache value: %w", err)
	}
	return []byte(data), nil
}
