package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Testing Audience Detection and Content Generation ===")
	
	// Create audience detector
	audienceDetector := services.NewAudienceDetector()
	
	// Test scenarios with different audience types
	testScenarios := []struct {
		name    string
		inquiry *domain.Inquiry
	}{
		{
			name: "Technical Inquiry - Infrastructure Migration",
			inquiry: &domain.Inquiry{
				ID:       "tech-001",
				Name:     "John Smith",
				Email:    "john.smith@techcorp.com",
				Company:  "TechCorp Solutions",
				Services: []string{"migration", "architecture"},
				Message:  "We need to migrate our microservices architecture from on-premises to AWS. Our current setup includes Docker containers, Kubernetes clusters, and a complex CI/CD pipeline. We're looking for guidance on VPC design, load balancer configuration, and auto-scaling strategies. Performance and scalability are critical requirements.",
			},
		},
		{
			name: "Business Inquiry - Cost Optimization",
			inquiry: &domain.Inquiry{
				ID:       "biz-001",
				Name:     "Sarah Johnson",
				Email:    "sarah.johnson@retailcorp.com",
				Company:  "RetailCorp Inc",
				Services: []string{"optimization", "assessment"},
				Message:  "Our company is looking to reduce cloud costs and improve operational efficiency. We need a business case for cloud optimization that shows clear ROI and cost savings. The board wants to understand the competitive advantages and how this aligns with our growth strategy.",
			},
		},
		{
			name: "Executive Inquiry - Digital Transformation",
			inquiry: &domain.Inquiry{
				ID:       "exec-001",
				Name:     "Michael Chen",
				Email:    "michael.chen@enterprise.com",
				Company:  "Enterprise Holdings",
				Services: []string{"assessment", "strategy"},
				Message:  "As CEO, I'm looking for a high-level strategic assessment of our digital transformation initiative. We need to understand the investment requirements, competitive positioning, and long-term vision for cloud adoption across our organization. This needs board approval and stakeholder buy-in.",
			},
		},
		{
			name: "Mixed Inquiry - Healthcare Compliance",
			inquiry: &domain.Inquiry{
				ID:       "mixed-001",
				Name:     "Dr. Lisa Wang",
				Email:    "lisa.wang@healthsystem.org",
				Company:  "Regional Health System",
				Services: []string{"migration", "compliance"},
				Message:  "We're a healthcare organization that needs to migrate patient data to the cloud while maintaining HIPAA compliance. We need both technical implementation details for our IT team and business justification for our executive leadership. The solution must address security, performance, and regulatory requirements.",
			},
		},
	}
	
	ctx := context.Background()
	
	for _, scenario := range testScenarios {
		fmt.Printf("\n--- %s ---\n", scenario.name)
		
		// Test audience detection
		profile, err := audienceDetector.DetectAudience(ctx, scenario.inquiry)
		if err != nil {
			log.Printf("Error detecting audience: %v", err)
			continue
		}
		
		// Print audience profile
		fmt.Printf("Detected Audience Profile:\n")
		fmt.Printf("  Primary Type: %s\n", profile.PrimaryType)
		if profile.SecondaryType != "" {
			fmt.Printf("  Secondary Type: %s\n", profile.SecondaryType)
		}
		fmt.Printf("  Technical Depth: %d/5\n", profile.TechnicalDepth)
		fmt.Printf("  Business Focus: %d/5\n", profile.BusinessFocus)
		fmt.Printf("  Executive Level: %t\n", profile.ExecutiveLevel)
		fmt.Printf("  Confidence: %.2f\n", profile.Confidence)
		fmt.Printf("  Indicators: %v\n", profile.Indicators)
		
		// Test content template retrieval
		template, err := audienceDetector.GetContentTemplate(profile.PrimaryType, profile.TechnicalDepth)
		if err != nil {
			log.Printf("Error getting content template: %v", err)
			continue
		}
		
		fmt.Printf("\nContent Template:\n")
		fmt.Printf("  Audience Type: %s\n", template.AudienceType)
		fmt.Printf("  Technical Depth: %d\n", template.TechnicalDepth)
		fmt.Printf("  Business Focus: %d\n", template.BusinessFocus)
		fmt.Printf("  Tone Guidelines: %v\n", template.ToneGuidelines)
		fmt.Printf("  Content Focus: %v\n", template.ContentFocus)
		
		// Test content adaptation
		sampleContent := `
EXECUTIVE SUMMARY
This report provides recommendations for cloud migration and optimization.

TECHNICAL ANALYSIS
The current architecture requires significant updates to support cloud deployment.

BUSINESS JUSTIFICATION
The proposed solution will deliver cost savings and improved efficiency.

RECOMMENDATIONS
1. Implement cloud-native architecture
2. Optimize cost structure
3. Enhance security posture
`
		
		adaptedContent, err := audienceDetector.AdaptContentForAudience(sampleContent, profile)
		if err != nil {
			log.Printf("Error adapting content: %v", err)
			continue
		}
		
		fmt.Printf("\nAdapted Content Preview:\n")
		fmt.Printf("%s\n", adaptedContent[:min(len(adaptedContent), 500)] + "...")
		
		// Test business/technical separation
		separated, err := audienceDetector.SeparateBusinessAndTechnical(sampleContent)
		if err != nil {
			log.Printf("Error separating content: %v", err)
			continue
		}
		
		fmt.Printf("\nContent Separation:\n")
		if separated.BusinessJustification != "" {
			fmt.Printf("  Business Section: %s\n", separated.BusinessJustification[:min(len(separated.BusinessJustification), 100)] + "...")
		}
		if separated.TechnicalExplanation != "" {
			fmt.Printf("  Technical Section: %s\n", separated.TechnicalExplanation[:min(len(separated.TechnicalExplanation), 100)] + "...")
		}
		if separated.SharedContent != "" {
			fmt.Printf("  Shared Content: %s\n", separated.SharedContent[:min(len(separated.SharedContent), 100)] + "...")
		}
		fmt.Printf("  Recommendations: %d items\n", len(separated.Recommendations))
	}
	
	// Test prompt architect integration
	fmt.Printf("\n--- Testing Prompt Architect Integration ---\n")
	
	promptArchitect := services.NewPromptArchitect()
	
	// Test with technical inquiry
	techInquiry := testScenarios[0].inquiry
	prompt, err := promptArchitect.BuildReportPrompt(ctx, techInquiry, nil)
	if err != nil {
		log.Printf("Error building prompt: %v", err)
	} else {
		fmt.Printf("Generated Prompt Length: %d characters\n", len(prompt))
		fmt.Printf("Prompt Preview: %s...\n", prompt[:min(len(prompt), 300)])
		
		// Check if audience-specific content is included
		if containsAudienceContent(prompt, "technical") {
			fmt.Println("✓ Technical audience content detected in prompt")
		} else {
			fmt.Println("✗ Technical audience content not detected in prompt")
		}
	}
	
	// Test with business inquiry
	bizInquiry := testScenarios[1].inquiry
	prompt, err = promptArchitect.BuildReportPrompt(ctx, bizInquiry, nil)
	if err != nil {
		log.Printf("Error building prompt: %v", err)
	} else {
		fmt.Printf("Generated Prompt Length: %d characters\n", len(prompt))
		
		// Check if audience-specific content is included
		if containsAudienceContent(prompt, "business") {
			fmt.Println("✓ Business audience content detected in prompt")
		} else {
			fmt.Println("✗ Business audience content not detected in prompt")
		}
	}
	
	fmt.Println("\n=== Audience Detection Test Complete ===")
}

