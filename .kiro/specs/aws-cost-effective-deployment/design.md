# Cost-Effective AWS Deployment Design

## Overview

This document provides a cost-effective architecture for deploying the Cloud Consulting Platform on AWS, designed for startups and early-stage businesses. The architecture prioritizes minimal costs while maintaining functionality and providing a clear upgrade path as the business grows.

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
              │  Let's Encrypt  │
              │   SSL (FREE)    │
              └─────────┬───────┘
                        │
    ┌───────────────────▼────────────────────┐
    │         Single EC2 Instance            │
    │           t3.micro (FREE)              │
    │  ┌─────────────────────────────────┐   │
    │  │           Nginx                 │   │
    │  │    (Reverse Proxy + SSL)        │   │
    │  └─────────┬───────────────────────┘   │
    │            │                           │
    │  ┌─────────▼─────────┐ ┌─────────────┐ │
    │  │   Frontend        │ │   Backend   │ │
    │  │   (React)         │ │    (Go)     │ │
    │  │   Port 3000       │ │  Port 8061  │ │
    │  └───────────────────┘ └─────┬───────┘ │
    │                              │         │
    │  ┌─────────────────────────────▼─────┐ │
    │  │        PostgreSQL               │ │
    │  │      (Local Install)            │ │
    │  └─────────────────────────────────┘ │
    │                                      │
    │  ┌─────────────────────────────────┐ │
    │  │          Redis                  │ │
    │  │      (Local Install)            │ │
    │  └─────────────────────────────────┘ │
    └──────────────────────────────────────┘
                        │
              ┌─────────▼───────┐
              │   S3 Backups    │
              │  ($1-2/month)   │
              └─────────────────┘

External Services (Already Configured):
├── AWS Bedrock (Pay per use)
├── AWS SES (Pay per email)
└── CloudWatch (Free tier)
```

### Cost Breakdown

#### Option 1: Ultra-Low Cost (Recommended for MVP)
- **EC2 t3.micro**: FREE (12 months free tier) then $8.50/month
- **EBS Storage**: 30GB FREE (free tier) then $3/month
- **Route 53**: $0.50/month
- **S3 Backups**: $1-2/month
- **Data Transfer**: 15GB FREE then minimal
- **Total Year 1**: ~$2-4/month
- **Total Year 2+**: ~$12-15/month

#### Option 2: Slightly More Robust
- **EC2 t3.small**: $15/month
- **EBS Storage**: $3/month
- **Route 53**: $0.50/month
- **S3 Backups**: $2/month
- **Total**: ~$20-25/month

## Components and Interfaces

### 1. EC2 Instance Configuration

#### Instance Specifications
```yaml
InstanceType: t3.micro  # Free tier eligible
AMI: Amazon Linux 2023
Storage: 30GB gp3 EBS (free tier eligible)
SecurityGroup: web-server-sg
KeyPair: your-key-pair
```

#### Security Group Configuration
```yaml
SecurityGroup: web-server-sg
InboundRules:
  - Port: 22 (SSH)
    Source: Your IP only
  - Port: 80 (HTTP)
    Source: 0.0.0.0/0
  - Port: 443 (HTTPS)
    Source: 0.0.0.0/0
OutboundRules:
  - All traffic to 0.0.0.0/0
```

### 2. Software Stack Installation

#### System Setup Script
```bash
#!/bin/bash
# install-stack.sh

# Update system
sudo yum update -y

# Install Docker
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Nginx
sudo yum install -y nginx
sudo systemctl enable nginx

# Install Certbot for Let's Encrypt
sudo yum install -y certbot python3-certbot-nginx

# Install PostgreSQL
sudo yum install -y postgresql15-server postgresql15
sudo postgresql-setup --initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Install Redis
sudo yum install -y redis
sudo systemctl enable redis
sudo systemctl start redis

# Install Node.js (for frontend builds)
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs

# Install Go (for backend builds)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

### 3. Application Deployment

#### Docker Compose Configuration
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "8061:8061"
    environment:
      - GIN_MODE=release
      - DATABASE_URL=postgresql://postgres:password@localhost:5432/consulting
      - REDIS_URL=redis://localhost:6379
      - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
      - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
      - AWS_BEARER_TOKEN_BEDROCK=${AWS_BEARER_TOKEN_BEDROCK}
      - SES_SENDER_EMAIL=info@cloudpartner.pro
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
    network_mode: host

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_API_URL=https://api.cloudpartner.pro
    restart: unless-stopped
    network_mode: host
