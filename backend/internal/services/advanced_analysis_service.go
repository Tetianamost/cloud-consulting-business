package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// AdvancedAnalysisService provides expert-level technical analysis for AWS architectures
type AdvancedAnalysisService struct {
	bedrockService interfaces.BedrockService
	knowledgeBase  interfaces.KnowledgeBase
}

// NewAdvancedAnalysisService creates a new advanced analysis service
func NewAdvancedAnalysisService(
	bedrock interfaces.BedrockService,
	kb interfaces.KnowledgeBase,
) *AdvancedAnalysisService {
	return &AdvancedAnalysisService{
		bedrockService: bedrock,
		knowledgeBase:  kb,
	}
}

// AdvancedAnalysisResult represents the result of advanced architecture analysis
type AdvancedAnalysisResult struct {
	ID                          string                       `json:"id"`
	InquiryID                   string                       `json:"inquiry_id"`
	OverallScore                float64                      `json:"overall_score"`
	TechnicalInsights           []TechnicalInsight           `json:"technical_insights"`
	CostOptimizations           []CostOptimization           `json:"cost_optimizations"`
	SecurityRecommendations     []SecurityRecommendation     `json:"security_recommendations"`
	PerformanceAnalysis         *PerformanceAnalysis         `json:"performance_analysis"`
	ArchitectureRecommendations []ArchitectureRecommendation `json:"architecture_recommendations"`
	// RiskAssessment can be added later if needed
	ROIAnalysis        *ROIAnalysis        `json:"roi_analysis"`
	ImplementationPlan *ImplementationPlan `json:"implementation_plan"`
	CreatedAt          time.Time           `json:"created_at"`
}

// TechnicalInsight represents expert-level technical insights
type TechnicalInsight struct {
	Category        string   `json:"category"` // "architecture", "scalability", "reliability", "security"
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Impact          string   `json:"impact"`
	Severity        string   `json:"severity"` // "low", "medium", "high", "critical"
	Evidence        []string `json:"evidence"`
	Recommendations []string `json:"recommendations"`
	References      []string `json:"references"`
}

