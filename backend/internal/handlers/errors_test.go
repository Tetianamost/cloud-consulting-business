package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_ValidateDownloadFormat(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests
	errorHandler := NewErrorHandler(logger)

	tests := []struct {
		name        string
		format      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid pdf format",
			format:      "pdf",
			expectError: false,
		},
		{
			name:        "valid html format",
			format:      "html",
			expectError: false,
		},
		{
			name:        "empty format",
			format:      "",
			expectError: true,
			errorMsg:    "format parameter is required",
		},
		{
			name:        "invalid format",
			format:      "xml",
			expectError: true,
			errorMsg:    "invalid format 'xml', supported formats: pdf, html",
		},
		{
			name:        "case sensitive format",
			format:      "PDF",
			expectError: true,
			errorMsg:    "invalid format 'PDF', supported formats: pdf, html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorHandler.ValidateDownloadFormat(tt.format)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErrorHandler_ValidateInquiryID(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	errorHandler := NewErrorHandler(logger)

	tests := []struct {
		name        string
		inquiryID   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid inquiry ID",
			inquiryID:   "inquiry-123",
			expectError: false,
		},
		{
			name:        "empty inquiry ID",
			inquiryID:   "",
			expectError: true,
			errorMsg:    "inquiry ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorHandler.ValidateInquiryID(tt.inquiryID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestErrorHandler_HandleDownloadError(t *testing.T) {
	// Capture logs
	var logBuffer bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&logBuffer)
	logger.SetFormatter(&logrus.JSONFormatter{})

	errorHandler := NewErrorHandler(logger)

	tests := []struct {
		name           string
		errorCode      ErrorCode
		context        *ErrorContext
		expectedStatus int
		expectedCode   ErrorCode
		expectedMsg    string
	}{
		{
			name:           "invalid format error",
			errorCode:      ErrCodeInvalidFormat,
			context:        &ErrorContext{InquiryID: "test-123", Format: "xml"},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrCodeInvalidFormat,
			expectedMsg:    "Invalid format parameter",
		},
		{
			name:           "inquiry not found error",
			errorCode:      ErrCodeInquiryNotFound,
			context:        &ErrorContext{InquiryID: "missing-123", Format: "pdf"},
			expectedStatus: http.StatusNotFound,
			expectedCode:   ErrCodeInquiryNotFound,
			expectedMsg:    "Inquiry not found",
		},
		{
			name:           "no reports error",
			errorCode:      ErrCodeNoReports,
			context:        &ErrorContext{InquiryID: "empty-123", Format: "html"},
			expectedStatus: http.StatusNotFound,
			expectedCode:   ErrCodeNoReports,
			expectedMsg:    "No reports available",
		},
		{
			name:           "report generation error",
			errorCode:      ErrCodeReportGeneration,
			context:        &ErrorContext{InquiryID: "fail-123", Format: "pdf"},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   ErrCodeReportGeneration,
			expectedMsg:    "Report generation failed",
		},
		{
			name:           "unauthorized error",
			errorCode:      ErrCodeUnauthorized,
			context:        &ErrorContext{InquiryID: "auth-123", Format: "pdf"},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   ErrCodeUnauthorized,
			expectedMsg:    "Authentication required",
		},
		{
			name:           "forbidden error",
			errorCode:      ErrCodeForbidden,
			context:        &ErrorContext{InquiryID: "forbidden-123", Format: "pdf"},
			expectedStatus: http.StatusForbidden,
			expectedCode:   ErrCodeForbidden,
			expectedMsg:    "Access denied",
		},
		{
			name:           "service unavailable error",
			errorCode:      ErrCodeServiceUnavailable,
			context:        &ErrorContext{InquiryID: "service-123", Format: "pdf"},
			expectedStatus: http.StatusServiceUnavailable,
			expectedCode:   ErrCodeServiceUnavailable,
			expectedMsg:    "Service temporarily unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log buffer
			logBuffer.Reset()

			// Create test context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)

			// Call error handler
			testError := fmt.Errorf("test error for %s", tt.name)
			errorHandler.HandleDownloadError(c, testError, tt.errorCode, tt.context)

			// Check HTTP status
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response APIErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check response structure
			assert.False(t, response.Success)
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Contains(t, response.Error, tt.expectedMsg)
			assert.NotEmpty(t, response.Timestamp)
			assert.NotEmpty(t, response.TraceID)

			// Check context in response
			if tt.context != nil {
				assert.NotNil(t, response.Context)
				if tt.context.InquiryID != "" {
					assert.Equal(t, tt.context.InquiryID, response.Context["inquiry_id"])
				}
				if tt.context.Format != "" {
					assert.Equal(t, tt.context.Format, response.Context["format"])
				}
			}

			// Check that log was written
			logOutput := logBuffer.String()
			assert.NotEmpty(t, logOutput)
			assert.Contains(t, logOutput, string(tt.errorCode))
			if tt.context != nil && tt.context.InquiryID != "" {
				assert.Contains(t, logOutput, tt.context.InquiryID)
			}
		})
	}
}

func TestErrorHandler_LogSuccessfulDownload(t *testing.T) {
	// Capture logs
	var logBuffer bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&logBuffer)
	logger.SetFormatter(&logrus.JSONFormatter{})

	errorHandler := NewErrorHandler(logger)

	// Create test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test/download", nil)

	context := &ErrorContext{
		InquiryID: "success-123",
		Format:    "pdf",
		UserID:    "user-456",
		RequestID: "req-789",
	}

	fileSize := int64(1024)

	// Call success logger
	errorHandler.LogSuccessfulDownload(c, context, fileSize)

	// Check log output
	logOutput := logBuffer.String()
	assert.NotEmpty(t, logOutput)
	assert.Contains(t, logOutput, "success-123")
	assert.Contains(t, logOutput, "pdf")
	assert.Contains(t, logOutput, "user-456")
	assert.Contains(t, logOutput, "req-789")
	assert.Contains(t, logOutput, "1024")
	assert.Contains(t, logOutput, "download_report")
	assert.Contains(t, logOutput, "success")
}

func TestGetUserIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		setValue interface{}
		expected string
	}{
		{
			name:     "valid string user ID",
			setValue: "user-123",
			expected: "user-123",
		},
		{
			name:     "non-string user ID",
			setValue: 123,
			expected: "",
		},
		{
			name:     "no user ID set",
			setValue: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setValue != nil {
				c.Set("user_id", tt.setValue)
			}

			result := GetUserIDFromContext(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRequestIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		setValue interface{}
		expected string
	}{
		{
			name:     "valid string request ID",
			setValue: "req-456",
			expected: "req-456",
		},
		{
			name:     "non-string request ID",
			setValue: 456,
			expected: "",
		},
		{
			name:     "no request ID set",
			setValue: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setValue != nil {
				c.Set("request_id", tt.setValue)
			}

			result := GetRequestIDFromContext(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("uses existing trace ID from header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("X-Trace-ID", "existing-trace-123")

		traceID := generateTraceID(c)
		assert.Equal(t, "existing-trace-123", traceID)
	})

	t.Run("generates new trace ID when none exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		traceID := generateTraceID(c)
		assert.NotEmpty(t, traceID)
		assert.Contains(t, traceID, "dl_")
	})
}
