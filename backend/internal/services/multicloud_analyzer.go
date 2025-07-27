package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MultiCloudAnalyzerService implements the MultiCloudAnalyzer interface
type MultiCloudAnalyzerService struct {
	knowledgeBase interfaces.KnowledgeBase
	docLibrary    interfaces.DocumentationLibrary
}

// NewMultiCloudAnalyzerService creates a new MultiCloudAnalyzerService
func NewMultiCloudAnalyzerService(kb interfaces.KnowledgeBase, docLib interfaces.DocumentationLibrary) *MultiCloudAnalyzerService {
	return &MultiCloudAnalyzerService{
		knowledgeBase: kb,
		docLibrary:    docLib,
	}
}

// CompareServices compares services across cloud providers based on requirements
func (m *MultiCloudAnalyzerService) CompareServices(ctx context.Context, requirement interfaces.ServiceRequirement) (*interfaces.ServiceComparison, error) {
	// Get service options from all major providers
	providers := []string{"aws", "azure", "gcp"}
	var providerOptions []interfaces.ProviderOption
	
	for _, provider := range providers {
		option, err := m.getProviderOption(ctx, provider, requirement)
		if err != nil {
			// Log error but continue with other providers
			continue
		}
		if option != nil {
			providerOptions = append(providerOptions, *option)
		}
	}

	if len(providerOptions) == 0 {
		return nil, fmt.Errorf("no suitable services found for category: %s", requirement.Category)
	}

	// Create comparison matrix
	comparisonMatrix := m.createComparisonMatrix(requirement, providerOptions)
	
	// Calculate scores and determine recommendation
	m.calculateProviderScores(&providerOptions, comparisonMatrix, requirement)
	
	// Sort by score (highest first)
	sort.Slice(providerOptions, func(i, j int) bool {
		return providerOptions[i].Score > providerOptions[j].Score
	})

	recommendation := m.generateRecommendation(providerOptions, requirement)

	return &interfaces.ServiceComparison{
		Category:         requirement.Category,
		Providers:        providerOptions,
		Recommendation:   recommendation,
		Reasoning:        m.generateReasoning(providerOptions[0], requirement),
		ComparisonMatrix: comparisonMatrix,
		CreatedAt:        time.Now(),
	}, nil
}

// AnalyzeCosts analyzes costs across providers for a given workload
func (m *MultiCloudAnalyzerService) AnalyzeCosts(ctx context.Context, workload interfaces.WorkloadSpec) (*interfaces.CostAnalysis, error) {
	providers := []string{"aws", "azure", "gcp"}
	var providerCosts []interfaces.ProviderCostBreakdown
	
	for _, provider := range providers {
		costBreakdown, err := m.calculateProviderCosts(ctx, provider, workload)
		if err != nil {
			// Log error but continue with other providers
			continue
		}
		providerCosts = append(providerCosts, *costBreakdown)
	}

	if len(providerCosts) == 0 {
		return nil, fmt.Errorf("unable to calculate costs for workload: %s", workload.Name)
	}

	// Sort by total monthly cost
	sort.Slice(providerCosts, func(i, j int) bool {
		return providerCosts[i].TotalMonthlyCost < providerCosts[j].TotalMonthlyCost
	})

	recommendation := m.generateCostRecommendation(providerCosts, workload)
	optimizationTips := m.generateCostOptimizationTips(providerCosts, workload)

	return &interfaces.CostAnalysis{
		WorkloadName:     workload.Name,
		AnalysisDate:     time.Now(),
		Currency:         "USD",
		ProviderCosts:    providerCosts,
		Recommendation:   recommendation,
		CostOptimization: optimizationTips,
		Assumptions:      m.getCostAssumptions(),
		Disclaimers:      m.getCostDisclaimers(),
		ValidUntil:       time.Now().AddDate(0, 1, 0), // Valid for 1 month
	}, nil
}

// EvaluateProviders evaluates providers based on specific criteria
func (m *MultiCloudAnalyzerService) EvaluateProviders(ctx context.Context, criteria interfaces.EvaluationCriteria) (*interfaces.ProviderEvaluation, error) {
	providers := []string{"aws", "azure", "gcp"}
	var providerScores []interfaces.ProviderScore
	
	for _, provider := range providers {
		score, err := m.evaluateProvider(ctx, provider, criteria)
		if err != nil {
			// Log error but continue with other providers
			continue
		}
		providerScores = append(providerScores, *score)
	}

	if len(providerScores) == 0 {
		return nil, fmt.Errorf("unable to evaluate providers for use case: %s", criteria.UseCase)
	}

	// Sort by overall score
	sort.Slice(providerScores, func(i, j int) bool {
		return providerScores[i].OverallScore > providerScores[j].OverallScore
	})

	comparisonMatrix := m.createEvaluationMatrix(criteria, providerScores)
	recommendation := m.generateProviderRecommendation(providerScores[0], criteria)

	return &interfaces.ProviderEvaluation{
		UseCase:          criteria.UseCase,
		EvaluationDate:   time.Now(),
		Providers:        providerScores,
		Recommendation:   recommendation,
		ComparisonMatrix: comparisonMatrix,
		Summary:          m.generateEvaluationSummary(providerScores, criteria),
		NextSteps:        m.generateNextSteps(providerScores[0], criteria),
	}, nil
}

// GetMigrationPaths provides migration paths between cloud providers
func (m *MultiCloudAnalyzerService) GetMigrationPaths(ctx context.Context, source, target interfaces.CloudProviderInfo) (*interfaces.MigrationPath, error) {
	if source.Code == target.Code {
		return nil, fmt.Errorf("source and target providers cannot be the same")
	}

	serviceMappings := m.getServiceMappings(source.Code, target.Code)
	migrationPhases := m.getMigrationPhases(source.Code, target.Code)
	toolsAndServices := m.getMigrationTools(source.Code, target.Code)

	complexity := m.calculateMigrationComplexity(serviceMappings)
	duration := m.estimateMigrationDuration(complexity, len(serviceMappings))
	cost := m.estimateMigrationCost(complexity, len(serviceMappings))

	return &interfaces.MigrationPath{
		SourceProvider:      source,
		TargetProvider:      target,
		MigrationStrategy:   m.getMigrationStrategy(source.Code, target.Code),
		EstimatedDuration:   duration,
		EstimatedCost:       cost,
		ComplexityLevel:     complexity,
		ServiceMappings:     serviceMappings,
		MigrationPhases:     migrationPhases,
		RisksAndChallenges:  m.getMigrationRisks(source.Code, target.Code),
		BestPractices:       m.getMigrationBestPractices(source.Code, target.Code),
		ToolsAndServices:    toolsAndServices,
		SuccessFactors:      m.getMigrationSuccessFactors(),
	}, nil
}

