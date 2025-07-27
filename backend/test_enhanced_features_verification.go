package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Enhanced Features Verification Test ===")
	
	// Create test inquiry for healthcare industry
	inquiry := &domain.Inquiry{
		ID:       "test-enhanced-features-001",
		Name:     "Dr. Sarah Johnson",
		Email:    "sarah.johnson@healthtech.com",
		Company:  "HealthTech Medical Systems",
		Phone:    "+1-555-0199",
		Services: []string{"migration", "security", "compliance"},
		Message:  "We need to migrate our patient management system to the cloud with full HIPAA compliance. We handle sensitive PHI data and need multi-cloud options for disaster recovery. This is urgent for our Q2 compliance audit.",
		Priority: domain.PriorityHigh,
		Status:   domain.InquiryStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fmt.Printf("Testing enhanced features with healthcare inquiry\n")
	fmt.Printf("Company: %s\n", inquiry.Company)
	fmt.Printf("Services: %v\n", inquiry.Services)
	fmt.Printf("Industry context: Healthcare (inferred from company name and message)\n\n")

	// Test individual enhanced services
	testPromptArchitect(inquiry)
	testKnowledgeBase()
	testMultiCloudAnalyzer(inquiry)
	testRiskAssessor(inquiry)
	testDocumentationLibrary()
	
	// Test integrated enhanced report generation
	testIntegratedEnhancedReport(inquiry)

	fmt.Println("\n=== Enhanced Features Verification Completed ===")
}

func testPromptArchitect(inquiry *domain.Inquiry) {
	fmt.Println("--- Testing PromptArchitect ---")
	
	promptArchitect := services.NewPromptArchitect()
	
	// Test different prompt types
	ctx := context.Background()
	
	// Test enhanced report prompt
	options := &interfaces.PromptOptions{
		TargetAudience:             "technical",
		IncludeDocumentationLinks:  true,
		IncludeCompetitiveAnalysis: true,
		IncludeRiskAssessment:      true,
		IncludeImplementationSteps: true,
		IndustryContext:            "healthcare",
		CloudProviders:             []string{"AWS", "Azure", "GCP"},
		MaxTokens:                  4000,
	}
	
	prompt, err := promptArchitect.BuildReportPrompt(ctx, inquiry, options)
	if err != nil {
		log.Printf("Error building report prompt: %v", err)
		return
	}
	
	fmt.Printf("✓ Enhanced report prompt generated (%d characters)\n", len(prompt))
	
	// Verify enhanced features are included
	if strings.Contains(prompt, "MULTI-CLOUD ANALYSIS") {
		fmt.Printf("✓ Multi-cloud analysis section included\n")
	}
	if strings.Contains(prompt, "DOCUMENTATION REQUIREMENTS") {
		fmt.Printf("✓ Documentation requirements section included\n")
	}
	if strings.Contains(prompt, "RISK ASSESSMENT") {
		fmt.Printf("✓ Risk assessment section included\n")
	}
	if strings.Contains(prompt, "healthcare") || strings.Contains(prompt, "HIPAA") {
		fmt.Printf("✓ Industry-specific context included\n")
	}
	
	// Test interview guide prompt
	interviewPrompt, err := promptArchitect.BuildInterviewPrompt(ctx, inquiry)
	if err != nil {
		log.Printf("Error building interview prompt: %v", err)
	} else {
		fmt.Printf("✓ Interview guide prompt generated (%d characters)\n", len(interviewPrompt))
	}
	
	// Test risk assessment prompt
	riskPrompt, err := promptArchitect.BuildRiskAssessmentPrompt(ctx, inquiry)
	if err != nil {
		log.Printf("Error building risk assessment prompt: %v", err)
	} else {
		fmt.Printf("✓ Risk assessment prompt generated (%d characters)\n", len(riskPrompt))
	}
	
	// Test competitive analysis prompt
	compPrompt, err := promptArchitect.BuildCompetitiveAnalysisPrompt(ctx, inquiry)
	if err != nil {
		log.Printf("Error building competitive analysis prompt: %v", err)
	} else {
		fmt.Printf("✓ Competitive analysis prompt generated (%d characters)\n", len(compPrompt))
	}
}

