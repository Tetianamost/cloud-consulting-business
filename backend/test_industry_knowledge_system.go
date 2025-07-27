package main

import (
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Testing Industry-Specific Knowledge System ===")
	
	// Initialize the knowledge base
	kb := services.NewInMemoryKnowledgeBase()
	
	// Test 1: Healthcare Industry Compliance Requirements
	fmt.Println("\n1. Testing Healthcare Compliance Requirements:")
	healthcareCompliance, err := kb.GetComplianceRequirements("healthcare")
	if err != nil {
		log.Printf("Error getting healthcare compliance: %v", err)
	} else {
		fmt.Printf("Found %d healthcare compliance requirements:\n", len(healthcareCompliance))
		for _, req := range healthcareCompliance {
			fmt.Printf("  - %s (%s): %s\n", req.Framework, req.Severity, req.Requirement)
		}
	}
	
	// Test 2: Financial Industry Best Practices
	fmt.Println("\n2. Testing Financial Industry Best Practices:")
	financialBestPractices, err := kb.GetIndustrySpecificBestPractices("financial")
	if err != nil {
		log.Printf("Error getting financial best practices: %v", err)
	} else {
		fmt.Printf("Found %d financial industry best practices:\n", len(financialBestPractices))
		for _, bp := range financialBestPractices {
			fmt.Printf("  - %s (%s): %s\n", bp.Title, bp.Priority, bp.Description)
		}
	}
	
	// Test 3: Manufacturing Architectural Patterns
	fmt.Println("\n3. Testing Manufacturing Architectural Patterns:")
	manufacturingPatterns, err := kb.GetIndustrySpecificArchitecturalPatterns("manufacturing")
	if err != nil {
		log.Printf("Error getting manufacturing patterns: %v", err)
	} else {
		fmt.Printf("Found %d manufacturing architectural patterns:\n", len(manufacturingPatterns))
		for _, pattern := range manufacturingPatterns {
			fmt.Printf("  - %s (%s): %s\n", pattern.Name, pattern.Complexity, pattern.Description)
		}
	}
	
	// Test 4: Retail Industry Risk Factors
	fmt.Println("\n4. Testing Retail Industry Risk Factors:")
	retailRisks, err := kb.GetIndustryRiskFactors("retail")
	if err != nil {
		log.Printf("Error getting retail risks: %v", err)
	} else {
		fmt.Printf("Found %d retail industry risk factors:\n", len(retailRisks))
		for i, risk := range retailRisks {
			if i < 5 { // Show first 5 risks
				fmt.Printf("  - %s\n", risk)
			}
		}
		if len(retailRisks) > 5 {
			fmt.Printf("  ... and %d more\n", len(retailRisks)-5)
		}
	}
	
	// Test 5: Healthcare Specific Recommendations
	fmt.Println("\n5. Testing Healthcare Data Migration Recommendations:")
	healthcareRecs, err := kb.GetIndustrySpecificRecommendations("healthcare", "data-migration")
	if err != nil {
		log.Printf("Error getting healthcare recommendations: %v", err)
	} else {
		fmt.Printf("Found %d healthcare data migration recommendations:\n", len(healthcareRecs))
		for _, rec := range healthcareRecs {
			fmt.Printf("  - %s\n", rec)
		}
	}
	
	// Test 6: Government Compliance Frameworks
	fmt.Println("\n6. Testing Government Compliance Frameworks:")
	govFrameworks, err := kb.GetComplianceFrameworksByIndustry("government")
	if err != nil {
		log.Printf("Error getting government frameworks: %v", err)
	} else {
		fmt.Printf("Found %d compliance frameworks for government:\n", len(govFrameworks))
		for _, framework := range govFrameworks {
			fmt.Printf("  - %s\n", framework)
		}
	}
	
	// Test 7: Education Industry Risk Factors
	fmt.Println("\n7. Testing Education Industry Risk Factors:")
	educationRisks, err := kb.GetIndustryRiskFactors("education")
	if err != nil {
		log.Printf("Error getting education risks: %v", err)
	} else {
		fmt.Printf("Found %d education industry risk factors:\n", len(educationRisks))
		for i, risk := range educationRisks {
			if i < 3 { // Show first 3 risks
				fmt.Printf("  - %s\n", risk)
			}
		}
		if len(educationRisks) > 3 {
			fmt.Printf("  ... and %d more\n", len(educationRisks)-3)
		}
	}
	
	// Test 8: Financial Trading Platform Recommendations
	fmt.Println("\n8. Testing Financial Trading Platform Recommendations:")
	tradingRecs, err := kb.GetIndustrySpecificRecommendations("financial", "trading-platform")
	if err != nil {
		log.Printf("Error getting trading recommendations: %v", err)
	} else {
		fmt.Printf("Found %d financial trading platform recommendations:\n", len(tradingRecs))
		for _, rec := range tradingRecs {
			fmt.Printf("  - %s\n", rec)
		}
	}
	
	// Test 9: PCI-DSS Compliance Requirements
	fmt.Println("\n9. Testing PCI-DSS Compliance Requirements:")
	allCompliance, err := kb.GetComplianceRequirements("")
	if err != nil {
		log.Printf("Error getting all compliance: %v", err)
	} else {
		pciCount := 0
		for _, req := range allCompliance {
			if req.Framework == "PCI-DSS" {
				pciCount++
				fmt.Printf("  - %s: %s\n", req.Requirement, req.Description)
			}
		}
		fmt.Printf("Found %d PCI-DSS requirements\n", pciCount)
	}
	
	// Test 10: HIPAA Compliance Requirements
	fmt.Println("\n10. Testing HIPAA Compliance Requirements:")
	hipaaCount := 0
	for _, req := range allCompliance {
		if req.Framework == "HIPAA" {
			hipaaCount++
			fmt.Printf("  - %s: %s\n", req.Requirement, req.Description)
			fmt.Printf("    Implementation: %v\n", req.Implementation)
		}
	}
	fmt.Printf("Found %d HIPAA requirements\n", hipaaCount)
	
	// Test 11: SOX Compliance Requirements
	fmt.Println("\n11. Testing SOX Compliance Requirements:")
	soxCount := 0
	for _, req := range allCompliance {
		if req.Framework == "SOX" {
			soxCount++
			fmt.Printf("  - %s: %s\n", req.Requirement, req.Description)
		}
	}
	fmt.Printf("Found %d SOX requirements\n", soxCount)
	
	// Test 12: Industry-Specific Architectural Patterns Summary
	fmt.Println("\n12. Testing All Industry-Specific Architectural Patterns:")
	industries := []string{"healthcare", "financial", "retail", "manufacturing", "government", "education"}
	for _, industry := range industries {
		patterns, err := kb.GetIndustrySpecificArchitecturalPatterns(industry)
		if err != nil {
			log.Printf("Error getting %s patterns: %v", industry, err)
		} else {
			fmt.Printf("  %s: %d patterns\n", industry, len(patterns))
		}
	}
	
	// Test 13: Knowledge Base Health Check
	fmt.Println("\n13. Testing Knowledge Base Health:")
	if kb.IsHealthy() {
		fmt.Println("✓ Knowledge base is healthy")
	} else {
		fmt.Println("✗ Knowledge base is not healthy")
	}
	
	// Test 14: General vs Industry-Specific Best Practices
	fmt.Println("\n14. Comparing General vs Industry-Specific Best Practices:")
	generalBP, _ := kb.GetBestPractices("", "")
	healthcareBP, _ := kb.GetIndustrySpecificBestPractices("healthcare")
	fmt.Printf("  General best practices: %d\n", len(generalBP))
	fmt.Printf("  Healthcare-specific best practices: %d\n", len(healthcareBP))
	
	// Test 15: Unknown Industry Fallback
	fmt.Println("\n15. Testing Unknown Industry Fallback:")
	unknownRisks, err := kb.GetIndustryRiskFactors("unknown-industry")
	if err != nil {
		log.Printf("Error getting unknown industry risks: %v", err)
	} else {
		fmt.Printf("Unknown industry returned %d general risk factors\n", len(unknownRisks))
	}
	
	unknownRecs, err := kb.GetIndustrySpecificRecommendations("unknown-industry", "unknown-use-case")
	if err != nil {
		log.Printf("Error getting unknown recommendations: %v", err)
	} else {
		fmt.Printf("Unknown industry/use case returned %d general recommendations\n", len(unknownRecs))
	}
	
	fmt.Println("\n=== Industry-Specific Knowledge System Test Complete ===")
}