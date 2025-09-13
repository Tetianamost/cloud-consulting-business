#!/bin/bash

echo "🔍 Checking RDS instance status..."

# Check if AWS CLI is available
if ! command -v aws &> /dev/null; then
    echo "❌ AWS CLI not found. Please install it first."
    exit 1
fi

DB_IDENTIFIER="consulting-prod"
REGION="us-east-1"

# Get the current status
STATUS=$(aws rds describe-db-instances \
    --db-instance-identifier $DB_IDENTIFIER \
    --region $REGION \
    --query 'DBInstances[0].DBInstanceStatus' \
    --output text 2>/dev/null)

if [ $? -ne 0 ]; then
    echo "❌ Failed to get RDS status. Check your AWS credentials and region."
    exit 1
fi

echo "📊 Current Status: $STATUS"

case $STATUS in
    "creating")
        echo "⏳ Database is still being created. Please wait..."
        echo "💡 This usually takes 10-20 minutes for a t3.micro instance."
        ;;
    "available")
        echo "✅ Database is ready!"
        
        # Get the endpoint
        ENDPOINT=$(aws rds describe-db-instances \
            --db-instance-identifier $DB_IDENTIFIER \
            --region $REGION \
            --query 'DBInstances[0].Endpoint.Address' \
            --output text)
        
        echo "🔗 Endpoint: $ENDPOINT"
        echo ""
        echo "📝 Add this to your .env file:"
        echo "DATABASE_URL=postgres://postgres:CloudConsulting2025@$ENDPOINT:5432/postgres?sslmode=require"
        echo "ENABLE_EMAIL_EVENTS=true"
        ;;
    "backing-up")
        echo "✅ Database is available (currently backing up)"
        
        # Get the endpoint
        ENDPOINT=$(aws rds describe-db-instances \
            --db-instance-identifier $DB_IDENTIFIER \
            --region $REGION \
            --query 'DBInstances[0].Endpoint.Address' \
            --output text)
        
        echo "🔗 Endpoint: $ENDPOINT"
        echo "💡 You can use it now, backup will continue in background"
        ;;
    "modifying")
        echo "⚠️  Database is available but being modified"
        ;;
    *)
        echo "⚠️  Unexpected status: $STATUS"
        ;;
esac

echo ""
echo "🔄 Run this script again to check status: ./check_rds_status.sh"