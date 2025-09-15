package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator with custom validation rules
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new custom validator instance
func NewValidator() *CustomValidator {
	v := validator.New()
	
	// Register custom validation rules
	v.RegisterValidation("service_type", validateServiceType)
	v.RegisterValidation("inquiry_status", validateInquiryStatus)
	v.RegisterValidation("priority", validatePriority)
	
	return &CustomValidator{validator: v}
}

// Validate validates a struct using the custom validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validateServiceType validates service type values
func validateServiceType(fl validator.FieldLevel) bool {
	serviceType := fl.Field().String()
	validTypes := []string{"assessment", "migration", "optimization", "architecture_review"}
	
	for _, validType := range validTypes {
		if serviceType == validType {
			return true
		}
	}
	return false
}

// validateInquiryStatus validates inquiry status values
func validateInquiryStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"pending", "processing", "reviewed", "responded", "closed"}
	
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// validatePriority validates priority values
func validatePriority(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	validPriorities := []string{"low", "medium", "high", "urgent"}
	
	for _, validPriority := range validPriorities {
		if priority == validPriority {
			return true
		}
	}
	return false
}

// SanitizeInput sanitizes user input by removing potentially harmful characters
func SanitizeInput(input string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	sanitized := re.ReplaceAllString(input, "")
	
	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)
	
	return sanitized
}

// ValidateEmail validates email format using regex
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	// Remove common phone number characters
	cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")
	
	// Check if it's a valid length (7-15 digits)
	return len(cleaned) >= 7 && len(cleaned) <= 15
}