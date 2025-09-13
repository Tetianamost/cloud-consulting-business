package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("🔧 Email Content Debug Test...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Initialize email service (the way it's done in production)
	var emailService interfaces.EmailService
	if cfg.SES.AccessKeyID != "" && cfg.SES.SecretAccessKey != "" && cfg.SES.SenderEmail != "" {
		var err error
		emailService, err = services.NewEmailServiceWithSES(cfg.SES, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize SES service")
			return
		} else {
			logger.Info("Email service initialized successfully")
		}
	} else {
		logger.Warn("SES configuration incomplete")
		return
	}

	// Create a test inquiry similar to what you're seeing in production
	inquiry := &domain.Inquiry{
		ID:        "inq_test_debug",
		Name:      "tania",
		Email:     "mosttn18@gmail.com",
		Company:   "partner",
		Phone:     "3035709119",
		Services:  []string{"assessment"},
		Message:   "Quote Request Details:\n- Service: Initial Assessment\n- Complexity: Simple\n- Servers/Applications: 1\n- Base Fee: $750\n- Per-Server/App Cost: $50\n- Complexity Multiplier: 1x\n- Total Estimate: $800\n\nAdditional Requirements:\nhelp me",
		Status:    domain.InquiryStatusPending,
		Priority:  domain.PriorityHigh,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a test report with the content that was generated
	report := &domain.Report{
		ID:     "test-report-debug",
		Type:   domain.ReportTypeAssessment,
		Status: domain.ReportStatusGenerated,
		Content: `# EXECUTIVE SUMMARY

This comprehensive assessment report provides detailed recommendations for partner's cloud transformation initiative. Based on our analysis of the Initial Assessment request, we recommend implementing a structured approach to cloud adoption with focus on cost optimization and scalability.

## KEY RECOMMENDATIONS

**Immediate Actions Required:**
- Conduct detailed infrastructure assessment
- Implement cloud-native architecture patterns
- Establish proper security and compliance frameworks

**Priority Level: HIGH** - Simple complexity allows for rapid implementation

## CURRENT STATE ASSESSMENT

### Infrastructure Overview
- **Service Type**: Initial Assessment
- **Complexity**: Simple (1x multiplier)
- **Estimated Investment**: $800
- **Business Impact**: Foundation for cloud transformation

### Identified Opportunities
- Cost optimization through right-sizing
- Improved scalability and reliability
- Enhanced security posture
- Operational efficiency gains

## DETAILED RECOMMENDATIONS

### Phase 1: Assessment and Planning
1. **Infrastructure Analysis**
   - Current state documentation
   - Dependency mapping
   - Performance baseline establishment

2. **Cloud Strategy Development**
   - Multi-cloud provider evaluation
   - Cost-benefit analysis
   - Risk assessment and mitigation

### Phase 2: Implementation Planning
1. **Architecture Design**
   - Cloud-native patterns adoption
   - Security framework implementation
   - Monitoring and observability setup

## NEXT STEPS

### Immediate Actions (Next 2 Weeks)
1. **Schedule detailed technical consultation**
2. **Finalize assessment scope and timeline**
3. **Begin infrastructure documentation**

### Contact Information
For immediate assistance: info@cloudpartner.pro`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inquiry.Reports = []*domain.Report{report}

	ctx := context.Background()

	fmt.Println("\n📧 Testing Internal Report Email Generation...")

	// Test the actual email generation that's happening in production
	err = emailService.SendReportEmail(ctx, inquiry, report)
	if err != nil {
		fmt.Printf("❌ Failed to send report email: %v\n", err)
	} else {
		fmt.Println("✅ Report email sent successfully!")
		fmt.Println("   Check your email at info@cloudpartner.pro to see the actual content")
	}

	fmt.Println("\n🎉 Email content debug test completed!")
}
