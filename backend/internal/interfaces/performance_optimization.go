package interfaces

import (
	"context"
	"time"
)

// PerformanceOptimizer defines the interface for performance optimization
type PerformanceOptimizer interface {
	// OptimizeRequest optimizes a request for better performance
	OptimizeRequest(ctx context.Context, request *OptimizationRequest) (*OptimizationResult, error)

	// GetPerformanceMetrics returns current performance metrics
	GetPerformanceMetrics() *PerformanceOptimizationMetrics

	// MonitorPerformance starts performance monitoring
	MonitorPerformance(ctx context.Context)
}

// IntelligentCache defines the interface for intelligent caching
type IntelligentCache interface {
	// GetCachedAnalysis retrieves a cached analysis if available
	GetCachedAnalysis(analysisType, content string) *CachedAnalysis

	// CacheAnalysis stores an analysis result in the cache
	CacheAnalysis(analysisType, content, result string, tokensUsed int, quality float64) error

	// OptimizeCacheStrategy optimizes caching strategy based on usage patterns
	OptimizeCacheStrategy()

	// WarmCache pre-loads cache with common analysis types
	WarmCache()

	// GetCacheStatistics returns detailed cache statistics
	GetCacheStatistics() *CacheStatistics
}

// SessionLoadBalancer defines the interface for load balancing consultant sessions
type SessionLoadBalancer interface {
	// AssignSession assigns a session to the most appropriate consultant
	AssignSession(sessionID, preferredConsultantID string) string

	// ReleaseSession releases a session from a consultant
	ReleaseSession(sessionID string) error

	// GetActiveSessionCount returns the number of active sessions
	GetActiveSessionCount() int

	// GetLoadBalancingMetrics returns load balancing metrics
	GetLoadBalancingMetrics() *LoadBalancingMetrics

	// OptimizeForResponseTime optimizes load balancer settings for response time
	OptimizeForResponseTime()

	// CleanupExpiredSessions removes expired sessions
	CleanupExpiredSessions(ctx context.Context)
}

// PerformanceMonitor defines the interface for performance monitoring and alerting
type PerformanceMonitor interface {
	// RecordRequest records a request metric
	RecordRequest(success bool, responseTime time.Duration)

	// RecordCacheMetrics records cache performance metrics
	RecordCacheMetrics(hits, misses, size, evictions int64, avgAge time.Duration)

	// RecordSystemMetrics records system performance metrics
	RecordSystemMetrics(cpu, memory float64, goroutines int, heapSize int64, gcPause time.Duration)

	// GetPerformanceReport returns a comprehensive performance report
	GetPerformanceReport() *SystemPerformanceReport

	// StartMonitoring starts the performance monitoring loop
	StartMonitoring(ctx context.Context)

	// RegisterAlertHandler registers a custom alert handler
	RegisterAlertHandler(name string, handler AlertHandler)

	// SetAlertThresholds updates alert thresholds
	SetAlertThresholds(thresholds *AlertThresholds)
}

