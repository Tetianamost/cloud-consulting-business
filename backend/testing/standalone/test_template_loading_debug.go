package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("🔧 Template Loading Debug Test...")

	// Create logger with debug level
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create template service
	fmt.Println("Creating template service...")
	templateService := services.NewTemplateService("templates", logger)

	// Create test data
	inquiry := &domain.Inquiry{
		ID:        "test-template-debug",
		Name:      "tania",
		Email:     "mosttn18@gmail.com",
		Company:   "partner",
		Phone:     "3035709119",
		Services:  []string{"optimization"},
		Message:   "Quote Request Details:\n- Service: Implementation Assistance\n- Complexity: Simple\n- Hours: 1\n- Base Fee: $0\n- Hourly Rate Cost: $125\n- Complexity Multiplier: 1x\n- Total Estimate: $125\n\nAdditional Requirements:\nhelm me",
		Status:    domain.InquiryStatusPending,
		Priority:  domain.PriorityHigh,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	report := &domain.Report{
		ID:     "e193ef17-4fec-44a4-b256-9504d4847b94",
		Type:   domain.ReportTypeOptimization,
		Status: domain.ReportStatusGenerated,
		Content: `# EXECUTIVE SUMMARY

This comprehensive assessment report provides detailed recommendations for partner's cloud transformation initiative.

## CURRENT STATE ASSESSMENT

### Infrastructure Overview
- **Service Type**: Implementation Assistance
- **Complexity**: Simple (1x multiplier)
- **Estimated Investment**: $125
- **Business Impact**: Foundation for cloud optimization

## DETAILED RECOMMENDATIONS

### Phase 1: Assessment and Planning
1. **Infrastructure Analysis**
   - Current state documentation
   - Dependency mapping
   - Performance baseline establishment

## NEXT STEPS

### Immediate Actions (Next 2 Weeks)
1. **Schedule detailed technical consultation**
2. **Finalize assessment scope and timeline**`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()

	// Test template data preparation
	fmt.Println("\n📋 Testing Template Data Preparation...")
	templateData := templateService.PrepareConsultantNotificationData(inquiry, report, true)
	fmt.Printf("✅ Template data prepared: %T\n", templateData)

	// Test template rendering with detailed error handling
	fmt.Println("\n🎨 Testing Template Rendering...")
	htmlContent, err := templateService.RenderEmailTemplate(ctx, "consultant_notification", templateData)
	if err != nil {
		fmt.Printf("❌ Template rendering failed: %v\n", err)

		// Try to get available templates
		templates := templateService.GetAvailableTemplates()
		fmt.Printf("Available templates: %v\n", templates)

		return
	}

	fmt.Printf("✅ Template rendered successfully (%d characters)\n", len(htmlContent))

	// Check if the content looks like the branded template or fallback
	if len(htmlContent) > 1000 &&
		(htmlContent[:100] != "<!DOCTYPE html>\n<html>\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\"") {
		fmt.Println("✅ Content appears to be using branded template")
	} else {
		fmt.Println("⚠️  Content appears to be using fallback template")
	}

	// Show first 200 characters
	fmt.Println("\n👀 Content Preview (first 200 characters):")
	if len(htmlContent) > 200 {
		fmt.Printf("%s...\n", htmlContent[:200])
	} else {
		fmt.Printf("%s\n", htmlContent)
	}

	fmt.Println("\n🎉 Template loading debug test completed!")
}
