# AI Consultant Live Chat Deployment Guide

## Overview

This guide covers the deployment and configuration of the AI Consultant Live Chat system across different environments (development, staging, production).

## Architecture Overview

The chat system consists of:
- **Backend API**: Go-based REST API with polling endpoints
- **Frontend**: React-based admin interface
- **Database**: PostgreSQL for persistent storage
- **Cache**: Redis for session management and caching
- **AI Service**: AWS Bedrock integration
- **Monitoring**: Prometheus, Grafana, and custom alerting
- **Load Balancer**: Nginx for HTTP traffic

## Prerequisites

### System Requirements

#### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 20GB SSD
- **Network**: 100Mbps connection

#### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 50GB+ SSD
- **Network**: 1Gbps connection

### Software Dependencies

#### Required Software
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Kubernetes**: 1.21+ (for K8s deployment)
- **Helm**: 3.0+ (for K8s deployment)

#### Optional Tools
- **kubectl**: For Kubernetes management
- **aws-cli**: For AWS resource management
- **terraform**: For infrastructure as code

## Environment Configuration

### Environment Variables

Create environment files for each deployment environment:

#### Development (.env.dev)
```bash
# Application
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=chat_dev
DB_USER=chat_user
DB_PASSWORD=dev_password
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# AWS Bedrock
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRY=24h

# Chat Configuration
CHAT_POLLING_INTERVAL=3000
CHAT_MAX_RETRIES=3
WS_MAX_MESSAGE_SIZE=4096

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
RATE_LIMIT_WS_MESSAGES=60

# Monitoring
PROMETHEUS_ENABLED=true
METRICS_PORT=9090
LOG_LEVEL=debug
```

#### Staging (.env.staging)
```bash
# Application
APP_ENV=staging
APP_PORT=8080
APP_DEBUG=false

# Database
DB_HOST=staging-db.internal
DB_PORT=5432
DB_NAME=chat_staging
DB_USER=chat_user
DB_PASSWORD=${DB_PASSWORD_SECRET}
DB_SSL_MODE=require

# Redis
REDIS_HOST=staging-redis.internal
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD_SECRET}
REDIS_DB=0

# AWS Bedrock
AWS_REGION=us-east-1
BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0

# JWT
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRY=8h

# Chat Configuration
CHAT_POLLING_INTERVAL=2000
CHAT_MAX_RETRIES=5
WS_MAX_MESSAGE_SIZE=8192

# Rate Limiting
RATE_LIMIT_REQUESTS=200
RATE_LIMIT_WINDOW=60s
RATE_LIMIT_WS_MESSAGES=120

# Monitoring
PROMETHEUS_ENABLED=true
METRICS_PORT=9090
LOG_LEVEL=info
```

#### Production (.env.prod)
```bash
# Application
APP_ENV=production
APP_PORT=8080
APP_DEBUG=false

# Database
DB_HOST=${RDS_ENDPOINT}
DB_PORT=5432
DB_NAME=chat_production
DB_USER=chat_user
DB_PASSWORD=${DB_PASSWORD_SECRET}
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Redis
REDIS_HOST=${ELASTICACHE_ENDPOINT}
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD_SECRET}
REDIS_DB=0

# AWS Bedrock
AWS_REGION=us-east-1
BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0

# JWT
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRY=4h

# Chat Configuration
CHAT_POLLING_INTERVAL=1000
CHAT_MAX_RETRIES=10
WS_MAX_MESSAGE_SIZE=16384

# Rate Limiting
RATE_LIMIT_REQUESTS=500
RATE_LIMIT_WINDOW=60s
RATE_LIMIT_WS_MESSAGES=300

# Monitoring
PROMETHEUS_ENABLED=true
METRICS_PORT=9090
LOG_LEVEL=warn

# Security
CORS_ALLOWED_ORIGINS=https://yourdomain.com
TLS_CERT_PATH=/etc/ssl/certs/app.crt
TLS_KEY_PATH=/etc/ssl/private/app.key
```

