# AWS Minimal Deployment Design

## Overview

This document provides a cost-optimized AWS deployment design that maximizes the use of $400 AWS credits while maintaining production-ready functionality. The architecture prioritizes simplicity, free tier usage, and existing service integration.

## Architecture

### High-Level Architecture Diagram

```
                    Internet
                       │
                       ▼
              ┌─────────────────┐
              │   Route 53 DNS  │
              │ ($0.50/month)   │
              └─────────┬───────┘
                        │
              ┌─────────▼───────┐
              │  ACM SSL Cert   │
              │     (FREE)      │
              └─────────┬───────┘
                        │
              ┌─────────▼───────┐
              │  App Runner     │
              │ (Pay per use)   │
              │ OR              │
              │ Elastic Beanstalk│
              │ (EC2 costs only)│
              └─────────┬───────┘
                        │
    ┌───────────────────▼────────────────────┐
    │         Application Container          │
    │  ┌─────────────────────────────────┐   │
    │  │   Frontend (React) +            │   │
    │  │   Backend (Go)                  │   │
    │  │   Single Container              │   │
    │  └─────────────────────────────────┘   │
    └────────────────┬───────────────────────┘
                     │
           ┌─────────▼─────────┐
           │  RDS PostgreSQL   │
           │   db.t3.micro     │
           │   (FREE TIER)     │
           └───────────────────┘

External Services (Already Configured):
├── AWS Bedrock (Pay per use)
├── AWS SES (Pay per email)
└── CloudWatch (Free tier)
```

### Cost Optimization Strategy

#### Free Tier Maximization
```yaml
# AWS Free Tier Usage (12 months)
EC2:
  - t3.micro: 750 hours/month (FREE)
  - EBS: 30GB General Purpose SSD (FREE)
  
RDS:
  - db.t3.micro: 750 hours/month (FREE)
  - Storage: 20GB (FREE)
  - Backups: 20GB (FREE)

CloudWatch:
  - 10 custom metrics (FREE)
  - 5GB log ingestion (FREE)
  - 1 million API requests (FREE)

Data Transfer:
  - 100GB outbound (FREE)
  - All inbound (FREE)

Route 53:
  - 1 hosted zone: $0.50/month
  - 1 billion queries/month (FREE)

Certificate Manager:
  - SSL certificates (FREE)
```

#### Credit Usage Projection
```yaml
# Monthly costs with $400 credit
Month 1-12: $0-5/month (Free tier covers most usage)
Month 13+: $15-25/month (After free tier expires)

# Credit exhaustion timeline
Conservative estimate: 24+ months
Aggressive usage: 12-18 months
```

## Components and Interfaces

### 1. Application Deployment Options

#### Option A: AWS App Runner (Recommended)
```yaml
# Serverless container service
AppRunner:
  Source: GitHub repository
  Runtime: Docker
  CPU: 0.25 vCPU
  Memory: 0.5 GB
  Scaling: 1-10 instances
  Cost: $0.007/hour when running + $0.000008/request
  
Benefits:
  - Automatic scaling to zero
  - Built-in load balancing
  - Automatic deployments
  - Minimal configuration
```

#### Option B: Elastic Beanstalk
```yaml
# Platform-as-a-Service
ElasticBeanstalk:
  Platform: Docker
  Instance: t3.micro (free tier)
  Load Balancer: Application Load Balancer
  Auto Scaling: 1-3 instances
  Cost: Only EC2 instance costs
  
Benefits:
  - Easy deployment
  - Automatic scaling
  - Health monitoring
  - Rolling deployments
```

### 2. Database Configuration

#### RDS PostgreSQL Setup
```yaml
Engine: postgres
EngineVersion: "15.4"
DBInstanceClass: db.t3.micro  # Free tier eligible
AllocatedStorage: 20          # Free tier limit
StorageType: gp2
StorageEncrypted: true
MultiAZ: false               # Keep costs low
BackupRetentionPeriod: 7
PreferredBackupWindow: "03:00-04:00"
PreferredMaintenanceWindow: "sun:04:00-sun:05:00"
DeletionProtection: true
```

