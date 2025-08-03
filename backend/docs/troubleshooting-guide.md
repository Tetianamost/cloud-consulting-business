# AI Consultant Live Chat Troubleshooting Guide

## Overview

This guide provides solutions for common issues encountered with the AI Consultant Live Chat system. Issues are organized by category with step-by-step resolution procedures.

## Quick Diagnostic Commands

### System Health Check
```bash
# Check all services status
docker-compose ps

# Check backend health
curl http://localhost:8080/health

# Check WebSocket connectivity
wscat -c ws://localhost:8080/api/v1/admin/chat/ws?token=your_token

# Check database connectivity
psql -h localhost -U chat_user -d chat_dev -c "SELECT 1;"

# Check Redis connectivity
redis-cli ping
```

### Log Analysis
```bash
# View backend logs
docker-compose logs -f backend

# View specific error logs
docker-compose logs backend | grep ERROR

# View WebSocket connection logs
docker-compose logs backend | grep "websocket"

# View database logs
docker-compose logs postgres
```

## Connection Issues

### WebSocket Connection Failures

#### Symptoms
- Red connection indicator in UI
- "Connection failed" error messages
- Messages not sending or receiving
- Frequent disconnections

#### Diagnostic Steps
```bash
# 1. Check WebSocket endpoint
curl -I http://localhost:8080/api/v1/admin/chat/ws

# 2. Test with wscat
wscat -c ws://localhost:8080/api/v1/admin/chat/ws?token=your_token

# 3. Check nginx configuration (if using nginx)
nginx -t
cat /etc/nginx/sites-available/chat

# 4. Check firewall rules
sudo ufw status
netstat -tlnp | grep 8080
```

#### Solutions

**Solution 1: Authentication Issues**
```bash
# Verify JWT token is valid
echo "your_token" | base64 -d

# Check token expiration
curl -H "Authorization: Bearer your_token" http://localhost:8080/api/v1/admin/chat/sessions
```

**Solution 2: Network Configuration**
```bash
# Check if port is accessible
telnet localhost 8080

# Verify CORS settings in backend
grep -r "CORS" backend/internal/server/

# Update nginx WebSocket configuration
location /api/v1/admin/chat/ws {
    proxy_pass http://backend:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_read_timeout 86400;
}
```

**Solution 3: Backend Service Issues**
```bash
# Restart backend service
docker-compose restart backend

# Check backend resource usage
docker stats backend

# Increase backend resources if needed
# Edit docker-compose.yml:
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
```

### Database Connection Issues

#### Symptoms
- "Database connection failed" errors
- Slow query responses
- Connection pool exhaustion
- Service startup failures

#### Diagnostic Steps
```bash
# 1. Check database status
docker-compose ps postgres

# 2. Test direct connection
psql -h localhost -p 5432 -U chat_user -d chat_dev

# 3. Check connection pool status
docker-compose exec backend ./backend -check-db

# 4. Monitor active connections
psql -h localhost -U chat_user -d chat_dev -c "
SELECT count(*) as active_connections 
FROM pg_stat_activity 
WHERE state = 'active';"
```

#### Solutions

**Solution 1: Connection Pool Configuration**
```bash
# Update environment variables
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Restart backend
docker-compose restart backend
```

**Solution 2: Database Performance**
```sql
-- Check for long-running queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query 
FROM pg_stat_activity 
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';

-- Kill long-running queries if needed
SELECT pg_terminate_backend(pid) FROM pg_stat_activity 
WHERE (now() - pg_stat_activity.query_start) > interval '10 minutes';

-- Check database locks
SELECT * FROM pg_locks WHERE NOT granted;
```

**Solution 3: Database Resource Issues**
```bash
# Check database disk space
df -h

# Check database memory usage
docker stats postgres

# Increase database resources
# Edit docker-compose.yml:
services:
  postgres:
    command: postgres -c shared_buffers=256MB -c max_connections=100
    deploy:
      resources:
        limits:
          memory: 1G
```

