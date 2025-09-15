package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// awsServiceIntelligenceService implements the AWSServiceIntelligenceService interface
type awsServiceIntelligenceService struct {
	bedrockService interfaces.BedrockService
	lastUpdate     time.Time
	serviceCache   map[string]*interfaces.AWSServiceStatus
}

// NewAWSServiceIntelligenceService creates a new AWS service intelligence service
func NewAWSServiceIntelligenceService(bedrockService interfaces.BedrockService) interfaces.AWSServiceIntelligenceService {
	return &awsServiceIntelligenceService{
		bedrockService: bedrockService,
		lastUpdate:     time.Now(),
		serviceCache:   make(map[string]*interfaces.AWSServiceStatus),
	}
}

// GetServiceStatus retrieves current AWS service status for a region
func (s *awsServiceIntelligenceService) GetServiceStatus(ctx context.Context, region string) (*interfaces.AWSServiceStatus, error) {
	// Check cache first
	if cached, exists := s.serviceCache[region]; exists {
		if time.Since(cached.LastUpdated) < 5*time.Minute {
			return cached, nil
		}
	}

	// Generate service status using Bedrock
	prompt := s.buildServiceStatusPrompt(region)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   2000,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get service status from Bedrock: %w", err)
	}

	// Parse the response into service status
	status, err := s.parseServiceStatusResponse(response.Content, region)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service status response: %w", err)
	}

	// Cache the result
	s.serviceCache[region] = status

	return status, nil
}

// GetServiceStatusHistory retrieves historical status events for a service
func (s *awsServiceIntelligenceService) GetServiceStatusHistory(ctx context.Context, service, region string, duration time.Duration) ([]*interfaces.ServiceStatusEvent, error) {
	prompt := s.buildServiceHistoryPrompt(service, region, duration)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   3000,
		Temperature: 0.2,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get service history from Bedrock: %w", err)
	}

	events, err := s.parseServiceHistoryResponse(response.Content, service, region)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service history response: %w", err)
	}

	return events, nil
}

// AnalyzeServiceImpact analyzes the impact of service issues on client recommendations
func (s *awsServiceIntelligenceService) AnalyzeServiceImpact(ctx context.Context, service, region string) (*interfaces.ServiceImpactAnalysis, error) {
	prompt := s.buildServiceImpactPrompt(service, region)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   2500,
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze service impact from Bedrock: %w", err)
	}

	analysis, err := s.parseServiceImpactResponse(response.Content, service, region)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service impact response: %w", err)
	}

	return analysis, nil
}

// GetNewServices retrieves newly announced AWS services since a given time
func (s *awsServiceIntelligenceService) GetNewServices(ctx context.Context, since time.Time) ([]*interfaces.NewAWSService, error) {
	prompt := s.buildNewServicesPrompt(since)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   4000,
		Temperature: 0.3,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get new services from Bedrock: %w", err)
	}

	services, err := s.parseNewServicesResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new services response: %w", err)
	}

	return services, nil
}

// EvaluateServiceForClient evaluates a new service for a specific client context
func (s *awsServiceIntelligenceService) EvaluateServiceForClient(ctx context.Context, service *interfaces.NewAWSService, clientContext *interfaces.ClientContext) (*interfaces.ServiceEvaluation, error) {
	prompt := s.buildServiceEvaluationPrompt(service, clientContext)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   3500,
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate service from Bedrock: %w", err)
	}

	evaluation, err := s.parseServiceEvaluationResponse(response.Content, service.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service evaluation response: %w", err)
	}

	return evaluation, nil
}

// GetServiceRecommendations gets service recommendations for a client context
func (s *awsServiceIntelligenceService) GetServiceRecommendations(ctx context.Context, clientContext *interfaces.ClientContext) ([]*interfaces.ServiceRecommendation, error) {
	prompt := s.buildServiceRecommendationsPrompt(clientContext)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   4000,
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get service recommendations from Bedrock: %w", err)
	}

	recommendations, err := s.parseServiceRecommendationsResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service recommendations response: %w", err)
	}

	return recommendations, nil
}

