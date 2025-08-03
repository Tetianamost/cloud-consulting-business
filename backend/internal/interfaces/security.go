package interfaces

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// ChatAuthService defines the interface for chat authentication and authorization
type ChatAuthService interface {
	// JWT Token Management
	ValidateToken(ctx context.Context, token string) (*ChatAuthContext, error)
	RefreshToken(ctx context.Context, token string) (*TokenPair, error)
	RevokeToken(ctx context.Context, token string) error
	IsTokenRevoked(ctx context.Context, token string) (bool, error)

	// Session-based Authorization
	AuthorizeSessionAccess(ctx context.Context, userID, sessionID string) error
	AuthorizeMessageAccess(ctx context.Context, userID, messageID string) error
	CheckPermission(ctx context.Context, userID string, permission string) (bool, error)

	// Role-based Access Control
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	HasRole(ctx context.Context, userID string, role string) (bool, error)
	GetRolePermissions(ctx context.Context, role string) ([]string, error)

	// Token Refresh for Long-running Sessions
	CreateRefreshToken(ctx context.Context, userID string) (*RefreshToken, error)
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*ChatAuthContext, error)
	ExtendSessionToken(ctx context.Context, sessionID string, duration time.Duration) error
}

// ChatSecurityService defines the interface for chat data protection and security
type ChatSecurityService interface {
	// Input Validation and Sanitization
	ValidateMessageContent(content string) error
	SanitizeMessageContent(content string) string
	ValidateSessionData(session interface{}) error

	// Content Filtering and Moderation
	FilterContent(content string) (*ContentFilterResult, error)
	ModerateMessage(content string) (*ModerationResult, error)
	IsContentAllowed(content string) (bool, error)

	// Rate Limiting
	CheckRateLimit(ctx context.Context, userID string, action string) (*RateLimitResult, error)
	IncrementRateLimit(ctx context.Context, userID string, action string) error
	ResetRateLimit(ctx context.Context, userID string, action string) error

	// Audit Logging
	LogSecurityEvent(ctx context.Context, event *SecurityEvent) error
	LogAuthenticationAttempt(ctx context.Context, userID string, success bool, details map[string]interface{}) error
	LogDataAccess(ctx context.Context, userID string, resource string, action string) error

	// Encryption and Data Protection
	EncryptSensitiveData(data string) (string, error)
	DecryptSensitiveData(encryptedData string) (string, error)
	HashSensitiveData(data string) string
}

// ChatRateLimiter defines the interface for rate limiting chat operations
type ChatRateLimiter interface {
	// Message Rate Limiting
	AllowMessage(ctx context.Context, userID string) (*RateLimitResult, error)
	AllowConnection(ctx context.Context, userID string) (*RateLimitResult, error)
	AllowSessionCreation(ctx context.Context, userID string) (*RateLimitResult, error)

	// Custom Rate Limits
	CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (*RateLimitResult, error)
	SetCustomLimit(ctx context.Context, userID string, action string, limit int, window time.Duration) error
	GetLimitInfo(ctx context.Context, key string) (*RateLimitInfo, error)

	// Rate Limit Management
	ResetUserLimits(ctx context.Context, userID string) error
	GetUserLimitStatus(ctx context.Context, userID string) (map[string]*RateLimitInfo, error)
}

