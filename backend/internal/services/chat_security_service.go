package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ChatSecurityService implements the ChatSecurityService interface
type ChatSecurityService struct {
	logger           *logrus.Logger
	encryptionKey    []byte
	rateLimiter      *ChatRateLimiter
	auditLogger      interfaces.ChatAuditLogger
	contentFilters   map[string]*regexp.Regexp
	bannedWords      []string
	maxMessageLength int
	mutex            sync.RWMutex
}

// NewChatSecurityService creates a new chat security service
func NewChatSecurityService(
	logger *logrus.Logger,
	encryptionKey string,
	rateLimiter *ChatRateLimiter,
	auditLogger interfaces.ChatAuditLogger,
) *ChatSecurityService {
	// Generate encryption key from string
	key := sha256.Sum256([]byte(encryptionKey))

	service := &ChatSecurityService{
		logger:           logger,
		encryptionKey:    key[:],
		rateLimiter:      rateLimiter,
		auditLogger:      auditLogger,
		contentFilters:   make(map[string]*regexp.Regexp),
		maxMessageLength: 4000, // 4KB max message length
	}

	// Initialize content filters and banned words
	service.initializeContentFilters()
	service.initializeBannedWords()

	return service
}

// initializeContentFilters sets up regex patterns for content filtering
func (s *ChatSecurityService) initializeContentFilters() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// XSS patterns
	s.contentFilters["xss"] = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>|javascript:|on\w+\s*=|<iframe|<object|<embed`)

	// SQL injection patterns
	s.contentFilters["sql_injection"] = regexp.MustCompile(`(?i)(union\s+select|drop\s+table|delete\s+from|insert\s+into|update\s+set|exec\s*\(|execute\s*\()`)

	// Command injection patterns
	s.contentFilters["command_injection"] = regexp.MustCompile(`(?i)(\||;|&|` + "`" + `|\$\(|<\(|>\(|{|}|\[|\])`)

	// Path traversal patterns
	s.contentFilters["path_traversal"] = regexp.MustCompile(`(\.\.\/|\.\.\\|%2e%2e%2f|%2e%2e%5c)`)

	// Sensitive data patterns (credit cards, SSNs, etc.)
	s.contentFilters["sensitive_data"] = regexp.MustCompile(`(?i)(\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b|\b\d{3}-\d{2}-\d{4}\b)`)
}

// initializeBannedWords sets up a list of banned words/phrases
func (s *ChatSecurityService) initializeBannedWords() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// This would typically be loaded from a database or configuration file
	s.bannedWords = []string{
		// Add inappropriate words/phrases here
		// For demo purposes, keeping it minimal
		"spam",
		"malicious",
	}
}

// ValidateMessageContent validates message content for security issues
func (s *ChatSecurityService) ValidateMessageContent(content string) error {
	// Check message length
	if len(content) > s.maxMessageLength {
		return interfaces.SecurityError{
			Code:    interfaces.ErrCodeValidationFailed,
			Message: fmt.Sprintf("Message too long. Maximum length is %d characters", s.maxMessageLength),
		}
	}

	// Check for empty content
	if strings.TrimSpace(content) == "" {
		return interfaces.SecurityError{
			Code:    interfaces.ErrCodeValidationFailed,
			Message: "Message content cannot be empty",
		}
	}

	// Run content filters
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for filterName, pattern := range s.contentFilters {
		if pattern.MatchString(content) {
			return interfaces.SecurityError{
				Code:    interfaces.ErrCodeContentViolation,
				Message: fmt.Sprintf("Content violates security policy: %s", filterName),
				Details: filterName,
			}
		}
	}

	return nil
}

// SanitizeMessageContent sanitizes message content to remove potentially harmful elements
func (s *ChatSecurityService) SanitizeMessageContent(content string) string {
	// HTML escape
	sanitized := html.EscapeString(content)

	// Remove excessive whitespace
	sanitized = regexp.MustCompile(`\s+`).ReplaceAllString(sanitized, " ")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Remove null bytes
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	return sanitized
}

// ValidateSessionData validates session data for security issues
func (s *ChatSecurityService) ValidateSessionData(session interface{}) error {
	// This would validate session data structure and content
	// For now, we'll implement basic validation

	if session == nil {
		return interfaces.SecurityError{
			Code:    interfaces.ErrCodeValidationFailed,
			Message: "Session data cannot be nil",
		}
	}

	return nil
}

// FilterContent filters content and returns the result
func (s *ChatSecurityService) FilterContent(content string) (*interfaces.ContentFilterResult, error) {
	result := &interfaces.ContentFilterResult{
		IsAllowed:    true,
		FilteredText: content,
		Violations:   []string{},
		Confidence:   1.0,
		Categories:   []string{},
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Check against content filters
	for filterName, pattern := range s.contentFilters {
		if pattern.MatchString(content) {
			result.IsAllowed = false
			result.Violations = append(result.Violations, filterName)
			result.Categories = append(result.Categories, filterName)
			result.Confidence = 0.9 // High confidence for regex matches
		}
	}

	// Check banned words
	lowerContent := strings.ToLower(content)
	for _, word := range s.bannedWords {
		if strings.Contains(lowerContent, strings.ToLower(word)) {
			result.IsAllowed = false
			result.Violations = append(result.Violations, "banned_word")
			result.Categories = append(result.Categories, "inappropriate_content")
			result.Confidence = 0.8

			// Replace banned words with asterisks
			result.FilteredText = strings.ReplaceAll(result.FilteredText, word, strings.Repeat("*", len(word)))
		}
	}

	return result, nil
}

// ModerateMessage performs content moderation on a message
func (s *ChatSecurityService) ModerateMessage(content string) (*interfaces.ModerationResult, error) {
	// First run content filtering
	filterResult, err := s.FilterContent(content)
	if err != nil {
		return nil, fmt.Errorf("failed to filter content: %w", err)
	}

	result := &interfaces.ModerationResult{
		IsApproved:   filterResult.IsAllowed,
		Confidence:   filterResult.Confidence,
		Categories:   filterResult.Categories,
		Reasons:      filterResult.Violations,
		Metadata:     make(map[string]interface{}),
		ReviewNeeded: false,
	}

	// Determine if manual review is needed
	if !result.IsApproved && result.Confidence < 0.9 {
		result.ReviewNeeded = true
	}

	// Add metadata
	result.Metadata["content_length"] = len(content)
	result.Metadata["filtered_text"] = filterResult.FilteredText
	result.Metadata["timestamp"] = time.Now()

	return result, nil
}

// IsContentAllowed checks if content is allowed
func (s *ChatSecurityService) IsContentAllowed(content string) (bool, error) {
	filterResult, err := s.FilterContent(content)
	if err != nil {
		return false, err
	}

	return filterResult.IsAllowed, nil
}

// CheckRateLimit checks if an action is within rate limits
func (s *ChatSecurityService) CheckRateLimit(ctx context.Context, userID string, action string) (*interfaces.RateLimitResult, error) {
	if s.rateLimiter == nil {
		// If no rate limiter configured, allow all requests
		return &interfaces.RateLimitResult{
			Allowed:       true,
			Remaining:     100,
			ResetTime:     time.Now().Add(time.Hour),
			RetryAfter:    0,
			TotalRequests: 0,
		}, nil
	}

	switch action {
	case "message":
		return s.rateLimiter.AllowMessage(ctx, userID)
	case "connection":
		return s.rateLimiter.AllowConnection(ctx, userID)
	case "session_creation":
		return s.rateLimiter.AllowSessionCreation(ctx, userID)
	default:
		// Default rate limit: 100 requests per hour
		return s.rateLimiter.CheckLimit(ctx, fmt.Sprintf("%s:%s", userID, action), 100, time.Hour)
	}
}

// IncrementRateLimit increments the rate limit counter for an action
func (s *ChatSecurityService) IncrementRateLimit(ctx context.Context, userID string, action string) error {
	// This would typically increment counters in Redis or similar
	// For now, we'll just log the action
	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"action":  action,
	}).Debug("Rate limit incremented")

	return nil
}