func testKnowledgeBase() {
	fmt.Println("\n--- Testing Knowledge Base ---")
	
	kb := services.NewInMemoryKnowledgeBase()
	
	// Test cloud service information
	serviceInfo, err := kb.GetCloudServiceInfo("aws", "ec2")
	if err != nil {
		log.Printf("Error getting service info: %v", err)
	} else {
		fmt.Printf("✓ AWS EC2 service info retrieved: %s\n", serviceInfo.Description)
		fmt.Printf("  Use cases: %v\n", serviceInfo.UseCases)
		fmt.Printf("  Alternatives: %v\n", serviceInfo.Alternatives)
	}
	
	// Test best practices
	bestPractices, err := kb.GetBestPractices("security", "aws")
	if err != nil {
		log.Printf("Error getting best practices: %v", err)
	} else {
		fmt.Printf("✓ Found %d security best practices for AWS\n", len(bestPractices))
		if len(bestPractices) > 0 {
			fmt.Printf("  Example: %s\n", bestPractices[0].Title)
		}
	}
	
	// Test compliance requirements
	complianceReqs, err := kb.GetComplianceRequirements("healthcare")
	if err != nil {
		log.Printf("Error getting compliance requirements: %v", err)
	} else {
		fmt.Printf("✓ Found %d compliance requirements for healthcare\n", len(complianceReqs))
		if len(complianceReqs) > 0 {
			fmt.Printf("  Example: %s (%s)\n", complianceReqs[0].Framework, complianceReqs[0].Severity)
		}
	}
	
	// Test architectural patterns
	patterns, err := kb.GetArchitecturalPatterns("healthcare", "aws")
	if err != nil {
		log.Printf("Error getting architectural patterns: %v", err)
	} else {
		fmt.Printf("✓ Found %d architectural patterns for healthcare on AWS\n", len(patterns))
	}
	
	// Test health check
	if kb.IsHealthy() {
		fmt.Printf("✓ Knowledge base is healthy\n")
	}
}

func testMultiCloudAnalyzer(inquiry *domain.Inquiry) {
	fmt.Println("\n--- Testing Multi-Cloud Analyzer ---")
	
	kb := services.NewInMemoryKnowledgeBase()
	docLib := services.NewDocumentationLibraryService()
	analyzer := services.NewMultiCloudAnalyzerService(kb, docLib)
	
	ctx := context.Background()
	
	// Test provider recommendation
	recommendation, err := analyzer.GetProviderRecommendation(ctx, inquiry)
	if err != nil {
		log.Printf("Error getting provider recommendation: %v", err)
	} else {
		fmt.Printf("✓ Provider recommendation generated\n")
		fmt.Printf("  Recommended: %s\n", recommendation.RecommendedProvider)
		fmt.Printf("  Confidence: %s\n", recommendation.Confidence)
		fmt.Printf("  Reasoning: %v\n", recommendation.Reasoning)
		fmt.Printf("  Alternatives: %v\n", recommendation.AlternativeOptions)
	}
	
	// Test service comparison
	requirement := interfaces.ServiceRequirement{
		Category:     "compute",
		Requirements: []string{"high-availability", "auto-scaling"},
		Performance: interfaces.PerformanceRequirements{
			CPU:          "4 vCPUs",
			Memory:       "16 GB",
			Availability: "99.99%",
		},
		Compliance: []string{"HIPAA", "SOC2"},
		Budget: interfaces.BudgetConstraints{
			MaxMonthlyCost:   1000.0,
			Currency:         "USD",
			CostOptimization: true,
		},
		Industry: "healthcare",
	}
	
	comparison, err := analyzer.CompareServices(ctx, requirement)
	if err != nil {
		log.Printf("Error comparing services: %v", err)
	} else {
		fmt.Printf("✓ Service comparison completed\n")
		fmt.Printf("  Category: %s\n", comparison.Category)
		fmt.Printf("  Providers analyzed: %d\n", len(comparison.Providers))
		fmt.Printf("  Recommendation: %s\n", comparison.Recommendation)
		
		// Show top provider
		if len(comparison.Providers) > 0 {
			top := comparison.Providers[0]
			fmt.Printf("  Top choice: %s %s (score: %.2f)\n", 
				top.Provider, top.ServiceName, top.Score)
		}
	}
}

