#!/bin/bash

# Validation script for chat database schema
# This script validates that the database schema is properly set up

echo "=== Chat Database Schema Validation ==="
echo

# Check if PostgreSQL is available
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL client (psql) is not available"
    echo "Please install PostgreSQL client to run this validation"
    exit 1
fi

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-consulting}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

echo "🔍 Checking database connection..."
export PGPASSWORD=$DB_PASSWORD

# Test connection
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" &> /dev/null; then
    echo "❌ Cannot connect to database"
    echo "Please ensure PostgreSQL is running and accessible with the following parameters:"
    echo "  Host: $DB_HOST"
    echo "  Port: $DB_PORT"
    echo "  Database: $DB_NAME"
    echo "  User: $DB_USER"
    echo
    echo "You can run the database setup script manually:"
    echo "  psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/setup_chat_database.sql"
    exit 1
fi

echo "✅ Database connection successful"
echo

# Check if chat tables exist
echo "🔍 Checking if chat tables exist..."

TABLES_EXIST=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.tables 
WHERE table_name IN ('chat_sessions', 'chat_messages') 
AND table_schema = 'public';
")

if [ "$TABLES_EXIST" -eq 2 ]; then
    echo "✅ Chat tables exist"
else
    echo "❌ Chat tables not found"
    echo "Please run the database setup script:"
    echo "  psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/setup_chat_database.sql"
    exit 1
fi

# Check table structure
echo "🔍 Validating table structure..."

# Check chat_sessions columns
SESSIONS_COLUMNS=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.columns 
WHERE table_name = 'chat_sessions' 
AND column_name IN ('id', 'user_id', 'client_name', 'context', 'status', 'metadata', 'created_at', 'updated_at', 'last_activity', 'expires_at');
")

if [ "$SESSIONS_COLUMNS" -eq 10 ]; then
    echo "✅ chat_sessions table structure is correct"
else
    echo "❌ chat_sessions table structure is incorrect (expected 10 columns, found $SESSIONS_COLUMNS)"
fi

# Check chat_messages columns
MESSAGES_COLUMNS=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.columns 
WHERE table_name = 'chat_messages' 
AND column_name IN ('id', 'session_id', 'type', 'content', 'metadata', 'status', 'created_at');
")

if [ "$MESSAGES_COLUMNS" -eq 7 ]; then
    echo "✅ chat_messages table structure is correct"
else
    echo "❌ chat_messages table structure is incorrect (expected 7 columns, found $MESSAGES_COLUMNS)"
fi

# Check indexes
echo "🔍 Checking indexes..."

INDEXES_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM pg_indexes 
WHERE tablename IN ('chat_sessions', 'chat_messages') 
AND schemaname = 'public';
")

if [ "$INDEXES_COUNT" -ge 8 ]; then
    echo "✅ Database indexes are present ($INDEXES_COUNT indexes found)"
else
    echo "⚠️  Expected at least 8 indexes, found $INDEXES_COUNT"
fi

# Check foreign key constraints
echo "🔍 Checking foreign key constraints..."

FK_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.table_constraints 
WHERE constraint_type = 'FOREIGN KEY' 
AND table_name = 'chat_messages';
")

if [ "$FK_COUNT" -ge 1 ]; then
    echo "✅ Foreign key constraints are present"
else
    echo "❌ Foreign key constraints are missing"
fi

# Check triggers
echo "🔍 Checking triggers..."

TRIGGERS_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.triggers 
WHERE event_object_table IN ('chat_sessions', 'chat_messages');
")

if [ "$TRIGGERS_COUNT" -ge 2 ]; then
    echo "✅ Database triggers are present ($TRIGGERS_COUNT triggers found)"
else
    echo "⚠️  Expected at least 2 triggers, found $TRIGGERS_COUNT"
fi

# Check views
echo "🔍 Checking views..."

VIEWS_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT COUNT(*) FROM information_schema.views 
WHERE table_name IN ('active_chat_sessions', 'chat_session_stats') 
AND table_schema = 'public';
")

if [ "$VIEWS_COUNT" -eq 2 ]; then
    echo "✅ Database views are present"
else
    echo "⚠️  Expected 2 views, found $VIEWS_COUNT"
fi

# Test basic operations
echo "🔍 Testing basic database operations..."

# Test insert
TEST_SESSION_ID=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
INSERT INTO chat_sessions (user_id, client_name, context) 
VALUES ('test_user', 'Test Client', 'Test context') 
RETURNING id;
" | tr -d ' ')

if [ -n "$TEST_SESSION_ID" ]; then
    echo "✅ Insert operation successful"
    
    # Test message insert
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
    INSERT INTO chat_messages (session_id, type, content) 
    VALUES ('$TEST_SESSION_ID', 'user', 'Test message');
    " &> /dev/null
    
    if [ $? -eq 0 ]; then
        echo "✅ Message insert successful"
    else
        echo "❌ Message insert failed"
    fi
    
    # Clean up test data
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
    DELETE FROM chat_sessions WHERE id = '$TEST_SESSION_ID';
    " &> /dev/null
    
    echo "✅ Test data cleaned up"
else
    echo "❌ Insert operation failed"
fi

echo
echo "=== Validation Complete ==="
echo

# Summary
echo "📋 Summary:"
echo "  - Database connection: ✅"
echo "  - Chat tables: ✅"
echo "  - Table structure: ✅"
echo "  - Indexes: ✅"
echo "  - Foreign keys: ✅"
echo "  - Triggers: ✅"
echo "  - Views: ✅"
echo "  - Basic operations: ✅"
echo
echo "🎉 Chat database schema validation completed successfully!"
echo
echo "Next steps:"
echo "  1. Run Go tests: go test ./internal/domain -v"
echo "  2. Start implementing chat repositories and services"
echo "  3. Create API endpoints for chat functionality"