// GetDeprecationAlerts retrieves current service deprecation alerts
func (s *awsServiceIntelligenceService) GetDeprecationAlerts(ctx context.Context) ([]*interfaces.DeprecationAlert, error) {
	prompt := s.buildDeprecationAlertsPrompt()

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   3500,
		Temperature: 0.2,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get deprecation alerts from Bedrock: %w", err)
	}

	alerts, err := s.parseDeprecationAlertsResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deprecation alerts response: %w", err)
	}

	return alerts, nil
}

// GenerateMigrationPlan generates a migration plan for deprecated services
func (s *awsServiceIntelligenceService) GenerateMigrationPlan(ctx context.Context, deprecatedService string, clientContext *interfaces.ClientContext) (*interfaces.MigrationPlan, error) {
	prompt := s.buildMigrationPlanPrompt(deprecatedService, clientContext)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   4000,
		Temperature: 0.3,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate migration plan from Bedrock: %w", err)
	}

	plan, err := s.parseMigrationPlanResponse(response.Content, deprecatedService)
	if err != nil {
		return nil, fmt.Errorf("failed to parse migration plan response: %w", err)
	}

	return plan, nil
}

// AnalyzePricingChanges analyzes AWS pricing changes since a given time
func (s *awsServiceIntelligenceService) AnalyzePricingChanges(ctx context.Context, since time.Time) ([]*interfaces.PricingChange, error) {
	prompt := s.buildPricingChangesPrompt(since)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   3500,
		Temperature: 0.2,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze pricing changes from Bedrock: %w", err)
	}

	changes, err := s.parsePricingChangesResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pricing changes response: %w", err)
	}

	return changes, nil
}

// CalculateCostImpact calculates cost impact of pricing changes for a client
func (s *awsServiceIntelligenceService) CalculateCostImpact(ctx context.Context, pricingChange *interfaces.PricingChange, clientContext *interfaces.ClientContext) (*interfaces.PricingCostImpactAnalysis, error) {
	prompt := s.buildCostImpactPrompt(pricingChange, clientContext)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   3000,
		Temperature: 0.3,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate cost impact from Bedrock: %w", err)
	}

	analysis, err := s.parseCostImpactResponse(response.Content, pricingChange.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cost impact response: %w", err)
	}

	return analysis, nil
}

// RefreshServiceIntelligence refreshes the service intelligence data
func (s *awsServiceIntelligenceService) RefreshServiceIntelligence(ctx context.Context) error {
	// Clear cache to force refresh
	s.serviceCache = make(map[string]*interfaces.AWSServiceStatus)
	s.lastUpdate = time.Now()
	return nil
}

// GetLastUpdateTime returns the last update time
func (s *awsServiceIntelligenceService) GetLastUpdateTime() time.Time {
	return s.lastUpdate
}

// IsHealthy checks if the service is healthy
func (s *awsServiceIntelligenceService) IsHealthy() bool {
	return s.bedrockService.IsHealthy()
}

// Prompt building methods

func (s *awsServiceIntelligenceService) buildServiceStatusPrompt(region string) string {
	return fmt.Sprintf(`As an AWS cloud consultant, provide current AWS service status information for the %s region.

Please provide a comprehensive status report including:
1. Overall health status of the region
2. Status of major AWS services (EC2, S3, RDS, Lambda, etc.)
3. Any current incidents or degraded performance
4. Recent service updates or improvements

Format the response as a structured analysis that includes:
- Service name and current status
- Any performance metrics or availability data
- Recent incidents or maintenance windows
- Impact assessment for client workloads

Focus on services commonly used in enterprise cloud consulting scenarios.`, region)
}

