#!/bin/bash

# Quick App Runner + RDS Deployment
# Run this after RDS is created

set -e

echo "🚀 Creating App Runner service..."

# Get RDS endpoint
RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier consulting-prod \
    --region us-east-1 \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

echo "Using RDS endpoint: $RDS_ENDPOINT"

# Your GitHub repository URL
GITHUB_REPO_URL="https://github.com/Tetianamost/cloud-consulting-business"

echo "Using GitHub repository: $GITHUB_REPO_URL"

# Create App Runner service
aws apprunner create-service \
    --service-name "cloud-consulting-prod" \
    --source-configuration '{
        "AutoDeploymentsEnabled": true,
        "CodeRepository": {
            "RepositoryUrl": "'$GITHUB_REPO_URL'",
            "SourceCodeVersion": {
                "Type": "BRANCH",
                "Value": "kiro-dev"
            },
            "CodeConfiguration": {
                "ConfigurationSource": "REPOSITORY"
            }
        }
    }' \
    --instance-configuration '{
        "Cpu": "0.25 vCPU",
        "Memory": "0.5 GB"
    }' \
    --health-check-configuration '{
        "Protocol": "HTTP",
        "Path": "/health",
        "Interval": 10,
        "Timeout": 5,
        "HealthyThreshold": 1,
        "UnhealthyThreshold": 5
    }' \
    --region us-east-1

echo "⏳ Waiting for App Runner service to be ready..."

# Get service ARN
SERVICE_ARN=$(aws apprunner list-services \
    --region us-east-1 \
    --query 'ServiceSummaryList[?ServiceName==`cloud-consulting-prod`].ServiceArn' \
    --output text)

echo "Service ARN: $SERVICE_ARN"

# Wait for service to be running
aws apprunner wait service-running --service-arn "$SERVICE_ARN" --region us-east-1

# Get service URL
SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region us-east-1 \
    --query 'Service.ServiceUrl' \
    --output text)

echo ""
echo "🎉 Deployment Complete!"
echo "================================"
echo "App Runner URL: https://$SERVICE_URL"
echo "Health Check: https://$SERVICE_URL/health"
echo "Database: $RDS_ENDPOINT"
echo ""
echo "Next steps:"
echo "1. Add environment variables to App Runner service"
echo "2. Run database migrations"
echo "3. Test the deployment"