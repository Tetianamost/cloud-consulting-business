-- Performance optimization indexes for chat system
-- This script creates indexes to optimize common chat queries

-- Chat Sessions Indexes
-- Index for session lookup by user and status
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_sessions_user_status 
ON chat_sessions (user_id, status) 
WHERE status = 'active';

-- Index for session cleanup (expired sessions)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_sessions_expires_at 
ON chat_sessions (expires_at) 
WHERE expires_at IS NOT NULL;

-- Index for session activity tracking
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_sessions_last_activity 
ON chat_sessions (last_activity DESC);

-- Composite index for user sessions with activity
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_sessions_user_activity 
ON chat_sessions (user_id, last_activity DESC);

-- Chat Messages Indexes
-- Primary index for message retrieval by session (most common query)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_session_created 
ON chat_messages (session_id, created_at DESC);

-- Index for message type filtering
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_session_type 
ON chat_messages (session_id, type, created_at DESC);

-- Index for message status updates
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_status 
ON chat_messages (status, created_at DESC);

-- Index for message search (full-text search on content)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_content_search 
ON chat_messages USING gin(to_tsvector('english', content));

-- Index for message pagination
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_session_pagination 
ON chat_messages (session_id, created_at ASC, id);

-- Composite index for filtering by session, type, and status
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_session_type_status 
ON chat_messages (session_id, type, status, created_at DESC);

-- Partial index for unread messages (performance optimization)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_unread 
ON chat_messages (session_id, created_at DESC) 
WHERE status IN ('sent', 'delivered');

-- Index for message metadata queries (if using JSONB)
-- Uncomment if metadata column is JSONB type
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chat_messages_metadata_gin 
-- ON chat_messages USING gin(metadata);

-- Performance Statistics Views
-- Create a view for session statistics
CREATE OR REPLACE VIEW chat_session_stats AS
SELECT 
    s.id as session_id,
    s.user_id,
    s.client_name,
    s.status,
    s.created_at,
    s.last_activity,
    COUNT(m.id) as message_count,
    COUNT(CASE WHEN m.type = 'user' THEN 1 END) as user_message_count,
    COUNT(CASE WHEN m.type = 'assistant' THEN 1 END) as assistant_message_count,
    MAX(m.created_at) as last_message_at,
    MIN(m.created_at) as first_message_at,
    EXTRACT(EPOCH FROM (MAX(m.created_at) - MIN(m.created_at))) as session_duration_seconds
FROM chat_sessions s
LEFT JOIN chat_messages m ON s.id = m.session_id
GROUP BY s.id, s.user_id, s.client_name, s.status, s.created_at, s.last_activity;

-- Create a view for message statistics
CREATE OR REPLACE VIEW chat_message_stats AS
SELECT 
    DATE(created_at) as message_date,
    COUNT(*) as total_messages,
    COUNT(CASE WHEN type = 'user' THEN 1 END) as user_messages,
    COUNT(CASE WHEN type = 'assistant' THEN 1 END) as assistant_messages,
    COUNT(CASE WHEN type = 'system' THEN 1 END) as system_messages,
    COUNT(CASE WHEN status = 'sent' THEN 1 END) as sent_messages,
    COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered_messages,
    COUNT(CASE WHEN status = 'read' THEN 1 END) as read_messages,
    AVG(LENGTH(content)) as avg_message_length
FROM chat_messages
GROUP BY DATE(created_at)
ORDER BY message_date DESC;

-- Optimize existing tables
-- Update table statistics
ANALYZE chat_sessions;
ANALYZE chat_messages;

-- Add comments for documentation
COMMENT ON INDEX idx_chat_sessions_user_status IS 'Optimizes user session lookup with status filtering';
COMMENT ON INDEX idx_chat_sessions_expires_at IS 'Optimizes session cleanup queries';
COMMENT ON INDEX idx_chat_sessions_last_activity IS 'Optimizes session activity tracking';
COMMENT ON INDEX idx_chat_messages_session_created IS 'Primary index for message retrieval by session';
COMMENT ON INDEX idx_chat_messages_content_search IS 'Full-text search index for message content';
COMMENT ON INDEX idx_chat_messages_unread IS 'Partial index for unread message queries';

-- Performance monitoring queries
-- Query to check index usage
CREATE OR REPLACE FUNCTION get_chat_index_usage() 
RETURNS TABLE(
    schemaname text,
    tablename text,
    indexname text,
    idx_scan bigint,
    idx_tup_read bigint,
    idx_tup_fetch bigint
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        s.schemaname::text,
        s.tablename::text,
        s.indexname::text,
        s.idx_scan,
        s.idx_tup_read,
        s.idx_tup_fetch
    FROM pg_stat_user_indexes s
    WHERE s.tablename IN ('chat_sessions', 'chat_messages')
    ORDER BY s.idx_scan DESC;
END;
$$ LANGUAGE plpgsql;

-- Query to check table statistics
CREATE OR REPLACE FUNCTION get_chat_table_stats() 
RETURNS TABLE(
    schemaname text,
    tablename text,
    n_tup_ins bigint,
    n_tup_upd bigint,
    n_tup_del bigint,
    n_live_tup bigint,
    n_dead_tup bigint,
    last_vacuum timestamp with time zone,
    last_autovacuum timestamp with time zone,
    last_analyze timestamp with time zone,
    last_autoanalyze timestamp with time zone
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        s.schemaname::text,
        s.tablename::text,
        s.n_tup_ins,
        s.n_tup_upd,
        s.n_tup_del,
        s.n_live_tup,
        s.n_dead_tup,
        s.last_vacuum,
        s.last_autovacuum,
        s.last_analyze,
        s.last_autoanalyze
    FROM pg_stat_user_tables s
    WHERE s.tablename IN ('chat_sessions', 'chat_messages');
END;
$$ LANGUAGE plpgsql;

-- Create a function to optimize chat message queries with proper pagination
CREATE OR REPLACE FUNCTION get_chat_messages_optimized(
    p_session_id text,
    p_limit integer DEFAULT 50,
    p_offset integer DEFAULT 0,
    p_message_type text DEFAULT NULL,
    p_status text DEFAULT NULL
) 
RETURNS TABLE(
    id text,
    session_id text,
    type text,
    content text,
    metadata jsonb,
    status text,
    created_at timestamp with time zone
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        m.id,
        m.session_id,
        m.type::text,
        m.content,
        m.metadata::jsonb,
        m.status::text,
        m.created_at
    FROM chat_messages m
    WHERE m.session_id = p_session_id
        AND (p_message_type IS NULL OR m.type = p_message_type)
        AND (p_status IS NULL OR m.status = p_status)
    ORDER BY m.created_at ASC
    LIMIT p_limit
    OFFSET p_offset;
END;
$$ LANGUAGE plpgsql;

-- Create a function for efficient message search
CREATE OR REPLACE FUNCTION search_chat_messages(
    p_session_id text,
    p_search_query text,
    p_limit integer DEFAULT 50
) 
RETURNS TABLE(
    id text,
    session_id text,
    type text,
    content text,
    metadata jsonb,
    status text,
    created_at timestamp with time zone,
    rank real
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        m.id,
        m.session_id,
        m.type::text,
        m.content,
        m.metadata::jsonb,
        m.status::text,
        m.created_at,
        ts_rank(to_tsvector('english', m.content), plainto_tsquery('english', p_search_query)) as rank
    FROM chat_messages m
    WHERE m.session_id = p_session_id
        AND to_tsvector('english', m.content) @@ plainto_tsquery('english', p_search_query)
    ORDER BY rank DESC, m.created_at DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;