// GetProviderRecommendation provides a recommendation based on an inquiry
func (m *MultiCloudAnalyzerService) GetProviderRecommendation(ctx context.Context, inquiry *domain.Inquiry) (*interfaces.ProviderRecommendation, error) {
	// Analyze inquiry to extract requirements
	criteria := m.extractCriteriaFromInquiry(inquiry)
	
	// Evaluate providers
	evaluation, err := m.EvaluateProviders(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate providers: %w", err)
	}

	if len(evaluation.Providers) == 0 {
		return nil, fmt.Errorf("no suitable providers found for inquiry")
	}

	topProvider := evaluation.Providers[0]
	
	// Generate scenarios
	scenarios := m.generateRecommendationScenarios(evaluation.Providers, criteria)
	
	// Generate implementation guidance
	implementation := m.generateImplementationGuidance(topProvider.Provider, criteria)

	return &interfaces.ProviderRecommendation{
		RecommendedProvider: topProvider.Provider,
		Confidence:          m.calculateConfidence(topProvider.OverallScore),
		Reasoning:           topProvider.Recommendations,
		AlternativeOptions:  m.getAlternativeOptions(evaluation.Providers),
		Scenarios:           scenarios,
		Implementation:      implementation,
		CostImplications:    m.generateCostImplications(topProvider.Provider, criteria),
		RiskAssessment:      m.generateRiskAssessment(topProvider.Provider, criteria),
		Timeline:            m.generateImplementationTimeline(topProvider.Provider, criteria),
	}, nil
}

// GetServiceEquivalents returns equivalent services across providers
func (m *MultiCloudAnalyzerService) GetServiceEquivalents(ctx context.Context, provider string, serviceName string) (map[string]string, error) {
	equivalents := make(map[string]string)
	
	// Define service mappings (this would typically come from a knowledge base)
	serviceMappings := m.getServiceEquivalentMappings()
	
	key := fmt.Sprintf("%s:%s", provider, serviceName)
	if mapping, exists := serviceMappings[key]; exists {
		equivalents = mapping
	}

	if len(equivalents) == 0 {
		return nil, fmt.Errorf("no equivalent services found for %s:%s", provider, serviceName)
	}

	return equivalents, nil
}

// Helper methods

func (m *MultiCloudAnalyzerService) getProviderOption(ctx context.Context, provider string, requirement interfaces.ServiceRequirement) (*interfaces.ProviderOption, error) {
	// This would typically query the knowledge base for service information
	serviceInfo := m.getServiceInfoForProvider(provider, requirement.Category)
	if serviceInfo == nil {
		return nil, fmt.Errorf("no service found for provider %s and category %s", provider, requirement.Category)
	}

	// Get documentation links
	docLinks, _ := m.docLibrary.GetDocumentationLinks(ctx, provider, requirement.Category)
	var docURL, pricingURL string
	if len(docLinks) > 0 {
		docURL = docLinks[0].URL
		if len(docLinks) > 1 {
			pricingURL = docLinks[1].URL
		}
	}

	return &interfaces.ProviderOption{
		Provider:                 provider,
		ServiceName:              serviceInfo.ServiceName,
		ServiceDescription:       serviceInfo.Description,
		Pros:                     serviceInfo.Pros,
		Cons:                     serviceInfo.Cons,
		EstimatedMonthlyCost:     serviceInfo.EstimatedCost,
		DocumentationURL:         docURL,
		PricingURL:               pricingURL,
		ImplementationComplexity: serviceInfo.Complexity,
		SupportLevel:             serviceInfo.SupportLevel,
		RegionalAvailability:     serviceInfo.Regions,
		ComplianceCertifications: serviceInfo.Compliance,
		Features:                 serviceInfo.Features,
		Score:                    0, // Will be calculated later
	}, nil
}

func (m *MultiCloudAnalyzerService) createComparisonMatrix(requirement interfaces.ServiceRequirement, options []interfaces.ProviderOption) []interfaces.ComparisonRow {
	criteria := []struct {
		name        string
		weight      float64
		description string
	}{
		{"Cost", 0.3, "Monthly cost estimation"},
		{"Performance", 0.25, "Performance capabilities"},
		{"Compliance", 0.2, "Compliance certifications"},
		{"Ease of Use", 0.15, "Implementation complexity"},
		{"Support", 0.1, "Support level and documentation"},
	}

	var matrix []interfaces.ComparisonRow
	
	for _, criterion := range criteria {
		row := interfaces.ComparisonRow{
			Criteria:    criterion.name,
			Weight:      criterion.weight,
			Scores:      make(map[string]float64),
			Notes:       make(map[string]string),
			Description: criterion.description,
		}

		for _, option := range options {
			score, note := m.calculateCriterionScore(criterion.name, option, requirement)
			row.Scores[option.Provider] = score
			row.Notes[option.Provider] = note
		}

		matrix = append(matrix, row)
	}

	return matrix
}

func (m *MultiCloudAnalyzerService) calculateProviderScores(options *[]interfaces.ProviderOption, matrix []interfaces.ComparisonRow, requirement interfaces.ServiceRequirement) {
	for i := range *options {
		var totalScore float64
		for _, row := range matrix {
			if score, exists := row.Scores[(*options)[i].Provider]; exists {
				totalScore += score * row.Weight
			}
		}
		(*options)[i].Score = totalScore
	}
}

func (m *MultiCloudAnalyzerService) calculateCriterionScore(criterion string, option interfaces.ProviderOption, requirement interfaces.ServiceRequirement) (float64, string) {
	switch criterion {
	case "Cost":
		return m.calculateCostScore(option, requirement), "Based on estimated monthly cost"
	case "Performance":
		return m.calculatePerformanceScore(option, requirement), "Based on performance capabilities"
	case "Compliance":
		return m.calculateComplianceScore(option, requirement), "Based on compliance certifications"
	case "Ease of Use":
		return m.calculateEaseOfUseScore(option), "Based on implementation complexity"
	case "Support":
		return m.calculateSupportScore(option), "Based on support level"
	default:
		return 0.5, "Default score"
	}
}

func (m *MultiCloudAnalyzerService) calculateCostScore(option interfaces.ProviderOption, requirement interfaces.ServiceRequirement) float64 {
	// Parse cost string and compare against budget
	// This is a simplified implementation
	if requirement.Budget.MaxMonthlyCost > 0 {
		// Estimate cost based on service name and complexity
		estimatedCost := m.estimateServiceCost(option.ServiceName, option.ImplementationComplexity)
		if estimatedCost <= requirement.Budget.MaxMonthlyCost {
			return 1.0 - (estimatedCost / requirement.Budget.MaxMonthlyCost)
		}
		return 0.1 // Over budget
	}
	return 0.7 // Default score when no budget specified
}

func (m *MultiCloudAnalyzerService) calculatePerformanceScore(option interfaces.ProviderOption, requirement interfaces.ServiceRequirement) float64 {
	// Evaluate performance based on features and requirements
	score := 0.5 // Base score
	
	// Check if performance requirements are met
	if requirement.Performance.CPU != "" {
		if m.meetsPerformanceRequirement(option.Features, "cpu", requirement.Performance.CPU) {
			score += 0.2
		}
	}
	
	if requirement.Performance.Memory != "" {
		if m.meetsPerformanceRequirement(option.Features, "memory", requirement.Performance.Memory) {
			score += 0.2
		}
	}
	
	if requirement.Performance.Availability != "" {
		if m.meetsPerformanceRequirement(option.Features, "availability", requirement.Performance.Availability) {
			score += 0.1
		}
	}
	
	return math.Min(score, 1.0)
}

