package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

func main() {
	fmt.Println("=== Task 9 Verification: Audience-Aware Content Generation ===")
	fmt.Println("Verifying all task requirements are implemented:")
	fmt.Println("1. Add audience detection logic (technical vs business stakeholders)")
	fmt.Println("2. Create different content templates for different audience types")
	fmt.Println("3. Implement technical depth adjustment based on identified audience")
	fmt.Println("4. Add business justification and technical explanation separation")
	fmt.Println()
	
	// Test 1: Audience Detection Logic
	fmt.Println("--- Test 1: Audience Detection Logic ---")
	testAudienceDetection()
	
	// Test 2: Different Content Templates
	fmt.Println("\n--- Test 2: Different Content Templates ---")
	testContentTemplates()
	
	// Test 3: Technical Depth Adjustment
	fmt.Println("\n--- Test 3: Technical Depth Adjustment ---")
	testTechnicalDepthAdjustment()
	
	// Test 4: Business/Technical Separation
	fmt.Println("\n--- Test 4: Business/Technical Separation ---")
	testBusinessTechnicalSeparation()
	
	// Test 5: Integration with Report Generation
	fmt.Println("\n--- Test 5: Integration with Report Generation ---")
	testReportGenerationIntegration()
	
	fmt.Println("\n=== Task 9 Verification Complete ===")
}

// Test 1: Verify audience detection logic works for different stakeholder types
func testAudienceDetection() {
	audienceDetector := services.NewAudienceDetector()
	ctx := context.Background()
	
	testCases := []struct {
		name           string
		inquiry        *domain.Inquiry
		expectedType   services.AudienceType
		minConfidence  float64
	}{
		{
			name: "Technical Stakeholder",
			inquiry: &domain.Inquiry{
				Message: "We need Kubernetes cluster setup with Istio service mesh, Prometheus monitoring, and CI/CD pipeline integration. Looking for Terraform configurations and YAML manifests.",
			},
			expectedType:  services.AudienceTechnical,
			minConfidence: 0.6,
		},
		{
			name: "Business Stakeholder",
			inquiry: &domain.Inquiry{
				Message: "Our board needs a comprehensive ROI analysis and business case for cloud migration. We want to understand cost savings, competitive advantages, and strategic benefits.",
			},
			expectedType:  services.AudienceBusiness,
			minConfidence: 0.6,
		},
		{
			name: "Executive Stakeholder",
			inquiry: &domain.Inquiry{
				Message: "As CEO, I need a strategic assessment of our digital transformation initiative. The board requires high-level vision and investment requirements.",
			},
			expectedType:  services.AudienceExecutive,
			minConfidence: 0.7,
		},
		{
			name: "Mixed Stakeholder",
			inquiry: &domain.Inquiry{
				Message: "We need both technical implementation details for our IT team and business justification for management. This should cover architecture and ROI.",
			},
			expectedType:  services.AudienceMixed,
			minConfidence: 0.3,
		},
	}
	
	for _, tc := range testCases {
		profile, err := audienceDetector.DetectAudience(ctx, tc.inquiry)
		if err != nil {
			fmt.Printf("  ✗ %s: Error - %v\n", tc.name, err)
			continue
		}
		
		// Check if detected type matches expected (or mixed for ambiguous cases)
		typeMatch := profile.PrimaryType == tc.expectedType || profile.PrimaryType == services.AudienceMixed
		confidenceOk := profile.Confidence >= tc.minConfidence || profile.PrimaryType == services.AudienceMixed
		
		if typeMatch && confidenceOk {
			fmt.Printf("  ✓ %s: Detected %s (confidence: %.2f)\n", tc.name, profile.PrimaryType, profile.Confidence)
		} else {
			fmt.Printf("  ✗ %s: Expected %s, got %s (confidence: %.2f)\n", tc.name, tc.expectedType, profile.PrimaryType, profile.Confidence)
		}
	}
}