## Deployment Methods

### 1. Docker Compose Deployment

#### Development Setup
```bash
# Clone repository
git clone <repository-url>
cd ai-consultant-chat

# Copy environment file
cp .env.example .env.dev

# Start services
docker-compose -f docker-compose.chat.yml --env-file .env.dev up -d

# Run database migrations
docker-compose exec backend ./migrate up

# Verify deployment
curl http://localhost:8080/health
```

#### Production Setup
```bash
# Copy production environment
cp .env.example .env.prod
# Edit .env.prod with production values

# Start production services
docker-compose -f docker-compose.chat.yml --env-file .env.prod up -d

# Run migrations
docker-compose exec backend ./migrate up

# Verify deployment
curl https://yourdomain.com/health
```

### 2. Kubernetes Deployment

#### Prerequisites
```bash
# Install Helm
curl https://get.helm.sh/helm-v3.10.0-linux-amd64.tar.gz | tar xz
sudo mv linux-amd64/helm /usr/local/bin/

# Add required Helm repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

#### Deploy Dependencies
```bash
# Deploy PostgreSQL
helm install chat-postgres bitnami/postgresql \
  --set auth.postgresPassword=your_password \
  --set auth.database=chat_production

# Deploy Redis
helm install chat-redis bitnami/redis \
  --set auth.password=your_redis_password
```

#### Deploy Chat System
```bash
# Create namespace
kubectl create namespace chat-system

# Create secrets
kubectl create secret generic chat-secrets \
  --from-literal=db-password=your_db_password \
  --from-literal=redis-password=your_redis_password \
  --from-literal=jwt-secret=your_jwt_secret \
  --namespace=chat-system

# Deploy using Helm
helm install chat-system ./k8s/chat-system \
  --namespace=chat-system \
  --values=./k8s/chat-system/values-production.yaml

# Verify deployment
kubectl get pods -n chat-system
kubectl get services -n chat-system
```

#### Kubernetes Configuration Files

**values-production.yaml**
```yaml
# Application configuration
app:
  name: chat-system
  version: "1.0.0"
  environment: production

# Backend configuration
backend:
  replicaCount: 3
  image:
    repository: your-registry/chat-backend
    tag: "latest"
    pullPolicy: Always
  
  resources:
    requests:
      cpu: 500m
      memory: 1Gi
    limits:
      cpu: 1000m
      memory: 2Gi
  
  env:
    APP_ENV: production
    LOG_LEVEL: info
    PROMETHEUS_ENABLED: "true"

# Frontend configuration
frontend:
  replicaCount: 2
  image:
    repository: your-registry/chat-frontend
    tag: "latest"
    pullPolicy: Always
  
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi

# Service configuration
service:
  type: LoadBalancer
  port: 80
  targetPort: 8080

# Ingress configuration
ingress:
  enabled: true
  className: nginx
  annotations:

    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
  hosts:
    - host: chat.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: chat-tls
      hosts:
        - chat.yourdomain.com

# Horizontal Pod Autoscaler
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80

# Monitoring
monitoring:
  enabled: true
  serviceMonitor:
    enabled: true
    interval: 30s
