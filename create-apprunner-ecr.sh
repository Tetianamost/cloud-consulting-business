#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="cloud-consulting-prod"
REGION="us-east-1"
ACCOUNT_ID="757742990331"
ECR_URI="757742990331.dkr.ecr.us-east-1.amazonaws.com/cloud-consulting-business"
ROLE_NAME="AppRunnerECRAccessRole"

echo -e "${YELLOW}🚀 Creating App Runner service with ECR image...${NC}"
echo "Service Name: $SERVICE_NAME"
echo "ECR Image: $ECR_URI:latest"
echo "Region: $REGION"
echo ""

# Step 1: Create IAM role for App Runner to access ECR
echo -e "${YELLOW}Step 1: Creating IAM role for ECR access...${NC}"

# Check if role already exists
if aws iam get-role --role-name $ROLE_NAME >/dev/null 2>&1; then
    echo "IAM role $ROLE_NAME already exists"
    ROLE_ARN=$(aws iam get-role --role-name $ROLE_NAME --query 'Role.Arn' --output text)
else
    # Create trust policy
    cat > /tmp/trust-policy.json << EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "build.apprunner.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

    # Create the role
    ROLE_ARN=$(aws iam create-role \
        --role-name $ROLE_NAME \
        --assume-role-policy-document file:///tmp/trust-policy.json \
        --query 'Role.Arn' \
        --output text)

    # Attach ECR access policy
    aws iam attach-role-policy \
        --role-name $ROLE_NAME \
        --policy-arn arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess

    echo -e "${GREEN}✅ IAM role created: $ROLE_ARN${NC}"
    
    # Clean up temp file
    rm /tmp/trust-policy.json
fi

# Step 2: Create App Runner service
echo -e "${YELLOW}Step 2: Creating App Runner service...${NC}"

# Check if service already exists
if aws apprunner describe-service --service-arn "arn:aws:apprunner:$REGION:$ACCOUNT_ID:service/$SERVICE_NAME" >/dev/null 2>&1; then
    echo -e "${YELLOW}Service $SERVICE_NAME already exists. Would you like to update it? (y/n)${NC}"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        SERVICE_ARN="arn:aws:apprunner:$REGION:$ACCOUNT_ID:service/$SERVICE_NAME"
        aws apprunner update-service \
            --service-arn "$SERVICE_ARN" \
            --source-configuration '{
                "ImageRepository": {
                    "ImageIdentifier": "'$ECR_URI:latest'",
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
    else
        echo "Skipping update."
        exit 0
    fi
else
    # Create new service
    aws apprunner create-service \
        --service-name "$SERVICE_NAME" \
        --source-configuration '{
            "ImageRepository": {
                "ImageIdentifier": "'$ECR_URI:latest'",
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
        --instance-configuration '{
            "Cpu": "0.25 vCPU",
            "Memory": "0.5 GB"
        }' \
        --region $REGION
    
    # Get the service ARN
    SERVICE_ARN=$(aws apprunner list-services \
        --region $REGION \
        --query 'ServiceSummaryList[?ServiceName==`'$SERVICE_NAME'`].ServiceArn' \
        --output text)
    
    echo -e "${GREEN}✅ Service created: $SERVICE_ARN${NC}"
fi

# Step 3: Wait for service to be ready
echo -e "${YELLOW}Step 3: Waiting for service to be running...${NC}"
echo "This may take several minutes..."

aws apprunner wait service-running --service-arn "$SERVICE_ARN" --region $REGION

# Step 4: Get service URL
echo -e "${YELLOW}Step 4: Getting service information...${NC}"

SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.ServiceUrl' \
    --output text)

SERVICE_STATUS=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region $REGION \
    --query 'Service.Status' \
    --output text)

echo ""
echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo -e "${GREEN}Service Name: $SERVICE_NAME${NC}"
echo -e "${GREEN}Status: $SERVICE_STATUS${NC}"
echo -e "${GREEN}Service URL: https://$SERVICE_URL${NC}"
echo -e "${GREEN}Service ARN: $SERVICE_ARN${NC}"
echo ""
echo -e "${YELLOW}You can now access your application at: https://$SERVICE_URL${NC}"
