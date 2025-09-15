package handlers

import (
	"net/http"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/gin-gonic/gin"
)

// ChatConfigResponse represents the chat configuration response
type ChatConfigResponse struct {
	Mode                    string `json:"mode"`
	EnableWebSocketFallback bool   `json:"enable_websocket_fallback"`
	WebSocketTimeout        int    `json:"websocket_timeout"`
	PollingInterval         int    `json:"polling_interval"`
	MaxReconnectAttempts    int    `json:"max_reconnect_attempts"`
	FallbackDelay           int    `json:"fallback_delay"`
}

// ChatConfigUpdateRequest represents the request to update chat configuration
type ChatConfigUpdateRequest struct {
	Mode                    *string `json:"mode,omitempty"`
	EnableWebSocketFallback *bool   `json:"enable_websocket_fallback,omitempty"`
	WebSocketTimeout        *int    `json:"websocket_timeout,omitempty"`
	PollingInterval         *int    `json:"polling_interval,omitempty"`
	MaxReconnectAttempts    *int    `json:"max_reconnect_attempts,omitempty"`
	FallbackDelay           *int    `json:"fallback_delay,omitempty"`
}

// ChatConfigHandler handles chat configuration endpoints
type ChatConfigHandler struct {
	config *config.Config
}

// NewChatConfigHandler creates a new chat configuration handler
func NewChatConfigHandler(cfg *config.Config) *ChatConfigHandler {
	return &ChatConfigHandler{
		config: cfg,
	}
}

// GetChatConfig returns the current chat configuration
func (h *ChatConfigHandler) GetChatConfig(c *gin.Context) {
	response := ChatConfigResponse{
		Mode:                    h.config.Chat.Mode,
		EnableWebSocketFallback: h.config.Chat.EnableWebSocketFallback,
		WebSocketTimeout:        h.config.Chat.WebSocketTimeout,
		PollingInterval:         h.config.Chat.PollingInterval,
		MaxReconnectAttempts:    h.config.Chat.MaxReconnectAttempts,
		FallbackDelay:           h.config.Chat.FallbackDelay,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"config":  response,
	})
}

// UpdateChatConfig updates the chat configuration (admin only)
func (h *ChatConfigHandler) UpdateChatConfig(c *gin.Context) {
	var request ChatConfigUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate mode if provided
	if request.Mode != nil {
		mode := *request.Mode
		if mode != "websocket" && mode != "polling" && mode != "auto" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid mode. Must be 'websocket', 'polling', or 'auto'",
			})
			return
		}
		h.config.Chat.Mode = mode
	}

	// Update other fields if provided
	if request.EnableWebSocketFallback != nil {
		h.config.Chat.EnableWebSocketFallback = *request.EnableWebSocketFallback
	}
	if request.WebSocketTimeout != nil {
		if *request.WebSocketTimeout < 1 || *request.WebSocketTimeout > 60 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "WebSocket timeout must be between 1 and 60 seconds",
			})
			return
		}
		h.config.Chat.WebSocketTimeout = *request.WebSocketTimeout
	}
	if request.PollingInterval != nil {
		if *request.PollingInterval < 1000 || *request.PollingInterval > 30000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Polling interval must be between 1000 and 30000 milliseconds",
			})
			return
		}
		h.config.Chat.PollingInterval = *request.PollingInterval
	}
	if request.MaxReconnectAttempts != nil {
		if *request.MaxReconnectAttempts < 1 || *request.MaxReconnectAttempts > 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Max reconnect attempts must be between 1 and 10",
			})
			return
		}
		h.config.Chat.MaxReconnectAttempts = *request.MaxReconnectAttempts
	}
	if request.FallbackDelay != nil {
		if *request.FallbackDelay < 1000 || *request.FallbackDelay > 30000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Fallback delay must be between 1000 and 30000 milliseconds",
			})
			return
		}
		h.config.Chat.FallbackDelay = *request.FallbackDelay
	}

	// Return updated configuration
	response := ChatConfigResponse{
		Mode:                    h.config.Chat.Mode,
		EnableWebSocketFallback: h.config.Chat.EnableWebSocketFallback,
		WebSocketTimeout:        h.config.Chat.WebSocketTimeout,
		PollingInterval:         h.config.Chat.PollingInterval,
		MaxReconnectAttempts:    h.config.Chat.MaxReconnectAttempts,
		FallbackDelay:           h.config.Chat.FallbackDelay,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Chat configuration updated successfully",
		"config":  response,
	})
}