func (m *MultiCloudAnalyzerService) calculateComplianceScore(option interfaces.ProviderOption, requirement interfaces.ServiceRequirement) float64 {
	if len(requirement.Compliance) == 0 {
		return 0.8 // Default score when no compliance requirements
	}
	
	matchedCompliance := 0
	for _, reqCompliance := range requirement.Compliance {
		for _, optCompliance := range option.ComplianceCertifications {
			if strings.EqualFold(reqCompliance, optCompliance) {
				matchedCompliance++
				break
			}
		}
	}
	
	return float64(matchedCompliance) / float64(len(requirement.Compliance))
}

func (m *MultiCloudAnalyzerService) calculateEaseOfUseScore(option interfaces.ProviderOption) float64 {
	switch strings.ToLower(option.ImplementationComplexity) {
	case "low":
		return 1.0
	case "medium":
		return 0.6
	case "high":
		return 0.3
	default:
		return 0.5
	}
}

func (m *MultiCloudAnalyzerService) calculateSupportScore(option interfaces.ProviderOption) float64 {
	switch strings.ToLower(option.SupportLevel) {
	case "enterprise":
		return 1.0
	case "business":
		return 0.8
	case "developer":
		return 0.6
	case "basic":
		return 0.4
	default:
		return 0.5
	}
}

func (m *MultiCloudAnalyzerService) generateRecommendation(options []interfaces.ProviderOption, requirement interfaces.ServiceRequirement) string {
	if len(options) == 0 {
		return "No suitable providers found"
	}
	
	top := options[0]
	return fmt.Sprintf("Based on your requirements for %s, %s is recommended with their %s service. This option provides the best balance of cost, performance, and compliance for your use case.",
		requirement.Category, strings.ToUpper(top.Provider), top.ServiceName)
}

func (m *MultiCloudAnalyzerService) generateReasoning(option interfaces.ProviderOption, requirement interfaces.ServiceRequirement) string {
	reasons := []string{
		fmt.Sprintf("Strong performance in %s category", requirement.Category),
		fmt.Sprintf("Implementation complexity is %s", option.ImplementationComplexity),
		fmt.Sprintf("Provides %s support level", option.SupportLevel),
	}
	
	if len(requirement.Compliance) > 0 {
		reasons = append(reasons, fmt.Sprintf("Meets %d of %d compliance requirements", 
			len(option.ComplianceCertifications), len(requirement.Compliance)))
	}
	
	return strings.Join(reasons, "; ")
}

// Additional helper methods for cost analysis, provider evaluation, and migration paths
// These would be implemented based on specific business logic and data sources

func (m *MultiCloudAnalyzerService) calculateProviderCosts(ctx context.Context, provider string, workload interfaces.WorkloadSpec) (*interfaces.ProviderCostBreakdown, error) {
	// Simplified cost calculation - in reality this would use pricing APIs
	var componentCosts []interfaces.ComponentCost
	totalMonthlyCost := 0.0
	
	for _, component := range workload.Components {
		cost := m.estimateComponentCost(provider, component)
		componentCosts = append(componentCosts, interfaces.ComponentCost{
			ComponentName: component.Name,
			ServiceName:   m.getServiceNameForComponent(provider, component.Type),
			MonthlyCost:   cost,
			Unit:          "hours",
			Quantity:      float64(component.Quantity) * 24 * 30, // Monthly hours
			UnitPrice:     cost / (float64(component.Quantity) * 24 * 30),
		})
		totalMonthlyCost += cost
	}
	
	return &interfaces.ProviderCostBreakdown{
		Provider:         provider,
		TotalMonthlyCost: totalMonthlyCost,
		TotalAnnualCost:  totalMonthlyCost * 12,
		ComponentCosts:   componentCosts,
		PricingModel:     "pay-as-you-go",
		LastUpdated:      time.Now(),
	}, nil
}

func (m *MultiCloudAnalyzerService) evaluateProvider(ctx context.Context, provider string, criteria interfaces.EvaluationCriteria) (*interfaces.ProviderScore, error) {
	// Simplified provider evaluation
	categoryScores := make(map[string]float64)
	
	// Evaluate different categories
	categoryScores["cost"] = m.evaluateCostCategory(provider, criteria)
	categoryScores["performance"] = m.evaluatePerformanceCategory(provider, criteria)
	categoryScores["compliance"] = m.evaluateComplianceCategory(provider, criteria)
	categoryScores["support"] = m.evaluateSupportCategory(provider, criteria)
	categoryScores["innovation"] = m.evaluateInnovationCategory(provider, criteria)
	
	// Calculate overall score based on priorities
	overallScore := 0.0
	for category, score := range categoryScores {
		if weight, exists := criteria.Priorities[category]; exists {
			overallScore += score * weight
		} else {
			overallScore += score * 0.2 // Default weight
		}
	}
	
	strengths, weaknesses := m.getProviderStrengthsWeaknesses(provider)
	recommendations := m.getProviderRecommendations(provider, criteria)
	fitScore := m.calculateFitScore(overallScore)
	
	return &interfaces.ProviderScore{
		Provider:        provider,
		OverallScore:    overallScore,
		CategoryScores:  categoryScores,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
		FitScore:        fitScore,
	}, nil
}

// Placeholder implementations for helper methods
// These would be expanded with real business logic and data

