package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// MultiCloudAnalyzer provides comparative analysis across cloud providers
type MultiCloudAnalyzer interface {
	CompareServices(ctx context.Context, requirement ServiceRequirement) (*ServiceComparison, error)
	AnalyzeCosts(ctx context.Context, workload WorkloadSpec) (*CostAnalysis, error)
	EvaluateProviders(ctx context.Context, criteria EvaluationCriteria) (*ProviderEvaluation, error)
	GetMigrationPaths(ctx context.Context, source, target CloudProviderInfo) (*MigrationPath, error)
	GetProviderRecommendation(ctx context.Context, inquiry *domain.Inquiry) (*ProviderRecommendation, error)
	GetServiceEquivalents(ctx context.Context, provider string, serviceName string) (map[string]string, error)
}

// ServiceRequirement defines requirements for a cloud service
type ServiceRequirement struct {
	Category     string                   `json:"category"` // "compute", "storage", "database", etc.
	Requirements []string                 `json:"requirements"`
	Performance  PerformanceRequirements  `json:"performance"`
	Compliance   []string                 `json:"compliance"`
	Budget       BudgetConstraints        `json:"budget"`
	Industry     string                   `json:"industry,omitempty"`
	Region       string                   `json:"region,omitempty"`
}

// PerformanceRequirements defines performance criteria
type PerformanceRequirements struct {
	CPU              string `json:"cpu,omitempty"`              // e.g., "4 vCPUs"
	Memory           string `json:"memory,omitempty"`           // e.g., "16 GB"
	Storage          string `json:"storage,omitempty"`          // e.g., "1 TB SSD"
	NetworkBandwidth string `json:"network_bandwidth,omitempty"` // e.g., "10 Gbps"
	IOPS             int    `json:"iops,omitempty"`             // Input/Output Operations Per Second
	Latency          string `json:"latency,omitempty"`          // e.g., "< 10ms"
	Availability     string `json:"availability,omitempty"`     // e.g., "99.99%"
	Scalability      string `json:"scalability,omitempty"`      // e.g., "auto-scaling"
}

// BudgetConstraints defines budget limitations
type BudgetConstraints struct {
	MaxMonthlyCost   float64 `json:"max_monthly_cost,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	CostOptimization bool    `json:"cost_optimization"`
	PaymentModel     string  `json:"payment_model,omitempty"` // "pay-as-you-go", "reserved", "spot"
}

// ServiceComparison represents a comparison of services across providers
type ServiceComparison struct {
	Category        string           `json:"category"`
	Providers       []ProviderOption `json:"providers"`
	Recommendation  string           `json:"recommendation"`
	Reasoning       string           `json:"reasoning"`
	ComparisonMatrix []ComparisonRow  `json:"comparison_matrix"`
	CreatedAt       time.Time        `json:"created_at"`
}

// ProviderOption represents a cloud provider option for a service
type ProviderOption struct {
	Provider                 string            `json:"provider"`
	ServiceName              string            `json:"service_name"`
	ServiceDescription       string            `json:"service_description"`
	Pros                     []string          `json:"pros"`
	Cons                     []string          `json:"cons"`
	EstimatedMonthlyCost     string            `json:"estimated_monthly_cost"`
	DocumentationURL         string            `json:"documentation_url"`
	PricingURL               string            `json:"pricing_url"`
	ImplementationComplexity string            `json:"implementation_complexity"` // "low", "medium", "high"
	SupportLevel             string            `json:"support_level"`
	RegionalAvailability     []string          `json:"regional_availability"`
	ComplianceCertifications []string          `json:"compliance_certifications"`
	Features                 map[string]string `json:"features"`
	Score                    float64           `json:"score"` // Overall score based on requirements
}

// ComparisonRow represents a single comparison criteria
type ComparisonRow struct {
	Criteria    string             `json:"criteria"`
	Weight      float64            `json:"weight"`
	Scores      map[string]float64 `json:"scores"` // provider -> score
	Notes       map[string]string  `json:"notes"`  // provider -> notes
	Description string             `json:"description"`
}

// WorkloadSpec defines a workload for cost analysis
type WorkloadSpec struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Components   []WorkloadComponent    `json:"components"`
	Usage        UsagePattern           `json:"usage"`
	Region       string                 `json:"region"`
	Duration     string                 `json:"duration"` // e.g., "1 year", "3 years"
	Environment  string                 `json:"environment"` // "dev", "staging", "production"
	Metadata     map[string]interface{} `json:"metadata"`
}

// WorkloadComponent represents a component of a workload
type WorkloadComponent struct {
	Type         string            `json:"type"` // "compute", "storage", "database", "network"
	Name         string            `json:"name"`
	Requirements map[string]string `json:"requirements"`
	Quantity     int               `json:"quantity"`
	Utilization  float64           `json:"utilization"` // 0.0 to 1.0
}

// UsagePattern defines usage patterns for cost estimation
type UsagePattern struct {
	PeakHours        []int   `json:"peak_hours"`        // Hours of day (0-23)
	PeakDays         []int   `json:"peak_days"`         // Days of week (0-6, Sunday=0)
	SeasonalFactors  []float64 `json:"seasonal_factors"` // Monthly multipliers
	GrowthRate       float64 `json:"growth_rate"`       // Annual growth rate
	BaselineLoad     float64 `json:"baseline_load"`     // Baseline utilization (0.0-1.0)
	BurstCapacity    float64 `json:"burst_capacity"`    // Maximum burst multiplier
}

// CostAnalysis represents cost analysis results
type CostAnalysis struct {
	WorkloadName     string                    `json:"workload_name"`
	AnalysisDate     time.Time                 `json:"analysis_date"`
	Currency         string                    `json:"currency"`
	ProviderCosts    []ProviderCostBreakdown   `json:"provider_costs"`
	Recommendation   string                    `json:"recommendation"`
	CostOptimization []CostOptimizationTip     `json:"cost_optimization"`
	Assumptions      []string                  `json:"assumptions"`
	Disclaimers      []string                  `json:"disclaimers"`
	ValidUntil       time.Time                 `json:"valid_until"`
}

// ProviderCostBreakdown represents cost breakdown for a provider
type ProviderCostBreakdown struct {
	Provider        string                 `json:"provider"`
	TotalMonthlyCost float64               `json:"total_monthly_cost"`
	TotalAnnualCost  float64               `json:"total_annual_cost"`
	ComponentCosts   []ComponentCost       `json:"component_costs"`
	Discounts        []CostDiscount        `json:"discounts"`
	AdditionalFees   []AdditionalFee       `json:"additional_fees"`
	PricingModel     string                `json:"pricing_model"`
	LastUpdated      time.Time             `json:"last_updated"`
}

// ComponentCost represents cost for a specific component
type ComponentCost struct {
	ComponentName   string  `json:"component_name"`
	ServiceName     string  `json:"service_name"`
	MonthlyCost     float64 `json:"monthly_cost"`
	Unit            string  `json:"unit"`
	Quantity        float64 `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	PricingTier     string  `json:"pricing_tier,omitempty"`
}

