package main

import (
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Testing Enhanced Download Error Handling...")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create error handler
	errorHandler := handlers.NewErrorHandler(logger)

	// Test format validation
	fmt.Println("\n1. Testing format validation:")

	validFormats := []string{"pdf", "html"}
	invalidFormats := []string{"", "xml", "PDF", "json", "txt"}

	for _, format := range validFormats {
		err := errorHandler.ValidateDownloadFormat(format)
		if err != nil {
			log.Fatalf("Expected valid format '%s' to pass validation, got error: %v", format, err)
		}
		fmt.Printf("✓ Format '%s' is valid\n", format)
	}

	for _, format := range invalidFormats {
		err := errorHandler.ValidateDownloadFormat(format)
		if err == nil {
			log.Fatalf("Expected invalid format '%s' to fail validation", format)
		}
		fmt.Printf("✓ Format '%s' correctly rejected: %v\n", format, err)
	}

	// Test inquiry ID validation
	fmt.Println("\n2. Testing inquiry ID validation:")

	validIDs := []string{"inquiry-123", "test-id", "abc123"}
	invalidIDs := []string{""}

	for _, id := range validIDs {
		err := errorHandler.ValidateInquiryID(id)
		if err != nil {
			log.Fatalf("Expected valid inquiry ID '%s' to pass validation, got error: %v", id, err)
		}
		fmt.Printf("✓ Inquiry ID '%s' is valid\n", id)
	}

	for _, id := range invalidIDs {
		err := errorHandler.ValidateInquiryID(id)
		if err == nil {
			log.Fatalf("Expected invalid inquiry ID '%s' to fail validation", id)
		}
		fmt.Printf("✓ Inquiry ID '%s' correctly rejected: %v\n", id, err)
	}

	// Test error codes
	fmt.Println("\n3. Testing error codes:")

	errorCodes := []handlers.ErrorCode{
		handlers.ErrCodeInvalidFormat,
		handlers.ErrCodeInquiryNotFound,
		handlers.ErrCodeNoReports,
		handlers.ErrCodeReportGeneration,
		handlers.ErrCodeUnauthorized,
		handlers.ErrCodeForbidden,
		handlers.ErrCodeInternalError,
		handlers.ErrCodeValidationError,
		handlers.ErrCodeServiceUnavailable,
	}

	for _, code := range errorCodes {
		fmt.Printf("✓ Error code '%s' is defined\n", string(code))
	}

	fmt.Println("\n✅ All validation tests passed!")
	fmt.Println("\nEnhanced error handling features implemented:")
	fmt.Println("- Structured error response types with error codes")
	fmt.Println("- Comprehensive error logging with contextual information")
	fmt.Println("- Format parameter validation with clear error messages")
	fmt.Println("- Proper HTTP status codes for different error scenarios")
	fmt.Println("- Trace ID generation for request tracking")
	fmt.Println("- Success logging for completed downloads")
}