func (m *MultiCloudAnalyzerService) getServiceInfoForProvider(provider, category string) *ServiceInfo {
	// This would query the knowledge base
	serviceMap := map[string]map[string]*ServiceInfo{
		"aws": {
			"compute": {
				ServiceName:     "EC2",
				Description:     "Elastic Compute Cloud - scalable virtual servers",
				Pros:           []string{"Highly scalable", "Wide instance types", "Mature ecosystem"},
				Cons:           []string{"Complex pricing", "Steep learning curve"},
				EstimatedCost:  "$50-500/month",
				Complexity:     "medium",
				SupportLevel:   "enterprise",
				Regions:        []string{"us-east-1", "us-west-2", "eu-west-1"},
				Compliance:     []string{"SOC2", "HIPAA", "PCI-DSS"},
				Features:       map[string]string{"auto-scaling": "yes", "load-balancing": "yes"},
			},
			"storage": {
				ServiceName:     "S3",
				Description:     "Simple Storage Service - object storage",
				Pros:           []string{"Highly durable", "Multiple storage classes", "Global availability"},
				Cons:           []string{"Complex permissions", "Data transfer costs"},
				EstimatedCost:  "$20-200/month",
				Complexity:     "low",
				SupportLevel:   "enterprise",
				Regions:        []string{"us-east-1", "us-west-2", "eu-west-1"},
				Compliance:     []string{"SOC2", "HIPAA", "PCI-DSS"},
				Features:       map[string]string{"versioning": "yes", "encryption": "yes"},
			},
		},
		"azure": {
			"compute": {
				ServiceName:     "Virtual Machines",
				Description:     "Scalable virtual machines in the cloud",
				Pros:           []string{"Good Windows integration", "Hybrid capabilities", "Enterprise features"},
				Cons:           []string{"Complex networking", "Limited Linux support"},
				EstimatedCost:  "$45-480/month",
				Complexity:     "medium",
				SupportLevel:   "enterprise",
				Regions:        []string{"eastus", "westus2", "westeurope"},
				Compliance:     []string{"SOC2", "HIPAA", "ISO27001"},
				Features:       map[string]string{"auto-scaling": "yes", "load-balancing": "yes"},
			},
			"storage": {
				ServiceName:     "Blob Storage",
				Description:     "Massively scalable object storage",
				Pros:           []string{"Good integration with Microsoft tools", "Tiered storage", "Strong security"},
				Cons:           []string{"Less mature than competitors", "Complex pricing tiers"},
				EstimatedCost:  "$18-180/month",
				Complexity:     "low",
				SupportLevel:   "enterprise",
				Regions:        []string{"eastus", "westus2", "westeurope"},
				Compliance:     []string{"SOC2", "HIPAA", "ISO27001"},
				Features:       map[string]string{"versioning": "yes", "encryption": "yes"},
			},
		},
		"gcp": {
			"compute": {
				ServiceName:     "Compute Engine",
				Description:     "High-performance virtual machines",
				Pros:           []string{"Competitive pricing", "Good performance", "Strong AI/ML integration"},
				Cons:           []string{"Smaller ecosystem", "Limited enterprise features", "Fewer regions"},
				EstimatedCost:  "$40-450/month",
				Complexity:     "medium",
				SupportLevel:   "business",
				Regions:        []string{"us-central1", "us-west1", "europe-west1"},
				Compliance:     []string{"SOC2", "HIPAA", "ISO27001"},
				Features:       map[string]string{"auto-scaling": "yes", "load-balancing": "yes"},
			},
			"storage": {
				ServiceName:     "Cloud Storage",
				Description:     "Unified object storage for developers and enterprises",
				Pros:           []string{"Simple pricing", "Good performance", "Strong consistency"},
				Cons:           []string{"Limited features", "Smaller ecosystem", "Less mature"},
				EstimatedCost:  "$15-150/month",
				Complexity:     "low",
				SupportLevel:   "business",
				Regions:        []string{"us-central1", "us-west1", "europe-west1"},
				Compliance:     []string{"SOC2", "HIPAA", "ISO27001"},
				Features:       map[string]string{"versioning": "yes", "encryption": "yes"},
			},
		},
	}
	
	if providerServices, exists := serviceMap[provider]; exists {
		if service, exists := providerServices[category]; exists {
			return service
		}
	}
	
	return nil
}

type ServiceInfo struct {
	ServiceName     string
	Description     string
	Pros           []string
	Cons           []string
	EstimatedCost  string
	Complexity     string
	SupportLevel   string
	Regions        []string
	Compliance     []string
	Features       map[string]string
}

func (m *MultiCloudAnalyzerService) estimateServiceCost(serviceName, complexity string) float64 {
	baseCosts := map[string]float64{
		"EC2":              100.0,
		"Virtual Machines": 95.0,
		"Compute Engine":   90.0,
		"S3":               50.0,
		"Blob Storage":     45.0,
		"Cloud Storage":    40.0,
	}
	
	complexityMultiplier := map[string]float64{
		"low":    1.0,
		"medium": 1.5,
		"high":   2.0,
	}
	
	baseCost := baseCosts[serviceName]
	if baseCost == 0 {
		baseCost = 100.0 // Default
	}
	
	multiplier := complexityMultiplier[complexity]
	if multiplier == 0 {
		multiplier = 1.0 // Default
	}
	
	return baseCost * multiplier
}

func (m *MultiCloudAnalyzerService) meetsPerformanceRequirement(features map[string]string, requirement, value string) bool {
	// Simplified performance requirement checking
	if featureValue, exists := features[requirement]; exists {
		return featureValue == "yes" || featureValue == value
	}
	return false
}

func (m *MultiCloudAnalyzerService) extractCriteriaFromInquiry(inquiry *domain.Inquiry) interfaces.EvaluationCriteria {
	// Extract evaluation criteria from inquiry
	criteria := interfaces.EvaluationCriteria{
		UseCase:         strings.Join(inquiry.Services, ", "),
		Industry:        m.inferIndustryFromCompany(inquiry.Company),
		ComplianceNeeds: m.inferComplianceFromMessage(inquiry.Message),
		TechnicalNeeds:  m.inferTechnicalNeeds(inquiry.Message),
		BusinessNeeds:   m.inferBusinessNeeds(inquiry.Message),
		Budget: interfaces.BudgetConstraints{
			CostOptimization: true,
		},
		Timeline:      "3-6 months",
		RiskTolerance: "medium",
		Priorities: map[string]float64{
			"cost":        0.3,
			"performance": 0.25,
			"compliance":  0.2,
			"support":     0.15,
			"innovation":  0.1,
		},
		TeamSkills: []string{"general"},
	}
	
	return criteria
}

// Additional placeholder methods that would be implemented with real business logic

func (m *MultiCloudAnalyzerService) generateCostRecommendation(costs []interfaces.ProviderCostBreakdown, workload interfaces.WorkloadSpec) string {
	if len(costs) == 0 {
		return "No cost analysis available"
	}
	
	cheapest := costs[0]
	return fmt.Sprintf("For workload '%s', %s offers the most cost-effective solution at $%.2f/month", 
		workload.Name, strings.ToUpper(cheapest.Provider), cheapest.TotalMonthlyCost)
}

func (m *MultiCloudAnalyzerService) generateCostOptimizationTips(costs []interfaces.ProviderCostBreakdown, workload interfaces.WorkloadSpec) []interfaces.CostOptimizationTip {
	return []interfaces.CostOptimizationTip{
		{
			Title:            "Use Reserved Instances",
			Description:      "Consider reserved instances for predictable workloads to save up to 60%",
			PotentialSavings: 500.0,
			Effort:          "low",
			Impact:          "high",
			Implementation:  []string{"Analyze usage patterns", "Purchase reserved capacity", "Monitor utilization"},
		},
		{
			Title:            "Implement Auto-scaling",
			Description:      "Use auto-scaling to match capacity with demand",
			PotentialSavings: 300.0,
			Effort:          "medium",
			Impact:          "medium",
			Implementation:  []string{"Configure scaling policies", "Set up monitoring", "Test scaling behavior"},
		},
	}
}

func (m *MultiCloudAnalyzerService) getCostAssumptions() []string {
	return []string{
		"Pricing based on current public rates",
		"Assumes pay-as-you-go pricing model",
		"Does not include data transfer costs",
		"Regional pricing may vary",
	}
}

