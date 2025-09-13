package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("Testing email integration with mock SES...")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg := &config.Config{
		SES: config.SESConfig{
			AccessKeyID:     "test-key",
			SecretAccessKey: "test-secret",
			Region:          "us-east-1",
			SenderEmail:     "noreply@cloudpartner.pro",
			ReplyToEmail:    "info@cloudpartner.pro",
			Timeout:         30,
		},
	}

	// Create template service
	templateService := services.NewTemplateService("templates", logger)

	// Create mock SES service
	mockSES := &MockSESService{}

	// Create email service
	emailService := services.NewEmailService(mockSES, templateService, cfg.SES, logger)

	// Create mock inquiry and report
	inquiry := &domain.Inquiry{
		ID:       "test_inquiry_123",
		Name:     "customer1",
		Email:    "customer1email",
		Company:  "Partner",
		Phone:    "customer1phone",
		Services: []string{"assessment"},
		Message:  "Quote Request Details: - Service: Initial Assessment - Complexity: Moderate - Servers/Applications: 5 - Base Fee: $750 - Per-Server/App Cost: $250 - Complexity Multiplier: 1.5x - Total Estimate: $1,500 Additional Requirements: test",
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

Client's Needs: customer1 from Partner has requested an initial assessment for a moderate complexity project involving 5 servers/applications. The base fee is $750, with a per-server/application cost of $250, and a complexity multiplier of 1.5x, resulting in a total estimate of $1,500.

**Key Recommendations Summary:**
- Conduct a thorough initial assessment of the current systems and applications.
- Identify potential areas for improvement and optimization.
- Provide a detailed report with actionable recommendations.

**PRIORITY LEVEL:** There are no urgency indicators or immediate meeting requests detected in the client's message. Therefore, the priority level is NORMAL.

## 2. CURRENT STATE ASSESSMENT

**Client's Current Situation:** Based on customer1's inquiry, Partner is seeking a comprehensive evaluation of their current IT infrastructure, which includes 5 servers/applications. The project has a moderate level of complexity, indicating that there might be some existing challenges or inefficiencies that need to be addressed.

**Identified Challenges and Opportunities:**
- Potential inefficiencies in server/application management.
- Opportunities for optimization and cost reduction.
- Need for a detailed assessment to inform future strategic decisions.

## 3. RECOMMENDATIONS

**Actionable Recommendations:**

### Initial Assessment:
- Conduct a detailed review of the current server/application setup.
- Evaluate performance metrics, resource utilization, and potential bottlenecks.
- Identify any security vulnerabilities or compliance issues.

### Optimization Plan:
- Develop a plan to optimize server performance and resource allocation.
- Recommend best practices for application management and maintenance.
- Suggest potential cost-saving measures and efficiency improvements.

**Prioritized Implementation Approach:**
- **Phase 1:** Initial Assessment and Reporting - Deliver a comprehensive assessment report within 2 weeks.
- **Phase 2:** Optimization and Recommendations - Provide a detailed optimization plan and recommendations within 4 weeks.

## 4. NEXT STEPS

**Immediate Actions:**
- Confirm the project scope and objectives with customer1.
- Schedule the initial assessment and begin the evaluation process.

**Proposed Engagement Timeline:**
- Week 1-2: Project kickoff and initial assessment.
- Week 3-4: Deliver assessment report and begin optimization planning.
- Week 5-6: Finalize optimization plan and recommendations.

**MEETING SCHEDULING:** There are no specific meeting requests or timeframes mentioned by the client. However, we recommend scheduling an initial meeting to discuss the project scope and objectives in detail.

## 5. URGENCY ASSESSMENT

**Urgency Indicators:**
- No urgency indicators were detected in the client's message.
- No immediate meeting requests were made.
- No specific dates/times were mentioned by the client.

**Recommended Response Timeline:** Given the normal priority level, we recommend responding within 2 business days to confirm the project scope and schedule the initial assessment.

## Contact Information

**Client Name:** customer1
**Email:** customer1email
**Company:** Partner
**Phone:** 303-570-9119

This report provides a comprehensive yet concise overview of customer1's needs and our proposed approach to addressing them.`,
		CreatedAt: time.Now(),
	}

	// Test consultant notification email
	fmt.Println("=== TESTING CONSULTANT NOTIFICATION EMAIL ===")
	
	err := emailService.SendReportEmail(context.Background(), inquiry, report)
	if err != nil {
		log.Fatalf("Failed to send report email: %v", err)
	}

	// Test customer confirmation email
	fmt.Println("\n=== TESTING CUSTOMER CONFIRMATION EMAIL ===")
	
	err = emailService.SendCustomerConfirmation(context.Background(), inquiry)
	if err != nil {
		log.Fatalf("Failed to send customer confirmation: %v", err)
	}

	fmt.Println("\n✅ Email integration test completed successfully!")
	fmt.Println("Check the generated files to see the actual email content that would be sent.")
}

// MockSESService for testing email content
type MockSESService struct{}

func (m *MockSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	fmt.Println("\n--- MOCK SES EMAIL CONTENT ---")
	fmt.Printf("From: %s\n", email.From)
	fmt.Printf("To: %s\n", strings.Join(email.To, ", "))
	fmt.Printf("Subject: %s\n", email.Subject)
	fmt.Printf("Reply-To: %s\n", email.ReplyTo)
	fmt.Printf("Has HTML Body: %t\n", email.HTMLBody != "")
	fmt.Printf("Has Text Body: %t\n", email.TextBody != "")
	fmt.Printf("HTML Body Length: %d characters\n", len(email.HTMLBody))
	fmt.Printf("Text Body Length: %d characters\n", len(email.TextBody))
	
	// Write the actual email content to files for inspection
	timestamp := time.Now().Unix()
	
	if email.HTMLBody != "" {
		filename := fmt.Sprintf("mock_email_html_%d.html", timestamp)
		os.WriteFile(filename, []byte(email.HTMLBody), 0644)
		fmt.Printf("HTML content written to: %s\n", filename)
	}
	
	if email.TextBody != "" {
		filename := fmt.Sprintf("mock_email_text_%d.txt", timestamp)
		os.WriteFile(filename, []byte(email.TextBody), 0644)
		fmt.Printf("Text content written to: %s\n", filename)
	}
	
	// Check if HTML content looks like it has proper structure
	if email.HTMLBody != "" {
		hasDoctype := strings.Contains(email.HTMLBody, "<!DOCTYPE html>")
		hasHtmlTag := strings.Contains(email.HTMLBody, "<html")
		hasBodyTag := strings.Contains(email.HTMLBody, "<body")
		hasStyles := strings.Contains(email.HTMLBody, "<style>")
		
		fmt.Printf("HTML Structure Check:\n")
		fmt.Printf("  - Has DOCTYPE: %t\n", hasDoctype)
		fmt.Printf("  - Has HTML tag: %t\n", hasHtmlTag)
		fmt.Printf("  - Has BODY tag: %t\n", hasBodyTag)
		fmt.Printf("  - Has STYLE tag: %t\n", hasStyles)
		
		if hasDoctype && hasHtmlTag && hasBodyTag && hasStyles {
			fmt.Printf("  ✅ HTML structure looks correct\n")
		} else {
			fmt.Printf("  ❌ HTML structure may have issues\n")
		}
	}
	
	fmt.Println("--- END MOCK EMAIL ---\n")
	return nil
}

func (m *MockSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (m *MockSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{
		Max24HourSend:   200,
		MaxSendRate:     1,
		SentLast24Hours: 0,
	}, nil
}