// OptimizationRequest represents a request for performance optimization
type OptimizationRequest struct {
	SessionID    string  `json:"session_id"`
	ConsultantID string  `json:"consultant_id"`
	AnalysisType string  `json:"analysis_type"`
	Content      string  `json:"content"`
	Prompt       string  `json:"prompt"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float64 `json:"temperature"`
	Priority     string  `json:"priority"`
}

// OptimizationResult represents the result of performance optimization
type OptimizationResult struct {
	Content          string               `json:"content"`
	TokensUsed       int                  `json:"tokens_used"`
	ResponseTime     time.Duration        `json:"response_time"`
	CacheHit         bool                 `json:"cache_hit"`
	Optimized        bool                 `json:"optimized"`
	SessionID        string               `json:"session_id"`
	OptimizedRequest *OptimizationRequest `json:"optimized_request,omitempty"`
}

// PerformanceOptimizationMetrics represents performance optimization metrics
type PerformanceOptimizationMetrics struct {
	TotalRequests        int64         `json:"total_requests"`
	OptimizedRequests    int64         `json:"optimized_requests"`
	CacheHits            int64         `json:"cache_hits"`
	LoadBalancedRequests int64         `json:"load_balanced_requests"`
	CacheHitRate         float64       `json:"cache_hit_rate"`
	OptimizationRate     float64       `json:"optimization_rate"`
	ActiveSessions       int           `json:"active_sessions"`
	AverageResponseTime  time.Duration `json:"average_response_time"`
	Timestamp            time.Time     `json:"timestamp"`
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

// LoadBalancingMetrics represents load balancing metrics
type LoadBalancingMetrics struct {
	TotalSessions        int64     `json:"total_sessions"`
	ActiveSessions       int64     `json:"active_sessions"`
	BalancedSessions     int64     `json:"balanced_sessions"`
	RejectedSessions     int64     `json:"rejected_sessions"`
	TotalConsultants     int64     `json:"total_consultants"`
	AvailableConsultants int64     `json:"available_consultants"`
	BusyConsultants      int64     `json:"busy_consultants"`
	AverageLoad          float64   `json:"average_load"`
	Strategy             string    `json:"strategy"`
	Timestamp            time.Time `json:"timestamp"`
}

// SystemPerformanceReport represents a comprehensive system performance report
type SystemPerformanceReport struct {
	RequestMetrics      SystemRequestMetrics      `json:"request_metrics"`
	ResponseTimeMetrics SystemResponseTimeMetrics `json:"response_time_metrics"`
	CacheMetrics        SystemCacheMetrics        `json:"cache_metrics"`
	SystemMetrics       SystemMetrics             `json:"system_metrics"`
	AlertThresholds     AlertThresholds           `json:"alert_thresholds"`
	GeneratedAt         time.Time                 `json:"generated_at"`
}

// SystemRequestMetrics tracks request-related metrics
type SystemRequestMetrics struct {
	TotalRequests         int64     `json:"total_requests"`
	SuccessfulRequests    int64     `json:"successful_requests"`
	FailedRequests        int64     `json:"failed_requests"`
	TimeoutRequests       int64     `json:"timeout_requests"`
	RequestsPerSecond     float64   `json:"requests_per_second"`
	ConcurrentRequests    int64     `json:"concurrent_requests"`
	MaxConcurrentRequests int64     `json:"max_concurrent_requests"`
	LastUpdated           time.Time `json:"last_updated"`
}

// SystemResponseTimeMetrics tracks response time statistics
type SystemResponseTimeMetrics struct {
	AverageResponseTime time.Duration   `json:"average_response_time"`
	MedianResponseTime  time.Duration   `json:"median_response_time"`
	P95ResponseTime     time.Duration   `json:"p95_response_time"`
	P99ResponseTime     time.Duration   `json:"p99_response_time"`
	MinResponseTime     time.Duration   `json:"min_response_time"`
	MaxResponseTime     time.Duration   `json:"max_response_time"`
	ResponseTimeHistory []time.Duration `json:"response_time_history"`
	LastUpdated         time.Time       `json:"last_updated"`
}

// SystemCacheMetrics tracks cache performance
type SystemCacheMetrics struct {
	CacheHits       int64         `json:"cache_hits"`
	CacheMisses     int64         `json:"cache_misses"`
	CacheHitRate    float64       `json:"cache_hit_rate"`
	CacheSize       int64         `json:"cache_size"`
	CacheEvictions  int64         `json:"cache_evictions"`
	AverageCacheAge time.Duration `json:"average_cache_age"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// SystemMetrics tracks system-level performance
type SystemMetrics struct {
	CPUUsage       float64       `json:"cpu_usage"`
	MemoryUsage    float64       `json:"memory_usage"`
	GoroutineCount int           `json:"goroutine_count"`
	HeapSize       int64         `json:"heap_size"`
	GCPauseTime    time.Duration `json:"gc_pause_time"`
	LastUpdated    time.Time     `json:"last_updated"`
}

// AlertThresholds defines thresholds for performance alerts
type AlertThresholds struct {
	MaxResponseTime       time.Duration `json:"max_response_time"`
	MinCacheHitRate       float64       `json:"min_cache_hit_rate"`
	MaxErrorRate          float64       `json:"max_error_rate"`
	MaxConcurrentRequests int64         `json:"max_concurrent_requests"`
	MaxCPUUsage           float64       `json:"max_cpu_usage"`
	MaxMemoryUsage        float64       `json:"max_memory_usage"`
}

// AlertHandler defines a function to handle performance alerts
type AlertHandler func(ctx context.Context, alert *PerformanceAlert) error

// PerformanceAlert represents a performance alert
type PerformanceAlert struct {
	ID        string                   `json:"id"`
	Type      PerformanceAlertType     `json:"type"`
	Severity  PerformanceAlertSeverity `json:"severity"`
	Message   string                   `json:"message"`
	Metric    string                   `json:"metric"`
	Value     interface{}              `json:"value"`
	Threshold interface{}              `json:"threshold"`
	Timestamp time.Time                `json:"timestamp"`
	Metadata  map[string]interface{}   `json:"metadata"`
}

// PerformanceAlertType defines types of performance alerts
type PerformanceAlertType string

const (
	PerformanceAlertTypeResponseTime   PerformanceAlertType = "response_time"
	PerformanceAlertTypeCacheHitRate   PerformanceAlertType = "cache_hit_rate"
	PerformanceAlertTypeErrorRate      PerformanceAlertType = "error_rate"
	PerformanceAlertTypeConcurrency    PerformanceAlertType = "concurrency"
	PerformanceAlertTypeSystemResource PerformanceAlertType = "system_resource"
)

// PerformanceAlertSeverity defines alert severity levels
type PerformanceAlertSeverity string

const (
	PerformanceAlertSeverityInfo     PerformanceAlertSeverity = "info"
	PerformanceAlertSeverityWarning  PerformanceAlertSeverity = "warning"
	PerformanceAlertSeverityCritical PerformanceAlertSeverity = "critical"
)