### Redis Connection Issues

#### Symptoms
- Session data not persisting
- Cache misses
- "Redis connection failed" errors
- Slow response times

#### Diagnostic Steps
```bash
# 1. Check Redis status
docker-compose ps redis

# 2. Test Redis connection
redis-cli ping

# 3. Check Redis memory usage
redis-cli info memory

# 4. Monitor Redis operations
redis-cli monitor
```

#### Solutions

**Solution 1: Redis Configuration**
```bash
# Check Redis configuration
redis-cli config get "*"

# Increase memory limit if needed
redis-cli config set maxmemory 512mb
redis-cli config set maxmemory-policy allkeys-lru

# Restart Redis
docker-compose restart redis
```

**Solution 2: Connection Pool Issues**
```bash
# Update Redis connection settings
REDIS_MAX_IDLE=10
REDIS_MAX_ACTIVE=100
REDIS_IDLE_TIMEOUT=240s

# Restart backend
docker-compose restart backend
```

## Performance Issues

### Slow AI Response Times

#### Symptoms
- AI responses taking >10 seconds
- Timeout errors
- High CPU usage on backend
- Users reporting slow chat experience

#### Diagnostic Steps
```bash
# 1. Check AI service metrics
curl http://localhost:8080/api/v1/admin/chat/metrics/ai

# 2. Monitor backend CPU/memory
docker stats backend

# 3. Check AWS Bedrock service status
aws bedrock get-model --model-identifier anthropic.claude-3-sonnet-20240229-v1:0

# 4. Analyze response time logs
docker-compose logs backend | grep "ai_response_time"
```

#### Solutions

**Solution 1: Optimize AI Requests**
```bash
# Enable response caching
ENABLE_AI_CACHE=true
AI_CACHE_TTL=3600

# Implement request batching
AI_BATCH_SIZE=5
AI_BATCH_TIMEOUT=2s

# Restart backend
docker-compose restart backend
```

**Solution 2: Scale Backend Services**
```bash
# Scale backend horizontally
docker-compose up -d --scale backend=3

# Or increase backend resources
# Edit docker-compose.yml:
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 4G
```

**Solution 3: AWS Bedrock Optimization**
```bash
# Check AWS region latency
aws bedrock list-foundation-models --region us-east-1
aws bedrock list-foundation-models --region us-west-2

# Switch to closer region if needed
AWS_REGION=us-west-2

# Implement retry logic with exponential backoff
AI_RETRY_ATTEMPTS=3
AI_RETRY_DELAY=1s
```

### High Memory Usage

#### Symptoms
- Out of memory errors
- Container restarts
- Slow performance
- System freezing

#### Diagnostic Steps
```bash
# 1. Check memory usage by service
docker stats

# 2. Check system memory
free -h

# 3. Analyze memory leaks
docker-compose exec backend pprof -http=:6060 http://localhost:6060/debug/pprof/heap

# 4. Check for memory-intensive queries
psql -h localhost -U chat_user -d chat_dev -c "
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 10;"
```

#### Solutions

**Solution 1: Optimize Memory Usage**
```bash
# Enable garbage collection tuning
GOGC=100
GOMEMLIMIT=1GiB

# Implement connection pooling
DB_MAX_OPEN_CONNS=10
REDIS_MAX_IDLE=5

# Restart services
docker-compose restart
```

**Solution 2: Increase Resource Limits**
```yaml
# docker-compose.yml
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 2G
        reservations:
          memory: 1G
  
  postgres:
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
```

### Database Performance Issues

#### Symptoms
- Slow query responses
- High database CPU usage
- Connection timeouts
- Query queue buildup

