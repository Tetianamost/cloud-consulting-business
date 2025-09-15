package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
		ID:       "test-inquiry-email-123",
		Name:     "Sarah Johnson",
		Email:    "sarah.johnson@techcorp.com",
		Company:  "TechCorp Solutions",
		Phone:    "+1-555-0199",
		Services: []string{"migration", "optimization"},
		Message:  "We need urgent help with our cloud migration project. The current system is experiencing performance issues and we need to migrate to AWS ASAP. Can we schedule a call this afternoon to discuss the timeline?",
	}

	// Create sample report
	report := &domain.Report{
		ID:        "test-report-email-456",
		InquiryID: inquiry.ID,
		Type:      domain.ReportTypeMigration,
		Title:     "Urgent Cloud Migration Assessment for TechCorp Solutions",
		Content: `### EXECUTIVE SUMMARY

**Client's Needs:**
TechCorp Solutions has requested urgent assistance with their cloud migration project. The client is experiencing performance issues with their current system and needs to migrate to AWS as soon as possible. They have requested an immediate call to discuss timeline and implementation.

**Key Recommendations Summary:**
Given the urgency and performance issues, we recommend an accelerated migration approach with immediate performance optimization. A phased migration strategy will minimize downtime while addressing current performance bottlenecks.

**PRIORITY LEVEL: HIGH PRIORITY**
The client's use of "urgent", "ASAP", and request for same-day meeting indicates critical business impact requiring immediate attention.

### CURRENT STATE ASSESSMENT

**Analysis:**
TechCorp Solutions is experiencing active performance issues that are impacting business operations. The urgency of the migration request suggests these issues may be causing significant business disruption.

**Identified Challenges:**
- Performance bottlenecks in current infrastructure
- Time-sensitive migration requirements
- Need for immediate technical consultation
- Business continuity concerns during migration

### RECOMMENDATIONS

**Immediate Actions:**
1. **Emergency Performance Assessment** - Conduct immediate analysis of current performance issues
2. **Accelerated Migration Planning** - Develop fast-track migration strategy for AWS
3. **Resource Allocation** - Assign senior migration specialists immediately

**Migration Strategy:**
- **Phase 1:** Critical system stabilization and performance fixes
- **Phase 2:** Core infrastructure migration to AWS
- **Phase 3:** Optimization and performance tuning

### NEXT STEPS

**Immediate Actions:**
- Schedule emergency consultation call within 2 hours
- Begin performance assessment of current systems
- Prepare AWS migration architecture proposal

**Timeline:**
- **Today:** Initial consultation and performance assessment
- **Week 1:** Migration planning and AWS environment setup
- **Week 2-3:** Phased migration execution

### URGENCY ASSESSMENT

**Urgent Language Detected:**
- "urgent help"
- "ASAP"
- "schedule a call this afternoon"

**Business Impact:**
- Active performance issues affecting operations
- Time-sensitive migration requirements
- Same-day consultation requested

**Recommended Response Time:** Within 1 hour

### CONTACT INFORMATION

**Client:** Sarah Johnson
**Email:** sarah.johnson@techcorp.com
**Company:** TechCorp Solutions
**Phone:** +1-555-0199`,
		Status:      domain.ReportStatusDraft,
		GeneratedBy: "bedrock-ai",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ctx := context.Background()

	// Test high priority email
	fmt.Println("Testing high priority consultant notification email...")
	
	// Prepare template data
	templateData := templateService.PrepareConsultantNotificationData(inquiry, report, true)

	// Render template
	emailHTML, err := templateService.RenderEmailTemplate(ctx, "consultant_notification", templateData)
	if err != nil {
		fmt.Printf("Error rendering email template: %v\n", err)
		return
	}

	// Save to file
	filename := "test_consultant_notification_high_priority.html"
	err = os.WriteFile(filename, []byte(emailHTML), 0644)
	if err != nil {
		fmt.Printf("Error saving email: %v\n", err)
		return
	}

	fmt.Printf("✓ High priority email generated: %s\n", filename)

	// Test normal priority email (without report)
	fmt.Println("Testing normal priority consultant notification email...")
	
	normalInquiry := &domain.Inquiry{
		ID:       "test-inquiry-normal-123",
		Name:     "John Smith",
		Email:    "john.smith@example.com",
		Company:  "Example Corp",
		Phone:    "+1-555-0100",
		Services: []string{"assessment"},
		Message:  "We are interested in getting a cloud assessment for our infrastructure. Please let us know your availability for a consultation next week.",
	}

	// Prepare template data for normal priority
	normalTemplateData := templateService.PrepareConsultantNotificationData(normalInquiry, nil, false)

	// Render template
	normalEmailHTML, err := templateService.RenderEmailTemplate(ctx, "consultant_notification", normalTemplateData)
	if err != nil {
		fmt.Printf("Error rendering normal email template: %v\n", err)
		return
	}

	// Save to file
	normalFilename := "test_consultant_notification_normal.html"
	err = os.WriteFile(normalFilename, []byte(normalEmailHTML), 0644)
	if err != nil {
		fmt.Printf("Error saving normal email: %v\n", err)
		return
	}

	fmt.Printf("✓ Normal priority email generated: %s\n", normalFilename)
	fmt.Println("\nEmail template test completed!")
	fmt.Println("Open the HTML files in a browser to verify the formatting.")
}