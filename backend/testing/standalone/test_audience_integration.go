package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Testing Audience-Aware Report Generation Integration ===")
	
	// Create services
	promptArchitect := services.NewPromptArchitect()
	audienceDetector := services.NewAudienceDetector()
	
	// Test different inquiry types
	testInquiries := []*domain.Inquiry{
		{
			ID:       "tech-detailed",
			Name:     "Alex Rodriguez",
			Email:    "alex.rodriguez@devops.com",
			Company:  "DevOps Solutions Inc",
			Services: []string{"migration", "architecture"},
			Message:  "We need detailed technical guidance for migrating our Kubernetes clusters to EKS. Our current setup includes Istio service mesh, Prometheus monitoring, and GitOps with ArgoCD. We need specific configuration examples for VPC networking, IAM roles, and auto-scaling policies. Performance benchmarks and security best practices are essential.",
		},
		{
			ID:       "business-roi",
			Name:     "Jennifer Williams",
			Email:    "jennifer.williams@retailchain.com",
			Company:  "National Retail Chain",
			Services: []string{"optimization", "assessment"},
			Message:  "Our board is requesting a comprehensive business case for cloud optimization. We need clear ROI projections, cost-benefit analysis, and competitive advantages. The proposal should demonstrate how cloud adoption will improve operational efficiency, reduce costs, and support our expansion strategy. Budget approval depends on showing measurable business value.",
		},
		{
			ID:       "executive-strategy",
			Name:     "Robert Thompson",
			Email:    "robert.thompson@globalcorp.com",
			Company:  "Global Manufacturing Corp",
			Services: []string{"strategy", "transformation"},
			Message:  "As Chief Digital Officer, I need a strategic assessment of our digital transformation roadmap. The board wants to understand our competitive positioning, investment requirements, and long-term vision for cloud adoption. This initiative will impact our entire organization and requires stakeholder alignment across all divisions.",
		},
	}
	
	ctx := context.Background()
	
	for i, inquiry := range testInquiries {
		fmt.Printf("\n--- Test Case %d: %s ---\n", i+1, inquiry.Company)
		
		// Test audience detection
		profile, err := audienceDetector.DetectAudience(ctx, inquiry)
		if err != nil {
			log.Printf("Error detecting audience: %v", err)
			continue
		}
		
		fmt.Printf("Detected Audience: %s (confidence: %.2f)\n", profile.PrimaryType, profile.Confidence)
		fmt.Printf("Technical Depth: %d/5, Business Focus: %d/5\n", profile.TechnicalDepth, profile.BusinessFocus)
		
		// Test prompt generation with audience awareness
		prompt, err := promptArchitect.BuildReportPrompt(ctx, inquiry, nil)
		if err != nil {
			log.Printf("Error building prompt: %v", err)
			continue
		}
		
		fmt.Printf("Generated Prompt Length: %d characters\n", len(prompt))
		
		// Analyze prompt content for audience-specific elements
		audienceElements := analyzePromptContent(prompt, string(profile.PrimaryType))
		fmt.Printf("Audience-Specific Elements Found: %d\n", len(audienceElements))
		for _, element := range audienceElements {
			fmt.Printf("  - %s\n", element)
		}
		
		// Test content separation for mixed audiences
		if profile.PrimaryType == services.AudienceMixed {
			sampleReport := generateSampleReport(inquiry)
			separated, err := audienceDetector.SeparateBusinessAndTechnical(sampleReport)
			if err != nil {
				log.Printf("Error separating content: %v", err)
				continue
			}
			
			fmt.Printf("Content Separation Results:\n")
			fmt.Printf("  Business Content: %d characters\n", len(separated.BusinessJustification))
			fmt.Printf("  Technical Content: %d characters\n", len(separated.TechnicalExplanation))
			fmt.Printf("  Shared Content: %d characters\n", len(separated.SharedContent))
			fmt.Printf("  Recommendations: %d items\n", len(separated.Recommendations))
		}
		
		// Test template customization
		template, err := audienceDetector.GetContentTemplate(profile.PrimaryType, profile.TechnicalDepth)
		if err != nil {
			log.Printf("Error getting template: %v", err)
			continue
		}
		
		fmt.Printf("Template Guidelines: %d tone guidelines, %d content focus areas\n", 
			len(template.ToneGuidelines), len(template.ContentFocus))
		
		// Verify audience-specific adaptations
		verifyAudienceAdaptations(profile, prompt)
	}
	
	// Test edge cases
	fmt.Printf("\n--- Testing Edge Cases ---\n")
	
	// Test with minimal content
	minimalInquiry := &domain.Inquiry{
		ID:       "minimal",
		Name:     "Test User",
		Email:    "test@example.com",
		Company:  "Test Corp",
		Services: []string{"general"},
		Message:  "Help needed.",
	}
	
	profile, err := audienceDetector.DetectAudience(ctx, minimalInquiry)
	if err != nil {
		log.Printf("Error with minimal inquiry: %v", err)
	} else {
		fmt.Printf("Minimal Inquiry - Detected: %s (confidence: %.2f)\n", profile.PrimaryType, profile.Confidence)
	}
	
	// Test with highly technical content
	technicalInquiry := &domain.Inquiry{
		ID:       "technical",
		Name:     "Senior Architect",
		Email:    "architect@tech.com",
		Company:  "Tech Solutions",
		Services: []string{"architecture"},
		Message:  "Need Terraform modules for EKS cluster with Istio service mesh, Prometheus monitoring, Grafana dashboards, ArgoCD GitOps, Vault secrets management, and Cilium CNI. Require YAML configurations for ingress controllers, network policies, and RBAC. Performance tuning for high-throughput microservices architecture.",
	}
	
	profile, err = audienceDetector.DetectAudience(ctx, technicalInquiry)
	if err != nil {
		log.Printf("Error with technical inquiry: %v", err)
	} else {
		fmt.Printf("Technical Inquiry - Detected: %s (depth: %d/5, confidence: %.2f)\n", 
			profile.PrimaryType, profile.TechnicalDepth, profile.Confidence)
	}
	
	// Test with business-heavy content
	businessInquiry := &domain.Inquiry{
		ID:       "business",
		Name:     "CFO",
		Email:    "cfo@business.com",
		Company:  "Business Corp",
		Services: []string{"optimization"},
		Message:  "Board requires comprehensive ROI analysis, cost-benefit projections, competitive advantage assessment, market positioning strategy, stakeholder value proposition, investment requirements, budget implications, and business case justification for cloud transformation initiative.",
	}
	
	profile, err = audienceDetector.DetectAudience(ctx, businessInquiry)
	if err != nil {
		log.Printf("Error with business inquiry: %v", err)
	} else {
		fmt.Printf("Business Inquiry - Detected: %s (focus: %d/5, confidence: %.2f)\n", 
			profile.PrimaryType, profile.BusinessFocus, profile.Confidence)
	}
	
	fmt.Println("\n=== Integration Test Complete ===")
}