// ResetRateLimit resets the rate limit for a user and action
func (s *ChatSecurityService) ResetRateLimit(ctx context.Context, userID string, action string) error {
	if s.rateLimiter != nil {
		return s.rateLimiter.ResetUserLimits(ctx, userID)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"action":  action,
	}).Info("Rate limit reset")

	return nil
}

// LogSecurityEvent logs a security event
func (s *ChatSecurityService) LogSecurityEvent(ctx context.Context, event *interfaces.SecurityEvent) error {
	s.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"event_type": event.Type,
		"user_id":    event.UserID,
		"resource":   event.Resource,
		"action":     event.Action,
		"result":     event.Result,
		"severity":   event.Severity,
		"ip_address": event.IPAddress,
		"user_agent": event.UserAgent,
	}).Warn("Security event logged")

	// If audit logger is available, use it
	if s.auditLogger != nil {
		violation := &interfaces.SecurityViolation{
			Type:        event.Type,
			Severity:    event.Severity,
			Description: event.Description,
			Resource:    event.Resource,
			Action:      event.Action,
			IPAddress:   event.IPAddress,
			UserAgent:   event.UserAgent,
			Metadata:    event.Metadata,
		}
		return s.auditLogger.LogSecurityViolation(ctx, event.UserID, violation)
	}

	return nil
}