// CostOptimization represents specific cost optimization opportunities with dollar amounts
type CostOptimization struct {
	ID                   string   `json:"id"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	CurrentMonthlyCost   float64  `json:"current_monthly_cost"`
	OptimizedMonthlyCost float64  `json:"optimized_monthly_cost"`
	MonthlySavings       float64  `json:"monthly_savings"`
	AnnualSavings        float64  `json:"annual_savings"`
	SavingsPercentage    float64  `json:"savings_percentage"`
	ImplementationSteps  []string `json:"implementation_steps"`
	Effort               string   `json:"effort"` // "low", "medium", "high"
	Risk                 string   `json:"risk"`   // "low", "medium", "high"
	Timeline             string   `json:"timeline"`
	AffectedServices     []string `json:"affected_services"`
}

// SecurityRecommendation represents actionable security remediation steps
type SecurityRecommendation struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Severity            string   `json:"severity"` // "low", "medium", "high", "critical"
	Category            string   `json:"category"` // "access_control", "encryption", "network_security", etc.
	RemediationSteps    []string `json:"remediation_steps"`
	EstimatedEffort     string   `json:"estimated_effort"`
	SecurityImprovement string   `json:"security_improvement"`
	ComplianceImpact    []string `json:"compliance_impact"`
	AffectedComponents  []string `json:"affected_components"`
	References          []string `json:"references"`
}

// PerformanceAnalysis represents performance bottleneck identification and scaling recommendations
type PerformanceAnalysis struct {
	OverallScore              float64                   `json:"overall_score"`
	Bottlenecks               []PerformanceBottleneck   `json:"bottlenecks"`
	ScalingRecommendations    []ScalingRecommendation   `json:"scaling_recommendations"`
	OptimizationOpportunities []PerformanceOptimization `json:"optimization_opportunities"`
}

// PerformanceBottleneck represents identified performance bottlenecks
type PerformanceBottleneck struct {
	ID                 string                 `json:"id"`
	Component          string                 `json:"component"`
	Type               string                 `json:"type"` // "cpu", "memory", "network", "storage", "database"
	Severity           string                 `json:"severity"`
	Description        string                 `json:"description"`
	Impact             string                 `json:"impact"`
	CurrentMetrics     map[string]interface{} `json:"current_metrics"`
	ThresholdMetrics   map[string]interface{} `json:"threshold_metrics"`
	RootCause          string                 `json:"root_cause"`
	AffectedOperations []string               `json:"affected_operations"`
}

// ScalingRecommendation represents scaling recommendations
type ScalingRecommendation struct {
	ID                  string                 `json:"id"`
	Component           string                 `json:"component"`
	ScalingType         string                 `json:"scaling_type"` // "horizontal", "vertical", "auto"
	CurrentCapacity     map[string]interface{} `json:"current_capacity"`
	RecommendedCapacity map[string]interface{} `json:"recommended_capacity"`
	Trigger             string                 `json:"trigger"`
	Justification       string                 `json:"justification"`
	ExpectedBenefit     string                 `json:"expected_benefit"`
	ImplementationSteps []string               `json:"implementation_steps"`
	CostImpact          string                 `json:"cost_impact"`
	Timeline            string                 `json:"timeline"`
}

// PerformanceOptimization represents performance optimization recommendations
type PerformanceOptimization struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	ExpectedImprovement string   `json:"expected_improvement"`
	ImplementationSteps []string `json:"implementation_steps"`
	EstimatedEffort     string   `json:"estimated_effort"`
	EstimatedCost       string   `json:"estimated_cost"`
	RiskLevel           string   `json:"risk_level"`
	ValidationMetrics   []string `json:"validation_metrics"`
}

// ArchitectureRecommendation represents architecture-level recommendations
type ArchitectureRecommendation struct {
	ID                  string   `json:"id"`
	Category            string   `json:"category"` // "reliability", "scalability", "maintainability"
	Priority            string   `json:"priority"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	BusinessImpact      string   `json:"business_impact"`
	TechnicalImpact     string   `json:"technical_impact"`
	ImplementationSteps []string `json:"implementation_steps"`
	EstimatedEffort     string   `json:"estimated_effort"`
	EstimatedCost       string   `json:"estimated_cost"`
	Timeline            string   `json:"timeline"`
	Dependencies        []string `json:"dependencies"`
	RiskLevel           string   `json:"risk_level"`
}

// ROIAnalysis represents return on investment analysis
type ROIAnalysis struct {
	TotalInvestment  float64 `json:"total_investment"`
	AnnualSavings    float64 `json:"annual_savings"`
	PaybackPeriod    string  `json:"payback_period"`
	ROIPercentage    float64 `json:"roi_percentage"`
	NPV              float64 `json:"npv"`
	ThreeYearSavings float64 `json:"three_year_savings"`
	BreakEvenPoint   string  `json:"break_even_point"`
}

// ImplementationPlan represents a comprehensive implementation plan
type ImplementationPlan struct {
	Phases               []ImplementationPhase `json:"phases"`
	TotalTimeline        string                `json:"total_timeline"`
	TotalEstimatedCost   float64               `json:"total_estimated_cost"`
	ResourceRequirements []string              `json:"resource_requirements"`
	RiskMitigation       []string              `json:"risk_mitigation"`
	SuccessMetrics       []string              `json:"success_metrics"`
	Dependencies         []string              `json:"dependencies"`
}

// ImplementationPhase represents a phase in the implementation plan
type ImplementationPhase struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Duration       string   `json:"duration"`
	Prerequisites  []string `json:"prerequisites"`
	Deliverables   []string `json:"deliverables"`
	EstimatedCost  float64  `json:"estimated_cost"`
	RiskLevel      string   `json:"risk_level"`
	SuccessMetrics []string `json:"success_metrics"`
	Dependencies   []string `json:"dependencies"`
}

