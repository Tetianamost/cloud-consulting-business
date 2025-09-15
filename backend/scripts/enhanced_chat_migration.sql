-- Enhanced chat system database migration with optimizations
-- This migration enhances the existing chat tables with additional indexes and optimizations

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create additional indexes for enhanced performance
-- These indexes support the enhanced repository queries

-- Additional indexes for chat_sessions table
CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_created ON chat_sessions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_status_activity ON chat_sessions(status, last_activity DESC);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_expires_status ON chat_sessions(expires_at, status) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_chat_sessions_client_name_gin ON chat_sessions USING gin(to_tsvector('english', client_name)) WHERE client_name IS NOT NULL;

-- Additional indexes for chat_messages table
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_type ON chat_messages(session_id, type);
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_status ON chat_messages(session_id, status);
CREATE INDEX IF NOT EXISTS idx_chat_messages_type_created ON chat_messages(type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_status_created ON chat_messages(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_content_gin ON chat_messages USING gin(to_tsvector('english', content));

-- Partial indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_chat_sessions_active_users ON chat_sessions(user_id, last_activity DESC) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_chat_messages_recent_by_session ON chat_messages(session_id, created_at DESC) WHERE created_at > NOW() - INTERVAL '7 days';

-- Create function for efficient session cleanup with better performance
CREATE OR REPLACE FUNCTION enhanced_cleanup_expired_chat_sessions()
RETURNS TABLE(deleted_sessions INTEGER, deleted_messages INTEGER) AS $
DECLARE
    session_count INTEGER;
    message_count INTEGER;
BEGIN
    -- Delete messages first (due to foreign key constraint)
    DELETE FROM chat_messages 
    WHERE session_id IN (
        SELECT id FROM chat_sessions 
        WHERE expires_at IS NOT NULL AND expires_at < NOW()
    );
    
    GET DIAGNOSTICS message_count = ROW_COUNT;
    
    -- Delete expired sessions
    DELETE FROM chat_sessions 
    WHERE expires_at IS NOT NULL AND expires_at < NOW();
    
    GET DIAGNOSTICS session_count = ROW_COUNT;
    
    -- Return counts
    deleted_sessions := session_count;
    deleted_messages := message_count;
    
    RETURN NEXT;
END;
$ language 'plpgsql';

-- Create function for cleaning up inactive sessions
CREATE OR REPLACE FUNCTION cleanup_inactive_chat_sessions(inactive_threshold INTERVAL DEFAULT '24 hours')
RETURNS TABLE(deleted_sessions INTEGER, deleted_messages INTEGER) AS $
DECLARE
    session_count INTEGER;
    message_count INTEGER;
    cutoff_time TIMESTAMP WITH TIME ZONE;
BEGIN
    cutoff_time := NOW() - inactive_threshold;
    
    -- Delete messages first (due to foreign key constraint)
    DELETE FROM chat_messages 
    WHERE session_id IN (
        SELECT id FROM chat_sessions 
        WHERE last_activity < cutoff_time
    );
    
    GET DIAGNOSTICS message_count = ROW_COUNT;
    
    -- Delete inactive sessions
    DELETE FROM chat_sessions 
    WHERE last_activity < cutoff_time;
    
    GET DIAGNOSTICS session_count = ROW_COUNT;
    
    -- Return counts
    deleted_sessions := session_count;
    deleted_messages := message_count;
    
    RETURN NEXT;
END;
$ language 'plpgsql';

-- Create function to get session statistics
CREATE OR REPLACE FUNCTION get_chat_session_statistics()
RETURNS TABLE(
    total_sessions BIGINT,
    active_sessions BIGINT,
    inactive_sessions BIGINT,
    expired_sessions BIGINT,
    terminated_sessions BIGINT,
    total_messages BIGINT,
    avg_messages_per_session NUMERIC,
    avg_session_duration_minutes NUMERIC
) AS $
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_sessions,
        COUNT(*) FILTER (WHERE status = 'active') as active_sessions,
        COUNT(*) FILTER (WHERE status = 'inactive') as inactive_sessions,
        COUNT(*) FILTER (WHERE status = 'expired') as expired_sessions,
        COUNT(*) FILTER (WHERE status = 'terminated') as terminated_sessions,
        (SELECT COUNT(*) FROM chat_messages) as total_messages,
        CASE 
            WHEN COUNT(*) > 0 THEN (SELECT COUNT(*) FROM chat_messages)::NUMERIC / COUNT(*)::NUMERIC
            ELSE 0
        END as avg_messages_per_session,
        AVG(EXTRACT(EPOCH FROM (COALESCE(expires_at, NOW()) - created_at))/60) as avg_session_duration_minutes
    FROM chat_sessions;
END;
$ language 'plpgsql';

-- Create materialized view for session analytics (refreshed periodically)
CREATE MATERIALIZED VIEW IF NOT EXISTS chat_session_analytics AS
SELECT 
    DATE_TRUNC('hour', created_at) as hour_bucket,
    status,
    COUNT(*) as session_count,
    AVG(EXTRACT(EPOCH FROM (last_activity - created_at))/60) as avg_duration_minutes,
    COUNT(DISTINCT user_id) as unique_users
FROM chat_sessions 
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY DATE_TRUNC('hour', created_at), status
ORDER BY hour_bucket DESC, status;

-- Create index on the materialized view
CREATE INDEX IF NOT EXISTS idx_chat_session_analytics_hour_status ON chat_session_analytics(hour_bucket, status);

-- Create materialized view for message analytics
CREATE MATERIALIZED VIEW IF NOT EXISTS chat_message_analytics AS
SELECT 
    DATE_TRUNC('hour', cm.created_at) as hour_bucket,
    cm.type,
    cm.status,
    COUNT(*) as message_count,
    AVG(LENGTH(cm.content)) as avg_content_length,
    COUNT(DISTINCT cm.session_id) as unique_sessions
FROM chat_messages cm
WHERE cm.created_at > NOW() - INTERVAL '30 days'
GROUP BY DATE_TRUNC('hour', cm.created_at), cm.type, cm.status
ORDER BY hour_bucket DESC, cm.type, cm.status;

-- Create index on the message analytics materialized view
CREATE INDEX IF NOT EXISTS idx_chat_message_analytics_hour_type ON chat_message_analytics(hour_bucket, type, status);

-- Create function to refresh analytics views
CREATE OR REPLACE FUNCTION refresh_chat_analytics()
RETURNS VOID AS $
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY chat_session_analytics;
    REFRESH MATERIALIZED VIEW CONCURRENTLY chat_message_analytics;
END;
$ language 'plpgsql';

-- Create enhanced view for active sessions with more details
CREATE OR REPLACE VIEW enhanced_active_chat_sessions AS
SELECT 
    cs.id,
    cs.user_id,
    cs.client_name,
    cs.context,
    cs.status,
    cs.created_at,
    cs.last_activity,
    cs.expires_at,
    COUNT(cm.id) as message_count,
    MAX(cm.created_at) as last_message_at,
    MIN(cm.created_at) as first_message_at,
    COUNT(cm.id) FILTER (WHERE cm.type = 'user') as user_message_count,
    COUNT(cm.id) FILTER (WHERE cm.type = 'assistant') as assistant_message_count,
    COUNT(cm.id) FILTER (WHERE cm.type = 'system') as system_message_count,
    COUNT(cm.id) FILTER (WHERE cm.type = 'error') as error_message_count,
    EXTRACT(EPOCH FROM (cs.last_activity - cs.created_at))/60 as session_duration_minutes
FROM chat_sessions cs
LEFT JOIN chat_messages cm ON cs.id = cm.session_id
WHERE cs.status = 'active'
GROUP BY cs.id, cs.user_id, cs.client_name, cs.context, cs.status, cs.created_at, cs.last_activity, cs.expires_at;

-- Add table constraints for data integrity
ALTER TABLE chat_sessions 
ADD CONSTRAINT IF NOT EXISTS chk_session_dates 
CHECK (created_at <= updated_at AND created_at <= last_activity);

ALTER TABLE chat_sessions 
ADD CONSTRAINT IF NOT EXISTS chk_expires_after_created 
CHECK (expires_at IS NULL OR expires_at > created_at);

ALTER TABLE chat_messages 
ADD CONSTRAINT IF NOT EXISTS chk_content_not_empty 
CHECK (LENGTH(TRIM(content)) > 0);

ALTER TABLE chat_messages 
ADD CONSTRAINT IF NOT EXISTS chk_content_max_length 
CHECK (LENGTH(content) <= 10000);

-- Add comments for enhanced documentation
COMMENT ON INDEX idx_chat_sessions_user_created IS 'Optimizes user session listing queries ordered by creation date';
COMMENT ON INDEX idx_chat_sessions_status_activity IS 'Optimizes queries filtering by status and ordering by activity';
COMMENT ON INDEX idx_chat_messages_content_gin IS 'Enables full-text search on message content';
COMMENT ON FUNCTION enhanced_cleanup_expired_chat_sessions() IS 'Efficiently cleans up expired sessions and their messages';
COMMENT ON FUNCTION cleanup_inactive_chat_sessions(INTERVAL) IS 'Cleans up sessions inactive for longer than the specified threshold';
COMMENT ON MATERIALIZED VIEW chat_session_analytics IS 'Hourly analytics for chat sessions over the last 30 days';
COMMENT ON MATERIALIZED VIEW chat_message_analytics IS 'Hourly analytics for chat messages over the last 30 days';

-- Create a scheduled job to refresh analytics (requires pg_cron extension)
-- This is optional and requires the pg_cron extension to be installed
-- SELECT cron.schedule('refresh-chat-analytics', '0 * * * *', 'SELECT refresh_chat_analytics();');

-- Success message
SELECT 'Enhanced chat database migration completed successfully!' AS migration_result,
       (SELECT COUNT(*) FROM pg_indexes WHERE tablename IN ('chat_sessions', 'chat_messages')) AS total_indexes_created;