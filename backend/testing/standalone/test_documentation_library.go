package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	// Test the documentation library service
	service := services.NewDocumentationLibraryService()
	
	fmt.Println("=== Testing Documentation Library Service ===")
	
	// Test 1: Check if service is healthy
	fmt.Printf("Service is healthy: %v\n", service.IsHealthy())
	
	// Test 2: Get stats
	stats := service.GetStats()
	fmt.Printf("Total links: %d\n", stats.TotalLinks)
	fmt.Printf("Valid links: %d\n", stats.ValidLinks)
	fmt.Printf("Providers: %v\n", stats.LinksByProvider)
	
	// Test 3: Get AWS links
	ctx := context.Background()
	awsLinks, err := service.GetDocumentationLinks(ctx, "aws", "")
	if err != nil {
		log.Fatalf("Error getting AWS links: %v", err)
	}
	fmt.Printf("AWS links found: %d\n", len(awsLinks))
	
	// Test 4: Search for security documentation
	securityLinks, err := service.SearchDocumentation(ctx, "security", nil)
	if err != nil {
		log.Fatalf("Error searching security links: %v", err)
	}
	fmt.Printf("Security links found: %d\n", len(securityLinks))
	
	// Test 5: Get links by category
	securityCategoryLinks, err := service.GetLinksByCategory(ctx, "security")
	if err != nil {
		log.Fatalf("Error getting security category links: %v", err)
	}
	fmt.Printf("Security category links found: %d\n", len(securityCategoryLinks))
	
	fmt.Println("=== All tests completed successfully! ===")
}