# AWS App Runner Deployment Guide

## Overview

This guide will help you deploy your Cloud Consulting Platform to AWS App Runner and configure it with your custom domain. Your project is already configured for App Runner deployment with Docker.

## Prerequisites

1. **AWS CLI installed and configured**

   ```bash
   aws configure
   ```

2. **GitHub repository** (I see you're using: `https://github.com/Tetianamost/cloud-consulting-business`)

3. **Custom domain** (e.g., `cloudpartner.pro`)

4. **AWS permissions** for App Runner, RDS, Route 53, and Certificate Manager

## Step 1: Deploy RDS Database

First, create your PostgreSQL database:

```bash
# Make the script executable
chmod +x deploy-apprunner-rds.sh

# Run the deployment script
./deploy-apprunner-rds.sh
```

This will:

- Create a PostgreSQL RDS instance (`db.t3.micro`)
- Set up the database with proper security groups
- Output the database endpoint for App Runner

## Step 2: Deploy App Runner Service

### Option A: Quick Deploy (Recommended)

```bash
# Make the script executable
chmod +x deploy-now.sh

# Deploy the App Runner service
./deploy-now.sh
```

### Option B: Manual Deploy

```bash
# Make the script executable
chmod +x create-apprunner-service.sh

# Update the GitHub URL in the script if needed
# Then run:
./create-apprunner-service.sh
```

## Step 3: Configure Environment Variables

After deployment, add environment variables to your App Runner service:

```bash
# Get your service ARN
SERVICE_ARN=$(aws apprunner list-services \
    --region us-east-1 \
    --query 'ServiceSummaryList[?ServiceName==`cloud-consulting-prod`].ServiceArn' \
    --output text)

# Update service with environment variables
aws apprunner update-service \
    --service-arn "$SERVICE_ARN" \
    --source-configuration '{
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
    }' \
    --region us-east-1
```

## Step 4: Set Up Custom Domain

### 4.1 Create SSL Certificate

```bash
# Request SSL certificate for your domain
aws acm request-certificate \
    --domain-name cloudpartner.pro \
    --subject-alternative-names "*.cloudpartner.pro" \
    --validation-method DNS \
    --region us-east-1
```

### 4.2 Validate Certificate

1. Go to AWS Certificate Manager console
2. Find your certificate
3. Add the DNS validation records to your domain's DNS settings

### 4.3 Associate Custom Domain with App Runner

```bash
# Get your service ARN
SERVICE_ARN=$(aws apprunner list-services \
    --region us-east-1 \
    --query 'ServiceSummaryList[?ServiceName==`cloud-consulting-prod`].ServiceArn' \
    --output text)

# Get your certificate ARN
CERT_ARN=$(aws acm list-certificates \
    --region us-east-1 \
    --query 'CertificateSummaryList[?DomainName==`cloudpartner.pro`].CertificateArn' \
    --output text)

# Associate custom domain
aws apprunner associate-custom-domain \
    --service-arn "$SERVICE_ARN" \
    --domain-name "cloudpartner.pro" \
    --enable-www-subdomain \
    --region us-east-1
```

### 4.4 Update DNS Records

After associating the domain, App Runner will provide DNS validation records. Add these to your domain's DNS:

```bash
# Get domain association details
aws apprunner describe-custom-domains \
    --service-arn "$SERVICE_ARN" \
    --region us-east-1
```

This will show you the CNAME records to add to your DNS.

## Step 5: Environment Variables Configuration

Create a script to update environment variables:

```bash
#!/bin/bash
# update-env-vars.sh

SERVICE_ARN=$(aws apprunner list-services \
    --region us-east-1 \
    --query 'ServiceSummaryList[?ServiceName==`cloud-consulting-prod`].ServiceArn' \
    --output text)

RDS_ENDPOINT=$(aws rds describe-db-instances \
    --db-instance-identifier consulting-prod \
    --region us-east-1 \
    --query 'DBInstances[0].Endpoint.Address' \
    --output text)

# Update service configuration
aws apprunner update-service \
    --service-arn "$SERVICE_ARN" \
    --source-configuration '{
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
                        "DATABASE_URL": "postgresql://postgres:CloudConsulting2024!@'$RDS_ENDPOINT':5432/postgres",
                        "GIN_MODE": "release",
                        "PORT": "80",
                        "AWS_REGION": "us-east-1",
                        "BEDROCK_REGION": "us-east-1",
                        "BEDROCK_MODEL_ID": "amazon.nova-lite-v1:0",
                        "AWS_SES_REGION": "us-east-1",
                        "SES_SENDER_EMAIL": "info@cloudpartner.pro",
                        "CORS_ALLOWED_ORIGINS": "https://cloudpartner.pro,https://www.cloudpartner.pro",
                        "CHAT_MODE": "polling",
                        "ENABLE_EMAIL_EVENTS": "true",
                        "JWT_SECRET": "your-secure-jwt-secret-change-this"
                    }
                }
            }
        }
    }' \
    --region us-east-1
```

## Step 6: Verify Deployment

### Check Service Status

```bash
# Get service details
aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region us-east-1
```

### Test Endpoints

```bash
# Get your App Runner URL
SERVICE_URL=$(aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region us-east-1 \
    --query 'Service.ServiceUrl' \
    --output text)

# Test health endpoint
curl https://$SERVICE_URL/health

# Test API endpoint
curl https://$SERVICE_URL/api/health
```

## Step 7: DNS Configuration for Custom Domain

### If using Route 53:

```bash
# Create hosted zone (if not exists)
aws route53 create-hosted-zone \
    --name cloudpartner.pro \
    --caller-reference $(date +%s)

# Get hosted zone ID
HOSTED_ZONE_ID=$(aws route53 list-hosted-zones \
    --query 'HostedZones[?Name==`cloudpartner.pro.`].Id' \
    --output text | cut -d'/' -f3)

# The CNAME records will be provided by App Runner after domain association
```

### If using external DNS provider:

1. Get the CNAME records from App Runner console
2. Add them to your DNS provider (GoDaddy, Namecheap, etc.)
3. Wait for DNS propagation (can take up to 48 hours)

## Step 8: Monitor Deployment

### View Logs

```bash
# View App Runner logs
aws logs describe-log-groups \
    --log-group-name-prefix "/aws/apprunner/cloud-consulting-prod" \
    --region us-east-1
```

### Check Service Health

```bash
# Monitor service status
watch -n 30 'aws apprunner describe-service \
    --service-arn "$SERVICE_ARN" \
    --region us-east-1 \
    --query "Service.Status"'
```

## Troubleshooting

### Common Issues:

1. **Build Failures**

   - Check your `apprunner.yaml` configuration
   - Verify Dockerfile builds locally
   - Check App Runner build logs

2. **Database Connection Issues**

   - Verify RDS security groups allow App Runner access
   - Check database credentials in environment variables
   - Ensure RDS is publicly accessible

3. **Custom Domain Issues**

   - Verify SSL certificate is validated
   - Check DNS propagation with `dig cloudpartner.pro`
   - Ensure CNAME records are correctly configured

4. **Environment Variables**
   - Use AWS Console to verify environment variables are set
   - Check for typos in variable names
   - Ensure sensitive values are properly escaped

## Cost Optimization

### App Runner Pricing:

- **Provisioned capacity**: $0.007/hour for 0.25 vCPU, 0.5 GB RAM
- **Request charges**: $0.40 per million requests
- **Data transfer**: Standard AWS rates

### RDS Pricing:

- **db.t3.micro**: ~$13/month
- **Storage**: $0.115/GB/month
- **Backup**: Free for 7 days

## Security Best Practices

1. **Use IAM roles** instead of access keys when possible
2. **Enable RDS encryption** at rest
3. **Use VPC** for database security
4. **Rotate secrets** regularly
5. **Monitor with CloudWatch** and set up alerts

## Next Steps

1. **Set up monitoring** with CloudWatch dashboards
2. **Configure auto-scaling** based on traffic
3. **Set up CI/CD pipeline** for automated deployments
4. **Implement backup strategy** for RDS
5. **Set up domain monitoring** for SSL certificate expiration

## Useful Commands

```bash
# Get service URL
aws apprunner describe-service --service-arn "$SERVICE_ARN" --query 'Service.ServiceUrl' --output text

# Update service
aws apprunner update-service --service-arn "$SERVICE_ARN" --source-configuration file://service-config.json

# Delete service (if needed)
aws apprunner delete-service --service-arn "$SERVICE_ARN"

# List all services
aws apprunner list-services --region us-east-1
```

Your Cloud Consulting Platform should now be accessible at both:

- App Runner URL: `https://[random-id].us-east-1.awsapprunner.com`
- Custom domain: `https://cloudpartner.pro` (after DNS propagation)
