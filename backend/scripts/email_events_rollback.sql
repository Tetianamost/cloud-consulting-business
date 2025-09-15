-- Email events tracking rollback migration
-- This script safely removes the email events tracking system
-- Use this to rollback the email_events_migration.sql changes

-- Drop views first (dependent objects)
DROP VIEW IF EXISTS daily_email_metrics;
DROP VIEW IF EXISTS email_event_stats;

-- Drop functions
DROP FUNCTION IF EXISTS cleanup_old_email_events(INTEGER);
DROP FUNCTION IF EXISTS get_email_status_by_inquiry(VARCHAR);
DROP FUNCTION IF EXISTS get_email_metrics(TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH TIME ZONE);
DROP FUNCTION IF EXISTS update_email_events_updated_at();

-- Drop triggers
DROP TRIGGER IF EXISTS update_email_events_updated_at ON email_events;

-- Drop indexes (will be automatically dropped with table, but explicit for clarity)
DROP INDEX IF EXISTS idx_email_events_recent_events;
DROP INDEX IF EXISTS idx_email_events_failed_status;
DROP INDEX IF EXISTS idx_email_events_ses_message_id;
DROP INDEX IF EXISTS idx_email_events_type_status;
DROP INDEX IF EXISTS idx_email_events_status_sent;
DROP INDEX IF EXISTS idx_email_events_inquiry_type;
DROP INDEX IF EXISTS idx_email_events_email_type;
DROP INDEX IF EXISTS idx_email_events_sent_at;
DROP INDEX IF EXISTS idx_email_events_status;
DROP INDEX IF EXISTS idx_email_events_inquiry_id;

-- Drop table (this will also drop all constraints and indexes)
DROP TABLE IF EXISTS email_events;

-- Drop enum types
DROP TYPE IF EXISTS bounce_type;
DROP TYPE IF EXISTS email_event_status;
DROP TYPE IF EXISTS email_event_type;

-- Success message
SELECT 'Email events tracking system rollback completed successfully!' AS rollback_result,
       'All email events data has been removed.' AS warning_message;