package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// PerformanceOptimizationHandler handles performance optimization endpoints
type PerformanceOptimizationHandler struct {
	logger               *logrus.Logger
	performanceOptimizer interfaces.PerformanceOptimizer
	intelligentCache     interfaces.IntelligentCache
	loadBalancer         interfaces.SessionLoadBalancer
	performanceMonitor   interfaces.PerformanceMonitor
}

// NewPerformanceOptimizationHandler creates a new performance optimization handler
func NewPerformanceOptimizationHandler(
	logger *logrus.Logger,
	optimizer interfaces.PerformanceOptimizer,
	cache interfaces.IntelligentCache,
	balancer interfaces.SessionLoadBalancer,
	monitor interfaces.PerformanceMonitor,
) *PerformanceOptimizationHandler {
	return &PerformanceOptimizationHandler{
		logger:               logger,
		performanceOptimizer: optimizer,
		intelligentCache:     cache,
		loadBalancer:         balancer,
		performanceMonitor:   monitor,
	}
}

// RegisterRoutes registers the performance optimization routes
func (h *PerformanceOptimizationHandler) RegisterRoutes(router *mux.Router) {
	// Performance optimization metrics
	router.HandleFunc("/api/v1/admin/performance/metrics", h.GetPerformanceMetrics).Methods("GET")
	router.HandleFunc("/api/v1/admin/performance/report", h.GetPerformanceReport).Methods("GET")

	// Cache management
	router.HandleFunc("/api/v1/admin/performance/cache/stats", h.GetCacheStatistics).Methods("GET")
	router.HandleFunc("/api/v1/admin/performance/cache/optimize", h.OptimizeCacheStrategy).Methods("POST")
	router.HandleFunc("/api/v1/admin/performance/cache/warm", h.WarmCache).Methods("POST")

	// Load balancing
	router.HandleFunc("/api/v1/admin/performance/load-balancer/metrics", h.GetLoadBalancingMetrics).Methods("GET")
	router.HandleFunc("/api/v1/admin/performance/load-balancer/optimize", h.OptimizeLoadBalancer).Methods("POST")

	// Alert management
	router.HandleFunc("/api/v1/admin/performance/alerts/thresholds", h.GetAlertThresholds).Methods("GET")
	router.HandleFunc("/api/v1/admin/performance/alerts/thresholds", h.SetAlertThresholds).Methods("PUT")

	// Health check
	router.HandleFunc("/api/v1/admin/performance/health", h.GetPerformanceHealth).Methods("GET")
}

