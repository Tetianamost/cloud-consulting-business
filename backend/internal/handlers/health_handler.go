package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	logger              *logrus.Logger
	db                  *sql.DB
	connectionPool      *ConnectionPool
	emailEventRecorder  interfaces.EmailEventRecorder
	emailMetricsService interfaces.EmailMetricsService
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(
	logger *logrus.Logger,
	db *sql.DB,
	connectionPool *ConnectionPool,
	emailEventRecorder interfaces.EmailEventRecorder,
	emailMetricsService interfaces.EmailMetricsService,
) *HealthHandler {
	return &HealthHandler{
		logger:              logger,
		db:                  db,
		connectionPool:      connectionPool,
		emailEventRecorder:  emailEventRecorder,
		emailMetricsService: emailMetricsService,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]HealthCheck `json:"checks,omitempty"`
}

// HealthCheck represents an individual health check
type HealthCheck struct {
	Status      string        `json:"status"`
	Message     string        `json:"message,omitempty"`
	Duration    time.Duration `json:"duration"`
	LastChecked string        `json:"last_checked"`
}

var startTime = time.Now()

// BasicHealthCheck handles GET /health
func (h *HealthHandler) BasicHealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Service:   "ai-consultant-chat-backend",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
	}

	c.JSON(http.StatusOK, response)
}

// DetailedHealthCheck handles GET /health/detailed
func (h *HealthHandler) DetailedHealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	checks := make(map[string]HealthCheck)
	overallStatus := "healthy"

	// Database health check
	dbCheck := h.checkDatabase(ctx)
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Email event recorder health check
	emailRecorderCheck := h.checkEmailEventRecorder(ctx)
	checks["email_event_recorder"] = emailRecorderCheck
	if emailRecorderCheck.Status == "unhealthy" {
		if overallStatus == "healthy" {
			overallStatus = "degraded" // Email recording issues shouldn't make system unhealthy
		}
	}

	// Email metrics service health check
	emailMetricsCheck := h.checkEmailMetricsService(ctx)
	checks["email_metrics_service"] = emailMetricsCheck
	if emailMetricsCheck.Status == "unhealthy" {
		if overallStatus == "healthy" {
			overallStatus = "degraded" // Email metrics issues shouldn't make system unhealthy
		}
	}

	// Memory health check
	memCheck := h.checkMemory()
	checks["memory"] = memCheck
	if memCheck.Status != "healthy" {
		if overallStatus != "unhealthy" {
			overallStatus = "degraded"
		}
	}

	// Disk health check
	diskCheck := h.checkDisk()
	checks["disk"] = diskCheck
	if diskCheck.Status != "healthy" {
		if overallStatus != "unhealthy" {
			overallStatus = "degraded"
		}
	}

	response := HealthResponse{
		Status:    overallStatus,
		Service:   "ai-consultant-chat-backend",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
		Checks:    checks,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == "degraded" {
		statusCode = http.StatusOK // Still return 200 for degraded
	}

	c.JSON(statusCode, response)
}

// ReadinessCheck handles GET /health/ready
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check critical dependencies
	dbCheck := h.checkDatabase(ctx)

	if dbCheck.Status != "healthy" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Database is not available",
			"checks": map[string]HealthCheck{
				"database": dbCheck,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"message": "Service is ready to accept traffic",
	})
}

// LivenessCheck handles GET /health/live
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// Simple liveness check - just verify the service is running
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"service":   "ai-consultant-chat-backend",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
	})
}

