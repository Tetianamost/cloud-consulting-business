#!/bin/bash

# Update App Runner Environment Variables
# This script updates your App Runner service with production environment variables

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

# Create temporary config file
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
            "ConfigurationSource": "REPOSITORY",
            "CodeConfigurationValues": {
                "Runtime": "DOCKER",
                "BuildCommand": "",
                "StartCommand": "",
                "RuntimeEnvironmentVariables": {
                    "DATABASE_URL": "postgresql://postgres:CloudConsulting2024!@$RDS_ENDPOINT:5432/postgres",
                    "GIN_MODE": "release",
                    "PORT": "80",
                    "LOG_LEVEL": "2",
                    "AWS_REGION": "$REGION",
                    "BEDROCK_REGION": "$REGION",
                    "BEDROCK_MODEL_ID": "amazon.nova-lite-v1:0",
                    "BEDROCK_BASE_URL": "https://bedrock-runtime.$REGION.amazonaws.com",
                    "BEDROCK_TIMEOUT_SECONDS": "30",
                    "AWS_SES_REGION": "$REGION",
                    "SES_SENDER_EMAIL": "info@$DOMAIN_NAME",
                    "SES_REPLY_TO_EMAIL": "info@$DOMAIN_NAME",
                    "SES_TIMEOUT_SECONDS": "30",
                    "CORS_ALLOWED_ORIGINS": "https://$DOMAIN_NAME,https://www.$DOMAIN_NAME",
                    "CHAT_MODE": "polling",
                    "ENABLE_EMAIL_EVENTS": "true",
                    "ENABLE_CHAT_METRICS": "true",
                    "ENABLE_PERFORMANCE_MONITORING": "true",
                    "JWT_SECRET": "$JWT_SECRET",
                    "REDIS_URL": "",
                    "CACHE_TTL_SECONDS": "300"
                }
            }
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

# Wait for deployment
echo -e "${YELLOW}⏳ Waiting for deployment to complete...${NC}"
aws apprunner wait service-updated \
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
echo -e "${GREEN}🎉 Environment variables updated successfully!${NC}"
echo ""
echo -e "${YELLOW}Service Details:${NC}"
echo "Service ARN: $SERVICE_ARN"
echo "Service URL: https://$SERVICE_URL"
echo "Database: $RDS_ENDPOINT"
echo "Custom Domain: https://$DOMAIN_NAME (when DNS is configured)"
echo ""
echo -e "${YELLOW}Test your deployment:${NC}"
echo "curl https://$SERVICE_URL/health"
echo ""
echo -e "${YELLOW}Important Security Note:${NC}"
echo "JWT Secret has been generated and set. Store this securely:"
echo "JWT_SECRET: $JWT_SECRET"