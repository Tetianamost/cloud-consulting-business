#!/bin/bash

# Fixed App Runner Deployment Script
# This script creates App Runner service with proper GitHub authentication

set -e

echo "🚀 Creating App Runner service with GitHub authentication..."

# Configuration
SERVICE_NAME="cloud-consulting-prod"
DB_INSTANCE_ID="consulting-prod"
REGION="us-east-1"
GITHUB_REPO_URL="https://github.com/Tetianamost/cloud-consulting-business"
BRANCH="kiro-dev"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if service already exists
echo "Checking if service already exists..."
EXISTING_SERVICE=$(aws apprunner list-services \
    --region $REGION \
    --query "ServiceSummaryList[?ServiceName=='$SERVICE_NAME'].ServiceArn" \
    --output text 2>/dev/null || echo "")

if [ ! -z "$EXISTING_SERVICE" ]; then
    echo -e "${YELLOW}⚠️  Service '$SERVICE_NAME' already exists!${NC}"
    echo "Service ARN: $EXISTING_SERVICE"
    echo ""
    echo "Options:"
    echo "1. Delete existing service and recreate"
    echo "2. Update existing service"
    echo "3. Exit"
    read -p "Choose option (1/2/3): " choice
    
    case $choice in
        1)
            echo "Deleting existing service..."
            aws apprunner delete-service \
                --service-arn "$EXISTING_SERVICE" \
                --region $REGION
            echo "Waiting for service deletion..."
            aws apprunner wait service-deleted \
                --service-arn "$EXISTING_SERVICE" \
                --region $REGION
            echo -e "${GREEN}✅ Service deleted${NC}"
            ;;
        2)
            echo "Use ./update-env-vars.sh to update existing service"
            exit 0
            ;;
        3)
            exit 0
            ;;
    esac
fi

# Get RDS endpoint
echo "Getting RDS endpoint..."
RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier $DB_INSTANCE_ID \
    --region $REGION \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

if [ -z "$RDS_ENDPOINT" ]; then
    echo -e "${RED}❌ RDS instance '$DB_INSTANCE_ID' not found!${NC}"
    echo "Please deploy RDS first using: ./deploy-apprunner-rds.sh"
    exit 1
fi

echo -e "${GREEN}✅ Found RDS endpoint: $RDS_ENDPOINT${NC}"

# Get existing GitHub connection
echo "Getting existing GitHub connection..."
CONNECTION_ARN=$(aws apprunner list-connections \
    --region $REGION \
    --query 'ConnectionSummaryList[?Status==`AVAILABLE`].ConnectionArn' \
    --output text)

if [ -z "$CONNECTION_ARN" ]; then
    echo -e "${RED}❌ No available GitHub connection found!${NC}"
    echo "Please set up GitHub connection first in AWS Console"
    exit 1
fi

echo -e "${GREEN}✅ Using existing GitHub connection: $CONNECTION_ARN${NC}"

# Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)

echo -e "${YELLOW}Creating App Runner service...${NC}"

# Create service configuration
cat > temp-apprunner-config.json << EOF
{
    "ServiceName": "$SERVICE_NAME",
    "SourceConfiguration": {
        "AutoDeploymentsEnabled": true,
        "CodeRepository": {
            "RepositoryUrl": "$GITHUB_REPO_URL",
            "SourceCodeVersion": {
                "Type": "BRANCH",
                "Value": "$BRANCH"
            },
            "CodeConfiguration": {
                "ConfigurationSource": "REPOSITORY"
            }
        },
        "AuthenticationConfiguration": {
            "ConnectionArn": "$CONNECTION_ARN"
        }
    },
    "InstanceConfiguration": {
        "Cpu": "0.25 vCPU",
        "Memory": "0.5 GB"
    },
    "HealthCheckConfiguration": {
        "Protocol": "HTTP",
        "Path": "/health",
        "Interval": 10,
        "Timeout": 5,
        "HealthyThreshold": 1,
        "UnhealthyThreshold": 5
    }
}
EOF

# Create the service
SERVICE_ARN=$(aws apprunner create-service \
    --cli-input-json file://temp-apprunner-config.json \
    --region $REGION \
    --query 'Service.ServiceArn' \
    --output text)

# Clean up temp file
rm temp-apprunner-config.json

echo -e "${GREEN}✅ App Runner service created!${NC}"
echo "Service ARN: $SERVICE_ARN"

# Wait for service to be running
echo -e "${YELLOW}⏳ Waiting for service to be running (this takes 3-5 minutes)...${NC}"
aws apprunner wait service-running \
    --service-arn "$SERVICE_ARN" \
    --region $REGION

# Get service URL
SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.ServiceUrl' \
    --output text)

echo -e "${GREEN}🎉 App Runner service is running!${NC}"
echo ""
echo -e "${YELLOW}Service Details:${NC}"
echo "Service Name: $SERVICE_NAME"
echo "Service ARN: $SERVICE_ARN"
echo "Service URL: https://$SERVICE_URL"
echo "Database: $RDS_ENDPOINT"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Update environment variables: ./update-env-vars.sh"
echo "2. Test the deployment: curl https://$SERVICE_URL/health"
echo "3. Set up custom domain: ./setup-custom-domain.sh"
echo ""
echo -e "${YELLOW}Important:${NC}"
echo "JWT Secret generated: $JWT_SECRET"
echo "Store this securely for your records."