package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	ctx := context.Background()

	// Initialize knowledge base service
	knowledgeBase := services.NewKnowledgeBaseService()

	fmt.Println("=== Testing Company Knowledge Base ===\n")

	// Test 1: Get service offerings
	fmt.Println("1. Testing service offerings...")
	offerings, err := knowledgeBase.GetServiceOfferings(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d service offerings:\n", len(offerings))
		for _, offering := range offerings {
			fmt.Printf("   - %s: %s\n", offering.Name, offering.Description)
			fmt.Printf("     Duration: %s, Team Size: %s\n", offering.TypicalDuration, offering.TeamSize)
		}
	}

	// Test 2: Get team expertise
	fmt.Println("\n2. Testing team expertise...")
	expertise, err := knowledgeBase.GetTeamExpertise(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d team members:\n", len(expertise))
		for _, expert := range expertise {
			fmt.Printf("   - %s (%s): %d years experience\n",
				expert.ConsultantName, expert.Role, expert.ExperienceYears)
			fmt.Printf("     Expertise: %v\n", expert.ExpertiseAreas)
		}
	}

	// Test 3: Get past solutions
	fmt.Println("\n3. Testing past solutions...")
	solutions, err := knowledgeBase.GetPastSolutions(ctx, "migration", "Technology")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d past solutions:\n", len(solutions))
		for _, solution := range solutions {
			fmt.Printf("   - %s (%s)\n", solution.Title, solution.Industry)
			fmt.Printf("     Cost Savings: $%.0f, Time to Value: %s\n",
				solution.CostSavings, solution.TimeToValue)
		}
	}

	// Test 4: Get methodology templates
	fmt.Println("\n4. Testing methodology templates...")
	templates, err := knowledgeBase.GetMethodologyTemplates(ctx, "migration")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d methodology templates:\n", len(templates))
		for _, template := range templates {
			fmt.Printf("   - %s: %s\n", template.Name, template.Description)
			fmt.Printf("     Phases: %d\n", len(template.Phases))
		}
	}

	// Test 5: Search knowledge
	fmt.Println("\n5. Testing knowledge search...")
	items, err := knowledgeBase.SearchKnowledge(ctx, "migration", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d knowledge items for 'migration':\n", len(items))
		for _, item := range items {
			fmt.Printf("   - %s (relevance: %.2f)\n", item.Title, item.Relevance)
		}
	}

	// Test 6: Get knowledge statistics
	fmt.Println("\n6. Testing knowledge statistics...")
	stats, err := knowledgeBase.GetKnowledgeStats(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Knowledge Base Statistics:\n")
		fmt.Printf("   - Services: %d\n", stats.TotalServices)
		fmt.Printf("   - Expertise: %d\n", stats.TotalExpertise)
		fmt.Printf("   - Engagements: %d\n", stats.TotalEngagements)
		fmt.Printf("   - Solutions: %d\n", stats.TotalSolutions)
		fmt.Printf("   - Methodologies: %d\n", stats.TotalMethodologies)
		fmt.Printf("   - Last Updated: %s\n", stats.LastUpdated.Format(time.RFC3339))
	}

	// Test 7: Test client history integration
	fmt.Println("\n7. Testing client history integration...")
	clientHistory := services.NewClientHistoryService(knowledgeBase)

	// Create test inquiry
	inquiry := &domain.Inquiry{
		ID:        "test-001",
		Name:      "John Smith",
		Email:     "john@techcorp.com",
		Company:   "TechCorp Solutions",
		Services:  []string{"migration"},
		Message:   "Need help with cloud migration",
		CreatedAt: time.Now(),
	}

	// Get client insights
	insights, err := clientHistory.GetClientInsights(ctx, inquiry.Company)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Client Insights for %s:\n", inquiry.Company)
		fmt.Printf("   - Total Engagements: %d\n", insights.TotalEngagements)
		fmt.Printf("   - Average Satisfaction: %.1f\n", insights.AverageSatisfaction)
		fmt.Printf("   - Total Value: $%.0f\n", insights.TotalValue)
		fmt.Printf("   - Recommended Approach: %s\n", insights.RecommendedApproach)
	}

	fmt.Println("\n=== Knowledge Base Test Complete ===")
}
