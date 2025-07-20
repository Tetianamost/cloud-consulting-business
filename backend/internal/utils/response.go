package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *interfaces.APIError `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
	TraceID   string      `json:"trace_id,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	APIResponse
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		TraceID:   getTraceID(c),
	}
	
	c.JSON(http.StatusOK, response)
}

// CreatedResponse sends a created response
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		TraceID:   getTraceID(c),
	}
	
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code string, message string, details string) {
	apiError := &interfaces.APIError{
		Code:       code,
		Message:    message,
		Details:    details,
		TraceID:    getTraceID(c),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		StatusCode: statusCode,
	}
	
	response := APIResponse{
		Success:   false,
		Error:     apiError,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		TraceID:   getTraceID(c),
	}
	
	c.JSON(statusCode, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, interfaces.ErrCodeValidation, 
		"Validation failed", err.Error())
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, interfaces.ErrCodeNotFound, 
		resource+" not found", "")
}

// InternalErrorResponse sends an internal server error response
func InternalErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, interfaces.ErrCodeInternal, 
		"Internal server error", message)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusUnauthorized, interfaces.ErrCodeUnauthorized, 
		"Unauthorized", "Authentication required")
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusForbidden, interfaces.ErrCodeForbidden, 
		"Forbidden", "Insufficient permissions")
}

// PaginatedSuccessResponse sends a paginated successful response
func PaginatedSuccessResponse(c *gin.Context, data interface{}, pagination *PaginationInfo, message string) {
	response := PaginatedResponse{
		APIResponse: APIResponse{
			Success:   true,
			Data:      data,
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			TraceID:   getTraceID(c),
		},
		Pagination: pagination,
	}
	
	c.JSON(http.StatusOK, response)
}

// CalculatePagination calculates pagination info
func CalculatePagination(page, limit int, total int64) *PaginationInfo {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	
	return &PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// getTraceID extracts trace ID from context
func getTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}