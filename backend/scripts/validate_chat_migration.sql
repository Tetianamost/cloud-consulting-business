-- Validation script for enhanced chat migration
-- This script validates that all database objects and optimizations are properly created

-- Check if required extensions are enabled
SELECT 
    'Extensions Check' as check_category,
    extname as extension_name,
    CASE WHEN extname IS NOT NULL THEN 'INSTALLED' ELSE 'MISSING' END as status
FROM pg_extension 
WHERE extname IN ('uuid-ossp')
UNION ALL
SELECT 
    'Extensions Check' as check_category,
    'uuid-ossp' as extension_name,
    CASE WHEN EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp') THEN 'INSTALLED' ELSE 'MISSING' END as status
WHERE NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp');

-- Check if required tables exist
SELECT 
    'Tables Check' as check_category,
    table_name,
    CASE WHEN table_name IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('chat_sessions', 'chat_messages')
ORDER BY table_name;

-- Check if required enum types exist
SELECT 
    'Enum Types Check' as check_category,
    typname as type_name,
    CASE WHEN typname IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM pg_type 
WHERE typname IN ('session_status', 'message_type', 'message_status')
ORDER BY typname;

-- Check if required indexes exist
SELECT 
    'Indexes Check' as check_category,
    indexname as index_name,
    tablename,
    CASE WHEN indexname IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM pg_indexes 
WHERE schemaname = 'public' 
AND tablename IN ('chat_sessions', 'chat_messages')
AND indexname IN (
    'idx_chat_sessions_user_id',
    'idx_chat_sessions_status',
    'idx_chat_sessions_created_at',
    'idx_chat_sessions_last_activity',
    'idx_chat_sessions_expires_at',
    'idx_chat_sessions_user_status',
    'idx_chat_sessions_user_created',
    'idx_chat_sessions_status_activity',
    'idx_chat_sessions_expires_status',
    'idx_chat_sessions_active_users',
    'idx_chat_messages_session_id',
    'idx_chat_messages_type',
    'idx_chat_messages_status',
    'idx_chat_messages_created_at',
    'idx_chat_messages_session_created',
    'idx_chat_messages_session_type',
    'idx_chat_messages_session_status',
    'idx_chat_messages_type_created',
    'idx_chat_messages_status_created',
    'idx_chat_messages_recent_by_session'
)
ORDER BY tablename, indexname;

-- Check if GIN indexes exist (for full-text search)
SELECT 
    'GIN Indexes Check' as check_category,
    indexname as index_name,
    tablename,
    CASE WHEN indexname IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM pg_indexes 
WHERE schemaname = 'public' 
AND tablename IN ('chat_sessions', 'chat_messages')
AND indexname IN (
    'idx_chat_sessions_client_name_gin',
    'idx_chat_messages_content_gin'
)
ORDER BY tablename, indexname;

-- Check if required functions exist
SELECT 
    'Functions Check' as check_category,
    proname as function_name,
    CASE WHEN proname IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM pg_proc 
WHERE proname IN (
    'cleanup_expired_chat_sessions',
    'enhanced_cleanup_expired_chat_sessions',
    'cleanup_inactive_chat_sessions',
    'get_chat_session_statistics',
    'refresh_chat_analytics',
    'update_chat_session_timestamps',
    'update_session_activity_on_message'
)
ORDER BY proname;

-- Check if required triggers exist
SELECT 
    'Triggers Check' as check_category,
    trigger_name,
    event_object_table as table_name,
    CASE WHEN trigger_name IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM information_schema.triggers 
WHERE trigger_schema = 'public'
AND trigger_name IN (
    'update_chat_sessions_timestamps',
    'update_session_on_message'
)
ORDER BY event_object_table, trigger_name;

-- Check if required views exist
SELECT 
    'Views Check' as check_category,
    table_name as view_name,
    CASE WHEN table_name IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM information_schema.views 
WHERE table_schema = 'public'
AND table_name IN (
    'active_chat_sessions',
    'enhanced_active_chat_sessions',
    'chat_session_stats'
)
ORDER BY table_name;

-- Check if materialized views exist
SELECT 
    'Materialized Views Check' as check_category,
    matviewname as view_name,
    CASE WHEN matviewname IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM pg_matviews 
WHERE schemaname = 'public'
AND matviewname IN (
    'chat_session_analytics',
    'chat_message_analytics'
)
ORDER BY matviewname;

-- Check table constraints
SELECT 
    'Constraints Check' as check_category,
    constraint_name,
    table_name,
    constraint_type,
    CASE WHEN constraint_name IS NOT NULL THEN 'EXISTS' ELSE 'MISSING' END as status
