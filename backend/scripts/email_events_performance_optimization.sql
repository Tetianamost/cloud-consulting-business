-- Email Events Performance Optimization Migration
-- This migration adds advanced indexing, caching, and retention policies
-- for optimal performance with large datasets

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create additional composite indexes for complex query patterns
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_email_events_time_status_type 
ON email_events(sent_at DESC, status, email_type) 
WHERE sent_at > NOW() - INTERVAL '90 days';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_email_events_inquiry_status_time 
ON email_events(inquiry_id, status, sent_at DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_email_events_recipient_time 
ON email_events(recipient_email, sent_at DESC) 
WHERE sent_at > NOW() - INTERVAL '30 days';

-- Partial indexes for failed email analysis
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_email_events_failed_recent 
ON email_events(email_type, error_message, sent_at DESC) 
WHERE status IN ('failed', 'bounced', 'spam') 
AND sent_at > NOW() - INTERVAL '7 days';

-- Index for SES message ID lookups (webhook processing)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_email_events_ses_message_lookup 
ON email_events(ses_message_id, status, updated_at) 
WHERE ses_message_id IS NOT NULL;

-- Covering index for metrics calculations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_email_events_metrics_covering 
ON email_events(sent_at, email_type) 
INCLUDE (status, inquiry_id) 
WHERE sent_at > NOW() - INTERVAL '1 year';

-- Create materialized view for daily metrics (performance optimization)
CREATE MATERIALIZED VIEW IF NOT EXISTS email_metrics_daily AS
SELECT 
    DATE_TRUNC('day', sent_at) as metric_date,
    email_type,
    COUNT(*) as total_sent,
    COUNT(*) FILTER (WHERE status = 'delivered') as delivered_count,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_count,
    COUNT(*) FILTER (WHERE status = 'bounced') as bounced_count,
    COUNT(*) FILTER (WHERE status = 'spam') as spam_count,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'delivered'))::NUMERIC / 
        NULLIF(COUNT(*), 0) * 100, 2
    ) as delivery_rate,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'bounced'))::NUMERIC / 
        NULLIF(COUNT(*), 0) * 100, 2
    ) as bounce_rate,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'spam'))::NUMERIC / 
        NULLIF(COUNT(*), 0) * 100, 2
    ) as spam_rate,
    AVG(
        EXTRACT(EPOCH FROM (COALESCE(delivered_at, NOW()) - sent_at))/60
    ) as avg_delivery_time_minutes
FROM email_events 
WHERE sent_at >= CURRENT_DATE - INTERVAL '1 year'
GROUP BY DATE_TRUNC('day', sent_at), email_type
ORDER BY metric_date DESC, email_type;

-- Create unique index on materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_email_metrics_daily_unique 
ON email_metrics_daily(metric_date, email_type);

-- Create materialized view for hourly metrics (recent performance)
CREATE MATERIALIZED VIEW IF NOT EXISTS email_metrics_hourly AS
SELECT 
    DATE_TRUNC('hour', sent_at) as metric_hour,
    email_type,
    COUNT(*) as total_sent,
    COUNT(*) FILTER (WHERE status = 'delivered') as delivered_count,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_count,
    COUNT(*) FILTER (WHERE status = 'bounced') as bounced_count,
    COUNT(*) FILTER (WHERE status = 'spam') as spam_count,
    ROUND(
        (COUNT(*) FILTER (WHERE status = 'delivered'))::NUMERIC / 
        NULLIF(COUNT(*), 0) * 100, 2
    ) as delivery_rate
FROM email_events 
WHERE sent_at >= CURRENT_TIMESTAMP - INTERVAL '7 days'
GROUP BY DATE_TRUNC('hour', sent_at), email_type
ORDER BY metric_hour DESC, email_type;

-- Create unique index on hourly materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_email_metrics_hourly_unique 
ON email_metrics_hourly(metric_hour, email_type);

