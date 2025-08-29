package services

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	configPkg "github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// emailService implements the EmailService interface
type emailService struct {
	sesService      interfaces.SESService
	templateService interfaces.TemplateService
	config          configPkg.SESConfig
	logger          *logrus.Logger
}

// NewEmailService creates a new email service instance
func NewEmailService(sesService interfaces.SESService, templateService interfaces.TemplateService, config configPkg.SESConfig, logger *logrus.Logger) interfaces.EmailService {
	return &emailService{
		sesService:      sesService,
		templateService: templateService,
		config:          config,
		logger:          logger,
	}
}

// SendReportEmail sends an internal email notification when a report is generated (internal only)
func (e *emailService) SendReportEmail(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) error {
	// Check if this is a high priority inquiry
	isHighPriority := e.detectHighPriority(inquiry.Message)

	var subject string
	if isHighPriority {
		subject = fmt.Sprintf("🚨 HIGH PRIORITY - New Cloud Consulting Report - %s", inquiry.Name)
	} else {
		subject = fmt.Sprintf("New Cloud Consulting Report Generated - %s", inquiry.Name)
	}

	// Use branded template if available, fallback to basic template
	var htmlBody string
	var textBody string

	if e.templateService != nil {
		// Use the template service to prepare data properly
		templateData := e.templateService.PrepareConsultantNotificationData(inquiry, report, isHighPriority)

		// Render branded HTML template
		brandedHTML, err := e.templateService.RenderEmailTemplate(ctx, "consultant_notification", templateData)
		if err != nil {
			e.logger.WithError(err).WithField("inquiry_id", inquiry.ID).Warn("Failed to render branded template, using fallback")
			htmlBody = e.buildReportEmailHTML(inquiry, report)
		} else {
			htmlBody = brandedHTML
		}
	} else {
		htmlBody = e.buildReportEmailHTML(inquiry, report)
	}

	// Always use the text version as fallback
	textBody = e.buildReportEmailText(inquiry, report)

	// Send only to internal address - no customer email included
	email := &interfaces.EmailMessage{
		From:     e.config.SenderEmail,
		To:       []string{"info@cloudpartner.pro"},
		Subject:  subject,
		TextBody: textBody,
		HTMLBody: htmlBody,
		ReplyTo:  e.config.ReplyToEmail,
	}

	err := e.sesService.SendEmail(ctx, email)
	if err != nil {
		e.logger.WithError(err).WithFields(logrus.Fields{
			"inquiry_id": inquiry.ID,
			"report_id":  report.ID,
		}).Error("Failed to send internal report email")
		return fmt.Errorf("failed to send internal report email: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"inquiry_id":    inquiry.ID,
		"report_id":     report.ID,
		"recipients":    email.To,
		"high_priority": isHighPriority,
		"template_used": "branded",
	}).Info("Internal report email sent successfully")

	return nil
}

// SendInquiryNotification sends an email notification when a new inquiry is received
func (e *emailService) SendInquiryNotification(ctx context.Context, inquiry *domain.Inquiry) error {
	// Check if this is a high priority inquiry
	isHighPriority := e.detectHighPriority(inquiry.Message)

	var subject string
	if isHighPriority {
		subject = fmt.Sprintf("🚨 HIGH PRIORITY - New Cloud Consulting Inquiry - %s", inquiry.Name)
	} else {
		subject = fmt.Sprintf("New Cloud Consulting Inquiry - %s", inquiry.Name)
	}

	textBody := e.buildInquiryEmailText(inquiry)
	htmlBody := e.buildInquiryEmailHTML(inquiry)

	email := &interfaces.EmailMessage{
		From:     e.config.SenderEmail,
		To:       []string{"info@cloudpartner.pro"},
		Subject:  subject,
		TextBody: textBody,
		HTMLBody: htmlBody,
		ReplyTo:  e.config.ReplyToEmail,
	}

	err := e.sesService.SendEmail(ctx, email)
	if err != nil {
		e.logger.WithError(err).WithField("inquiry_id", inquiry.ID).Error("Failed to send inquiry notification email")
		return fmt.Errorf("failed to send inquiry notification email: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"inquiry_id":    inquiry.ID,
		"recipients":    email.To,
		"high_priority": isHighPriority,
	}).Info("Inquiry notification email sent successfully")

	return nil
}

