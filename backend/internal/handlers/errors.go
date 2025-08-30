package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorCode represents standardized error codes for the API
type ErrorCode string

const (
	// Download-specific error codes
	ErrCodeInvalidFormat      ErrorCode = "INVALID_FORMAT"
	ErrCodeInquiryNotFound    ErrorCode = "INQUIRY_NOT_FOUND"
	ErrCodeNoReports          ErrorCode = "NO_REPORTS"
	ErrCodeReportGeneration   ErrorCode = "REPORT_GENERATION_ERROR"
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidationError    ErrorCode = "VALIDATION_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// APIErrorResponse represents a structured error response
type APIErrorResponse struct {
	Success   bool                   `json:"success"`
	Error     string                 `json:"error"`
	Code      ErrorCode              `json:"code"`
	Details   string                 `json:"details,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// ErrorContext contains contextual information for error logging
type ErrorContext struct {
	InquiryID string
	Format    string
	UserID    string
	Action    string
	RequestID string
}

// ErrorHandler provides centralized error handling for handlers
type ErrorHandler struct {
	logger *logrus.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleDownloadError handles download-specific errors with proper logging and response
func (eh *ErrorHandler) HandleDownloadError(c *gin.Context, err error, code ErrorCode, context *ErrorContext) {
	// Generate trace ID for request tracking
	traceID := generateTraceID(c)

	// Create structured log entry
	logFields := logrus.Fields{
		"error_code": string(code),
		"trace_id":   traceID,
		"action":     "download_report",
	}

	// Add context information if provided
	if context != nil {
		if context.InquiryID != "" {
			logFields["inquiry_id"] = context.InquiryID
		}
		if context.Format != "" {
			logFields["format"] = context.Format
		}
		if context.UserID != "" {
			logFields["user_id"] = context.UserID
		}
		if context.RequestID != "" {
			logFields["request_id"] = context.RequestID
		}
	}

	// Add request information
	logFields["method"] = c.Request.Method
	logFields["path"] = c.Request.URL.Path
	logFields["remote_addr"] = c.ClientIP()
	logFields["user_agent"] = c.Request.UserAgent()

	// Determine HTTP status code and error message based on error code
	var statusCode int
	var errorMessage string
	var details string

	switch code {
	case ErrCodeInvalidFormat:
		statusCode = http.StatusBadRequest
		errorMessage = "Invalid format parameter"
		details = "Supported formats are: pdf, html"
		eh.logger.WithFields(logFields).WithError(err).Warn("Invalid format parameter provided")

	case ErrCodeInquiryNotFound:
		statusCode = http.StatusNotFound
		errorMessage = "Inquiry not found"
		details = "The specified inquiry ID does not exist or you don't have access to it"
		eh.logger.WithFields(logFields).WithError(err).Warn("Inquiry not found for download request")

	case ErrCodeNoReports:
		statusCode = http.StatusNotFound
		errorMessage = "No reports available"
		details = "No reports have been generated for this inquiry yet"
		eh.logger.WithFields(logFields).WithError(err).Info("No reports available for download")

	case ErrCodeReportGeneration:
		statusCode = http.StatusInternalServerError
		errorMessage = "Report generation failed"
		details = "Unable to generate the requested report format. Please try again later"
		eh.logger.WithFields(logFields).WithError(err).Error("Report generation failed")

	case ErrCodeUnauthorized:
		statusCode = http.StatusUnauthorized
		errorMessage = "Authentication required"
		details = "Please provide valid authentication credentials"
		eh.logger.WithFields(logFields).WithError(err).Warn("Unauthorized download attempt")

	case ErrCodeForbidden:
		statusCode = http.StatusForbidden
		errorMessage = "Access denied"
		details = "You don't have permission to download this report"
		eh.logger.WithFields(logFields).WithError(err).Warn("Forbidden download attempt")

	case ErrCodeServiceUnavailable:
		statusCode = http.StatusServiceUnavailable
		errorMessage = "Service temporarily unavailable"
		details = "The report generation service is currently unavailable. Please try again later"
		eh.logger.WithFields(logFields).WithError(err).Error("Service unavailable for report generation")

	default:
		statusCode = http.StatusInternalServerError
		errorMessage = "Internal server error"
		details = "An unexpected error occurred. Please try again later"
		eh.logger.WithFields(logFields).WithError(err).Error("Unhandled error in download request")
	}

	// Create error response
	errorResponse := APIErrorResponse{
		Success:   false,
		Error:     errorMessage,
		Code:      code,
		Details:   details,
		TraceID:   traceID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Add context information to response if provided
	if context != nil {
		errorResponse.Context = make(map[string]interface{})
		if context.InquiryID != "" {
			errorResponse.Context["inquiry_id"] = context.InquiryID
		}
		if context.Format != "" {
			errorResponse.Context["format"] = context.Format
		}
	}

	c.JSON(statusCode, errorResponse)
}

// ValidateDownloadFormat validates the format parameter for report downloads
func (eh *ErrorHandler) ValidateDownloadFormat(format string) error {
	validFormats := map[string]bool{
		"pdf":  true,
		"html": true,
	}

	if format == "" {
		return fmt.Errorf("format parameter is required")
	}

	if !validFormats[format] {
		return fmt.Errorf("invalid format '%s', supported formats: pdf, html", format)
	}

	return nil
}

// ValidateInquiryID validates the inquiry ID parameter
func (eh *ErrorHandler) ValidateInquiryID(inquiryID string) error {
	if inquiryID == "" {
		return fmt.Errorf("inquiry ID is required")
	}

	// Basic format validation - inquiry IDs should be non-empty strings
	if len(inquiryID) < 1 {
		return fmt.Errorf("invalid inquiry ID format")
	}

	return nil
}

// LogSuccessfulDownload logs successful download operations
func (eh *ErrorHandler) LogSuccessfulDownload(c *gin.Context, context *ErrorContext, fileSize int64) {
	logFields := logrus.Fields{
		"action":      "download_report",
		"status":      "success",
		"file_size":   fileSize,
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"remote_addr": c.ClientIP(),
		"user_agent":  c.Request.UserAgent(),
	}

	if context != nil {
		if context.InquiryID != "" {
			logFields["inquiry_id"] = context.InquiryID
		}
		if context.Format != "" {
			logFields["format"] = context.Format
		}
		if context.UserID != "" {
			logFields["user_id"] = context.UserID
		}
		if context.RequestID != "" {
			logFields["request_id"] = context.RequestID
		}
	}

	eh.logger.WithFields(logFields).Info("Report download completed successfully")
}

// generateTraceID generates a unique trace ID for request tracking
func generateTraceID(c *gin.Context) string {
	// Try to get existing trace ID from headers
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		return traceID
	}

	// Generate a simple trace ID based on timestamp and request info
	return fmt.Sprintf("dl_%d_%s", time.Now().UnixNano(), c.ClientIP())
}

// GetUserIDFromContext extracts user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if userIDStr, ok := userID.(string); ok {
			return userIDStr
		}
	}
	return ""
}

// GetRequestIDFromContext extracts request ID from the Gin context
func GetRequestIDFromContext(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if requestIDStr, ok := requestID.(string); ok {
			return requestIDStr
		}
	}
	return ""
}