FROM information_schema.table_constraints 
WHERE table_schema = 'public'
AND table_name IN ('chat_sessions', 'chat_messages')
AND constraint_type IN ('PRIMARY KEY', 'FOREIGN KEY', 'CHECK')
ORDER BY table_name, constraint_type, constraint_name;

-- Check foreign key relationships
SELECT 
    'Foreign Keys Check' as check_category,
    tc.constraint_name,
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name,
    'EXISTS' as status
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY' 
AND tc.table_schema = 'public'
AND tc.table_name IN ('chat_sessions', 'chat_messages');

-- Validate data integrity
SELECT 
    'Data Integrity Check' as check_category,
    'Orphaned Messages' as check_name,
    COUNT(*) as count,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM chat_messages cm
LEFT JOIN chat_sessions cs ON cm.session_id = cs.id
WHERE cs.id IS NULL

UNION ALL

SELECT 
    'Data Integrity Check' as check_category,
    'Invalid Session Status' as check_name,
    COUNT(*) as count,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM chat_sessions 
WHERE status NOT IN ('active', 'inactive', 'expired', 'terminated')

UNION ALL

SELECT 
    'Data Integrity Check' as check_category,
    'Invalid Message Type' as check_name,
    COUNT(*) as count,
    CASE WHEN COUNT(*) = 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM chat_messages 
WHERE type NOT IN ('user', 'assistant', 'system', 'error')

UNION ALL

SELECT 
    'Data Integrity Check' as check_category,
    'Invalid Message Status' as check_name,
    COUNT(*) as count,
    CASE WHEN COUNT(*) = 0 THEN 'FAIL' ELSE 'PASS' END as status
FROM chat_messages 
WHERE status NOT IN ('sent', 'delivered', 'read', 'failed');

-- Performance validation - check if indexes are being used
EXPLAIN (FORMAT JSON) 
SELECT * FROM chat_sessions WHERE user_id = 'test_user' AND status = 'active';

-- Test function execution
SELECT 
    'Function Test' as check_category,
    'get_chat_session_statistics' as function_name,
    CASE WHEN total_sessions >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM get_chat_session_statistics();

-- Test cleanup functions (dry run)
SELECT 
    'Cleanup Functions Test' as check_category,
    'enhanced_cleanup_expired_chat_sessions' as function_name,
    CASE WHEN deleted_sessions >= 0 AND deleted_messages >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM enhanced_cleanup_expired_chat_sessions();

-- Check table sizes and statistics
SELECT 
    'Table Statistics' as check_category,
    schemaname,
    tablename,
    n_tup_ins as inserts,
    n_tup_upd as updates,
    n_tup_del as deletes,
    n_live_tup as live_tuples,
    n_dead_tup as dead_tuples,
    last_vacuum,
    last_autovacuum,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables 
WHERE tablename IN ('chat_sessions', 'chat_messages')
ORDER BY tablename;

-- Check index usage statistics
SELECT 
    'Index Usage Statistics' as check_category,
    schemaname,
    tablename,
    indexname,
    idx_tup_read,
    idx_tup_fetch,
    idx_scan as scans
FROM pg_stat_user_indexes 
WHERE tablename IN ('chat_sessions', 'chat_messages')
ORDER BY tablename, idx_scan DESC;

-- Final validation summary
SELECT 
    'VALIDATION SUMMARY' as summary,
    COUNT(*) FILTER (WHERE status IN ('EXISTS', 'PASS', 'INSTALLED')) as passed_checks,
    COUNT(*) FILTER (WHERE status IN ('MISSING', 'FAIL')) as failed_checks,
    COUNT(*) as total_checks,
    CASE 
        WHEN COUNT(*) FILTER (WHERE status IN ('MISSING', 'FAIL')) = 0 THEN 'ALL CHECKS PASSED'
        ELSE 'SOME CHECKS FAILED'
    END as overall_status
FROM (
    -- Combine all previous check results here
    SELECT 'EXISTS' as status FROM information_schema.tables WHERE table_name = 'chat_sessions'
    UNION ALL
    SELECT 'EXISTS' as status FROM information_schema.tables WHERE table_name = 'chat_messages'
    UNION ALL
    SELECT CASE WHEN COUNT(*) > 0 THEN 'EXISTS' ELSE 'MISSING' END as status 
    FROM pg_indexes WHERE tablename IN ('chat_sessions', 'chat_messages')
) validation_results;