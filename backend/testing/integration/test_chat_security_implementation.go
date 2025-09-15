package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	fmt.Println("Testing Chat Security Implementation...")

	// Test 1: Authentication Service
	fmt.Println("\n=== Testing Authentication Service ===")
	testAuthService(logger)

	// Test 2: Security Service
	fmt.Println("\n=== Testing Security Service ===")
	testSecurityService(logger)

	// Test 3: Rate Limiter
	fmt.Println("\n=== Testing Rate Limiter ===")
	testRateLimiter(logger)

	// Test 4: Audit Logger
	fmt.Println("\n=== Testing Audit Logger ===")
	testAuditLogger(logger)

	fmt.Println("\n=== All Security Tests Completed ===")
}

func testAuthService(logger *logrus.Logger) {
	authService := services.NewChatAuthService("test-secret-key", logger)

	// Test token creation and validation
	token, err := authService.CreateAccessToken("test-user", "testuser", "test@example.com", "")
	if err != nil {
		log.Fatalf("Failed to create access token: %v", err)
	}
	fmt.Printf("✓ Created access token: %s...\n", token[:20])

	// Validate the token
	authContext, err := authService.ValidateToken(context.Background(), token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	fmt.Printf("✓ Token validated for user: %s\n", authContext.UserID)

	// Test role checking
	hasRole, err := authService.HasRole(context.Background(), "test-user", interfaces.RoleUser)
	if err != nil {
		log.Fatalf("Failed to check role: %v", err)
	}
	fmt.Printf("✓ User has role 'user': %v\n", hasRole)

	// Test permission checking
	hasPermission, err := authService.CheckPermission(context.Background(), "test-user", interfaces.PermissionChatRead)
	if err != nil {
		log.Fatalf("Failed to check permission: %v", err)
	}
	fmt.Printf("✓ User has chat:read permission: %v\n", hasPermission)

	// Test refresh token
	refreshToken, err := authService.CreateRefreshToken(context.Background(), "test-user")
	if err != nil {
		log.Fatalf("Failed to create refresh token: %v", err)
	}
	fmt.Printf("✓ Created refresh token: %s...\n", refreshToken.Token[:20])
}

func testSecurityService(logger *logrus.Logger) {
	rateLimiter := services.NewChatRateLimiter(logger)
	auditLogger := services.NewChatAuditLogger(logger)
	securityService := services.NewChatSecurityService(logger, "test-encryption-key", rateLimiter, auditLogger)

	// Test message validation
	validMessage := "Hello, this is a valid message"
	err := securityService.ValidateMessageContent(validMessage)
	if err != nil {
		log.Fatalf("Valid message failed validation: %v", err)
	}
	fmt.Printf("✓ Valid message passed validation\n")

	// Test message sanitization
	unsafeMessage := "<script>alert('xss')</script>Hello World"
	sanitized := securityService.SanitizeMessageContent(unsafeMessage)
	fmt.Printf("✓ Sanitized message: %s\n", sanitized)

	// Test content filtering
	filterResult, err := securityService.FilterContent("This is a test message")
	if err != nil {
		log.Fatalf("Failed to filter content: %v", err)
	}
	fmt.Printf("✓ Content filtering result - allowed: %v\n", filterResult.IsAllowed)

	// Test encryption/decryption
	sensitiveData := "This is sensitive information"
	encrypted, err := securityService.EncryptSensitiveData(sensitiveData)
	if err != nil {
		log.Fatalf("Failed to encrypt data: %v", err)
	}
	fmt.Printf("✓ Data encrypted successfully\n")

	decrypted, err := securityService.DecryptSensitiveData(encrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err)
	}
	if decrypted != sensitiveData {
		log.Fatalf("Decrypted data doesn't match original")
	}
	fmt.Printf("✓ Data decrypted successfully\n")

	// Test hashing
	hash := securityService.HashSensitiveData(sensitiveData)
	fmt.Printf("✓ Data hashed: %s...\n", hash[:20])
}

func testRateLimiter(logger *logrus.Logger) {
	rateLimiter := services.NewChatRateLimiter(logger)

	// Test message rate limiting
	ctx := context.Background()
	userID := "test-user"

	// First request should be allowed
	result, err := rateLimiter.AllowMessage(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to check message rate limit: %v", err)
	}
	fmt.Printf("✓ First message allowed: %v, remaining: %d\n", result.Allowed, result.Remaining)

	// Test connection rate limiting
	result, err = rateLimiter.AllowConnection(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to check connection rate limit: %v", err)
	}
	fmt.Printf("✓ Connection allowed: %v, remaining: %d\n", result.Allowed, result.Remaining)

	// Test session creation rate limiting
	result, err = rateLimiter.AllowSessionCreation(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to check session creation rate limit: %v", err)
	}
	fmt.Printf("✓ Session creation allowed: %v, remaining: %d\n", result.Allowed, result.Remaining)

	// Test custom rate limit
	result, err = rateLimiter.CheckLimit(ctx, "custom:test", 5, time.Minute)
	if err != nil {
		log.Fatalf("Failed to check custom rate limit: %v", err)
	}
	fmt.Printf("✓ Custom rate limit check: %v, remaining: %d\n", result.Allowed, result.Remaining)
}

func testAuditLogger(logger *logrus.Logger) {
	auditLogger := services.NewChatAuditLogger(logger)

	ctx := context.Background()
	userID := "test-user"
	sessionID := "test-session"

	// Test login logging
	metadata := map[string]interface{}{
		"ip_address": "127.0.0.1",
		"user_agent": "test-agent",
	}
	err := auditLogger.LogLogin(ctx, userID, true, metadata)
	if err != nil {
		log.Fatalf("Failed to log login: %v", err)
	}
	fmt.Printf("✓ Login event logged\n")

	// Test session creation logging
	err = auditLogger.LogSessionCreated(ctx, userID, sessionID, metadata)
	if err != nil {
		log.Fatalf("Failed to log session creation: %v", err)
	}
	fmt.Printf("✓ Session creation logged\n")

	// Test message sent logging
	messageID := "test-message"
	err = auditLogger.LogMessageSent(ctx, userID, sessionID, messageID, metadata)
	if err != nil {
		log.Fatalf("Failed to log message sent: %v", err)
	}
	fmt.Printf("✓ Message sent logged\n")

	// Test security violation logging
	violation := &interfaces.SecurityViolation{
		Type:        "content_violation",
		Severity:    "medium",
		Description: "Test security violation",
		Resource:    "chat_message",
		Action:      "send_message",
		IPAddress:   "127.0.0.1",
		UserAgent:   "test-agent",
		Metadata:    metadata,
	}
	err = auditLogger.LogSecurityViolation(ctx, userID, violation)
	if err != nil {
		log.Fatalf("Failed to log security violation: %v", err)
	}
	fmt.Printf("✓ Security violation logged\n")

	// Test audit log retrieval
	filters := &interfaces.AuditLogFilters{
		UserID: userID,
		Limit:  10,
	}
	logs, err := auditLogger.GetAuditLogs(ctx, filters)
	if err != nil {
		log.Fatalf("Failed to get audit logs: %v", err)
	}
	fmt.Printf("✓ Retrieved %d audit logs\n", len(logs))

	// Display audit log count
	count := auditLogger.GetAuditLogCount()
	fmt.Printf("✓ Total audit logs in memory: %d\n", count)
}
