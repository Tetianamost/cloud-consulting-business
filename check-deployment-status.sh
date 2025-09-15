#!/bin/bash

# Check App Runner Deployment Status
# This script checks the status of your App Runner service and custom domain

set -e

# Configuration
SERVICE_NAME="cloud-consulting-prod"
DOMAIN_NAME="cloudpartner.pro"
REGION="us-east-1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Checking App Runner deployment status...${NC}"

# Get service ARN
SERVICE_ARN=$(aws apprunner list-services \
    --region $REGION \
    --query "ServiceSummaryList[?ServiceName=='$SERVICE_NAME'].ServiceArn" \
    --output text)

if [ -z "$SERVICE_ARN" ]; then
    echo -e "${RED}❌ App Runner service '$SERVICE_NAME' not found!${NC}"
    exit 1
fi

# Get service details
echo -e "${YELLOW}📊 Service Status:${NC}"
SERVICE_STATUS=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.Status' \
    --output text)

SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.ServiceUrl' \
    --output text)

echo "Status: $SERVICE_STATUS"
echo "URL: https://$SERVICE_URL"

# Check service health
echo ""
echo -e "${YELLOW}🏥 Health Check:${NC}"
if curl -s -f "https://$SERVICE_URL/health" > /dev/null; then
    echo -e "${GREEN}✅ Service is healthy${NC}"
    echo "Health endpoint: https://$SERVICE_URL/health"
else
    echo -e "${RED}❌ Service health check failed${NC}"
    echo "Check logs: aws logs tail /aws/apprunner/$SERVICE_NAME --region $REGION"
fi

# Check API endpoint
echo ""
echo -e "${YELLOW}🔌 API Check:${NC}"
if curl -s -f "https://$SERVICE_URL/api/health" > /dev/null; then
    echo -e "${GREEN}✅ API is responding${NC}"
    echo "API endpoint: https://$SERVICE_URL/api/health"
else
    echo -e "${RED}❌ API health check failed${NC}"
fi

# Check custom domain status
echo ""
echo -e "${YELLOW}🌐 Custom Domain Status:${NC}"
CUSTOM_DOMAINS=$(aws apprunner describe-custom-domains \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'CustomDomains' \
    --output json 2>/dev/null || echo "[]")

if [ "$CUSTOM_DOMAINS" = "[]" ]; then
    echo -e "${YELLOW}⚠️  No custom domain configured${NC}"
    echo "Run ./setup-custom-domain.sh to set up your custom domain"
else
    echo -e "${GREEN}✅ Custom domain configured${NC}"
    aws apprunner describe-custom-domains \
        --service-arn "$SERVICE_ARN" \
        --region $REGION \
        --query 'CustomDomains[*].[DomainName,Status]' \
        --output table
fi

# Check DNS resolution
echo ""
echo -e "${YELLOW}🔍 DNS Check:${NC}"
if nslookup $DOMAIN_NAME > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Domain resolves${NC}"
    echo "Testing HTTPS connection to $DOMAIN_NAME..."
    if curl -s -f "https://$DOMAIN_NAME/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Custom domain is working!${NC}"
        echo "Your site is live at: https://$DOMAIN_NAME"
    else
        echo -e "${YELLOW}⚠️  Domain resolves but HTTPS not working yet${NC}"
        echo "This is normal during DNS propagation (can take up to 48 hours)"
    fi
else
    echo -e "${YELLOW}⚠️  Domain not resolving yet${NC}"
    echo "DNS propagation in progress..."
fi

# Check recent deployments
echo ""
echo -e "${YELLOW}📈 Recent Activity:${NC}"
aws apprunner list-operations \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'OperationSummaryList[0:3].[Type,Status,StartedAt]' \
    --output table

# Show environment variables (without sensitive values)
echo ""
echo -e "${YELLOW}⚙️  Environment Configuration:${NC}"
aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.SourceConfiguration.CodeRepository.CodeConfiguration.CodeConfigurationValues.RuntimeEnvironmentVariables' \
    --output json | jq -r 'to_entries[] | select(.key | test("SECRET|PASSWORD|TOKEN") | not) | "\(.key): \(.value)"' 2>/dev/null || echo "Environment variables configured (use AWS console to view)"

# Show useful commands
echo ""
echo -e "${BLUE}📋 Useful Commands:${NC}"
echo ""
echo -e "${YELLOW}View logs:${NC}"
echo "aws logs tail /aws/apprunner/$SERVICE_NAME --region $REGION --follow"
echo ""
echo -e "${YELLOW}Update service:${NC}"
echo "./update-env-vars.sh"
echo ""
echo -e "${YELLOW}Set up custom domain:${NC}"
echo "./setup-custom-domain.sh"
echo ""
echo -e "${YELLOW}Test endpoints:${NC}"
echo "curl https://$SERVICE_URL/health"
echo "curl https://$SERVICE_URL/api/health"
if [ "$CUSTOM_DOMAINS" != "[]" ]; then
    echo "curl https://$DOMAIN_NAME/health"
fi

echo ""
echo -e "${GREEN}✅ Status check complete!${NC}"