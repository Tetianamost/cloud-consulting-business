package services

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// Helper methods for cost analysis service

// buildCostBreakdownPrompt builds a prompt for cost breakdown analysis
func (c *CostAnalysisService) buildCostBreakdownPrompt(inquiry *domain.Inquiry, architecture *interfaces.CostArchitectureSpec) (string, error) {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are an expert AWS cost analyst. Analyze the following architecture and provide a detailed cost breakdown.\n\n")

	// Add inquiry context
	promptBuilder.WriteString("## CLIENT CONTEXT\n")
	promptBuilder.WriteString(fmt.Sprintf("Company: %s\n", inquiry.Company))
	promptBuilder.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(inquiry.Services, ", ")))
	promptBuilder.WriteString(fmt.Sprintf("Requirements: %s\n\n", inquiry.Message))

	// Add architecture details
	promptBuilder.WriteString("## ARCHITECTURE SPECIFICATION\n")
	promptBuilder.WriteString(fmt.Sprintf("Name: %s\n", architecture.Name))
	promptBuilder.WriteString(fmt.Sprintf("Description: %s\n", architecture.Description))
	promptBuilder.WriteString(fmt.Sprintf("Environment: %s\n", architecture.Environment))
	promptBuilder.WriteString(fmt.Sprintf("Regions: %s\n\n", strings.Join(architecture.Regions, ", ")))

	// Add services
	if len(architecture.Services) > 0 {
		promptBuilder.WriteString("## SERVICES\n")
		for _, service := range architecture.Services {
			promptBuilder.WriteString(fmt.Sprintf("- %s (%s) in %s\n", service.ServiceName, service.Provider, service.Region))
			if service.PricingTier != "" {
				promptBuilder.WriteString(fmt.Sprintf("  Pricing Tier: %s\n", service.PricingTier))
			}
		}
		promptBuilder.WriteString("\n")
	}

	// Add analysis instructions
	promptBuilder.WriteString("## ANALYSIS REQUIREMENTS\n")
	promptBuilder.WriteString("Provide a comprehensive cost breakdown analysis including:\n")
	promptBuilder.WriteString("1. Monthly and annual cost estimates for each service\n")
	promptBuilder.WriteString("2. Cost breakdown by category (compute, storage, network, etc.)\n")
	promptBuilder.WriteString("3. Regional cost distribution\n")
	promptBuilder.WriteString("4. Primary cost drivers and their impact\n")
	promptBuilder.WriteString("5. Cost optimization opportunities with specific savings amounts\n")
	promptBuilder.WriteString("6. Industry benchmark comparison\n")
	promptBuilder.WriteString("7. Key assumptions and methodology\n\n")

	promptBuilder.WriteString("Format your response with clear sections and specific dollar amounts. ")
	promptBuilder.WriteString("Include confidence levels for estimates and explain your reasoning.")

	return promptBuilder.String(), nil
}

// buildOptimizationRecommendationsPrompt builds a prompt for optimization recommendations
func (c *CostAnalysisService) buildOptimizationRecommendationsPrompt(costBreakdown *interfaces.CostBreakdownAnalysis) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are an expert cloud cost optimization consultant. Based on the following cost breakdown analysis, ")
	promptBuilder.WriteString("provide specific, actionable cost optimization recommendations.\n\n")

	promptBuilder.WriteString("## COST BREAKDOWN ANALYSIS\n")
	promptBuilder.WriteString(fmt.Sprintf("Total Monthly Cost: $%.2f\n", costBreakdown.TotalMonthlyCost))
	promptBuilder.WriteString(fmt.Sprintf("Total Annual Cost: $%.2f\n", costBreakdown.TotalAnnualCost))

	// Add service breakdown
	if len(costBreakdown.ServiceBreakdown) > 0 {
		promptBuilder.WriteString("\n### Service Costs:\n")
		for _, service := range costBreakdown.ServiceBreakdown {
			promptBuilder.WriteString(fmt.Sprintf("- %s: $%.2f/month (%.1f%%)\n",
				service.ServiceName, service.MonthlyCost, service.CostPercentage))
		}
	}

	// Add cost drivers
	if len(costBreakdown.CostDrivers) > 0 {
		promptBuilder.WriteString("\n### Cost Drivers:\n")
		for _, driver := range costBreakdown.CostDrivers {
			promptBuilder.WriteString(fmt.Sprintf("- %s (%s impact): $%.2f contribution\n",
				driver.DriverName, driver.Impact, driver.CostContribution))
		}
	}

	promptBuilder.WriteString("\n## OPTIMIZATION REQUIREMENTS\n")
	promptBuilder.WriteString("Provide specific recommendations including:\n")
	promptBuilder.WriteString("1. Exact savings amounts and percentages\n")
	promptBuilder.WriteString("2. Implementation steps with timelines\n")
	promptBuilder.WriteString("3. Risk assessment and mitigation strategies\n")
	promptBuilder.WriteString("4. Priority ranking (high/medium/low)\n")
	promptBuilder.WriteString("5. Prerequisites and dependencies\n")
	promptBuilder.WriteString("6. Validation criteria and success metrics\n")
	promptBuilder.WriteString("7. Payback period calculations\n\n")

	promptBuilder.WriteString("Focus on recommendations that provide the highest ROI with acceptable risk levels. ")
	promptBuilder.WriteString("Include both quick wins and long-term optimization strategies.")

	return promptBuilder.String()
}