func testRiskAssessor(inquiry *domain.Inquiry) {
	fmt.Println("\n--- Testing Risk Assessor ---")
	
	kb := services.NewInMemoryKnowledgeBase()
	docLib := services.NewDocumentationLibraryService()
	riskAssessor := services.NewRiskAssessorService(kb, docLib)
	
	ctx := context.Background()
	
	// Create a sample proposed solution
	solution := &interfaces.ProposedSolution{
		ID:             "test-solution-001",
		InquiryID:      inquiry.ID,
		CloudProviders: []string{"aws", "azure"},
		Services: []interfaces.CloudService{
			{
				Provider:     "aws",
				ServiceName:  "EC2",
				ServiceType:  "compute",
				Configuration: make(map[string]interface{}),
				Dependencies: []string{},
				CriticalPath: true,
			},
			{
				Provider:     "aws",
				ServiceName:  "RDS",
				ServiceType:  "database",
				Configuration: make(map[string]interface{}),
				Dependencies: []string{"EC2"},
				CriticalPath: true,
			},
		},
		Architecture: &interfaces.Architecture{
			ID:   "test-arch-001",
			Type: "microservices",
			Components: []interfaces.ArchitectureComponent{
				{
					Name:         "web-tier",
					Type:         "compute",
					Layer:        "presentation",
					Dependencies: []string{"app-tier"},
					Criticality:  "high",
					Configuration: make(map[string]interface{}),
				},
			},
			NetworkTopology: interfaces.NetworkTopology{
				VPCConfiguration: make(map[string]interface{}),
				SubnetStrategy:   "multi-az",
				SecurityGroups:   []interfaces.SecurityGroup{},
				LoadBalancers:    []interfaces.LoadBalancer{},
			},
			DataStorage: []interfaces.DataStorageComponent{
				{
					Type:            "database",
					Provider:        "aws",
					ServiceName:     "RDS",
					DataType:        "personal",
					SensitivityLevel: "high",
					BackupStrategy:  "automated",
					Configuration:   make(map[string]interface{}),
				},
			},
			SecurityLayers:   []interfaces.SecurityLayer{},
			HighAvailability: true,
			DisasterRecovery: false,
		},
		EstimatedCost: "$2000/month",
		Timeline:      "3 months",
	}
	
	// Perform risk assessment
	assessment, err := riskAssessor.AssessRisks(ctx, inquiry, solution)
	if err != nil {
		log.Printf("Error performing risk assessment: %v", err)
	} else {
		fmt.Printf("✓ Risk assessment completed\n")
		fmt.Printf("  Overall risk level: %s\n", assessment.OverallRiskLevel)
		fmt.Printf("  Technical risks: %d\n", len(assessment.TechnicalRisks))
		fmt.Printf("  Security risks: %d\n", len(assessment.SecurityRisks))
		fmt.Printf("  Compliance risks: %d\n", len(assessment.ComplianceRisks))
		fmt.Printf("  Business risks: %d\n", len(assessment.BusinessRisks))
		fmt.Printf("  Mitigation strategies: %d\n", len(assessment.MitigationStrategies))
		fmt.Printf("  Recommended actions: %d\n", len(assessment.RecommendedActions))
		
		// Show sample risks
		if len(assessment.SecurityRisks) > 0 {
			fmt.Printf("  Sample security risk: %s (%s impact)\n", 
				assessment.SecurityRisks[0].Title, assessment.SecurityRisks[0].Impact)
		}
		if len(assessment.ComplianceRisks) > 0 {
			fmt.Printf("  Sample compliance risk: %s (%s)\n", 
				assessment.ComplianceRisks[0].Title, assessment.ComplianceRisks[0].Framework)
		}
	}
}

func testDocumentationLibrary() {
	fmt.Println("\n--- Testing Documentation Library ---")
	
	docLib := services.NewDocumentationLibraryService()
	ctx := context.Background()
	
	// Test getting documentation links
	awsLinks, err := docLib.GetDocumentationLinks(ctx, "aws", "security")
	if err != nil {
		log.Printf("Error getting AWS security links: %v", err)
	} else {
		fmt.Printf("✓ Found %d AWS security documentation links\n", len(awsLinks))
		if len(awsLinks) > 0 {
			fmt.Printf("  Example: %s (%s)\n", awsLinks[0].Title, awsLinks[0].URL)
		}
	}
	
	// Test searching documentation
	searchResults, err := docLib.SearchDocumentation(ctx, "compliance", []string{"aws", "azure"})
	if err != nil {
		log.Printf("Error searching documentation: %v", err)
	} else {
		fmt.Printf("✓ Found %d compliance-related documentation links\n", len(searchResults))
	}
	
	// Test getting links by category
	bestPracticeLinks, err := docLib.GetLinksByCategory(ctx, "best-practices")
	if err != nil {
		log.Printf("Error getting best practice links: %v", err)
	} else {
		fmt.Printf("✓ Found %d best practice documentation links\n", len(bestPracticeLinks))
	}
	
	// Test health check
	if docLib.IsHealthy() {
		fmt.Printf("✓ Documentation library is healthy\n")
	}
	
	// Test statistics
	stats := docLib.GetStats()
	fmt.Printf("✓ Documentation library stats:\n")
	fmt.Printf("  Total links: %d\n", stats.TotalLinks)
	fmt.Printf("  Valid links: %d\n", stats.ValidLinks)
	fmt.Printf("  Links by provider: %v\n", stats.LinksByProvider)
}

