package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ChatAuthService implements the ChatAuthService interface
type ChatAuthService struct {
	jwtSecret       []byte
	tokenExpiry     time.Duration
	refreshExpiry   time.Duration
	logger          *logrus.Logger
	revokedTokens   map[string]time.Time // In production, use Redis or database
	refreshTokens   map[string]*interfaces.RefreshToken
	userRoles       map[string][]string // In production, use database
	rolePermissions map[string][]string
	mutex           sync.RWMutex
}

// NewChatAuthService creates a new chat authentication service
func NewChatAuthService(jwtSecret string, logger *logrus.Logger) *ChatAuthService {
	service := &ChatAuthService{
		jwtSecret:       []byte(jwtSecret),
		tokenExpiry:     24 * time.Hour,
		refreshExpiry:   7 * 24 * time.Hour,
		logger:          logger,
		revokedTokens:   make(map[string]time.Time),
		refreshTokens:   make(map[string]*interfaces.RefreshToken),
		userRoles:       make(map[string][]string),
		rolePermissions: make(map[string][]string),
	}

	// Initialize default roles and permissions
	service.initializeDefaultRolesAndPermissions()

	// Start cleanup routine for expired tokens
	go service.cleanupExpiredTokens()

	return service
}

// initializeDefaultRolesAndPermissions sets up default roles and permissions
func (s *ChatAuthService) initializeDefaultRolesAndPermissions() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Define role permissions
	s.rolePermissions[interfaces.RoleAdmin] = []string{
		interfaces.PermissionChatRead,
		interfaces.PermissionChatWrite,
		interfaces.PermissionChatDelete,
		interfaces.PermissionSessionCreate,
		interfaces.PermissionSessionRead,
		interfaces.PermissionSessionUpdate,
		interfaces.PermissionSessionDelete,
		interfaces.PermissionAdminAccess,
		interfaces.PermissionAuditRead,
	}

	s.rolePermissions[interfaces.RoleConsultant] = []string{
		interfaces.PermissionChatRead,
		interfaces.PermissionChatWrite,
		interfaces.PermissionSessionCreate,
		interfaces.PermissionSessionRead,
		interfaces.PermissionSessionUpdate,
	}

	s.rolePermissions[interfaces.RoleUser] = []string{
		interfaces.PermissionChatRead,
		interfaces.PermissionChatWrite,
		interfaces.PermissionSessionCreate,
		interfaces.PermissionSessionRead,
	}

	s.rolePermissions[interfaces.RoleGuest] = []string{
		interfaces.PermissionChatRead,
	}

	// Set default user roles (in production, this would come from database)
	s.userRoles["admin"] = []string{interfaces.RoleAdmin}
	s.userRoles["consultant"] = []string{interfaces.RoleConsultant}
}

// ValidateToken validates a JWT token and returns the authentication context
func (s *ChatAuthService) ValidateToken(ctx context.Context, tokenString string) (*interfaces.ChatAuthContext, error) {
	// Check if token is revoked
	if revoked, err := s.IsTokenRevoked(ctx, tokenString); err != nil {
		return nil, fmt.Errorf("failed to check token revocation: %w", err)
	} else if revoked {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeRevokedToken,
			Message: "Token has been revoked",
		}
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &interfaces.ChatJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	var claims *interfaces.ChatJWTClaims

	if err != nil {
		// Try parsing with MapClaims for backward compatibility
		s.logger.WithError(err).Debug("Failed to parse with ChatJWTClaims, trying MapClaims for backward compatibility")

		mapToken, mapErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.jwtSecret, nil
		})

		if mapErr != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					return nil, interfaces.SecurityError{
						Code:    interfaces.ErrCodeExpiredToken,
						Message: "Token has expired",
					}
				}
			}
			return nil, interfaces.SecurityError{
				Code:    interfaces.ErrCodeInvalidToken,
				Message: "Invalid token",
				Details: err.Error(),
			}
		}

		if mapClaims, ok := mapToken.Claims.(jwt.MapClaims); ok && mapToken.Valid {
			// Convert MapClaims to ChatJWTClaims
			claims = &interfaces.ChatJWTClaims{
				UserID:    getStringFromClaims(mapClaims, "user_id"),
				Username:  getStringFromClaims(mapClaims, "username"),
				Email:     getStringFromClaims(mapClaims, "email"),
				TokenType: getStringFromClaims(mapClaims, "token_type"),
			}

			// Handle roles - could be string or []string
			if roleInterface, exists := mapClaims["role"]; exists {
				if roleStr, ok := roleInterface.(string); ok {
					claims.Roles = []string{roleStr}
				}
			}
			if rolesInterface, exists := mapClaims["roles"]; exists {
				if rolesSlice, ok := rolesInterface.([]interface{}); ok {
					roles := make([]string, len(rolesSlice))
					for i, role := range rolesSlice {
						if roleStr, ok := role.(string); ok {
							roles[i] = roleStr
						}
					}
					claims.Roles = roles
				}
			}

			// Handle standard JWT claims
			if exp, exists := mapClaims["exp"]; exists {
				if expFloat, ok := exp.(float64); ok {
					claims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(expFloat), 0))
				}
			}
			if iat, exists := mapClaims["iat"]; exists {
				if iatFloat, ok := iat.(float64); ok {
					claims.IssuedAt = jwt.NewNumericDate(time.Unix(int64(iatFloat), 0))
				}
			}

			s.logger.Debug("Successfully parsed token using MapClaims for backward compatibility")
		} else {
			return nil, interfaces.SecurityError{
				Code:    interfaces.ErrCodeInvalidToken,
				Message: "Invalid token claims",
			}
		}
	} else {
		// Extract claims from ChatJWTClaims
		var ok bool
		claims, ok = token.Claims.(*interfaces.ChatJWTClaims)
		if !ok || !token.Valid {
			return nil, interfaces.SecurityError{
				Code:    interfaces.ErrCodeInvalidToken,
				Message: "Invalid token claims",
			}
		}
	}

	// Create authentication context with safe time handling
	var issuedAt, expiresAt time.Time
	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Time
	}
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	// Handle backward compatibility: if UserID is empty but Username is set, use Username as UserID
	userID := claims.UserID
	if userID == "" && claims.Username != "" {
		userID = claims.Username
		s.logger.WithField("username", claims.Username).Debug("Using username as user_id for backward compatibility")
	}

	authContext := &interfaces.ChatAuthContext{
		UserID:      userID,
		Username:    claims.Username,
		Email:       claims.Email,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		TokenType:   claims.TokenType,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		SessionID:   claims.SessionID,
		Metadata:    make(map[string]interface{}),
	}

	return authContext, nil
}

