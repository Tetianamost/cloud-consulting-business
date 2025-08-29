---
inclusion: always
---

# Email System Implementation Guidelines

## Overview

This document provides comprehensive guidance for working with the email system in the Cloud Consulting Platform. It covers architecture, implementation patterns, testing approaches, and common pitfalls to avoid duplicating work.

## Architecture Overview

### Core Components

#### 1. **Email Service Layer** (`backend/internal/services/email.go`)
- **Purpose**: High-level email business logic
- **Interface**: `interfaces.EmailService`
- **Key Methods**:
  - `SendReportEmail()` - Internal consultant notifications with AI reports
  - `SendCustomerConfirmation()` - Customer acknowledgment (NO reports)
  - `SendInquiryNotification()` - New inquiry alerts
  - `IsHealthy()` - Service health check

#### 2. **SES Service Layer** (`backend/internal/services/ses.go`)
- **Purpose**: AWS SES integration and MIME email handling
- **Interface**: `interfaces.SESService`
- **Key Methods**:
  - `SendEmail()` - Raw email delivery via AWS SES
  - `VerifyEmailAddress()` - Email verification
  - `GetSendingQuota()` - SES quota management

#### 3. **Template Service** (`backend/internal/services/template.go`)
- **Purpose**: HTML email template rendering and management
- **Interface**: `interfaces.TemplateService`
- **Key Methods**:
  - `RenderEmailTemplate()` - Render branded HTML emails
  - `PrepareCustomerConfirmationData()` - Customer email data
  - `PrepareConsultantNotificationData()` - Internal email data

#### 4. **Email Templates** (`backend/templates/email/`)
- **Customer Confirmation**: `customer_confirmation.html`
- **Consultant Notification**: `consultant_notification.html`
- **Professional Design**: Responsive, branded, modern styling

## Email Types and Security

### 🔒 **CRITICAL SECURITY RULE**
**CUSTOMERS NEVER RECEIVE AI-GENERATED REPORTS**

#### Customer Emails (`SendCustomerConfirmation`)
- ✅ **Purpose**: Professional acknowledgment and next steps
- ✅ **Content**: Thank you, inquiry summary, what happens next
- ❌ **NEVER Include**: AI reports, internal analysis, sensitive data
- 📧 **Template**: `customer_confirmation.html`

#### Internal Emails (`SendReportEmail`)
- ✅ **Purpose**: Consultant notifications with full AI reports
- ✅ **Content**: Client info, AI analysis, priority flags, action items
- ✅ **Recipients**: Only `info@cloudpartner.pro`
- 📧 **Template**: `consultant_notification.html`

## Configuration

### Environment Variables
```bash
# AWS SES Configuration
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_SES_REGION=us-east-1
SES_SENDER_EMAIL=info@cloudpartner.pro
SES_REPLY_TO_EMAIL=info@cloudpartner.pro
SES_TIMEOUT_SECONDS=30
```

### SES Setup Requirements
1. **Verify Sender Email** in AWS SES Console
2. **Check Sandbox Mode** - limits to verified emails only
3. **Production Access** - request if sending to unverified emails
4. **DNS Records** - SPF, DKIM, DMARC for deliverability

## Implementation Patterns

### ✅ **Correct Usage**

#### Creating Email Service (RECOMMENDED)
```go
// Load configuration
cfg, err := config.Load()
if err != nil {
    return err
}

// Create logger
logger := logrus.New()

// Use the email service factory (recommended)
emailService, err := services.NewEmailServiceWithSES(cfg.SES, logger)
if err != nil {
    return err
}
```

#### Creating Email Service (Manual - Advanced)
```go
// Load configuration
cfg, err := config.Load()
if err != nil {
    return err
}

// Create logger
logger := logrus.New()

// Create SES service
sesService, err := services.NewSESService(cfg.SES, logger)
if err != nil {
    return err
}

// Create template service
templateService := services.NewTemplateService("templates", logger)

// Create email service
emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)
```

#### Sending Customer Confirmation
```go
// CORRECT - Only acknowledgment, no reports
err := emailService.SendCustomerConfirmation(ctx, inquiry)
if err != nil {
    logger.WithError(err).Error("Failed to send customer confirmation")
    // Don't fail the inquiry creation for email issues
}
```

#### Sending Internal Report
```go
// CORRECT - Full report to consultants only
err := emailService.SendReportEmail(ctx, inquiry, report)
if err != nil {
    logger.WithError(err).Error("Failed to send internal report")
    return err
}
```

