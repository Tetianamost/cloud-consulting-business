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

	"github.com/cloud-consulting/backend/internal/interfaces"
	configPkg "github.com/cloud-consulting/backend/internal/config"
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

// SendEmail sends an email using AWS SES
func (s *sesService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	// Always use SendRawEmail for better email client compatibility
	// This ensures proper MIME headers and multipart/alternative structure
	// which improves HTML rendering across different email clients
	return s.sendRawEmail(timeoutCtx, email)
}

// sendSimpleEmail sends an email without attachments using the simple SES API
func (s *sesService) sendSimpleEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Build the email input
	input := &ses.SendEmailInput{
		Source: aws.String(email.From),
		Destination: &types.Destination{
			ToAddresses: email.To,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(email.Subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data:    aws.String(email.TextBody),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	// Add HTML body if provided
	if email.HTMLBody != "" {
		input.Message.Body.Html = &types.Content{
			Data:    aws.String(email.HTMLBody),
			Charset: aws.String("UTF-8"),
		}
	}

	// Add reply-to if provided
	if email.ReplyTo != "" {
		input.ReplyToAddresses = []string{email.ReplyTo}
	}

	// Send the email
	result, err := s.client.SendEmail(ctx, input)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"to":      email.To,
			"subject": email.Subject,
		}).Error("Failed to send email via SES")
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"message_id": aws.ToString(result.MessageId),
		"to":         email.To,
		"subject":    email.Subject,
	}).Info("Email sent successfully via SES")

	return nil
}

// sendRawEmail sends an email with attachments using the raw SES API
func (s *sesService) sendRawEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Build the raw email message
	rawMessage, err := s.buildRawMessage(email)
	if err != nil {
		return fmt.Errorf("failed to build raw email message: %w", err)
	}

	// Create the raw email input
	input := &ses.SendRawEmailInput{
		Source: aws.String(email.From),
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
	}).Info("Raw email with attachments sent successfully via SES")

	return nil
}

// buildRawMessage builds a raw MIME email message with attachments
func (s *sesService) buildRawMessage(email *interfaces.EmailMessage) ([]byte, error) {
	var buf bytes.Buffer

	// Write headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	if email.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", email.ReplyTo))
	}
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Create multipart writer
	writer := multipart.NewWriter(&buf)
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary()))

	// Add text/HTML body as multipart alternative
	if email.HTMLBody != "" || email.TextBody != "" {
		// Create alternative part for text and HTML
		altHeader := textproto.MIMEHeader{}
		altHeader.Set("Content-Type", fmt.Sprintf("multipart/alternative; boundary=%s", writer.Boundary()))
		
		part, err := writer.CreatePart(altHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to create alternative part: %w", err)
		}
		
		// Create a separate writer for the alternative content
		altWriter := multipart.NewWriter(part)
		
		// Add text part
		if email.TextBody != "" {
			textHeader := textproto.MIMEHeader{}
			textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
			textHeader.Set("Content-Transfer-Encoding", "7bit")
			
			textPart, err := altWriter.CreatePart(textHeader)
			if err != nil {
				return nil, fmt.Errorf("failed to create text part: %w", err)
			}
			
			textPart.Write([]byte(email.TextBody))
		}
		
		// Add HTML part
		if email.HTMLBody != "" {
			htmlHeader := textproto.MIMEHeader{}
			htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
			htmlHeader.Set("Content-Transfer-Encoding", "7bit")
			
			htmlPart, err := altWriter.CreatePart(htmlHeader)
			if err != nil {
				return nil, fmt.Errorf("failed to create HTML part: %w", err)
			}
			
			htmlPart.Write([]byte(email.HTMLBody))
		}
		
		altWriter.Close()
	}

	// Add attachments
	for _, attachment := range email.Attachments {
		header := textproto.MIMEHeader{}
		header.Set("Content-Type", attachment.ContentType)
		header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", mime.QEncoding.Encode("UTF-8", attachment.Filename)))
		header.Set("Content-Transfer-Encoding", "base64")
		
		part, err := writer.CreatePart(header)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachment part: %w", err)
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
	}

	// Close the multipart writer
	writer.Close()

	return buf.Bytes(), nil
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