package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Testing email formatting...")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create template service
	templateService := services.NewTemplateService("templates", logger)

	// Create mock inquiry and report
	inquiry := &domain.Inquiry{
		ID:        "test_inquiry_123",
		Name:      "customer1",
		Email:     "customer1email",
		Company:   "Partner",
		Phone:     "customer1phone",
		Services:  []string{"assessment"},
		Message:   "Quote Request Details: - Service: Initial Assessment - Complexity: Moderate - Servers/Applications: 5 - Base Fee: $750 - Per-Server/App Cost: $250 - Complexity Multiplier: 1.5x - Total Estimate: $1,500 Additional Requirements: test",
		CreatedAt: time.Now(),
	}

	report := &domain.Report{
		ID:        "dee33b60-f205-47bd-b34f-bb470e895728",
		InquiryID: inquiry.ID,
		Title:     "Professional Consulting Report Draft",
		Content: `# Professional Consulting Report Draft

Client: customer1
Email: customer1email
Company: Partner
Phone: 303-570-9119

## 1. EXECUTIVE SUMMARY

Client's Needs: customer1 from Partner has requested an initial assessment for a moderate complexity project involving 5 servers/applications.

**Key Recommendations Summary:**
- Conduct a thorough initial assessment of the current systems and applications.
- Identify potential areas for improvement and optimization.
- Provide a detailed report with actionable recommendations.

**PRIORITY LEVEL:** NORMAL

## 2. CURRENT STATE ASSESSMENT

**Client's Current Situation:** Based on customer1's inquiry, Partner is seeking a comprehensive evaluation of their current IT infrastructure.

## 3. RECOMMENDATIONS

**Actionable Recommendations:**

### Initial Assessment:
- Conduct a detailed review of the current server/application setup.
- Evaluate performance metrics, resource utilization, and potential bottlenecks.

### Optimization Plan:
- Develop a plan to optimize server performance and resource allocation.
- Recommend best practices for application management and maintenance.

## 4. NEXT STEPS

**Immediate Actions:**
- Confirm the project scope and objectives with customer1.
- Schedule the initial assessment and begin the evaluation process.

## Contact Information

**Client Name:** customer1
**Email:** customer1email
**Company:** Partner
**Phone:** 303-570-9119`,
		CreatedAt: time.Now(),
	}

	// Test consultant notification email
	fmt.Println("=== TESTING CONSULTANT NOTIFICATION EMAIL ===")

	// Prepare template data
	templateData := templateService.PrepareConsultantNotificationData(inquiry, report, false)

	// Render the HTML template
	htmlContent, err := templateService.RenderEmailTemplate(context.Background(), "consultant_notification", templateData)
	if err != nil {
		log.Fatalf("Failed to render consultant notification template: %v", err)
	}

	// Write HTML to file for inspection
	err = os.WriteFile("test_consultant_formatted.html", []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write HTML file: %v", err)
	}

	fmt.Printf("✅ Consultant notification HTML written to: test_consultant_formatted.html (%d bytes)\n", len(htmlContent))

	// Test customer confirmation email
	fmt.Println("\n=== TESTING CUSTOMER CONFIRMATION EMAIL ===")

	// Prepare customer template data using the template service method
	customerData := &interfaces.CustomerConfirmationData{
		Name:     inquiry.Name,
		Company:  inquiry.Company,
		Services: "assessment",
		ID:       inquiry.ID,
	}

	// Render the customer HTML template
	customerHTML, err := templateService.RenderEmailTemplate(context.Background(), "customer_confirmation", customerData)
	if err != nil {
		log.Fatalf("Failed to render customer confirmation template: %v", err)
	}

	// Write customer HTML to file for inspection
	err = os.WriteFile("test_customer_formatted.html", []byte(customerHTML), 0644)
	if err != nil {
		log.Fatalf("Failed to write customer HTML file: %v", err)
	}

	fmt.Printf("✅ Customer confirmation HTML written to: test_customer_formatted.html (%d bytes)\n", len(customerHTML))

	fmt.Println("\n✅ Email formatting test completed successfully!")
	fmt.Println("Open the generated HTML files in a browser to verify proper formatting.")
}