// GetPerformanceMetrics returns current performance optimization metrics
func (h *PerformanceOptimizationHandler) GetPerformanceMetrics(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting performance optimization metrics")

	metrics := h.performanceOptimizer.GetPerformanceMetrics()

	response := map[string]interface{}{
		"success": true,
		"data":    metrics,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetPerformanceReport returns a comprehensive performance report
func (h *PerformanceOptimizationHandler) GetPerformanceReport(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting comprehensive performance report")

	report := h.performanceMonitor.GetPerformanceReport()

	response := map[string]interface{}{
		"success": true,
		"data":    report,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetCacheStatistics returns detailed cache statistics
func (h *PerformanceOptimizationHandler) GetCacheStatistics(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting cache statistics")

	stats := h.intelligentCache.GetCacheStatistics()

	response := map[string]interface{}{
		"success": true,
		"data":    stats,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// OptimizeCacheStrategy triggers cache strategy optimization
func (h *PerformanceOptimizationHandler) OptimizeCacheStrategy(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Optimizing cache strategy")

	h.intelligentCache.OptimizeCacheStrategy()

	response := map[string]interface{}{
		"success": true,
		"message": "Cache strategy optimization triggered",
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// WarmCache triggers cache warming
func (h *PerformanceOptimizationHandler) WarmCache(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Warming cache")

	h.intelligentCache.WarmCache()

	response := map[string]interface{}{
		"success": true,
		"message": "Cache warming triggered",
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetLoadBalancingMetrics returns load balancing metrics
func (h *PerformanceOptimizationHandler) GetLoadBalancingMetrics(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting load balancing metrics")

	metrics := h.loadBalancer.GetLoadBalancingMetrics()

	response := map[string]interface{}{
		"success": true,
		"data":    metrics,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// OptimizeLoadBalancer triggers load balancer optimization
func (h *PerformanceOptimizationHandler) OptimizeLoadBalancer(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Optimizing load balancer for response time")

	h.loadBalancer.OptimizeForResponseTime()

	response := map[string]interface{}{
		"success": true,
		"message": "Load balancer optimization triggered",
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetAlertThresholds returns current alert thresholds
func (h *PerformanceOptimizationHandler) GetAlertThresholds(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting alert thresholds")

	report := h.performanceMonitor.GetPerformanceReport()

	response := map[string]interface{}{
		"success": true,
		"data":    report.AlertThresholds,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// SetAlertThresholds updates alert thresholds
func (h *PerformanceOptimizationHandler) SetAlertThresholds(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Updating alert thresholds")

	var thresholds interfaces.AlertThresholds
	if err := json.NewDecoder(r.Body).Decode(&thresholds); err != nil {
		h.logger.WithError(err).Error("Failed to decode alert thresholds")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	h.performanceMonitor.SetAlertThresholds(&thresholds)

	response := map[string]interface{}{
		"success": true,
		"message": "Alert thresholds updated successfully",
		"data":    thresholds,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetPerformanceHealth returns overall performance health status
func (h *PerformanceOptimizationHandler) GetPerformanceHealth(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting performance health status")

	// Collect health data from all components
	optimizerMetrics := h.performanceOptimizer.GetPerformanceMetrics()
	cacheStats := h.intelligentCache.GetCacheStatistics()
	loadBalancerMetrics := h.loadBalancer.GetLoadBalancingMetrics()
	performanceReport := h.performanceMonitor.GetPerformanceReport()

	// Determine overall health status
	healthy := true
	issues := make([]string, 0)

	// Check cache hit rate
	if cacheStats.HitRate < 0.5 && cacheStats.TotalRequests > 100 {
		healthy = false
		issues = append(issues, "Low cache hit rate")
	}

	// Check response time
	if performanceReport.ResponseTimeMetrics.AverageResponseTime.Seconds() > 5.0 {
		healthy = false
		issues = append(issues, "High average response time")
	}

	// Check error rate
	errorRate := float64(performanceReport.RequestMetrics.FailedRequests) / float64(performanceReport.RequestMetrics.TotalRequests)
	if errorRate > 0.05 && performanceReport.RequestMetrics.TotalRequests > 100 {
		healthy = false
		issues = append(issues, "High error rate")
	}

	// Check load balancer
	if loadBalancerMetrics.RejectedSessions > 0 {
		healthy = false
		issues = append(issues, "Load balancer rejecting sessions")
	}

	healthStatus := map[string]interface{}{
		"healthy":               healthy,
		"issues":                issues,
		"optimizer_metrics":     optimizerMetrics,
		"cache_statistics":      cacheStats,
		"load_balancer_metrics": loadBalancerMetrics,
		"performance_report":    performanceReport,
		"timestamp":             performanceReport.GeneratedAt,
	}

	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	response := map[string]interface{}{
		"success": healthy,
		"data":    healthStatus,
	}

	utils.WriteJSONResponse(w, status, response)
}

// GetPerformanceMetricsPrometheus returns metrics in Prometheus format
func (h *PerformanceOptimizationHandler) GetPerformanceMetricsPrometheus(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Getting performance metrics in Prometheus format")

	optimizerMetrics := h.performanceOptimizer.GetPerformanceMetrics()
	cacheStats := h.intelligentCache.GetCacheStatistics()
	loadBalancerMetrics := h.loadBalancer.GetLoadBalancingMetrics()
	performanceReport := h.performanceMonitor.GetPerformanceReport()

	// Generate Prometheus format metrics
	prometheusMetrics := generatePrometheusMetrics(optimizerMetrics, cacheStats, loadBalancerMetrics, performanceReport)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(prometheusMetrics))
}

// Helper function to generate Prometheus format metrics
func generatePrometheusMetrics(
	optimizerMetrics *interfaces.PerformanceOptimizationMetrics,
	cacheStats *interfaces.CacheStatistics,
	loadBalancerMetrics *interfaces.LoadBalancingMetrics,
	performanceReport *interfaces.SystemPerformanceReport,
) string {
	metrics := ""

	// Optimizer metrics
	metrics += "# HELP enhanced_bedrock_total_requests Total number of requests processed\n"
	metrics += "# TYPE enhanced_bedrock_total_requests counter\n"
	metrics += "enhanced_bedrock_total_requests " + strconv.FormatInt(optimizerMetrics.TotalRequests, 10) + "\n"

	metrics += "# HELP enhanced_bedrock_cache_hit_rate Cache hit rate\n"
	metrics += "# TYPE enhanced_bedrock_cache_hit_rate gauge\n"
	metrics += "enhanced_bedrock_cache_hit_rate " + strconv.FormatFloat(optimizerMetrics.CacheHitRate, 'f', 2, 64) + "\n"

	metrics += "# HELP enhanced_bedrock_active_sessions Number of active sessions\n"
	metrics += "# TYPE enhanced_bedrock_active_sessions gauge\n"
	metrics += "enhanced_bedrock_active_sessions " + strconv.Itoa(optimizerMetrics.ActiveSessions) + "\n"

	metrics += "# HELP enhanced_bedrock_average_response_time Average response time in seconds\n"
	metrics += "# TYPE enhanced_bedrock_average_response_time gauge\n"
	metrics += "enhanced_bedrock_average_response_time " + strconv.FormatFloat(optimizerMetrics.AverageResponseTime.Seconds(), 'f', 3, 64) + "\n"

	// Cache metrics
	metrics += "# HELP enhanced_bedrock_cache_size Current cache size\n"
	metrics += "# TYPE enhanced_bedrock_cache_size gauge\n"
	metrics += "enhanced_bedrock_cache_size " + strconv.Itoa(cacheStats.CacheSize) + "\n"

	metrics += "# HELP enhanced_bedrock_cache_evictions Total cache evictions\n"
	metrics += "# TYPE enhanced_bedrock_cache_evictions counter\n"
	metrics += "enhanced_bedrock_cache_evictions " + strconv.FormatInt(cacheStats.Evictions, 10) + "\n"

	// Load balancer metrics
	metrics += "# HELP enhanced_bedrock_load_balancer_active_sessions Active load balanced sessions\n"
	metrics += "# TYPE enhanced_bedrock_load_balancer_active_sessions gauge\n"
	metrics += "enhanced_bedrock_load_balancer_active_sessions " + strconv.FormatInt(loadBalancerMetrics.ActiveSessions, 10) + "\n"

	metrics += "# HELP enhanced_bedrock_load_balancer_average_load Average consultant load\n"
	metrics += "# TYPE enhanced_bedrock_load_balancer_average_load gauge\n"
	metrics += "enhanced_bedrock_load_balancer_average_load " + strconv.FormatFloat(loadBalancerMetrics.AverageLoad, 'f', 2, 64) + "\n"

	// System metrics
	metrics += "# HELP enhanced_bedrock_cpu_usage CPU usage percentage\n"
	metrics += "# TYPE enhanced_bedrock_cpu_usage gauge\n"
	metrics += "enhanced_bedrock_cpu_usage " + strconv.FormatFloat(performanceReport.SystemMetrics.CPUUsage, 'f', 2, 64) + "\n"

	metrics += "# HELP enhanced_bedrock_memory_usage Memory usage percentage\n"
	metrics += "# TYPE enhanced_bedrock_memory_usage gauge\n"
	metrics += "enhanced_bedrock_memory_usage " + strconv.FormatFloat(performanceReport.SystemMetrics.MemoryUsage, 'f', 2, 64) + "\n"

	return metrics
}
