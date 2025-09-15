package interfaces

import (
	"context"
	"time"
)

// AWSServiceIntelligenceService defines the interface for AWS service intelligence
type AWSServiceIntelligenceService interface {
	// Service Status Monitoring
	GetServiceStatus(ctx context.Context, region string) (*AWSServiceStatus, error)
	GetServiceStatusHistory(ctx context.Context, service, region string, duration time.Duration) ([]*ServiceStatusEvent, error)
	AnalyzeServiceImpact(ctx context.Context, service, region string) (*ServiceImpactAnalysis, error)

	// New Service Intelligence
	GetNewServices(ctx context.Context, since time.Time) ([]*NewAWSService, error)
	EvaluateServiceForClient(ctx context.Context, service *NewAWSService, clientContext *ClientContext) (*ServiceEvaluation, error)
	GetServiceRecommendations(ctx context.Context, clientContext *ClientContext) ([]*ServiceRecommendation, error)

	// Deprecation and Migration
	GetDeprecationAlerts(ctx context.Context) ([]*DeprecationAlert, error)
	GenerateMigrationPlan(ctx context.Context, deprecatedService string, clientContext *ClientContext) (*MigrationPlan, error)

	// Cost Impact Analysis
	AnalyzePricingChanges(ctx context.Context, since time.Time) ([]*PricingChange, error)
	CalculateCostImpact(ctx context.Context, pricingChange *PricingChange, clientContext *ClientContext) (*PricingCostImpactAnalysis, error)

	// Intelligence Updates
	RefreshServiceIntelligence(ctx context.Context) error
	GetLastUpdateTime() time.Time
	IsHealthy() bool
}

// AWSServiceStatus represents the current status of AWS services
type AWSServiceStatus struct {
	Region        string                    `json:"region"`
	LastUpdated   time.Time                 `json:"last_updated"`
	Services      map[string]*ServiceHealth `json:"services"`
	OverallHealth string                    `json:"overall_health"` // "healthy", "degraded", "unhealthy"
}

// ServiceHealth represents the health status of a specific AWS service
type ServiceHealth struct {
	ServiceName  string                 `json:"service_name"`
	Status       string                 `json:"status"` // "operational", "degraded", "disrupted"
	LastIncident *ServiceIncident       `json:"last_incident,omitempty"`
	Metrics      *ServiceMetrics        `json:"metrics,omitempty"`
	Regions      map[string]string      `json:"regions"` // region -> status
	Metadata     map[string]interface{} `json:"metadata"`
}

// ServiceIncident represents a service incident
type ServiceIncident struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Impact      string     `json:"impact"` // "low", "medium", "high", "critical"
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Updates     []string   `json:"updates"`
}

// ServiceMetrics represents performance metrics for a service
type ServiceMetrics struct {
	Availability  float64   `json:"availability"`   // percentage
	ResponseTime  float64   `json:"response_time"`  // milliseconds
	ErrorRate     float64   `json:"error_rate"`     // percentage
	ThroughputRPS float64   `json:"throughput_rps"` // requests per second
	LastMeasured  time.Time `json:"last_measured"`
}

// ServiceStatusEvent represents a historical status event
type ServiceStatusEvent struct {
	ServiceName string         `json:"service_name"`
	Region      string         `json:"region"`
	Status      string         `json:"status"`
	Message     string         `json:"message"`
	Timestamp   time.Time      `json:"timestamp"`
	Duration    *time.Duration `json:"duration,omitempty"`
}

// ServiceImpactAnalysis represents analysis of service impact on client recommendations
type ServiceImpactAnalysis struct {
	ServiceName          string                 `json:"service_name"`
	Region               string                 `json:"region"`
	CurrentStatus        string                 `json:"current_status"`
	ImpactLevel          string                 `json:"impact_level"` // "none", "low", "medium", "high", "critical"
	AffectedClients      []string               `json:"affected_clients"`
	RecommendationImpact []RecommendationImpact `json:"recommendation_impact"`
	AlternativeServices  []string               `json:"alternative_services"`
	ActionRequired       bool                   `json:"action_required"`
	RecommendedActions   []string               `json:"recommended_actions"`
}

