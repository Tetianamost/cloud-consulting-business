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

// SendEmail sends an email using AWS SES with proper MIME structure
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

	s.logger.WithFields(logrus.Fields{
		"message_id":  aws.ToString(result.MessageId),
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