func (m *MultiCloudAnalyzerService) getCostDisclaimers() []string {
	return []string{
		"Costs are estimates and may vary based on actual usage",
		"Pricing subject to change by cloud providers",
		"Additional costs may apply for premium support",
		"Consult official pricing calculators for accurate quotes",
	}
}

func (m *MultiCloudAnalyzerService) estimateComponentCost(provider string, component interfaces.WorkloadComponent) float64 {
	// Simplified cost estimation
	baseCosts := map[string]map[string]float64{
		"aws": {
			"compute":  0.10,
			"storage":  0.023,
			"database": 0.15,
			"network":  0.05,
		},
		"azure": {
			"compute":  0.096,
			"storage":  0.021,
			"database": 0.14,
			"network":  0.048,
		},
		"gcp": {
			"compute":  0.095,
			"storage":  0.020,
			"database": 0.13,
			"network":  0.045,
		},
	}
	
	if providerCosts, exists := baseCosts[provider]; exists {
		if unitCost, exists := providerCosts[component.Type]; exists {
			return unitCost * float64(component.Quantity) * component.Utilization * 24 * 30 // Monthly cost
		}
	}
	
	return 100.0 // Default monthly cost
}

func (m *MultiCloudAnalyzerService) getServiceNameForComponent(provider, componentType string) string {
	serviceMap := map[string]map[string]string{
		"aws": {
			"compute":  "EC2",
			"storage":  "S3",
			"database": "RDS",
			"network":  "VPC",
		},
		"azure": {
			"compute":  "Virtual Machines",
			"storage":  "Blob Storage",
			"database": "SQL Database",
			"network":  "Virtual Network",
		},
		"gcp": {
			"compute":  "Compute Engine",
			"storage":  "Cloud Storage",
			"database": "Cloud SQL",
			"network":  "VPC",
		},
	}
	
	if providerServices, exists := serviceMap[provider]; exists {
		if service, exists := providerServices[componentType]; exists {
			return service
		}
	}
	
	return componentType
}

func (m *MultiCloudAnalyzerService) evaluateCostCategory(provider string, criteria interfaces.EvaluationCriteria) float64 {
	// Simplified cost evaluation
	costScores := map[string]float64{
		"aws":   0.7,
		"azure": 0.75,
		"gcp":   0.8,
	}
	
	if score, exists := costScores[provider]; exists {
		return score
	}
	return 0.6
}

func (m *MultiCloudAnalyzerService) evaluatePerformanceCategory(provider string, criteria interfaces.EvaluationCriteria) float64 {
	performanceScores := map[string]float64{
		"aws":   0.9,
		"azure": 0.8,
		"gcp":   0.85,
	}
	
	if score, exists := performanceScores[provider]; exists {
		return score
	}
	return 0.7
}

func (m *MultiCloudAnalyzerService) evaluateComplianceCategory(provider string, criteria interfaces.EvaluationCriteria) float64 {
	complianceScores := map[string]float64{
		"aws":   0.95,
		"azure": 0.9,
		"gcp":   0.85,
	}
	
	if score, exists := complianceScores[provider]; exists {
		return score
	}
	return 0.7
}

func (m *MultiCloudAnalyzerService) evaluateSupportCategory(provider string, criteria interfaces.EvaluationCriteria) float64 {
	supportScores := map[string]float64{
		"aws":   0.9,
		"azure": 0.85,
		"gcp":   0.75,
	}
	
	if score, exists := supportScores[provider]; exists {
		return score
	}
	return 0.7
}

func (m *MultiCloudAnalyzerService) evaluateInnovationCategory(provider string, criteria interfaces.EvaluationCriteria) float64 {
	innovationScores := map[string]float64{
		"aws":   0.9,
		"azure": 0.8,
		"gcp":   0.95,
	}
	
	if score, exists := innovationScores[provider]; exists {
		return score
	}
	return 0.7
}

func (m *MultiCloudAnalyzerService) getProviderStrengthsWeaknesses(provider string) ([]string, []string) {
	strengthsWeaknesses := map[string]struct {
		strengths   []string
		weaknesses  []string
	}{
		"aws": {
			strengths:  []string{"Market leader", "Extensive service portfolio", "Strong ecosystem", "Global presence"},
			weaknesses: []string{"Complex pricing", "Steep learning curve", "Vendor lock-in risk"},
		},
		"azure": {
			strengths:  []string{"Strong enterprise integration", "Hybrid capabilities", "Microsoft ecosystem", "Good compliance"},
			weaknesses: []string{"Less mature than AWS", "Complex licensing", "Limited Linux support"},
		},
		"gcp": {
			strengths:  []string{"Competitive pricing", "Strong AI/ML capabilities", "Simple pricing", "Good performance"},
			weaknesses: []string{"Smaller market share", "Limited enterprise features", "Fewer regions"},
		},
	}
	
	if data, exists := strengthsWeaknesses[provider]; exists {
		return data.strengths, data.weaknesses
	}
	
	return []string{"General cloud capabilities"}, []string{"Limited information available"}
}

func (m *MultiCloudAnalyzerService) getProviderRecommendations(provider string, criteria interfaces.EvaluationCriteria) []string {
	recommendations := map[string][]string{
		"aws": {
			"Start with core services like EC2 and S3",
			"Leverage AWS Well-Architected Framework",
			"Consider AWS Training and Certification",
			"Use AWS Cost Explorer for cost optimization",
		},
		"azure": {
			"Leverage existing Microsoft investments",
			"Use Azure Hybrid Benefit for cost savings",
			"Consider Azure Arc for hybrid management",
			"Implement Azure Policy for governance",
		},
		"gcp": {
			"Take advantage of sustained use discounts",
			"Leverage Google's AI/ML capabilities",
			"Use Google Cloud's simple pricing model",
			"Consider multi-region deployments for availability",
		},
	}
	
	if recs, exists := recommendations[provider]; exists {
		return recs
	}
	
	return []string{"Follow cloud best practices", "Implement proper monitoring", "Plan for scalability"}
}

func (m *MultiCloudAnalyzerService) calculateFitScore(overallScore float64) string {
	if overallScore >= 0.8 {
		return "excellent"
	} else if overallScore >= 0.6 {
		return "good"
	} else if overallScore >= 0.4 {
		return "fair"
	}
	return "poor"
}

func (m *MultiCloudAnalyzerService) createEvaluationMatrix(criteria interfaces.EvaluationCriteria, scores []interfaces.ProviderScore) []interfaces.ComparisonRow {
	var matrix []interfaces.ComparisonRow
	
	categories := []string{"cost", "performance", "compliance", "support", "innovation"}
	
	for _, category := range categories {
		row := interfaces.ComparisonRow{
			Criteria:    category,
			Weight:      criteria.Priorities[category],
			Scores:      make(map[string]float64),
			Notes:       make(map[string]string),
			Description: fmt.Sprintf("Evaluation of %s capabilities", category),
		}
		
		for _, score := range scores {
			if categoryScore, exists := score.CategoryScores[category]; exists {
				row.Scores[score.Provider] = categoryScore
				row.Notes[score.Provider] = fmt.Sprintf("Score: %.2f", categoryScore)
			}
		}
		
		matrix = append(matrix, row)
	}
	
	return matrix
}