func (s *awsServiceIntelligenceService) buildServiceHistoryPrompt(service, region string, duration time.Duration) string {
	return fmt.Sprintf(`As an AWS cloud consultant, provide historical status information for the %s service in the %s region over the past %v.

Please include:
1. Major incidents or outages
2. Performance degradation events
3. Maintenance windows
4. Service improvements or updates
5. Impact on client workloads

Format as a chronological list of events with:
- Timestamp
- Event type and severity
- Duration of impact
- Root cause (if known)
- Resolution details`, service, region, duration)
}

func (s *awsServiceIntelligenceService) buildServiceImpactPrompt(service, region string) string {
	return fmt.Sprintf(`As an AWS cloud consultant, analyze the current impact of %s service issues in %s region on client recommendations and workloads.

Please provide:
1. Current service status and any ongoing issues
2. Impact level on different types of workloads
3. Affected client scenarios and use cases
4. Alternative services or workarounds
5. Recommended actions for consultants
6. Timeline for resolution (if applicable)

Focus on practical consulting implications and client communication needs.`, service, region)
}

func (s *awsServiceIntelligenceService) buildNewServicesPrompt(since time.Time) string {
	return fmt.Sprintf(`As an AWS cloud consultant, provide information about new AWS services and features announced since %s.

Please include:
1. Service name and category
2. Key features and capabilities
3. Target use cases and industries
4. Pricing model and availability
5. Competitive advantages
6. Integration with existing services

Focus on services that would be relevant for enterprise cloud consulting and client recommendations.`, since.Format("2006-01-02"))
}
func (s *awsServiceIntelligenceService) buildServiceEvaluationPrompt(service *interfaces.NewAWSService, clientContext *interfaces.ClientContext) string {
	return fmt.Sprintf(`As an AWS cloud consultant, evaluate the new AWS service "%s" for a client with the following context:

Client Profile:
- Industry: %s
- Company Size: %s
- Current Services: %s
- Preferred Regions: %s
- Compliance Requirements: %s
- Technical Maturity: %s

Service Details:
- Description: %s
- Key Features: %s
- Use Cases: %s

Please provide:
1. Relevance score (0-100) for this client
2. Adoption priority (immediate/short-term/long-term/not-applicable)
3. Expected benefits for this client
4. Implementation challenges
5. Prerequisites and dependencies
6. Cost implications and ROI analysis
7. Risk assessment
8. Specific recommendations and next steps

Focus on practical consulting advice tailored to this client's context.`,
		service.ServiceName,
		clientContext.IndustryVertical,
		clientContext.CompanySize,
		strings.Join(clientContext.CurrentServices, ", "),
		strings.Join(clientContext.PreferredRegions, ", "),
		strings.Join(clientContext.ComplianceRequirements, ", "),
		clientContext.TechnicalMaturity,
		service.Description,
		strings.Join(service.KeyFeatures, ", "),
		strings.Join(service.UseCases, ", "))
}

func (s *awsServiceIntelligenceService) buildServiceRecommendationsPrompt(clientContext *interfaces.ClientContext) string {
	return fmt.Sprintf(`As an AWS cloud consultant, provide service recommendations for a client with the following profile:

Client Context:
- Industry: %s
- Company Size: %s
- Current Services: %s
- Preferred Regions: %s
- Compliance Requirements: %s
- Technical Maturity: %s
- Workloads: %d different workload types

Please recommend:
1. New services to adopt
2. Existing services to upgrade or optimize
3. Services to replace with better alternatives
4. Cost optimization opportunities

For each recommendation, include:
- Service name and type
- Priority level and rationale
- Expected benefits
- Implementation timeline
- Cost-benefit analysis
- Dependencies and prerequisites

Focus on practical, actionable recommendations that align with the client's maturity and requirements.`,
		clientContext.IndustryVertical,
		clientContext.CompanySize,
		strings.Join(clientContext.CurrentServices, ", "),
		strings.Join(clientContext.PreferredRegions, ", "),
		strings.Join(clientContext.ComplianceRequirements, ", "),
		clientContext.TechnicalMaturity,
		len(clientContext.Workloads))
}