```

### 3. AWS ECS Deployment

#### Task Definition
```json
{
  "family": "chat-system",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "1024",
  "memory": "2048",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "taskRoleArn": "arn:aws:iam::account:role/chatSystemTaskRole",
  "containerDefinitions": [
    {
      "name": "chat-backend",
      "image": "your-account.dkr.ecr.region.amazonaws.com/chat-backend:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "APP_ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:chat/db-password"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/chat-system",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

#### Service Configuration
```json
{
  "serviceName": "chat-system",
  "cluster": "production-cluster",
  "taskDefinition": "chat-system:1",
  "desiredCount": 3,
  "launchType": "FARGATE",
  "networkConfiguration": {
    "awsvpcConfiguration": {
      "subnets": [
        "subnet-12345678",
        "subnet-87654321"
      ],
      "securityGroups": [
        "sg-12345678"
      ],
      "assignPublicIp": "DISABLED"
    }
  },
  "loadBalancers": [
    {
      "targetGroupArn": "arn:aws:elasticloadbalancing:region:account:targetgroup/chat-system/1234567890123456",
      "containerName": "chat-backend",
      "containerPort": 8080
    }
  ]
}
```

## Database Setup

### PostgreSQL Configuration

#### Development Database
```sql
-- Create database and user
CREATE DATABASE chat_dev;
CREATE USER chat_user WITH PASSWORD 'dev_password';
GRANT ALL PRIVILEGES ON DATABASE chat_dev TO chat_user;

-- Connect to database and run migrations
\c chat_dev;
-- Run migration files from backend/scripts/
```

#### Production Database Setup
```bash
# Using AWS RDS
aws rds create-db-instance \
  --db-instance-identifier chat-production \
  --db-instance-class db.t3.medium \
  --engine postgres \
  --engine-version 13.7 \
  --allocated-storage 100 \
  --storage-type gp2 \
  --storage-encrypted \
  --master-username chat_admin \
  --master-user-password your_secure_password \
  --vpc-security-group-ids sg-12345678 \
  --db-subnet-group-name production-subnet-group \
  --backup-retention-period 7 \
  --multi-az \
  --auto-minor-version-upgrade
```

### Redis Configuration

#### Development Redis
```bash
# Using Docker
docker run -d \
  --name chat-redis \
  -p 6379:6379 \
  redis:7-alpine
```

#### Production Redis Setup
```bash
# Using AWS ElastiCache
aws elasticache create-cache-cluster \
  --cache-cluster-id chat-production \
  --cache-node-type cache.t3.micro \
  --engine redis \
  --engine-version 7.0 \
  --num-cache-nodes 1 \
  --cache-subnet-group-name production-subnet-group \
  --security-group-ids sg-87654321
```

## Load Balancer Configuration

### Nginx Configuration

#### nginx.conf
```nginx
upstream chat_backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    listen 80;
    server_name chat.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name chat.yourdomain.com;

    ssl_certificate /etc/ssl/certs/chat.crt;
    ssl_certificate_key /etc/ssl/private/chat.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

    # Chat polling endpoints
    location /api/v1/admin/chat/polling {
        proxy_pass http://chat_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 30;
    }

    # API endpoints
    location /api/ {
        proxy_pass http://chat_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Frontend
    location / {
        root /var/www/html;
        try_files $uri $uri/ /index.html;
    }
}
```

### AWS Application Load Balancer

#### Target Group Configuration
```bash
# Create target group
aws elbv2 create-target-group \
  --name chat-system-tg \
  --protocol HTTP \
  --port 8080 \
  --vpc-id vpc-12345678 \
  --health-check-path /health \
  --health-check-interval-seconds 30 \
  --health-check-timeout-seconds 5 \
  --healthy-threshold-count 2 \
  --unhealthy-threshold-count 3
```

#### Load Balancer Configuration
```bash
# Create load balancer
aws elbv2 create-load-balancer \
  --name chat-system-alb \
  --subnets subnet-12345678 subnet-87654321 \
  --security-groups sg-12345678 \
  --scheme internet-facing \
  --type application \
  --ip-address-type ipv4

# Create listener
aws elbv2 create-listener \
  --load-balancer-arn arn:aws:elasticloadbalancing:region:account:loadbalancer/app/chat-system-alb/1234567890123456 \
  --protocol HTTPS \
  --port 443 \
  --certificates CertificateArn=arn:aws:acm:region:account:certificate/12345678-1234-1234-1234-123456789012 \
  --default-actions Type=forward,TargetGroupArn=arn:aws:elasticloadbalancing:region:account:targetgroup/chat-system-tg/1234567890123456
```

## Monitoring Setup

### Prometheus Configuration

#### prometheus.yml
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "chat-system-alerts.yml"

scrape_configs:
  - job_name: 'chat-backend'
    static_configs:
      - targets: ['backend:9090']
    metrics_path: /metrics
    scrape_interval: 30s

  - job_name: 'chat-system-health'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: /api/v1/admin/chat/health
    scrape_interval: 60s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

### Grafana Dashboard

Import the dashboard configuration from `monitoring/grafana/dashboards/chat-system-dashboard.json`

Key metrics to monitor:
- HTTP polling requests
- Message throughput
- AI response times
- Error rates
- Database performance
- Cache hit rates

## Security Configuration

### SSL/TLS Setup

#### Certificate Generation
```bash
# Using Let's Encrypt
certbot certonly --nginx -d chat.yourdomain.com

# Using AWS Certificate Manager
aws acm request-certificate \
  --domain-name chat.yourdomain.com \
  --validation-method DNS \
  --subject-alternative-names *.yourdomain.com
```

### Firewall Rules

#### Security Group (AWS)
```bash
# Create security group
aws ec2 create-security-group \
  --group-name chat-system-sg \
  --description "Security group for chat system"

# Allow HTTPS
aws ec2 authorize-security-group-ingress \
  --group-id sg-12345678 \
  --protocol tcp \
  --port 443 \
  --cidr 0.0.0.0/0

# Allow HTTP (redirect to HTTPS)
aws ec2 authorize-security-group-ingress \
  --group-id sg-12345678 \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0
```

## Backup and Recovery

### Database Backup

#### Automated Backup Script
```bash
#!/bin/bash
# backup-chat-database.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/chat-system"
DB_NAME="chat_production"

# Create backup directory
mkdir -p $BACKUP_DIR

# Perform backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > $BACKUP_DIR/chat_backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/chat_backup_$DATE.sql

# Remove backups older than 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: chat_backup_$DATE.sql.gz"
```

#### Cron Job Setup
```bash
# Add to crontab
0 2 * * * /scripts/backup-chat-database.sh
```

### Recovery Procedures

#### Database Recovery
```bash
# Stop application
docker-compose stop backend

# Restore database
gunzip -c /backups/chat-system/chat_backup_20250208_020000.sql.gz | \
  psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Start application
docker-compose start backend
```

## Deployment Checklist

### Pre-Deployment
- [ ] Environment variables configured
- [ ] Database migrations tested
- [ ] SSL certificates installed
- [ ] Monitoring configured
- [ ] Backup procedures tested
- [ ] Security groups configured
- [ ] Load balancer configured

### Deployment
- [ ] Deploy database and Redis
- [ ] Run database migrations
- [ ] Deploy backend services
- [ ] Deploy frontend
- [ ] Configure load balancer
- [ ] Test polling endpoints
- [ ] Verify API endpoints
- [ ] Check monitoring dashboards

### Post-Deployment
- [ ] Verify all services are running
- [ ] Test chat functionality
- [ ] Check logs for errors
- [ ] Verify metrics collection
- [ ] Test backup procedures
- [ ] Update documentation
- [ ] Notify stakeholders

## Troubleshooting

### Common Deployment Issues

#### Database Connection Issues
```bash
# Check database connectivity
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;"

# Check connection pool
docker-compose logs backend | grep "database"
```

#### Chat Polling Issues
```bash
# Test polling endpoint
curl -X POST http://localhost:8080/api/v1/admin/chat/polling

# Check nginx configuration
nginx -t
systemctl reload nginx
```

#### Performance Issues
```bash
# Check resource usage
docker stats
kubectl top pods -n chat-system

# Check database performance
SELECT * FROM pg_stat_activity WHERE state = 'active';
```

For additional troubleshooting, see the [Troubleshooting Guide](troubleshooting-guide.md).