func (m *MultiCloudAnalyzerService) generateProviderRecommendation(topProvider interfaces.ProviderScore, criteria interfaces.EvaluationCriteria) interfaces.ProviderRecommendation {
	return interfaces.ProviderRecommendation{
		RecommendedProvider: topProvider.Provider,
		Confidence:          m.calculateConfidence(topProvider.OverallScore),
		Reasoning:           topProvider.Recommendations,
		AlternativeOptions:  []string{}, // Would be populated with other providers
		Scenarios:           []interfaces.RecommendationScenario{},
		Implementation:      interfaces.ImplementationGuidance{},
		CostImplications:    fmt.Sprintf("Estimated cost range based on %s pricing model", topProvider.Provider),
		RiskAssessment:      fmt.Sprintf("Risk level: %s based on evaluation criteria", m.calculateRiskLevel(topProvider.OverallScore)),
		Timeline:            criteria.Timeline,
	}
}

func (m *MultiCloudAnalyzerService) generateEvaluationSummary(scores []interfaces.ProviderScore, criteria interfaces.EvaluationCriteria) string {
	if len(scores) == 0 {
		return "No providers evaluated"
	}
	
	top := scores[0]
	return fmt.Sprintf("Based on evaluation criteria for %s, %s scored highest with %.2f overall score, particularly strong in %s capabilities.",
		criteria.UseCase, strings.ToUpper(top.Provider), top.OverallScore, m.getTopCategory(top.CategoryScores))
}

func (m *MultiCloudAnalyzerService) generateNextSteps(topProvider interfaces.ProviderScore, criteria interfaces.EvaluationCriteria) []string {
	return []string{
		fmt.Sprintf("Schedule proof-of-concept with %s", strings.ToUpper(topProvider.Provider)),
		"Conduct detailed cost analysis with actual workload requirements",
		"Review compliance and security requirements in detail",
		"Plan migration strategy and timeline",
		"Identify training and skill development needs",
	}
}

func (m *MultiCloudAnalyzerService) getTopCategory(categoryScores map[string]float64) string {
	var topCategory string
	var topScore float64
	
	for category, score := range categoryScores {
		if score > topScore {
			topScore = score
			topCategory = category
		}
	}
	
	return topCategory
}

func (m *MultiCloudAnalyzerService) calculateConfidence(score float64) string {
	if score >= 0.8 {
		return "high"
	} else if score >= 0.6 {
		return "medium"
	}
	return "low"
}

func (m *MultiCloudAnalyzerService) calculateRiskLevel(score float64) string {
	if score >= 0.8 {
		return "low"
	} else if score >= 0.6 {
		return "medium"
	}
	return "high"
}

// Migration-related helper methods

func (m *MultiCloudAnalyzerService) getServiceMappings(source, target string) []interfaces.ServiceMapping {
	// This would be populated from a comprehensive service mapping database
	mappings := []interfaces.ServiceMapping{
		{
			SourceService:   "EC2",
			TargetService:   "Virtual Machines",
			MappingType:     "direct",
			Compatibility:   "high",
			MigrationNotes:  []string{"Instance types may differ", "Security groups vs NSGs"},
			DataTransfer:    "VM images can be migrated",
			ConfigChanges:   []string{"Update security group rules", "Modify instance metadata"},
		},
		{
			SourceService:   "S3",
			TargetService:   "Blob Storage",
			MappingType:     "direct",
			Compatibility:   "high",
			MigrationNotes:  []string{"Different API endpoints", "Access control differences"},
			DataTransfer:    "Use Azure Data Factory or third-party tools",
			ConfigChanges:   []string{"Update application endpoints", "Modify access policies"},
		},
	}
	
	return mappings
}

func (m *MultiCloudAnalyzerService) getMigrationPhases(source, target string) []interfaces.MigrationPhase {
	return []interfaces.MigrationPhase{
		{
			Name:            "Assessment and Planning",
			Description:     "Analyze current infrastructure and plan migration strategy",
			Duration:        "2-4 weeks",
			Prerequisites:   []string{"Access to current infrastructure", "Stakeholder alignment"},
			Tasks:           []string{"Inventory current resources", "Assess dependencies", "Create migration plan"},
			Deliverables:    []string{"Migration assessment report", "Detailed migration plan", "Risk assessment"},
			RiskLevel:       "low",
			Dependencies:    []string{},
			SuccessCriteria: []string{"Complete inventory", "Approved migration plan", "Risk mitigation strategies"},
		},
		{
			Name:            "Proof of Concept",
			Description:     "Validate migration approach with non-critical workloads",
			Duration:        "2-3 weeks",
			Prerequisites:   []string{"Approved migration plan", "Target environment setup"},
			Tasks:           []string{"Migrate test workload", "Validate functionality", "Performance testing"},
			Deliverables:    []string{"PoC results", "Performance benchmarks", "Lessons learned"},
			RiskLevel:       "medium",
			Dependencies:    []string{"Assessment and Planning"},
			SuccessCriteria: []string{"Successful test migration", "Performance validation", "Process refinement"},
		},
		{
			Name:            "Production Migration",
			Description:     "Migrate production workloads in planned waves",
			Duration:        "4-12 weeks",
			Prerequisites:   []string{"Successful PoC", "Change management approval"},
			Tasks:           []string{"Execute migration waves", "Monitor performance", "Validate functionality"},
			Deliverables:    []string{"Migrated workloads", "Performance reports", "Documentation"},
			RiskLevel:       "high",
			Dependencies:    []string{"Proof of Concept"},
			SuccessCriteria: []string{"All workloads migrated", "Performance targets met", "Zero data loss"},
		},
	}
}

func (m *MultiCloudAnalyzerService) getMigrationTools(source, target string) []interfaces.MigrationTool {
	return []interfaces.MigrationTool{
		{
			Name:              "Azure Migrate",
			Type:              "native",
			Provider:          "Microsoft",
			Description:       "Comprehensive migration service for Azure",
			SupportedServices: []string{"Virtual Machines", "Databases", "Web Apps"},
			Cost:              "Free",
			Documentation:     "https://docs.microsoft.com/azure/migrate/",
			Limitations:       []string{"Azure-specific", "Limited customization"},
		},
		{
			Name:              "CloudEndure",
			Type:              "third-party",
			Provider:          "AWS",
			Description:       "Live migration service for minimal downtime",
			SupportedServices: []string{"Virtual Machines", "Physical Servers"},
			Cost:              "Usage-based pricing",
			Documentation:     "https://docs.cloudendure.com/",
			Limitations:       []string{"Requires agent installation", "Network bandwidth dependent"},
		},
	}
}