func (s *awsServiceIntelligenceService) buildDeprecationAlertsPrompt() string {
	return `As an AWS cloud consultant, provide current information about AWS service deprecations, end-of-life announcements, and feature removals.

Please include:
1. Services being deprecated or reaching end-of-life
2. Features being removed or changed
3. API versions being deprecated
4. Timeline and effective dates
5. Impact assessment for different client types
6. Recommended migration paths
7. Alternative services or solutions

Focus on deprecations that would impact enterprise clients and require proactive consulting intervention.`
}

func (s *awsServiceIntelligenceService) buildMigrationPlanPrompt(deprecatedService string, clientContext *interfaces.ClientContext) string {
	return fmt.Sprintf(`As an AWS cloud consultant, create a detailed migration plan for moving away from the deprecated service "%s" for a client with the following context:

Client Profile:
- Industry: %s
- Company Size: %s
- Current Services: %s
- Technical Maturity: %s
- Compliance Requirements: %s

Please provide:
1. Recommended target service(s)
2. Migration strategy (lift-and-shift, re-architect, hybrid)
3. Detailed migration steps with timeline
4. Cost analysis and comparison
5. Risk assessment and mitigation strategies
6. Testing and validation plan
7. Rollback procedures
8. Resource requirements
9. Training and documentation needs

Focus on minimizing business disruption while ensuring a smooth transition.`,
		deprecatedService,
		clientContext.IndustryVertical,
		clientContext.CompanySize,
		strings.Join(clientContext.CurrentServices, ", "),
		clientContext.TechnicalMaturity,
		strings.Join(clientContext.ComplianceRequirements, ", "))
}
func (s *awsServiceIntelligenceService) buildPricingChangesPrompt(since time.Time) string {
	return fmt.Sprintf(`As an AWS cloud consultant, provide information about AWS pricing changes announced since %s.

Please include:
1. Services with pricing changes
2. Type of change (increase, decrease, new pricing tiers)
3. Effective dates and announcement dates
4. Percentage changes and impact levels
5. Affected regions and service tiers
6. Rationale for changes (if provided by AWS)
7. Impact on different customer segments

Focus on changes that would significantly impact enterprise clients and require proactive cost management.`, since.Format("2006-01-02"))
}

func (s *awsServiceIntelligenceService) buildCostImpactPrompt(pricingChange *interfaces.PricingChange, clientContext *interfaces.ClientContext) string {
	return fmt.Sprintf(`As an AWS cloud consultant, analyze the cost impact of the following pricing change for a specific client:

Pricing Change:
- Service: %s
- Region: %s
- Change Type: %s
- Description: %s
- Effective Date: %s

Client Context:
- Industry: %s
- Company Size: %s
- Current Services: %s
- Technical Maturity: %s

Please provide:
1. Estimated current monthly cost for this service
2. Projected new monthly cost after the change
3. Absolute and percentage cost difference
4. Annual impact projection
5. Cost optimization recommendations
6. Alternative service options
7. Action items and timeline
8. Risk assessment if no action is taken

Focus on practical cost management strategies and client communication points.`,
		pricingChange.ServiceName,
		pricingChange.Region,
		pricingChange.ChangeType,
		pricingChange.ChangeDescription,
		pricingChange.EffectiveDate.Format("2006-01-02"),
		clientContext.IndustryVertical,
		clientContext.CompanySize,
		strings.Join(clientContext.CurrentServices, ", "),
		clientContext.TechnicalMaturity)
}

// Response parsing methods