-- Function to refresh materialized views (for scheduled updates)
CREATE OR REPLACE FUNCTION refresh_email_metrics_views()
RETURNS VOID AS $
BEGIN
    -- Refresh daily metrics (full refresh)
    REFRESH MATERIALIZED VIEW CONCURRENTLY email_metrics_daily;
    
    -- Refresh hourly metrics (full refresh)
    REFRESH MATERIALIZED VIEW CONCURRENTLY email_metrics_hourly;
    
    -- Log the refresh
    INSERT INTO email_metrics_refresh_log (refreshed_at, view_name, status)
    VALUES 
        (NOW(), 'email_metrics_daily', 'success'),
        (NOW(), 'email_metrics_hourly', 'success');
        
EXCEPTION WHEN OTHERS THEN
    -- Log the error
    INSERT INTO email_metrics_refresh_log (refreshed_at, view_name, status, error_message)
    VALUES (NOW(), 'materialized_views', 'error', SQLERRM);
    RAISE;
END;
$ language 'plpgsql';

-- Create log table for materialized view refreshes
CREATE TABLE IF NOT EXISTS email_metrics_refresh_log (
    id SERIAL PRIMARY KEY,
    refreshed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    view_name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'error')),
    error_message TEXT,
    duration_ms INTEGER
);

-- Create index on refresh log
CREATE INDEX IF NOT EXISTS idx_email_metrics_refresh_log_time 
ON email_metrics_refresh_log(refreshed_at DESC);

-- Enhanced function for fast metrics calculation using materialized views
CREATE OR REPLACE FUNCTION get_email_metrics_fast(
    start_time TIMESTAMP WITH TIME ZONE DEFAULT NOW() - INTERVAL '30 days',
    end_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    email_type_filter email_event_type DEFAULT NULL
)
RETURNS TABLE(
    total_emails BIGINT,
    delivered_emails BIGINT,
    failed_emails BIGINT,
    bounced_emails BIGINT,
    spam_emails BIGINT,
    delivery_rate NUMERIC,
    bounce_rate NUMERIC,
    spam_rate NUMERIC,
    time_range TEXT
) AS $
DECLARE
    use_daily_view BOOLEAN;
    use_hourly_view BOOLEAN;
