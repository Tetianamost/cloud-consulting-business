package main

import (
	"context"
	"fmt"
	"log"
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
		ID:       "test-inquiry-123",
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Company:  "Tech Corp",
		Phone:    "+1-555-0123",
		Services: []string{"assessment", "migration"},
		Message:  "We need urgent help with our cloud migration. Can we schedule a meeting this week?",
	}

	// Create sample report
	report := &domain.Report{
		ID:        "test-report-456",
		InquiryID: inquiry.ID,
		Type:      domain.ReportTypeAssessment,
		Title:     "Cloud Assessment Report for Tech Corp",
		Content: `EXECUTIVE SUMMARY

This report provides a comprehensive assessment of Tech Corp's cloud migration requirements based on the urgent inquiry received.

PRIORITY LEVEL: HIGH PRIORITY - Client has requested immediate meeting and expressed urgency.

CURRENT STATE ASSESSMENT

The client has indicated an urgent need for cloud migration assistance. Key challenges identified:
- Time-sensitive migration requirements
- Need for immediate consultation
- Request for meeting scheduling this week

RECOMMENDATIONS

1. **Immediate Response Required**
   - Schedule consultation call within 24 hours
   - Prepare migration assessment framework
   - Allocate senior architect resources

2. **Migration Planning**
   - Conduct detailed infrastructure audit
   - Develop phased migration strategy
   - Establish timeline and milestones

NEXT STEPS

- Contact client immediately to schedule meeting
- Prepare migration assessment questionnaire
- Assign dedicated migration specialist

URGENCY ASSESSMENT

**URGENT LANGUAGE DETECTED**: "urgent help", "schedule a meeting this week"
**RECOMMENDED RESPONSE TIME**: Within 4 hours
**PRIORITY LEVEL**: HIGH

CONTACT INFORMATION

Client: John Doe
Email: john.doe@example.com
Company: Tech Corp
Phone: +1-555-0123`,
		Status:      domain.ReportStatusDraft,
		GeneratedBy: "test-system",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test different report types
	reportTypes := []struct {
		name       string
		reportType domain.ReportType
		template   string
	}{
		{"Assessment", domain.ReportTypeAssessment, "assessment"},
		{"Migration", domain.ReportTypeMigration, "migration"},
		{"Optimization", domain.ReportTypeOptimization, "optimization"},
		{"Architecture", domain.ReportTypeArchitecture, "architecture"},
	}

	ctx := context.Background()

	for _, rt := range reportTypes {
		fmt.Printf("Testing %s report template...\n", rt.name)

		// Update report type
		testReport := *report
		testReport.Type = rt.reportType
		testReport.Title = fmt.Sprintf("%s Report for %s", rt.name, inquiry.Company)

		// Prepare template data
		templateData := templateService.PrepareReportTemplateData(inquiry, &testReport)

		// Render template
		htmlContent, err := templateService.RenderReportTemplate(ctx, rt.template, templateData)
		if err != nil {
			log.Printf("Error rendering %s template: %v", rt.name, err)
			continue
		}

		// Save to file
		filename := fmt.Sprintf("test_%s_report.html", rt.template)
		err = os.WriteFile(filename, []byte(htmlContent), 0644)
		if err != nil {
			log.Printf("Error saving %s report: %v", rt.name, err)
			continue
		}

		fmt.Printf("✓ %s report generated successfully: %s\n", rt.name, filename)
	}

	fmt.Println("\nHTML report generation test completed!")
	fmt.Println("Check the generated HTML files to verify the formatting.")
}