// CostDiscount represents available discounts
type CostDiscount struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "volume", "commitment", "promotional"
	Amount      float64 `json:"amount"`
	Percentage  float64 `json:"percentage"`
	Description string  `json:"description"`
	Requirements []string `json:"requirements"`
}

// AdditionalFee represents additional fees
type AdditionalFee struct {
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // "one-time", "monthly", "per-transaction"
	Description string  `json:"description"`
}

// CostOptimizationTip represents a cost optimization recommendation
type CostOptimizationTip struct {
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	PotentialSavings float64 `json:"potential_savings"`
	Effort          string  `json:"effort"` // "low", "medium", "high"
	Impact          string  `json:"impact"` // "low", "medium", "high"
	Implementation  []string `json:"implementation"`
	Risks           []string `json:"risks,omitempty"`
}

// EvaluationCriteria defines criteria for provider evaluation
type EvaluationCriteria struct {
	UseCase          string                 `json:"use_case"`
	Industry         string                 `json:"industry"`
	ComplianceNeeds  []string               `json:"compliance_needs"`
	TechnicalNeeds   []string               `json:"technical_needs"`
	BusinessNeeds    []string               `json:"business_needs"`
	Budget           BudgetConstraints      `json:"budget"`
	Timeline         string                 `json:"timeline"`
	RiskTolerance    string                 `json:"risk_tolerance"` // "low", "medium", "high"
	Priorities       map[string]float64     `json:"priorities"`     // criteria -> weight
	ExistingProvider string                 `json:"existing_provider,omitempty"`
	TeamSkills       []string               `json:"team_skills"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// ProviderEvaluation represents evaluation results for providers
type ProviderEvaluation struct {
	UseCase         string                  `json:"use_case"`
	EvaluationDate  time.Time               `json:"evaluation_date"`
	Providers       []ProviderScore         `json:"providers"`
	Recommendation  ProviderRecommendation  `json:"recommendation"`
	ComparisonMatrix []ComparisonRow        `json:"comparison_matrix"`
	Summary         string                  `json:"summary"`
	NextSteps       []string                `json:"next_steps"`
}

// ProviderScore represents a provider's evaluation score
type ProviderScore struct {
	Provider        string             `json:"provider"`
	OverallScore    float64            `json:"overall_score"`
	CategoryScores  map[string]float64 `json:"category_scores"`
	Strengths       []string           `json:"strengths"`
	Weaknesses      []string           `json:"weaknesses"`
	Recommendations []string           `json:"recommendations"`
	FitScore        string             `json:"fit_score"` // "excellent", "good", "fair", "poor"
}

// CloudProviderInfo represents detailed information about a cloud provider
type CloudProviderInfo struct {
	Name            string   `json:"name"`
	Code            string   `json:"code"` // "aws", "azure", "gcp", etc.
	DisplayName     string   `json:"display_name"`
	Regions         []string `json:"regions"`
	Strengths       []string `json:"strengths"`
	Specializations []string `json:"specializations"`
	MarketShare     float64  `json:"market_share"`
	Founded         int      `json:"founded"`
	Headquarters    string   `json:"headquarters"`
}

// MigrationPath represents a migration path between providers
type MigrationPath struct {
	SourceProvider      CloudProviderInfo   `json:"source_provider"`
	TargetProvider      CloudProviderInfo   `json:"target_provider"`
	MigrationStrategy   string              `json:"migration_strategy"`
	EstimatedDuration   string              `json:"estimated_duration"`
	EstimatedCost       string              `json:"estimated_cost"`
	ComplexityLevel     string              `json:"complexity_level"` // "low", "medium", "high"
	ServiceMappings     []ServiceMapping    `json:"service_mappings"`
	MigrationPhases     []MigrationPhase    `json:"migration_phases"`
	RisksAndChallenges  []string            `json:"risks_and_challenges"`
	BestPractices       []string            `json:"best_practices"`
	ToolsAndServices    []MigrationTool     `json:"tools_and_services"`
	SuccessFactors      []string            `json:"success_factors"`
}

// ServiceMapping represents mapping between services of different providers
type ServiceMapping struct {
	SourceService   string   `json:"source_service"`
	TargetService   string   `json:"target_service"`
	MappingType     string   `json:"mapping_type"` // "direct", "partial", "alternative", "custom"
	Compatibility   string   `json:"compatibility"` // "high", "medium", "low"
	MigrationNotes  []string `json:"migration_notes"`
	DataTransfer    string   `json:"data_transfer,omitempty"`
	ConfigChanges   []string `json:"config_changes,omitempty"`
}

// MigrationPhase represents a phase in the migration process
type MigrationPhase struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Duration        string   `json:"duration"`
	Prerequisites   []string `json:"prerequisites"`
	Tasks           []string `json:"tasks"`
	Deliverables    []string `json:"deliverables"`
	RiskLevel       string   `json:"risk_level"`
	Dependencies    []string `json:"dependencies"`
	SuccessCriteria []string `json:"success_criteria"`
}

// MigrationTool represents a tool or service for migration
type MigrationTool struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"` // "native", "third-party", "open-source"
	Provider     string   `json:"provider"`
	Description  string   `json:"description"`
	SupportedServices []string `json:"supported_services"`
	Cost         string   `json:"cost"`
	Documentation string  `json:"documentation"`
	Limitations  []string `json:"limitations,omitempty"`
}