#### Diagnostic Steps
```sql
-- Check slow queries
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements 
WHERE mean_time > 1000 
ORDER BY mean_time DESC;

-- Check database locks
SELECT * FROM pg_locks l 
JOIN pg_stat_activity a ON l.pid = a.pid 
WHERE NOT l.granted;

-- Check index usage
SELECT schemaname, tablename, attname, n_distinct, correlation 
FROM pg_stats 
WHERE tablename IN ('chat_sessions', 'chat_messages');
```

#### Solutions

**Solution 1: Database Optimization**
```sql
-- Add missing indexes
CREATE INDEX CONCURRENTLY idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX CONCURRENTLY idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX CONCURRENTLY idx_chat_sessions_user_id ON chat_sessions(user_id);

-- Update table statistics
ANALYZE chat_sessions;
ANALYZE chat_messages;

-- Vacuum tables
VACUUM ANALYZE chat_sessions;
VACUUM ANALYZE chat_messages;
```

**Solution 2: Query Optimization**
```sql
-- Optimize pagination queries
-- Instead of OFFSET/LIMIT, use cursor-based pagination
SELECT * FROM chat_messages 
WHERE session_id = $1 AND created_at < $2 
ORDER BY created_at DESC 
LIMIT 50;
```

## Authentication and Authorization Issues

### JWT Token Issues

#### Symptoms
- "Token expired" errors
- "Invalid token" messages
- Frequent re-authentication required
- Authorization failures

#### Diagnostic Steps
```bash
# 1. Decode JWT token
echo "your_token" | cut -d. -f2 | base64 -d | jq

# 2. Check token expiration
date -d @$(echo "your_token" | cut -d. -f2 | base64 -d | jq -r .exp)

# 3. Verify JWT secret
grep JWT_SECRET .env

# 4. Test token validation
curl -H "Authorization: Bearer your_token" http://localhost:8080/api/v1/admin/chat/sessions
```

#### Solutions

**Solution 1: Token Configuration**
```bash
# Increase token expiry time
JWT_EXPIRY=24h

# Ensure JWT secret is secure and consistent
JWT_SECRET=$(openssl rand -base64 32)

# Restart backend
docker-compose restart backend
```

**Solution 2: Token Refresh Implementation**
```javascript
// Frontend token refresh logic
const refreshToken = async () => {
  try {
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('refreshToken')}`
      }
    });
    const data = await response.json();
    localStorage.setItem('accessToken', data.accessToken);
    return data.accessToken;
  } catch (error) {
    // Redirect to login
    window.location.href = '/login';
  }
};
```

### Permission Issues

#### Symptoms
- "Access denied" errors
- Users can't access chat features
- Admin functions not available
- Role-based restrictions not working

#### Diagnostic Steps
```bash
# 1. Check user roles in database
psql -h localhost -U chat_user -d chat_dev -c "
SELECT id, email, role, is_active FROM users WHERE email = 'user@example.com';"

# 2. Verify middleware configuration
grep -r "AuthMiddleware" backend/internal/server/

# 3. Check role assignments
curl -H "Authorization: Bearer your_token" http://localhost:8080/api/v1/admin/users/me
```

#### Solutions

**Solution 1: Update User Roles**
```sql
-- Grant admin role to user
UPDATE users SET role = 'admin' WHERE email = 'user@example.com';

-- Activate user account
UPDATE users SET is_active = true WHERE email = 'user@example.com';
```

**Solution 2: Fix Middleware Configuration**
```go
// Ensure proper middleware order
r.Use(authMiddleware.AuthMiddleware())
r.Use(authMiddleware.AdminMiddleware())
```

## Error Handling and Logging

### Application Errors

#### Symptoms
- 500 Internal Server Error
- Panic errors in logs
- Service crashes
- Unexpected behavior

#### Diagnostic Steps
```bash
# 1. Check error logs
docker-compose logs backend | grep -E "(ERROR|FATAL|panic)"

# 2. Check stack traces
docker-compose logs backend | grep -A 10 "panic"