func (s *awsServiceIntelligenceService) parseServiceStatusResponse(content, region string) (*interfaces.AWSServiceStatus, error) {
	// Create a basic service status structure
	// In a real implementation, this would parse structured data from Bedrock
	status := &interfaces.AWSServiceStatus{
		Region:        region,
		LastUpdated:   time.Now(),
		Services:      make(map[string]*interfaces.ServiceHealth),
		OverallHealth: "healthy", // Default assumption
	}

	// Parse common AWS services from the response
	services := []string{"EC2", "S3", "RDS", "Lambda", "ECS", "EKS", "CloudFormation", "IAM"}

	for _, serviceName := range services {
		serviceHealth := &interfaces.ServiceHealth{
			ServiceName: serviceName,
			Status:      "operational", // Default
			Regions:     map[string]string{region: "operational"},
			Metadata:    make(map[string]interface{}),
		}

		// Simple parsing logic - in production this would be more sophisticated
		if strings.Contains(strings.ToLower(content), strings.ToLower(serviceName)) {
			if strings.Contains(strings.ToLower(content), "degraded") ||
				strings.Contains(strings.ToLower(content), "issue") {
				serviceHealth.Status = "degraded"
				status.OverallHealth = "degraded"
			}
		}

		status.Services[serviceName] = serviceHealth
	}

	return status, nil
}

func (s *awsServiceIntelligenceService) parseServiceHistoryResponse(content, service, region string) ([]*interfaces.ServiceStatusEvent, error) {
	// Create sample historical events
	// In a real implementation, this would parse structured data from Bedrock
	events := []*interfaces.ServiceStatusEvent{
		{
			ServiceName: service,
			Region:      region,
			Status:      "operational",
			Message:     "Service operating normally",
			Timestamp:   time.Now().Add(-24 * time.Hour),
		},
	}

	// Parse content for incidents or issues
	if strings.Contains(strings.ToLower(content), "incident") ||
		strings.Contains(strings.ToLower(content), "outage") {
		duration := 2 * time.Hour
		events = append(events, &interfaces.ServiceStatusEvent{
			ServiceName: service,
			Region:      region,
			Status:      "disrupted",
			Message:     "Service incident detected",
			Timestamp:   time.Now().Add(-48 * time.Hour),
			Duration:    &duration,
		})
	}

	return events, nil
}
func (s *awsServiceIntelligenceService) parseServiceImpactResponse(content, service, region string) (*interfaces.ServiceImpactAnalysis, error) {
	analysis := &interfaces.ServiceImpactAnalysis{
		ServiceName:          service,
		Region:               region,
		CurrentStatus:        "operational",
		ImpactLevel:          "low",
		AffectedClients:      []string{},
		RecommendationImpact: []interfaces.RecommendationImpact{},
		AlternativeServices:  []string{},
		ActionRequired:       false,
		RecommendedActions:   []string{},
	}

	// Parse impact level from content
	contentLower := strings.ToLower(content)
	if strings.Contains(contentLower, "critical") {
		analysis.ImpactLevel = "critical"
		analysis.ActionRequired = true
	} else if strings.Contains(contentLower, "high") {
		analysis.ImpactLevel = "high"
		analysis.ActionRequired = true
	} else if strings.Contains(contentLower, "medium") {
		analysis.ImpactLevel = "medium"
	}

	// Add basic recommendations based on service
	if analysis.ActionRequired {
		analysis.RecommendedActions = []string{
			"Monitor service status closely",
			"Prepare contingency plans",
			"Communicate with affected clients",
		}
	}

	return analysis, nil
}

func (s *awsServiceIntelligenceService) parseNewServicesResponse(content string) ([]*interfaces.NewAWSService, error) {
	// Create sample new services
	// In a real implementation, this would parse structured data from Bedrock
	services := []*interfaces.NewAWSService{
		{
			ServiceName:        "AWS Example Service",
			Category:           "Analytics",
			Description:        "New analytics service for real-time data processing",
			AnnouncementDate:   time.Now().Add(-30 * 24 * time.Hour),
			Regions:            []string{"us-east-1", "us-west-2", "eu-west-1"},
			PricingModel:       "pay-per-use",
			KeyFeatures:        []string{"Real-time processing", "Serverless", "Auto-scaling"},
			UseCases:           []string{"Data analytics", "Real-time dashboards", "IoT processing"},
			CompetitorServices: []string{"Google Cloud Dataflow", "Azure Stream Analytics"},
			DocumentationURL:   "https://docs.aws.amazon.com/example-service/",
			Metadata:           make(map[string]interface{}),
		},
	}

	return services, nil
}

