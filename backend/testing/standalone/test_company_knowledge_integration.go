package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	ctx := context.Background()

	// Initialize services
	knowledgeBase := services.NewKnowledgeBaseService()
	clientHistory := services.NewClientHistoryService(knowledgeBase)
	companyKnowledgeInteg := services.NewCompanyKnowledgeIntegrationService(knowledgeBase, clientHistory)

	// Test inquiry
	inquiry := &domain.Inquiry{
		ID:        "test-001",
		Name:      "John Smith",
		Email:     "john@techcorp.com",
		Company:   "TechCorp Solutions",
		Services:  []string{"migration", "architecture"},
		Message:   "We need help migrating our legacy applications to AWS cloud",
		CreatedAt: time.Now(),
	}

	fmt.Println("=== Testing Company Knowledge Integration ===\n")

	// Test 1: Get company context for inquiry
	fmt.Println("1. Getting company context for inquiry...")
	context, err := companyKnowledgeInteg.GetCompanyContextForInquiry(ctx, inquiry)
	if err != nil {
		log.Printf("Error getting company context: %v", err)
	} else {
		fmt.Printf("   - Found %d service offerings\n", len(context.ServiceOfferings))
		fmt.Printf("   - Found %d team experts\n", len(context.TeamExpertise))
		fmt.Printf("   - Found %d past solutions\n", len(context.PastSolutions))
		fmt.Printf("   - Found %d client history records\n", len(context.ClientHistory))
		fmt.Printf("   - Found %d project patterns\n", len(context.ProjectPatterns))
	}

	// Test 2: Generate contextual prompt
	fmt.Println("\n2. Generating contextual prompt...")
	basePrompt := "Generate a professional response for this client inquiry."
	enhancedPrompt, err := companyKnowledgeInteg.GenerateContextualPrompt(ctx, inquiry, basePrompt)
	if err != nil {
		log.Printf("Error generating contextual prompt: %v", err)
	} else {
		fmt.Printf("   - Enhanced prompt length: %d characters\n", len(enhancedPrompt))
		fmt.Printf("   - Contains service offerings: %t\n", containsString(enhancedPrompt, "Service Offerings"))
		fmt.Printf("   - Contains team expertise: %t\n", containsString(enhancedPrompt, "Team Expertise"))
		fmt.Printf("   - Contains past solutions: %t\n", containsString(enhancedPrompt, "Past Solutions"))
	}

	// Test 3: Get recommendations for inquiry
	fmt.Println("\n3. Getting recommendations for inquiry...")
	recommendations, err := companyKnowledgeInteg.GetRecommendationsForInquiry(ctx, inquiry)
	if err != nil {
		log.Printf("Error getting recommendations: %v", err)
	} else {
		fmt.Printf("   - Recommended services: %d\n", len(recommendations.RecommendedServices))
		fmt.Printf("   - Recommended team members: %d\n", len(recommendations.RecommendedTeam))
		fmt.Printf("   - Similar patterns: %d\n", len(recommendations.SimilarPatterns))
		fmt.Printf("   - Methodology templates: %d\n", len(recommendations.MethodologyTemplates))
	}

	// Test 4: Test client history integration
	fmt.Println("\n4. Testing client history integration...")

	// Record engagement
	err = clientHistory.RecordEngagement(ctx, inquiry)
	if err != nil {
		log.Printf("Error recording engagement: %v", err)
	} else {
		fmt.Println("   - Successfully recorded engagement")
	}

	// Get client insights
	insights, err := clientHistory.GetClientInsights(ctx, inquiry.Company)
	if err != nil {
		log.Printf("Error getting client insights: %v", err)
	} else {
		fmt.Printf("   - Total engagements: %d\n", insights.TotalEngagements)
		fmt.Printf("   - Average satisfaction: %.1f\n", insights.AverageSatisfaction)
		fmt.Printf("   - Total value: $%.0f\n", insights.TotalValue)
		fmt.Printf("   - Recommended approach: %s\n", insights.RecommendedApproach)
	}

	// Test 5: Get recommended services
	fmt.Println("\n5. Testing service recommendations...")
	recommendedServices, err := clientHistory.GetRecommendedServices(ctx, inquiry.Company, "Technology")
	if err != nil {
		log.Printf("Error getting recommended services: %v", err)
	} else {
		fmt.Printf("   - Found %d recommended services\n", len(recommendedServices))
		for _, service := range recommendedServices {
			fmt.Printf("     • %s (%s)\n", service.Name, service.Category)
		}
	}

	// Test 6: Get recommended team
	fmt.Println("\n6. Testing team recommendations...")
	recommendedTeam, err := clientHistory.GetRecommendedTeam(ctx, inquiry.Company, "migration")
	if err != nil {
		log.Printf("Error getting recommended team: %v", err)
	} else {
		fmt.Printf("   - Found %d recommended team members\n", len(recommendedTeam))
		for _, member := range recommendedTeam {
			fmt.Printf("     • %s (%s) - %s\n", member.ConsultantName, member.Role, member.ExpertiseAreas[0])
		}
	}

	// Test 7: Knowledge base statistics
	fmt.Println("\n7. Knowledge base statistics...")
	stats, err := knowledgeBase.GetKnowledgeStats(ctx)
	if err != nil {
		log.Printf("Error getting knowledge stats: %v", err)
	} else {
		fmt.Printf("   - Total services: %d\n", stats.TotalServices)
		fmt.Printf("   - Total expertise records: %d\n", stats.TotalExpertise)
		fmt.Printf("   - Total engagements: %d\n", stats.TotalEngagements)
		fmt.Printf("   - Total solutions: %d\n", stats.TotalSolutions)
		fmt.Printf("   - Total methodologies: %d\n", stats.TotalMethodologies)
	}

	fmt.Println("\n=== Company Knowledge Integration Test Complete ===")
}

func containsString(text, substring string) bool {
	return len(text) > 0 && len(substring) > 0 &&
		text != substring &&
		len(text) >= len(substring) &&
		findSubstring(text, substring)
}

func findSubstring(text, substring string) bool {
	if len(substring) > len(text) {
		return false
	}

	for i := 0; i <= len(text)-len(substring); i++ {
		match := true
		for j := 0; j < len(substring); j++ {
			if text[i+j] != substring[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