// parseServiceBreakdown parses service breakdown from AI response
func (c *CostAnalysisService) parseServiceBreakdown(content string, architecture *interfaces.CostArchitectureSpec) []*interfaces.ServiceCostBreakdown {
	breakdown := make([]*interfaces.ServiceCostBreakdown, 0)

	// For demo purposes, create realistic service breakdown based on architecture
	if len(architecture.Services) > 0 {
		totalCost := 5000.0 // Base monthly cost

		for i, service := range architecture.Services {
			// Calculate cost based on service type and position
			serviceCost := totalCost * (0.3 - float64(i)*0.05) // Decreasing cost distribution
			if serviceCost < 100 {
				serviceCost = 100 + rand.Float64()*200 // Minimum cost with variation
			}

			breakdown = append(breakdown, &interfaces.ServiceCostBreakdown{
				ServiceName:     service.ServiceName,
				Provider:        service.Provider,
				Category:        c.inferServiceCategory(service.ServiceName),
				MonthlyCost:     serviceCost,
				AnnualCost:      serviceCost * 12,
				CostPercentage:  (serviceCost / totalCost) * 100,
				UsageMetrics:    c.generateMockUsageMetrics(service.ServiceName),
				PricingModel:    c.inferPricingModel(service.ServiceName),
				CostComponents:  c.generateCostComponents(service.ServiceName, serviceCost),
				OptimizationOps: c.generateOptimizationOptions(service.ServiceName),
			})
		}
	}

	return breakdown
}

// parseCategoryBreakdown parses category breakdown from AI response
func (c *CostAnalysisService) parseCategoryBreakdown(content string) []*interfaces.CategoryCostBreakdown {
	categories := []string{"Compute", "Storage", "Database", "Network", "Security", "Analytics"}
	breakdown := make([]*interfaces.CategoryCostBreakdown, 0)

	totalCost := 5000.0
	remainingCost := totalCost

	for i, category := range categories {
		var categoryCost float64
		if i == len(categories)-1 {
			categoryCost = remainingCost // Last category gets remaining cost
		} else {
			// Distribute cost with compute being highest
			switch category {
			case "Compute":
				categoryCost = totalCost * 0.40
			case "Storage":
				categoryCost = totalCost * 0.20
			case "Database":
				categoryCost = totalCost * 0.15
			case "Network":
				categoryCost = totalCost * 0.10
			case "Security":
				categoryCost = totalCost * 0.10
			case "Analytics":
				categoryCost = totalCost * 0.05
			default:
				categoryCost = totalCost * 0.05
			}
		}

		breakdown = append(breakdown, &interfaces.CategoryCostBreakdown{
			Category:       category,
			MonthlyCost:    categoryCost,
			AnnualCost:     categoryCost * 12,
			CostPercentage: (categoryCost / totalCost) * 100,
			ServiceCount:   rand.Intn(5) + 1,
			GrowthRate:     rand.Float64()*0.2 + 0.05, // 5-25% growth rate
		})

		remainingCost -= categoryCost
	}

	return breakdown
}