func (s *awsServiceIntelligenceService) parseServiceEvaluationResponse(content, serviceName string) (*interfaces.ServiceEvaluation, error) {
	evaluation := &interfaces.ServiceEvaluation{
		ServiceName:      serviceName,
		RelevanceScore:   75.0, // Default score
		AdoptionPriority: "short-term",
		Benefits:         []string{"Improved performance", "Cost optimization", "Better scalability"},
		Challenges:       []string{"Learning curve", "Integration complexity"},
		Prerequisites:    []string{"Technical training", "Architecture review"},
		CostImplications: &interfaces.CostImplications{
			EstimatedMonthlyCost: 1000.0,
			CostSavings:          500.0,
			ROITimeframe:         "6 months",
			CostFactors:          []string{"Usage volume", "Data transfer", "Storage"},
			OptimizationTips:     []string{"Use reserved capacity", "Optimize data transfer"},
		},
		ImplementationComplexity: "medium",
		TimeToValue:              "3-6 months",
		RiskAssessment: &interfaces.ServiceRiskAssessment{
			OverallRisk:          "medium",
			TechnicalRisks:       []string{"Integration challenges", "Performance impact"},
			BusinessRisks:        []string{"Training costs", "Change management"},
			ComplianceRisks:      []string{"Data governance", "Security review"},
			MitigationStrategies: []string{"Phased rollout", "Comprehensive testing"},
		},
		Recommendation: "Proceed with pilot implementation",
		NextSteps:      []string{"Conduct proof of concept", "Develop implementation plan", "Train team"},
	}

	// Parse relevance score from content
	contentLower := strings.ToLower(content)
	if strings.Contains(contentLower, "highly relevant") || strings.Contains(contentLower, "excellent fit") {
		evaluation.RelevanceScore = 90.0
		evaluation.AdoptionPriority = "immediate"
	} else if strings.Contains(contentLower, "not relevant") || strings.Contains(contentLower, "poor fit") {
		evaluation.RelevanceScore = 25.0
		evaluation.AdoptionPriority = "not-applicable"
	}

	return evaluation, nil
}
func (s *awsServiceIntelligenceService) parseServiceRecommendationsResponse(content string) ([]*interfaces.ServiceRecommendation, error) {
	// Create sample recommendations
	recommendations := []*interfaces.ServiceRecommendation{
		{
			ServiceName:        "Amazon EKS",
			RecommendationType: "new",
			Priority:           "high",
			Rationale:          "Container orchestration needs and scalability requirements",
			ExpectedBenefits:   []string{"Better container management", "Improved scalability", "Reduced operational overhead"},
			ImplementationPlan: &interfaces.ServiceImplementationPlan{
				Phases: []interfaces.ImplementationPhase{
					{
						PhaseNumber:   1,
						PhaseName:     "Planning and Design",
						Description:   "Architecture design and planning",
						Duration:      "2 weeks",
						Deliverables:  []string{"Architecture document", "Implementation plan"},
						Prerequisites: []string{"Team training"},
					},
					{
						PhaseNumber:   2,
						PhaseName:     "Pilot Implementation",
						Description:   "Deploy pilot EKS cluster",
						Duration:      "3 weeks",
						Deliverables:  []string{"Pilot cluster", "Testing results"},
						Prerequisites: []string{"Architecture approval"},
					},
				},
				Prerequisites:        []string{"Kubernetes knowledge", "Container expertise"},
				ResourceRequirements: []string{"DevOps engineer", "Cloud architect"},
				RiskMitigation:       []string{"Comprehensive testing", "Gradual migration"},
				TestingStrategy:      "Automated testing with staging environment",
				RollbackPlan:         "Maintain existing infrastructure during transition",
			},
			CostBenefit: &interfaces.CostBenefitAnalysis{
				InitialCost:     15000.0,
				OngoingCosts:    3000.0,
				ExpectedSavings: 5000.0,
				ROI:             25.0,
				PaybackPeriod:   "8 months",
				NetPresentValue: 45000.0,
			},
			Timeline:     "3-4 months",
			Dependencies: []string{"Team training", "Network setup"},
			Alternatives: []string{"Amazon ECS", "AWS Fargate"},
		},
	}

	return recommendations, nil
}