### ❌ **Incorrect Patterns to Avoid**

#### DON'T: Send Reports to Customers
```go
// WRONG - Never do this!
err := emailService.SendReportEmail(ctx, inquiry, report)
// This goes to customers - SECURITY VIOLATION
```

#### DON'T: Create Multiple Email Services
```go
// WRONG - Creates duplicate services
emailService1 := services.NewEmailService(...)
emailService2 := services.NewEmailService(...) // Duplicate!
```

#### DON'T: Ignore Email Failures for Internal Emails
```go
// WRONG - Internal emails are critical
err := emailService.SendReportEmail(ctx, inquiry, report)
// Ignoring error - consultants won't get reports!
```

## Testing Patterns

### Mock Testing (Development)
```go
// Create mock SES service
mockSES := &MockSESService{
    sentEmails: make([]*interfaces.EmailMessage, 0),
}

// Create email service with mock
emailService := services.NewEmailService(mockSES, templateService, cfg.SES, logger)

// Test email sending
err := emailService.SendCustomerConfirmation(ctx, inquiry)
assert.NoError(t, err)

// Verify email content
assert.Len(t, mockSES.sentEmails, 1)
email := mockSES.sentEmails[0]
assert.Contains(t, email.HTMLBody, "CloudPartner Pro")
assert.NotContains(t, email.HTMLBody, "Generated Report") // Security check
```

### Real SES Testing (Staging/Production)
```go
// Test SES connectivity
sesService, err := services.NewSESService(cfg.SES, logger)
require.NoError(t, err)

quota, err := sesService.GetSendingQuota(ctx)
require.NoError(t, err)
assert.Greater(t, quota.Max24HourSend, float64(0))
```

### Template Testing
```go
// Test template rendering
templateService := services.NewTemplateService("templates", logger)
data := templateService.PrepareCustomerConfirmationData(inquiry)

html, err := templateService.RenderEmailTemplate(ctx, "customer_confirmation", data)
assert.NoError(t, err)
assert.Contains(t, html, inquiry.Name)
assert.Contains(t, html, "CloudPartner Pro")
```

## Common Issues and Solutions

### Issue: "SES Connection Failed"
**Symptoms**: `GetSendingQuota` fails, emails not sent
**Solutions**:
1. Verify AWS credentials in environment
2. Check SES region configuration
3. Ensure sender email is verified in SES Console
4. Check AWS account has SES access

### Issue: "Missing start boundary" Error
**Symptoms**: AWS SES returns `InvalidParameterValue: Missing start boundary` error
**Root Cause**: Incorrect MIME multipart message construction in `buildRawMessage`
**Solutions**:
1. **Use Email Factory**: Use `NewEmailServiceWithSES()` for proper implementation
2. **Verify MIME Structure**: Ensure proper boundary placement and multipart nesting
3. **Check SES Service**: The current implementation has resolved MIME boundary issues

### Issue: "Emails Not Delivered"
**Symptoms**: SES succeeds but emails don't reach inbox
**Solutions**:
1. Check spam/junk folders
2. Verify recipient email addresses
3. Check SES sandbox mode restrictions
4. Review DNS records (SPF, DKIM, DMARC)

### Issue: "Template Not Found"
**Symptoms**: `RenderEmailTemplate` fails
**Solutions**:
1. Verify template files exist in `backend/templates/email/`
2. Check template service initialization
3. Ensure correct template names in code

### Issue: "Broken Email Formatting"
**Symptoms**: Emails display poorly in email clients
**Solutions**:
1. Test in multiple email clients (Gmail, Outlook, Apple Mail)
2. Validate HTML structure and CSS
3. Check responsive design on mobile
4. Use email-safe CSS properties

### Issue: "Production Emails Unformatted"
**Symptoms**: Test emails look great, but production emails are plain text or unformatted
**Root Cause**: Production code using incorrect SES implementation
**Solutions**:
1. **Check server.go**: Ensure `NewEmailServiceWithSES()` is used in production
2. **Restart Application**: After fixing, restart the server to load new implementation
3. **Verify Logs**: Look for successful email service initialization in startup logs

## File Locations Reference

### Core Implementation
- **Email Service**: `backend/internal/services/email.go`
- **SES Service**: `backend/internal/services/ses.go`
- **Template Service**: `backend/internal/services/template.go`
- **Interfaces**: `backend/internal/interfaces/services.go`
- **Configuration**: `backend/internal/config/config.go`