func testIntegratedEnhancedReport(inquiry *domain.Inquiry) {
	fmt.Println("\n--- Testing Integrated Enhanced Report Generation ---")
	
	// Create all services
	bedrockService := &MockBedrockService{}
	templateService := &MockTemplateService{}
	pdfService := &MockPDFService{}
	promptArchitect := services.NewPromptArchitect()
	knowledgeBase := services.NewInMemoryKnowledgeBase()
	docLibrary := services.NewDocumentationLibraryService()
	multiCloudAnalyzer := services.NewMultiCloudAnalyzerService(knowledgeBase, docLibrary)
	riskAssessor := services.NewRiskAssessorService(knowledgeBase, docLibrary)

	// Create enhanced report generator
	reportGen := services.NewReportGenerator(
		bedrockService,
		templateService,
		pdfService,
		promptArchitect,
		knowledgeBase,
		multiCloudAnalyzer,
		riskAssessor,
		docLibrary,
	)

	// Generate enhanced report
	ctx := context.Background()
	report, err := reportGen.GenerateReport(ctx, inquiry)
	if err != nil {
		log.Printf("Error generating integrated enhanced report: %v", err)
		return
	}

	fmt.Printf("✓ Integrated enhanced report generated successfully\n")
	fmt.Printf("  Report ID: %s\n", report.ID)
	fmt.Printf("  Title: %s\n", report.Title)
	fmt.Printf("  Type: %s\n", report.Type)
	fmt.Printf("  Generated by: %s\n", report.GeneratedBy)
	fmt.Printf("  Content length: %d characters\n", len(report.Content))

	// Analyze content for enhanced features
	content := report.Content
	enhancedFeatures := 0
	
	if strings.Contains(content, "MULTI-CLOUD ANALYSIS") {
		fmt.Printf("✓ Multi-cloud analysis included in report\n")
		enhancedFeatures++
	}
	if strings.Contains(content, "RISK ASSESSMENT") {
		fmt.Printf("✓ Risk assessment included in report\n")
		enhancedFeatures++
	}
	if strings.Contains(content, "DOCUMENTATION REFERENCES") {
		fmt.Printf("✓ Documentation references included in report\n")
		enhancedFeatures++
	}
	if strings.Contains(content, "KNOWLEDGE BASE CONTEXT") {
		fmt.Printf("✓ Knowledge base context included in report\n")
		enhancedFeatures++
	}
	if strings.Contains(content, "HIPAA") || strings.Contains(content, "healthcare") {
		fmt.Printf("✓ Industry-specific content included in report\n")
		enhancedFeatures++
	}
	if strings.Contains(content, "AWS") && strings.Contains(content, "Azure") {
		fmt.Printf("✓ Multi-provider analysis included in report\n")
		enhancedFeatures++
	}
	
	fmt.Printf("✓ Enhanced features detected: %d/6\n", enhancedFeatures)
	
	if enhancedFeatures >= 4 {
		fmt.Printf("✓ EXCELLENT: Report shows strong integration of enhanced AI assistant capabilities\n")
	} else if enhancedFeatures >= 2 {
		fmt.Printf("✓ GOOD: Report shows moderate integration of enhanced capabilities\n")
	} else {
		fmt.Printf("⚠ LIMITED: Report shows limited integration of enhanced capabilities\n")
	}
}