# 3. Monitor error rates
curl http://localhost:8080/api/v1/admin/chat/metrics/errors

# 4. Check application health
curl http://localhost:8080/health
```

#### Solutions

**Solution 1: Error Recovery**
```bash
# Restart failed services
docker-compose restart backend

# Check for resource constraints
docker stats

# Review recent changes
git log --oneline -10
```

**Solution 2: Improve Error Handling**
```go
// Add proper error handling
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
    defer func() {
        if r := recover(); r != nil {
            h.logger.WithField("panic", r).Error("WebSocket handler panic")
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
    }()
    // ... handler logic
}
```

### Logging Issues

#### Symptoms
- Missing log entries
- Log files growing too large
- Insufficient log detail
- Log parsing errors

#### Diagnostic Steps
```bash
# 1. Check log configuration
grep -r "LOG_LEVEL" .env

# 2. Check log file sizes
du -sh /var/log/chat-system/

# 3. Test log output
docker-compose logs --tail=100 backend

# 4. Check log rotation
ls -la /var/log/chat-system/
```

#### Solutions

**Solution 1: Configure Log Levels**
```bash
# Set appropriate log level
LOG_LEVEL=info  # debug, info, warn, error

# Enable structured logging
LOG_FORMAT=json

# Configure log rotation
LOG_MAX_SIZE=100MB
LOG_MAX_BACKUPS=5
LOG_MAX_AGE=30
```

**Solution 2: Centralized Logging**
```yaml
# docker-compose.yml
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Monitoring and Alerting Issues

### Metrics Collection Problems

#### Symptoms
- Missing metrics in Grafana
- Prometheus scraping failures
- Incomplete monitoring data
- Alert notifications not working

#### Diagnostic Steps
```bash
# 1. Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# 2. Test metrics endpoint
curl http://localhost:8080/metrics

# 3. Check Grafana data sources
curl http://admin:admin@localhost:3000/api/datasources

# 4. Verify alert manager
curl http://localhost:9093/api/v1/alerts
```

#### Solutions

**Solution 1: Fix Prometheus Configuration**
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'chat-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: /metrics
    scrape_interval: 30s
```

**Solution 2: Enable Metrics in Backend**
```bash
# Ensure metrics are enabled
PROMETHEUS_ENABLED=true
METRICS_PORT=8080

# Restart services
docker-compose restart backend prometheus
```

## Backup and Recovery Issues

### Backup Failures

#### Symptoms
- Backup scripts failing
- Incomplete backups
- Storage space issues
- Backup corruption

#### Diagnostic Steps
```bash
# 1. Check backup script logs
tail -f /var/log/backup.log

# 2. Test backup manually
./scripts/backup-chat-system.sh

# 3. Check storage space
df -h /backups/

# 4. Verify backup integrity
pg_restore --list /backups/chat_backup_latest.sql.gz
```

#### Solutions

**Solution 1: Fix Backup Script**
```bash
#!/bin/bash
# Improved backup script with error handling
set -e

BACKUP_DIR="/backups/chat-system"
DATE=$(date +%Y%m%d_%H%M%S)

# Check available space
AVAILABLE=$(df /backups | tail -1 | awk '{print $4}')
if [ $AVAILABLE -lt 1000000 ]; then
    echo "Insufficient disk space for backup"
    exit 1
fi

# Perform backup with error checking
if pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > $BACKUP_DIR/chat_backup_$DATE.sql; then
    gzip $BACKUP_DIR/chat_backup_$DATE.sql
    echo "Backup completed successfully"
else
    echo "Backup failed"
    exit 1
fi
```

### Recovery Issues

#### Symptoms
- Cannot restore from backup
- Data corruption during recovery
- Service won't start after recovery
- Partial data recovery

#### Diagnostic Steps
```bash
# 1. Verify backup file
file /backups/chat_backup_latest.sql.gz
gunzip -t /backups/chat_backup_latest.sql.gz