// RefreshToken creates a new access token using a refresh token
func (s *ChatAuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*interfaces.TokenPair, error) {
	s.mutex.RLock()
	refreshToken, exists := s.refreshTokens[refreshTokenString]
	s.mutex.RUnlock()

	if !exists {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeInvalidToken,
			Message: "Invalid refresh token",
		}
	}

	if refreshToken.IsRevoked {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeRevokedToken,
			Message: "Refresh token has been revoked",
		}
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeExpiredToken,
			Message: "Refresh token has expired",
		}
	}

	// Get user roles and permissions
	roles := s.getUserRoles(refreshToken.UserID)
	permissions := s.getUserPermissions(refreshToken.UserID)

	// Create new access token
	accessToken, err := s.createAccessToken(refreshToken.UserID, refreshToken.UserID, "", roles, permissions, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Create new refresh token
	newRefreshToken, err := s.CreateRefreshToken(ctx, refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Revoke old refresh token
	s.mutex.Lock()
	refreshToken.IsRevoked = true
	s.mutex.Unlock()

	return &interfaces.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenExpiry.Seconds()),
		ExpiresAt:    time.Now().Add(s.tokenExpiry),
	}, nil
}

// RevokeToken revokes a JWT token
func (s *ChatAuthService) RevokeToken(ctx context.Context, tokenString string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Add token to revoked list with expiration time
	s.revokedTokens[tokenString] = time.Now().Add(s.tokenExpiry)

	s.logger.WithField("token_hash", s.hashToken(tokenString)).Info("Token revoked")
	return nil
}

// IsTokenRevoked checks if a token has been revoked
func (s *ChatAuthService) IsTokenRevoked(ctx context.Context, tokenString string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, revoked := s.revokedTokens[tokenString]
	return revoked, nil
}

// AuthorizeSessionAccess checks if a user can access a specific session
func (s *ChatAuthService) AuthorizeSessionAccess(ctx context.Context, userID, sessionID string) error {
	// In a real implementation, this would check the database to ensure
	// the session belongs to the user or the user has admin privileges

	// For now, we'll implement basic logic
	roles := s.getUserRoles(userID)

	// Admins can access any session
	for _, role := range roles {
		if role == interfaces.RoleAdmin {
			return nil
		}
	}

	// For non-admins, we would need to check session ownership
	// This is a simplified implementation
	return nil
}

// AuthorizeMessageAccess checks if a user can access a specific message
func (s *ChatAuthService) AuthorizeMessageAccess(ctx context.Context, userID, messageID string) error {
	// Similar to session access, this would check message ownership
	// or admin privileges in a real implementation

	roles := s.getUserRoles(userID)

	// Admins can access any message
	for _, role := range roles {
		if role == interfaces.RoleAdmin {
			return nil
		}
	}

	return nil
}

// CheckPermission checks if a user has a specific permission
func (s *ChatAuthService) CheckPermission(ctx context.Context, userID string, permission string) (bool, error) {
	permissions := s.getUserPermissions(userID)

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}

// GetUserRoles returns the roles for a user
func (s *ChatAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return s.getUserRoles(userID), nil
}

// HasRole checks if a user has a specific role
func (s *ChatAuthService) HasRole(ctx context.Context, userID string, role string) (bool, error) {
	roles := s.getUserRoles(userID)

	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}

	return false, nil
}

