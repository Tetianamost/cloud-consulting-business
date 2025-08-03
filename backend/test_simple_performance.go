package main

import (
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("=== Simple Performance Test ===")

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Test intelligent cache
	fmt.Println("Testing intelligent cache...")
	cache := services.NewIntelligentAnalysisCache(logger)

	// Test cache miss
	result := cache.GetCachedAnalysis("test", "test content")
	if result != nil {
		log.Printf("Unexpected cache hit")
		return
	}
	fmt.Println("✓ Cache miss handled correctly")

	// Test caching
	err := cache.CacheAnalysis("test", "test content", "test result", 100, 0.8)
	if err != nil {
		log.Printf("Error caching: %v", err)
		return
	}
	fmt.Println("✓ Analysis cached successfully")

	// Test cache hit
	result = cache.GetCachedAnalysis("test", "test content")
	if result == nil {
		log.Printf("Expected cache hit")
		return
	}
	fmt.Printf("✓ Cache hit successful - Content: %s\n", result.Content)

	fmt.Println("=== Simple Performance Test Completed ===")
}