### Templates
- **Customer Template**: `backend/templates/email/customer_confirmation.html`
- **Consultant Template**: `backend/templates/email/consultant_notification.html`

### Test Files
- **Mock Testing**: `backend/test_email_simple_verification.go`
- **SES Testing**: `backend/test_ses_connectivity.go`
- **Real Email Testing**: `backend/test_real_ses_email.go`
- **Fixed SES Testing**: `backend/test_ses_fixed.go`

### Implementation Files
- **SES Service**: `backend/internal/services/ses.go`
- **Email Factory**: `backend/internal/services/email_factory.go`

## Priority Detection System

The email system automatically detects high-priority inquiries based on:

### Urgent Keywords
- "urgent", "asap", "immediately", "emergency", "critical"
- "today", "tomorrow", "this week", "deadline"
- "meeting", "schedule", "call", "discuss"

### High Priority Email Features
- 🚨 **Subject**: Prefixed with "HIGH PRIORITY"
- 🎨 **Styling**: Red gradient header, pulsing animations
- 📧 **Content**: Emphasizes immediate action required
- 🔔 **Alerts**: Visual indicators for consultants

## Professional Email Standards

### Design Requirements
- ✅ **Responsive Design**: Works on desktop and mobile
- ✅ **Professional Branding**: CloudPartner Pro logo and colors
- ✅ **Modern Styling**: Gradients, shadows, clean typography
- ✅ **Accessibility**: Proper contrast, readable fonts
- ✅ **Email Client Compatibility**: Gmail, Outlook, Apple Mail

### Content Standards
- ✅ **Clear Subject Lines**: Descriptive and actionable
- ✅ **Professional Tone**: Friendly but business-appropriate
- ✅ **Contact Information**: Always include support email
- ✅ **Unsubscribe Info**: For automated emails
- ✅ **Company Footer**: Branding and legal information

## Monitoring and Observability

### Health Checks
```go
// Check email service health
if !emailService.IsHealthy() {
    logger.Error("Email service is unhealthy")
    // Alert monitoring system
}
```

### Logging Standards
```go
// Log email events with context
logger.WithFields(logrus.Fields{
    "inquiry_id": inquiry.ID,
    "email_type": "customer_confirmation",
    "recipient":  inquiry.Email,
}).Info("Customer confirmation email sent")
```

### Metrics to Track
- Email delivery success rate
- Template rendering performance
- SES quota usage
- High priority email frequency
- Customer vs internal email volume

## Migration and Deployment

### Database Considerations
- No database schema changes required for email system
- Email content is generated dynamically
- Templates are file-based, not database-stored

### Deployment Checklist
1. ✅ Verify AWS SES credentials in production
2. ✅ Confirm sender email verification
3. ✅ Test email delivery to real addresses
4. ✅ Check template file deployment
5. ✅ Verify DNS records for deliverability
6. ✅ Monitor email delivery rates
7. ✅ Set up alerting for email failures

## Best Practices Summary

### DO ✅
- Use the existing email service interfaces
- Test with both mock and real SES
- Verify templates render correctly
- Check email client compatibility
- Monitor delivery rates and errors
- Follow security guidelines (no reports to customers)
- Use structured logging for email events

### DON'T ❌
- Create duplicate email services
- Send AI reports to customers
- Ignore email delivery failures
- Skip template testing
- Hardcode email addresses
- Use non-responsive email designs
- Forget to verify sender emails in SES

## Getting Help

### Debugging Email Issues
1. **Check Logs**: Look for structured log entries with email context
2. **Test Templates**: Use `test_email_simple_verification.go`
3. **Verify SES**: Use `test_ses_connectivity.go`
4. **Check Configuration**: Ensure all environment variables are set
5. **Review Templates**: Open generated HTML files in browser

### Common Commands
```bash
# Test email system (mock)
cd backend && go run test_email_simple_verification.go

# Test SES connectivity
cd backend && go run test_ses_connectivity.go

# Test real email delivery (with confirmation)
cd backend && go run test_real_ses_email.go

# Test SES implementation
cd backend && go run test_ses_fixed.go

# Test production email system
cd backend && go run test_production_email_fix.go
```

This comprehensive guide should prevent duplicate implementations and ensure consistent, secure email handling across the platform.