// GetRolePermissions returns the permissions for a role
func (s *ChatAuthService) GetRolePermissions(ctx context.Context, role string) ([]string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	permissions, exists := s.rolePermissions[role]
	if !exists {
		return []string{}, nil
	}

	// Return a copy to prevent modification
	result := make([]string, len(permissions))
	copy(result, permissions)
	return result, nil
}

// CreateRefreshToken creates a new refresh token for a user
func (s *ChatAuthService) CreateRefreshToken(ctx context.Context, userID string) (*interfaces.RefreshToken, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	tokenString := hex.EncodeToString(tokenBytes)

	refreshToken := &interfaces.RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.refreshExpiry),
		CreatedAt: time.Now(),
		IsRevoked: false,
	}

	s.mutex.Lock()
	s.refreshTokens[tokenString] = refreshToken
	s.mutex.Unlock()

	return refreshToken, nil
}

// ValidateRefreshToken validates a refresh token and returns auth context
func (s *ChatAuthService) ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*interfaces.ChatAuthContext, error) {
	s.mutex.RLock()
	refreshToken, exists := s.refreshTokens[refreshTokenString]
	s.mutex.RUnlock()

	if !exists {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeInvalidToken,
			Message: "Invalid refresh token",
		}
	}

	if refreshToken.IsRevoked {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeRevokedToken,
			Message: "Refresh token has been revoked",
		}
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, interfaces.SecurityError{
			Code:    interfaces.ErrCodeExpiredToken,
			Message: "Refresh token has expired",
		}
	}

	// Create authentication context
	roles := s.getUserRoles(refreshToken.UserID)
	permissions := s.getUserPermissions(refreshToken.UserID)

	authContext := &interfaces.ChatAuthContext{
		UserID:      refreshToken.UserID,
		Username:    refreshToken.UserID, // Using userID as username for now
		Roles:       roles,
		Permissions: permissions,
		TokenType:   "refresh",
		IssuedAt:    refreshToken.CreatedAt,
		ExpiresAt:   refreshToken.ExpiresAt,
		Metadata:    make(map[string]interface{}),
	}

	return authContext, nil
}

// ExtendSessionToken extends the expiration of a session token
func (s *ChatAuthService) ExtendSessionToken(ctx context.Context, sessionID string, duration time.Duration) error {
	// In a real implementation, this would update the token expiration
	// For now, we'll just log the operation
	s.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"duration":   duration,
	}).Info("Session token extended")

	return nil
}

// Helper methods

// getUserRoles returns the roles for a user
func (s *ChatAuthService) getUserRoles(userID string) []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	roles, exists := s.userRoles[userID]
	if !exists {
		// Default role for unknown users
		return []string{interfaces.RoleUser}
	}

	// Return a copy to prevent modification
	result := make([]string, len(roles))
	copy(result, roles)
	return result
}

// getUserPermissions returns all permissions for a user based on their roles
func (s *ChatAuthService) getUserPermissions(userID string) []string {
	roles := s.getUserRoles(userID)
	permissionSet := make(map[string]bool)

	s.mutex.RLock()
	for _, role := range roles {
		if permissions, exists := s.rolePermissions[role]; exists {
			for _, permission := range permissions {
				permissionSet[permission] = true
			}
		}
	}
	s.mutex.RUnlock()

	// Convert set to slice
	permissions := make([]string, 0, len(permissionSet))
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	return permissions
}

// createAccessToken creates a new JWT access token
func (s *ChatAuthService) createAccessToken(userID, username, email string, roles, permissions []string, sessionID string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.tokenExpiry)

	claims := &interfaces.ChatJWTClaims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		SessionID:   sessionID,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "chat-service",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// hashToken creates a hash of the token for logging (security)
func (s *ChatAuthService) hashToken(token string) string {
	if len(token) < 10 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// cleanupExpiredTokens periodically removes expired revoked tokens
func (s *ChatAuthService) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mutex.Lock()
			now := time.Now()

			// Clean up expired revoked tokens
			for token, expiry := range s.revokedTokens {
				if now.After(expiry) {
					delete(s.revokedTokens, token)
				}
			}

			// Clean up expired refresh tokens
			for token, refreshToken := range s.refreshTokens {
				if now.After(refreshToken.ExpiresAt) {
					delete(s.refreshTokens, token)
				}
			}

			s.mutex.Unlock()

			s.logger.Debug("Cleaned up expired tokens")
		}
	}
}

// SetUserRoles sets the roles for a user (for testing/admin purposes)
func (s *ChatAuthService) SetUserRoles(userID string, roles []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.userRoles[userID] = roles
}

// getStringFromClaims safely extracts a string value from JWT MapClaims
func getStringFromClaims(claims jwt.MapClaims, key string) string {
	if value, exists := claims[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// CreateAccessToken creates a new access token for a user (public method for login)
func (s *ChatAuthService) CreateAccessToken(userID, username, email string, sessionID string) (string, error) {
	roles := s.getUserRoles(userID)
	permissions := s.getUserPermissions(userID)
	return s.createAccessToken(userID, username, email, roles, permissions, sessionID)
}
