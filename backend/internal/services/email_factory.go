package services

import (
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// NewEmailServiceWithSES creates a new email service with proper SES implementation
func NewEmailServiceWithSES(cfg config.SESConfig, logger *logrus.Logger) (interfaces.EmailService, error) {
	// Create the FIXED SES service (resolves MIME boundary issues)
	sesService, err := NewSESService(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Create template service
	templateService := NewTemplateService("templates", logger)

	// Create email service with fixed SES implementation
	emailService := NewEmailService(sesService, templateService, cfg, logger)

	return emailService, nil
}

// NewEmailServiceForTesting creates an email service with mock SES for testing
func NewEmailServiceForTesting(mockSES interfaces.SESService, logger *logrus.Logger) interfaces.EmailService {
	// Create template service
	templateService := NewTemplateService("templates", logger)

	// Create minimal SES config for testing
	cfg := config.SESConfig{
		SenderEmail:  "test@example.com",
		ReplyToEmail: "test@example.com",
		Region:       "us-east-1",
		Timeout:      30,
	}

	// Create email service with mock SES
	return NewEmailService(mockSES, templateService, cfg, logger)
}
