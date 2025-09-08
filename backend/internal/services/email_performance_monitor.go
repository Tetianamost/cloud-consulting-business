package services

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// PerformanceMetric represents a performance measurement
type PerformanceMetric struct {
	Name      string                 `json:"name"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit"`
	Timestamp time.Time              `json:"timestamp"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// QueryPerformanceStats represents database query performance statistics
type QueryPerformanceStats struct {
	QueryType       string        `json:"query_type"`
	AverageDuration time.Duration `json:"average_duration"`
	TotalCalls      int64         `json:"total_calls"`
	SlowQueries     int64         `json:"slow_queries"`
	ErrorRate       float64       `json:"error_rate"`
	LastExecuted    time.Time     `json:"last_executed"`
	IndexUsageRatio float64       `json:"index_usage_ratio"`
}

// TablePerformanceStats represents table-level performance statistics
type TablePerformanceStats struct {
	TableName       string     `json:"table_name"`
	TotalSize       int64      `json:"total_size_bytes"`
	IndexSize       int64      `json:"index_size_bytes"`
	RowCount        int64      `json:"row_count"`
	SequentialScans int64      `json:"sequential_scans"`
	IndexScans      int64      `json:"index_scans"`
	IndexUsageRatio float64    `json:"index_usage_ratio"`
	LastVacuum      *time.Time `json:"last_vacuum,omitempty"`
	LastAnalyze     *time.Time `json:"last_analyze,omitempty"`
}

// EmailPerformanceMonitor monitors email system performance
type EmailPerformanceMonitor struct {
	db          *sql.DB
	logger      *logrus.Logger
	metrics     []PerformanceMetric
	metricsLock sync.RWMutex

	// Configuration
	maxMetricsHistory  int
	slowQueryThreshold time.Duration

	// Monitoring state
	isMonitoring   bool
	monitoringStop chan struct{}
	lastCollection time.Time
}

// NewEmailPerformanceMonitor creates a new performance monitor
func NewEmailPerformanceMonitor(db *sql.DB, logger *logrus.Logger) *EmailPerformanceMonitor {
	return &EmailPerformanceMonitor{
		db:                 db,
		logger:             logger,
		metrics:            make([]PerformanceMetric, 0),
		maxMetricsHistory:  1000,
		slowQueryThreshold: 1 * time.Second,
		monitoringStop:     make(chan struct{}),
	}
}

// StartMonitoring begins continuous performance monitoring
func (m *EmailPerformanceMonitor) StartMonitoring(ctx context.Context, interval time.Duration) {
	if m.isMonitoring {
		m.logger.Warn("Performance monitoring is already running")
		return
	}

	m.isMonitoring = true
	m.logger.WithField("interval", interval).Info("Starting email performance monitoring")

	go m.monitoringLoop(ctx, interval)
}

// StopMonitoring stops continuous performance monitoring
func (m *EmailPerformanceMonitor) StopMonitoring() {
	if !m.isMonitoring {
		return
	}

	m.logger.Info("Stopping email performance monitoring")
	close(m.monitoringStop)
	m.isMonitoring = false
}

// monitoringLoop runs the continuous monitoring process
func (m *EmailPerformanceMonitor) monitoringLoop(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Performance monitoring stopped due to context cancellation")
			return
		case <-m.monitoringStop:
			m.logger.Info("Performance monitoring stopped")
			return
		case <-ticker.C:
			m.collectMetrics(ctx)
		}
	}
}

// collectMetrics collects all performance metrics
func (m *EmailPerformanceMonitor) collectMetrics(ctx context.Context) {
	m.logger.Debug("Collecting performance metrics")

	// Collect database performance metrics
	if err := m.collectDatabaseMetrics(ctx); err != nil {
		m.logger.WithError(err).Error("Failed to collect database metrics")
	}

	// Collect query performance metrics
	if err := m.collectQueryMetrics(ctx); err != nil {
		m.logger.WithError(err).Error("Failed to collect query metrics")
	}

	// Collect table performance metrics
	if err := m.collectTableMetrics(ctx); err != nil {
		m.logger.WithError(err).Error("Failed to collect table metrics")
	}

	// Collect cache performance metrics
	if err := m.collectCacheMetrics(ctx); err != nil {
		m.logger.WithError(err).Error("Failed to collect cache metrics")
	}

	m.lastCollection = time.Now()
}

// collectDatabaseMetrics collects general database performance metrics
func (m *EmailPerformanceMonitor) collectDatabaseMetrics(ctx context.Context) error {
	// Database connection count
	var activeConnections int
	err := m.db.QueryRowContext(ctx,
		"SELECT count(*) FROM pg_stat_activity WHERE state = 'active'").Scan(&activeConnections)
	if err != nil {
		return fmt.Errorf("failed to get active connections: %w", err)
	}

	m.addMetric(PerformanceMetric{
		Name:      "database_active_connections",
		Value:     float64(activeConnections),
		Unit:      "count",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "database"},
	})

	// Database size
	var dbSize int64
	err = m.db.QueryRowContext(ctx,
		"SELECT pg_database_size(current_database())").Scan(&dbSize)
	if err != nil {
		return fmt.Errorf("failed to get database size: %w", err)
	}

	m.addMetric(PerformanceMetric{
		Name:      "database_size",
		Value:     float64(dbSize),
		Unit:      "bytes",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "database"},
	})

	// Transaction rate
	var xactCommit, xactRollback int64
	err = m.db.QueryRowContext(ctx,
		"SELECT xact_commit, xact_rollback FROM pg_stat_database WHERE datname = current_database()").
		Scan(&xactCommit, &xactRollback)
	if err != nil {
		return fmt.Errorf("failed to get transaction stats: %w", err)
	}

	m.addMetric(PerformanceMetric{
		Name:      "database_transactions_committed",
		Value:     float64(xactCommit),
		Unit:      "count",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "database"},
	})

	m.addMetric(PerformanceMetric{
		Name:      "database_transactions_rolled_back",
		Value:     float64(xactRollback),
		Unit:      "count",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "database"},
	})

	return nil
}

// collectQueryMetrics collects query performance metrics
func (m *EmailPerformanceMonitor) collectQueryMetrics(ctx context.Context) error {
	// Check if pg_stat_statements is available
	var extensionExists bool
	err := m.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'pg_stat_statements')").
		Scan(&extensionExists)
	if err != nil || !extensionExists {
		m.logger.Debug("pg_stat_statements extension not available, skipping query metrics")
		return nil
	}

	// Get email-related query statistics
	query := `
		SELECT 
			'email_events_queries' as query_type,
			avg(mean_exec_time) as avg_duration,
			sum(calls) as total_calls,
			sum(CASE WHEN mean_exec_time > $1 THEN calls ELSE 0 END) as slow_queries
		FROM pg_stat_statements 
		WHERE query LIKE '%email_events%'
		AND query NOT LIKE '%pg_stat_statements%'`

	var queryType string
	var avgDuration float64
	var totalCalls, slowQueries int64

	err = m.db.QueryRowContext(ctx, query, m.slowQueryThreshold.Milliseconds()).
		Scan(&queryType, &avgDuration, &totalCalls, &slowQueries)
	if err != nil {
		if err == sql.ErrNoRows {
			// No email queries found, which is fine
			return nil
		}
		return fmt.Errorf("failed to get query metrics: %w", err)
	}

	m.addMetric(PerformanceMetric{
		Name:      "query_average_duration",
		Value:     avgDuration,
		Unit:      "milliseconds",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "query", "query_type": queryType},
	})

	m.addMetric(PerformanceMetric{
		Name:      "query_total_calls",
		Value:     float64(totalCalls),
		Unit:      "count",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "query", "query_type": queryType},
	})

	m.addMetric(PerformanceMetric{
		Name:      "query_slow_queries",
		Value:     float64(slowQueries),
		Unit:      "count",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "query", "query_type": queryType},
	})

	return nil
}

// collectTableMetrics collects table-level performance metrics
func (m *EmailPerformanceMonitor) collectTableMetrics(ctx context.Context) error {
	// Get email_events table statistics
	query := `
		SELECT 
			schemaname,
			tablename,
			pg_total_relation_size(schemaname||'.'||tablename) as total_size,
			pg_indexes_size(schemaname||'.'||tablename) as index_size,
			n_tup_ins + n_tup_upd + n_tup_del as total_operations,
			seq_scan,
			idx_scan,
			CASE 
				WHEN (seq_scan + idx_scan) > 0 THEN 
					idx_scan::float / (seq_scan + idx_scan)::float * 100
				ELSE 0 
			END as index_usage_ratio,
			last_vacuum,
			last_analyze
		FROM pg_stat_user_tables 
		WHERE tablename IN ('email_events', 'email_events_archive')`

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get table metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schemaName, tableName string
		var totalSize, indexSize, totalOps, seqScan, idxScan int64
		var indexUsageRatio float64
		var lastVacuum, lastAnalyze sql.NullTime

		err := rows.Scan(&schemaName, &tableName, &totalSize, &indexSize,
			&totalOps, &seqScan, &idxScan, &indexUsageRatio, &lastVacuum, &lastAnalyze)
		if err != nil {
			m.logger.WithError(err).Error("Failed to scan table metrics row")
			continue
		}

		tags := map[string]string{
			"type":   "table",
			"table":  tableName,
			"schema": schemaName,
		}

		m.addMetric(PerformanceMetric{
			Name:      "table_total_size",
			Value:     float64(totalSize),
			Unit:      "bytes",
			Timestamp: time.Now(),
			Tags:      tags,
		})

		m.addMetric(PerformanceMetric{
			Name:      "table_index_size",
			Value:     float64(indexSize),
			Unit:      "bytes",
			Timestamp: time.Now(),
			Tags:      tags,
		})

		m.addMetric(PerformanceMetric{
			Name:      "table_sequential_scans",
			Value:     float64(seqScan),
			Unit:      "count",
			Timestamp: time.Now(),
			Tags:      tags,
		})

		m.addMetric(PerformanceMetric{
			Name:      "table_index_scans",
			Value:     float64(idxScan),
			Unit:      "count",
			Timestamp: time.Now(),
			Tags:      tags,
		})

		m.addMetric(PerformanceMetric{
			Name:      "table_index_usage_ratio",
			Value:     indexUsageRatio,
			Unit:      "percent",
			Timestamp: time.Now(),
			Tags:      tags,
		})
	}

	return nil
}

// collectCacheMetrics collects cache-related performance metrics
func (m *EmailPerformanceMonitor) collectCacheMetrics(ctx context.Context) error {
	// Get buffer cache hit ratio
	var bufferHitRatio float64
	err := m.db.QueryRowContext(ctx, `
		SELECT 
			CASE 
				WHEN (blks_hit + blks_read) > 0 THEN 
					blks_hit::float / (blks_hit + blks_read)::float * 100
				ELSE 0 
			END as buffer_hit_ratio
		FROM pg_stat_database 
		WHERE datname = current_database()`).Scan(&bufferHitRatio)
	if err != nil {
		return fmt.Errorf("failed to get buffer hit ratio: %w", err)
	}

	m.addMetric(PerformanceMetric{
		Name:      "database_buffer_hit_ratio",
		Value:     bufferHitRatio,
		Unit:      "percent",
		Timestamp: time.Now(),
		Tags:      map[string]string{"type": "cache"},
	})

	return nil
}

// addMetric adds a metric to the collection
func (m *EmailPerformanceMonitor) addMetric(metric PerformanceMetric) {
	m.metricsLock.Lock()
	defer m.metricsLock.Unlock()

	m.metrics = append(m.metrics, metric)

	// Keep only the most recent metrics
	if len(m.metrics) > m.maxMetricsHistory {
		m.metrics = m.metrics[len(m.metrics)-m.maxMetricsHistory:]
	}
}

// GetMetrics returns collected performance metrics
func (m *EmailPerformanceMonitor) GetMetrics(since *time.Time) []PerformanceMetric {
	m.metricsLock.RLock()
	defer m.metricsLock.RUnlock()

	if since == nil {
		// Return all metrics
		metrics := make([]PerformanceMetric, len(m.metrics))
		copy(metrics, m.metrics)
		return metrics
	}

	// Filter metrics by timestamp
	var filteredMetrics []PerformanceMetric
	for _, metric := range m.metrics {
		if metric.Timestamp.After(*since) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}

	return filteredMetrics
}

// GetQueryPerformanceStats returns detailed query performance statistics
func (m *EmailPerformanceMonitor) GetQueryPerformanceStats(ctx context.Context) (*QueryPerformanceStats, error) {
	// Check if pg_stat_statements is available
	var extensionExists bool
	err := m.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'pg_stat_statements')").
		Scan(&extensionExists)
	if err != nil || !extensionExists {
		return nil, fmt.Errorf("pg_stat_statements extension not available")
	}

	query := `
		SELECT 
			'email_events_queries' as query_type,
			avg(mean_exec_time) as avg_duration,
			sum(calls) as total_calls,
			sum(CASE WHEN mean_exec_time > $1 THEN calls ELSE 0 END) as slow_queries,
			max(last_exec) as last_executed
		FROM pg_stat_statements 
		WHERE query LIKE '%email_events%'
		AND query NOT LIKE '%pg_stat_statements%'`

	var stats QueryPerformanceStats
	var avgDurationMs float64
	var lastExecuted sql.NullTime

	err = m.db.QueryRowContext(ctx, query, m.slowQueryThreshold.Milliseconds()).
		Scan(&stats.QueryType, &avgDurationMs, &stats.TotalCalls,
			&stats.SlowQueries, &lastExecuted)
	if err != nil {
		if err == sql.ErrNoRows {
			return &QueryPerformanceStats{QueryType: "email_events_queries"}, nil
		}
		return nil, fmt.Errorf("failed to get query performance stats: %w", err)
	}

	stats.AverageDuration = time.Duration(avgDurationMs * float64(time.Millisecond))
	if lastExecuted.Valid {
		stats.LastExecuted = lastExecuted.Time
	}

	// Calculate error rate (simplified - would need more complex tracking in production)
	if stats.TotalCalls > 0 {
		stats.ErrorRate = float64(stats.SlowQueries) / float64(stats.TotalCalls) * 100
	}

	// Get index usage ratio for email_events table
	err = m.db.QueryRowContext(ctx, `
		SELECT 
			CASE 
				WHEN (seq_scan + idx_scan) > 0 THEN 
					idx_scan::float / (seq_scan + idx_scan)::float * 100
				ELSE 0 
			END as index_usage_ratio
		FROM pg_stat_user_tables 
		WHERE tablename = 'email_events'`).Scan(&stats.IndexUsageRatio)
	if err != nil {
		m.logger.WithError(err).Warn("Failed to get index usage ratio")
	}

	return &stats, nil
}

// GetTablePerformanceStats returns detailed table performance statistics
func (m *EmailPerformanceMonitor) GetTablePerformanceStats(ctx context.Context, tableName string) (*TablePerformanceStats, error) {
	query := `
		SELECT 
			schemaname||'.'||tablename as full_table_name,
			pg_total_relation_size(schemaname||'.'||tablename) as total_size,
			pg_indexes_size(schemaname||'.'||tablename) as index_size,
			n_tup_ins + n_tup_upd + n_tup_del as row_count_estimate,
			seq_scan,
			idx_scan,
			CASE 
				WHEN (seq_scan + idx_scan) > 0 THEN 
					idx_scan::float / (seq_scan + idx_scan)::float * 100
				ELSE 0 
			END as index_usage_ratio,
			last_vacuum,
			last_analyze
		FROM pg_stat_user_tables 
		WHERE tablename = $1`

	var stats TablePerformanceStats
	var fullTableName string
	var lastVacuum, lastAnalyze sql.NullTime

	err := m.db.QueryRowContext(ctx, query, tableName).
		Scan(&fullTableName, &stats.TotalSize, &stats.IndexSize, &stats.RowCount,
			&stats.SequentialScans, &stats.IndexScans, &stats.IndexUsageRatio,
			&lastVacuum, &lastAnalyze)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("table %s not found", tableName)
		}
		return nil, fmt.Errorf("failed to get table performance stats: %w", err)
	}

	stats.TableName = tableName
	if lastVacuum.Valid {
		stats.LastVacuum = &lastVacuum.Time
	}
	if lastAnalyze.Valid {
		stats.LastAnalyze = &lastAnalyze.Time
	}

	return &stats, nil
}

// GetPerformanceSummary returns a summary of current performance status
func (m *EmailPerformanceMonitor) GetPerformanceSummary(ctx context.Context) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Get query performance
	queryStats, err := m.GetQueryPerformanceStats(ctx)
	if err == nil {
		summary["query_performance"] = queryStats
	}

	// Get table performance
	tableStats, err := m.GetTablePerformanceStats(ctx, "email_events")
	if err == nil {
		summary["table_performance"] = tableStats
	}

	// Get recent metrics
	since := time.Now().Add(-1 * time.Hour)
	recentMetrics := m.GetMetrics(&since)
	summary["recent_metrics_count"] = len(recentMetrics)

	// Calculate average response times from recent metrics
	var totalDuration float64
	var durationCount int
	for _, metric := range recentMetrics {
		if metric.Name == "query_average_duration" {
			totalDuration += metric.Value
			durationCount++
		}
	}

	if durationCount > 0 {
		summary["average_query_duration_ms"] = totalDuration / float64(durationCount)
	}

	summary["monitoring_active"] = m.isMonitoring
	summary["last_collection"] = m.lastCollection

	return summary, nil
}

// AnalyzePerformance provides performance analysis and recommendations
func (m *EmailPerformanceMonitor) AnalyzePerformance(ctx context.Context) (map[string]interface{}, error) {
	analysis := make(map[string]interface{})
	recommendations := make([]string, 0)

	// Analyze query performance
	queryStats, err := m.GetQueryPerformanceStats(ctx)
	if err == nil {
		if queryStats.AverageDuration > m.slowQueryThreshold {
			recommendations = append(recommendations,
				fmt.Sprintf("Average query duration (%.2fms) exceeds threshold (%.2fms). Consider query optimization.",
					queryStats.AverageDuration.Seconds()*1000, m.slowQueryThreshold.Seconds()*1000))
		}

		if queryStats.IndexUsageRatio < 80 {
			recommendations = append(recommendations,
				fmt.Sprintf("Index usage ratio (%.1f%%) is low. Consider adding or optimizing indexes.",
					queryStats.IndexUsageRatio))
		}
	}

	// Analyze table performance
	tableStats, err := m.GetTablePerformanceStats(ctx, "email_events")
	if err == nil {
		if tableStats.SequentialScans > tableStats.IndexScans {
			recommendations = append(recommendations,
				"Sequential scans exceed index scans. Review query patterns and index coverage.")
		}

		if tableStats.TotalSize > 1024*1024*1024 { // 1GB
			recommendations = append(recommendations,
				"Table size is large. Consider implementing data archival or partitioning.")
		}

		if tableStats.LastVacuum == nil || time.Since(*tableStats.LastVacuum) > 7*24*time.Hour {
			recommendations = append(recommendations,
				"Table hasn't been vacuumed recently. Consider running VACUUM ANALYZE.")
		}
	}

	analysis["recommendations"] = recommendations
	analysis["query_stats"] = queryStats
	analysis["table_stats"] = tableStats
	analysis["analysis_timestamp"] = time.Now()

	return analysis, nil
}