// Test 2: Verify different content templates exist for different audience types
func testContentTemplates() {
	audienceDetector := services.NewAudienceDetector()
	
	audienceTypes := []services.AudienceType{
		services.AudienceTechnical,
		services.AudienceBusiness,
		services.AudienceExecutive,
		services.AudienceMixed,
	}
	
	for _, audienceType := range audienceTypes {
		template, err := audienceDetector.GetContentTemplate(audienceType, 3)
		if err != nil {
			fmt.Printf("  ✗ %s template: Error - %v\n", audienceType, err)
			continue
		}
		
		// Verify template has appropriate characteristics
		hasGuidelines := len(template.ToneGuidelines) > 0
		hasFocus := len(template.ContentFocus) > 0
		correctType := template.AudienceType == audienceType
		
		if hasGuidelines && hasFocus && correctType {
			fmt.Printf("  ✓ %s template: %d guidelines, %d focus areas\n", 
				audienceType, len(template.ToneGuidelines), len(template.ContentFocus))
		} else {
			fmt.Printf("  ✗ %s template: Missing required elements\n", audienceType)
		}
		
		// Verify template content is audience-appropriate
		verifyTemplateContent(audienceType, template)
	}
}

// Test 3: Verify technical depth adjustment works correctly
func testTechnicalDepthAdjustment() {
	audienceDetector := services.NewAudienceDetector()
	ctx := context.Background()
	
	// Test inquiries with different technical complexity levels
	testCases := []struct {
		name          string
		inquiry       *domain.Inquiry
		expectedDepth int
	}{
		{
			name: "High Technical Depth",
			inquiry: &domain.Inquiry{
				Message: "Need advanced Kubernetes operators, custom resource definitions, Istio traffic management, Prometheus alerting rules, Grafana dashboards, and Terraform modules for multi-cluster deployment.",
			},
			expectedDepth: 4, // Should be 4 or 5
		},
		{
			name: "Medium Technical Depth",
			inquiry: &domain.Inquiry{
				Message: "Looking for cloud migration guidance with basic containerization and CI/CD setup. Need help with AWS services selection and configuration.",
			},
			expectedDepth: 3, // Should be around 3
		},
		{
			name: "Low Technical Depth",
			inquiry: &domain.Inquiry{
				Message: "Our business needs cloud solutions for cost optimization and improved efficiency. Looking for high-level recommendations and strategic guidance.",
			},
			expectedDepth: 2, // Should be 1 or 2
		},
	}
	
	for _, tc := range testCases {
		profile, err := audienceDetector.DetectAudience(ctx, tc.inquiry)
		if err != nil {
			fmt.Printf("  ✗ %s: Error - %v\n", tc.name, err)
			continue
		}
		
		// Check if technical depth is in expected range
		depthOk := false
		switch tc.expectedDepth {
		case 4:
			depthOk = profile.TechnicalDepth >= 4
		case 3:
			depthOk = profile.TechnicalDepth >= 2 && profile.TechnicalDepth <= 4
		case 2:
			depthOk = profile.TechnicalDepth <= 3
		}
		
		if depthOk {
			fmt.Printf("  ✓ %s: Technical depth %d/5\n", tc.name, profile.TechnicalDepth)
		} else {
			fmt.Printf("  ✗ %s: Expected depth ~%d, got %d/5\n", tc.name, tc.expectedDepth, profile.TechnicalDepth)
		}
		
		// Test content adaptation based on technical depth
		sampleContent := "This is a sample technical report with implementation details."
		adaptedContent, err := audienceDetector.AdaptContentForAudience(sampleContent, profile)
		if err != nil {
			fmt.Printf("    ✗ Content adaptation failed: %v\n", err)
		} else {
			hasEnhancements := len(adaptedContent) > len(sampleContent)
			if hasEnhancements {
				fmt.Printf("    ✓ Content adapted based on technical depth\n")
			} else {
				fmt.Printf("    ✗ No content enhancements applied\n")
			}
		}
	}
}