// ChatAuditLogger defines the interface for security audit logging
type ChatAuditLogger interface {
	// Authentication Events
	LogLogin(ctx context.Context, userID string, success bool, metadata map[string]interface{}) error
	LogLogout(ctx context.Context, userID string, metadata map[string]interface{}) error
	LogTokenRefresh(ctx context.Context, userID string, metadata map[string]interface{}) error

	// Session Events
	LogSessionCreated(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error
	LogSessionAccessed(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error
	LogSessionDeleted(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error

	// Message Events
	LogMessageSent(ctx context.Context, userID, sessionID, messageID string, metadata map[string]interface{}) error
	LogMessageAccessed(ctx context.Context, userID, messageID string, metadata map[string]interface{}) error

	// Security Events
	LogSecurityViolation(ctx context.Context, userID string, violation *SecurityViolation) error
	LogRateLimitExceeded(ctx context.Context, userID string, action string, metadata map[string]interface{}) error
	LogUnauthorizedAccess(ctx context.Context, userID string, resource string, metadata map[string]interface{}) error

	// Query Audit Logs
	GetAuditLogs(ctx context.Context, filters *AuditLogFilters) ([]*AuditLog, error)
	GetUserAuditLogs(ctx context.Context, userID string, limit int) ([]*AuditLog, error)
}

// Supporting types for security interfaces

// ChatAuthContext represents the authentication context for chat operations
type ChatAuthContext struct {
	UserID      string                 `json:"user_id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Roles       []string               `json:"roles"`
	Permissions []string               `json:"permissions"`
	TokenType   string                 `json:"token_type"`
	IssuedAt    time.Time              `json:"issued_at"`
	ExpiresAt   time.Time              `json:"expires_at"`
	SessionID   string                 `json:"session_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	IsRevoked bool      `json:"is_revoked"`
}

// ContentFilterResult represents the result of content filtering
type ContentFilterResult struct {
	IsAllowed    bool     `json:"is_allowed"`
	FilteredText string   `json:"filtered_text"`
	Violations   []string `json:"violations"`
	Confidence   float64  `json:"confidence"`
	Categories   []string `json:"categories"`
}

// ModerationResult represents the result of content moderation
type ModerationResult struct {
	IsApproved   bool                   `json:"is_approved"`
	Confidence   float64                `json:"confidence"`
	Categories   []string               `json:"categories"`
	Reasons      []string               `json:"reasons"`
	Metadata     map[string]interface{} `json:"metadata"`
	ReviewNeeded bool                   `json:"review_needed"`
}

// SecurityEvent represents a security event for audit logging
type SecurityEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	UserID      string                 `json:"user_id"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Result      string                 `json:"result"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
}

// SecurityViolation represents a security violation
type SecurityViolation struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Result    string                 `json:"result"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
	SessionID string                 `json:"session_id,omitempty"`
	MessageID string                 `json:"message_id,omitempty"`
}

// AuditLogFilters represents filters for querying audit logs
type AuditLogFilters struct {
	UserID    string    `json:"user_id,omitempty"`
	Action    string    `json:"action,omitempty"`
	Resource  string    `json:"resource,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// JWT Claims for chat authentication
type ChatJWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	SessionID   string   `json:"session_id,omitempty"`
	TokenType   string   `json:"token_type"`
	jwt.RegisteredClaims
}

// Middleware interfaces for chat security
type ChatAuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
	RequireRole(role string) gin.HandlerFunc
	RequirePermission(permission string) gin.HandlerFunc
	RequireSessionAccess() gin.HandlerFunc
	OptionalAuth() gin.HandlerFunc
}

// Error types for security operations
type SecurityError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e SecurityError) Error() string {
	return e.Message
}

// Security error codes
const (
	ErrCodeInvalidToken           = "INVALID_TOKEN"
	ErrCodeExpiredToken           = "EXPIRED_TOKEN"
	ErrCodeRevokedToken           = "REVOKED_TOKEN"
	ErrCodeInsufficientRole       = "INSUFFICIENT_ROLE"
	ErrCodeInsufficientPermission = "INSUFFICIENT_PERMISSION"
	ErrCodeUnauthorizedAccess     = "UNAUTHORIZED_ACCESS"
	ErrCodeRateLimitExceeded      = "RATE_LIMIT_EXCEEDED"
	ErrCodeContentViolation       = "CONTENT_VIOLATION"
	ErrCodeValidationFailed       = "VALIDATION_FAILED"
	ErrCodeEncryptionFailed       = "ENCRYPTION_FAILED"
	ErrCodeDecryptionFailed       = "DECRYPTION_FAILED"
)

// Permission constants for chat operations
const (
	PermissionChatRead      = "chat:read"
	PermissionChatWrite     = "chat:write"
	PermissionChatDelete    = "chat:delete"
	PermissionSessionCreate = "session:create"
	PermissionSessionRead   = "session:read"
	PermissionSessionUpdate = "session:update"
	PermissionSessionDelete = "session:delete"
	PermissionAdminAccess   = "admin:access"
	PermissionAuditRead     = "audit:read"
)

// Role constants for chat system
const (
	RoleAdmin      = "admin"
	RoleConsultant = "consultant"
	RoleUser       = "user"
	RoleGuest      = "guest"
)