#### Database Connection Configuration
```go
// Optimized for t3.micro
type DatabaseConfig struct {
    URL                string
    MaxOpenConnections int `default:"5"`   // Low for micro instance
    MaxIdleConnections int `default:"2"`   
    ConnMaxLifetime    time.Duration `default:"30m"`
}
```

### 3. Container Configuration

#### Single Container Approach
```dockerfile
# Multi-stage build for efficiency
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM nginx:alpine
# Copy frontend build
COPY --from=frontend-builder /app/frontend/build /usr/share/nginx/html
# Copy backend binary
COPY --from=backend-builder /app/server /usr/local/bin/
# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf
# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 80
CMD ["/start.sh"]
```

#### Startup Script
```bash
#!/bin/sh
# start.sh - Run both frontend and backend in single container

# Start backend in background
/usr/local/bin/server &

# Start nginx in foreground
nginx -g 'daemon off;'
```

#### Nginx Configuration
```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server localhost:8061;
    }
    
    server {
        listen 80;
        
        # Frontend
        location / {
            root /usr/share/nginx/html;
            try_files $uri $uri/ /index.html;
        }
        
        # Backend API
        location /api/ {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
        
        # Health check
        location /health {
            proxy_pass http://backend/health;
        }
    }
}
```

### 4. Environment Configuration

#### App Runner Configuration
```yaml
# apprunner.yaml
version: 1.0
runtime: docker
build:
  commands:
    build:
      - echo "Building application..."
run:
  runtime-version: latest
  command: /start.sh
  network:
    port: 80
    env: PORT
  env:
    - name: PORT
      value: "80"
    - name: DATABASE_URL
      value: "postgresql://user:pass@rds-endpoint:5432/dbname"
    - name: GIN_MODE
      value: "release"
```

#### Environment Variables
```bash
# Production environment variables
PORT=80
GIN_MODE=release
LOG_LEVEL=info

# Database (RDS connection)
DATABASE_URL=postgresql://postgres:password@rds-endpoint.region.rds.amazonaws.com:5432/consulting

# Your existing AWS configuration
AWS_ACCESS_KEY_ID=your_existing_key
AWS_SECRET_ACCESS_KEY=your_existing_secret
AWS_BEARER_TOKEN_BEDROCK=your_existing_token
BEDROCK_REGION=us-east-1
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0
AWS_SES_REGION=us-east-1
SES_SENDER_EMAIL=info@cloudpartner.pro

# Application configuration
JWT_SECRET=your_production_jwt_secret
CORS_ALLOWED_ORIGINS=https://your-domain.com
CHAT_MODE=polling
CHAT_POLLING_INTERVAL=3000
ENABLE_EMAIL_EVENTS=true
```

## Data Models

### Cost Tracking Configuration

#### CloudWatch Billing Alarms
```yaml
# Billing alerts for credit monitoring
BillingAlerts:
  - Name: "75PercentCreditUsage"
    Threshold: 300  # $300 of $400 credits
    ComparisonOperator: GreaterThanThreshold
    
  - Name: "90PercentCreditUsage"
    Threshold: 360  # $360 of $400 credits
    ComparisonOperator: GreaterThanThreshold
    
  - Name: "DailyCostAlert"
    Threshold: 5    # $5/day
    ComparisonOperator: GreaterThanThreshold
```

#### Cost Optimization Rules
```go
// Cost monitoring service
type CostMonitor struct {
    CloudWatchClient *cloudwatch.CloudWatch
    BillingClient    *billing.Billing
}

func (c *CostMonitor) GetDailyCosts() (*CostReport, error) {
    // Get daily cost breakdown
    // Alert if approaching credit limits
    // Suggest optimizations
}
```

### Resource Tagging Strategy
```yaml
# Consistent tagging for cost tracking
Tags:
  Project: "cloud-consulting"
  Environment: "production"
  CostCenter: "startup"
  Owner: "founder"
  BillingType: "credits"
```

## Error Handling