// AnalyzeArchitectureAdvanced performs comprehensive advanced analysis
func (a *AdvancedAnalysisService) AnalyzeArchitectureAdvanced(ctx context.Context, inquiry *domain.Inquiry) (*AdvancedAnalysisResult, error) {
	analysisID := fmt.Sprintf("advanced-analysis-%d", time.Now().Unix())

	// Generate AI-powered comprehensive analysis
	analysisPrompt := a.buildComprehensiveAnalysisPrompt(inquiry)

	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.1,
		TopP:        0.9,
	}

	response, err := a.bedrockService.GenerateText(ctx, analysisPrompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI analysis: %w", err)
	}

	// Parse AI response into structured analysis
	analysis, err := a.parseAdvancedAnalysis(response.Content, analysisID, inquiry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI analysis: %w", err)
	}

	// Risk assessment can be added later if needed

	// Calculate overall score
	analysis.OverallScore = a.calculateOverallScore(analysis)

	return analysis, nil
}

// buildComprehensiveAnalysisPrompt builds a comprehensive analysis prompt
func (a *AdvancedAnalysisService) buildComprehensiveAnalysisPrompt(inquiry *domain.Inquiry) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are an expert AWS cloud architect and consultant providing advanced technical analysis. ")
	promptBuilder.WriteString("Analyze the following client inquiry and provide expert-level insights that match the sophistication of experienced AWS consultants.\n\n")

	// Add client context
	promptBuilder.WriteString("## CLIENT INQUIRY\n")
	promptBuilder.WriteString(fmt.Sprintf("Client: %s (%s)\n", inquiry.Name, inquiry.Company))
	promptBuilder.WriteString(fmt.Sprintf("Services Requested: %s\n", strings.Join(inquiry.Services, ", ")))
	promptBuilder.WriteString(fmt.Sprintf("Message: %s\n", inquiry.Message))
	promptBuilder.WriteString(fmt.Sprintf("Priority: %s\n\n", inquiry.Priority))

	// Add analysis requirements
	promptBuilder.WriteString("## ANALYSIS REQUIREMENTS\n")
	promptBuilder.WriteString("Provide expert-level analysis covering:\n")
	promptBuilder.WriteString("1. **Technical Insights**: Deep architectural analysis with specific AWS service recommendations\n")
	promptBuilder.WriteString("2. **Cost Optimizations**: Specific savings opportunities with dollar amounts and percentages\n")
	promptBuilder.WriteString("3. **Security Recommendations**: Actionable remediation steps for experienced consultants\n")
	promptBuilder.WriteString("4. **Performance Analysis**: Bottleneck identification and scaling recommendations\n")
	promptBuilder.WriteString("5. **Architecture Recommendations**: Best practices and architectural improvements\n")
	promptBuilder.WriteString("6. **ROI Analysis**: Financial impact and return on investment calculations\n\n")

	promptBuilder.WriteString("## OUTPUT FORMAT\n")
	promptBuilder.WriteString("Provide analysis in the following JSON format:\n")
	promptBuilder.WriteString("{\n")
	promptBuilder.WriteString("  \"technical_insights\": [\n")
	promptBuilder.WriteString("    {\n")
	promptBuilder.WriteString("      \"category\": \"architecture|scalability|reliability|security\",\n")
	promptBuilder.WriteString("      \"title\": \"Insight title\",\n")
	promptBuilder.WriteString("      \"description\": \"Detailed technical description\",\n")
	promptBuilder.WriteString("      \"impact\": \"Business and technical impact\",\n")
	promptBuilder.WriteString("      \"severity\": \"low|medium|high|critical\",\n")
	promptBuilder.WriteString("      \"evidence\": [\"Evidence point 1\", \"Evidence point 2\"],\n")
	promptBuilder.WriteString("      \"recommendations\": [\"Recommendation 1\", \"Recommendation 2\"],\n")
	promptBuilder.WriteString("      \"references\": [\"AWS documentation links\"]\n")
	promptBuilder.WriteString("    }\n")
	promptBuilder.WriteString("  ],\n")
	promptBuilder.WriteString("  \"cost_optimizations\": [\n")
	promptBuilder.WriteString("    {\n")
	promptBuilder.WriteString("      \"title\": \"Optimization title\",\n")
	promptBuilder.WriteString("      \"description\": \"Detailed description\",\n")
	promptBuilder.WriteString("      \"current_monthly_cost\": 5000.0,\n")
	promptBuilder.WriteString("      \"optimized_monthly_cost\": 3500.0,\n")
	promptBuilder.WriteString("      \"monthly_savings\": 1500.0,\n")
	promptBuilder.WriteString("      \"annual_savings\": 18000.0,\n")
	promptBuilder.WriteString("      \"savings_percentage\": 30.0,\n")
	promptBuilder.WriteString("      \"implementation_steps\": [\"Step 1\", \"Step 2\"],\n")
	promptBuilder.WriteString("      \"effort\": \"low|medium|high\",\n")
	promptBuilder.WriteString("      \"risk\": \"low|medium|high\",\n")
	promptBuilder.WriteString("      \"timeline\": \"Implementation timeline\",\n")
	promptBuilder.WriteString("      \"affected_services\": [\"EC2\", \"RDS\"]\n")
	promptBuilder.WriteString("    }\n")
	promptBuilder.WriteString("  ],\n")
	promptBuilder.WriteString("  \"security_recommendations\": [\n")
	promptBuilder.WriteString("    {\n")
	promptBuilder.WriteString("      \"title\": \"Security recommendation title\",\n")
	promptBuilder.WriteString("      \"description\": \"Detailed security issue description\",\n")
	promptBuilder.WriteString("      \"severity\": \"low|medium|high|critical\",\n")
	promptBuilder.WriteString("      \"category\": \"access_control|encryption|network_security|monitoring\",\n")
	promptBuilder.WriteString("      \"remediation_steps\": [\"Step 1\", \"Step 2\"],\n")
	promptBuilder.WriteString("      \"estimated_effort\": \"low|medium|high\",\n")
	promptBuilder.WriteString("      \"security_improvement\": \"Expected security improvement\",\n")
	promptBuilder.WriteString("      \"compliance_impact\": [\"HIPAA\", \"SOC2\"],\n")
	promptBuilder.WriteString("      \"affected_components\": [\"Component 1\", \"Component 2\"],\n")
	promptBuilder.WriteString("      \"references\": [\"Security documentation links\"]\n")
	promptBuilder.WriteString("    }\n")
	promptBuilder.WriteString("  ],\n")
	promptBuilder.WriteString("  \"performance_analysis\": {\n")
	promptBuilder.WriteString("    \"overall_score\": 75.0,\n")
	promptBuilder.WriteString("    \"bottlenecks\": [\n")
	promptBuilder.WriteString("      {\n")
	promptBuilder.WriteString("        \"component\": \"Database\",\n")
	promptBuilder.WriteString("        \"type\": \"cpu|memory|network|storage|database\",\n")
	promptBuilder.WriteString("        \"severity\": \"low|medium|high|critical\",\n")
	promptBuilder.WriteString("        \"description\": \"Bottleneck description\",\n")
	promptBuilder.WriteString("        \"impact\": \"Performance impact\",\n")
	promptBuilder.WriteString("        \"root_cause\": \"Root cause analysis\",\n")
	promptBuilder.WriteString("        \"affected_operations\": [\"Operation 1\", \"Operation 2\"]\n")
	promptBuilder.WriteString("      }\n")
	promptBuilder.WriteString("    ],\n")
	promptBuilder.WriteString("    \"scaling_recommendations\": [\n")
	promptBuilder.WriteString("      {\n")
	promptBuilder.WriteString("        \"component\": \"Web Servers\",\n")
	promptBuilder.WriteString("        \"scaling_type\": \"horizontal|vertical|auto\",\n")
	promptBuilder.WriteString("        \"trigger\": \"Scaling trigger condition\",\n")
	promptBuilder.WriteString("        \"justification\": \"Why scaling is needed\",\n")
	promptBuilder.WriteString("        \"expected_benefit\": \"Expected performance benefit\",\n")
	promptBuilder.WriteString("        \"implementation_steps\": [\"Step 1\", \"Step 2\"],\n")
	promptBuilder.WriteString("        \"cost_impact\": \"Cost impact description\",\n")
	promptBuilder.WriteString("        \"timeline\": \"Implementation timeline\"\n")
	promptBuilder.WriteString("      }\n")
	promptBuilder.WriteString("    ]\n")
	promptBuilder.WriteString("  },\n")
	promptBuilder.WriteString("  \"architecture_recommendations\": [\n")
	promptBuilder.WriteString("    {\n")
	promptBuilder.WriteString("      \"category\": \"reliability|scalability|maintainability\",\n")
	promptBuilder.WriteString("      \"priority\": \"low|medium|high|critical\",\n")
	promptBuilder.WriteString("      \"title\": \"Recommendation title\",\n")
	promptBuilder.WriteString("      \"description\": \"Detailed recommendation\",\n")
	promptBuilder.WriteString("      \"business_impact\": \"Business impact\",\n")
	promptBuilder.WriteString("      \"technical_impact\": \"Technical impact\",\n")
	promptBuilder.WriteString("      \"implementation_steps\": [\"Step 1\", \"Step 2\"],\n")
	promptBuilder.WriteString("      \"estimated_effort\": \"low|medium|high\",\n")
	promptBuilder.WriteString("      \"estimated_cost\": \"Cost estimate\",\n")
	promptBuilder.WriteString("      \"timeline\": \"Implementation timeline\",\n")
	promptBuilder.WriteString("      \"risk_level\": \"low|medium|high\"\n")
	promptBuilder.WriteString("    }\n")
	promptBuilder.WriteString("  ],\n")
	promptBuilder.WriteString("  \"roi_analysis\": {\n")
	promptBuilder.WriteString("    \"total_investment\": 50000.0,\n")
	promptBuilder.WriteString("    \"annual_savings\": 75000.0,\n")
	promptBuilder.WriteString("    \"payback_period\": \"8 months\",\n")
	promptBuilder.WriteString("    \"roi_percentage\": 150.0,\n")
	promptBuilder.WriteString("    \"three_year_savings\": 225000.0\n")
	promptBuilder.WriteString("  }\n")
	promptBuilder.WriteString("}\n\n")

	promptBuilder.WriteString("Generate comprehensive expert-level analysis:")

	return promptBuilder.String()
}

