package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	configPkg "github.com/cloud-consulting/backend/internal/config"
)

// emailService implements the EmailService interface
type emailService struct {
	sesService interfaces.SESService
	config     configPkg.SESConfig
	logger     *logrus.Logger
}

// NewEmailService creates a new email service instance
func NewEmailService(sesService interfaces.SESService, config configPkg.SESConfig, logger *logrus.Logger) interfaces.EmailService {
	return &emailService{
		sesService: sesService,
		config:     config,
		logger:     logger,
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
	
	textBody := e.buildReportEmailText(inquiry, report)
	htmlBody := e.buildReportEmailHTML(inquiry, report)

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
		"inquiry_id":      inquiry.ID,
		"report_id":       report.ID,
		"recipients":      email.To,
		"high_priority":   isHighPriority,
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
func (e *emailService) SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error {
	// Validate and clean customer email
	customerEmail := e.validateAndCleanEmail(inquiry.Email)
	if customerEmail == "" {
		e.logger.WithField("inquiry_id", inquiry.ID).Warn("Invalid customer email, skipping confirmation")
		return nil // Don't fail the inquiry creation for invalid email
	}

	subject := "Thank you for your cloud consulting inquiry"
	
	textBody := e.buildCustomerConfirmationText(inquiry)
	htmlBody := e.buildCustomerConfirmationHTML(inquiry)

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
	builder.WriteString(fmt.Sprintf("%s\n\n", report.Content))
	
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
                <pre>%s</pre>
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
		inquiry.Name,
		inquiry.Email,
		inquiry.Company,
		inquiry.Phone,
		strings.Join(inquiry.Services, ", "),
		inquiry.ID,
		report.ID,
		inquiry.Message,
		report.Content,
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
		inquiry.Name,
		inquiry.Email,
		inquiry.Company,
		inquiry.Phone,
		strings.Join(inquiry.Services, ", "),
		inquiry.ID,
		inquiry.Message)
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

// buildCustomerConfirmationHTML creates the HTML version of the customer confirmation email
func (e *emailService) buildCustomerConfirmationHTML(inquiry *domain.Inquiry) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Thank you for your inquiry</title>
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
        .confirmation-section { 
            background-color: #e7f3ff; 
            padding: 20px; 
            margin: 20px 0; 
            border-left: 5px solid #007cba; 
            border-radius: 5px;
        }
        .details-section { 
            background-color: #f8f9fa; 
            padding: 20px; 
            margin: 20px 0; 
            border-radius: 5px;
        }
        .next-steps { 
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
        .detail-item {
            margin: 10px 0;
            padding: 8px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .detail-item:last-child { border-bottom: none; }
        .detail-item strong { color: #007cba; }
        h2, h3 { color: #007cba; margin-top: 0; }
        ul { padding-left: 20px; }
        li { margin: 8px 0; }
        .checkmark { color: #28a745; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Thank You for Your Inquiry!</h1>
        </div>
        
        <div class="content">
            <div class="confirmation-section">
                <h2>Dear %s,</h2>
                <p>We have received your inquiry for cloud consulting services and wanted to confirm that it has been <strong>successfully submitted</strong>.</p>
            </div>
            
            <div class="details-section">
                <h3>📋 Inquiry Details</h3>
                <div class="detail-item"><strong>Services Requested:</strong> %s</div>
                <div class="detail-item"><strong>Company:</strong> %s</div>
                <div class="detail-item"><strong>Reference ID:</strong> %s</div>
            </div>
            
            <div class="next-steps">
                <h3>🚀 What happens next?</h3>
                <ul>
                    <li><span class="checkmark">✓</span> Our team will review your inquiry within 24 hours</li>
                    <li><span class="checkmark">✓</span> We'll prepare a customized assessment based on your requirements</li>
                    <li><span class="checkmark">✓</span> A cloud consultant will reach out to discuss your project in detail</li>
                    <li><span class="checkmark">✓</span> We'll provide you with a detailed proposal and timeline</li>
                </ul>
            </div>
            
            <div style="background-color: #fff3cd; padding: 15px; border-radius: 5px; margin: 20px 0; border-left: 5px solid #ffc107;">
                <p style="margin: 0;"><strong>💡 Need to add something?</strong> If you have any immediate questions or need to provide additional information, please don't hesitate to contact us at <strong>info@cloudpartner.pro</strong>.</p>
            </div>
            
            <div style="text-align: center; margin: 30px 0;">
                <p><strong>Best regards,</strong><br>
                Cloud Consulting Team<br>
                <a href="mailto:info@cloudpartner.pro" style="color: #007cba;">info@cloudpartner.pro</a></p>
            </div>
        </div>
        
        <div class="footer">
            This is an automated confirmation email. Please do not reply to this message.
        </div>
    </div>
</body>
</html>`,
		inquiry.Name,
		strings.Join(inquiry.Services, ", "),
		inquiry.Company,
		inquiry.ID)
}