// Mock services (reusing from previous test)
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	// Generate more sophisticated content based on prompt analysis
	var content string
	
	// Analyze prompt for enhanced features
	hasMultiCloud := strings.Contains(prompt, "MULTI-CLOUD") || strings.Contains(prompt, "multi-cloud")
	hasRiskAssessment := strings.Contains(prompt, "RISK ASSESSMENT") || strings.Contains(prompt, "risk")
	hasDocumentation := strings.Contains(prompt, "DOCUMENTATION") || strings.Contains(prompt, "documentation")
	hasKnowledgeBase := strings.Contains(prompt, "KNOWLEDGE BASE") || strings.Contains(prompt, "best practices")
	hasHealthcare := strings.Contains(prompt, "healthcare") || strings.Contains(prompt, "HIPAA")
	
	// Build enhanced content based on detected features
	content = "# EXECUTIVE SUMMARY\n\n"
	content += "This enhanced cloud consulting report provides comprehensive recommendations for HealthTech Medical Systems' cloud migration initiative with full HIPAA compliance support.\n\n"
	
	if hasHealthcare {
		content += "**INDUSTRY FOCUS: HEALTHCARE** - Specialized recommendations for healthcare industry with PHI data protection.\n\n"
	}
	
	content += "**PRIORITY LEVEL: HIGH PRIORITY** - Urgent Q2 compliance audit timeline detected.\n\n"
	
	content += "# CURRENT STATE ASSESSMENT\n\n"
	content += "HealthTech Medical Systems requires:\n"
	content += "- Patient management system cloud migration\n"
	content += "- Full HIPAA compliance implementation\n"
	content += "- Multi-cloud disaster recovery strategy\n"
	content += "- Q2 compliance audit preparation\n\n"
	
	if hasMultiCloud {
		content += "# MULTI-CLOUD ANALYSIS\n\n"
		content += "**RECOMMENDED CLOUD PROVIDER: AWS**\n\n"
		content += "**REASONING:**\n"
		content += "- Superior healthcare compliance support with HIPAA BAA\n"
		content += "- Comprehensive security services portfolio\n"
		content += "- Strong disaster recovery capabilities\n\n"
		content += "**ALTERNATIVE OPTIONS:** Azure (strong hybrid capabilities), GCP (competitive pricing)\n\n"
		content += "**COST IMPLICATIONS:** Estimated 15-20% premium for healthcare-specific compliance features\n\n"
	}
	
	content += "# RECOMMENDATIONS\n\n"
	content += "## 1. Primary Cloud Strategy\n"
	content += "- **AWS as primary provider** for healthcare workloads\n"
	content += "- **Azure as secondary** for disaster recovery\n"
	content += "- Implement cross-cloud data replication\n\n"
	
	content += "## 2. HIPAA Compliance Implementation\n"
	content += "- Deploy AWS HealthLake for FHIR compliance\n"
	content += "- Implement comprehensive audit logging\n"
	content += "- Establish Business Associate Agreements\n"
	content += "- Enable end-to-end encryption for PHI data\n\n"
	
	if hasRiskAssessment {
		content += "# RISK ASSESSMENT\n\n"
		content += "**OVERALL RISK LEVEL: MEDIUM**\n\n"
		content += "**KEY TECHNICAL RISKS:**\n"
		content += "- Data migration complexity (medium impact)\n"
		content += "- Multi-cloud integration challenges (medium impact)\n\n"
		content += "**KEY SECURITY RISKS:**\n"
		content += "- PHI data exposure during migration (high impact)\n"
		content += "- Inadequate access controls (high impact)\n\n"
		content += "**KEY COMPLIANCE RISKS:**\n"
		content += "- HIPAA violation during transition (critical impact)\n"
		content += "- Audit trail gaps (high impact)\n\n"
		content += "**RECOMMENDED RISK MITIGATION ACTIONS:**\n"
		content += "- Implement comprehensive encryption strategy\n"
		content += "- Establish detailed audit logging\n"
		content += "- Deploy multi-factor authentication\n\n"
	}
	
	if hasKnowledgeBase {
		content += "# KNOWLEDGE BASE INSIGHTS\n\n"
		content += "**BEST PRACTICES FOR HEALTHCARE MIGRATION:**\n"
		content += "- Implement zero-trust network architecture\n"
		content += "- Use dedicated tenancy for sensitive workloads\n"
		content += "- Deploy comprehensive monitoring and alerting\n\n"
		content += "**COMPLIANCE REQUIREMENTS FOR HEALTHCARE INDUSTRY:**\n"
		content += "- HIPAA (critical): Protect PHI data with encryption and access controls\n"
		content += "- HITECH (high): Implement breach notification procedures\n"
		content += "- SOC2 (medium): Establish security and availability controls\n\n"
	}
	
	if hasDocumentation {
		content += "# DOCUMENTATION REFERENCES\n\n"
		content += "**HEALTHCARE COMPLIANCE DOCUMENTATION:**\n"
		content += "- AWS HIPAA Compliance Guide: https://aws.amazon.com/compliance/hipaa-compliance/\n"
		content += "- Azure Healthcare Compliance: https://docs.microsoft.com/en-us/azure/compliance/\n"
		content += "- NIST Cloud Security Framework: https://www.nist.gov/programs-projects/nist-cloud-computing-program-nccp\n\n"
		content += "**BEST PRACTICES DOCUMENTATION:**\n"
		content += "- AWS Well-Architected Framework: https://docs.aws.amazon.com/wellarchitected/\n"
		content += "- Azure Architecture Center: https://docs.microsoft.com/en-us/azure/architecture/\n"
		content += "- Multi-Cloud Security Best Practices: https://cloud.google.com/security/best-practices\n\n"
	}
	
	content += "# IMPLEMENTATION ROADMAP\n\n"
	content += "## Phase 1: Foundation (Weeks 1-4)\n"
	content += "- Establish AWS and Azure environments\n"
	content += "- Implement basic security controls\n"
	content += "- Set up compliance monitoring\n\n"
	content += "## Phase 2: Migration (Weeks 5-8)\n"
	content += "- Migrate non-critical systems first\n"
	content += "- Implement data replication\n"
	content += "- Test disaster recovery procedures\n\n"
	content += "## Phase 3: Optimization (Weeks 9-12)\n"
	content += "- Fine-tune performance and costs\n"
	content += "- Complete compliance documentation\n"
	content += "- Prepare for Q2 audit\n\n"
	
	content += "# NEXT STEPS\n\n"
	content += "1. **IMMEDIATE ACTIONS (This Week)**\n"
	content += "   - Schedule technical discovery workshop\n"
	content += "   - Begin compliance gap analysis\n"
	content += "   - Establish project timeline\n\n"
	content += "2. **SHORT TERM (2-4 Weeks)**\n"
	content += "   - Finalize multi-cloud architecture\n"
	content += "   - Establish cloud environments\n"
	content += "   - Begin pilot migration\n\n"
	
	content += "**MEETING SCHEDULING:** Urgent Q2 compliance audit timeline requires immediate engagement.\n\n"
	
	content += "# URGENCY ASSESSMENT\n\n"
	content += "**Urgent language detected:** \"urgent\", \"critical for our Q2 compliance audit\"\n"
	content += "**Compliance deadline:** Q2 audit timeline\n"
	content += "**Recommended response timeline:** Within 24 hours\n"
	content += "**Business impact:** High - compliance audit failure could result in significant penalties\n\n"
	
	content += "# CONTACT INFORMATION\n\n"
	content += "- Client: Dr. Sarah Johnson (sarah.johnson@healthtech.com)\n"
	content += "- Company: HealthTech Medical Systems\n"
	content += "- Phone: +1-555-0199\n"
	
	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(prompt) / 4,
			OutputTokens: len(content) / 4,
		},
		Metadata: map[string]string{
			"model": options.ModelID,
		},
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "amazon.nova-lite-v1:0",
		ModelName:   "Nova Lite",
		Provider:    "Amazon",
		MaxTokens:   4000,
		IsAvailable: true,
	}
}