// parseAdvancedAnalysis parses AI response into structured analysis
func (a *AdvancedAnalysisService) parseAdvancedAnalysis(aiResponse, analysisID, inquiryID string) (*AdvancedAnalysisResult, error) {
	// Extract JSON from response
	jsonStart := strings.Index(aiResponse, "{")
	jsonEnd := strings.LastIndex(aiResponse, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no valid JSON found in AI response")
	}

	jsonStr := aiResponse[jsonStart : jsonEnd+1]

	// Parse JSON response
	var aiAnalysis struct {
		TechnicalInsights           []TechnicalInsight           `json:"technical_insights"`
		CostOptimizations           []CostOptimization           `json:"cost_optimizations"`
		SecurityRecommendations     []SecurityRecommendation     `json:"security_recommendations"`
		PerformanceAnalysis         *PerformanceAnalysis         `json:"performance_analysis"`
		ArchitectureRecommendations []ArchitectureRecommendation `json:"architecture_recommendations"`
		ROIAnalysis                 *ROIAnalysis                 `json:"roi_analysis"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &aiAnalysis); err != nil {
		return nil, fmt.Errorf("failed to parse AI analysis JSON: %w", err)
	}

	// Generate unique IDs for all items
	for i := range aiAnalysis.CostOptimizations {
		aiAnalysis.CostOptimizations[i].ID = fmt.Sprintf("cost-opt-%d-%d", time.Now().Unix(), i)
	}

	for i := range aiAnalysis.SecurityRecommendations {
		aiAnalysis.SecurityRecommendations[i].ID = fmt.Sprintf("sec-rec-%d-%d", time.Now().Unix(), i)
	}

	if aiAnalysis.PerformanceAnalysis != nil {
		for i := range aiAnalysis.PerformanceAnalysis.Bottlenecks {
			aiAnalysis.PerformanceAnalysis.Bottlenecks[i].ID = fmt.Sprintf("bottleneck-%d-%d", time.Now().Unix(), i)
		}
		for i := range aiAnalysis.PerformanceAnalysis.ScalingRecommendations {
			aiAnalysis.PerformanceAnalysis.ScalingRecommendations[i].ID = fmt.Sprintf("scaling-%d-%d", time.Now().Unix(), i)
		}
	}

	for i := range aiAnalysis.ArchitectureRecommendations {
		aiAnalysis.ArchitectureRecommendations[i].ID = fmt.Sprintf("arch-rec-%d-%d", time.Now().Unix(), i)
	}

	// Generate implementation plan
	implementationPlan := a.generateImplementationPlan(aiAnalysis.CostOptimizations, aiAnalysis.SecurityRecommendations, aiAnalysis.ArchitectureRecommendations)

	return &AdvancedAnalysisResult{
		ID:                          analysisID,
		InquiryID:                   inquiryID,
		TechnicalInsights:           aiAnalysis.TechnicalInsights,
		CostOptimizations:           aiAnalysis.CostOptimizations,
		SecurityRecommendations:     aiAnalysis.SecurityRecommendations,
		PerformanceAnalysis:         aiAnalysis.PerformanceAnalysis,
		ArchitectureRecommendations: aiAnalysis.ArchitectureRecommendations,
		ROIAnalysis:                 aiAnalysis.ROIAnalysis,
		ImplementationPlan:          implementationPlan,
		CreatedAt:                   time.Now(),
	}, nil
}

// generateImplementationPlan generates a comprehensive implementation plan
func (a *AdvancedAnalysisService) generateImplementationPlan(costOpts []CostOptimization, secRecs []SecurityRecommendation, archRecs []ArchitectureRecommendation) *ImplementationPlan {
	var phases []ImplementationPhase
	totalCost := 0.0

	// Phase 1: Quick Wins (Security and Low-effort Cost Optimizations)
	phase1Items := []string{}
	phase1Cost := 0.0

	for _, secRec := range secRecs {
		if secRec.Severity == "critical" || secRec.Severity == "high" {
			phase1Items = append(phase1Items, fmt.Sprintf("Security: %s", secRec.Title))
			phase1Cost += 2000.0 // Estimated cost for security fixes
		}
	}

	for _, costOpt := range costOpts {
		if costOpt.Effort == "low" && costOpt.Risk == "low" {
			phase1Items = append(phase1Items, fmt.Sprintf("Cost Optimization: %s", costOpt.Title))
		}
	}

	if len(phase1Items) > 0 {
		phases = append(phases, ImplementationPhase{
			ID:             "phase-1-quick-wins",
			Name:           "Quick Wins & Critical Security",
			Description:    "Address critical security issues and implement low-effort cost optimizations",
			Duration:       "2-4 weeks",
			Prerequisites:  []string{"Stakeholder approval", "Access to AWS console"},
			Deliverables:   phase1Items,
			EstimatedCost:  phase1Cost,
			RiskLevel:      "low",
			SuccessMetrics: []string{"All critical security issues resolved", "Cost optimizations implemented"},
			Dependencies:   []string{},
		})
		totalCost += phase1Cost
	}

	// Phase 2: Architecture Improvements
	phase2Items := []string{}
	phase2Cost := 0.0

	for _, archRec := range archRecs {
		if archRec.Priority == "high" || archRec.Priority == "medium" {
			phase2Items = append(phase2Items, archRec.Title)
			phase2Cost += 5000.0 // Estimated cost for architecture changes
		}
	}

	if len(phase2Items) > 0 {
		phases = append(phases, ImplementationPhase{
			ID:             "phase-2-architecture",
			Name:           "Architecture Improvements",
			Description:    "Implement architectural enhancements for scalability and reliability",
			Duration:       "4-8 weeks",
			Prerequisites:  []string{"Phase 1 completion", "Architecture review"},
			Deliverables:   phase2Items,
			EstimatedCost:  phase2Cost,
			RiskLevel:      "medium",
			SuccessMetrics: []string{"Architecture improvements deployed", "Performance metrics improved"},
			Dependencies:   []string{"phase-1-quick-wins"},
		})
		totalCost += phase2Cost
	}

	// Phase 3: Advanced Optimizations
	phase3Items := []string{}
	phase3Cost := 0.0

	for _, costOpt := range costOpts {
		if costOpt.Effort == "high" || costOpt.Risk == "medium" {
			phase3Items = append(phase3Items, fmt.Sprintf("Advanced Cost Optimization: %s", costOpt.Title))
			phase3Cost += 3000.0
		}
	}

	if len(phase3Items) > 0 {
		phases = append(phases, ImplementationPhase{
			ID:             "phase-3-advanced",
			Name:           "Advanced Optimizations",
			Description:    "Implement complex optimizations and advanced features",
			Duration:       "6-12 weeks",
			Prerequisites:  []string{"Phase 2 completion", "Advanced planning"},
			Deliverables:   phase3Items,
			EstimatedCost:  phase3Cost,
			RiskLevel:      "medium",
			SuccessMetrics: []string{"Advanced optimizations implemented", "Target cost savings achieved"},
			Dependencies:   []string{"phase-2-architecture"},
		})
		totalCost += phase3Cost
	}

	return &ImplementationPlan{
		Phases:             phases,
		TotalTimeline:      "3-6 months",
		TotalEstimatedCost: totalCost,
		ResourceRequirements: []string{
			"AWS Solutions Architect",
			"DevOps Engineer",
			"Security Specialist",
			"Project Manager",
		},
		RiskMitigation: []string{
			"Comprehensive testing in staging environment",
			"Gradual rollout with monitoring",
			"Rollback procedures for each change",
			"Regular stakeholder communication",
		},
		SuccessMetrics: []string{
			"All security vulnerabilities addressed",
			"Target cost savings achieved",
			"Performance improvements measured",
			"Zero downtime during implementation",
		},
		Dependencies: []string{
			"Stakeholder approval and budget allocation",
			"Access to AWS accounts and resources",
			"Dedicated implementation team",
		},
	}
}

// calculateOverallScore calculates overall analysis score
func (a *AdvancedAnalysisService) calculateOverallScore(analysis *AdvancedAnalysisResult) float64 {
	baseScore := 70.0 // Start with a baseline

	// Deduct points for critical issues
	for _, insight := range analysis.TechnicalInsights {
		switch insight.Severity {
		case "critical":
			baseScore -= 15.0
		case "high":
			baseScore -= 10.0
		case "medium":
			baseScore -= 5.0
		}
	}

	// Deduct points for security issues
	for _, secRec := range analysis.SecurityRecommendations {
		switch secRec.Severity {
		case "critical":
			baseScore -= 20.0
		case "high":
			baseScore -= 10.0
		case "medium":
			baseScore -= 5.0
		}
	}

	// Add points for cost optimization opportunities (indicates room for improvement)
	totalSavings := 0.0
	for _, costOpt := range analysis.CostOptimizations {
		totalSavings += costOpt.AnnualSavings
	}

	// If there are significant savings opportunities, it indicates current inefficiency
	if totalSavings > 50000 {
		baseScore -= 10.0
	} else if totalSavings > 20000 {
		baseScore -= 5.0
	}

	// Factor in performance analysis if available
	if analysis.PerformanceAnalysis != nil {
		performanceWeight := 0.3
		baseScore = baseScore*(1-performanceWeight) + analysis.PerformanceAnalysis.OverallScore*performanceWeight
	}

	// Ensure score is within bounds
	if baseScore < 0 {
		baseScore = 0
	}
	if baseScore > 100 {
		baseScore = 100
	}

	return baseScore
}