func (m *MultiCloudAnalyzerService) calculateMigrationComplexity(mappings []interfaces.ServiceMapping) string {
	if len(mappings) == 0 {
		return "low"
	}
	
	complexityScore := 0
	for _, mapping := range mappings {
		switch mapping.Compatibility {
		case "high":
			complexityScore += 1
		case "medium":
			complexityScore += 2
		case "low":
			complexityScore += 3
		}
	}
	
	avgComplexity := float64(complexityScore) / float64(len(mappings))
	
	if avgComplexity <= 1.5 {
		return "low"
	} else if avgComplexity <= 2.5 {
		return "medium"
	}
	return "high"
}

func (m *MultiCloudAnalyzerService) estimateMigrationDuration(complexity string, serviceCount int) string {
	baseDuration := map[string]int{
		"low":    2,
		"medium": 4,
		"high":   8,
	}
	
	base := baseDuration[complexity]
	if base == 0 {
		base = 4
	}
	
	// Add time based on number of services
	additional := serviceCount / 5 // 1 week per 5 services
	total := base + additional
	
	return fmt.Sprintf("%d-%d weeks", total, total+2)
}

func (m *MultiCloudAnalyzerService) estimateMigrationCost(complexity string, serviceCount int) string {
	baseCost := map[string]int{
		"low":    5000,
		"medium": 15000,
		"high":   30000,
	}
	
	base := baseCost[complexity]
	if base == 0 {
		base = 15000
	}
	
	// Add cost based on number of services
	additional := serviceCount * 2000 // $2k per service
	total := base + additional
	
	return fmt.Sprintf("$%d-$%d", total, int(float64(total)*1.3))
}

func (m *MultiCloudAnalyzerService) getMigrationStrategy(source, target string) string {
	strategies := map[string]map[string]string{
		"aws": {
			"azure": "Lift-and-shift with service mapping and gradual optimization",
			"gcp":   "Hybrid approach with containerization where applicable",
		},
		"azure": {
			"aws": "Direct migration with service equivalency mapping",
			"gcp": "Modernization-focused migration with cloud-native services",
		},
		"gcp": {
			"aws": "Service-by-service migration with performance optimization",
			"azure": "Enterprise-focused migration with hybrid considerations",
		},
	}
	
	if sourceStrategies, exists := strategies[source]; exists {
		if strategy, exists := sourceStrategies[target]; exists {
			return strategy
		}
	}
	
	return "Standard lift-and-shift migration approach"
}

func (m *MultiCloudAnalyzerService) getMigrationRisks(source, target string) []string {
	return []string{
		"Data loss during migration",
		"Extended downtime beyond planned windows",
		"Performance degradation post-migration",
		"Compatibility issues with existing applications",
		"Cost overruns due to unexpected complexity",
		"Security vulnerabilities during transition",
		"Skill gaps in target platform",
	}
}

func (m *MultiCloudAnalyzerService) getMigrationBestPractices(source, target string) []string {
	return []string{
		"Conduct thorough assessment before migration",
		"Start with non-critical workloads for proof of concept",
		"Implement comprehensive backup and rollback procedures",
		"Plan for adequate testing and validation phases",
		"Ensure team training on target platform",
		"Establish monitoring and alerting early",
		"Document all changes and configurations",
		"Plan for post-migration optimization",
	}
}

func (m *MultiCloudAnalyzerService) getMigrationSuccessFactors() []string {
	return []string{
		"Strong executive sponsorship and change management",
		"Dedicated migration team with clear responsibilities",
		"Comprehensive testing strategy and validation procedures",
		"Effective communication plan for all stakeholders",
		"Proper risk management and contingency planning",
		"Adequate budget and timeline allocation",
		"Post-migration support and optimization planning",
	}
}

func (m *MultiCloudAnalyzerService) generateRecommendationScenarios(providers []interfaces.ProviderScore, criteria interfaces.EvaluationCriteria) []interfaces.RecommendationScenario {
	var scenarios []interfaces.RecommendationScenario
	
	for i, provider := range providers {
		if i >= 3 { // Limit to top 3 providers
			break
		}
		
		scenario := interfaces.RecommendationScenario{
			Name:        fmt.Sprintf("%s for %s", strings.ToUpper(provider.Provider), criteria.UseCase),
			Description: fmt.Sprintf("Using %s as primary cloud provider", strings.ToUpper(provider.Provider)),
			Provider:    provider.Provider,
			Conditions:  []string{fmt.Sprintf("Budget allows for %s pricing model", provider.Provider)},
			Benefits:    provider.Strengths,
			Considerations: provider.Weaknesses,
		}
		
		scenarios = append(scenarios, scenario)
	}
	
	return scenarios
}

func (m *MultiCloudAnalyzerService) generateImplementationGuidance(provider string, criteria interfaces.EvaluationCriteria) interfaces.ImplementationGuidance {
	guidanceMap := map[string]interfaces.ImplementationGuidance{
		"aws": {
			QuickStart:      []string{"Create AWS account", "Set up IAM users and roles", "Launch first EC2 instance", "Configure basic monitoring"},
			BestPractices:   []string{"Follow AWS Well-Architected Framework", "Implement least privilege access", "Use Infrastructure as Code", "Enable CloudTrail logging"},
			CommonPitfalls:  []string{"Overly permissive IAM policies", "Not using reserved instances", "Ignoring cost monitoring", "Poor tagging strategy"},
			ResourceLinks:   []string{"https://aws.amazon.com/getting-started/", "https://docs.aws.amazon.com/wellarchitected/"},
			SupportOptions:  []string{"AWS Support Plans", "AWS Professional Services", "AWS Partner Network"},
			TrainingNeeds:   []string{"AWS Cloud Practitioner", "Solutions Architect Associate", "Security Specialty"},
		},
		"azure": {
			QuickStart:      []string{"Create Azure subscription", "Set up Azure AD", "Deploy first virtual machine", "Configure monitoring"},
			BestPractices:   []string{"Use Azure Resource Manager templates", "Implement Azure Policy", "Enable Azure Security Center", "Use managed identities"},
			CommonPitfalls:  []string{"Complex licensing models", "Not leveraging hybrid benefits", "Poor resource organization", "Inadequate backup strategy"},
			ResourceLinks:   []string{"https://docs.microsoft.com/azure/", "https://azure.microsoft.com/architecture/"},
			SupportOptions:  []string{"Azure Support Plans", "Microsoft Consulting Services", "Azure Expert MSPs"},
			TrainingNeeds:   []string{"Azure Fundamentals", "Azure Administrator", "Azure Solutions Architect"},
		},
		"gcp": {
			QuickStart:      []string{"Create GCP project", "Set up billing account", "Launch Compute Engine instance", "Enable monitoring"},
			BestPractices:   []string{"Use Google Cloud Deployment Manager", "Implement IAM best practices", "Enable audit logging", "Use sustained use discounts"},
			CommonPitfalls:  []string{"Not understanding pricing model", "Poor project organization", "Inadequate monitoring", "Missing security configurations"},
			ResourceLinks:   []string{"https://cloud.google.com/docs/", "https://cloud.google.com/architecture/"},
			SupportOptions:  []string{"Google Cloud Support", "Google Professional Services", "Google Cloud Partners"},
			TrainingNeeds:   []string{"Google Cloud Digital Leader", "Professional Cloud Architect", "Professional Cloud Security Engineer"},
		},
	}
	
	if guidance, exists := guidanceMap[provider]; exists {
		return guidance
	}
	
	return interfaces.ImplementationGuidance{
		QuickStart:      []string{"Set up account", "Deploy first resource", "Configure monitoring"},
		BestPractices:   []string{"Follow security best practices", "Implement cost controls", "Use infrastructure as code"},
		CommonPitfalls:  []string{"Poor security configuration", "Lack of monitoring", "Cost overruns"},
		ResourceLinks:   []string{"Official documentation", "Best practices guides"},
		SupportOptions:  []string{"Official support", "Community forums", "Professional services"},
		TrainingNeeds:   []string{"Platform fundamentals", "Architecture certification", "Security training"},
	}
}

