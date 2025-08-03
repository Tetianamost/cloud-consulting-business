package services

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// IntelligentAnalysisCache provides intelligent caching for frequently requested analysis types
type IntelligentAnalysisCache struct {
	logger    *logrus.Logger
	cache     map[string]*CachedAnalysis
	analytics map[string]*AnalysisTypeMetrics
	mu        sync.RWMutex

	// Cache configuration
	maxCacheSize       int
	defaultTTL         time.Duration
	intelligentTTL     bool
	compressionEnabled bool

	// Analytics for cache optimization
	totalRequests int64
	cacheHits     int64
	cacheMisses   int64
	evictions     int64
}

// CachedAnalysis represents a cached analysis result
type CachedAnalysis struct {
	Content      string                 `json:"content"`
	AnalysisType string                 `json:"analysis_type"`
	TokensUsed   int                    `json:"tokens_used"`
	Quality      float64                `json:"quality"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	LastAccessed time.Time              `json:"last_accessed"`
	AccessCount  int                    `json:"access_count"`
	TTL          time.Duration          `json:"ttl"`
	Compressed   bool                   `json:"compressed"`
}

// AnalysisTypeMetrics tracks metrics for different analysis types
type AnalysisTypeMetrics struct {
	AnalysisType   string        `json:"analysis_type"`
	RequestCount   int64         `json:"request_count"`
	CacheHitRate   float64       `json:"cache_hit_rate"`
	AverageTokens  float64       `json:"average_tokens"`
	AverageQuality float64       `json:"average_quality"`
	OptimalTTL     time.Duration `json:"optimal_ttl"`
	LastOptimized  time.Time     `json:"last_optimized"`
}

// NewIntelligentAnalysisCache creates a new intelligent analysis cache
func NewIntelligentAnalysisCache(logger *logrus.Logger) *IntelligentAnalysisCache {
	return &IntelligentAnalysisCache{
		logger:             logger,
		cache:              make(map[string]*CachedAnalysis),
		analytics:          make(map[string]*AnalysisTypeMetrics),
		maxCacheSize:       1000,
		defaultTTL:         30 * time.Minute,
		intelligentTTL:     true,
		compressionEnabled: true,
	}
}

// GetCachedAnalysis retrieves a cached analysis if available
func (c *IntelligentAnalysisCache) GetCachedAnalysis(analysisType, content string) *CachedAnalysis {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalRequests++

	// Generate cache key
	cacheKey := c.generateCacheKey(analysisType, content)

	// Check if cached analysis exists
	cached, exists := c.cache[cacheKey]
	if !exists {
		c.cacheMisses++
		c.updateAnalysisTypeMetrics(analysisType, false, 0, 0)
		return nil
	}

	// Check if cache entry has expired
	if c.isCacheExpired(cached) {
		delete(c.cache, cacheKey)
		c.cacheMisses++
		c.updateAnalysisTypeMetrics(analysisType, false, 0, 0)
		return nil
	}

	// Update access statistics
	cached.LastAccessed = time.Now()
	cached.AccessCount++

	c.cacheHits++
	c.updateAnalysisTypeMetrics(analysisType, true, cached.TokensUsed, cached.Quality)

	c.logger.WithFields(logrus.Fields{
		"analysis_type": analysisType,
		"cache_key":     cacheKey[:16], // Log first 16 chars for debugging
		"access_count":  cached.AccessCount,
		"age":           time.Since(cached.CreatedAt).String(),
	}).Debug("Cache hit for analysis")

	return cached
}

// CacheAnalysis stores an analysis result in the cache
func (c *IntelligentAnalysisCache) CacheAnalysis(analysisType, content, result string, tokensUsed int, quality float64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check cache size and evict if necessary
	if len(c.cache) >= c.maxCacheSize {
		c.evictLeastUsed()
	}

	// Generate cache key
	cacheKey := c.generateCacheKey(analysisType, content)

	// Determine optimal TTL for this analysis type
	ttl := c.calculateOptimalTTL(analysisType, quality)

	// Compress content if enabled and beneficial
	finalContent := result
	compressed := false
	if c.compressionEnabled && len(result) > 1000 {
		if compressedContent := c.compressContent(result); len(compressedContent) < len(result) {
			finalContent = compressedContent
			compressed = true
		}
	}

	// Create cached analysis
	cached := &CachedAnalysis{
		Content:      finalContent,
		AnalysisType: analysisType,
		TokensUsed:   tokensUsed,
		Quality:      quality,
		Metadata: map[string]interface{}{
			"original_length":   len(result),
			"compressed_length": len(finalContent),
			"compression_ratio": float64(len(finalContent)) / float64(len(result)),
		},
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		AccessCount:  0,
		TTL:          ttl,
		Compressed:   compressed,
	}

	// Store in cache
	c.cache[cacheKey] = cached

	c.logger.WithFields(logrus.Fields{
		"analysis_type": analysisType,
		"cache_key":     cacheKey[:16],
		"tokens_used":   tokensUsed,
		"quality":       quality,
		"ttl":           ttl.String(),
		"compressed":    compressed,
	}).Debug("Analysis cached")

	return nil
}

// OptimizeCacheStrategy optimizes caching strategy based on usage patterns
func (c *IntelligentAnalysisCache) OptimizeCacheStrategy() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info("Optimizing cache strategy based on usage patterns")

	// Analyze usage patterns for each analysis type
	for _, metrics := range c.analytics {
		// Calculate optimal TTL based on access patterns
		if metrics.RequestCount > 10 {
			if metrics.CacheHitRate > 0.8 {
				// High hit rate - increase TTL
				metrics.OptimalTTL = time.Duration(float64(metrics.OptimalTTL) * 1.2)
			} else if metrics.CacheHitRate < 0.3 {
				// Low hit rate - decrease TTL
				metrics.OptimalTTL = time.Duration(float64(metrics.OptimalTTL) * 0.8)
			}

			// Ensure TTL stays within reasonable bounds
			if metrics.OptimalTTL > 2*time.Hour {
				metrics.OptimalTTL = 2 * time.Hour
			} else if metrics.OptimalTTL < 5*time.Minute {
				metrics.OptimalTTL = 5 * time.Minute
			}

			metrics.LastOptimized = time.Now()
		}
	}

	// Adjust cache size based on hit rate
	overallHitRate := float64(c.cacheHits) / float64(c.totalRequests)
	if overallHitRate < 0.5 && c.maxCacheSize > 500 {
		c.maxCacheSize = int(float64(c.maxCacheSize) * 0.9)
		c.logger.WithField("new_max_size", c.maxCacheSize).Info("Reduced cache size due to low hit rate")
	} else if overallHitRate > 0.8 && c.maxCacheSize < 2000 {
		c.maxCacheSize = int(float64(c.maxCacheSize) * 1.1)
		c.logger.WithField("new_max_size", c.maxCacheSize).Info("Increased cache size due to high hit rate")
	}
}

// WarmCache pre-loads cache with common analysis types
func (c *IntelligentAnalysisCache) WarmCache() {
	c.logger.Info("Starting cache warming for common analysis types")

	// Common analysis patterns to warm up
	commonPatterns := []struct {
		analysisType string
		content      string
		mockResult   string
	}{
		{
			analysisType: "cost_analysis",
			content:      "AWS cost optimization",
			mockResult:   "Cost analysis template for AWS optimization scenarios",
		},
		{
			analysisType: "architecture_review",
			content:      "cloud architecture assessment",
			mockResult:   "Architecture review template for cloud assessments",
		},
		{
			analysisType: "security_assessment",
			content:      "security compliance review",
			mockResult:   "Security assessment template for compliance reviews",
		},
		{
			analysisType: "migration_planning",
			content:      "cloud migration strategy",
			mockResult:   "Migration planning template for cloud strategies",
		},
	}

	for _, pattern := range commonPatterns {
		// Only warm if not already cached
		if c.GetCachedAnalysis(pattern.analysisType, pattern.content) == nil {
			c.CacheAnalysis(pattern.analysisType, pattern.content, pattern.mockResult, 100, 0.8)
		}
	}

	c.logger.WithField("patterns_warmed", len(commonPatterns)).Info("Cache warming completed")
}

// GetCacheStatistics returns detailed cache statistics
func (c *IntelligentAnalysisCache) GetCacheStatistics() *CacheStatistics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hitRate := float64(0)
	if c.totalRequests > 0 {
		hitRate = float64(c.cacheHits) / float64(c.totalRequests)
	}

	// Calculate average cache age
	var totalAge time.Duration
	var validEntries int
	for _, cached := range c.cache {
		if !c.isCacheExpired(cached) {
			totalAge += time.Since(cached.CreatedAt)
			validEntries++
		}
	}

	avgAge := time.Duration(0)
	if validEntries > 0 {
		avgAge = totalAge / time.Duration(validEntries)
	}

	return &CacheStatistics{
		TotalRequests:      c.totalRequests,
		CacheHits:          c.cacheHits,
		CacheMisses:        c.cacheMisses,
		HitRate:            hitRate,
		CacheSize:          len(c.cache),
		MaxCacheSize:       c.maxCacheSize,
		ValidEntries:       validEntries,
		ExpiredEntries:     len(c.cache) - validEntries,
		AverageAge:         avgAge,
		Evictions:          c.evictions,
		AnalysisTypes:      len(c.analytics),
		CompressionEnabled: c.compressionEnabled,
		IntelligentTTL:     c.intelligentTTL,
		Timestamp:          time.Now(),
	}
}

// Private helper methods

func (c *IntelligentAnalysisCache) generateCacheKey(analysisType, content string) string {
	// Create a hash of analysis type and content for consistent caching
	hasher := sha256.New()
	hasher.Write([]byte(analysisType + ":" + strings.ToLower(strings.TrimSpace(content))))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *IntelligentAnalysisCache) isCacheExpired(cached *CachedAnalysis) bool {
	return time.Since(cached.CreatedAt) > cached.TTL
}

func (c *IntelligentAnalysisCache) calculateOptimalTTL(analysisType string, quality float64) time.Duration {
	if !c.intelligentTTL {
		return c.defaultTTL
	}

	// Get metrics for this analysis type
	metrics, exists := c.analytics[analysisType]
	if !exists {
		return c.defaultTTL
	}

	// Base TTL on quality and historical performance
	baseTTL := metrics.OptimalTTL
	if baseTTL == 0 {
		baseTTL = c.defaultTTL
	}

	// Adjust based on quality
	qualityMultiplier := 0.5 + (quality * 0.5) // Range: 0.5 to 1.0
	adjustedTTL := time.Duration(float64(baseTTL) * qualityMultiplier)

	// Ensure reasonable bounds
	if adjustedTTL < 5*time.Minute {
		adjustedTTL = 5 * time.Minute
	} else if adjustedTTL > 2*time.Hour {
		adjustedTTL = 2 * time.Hour
	}

	return adjustedTTL
}

func (c *IntelligentAnalysisCache) updateAnalysisTypeMetrics(analysisType string, hit bool, tokens int, quality float64) {
	metrics, exists := c.analytics[analysisType]
	if !exists {
		metrics = &AnalysisTypeMetrics{
			AnalysisType: analysisType,
			OptimalTTL:   c.defaultTTL,
		}
		c.analytics[analysisType] = metrics
	}

	metrics.RequestCount++

	// Update hit rate
	if hit {
		metrics.CacheHitRate = (metrics.CacheHitRate*float64(metrics.RequestCount-1) + 1.0) / float64(metrics.RequestCount)
	} else {
		metrics.CacheHitRate = (metrics.CacheHitRate * float64(metrics.RequestCount-1)) / float64(metrics.RequestCount)
	}

	// Update averages if we have data
	if tokens > 0 {
		metrics.AverageTokens = (metrics.AverageTokens*float64(metrics.RequestCount-1) + float64(tokens)) / float64(metrics.RequestCount)
	}

	if quality > 0 {
		metrics.AverageQuality = (metrics.AverageQuality*float64(metrics.RequestCount-1) + quality) / float64(metrics.RequestCount)
	}
}

func (c *IntelligentAnalysisCache) evictLeastUsed() {
	if len(c.cache) == 0 {
		return
	}

	// Find the least recently used entry
	var oldestKey string
	var oldestTime time.Time = time.Now()

	for key, cached := range c.cache {
		if cached.LastAccessed.Before(oldestTime) {
			oldestTime = cached.LastAccessed
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.cache, oldestKey)
		c.evictions++

		c.logger.WithFields(logrus.Fields{
			"evicted_key": oldestKey[:16],
			"age":         time.Since(oldestTime).String(),
		}).Debug("Evicted least used cache entry")
	}
}

func (c *IntelligentAnalysisCache) compressContent(content string) string {
	// Simple compression - in production, use proper compression library
	// This is a placeholder implementation
	return content
}

// CacheStatistics represents cache statistics
type CacheStatistics struct {
	TotalRequests      int64         `json:"total_requests"`
	CacheHits          int64         `json:"cache_hits"`
	CacheMisses        int64         `json:"cache_misses"`
	HitRate            float64       `json:"hit_rate"`
	CacheSize          int           `json:"cache_size"`
	MaxCacheSize       int           `json:"max_cache_size"`
	ValidEntries       int           `json:"valid_entries"`
	ExpiredEntries     int           `json:"expired_entries"`
	AverageAge         time.Duration `json:"average_age"`
	Evictions          int64         `json:"evictions"`
	AnalysisTypes      int           `json:"analysis_types"`
	CompressionEnabled bool          `json:"compression_enabled"`
	IntelligentTTL     bool          `json:"intelligent_ttl"`
	Timestamp          time.Time     `json:"timestamp"`
}
