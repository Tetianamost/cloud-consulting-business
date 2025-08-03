package services

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockChatRateLimiter for testing
type MockChatRateLimiter struct {
	mock.Mock
}

func (m *MockChatRateLimiter) CheckRateLimit(ctx context.Context, userID string, action string) (*interfaces.RateLimitResult, error) {
	args := m.Called(ctx, userID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.RateLimitResult), args.Error(1)
}

func (m *MockChatRateLimiter) ResetRateLimit(ctx context.Context, userID string, action string) error {
	args := m.Called(ctx, userID, action)
	return args.Error(0)
}

func (m *MockChatRateLimiter) GetRateLimitStatus(ctx context.Context, userID string, action string) (*interfaces.RateLimitStatus, error) {
	args := m.Called(ctx, userID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.RateLimitStatus), args.Error(1)
}

// MockChatAuditLogger for testing
type MockChatAuditLogger struct {
	mock.Mock
}

func (m *MockChatAuditLogger) LogLogin(ctx context.Context, userID string, success bool, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, success, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogLogout(ctx context.Context, userID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogMessageSent(ctx context.Context, userID string, sessionID string, messageID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, messageID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSessionCreated(ctx context.Context, userID string, sessionID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSessionDeleted(ctx context.Context, userID string, sessionID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogRateLimitExceeded(ctx context.Context, userID string, action string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, action, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSecurityEvent(ctx context.Context, userID string, eventType string, severity string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, eventType, severity, metadata)
	return args.Error(0)
}

// Test helper functions
func createTestChatSecurityService() (*ChatSecurityService, *MockChatRateLimiter, *MockChatAuditLogger) {
	mockRateLimiter := &MockChatRateLimiter{}
	mockAuditLogger := &MockChatAuditLogger{}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests

	service := NewChatSecurityService(logger, "test-jwt-secret", mockRateLimiter, mockAuditLogger)
	return service, mockRateLimiter, mockAuditLogger
}

// Test CheckRateLimit
func TestChatSecurityService_CheckRateLimit_Success(t *testing.T) {
	service, mockRateLimiter, _ := createTestChatSecurityService()
	ctx := context.Background()

	expectedResult := &interfaces.RateLimitResult{
		Allowed:    true,
		Remaining:  9,
		ResetTime:  time.Now().Add(time.Minute),
		RetryAfter: 0,
	}

	mockRateLimiter.On("CheckRateLimit", ctx, "test-user", "message").Return(expectedResult, nil)

	result, err := service.CheckRateLimit(ctx, "test-user", "message")

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	mockRateLimiter.AssertExpectations(t)
}

func TestChatSecurityService_CheckRateLimit_Exceeded(t *testing.T) {
	service, mockRateLimiter, _ := createTestChatSecurityService()
	ctx := context.Background()

	expectedResult := &interfaces.RateLimitResult{
		Allowed:    false,
		Remaining:  0,
		ResetTime:  time.Now().Add(time.Minute),
		RetryAfter: time.Minute,
	}

	mockRateLimiter.On("CheckRateLimit", ctx, "test-user", "message").Return(expectedResult, nil)

	result, err := service.CheckRateLimit(ctx, "test-user", "message")

	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, 0, result.Remaining)
	assert.True(t, result.RetryAfter > 0)

	mockRateLimiter.AssertExpectations(t)
}

func TestChatSecurityService_CheckRateLimit_Error(t *testing.T) {
	service, mockRateLimiter, _ := createTestChatSecurityService()
	ctx := context.Background()

	mockRateLimiter.On("CheckRateLimit", ctx, "test-user", "message").Return(nil, assert.AnError)

	result, err := service.CheckRateLimit(ctx, "test-user", "message")

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRateLimiter.AssertExpectations(t)
}

// Test ValidateMessageContent
func TestChatSecurityService_ValidateMessageContent_Valid(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	tests := []string{
		"Normal message content",
		"Message with numbers 123",
		"Message with special chars !@#$%",
		"Multi-line\nmessage\ncontent",
	}

	for _, content := range tests {
		err := service.ValidateMessageContent(content)
		assert.NoError(t, err, "Content: %s", content)
	}
}

func TestChatSecurityService_ValidateMessageContent_Invalid(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	tests := []struct {
		content string
		reason  string
	}{
		{"", "empty content"},
		{string(make([]byte, 10001)), "content too long"},
		{"<script>alert('xss')</script>", "contains script tag"},
		{"javascript:alert('xss')", "contains javascript protocol"},
		{"<iframe src='evil.com'></iframe>", "contains iframe tag"},
		{"SELECT * FROM users", "potential SQL injection"},
		{"DROP TABLE users", "potential SQL injection"},
	}

	for _, test := range tests {
		err := service.ValidateMessageContent(test.content)
		assert.Error(t, err, "Content should be invalid: %s (%s)", test.content, test.reason)
	}
}

// Test SanitizeMessageContent
func TestChatSecurityService_SanitizeMessageContent(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Normal message",
			expected: "Normal message",
		},
		{
			input:    "<script>alert('xss')</script>Hello",
			expected: "Hello",
		},
		{
			input:    "Hello <b>world</b>",
			expected: "Hello world",
		},
		{
			input:    "  Multiple   spaces  ",
			expected: "Multiple spaces",
		},
		{
			input:    "<div>Content</div>",
			expected: "Content",
		},
		{
			input:    "Line 1\n\nLine 2\n\n\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, test := range tests {
		result := service.SanitizeMessageContent(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}
}

// Test ValidateSessionContext
func TestChatSecurityService_ValidateSessionContext_Valid(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	validContexts := []map[string]interface{}{
		{
			"client_name":     "Valid Client Name",
			"meeting_type":    "consultation",
			"project_context": "Valid project description",
		},
		{
			"client_name":   "Client with Numbers 123",
			"service_types": []string{"migration", "optimization"},
		},
		{
			"custom_fields": map[string]string{
				"priority": "high",
				"budget":   "100k",
			},
		},
	}

	for _, context := range validContexts {
		err := service.ValidateSessionContext(context)
		assert.NoError(t, err, "Context: %v", context)
	}
}

func TestChatSecurityService_ValidateSessionContext_Invalid(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	invalidContexts := []struct {
		context map[string]interface{}
		reason  string
	}{
		{
			context: map[string]interface{}{
				"client_name": "<script>alert('xss')</script>",
			},
			reason: "XSS in client name",
		},
		{
			context: map[string]interface{}{
				"client_name": string(make([]byte, 256)),
			},
			reason: "client name too long",
		},
		{
			context: map[string]interface{}{
				"project_context": string(make([]byte, 1001)),
			},
			reason: "project context too long",
		},
	}

	for _, test := range invalidContexts {
		err := service.ValidateSessionContext(test.context)
		assert.Error(t, err, "Context should be invalid: %s", test.reason)
	}
}

// Test DetectSuspiciousActivity
func TestChatSecurityService_DetectSuspiciousActivity_Normal(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	normalMessages := []string{
		"What is AWS EC2?",
		"How much does S3 cost?",
		"Can you help with migration planning?",
		"What are the best practices for security?",
	}

	for _, message := range normalMessages {
		suspicious, reason := service.DetectSuspiciousActivity("test-user", message, map[string]interface{}{})
		assert.False(t, suspicious, "Message should not be suspicious: %s", message)
		assert.Empty(t, reason)
	}
}

func TestChatSecurityService_DetectSuspiciousActivity_Suspicious(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	suspiciousTests := []struct {
		message  string
		metadata map[string]interface{}
		reason   string
	}{
		{
			message: "password123",
			reason:  "contains potential password",
		},
		{
			message: "My credit card number is 4111-1111-1111-1111",
			reason:  "contains potential credit card",
		},
		{
			message: "My SSN is 123-45-6789",
			reason:  "contains potential SSN",
		},
		{
			message: "SELECT * FROM users WHERE password = 'admin'",
			reason:  "contains SQL injection patterns",
		},
		{
			message: "../../etc/passwd",
			reason:  "contains path traversal patterns",
		},
		{
			message: "eval(malicious_code)",
			reason:  "contains code injection patterns",
		},
	}

	for _, test := range suspiciousTests {
		suspicious, reason := service.DetectSuspiciousActivity("test-user", test.message, test.metadata)
		assert.True(t, suspicious, "Message should be suspicious: %s", test.message)
		assert.NotEmpty(t, reason)
	}
}

func TestChatSecurityService_DetectSuspiciousActivity_RapidMessages(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	// Simulate rapid message sending
	metadata := map[string]interface{}{
		"messages_in_last_minute": 25, // Above threshold
	}

	suspicious, reason := service.DetectSuspiciousActivity("test-user", "Normal message", metadata)
	assert.True(t, suspicious)
	assert.Contains(t, reason, "rapid messaging")
}

// Test LogSecurityEvent
func TestChatSecurityService_LogSecurityEvent_Success(t *testing.T) {
	service, _, mockAuditLogger := createTestChatSecurityService()
	ctx := context.Background()

	metadata := map[string]interface{}{
		"ip_address": "192.168.1.1",
		"user_agent": "Test Agent",
	}

	mockAuditLogger.On("LogSecurityEvent", ctx, "test-user", "suspicious_message", "medium", metadata).Return(nil)

	err := service.LogSecurityEvent(ctx, "test-user", "suspicious_message", "medium", metadata)

	assert.NoError(t, err)
	mockAuditLogger.AssertExpectations(t)
}

// Test EncryptSensitiveData
func TestChatSecurityService_EncryptSensitiveData(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	originalData := "sensitive information"

	encryptedData, err := service.EncryptSensitiveData(originalData)
	assert.NoError(t, err)
	assert.NotEqual(t, originalData, encryptedData)
	assert.NotEmpty(t, encryptedData)
}

func TestChatSecurityService_EncryptSensitiveData_EmptyData(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	encryptedData, err := service.EncryptSensitiveData("")
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData) // Should still return encrypted empty string
}

// Test DecryptSensitiveData
func TestChatSecurityService_DecryptSensitiveData(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	originalData := "sensitive information"

	// Encrypt first
	encryptedData, err := service.EncryptSensitiveData(originalData)
	assert.NoError(t, err)

	// Then decrypt
	decryptedData, err := service.DecryptSensitiveData(encryptedData)
	assert.NoError(t, err)
	assert.Equal(t, originalData, decryptedData)
}

func TestChatSecurityService_DecryptSensitiveData_InvalidData(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	// Try to decrypt invalid data
	_, err := service.DecryptSensitiveData("invalid-encrypted-data")
	assert.Error(t, err)
}

// Test GenerateSecureToken
func TestChatSecurityService_GenerateSecureToken(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	token1 := service.GenerateSecureToken()
	token2 := service.GenerateSecureToken()

	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2) // Should generate unique tokens
	assert.Len(t, token1, 64)          // Should be 32 bytes hex encoded = 64 chars
	assert.Len(t, token2, 64)
}

// Test ValidateIPAddress
func TestChatSecurityService_ValidateIPAddress(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	validIPs := []string{
		"192.168.1.1",
		"10.0.0.1",
		"127.0.0.1",
		"8.8.8.8",
		"2001:db8::1",
		"::1",
	}

	for _, ip := range validIPs {
		valid := service.ValidateIPAddress(ip)
		assert.True(t, valid, "IP should be valid: %s", ip)
	}

	invalidIPs := []string{
		"256.256.256.256",
		"192.168.1",
		"not-an-ip",
		"",
		"192.168.1.1.1",
	}

	for _, ip := range invalidIPs {
		valid := service.ValidateIPAddress(ip)
		assert.False(t, valid, "IP should be invalid: %s", ip)
	}
}

// Test IsAllowedUserAgent
func TestChatSecurityService_IsAllowedUserAgent(t *testing.T) {
	service, _, _ := createTestChatSecurityService()

	allowedUserAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
	}

	for _, ua := range allowedUserAgents {
		allowed := service.IsAllowedUserAgent(ua)
		assert.True(t, allowed, "User agent should be allowed: %s", ua)
	}

	blockedUserAgents := []string{
		"curl/7.68.0",
		"wget/1.20.3",
		"python-requests/2.25.1",
		"bot/1.0",
		"scanner/1.0",
		"",
	}

	for _, ua := range blockedUserAgents {
		allowed := service.IsAllowedUserAgent(ua)
		assert.False(t, allowed, "User agent should be blocked: %s", ua)
	}
}

// Test comprehensive security validation
func TestChatSecurityService_ComprehensiveValidation(t *testing.T) {
	service, mockRateLimiter, mockAuditLogger := createTestChatSecurityService()
	ctx := context.Background()

	// Mock rate limiter to allow request
	rateLimitResult := &interfaces.RateLimitResult{
		Allowed:   true,
		Remaining: 9,
		ResetTime: time.Now().Add(time.Minute),
	}
	mockRateLimiter.On("CheckRateLimit", ctx, "test-user", "message").Return(rateLimitResult, nil)

	// Test valid message
	message := "What are the best practices for AWS security?"
	metadata := map[string]interface{}{
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	// Check rate limit
	rateLimitCheck, err := service.CheckRateLimit(ctx, "test-user", "message")
	assert.NoError(t, err)
	assert.True(t, rateLimitCheck.Allowed)

	// Validate message content
	err = service.ValidateMessageContent(message)
	assert.NoError(t, err)

	// Sanitize message content
	sanitized := service.SanitizeMessageContent(message)
	assert.Equal(t, message, sanitized) // Should be unchanged for clean message

	// Check for suspicious activity
	suspicious, reason := service.DetectSuspiciousActivity("test-user", message, metadata)
	assert.False(t, suspicious)
	assert.Empty(t, reason)

	// Validate IP address
	valid := service.ValidateIPAddress(metadata["ip_address"].(string))
	assert.True(t, valid)

	// Check user agent
	allowed := service.IsAllowedUserAgent(metadata["user_agent"].(string))
	assert.True(t, allowed)

	mockRateLimiter.AssertExpectations(t)
}

// Test security event logging integration
func TestChatSecurityService_SecurityEventLogging(t *testing.T) {
	service, _, mockAuditLogger := createTestChatSecurityService()
	ctx := context.Background()

	metadata := map[string]interface{}{
		"ip_address": "192.168.1.1",
		"message":    "suspicious content",
	}

	// Test different severity levels
	severities := []string{"low", "medium", "high", "critical"}

	for _, severity := range severities {
		mockAuditLogger.On("LogSecurityEvent", ctx, "test-user", "test_event", severity, metadata).Return(nil)

		err := service.LogSecurityEvent(ctx, "test-user", "test_event", severity, metadata)
		assert.NoError(t, err)
	}

	mockAuditLogger.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkChatSecurityService_ValidateMessageContent(b *testing.B) {
	service, _, _ := createTestChatSecurityService()
	message := "This is a normal message that should pass validation checks"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateMessageContent(message)
	}
}

func BenchmarkChatSecurityService_SanitizeMessageContent(b *testing.B) {
	service, _, _ := createTestChatSecurityService()
	message := "This is a message with <script>alert('test')</script> some HTML tags <b>bold</b> and extra   spaces"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.SanitizeMessageContent(message)
	}
}

func BenchmarkChatSecurityService_DetectSuspiciousActivity(b *testing.B) {
	service, _, _ := createTestChatSecurityService()
	message := "This is a normal message asking about AWS services and best practices"
	metadata := map[string]interface{}{
		"ip_address": "192.168.1.1",
		"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.DetectSuspiciousActivity("test-user", message, metadata)
	}
}
