#!/bin/bash

# Email Events Performance Optimization Script
# This script applies performance optimizations to the email events system

set -e

echo "=== Email Events Performance Optimization ==="

# Configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-cloud_consulting}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Function to check if database is accessible
check_database() {
    print_info "Checking database connection..."
    
    if ! command -v psql &> /dev/null; then
        print_error "psql command not found. Please install PostgreSQL client."
        exit 1
    fi
    
    if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" &> /dev/null; then
        print_error "Cannot connect to database. Please check your connection parameters."
        print_info "Host: $DB_HOST, Port: $DB_PORT, Database: $DB_NAME, User: $DB_USER"
        exit 1
    fi
    
    print_status "Database connection successful"
}

# Function to backup database
backup_database() {
    print_info "Creating database backup..."
    
    BACKUP_FILE="email_events_backup_$(date +%Y%m%d_%H%M%S).sql"
    
    if PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        --table=email_events --table=email_events_archive > "$BACKUP_FILE" 2>/dev/null; then
        print_status "Database backup created: $BACKUP_FILE"
    else
        print_warning "Failed to create backup, but continuing with optimization"
    fi
}

# Function to check if email_events table exists
check_email_events_table() {
    print_info "Checking if email_events table exists..."
    
    TABLE_EXISTS=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -t -c "SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'email_events');" | tr -d ' ')
    
    if [ "$TABLE_EXISTS" = "t" ]; then
        print_status "email_events table exists"
        return 0
    else
        print_error "email_events table does not exist. Please run the main email events migration first."
        exit 1
    fi
}

# Function to apply performance optimization migration
apply_optimization() {
    print_info "Applying performance optimization migration..."
    
    MIGRATION_FILE="email_events_performance_optimization.sql"
    
    if [ ! -f "$MIGRATION_FILE" ]; then
        print_error "Migration file $MIGRATION_FILE not found"
        exit 1
    fi
    
    print_info "Running performance optimization migration..."
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -f "$MIGRATION_FILE" > optimization_output.log 2>&1; then
        print_status "Performance optimization migration completed successfully"
        
        # Show summary from migration output
        if grep -q "optimization_result" optimization_output.log; then
            print_info "Migration summary:"
            grep -A 10 "optimization_result" optimization_output.log | tail -n +2
        fi
    else
        print_error "Performance optimization migration failed"
        print_info "Check optimization_output.log for details"
        exit 1
    fi
}

# Function to verify optimization
verify_optimization() {
    print_info "Verifying performance optimizations..."
    
    # Check indexes
    print_info "Checking indexes..."
    INDEX_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -t -c "SELECT COUNT(*) FROM pg_indexes WHERE tablename = 'email_events';" | tr -d ' ')
    
    if [ "$INDEX_COUNT" -gt 5 ]; then
        print_status "Found $INDEX_COUNT indexes on email_events table"
    else
        print_warning "Only $INDEX_COUNT indexes found, expected more"
    fi
    
    # Check materialized views
    print_info "Checking materialized views..."
    MATVIEW_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -t -c "SELECT COUNT(*) FROM pg_matviews WHERE matviewname LIKE 'email_metrics_%';" | tr -d ' ')
    
    if [ "$MATVIEW_COUNT" -gt 0 ]; then
        print_status "Found $MATVIEW_COUNT materialized views"
    else
        print_warning "No materialized views found"
    fi
    
    # Check functions
    print_info "Checking performance functions..."
    FUNCTION_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -t -c "SELECT COUNT(*) FROM pg_proc WHERE proname LIKE '%email%';" | tr -d ' ')
    
    if [ "$FUNCTION_COUNT" -gt 0 ]; then
        print_status "Found $FUNCTION_COUNT email-related functions"
    else
        print_warning "No email-related functions found"
    fi
    
    # Test fast metrics function
    print_info "Testing fast metrics function..."
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -c "SELECT * FROM get_email_metrics_fast(NOW() - INTERVAL '7 days', NOW()) LIMIT 1;" &> /dev/null; then
        print_status "Fast metrics function is working"
    else
        print_warning "Fast metrics function test failed"
    fi
}

# Function to refresh materialized views
refresh_materialized_views() {
    print_info "Refreshing materialized views..."
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -c "SELECT refresh_email_metrics_views();" &> /dev/null; then
        print_status "Materialized views refreshed successfully"
    else
        print_warning "Failed to refresh materialized views"
    fi
}

# Function to analyze table statistics
analyze_tables() {
    print_info "Analyzing table statistics..."
    
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -c "ANALYZE email_events;" &> /dev/null; then
        print_status "Table statistics updated"
    else
        print_warning "Failed to analyze table statistics"
    fi
}

# Function to show performance summary
show_performance_summary() {
    print_info "Performance optimization summary:"
    
    # Table size
    TABLE_SIZE=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -t -c "SELECT pg_size_pretty(pg_total_relation_size('email_events'));" | tr -d ' ')
    echo "  - Email events table size: $TABLE_SIZE"
    
    # Record count
    RECORD_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME \
        -t -c "SELECT COUNT(*) FROM email_events;" | tr -d ' ')
    echo "  - Total email events: $RECORD_COUNT"
    
    # Index count
    echo "  - Indexes created: $INDEX_COUNT"
    
    # Materialized views
    echo "  - Materialized views: $MATVIEW_COUNT"
    
    # Functions
    echo "  - Performance functions: $FUNCTION_COUNT"
    
    print_status "Performance optimization completed successfully!"
}

# Function to run performance test
run_performance_test() {
    print_info "Running performance test..."
    
    if [ -f "../test_email_performance_optimization.go" ]; then
        print_info "Performance test file found, you can run it with:"
        echo "  cd .. && go run test_email_performance_optimization.go"
    else
        print_warning "Performance test file not found"
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --backup-only     Create backup only, don't apply optimizations"
    echo "  --verify-only     Verify existing optimizations only"
    echo "  --no-backup       Skip database backup"
    echo "  --help           Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  DB_HOST          Database host (default: localhost)"
    echo "  DB_PORT          Database port (default: 5432)"
    echo "  DB_NAME          Database name (default: cloud_consulting)"
    echo "  DB_USER          Database user (default: postgres)"
    echo "  DB_PASSWORD      Database password (default: password)"
}

# Main execution
main() {
    local backup_only=false
    local verify_only=false
    local no_backup=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --backup-only)
                backup_only=true
                shift
                ;;
            --verify-only)
                verify_only=true
                shift
                ;;
            --no-backup)
                no_backup=true
                shift
                ;;
            --help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Check database connection
    check_database
    
    if [ "$backup_only" = true ]; then
        backup_database
        exit 0
    fi
    
    if [ "$verify_only" = true ]; then
        verify_optimization
        exit 0
    fi
    
    # Check prerequisites
    check_email_events_table
    
    # Create backup unless skipped
    if [ "$no_backup" = false ]; then
        backup_database
    fi
    
    # Apply optimizations
    apply_optimization
    
    # Refresh materialized views
    refresh_materialized_views
    
    # Analyze tables
    analyze_tables
    
    # Verify optimizations
    verify_optimization
    
    # Show summary
    show_performance_summary
    
    # Suggest performance test
    run_performance_test
    
    print_status "All performance optimizations completed successfully!"
    print_info "You can now run the performance test to verify improvements."
}

# Run main function with all arguments
main "$@"