BEGIN
    -- Determine which materialized view to use based on time range
    use_daily_view := (end_time - start_time) > INTERVAL '2 days';
    use_hourly_view := (end_time - start_time) <= INTERVAL '2 days' AND 
                       start_time >= CURRENT_TIMESTAMP - INTERVAL '7 days';
    
    IF use_daily_view THEN
        -- Use daily materialized view for longer time ranges
        RETURN QUERY
        SELECT 
            SUM(total_sent) as total_emails,
            SUM(delivered_count) as delivered_emails,
            SUM(failed_count) as failed_emails,
            SUM(bounced_count) as bounced_emails,
            SUM(spam_count) as spam_emails,
            CASE 
                WHEN SUM(total_sent) > 0 THEN 
                    ROUND(SUM(delivered_count)::NUMERIC / SUM(total_sent)::NUMERIC * 100, 2)
                ELSE 0
            END as delivery_rate,
            CASE 
                WHEN SUM(total_sent) > 0 THEN 
                    ROUND(SUM(bounced_count)::NUMERIC / SUM(total_sent)::NUMERIC * 100, 2)
                ELSE 0
            END as bounce_rate,
            CASE 
                WHEN SUM(total_sent) > 0 THEN 
                    ROUND(SUM(spam_count)::NUMERIC / SUM(total_sent)::NUMERIC * 100, 2)
                ELSE 0
            END as spam_rate,
            CONCAT(start_time::DATE, ' to ', end_time::DATE) as time_range
        FROM email_metrics_daily 
        WHERE metric_date >= start_time::DATE 
        AND metric_date <= end_time::DATE
        AND (email_type_filter IS NULL OR email_type = email_type_filter);
        
    ELSIF use_hourly_view THEN
        -- Use hourly materialized view for recent short time ranges
        RETURN QUERY
        SELECT 
            SUM(total_sent) as total_emails,
            SUM(delivered_count) as delivered_emails,
            SUM(failed_count) as failed_emails,
            SUM(bounced_count) as bounced_emails,
            SUM(spam_count) as spam_emails,
            CASE 
                WHEN SUM(total_sent) > 0 THEN 
                    ROUND(SUM(delivered_count)::NUMERIC / SUM(total_sent)::NUMERIC * 100, 2)
                ELSE 0
            END as delivery_rate,
            CASE 
                WHEN SUM(total_sent) > 0 THEN 
                    ROUND(SUM(bounced_count)::NUMERIC / SUM(total_sent)::NUMERIC * 100, 2)
                ELSE 0
            END as bounce_rate,
            CASE 
                WHEN SUM(total_sent) > 0 THEN 
                    ROUND(SUM(spam_count)::NUMERIC / SUM(total_sent)::NUMERIC * 100, 2)
                ELSE 0
            END as spam_rate,
            CONCAT(start_time::DATE, ' to ', end_time::DATE) as time_range
        FROM email_metrics_hourly 
        WHERE metric_hour >= start_time 
        AND metric_hour <= end_time
        AND (email_type_filter IS NULL OR email_type = email_type_filter);
        
    ELSE
        -- Fall back to real-time calculation for edge cases
        RETURN QUERY
        SELECT 
            COUNT(*) as total_emails,
            COUNT(*) FILTER (WHERE status = 'delivered') as delivered_emails,
            COUNT(*) FILTER (WHERE status = 'failed') as failed_emails,
            COUNT(*) FILTER (WHERE status = 'bounced') as bounced_emails,
            COUNT(*) FILTER (WHERE status = 'spam') as spam_emails,
            CASE 
                WHEN COUNT(*) > 0 THEN 
                    ROUND((COUNT(*) FILTER (WHERE status = 'delivered'))::NUMERIC / COUNT(*)::NUMERIC * 100, 2)
                ELSE 0
            END as delivery_rate,
            CASE 
                WHEN COUNT(*) > 0 THEN 
                    ROUND((COUNT(*) FILTER (WHERE status = 'bounced'))::NUMERIC / COUNT(*)::NUMERIC * 100, 2)
                ELSE 0
            END as bounce_rate,
            CASE 
                WHEN COUNT(*) > 0 THEN 
                    ROUND((COUNT(*) FILTER (WHERE status = 'spam'))::NUMERIC / COUNT(*)::NUMERIC * 100, 2)
                ELSE 0
            END as spam_rate,
            CONCAT(start_time::DATE, ' to ', end_time::DATE) as time_range
        FROM email_events 
        WHERE sent_at >= start_time 
        AND sent_at <= end_time
        AND (email_type_filter IS NULL OR email_type = email_type_filter);
    END IF;
END;
$ language 'plpgsql';

-- Enhanced cleanup function with performance optimizations
CREATE OR REPLACE FUNCTION cleanup_old_email_events_optimized(
    retention_days INTEGER DEFAULT 365,
    batch_size INTEGER DEFAULT 1000
)
RETURNS TABLE(
    deleted_count INTEGER,
    batches_processed INTEGER,
    duration_seconds NUMERIC
) AS $
DECLARE
    cutoff_date TIMESTAMP WITH TIME ZONE;
    batch_count INTEGER := 0;
    total_deleted INTEGER := 0;
    current_batch_deleted INTEGER;
    start_time TIMESTAMP WITH TIME ZONE;
BEGIN
    start_time := NOW();
    cutoff_date := NOW() - (retention_days || ' days')::INTERVAL;
    
    -- Process in batches to avoid long-running transactions
    LOOP
        DELETE FROM email_events 
        WHERE id IN (
            SELECT id FROM email_events 
            WHERE sent_at < cutoff_date 
            ORDER BY sent_at 
            LIMIT batch_size
        );
        
        GET DIAGNOSTICS current_batch_deleted = ROW_COUNT;
        
        IF current_batch_deleted = 0 THEN
            EXIT;
        END IF;
        
        total_deleted := total_deleted + current_batch_deleted;
        batch_count := batch_count + 1;
        
        -- Commit batch and brief pause to avoid blocking
        COMMIT;
        PERFORM pg_sleep(0.1);
    END LOOP;
    
    -- Refresh materialized views after cleanup
    PERFORM refresh_email_metrics_views();
    
    RETURN QUERY SELECT 
        total_deleted,
        batch_count,
        EXTRACT(EPOCH FROM (NOW() - start_time))::NUMERIC;
