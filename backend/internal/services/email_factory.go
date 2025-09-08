package services

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// NewEmailServiceWithSES creates a new email service with proper SES implementation
// This method is now deprecated in favor of NewEmailServiceWithEventRecorder
// It's maintained for backward compatibility but should not be used in production
func NewEmailServiceWithSES(cfg config.SESConfig, logger *logrus.Logger) (interfaces.EmailService, error) {
	logger.Warn("Using deprecated NewEmailServiceWithSES - consider using NewEmailServiceWithEventRecorder for production")

	// Create the FIXED SES service (resolves MIME boundary issues)
	sesService, err := NewSESService(cfg, logger)
	if err != nil {
		return nil, err
	}

	// Create template service
	templateService := NewTemplateService("templates", logger)

	// Create email service with fixed SES implementation (no event recorder)
	emailService := NewEmailService(sesService, templateService, nil, cfg, logger)

	return emailService, nil
}

// NewEmailServiceWithEventRecorder creates a new email service with event recording capability
// This is the recommended method for production use as it includes email event tracking
func NewEmailServiceWithEventRecorder(cfg config.SESConfig, eventRecorder interfaces.EmailEventRecorder, logger *logrus.Logger) (interfaces.EmailService, error) {
	// Validate configuration for email event tracking
	if err := validateEmailEventConfig(cfg, eventRecorder, logger); err != nil {
		return nil, fmt.Errorf("email event configuration validation failed: %w", err)
	}

	// Create the FIXED SES service (resolves MIME boundary issues)
	sesService, err := NewSESService(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SES service: %w", err)
	}

	// Create template service
	templateService := NewTemplateService("templates", logger)

	// Create email service with fixed SES implementation and event recorder
	emailService := NewEmailService(sesService, templateService, eventRecorder, cfg, logger)

	logger.WithFields(logrus.Fields{
		"sender_email":    cfg.SenderEmail,
		"reply_to_email":  cfg.ReplyToEmail,
		"region":          cfg.Region,
		"event_recording": eventRecorder != nil,
		"timeout":         cfg.Timeout,
	}).Info("Email service with event recording initialized successfully")

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

	// Create email service with mock SES (no event recorder for testing)
	return NewEmailService(mockSES, templateService, nil, cfg, logger)
}

// NewEmailServiceForTestingWithEventRecorder creates an email service with mock SES and event recorder for testing
func NewEmailServiceForTestingWithEventRecorder(mockSES interfaces.SESService, eventRecorder interfaces.EmailEventRecorder, logger *logrus.Logger) interfaces.EmailService {
	// Create template service
	templateService := NewTemplateService("templates", logger)

	// Create minimal SES config for testing
	cfg := config.SESConfig{
		SenderEmail:  "test@example.com",
		ReplyToEmail: "test@example.com",
		Region:       "us-east-1",
		Timeout:      30,
	}

	// Create email service with mock SES and event recorder
	return NewEmailService(mockSES, templateService, eventRecorder, cfg, logger)
}

// NewEmailServiceFactory creates a new email service with automatic event recording detection
// This is the primary factory method that should be used in production and test environments
func NewEmailServiceFactory(cfg config.SESConfig, eventRecorder interfaces.EmailEventRecorder, logger *logrus.Logger) (interfaces.EmailService, error) {
	// Validate basic SES configuration
	if err := ValidateSESConfig(cfg); err != nil {
		return nil, fmt.Errorf("SES configuration validation failed: %w", err)
	}

	// If event recorder is provided, use the event recording version
	if eventRecorder != nil {
		logger.WithField("event_recording", "enabled").Info("Creating email service with event recording enabled")
		return NewEmailServiceWithEventRecorder(cfg, eventRecorder, logger)
	}

	// Fall back to basic email service without event recording
	logger.WithField("event_recording", "disabled").Warn("Creating email service without event recording - consider enabling email event tracking for production")
	return NewEmailServiceWithSES(cfg, logger)
}

// NewEmailServiceForProduction creates a new email service optimized for production environments
// This method includes comprehensive validation and proper error handling for production use
func NewEmailServiceForProduction(cfg config.SESConfig, eventRecorder interfaces.EmailEventRecorder, logger *logrus.Logger) (interfaces.EmailService, error) {
	// Comprehensive production validation
	if err := validateProductionConfig(cfg, eventRecorder, logger); err != nil {
		return nil, fmt.Errorf("production configuration validation failed: %w", err)
	}

	// Always require event recording in production
	if eventRecorder == nil {
		return nil, fmt.Errorf("event recorder is required for production email service - email monitoring is mandatory")
	}

	// Create email service with event recording
	emailService, err := NewEmailServiceWithEventRecorder(cfg, eventRecorder, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create production email service: %w", err)
	}

	// Verify service health before returning
	if !emailService.IsHealthy() {
		logger.Error("Email service health check failed after initialization")
		return nil, fmt.Errorf("email service failed health check after initialization")
	}

	logger.WithFields(logrus.Fields{
		"sender_email":    cfg.SenderEmail,
		"region":          cfg.Region,
		"event_recording": true,
		"environment":     "production",
	}).Info("Production email service initialized successfully with full monitoring")

	return emailService, nil
}

