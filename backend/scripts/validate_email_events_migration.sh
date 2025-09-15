#!/bin/bash

# Email events migration validation script
# This script validates the SQL syntax and structure of the migration files

echo "=== Email Events Migration Validation ==="
echo

# Check if migration files exist
MIGRATION_FILE="scripts/email_events_migration.sql"
ROLLBACK_FILE="scripts/email_events_rollback.sql"
TEST_FILE="scripts/test_email_events_migration.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo "❌ Migration file not found: $MIGRATION_FILE"
    exit 1
fi

if [ ! -f "$ROLLBACK_FILE" ]; then
    echo "❌ Rollback file not found: $ROLLBACK_FILE"
    exit 1
fi

if [ ! -f "$TEST_FILE" ]; then
    echo "❌ Test file not found: $TEST_FILE"
    exit 1
fi

echo "✅ All migration files found"
echo

# Basic syntax validation (if psql is available)
if command -v psql >/dev/null 2>&1; then
    echo "🔍 Validating SQL syntax..."
    
    # Check migration file syntax
    if psql --set ON_ERROR_STOP=1 --quiet --no-psqlrc -f "$MIGRATION_FILE" --dry-run 2>/dev/null; then
        echo "✅ Migration SQL syntax is valid"
    else
        echo "❌ Migration SQL syntax has errors"
        echo "   Run: psql -f $MIGRATION_FILE --dry-run"
        exit 1
    fi
    
    # Check rollback file syntax
    if psql --set ON_ERROR_STOP=1 --quiet --no-psqlrc -f "$ROLLBACK_FILE" --dry-run 2>/dev/null; then
        echo "✅ Rollback SQL syntax is valid"
    else
        echo "❌ Rollback SQL syntax has errors"
        echo "   Run: psql -f $ROLLBACK_FILE --dry-run"
        exit 1
    fi
else
    echo "⚠️  psql not available, skipping syntax validation"
fi

echo

# Check file structure and content
echo "🔍 Validating migration structure..."

# Check for required elements in migration file
REQUIRED_ELEMENTS=(
    "CREATE TABLE.*email_events"
    "CREATE INDEX.*email_events"
    "CREATE TYPE.*email_event_type"
    "CREATE TYPE.*email_event_status"
    "CREATE FUNCTION.*get_email_metrics"
    "CREATE FUNCTION.*get_email_status_by_inquiry"
    "CREATE VIEW.*email_event_stats"
)

for element in "${REQUIRED_ELEMENTS[@]}"; do
    if grep -q "$element" "$MIGRATION_FILE"; then
        echo "✅ Found: $element"
    else
        echo "❌ Missing: $element"
        exit 1
    fi
done

echo

# Check rollback file has corresponding DROP statements
REQUIRED_DROPS=(
    "DROP TABLE.*email_events"
    "DROP TYPE.*email_event_type"
    "DROP TYPE.*email_event_status"
    "DROP FUNCTION.*get_email_metrics"
    "DROP VIEW.*email_event_stats"
)

echo "🔍 Validating rollback structure..."

for drop in "${REQUIRED_DROPS[@]}"; do
    if grep -q "$drop" "$ROLLBACK_FILE"; then
        echo "✅ Found: $drop"
    else
        echo "❌ Missing: $drop"
        exit 1
    fi
done

echo

# Check for best practices
echo "🔍 Checking best practices..."

# Check for IF NOT EXISTS
if grep -q "IF NOT EXISTS" "$MIGRATION_FILE"; then
    echo "✅ Uses IF NOT EXISTS for safe execution"
else
    echo "⚠️  Consider using IF NOT EXISTS for safer migrations"
fi

# Check for proper indexing
INDEX_COUNT=$(grep -c "CREATE INDEX" "$MIGRATION_FILE")
if [ "$INDEX_COUNT" -ge 4 ]; then
    echo "✅ Has adequate indexing ($INDEX_COUNT indexes)"
else
    echo "⚠️  Consider adding more indexes for performance"
fi

# Check for constraints
if grep -q "ADD CONSTRAINT" "$MIGRATION_FILE"; then
    echo "✅ Includes data integrity constraints"
else
    echo "⚠️  Consider adding constraints for data integrity"
fi

# Check for comments
if grep -q "COMMENT ON" "$MIGRATION_FILE"; then
    echo "✅ Includes documentation comments"
else
    echo "⚠️  Consider adding comments for documentation"
fi

echo

# File size check
MIGRATION_SIZE=$(wc -l < "$MIGRATION_FILE")
ROLLBACK_SIZE=$(wc -l < "$ROLLBACK_FILE")

echo "📊 Migration Statistics:"
echo "   Migration file: $MIGRATION_SIZE lines"
echo "   Rollback file: $ROLLBACK_SIZE lines"
echo "   Tables created: $(grep -c "CREATE TABLE" "$MIGRATION_FILE")"
echo "   Indexes created: $(grep -c "CREATE INDEX" "$MIGRATION_FILE")"
echo "   Functions created: $(grep -c "CREATE.*FUNCTION" "$MIGRATION_FILE")"
echo "   Views created: $(grep -c "CREATE.*VIEW" "$MIGRATION_FILE")"

echo
echo "✅ Email events migration validation completed successfully!"
echo
echo "Next steps:"
echo "1. Review the migration files"
echo "2. Test with: docker-compose --profile database up -d db"
echo "3. Run migration: psql -h localhost -U postgres -d consulting -f $MIGRATION_FILE"
echo "4. Test rollback: psql -h localhost -U postgres -d consulting -f $ROLLBACK_FILE"