#!/bin/bash

# Get RDS endpoint
RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier consulting-prod \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

# Replace with your GitHub repository URL
GITHUB_REPO_URL="https://github.com/YOUR_USERNAME/YOUR_REPO"

echo "Creating App Runner service with RDS endpoint: $RDS_ENDPOINT"

# Create App Runner service
aws apprunner create-service \
    --service-name "cloud-consulting-prod" \
    --source-configuration '{
        "AutoDeploymentsEnabled": true,
        "CodeRepository": {
            "RepositoryUrl": "'$GITHUB_REPO_URL'",
            "SourceCodeVersion": {
                "Type": "BRANCH",
                "Value": "main"
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
    }'

# Get the service ARN
SERVICE_ARN=$(aws apprunner list-services \
    --query 'ServiceSummaryList[?ServiceName==`cloud-consulting-prod`].ServiceArn' \
    --output text)

echo "Service ARN: $SERVICE_ARN"

# Wait for service to be running
echo "Waiting for App Runner service to be ready..."
aws apprunner wait service-running --service-arn "$SERVICE_ARN"

# Get the service URL
SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --query 'Service.ServiceUrl' \
    --output text)

echo "🎉 Deployment complete!"
echo "Service URL: https://$SERVICE_URL"
echo "Health Check: https://$SERVICE_URL/health"