package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/sirupsen/logrus"

	configPkg "github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// sesService implements the SESService interface
type sesService struct {
	client *ses.Client
	config configPkg.SESConfig
	logger *logrus.Logger
}

// NewSESService creates a new SES service instance
func NewSESService(cfg configPkg.SESConfig, logger *logrus.Logger) (interfaces.SESService, error) {
	// Create AWS config with static credentials
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create SES client
	client := ses.NewFromConfig(awsConfig)

	return &sesService{
		client: client,
		config: cfg,
		logger: logger,
	}, nil
}

// SendEmail sends an email using AWS SES with proper MIME structure and returns message ID
func (s *sesService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	// Always use SendRawEmail for better compatibility and proper MIME structure
	return s.sendRawEmail(timeoutCtx, email)
}

// sendRawEmail sends an email using the raw SES API with corrected MIME structure
func (s *sesService) sendRawEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Build the raw email message with proper MIME structure
	rawMessage, err := s.buildRawMessage(email)
	if err != nil {
		return fmt.Errorf("failed to build raw email message: %w", err)
	}

	// Create the raw email input
	input := &ses.SendRawEmailInput{
		Source:       aws.String(email.From),
		Destinations: email.To,
		RawMessage: &types.RawMessage{
			Data: rawMessage,
		},
	}

	// Send the raw email
	result, err := s.client.SendRawEmail(ctx, input)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"to":          email.To,
			"subject":     email.Subject,
			"attachments": len(email.Attachments),
		}).Error("Failed to send raw email via SES")
		return fmt.Errorf("failed to send raw email: %w", err)
	}

	// Set the message ID on the email message for event recording
	email.MessageID = aws.ToString(result.MessageId)

	s.logger.WithFields(logrus.Fields{
		"message_id":  email.MessageID,
		"to":          email.To,
		"subject":     email.Subject,
		"attachments": len(email.Attachments),
	}).Info("Raw email sent successfully via SES")

	return nil
}

// buildRawMessage builds a properly structured MIME email message
func (s *sesService) buildRawMessage(email *interfaces.EmailMessage) ([]byte, error) {
	var buf bytes.Buffer

	// Write email headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	if email.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", email.ReplyTo))
	}
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Determine the main content type based on whether we have attachments
	if len(email.Attachments) > 0 {
		// With attachments: use multipart/mixed
		return s.buildMixedMessage(&buf, email)
	} else if email.HTMLBody != "" && email.TextBody != "" {
		// Both text and HTML: use multipart/alternative
		return s.buildAlternativeMessage(&buf, email)
	} else if email.HTMLBody != "" {
		// HTML only: simple message
		return s.buildSimpleHTMLMessage(&buf, email)
	} else {
		// Text only: simple message
		return s.buildSimpleTextMessage(&buf, email)
	}
}

// buildMixedMessage builds a multipart/mixed message (with attachments)
func (s *sesService) buildMixedMessage(buf *bytes.Buffer, email *interfaces.EmailMessage) ([]byte, error) {
	writer := multipart.NewWriter(buf)

	// Set content type header
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary()))

	// Add the email body as multipart/alternative if we have both text and HTML
	if email.HTMLBody != "" && email.TextBody != "" {
		// Create nested multipart/alternative for text and HTML
		altHeader := textproto.MIMEHeader{}
		altHeader.Set("Content-Type", "multipart/alternative")

		altPart, err := writer.CreatePart(altHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to create alternative part: %w", err)
		}

		altWriter := multipart.NewWriter(altPart)
		altPart.Write([]byte(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", altWriter.Boundary())))

		// Add text part
		if err := s.addTextPart(altWriter, email.TextBody); err != nil {
			return nil, err
		}

		// Add HTML part
		if err := s.addHTMLPart(altWriter, email.HTMLBody); err != nil {
			return nil, err
		}

		altWriter.Close()
	} else if email.HTMLBody != "" {
		// HTML only
		if err := s.addHTMLPart(writer, email.HTMLBody); err != nil {
			return nil, err
		}
	} else if email.TextBody != "" {
		// Text only
		if err := s.addTextPart(writer, email.TextBody); err != nil {
			return nil, err
		}
	}

	// Add attachments
	for _, attachment := range email.Attachments {
		if err := s.addAttachment(writer, attachment); err != nil {
			return nil, err
		}
	}

	writer.Close()
	return buf.Bytes(), nil
}

// buildAlternativeMessage builds a multipart/alternative message (text + HTML, no attachments)
func (s *sesService) buildAlternativeMessage(buf *bytes.Buffer, email *interfaces.EmailMessage) ([]byte, error) {
	writer := multipart.NewWriter(buf)

	// Set content type header
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", writer.Boundary()))

	// Add text part
	if err := s.addTextPart(writer, email.TextBody); err != nil {
		return nil, err
	}

	// Add HTML part
	if err := s.addHTMLPart(writer, email.HTMLBody); err != nil {
		return nil, err
	}

	writer.Close()
	return buf.Bytes(), nil
}

// buildSimpleHTMLMessage builds a simple HTML message
func (s *sesService) buildSimpleHTMLMessage(buf *bytes.Buffer, email *interfaces.EmailMessage) ([]byte, error) {
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	buf.WriteString(email.HTMLBody)
	return buf.Bytes(), nil
}

// buildSimpleTextMessage builds a simple text message
func (s *sesService) buildSimpleTextMessage(buf *bytes.Buffer, email *interfaces.EmailMessage) ([]byte, error) {
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	buf.WriteString(email.TextBody)
	return buf.Bytes(), nil
}

// addTextPart adds a text part to the multipart writer
func (s *sesService) addTextPart(writer *multipart.Writer, textBody string) error {
	textHeader := textproto.MIMEHeader{}
	textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	textHeader.Set("Content-Transfer-Encoding", "7bit")

	textPart, err := writer.CreatePart(textHeader)
	if err != nil {
		return fmt.Errorf("failed to create text part: %w", err)
	}

	_, err = textPart.Write([]byte(textBody))
	if err != nil {
		return fmt.Errorf("failed to write text content: %w", err)
	}

	return nil
}

// addHTMLPart adds an HTML part to the multipart writer
func (s *sesService) addHTMLPart(writer *multipart.Writer, htmlBody string) error {
	htmlHeader := textproto.MIMEHeader{}
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlHeader.Set("Content-Transfer-Encoding", "7bit")

	htmlPart, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %w", err)
	}

	_, err = htmlPart.Write([]byte(htmlBody))
	if err != nil {
		return fmt.Errorf("failed to write HTML content: %w", err)
	}

	return nil
}

