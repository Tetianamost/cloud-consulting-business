# AWS Setup Guide for Bedrock Integration

## Overview

This guide helps you set up AWS credentials for the Amazon Bedrock integration in your cloud consulting backend.

## Required AWS Credentials

### 1. AWS Access Keys (Required)
- **AWS_ACCESS_KEY_ID**: Your AWS access key ID
- **AWS_SECRET_ACCESS_KEY**: Your AWS secret access key

### 2. Amazon Bedrock API Key (Required)
- **AWS_BEARER_TOKEN_BEDROCK**: Your Bedrock API key

### 3. AWS Session Token (Optional)
- **AWS_SESSION_TOKEN**: Only needed for temporary credentials
- **When to use**: If you're using AWS STS, AWS SSO, or assume-role
- **When to skip**: For regular IAM user credentials (most common)

## How to Get Your Credentials

### Step 1: Get AWS Access Keys

1. **Log in to AWS Console**
   - Go to https://console.aws.amazon.com
   - Sign in with your AWS account

2. **Navigate to IAM**
   - Search for "IAM" in the AWS Console
   - Click on "Identity and Access Management"

3. **Create or Select User**
   - Go to "Users" in the left sidebar
   - Either select existing user or create new one
   - Click "Create access key" under "Access keys" tab

4. **Download Credentials**
   - Choose "Application running outside AWS"
   - Download the CSV file or copy the keys
   - **Important**: Save these securely, you can't retrieve them later

### Step 2: Get Bedrock API Key

1. **Navigate to Amazon Bedrock**
   - In AWS Console, search for "Bedrock"
   - Click on "Amazon Bedrock"

2. **Generate API Key**
   - Go to "API keys" in the left navigation
   - Click "Generate long-term API keys"
   - Set expiration (e.g., 30 days for development)
   - Click "Generate" and copy the key immediately
   - **Important**: You cannot retrieve this key later

### Step 3: Set Up Permissions

Your AWS user needs these permissions:
- `AmazonBedrockFullAccess` (or custom policy with Bedrock permissions)

## Environment Configuration

### For Development (.env file)

```bash
# Required - AWS Access Keys
AWS_ACCESS_KEY_ID=AKIA...your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key_here
AWS_REGION=us-east-1

# Required - Bedrock API Key
AWS_BEARER_TOKEN_BEDROCK=your_bedrock_api_key_here

# Optional - Only for temporary credentials
# AWS_SESSION_TOKEN=your_session_token_here

# Bedrock Configuration (optional - has defaults)
BEDROCK_REGION=us-east-1
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0
BEDROCK_BASE_URL=https://bedrock-runtime.us-east-1.amazonaws.com
BEDROCK_TIMEOUT_SECONDS=30
```

### For Production

Use AWS Secrets Manager or environment variables in your deployment platform:

```bash
# In your deployment environment
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
export AWS_BEARER_TOKEN_BEDROCK="your_bedrock_api_key"
```

## Testing Your Setup

### 1. Test AWS Credentials

```bash
# Using AWS CLI (if installed)
aws sts get-caller-identity

# Should return your user information
```

### 2. Test Bedrock Access

```bash
# Start your backend
docker-compose up -d backend

# Create a test inquiry (should generate AI report)
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "services": ["assessment"],
    "message": "Test inquiry for Bedrock integration"
  }'

# Check logs for Bedrock API calls
docker-compose logs backend
```

## Common Issues and Solutions

### Issue: "Bedrock API key not configured"
**Solution**: Make sure `AWS_BEARER_TOKEN_BEDROCK` is set in your `.env` file

### Issue: "Access denied" errors
**Solution**: 
- Check your AWS user has Bedrock permissions
- Verify your access keys are correct
- Ensure your Bedrock API key is valid and not expired

### Issue: "Session token required"
**Solution**: 
- If you're using regular IAM user credentials, leave `AWS_SESSION_TOKEN` empty
- Only set it if you're using temporary credentials from STS/SSO

### Issue: "Region not supported"
**Solution**: 
- Bedrock is not available in all regions
- Use `us-east-1` or `us-west-2` for best availability
- Check AWS documentation for current region support

## Security Best Practices

### Development
- ✅ Use `.env` file (never commit to git)
- ✅ Set short expiration on Bedrock API keys
- ✅ Use least-privilege IAM policies
- ✅ Rotate keys regularly

### Production
- ✅ Use AWS Secrets Manager
- ✅ Use IAM roles when possible (ECS, Lambda, EC2)
- ✅ Enable CloudTrail logging
- ✅ Monitor API usage and costs
- ✅ Set up billing alerts

## Cost Considerations

- Bedrock charges per API call and tokens processed
- Monitor usage in AWS Cost Explorer
- Set up billing alerts for unexpected charges
- Consider caching reports to reduce API calls

## Support

If you're still having issues:

1. **Check AWS Documentation**: https://docs.aws.amazon.com/bedrock/
2. **AWS Support**: Use AWS Support if you have a support plan
3. **Community**: AWS re:Post community forums
4. **Application Logs**: Check `docker-compose logs backend` for detailed error messages

## Example Working Configuration

Here's a complete working `.env` example (with fake values):

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION=us-east-1

# Bedrock API Key
AWS_BEARER_TOKEN_BEDROCK=bedrock_api_key_example_12345

# Application Settings
PORT=8061
LOG_LEVEL=4
GIN_MODE=debug
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Frontend
REACT_APP_API_URL=http://localhost:8061
```

Remember to replace all example values with your actual credentials!