# 2. Check database state
psql -h localhost -U chat_user -d chat_dev -c "\dt"

# 3. Test recovery in staging
./scripts/restore-chat-system.sh /backups/chat_backup_latest.sql.gz
```

#### Solutions

**Solution 1: Safe Recovery Process**
```bash
#!/bin/bash
# Safe recovery script
set -e

BACKUP_FILE=$1
if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# Stop application
docker-compose stop backend

# Create database backup before recovery
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > /tmp/pre_recovery_backup.sql

# Restore from backup
gunzip -c $BACKUP_FILE | psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Start application
docker-compose start backend

# Verify recovery
curl http://localhost:8080/health
```

## Emergency Procedures

### System Down Recovery

#### Complete System Failure
```bash
# 1. Check system resources
free -h
df -h
ps aux | head -20

# 2. Restart all services
docker-compose down
docker-compose up -d

# 3. Verify services
docker-compose ps
curl http://localhost:8080/health

# 4. Check logs for errors
docker-compose logs --tail=50
```

#### Database Corruption
```bash
# 1. Stop all services
docker-compose stop

# 2. Check database integrity
docker-compose run --rm postgres pg_dump -h postgres -U chat_user -d chat_dev --schema-only > /tmp/schema_check.sql

# 3. Restore from latest backup
./scripts/restore-chat-system.sh /backups/chat_backup_latest.sql.gz

# 4. Restart services
docker-compose up -d
```

### Data Loss Recovery

#### Session Data Recovery
```sql
-- Check for recoverable sessions
SELECT id, user_id, created_at, status 
FROM chat_sessions 
WHERE status != 'deleted' 
ORDER BY created_at DESC;

-- Recover soft-deleted sessions
UPDATE chat_sessions 
SET status = 'active', updated_at = NOW() 
WHERE status = 'deleted' AND created_at > NOW() - INTERVAL '24 hours';
```

#### Message Recovery
```sql
-- Check message integrity
SELECT session_id, COUNT(*) as message_count 
FROM chat_messages 
GROUP BY session_id 
HAVING COUNT(*) = 0;

-- Recover from backup if needed
-- (Use backup restoration procedures above)
```

## Getting Help

### Support Escalation

#### Level 1: Self-Service
1. Check this troubleshooting guide
2. Review system logs
3. Check monitoring dashboards
4. Try basic restart procedures

#### Level 2: Technical Support
1. Gather diagnostic information
2. Create support ticket with:
   - Error messages
   - Log excerpts
   - System metrics
   - Steps to reproduce

#### Level 3: Emergency Support
1. Critical system down
2. Data loss or corruption
3. Security incidents
4. Contact emergency hotline

### Diagnostic Information Collection

#### System Information
```bash
# Create diagnostic bundle
mkdir -p /tmp/chat-system-diagnostics

# System info
uname -a > /tmp/chat-system-diagnostics/system_info.txt
free -h > /tmp/chat-system-diagnostics/memory_info.txt
df -h > /tmp/chat-system-diagnostics/disk_info.txt

# Docker info
docker-compose ps > /tmp/chat-system-diagnostics/docker_status.txt
docker stats --no-stream > /tmp/chat-system-diagnostics/docker_stats.txt

# Application logs
docker-compose logs --tail=1000 > /tmp/chat-system-diagnostics/application_logs.txt

# Configuration
cp .env /tmp/chat-system-diagnostics/environment.txt
cp docker-compose.yml /tmp/chat-system-diagnostics/

# Create archive
tar -czf chat-system-diagnostics-$(date +%Y%m%d_%H%M%S).tar.gz -C /tmp chat-system-diagnostics/
```

### Contact Information

- **Technical Support**: support@yourcompany.com
- **Emergency Hotline**: +1-555-EMERGENCY
- **Documentation**: https://docs.yourcompany.com/chat-system
- **Status Page**: https://status.yourcompany.com