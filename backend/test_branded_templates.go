package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Initialize template service
	templateService := services.NewTemplateService("templates", logger)

	// Test customer confirmation template
	fmt.Println("=== Testing Customer Confirmation Template ===")
	
	inquiry := &domain.Inquiry{
		ID:       "test-inquiry-123",
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Company:  "Acme Corporation",
		Services: []string{"assessment", "migration"},
		Message:  "We need help migrating our infrastructure to AWS.",
	}

	// Prepare template data
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

	// Render customer confirmation template
	customerHTML, err := templateService.RenderEmailTemplate(context.Background(), "customer_confirmation", templateData)
	if err != nil {
		log.Printf("Error rendering customer confirmation template: %v", err)
	} else {
		fmt.Printf("Customer confirmation template rendered successfully (%d characters)\n", len(customerHTML))
		
		// Save to file for inspection
		err = os.WriteFile("customer_confirmation_test.html", []byte(customerHTML), 0644)
		if err != nil {
			log.Printf("Error saving customer confirmation HTML: %v", err)
		} else {
			fmt.Println("Customer confirmation HTML saved to customer_confirmation_test.html")
		}
	}

	// Test consultant notification template
	fmt.Println("\n=== Testing Consultant Notification Template ===")
	
	report := &domain.Report{
		ID:      "test-report-456",
		Content: "# Cloud Assessment Report\n\n## Executive Summary\n\nThis report provides a comprehensive assessment of your current cloud infrastructure.\n\n## Recommendations\n\n- **Immediate Actions**: Implement security best practices\n- **Short-term Goals**: Optimize cost structure\n- **Long-term Strategy**: Adopt microservices architecture",
	}

	// Prepare consultant notification template data
	consultantData := struct {
		Name           string
		Email          string
		Company        string
		Phone          string
		Services       string
		Message        string
		ID             string
		IsHighPriority bool
		Priority       string
		Report         *struct {
			ID          string
			Content     string
			HTMLContent string
		}
	}{
		Name:           inquiry.Name,
		Email:          inquiry.Email,
		Company:        inquiry.Company,
		Phone:          "555-123-4567",
		Services:       strings.Join(inquiry.Services, ", "),
		Message:        inquiry.Message,
		ID:             inquiry.ID,
		IsHighPriority: false,
		Priority:       "NORMAL",
		Report: &struct {
			ID          string
			Content     string
			HTMLContent string
		}{
			ID:          report.ID,
			Content:     report.Content,
			HTMLContent: convertMarkdownToHTML(report.Content),
		},
	}

	// Render consultant notification template
	consultantHTML, err := templateService.RenderEmailTemplate(context.Background(), "consultant_notification", consultantData)
	if err != nil {
		log.Printf("Error rendering consultant notification template: %v", err)
	} else {
		fmt.Printf("Consultant notification template rendered successfully (%d characters)\n", len(consultantHTML))
		
		// Save to file for inspection
		err = os.WriteFile("consultant_notification_test.html", []byte(consultantHTML), 0644)
		if err != nil {
			log.Printf("Error saving consultant notification HTML: %v", err)
		} else {
			fmt.Println("Consultant notification HTML saved to consultant_notification_test.html")
		}
	}

	// Test high priority version
	fmt.Println("\n=== Testing High Priority Consultant Notification ===")
	
	consultantData.IsHighPriority = true
	consultantData.Priority = "HIGH"
	consultantData.Message = "URGENT: We need immediate help with our production environment. Our systems are down and we need assistance ASAP!"

	highPriorityHTML, err := templateService.RenderEmailTemplate(context.Background(), "consultant_notification", consultantData)
	if err != nil {
		log.Printf("Error rendering high priority template: %v", err)
	} else {
		fmt.Printf("High priority template rendered successfully (%d characters)\n", len(highPriorityHTML))
		
		// Save to file for inspection
		err = os.WriteFile("consultant_notification_high_priority_test.html", []byte(highPriorityHTML), 0644)
		if err != nil {
			log.Printf("Error saving high priority HTML: %v", err)
		} else {
			fmt.Println("High priority HTML saved to consultant_notification_high_priority_test.html")
		}
	}

	fmt.Println("\n=== Template Testing Complete ===")
	fmt.Println("All HTML files have been generated for visual inspection.")
	fmt.Println("Open the generated .html files in a web browser to see the branded templates.")
}

// Simple markdown to HTML conversion for testing
func convertMarkdownToHTML(markdown string) string {
	html := markdown
	
	// Convert headers
	html = strings.ReplaceAll(html, "# ", "<h1>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "### ", "<h3>")
	html = strings.ReplaceAll(html, "#### ", "<h4>")
	
	// Convert bold text
	html = strings.ReplaceAll(html, "**", "<strong>")
	
	// Convert line breaks to paragraphs
	paragraphs := strings.Split(html, "\n\n")
	for i, p := range paragraphs {
		if strings.TrimSpace(p) != "" && !strings.HasPrefix(p, "<h") {
			paragraphs[i] = "<p>" + strings.ReplaceAll(p, "\n", "<br>") + "</p>"
		}
	}
	
	return strings.Join(paragraphs, "\n")
}