// addAttachment adds an attachment to the multipart writer
func (s *sesService) addAttachment(writer *multipart.Writer, attachment interfaces.EmailAttachment) error {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", attachment.ContentType)
	header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", mime.QEncoding.Encode("UTF-8", attachment.Filename)))
	header.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(header)
	if err != nil {
		return fmt.Errorf("failed to create attachment part: %w", err)
	}

	// Encode attachment data as base64
	encoded := base64.StdEncoding.EncodeToString(attachment.Data)

	// Write base64 data with line breaks every 76 characters (RFC 2045)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		part.Write([]byte(encoded[i:end] + "\r\n"))
	}

	return nil
}

// VerifyEmailAddress verifies an email address with SES
func (s *sesService) VerifyEmailAddress(ctx context.Context, email string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	input := &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(email),
	}

	_, err := s.client.VerifyEmailIdentity(timeoutCtx, input)
	if err != nil {
		s.logger.WithError(err).WithField("email", email).Error("Failed to verify email address")
		return fmt.Errorf("failed to verify email address: %w", err)
	}

	s.logger.WithField("email", email).Info("Email address verification initiated")
	return nil
}

// GetSendingQuota retrieves the current SES sending quota
func (s *sesService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	result, err := s.client.GetSendQuota(timeoutCtx, &ses.GetSendQuotaInput{})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get SES sending quota")
		return nil, fmt.Errorf("failed to get sending quota: %w", err)
	}

	quota := &interfaces.SendingQuota{
		Max24HourSend:   result.Max24HourSend,
		MaxSendRate:     result.MaxSendRate,
		SentLast24Hours: result.SentLast24Hours,
	}

	return quota, nil
}

// GetDeliveryStatus retrieves the delivery status of an email by message ID
// Note: This is a placeholder implementation as SES doesn't provide direct message status lookup
// In practice, delivery status is tracked through SNS notifications
func (s *sesService) GetDeliveryStatus(ctx context.Context, messageID string) (*interfaces.EmailDeliveryStatus, error) {
	s.logger.WithField("message_id", messageID).Debug("Getting delivery status for message")

	// SES doesn't provide a direct API to get message status by ID
	// Status tracking is typically done through SNS notifications
	// This method returns a basic status indicating the message was sent
	status := &interfaces.EmailDeliveryStatus{
		MessageID:   messageID,
		Status:      "sent", // We can only confirm it was sent via SES
		Timestamp:   time.Now(),
		Destination: "", // Would need to be tracked separately
	}

	s.logger.WithFields(logrus.Fields{
		"message_id": messageID,
		"status":     status.Status,
	}).Debug("Retrieved delivery status")

	return status, nil
}