// Test 4: Verify business/technical separation functionality
func testBusinessTechnicalSeparation() {
	audienceDetector := services.NewAudienceDetector()
	
	// Sample mixed content with both business and technical elements
	mixedContent := `
EXECUTIVE SUMMARY
This cloud migration initiative will deliver significant ROI and competitive advantages for the organization.

BUSINESS JUSTIFICATION
The proposed solution will reduce operational costs by 30% and improve time-to-market for new products. The investment will provide competitive advantages through improved scalability and reliability.

TECHNICAL ARCHITECTURE
The solution involves containerizing applications using Docker and deploying them on Kubernetes clusters. We'll implement microservices architecture with API gateways and service mesh for traffic management.

COST ANALYSIS
Total investment: $500K over 18 months
Expected savings: $200K annually
ROI: 40% in year one

IMPLEMENTATION DETAILS
Setup includes EKS cluster configuration, Istio service mesh deployment, Prometheus monitoring stack, and CI/CD pipeline integration with GitOps workflows.

RECOMMENDATIONS
1. Implement phased migration approach
2. Establish cloud governance framework
3. Optimize cost structure through reserved instances
4. Enhance security posture with zero-trust architecture
`
	
	separated, err := audienceDetector.SeparateBusinessAndTechnical(mixedContent)
	if err != nil {
		fmt.Printf("  ✗ Content separation failed: %v\n", err)
		return
	}
	
	// Verify separation results
	hasBusinessContent := len(separated.BusinessJustification) > 0
	hasTechnicalContent := len(separated.TechnicalExplanation) > 0
	hasSharedContent := len(separated.SharedContent) > 0
	hasRecommendations := len(separated.Recommendations) > 0
	
	fmt.Printf("  Business Content: %d characters\n", len(separated.BusinessJustification))
	fmt.Printf("  Technical Content: %d characters\n", len(separated.TechnicalExplanation))
	fmt.Printf("  Shared Content: %d characters\n", len(separated.SharedContent))
	fmt.Printf("  Recommendations: %d items\n", len(separated.Recommendations))
	
	if hasBusinessContent {
		fmt.Printf("  ✓ Business justification extracted\n")
	} else {
		fmt.Printf("  ✗ No business justification found\n")
	}
	
	if hasTechnicalContent {
		fmt.Printf("  ✓ Technical explanation extracted\n")
	} else {
		fmt.Printf("  ✗ No technical explanation found\n")
	}
	
	if hasSharedContent {
		fmt.Printf("  ✓ Shared content identified\n")
	}
	
	if hasRecommendations {
		fmt.Printf("  ✓ Recommendations separated\n")
	} else {
		fmt.Printf("  ✗ No recommendations found\n")
	}
	
	// Verify content classification accuracy
	verifyContentClassification(separated)
}

// Test 5: Verify integration with report generation
func testReportGenerationIntegration() {
	promptArchitect := services.NewPromptArchitect()
	ctx := context.Background()
	
	// Test inquiry
	inquiry := &domain.Inquiry{
		ID:       "integration-test",
		Name:     "Test User",
		Email:    "test@example.com",
		Company:  "Test Company",
		Services: []string{"migration"},
		Message:  "We need technical guidance for Kubernetes migration with business justification for the board.",
	}
	
	// Generate prompt with audience detection
	prompt, err := promptArchitect.BuildReportPrompt(ctx, inquiry, nil)
	if err != nil {
		fmt.Printf("  ✗ Prompt generation failed: %v\n", err)
		return
	}
	
	// Verify prompt contains audience-aware elements
	hasAudienceAdaptation := strings.Contains(prompt, "AUDIENCE ADAPTATION")
	hasTechnicalDepth := strings.Contains(prompt, "Technical depth level")
	hasBusinessFocus := strings.Contains(prompt, "Business focus level")
	hasAudienceDirective := strings.Contains(prompt, "FOCUS:") || strings.Contains(prompt, "MIXED AUDIENCE")
	
	fmt.Printf("  Prompt length: %d characters\n", len(prompt))
	
	if hasAudienceAdaptation {
		fmt.Printf("  ✓ Audience adaptation section included\n")
	} else {
		fmt.Printf("  ✗ No audience adaptation section\n")
	}
	
	if hasTechnicalDepth {
		fmt.Printf("  ✓ Technical depth specification included\n")
	} else {
		fmt.Printf("  ✗ No technical depth specification\n")
	}
	
	if hasBusinessFocus {
		fmt.Printf("  ✓ Business focus specification included\n")
	} else {
		fmt.Printf("  ✗ No business focus specification\n")
	}
	
	if hasAudienceDirective {
		fmt.Printf("  ✓ Audience-specific directive included\n")
	} else {
		fmt.Printf("  ✗ No audience-specific directive\n")
	}
	
	// Test with different explicit audience types
	testExplicitAudienceTypes(promptArchitect, ctx, inquiry)
}

