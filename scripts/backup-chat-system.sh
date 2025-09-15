#!/bin/bash

# AI Consultant Live Chat System Backup Script
# This script creates backups of the chat system data

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="${BACKUP_DIR:-/var/backups/chat-system}"
ENVIRONMENT="${ENVIRONMENT:-production}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
S3_BUCKET="${S3_BUCKET:-chat-system-backups}"
AWS_REGION="${AWS_REGION:-us-east-1}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
AI Consultant Live Chat System Backup Script

Usage: $0 [OPTIONS]

Options:
    -e, --environment    Environment (development, staging, production)
    -d, --backup-dir     Backup directory (default: /var/backups/chat-system)
    -r, --retention      Retention period in days (default: 30)
    -s, --s3-bucket      S3 bucket for remote backups
    -h, --help           Show this help message

Environment Variables:
    BACKUP_DIR          Local backup directory
    RETENTION_DAYS      Number of days to retain backups
    S3_BUCKET          S3 bucket name for remote storage
    AWS_REGION         AWS region for S3 bucket

Examples:
    # Create local backup
    $0 -e production

    # Create backup with custom retention
    $0 -e production -r 60

    # Create backup and upload to S3
    $0 -e production -s my-backup-bucket
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -d|--backup-dir)
            BACKUP_DIR="$2"
            shift 2
            ;;
        -r|--retention)
            RETENTION_DAYS="$2"
            shift 2
            ;;
        -s|--s3-bucket)
            S3_BUCKET="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(development|staging|production)$ ]]; then
    log_error "Invalid environment: $ENVIRONMENT. Must be development, staging, or production."
    exit 1
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Generate backup timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="chat-system-${ENVIRONMENT}-${TIMESTAMP}"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_NAME}"

log_info "Starting backup of AI Consultant Live Chat System"
log_info "Environment: $ENVIRONMENT"
log_info "Backup Path: $BACKUP_PATH"
log_info "Retention: $RETENTION_DAYS days"

# Create backup directory
mkdir -p "$BACKUP_PATH"

# Database backup function
backup_database() {
    log_info "Creating database backup..."
    
    local db_backup_file="${BACKUP_PATH}/database.sql"
    
    if [[ "$ENVIRONMENT" == "development" ]]; then
        # Local development backup
        pg_dump -h localhost -p 5432 -U postgres -d consulting > "$db_backup_file"
    else
        # Production/staging backup (assuming Kubernetes)
        kubectl exec -n ai-consultant-chat deployment/ai-consultant-chat-postgresql -- \
            pg_dump -U consulting consulting > "$db_backup_file"
    fi
    
    if [[ $? -eq 0 ]]; then
        log_success "Database backup completed: $(du -h "$db_backup_file" | cut -f1)"
        
        # Compress database backup
        gzip "$db_backup_file"
        log_success "Database backup compressed: ${db_backup_file}.gz"
    else
        log_error "Database backup failed"
        return 1
    fi
}

# Redis backup function
backup_redis() {
    log_info "Creating Redis backup..."
    
    local redis_backup_file="${BACKUP_PATH}/redis.rdb"
    
    if [[ "$ENVIRONMENT" == "development" ]]; then
        # Local development backup
        redis-cli --rdb "$redis_backup_file"
    else
        # Production/staging backup (assuming Kubernetes)
        kubectl exec -n ai-consultant-chat deployment/ai-consultant-chat-redis-master -- \
            redis-cli BGSAVE
        
        # Wait for background save to complete
        sleep 5
        
        kubectl cp ai-consultant-chat/ai-consultant-chat-redis-master-0:/data/dump.rdb "$redis_backup_file"
    fi
    
    if [[ -f "$redis_backup_file" ]]; then
        log_success "Redis backup completed: $(du -h "$redis_backup_file" | cut -f1)"
        
        # Compress Redis backup
        gzip "$redis_backup_file"
        log_success "Redis backup compressed: ${redis_backup_file}.gz"
    else
        log_error "Redis backup failed"
        return 1
    fi
}

