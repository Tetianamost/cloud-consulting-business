# AWS Minimal Deployment Implementation Plan

## Task Overview

This implementation plan leverages your **$400 AWS credits** to deploy your Cloud Consulting Platform with minimal setup and maximum credit efficiency. **Estimated monthly cost: $0-5** (covered by credits for 12+ months).

## **🎯 Perfect for Your $400 Credits:**

✅ **Free for 12+ months** - Credits cover everything  
✅ **15-30 minute setup** - Much faster than complex AWS setups  
✅ **Your existing SES/Bedrock work** - No reconfiguration needed  
✅ **Production-ready** - Auto-scaling, monitoring, backups included  
✅ **Simple architecture** - Easy to manage and debug  

## **Cost Breakdown with Credits:**

| Service | Normal Cost | With Free Tier | Credit Usage |
|---------|-------------|----------------|--------------|
| **App Runner** | $5-15/month | $0-5/month | Minimal |
| **RDS db.t3.micro** | $15/month | **FREE** (12 months) | $0 |
| **Route 53** | $0.50/month | $0.50/month | $6/year |
| **SSL Certificate** | $0/month | **FREE** | $0 |
| **CloudWatch** | $10/month | **FREE** (basic) | $0 |
| **Total Year 1** | **$0-6/month** | **Credits last 24+ months** |

## Implementation Tasks

- [ ] 1. Set up AWS billing alerts and credit monitoring
  - Configure billing alerts at $50, $200, $300 of credit usage
  - Set up daily cost monitoring dashboard
  - Enable detailed billing reports
  - Configure SNS notifications for cost alerts
  - _Requirements: 1.4, 5.1, 5.2_

- [ ] 2. Create RDS PostgreSQL database (Free Tier)
  - Launch db.t3.micro PostgreSQL instance (free tier eligible)
  - Configure 20GB storage (free tier limit)
  - Set up automated backups (7-day retention)
  - Configure security group for database access
  - Enable encryption at rest
  - _Requirements: 2.2, 6.1, 6.2, 6.3_

- [ ] 3. Prepare application for containerized deployment
  - Create single-container Dockerfile combining frontend and backend
  - Configure Nginx to serve React frontend and proxy API calls
  - Create startup script to run both services
  - Test container build and functionality locally
  - _Requirements: 2.1, 4.3_

- [ ] 4. Set up AWS App Runner service
  - Create App Runner service from GitHub repository
  - Configure automatic deployments on git push
  - Set environment variables for production
  - Configure scaling settings (1-3 instances)
  - Test initial deployment and health checks
  - _Requirements: 4.1, 4.2, 4.4_

- [ ] 5. Configure environment variables and secrets
  - Add DATABASE_URL from RDS instance
  - Configure existing AWS credentials (SES, Bedrock)
  - Set production application configuration
  - Configure CORS origins for custom domain
  - Test all AWS service integrations
  - _Requirements: 3.1, 3.2, 3.4_

- [ ] 6. Run database migrations
  - Connect to RDS instance from App Runner
  - Execute init.sql migration for core tables
  - Execute chat_migration.sql for chat system
  - Execute email_events_migration.sql for email tracking
  - Verify all tables and indexes are created
  - _Requirements: 6.4_

- [ ] 7. Set up custom domain and SSL (optional)
  - Configure Route 53 hosted zone for your domain
  - Request SSL certificate through AWS Certificate Manager
  - Configure App Runner custom domain
  - Update DNS records to point to App Runner
  - Test SSL certificate and domain routing
  - _Requirements: 7.1, 7.2, 7.3_

- [ ] 8. Configure monitoring and logging
  - Set up CloudWatch monitoring for App Runner
  - Configure application log aggregation
  - Create basic health check alarms
  - Set up error rate and response time monitoring
  - Test monitoring dashboard and alerts
  - _Requirements: 8.1, 8.2, 8.3_

