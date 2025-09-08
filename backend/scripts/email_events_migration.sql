-- Email events tracking database migration
-- This migration adds email_events table for real email monitoring data
-- Replaces mock data with actual email delivery tracking

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for better type safety and performance
CREATE TYPE email_event_type AS ENUM ('customer_confirmation', 'consultant_notification', 'inquiry_notification');
CREATE TYPE email_event_status AS ENUM ('sent', 'delivered', 'failed', 'bounced', 'spam');
CREATE TYPE bounce_type AS ENUM ('permanent', 'temporary', 'complaint');

-- Create email_events table
CREATE TABLE IF NOT EXISTS email_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inquiry_id VARCHAR(255) NOT NULL,
    email_type email_event_type NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    sender_email VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    status email_event_status NOT NULL DEFAULT 'sent',
    sent_at TIMESTAMP WITH TIME ZONE NOT NULL,
    delivered_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    bounce_type bounce_type,
    ses_message_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for optimal query performance
CREATE INDEX IF NOT EXISTS idx_email_events_inquiry_id ON email_events(inquiry_id);
CREATE INDEX IF NOT EXISTS idx_email_events_status ON email_events(status);
CREATE INDEX IF NOT EXISTS idx_email_events_sent_at ON email_events(sent_at);
CREATE INDEX IF NOT EXISTS idx_email_events_email_type ON email_events(email_type);

-- Additional composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_email_events_inquiry_type ON email_events(inquiry_id, email_type);
CREATE INDEX IF NOT EXISTS idx_email_events_status_sent ON email_events(status, sent_at DESC);
CREATE INDEX IF NOT EXISTS idx_email_events_type_status ON email_events(email_type, status);
CREATE INDEX IF NOT EXISTS idx_email_events_ses_message_id ON email_events(ses_message_id) WHERE ses_message_id IS NOT NULL;

-- Partial indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_email_events_failed_status ON email_events(sent_at DESC, error_message) WHERE status IN ('failed', 'bounced', 'spam');
CREATE INDEX IF NOT EXISTS idx_email_events_recent_events ON email_events(sent_at DESC, status) WHERE sent_at > NOW() - INTERVAL '30 days';

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_email_events_updated_at()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$ language 'plpgsql';

CREATE TRIGGER update_email_events_updated_at 
    BEFORE UPDATE ON email_events 
    FOR EACH ROW EXECUTE FUNCTION update_email_events_updated_at();