# Configuration backup function
backup_configuration() {
    log_info "Creating configuration backup..."
    
    local config_backup_dir="${BACKUP_PATH}/config"
    mkdir -p "$config_backup_dir"
    
    # Backup environment configurations
    if [[ -d "${PROJECT_ROOT}/config/environments" ]]; then
        cp -r "${PROJECT_ROOT}/config/environments" "$config_backup_dir/"
    fi
    
    # Backup Kubernetes manifests
    if [[ -d "${PROJECT_ROOT}/k8s" ]]; then
        cp -r "${PROJECT_ROOT}/k8s" "$config_backup_dir/"
    fi
    
    # Backup Docker configurations
    if [[ -f "${PROJECT_ROOT}/docker-compose.chat.yml" ]]; then
        cp "${PROJECT_ROOT}/docker-compose.chat.yml" "$config_backup_dir/"
    fi
    
    # Backup monitoring configurations
    if [[ -d "${PROJECT_ROOT}/monitoring" ]]; then
        cp -r "${PROJECT_ROOT}/monitoring" "$config_backup_dir/"
    fi
    
    log_success "Configuration backup completed"
}

# Application logs backup function
backup_logs() {
    log_info "Creating logs backup..."
    
    local logs_backup_dir="${BACKUP_PATH}/logs"
    mkdir -p "$logs_backup_dir"
    
    if [[ "$ENVIRONMENT" == "development" ]]; then
        # Local development logs
        if [[ -d "/var/log/chat-system" ]]; then
            cp -r /var/log/chat-system/* "$logs_backup_dir/" 2>/dev/null || true
        fi
    else
        # Production/staging logs (assuming Kubernetes)
        kubectl logs -n ai-consultant-chat deployment/ai-consultant-chat-backend --tail=10000 > "${logs_backup_dir}/backend.log" 2>/dev/null || true
        kubectl logs -n ai-consultant-chat deployment/ai-consultant-chat-frontend --tail=10000 > "${logs_backup_dir}/frontend.log" 2>/dev/null || true
    fi
    
    log_success "Logs backup completed"
}

# Create backup manifest
create_manifest() {
    log_info "Creating backup manifest..."
    
    local manifest_file="${BACKUP_PATH}/manifest.json"
    
    cat > "$manifest_file" << EOF
{
    "backup_name": "$BACKUP_NAME",
    "environment": "$ENVIRONMENT",
    "timestamp": "$TIMESTAMP",
    "created_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
    "version": "1.0.0",
    "components": {
        "database": {
            "included": true,
            "file": "database.sql.gz",
            "size": "$(stat -f%z "${BACKUP_PATH}/database.sql.gz" 2>/dev/null || echo "0")"
        },
        "redis": {
            "included": true,
            "file": "redis.rdb.gz",
            "size": "$(stat -f%z "${BACKUP_PATH}/redis.rdb.gz" 2>/dev/null || echo "0")"
        },
        "configuration": {
            "included": true,
            "directory": "config"
        },
        "logs": {
            "included": true,
            "directory": "logs"
        }
    },
    "retention_days": $RETENTION_DAYS,
    "s3_bucket": "$S3_BUCKET"
}
EOF
    
    log_success "Backup manifest created"
}

# Upload to S3 function
upload_to_s3() {
    if [[ -z "$S3_BUCKET" ]]; then
        log_info "No S3 bucket specified, skipping remote upload"
        return 0
    fi
    
    log_info "Uploading backup to S3 bucket: $S3_BUCKET"
    
    # Check if AWS CLI is available
    if ! command -v aws >/dev/null 2>&1; then
        log_error "AWS CLI is not installed. Cannot upload to S3."
        return 1
    fi
    
    # Create tar archive
    local archive_file="${BACKUP_DIR}/${BACKUP_NAME}.tar.gz"
    tar -czf "$archive_file" -C "$BACKUP_DIR" "$BACKUP_NAME"
    
    # Upload to S3
    aws s3 cp "$archive_file" "s3://${S3_BUCKET}/chat-system/${ENVIRONMENT}/" \
        --region "$AWS_REGION" \
        --storage-class STANDARD_IA
    
    if [[ $? -eq 0 ]]; then
        log_success "Backup uploaded to S3: s3://${S3_BUCKET}/chat-system/${ENVIRONMENT}/${BACKUP_NAME}.tar.gz"
        
        # Remove local archive after successful upload
        rm "$archive_file"
        log_info "Local archive removed after successful S3 upload"
    else
        log_error "Failed to upload backup to S3"
        return 1
    fi
}

# Cleanup old backups function
cleanup_old_backups() {
    log_info "Cleaning up backups older than $RETENTION_DAYS days..."
    
    # Local cleanup
    find "$BACKUP_DIR" -name "chat-system-${ENVIRONMENT}-*" -type d -mtime +$RETENTION_DAYS -exec rm -rf {} \; 2>/dev/null || true
    
    # S3 cleanup (if S3 bucket is specified)
    if [[ -n "$S3_BUCKET" ]]; then
        local cutoff_date=$(date -d "$RETENTION_DAYS days ago" +%Y-%m-%d)
        aws s3 ls "s3://${S3_BUCKET}/chat-system/${ENVIRONMENT}/" --recursive | \
        while read -r line; do
            local file_date=$(echo "$line" | awk '{print $1}')
            local file_path=$(echo "$line" | awk '{print $4}')
            
            if [[ "$file_date" < "$cutoff_date" ]]; then
                aws s3 rm "s3://${S3_BUCKET}/${file_path}"
                log_info "Removed old S3 backup: $file_path"
            fi
        done
    fi
    
    log_success "Old backups cleanup completed"
}

# Verify backup function
verify_backup() {
    log_info "Verifying backup integrity..."
    
    local errors=0
    
    # Check database backup
    if [[ -f "${BACKUP_PATH}/database.sql.gz" ]]; then
        if gzip -t "${BACKUP_PATH}/database.sql.gz"; then
            log_success "Database backup integrity verified"
        else
            log_error "Database backup is corrupted"
            ((errors++))
        fi
    else
        log_error "Database backup file not found"
        ((errors++))
    fi
    
    # Check Redis backup
    if [[ -f "${BACKUP_PATH}/redis.rdb.gz" ]]; then
        if gzip -t "${BACKUP_PATH}/redis.rdb.gz"; then
            log_success "Redis backup integrity verified"
        else
            log_error "Redis backup is corrupted"
            ((errors++))
        fi
    else
        log_error "Redis backup file not found"
        ((errors++))
    fi
    
    # Check manifest
    if [[ -f "${BACKUP_PATH}/manifest.json" ]]; then
        if python3 -m json.tool "${BACKUP_PATH}/manifest.json" >/dev/null 2>&1; then
            log_success "Backup manifest is valid JSON"
        else
            log_error "Backup manifest is invalid JSON"
            ((errors++))
        fi
    else
        log_error "Backup manifest not found"
        ((errors++))
    fi
    
    if [[ $errors -eq 0 ]]; then
        log_success "Backup verification completed successfully"
        return 0
    else
        log_error "Backup verification failed with $errors errors"
        return 1
    fi
}

# Main backup function
main() {
    local start_time=$(date +%s)
    
    # Create backups
    backup_database || exit 1
    backup_redis || exit 1
    backup_configuration
    backup_logs
    create_manifest
    
    # Verify backup
    verify_backup || exit 1
    
    # Upload to S3 if configured
    upload_to_s3
    
    # Cleanup old backups
    cleanup_old_backups
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    # Calculate backup size
    local backup_size=$(du -sh "$BACKUP_PATH" | cut -f1)
    
    log_success "Backup completed successfully!"
    log_info "Backup Size: $backup_size"
    log_info "Duration: ${duration}s"
    log_info "Backup Location: $BACKUP_PATH"
    
    # Send notification (if configured)
    if command -v curl >/dev/null 2>&1 && [[ -n "$WEBHOOK_URL" ]]; then
        curl -X POST "$WEBHOOK_URL" \
            -H "Content-Type: application/json" \
            -d "{\"text\":\"✅ Chat System Backup Completed\\nEnvironment: $ENVIRONMENT\\nSize: $backup_size\\nDuration: ${duration}s\"}" \
            >/dev/null 2>&1 || true
    fi
}

# Trap to cleanup on exit
cleanup() {
    log_info "Cleaning up temporary files..."
    # Remove any temporary files if needed
}

trap cleanup EXIT

# Run main function
main