// analyzePromptContent analyzes prompt for audience-specific elements
func analyzePromptContent(prompt, audienceType string) []string {
	elements := []string{}
	
	// Common audience indicators
	if contains(prompt, "AUDIENCE ADAPTATION") {
		elements = append(elements, "Audience adaptation section found")
	}
	
	if contains(prompt, "Technical depth level") {
		elements = append(elements, "Technical depth specification")
	}
	
	if contains(prompt, "Business focus level") {
		elements = append(elements, "Business focus specification")
	}
	
	// Audience-specific indicators
	switch audienceType {
	case "technical":
		if contains(prompt, "TECHNICAL FOCUS") {
			elements = append(elements, "Technical focus directive")
		}
		if contains(prompt, "implementation details") {
			elements = append(elements, "Implementation details emphasis")
		}
		if contains(prompt, "architecture") {
			elements = append(elements, "Architecture focus")
		}
		
	case "business":
		if contains(prompt, "BUSINESS FOCUS") {
			elements = append(elements, "Business focus directive")
		}
		if contains(prompt, "ROI") {
			elements = append(elements, "ROI emphasis")
		}
		if contains(prompt, "business value") {
			elements = append(elements, "Business value focus")
		}
		
	case "executive":
		if contains(prompt, "EXECUTIVE FOCUS") {
			elements = append(elements, "Executive focus directive")
		}
		if contains(prompt, "strategic") {
			elements = append(elements, "Strategic emphasis")
		}
		if contains(prompt, "transformation") {
			elements = append(elements, "Transformation focus")
		}
		
	case "mixed":
		if contains(prompt, "MIXED AUDIENCE") {
			elements = append(elements, "Mixed audience directive")
		}
		if contains(prompt, "balance") {
			elements = append(elements, "Balance directive")
		}
	}
	
	return elements
}

