# Fixed App Runner Deployment Guide

## Issue Resolution

The original error occurred because you were trying to use `"Runtime": "DOCKER"` in the App Runner API configuration, but App Runner expects specific runtime values like `NODEJS_18`, `GO_1`, etc. when using `ConfigurationSource: "API"`.

## Solution

The deployment has been fixed to use `ConfigurationSource: "REPOSITORY"`, which reads the `apprunner.yaml` file for Docker configuration.

## Deployment Steps

### Step 1: Set up IAM Role (Optional but Recommended)

```bash
./setup-apprunner-role.sh
```

This creates the necessary IAM role for App Runner to access AWS services (Bedrock, SES).

### Step 2: Deploy App Runner Service

```bash
./deploy-apprunner-no-yaml.sh
```

This script will:
- Check for existing services and offer to delete/update
- Verify RDS endpoint exists
- Use existing GitHub connection
- Create App Runner service using `apprunner.yaml` configuration
- Wait for deployment to complete

### Step 3: Update Environment Variables

```bash
./update-apprunner-env.sh
```

This script will:
- Get the real RDS endpoint
- Generate a secure JWT secret
- Update `apprunner.yaml` with real values
- Trigger a service update to apply changes

## Key Files

### `apprunner.yaml`
- Docker runtime configuration
- Environment variables
- Build and run commands
- Port configuration

### `deploy-apprunner-no-yaml.sh`
- Main deployment script
- Handles service creation
- Manages existing services
- Uses JSON configuration for AWS CLI

### `update-apprunner-env.sh`
- Updates environment variables
- Replaces placeholder values
- Triggers service redeployment

### `setup-apprunner-role.sh`
- Creates IAM role for AWS service access
- Sets up policies for Bedrock and SES
- Provides role ARN for deployment

## Environment Variables Set

The following environment variables are configured:

- `PORT`: 80
- `GIN_MODE`: release
- `CHAT_MODE`: polling
- `AWS_REGION`: us-east-1
- `BEDROCK_REGION`: us-east-1
- `BEDROCK_MODEL_ID`: amazon.nova-lite-v1:0
- `AWS_SES_REGION`: us-east-1
- `SES_SENDER_EMAIL`: info@cloudpartner.pro
- `CORS_ALLOWED_ORIGINS`: https://cloudpartner.pro,https://www.cloudpartner.pro
- `DATABASE_URL`: (Real RDS endpoint)
- `JWT_SECRET`: (Generated secure secret)

## Testing

After deployment, test your service:

```bash
# Get service URL from deployment output, then:
curl https://YOUR_SERVICE_URL/health
```

## Troubleshooting

### Service Already Exists
The deployment script will detect existing services and offer options to delete or update.

### GitHub Connection Issues
Ensure you have a GitHub connection set up in the AWS App Runner console.

### RDS Not Found
Make sure your RDS instance is deployed first:
```bash
./deploy-apprunner-rds.sh
```

### IAM Permission Issues
Run the IAM setup script:
```bash
./setup-apprunner-role.sh
```

## Architecture

```
GitHub Repository (kiro-dev branch)
    ↓
App Runner Service
    ├── Docker Build (multi-stage)
    ├── Frontend (React) → Nginx
    ├── Backend (Go) → Port 8080
    └── Supervisor → Port 80
    ↓
RDS PostgreSQL Database
AWS Bedrock (Nova model)
AWS SES (Email service)
```

## Next Steps

1. **Custom Domain**: Use `./setup-custom-domain.sh`
2. **Monitoring**: Set up CloudWatch dashboards
3. **SSL Certificate**: Configure through AWS Certificate Manager
4. **Auto-scaling**: Configure in App Runner console

## Cost Optimization

- **Instance Size**: 0.25 vCPU, 0.5 GB RAM (minimal for testing)
- **Auto-scaling**: Scales to zero when not in use
- **Pay-per-use**: Only pay for active requests

## Security Features

- **IAM Roles**: Secure access to AWS services
- **VPC**: Can be configured for private networking
- **HTTPS**: Automatic SSL/TLS termination
- **Environment Variables**: Secure configuration management

This fixed deployment approach resolves the runtime validation error and provides a robust, scalable deployment solution.