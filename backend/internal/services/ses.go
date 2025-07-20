package services

import (
	"context"
	"fmt"
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
	result, err := s.client.SendEmail(timeoutCtx, input)
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