// ProviderRecommendation represents a recommendation for a specific provider
type ProviderRecommendation struct {
	RecommendedProvider string            `json:"recommended_provider"`
	Confidence          string            `json:"confidence"` // "high", "medium", "low"
	Reasoning           []string          `json:"reasoning"`
	AlternativeOptions  []string          `json:"alternative_options"`
	Scenarios           []RecommendationScenario `json:"scenarios"`
	Implementation      ImplementationGuidance   `json:"implementation"`
	CostImplications    string            `json:"cost_implications"`
	RiskAssessment      string            `json:"risk_assessment"`
	Timeline            string            `json:"timeline"`
}

// RecommendationScenario represents different scenarios for recommendations
type RecommendationScenario struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Provider       string   `json:"provider"`
	Conditions     []string `json:"conditions"`
	Benefits       []string `json:"benefits"`
	Considerations []string `json:"considerations"`
}

// ImplementationGuidance provides guidance for implementation
type ImplementationGuidance struct {
	QuickStart      []string `json:"quick_start"`
	BestPractices   []string `json:"best_practices"`
	CommonPitfalls  []string `json:"common_pitfalls"`
	ResourceLinks   []string `json:"resource_links"`
	SupportOptions  []string `json:"support_options"`
	TrainingNeeds   []string `json:"training_needs"`
}