func (s *awsServiceIntelligenceService) parseDeprecationAlertsResponse(content string) ([]*interfaces.DeprecationAlert, error) {
	// Create sample deprecation alerts
	alerts := []*interfaces.DeprecationAlert{
		{
			ServiceName:       "AWS Example Legacy Service",
			DeprecationType:   "end-of-life",
			AnnouncementDate:  time.Now().Add(-60 * 24 * time.Hour),
			EffectiveDate:     time.Now().Add(365 * 24 * time.Hour),
			GracePeriod:       365 * 24 * time.Hour,
			Severity:          "high",
			ImpactDescription: "Service will be discontinued, requiring migration to alternative solutions",
			RecommendedActions: []string{
				"Assess current usage",
				"Plan migration to alternative service",
				"Update documentation and processes",
			},
			AlternativeServices: []string{"AWS Modern Service", "Third-party alternative"},
			MigrationResources:  []string{"Migration guide", "Support documentation"},
			AffectedRegions:     []string{"all"},
			ClientImpact: []interfaces.ClientImpact{
				{
					ClientID:           "example-client",
					ImpactLevel:        "high",
					AffectedWorkloads:  []string{"data processing", "analytics"},
					ActionRequired:     true,
					RecommendedActions: []string{"Schedule migration", "Test alternatives"},
				},
			},
		},
	}

	return alerts, nil
}
func (s *awsServiceIntelligenceService) parseMigrationPlanResponse(content, deprecatedService string) (*interfaces.MigrationPlan, error) {
	plan := &interfaces.MigrationPlan{
		DeprecatedService: deprecatedService,
		TargetService:     "AWS Modern Alternative",
		MigrationStrategy: "re-architect",
		EstimatedDuration: "3-6 months",
		EstimatedCost:     25000.0,
		RiskLevel:         "medium",
		Prerequisites:     []string{"Team training", "Architecture review", "Testing environment"},
		MigrationSteps: []interfaces.MigrationStep{
			{
				StepNumber:  1,
				Name:        "Assessment and Planning",
				Description: "Assess current usage and plan migration approach",
				Duration:    "2 weeks",
				Resources:   []string{"Migration specialist", "Stakeholder alignment"},
			},
			{
				StepNumber:  2,
				Name:        "Environment Setup",
				Description: "Set up target service environment",
				Duration:    "3 weeks",
				Resources:   []string{"Cloud engineer", "Security specialist"},
			},
			{
				StepNumber:  3,
				Name:        "Data Migration",
				Description: "Migrate data to target service",
				Duration:    "4 weeks",
				Resources:   []string{"Data engineer", "Migration tools"},
			},
		},
		TestingPlan: &interfaces.TestingPlan{
			TestingPhases:       []string{"Unit testing", "Integration testing", "Performance testing"},
			TestEnvironments:    []string{"Development", "Staging", "Pre-production"},
			TestCriteria:        []string{"Functional correctness", "Performance benchmarks"},
			PerformanceTests:    []string{"Load testing", "Stress testing"},
			SecurityTests:       []string{"Vulnerability scanning", "Access control testing"},
			UserAcceptanceTests: []string{"Business workflow validation", "User experience testing"},
		},
		RollbackPlan: &interfaces.RollbackPlan{
			TriggerConditions: []string{"Critical failures", "Performance degradation"},
			RollbackSteps:     []string{"Stop migration", "Restore original service", "Validate restoration"},
			ValidationSteps:   []string{"Data integrity checks", "Service functionality tests"},
			Timeline:          "4-8 hours",
			Resources:         []string{"Migration team", "Backup systems"},
		},
		CostComparison: &interfaces.MigrationCostComparison{
			CurrentServiceCost:  5000.0,
			NewServiceCost:      4000.0,
			MigrationCost:       25000.0,
			TotalCostDifference: -1000.0,
			BreakEvenPoint:      "25 months",
		},
		Timeline: &interfaces.MigrationTimeline{
			StartDate:    time.Now(),
			EndDate:      time.Now().AddDate(0, 4, 0), // 4 months from now
			Milestones:   []string{"Planning complete", "Environment ready", "Data migrated", "Validation complete"},
			Dependencies: []string{"Stakeholder approval", "Resource allocation", "Environment setup"},
			CriticalPath: []string{"Planning", "Environment setup", "Data migration", "Validation"},
		},
		ResourceRequirements: []string{"Cloud architect", "DevOps engineer", "Data engineer", "QA engineer"},
	}

	return plan, nil
}

