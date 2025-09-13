#!/bin/bash

# AWS App Runner + RDS Deployment Script
# This script creates RDS database and App Runner service

set -e

echo "🚀 Starting AWS App Runner + RDS Deployment..."

# Configuration
DB_INSTANCE_ID="consulting-prod"
DB_PASSWORD="CloudConsulting2024!"
APP_RUNNER_SERVICE="cloud-consulting-prod"
REGION="us-east-1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Step 1: Creating RDS PostgreSQL Database...${NC}"

# Create RDS instance
aws rds create-db-instance \
    --db-instance-identifier $DB_INSTANCE_ID \
    --db-instance-class db.t3.micro \
    --engine postgres \
    --engine-version 15.4 \
    --master-username postgres \
    --master-user-password "$DB_PASSWORD" \
    --allocated-storage 20 \
    --storage-type gp2 \
    --storage-encrypted \
    --backup-retention-period 7 \
    --no-multi-az \
    --publicly-accessible \
    --region $REGION

echo -e "${YELLOW}Waiting for RDS instance to be available (this takes 5-10 minutes)...${NC}"

# Wait for RDS to be available
aws rds wait db-instance-available \
    --db-instance-identifier $DB_INSTANCE_ID \
    --region $REGION

# Get RDS endpoint
RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier $DB_INSTANCE_ID \
    --region $REGION \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

echo -e "${GREEN}✅ RDS Database created successfully!${NC}"
echo -e "Database Endpoint: ${GREEN}$RDS_ENDPOINT${NC}"

echo -e "${YELLOW}Step 2: Creating App Runner Service...${NC}"

# Create service.json for App Runner
cat > service.json << EOF
{
  "ServiceName": "$APP_RUNNER_SERVICE",
  "SourceConfiguration": {
    "AutoDeploymentsEnabled": true,
    "CodeRepository": {
      "RepositoryUrl": "https://github.com/YOUR_USERNAME/YOUR_REPO",
      "SourceCodeVersion": {
        "Type": "BRANCH",
        "Value": "main"
      },
      "CodeConfiguration": {
        "ConfigurationSource": "REPOSITORY"
      }
    }
  },
  "InstanceConfiguration": {
    "Cpu": "0.25 vCPU",
    "Memory": "0.5 GB",
    "InstanceRoleArn": ""
  },
  "HealthCheckConfiguration": {
    "Protocol": "HTTP",
    "Path": "/health",
    "Interval": 10,
    "Timeout": 5,
    "HealthyThreshold": 1,
    "UnhealthyThreshold": 5
  },
  "EnvironmentVariables": {
    "DATABASE_URL": "postgresql://postgres:$DB_PASSWORD@$RDS_ENDPOINT:5432/postgres",
    "GIN_MODE": "release",
    "JWT_SECRET": "your-jwt-secret-change-this",
    "AWS_REGION": "$REGION",
    "BEDROCK_REGION": "$REGION",
    "AWS_SES_REGION": "$REGION",
    "SES_SENDER_EMAIL": "info@cloudpartner.pro",
    "CORS_ALLOWED_ORIGINS": "*",
    "CHAT_MODE": "polling",
    "ENABLE_EMAIL_EVENTS": "true"
  }
}
EOF

echo -e "${GREEN}✅ Deployment files created!${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Push your code to GitHub"
echo "2. Update the repository URL in service.json"
echo "3. Run: aws apprunner create-service --cli-input-json file://service.json"
echo ""
echo -e "${GREEN}Database Connection String:${NC}"
echo "postgresql://postgres:$DB_PASSWORD@$RDS_ENDPOINT:5432/postgres"