-- Add foreign key constraint to inquiries table if it exists
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'inquiries') THEN
        ALTER TABLE email_events 
        ADD CONSTRAINT IF NOT EXISTS fk_email_events_inquiry_id 
        FOREIGN KEY (inquiry_id) REFERENCES inquiries(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Add table constraints for data integrity
ALTER TABLE email_events 
ADD CONSTRAINT IF NOT EXISTS chk_email_events_dates 
CHECK (sent_at <= COALESCE(delivered_at, sent_at) AND created_at <= updated_at);

ALTER TABLE email_events 
ADD CONSTRAINT IF NOT EXISTS chk_email_events_recipient_email 
CHECK (recipient_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE email_events 
ADD CONSTRAINT IF NOT EXISTS chk_email_events_sender_email 
CHECK (sender_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE email_events 
ADD CONSTRAINT IF NOT EXISTS chk_email_events_subject_length 
CHECK (LENGTH(TRIM(COALESCE(subject, ''))) <= 500);

ALTER TABLE email_events 
ADD CONSTRAINT IF NOT EXISTS chk_email_events_error_message_length 
CHECK (LENGTH(COALESCE(error_message, '')) <= 10000);

-- Create function to get email metrics for a time range
CREATE OR REPLACE FUNCTION get_email_metrics(
    start_time TIMESTAMP WITH TIME ZONE DEFAULT NOW() - INTERVAL '30 days',
    end_time TIMESTAMP WITH TIME ZONE DEFAULT NOW()
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
BEGIN
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
    WHERE sent_at >= start_time AND sent_at <= end_time;
END;
$ language 'plpgsql';

-- Create function to get email status by inquiry
CREATE OR REPLACE FUNCTION get_email_status_by_inquiry(inquiry_uuid VARCHAR(255))
RETURNS TABLE(
    inquiry_id VARCHAR(255),
    email_type email_event_type,
    recipient_email VARCHAR(255),
    status email_event_status,
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT
) AS $
BEGIN
    RETURN QUERY
    SELECT 
        ee.inquiry_id,
        ee.email_type,
        ee.recipient_email,
        ee.status,
        ee.sent_at,
        ee.delivered_at,
        ee.error_message
    FROM email_events ee
    WHERE ee.inquiry_id = inquiry_uuid
    ORDER BY ee.sent_at DESC;
END;
$ language 'plpgsql';

-- Create function to cleanup old email events (retention policy)
CREATE OR REPLACE FUNCTION cleanup_old_email_events(retention_days INTEGER DEFAULT 365)
RETURNS INTEGER AS $
DECLARE
    deleted_count INTEGER;
    cutoff_date TIMESTAMP WITH TIME ZONE;
BEGIN
    cutoff_date := NOW() - (retention_days || ' days')::INTERVAL;
    
    DELETE FROM email_events 
    WHERE sent_at < cutoff_date;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$ language 'plpgsql';

-- Create view for email event statistics
CREATE OR REPLACE VIEW email_event_stats AS
SELECT 
    email_type,
    status,
    COUNT(*) as event_count,
    AVG(EXTRACT(EPOCH FROM (COALESCE(delivered_at, NOW()) - sent_at))/60) as avg_delivery_time_minutes,
    MIN(sent_at) as earliest_event,
    MAX(sent_at) as latest_event
FROM email_events 
WHERE sent_at > NOW() - INTERVAL '30 days'
GROUP BY email_type, status
ORDER BY email_type, status;

-- Create view for daily email metrics
CREATE OR REPLACE VIEW daily_email_metrics AS
SELECT 
    DATE_TRUNC('day', sent_at) as date,
    email_type,
    COUNT(*) as total_sent,
    COUNT(*) FILTER (WHERE status = 'delivered') as delivered,
    COUNT(*) FILTER (WHERE status = 'failed') as failed,
    COUNT(*) FILTER (WHERE status = 'bounced') as bounced,
    COUNT(*) FILTER (WHERE status = 'spam') as spam,
    CASE 
        WHEN COUNT(*) > 0 THEN 
            ROUND((COUNT(*) FILTER (WHERE status = 'delivered'))::NUMERIC / COUNT(*)::NUMERIC * 100, 2)
        ELSE 0
    END as delivery_rate_percent
FROM email_events 
WHERE sent_at > NOW() - INTERVAL '90 days'
GROUP BY DATE_TRUNC('day', sent_at), email_type
ORDER BY date DESC, email_type;

-- Add table and column comments for documentation
COMMENT ON TABLE email_events IS 'Stores email event tracking data for real email monitoring instead of mock data';
COMMENT ON COLUMN email_events.inquiry_id IS 'Reference to the inquiry that triggered the email';
COMMENT ON COLUMN email_events.email_type IS 'Type of email: customer_confirmation, consultant_notification, or inquiry_notification';
COMMENT ON COLUMN email_events.recipient_email IS 'Email address of the recipient';
COMMENT ON COLUMN email_events.sender_email IS 'Email address of the sender (usually system email)';
COMMENT ON COLUMN email_events.subject IS 'Email subject line';
COMMENT ON COLUMN email_events.status IS 'Current delivery status of the email';
COMMENT ON COLUMN email_events.sent_at IS 'Timestamp when the email was sent';
COMMENT ON COLUMN email_events.delivered_at IS 'Timestamp when the email was delivered (if available)';
COMMENT ON COLUMN email_events.error_message IS 'Error message if email delivery failed';
COMMENT ON COLUMN email_events.bounce_type IS 'Type of bounce if email bounced';
COMMENT ON COLUMN email_events.ses_message_id IS 'AWS SES message ID for tracking';

COMMENT ON FUNCTION get_email_metrics(TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH TIME ZONE) IS 'Returns comprehensive email metrics for a given time range';
COMMENT ON FUNCTION get_email_status_by_inquiry(VARCHAR) IS 'Returns all email events for a specific inquiry';
COMMENT ON FUNCTION cleanup_old_email_events(INTEGER) IS 'Removes email events older than specified retention period';

COMMENT ON VIEW email_event_stats IS 'Statistics view showing email event counts and delivery times by type and status';
COMMENT ON VIEW daily_email_metrics IS 'Daily aggregated email metrics for the last 90 days';

-- Success message with statistics
SELECT 
    'Email events migration completed successfully!' AS migration_result,
    (SELECT COUNT(*) FROM pg_indexes WHERE tablename = 'email_events') AS indexes_created,
    (SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_name = 'email_events') AS constraints_created,
    (SELECT COUNT(*) FROM information_schema.routines WHERE routine_name LIKE '%email%') AS functions_created;