END;
$ language 'plpgsql';

-- Create table for email event archival (cold storage)
CREATE TABLE IF NOT EXISTS email_events_archive (
    LIKE email_events INCLUDING ALL
);

-- Partition the archive table by year for better performance
CREATE TABLE IF NOT EXISTS email_events_archive_2024 
PARTITION OF email_events_archive 
FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE IF NOT EXISTS email_events_archive_2025 
PARTITION OF email_events_archive 
FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

-- Function to archive old email events instead of deleting them
CREATE OR REPLACE FUNCTION archive_old_email_events(
    archive_days INTEGER DEFAULT 90,
    batch_size INTEGER DEFAULT 1000
)
RETURNS TABLE(
    archived_count INTEGER,
    batches_processed INTEGER,
    duration_seconds NUMERIC
) AS $
DECLARE
    cutoff_date TIMESTAMP WITH TIME ZONE;
    batch_count INTEGER := 0;
    total_archived INTEGER := 0;
    current_batch_archived INTEGER;
    start_time TIMESTAMP WITH TIME ZONE;
BEGIN
    start_time := NOW();
    cutoff_date := NOW() - (archive_days || ' days')::INTERVAL;
    
    -- Process in batches
    LOOP
        -- Move old records to archive
        WITH archived_events AS (
            DELETE FROM email_events 
            WHERE id IN (
                SELECT id FROM email_events 
                WHERE sent_at < cutoff_date 
                ORDER BY sent_at 
                LIMIT batch_size
            )
            RETURNING *
        )
        INSERT INTO email_events_archive 
        SELECT * FROM archived_events;
        
        GET DIAGNOSTICS current_batch_archived = ROW_COUNT;
        
        IF current_batch_archived = 0 THEN
            EXIT;
        END IF;
        
        total_archived := total_archived + current_batch_archived;
        batch_count := batch_count + 1;
        
        -- Commit batch and brief pause
        COMMIT;
        PERFORM pg_sleep(0.1);
    END LOOP;
    
    -- Refresh materialized views after archival
    PERFORM refresh_email_metrics_views();
    
    RETURN QUERY SELECT 
        total_archived,
        batch_count,
        EXTRACT(EPOCH FROM (NOW() - start_time))::NUMERIC;
END;
$ language 'plpgsql';

-- Create function to analyze email event query performance
CREATE OR REPLACE FUNCTION analyze_email_event_performance()
RETURNS TABLE(
    query_type TEXT,
    avg_duration_ms NUMERIC,
    total_calls BIGINT,
    table_size_mb NUMERIC,
    index_usage_ratio NUMERIC
) AS $
BEGIN
    RETURN QUERY
    WITH table_stats AS (
        SELECT 
            pg_size_pretty(pg_total_relation_size('email_events'))::TEXT as size_pretty,
            (pg_total_relation_size('email_events') / 1024.0 / 1024.0)::NUMERIC as size_mb
    ),
    query_stats AS (
        SELECT 
            'email_events_queries' as query_type,
            avg(mean_exec_time)::NUMERIC as avg_duration,
            sum(calls)::BIGINT as total_calls
        FROM pg_stat_statements 
        WHERE query LIKE '%email_events%'
    ),
    index_stats AS (
        SELECT 
            CASE 
                WHEN (idx_scan + seq_scan) > 0 THEN 
                    (idx_scan::NUMERIC / (idx_scan + seq_scan)::NUMERIC * 100)
                ELSE 0 
            END as usage_ratio
        FROM pg_stat_user_tables 
        WHERE relname = 'email_events'
    )
    SELECT 
        qs.query_type,
        qs.avg_duration,
        qs.total_calls,
        ts.size_mb,
        ist.usage_ratio
    FROM query_stats qs, table_stats ts, index_stats ist;
