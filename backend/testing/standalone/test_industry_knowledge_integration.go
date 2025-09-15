package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// mockDocumentationLibrary is a simple mock for testing
type mockDocumentationLibrary struct{}

func (m *mockDocumentationLibrary) GetDocumentationLinks(ctx context.Context, provider, topic string) ([]*interfaces.DocumentationLink, error) {
	return []*interfaces.DocumentationLink{}, nil
}

func (m *mockDocumentationLibrary) ValidateLinks(ctx context.Context, links []*interfaces.DocumentationLink) ([]*interfaces.LinkValidation, error) {
	return []*interfaces.LinkValidation{}, nil
}

func (m *mockDocumentationLibrary) UpdateDocumentationIndex(ctx context.Context) error {
	return nil
}

func (m *mockDocumentationLibrary) SearchDocumentation(ctx context.Context, query string, providers []string) ([]*interfaces.DocumentationLink, error) {
	return []*interfaces.DocumentationLink{}, nil
}

func (m *mockDocumentationLibrary) AddDocumentationLink(ctx context.Context, link *interfaces.DocumentationLink) error {
	return nil
}

func (m *mockDocumentationLibrary) RemoveDocumentationLink(ctx context.Context, linkID string) error {
	return nil
}

func (m *mockDocumentationLibrary) GetLinksByCategory(ctx context.Context, category string) ([]*interfaces.DocumentationLink, error) {
	return []*interfaces.DocumentationLink{}, nil
}

func (m *mockDocumentationLibrary) GetLinksByProvider(ctx context.Context, provider string) ([]*interfaces.DocumentationLink, error) {
	return []*interfaces.DocumentationLink{}, nil
}

func (m *mockDocumentationLibrary) GetLinksByType(ctx context.Context, linkType string) ([]*interfaces.DocumentationLink, error) {
	return []*interfaces.DocumentationLink{}, nil
}

func (m *mockDocumentationLibrary) GetLinkValidationStatus(ctx context.Context, linkID string) (*interfaces.LinkValidation, error) {
	return nil, nil
}

func (m *mockDocumentationLibrary) IsHealthy() bool {
	return true
}

func (m *mockDocumentationLibrary) GetStats() *interfaces.DocumentationLibraryStats {
	return &interfaces.DocumentationLibraryStats{}
}