// ProcessSESNotification processes SES notifications (bounces, complaints, deliveries)
func (s *sesService) ProcessSESNotification(ctx context.Context, notification *interfaces.SESNotification) (*interfaces.SESNotificationResult, error) {
	s.logger.WithFields(logrus.Fields{
		"message_id":        notification.MessageID,
		"notification_type": notification.NotificationType,
		"source":            notification.Source,
		"destinations":      notification.Destination,
	}).Info("Processing SES notification")

	result := &interfaces.SESNotificationResult{
		MessageID:        notification.MessageID,
		NotificationType: notification.NotificationType,
		ProcessedAt:      time.Now(),
		UpdatedEvents:    0,
	}

	// Process different notification types
	switch notification.NotificationType {
	case "Bounce":
		result.Status = "bounced"
		if notification.Bounce != nil {
			s.logger.WithFields(logrus.Fields{
				"bounce_type":    notification.Bounce.BounceType,
				"bounce_subtype": notification.Bounce.BounceSubType,
				"bounced_count":  len(notification.Bounce.BouncedRecipients),
			}).Info("Processing bounce notification")
		}

	case "Complaint":
		result.Status = "spam"
		if notification.Complaint != nil {
			s.logger.WithFields(logrus.Fields{
				"complaint_type":   notification.Complaint.ComplaintFeedbackType,
				"complained_count": len(notification.Complaint.ComplainedRecipients),
			}).Info("Processing complaint notification")
		}

	case "Delivery":
		result.Status = "delivered"
		if notification.Delivery != nil {
			s.logger.WithFields(logrus.Fields{
				"processing_time_ms": notification.Delivery.ProcessingTimeMillis,
				"recipients_count":   len(notification.Delivery.Recipients),
				"smtp_response":      notification.Delivery.SMTPResponse,
			}).Info("Processing delivery notification")
		}

	default:
		result.Error = fmt.Sprintf("Unknown notification type: %s", notification.NotificationType)
		s.logger.WithField("notification_type", notification.NotificationType).Warn("Unknown SES notification type")
		return result, fmt.Errorf("unknown notification type: %s", notification.NotificationType)
	}

	s.logger.WithFields(logrus.Fields{
		"message_id":        result.MessageID,
		"notification_type": result.NotificationType,
		"status":            result.Status,
		"updated_events":    result.UpdatedEvents,
	}).Info("SES notification processed successfully")

	return result, nil
}

// CategorizeError categorizes email errors for better handling and reporting
func (s *sesService) CategorizeError(errorType string, errorMessage string) *interfaces.EmailErrorCategory {
	category := &interfaces.EmailErrorCategory{
		Category:   "unknown",
		Severity:   "permanent",
		Reason:     errorMessage,
		Actionable: false,
	}

	// Normalize error type for comparison
	errorTypeLower := strings.ToLower(errorType)
	errorMessageLower := strings.ToLower(errorMessage)

	// Categorize bounces
	if strings.Contains(errorTypeLower, "bounce") {
		category.Category = "bounce"

		// Determine bounce severity
		if strings.Contains(errorMessageLower, "permanent") ||
			strings.Contains(errorMessageLower, "5.") ||
			strings.Contains(errorMessageLower, "mailbox") && strings.Contains(errorMessageLower, "not") && strings.Contains(errorMessageLower, "exist") ||
			strings.Contains(errorMessageLower, "user unknown") ||
			strings.Contains(errorMessageLower, "invalid") && strings.Contains(errorMessageLower, "recipient") {
			category.Severity = "permanent"
			category.Reason = "Permanent bounce - recipient address is invalid or doesn't exist"
			category.Actionable = true // Can remove from mailing list
		} else if strings.Contains(errorMessageLower, "temporary") ||
			strings.Contains(errorMessageLower, "4.") ||
			strings.Contains(errorMessageLower, "mailbox full") ||
			strings.Contains(errorMessageLower, "quota exceeded") ||
			strings.Contains(errorMessageLower, "try again") {
			category.Severity = "temporary"
			category.Reason = "Temporary bounce - recipient mailbox may be full or temporarily unavailable"
			category.Actionable = true
			retryAfter := 3600 // Retry after 1 hour
			category.RetryAfter = &retryAfter
		}
	}

	// Categorize complaints (spam reports)
	if strings.Contains(errorTypeLower, "complaint") || strings.Contains(errorMessageLower, "spam") {
		category.Category = "complaint"
		category.Severity = "permanent"
		category.Reason = "Recipient marked email as spam"
		category.Actionable = true // Should unsubscribe recipient
	}

	// Categorize delivery delays
	if strings.Contains(errorMessageLower, "delay") || strings.Contains(errorMessageLower, "deferred") {
		category.Category = "delivery_delay"
		category.Severity = "warning"
		category.Reason = "Email delivery delayed but will be retried"
		category.Actionable = false
	}

	// Categorize rate limiting
	if strings.Contains(errorMessageLower, "throttl") || strings.Contains(errorMessageLower, "rate limit") {
		category.Category = "rate_limit"
		category.Severity = "temporary"
		category.Reason = "Sending rate limit exceeded"
		category.Actionable = true
		retryAfter := 300 // Retry after 5 minutes
		category.RetryAfter = &retryAfter
	}

	// Categorize authentication/configuration errors
	if strings.Contains(errorMessageLower, "authentication") ||
		strings.Contains(errorMessageLower, "credential") ||
		strings.Contains(errorMessageLower, "unauthorized") {
		category.Category = "configuration"
		category.Severity = "permanent"
		category.Reason = "Email service configuration error"
		category.Actionable = true
	}

	s.logger.WithFields(logrus.Fields{
		"error_type":  errorType,
		"category":    category.Category,
		"severity":    category.Severity,
		"actionable":  category.Actionable,
		"retry_after": category.RetryAfter,
	}).Debug("Categorized email error")

	return category
}
