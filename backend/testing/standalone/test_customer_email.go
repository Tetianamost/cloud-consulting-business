package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create template service
	templateService := services.NewTemplateService("templates", logger)

	// Create sample inquiry
	inquiry := &domain.Inquiry{
		ID:       "test-customer-inquiry-789",
		Name:     "Sarah Johnson",
		Email:    "sarah.johnson@techcorp.com",
		Company:  "TechCorp Solutions",
		Phone:    "+1-555-0199",
		Services: []string{"migration", "optimization", "assessment"},
		Message:  "We need help with our cloud migration project. The current system is experiencing performance issues and we need to migrate to AWS. Can we schedule a consultation to discuss the timeline and approach?",
	}

	ctx := context.Background()

	// Test customer confirmation email
	fmt.Println("Testing customer confirmation email...")
	
	// Prepare template data manually
	templateData := struct {
		Name     string
		Company  string
		Services string
		ID       string
	}{
		Name:     inquiry.Name,
		Company:  inquiry.Company,
		Services: strings.Join(inquiry.Services, ", "),
		ID:       inquiry.ID,
	}

	// Render template
	emailHTML, err := templateService.RenderEmailTemplate(ctx, "customer_confirmation", templateData)
	if err != nil {
		fmt.Printf("Error rendering customer confirmation template: %v\n", err)
		return
	}

	// Save to file
	filename := "test_customer_confirmation.html"
	err = os.WriteFile(filename, []byte(emailHTML), 0644)
	if err != nil {
		fmt.Printf("Error saving customer email: %v\n", err)
		return
	}

	fmt.Printf("✓ Customer confirmation email generated: %s\n", filename)

	// Test with minimal data (no company)
	fmt.Println("Testing customer confirmation email with minimal data...")
	
	minimalInquiry := &domain.Inquiry{
		ID:       "test-minimal-inquiry-999",
		Name:     "John Doe",
		Email:    "john.doe@email.com",
		Company:  "", // No company
		Services: []string{"assessment"},
		Message:  "I need a cloud assessment for my small business.",
	}

	// Prepare template data for minimal inquiry
	minimalTemplateData := struct {
		Name     string
		Company  string
		Services string
		ID       string
	}{
		Name:     minimalInquiry.Name,
		Company:  minimalInquiry.Company,
		Services: strings.Join(minimalInquiry.Services, ", "),
		ID:       minimalInquiry.ID,
	}

	// Render template
	minimalEmailHTML, err := templateService.RenderEmailTemplate(ctx, "customer_confirmation", minimalTemplateData)
	if err != nil {
		fmt.Printf("Error rendering minimal customer confirmation template: %v\n", err)
		return
	}

	// Save to file
	minimalFilename := "test_customer_confirmation_minimal.html"
	err = os.WriteFile(minimalFilename, []byte(minimalEmailHTML), 0644)
	if err != nil {
		fmt.Printf("Error saving minimal customer email: %v\n", err)
		return
	}

	fmt.Printf("✓ Minimal customer confirmation email generated: %s\n", minimalFilename)
	fmt.Println("\nCustomer email template test completed!")
	fmt.Println("Open the HTML files in a browser to verify the formatting.")
}