// Helper function to check if prompt contains audience-specific content
func containsAudienceContent(prompt, audienceType string) bool {
	switch audienceType {
	case "technical":
		return contains(prompt, "TECHNICAL") || contains(prompt, "architecture") || contains(prompt, "implementation")
	case "business":
		return contains(prompt, "BUSINESS") || contains(prompt, "ROI") || contains(prompt, "cost")
	case "executive":
		return contains(prompt, "EXECUTIVE") || contains(prompt, "strategic") || contains(prompt, "transformation")
	}
	return false
}

func contains(text, substr string) bool {
	return len(text) > 0 && len(substr) > 0 && 
		   (text == substr || 
		    (len(text) >= len(substr) && 
		     (text[:len(substr)] == substr || 
		      text[len(text)-len(substr):] == substr || 
		      findInString(text, substr))))
}

func findInString(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Test JSON serialization of audience profile
func testJSONSerialization() {
	profile := &services.AudienceProfile{
		PrimaryType:    services.AudienceTechnical,
		TechnicalDepth: 4,
		BusinessFocus:  2,
		ExecutiveLevel: false,
		Confidence:     0.85,
		Indicators:     []string{"architecture", "microservices", "kubernetes"},
	}
	
	jsonData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	
	fmt.Printf("JSON Serialization Test:\n%s\n", string(jsonData))
}