// SendCustomerConfirmation sends a confirmation email to the customer
// NOTE: This function NEVER includes AI-generated reports - those are for internal use only
func (e *emailService) SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error {
	// Validate and clean customer email
	customerEmail := e.validateAndCleanEmail(inquiry.Email)
	if customerEmail == "" {
		e.logger.WithField("inquiry_id", inquiry.ID).Warn("Invalid customer email, skipping confirmation")
		return nil // Don't fail the inquiry creation for invalid email
	}

	subject := "Thank you for your cloud consulting inquiry - CloudPartner Pro"

	// Use branded template if available, fallback to basic template
	var htmlBody string
	var textBody string

	if e.templateService != nil {
		// Use the template service to prepare data properly - NO REPORT INFORMATION FOR CUSTOMERS
		templateData := e.templateService.PrepareCustomerConfirmationData(inquiry)

		// Render branded HTML template
		brandedHTML, err := e.templateService.RenderEmailTemplate(ctx, "customer_confirmation", templateData)
		if err != nil {
			e.logger.WithError(err).WithField("inquiry_id", inquiry.ID).Warn("Failed to render branded template, using fallback")
			htmlBody = e.buildCustomerConfirmationHTML(inquiry)
		} else {
			htmlBody = brandedHTML
		}
	} else {
		htmlBody = e.buildCustomerConfirmationHTML(inquiry)
	}

	// Always use the text version as fallback - NO REPORTS
	textBody = e.buildCustomerConfirmationText(inquiry)

	email := &interfaces.EmailMessage{
		From:     e.config.SenderEmail,
		To:       []string{customerEmail},
		Subject:  subject,
		TextBody: textBody,
		HTMLBody: htmlBody,
		ReplyTo:  e.config.ReplyToEmail,
	}

	err := e.sesService.SendEmail(ctx, email)
	if err != nil {
		e.logger.WithError(err).WithFields(logrus.Fields{
			"inquiry_id":     inquiry.ID,
			"customer_email": customerEmail,
		}).Error("Failed to send customer confirmation email")
		return fmt.Errorf("failed to send customer confirmation email: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"inquiry_id":     inquiry.ID,
		"customer_email": customerEmail,
		"template_used":  "branded",
	}).Info("Customer confirmation email sent successfully")

	return nil
}

// IsHealthy checks if the email service is healthy
func (e *emailService) IsHealthy() bool {
	// Basic health check - verify we have required configuration
	return e.config.SenderEmail != "" && e.config.AccessKeyID != "" && e.config.SecretAccessKey != ""
}

// validateAndCleanEmail validates and cleans an email address
func (e *emailService) validateAndCleanEmail(email string) string {
	if email == "" {
		return ""
	}

	// Clean the email
	cleaned := strings.TrimSpace(strings.ToLower(email))

	// Check for placeholder emails or invalid formats
	placeholders := []string{
		"test@example.com",
		"user@example.com",
		"admin@example.com",
		"noreply@example.com",
		"example@example.com",
		"test@test.com",
		"user@test.com",
	}

	for _, placeholder := range placeholders {
		if cleaned == placeholder {
			return ""
		}
	}

	// Basic email validation (contains @ and .)
	if !strings.Contains(cleaned, "@") || !strings.Contains(cleaned, ".") {
		return ""
	}

	// Check for minimum length and basic structure
	parts := strings.Split(cleaned, "@")
	if len(parts) != 2 || len(parts[0]) < 1 || len(parts[1]) < 3 {
		return ""
	}

	return cleaned
}