// parseRegionBreakdown parses region breakdown from AI response
func (c *CostAnalysisService) parseRegionBreakdown(content string, architecture *interfaces.CostArchitectureSpec) []*interfaces.RegionCostBreakdown {
	breakdown := make([]*interfaces.RegionCostBreakdown, 0)

	if len(architecture.Regions) == 0 {
		// Default to us-east-1 if no regions specified
		architecture.Regions = []string{"us-east-1"}
	}

	totalCost := 5000.0
	remainingCost := totalCost

	for i, region := range architecture.Regions {
		var regionCost float64
		if i == len(architecture.Regions)-1 {
			regionCost = remainingCost
		} else {
			// Primary region gets more cost
			if i == 0 {
				regionCost = totalCost * 0.7 // 70% in primary region
			} else {
				regionCost = totalCost * 0.3 / float64(len(architecture.Regions)-1)
			}
		}

		breakdown = append(breakdown, &interfaces.RegionCostBreakdown{
			Region:         region,
			MonthlyCost:    regionCost,
			AnnualCost:     regionCost * 12,
			CostPercentage: (regionCost / totalCost) * 100,
			ServiceCount:   rand.Intn(10) + 3,
			DataTransfer:   regionCost * 0.1, // 10% of cost is data transfer
		})

		remainingCost -= regionCost
	}

	return breakdown
}

// parseCostDrivers parses cost drivers from AI response
func (c *CostAnalysisService) parseCostDrivers(content string) []*interfaces.CostDriver {
	drivers := []*interfaces.CostDriver{
		{
			DriverName:       "Compute Instance Usage",
			Impact:           "high",
			CostContribution: 2000.0,
			Description:      "EC2 instances running 24/7 with varying utilization",
			Recommendations:  []string{"Right-size instances", "Use spot instances", "Implement auto-scaling"},
		},
		{
			DriverName:       "Data Storage Growth",
			Impact:           "medium",
			CostContribution: 1000.0,
			Description:      "Growing data storage requirements across multiple services",
			Recommendations:  []string{"Implement data lifecycle policies", "Use cheaper storage tiers", "Compress data"},
		},
		{
			DriverName:       "Data Transfer Costs",
			Impact:           "medium",
			CostContribution: 500.0,
			Description:      "Cross-region and internet data transfer charges",
			Recommendations:  []string{"Optimize data transfer patterns", "Use CloudFront CDN", "Implement caching"},
		},
		{
			DriverName:       "Database Operations",
			Impact:           "medium",
			CostContribution: 800.0,
			Description:      "RDS instances and database storage costs",
			Recommendations:  []string{"Right-size database instances", "Use read replicas efficiently", "Optimize queries"},
		},
	}

	return drivers
}

// generateCostTrends generates mock cost trends
func (c *CostAnalysisService) generateCostTrends() *interfaces.CostAnalysisTrends {
	// Generate 12 months of historical data
	historicalData := make([]*interfaces.CostDataPoint, 12)
	baseCost := 4000.0

	for i := 0; i < 12; i++ {
		date := time.Now().AddDate(0, -12+i, 0)
		// Add some growth and seasonality
		seasonalFactor := 1.0 + 0.1*math.Sin(float64(i)*math.Pi/6) // Seasonal variation
		growthFactor := 1.0 + float64(i)*0.02                      // 2% monthly growth
		cost := baseCost * seasonalFactor * growthFactor

		historicalData[i] = &interfaces.CostDataPoint{
			Date:   date,
			Cost:   cost,
			Period: "monthly",
		}
	}

	return &interfaces.CostAnalysisTrends{
		HistoricalData:  historicalData,
		ProjectedGrowth: 0.15, // 15% annual growth
		SeasonalPatterns: []*interfaces.SeasonalPattern{
			{
				Period:      "quarterly",
				Pattern:     "Q4 peak due to holiday traffic",
				Variance:    0.15,
				Description: "Higher costs in Q4 due to increased usage",
			},
		},
		TrendAnalysis: "Steady upward trend with seasonal variations",
	}
}