func (m *MultiCloudAnalyzerService) getAlternativeOptions(providers []interfaces.ProviderScore) []string {
	var alternatives []string
	
	for i, provider := range providers {
		if i > 0 && i < 3 { // Skip first (recommended) and limit to 2 alternatives
			alternatives = append(alternatives, strings.ToUpper(provider.Provider))
		}
	}
	
	return alternatives
}

func (m *MultiCloudAnalyzerService) generateCostImplications(provider string, criteria interfaces.EvaluationCriteria) string {
	return fmt.Sprintf("Based on %s pricing model, expect initial costs to be moderate with potential for optimization through reserved capacity and right-sizing. Consider %s-specific cost management tools for ongoing optimization.",
		strings.ToUpper(provider), strings.ToUpper(provider))
}

func (m *MultiCloudAnalyzerService) generateRiskAssessment(provider string, criteria interfaces.EvaluationCriteria) string {
	riskLevel := "medium"
	if criteria.RiskTolerance == "low" {
		riskLevel = "low to medium"
	} else if criteria.RiskTolerance == "high" {
		riskLevel = "medium to high"
	}
	
	return fmt.Sprintf("Risk level assessed as %s based on %s maturity, your team's experience, and project complexity. Key risks include vendor lock-in, skill gaps, and potential cost overruns.",
		riskLevel, strings.ToUpper(provider))
}

func (m *MultiCloudAnalyzerService) generateImplementationTimeline(provider string, criteria interfaces.EvaluationCriteria) string {
	return fmt.Sprintf("Estimated implementation timeline: %s. This includes initial setup, migration/deployment, testing, and optimization phases. Timeline may vary based on complexity and team experience with %s.",
		criteria.Timeline, strings.ToUpper(provider))
}

func (m *MultiCloudAnalyzerService) getServiceEquivalentMappings() map[string]map[string]string {
	return map[string]map[string]string{
		"aws:EC2": {
			"azure": "Virtual Machines",
			"gcp":   "Compute Engine",
		},
		"aws:S3": {
			"azure": "Blob Storage",
			"gcp":   "Cloud Storage",
		},
		"aws:RDS": {
			"azure": "SQL Database",
			"gcp":   "Cloud SQL",
		},
		"azure:Virtual Machines": {
			"aws": "EC2",
			"gcp": "Compute Engine",
		},
		"azure:Blob Storage": {
			"aws": "S3",
			"gcp": "Cloud Storage",
		},
		"gcp:Compute Engine": {
			"aws":   "EC2",
			"azure": "Virtual Machines",
		},
		"gcp:Cloud Storage": {
			"aws":   "S3",
			"azure": "Blob Storage",
		},
	}
}

// Helper methods for inquiry analysis

func (m *MultiCloudAnalyzerService) inferIndustryFromCompany(company string) string {
	// Simple industry inference based on company name
	company = strings.ToLower(company)
	
	if strings.Contains(company, "hospital") || strings.Contains(company, "health") || strings.Contains(company, "medical") {
		return "healthcare"
	} else if strings.Contains(company, "bank") || strings.Contains(company, "financial") || strings.Contains(company, "credit") {
		return "financial"
	} else if strings.Contains(company, "retail") || strings.Contains(company, "store") || strings.Contains(company, "shop") {
		return "retail"
	} else if strings.Contains(company, "tech") || strings.Contains(company, "software") || strings.Contains(company, "digital") {
		return "technology"
	}
	
	return "general"
}

func (m *MultiCloudAnalyzerService) inferComplianceFromMessage(message string) []string {
	var compliance []string
	message = strings.ToLower(message)
	
	if strings.Contains(message, "hipaa") || strings.Contains(message, "healthcare") || strings.Contains(message, "patient") {
		compliance = append(compliance, "HIPAA")
	}
	if strings.Contains(message, "pci") || strings.Contains(message, "payment") || strings.Contains(message, "credit card") {
		compliance = append(compliance, "PCI-DSS")
	}
	if strings.Contains(message, "sox") || strings.Contains(message, "financial reporting") {
		compliance = append(compliance, "SOX")
	}
	if strings.Contains(message, "gdpr") || strings.Contains(message, "privacy") || strings.Contains(message, "personal data") {
		compliance = append(compliance, "GDPR")
	}
	
	return compliance
}

func (m *MultiCloudAnalyzerService) inferTechnicalNeeds(message string) []string {
	var needs []string
	message = strings.ToLower(message)
	
	if strings.Contains(message, "scale") || strings.Contains(message, "scalability") {
		needs = append(needs, "scalability")
	}
	if strings.Contains(message, "performance") || strings.Contains(message, "fast") || strings.Contains(message, "speed") {
		needs = append(needs, "high-performance")
	}
	if strings.Contains(message, "backup") || strings.Contains(message, "disaster recovery") || strings.Contains(message, "availability") {
		needs = append(needs, "high-availability")
	}
	if strings.Contains(message, "security") || strings.Contains(message, "secure") {
		needs = append(needs, "security")
	}
	if strings.Contains(message, "integration") || strings.Contains(message, "api") {
		needs = append(needs, "integration")
	}
	
	return needs
}

func (m *MultiCloudAnalyzerService) inferBusinessNeeds(message string) []string {
	var needs []string
	message = strings.ToLower(message)
	
	if strings.Contains(message, "cost") || strings.Contains(message, "budget") || strings.Contains(message, "save") {
		needs = append(needs, "cost-optimization")
	}
	if strings.Contains(message, "quick") || strings.Contains(message, "fast") || strings.Contains(message, "urgent") {
		needs = append(needs, "rapid-deployment")
	}
	if strings.Contains(message, "support") || strings.Contains(message, "help") {
		needs = append(needs, "managed-services")
	}
	if strings.Contains(message, "grow") || strings.Contains(message, "expansion") {
		needs = append(needs, "growth-support")
	}
	
	return needs
}