// detectHighPriority analyzes the message content for urgency indicators
func (e *emailService) detectHighPriority(message string) bool {
	messageLower := strings.ToLower(message)

	// Time-sensitive keywords
	urgentKeywords := []string{
		"urgent", "asap", "immediately", "emergency", "critical",
		"today", "tomorrow", "this week", "right away", "quickly",
		"deadline", "time sensitive", "time-sensitive", "rush",
		"priority", "important", "blocking", "stuck", "help",
	}

	// Meeting request keywords that suggest urgency
	meetingKeywords := []string{
		"schedule today", "schedule tomorrow", "meet today", "meet tomorrow",
		"call today", "call tomorrow", "available today", "available tomorrow",
		"schedule asap", "meet asap", "call asap", "discuss today",
		"discuss tomorrow", "talk today", "talk tomorrow",
	}

	// Check for urgent keywords
	for _, keyword := range urgentKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}

	// Check for urgent meeting requests
	for _, keyword := range meetingKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}

	// Check for specific time patterns that suggest urgency
	timePatterns := []string{
		"within", "by end of", "before", "need by", "due",
		"this morning", "this afternoon", "this evening",
		"first thing", "end of day", "eod",
	}

	for _, pattern := range timePatterns {
		if strings.Contains(messageLower, pattern) {
			return true
		}
	}

	return false
}

// sanitizeMarkdownForPlainText cleans up markdown for plain text display
func (e *emailService) sanitizeMarkdownForPlainText(markdown string) string {
	// Remove markdown formatting for plain text
	text := markdown

	// Remove markdown headers
	text = regexp.MustCompile(`#{1,6}\s*`).ReplaceAllString(text, "")

	// Remove bold/italic markers
	text = regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`__([^_]+)__`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`_([^_]+)_`).ReplaceAllString(text, "$1")

	// Remove markdown links, keep just the text
	text = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(text, "$1")

	// Clean up extra whitespace
	text = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)

	return text
}

// buildReportEmailText creates the plain text version of the report email
func (e *emailService) buildReportEmailText(inquiry *domain.Inquiry, report *domain.Report) string {
	var builder strings.Builder

	isHighPriority := e.detectHighPriority(inquiry.Message)

	if isHighPriority {
		builder.WriteString("🚨 HIGH PRIORITY - NEW CLOUD CONSULTING REPORT GENERATED\n\n")
		builder.WriteString("⚠️  URGENT ATTENTION REQUIRED ⚠️\n")
		builder.WriteString("This inquiry contains urgent language or meeting requests.\n")
		builder.WriteString("Please prioritize response and review immediately.\n\n")
	} else {
		builder.WriteString("🔔 NEW CLOUD CONSULTING REPORT GENERATED\n\n")
	}

	builder.WriteString("👤 CLIENT INFORMATION\n")
	builder.WriteString("=====================\n")
	builder.WriteString(fmt.Sprintf("Name: %s\n", inquiry.Name))
	builder.WriteString(fmt.Sprintf("Email: %s\n", inquiry.Email))
	builder.WriteString(fmt.Sprintf("Company: %s\n", inquiry.Company))
	builder.WriteString(fmt.Sprintf("Phone: %s\n", inquiry.Phone))
	builder.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(inquiry.Services, ", ")))
	builder.WriteString(fmt.Sprintf("Inquiry ID: %s\n", inquiry.ID))
	builder.WriteString(fmt.Sprintf("Report ID: %s\n\n", report.ID))

	builder.WriteString("💬 ORIGINAL MESSAGE\n")
	builder.WriteString("===================\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", inquiry.Message))

	builder.WriteString("📋 GENERATED REPORT\n")
	builder.WriteString("===================\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", e.sanitizeMarkdownForPlainText(report.Content)))

	builder.WriteString("📌 ACTION REQUIRED\n")
	builder.WriteString("==================\n")
	if isHighPriority {
		builder.WriteString("🚨 HIGH PRIORITY: Please review and respond to the client IMMEDIATELY.\n")
		builder.WriteString("Check for meeting requests, urgent timelines, or critical business needs.\n\n")
	} else {
		builder.WriteString("Please review and respond to the client accordingly.\n\n")
	}

	builder.WriteString("---\n")
	builder.WriteString("This is an automated notification from Cloud Consulting Business.\n")
	builder.WriteString("Contact: info@cloudpartner.pro")

	return builder.String()
}