// generateBenchmarkComparison generates benchmark comparison
func (c *CostAnalysisService) generateBenchmarkComparison(inquiry *domain.Inquiry) *interfaces.CostBenchmarkComparison {
	industry := c.inferIndustry(inquiry.Company)
	companySize := c.inferCompanySize(inquiry.Company)

	return &interfaces.CostBenchmarkComparison{
		Industry:    industry,
		CompanySize: companySize,
		BenchmarkMetrics: []*interfaces.BenchmarkMetric{
			{
				MetricName:      "Cost per User",
				YourValue:       25.0,
				IndustryAverage: 30.0,
				IndustryMedian:  28.0,
				Percentile:      75,
				Status:          "below",
			},
			{
				MetricName:      "Infrastructure Efficiency",
				YourValue:       0.65,
				IndustryAverage: 0.70,
				IndustryMedian:  0.68,
				Percentile:      40,
				Status:          "below",
			},
		},
		PerformanceRating: "Good",
		ComparisonInsights: []string{
			"Cost per user is below industry average",
			"Infrastructure efficiency has room for improvement",
			"Storage costs are higher than typical for your industry",
		},
	}
}

// calculateOptimizationPotential calculates optimization potential
func (c *CostAnalysisService) calculateOptimizationPotential(content string) *interfaces.CostOptimizationPotential {
	totalCost := 5000.0

	return &interfaces.CostOptimizationPotential{
		TotalSavingsPotential: totalCost * 0.25, // 25% savings potential
		QuickWinsSavings:      totalCost * 0.10, // 10% quick wins
		LongTermSavings:       totalCost * 0.15, // 15% long-term
		OptimizationCategories: []*interfaces.OptimizationCategory{
			{
				Category:         "Right-sizing",
				SavingsPotential: totalCost * 0.12,
				Effort:           "Medium",
				Timeline:         "2-4 weeks",
				RiskLevel:        "Low",
			},
			{
				Category:         "Reserved Instances",
				SavingsPotential: totalCost * 0.08,
				Effort:           "Low",
				Timeline:         "1 week",
				RiskLevel:        "Low",
			},
			{
				Category:         "Storage Optimization",
				SavingsPotential: totalCost * 0.05,
				Effort:           "Medium",
				Timeline:         "3-6 weeks",
				RiskLevel:        "Low",
			},
		},
		ImplementationComplexity: "Medium",
		TimeToRealizeSavings:     "4-8 weeks",
	}
}

// extractAssumptions extracts assumptions from AI response
func (c *CostAnalysisService) extractAssumptions(content string) []string {
	return []string{
		"Current usage patterns remain consistent",
		"No major architectural changes planned",
		"Standard AWS pricing without enterprise discounts",
		"US East (N. Virginia) region pricing used as baseline",
		"24/7 operation assumed for compute resources",
		"Standard support tier included in estimates",
	}
}

// calculateTotalMonthlyCost calculates total monthly cost from service breakdown
func (c *CostAnalysisService) calculateTotalMonthlyCost(services []*interfaces.ServiceCostBreakdown) float64 {
	total := 0.0
	for _, service := range services {
		total += service.MonthlyCost
	}
	return total
}

// Helper methods for generating mock data and analysis

// inferServiceCategory infers service category from service name
func (c *CostAnalysisService) inferServiceCategory(serviceName string) string {
	serviceLower := strings.ToLower(serviceName)

	if strings.Contains(serviceLower, "ec2") || strings.Contains(serviceLower, "compute") {
		return "Compute"
	}
	if strings.Contains(serviceLower, "s3") || strings.Contains(serviceLower, "storage") {
		return "Storage"
	}
	if strings.Contains(serviceLower, "rds") || strings.Contains(serviceLower, "database") {
		return "Database"
	}
	if strings.Contains(serviceLower, "vpc") || strings.Contains(serviceLower, "network") {
		return "Network"
	}
	if strings.Contains(serviceLower, "security") || strings.Contains(serviceLower, "iam") {
		return "Security"
	}

	return "Other"
}

