# Email System Tests

This directory contains email system integration tests, SES connectivity tests, and template rendering tests.

## Test Files

Email tests will be moved here from the backend root directory in subsequent tasks:

- SES connectivity tests (`test_ses_*.go`)
- Email integration tests (`test_email_*.go`)
- Template rendering tests (`test_template_*.go`)
- Email delivery verification tests
- MIME formatting tests

## Running Email Tests

```bash
# Run specific email test
cd backend && go run testing/email/test_name.go

# Test SES connectivity
cd backend && go run testing/email/test_ses_connectivity.go

# Test email templates
cd backend && go run testing/email/test_template_rendering.go
```

## Requirements

- AWS SES credentials configured
- Verified sender email addresses
- Test email addresses for delivery testing
- Email templates in `backend/templates/email/`

## Environment Variables

- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key  
- `AWS_SES_REGION`: SES region (default: us-east-1)
- `SES_SENDER_EMAIL`: Verified sender email
- `TEST_RECIPIENT_EMAIL`: Test recipient email