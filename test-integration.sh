#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Testing Cloud Consulting Backend with Bedrock Integration${NC}"
echo -e "${YELLOW}=====================================================${NC}"

# Check if backend is running
echo -e "\n${GREEN}Checking if backend is running...${NC}"
if curl -s http://localhost:8061/health > /dev/null; then
  echo -e "${GREEN}✓ Backend is running${NC}"
else
  echo -e "${RED}✗ Backend is not running. Please start it with: cd backend && go run ./cmd/server/main.go${NC}"
  exit 1
fi

# Check if frontend is running
echo -e "\n${GREEN}Checking if frontend is running...${NC}"
if curl -s http://localhost:3001 > /dev/null; then
  echo -e "${GREEN}✓ Frontend is running${NC}"
else
  echo -e "${YELLOW}⚠ Frontend might not be running on port 3001. Please start it with: cd frontend && npm start${NC}"
fi

# Test backend API endpoints
echo -e "\n${GREEN}Testing backend API endpoints...${NC}"

# Test health endpoint
echo -e "\n${YELLOW}Testing health endpoint:${NC}"
curl -s http://localhost:8061/health | jq .

# Test services endpoint
echo -e "\n${YELLOW}Testing services endpoint:${NC}"
curl -s http://localhost:8061/api/v1/config/services | jq .

# Create a test inquiry
echo -e "\n${YELLOW}Creating a test inquiry with Bedrock integration:${NC}"
INQUIRY_RESPONSE=$(curl -s -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "company": "Test Company",
    "services": ["assessment"],
    "message": "This is a test inquiry to verify Bedrock integration."
  }')

echo $INQUIRY_RESPONSE | jq .

# Extract inquiry ID
INQUIRY_ID=$(echo $INQUIRY_RESPONSE | jq -r '.data.id')

if [ "$INQUIRY_ID" != "null" ]; then
  echo -e "\n${GREEN}✓ Inquiry created with ID: $INQUIRY_ID${NC}"
  
  # Wait for report generation
  echo -e "\n${YELLOW}Waiting 2 seconds for report generation...${NC}"
  sleep 2
  
  # Get the inquiry with report
  echo -e "\n${YELLOW}Getting inquiry details:${NC}"
  curl -s http://localhost:8061/api/v1/inquiries/$INQUIRY_ID | jq .
  
  # Get the report
  echo -e "\n${YELLOW}Getting generated report:${NC}"
  curl -s http://localhost:8061/api/v1/inquiries/$INQUIRY_ID/report | jq .
else
  echo -e "\n${RED}✗ Failed to create inquiry${NC}"
fi

echo -e "\n${GREEN}Integration test complete!${NC}"
echo -e "${YELLOW}=====================================================${NC}"
echo -e "${GREEN}To test manually:${NC}"
echo -e "1. Open ${YELLOW}http://localhost:3001${NC} in your browser"
echo -e "2. Fill out the contact form"
echo -e "3. Submit the form"
echo -e "4. Check backend logs for Bedrock API calls"
echo -e "${YELLOW}=====================================================${NC}"