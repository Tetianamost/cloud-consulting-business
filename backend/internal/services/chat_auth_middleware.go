package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ChatAuthMiddleware implements the ChatAuthMiddleware interface
type ChatAuthMiddleware struct {
	authService interfaces.ChatAuthService
	auditLogger interfaces.ChatAuditLogger
	logger      *logrus.Logger
}

// NewChatAuthMiddleware creates a new chat authentication middleware
func NewChatAuthMiddleware(
	authService interfaces.ChatAuthService,
	auditLogger interfaces.ChatAuditLogger,
	logger *logrus.Logger,
) *ChatAuthMiddleware {
	return &ChatAuthMiddleware{
		authService: authService,
		auditLogger: auditLogger,
		logger:      logger,
	}
}

// RequireAuth returns a middleware that requires authentication
func (m *ChatAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, err := m.authenticateRequest(c)
		if err != nil {
			m.handleAuthError(c, err, "authentication_required")
			return
		}

		// Set authentication context in Gin context
		m.setAuthContext(c, authContext)

		// Log successful authentication
		m.logAuthEvent(c, authContext.UserID, "authentication_success", true)

		c.Next()
	}
}

// RequireRole returns a middleware that requires a specific role
func (m *ChatAuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure authentication
		authContext, err := m.getOrAuthenticateContext(c)
		if err != nil {
			m.handleAuthError(c, err, "role_check_auth_failed")
			return
		}

		// Check if user has the required role
		hasRole, err := m.authService.HasRole(context.Background(), authContext.UserID, role)
		if err != nil {
			m.logger.WithError(err).WithFields(logrus.Fields{
				"user_id": authContext.UserID,
				"role":    role,
			}).Error("Failed to check user role")

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to verify user role",
			})
			c.Abort()
			return
		}

		if !hasRole {
			m.logger.WithFields(logrus.Fields{
				"user_id":       authContext.UserID,
				"required_role": role,
				"user_roles":    authContext.Roles,
			}).Warn("User lacks required role")

			// Log unauthorized access attempt
			m.logAuthEvent(c, authContext.UserID, "insufficient_role", false)

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Insufficient role privileges",
				"code":    interfaces.ErrCodeInsufficientRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission returns a middleware that requires a specific permission
func (m *ChatAuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure authentication
		authContext, err := m.getOrAuthenticateContext(c)
		if err != nil {
			m.handleAuthError(c, err, "permission_check_auth_failed")
			return
		}

		// Check if user has the required permission
		hasPermission, err := m.authService.CheckPermission(context.Background(), authContext.UserID, permission)
		if err != nil {
			m.logger.WithError(err).WithFields(logrus.Fields{
				"user_id":    authContext.UserID,
				"permission": permission,
			}).Error("Failed to check user permission")

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to verify user permission",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			m.logger.WithFields(logrus.Fields{
				"user_id":             authContext.UserID,
				"required_permission": permission,
				"user_permissions":    authContext.Permissions,
			}).Warn("User lacks required permission")

			// Log unauthorized access attempt
			m.logAuthEvent(c, authContext.UserID, "insufficient_permission", false)

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Insufficient permissions",
				"code":    interfaces.ErrCodeInsufficientPermission,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSessionAccess returns a middleware that requires access to a specific session
func (m *ChatAuthMiddleware) RequireSessionAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure authentication
		authContext, err := m.getOrAuthenticateContext(c)
		if err != nil {
			m.handleAuthError(c, err, "session_access_auth_failed")
			return
		}

		// Get session ID from URL parameter or request body
		sessionID := c.Param("sessionId")
		if sessionID == "" {
			sessionID = c.Param("id")
		}

		if sessionID == "" {
			// Try to get from request body for POST requests
			var requestBody struct {
				SessionID string `json:"session_id"`
			}
			if err := c.ShouldBindJSON(&requestBody); err == nil {
				sessionID = requestBody.SessionID
			}
		}

		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Session ID is required",
			})
			c.Abort()
			return
		}

		// Check if user can access this session
		err = m.authService.AuthorizeSessionAccess(context.Background(), authContext.UserID, sessionID)
		if err != nil {
			m.logger.WithError(err).WithFields(logrus.Fields{
				"user_id":    authContext.UserID,
				"session_id": sessionID,
			}).Warn("User denied access to session")

			// Log unauthorized access attempt
			metadata := map[string]interface{}{
				"session_id": sessionID,
				"ip_address": c.ClientIP(),
				"user_agent": c.GetHeader("User-Agent"),
			}
			if m.auditLogger != nil {
				m.auditLogger.LogUnauthorizedAccess(context.Background(), authContext.UserID, "chat_session", metadata)
			}

			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Access denied to session",
				"code":    interfaces.ErrCodeUnauthorizedAccess,
			})
			c.Abort()
			return
		}

		// Set session ID in context for use by handlers
		c.Set("session_id", sessionID)

		c.Next()
	}
}