// generateSampleReport creates a sample report for testing content separation
func generateSampleReport(inquiry *domain.Inquiry) string {
	return fmt.Sprintf(`
EXECUTIVE SUMMARY
This report provides strategic recommendations for %s's cloud initiative.

BUSINESS JUSTIFICATION
The proposed cloud migration will deliver significant ROI through cost savings and operational efficiency improvements. The investment will provide competitive advantages and support business growth objectives.

TECHNICAL ANALYSIS
The current architecture requires modernization to support cloud-native deployment patterns. Key technical considerations include containerization, microservices architecture, and automated CI/CD pipelines.

RECOMMENDATIONS
1. Implement phased migration approach
2. Establish cloud governance framework
3. Optimize cost structure through reserved instances
4. Enhance security posture with zero-trust architecture

RISK ASSESSMENT
Technical risks include potential downtime during migration and integration complexity. Business risks involve budget overruns and timeline delays.

IMPLEMENTATION ROADMAP
Phase 1: Assessment and planning (4 weeks)
Phase 2: Pilot migration (8 weeks)
Phase 3: Full migration (16 weeks)
Phase 4: Optimization (ongoing)
`, inquiry.Company)
}

// verifyAudienceAdaptations checks if prompt contains appropriate audience adaptations
func verifyAudienceAdaptations(profile *services.AudienceProfile, prompt string) {
	fmt.Printf("Audience Adaptation Verification:\n")
	
	// Check for audience type mention
	if contains(prompt, string(profile.PrimaryType)) {
		fmt.Printf("  ✓ Audience type mentioned in prompt\n")
	} else {
		fmt.Printf("  ✗ Audience type not mentioned in prompt\n")
	}
	
	// Check for technical depth adaptation
	if profile.TechnicalDepth >= 4 && contains(prompt, "technical") {
		fmt.Printf("  ✓ High technical depth reflected\n")
	} else if profile.TechnicalDepth <= 2 && !contains(prompt, "technical details") {
		fmt.Printf("  ✓ Low technical depth reflected\n")
	}
	
	// Check for business focus adaptation
	if profile.BusinessFocus >= 4 && contains(prompt, "business") {
		fmt.Printf("  ✓ High business focus reflected\n")
	} else if profile.BusinessFocus <= 2 && !contains(prompt, "business case") {
		fmt.Printf("  ✓ Low business focus reflected\n")
	}
	
	// Check for executive level adaptation
	if profile.ExecutiveLevel && contains(prompt, "strategic") {
		fmt.Printf("  ✓ Executive level adaptation present\n")
	}
}

// Helper function to check if string contains substring (case-insensitive)
func contains(text, substr string) bool {
	return len(text) > 0 && len(substr) > 0 && findInString(strings.ToLower(text), strings.ToLower(substr))
}

func findInString(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