```

### 4. Nginx Configuration

#### Main Nginx Config
```nginx
# /etc/nginx/nginx.conf
server {
    listen 80;
    server_name cloudpartner.pro www.cloudpartner.pro;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name cloudpartner.pro www.cloudpartner.pro;

    ssl_certificate /etc/letsencrypt/live/cloudpartner.pro/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/cloudpartner.pro/privkey.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # Frontend (React)
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend API (Go)
    location /api/ {
        proxy_pass http://localhost:8061;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health check
    location /health {
        proxy_pass http://localhost:8061/health;
    }
}
```

### 5. Database Configuration

#### PostgreSQL Setup
```bash
# Setup PostgreSQL
sudo -u postgres createdb consulting
sudo -u postgres createuser --interactive --pwprompt appuser

# Configure PostgreSQL
sudo tee -a /var/lib/pgsql/data/postgresql.conf << EOF
listen_addresses = 'localhost'
port = 5432
max_connections = 100
shared_buffers = 128MB
EOF

# Configure authentication
sudo tee /var/lib/pgsql/data/pg_hba.conf << EOF
local   all             postgres                                peer
local   all             all                                     md5
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5
EOF

sudo systemctl restart postgresql
```

### 6. Backup Strategy

#### Automated Backup Script
```bash
#!/bin/bash
# backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/tmp/backups"
S3_BUCKET="your-backup-bucket"

mkdir -p $BACKUP_DIR

# Database backup
sudo -u postgres pg_dump consulting > $BACKUP_DIR/db_backup_$DATE.sql

# Application backup
tar -czf $BACKUP_DIR/app_backup_$DATE.tar.gz /home/ec2-user/cloud-consulting

# Upload to S3
aws s3 cp $BACKUP_DIR/db_backup_$DATE.sql s3://$S3_BUCKET/database/
aws s3 cp $BACKUP_DIR/app_backup_$DATE.tar.gz s3://$S3_BUCKET/application/

# Cleanup local backups older than 7 days
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete

# Cleanup S3 backups older than 30 days
aws s3api list-objects-v2 --bucket $S3_BUCKET --prefix database/ --query 'Contents[?LastModified<=`'$(date -d '30 days ago' --iso-8601)'`].Key' --output text | xargs -I {} aws s3 rm s3://$S3_BUCKET/{}
```

#### Cron Job Setup
```bash
# Add to crontab
0 2 * * * /home/ec2-user/scripts/backup.sh >> /var/log/backup.log 2>&1
```

## Data Models

### Environment Configuration

#### Production Environment Variables
```bash
# /home/ec2-user/.env
PORT=8061
GIN_MODE=release
LOG_LEVEL=info

# Database
DATABASE_URL=postgresql://appuser:password@localhost:5432/consulting

# Redis
REDIS_URL=redis://localhost:6379

# AWS Configuration (use your existing values)
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_BEARER_TOKEN_BEDROCK=your_bedrock_token
BEDROCK_REGION=us-east-1
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0
AWS_SES_REGION=us-east-1
SES_SENDER_EMAIL=info@cloudpartner.pro

# Application
CORS_ALLOWED_ORIGINS=https://cloudpartner.pro,https://www.cloudpartner.pro
JWT_SECRET=your_production_jwt_secret
CHAT_MODE=polling
CHAT_POLLING_INTERVAL=3000
```

## Error Handling

### Application-Level Error Handling
```go
// Implement graceful degradation
func (s *Service) handleDatabaseError(err error) error {
    if isConnectionError(err) {
        // Log error and attempt reconnection
        s.logger.Error("Database connection lost, attempting reconnection")
        return s.reconnectDatabase()
    }
    return err
}

// Circuit breaker for external services
func (s *Service) callExternalService(fn func() error) error {
    if s.circuitBreaker.IsOpen() {
        return errors.New("service temporarily unavailable")
    }
    
    err := fn()
    if err != nil {
        s.circuitBreaker.RecordFailure()
        return err
    }
    
    s.circuitBreaker.RecordSuccess()
    return nil
}
```

### System-Level Error Handling
```bash
# Systemd service for auto-restart
sudo tee /etc/systemd/system/cloud-consulting.service << EOF
[Unit]
Description=Cloud Consulting Platform
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/home/ec2-user/cloud-consulting
ExecStart=/usr/local/bin/docker-compose up
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable cloud-consulting
sudo systemctl start cloud-consulting
```

## Testing Strategy

### Health Check Implementation
```go
// Health check endpoint
func (h *HealthHandler) CheckHealth(c *gin.Context) {
    health := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "services": map[string]string{
            "database": h.checkDatabase(),
            "redis": h.checkRedis(),
            "bedrock": h.checkBedrock(),
            "ses": h.checkSES(),
        },
    }
    
    c.JSON(http.StatusOK, health)
}
```

### Monitoring Script
```bash
#!/bin/bash
# monitor.sh

# Check application health
curl -f http://localhost:8061/health || echo "Application health check failed"

# Check disk space
df -h | awk '$5 > 80 {print "Disk usage high: " $0}'

# Check memory usage
free -m | awk 'NR==2{printf "Memory Usage: %s/%sMB (%.2f%%)\n", $3,$2,$3*100/$2 }'

# Check PostgreSQL
sudo systemctl is-active postgresql || echo "PostgreSQL is not running"

# Check Redis
redis-cli ping || echo "Redis is not responding"
```

## Security Considerations

### SSL Certificate Management
```bash
# Initial certificate setup
sudo certbot --nginx -d cloudpartner.pro -d www.cloudpartner.pro

# Auto-renewal setup
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

### Security Hardening
```bash
# Disable password authentication
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# Setup fail2ban
sudo yum install -y fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# Configure firewall
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --reload
```

## Upgrade Path

### Migration to RDS (When Ready)
```bash
# 1. Create RDS instance
aws rds create-db-instance \
    --db-instance-identifier cloud-consulting-prod \
    --db-instance-class db.t3.micro \
    --engine postgres \
    --master-username postgres \
    --master-user-password "secure_password" \
    --allocated-storage 20

# 2. Migrate data
pg_dump consulting | psql -h rds-endpoint -U postgres -d consulting

# 3. Update application configuration
# Change DATABASE_URL to point to RDS endpoint
```

### Migration to Load Balancer (When Traffic Grows)
```bash
# 1. Create Application Load Balancer
aws elbv2 create-load-balancer \
    --name cloud-consulting-alb \
    --subnets subnet-xxx subnet-yyy \
    --security-groups sg-xxx

# 2. Create multiple EC2 instances
# 3. Update DNS to point to ALB
```

This cost-effective design provides a solid foundation for your startup while keeping costs minimal. You can start with the ultra-low cost option and scale up as your business grows and generates revenue.