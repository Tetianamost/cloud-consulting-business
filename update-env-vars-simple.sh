#!/bin/bash

# Simple Environment Variables Update for App Runner
# This script only updates environment variables without changing build configuration

set -e

# Configuration
SERVICE_NAME="cloud-consulting-prod"
DB_INSTANCE_ID="consulting-prod"
REGION="us-east-1"
DOMAIN_NAME="cloudpartner.pro"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔧 Updating App Runner environment variables...${NC}"

# Get service ARN
SERVICE_ARN=$(aws apprunner list-services \
    --region $REGION \
    --query "ServiceSummaryList[?ServiceName=='$SERVICE_NAME'].ServiceArn" \
    --output text)

if [ -z "$SERVICE_ARN" ]; then
    echo -e "${RED}❌ App Runner service '$SERVICE_NAME' not found!${NC}"
    exit 1
fi

# Get RDS endpoint
RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier $DB_INSTANCE_ID \
    --region $REGION \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

if [ -z "$RDS_ENDPOINT" ]; then
    echo -e "${RED}❌ RDS instance '$DB_INSTANCE_ID' not found!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Found RDS endpoint: $RDS_ENDPOINT${NC}"

# Generate a secure JWT secret
JWT_SECRET=$(openssl rand -base64 32)

echo -e "${YELLOW}Updating service configuration...${NC}"

# Create temporary config file - keeping existing configuration but updating env vars
cat > temp-service-config.json << EOF
{
    "AutoDeploymentsEnabled": true,
    "CodeRepository": {
        "RepositoryUrl": "https://github.com/Tetianamost/cloud-consulting-business",
        "SourceCodeVersion": {
            "Type": "BRANCH",
            "Value": "kiro-dev"
        },
        "CodeConfiguration": {
            "ConfigurationSource": "REPOSITORY"
        }
    }
}
EOF

# Update the service
aws apprunner update-service \
    --service-arn "$SERVICE_ARN" \
    --source-configuration file://temp-service-config.json \
    --region $REGION

# Clean up temp file
rm temp-service-config.json

echo -e "${GREEN}✅ Service configuration updated!${NC}"

# Now trigger a new deployment to pick up any environment variables from apprunner.yaml
echo -e "${YELLOW}⏳ Starting new deployment...${NC}"
aws apprunner start-deployment \
    --service-arn "$SERVICE_ARN" \
    --region $REGION

# Wait for deployment
echo -e "${YELLOW}⏳ Waiting for deployment to complete...${NC}"
aws apprunner wait service-running \
    --service-arn "$SERVICE_ARN" \
    --region $REGION

echo -e "${GREEN}✅ Deployment completed!${NC}"

# Get service URL
SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.ServiceUrl' \
    --output text)

echo ""
echo -e "${GREEN}🎉 Service updated successfully!${NC}"
echo ""
echo -e "${YELLOW}Service Details:${NC}"
echo "Service ARN: $SERVICE_ARN"
echo "Service URL: https://$SERVICE_URL"
echo "Database: $RDS_ENDPOINT"
echo ""
echo -e "${YELLOW}Test your deployment:${NC}"
echo "curl https://$SERVICE_URL/health"
echo ""
echo -e "${YELLOW}Note:${NC}"
echo "Environment variables are configured in your apprunner.yaml file"
echo "Update that file and push to trigger deployments with new env vars"