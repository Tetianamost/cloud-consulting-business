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
	fmt.Println("🔧 Email Template Debug Test...")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create template service
	templateService := services.NewTemplateService("templates", logger)

	// Create a test inquiry with rich data
	inquiry := &domain.Inquiry{
		ID:        "test-debug-" + fmt.Sprintf("%d", time.Now().Unix()),
		Name:      "John Smith",
		Email:     "john.smith@techcorp.com",
		Company:   "TechCorp Solutions",
		Phone:     "+1 (555) 123-4567",
		Services:  []string{"assessment", "migration", "optimization"},
		Message:   "We need help migrating our legacy applications to AWS. We have about 15 applications running on on-premises servers and need a comprehensive assessment and migration plan. Timeline is critical as our current infrastructure contract expires in 6 months.",
		Status:    domain.InquiryStatusPending,
		Priority:  domain.PriorityHigh,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a comprehensive test report
	report := &domain.Report{
		ID:     "test-report-debug-" + fmt.Sprintf("%d", time.Now().Unix()),
		Type:   domain.ReportTypeAssessment,
		Status: domain.ReportStatusGenerated,
		Content: `# EXECUTIVE SUMMARY

This comprehensive assessment report provides detailed recommendations for TechCorp Solutions' cloud migration initiative. Based on our analysis, we recommend a phased approach to migrate 15 applications to AWS over a 4-month timeline.

## KEY RECOMMENDATIONS

**Immediate Actions Required:**
- Conduct detailed application dependency mapping
- Implement AWS Landing Zone with proper security controls
- Begin with pilot migration of 3 low-risk applications

**Priority Level: HIGH** - Timeline constraints require immediate action

## CURRENT STATE ASSESSMENT

### Infrastructure Overview
- **Applications**: 15 legacy applications
- **Current Environment**: On-premises data center
- **Contract Expiration**: 6 months
- **Business Impact**: Critical timeline constraints

### Identified Challenges
- Legacy application dependencies
- Data migration complexity
- Security and compliance requirements
- Staff training and change management

## DETAILED RECOMMENDATIONS

### Phase 1: Foundation (Month 1)
1. **AWS Account Setup and Security**
   - Implement AWS Organizations with multi-account strategy
   - Configure AWS SSO and IAM policies
   - Set up CloudTrail and Config for compliance

2. **Network Architecture**
   - Design VPC architecture with proper subnetting
   - Implement Direct Connect for hybrid connectivity
   - Configure security groups and NACLs

### Phase 2: Pilot Migration (Month 2)
1. **Application Assessment**
   - Detailed analysis of 3 pilot applications
   - Performance baseline establishment
   - Migration strategy validation

2. **Migration Execution**
   - Database migration using AWS DMS
   - Application server migration
   - Testing and validation

### Phase 3: Production Migration (Months 3-4)
1. **Remaining Applications**
   - Batch migration of remaining 12 applications
   - Performance optimization
   - Disaster recovery implementation

## RISK ASSESSMENT

### High Priority Risks
- **Timeline Risk**: 6-month contract expiration
- **Data Loss Risk**: Complex database migrations
- **Downtime Risk**: Business-critical applications

### Mitigation Strategies
- Parallel migration approach to reduce timeline
- Comprehensive backup and rollback procedures
- Phased cutover with minimal downtime windows

## FINANCIAL ANALYSIS

### Estimated Costs
- **Migration Services**: $75,000 - $100,000
- **AWS Infrastructure**: $8,000 - $12,000/month
- **Training and Support**: $15,000 - $20,000

### ROI Projections
- **Year 1 Savings**: $50,000 (reduced infrastructure costs)
- **Year 2+ Savings**: $80,000/year (operational efficiency)
- **Payback Period**: 18-24 months

## NEXT STEPS

### Immediate Actions (Next 2 Weeks)
1. **Schedule detailed technical assessment meeting**
2. **Finalize migration timeline and resource allocation**
3. **Begin AWS account setup and security configuration**
4. **Identify pilot applications for Phase 1**

### Meeting Request
**URGENT**: We recommend scheduling a technical deep-dive meeting within the next week to discuss:
- Detailed application inventory and dependencies
- Migration timeline and resource requirements
- Risk mitigation strategies and contingency planning

## CONTACT INFORMATION

For immediate assistance or to schedule the recommended meeting:
- **Primary Contact**: Senior Cloud Architect
- **Email**: info@cloudpartner.pro
- **Phone**: Available upon request
- **Response Time**: Within 4 hours for urgent matters

---

*This report was generated using advanced AI analysis and validated by our certified cloud architects. All recommendations follow AWS Well-Architected Framework principles and industry best practices.*`,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inquiry.Reports = []*domain.Report{report}

	ctx := context.Background()

	// Test template data preparation
	fmt.Println("\n📋 Testing Template Data Preparation...")
	templateData := templateService.PrepareConsultantNotificationData(inquiry, report, true)
	fmt.Printf("✅ Template data prepared successfully\n")
	fmt.Printf("   Data type: %T\n", templateData)

	// Test template rendering
	fmt.Println("\n🎨 Testing Template Rendering...")
	htmlContent, err := templateService.RenderEmailTemplate(ctx, "consultant_notification", templateData)
	if err != nil {
		log.Fatalf("❌ Failed to render template: %v", err)
	}

	fmt.Printf("✅ Template rendered successfully (%d characters)\n", len(htmlContent))

	// Save the rendered HTML to a file for inspection
	outputFile := "debug_email_output.html"
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write output file: %v", err)
	}

	fmt.Printf("✅ Email HTML saved to: %s\n", outputFile)

	// Test customer confirmation template as well
	fmt.Println("\n📧 Testing Customer Confirmation Template...")
	customerData := templateService.PrepareCustomerConfirmationData(inquiry)
	customerHTML, err := templateService.RenderEmailTemplate(ctx, "customer_confirmation", customerData)
	if err != nil {
		log.Fatalf("❌ Failed to render customer template: %v", err)
	}

	customerOutputFile := "debug_customer_email_output.html"
	err = os.WriteFile(customerOutputFile, []byte(customerHTML), 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write customer output file: %v", err)
	}

	fmt.Printf("✅ Customer email HTML saved to: %s\n", customerOutputFile)

	fmt.Println("\n🎉 Email template debug test completed!")
	fmt.Println("\n📝 Summary:")
	fmt.Printf("   - Consultant email: %d characters\n", len(htmlContent))
	fmt.Printf("   - Customer email: %d characters\n", len(customerHTML))
	fmt.Println("   - Both templates rendered successfully")
	fmt.Println("   - Check the generated HTML files to verify content quality")

	// Show a preview of the content
	fmt.Println("\n👀 Content Preview (first 500 characters):")
	if len(htmlContent) > 500 {
		fmt.Printf("Consultant Email: %s...\n", htmlContent[:500])
	} else {
		fmt.Printf("Consultant Email: %s\n", htmlContent)
	}
}
