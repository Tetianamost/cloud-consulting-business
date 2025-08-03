package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/storage"
	"github.com/sirupsen/logrus"
)

// SimpleChatResponseCache provides caching for AI responses using Redis
type SimpleChatResponseCache struct {
	redis  *storage.RedisCache
	logger *logrus.Logger

	// Cache configuration
	defaultTTL  time.Duration
	enableCache bool
}

// SimpleCachedResponse represents a cached AI response
type SimpleCachedResponse struct {
	Content    string                 `json:"content"`
	TokensUsed int                    `json:"tokens_used"`
	Metadata   map[string]interface{} `json:"metadata"`
	CachedAt   time.Time              `json:"cached_at"`
	HitCount   int                    `json:"hit_count"`
	LastHit    time.Time              `json:"last_hit"`
}

// NewSimpleChatResponseCache creates a new simple chat response cache
func NewSimpleChatResponseCache(redis *storage.RedisCache, logger *logrus.Logger) *SimpleChatResponseCache {
	return &SimpleChatResponseCache{
		redis:       redis,
		logger:      logger,
		defaultTTL:  24 * time.Hour, // Cache responses for 24 hours
		enableCache: true,
	}
}

// GetCachedResponse retrieves a cached response if available
func (c *SimpleChatResponseCache) GetCachedResponse(ctx context.Context, content string, contextStr string) (*SimpleCachedResponse, error) {
	if !c.enableCache || c.redis == nil {
		return nil, nil
	}

	cacheKey := c.generateCacheKey(content, contextStr)

	// Try to get from Redis using a generic key
	data, err := c.redis.GetGeneric(ctx, cacheKey)
	if err != nil {
		// Cache miss is not an error
		return nil, nil
	}

	// Deserialize cached response
	var cachedResponse SimpleCachedResponse
	if err := json.Unmarshal([]byte(data), &cachedResponse); err != nil {
		c.logger.WithError(err).WithField("cache_key", cacheKey).Warn("Failed to deserialize cached response")
		return nil, nil
	}

	// Update hit statistics
	cachedResponse.HitCount++
	cachedResponse.LastHit = time.Now()

	// Update cache with new statistics (async)
	go func() {
		if updatedData, err := json.Marshal(cachedResponse); err == nil {
			c.redis.SetGeneric(context.Background(), cacheKey, string(updatedData), c.defaultTTL)
		}
	}()

	c.logger.WithFields(logrus.Fields{
		"cache_key": cacheKey,
		"hit_count": cachedResponse.HitCount,
		"cached_at": cachedResponse.CachedAt,
	}).Debug("Cache hit for chat response")

	return &cachedResponse, nil
}

// CacheResponse stores a response in the cache
func (c *SimpleChatResponseCache) CacheResponse(ctx context.Context, content string, contextStr string, response string, tokensUsed int) error {
	if !c.enableCache || c.redis == nil {
		return nil
	}

	cacheKey := c.generateCacheKey(content, contextStr)

	cachedResponse := SimpleCachedResponse{
		Content:    response,
		TokensUsed: tokensUsed,
		Metadata:   map[string]interface{}{"cached": true},
		CachedAt:   time.Now(),
		HitCount:   0,
		LastHit:    time.Time{},
	}

	data, err := json.Marshal(cachedResponse)
	if err != nil {
		c.logger.WithError(err).WithField("cache_key", cacheKey).Error("Failed to serialize response for caching")
		return nil // Don't fail on cache errors
	}

	if err := c.redis.SetGeneric(ctx, cacheKey, string(data), c.defaultTTL); err != nil {
		c.logger.WithError(err).WithField("cache_key", cacheKey).Warn("Failed to cache response")
		return nil // Don't fail on cache errors
	}

	c.logger.WithFields(logrus.Fields{
		"cache_key":   cacheKey,
		"tokens_used": tokensUsed,
	}).Debug("Cached chat response")

	return nil
}

// generateCacheKey creates a cache key based on content and context
func (c *SimpleChatResponseCache) generateCacheKey(content string, contextStr string) string {
	// Hash the content for consistent key generation
	contentHash := c.hashString(strings.ToLower(strings.TrimSpace(content)))
	contextHash := c.hashString(contextStr)

	return fmt.Sprintf("chat_response:%s:%s", contentHash, contextHash)
}

// hashString creates a SHA256 hash of a string
func (c *SimpleChatResponseCache) hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])[:16] // Use first 16 characters for shorter keys
}

// SetCacheEnabled enables or disables caching
func (c *SimpleChatResponseCache) SetCacheEnabled(enabled bool) {
	c.enableCache = enabled
	c.logger.WithField("cache_enabled", enabled).Info("Cache enabled status changed")
}

// GetCacheStats returns basic cache statistics
func (c *SimpleChatResponseCache) GetCacheStats(ctx context.Context) map[string]interface{} {
	stats := map[string]interface{}{
		"cache_enabled": c.enableCache,
		"default_ttl":   c.defaultTTL.String(),
	}

	if c.redis != nil {
		if redisStats, err := c.redis.GetStats(ctx); err == nil {
			stats["redis_healthy"] = redisStats.IsHealthy
			stats["last_checked"] = redisStats.LastChecked
		}
	}

	return stats
}
