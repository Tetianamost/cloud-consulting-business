#!/bin/bash

# Update App Runner environment variables to fix backend crash
# Service: cloud-consulting-prod
# ARN: arn:aws:apprunner:us-east-1:757742990331:service/cloud-consulting-prod/5d093655c4ba4030bafebd36900fee75

SERVICE_ARN="arn:aws:apprunner:us-east-1:757742990331:service/cloud-consulting-prod/5d093655c4ba4030bafebd36900fee75"

# Get current service configuration
echo "Getting current App Runner service configuration..."
aws apprunner describe-service --service-arn "$SERVICE_ARN" --region us-east-1 > current-config.json

# Extract current environment variables and add/update required ones
cat > env-vars-update.json << 'EOF'
{
  "ServiceArn": "arn:aws:apprunner:us-east-1:757742990331:service/cloud-consulting-prod/5d093655c4ba4030bafebd36900fee75",
  "SourceConfiguration": {
    "ImageRepository": {
      "ImageConfiguration": {
        "RuntimeEnvironmentVariables": {
          "PORT": "8080",
          "GIN_MODE": "release",
          "ENABLE_EMAIL_EVENTS": "false",
          "LOG_LEVEL": "4",
          "JWT_SECRET": "cloud-consulting-demo-secret",
          "CORS_ALLOWED_ORIGINS": "https://5d093655c4ba4030bafebd36900fee75.us-east-1.awsapprunner.com"
        }
      }
    }
  }
}
EOF

# Update the service
echo "Updating App Runner service with environment variables..."
aws apprunner update-service \
  --service-arn "$SERVICE_ARN" \
  --region us-east-1 \
  --source-configuration '{
    "ImageRepository": {
      "ImageConfiguration": {
        "RuntimeEnvironmentVariables": {
          "PORT": "8080",
          "GIN_MODE": "release", 
          "ENABLE_EMAIL_EVENTS": "false",
          "LOG_LEVEL": "4",
          "JWT_SECRET": "cloud-consulting-demo-secret",
          "CORS_ALLOWED_ORIGINS": "https://5d093655c4ba4030bafebd36900fee75.us-east-1.awsapprunner.com"
        }
      }
    }
  }'

echo "Environment variables updated. App Runner will automatically redeploy."
echo "Check status with: aws apprunner describe-service --service-arn $SERVICE_ARN --region us-east-1"