// LogAuthenticationAttempt logs an authentication attempt
func (s *ChatSecurityService) LogAuthenticationAttempt(ctx context.Context, userID string, success bool, details map[string]interface{}) error {
	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"success": success,
		"details": details,
	}).Info("Authentication attempt logged")

	// If audit logger is available, use it
	if s.auditLogger != nil {
		if success {
			return s.auditLogger.LogLogin(ctx, userID, true, details)
		} else {
			return s.auditLogger.LogLogin(ctx, userID, false, details)
		}
	}

	return nil
}

// LogDataAccess logs data access events
func (s *ChatSecurityService) LogDataAccess(ctx context.Context, userID string, resource string, action string) error {
	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"resource": resource,
		"action":   action,
	}).Debug("Data access logged")

	// If audit logger is available, use it
	if s.auditLogger != nil {
		metadata := map[string]interface{}{
			"resource": resource,
			"action":   action,
		}

		switch resource {
		case "session":
			return s.auditLogger.LogSessionAccessed(ctx, userID, resource, metadata)
		case "message":
			return s.auditLogger.LogMessageAccessed(ctx, userID, resource, metadata)
		default:
			// Generic audit log entry
			return nil
		}
	}

	return nil
}

// EncryptSensitiveData encrypts sensitive data using AES-GCM
func (s *ChatSecurityService) EncryptSensitiveData(data string) (string, error) {
	if data == "" {
		return "", nil
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeEncryptionFailed,
			Message: "Failed to create cipher",
			Details: err.Error(),
		}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeEncryptionFailed,
			Message: "Failed to create GCM",
			Details: err.Error(),
		}
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeEncryptionFailed,
			Message: "Failed to generate nonce",
			Details: err.Error(),
		}
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSensitiveData decrypts sensitive data using AES-GCM
func (s *ChatSecurityService) DecryptSensitiveData(encryptedData string) (string, error) {
	if encryptedData == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeDecryptionFailed,
			Message: "Failed to decode base64",
			Details: err.Error(),
		}
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeDecryptionFailed,
			Message: "Failed to create cipher",
			Details: err.Error(),
		}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeDecryptionFailed,
			Message: "Failed to create GCM",
			Details: err.Error(),
		}
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeDecryptionFailed,
			Message: "Ciphertext too short",
		}
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeDecryptionFailed,
			Message: "Failed to decrypt data",
			Details: err.Error(),
		}
	}

	return string(plaintext), nil
}

// HashSensitiveData creates a hash of sensitive data for comparison
func (s *ChatSecurityService) HashSensitiveData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// AddBannedWord adds a word to the banned words list
func (s *ChatSecurityService) AddBannedWord(word string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.bannedWords = append(s.bannedWords, strings.ToLower(word))
}

// RemoveBannedWord removes a word from the banned words list
func (s *ChatSecurityService) RemoveBannedWord(word string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, bannedWord := range s.bannedWords {
		if strings.ToLower(bannedWord) == strings.ToLower(word) {
			s.bannedWords = append(s.bannedWords[:i], s.bannedWords[i+1:]...)
			break
		}
	}
}

// UpdateContentFilter updates or adds a content filter
func (s *ChatSecurityService) UpdateContentFilter(name, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.contentFilters[name] = regex
	return nil
}

// RemoveContentFilter removes a content filter
func (s *ChatSecurityService) RemoveContentFilter(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.contentFilters, name)
}
