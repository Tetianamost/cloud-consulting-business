package interfaces

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware defines the interface for HTTP middleware
type Middleware interface {
	Handler() gin.HandlerFunc
	GetName() string
	GetOrder() int
}

// AuthMiddleware defines the interface for authentication middleware
type AuthMiddleware interface {
	Middleware
	ValidateToken(ctx context.Context, token string) (*AuthContext, error)
	ExtractToken(c *gin.Context) (string, error)
	RequireAuth() gin.HandlerFunc
	OptionalAuth() gin.HandlerFunc
}

// RateLimitMiddleware defines the interface for rate limiting middleware
type RateLimitMiddleware interface {
	Middleware
	CheckLimit(ctx context.Context, key string) (*RateLimitResult, error)
	GetLimitInfo(ctx context.Context, key string) (*RateLimitInfo, error)
	ResetLimit(ctx context.Context, key string) error
}

// CacheMiddleware defines the interface for caching middleware
type CacheMiddleware interface {
	Middleware
	CacheResponse(duration time.Duration) gin.HandlerFunc
	InvalidateCache(pattern string) error
	GetCacheStats() *CacheStats
}

// LoggingMiddleware defines the interface for logging middleware
type LoggingMiddleware interface {
	Middleware
	SetLogLevel(level string)
	SetLogFormat(format string)
	AddSensitiveField(field string)
	RemoveSensitiveField(field string)
}

// MetricsMiddleware defines the interface for metrics collection middleware
type MetricsMiddleware interface {
	Middleware
	RecordRequest(method, path string, statusCode int, duration time.Duration)
	GetMetrics() map[string]interface{}
	ResetMetrics()
}

// CORSMiddleware defines the interface for CORS middleware
type CORSMiddleware interface {
	Middleware
	SetAllowedOrigins(origins []string)
	SetAllowedMethods(methods []string)
	SetAllowedHeaders(headers []string)
	SetMaxAge(maxAge time.Duration)
}

// ValidationMiddleware defines the interface for request validation middleware
type ValidationMiddleware interface {
	Middleware
	ValidateJSON(schema interface{}) gin.HandlerFunc
	ValidateQuery(schema interface{}) gin.HandlerFunc
	ValidateParams(schema interface{}) gin.HandlerFunc
}

// Supporting types for middleware

// AuthContext represents authentication context
type AuthContext struct {
	UserID      string                 `json:"user_id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Roles       []string               `json:"roles"`
	Permissions []string               `json:"permissions"`
	TokenType   string                 `json:"token_type"`
	ExpiresAt   time.Time              `json:"expires_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// RateLimitResult represents the result of a rate limit check
type RateLimitResult struct {
	Allowed       bool          `json:"allowed"`
	Remaining     int           `json:"remaining"`
	ResetTime     time.Time     `json:"reset_time"`
	RetryAfter    time.Duration `json:"retry_after"`
	TotalRequests int           `json:"total_requests"`
}

// RateLimitInfo represents rate limit information
type RateLimitInfo struct {
	Key           string        `json:"key"`
	Limit         int           `json:"limit"`
	Remaining     int           `json:"remaining"`
	ResetTime     time.Time     `json:"reset_time"`
	WindowSize    time.Duration `json:"window_size"`
	TotalRequests int           `json:"total_requests"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	HitRate     float64 `json:"hit_rate"`
	Size        int64   `json:"size"`
	MaxSize     int64   `json:"max_size"`
	Evictions   int64   `json:"evictions"`
	LastCleared time.Time `json:"last_cleared"`
}

// RequestMetrics represents request metrics
type RequestMetrics struct {
	TotalRequests    int64                    `json:"total_requests"`
	RequestsByMethod map[string]int64         `json:"requests_by_method"`
	RequestsByPath   map[string]int64         `json:"requests_by_path"`
	RequestsByStatus map[int]int64            `json:"requests_by_status"`
	AverageLatency   time.Duration            `json:"average_latency"`
	P95Latency       time.Duration            `json:"p95_latency"`
	P99Latency       time.Duration            `json:"p99_latency"`
	ErrorRate        float64                  `json:"error_rate"`
	LastReset        time.Time                `json:"last_reset"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// MiddlewareConfig represents middleware configuration
type MiddlewareConfig struct {
	Name     string                 `json:"name"`
	Enabled  bool                   `json:"enabled"`
	Order    int                    `json:"order"`
	Settings map[string]interface{} `json:"settings"`
}

// MiddlewareChain represents a chain of middleware
type MiddlewareChain interface {
	Add(middleware Middleware) MiddlewareChain
	Remove(name string) MiddlewareChain
	GetMiddleware(name string) Middleware
	GetAll() []Middleware
	Apply(handler gin.HandlerFunc) gin.HandlerFunc
}