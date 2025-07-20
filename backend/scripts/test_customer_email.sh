#!/bin/bash

# Customer Confirmation Email Test Script
# Usage: ./test_customer_email.sh your-email@example.com

if [ -z "$1" ]; then
    echo "Usage: $0 <your-test-email>"
    echo "Example: $0 john@example.com"
    exit 1
fi

TEST_EMAIL="$1"
BASE_URL="http://localhost:8061"

echo "🧪 Testing Customer Confirmation Email Functionality"
echo "📧 Test email: $TEST_EMAIL"
echo "🌐 Backend URL: $BASE_URL"
echo ""

# Test 1: Basic Quote Request
echo "📋 Test 1: Basic Quote Request"
echo "Submitting quote request..."

RESPONSE1=$(curl -s -X POST "$BASE_URL/api/v1/inquiries" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Test Customer\",
    \"email\": \"$TEST_EMAIL\",
    \"company\": \"Test Company Inc\",
    \"phone\": \"555-123-4567\",
    \"services\": [\"assessment\"],
    \"message\": \"Quote Request Details: - Service: Initial Assessment - Complexity: Moderate - Servers/Applications: 5 - Base Fee: $750 - Total Estimate: $1,500 Additional Requirements: I need help with cloud migration assessment.\"
  }")

if echo "$RESPONSE1" | grep -q "id"; then
    INQUIRY_ID1=$(echo "$RESPONSE1" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "✅ Quote request submitted successfully"
    echo "📝 Inquiry ID: $INQUIRY_ID1"
    echo "📧 Customer confirmation email should be sent to: $TEST_EMAIL"
else
    echo "❌ Quote request failed"
    echo "Response: $RESPONSE1"
fi

echo ""
sleep 2

# Test 2: Contact Us Form
echo "📞 Test 2: Contact Us Form"
echo "Submitting contact inquiry..."

RESPONSE2=$(curl -s -X POST "$BASE_URL/api/v1/inquiries" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Jane Smith\",
    \"email\": \"$TEST_EMAIL\",
    \"company\": \"ABC Corp\",
    \"phone\": \"555-987-6543\",
    \"services\": [\"optimization\"],
    \"message\": \"I would like to discuss cloud optimization opportunities for our company. Can we schedule a meeting this week?\"
  }")

if echo "$RESPONSE2" | grep -q "id"; then
    INQUIRY_ID2=$(echo "$RESPONSE2" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "✅ Contact inquiry submitted successfully"
    echo "📝 Inquiry ID: $INQUIRY_ID2"
    echo "📧 Customer confirmation email should be sent to: $TEST_EMAIL"
else
    echo "❌ Contact inquiry failed"
    echo "Response: $RESPONSE2"
fi

echo ""
sleep 2

# Test 3: High Priority Request
echo "🚨 Test 3: High Priority Request"
echo "Submitting urgent inquiry..."

RESPONSE3=$(curl -s -X POST "$BASE_URL/api/v1/inquiries" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Urgent Client\",
    \"email\": \"$TEST_EMAIL\",
    \"company\": \"Emergency Corp\",
    \"phone\": \"555-911-1234\",
    \"services\": [\"migration\"],
    \"message\": \"URGENT: We need help immediately! Our current system is down and we need to migrate to cloud ASAP. Can we schedule a meeting today or tomorrow?\"
  }")

if echo "$RESPONSE3" | grep -q "id"; then
    INQUIRY_ID3=$(echo "$RESPONSE3" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    echo "✅ Urgent inquiry submitted successfully"
    echo "📝 Inquiry ID: $INQUIRY_ID3"
    echo "📧 Customer confirmation email should be sent to: $TEST_EMAIL"
    echo "🚨 Internal team should receive HIGH PRIORITY notification"
else
    echo "❌ Urgent inquiry failed"
    echo "Response: $RESPONSE3"
fi

echo ""
echo "🔍 Verification Steps:"
echo "1. Check your email inbox ($TEST_EMAIL) for confirmation emails"
echo "2. Check spam/junk folder if not in inbox"
echo "3. Verify emails have subject: 'Thank you for your cloud consulting inquiry'"
echo "4. Confirm emails contain inquiry details and next steps"
echo "5. Verify NO AI report content is included in customer emails"
echo ""
echo "📊 Expected Results:"
echo "• 3 customer confirmation emails sent to $TEST_EMAIL"
echo "• 3 internal notifications sent to info@cloudpartner.pro"
echo "• 1 high priority email with 🚨 in subject (internal only)"
echo ""
echo "🔧 If emails not received, check:"
echo "• AWS SES configuration and sender email verification"
echo "• Server logs for email sending errors"
echo "• Email validation (placeholder emails are filtered)"
echo "• SES sending quota and sandbox mode"