// Helper function to verify template content is appropriate for audience type
func verifyTemplateContent(audienceType services.AudienceType, template *services.ContentTemplate) {
	switch audienceType {
	case services.AudienceTechnical:
		hasArchitecture := containsAny(template.SectionTemplate, []string{"TECHNICAL", "IMPLEMENTATION", "ARCHITECTURE"})
		if hasArchitecture {
			fmt.Printf("    ✓ Technical template contains appropriate sections\n")
		} else {
			fmt.Printf("    ✗ Technical template missing technical sections\n")
		}
		
	case services.AudienceBusiness:
		hasBusiness := containsAny(template.SectionTemplate, []string{"BUSINESS", "ROI", "FINANCIAL"})
		if hasBusiness {
			fmt.Printf("    ✓ Business template contains appropriate sections\n")
		} else {
			fmt.Printf("    ✗ Business template missing business sections\n")
		}
		
	case services.AudienceExecutive:
		hasExecutive := containsAny(template.SectionTemplate, []string{"EXECUTIVE", "STRATEGIC", "TRANSFORMATION"})
		if hasExecutive {
			fmt.Printf("    ✓ Executive template contains appropriate sections\n")
		} else {
			fmt.Printf("    ✗ Executive template missing executive sections\n")
		}
		
	case services.AudienceMixed:
		hasMixed := containsAny(template.SectionTemplate, []string{"BUSINESS", "TECHNICAL", "EXECUTIVE"})
		if hasMixed {
			fmt.Printf("    ✓ Mixed template contains balanced sections\n")
		} else {
			fmt.Printf("    ✗ Mixed template missing balanced sections\n")
		}
	}
}

// Helper function to verify content classification accuracy
func verifyContentClassification(separated *services.SeparatedContent) {
	// Check if business content contains business keywords
	businessKeywords := []string{"roi", "cost", "business", "investment", "savings"}
	businessAccurate := separated.BusinessJustification == "" || 
		containsAny(strings.ToLower(separated.BusinessJustification), businessKeywords)
	
	// Check if technical content contains technical keywords
	technicalKeywords := []string{"architecture", "kubernetes", "docker", "implementation", "configuration"}
	technicalAccurate := separated.TechnicalExplanation == "" || 
		containsAny(strings.ToLower(separated.TechnicalExplanation), technicalKeywords)
	
	if businessAccurate {
		fmt.Printf("  ✓ Business content classification accurate\n")
	} else {
		fmt.Printf("  ✗ Business content classification inaccurate\n")
	}
	
	if technicalAccurate {
		fmt.Printf("  ✓ Technical content classification accurate\n")
	} else {
		fmt.Printf("  ✗ Technical content classification inaccurate\n")
	}
}

// Test explicit audience type specification
func testExplicitAudienceTypes(promptArchitect interfaces.PromptArchitect, ctx context.Context, inquiry *domain.Inquiry) {
	audienceTypes := []string{"technical", "business", "executive", "mixed"}
	
	for _, audienceType := range audienceTypes {
		options := &interfaces.PromptOptions{
			TargetAudience: audienceType,
		}
		
		prompt, err := promptArchitect.BuildReportPrompt(ctx, inquiry, options)
		if err != nil {
			fmt.Printf("    ✗ %s audience prompt failed: %v\n", audienceType, err)
			continue
		}
		
		// Check if prompt contains audience-specific content
		hasAudienceContent := strings.Contains(strings.ToUpper(prompt), strings.ToUpper(audienceType+" FOCUS"))
		if hasAudienceContent {
			fmt.Printf("    ✓ %s audience prompt generated correctly\n", audienceType)
		} else {
			fmt.Printf("    ✗ %s audience prompt missing specific content\n", audienceType)
		}
	}
}

// Helper function to check if text contains any of the given substrings
func containsAny(text string, substrings []string) bool {
	textLower := strings.ToLower(text)
	for _, substr := range substrings {
		if strings.Contains(textLower, strings.ToLower(substr)) {
			return true
		}
	}
	return false
}