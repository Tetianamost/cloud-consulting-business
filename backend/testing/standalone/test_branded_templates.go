package main

import (
	"context"
	"fmt"
	"os"
	"strings"
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

	ctx := context.Background()

	fmt.Println("🎨 Testing Branded Email Templates")
	fmt.Println("==================================")

	// Test Case 1: High Priority Consultant Notification with Report
	fmt.Println("\n1. Testing HIGH PRIORITY consultant notification with AI report...")
	
	urgentInquiry := &domain.Inquiry{
		ID:       "urgent-inquiry-001",
		Name:     "Sarah Johnson",
		Email:    "sarah.johnson@techcorp.com",
		Company:  "TechCorp Solutions Inc.",
		Phone:    "+1-555-0199",
		Services: []string{"migration", "optimization", "assessment"},
		Message:  "URGENT: We need immediate help with our cloud migration project. Our current system is experiencing critical performance issues and we need to migrate to AWS ASAP. Can we schedule an emergency call this afternoon to discuss the timeline? This is blocking our entire development team and affecting customer experience.",
	}

	urgentReport := &domain.Report{
		ID:        "urgent-report-001",
		InquiryID: urgentInquiry.ID,
		Type:      domain.ReportTypeMigration,
		Title:     "🚨 URGENT: Critical Cloud Migration Assessment for TechCorp Solutions",
		Content: `# EXECUTIVE SUMMARY

**PRIORITY LEVEL: CRITICAL - IMMEDIATE ACTION REQUIRED**

TechCorp Solutions Inc. has submitted an urgent request for cloud migration assistance due to critical performance issues affecting their entire development team and customer experience. The client has requested an emergency consultation call this afternoon.

## URGENCY INDICATORS
- **Critical system performance issues** affecting operations
- **Blocking entire development team** - business impact
- **Customer experience degradation** - revenue impact
- **Emergency consultation requested** for same day
- **ASAP migration timeline** required

## CURRENT STATE ASSESSMENT

**Critical Issues Identified:**
- System performance degradation causing business disruption
- Development team productivity blocked
- Customer experience negatively impacted
- Urgent timeline requirements for AWS migration

**Business Impact:**
- **HIGH** - Active operational disruption
- **HIGH** - Development team productivity loss
- **HIGH** - Customer satisfaction risk
- **CRITICAL** - Revenue impact potential

## IMMEDIATE RECOMMENDATIONS

### Phase 1: Emergency Response (Today)
1. **Schedule emergency consultation call within 2 hours**
2. **Conduct rapid performance assessment** of current infrastructure
3. **Identify immediate stabilization measures** to reduce business impact
4. **Prepare emergency migration roadmap** for AWS

### Phase 2: Rapid Migration Planning (Week 1)
1. **Accelerated AWS architecture design** for critical systems
2. **Priority-based migration sequencing** to minimize downtime
3. **Performance optimization strategy** for immediate improvements
4. **Risk mitigation planning** for business continuity

### Phase 3: Expedited Implementation (Weeks 2-4)
1. **Critical system migration** with minimal downtime
2. **Performance monitoring and optimization** throughout migration
3. **Team training and knowledge transfer** for ongoing operations
4. **Post-migration support and optimization**

## NEXT STEPS - IMMEDIATE ACTION REQUIRED

**Within 1 Hour:**
- Contact client immediately at +1-555-0199
- Schedule emergency consultation call
- Assign senior migration architect

**Within 2 Hours:**
- Begin remote performance assessment
- Prepare emergency migration proposal
- Identify immediate stabilization options

**Within 24 Hours:**
- Deliver comprehensive migration roadmap
- Begin AWS environment preparation
- Establish dedicated support channel

## CONTACT INFORMATION

**Primary Contact:** Sarah Johnson  
**Email:** sarah.johnson@techcorp.com  
**Phone:** +1-555-0199  
**Company:** TechCorp Solutions Inc.  
**Services:** Migration, Optimization, Assessment  

**Response Required:** IMMEDIATE (within 1 hour)  
**Business Impact:** CRITICAL  
**Timeline:** URGENT - ASAP migration required`,
		Status:      domain.ReportStatusDraft,
		GeneratedBy: "bedrock-ai-urgent",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test high priority consultant notification
	consultantData := templateService.PrepareConsultantNotificationData(urgentInquiry, urgentReport, true)
	consultantHTML, err := templateService.RenderEmailTemplate(ctx, "consultant_notification", consultantData)
	if err != nil {
		fmt.Printf("❌ Error rendering high priority consultant template: %v\n", err)
		return
	}

	err = os.WriteFile("test_branded_consultant_urgent.html", []byte(consultantHTML), 0644)
	if err != nil {
		fmt.Printf("❌ Error saving urgent consultant email: %v\n", err)
		return
	}
	fmt.Printf("✅ High priority consultant notification generated: test_branded_consultant_urgent.html\n")

	// Test Case 2: Normal Priority Consultant Notification (No Report)
	fmt.Println("\n2. Testing NORMAL priority consultant notification (inquiry only)...")
	
	normalInquiry := &domain.Inquiry{
		ID:       "normal-inquiry-002",
		Name:     "Michael Chen",
		Email:    "m.chen@startupco.io",
		Company:  "StartupCo",
		Phone:    "+1-555-0234",
		Services: []string{"assessment", "architecture"},
		Message:  "Hi there! We're a growing startup and we're looking to get a cloud architecture assessment for our platform. We're currently running on a single server and expecting significant growth in the next 6 months. Would love to schedule a consultation next week to discuss our options and get recommendations for scaling our infrastructure. No immediate rush, but would appreciate your expertise in planning our cloud journey.",
	}

	normalConsultantData := templateService.PrepareConsultantNotificationData(normalInquiry, nil, false)
	normalConsultantHTML, err := templateService.RenderEmailTemplate(ctx, "consultant_notification", normalConsultantData)
	if err != nil {
		fmt.Printf("❌ Error rendering normal consultant template: %v\n", err)
		return
	}

	err = os.WriteFile("test_branded_consultant_normal.html", []byte(normalConsultantHTML), 0644)
	if err != nil {
		fmt.Printf("❌ Error saving normal consultant email: %v\n", err)
		return
	}
	fmt.Printf("✅ Normal priority consultant notification generated: test_branded_consultant_normal.html\n")

	// Test Case 3: Customer Confirmation - Enterprise Client
	fmt.Println("\n3. Testing customer confirmation for enterprise client...")
	
	enterpriseInquiry := &domain.Inquiry{
		ID:       "enterprise-inquiry-003",
		Name:     "Jennifer Rodriguez",
		Email:    "j.rodriguez@globalcorp.com",
		Company:  "Global Corporation Ltd.",
		Phone:    "+1-555-0567",
		Services: []string{"migration", "optimization", "assessment", "architecture"},
		Message:  "We are a large enterprise looking to modernize our entire cloud infrastructure. We currently have a hybrid setup with multiple data centers and are looking to consolidate and optimize our cloud presence across AWS, Azure, and GCP. This is a multi-million dollar initiative that will span 18-24 months. We need a comprehensive assessment and migration strategy.",
	}

	enterpriseCustomerData := struct {
		Name     string
		Company  string
		Services string
		ID       string
	}{
		Name:     enterpriseInquiry.Name,
		Company:  enterpriseInquiry.Company,
		Services: strings.Join(enterpriseInquiry.Services, ", "),
		ID:       enterpriseInquiry.ID,
	}

	enterpriseCustomerHTML, err := templateService.RenderEmailTemplate(ctx, "customer_confirmation", enterpriseCustomerData)
	if err != nil {
		fmt.Printf("❌ Error rendering enterprise customer template: %v\n", err)
		return
	}

	err = os.WriteFile("test_branded_customer_enterprise.html", []byte(enterpriseCustomerHTML), 0644)
	if err != nil {
		fmt.Printf("❌ Error saving enterprise customer email: %v\n", err)
		return
	}
	fmt.Printf("✅ Enterprise customer confirmation generated: test_branded_customer_enterprise.html\n")

	// Test Case 4: Customer Confirmation - Individual/Small Business
	fmt.Println("\n4. Testing customer confirmation for individual client...")
	
	individualInquiry := &domain.Inquiry{
		ID:       "individual-inquiry-004",
		Name:     "Alex Thompson",
		Email:    "alex@freelancer.dev",
		Company:  "", // No company - individual
		Phone:    "",  // No phone
		Services: []string{"assessment"},
		Message:  "I'm a freelance developer and I'm looking to move my personal projects and client work to the cloud. I'm currently using shared hosting but need something more scalable and professional. Looking for guidance on the best approach for a small-scale setup.",
	}

	individualCustomerData := struct {
		Name     string
		Company  string
		Services string
		ID       string
	}{
		Name:     individualInquiry.Name,
		Company:  individualInquiry.Company,
		Services: strings.Join(individualInquiry.Services, ", "),
		ID:       individualInquiry.ID,
	}

	individualCustomerHTML, err := templateService.RenderEmailTemplate(ctx, "customer_confirmation", individualCustomerData)
	if err != nil {
		fmt.Printf("❌ Error rendering individual customer template: %v\n", err)
		return
	}

	err = os.WriteFile("test_branded_customer_individual.html", []byte(individualCustomerHTML), 0644)
	if err != nil {
		fmt.Printf("❌ Error saving individual customer email: %v\n", err)
		return
	}
	fmt.Printf("✅ Individual customer confirmation generated: test_branded_customer_individual.html\n")

	// Summary
	fmt.Println("\n🎉 Branded Email Template Testing Complete!")
	fmt.Println("==========================================")
	fmt.Println("Generated test files:")
	fmt.Println("  📧 test_branded_consultant_urgent.html    - High priority consultant notification with AI report")
	fmt.Println("  📧 test_branded_consultant_normal.html    - Normal priority consultant notification")
	fmt.Println("  📧 test_branded_customer_enterprise.html  - Enterprise customer confirmation")
	fmt.Println("  📧 test_branded_customer_individual.html  - Individual customer confirmation")
	fmt.Println("\n💡 Open these HTML files in a browser to verify:")
	fmt.Println("  ✨ Modern, professional design")
	fmt.Println("  🎨 Consistent branding and colors")
	fmt.Println("  📱 Mobile-responsive layout")
	fmt.Println("  🚨 Priority-based styling (urgent vs normal)")
	fmt.Println("  🏢 Proper handling of company vs individual clients")
	fmt.Println("  📊 Clear information hierarchy")
	fmt.Println("  🎯 Call-to-action clarity")
}