// OptionalAuth returns a middleware that optionally authenticates requests
func (m *ChatAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, err := m.authenticateRequest(c)
		if err != nil {
			// For optional auth, we don't abort on auth failure
			// Just log the attempt and continue
			m.logger.WithError(err).Debug("Optional authentication failed")
		} else {
			// Set authentication context if successful
			m.setAuthContext(c, authContext)
			m.logAuthEvent(c, authContext.UserID, "optional_authentication_success", true)
		}

		c.Next()
	}
}

// Helper methods

// authenticateRequest authenticates a request and returns the auth context
func (m *ChatAuthMiddleware) authenticateRequest(c *gin.Context) (*interfaces.ChatAuthContext, error) {
	// Extract token from Authorization header
	token, err := m.extractToken(c)
	if err != nil {
		return nil, err
	}

	// Validate token
	authContext, err := m.authService.ValidateToken(context.Background(), token)
	if err != nil {
		return nil, err
	}

	return authContext, nil
}

// extractToken extracts the JWT token from the request
func (m *ChatAuthMiddleware) extractToken(c *gin.Context) (string, error) {
	// Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
			return authHeader[7:], nil
		}
		return "", interfaces.SecurityError{
			Code:    interfaces.ErrCodeInvalidToken,
			Message: "Invalid authorization header format",
		}
	}

	// Try query parameter (for WebSocket connections)
	token := c.Query("token")
	if token != "" {
		return token, nil
	}

	return "", interfaces.SecurityError{
		Code:    interfaces.ErrCodeInvalidToken,
		Message: "No authentication token provided",
	}
}

// getOrAuthenticateContext gets existing auth context or authenticates the request
func (m *ChatAuthMiddleware) getOrAuthenticateContext(c *gin.Context) (*interfaces.ChatAuthContext, error) {
	// Try to get existing auth context from Gin context
	if authContextInterface, exists := c.Get("auth_context"); exists {
		if authContext, ok := authContextInterface.(*interfaces.ChatAuthContext); ok {
			return authContext, nil
		}
	}

	// If not found, authenticate the request
	authContext, err := m.authenticateRequest(c)
	if err != nil {
		return nil, err
	}

	// Set the context for future use
	m.setAuthContext(c, authContext)

	return authContext, nil
}

// setAuthContext sets the authentication context in Gin context
func (m *ChatAuthMiddleware) setAuthContext(c *gin.Context, authContext *interfaces.ChatAuthContext) {
	c.Set("auth_context", authContext)
	c.Set("user_id", authContext.UserID)
	c.Set("username", authContext.Username)
	c.Set("user_roles", authContext.Roles)
	c.Set("user_permissions", authContext.Permissions)
}

