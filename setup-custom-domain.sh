#!/bin/bash

# Custom Domain Setup for App Runner
# This script sets up SSL certificate and custom domain for your App Runner service

set -e

# Configuration
DOMAIN_NAME="cloudpartner.pro"
SERVICE_NAME="cloud-consulting-prod"
REGION="us-east-1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔒 Setting up custom domain for App Runner...${NC}"

# Get service ARN
echo "Getting App Runner service ARN..."
SERVICE_ARN=$(aws apprunner list-services \
    --region $REGION \
    --query "ServiceSummaryList[?ServiceName=='$SERVICE_NAME'].ServiceArn" \
    --output text)

if [ -z "$SERVICE_ARN" ]; then
    echo -e "${RED}❌ App Runner service '$SERVICE_NAME' not found!${NC}"
    echo "Please deploy your App Runner service first using deploy-now.sh"
    exit 1
fi

echo -e "${GREEN}✅ Found service: $SERVICE_ARN${NC}"

# Step 1: Request SSL Certificate
echo -e "${YELLOW}Step 1: Requesting SSL certificate...${NC}"

CERT_ARN=$(aws acm request-certificate \
    --domain-name $DOMAIN_NAME \
    --subject-alternative-names "www.$DOMAIN_NAME" \
    --validation-method DNS \
    --region $REGION \
    --query 'CertificateArn' \
    --output text)

echo -e "${GREEN}✅ Certificate requested: $CERT_ARN${NC}"

# Step 2: Get DNS validation records
echo -e "${YELLOW}Step 2: Getting DNS validation records...${NC}"
echo "Waiting for certificate details..."
sleep 10

aws acm describe-certificate \
    --certificate-arn $CERT_ARN \
    --region $REGION \
    --query 'Certificate.DomainValidationOptions[*].[DomainName,ResourceRecord.Name,ResourceRecord.Value]' \
    --output table

echo -e "${YELLOW}📋 IMPORTANT: Add the DNS validation records above to your domain's DNS settings${NC}"
echo -e "${YELLOW}This is required before proceeding to the next step.${NC}"
echo ""
echo "Press Enter when you have added the DNS validation records..."
read -r

# Step 3: Wait for certificate validation
echo -e "${YELLOW}Step 3: Waiting for certificate validation...${NC}"
echo "This may take several minutes..."

aws acm wait certificate-validated \
    --certificate-arn $CERT_ARN \
    --region $REGION

echo -e "${GREEN}✅ Certificate validated successfully!${NC}"

# Step 4: Associate custom domain with App Runner
echo -e "${YELLOW}Step 4: Associating custom domain with App Runner...${NC}"

aws apprunner associate-custom-domain \
    --service-arn "$SERVICE_ARN" \
    --domain-name "$DOMAIN_NAME" \
    --enable-www-subdomain \
    --region $REGION

echo -e "${GREEN}✅ Custom domain association initiated!${NC}"

# Step 5: Get CNAME records for domain
echo -e "${YELLOW}Step 5: Getting CNAME records for your domain...${NC}"
echo "Waiting for domain association details..."
sleep 15

echo -e "${YELLOW}📋 Add these CNAME records to your domain's DNS:${NC}"
aws apprunner describe-custom-domains \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'CustomDomains[*].DomainName' \
    --output table

# Get detailed DNS records
aws apprunner describe-custom-domains \
    --service-arn "$SERVICE_ARN" \
    --region $REGION

echo ""
echo -e "${GREEN}🎉 Custom domain setup initiated!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Add the CNAME records shown above to your domain's DNS"
echo "2. Wait for DNS propagation (up to 48 hours)"
echo "3. Your site will be available at: https://$DOMAIN_NAME"
echo "4. And also at: https://www.$DOMAIN_NAME"
echo ""
echo -e "${YELLOW}To check status:${NC}"
echo "aws apprunner describe-custom-domains --service-arn \"$SERVICE_ARN\" --region $REGION"