- [ ] 9. Implement backup and recovery procedures
  - Verify RDS automated backups are working
  - Set up application code backup to S3
  - Create disaster recovery documentation
  - Test database restore procedures
  - Configure backup monitoring and alerts
  - _Requirements: 10.1, 10.3, 10.4_

- [ ] 10. Security hardening and access control
  - Configure IAM roles with least privilege access
  - Set up security groups with minimal required access
  - Enable CloudTrail for audit logging
  - Configure AWS Systems Manager for secrets
  - Review and test security configurations
  - _Requirements: 9.1, 9.2, 9.3_

- [ ] 11. Performance optimization for free tier
  - Optimize database connection pooling for t3.micro
  - Configure application caching strategies
  - Optimize container resource usage
  - Set up CloudFront for static asset delivery (optional)
  - Test performance under expected load
  - _Requirements: 1.1, 5.3_

- [ ] 12. Final testing and validation
  - Perform end-to-end testing of all features
  - Test email delivery and AI report generation
  - Verify chat functionality and admin dashboard
  - Test automatic scaling and deployment
  - Validate cost monitoring and alerts
  - _Requirements: All requirements validation_

## Detailed Implementation Steps

### Task 1: Set Up Billing Alerts

**Configure cost monitoring:**
```bash
# Set up billing alerts using AWS CLI
aws budgets create-budget \
    --account-id $(aws sts get-caller-identity --query Account --output text) \
    --budget '{
        "BudgetName": "AWS-Credits-Monitor",
        "BudgetLimit": {
            "Amount": "400",
            "Unit": "USD"
        },
        "TimeUnit": "MONTHLY",
        "BudgetType": "COST"
    }' \
    --notifications-with-subscribers '[
        {
            "Notification": {
                "NotificationType": "ACTUAL",
                "ComparisonOperator": "GREATER_THAN",
                "Threshold": 75
            },
            "Subscribers": [
                {
                    "SubscriptionType": "EMAIL",
                    "Address": "your-email@example.com"
                }
            ]
        }
    ]'

# Create CloudWatch billing alarm
aws cloudwatch put-metric-alarm \
    --alarm-name "AWS-Credits-75-Percent" \
    --alarm-description "Alert when 75% of AWS credits used" \
    --metric-name EstimatedCharges \
    --namespace AWS/Billing \
    --statistic Maximum \
    --period 86400 \
    --threshold 300 \
    --comparison-operator GreaterThanThreshold \
    --dimensions Name=Currency,Value=USD \
    --evaluation-periods 1
```

### Task 2: Create RDS Database

**Launch free tier RDS instance:**
```bash
# Create RDS subnet group (if using custom VPC)
aws rds create-db-subnet-group \
    --db-subnet-group-name consulting-db-subnet-group \
    --db-subnet-group-description "Subnet group for consulting database" \
    --subnet-ids subnet-12345 subnet-67890

# Create RDS instance (free tier)
aws rds create-db-instance \
    --db-instance-identifier consulting-prod \
    --db-instance-class db.t3.micro \
    --engine postgres \
    --engine-version 15.4 \
    --master-username postgres \
    --master-user-password "$(openssl rand -base64 32)" \
    --allocated-storage 20 \
    --storage-type gp2 \
    --storage-encrypted \
    --backup-retention-period 7 \
    --preferred-backup-window "03:00-04:00" \
    --preferred-maintenance-window "sun:04:00-sun:05:00" \
    --no-multi-az \
    --no-publicly-accessible \
    --vpc-security-group-ids sg-12345 \
    --db-subnet-group-name consulting-db-subnet-group \
    --no-deletion-protection

# Get RDS endpoint
aws rds describe-db-instances \
    --db-instance-identifier consulting-prod \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text
```

### Task 3: Prepare Container

