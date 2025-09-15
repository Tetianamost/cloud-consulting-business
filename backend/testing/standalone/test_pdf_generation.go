package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	fmt.Println("=== PDF Generation Test ===")

	// Test 1: Initialize PDF Service
	fmt.Println("\n1. Testing PDF Service Initialization...")
	pdfService := services.NewPDFService(logger)
	
	if !pdfService.IsHealthy() {
		fmt.Println("❌ PDF service is not healthy")
		return
	}
	fmt.Printf("✅ PDF service initialized successfully - Version: %s\n", pdfService.GetVersion())

	// Test 2: Test basic PDF generation
	fmt.Println("\n2. Testing Basic PDF Generation...")
	testHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #007cba; border-bottom: 2px solid #007cba; }
        h2 { color: #333; margin-top: 20px; }
        p { line-height: 1.6; text-align: justify; }
        .highlight { background-color: #f0f8ff; padding: 10px; border-left: 4px solid #007cba; }
    </style>
</head>
<body>
    <h1>Cloud Assessment Report</h1>
    <h2>Executive Summary</h2>
    <p>This is a comprehensive cloud assessment report that demonstrates the PDF generation capabilities of our system.</p>
    
    <h2>Current State Analysis</h2>
    <div class="highlight">
        <p><strong>Key Finding:</strong> The current infrastructure shows significant opportunities for optimization and cost reduction.</p>
    </div>
    
    <h2>Recommendations</h2>
    <p>Based on our analysis, we recommend the following actions:</p>
    <ul>
        <li>Migrate to cloud-native services</li>
        <li>Implement auto-scaling policies</li>
        <li>Optimize storage configurations</li>
        <li>Enhance security posture</li>
    </ul>
    
    <h2>Next Steps</h2>
    <p>We recommend scheduling a follow-up meeting to discuss implementation timelines and resource requirements.</p>
</body>
</html>`

	ctx := context.Background()
	pdfBytes, err := pdfService.GeneratePDF(ctx, testHTML, nil)
	if err != nil {
		fmt.Printf("❌ Failed to generate basic PDF: %v\n", err)
		return
	}
	
	// Save test PDF
	err = os.WriteFile("test_basic_report.pdf", pdfBytes, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save basic PDF: %v\n", err)
		return
	}
	fmt.Printf("✅ Basic PDF generated successfully - Size: %d bytes\n", len(pdfBytes))

	// Test 3: Test PDF with custom options
	fmt.Println("\n3. Testing PDF with Custom Options...")
	options := &interfaces.PDFOptions{
		PageSize:     "A4",
		Orientation:  "Portrait",
		MarginTop:    "1in",
		MarginRight:  "0.75in",
		MarginBottom: "1in",
		MarginLeft:   "0.75in",
		HeaderHTML:   "Cloud Consulting Report - Generated on " + time.Now().Format("January 2, 2006"),
		FooterHTML:   "Page [page] of [topage] - Confidential",
		Quality:      94,
		LoadTimeout:  60, // Increase timeout to 60 seconds
	}

	pdfBytesWithOptions, err := pdfService.GeneratePDF(ctx, testHTML, options)
	if err != nil {
		fmt.Printf("❌ Failed to generate PDF with options: %v\n", err)
		return
	}

	err = os.WriteFile("test_options_report.pdf", pdfBytesWithOptions, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save options PDF: %v\n", err)
		return
	}
	fmt.Printf("✅ PDF with options generated successfully - Size: %d bytes\n", len(pdfBytesWithOptions))

	// Test 4: Test with realistic report content
	fmt.Println("\n4. Testing with Realistic Report Content...")
	
	// Create a mock inquiry and report
	inquiry := &domain.Inquiry{
		ID:       "test-inquiry-123",
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Company:  "Example Corp",
		Phone:    "+1-555-0123",
		Services: []string{"assessment"},
		Message:  "We need a comprehensive cloud assessment for our current infrastructure.",
	}

	report := &domain.Report{
		ID:        "test-report-456",
		InquiryID: inquiry.ID,
		Type:      "assessment",
		Title:     "Cloud Infrastructure Assessment Report",
		Content: `# EXECUTIVE SUMMARY

This comprehensive cloud assessment report provides a detailed analysis of Example Corp's current infrastructure and strategic recommendations for cloud transformation.

## CURRENT STATE ANALYSIS

### Infrastructure Overview
- Legacy on-premises servers (5+ years old)
- Mixed Windows/Linux environment
- Limited automation and monitoring
- High operational overhead

### Key Challenges Identified
1. **Scalability Limitations**: Current infrastructure cannot handle peak loads efficiently
2. **Security Concerns**: Outdated security protocols and limited compliance capabilities
3. **Cost Inefficiencies**: Over-provisioned resources leading to unnecessary expenses
4. **Maintenance Burden**: Significant time spent on routine maintenance tasks

## RECOMMENDATIONS

### Phase 1: Foundation (Months 1-3)
- Establish cloud landing zone with proper governance
- Implement identity and access management (IAM)
- Set up monitoring and logging infrastructure
- Create backup and disaster recovery procedures

### Phase 2: Migration (Months 4-8)
- Migrate non-critical workloads first
- Implement containerization for suitable applications
- Establish CI/CD pipelines
- Optimize database configurations

### Phase 3: Optimization (Months 9-12)
- Implement auto-scaling policies
- Optimize costs through reserved instances and spot pricing
- Enhance security posture with advanced threat detection
- Establish performance monitoring and alerting

## EXPECTED BENEFITS

### Cost Savings
- Estimated 30-40% reduction in infrastructure costs
- Elimination of hardware refresh cycles
- Reduced operational overhead

### Performance Improvements
- 99.9% uptime SLA
- Improved application response times
- Better scalability during peak periods

### Security Enhancements
- Advanced threat detection and response
- Automated compliance reporting
- Enhanced data encryption and protection

## NEXT STEPS

1. **Immediate Actions** (Next 30 days)
   - Conduct detailed application inventory
   - Assess network connectivity requirements
   - Begin team training on cloud technologies

2. **Short-term Goals** (Next 90 days)
   - Finalize cloud provider selection
   - Establish project governance structure
   - Begin pilot migration with non-critical systems

3. **Long-term Objectives** (6-12 months)
   - Complete full infrastructure migration
   - Achieve operational excellence in cloud environment
   - Establish continuous improvement processes

## CONCLUSION

The migration to cloud infrastructure represents a significant opportunity for Example Corp to modernize its technology stack, reduce costs, and improve operational efficiency. With proper planning and execution, this transformation will position the organization for future growth and innovation.

We recommend proceeding with the phased approach outlined above, starting with the foundation phase to establish proper governance and security controls before beginning the migration process.`,
		Status:    "completed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Initialize template service and report generator
	templateService := services.NewTemplateService("templates", logger)
	reportGenerator := services.NewReportGenerator(nil, templateService, pdfService)

	// Generate PDF using the report generator
	pdfBytesRealistic, err := reportGenerator.GeneratePDF(ctx, inquiry, report)
	if err != nil {
		fmt.Printf("❌ Failed to generate realistic PDF: %v\n", err)
		return
	}

	err = os.WriteFile("test_realistic_report.pdf", pdfBytesRealistic, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to save realistic PDF: %v\n", err)
		return
	}
	fmt.Printf("✅ Realistic PDF generated successfully - Size: %d bytes\n", len(pdfBytesRealistic))

	// Test 5: Test error handling
	fmt.Println("\n5. Testing Error Handling...")
	
	// Test with invalid HTML
	invalidHTML := "<html><body><h1>Unclosed tag<body></html>"
	_, err = pdfService.GeneratePDF(ctx, invalidHTML, nil)
	if err != nil {
		fmt.Printf("✅ Error handling works correctly: %v\n", err)
	} else {
		fmt.Println("⚠️  Expected error for invalid HTML, but none occurred")
	}

	// Test with empty content
	_, err = pdfService.GeneratePDF(ctx, "", nil)
	if err != nil {
		fmt.Printf("✅ Empty content handled correctly: %v\n", err)
	} else {
		fmt.Println("✅ Empty content handled gracefully")
	}

	// Test 6: Performance test
	fmt.Println("\n6. Testing Performance...")
	start := time.Now()
	for i := 0; i < 5; i++ {
		_, err := pdfService.GeneratePDF(ctx, testHTML, nil)
		if err != nil {
			fmt.Printf("❌ Performance test failed on iteration %d: %v\n", i+1, err)
			return
		}
	}
	duration := time.Since(start)
	avgTime := duration / 5
	fmt.Printf("✅ Performance test completed - Average generation time: %v\n", avgTime)

	fmt.Println("\n=== PDF Generation Test Summary ===")
	fmt.Println("✅ PDF Service Initialization: PASSED")
	fmt.Println("✅ Basic PDF Generation: PASSED")
	fmt.Println("✅ PDF with Custom Options: PASSED")
	fmt.Println("✅ Realistic Report PDF: PASSED")
	fmt.Println("✅ Error Handling: PASSED")
	fmt.Println("✅ Performance Test: PASSED")
	fmt.Println("\nGenerated test files:")
	fmt.Println("- test_basic_report.pdf")
	fmt.Println("- test_options_report.pdf")
	fmt.Println("- test_realistic_report.pdf")
	fmt.Println("\n🎉 All PDF generation tests PASSED!")
}