func (s *awsServiceIntelligenceService) parsePricingChangesResponse(content string) ([]*interfaces.PricingChange, error) {
	// Create sample pricing changes
	changes := []*interfaces.PricingChange{
		{
			ServiceName:       "Amazon EC2",
			Region:            "us-east-1",
			ChangeType:        "increase",
			EffectiveDate:     time.Now().Add(30 * 24 * time.Hour),
			AnnouncementDate:  time.Now().Add(-7 * 24 * time.Hour),
			ChangeDescription: "Price increase for m5.large instances",
			OldPricing: &interfaces.PricingStructure{
				PricingModel:      "on-demand",
				BasePrice:         0.096,
				Unit:              "hour",
				Tiers:             []interfaces.PricingTier{},
				AdditionalCharges: make(map[string]float64),
				Discounts:         []interfaces.PricingDiscount{},
			},
			NewPricing: &interfaces.PricingStructure{
				PricingModel:      "on-demand",
				BasePrice:         0.105,
				Unit:              "hour",
				Tiers:             []interfaces.PricingTier{},
				AdditionalCharges: make(map[string]float64),
				Discounts:         []interfaces.PricingDiscount{},
			},
			ImpactLevel:   "medium",
			AffectedTiers: []string{"m5.large"},
			Metadata:      make(map[string]interface{}),
		},
	}

	return changes, nil
}

func (s *awsServiceIntelligenceService) parseCostImpactResponse(content, serviceName string) (*interfaces.PricingCostImpactAnalysis, error) {
	analysis := &interfaces.PricingCostImpactAnalysis{
		ServiceName:        serviceName,
		PricingChangeID:    "example-change-id",
		CurrentMonthlyCost: 2000.0,
		NewMonthlyCost:     2200.0,
		CostDifference:     200.0,
		PercentageChange:   10.0,
		AnnualImpact:       2400.0,
		ImpactCategory:     "increase",
		Recommendations: []interfaces.CostOptimizationRec{
			{
				Type:                 "reserved-instance",
				Description:          "Purchase reserved instances to offset price increase",
				EstimatedSavings:     300.0,
				ImplementationEffort: "low",
				RiskLevel:            "low",
				Prerequisites:        []string{"Commitment to 1-year term"},
			},
		},
		AlternativeOptions: []interfaces.AlternativeOption{
			{
				ServiceName:       "Alternative Instance Type",
				Description:       "Switch to more cost-effective instance type",
				CostComparison:    -150.0,
				FeatureComparison: "Similar performance with lower cost",
				MigrationEffort:   "medium",
				Pros:              []string{"Lower cost", "Similar performance"},
				Cons:              []string{"Migration effort", "Testing required"},
			},
		},
		ActionRequired: true,
		Timeline:       "30 days before effective date",
	}

	return analysis, nil
}