**Create optimized Dockerfile:**
```dockerfile
# Dockerfile
# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Stage 3: Final container
FROM nginx:alpine
# Install supervisor to run multiple processes
RUN apk add --no-cache supervisor

# Copy frontend build
COPY --from=frontend-builder /app/frontend/build /usr/share/nginx/html

# Copy backend binary
COPY --from=backend-builder /app/server /usr/local/bin/server
COPY --from=backend-builder /app/templates /templates

# Copy configuration files
COPY nginx.conf /etc/nginx/nginx.conf
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

EXPOSE 80
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
```

**Create supervisor configuration:**
```ini
# supervisord.conf
[supervisord]
nodaemon=true
user=root

[program:backend]
command=/usr/local/bin/server
autostart=true
autorestart=true
stderr_logfile=/var/log/backend.err.log
stdout_logfile=/var/log/backend.out.log

[program:nginx]
command=nginx -g 'daemon off;'
autostart=true
autorestart=true
stderr_logfile=/var/log/nginx.err.log
stdout_logfile=/var/log/nginx.out.log
```

**Create nginx configuration:**
```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    upstream backend {
        server localhost:8061;
    }
    
    server {
        listen 80;
        
        # Frontend
        location / {
            root /usr/share/nginx/html;
            index index.html;
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

### Task 4: Deploy with App Runner

**Create App Runner service:**
```bash
# Create apprunner.yaml in your repository root
cat > apprunner.yaml << EOF
version: 1.0
runtime: docker
build:
  commands:
    build:
      - echo "Building with Docker..."
run:
  runtime-version: latest
  command: /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
  network:
    port: 80
    env: PORT
  env:
    - name: PORT
      value: "80"
EOF

# Create App Runner service via AWS Console or CLI
aws apprunner create-service \
    --service-name "cloud-consulting-prod" \
    --source-configuration '{
        "ImageRepository": {
            "ImageIdentifier": "your-account.dkr.ecr.region.amazonaws.com/cloud-consulting:latest",
            "ImageConfiguration": {
                "Port": "80"
            },
            "ImageRepositoryType": "ECR"
        },
        "AutoDeploymentsEnabled": true
    }' \
    --instance-configuration '{
        "Cpu": "0.25 vCPU",
        "Memory": "0.5 GB"
    }'
```

### Task 5: Configure Environment Variables

**Set environment variables in App Runner:**
```bash
# Database connection
DATABASE_URL=postgresql://postgres:password@rds-endpoint.region.rds.amazonaws.com:5432/consulting

# Your existing AWS configuration
AWS_ACCESS_KEY_ID=your_existing_access_key
AWS_SECRET_ACCESS_KEY=your_existing_secret_key
AWS_BEARER_TOKEN_BEDROCK=your_existing_bedrock_token
BEDROCK_REGION=us-east-1
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0
AWS_SES_REGION=us-east-1
SES_SENDER_EMAIL=info@cloudpartner.pro
SES_REPLY_TO_EMAIL=info@cloudpartner.pro

# Application configuration
PORT=80
GIN_MODE=release
LOG_LEVEL=info
JWT_SECRET=your_new_production_jwt_secret
CORS_ALLOWED_ORIGINS=https://your-app-runner-url.region.awsapprunner.com
CHAT_MODE=polling
CHAT_POLLING_INTERVAL=3000
ENABLE_EMAIL_EVENTS=true
```

### Task 6: Run Database Migrations

**Execute migrations:**
```bash
# Connect to RDS and run migrations
# Option 1: From local machine (if RDS is publicly accessible)
psql "postgresql://postgres:password@rds-endpoint.region.rds.amazonaws.com:5432/consulting" \
    -f backend/scripts/init.sql

