# Email Integration Test Guide

This guide helps you test the new dual email notification system.

## Prerequisites

1. AWS SES configured with verified sender email
2. Backend server running with proper environment variables
3. Valid customer email for testing

## Test Scenarios

### Scenario 1: Valid Customer Email

**Test**: Create inquiry with valid customer email

```bash
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "company": "Test Company",
    "phone": "555-123-4567",
    "services": ["assessment"],
    "message": "I need help with cloud migration assessment."
  }'
```

**Expected Results**:
1. Inquiry created successfully
2. Customer confirmation email sent to john.doe@example.com
3. Internal inquiry notification sent to info@cloudpartner.pro
4. AI report generated (if Bedrock configured)
5. Internal report email sent to info@cloudpartner.pro (if report generated)

### Scenario 2: Invalid/Placeholder Email

**Test**: Create inquiry with placeholder email

```bash
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "company": "Test Company",
    "phone": "555-123-4567",
    "services": ["assessment"],
    "message": "Test inquiry with placeholder email."
  }'
```

**Expected Results**:
1. Inquiry created successfully
2. NO customer confirmation email sent (placeholder email filtered out)
3. Internal inquiry notification sent to info@cloudpartner.pro
4. AI report generated (if Bedrock configured)
5. Internal report email sent to info@cloudpartner.pro (if report generated)

### Scenario 3: Empty Email

**Test**: Create inquiry with empty email

```bash
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anonymous User",
    "email": "",
    "company": "Test Company",
    "phone": "555-123-4567",
    "services": ["assessment"],
    "message": "Test inquiry with no email."
  }'
```

**Expected Results**:
1. Inquiry created successfully
2. NO customer confirmation email sent (empty email)
3. Internal inquiry notification sent to info@cloudpartner.pro
4. AI report generated (if Bedrock configured)
5. Internal report email sent to info@cloudpartner.pro (if report generated)

## Verification

### Check Server Logs

Monitor server logs for email-related messages:

```bash
# Look for these log patterns:
# - "Customer confirmation email sent successfully"
# - "Internal report email sent successfully"
# - "Inquiry notification email sent successfully"
# - "Warning: Invalid customer email, skipping confirmation"
# - "Warning: Failed to send [email type] for inquiry [id]"
```

### Check Email Delivery

1. **Customer Confirmation Email** should contain:
   - Professional thank you message
   - Inquiry reference ID
   - Next steps timeline
   - Contact information
   - NO AI report content

2. **Internal Emails** should contain:
   - Complete customer information
   - AI-generated report (in report email)
   - Original customer message
   - Professional formatting

### Email Content Validation

**Customer Confirmation Email**:
- Subject: "Thank you for your cloud consulting inquiry"
- Professional, customer-facing tone
- Clear next steps
- Reference ID for tracking

**Internal Report Email**:
- Subject: "New Cloud Consulting Report Generated - [Customer Name]"
- Complete customer details
- AI-generated report content
- Action required notice

## Troubleshooting

### Common Issues

1. **No emails sent**: Check AWS SES configuration and credentials
2. **Customer email not sent**: Verify email validation logic
3. **Internal emails not sent**: Check sender email verification in SES
4. **Malformed emails**: Check HTML/text template formatting

### Debug Steps

1. Check environment variables are set correctly
2. Verify AWS SES sender email is verified
3. Check server logs for specific error messages
4. Test with different email formats
5. Verify SES sending quota and limits

## Expected Log Output

```
INFO[...] Customer confirmation email sent successfully inquiry_id=... customer_email=...
INFO[...] Inquiry notification email sent successfully inquiry_id=... recipients=[info@cloudpartner.pro]
INFO[...] Internal report email sent successfully inquiry_id=... report_id=... recipients=[info@cloudpartner.pro]
```

Or for invalid emails:
```
WARN[...] Invalid customer email, skipping confirmation inquiry_id=...
```