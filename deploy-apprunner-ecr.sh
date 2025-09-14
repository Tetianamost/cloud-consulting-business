#!/bin/bash

# App Runner ECR Deployment Script
# Builds Docker image, pushes to ECR, and deploys via App Runner

set -e

echo "🚀 Building and deploying to App Runner via ECR..."

# Configuration
SERVICE_NAME="cloud-consulting-prod"
ECR_REPO_NAME="cloud-consulting-business"
REGION="us-east-1"
IMAGE_TAG="latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get AWS account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_URI="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/${ECR_REPO_NAME}"

echo -e "${YELLOW}Configuration:${NC}"
echo "Service Name: $SERVICE_NAME"
echo "ECR Repository: $ECR_REPO_NAME"
echo "Region: $REGION"
echo "Account ID: $ACCOUNT_ID"
echo "ECR URI: $ECR_URI"
echo ""

# Step 1: Create ECR repository if it doesn't exist
echo -e "${YELLOW}Step 1: Creating ECR repository...${NC}"
aws ecr describe-repositories --repository-names $ECR_REPO_NAME --region $REGION >/dev/null 2>&1 || {
    echo "Creating ECR repository: $ECR_REPO_NAME"
    aws ecr create-repository --repository-name $ECR_REPO_NAME --region $REGION
    echo -e "${GREEN}✅ ECR repository created${NC}"
}

# Step 2: Login to ECR
echo -e "${YELLOW}Step 2: Logging into ECR...${NC}"
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.$REGION.amazonaws.com
echo -e "${GREEN}✅ Logged into ECR${NC}"

# Step 3: Build Docker image
echo -e "${YELLOW}Step 3: Building Docker image...${NC}"
docker build -t $ECR_REPO_NAME:$IMAGE_TAG .
echo -e "${GREEN}✅ Docker image built${NC}"

# Step 4: Tag image for ECR
echo -e "${YELLOW}Step 4: Tagging image for ECR...${NC}"
docker tag $ECR_REPO_NAME:$IMAGE_TAG $ECR_URI:$IMAGE_TAG
echo -e "${GREEN}✅ Image tagged${NC}"

# Step 5: Push to ECR
echo -e "${YELLOW}Step 5: Pushing image to ECR...${NC}"
docker push $ECR_URI:$IMAGE_TAG
echo -e "${GREEN}✅ Image pushed to ECR${NC}"

# Step 6: Check if App Runner service exists
echo -e "${YELLOW}Step 6: Checking existing App Runner service...${NC}"
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
    echo "2. Update existing service to use new image"
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
            echo "Updating existing service with new image..."
            aws apprunner update-service \
                --service-arn "$EXISTING_SERVICE" \
                --source-configuration '{
                    "ImageRepository": {
                        "ImageIdentifier": "'$ECR_URI:$IMAGE_TAG'",
                        "ImageConfiguration": {
                            "Port": "80",
                            "RuntimeEnvironmentVariables": {
                                "PORT": "80",
                                "GIN_MODE": "release",
                                "CHAT_MODE": "polling",
                                "AWS_REGION": "us-east-1",
                                "BEDROCK_REGION": "us-east-1",
                                "BEDROCK_MODEL_ID": "amazon.nova-lite-v1:0",
                                "AWS_SES_REGION": "us-east-1",
                                "SES_SENDER_EMAIL": "info@cloudpartner.pro",
                                "CORS_ALLOWED_ORIGINS": "https://cloudpartner.pro,https://www.cloudpartner.pro"
                            }
                        },
                        "ImageRepositoryType": "ECR"
                    },
                    "AutoDeploymentsEnabled": false
                }' \
                --region $REGION
            
            echo -e "${GREEN}✅ Service update initiated${NC}"
            
            # Wait for service to be running
            echo -e "${YELLOW}⏳ Waiting for service to be running...${NC}"
            aws apprunner wait service-running \
                --service-arn "$EXISTING_SERVICE" \
                --region $REGION
            
            # Get service URL
            SERVICE_URL=$(aws apprunner describe-service \
                --service-arn "$EXISTING_SERVICE" \
                --region $REGION \
                --query 'Service.ServiceUrl' \
                --output text)
            
            echo -e "${GREEN}🎉 Service updated successfully!${NC}"
            echo "Service URL: https://$SERVICE_URL"
            exit 0
            ;;
        3)
            exit 0
            ;;
    esac
fi

# Step 7: Create new App Runner service with ECR image
echo -e "${YELLOW}Step 7: Creating App Runner service with ECR image...${NC}"

# Generate JWT secret
JWT_SECRET=$(openssl rand -base64 32)

# Create service configuration
cat > temp-apprunner-ecr-config.json << EOF
{
    "ServiceName": "$SERVICE_NAME",
    "SourceConfiguration": {
        "ImageRepository": {
            "ImageIdentifier": "$ECR_URI:$IMAGE_TAG",
            "ImageConfiguration": {
                "Port": "80",
                "RuntimeEnvironmentVariables": {
                    "PORT": "80",
                    "GIN_MODE": "release",
                    "CHAT_MODE": "polling",
                    "AWS_REGION": "us-east-1",
                    "BEDROCK_REGION": "us-east-1",
                    "BEDROCK_MODEL_ID": "amazon.nova-lite-v1:0",
                    "AWS_SES_REGION": "us-east-1",
                    "SES_SENDER_EMAIL": "info@cloudpartner.pro",
                    "CORS_ALLOWED_ORIGINS": "https://cloudpartner.pro,https://www.cloudpartner.pro",
                    "DATABASE_URL": "postgresql://postgres:CloudConsulting2024!@consulting-prod.c2v4ggsygj4q.us-east-1.rds.amazonaws.com:5432/postgres",
                    "JWT_SECRET": "$JWT_SECRET"
                }
            },
            "ImageRepositoryType": "ECR"
        },
        "AutoDeploymentsEnabled": false
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
    --cli-input-json file://temp-apprunner-ecr-config.json \
    --region $REGION \
    --query 'Service.ServiceArn' \
    --output text)

# Clean up temp file
rm temp-apprunner-ecr-config.json

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
echo "ECR Image: $ECR_URI:$IMAGE_TAG"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Test the deployment: curl https://$SERVICE_URL/health"
echo "2. Set up custom domain if needed"
echo ""
echo -e "${YELLOW}For future updates:${NC}"
echo "Run this script again to build and deploy new versions"
echo ""
echo -e "${YELLOW}Important:${NC}"
echo "JWT Secret generated: $JWT_SECRET"
echo "Store this securely for your records."