// ValidateSESConfig validates the basic SES configuration
func ValidateSESConfig(cfg config.SESConfig) error {
	if cfg.SenderEmail == "" {
		return fmt.Errorf("sender email is required")
	}

	if cfg.AccessKeyID == "" {
		return fmt.Errorf("AWS access key ID is required")
	}

	if cfg.SecretAccessKey == "" {
		return fmt.Errorf("AWS secret access key is required")
	}

	if cfg.Region == "" {
		return fmt.Errorf("AWS region is required")
	}

	// Validate email format (basic check)
	if !IsValidEmail(cfg.SenderEmail) {
		return fmt.Errorf("sender email format is invalid: %s", cfg.SenderEmail)
	}

	if cfg.ReplyToEmail != "" && !IsValidEmail(cfg.ReplyToEmail) {
		return fmt.Errorf("reply-to email format is invalid: %s", cfg.ReplyToEmail)
	}

	return nil
}

// validateEmailEventConfig validates the email event tracking configuration
func validateEmailEventConfig(cfg config.SESConfig, eventRecorder interfaces.EmailEventRecorder, logger *logrus.Logger) error {
	// First validate basic SES config
	if err := ValidateSESConfig(cfg); err != nil {
		return fmt.Errorf("SES configuration validation failed: %w", err)
	}

	// Validate event recorder
	if eventRecorder == nil {
		return fmt.Errorf("event recorder is required for email event tracking")
	}

	// Test event recorder health if possible
	if healthChecker, ok := eventRecorder.(interface{ IsHealthy() bool }); ok {
		if !healthChecker.IsHealthy() {
			logger.Warn("Email event recorder health check failed, but continuing with initialization")
			// Don't fail initialization, just warn - the recorder might recover
		}
	}

	// Validate timeout configuration
	if cfg.Timeout <= 0 {
		logger.Warn("SES timeout not configured, using default of 30 seconds")
	} else if cfg.Timeout > 300 {
		logger.Warn("SES timeout is very high (>5 minutes), this may cause request timeouts")
	}

	// Log configuration for debugging
	logger.WithFields(logrus.Fields{
		"sender_email":     cfg.SenderEmail,
		"reply_to_email":   cfg.ReplyToEmail,
		"region":           cfg.Region,
		"timeout":          cfg.Timeout,
		"event_recording":  true,
		"validation_stage": "email_event_config",
	}).Debug("Email event tracking configuration validated successfully")

	return nil
}

// validateProductionConfig validates configuration specifically for production environments
func validateProductionConfig(cfg config.SESConfig, eventRecorder interfaces.EmailEventRecorder, logger *logrus.Logger) error {
	// First validate basic SES config
	if err := ValidateSESConfig(cfg); err != nil {
		return fmt.Errorf("basic SES configuration validation failed: %w", err)
	}

	// Production-specific validations
	if eventRecorder == nil {
		return fmt.Errorf("event recorder is mandatory for production environments")
	}

	// Validate sender email is not a test/development email
	testDomains := []string{"example.com", "test.com", "localhost", "127.0.0.1"}
	for _, testDomain := range testDomains {
		if strings.Contains(cfg.SenderEmail, testDomain) {
			return fmt.Errorf("sender email contains test domain '%s' - not suitable for production", testDomain)
		}
	}

	// Validate reply-to email is configured for production
	if cfg.ReplyToEmail == "" {
		logger.Warn("Reply-to email not configured - using sender email as reply-to")
	}

	// Validate timeout is reasonable for production
	if cfg.Timeout < 10 {
		return fmt.Errorf("SES timeout too low for production (%d seconds) - minimum 10 seconds recommended", cfg.Timeout)
	}
	if cfg.Timeout > 120 {
		logger.Warn("SES timeout is very high (%d seconds) - may cause slow email delivery", cfg.Timeout)
	}

	// Test event recorder health
	if healthChecker, ok := eventRecorder.(interface{ IsHealthy() bool }); ok {
		if !healthChecker.IsHealthy() {
			return fmt.Errorf("email event recorder failed health check - email monitoring is not functional")
		}
	}

	logger.WithFields(logrus.Fields{
		"sender_email":     cfg.SenderEmail,
		"reply_to_email":   cfg.ReplyToEmail,
		"region":           cfg.Region,
		"timeout":          cfg.Timeout,
		"validation_stage": "production_config",
	}).Debug("Production email configuration validated successfully")

	return nil
}

// IsValidEmail performs basic email validation
func IsValidEmail(email string) bool {
	// Basic email validation - contains @ and at least one dot after @
	if email == "" {
		return false
	}

	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return false // Multiple @ symbols
			}
			atIndex = i
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false // No @, @ at start, or @ at end
	}

	// Check for at least one dot after @
	domainPart := email[atIndex+1:]
	hasDot := false
	for _, char := range domainPart {
		if char == '.' {
			hasDot = true
			break
		}
	}

	return hasDot && len(domainPart) > 2
}