// RecommendationImpact represents how a service issue impacts specific recommendations
type RecommendationImpact struct {
	RecommendationType string   `json:"recommendation_type"`
	ImpactDescription  string   `json:"impact_description"`
	Severity           string   `json:"severity"`
	Alternatives       []string `json:"alternatives"`
}

// NewAWSService represents a newly announced AWS service
type NewAWSService struct {
	ServiceName        string                 `json:"service_name"`
	Category           string                 `json:"category"`
	Description        string                 `json:"description"`
	AnnouncementDate   time.Time              `json:"announcement_date"`
	GADate             *time.Time             `json:"ga_date,omitempty"`
	PreviewDate        *time.Time             `json:"preview_date,omitempty"`
	Regions            []string               `json:"regions"`
	PricingModel       string                 `json:"pricing_model"`
	KeyFeatures        []string               `json:"key_features"`
	UseCases           []string               `json:"use_cases"`
	CompetitorServices []string               `json:"competitor_services"`
	DocumentationURL   string                 `json:"documentation_url"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// ClientContext represents client-specific context for service evaluation
type ClientContext struct {
	IndustryVertical       string                 `json:"industry_vertical"`
	CompanySize            string                 `json:"company_size"`
	CurrentServices        []string               `json:"current_services"`
	PreferredRegions       []string               `json:"preferred_regions"`
	ComplianceRequirements []string               `json:"compliance_requirements"`
	BudgetConstraints      *BudgetConstraints     `json:"budget_constraints,omitempty"`
	TechnicalMaturity      string                 `json:"technical_maturity"` // "beginner", "intermediate", "advanced"
	Workloads              []WorkloadProfile      `json:"workloads"`
	Metadata               map[string]interface{} `json:"metadata"`
}

// Note: BudgetConstraints and WorkloadProfile are imported from existing interfaces

// ServiceEvaluation represents evaluation of a new service for a specific client
type ServiceEvaluation struct {
	ServiceName              string                 `json:"service_name"`
	ClientID                 string                 `json:"client_id,omitempty"`
	RelevanceScore           float64                `json:"relevance_score"`   // 0-100
	AdoptionPriority         string                 `json:"adoption_priority"` // "immediate", "short-term", "long-term", "not-applicable"
	Benefits                 []string               `json:"benefits"`
	Challenges               []string               `json:"challenges"`
	Prerequisites            []string               `json:"prerequisites"`
	CostImplications         *CostImplications      `json:"cost_implications"`
	ImplementationComplexity string                 `json:"implementation_complexity"` // "low", "medium", "high"
	TimeToValue              string                 `json:"time_to_value"`
	RiskAssessment           *ServiceRiskAssessment `json:"risk_assessment"`
	Recommendation           string                 `json:"recommendation"`
	NextSteps                []string               `json:"next_steps"`
}

// CostImplications represents cost implications of adopting a service
type CostImplications struct {
	EstimatedMonthlyCost float64  `json:"estimated_monthly_cost"`
	CostSavings          float64  `json:"cost_savings"`
	ROITimeframe         string   `json:"roi_timeframe"`
	CostFactors          []string `json:"cost_factors"`
	OptimizationTips     []string `json:"optimization_tips"`
}

// ServiceRiskAssessment represents risk assessment for service adoption
type ServiceRiskAssessment struct {
	OverallRisk          string   `json:"overall_risk"` // "low", "medium", "high"
	TechnicalRisks       []string `json:"technical_risks"`
	BusinessRisks        []string `json:"business_risks"`
	ComplianceRisks      []string `json:"compliance_risks"`
	MitigationStrategies []string `json:"mitigation_strategies"`
}

// ServiceRecommendation represents a service recommendation for a client
type ServiceRecommendation struct {
	ServiceName        string                     `json:"service_name"`
	RecommendationType string                     `json:"recommendation_type"` // "new", "upgrade", "replacement", "optimization"
	Priority           string                     `json:"priority"`            // "critical", "high", "medium", "low"
	Rationale          string                     `json:"rationale"`
	ExpectedBenefits   []string                   `json:"expected_benefits"`
	ImplementationPlan *ServiceImplementationPlan `json:"implementation_plan"`
	CostBenefit        *CostBenefitAnalysis       `json:"cost_benefit"`
	Timeline           string                     `json:"timeline"`
	Dependencies       []string                   `json:"dependencies"`
	Alternatives       []string                   `json:"alternatives"`
}

// ServiceImplementationPlan represents implementation plan for a service recommendation
type ServiceImplementationPlan struct {
	Phases               []ImplementationPhase `json:"phases"`
	Prerequisites        []string              `json:"prerequisites"`
	ResourceRequirements []string              `json:"resource_requirements"`
	RiskMitigation       []string              `json:"risk_mitigation"`
	TestingStrategy      string                `json:"testing_strategy"`
	RollbackPlan         string                `json:"rollback_plan"`
}

// CostBenefitAnalysis represents cost-benefit analysis for a recommendation
type CostBenefitAnalysis struct {
	InitialCost     float64 `json:"initial_cost"`
	OngoingCosts    float64 `json:"ongoing_costs"`
	ExpectedSavings float64 `json:"expected_savings"`
	ROI             float64 `json:"roi"`
	PaybackPeriod   string  `json:"payback_period"`
	NetPresentValue float64 `json:"net_present_value"`
}

// DeprecationAlert represents a service deprecation alert
type DeprecationAlert struct {
	ServiceName         string         `json:"service_name"`
	DeprecationType     string         `json:"deprecation_type"` // "end-of-life", "feature-removal", "version-deprecation"
	AnnouncementDate    time.Time      `json:"announcement_date"`
	EffectiveDate       time.Time      `json:"effective_date"`
	GracePeriod         time.Duration  `json:"grace_period"`
	Severity            string         `json:"severity"` // "critical", "high", "medium", "low"
	ImpactDescription   string         `json:"impact_description"`
	RecommendedActions  []string       `json:"recommended_actions"`
	AlternativeServices []string       `json:"alternative_services"`
	MigrationResources  []string       `json:"migration_resources"`
	AffectedRegions     []string       `json:"affected_regions"`
	ClientImpact        []ClientImpact `json:"client_impact"`
}

// ClientImpact represents impact on specific clients
type ClientImpact struct {
	ClientID           string   `json:"client_id"`
	ImpactLevel        string   `json:"impact_level"` // "none", "low", "medium", "high", "critical"
	AffectedWorkloads  []string `json:"affected_workloads"`
	ActionRequired     bool     `json:"action_required"`
	RecommendedActions []string `json:"recommended_actions"`
}

// MigrationPlan represents a migration plan for deprecated services
type MigrationPlan struct {
	DeprecatedService    string                   `json:"deprecated_service"`
	TargetService        string                   `json:"target_service"`
	MigrationStrategy    string                   `json:"migration_strategy"` // "lift-and-shift", "re-architect", "hybrid"
	EstimatedDuration    string                   `json:"estimated_duration"`
	EstimatedCost        float64                  `json:"estimated_cost"`
	RiskLevel            string                   `json:"risk_level"`
	Prerequisites        []string                 `json:"prerequisites"`
	MigrationSteps       []MigrationStep          `json:"migration_steps"`
	TestingPlan          *TestingPlan             `json:"testing_plan"`
	RollbackPlan         *RollbackPlan            `json:"rollback_plan"`
	CostComparison       *MigrationCostComparison `json:"cost_comparison"`
	Timeline             *MigrationTimeline       `json:"timeline"`
	ResourceRequirements []string                 `json:"resource_requirements"`
}

// Note: MigrationStep is imported from existing interfaces

// TestingPlan represents testing plan for migration
type TestingPlan struct {
	TestingPhases       []string `json:"testing_phases"`
	TestEnvironments    []string `json:"test_environments"`
	TestCriteria        []string `json:"test_criteria"`
	PerformanceTests    []string `json:"performance_tests"`
	SecurityTests       []string `json:"security_tests"`
	UserAcceptanceTests []string `json:"user_acceptance_tests"`
}

// Note: RollbackPlan is imported from existing interfaces

// MigrationCostComparison represents cost comparison for migration
type MigrationCostComparison struct {
	CurrentServiceCost  float64 `json:"current_service_cost"`
	NewServiceCost      float64 `json:"new_service_cost"`
	MigrationCost       float64 `json:"migration_cost"`
	TotalCostDifference float64 `json:"total_cost_difference"`
	BreakEvenPoint      string  `json:"break_even_point"`
}

// Note: MigrationTimeline is imported from existing interfaces

// PricingChange represents a change in AWS service pricing
type PricingChange struct {
	ServiceName       string                 `json:"service_name"`
	Region            string                 `json:"region"`
	ChangeType        string                 `json:"change_type"` // "increase", "decrease", "new-tier", "restructure"
	EffectiveDate     time.Time              `json:"effective_date"`
	AnnouncementDate  time.Time              `json:"announcement_date"`
	ChangeDescription string                 `json:"change_description"`
	OldPricing        *PricingStructure      `json:"old_pricing"`
	NewPricing        *PricingStructure      `json:"new_pricing"`
	ImpactLevel       string                 `json:"impact_level"` // "low", "medium", "high", "critical"
	AffectedTiers     []string               `json:"affected_tiers"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// PricingStructure represents pricing structure for a service
type PricingStructure struct {
	PricingModel      string             `json:"pricing_model"` // "on-demand", "reserved", "spot", "savings-plan"
	BasePrice         float64            `json:"base_price"`
	Unit              string             `json:"unit"` // "hour", "request", "GB", etc.
	Tiers             []PricingTier      `json:"tiers"`
	AdditionalCharges map[string]float64 `json:"additional_charges"`
	Discounts         []PricingDiscount  `json:"discounts"`
}

// PricingTier represents a pricing tier
type PricingTier struct {
	Name         string  `json:"name"`
	MinUsage     float64 `json:"min_usage"`
	MaxUsage     float64 `json:"max_usage"`
	PricePerUnit float64 `json:"price_per_unit"`
}

// PricingDiscount represents a pricing discount
type PricingDiscount struct {
	Name            string   `json:"name"`
	DiscountPercent float64  `json:"discount_percent"`
	Conditions      []string `json:"conditions"`
}

// PricingCostImpactAnalysis represents analysis of cost impact from pricing changes
type PricingCostImpactAnalysis struct {
	ServiceName        string                `json:"service_name"`
	ClientID           string                `json:"client_id,omitempty"`
	PricingChangeID    string                `json:"pricing_change_id"`
	CurrentMonthlyCost float64               `json:"current_monthly_cost"`
	NewMonthlyCost     float64               `json:"new_monthly_cost"`
	CostDifference     float64               `json:"cost_difference"`
	PercentageChange   float64               `json:"percentage_change"`
	AnnualImpact       float64               `json:"annual_impact"`
	ImpactCategory     string                `json:"impact_category"` // "savings", "increase", "neutral"
	Recommendations    []CostOptimizationRec `json:"recommendations"`
	AlternativeOptions []AlternativeOption   `json:"alternative_options"`
	ActionRequired     bool                  `json:"action_required"`
	Timeline           string                `json:"timeline"`
}

// CostOptimizationRec represents a cost optimization recommendation
type CostOptimizationRec struct {
	Type                 string   `json:"type"` // "reserved-instance", "savings-plan", "right-sizing", "alternative-service"
	Description          string   `json:"description"`
	EstimatedSavings     float64  `json:"estimated_savings"`
	ImplementationEffort string   `json:"implementation_effort"` // "low", "medium", "high"
	RiskLevel            string   `json:"risk_level"`
	Prerequisites        []string `json:"prerequisites"`
}

// AlternativeOption represents an alternative service option
type AlternativeOption struct {
	ServiceName       string   `json:"service_name"`
	Description       string   `json:"description"`
	CostComparison    float64  `json:"cost_comparison"`
	FeatureComparison string   `json:"feature_comparison"`
	MigrationEffort   string   `json:"migration_effort"`
	Pros              []string `json:"pros"`
	Cons              []string `json:"cons"`
}