// buildReportEmailHTML creates the HTML version of the report email
func (e *emailService) buildReportEmailHTML(inquiry *domain.Inquiry, report *domain.Report) string {
	isHighPriority := e.detectHighPriority(inquiry.Message)

	var headerTitle, priorityAlert, actionText string
	var headerStyle string

	if isHighPriority {
		headerTitle = "🚨 HIGH PRIORITY - New Cloud Consulting Report Generated"
		headerStyle = "background: linear-gradient(135deg, #dc3545, #c82333);"
		priorityAlert = `<div style="background-color: #f8d7da; color: #721c24; padding: 20px; margin: 20px 0; border: 2px solid #f5c6cb; border-radius: 5px; text-align: center;">
                <h2 style="color: #721c24; margin: 0 0 10px 0;">⚠️ URGENT ATTENTION REQUIRED ⚠️</h2>
                <p style="margin: 0; font-weight: bold;">This inquiry contains urgent language or meeting requests. Please prioritize response and review immediately.</p>
            </div>`
		actionText = "🚨 <strong>HIGH PRIORITY:</strong> Please review and respond to the client IMMEDIATELY. Check for meeting requests, urgent timelines, or critical business needs."
	} else {
		headerTitle = "🔔 New Cloud Consulting Report Generated"
		headerStyle = "background: linear-gradient(135deg, #007cba, #005a8b);"
		priorityAlert = ""
		actionText = "<strong>📌 Action Required:</strong> Please review and respond to the client accordingly."
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Cloud Consulting Report Generated</title>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 0;
            background-color: #f4f4f4;
        }
        .container { max-width: 800px; margin: 0 auto; background-color: white; }
        .header { 
            %s
            color: white; 
            padding: 30px 20px; 
            text-align: center; 
        }
        .header h1 { margin: 0; font-size: 24px; }
        .content { padding: 30px; }
        .info-section { 
            background-color: #f8f9fa; 
            padding: 20px; 
            margin: 20px 0; 
            border-left: 5px solid #007cba; 
            border-radius: 5px;
        }
        .report-section { 
            background-color: #f0f8f0; 
            padding: 20px; 
            margin: 20px 0; 
            border-left: 5px solid #28a745; 
            border-radius: 5px;
        }
        .footer { 
            background-color: #6c757d; 
            color: white; 
            padding: 15px; 
            text-align: center; 
            font-size: 12px; 
        }
        .info-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
            margin-top: 15px;
        }
        .info-item {
            background-color: white;
            padding: 10px;
            border-radius: 3px;
            border: 1px solid #e9ecef;
        }
        .info-item strong { color: #007cba; }
        pre { 
            white-space: pre-wrap; 
            word-wrap: break-word; 
            background-color: white;
            padding: 15px;
            border-radius: 5px;
            border: 1px solid #ddd;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            line-height: 1.4;
        }
        .report-content {
            background-color: white;
            padding: 20px;
            border-radius: 5px;
            border: 1px solid #ddd;
            line-height: 1.6;
        }
        .report-content h1, .report-content h2, .report-content h3, .report-content h4 {
            color: #007cba;
            margin-top: 25px;
            margin-bottom: 15px;
        }
        .report-content h1 { font-size: 24px; border-bottom: 2px solid #007cba; padding-bottom: 10px; }
        .report-content h2 { font-size: 20px; border-bottom: 1px solid #e9ecef; padding-bottom: 8px; }
        .report-content h3 { font-size: 18px; }
        .report-content h4 { font-size: 16px; }
        .report-content p { margin: 12px 0; }
        .report-content ul, .report-content ol { margin: 12px 0; padding-left: 25px; }
        .report-content li { margin: 6px 0; }
        .report-content strong { color: #007cba; font-weight: 600; }
        .report-content em { font-style: italic; color: #666; }
        .report-content blockquote {
            border-left: 4px solid #007cba;
            margin: 15px 0;
            padding: 10px 20px;
            background-color: #f8f9fa;
            font-style: italic;
        }
        .report-content code {
            background-color: #f8f9fa;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            font-size: 13px;
        }
        .report-content hr {
            border: none;
            border-top: 2px solid #e9ecef;
            margin: 25px 0;
        }
        h2, h3 { color: #007cba; margin-top: 0; }
        .message-box {
            background-color: white;
            padding: 15px;
            border-radius: 5px;
            border: 1px solid #ddd;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
        
        <div class="content">
            %s
            
            <div class="info-section">
                <h2>👤 Client Information</h2>
                <div class="info-grid">
                    <div class="info-item"><strong>Name:</strong> %s</div>
                    <div class="info-item"><strong>Email:</strong> %s</div>
                    <div class="info-item"><strong>Company:</strong> %s</div>
                    <div class="info-item"><strong>Phone:</strong> %s</div>
                </div>
                <div style="margin-top: 15px;">
                    <div class="info-item"><strong>Services Requested:</strong> %s</div>
                    <div class="info-item" style="margin-top: 10px;"><strong>Inquiry ID:</strong> %s</div>
                    <div class="info-item" style="margin-top: 10px;"><strong>Report ID:</strong> %s</div>
                </div>
            </div>
            
            <div class="info-section">
                <h3>💬 Original Message</h3>
                <div class="message-box">%s</div>
            </div>
            
            <div class="report-section">
                <h3>📋 Generated Report</h3>
                <div class="report-content">%s</div>
            </div>
            
            <div style="background-color: %s; padding: 15px; border-radius: 5px; margin: 20px 0;">
                <p style="margin: 0;">%s</p>
            </div>
        </div>
        
        <div class="footer">
            This is an automated notification from the Cloud Consulting Business.
        </div>
    </div>
</body>
</html>`,
		headerStyle,
		headerTitle,
		priorityAlert,
		html.EscapeString(inquiry.Name),
		html.EscapeString(inquiry.Email),
		html.EscapeString(inquiry.Company),
		html.EscapeString(inquiry.Phone),
		html.EscapeString(strings.Join(inquiry.Services, ", ")),
		inquiry.ID,
		report.ID,
		html.EscapeString(inquiry.Message),
		html.EscapeString(report.Content),
		func() string {
			if isHighPriority {
				return "#f8d7da"
			}
			return "#e7f3ff"
		}(),
		actionText)
}

// buildInquiryEmailText creates the plain text version of the inquiry notification email
func (e *emailService) buildInquiryEmailText(inquiry *domain.Inquiry) string {
	var builder strings.Builder

	builder.WriteString("🔔 NEW CLOUD CONSULTING INQUIRY RECEIVED\n\n")

	builder.WriteString("👤 CLIENT INFORMATION\n")
	builder.WriteString("=====================\n")
	builder.WriteString(fmt.Sprintf("Name: %s\n", inquiry.Name))
	builder.WriteString(fmt.Sprintf("Email: %s\n", inquiry.Email))
	builder.WriteString(fmt.Sprintf("Company: %s\n", inquiry.Company))
	builder.WriteString(fmt.Sprintf("Phone: %s\n", inquiry.Phone))
	builder.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(inquiry.Services, ", ")))
	builder.WriteString(fmt.Sprintf("Inquiry ID: %s\n\n", inquiry.ID))

	builder.WriteString("💬 MESSAGE\n")
	builder.WriteString("===========\n")
	builder.WriteString(fmt.Sprintf("%s\n\n", inquiry.Message))

	builder.WriteString("📌 ACTION REQUIRED\n")
	builder.WriteString("==================\n")
	builder.WriteString("Please review this inquiry and respond accordingly.\n\n")

	builder.WriteString("---\n")
	builder.WriteString("This is an automated notification from Cloud Consulting Business.\n")
	builder.WriteString("Contact: info@cloudpartner.pro")

	return builder.String()
}

// buildInquiryEmailHTML creates the HTML version of the inquiry notification email
func (e *emailService) buildInquiryEmailHTML(inquiry *domain.Inquiry) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>New Cloud Consulting Inquiry</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .header { background-color: #007cba; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .info-section { background-color: #f8f9fa; padding: 15px; margin: 10px 0; border-left: 4px solid #007cba; }
        .footer { background-color: #6c757d; color: white; padding: 10px; text-align: center; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>New Cloud Consulting Inquiry</h1>
    </div>
    
    <div class="content">
        <div class="info-section">
            <h2>Client Information</h2>
            <p><strong>Name:</strong> %s</p>
            <p><strong>Email:</strong> %s</p>
            <p><strong>Company:</strong> %s</p>
            <p><strong>Phone:</strong> %s</p>
            <p><strong>Services Requested:</strong> %s</p>
            <p><strong>Inquiry ID:</strong> %s</p>
        </div>
        
        <div class="info-section">
            <h3>Message</h3>
            <p>%s</p>
        </div>
        
        <p>Please review this inquiry and respond accordingly.</p>
    </div>
    
    <div class="footer">
        This is an automated notification from the Cloud Consulting Business.
    </div>
</body>
</html>`,
		html.EscapeString(inquiry.Name),
		html.EscapeString(inquiry.Email),
		html.EscapeString(inquiry.Company),
		html.EscapeString(inquiry.Phone),
		html.EscapeString(strings.Join(inquiry.Services, ", ")),
		inquiry.ID,
		html.EscapeString(inquiry.Message))
}

// Template data structures for branded email templates

// CustomerConfirmationTemplateData represents the data structure for customer confirmation emails
// NOTE: This struct NEVER includes report information - reports are for internal use only
type CustomerConfirmationTemplateData struct {
	Name     string
	Company  string
	Services string
	ID       string
	// Report fields have been REMOVED - customers should never receive AI-generated reports
}

// buildCustomerConfirmationText creates the plain text version of the customer confirmation email
func (e *emailService) buildCustomerConfirmationText(inquiry *domain.Inquiry) string {
	var builder strings.Builder

	builder.WriteString("Thank you for your cloud consulting inquiry!\r\n\r\n")

	builder.WriteString(fmt.Sprintf("Dear %s,\r\n\r\n", inquiry.Name))

	builder.WriteString("We have received your inquiry for cloud consulting services and wanted to confirm that it has been successfully submitted.\r\n\r\n")

	builder.WriteString("Inquiry Details\r\n")
	builder.WriteString("===============\r\n")
	builder.WriteString(fmt.Sprintf("Services Requested: %s\r\n", strings.Join(inquiry.Services, ", ")))
	builder.WriteString(fmt.Sprintf("Company: %s\r\n", inquiry.Company))
	builder.WriteString(fmt.Sprintf("Reference ID: %s\r\n\r\n", inquiry.ID))

	builder.WriteString("What happens next?\r\n")
	builder.WriteString("==================\r\n")
	builder.WriteString("• Our team will review your inquiry within 24 hours\r\n")
	builder.WriteString("• We'll prepare a customized assessment based on your requirements\r\n")
	builder.WriteString("• A cloud consultant will reach out to discuss your project in detail\r\n")
	builder.WriteString("• We'll provide you with a detailed proposal and timeline\r\n\r\n")

	builder.WriteString("If you have any immediate questions or need to provide additional information, please don't hesitate to contact us.\r\n\r\n")

	builder.WriteString("Best regards,\r\n")
	builder.WriteString("Cloud Consulting Team\r\n")
	builder.WriteString("info@cloudpartner.pro\r\n\r\n")

	builder.WriteString("---\r\n")
	builder.WriteString("This is an automated confirmation email. Please do not reply to this message.")

	return builder.String()
}

// NOTE: buildCustomerConfirmationTextWithReport function has been REMOVED
// AI-generated reports should NEVER be sent to customers - they are for internal use only

// buildCustomerConfirmationHTML creates the HTML version of the customer confirmation email
func (e *emailService) buildCustomerConfirmationHTML(inquiry *domain.Inquiry) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Thank You for Your Cloud Consulting Inquiry</title>
    <style>
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 0;
            background-color: #f4f4f4;
        }
        .container { max-width: 600px; margin: 0 auto; background-color: white; }
        .header { 
            background: linear-gradient(135deg, #007cba, #005a8b);
            color: white; 
            padding: 30px 20px; 
            text-align: center; 
        }
        .header h1 { margin: 0; font-size: 24px; }
        .content { padding: 30px; }
        .success-banner {
            background-color: #28a745;
            color: white;
            padding: 15px;
            text-align: center;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .inquiry-details { 
            background-color: #f8f9fa; 
            padding: 20px; 
            margin: 20px 0; 
            border-left: 5px solid #007cba; 
            border-radius: 5px;
        }
        .next-steps { 
            background-color: #e7f3ff; 
            padding: 20px; 
            margin: 20px 0; 
            border-left: 5px solid #28a745; 
            border-radius: 5px;
        }
        .steps-list {
            list-style: none;
            padding: 0;
        }
        .steps-list li {
            margin-bottom: 10px;
            padding-left: 25px;
            position: relative;
        }
        .steps-list li:before {
            content: "•";
            color: #28a745;
            font-weight: bold;
            position: absolute;
            left: 0;
        }
        .footer { 
            background-color: #6c757d; 
            color: white; 
            padding: 15px; 
            text-align: center; 
            font-size: 12px; 
        }
        h2, h3 { color: #007cba; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>CloudPartner Pro</h1>
            <p>Your Trusted Cloud Consulting Partner</p>
        </div>
        
        <div class="success-banner">
            ✅ Inquiry Successfully Received!
        </div>
        
        <div class="content">
            <h2>Dear %s,</h2>
            
            <p>Thank you for reaching out to CloudPartner Pro! We have successfully received your cloud consulting inquiry and are excited to help you achieve your cloud transformation goals.</p>
            
            <div class="inquiry-details">
                <h3>Your Inquiry Details</h3>
                <p><strong>Services Requested:</strong> %s</p>
                <p><strong>Company:</strong> %s</p>
                <p><strong>Reference ID:</strong> %s</p>
            </div>
            
            <div class="next-steps">
                <h3>What Happens Next?</h3>
                <ul class="steps-list">
                    <li>Our expert team will review your inquiry within <strong>24 hours</strong></li>
                    <li>We'll prepare a <strong>customized assessment</strong> based on your specific requirements</li>
                    <li>A senior cloud consultant will <strong>reach out personally</strong> to discuss your project</li>
                    <li>You'll receive a <strong>detailed proposal</strong> with timeline and next steps</li>
                </ul>
            </div>
            
            <p>If you have any immediate questions or additional information to share, please don't hesitate to contact us at <a href="mailto:info@cloudpartner.pro">info@cloudpartner.pro</a>.</p>
            
            <p>Best regards,<br>
            The CloudPartner Pro Team<br>
            <a href="mailto:info@cloudpartner.pro">info@cloudpartner.pro</a></p>
        </div>
        
        <div class="footer">
            <p>This is an automated confirmation email. Please do not reply directly to this message.<br>
            For support, please contact us at info@cloudpartner.pro</p>
        </div>
    </div>
</body>
</html>`,
		html.EscapeString(inquiry.Name),
		html.EscapeString(strings.Join(inquiry.Services, ", ")),
		html.EscapeString(inquiry.Company),
		inquiry.ID)
}

// NOTE: buildCustomerConfirmationHTMLWithReport function has been REMOVED
// AI-generated reports should NEVER be sent to customers - they are for internal use only

// NOTE: SendCustomerConfirmationWithPDF function has been REMOVED
// AI-generated reports should NEVER be sent to customers - they are for internal use only
// Customers only receive simple confirmation emails via SendCustomerConfirmation

// Additional helper methods for the email service