// handleAuthError handles authentication errors
func (m *ChatAuthMiddleware) handleAuthError(c *gin.Context, err error, eventType string) {
	var statusCode int
	var errorCode string
	var message string

	// Handle different types of security errors
	if secErr, ok := err.(interfaces.SecurityError); ok {
		switch secErr.Code {
		case interfaces.ErrCodeExpiredToken:
			statusCode = http.StatusUnauthorized
			errorCode = secErr.Code
			message = "Token has expired"
		case interfaces.ErrCodeRevokedToken:
			statusCode = http.StatusUnauthorized
			errorCode = secErr.Code
			message = "Token has been revoked"
		case interfaces.ErrCodeInvalidToken:
			statusCode = http.StatusUnauthorized
			errorCode = secErr.Code
			message = "Invalid authentication token"
		default:
			statusCode = http.StatusUnauthorized
			errorCode = secErr.Code
			message = secErr.Message
		}
	} else {
		statusCode = http.StatusUnauthorized
		errorCode = "AUTHENTICATION_FAILED"
		message = "Authentication failed"
	}

	// Log authentication failure
	userID := "unknown"
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		// Try to extract user ID from token for logging (even if invalid)
		if token, tokenErr := m.extractToken(c); tokenErr == nil {
			if authContext, validateErr := m.authService.ValidateToken(context.Background(), token); validateErr == nil {
				userID = authContext.UserID
			}
		}
	}

	m.logAuthEvent(c, userID, eventType, false)

	m.logger.WithError(err).WithFields(logrus.Fields{
		"user_id":    userID,
		"event_type": eventType,
		"ip_address": c.ClientIP(),
		"user_agent": c.GetHeader("User-Agent"),
	}).Warn("Authentication failed")

	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
		"code":    errorCode,
	})
	c.Abort()
}

// logAuthEvent logs authentication events
func (m *ChatAuthMiddleware) logAuthEvent(c *gin.Context, userID string, eventType string, success bool) {
	if m.auditLogger == nil {
		return
	}

	metadata := map[string]interface{}{
		"event_type": eventType,
		"ip_address": c.ClientIP(),
		"user_agent": c.GetHeader("User-Agent"),
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
	}

	// Log based on event type
	switch eventType {
	case "authentication_success", "optional_authentication_success":
		m.auditLogger.LogLogin(context.Background(), userID, success, metadata)
	case "authentication_required", "role_check_auth_failed", "permission_check_auth_failed", "session_access_auth_failed":
		m.auditLogger.LogLogin(context.Background(), userID, success, metadata)
	case "insufficient_role", "insufficient_permission":
		m.auditLogger.LogUnauthorizedAccess(context.Background(), userID, "role_or_permission", metadata)
	}
}

// GetAuthContext is a helper function to get auth context from Gin context
func GetAuthContext(c *gin.Context) (*interfaces.ChatAuthContext, bool) {
	if authContextInterface, exists := c.Get("auth_context"); exists {
		if authContext, ok := authContextInterface.(*interfaces.ChatAuthContext); ok {
			return authContext, true
		}
	}
	return nil, false
}

// GetUserID is a helper function to get user ID from Gin context
func GetUserID(c *gin.Context) (string, bool) {
	if userIDInterface, exists := c.Get("user_id"); exists {
		if userID, ok := userIDInterface.(string); ok {
			return userID, true
		}
	}
	return "", false
}

// RequirePermissions returns a middleware that requires multiple permissions (AND logic)
func (m *ChatAuthMiddleware) RequirePermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure authentication
		authContext, err := m.getOrAuthenticateContext(c)
		if err != nil {
			m.handleAuthError(c, err, "permissions_check_auth_failed")
			return
		}

		// Check each required permission
		for _, permission := range permissions {
			hasPermission, err := m.authService.CheckPermission(context.Background(), authContext.UserID, permission)
			if err != nil {
				m.logger.WithError(err).WithFields(logrus.Fields{
					"user_id":    authContext.UserID,
					"permission": permission,
				}).Error("Failed to check user permission")

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "Failed to verify user permissions",
				})
				c.Abort()
				return
			}

			if !hasPermission {
				m.logger.WithFields(logrus.Fields{
					"user_id":             authContext.UserID,
					"required_permission": permission,
					"user_permissions":    authContext.Permissions,
				}).Warn("User lacks required permission")

				// Log unauthorized access attempt
				m.logAuthEvent(c, authContext.UserID, "insufficient_permission", false)

				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"error":   fmt.Sprintf("Missing required permission: %s", permission),
					"code":    interfaces.ErrCodeInsufficientPermission,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