END;
$ language 'plpgsql';

-- Create automated maintenance job scheduling table
CREATE TABLE IF NOT EXISTS email_maintenance_schedule (
    id SERIAL PRIMARY KEY,
    job_name VARCHAR(100) NOT NULL UNIQUE,
    job_function VARCHAR(200) NOT NULL,
    schedule_cron VARCHAR(50) NOT NULL,
    last_run TIMESTAMP WITH TIME ZONE,
    next_run TIMESTAMP WITH TIME ZONE,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert default maintenance jobs
INSERT INTO email_maintenance_schedule (job_name, job_function, schedule_cron, next_run) 
VALUES 
    ('refresh_daily_metrics', 'refresh_email_metrics_views()', '0 1 * * *', NOW() + INTERVAL '1 day'),
    ('cleanup_old_events', 'cleanup_old_email_events_optimized(365)', '0 2 * * 0', NOW() + INTERVAL '7 days'),
    ('archive_old_events', 'archive_old_email_events(90)', '0 3 * * 0', NOW() + INTERVAL '7 days')
ON CONFLICT (job_name) DO NOTHING;

-- Create performance monitoring view
CREATE OR REPLACE VIEW email_performance_summary AS
SELECT 
    'email_events' as table_name,
    (SELECT COUNT(*) FROM email_events) as total_records,
    (SELECT COUNT(*) FROM email_events WHERE sent_at > NOW() - INTERVAL '24 hours') as records_last_24h,
    (SELECT COUNT(*) FROM email_events WHERE sent_at > NOW() - INTERVAL '7 days') as records_last_7d,
    (SELECT COUNT(*) FROM email_events WHERE sent_at > NOW() - INTERVAL '30 days') as records_last_30d,
    pg_size_pretty(pg_total_relation_size('email_events')) as table_size,
    (SELECT COUNT(*) FROM pg_indexes WHERE tablename = 'email_events') as index_count,
    (SELECT schemaname||'.'||tablename as table_name, 
            attname as column_name, 
            n_distinct, 
            correlation 
     FROM pg_stats 
     WHERE tablename = 'email_events' 
     ORDER BY n_distinct DESC 
     LIMIT 1) as most_selective_column;

-- Add table and function comments
COMMENT ON MATERIALIZED VIEW email_metrics_daily IS 'Pre-calculated daily email metrics for fast dashboard queries';
COMMENT ON MATERIALIZED VIEW email_metrics_hourly IS 'Pre-calculated hourly email metrics for recent performance monitoring';
COMMENT ON FUNCTION get_email_metrics_fast(TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH TIME ZONE, email_event_type) IS 'Optimized metrics calculation using materialized views';
COMMENT ON FUNCTION cleanup_old_email_events_optimized(INTEGER, INTEGER) IS 'Batch-based cleanup of old email events with performance optimization';
COMMENT ON FUNCTION archive_old_email_events(INTEGER, INTEGER) IS 'Archive old email events to separate table instead of deletion';
COMMENT ON FUNCTION refresh_email_metrics_views() IS 'Refresh all materialized views for email metrics';
COMMENT ON TABLE email_events_archive IS 'Archive table for old email events (cold storage)';
COMMENT ON VIEW email_performance_summary IS 'Performance monitoring summary for email events system';

-- Success message
SELECT 
    'Email events performance optimization completed!' AS optimization_result,
    (SELECT COUNT(*) FROM pg_indexes WHERE tablename = 'email_events') AS total_indexes,
    (SELECT COUNT(*) FROM information_schema.views WHERE table_name LIKE '%email%') AS views_created,
    (SELECT COUNT(*) FROM information_schema.routines WHERE routine_name LIKE '%email%') AS functions_created;