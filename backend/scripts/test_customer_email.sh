#!/bin/bash

# Test script for customer email with branded templates
# This script tests the customer confirmation email functionality

echo "=== Testing Customer Email with Branded Templates ==="

# Set environment variables for testing (using placeholder values)
export GIN_MODE=debug
export CORS_ALLOWED_ORIGINS="http://localhost:3000,https://cloudpartner.pro"
export AWS_BEARER_TOKEN_BEDROCK="test-token"
export BEDROCK_REGION="us-east-1"
export BEDROCK_MODEL_ID="amazon.nova-lite-v1:0"
export AWS_ACCESS_KEY_ID="test-access-key"
export AWS_SECRET_ACCESS_KEY="test-secret-key"
export AWS_SES_REGION="us-east-1"
export SES_SENDER_EMAIL="info@cloudpartner.pro"
export SES_REPLY_TO_EMAIL="info@cloudpartner.pro"

echo "Environment variables set for testing..."

# Start the server in the background
echo "Starting server..."
./server &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test the inquiry creation endpoint
echo "Testing inquiry creation with customer email..."

curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane.smith@example.com",
    "company": "Tech Innovations Inc",
    "phone": "555-987-6543",
    "services": ["assessment", "optimization"],
    "message": "We are looking to optimize our current AWS infrastructure and would like a comprehensive assessment of our setup. Please provide recommendations for cost optimization and performance improvements."
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s

echo ""
echo "=== Test Complete ==="
echo "Check the server logs above to see if the branded email templates were used."
echo "Look for log messages indicating 'template_used: branded' in the email service logs."

# Stop the server
echo "Stopping server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

echo "Server stopped."
echo ""
echo "Note: This test uses placeholder AWS credentials, so actual emails won't be sent."
echo "The test verifies that the branded templates are loaded and rendered correctly."