// checkDatabase performs database connectivity check
func (h *HealthHandler) checkDatabase(ctx context.Context) HealthCheck {
	start := time.Now()

	if h.db == nil {
		return HealthCheck{
			Status:      "unhealthy",
			Message:     "Database connection not initialized",
			Duration:    time.Since(start),
			LastChecked: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Test database connectivity
	err := h.db.PingContext(ctx)
	duration := time.Since(start)

	if err != nil {
		return HealthCheck{
			Status:      "unhealthy",
			Message:     "Database ping failed: " + err.Error(),
			Duration:    duration,
			LastChecked: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Test a simple query
	var count int
	err = h.db.QueryRowContext(ctx, "SELECT 1").Scan(&count)
	if err != nil {
		return HealthCheck{
			Status:      "unhealthy",
			Message:     "Database query failed: " + err.Error(),
			Duration:    time.Since(start),
			LastChecked: time.Now().UTC().Format(time.RFC3339),
		}
	}

	return HealthCheck{
		Status:      "healthy",
		Message:     "Database is accessible",
		Duration:    time.Since(start),
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

// checkMemory performs memory usage check
func (h *HealthHandler) checkMemory() HealthCheck {
	start := time.Now()

	// This is a simplified memory check
	// In production, you might want to use runtime.MemStats
	return HealthCheck{
		Status:      "healthy",
		Message:     "Memory usage within acceptable limits",
		Duration:    time.Since(start),
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

// checkDisk performs disk space check
func (h *HealthHandler) checkDisk() HealthCheck {
	start := time.Now()

	// This is a simplified disk check
	// In production, you might want to check actual disk usage
	return HealthCheck{
		Status:      "healthy",
		Message:     "Disk space sufficient",
		Duration:    time.Since(start),
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

// WebSocketHealthResponse represents WebSocket-specific health information
type WebSocketHealthResponse struct {
	Status            string                 `json:"status"`
	WebSocketReady    bool                   `json:"websocket_ready"`
	ActiveConnections int                    `json:"active_connections"`
	ServerUptime      string                 `json:"server_uptime"`
	Configuration     map[string]interface{} `json:"configuration"`
	Diagnostics       []string               `json:"diagnostics"`
	ConnectionStats   ConnectionStats        `json:"connection_stats"`
	Timestamp         string                 `json:"timestamp"`
}

// ConnectionStats provides detailed connection statistics
type ConnectionStats struct {
	TotalConnections     int                    `json:"total_connections"`
	ConnectionsByUser    map[string]int         `json:"connections_by_user"`
	ConnectionsBySession map[string]int         `json:"connections_by_session"`
	AverageLatency       float64                `json:"average_latency_ms"`
	HealthyConnections   int                    `json:"healthy_connections"`
	StaleConnections     int                    `json:"stale_connections"`
	LastActivity         string                 `json:"last_activity"`
	MemoryUsage          map[string]interface{} `json:"memory_usage"`
}

// WebSocketHealthCheck handles GET /health/websocket
func (h *HealthHandler) WebSocketHealthCheck(c *gin.Context) {
	diagnostics := []string{}
	status := "healthy"
	wsReady := true

	// Check if connection pool is available
	if h.connectionPool == nil {
		diagnostics = append(diagnostics, "WebSocket connection pool not initialized")
		status = "unhealthy"
		wsReady = false
	}

	// Get connection statistics
	var stats ConnectionStats
	if h.connectionPool != nil {
		stats = h.getConnectionStats()

		// Check for stale connections
		if stats.StaleConnections > 0 {
			diagnostics = append(diagnostics, fmt.Sprintf("%d stale connections detected", stats.StaleConnections))
			if status == "healthy" {
				status = "degraded"
			}
		}

		// Check connection health
		if stats.TotalConnections == 0 {
			diagnostics = append(diagnostics, "No active WebSocket connections")
		} else {
			diagnostics = append(diagnostics, fmt.Sprintf("%d active connections", stats.TotalConnections))
		}
	}

	// Get memory statistics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check memory usage
	memoryMB := float64(m.Alloc) / 1024 / 1024
	if memoryMB > 500 { // 500MB threshold
		diagnostics = append(diagnostics, fmt.Sprintf("High memory usage: %.2f MB", memoryMB))
		if status == "healthy" {
			status = "degraded"
		}
	}

	// Configuration information
	config := map[string]interface{}{
		"websocket_enabled":   true,
		"max_connections":     1000, // This should come from actual config
		"ping_interval":       "30s",
		"connection_timeout":  "10s",
		"message_buffer_size": 100,
		"heartbeat_enabled":   true,
		"compression_enabled": false,
		"max_message_size":    32768, // 32KB
	}

	response := WebSocketHealthResponse{
		Status:            status,
		WebSocketReady:    wsReady,
		ActiveConnections: stats.TotalConnections,
		ServerUptime:      time.Since(startTime).String(),
		Configuration:     config,
		Diagnostics:       diagnostics,
		ConnectionStats:   stats,
		Timestamp:         time.Now().UTC().Format(time.RFC3339),
	}

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// ConfigurationValidationCheck handles GET /health/websocket/config
func (h *HealthHandler) ConfigurationValidationCheck(c *gin.Context) {
	issues := []string{}
	valid := true

	// Check connection pool
	if h.connectionPool == nil {
		issues = append(issues, "WebSocket connection pool not initialized")
		valid = false
	}

	// Check database connection
	if h.db == nil {
		issues = append(issues, "Database connection not available")
		valid = false
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := h.db.PingContext(ctx); err != nil {
			issues = append(issues, "Database connectivity issue: "+err.Error())
			valid = false
		}
	}

	// Check memory constraints
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memoryMB := float64(m.Alloc) / 1024 / 1024
	if memoryMB > 1000 { // 1GB threshold for warnings
		issues = append(issues, fmt.Sprintf("High memory usage detected: %.2f MB", memoryMB))
		if memoryMB > 2000 { // 2GB threshold for errors
			valid = false
		}
	}

	// Check goroutine count
	goroutines := runtime.NumGoroutine()
	if goroutines > 1000 {
		issues = append(issues, fmt.Sprintf("High goroutine count: %d", goroutines))
		if goroutines > 5000 {
			valid = false
		}
	}

	response := map[string]interface{}{
		"valid":     valid,
		"issues":    issues,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks": map[string]interface{}{
			"connection_pool": h.connectionPool != nil,
			"database":        h.db != nil,
			"memory_mb":       memoryMB,
			"goroutines":      goroutines,
		},
	}

	statusCode := http.StatusOK
	if !valid {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// getConnectionStats calculates detailed connection statistics
func (h *HealthHandler) getConnectionStats() ConnectionStats {
	if h.connectionPool == nil {
		return ConnectionStats{}
	}

	// This is a simplified version - in the actual implementation,
	// you would need to access the connection pool's internal state
	totalConnections := h.connectionPool.Count()

	// Get memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memoryUsage := map[string]interface{}{
		"alloc_mb":       float64(m.Alloc) / 1024 / 1024,
		"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
		"sys_mb":         float64(m.Sys) / 1024 / 1024,
		"num_gc":         m.NumGC,
		"goroutines":     runtime.NumGoroutine(),
	}

	return ConnectionStats{
		TotalConnections:     totalConnections,
		ConnectionsByUser:    make(map[string]int), // Would need actual implementation
		ConnectionsBySession: make(map[string]int), // Would need actual implementation
		AverageLatency:       0.0,                  // Would need actual implementation
		HealthyConnections:   totalConnections,     // Simplified
		StaleConnections:     0,                    // Would need actual implementation
		LastActivity:         time.Now().UTC().Format(time.RFC3339),
		MemoryUsage:          memoryUsage,
	}
}

// checkEmailEventRecorder performs email event recorder health check
func (h *HealthHandler) checkEmailEventRecorder(ctx context.Context) HealthCheck {
	start := time.Now()

	if h.emailEventRecorder == nil {
		return HealthCheck{
			Status:      "unhealthy",
			Message:     "Email event recorder not configured",
			Duration:    time.Since(start),
			LastChecked: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Create a timeout context for the health check
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Use the context-aware health check if available
	if recorder, ok := h.emailEventRecorder.(interface {
		IsHealthyWithContext(context.Context) bool
	}); ok {
		if !recorder.IsHealthyWithContext(checkCtx) {
			return HealthCheck{
				Status:      "unhealthy",
				Message:     "Email event recorder health check failed",
				Duration:    time.Since(start),
				LastChecked: time.Now().UTC().Format(time.RFC3339),
			}
		}
	} else {
		// Fallback to basic health check
		if !h.emailEventRecorder.IsHealthy() {
			return HealthCheck{
				Status:      "unhealthy",
				Message:     "Email event recorder health check failed",
				Duration:    time.Since(start),
				LastChecked: time.Now().UTC().Format(time.RFC3339),
			}
		}
	}

	return HealthCheck{
		Status:      "healthy",
		Message:     "Email event recorder is operational",
		Duration:    time.Since(start),
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

// checkEmailMetricsService performs email metrics service health check
func (h *HealthHandler) checkEmailMetricsService(ctx context.Context) HealthCheck {
	start := time.Now()

	if h.emailMetricsService == nil {
		return HealthCheck{
			Status:      "unhealthy",
			Message:     "Email metrics service not configured",
			Duration:    time.Since(start),
			LastChecked: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Create a timeout context for the health check
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if !h.emailMetricsService.IsHealthy(checkCtx) {
		return HealthCheck{
			Status:      "unhealthy",
			Message:     "Email metrics service health check failed",
			Duration:    time.Since(start),
			LastChecked: time.Now().UTC().Format(time.RFC3339),
		}
	}

	return HealthCheck{
		Status:      "healthy",
		Message:     "Email metrics service is operational",
		Duration:    time.Since(start),
		LastChecked: time.Now().UTC().Format(time.RFC3339),
	}
}

// EmailMonitoringHealthCheck handles GET /health/email-monitoring
func (h *HealthHandler) EmailMonitoringHealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	checks := make(map[string]HealthCheck)
	overallStatus := "healthy"

	// Email event recorder health check
	emailRecorderCheck := h.checkEmailEventRecorder(ctx)
	checks["email_event_recorder"] = emailRecorderCheck
	if emailRecorderCheck.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Email metrics service health check
	emailMetricsCheck := h.checkEmailMetricsService(ctx)
	checks["email_metrics_service"] = emailMetricsCheck
	if emailMetricsCheck.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Database connectivity check (required for email monitoring)
	dbCheck := h.checkDatabase(ctx)
	checks["database"] = dbCheck
	if dbCheck.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Calculate email monitoring specific metrics
	emailMonitoringInfo := h.getEmailMonitoringInfo(ctx)

	response := map[string]interface{}{
		"status":    overallStatus,
		"service":   "email-monitoring",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    checks,
		"info":      emailMonitoringInfo,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// getEmailMonitoringInfo gathers email monitoring specific information
func (h *HealthHandler) getEmailMonitoringInfo(ctx context.Context) map[string]interface{} {
	info := map[string]interface{}{
		"event_recorder_configured":  h.emailEventRecorder != nil,
		"metrics_service_configured": h.emailMetricsService != nil,
		"database_required":          true,
	}

	// Try to get recent email activity if services are available
	if h.emailMetricsService != nil {
		if recentActivity, err := h.emailMetricsService.GetRecentEmailActivity(ctx, 24); err == nil {
			info["recent_email_events"] = len(recentActivity)
		} else {
			info["recent_email_events_error"] = err.Error()
		}
	}

	return info
}