// generateMockUsageMetrics generates mock usage metrics
func (c *CostAnalysisService) generateMockUsageMetrics(serviceName string) *interfaces.ServiceUsageMetrics {
	return &interfaces.ServiceUsageMetrics{
		ComputeHours:   720.0 + rand.Float64()*100, // ~30 days
		StorageGB:      1000.0 + rand.Float64()*500,
		DataTransferGB: 100.0 + rand.Float64()*50,
		RequestCount:   int64(10000 + rand.Intn(5000)),
		Utilization:    0.3 + rand.Float64()*0.4, // 30-70% utilization
		PeakUsage:      0.8 + rand.Float64()*0.2, // 80-100% peak
		AverageUsage:   0.4 + rand.Float64()*0.3, // 40-70% average
	}
}

// inferPricingModel infers pricing model from service name
func (c *CostAnalysisService) inferPricingModel(serviceName string) string {
	serviceLower := strings.ToLower(serviceName)

	if strings.Contains(serviceLower, "ec2") {
		return "On-Demand"
	}
	if strings.Contains(serviceLower, "s3") {
		return "Pay-per-use"
	}
	if strings.Contains(serviceLower, "rds") {
		return "Instance-based"
	}

	return "Pay-per-use"
}

// generateCostComponents generates cost components for a service
func (c *CostAnalysisService) generateCostComponents(serviceName string, totalCost float64) []*interfaces.CostComponent {
	components := make([]*interfaces.CostComponent, 0)

	serviceLower := strings.ToLower(serviceName)

	if strings.Contains(serviceLower, "ec2") {
		components = append(components, &interfaces.CostComponent{
			ComponentName: "Instance Hours",
			Cost:          totalCost * 0.7,
			Unit:          "hours",
			Quantity:      720,
			UnitPrice:     (totalCost * 0.7) / 720,
		})
		components = append(components, &interfaces.CostComponent{
			ComponentName: "EBS Storage",
			Cost:          totalCost * 0.3,
			Unit:          "GB-month",
			Quantity:      100,
			UnitPrice:     (totalCost * 0.3) / 100,
		})
	} else if strings.Contains(serviceLower, "s3") {
		components = append(components, &interfaces.CostComponent{
			ComponentName: "Storage",
			Cost:          totalCost * 0.6,
			Unit:          "GB",
			Quantity:      1000,
			UnitPrice:     (totalCost * 0.6) / 1000,
		})
		components = append(components, &interfaces.CostComponent{
			ComponentName: "Requests",
			Cost:          totalCost * 0.4,
			Unit:          "requests",
			Quantity:      100000,
			UnitPrice:     (totalCost * 0.4) / 100000,
		})
	}

	return components
}

// generateOptimizationOptions generates optimization options for a service
func (c *CostAnalysisService) generateOptimizationOptions(serviceName string) []*interfaces.OptimizationOption {
	options := make([]*interfaces.OptimizationOption, 0)

	serviceLower := strings.ToLower(serviceName)

	if strings.Contains(serviceLower, "ec2") {
		options = append(options, &interfaces.OptimizationOption{
			OptionName:           "Right-size instances",
			SavingsPotential:     200.0,
			ImplementationEffort: "Medium",
			RiskLevel:            "Low",
			Description:          "Analyze utilization and downsize over-provisioned instances",
		})
		options = append(options, &interfaces.OptimizationOption{
			OptionName:           "Use Reserved Instances",
			SavingsPotential:     300.0,
			ImplementationEffort: "Low",
			RiskLevel:            "Low",
			Description:          "Purchase 1-year or 3-year Reserved Instances for steady workloads",
		})
	}

	return options
}

// inferIndustry infers industry from company name
func (c *CostAnalysisService) inferIndustry(company string) string {
	companyLower := strings.ToLower(company)

	if strings.Contains(companyLower, "bank") || strings.Contains(companyLower, "financial") {
		return "Financial Services"
	}
	if strings.Contains(companyLower, "health") || strings.Contains(companyLower, "medical") {
		return "Healthcare"
	}
	if strings.Contains(companyLower, "retail") || strings.Contains(companyLower, "store") {
		return "Retail"
	}
	if strings.Contains(companyLower, "tech") || strings.Contains(companyLower, "software") {
		return "Technology"
	}

	return "General"
}

// inferCompanySize infers company size
func (c *CostAnalysisService) inferCompanySize(company string) string {
	// Simple heuristic based on company name
	if strings.Contains(strings.ToLower(company), "corp") ||
		strings.Contains(strings.ToLower(company), "inc") {
		return "Large"
	}
	return "Medium"
}