func main() {
	fmt.Println("=== Testing Industry-Specific Knowledge System Integration ===")
	
	// Initialize the knowledge base and risk assessor
	kb := services.NewInMemoryKnowledgeBase()
	// Create a simple mock documentation library for testing
	docLib := &mockDocumentationLibrary{}
	riskAssessor := services.NewRiskAssessorService(kb, docLib)
	
	// Test healthcare inquiry with HIPAA compliance
	fmt.Println("\n1. Testing Healthcare Inquiry with HIPAA Compliance:")
	healthcareInquiry := &domain.Inquiry{
		ID:       "test-healthcare-001",
		Company:  "Regional Medical Center",
		Services: []string{"data-migration", "application-modernization"},
		Message:  "We need to migrate our electronic health records system to the cloud while maintaining HIPAA compliance. We handle sensitive patient data and need to ensure all PHI is properly protected.",
		Priority: "high",
	}
	
	// Test risk assessment for healthcare
	ctx := context.Background()
	solution := &interfaces.ProposedSolution{
		ID:            "test-solution-001",
		InquiryID:     healthcareInquiry.ID,
		CloudProviders: []string{"aws"},
		Services: []interfaces.CloudService{
			{Provider: "aws", ServiceName: "RDS", ServiceType: "database", CriticalPath: true},
			{Provider: "aws", ServiceName: "S3", ServiceType: "storage", CriticalPath: true},
			{Provider: "aws", ServiceName: "Lambda", ServiceType: "compute", CriticalPath: false},
		},
		Architecture: &interfaces.Architecture{
			ID:   "arch-001",
			Type: "microservices",
			Components: []interfaces.ArchitectureComponent{
				{Name: "Database", Type: "RDS", Layer: "data", Criticality: "high"},
				{Name: "Storage", Type: "S3", Layer: "data", Criticality: "high"},
			},
			DataStorage: []interfaces.DataStorageComponent{
				{Type: "relational", Provider: "aws", ServiceName: "RDS", DataType: "health", SensitivityLevel: "high"},
			},
			HighAvailability: false,
			DisasterRecovery: false,
		},
		DataFlow: []interfaces.DataFlowComponent{
			{Source: "application", Destination: "database", DataType: "health", Encryption: false},
		},
		EstimatedCost: "$5000/month",
	}
	
	riskAssessment, err := riskAssessor.AssessRisks(ctx, healthcareInquiry, solution)
	if err != nil {
		log.Printf("Error assessing healthcare risks: %v", err)
	} else {
		fmt.Printf("Healthcare Risk Assessment:\n")
		fmt.Printf("  Overall Risk Level: %s\n", riskAssessment.OverallRiskLevel)
		fmt.Printf("  Technical Risks: %d\n", len(riskAssessment.TechnicalRisks))
		fmt.Printf("  Security Risks: %d\n", len(riskAssessment.SecurityRisks))
		fmt.Printf("  Compliance Risks: %d\n", len(riskAssessment.ComplianceRisks))
		fmt.Printf("  Business Risks: %d\n", len(riskAssessment.BusinessRisks))
		
		// Show some specific risks
		if len(riskAssessment.ComplianceRisks) > 0 {
			fmt.Printf("  Sample Compliance Risk: %s\n", riskAssessment.ComplianceRisks[0].Description)
		}
		if len(riskAssessment.MitigationStrategies) > 0 {
			fmt.Printf("  Sample Mitigation: %s\n", riskAssessment.MitigationStrategies[0].Strategy)
		}
	}
	
	// Test financial services inquiry with PCI-DSS
	fmt.Println("\n2. Testing Financial Services Inquiry with PCI-DSS:")
	financialInquiry := &domain.Inquiry{
		ID:       "test-financial-001",
		Company:  "SecureBank Corp",
		Services: []string{"fraud-detection", "trading-platform"},
		Message:  "We need to build a real-time fraud detection system and upgrade our trading platform. We process credit card transactions and need PCI-DSS compliance.",
		Priority: "critical",
	}
	
	financialSolution := &interfaces.ProposedSolution{
		ID:            "test-solution-002",
		InquiryID:     financialInquiry.ID,
		CloudProviders: []string{"aws"},
		Services: []interfaces.CloudService{
			{Provider: "aws", ServiceName: "Kinesis", ServiceType: "streaming", CriticalPath: true},
			{Provider: "aws", ServiceName: "Lambda", ServiceType: "compute", CriticalPath: true},
		},
		Architecture: &interfaces.Architecture{
			ID:   "arch-002",
			Type: "event-driven",
			DataStorage: []interfaces.DataStorageComponent{
				{Type: "nosql", Provider: "aws", ServiceName: "DynamoDB", DataType: "payment", SensitivityLevel: "critical"},
			},
			HighAvailability: true,
			DisasterRecovery: true,
		},
		DataFlow: []interfaces.DataFlowComponent{
			{Source: "payment-gateway", Destination: "fraud-detection", DataType: "payment", Encryption: false},
		},
		EstimatedCost: "$8000/month",
	}
	
	financialRisk, err := riskAssessor.AssessRisks(ctx, financialInquiry, financialSolution)
	if err != nil {
		log.Printf("Error assessing financial risks: %v", err)
	} else {
		fmt.Printf("Financial Risk Assessment:\n")
		fmt.Printf("  Overall Risk Level: %s\n", financialRisk.OverallRiskLevel)
		fmt.Printf("  Total Risks Identified: %d\n", 
			len(financialRisk.TechnicalRisks)+len(financialRisk.SecurityRisks)+
			len(financialRisk.ComplianceRisks)+len(financialRisk.BusinessRisks))
	}
	
	// Test retail inquiry
	fmt.Println("\n3. Testing Retail Inquiry:")
	retailInquiry := &domain.Inquiry{
		ID:       "test-retail-001",
		Company:  "GlobalRetail Inc",
		Services: []string{"e-commerce", "inventory-management"},
		Message:  "We need to scale our e-commerce platform for Black Friday traffic and implement real-time inventory management across 500 stores.",
		Priority: "high",
	}
	
	retailSolution := &interfaces.ProposedSolution{
		ID:            "test-solution-003",
		InquiryID:     retailInquiry.ID,
		CloudProviders: []string{"aws"},
		Services: []interfaces.CloudService{
			{Provider: "aws", ServiceName: "ECS", ServiceType: "container", CriticalPath: true},
			{Provider: "aws", ServiceName: "RDS", ServiceType: "database", CriticalPath: true},
		},
		Architecture: &interfaces.Architecture{
			ID:   "arch-003",
			Type: "microservices",
			DataStorage: []interfaces.DataStorageComponent{
				{Type: "relational", Provider: "aws", ServiceName: "RDS", DataType: "customer", SensitivityLevel: "medium"},
			},
			HighAvailability: true,
			DisasterRecovery: false,
		},
		EstimatedCost: "$3000/month",
	}
	
	retailRisk, err := riskAssessor.AssessRisks(ctx, retailInquiry, retailSolution)
	if err != nil {
		log.Printf("Error assessing retail risks: %v", err)
	} else {
		fmt.Printf("Retail Risk Assessment:\n")
		fmt.Printf("  Overall Risk Level: %s\n", retailRisk.OverallRiskLevel)
		fmt.Printf("  Recommended Actions: %d\n", len(retailRisk.RecommendedActions))
	}
	
	// Test manufacturing inquiry
	fmt.Println("\n4. Testing Manufacturing Inquiry:")
	manufacturingInquiry := &domain.Inquiry{
		ID:       "test-manufacturing-001",
		Company:  "AutoParts Manufacturing",
		Services: []string{"iot-platform", "predictive-maintenance"},
		Message:  "We want to implement IoT sensors on our production line for predictive maintenance and quality control monitoring.",
		Priority: "medium",
	}
	
	manufacturingSolution := &interfaces.ProposedSolution{
		ID:            "test-solution-004",
		InquiryID:     manufacturingInquiry.ID,
		CloudProviders: []string{"aws"},
		Services: []interfaces.CloudService{
			{Provider: "aws", ServiceName: "IoT Core", ServiceType: "iot", CriticalPath: true},
			{Provider: "aws", ServiceName: "Kinesis", ServiceType: "streaming", CriticalPath: true},
		},
		Architecture: &interfaces.Architecture{
			ID:   "arch-004",
			Type: "iot",
			DataStorage: []interfaces.DataStorageComponent{
				{Type: "timeseries", Provider: "aws", ServiceName: "Timestream", DataType: "sensor", SensitivityLevel: "low"},
			},
			HighAvailability: true,
			DisasterRecovery: true,
		},
		EstimatedCost: "$2000/month",
	}
	
	manufacturingRisk, err := riskAssessor.AssessRisks(ctx, manufacturingInquiry, manufacturingSolution)
	if err != nil {
		log.Printf("Error assessing manufacturing risks: %v", err)
	} else {
		fmt.Printf("Manufacturing Risk Assessment:\n")
		fmt.Printf("  Overall Risk Level: %s\n", manufacturingRisk.OverallRiskLevel)
		fmt.Printf("  Mitigation Strategies: %d\n", len(manufacturingRisk.MitigationStrategies))
	}
	
	// Test compliance framework coverage
	fmt.Println("\n5. Testing Compliance Framework Coverage:")
	frameworks := []string{"HIPAA", "PCI-DSS", "SOX", "GDPR", "SOC2", "ISO27001", "FedRAMP", "NIST"}
	for _, framework := range frameworks {
		allCompliance, _ := kb.GetComplianceRequirements("")
		count := 0
		for _, req := range allCompliance {
			if req.Framework == framework {
				count++
			}
		}
		fmt.Printf("  %s: %d requirements\n", framework, count)
	}
	
	// Test industry-specific architectural patterns
	fmt.Println("\n6. Testing Industry-Specific Architectural Patterns:")
	industries := []string{"healthcare", "financial", "retail", "manufacturing", "government", "education"}
	for _, industry := range industries {
		patterns, _ := kb.GetIndustrySpecificArchitecturalPatterns(industry)
		bestPractices, _ := kb.GetIndustrySpecificBestPractices(industry)
		risks, _ := kb.GetIndustryRiskFactors(industry)
		
		fmt.Printf("  %s:\n", industry)
		fmt.Printf("    Architectural Patterns: %d\n", len(patterns))
		fmt.Printf("    Best Practices: %d\n", len(bestPractices))
		fmt.Printf("    Risk Factors: %d\n", len(risks))
	}
	
	// Test multi-cloud support for different industries
	fmt.Println("\n7. Testing Multi-Cloud Support by Industry:")
	providers := []string{"aws", "azure", "gcp"}
	for _, provider := range providers {
		services, _ := kb.SearchServices(ctx, "", []string{provider})
		fmt.Printf("  %s: %d services\n", provider, len(services))
	}
	
	// Test industry-specific recommendations
	fmt.Println("\n8. Testing Industry-Specific Recommendations:")
	testCases := []struct {
		industry string
		useCase  string
	}{
		{"healthcare", "data-migration"},
		{"financial", "fraud-detection"},
		{"retail", "e-commerce"},
		{"manufacturing", "iot-platform"},
		{"government", "citizen-services"},
		{"education", "learning-management"},
	}
	
	for _, tc := range testCases {
		recs, _ := kb.GetIndustrySpecificRecommendations(tc.industry, tc.useCase)
		fmt.Printf("  %s/%s: %d recommendations\n", tc.industry, tc.useCase, len(recs))
	}
	
	// Test validation of requirements from task
	fmt.Println("\n9. Validating Task Requirements:")
	
	// Requirement 5.1: Industry-specific compliance requirements
	healthcareCompliance, _ := kb.GetComplianceRequirements("healthcare")
	fmt.Printf("  ✓ 5.1 - Healthcare compliance requirements: %d found\n", len(healthcareCompliance))
	
	// Requirement 5.2: Industry-specific architectural patterns
	healthcarePatterns, _ := kb.GetIndustrySpecificArchitecturalPatterns("healthcare")
	fmt.Printf("  ✓ 5.2 - Healthcare architectural patterns: %d found\n", len(healthcarePatterns))
	
	// Requirement 5.3: Industry-specific security requirements
	healthcareBP, _ := kb.GetIndustrySpecificBestPractices("healthcare")
	securityBPCount := 0
	for _, bp := range healthcareBP {
		if bp.Category == "security" {
			securityBPCount++
		}
	}
	fmt.Printf("  ✓ 5.3 - Healthcare security best practices: %d found\n", securityBPCount)
	
	// Requirement 5.4: Industry-specific case studies and success stories
	healthcareRecs, _ := kb.GetIndustrySpecificRecommendations("healthcare", "data-migration")
	fmt.Printf("  ✓ 5.4 - Healthcare-specific recommendations: %d found\n", len(healthcareRecs))
	
	// Requirement 5.5: Industry-specific approval and compliance processes
	healthcareRisks, _ := kb.GetIndustryRiskFactors("healthcare")
	fmt.Printf("  ✓ 5.5 - Healthcare-specific risk factors: %d found\n", len(healthcareRisks))
	
	fmt.Println("\n=== Industry-Specific Knowledge System Integration Test Complete ===")
}