func (m *MockBedrockService) IsHealthy() bool {
	return true
}

type MockTemplateService struct{}

func (m *MockTemplateService) RenderReportTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><title>Enhanced Report</title></head><body>%s</body></html>`, data), nil
}

func (m *MockTemplateService) RenderEmailTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	return "Mock email", nil
}

func (m *MockTemplateService) LoadTemplate(templateName string) (*template.Template, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MockTemplateService) ValidateTemplate(templateContent string) error { return nil }
func (m *MockTemplateService) GetAvailableTemplates() []string { return []string{} }
func (m *MockTemplateService) ReloadTemplates() error { return nil }
func (m *MockTemplateService) PrepareReportTemplateData(inquiry *domain.Inquiry, report *domain.Report) interface{} { return "" }
func (m *MockTemplateService) PrepareConsultantNotificationData(inquiry *domain.Inquiry, report *domain.Report, isHighPriority bool) interface{} { return "" }

type MockPDFService struct{}

func (m *MockPDFService) GeneratePDF(ctx context.Context, htmlContent string, options *interfaces.PDFOptions) ([]byte, error) {
	return []byte("Mock PDF"), nil
}

func (m *MockPDFService) GeneratePDFFromURL(ctx context.Context, url string, options *interfaces.PDFOptions) ([]byte, error) {
	return []byte("Mock PDF"), nil
}

func (m *MockPDFService) IsHealthy() bool { return true }
func (m *MockPDFService) GetVersion() string { return "mock-1.0.0" }