# Option 2: From App Runner container (recommended)
# Add migration execution to your application startup
```

**Add migration to Go application:**
```go
// Add to main.go
func runMigrations(db *sql.DB) error {
    migrationFiles := []string{
        "scripts/init.sql",
        "scripts/chat_migration.sql",
        "scripts/email_events_migration.sql",
    }
    
    for _, file := range migrationFiles {
        content, err := os.ReadFile(file)
        if err != nil {
            log.Printf("Migration file %s not found, skipping", file)
            continue
        }
        
        if _, err := db.Exec(string(content)); err != nil {
            log.Printf("Migration %s failed: %v", file, err)
            return err
        }
        
        log.Printf("Migration %s completed successfully", file)
    }
    
    return nil
}
```

### Task 7: Custom Domain Setup (Optional)

**Configure custom domain:**
```bash
# Create hosted zone
aws route53 create-hosted-zone \
    --name your-domain.com \
    --caller-reference $(date +%s)

# Request SSL certificate
aws acm request-certificate \
    --domain-name your-domain.com \
    --subject-alternative-names "*.your-domain.com" \
    --validation-method DNS \
    --region us-east-1

# Associate custom domain with App Runner (via console)
# Update DNS records to point to App Runner
```

## Validation Steps

### Cost Monitoring Validation
```bash
# Check current AWS costs
aws ce get-cost-and-usage \
    --time-period Start=2024-01-01,End=2024-01-31 \
    --granularity DAILY \
    --metrics BlendedCost

# Check free tier usage
aws support describe-trusted-advisor-checks \
    --language en
```

### Application Testing
```bash
# Test App Runner health
curl https://your-app-runner-url.region.awsapprunner.com/health

# Test API functionality
curl -X POST https://your-app-runner-url.region.awsapprunner.com/api/v1/inquiries \
    -H "Content-Type: application/json" \
    -d '{"name":"Test","email":"test@example.com","services":["assessment"],"message":"Test"}'

# Test admin login
curl -X POST https://your-app-runner-url.region.awsapprunner.com/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"cloudadmin"}'
```

### AWS Service Integration Testing
```bash
# Test SES integration (check App Runner logs)
# Test Bedrock integration (try AI chat)
# Verify database connectivity
```

## Credit Usage Optimization Tips

### Maximize Free Tier Benefits
```yaml
# Services to prioritize (free tier)
Always_Free:
  - Lambda: 1M requests/month
  - DynamoDB: 25GB storage
  - CloudWatch: 10 metrics, 5GB logs
  - S3: 5GB storage
  - CloudFront: 1TB data transfer

12_Month_Free:
  - EC2: 750 hours t3.micro
  - RDS: 750 hours db.t3.micro + 20GB storage
  - EBS: 30GB General Purpose SSD
  - Elastic Load Balancer: 750 hours
```

### Cost Optimization Strategies
```bash
# Monitor daily costs
aws ce get-cost-and-usage \
    --time-period Start=$(date -d '1 day ago' +%Y-%m-%d),End=$(date +%Y-%m-%d) \
    --granularity DAILY \
    --metrics BlendedCost

# Set up cost anomaly detection
aws ce create-anomaly-detector \
    --anomaly-detector MonitorType=DIMENSIONAL,DimensionKey=SERVICE,MatchOptions=EQUALS,Values=AmazonRDS

# Use Spot instances for batch processing (when needed)
# Implement auto-scaling to minimize resource usage
# Use CloudFront for static content delivery
```

## Expected Timeline and Costs

### Implementation Timeline
- **Day 1**: Tasks 1-3 (Setup and preparation) - 2 hours
- **Day 2**: Tasks 4-6 (Deployment and database) - 2 hours  
- **Day 3**: Tasks 7-9 (Domain and monitoring) - 1 hour
- **Day 4**: Tasks 10-12 (Security and testing) - 1 hour

### Credit Usage Projection
```yaml
# Conservative estimate
Month 1-12: $0-5/month (Free tier covers everything)
Month 13-24: $15-25/month (After free tier expires)
Total credit usage: $180-300 over 24 months

# Your $400 credits should last: 18-24+ months
```

This approach maximizes your AWS credits while providing a production-ready deployment that can easily scale as your business grows!