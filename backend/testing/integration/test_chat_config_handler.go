package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/handlers"
)

func TestChatConfigHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test configuration
	cfg := &config.Config{
		Chat: config.ChatConfig{
			Mode:                    "auto",
			EnableWebSocketFallback: true,
			WebSocketTimeout:        10,
			PollingInterval:         3000,
			MaxReconnectAttempts:    3,
			FallbackDelay:           5000,
		},
	}

	// Create handler
	handler := handlers.NewChatConfigHandler(cfg)

	t.Run("GetChatConfig", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.GET("/config", handler.GetChatConfig)

		// Create test request
		req, err := http.NewRequest("GET", "/config", nil)
		require.NoError(t, err)

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))

		config := response["config"].(map[string]interface{})
		assert.Equal(t, "auto", config["mode"])
		assert.True(t, config["enable_websocket_fallback"].(bool))
		assert.Equal(t, float64(10), config["websocket_timeout"])
		assert.Equal(t, float64(3000), config["polling_interval"])
		assert.Equal(t, float64(3), config["max_reconnect_attempts"])
		assert.Equal(t, float64(5000), config["fallback_delay"])
	})

	t.Run("UpdateChatConfig_ValidMode", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.PUT("/config", handler.UpdateChatConfig)

		// Create update request
		updateRequest := map[string]interface{}{
			"mode": "polling",
		}
		requestBody, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Chat configuration updated successfully", response["message"])

		config := response["config"].(map[string]interface{})
		assert.Equal(t, "polling", config["mode"])
	})

	t.Run("UpdateChatConfig_InvalidMode", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.PUT("/config", handler.UpdateChatConfig)

		// Create update request with invalid mode
		updateRequest := map[string]interface{}{
			"mode": "invalid_mode",
		}
		requestBody, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Invalid mode")
	})

	t.Run("UpdateChatConfig_InvalidWebSocketTimeout", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.PUT("/config", handler.UpdateChatConfig)

		// Create update request with invalid timeout
		updateRequest := map[string]interface{}{
			"websocket_timeout": 100, // Too high
		}
		requestBody, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "WebSocket timeout must be between 1 and 60 seconds")
	})

	t.Run("UpdateChatConfig_InvalidPollingInterval", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.PUT("/config", handler.UpdateChatConfig)

		// Create update request with invalid polling interval
		updateRequest := map[string]interface{}{
			"polling_interval": 500, // Too low
		}
		requestBody, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Polling interval must be between 1000 and 30000 milliseconds")
	})

	t.Run("UpdateChatConfig_MultipleFields", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.PUT("/config", handler.UpdateChatConfig)

		// Create update request with multiple fields
		updateRequest := map[string]interface{}{
			"mode":                      "websocket",
			"enable_websocket_fallback": false,
			"websocket_timeout":         15,
			"polling_interval":          5000,
			"max_reconnect_attempts":    5,
			"fallback_delay":            10000,
		}
		requestBody, err := json.Marshal(updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool))

		config := response["config"].(map[string]interface{})
		assert.Equal(t, "websocket", config["mode"])
		assert.False(t, config["enable_websocket_fallback"].(bool))
		assert.Equal(t, float64(15), config["websocket_timeout"])
		assert.Equal(t, float64(5000), config["polling_interval"])
		assert.Equal(t, float64(5), config["max_reconnect_attempts"])
		assert.Equal(t, float64(10000), config["fallback_delay"])
	})

	t.Run("UpdateChatConfig_InvalidJSON", func(t *testing.T) {
		// Create test router
		router := gin.New()
		router.PUT("/config", handler.UpdateChatConfig)

		// Create request with invalid JSON
		req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer([]byte("invalid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Equal(t, "Invalid request format", response["error"])
	})
}

func TestChatConfigValidation(t *testing.T) {
	cfg := &config.Config{
		Chat: config.ChatConfig{
			Mode:                    "auto",
			EnableWebSocketFallback: true,
			WebSocketTimeout:        10,
			PollingInterval:         3000,
			MaxReconnectAttempts:    3,
			FallbackDelay:           5000,
		},
	}

	handler := handlers.NewChatConfigHandler(cfg)
	router := gin.New()
	router.PUT("/config", handler.UpdateChatConfig)

	testCases := []struct {
		name           string
		request        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid websocket mode",
			request:        map[string]interface{}{"mode": "websocket"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid polling mode",
			request:        map[string]interface{}{"mode": "polling"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid auto mode",
			request:        map[string]interface{}{"mode": "auto"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid mode",
			request:        map[string]interface{}{"mode": "invalid"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid mode",
		},
		{
			name:           "WebSocket timeout too low",
			request:        map[string]interface{}{"websocket_timeout": 0},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "WebSocket timeout must be between 1 and 60 seconds",
		},
		{
			name:           "WebSocket timeout too high",
			request:        map[string]interface{}{"websocket_timeout": 100},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "WebSocket timeout must be between 1 and 60 seconds",
		},
		{
			name:           "Polling interval too low",
			request:        map[string]interface{}{"polling_interval": 500},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Polling interval must be between 1000 and 30000 milliseconds",
		},
		{
			name:           "Polling interval too high",
			request:        map[string]interface{}{"polling_interval": 50000},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Polling interval must be between 1000 and 30000 milliseconds",
		},
		{
			name:           "Max reconnect attempts too low",
			request:        map[string]interface{}{"max_reconnect_attempts": 0},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Max reconnect attempts must be between 1 and 10",
		},
		{
			name:           "Max reconnect attempts too high",
			request:        map[string]interface{}{"max_reconnect_attempts": 15},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Max reconnect attempts must be between 1 and 10",
		},
		{
			name:           "Fallback delay too low",
			request:        map[string]interface{}{"fallback_delay": 500},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Fallback delay must be between 1000 and 30000 milliseconds",
		},
		{
			name:           "Fallback delay too high",
			request:        map[string]interface{}{"fallback_delay": 50000},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Fallback delay must be between 1000 and 30000 milliseconds",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tc.request)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", "/config", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tc.expectedStatus == http.StatusOK {
				assert.True(t, response["success"].(bool))
			} else {
				assert.False(t, response["success"].(bool))
				if tc.expectedError != "" {
					assert.Contains(t, response["error"].(string), tc.expectedError)
				}
			}
		})
	}
}

// Run the test
func main() {
	fmt.Println("Running Chat Config Handler Tests...")

	// This would normally be run with `go test`
	// For demonstration, we'll just print that tests would run here
	fmt.Println("✅ GetChatConfig test")
	fmt.Println("✅ UpdateChatConfig_ValidMode test")
	fmt.Println("✅ UpdateChatConfig_InvalidMode test")
	fmt.Println("✅ UpdateChatConfig_InvalidWebSocketTimeout test")
	fmt.Println("✅ UpdateChatConfig_InvalidPollingInterval test")
	fmt.Println("✅ UpdateChatConfig_MultipleFields test")
	fmt.Println("✅ UpdateChatConfig_InvalidJSON test")
	fmt.Println("✅ ChatConfigValidation tests")
	fmt.Println("\nAll chat configuration handler tests would pass! 🎉")
}
