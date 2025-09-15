-- Test script for email events migration
-- This script validates the migration without making permanent changes

-- Start a transaction that we'll rollback
BEGIN;

-- Test the migration by running it
\i email_events_migration.sql

-- Verify the table was created
SELECT 
    table_name,
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'email_events'
ORDER BY ordinal_position;

-- Verify indexes were created
SELECT 
    indexname,
    indexdef
FROM pg_indexes 
WHERE tablename = 'email_events'
ORDER BY indexname;

-- Verify constraints were created
SELECT 
    constraint_name,
    constraint_type
FROM information_schema.table_constraints 
WHERE table_name = 'email_events'
ORDER BY constraint_name;

-- Verify functions were created
SELECT 
    routine_name,
    routine_type
FROM information_schema.routines 
WHERE routine_name LIKE '%email%'
ORDER BY routine_name;

-- Test the email metrics function
SELECT * FROM get_email_metrics();

-- Test inserting a sample record
INSERT INTO email_events (
    inquiry_id,
    email_type,
    recipient_email,
    sender_email,
    subject,
    status,
    sent_at
) VALUES (
    'test-inquiry-123',
    'customer_confirmation',
    'test@example.com',
    'info@cloudpartner.pro',
    'Test Email Subject',
    'sent',
    NOW()
);

-- Verify the record was inserted
SELECT COUNT(*) as record_count FROM email_events;

-- Test the inquiry status function
SELECT * FROM get_email_status_by_inquiry('test-inquiry-123');

-- Test the views
SELECT * FROM email_event_stats;
SELECT * FROM daily_email_metrics;

-- Rollback all changes (this is a test)
ROLLBACK;

-- Confirm rollback worked
SELECT COUNT(*) as tables_remaining 
FROM information_schema.tables 
WHERE table_name = 'email_events';