### Cost Overrun Protection
```go
// Implement cost protection
func (s *Service) checkCostLimits() error {
    dailyCost, err := s.costMonitor.GetDailyCosts()
    if err != nil {
        return err
    }
    
    if dailyCost > s.config.MaxDailyCost {
        // Alert and potentially scale down
        s.alerting.SendCostAlert(dailyCost)
        return s.scaleDown()
    }
    
    return nil
}
```

### Resource Optimization
```go
// Optimize database connections for t3.micro
func optimizeForMicroInstance() *sql.DB {
    db.SetMaxOpenConns(5)    // Low for micro instance
    db.SetMaxIdleConns(2)    
    db.SetConnMaxLifetime(30 * time.Minute)
    
    return db
}
```

## Testing Strategy

### Cost Testing
```bash
# Test cost optimization
aws ce get-cost-and-usage \
    --time-period Start=2024-01-01,End=2024-01-31 \
    --granularity DAILY \
    --metrics BlendedCost

# Monitor free tier usage
aws support describe-trusted-advisor-checks \
    --language en \
    --query 'checks[?name==`Service Limits`]'
```

### Load Testing for Micro Instances
```yaml
# K6 test optimized for small instances
export let options = {
  stages: [
    { duration: '2m', target: 5 },   # Low load for micro
    { duration: '5m', target: 10 },  # Sustained low load
    { duration: '2m', target: 0 },
  ],
};
```

## Security Considerations

### IAM Roles and Policies
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "bedrock:InvokeModel",
        "ses:SendEmail",
        "ses:SendRawEmail",
        "rds:DescribeDBInstances",
        "cloudwatch:PutMetricData"
      ],
      "Resource": "*"
    }
  ]
}
```

### Security Groups
```yaml
# App Runner Security Group
AppRunnerSG:
  Ingress:
    - Port: 80
      Source: 0.0.0.0/0
    - Port: 443  
      Source: 0.0.0.0/0
  Egress:
    - All traffic to 0.0.0.0/0

# RDS Security Group  
RDSSG:
  Ingress:
    - Port: 5432
      Source: AppRunnerSG
  Egress: None
```

## Monitoring and Observability

### CloudWatch Configuration
```yaml
# Free tier monitoring
CloudWatch:
  Metrics:
    - ApplicationRequests
    - DatabaseConnections
    - ErrorRate
    - ResponseTime
  
  Logs:
    - ApplicationLogs
    - DatabaseLogs
    - AccessLogs
  
  Alarms:
    - HighErrorRate
    - DatabaseConnectionFailures
    - CostThresholds
```

### Application Metrics
```go
// Custom metrics within free tier limits
func (m *MetricsCollector) RecordMetrics() {
    // Use only essential metrics to stay within free tier
    m.cloudwatch.PutMetricData(&cloudwatch.PutMetricDataInput{
        Namespace: aws.String("CloudConsulting/Application"),
        MetricData: []*cloudwatch.MetricDatum{
            {
                MetricName: aws.String("RequestCount"),
                Value:      aws.Float64(float64(m.requestCount)),
                Unit:       aws.String("Count"),
            },
        },
    })
}
```

## Cost Optimization Strategies

### Automatic Scaling Policies
```yaml
# App Runner auto-scaling
AutoScaling:
  MinInstances: 1
  MaxInstances: 3
  TargetCPU: 70
  TargetMemory: 80
  ScaleDownDelay: 300s  # Quick scale down to save costs
```

### Resource Scheduling
```go
// Schedule non-critical tasks during low-cost periods
func (s *Scheduler) ScheduleBackups() {
    // Run backups during free tier hours
    // Batch operations to minimize costs
    // Use spot instances for batch processing if needed
}
```

### Credit Usage Optimization
```yaml
# Strategies to maximize credit value
Optimization:
  - Use free tier services first
  - Batch operations to minimize API calls
  - Implement caching to reduce database load
  - Use CloudFront for static content delivery
  - Optimize images and assets for faster loading
  - Implement connection pooling
  - Use compression for data transfer
```

This design maximizes your $400 AWS credits while providing a production-ready deployment that can easily scale as your business grows.