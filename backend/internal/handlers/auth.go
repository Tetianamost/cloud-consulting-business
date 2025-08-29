package handlers

import (
	"net/http"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	jwtSecret string
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		jwtSecret: jwtSecret,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// For demo purposes, use a simple hardcoded admin/password check
	// In a real application, this would use a secure password hash comparison
	if req.Username != "admin" || req.Password != "cloudadmin" {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}

	// Create a JWT token with ChatJWTClaims structure
	now := time.Now()
	expiresAt := now.Add(time.Hour * 24)

	claims := &interfaces.ChatJWTClaims{
		UserID:      req.Username, // Use username as UserID for now
		Username:    req.Username,
		Email:       "", // Empty for now
		Roles:       []string{"admin"},
		Permissions: []string{}, // Will be populated by ChatAuthService
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "auth-service",
			Subject:   req.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   tokenString,
	})
}

// AuthMiddleware creates a middleware for JWT authentication
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid authorization format, expected 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := authHeader[7:]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(h.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Set the username and role in the context
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
