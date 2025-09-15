package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// CostAnalysisEngine defines the interface for advanced cost analysis and optimization
type CostAnalysisEngine interface {
	// Cost breakdown analysis
	AnalyzeCostBreakdown(ctx context.Context, inquiry *domain.Inquiry, architecture *CostArchitectureSpec) (*CostBreakdownAnalysis, error)
	GenerateCostOptimizationRecommendations(ctx context.Context, costBreakdown *CostBreakdownAnalysis) (*CostOptimizationRecommendations, error)

	// Reserved instance and savings plan optimization
	AnalyzeReservedInstanceOpportunities(ctx context.Context, usageData *UsageData) (*ReservedInstanceAnalysis, error)
	AnalyzeSavingsPlansOpportunities(ctx context.Context, usageData *UsageData) (*SavingsPlansAnalysis, error)
	CalculateExactSavings(ctx context.Context, recommendations *PurchaseRecommendations) (*SavingsCalculation, error)

	// Right-sizing recommendations
	AnalyzeResourceUtilization(ctx context.Context, resources *ResourceUtilizationData) (*RightSizingAnalysis, error)
	GenerateRightSizingRecommendations(ctx context.Context, utilizationAnalysis *RightSizingAnalysis) (*RightSizingRecommendations, error)

	// Cost forecasting
	GenerateCostForecast(ctx context.Context, architecture *CostArchitectureSpec, forecastParams *ForecastParameters) (*CostForecast, error)
	CalculateConfidenceIntervals(ctx context.Context, forecast *CostForecast) (*ForecastConfidenceIntervals, error)

	// Comprehensive cost analysis
	GenerateComprehensiveCostAnalysis(ctx context.Context, inquiry *domain.Inquiry) (*ComprehensiveCostAnalysis, error)
}

// CostBreakdownAnalysis represents detailed cost breakdown analysis
type CostBreakdownAnalysis struct {
	ID                    string                     `json:"id"`
	InquiryID             string                     `json:"inquiry_id"`
	AnalysisDate          time.Time                  `json:"analysis_date"`
	TotalMonthlyCost      float64                    `json:"total_monthly_cost"`
	TotalAnnualCost       float64                    `json:"total_annual_cost"`
	Currency              string                     `json:"currency"`
	ServiceBreakdown      []*ServiceCostBreakdown    `json:"service_breakdown"`
	CategoryBreakdown     []*CategoryCostBreakdown   `json:"category_breakdown"`
	RegionBreakdown       []*RegionCostBreakdown     `json:"region_breakdown"`
	CostDrivers           []*CostDriver              `json:"cost_drivers"`
	CostTrends            *CostAnalysisTrends        `json:"cost_trends"`
	BenchmarkComparison   *CostBenchmarkComparison   `json:"benchmark_comparison"`
	OptimizationPotential *CostOptimizationPotential `json:"optimization_potential"`
	Assumptions           []string                   `json:"assumptions"`
	Methodology           string                     `json:"methodology"`
	CreatedAt             time.Time                  `json:"created_at"`
}

// ServiceCostBreakdown represents cost breakdown by service
type ServiceCostBreakdown struct {
	ServiceName     string                `json:"service_name"`
	Provider        string                `json:"provider"`
	Category        string                `json:"category"`
	MonthlyCost     float64               `json:"monthly_cost"`
	AnnualCost      float64               `json:"annual_cost"`
	CostPercentage  float64               `json:"cost_percentage"`
	UsageMetrics    *ServiceUsageMetrics  `json:"usage_metrics"`
	PricingModel    string                `json:"pricing_model"`
	CostComponents  []*CostComponent      `json:"cost_components"`
	OptimizationOps []*OptimizationOption `json:"optimization_opportunities"`
}

// CategoryCostBreakdown represents cost breakdown by category
type CategoryCostBreakdown struct {
	Category       string  `json:"category"`
	MonthlyCost    float64 `json:"monthly_cost"`
	AnnualCost     float64 `json:"annual_cost"`
	CostPercentage float64 `json:"cost_percentage"`
	ServiceCount   int     `json:"service_count"`
	GrowthRate     float64 `json:"growth_rate"`
}

// RegionCostBreakdown represents cost breakdown by region
type RegionCostBreakdown struct {
	Region         string  `json:"region"`
	MonthlyCost    float64 `json:"monthly_cost"`
	AnnualCost     float64 `json:"annual_cost"`
	CostPercentage float64 `json:"cost_percentage"`
	ServiceCount   int     `json:"service_count"`
	DataTransfer   float64 `json:"data_transfer_cost"`
}

// CostDriver represents factors driving costs
type CostDriver struct {
	DriverName       string   `json:"driver_name"`
	Impact           string   `json:"impact"` // "high", "medium", "low"
	CostContribution float64  `json:"cost_contribution"`
	Description      string   `json:"description"`
	Recommendations  []string `json:"recommendations"`
}

// CostAnalysisTrends represents historical and projected cost trends
type CostAnalysisTrends struct {
	HistoricalData   []*CostDataPoint   `json:"historical_data"`
	ProjectedGrowth  float64            `json:"projected_growth"`
	SeasonalPatterns []*SeasonalPattern `json:"seasonal_patterns"`
	TrendAnalysis    string             `json:"trend_analysis"`
}

// CostDataPoint represents a cost data point in time
type CostDataPoint struct {
	Date   time.Time `json:"date"`
	Cost   float64   `json:"cost"`
	Period string    `json:"period"` // "monthly", "weekly", "daily"
}

// SeasonalPattern represents seasonal cost patterns
type SeasonalPattern struct {
	Period      string  `json:"period"` // "monthly", "quarterly"
	Pattern     string  `json:"pattern"`
	Variance    float64 `json:"variance"`
	Description string  `json:"description"`
}

// CostBenchmarkComparison represents cost comparison with industry benchmarks
type CostBenchmarkComparison struct {
	Industry           string             `json:"industry"`
	CompanySize        string             `json:"company_size"`
	BenchmarkMetrics   []*BenchmarkMetric `json:"benchmark_metrics"`
	PerformanceRating  string             `json:"performance_rating"`
	ComparisonInsights []string           `json:"comparison_insights"`
}

// BenchmarkMetric represents a benchmark metric
type BenchmarkMetric struct {
	MetricName      string  `json:"metric_name"`
	YourValue       float64 `json:"your_value"`
	IndustryAverage float64 `json:"industry_average"`
	IndustryMedian  float64 `json:"industry_median"`
	Percentile      int     `json:"percentile"`
	Status          string  `json:"status"` // "above", "below", "at" benchmark
}

// CostOptimizationPotential represents potential cost optimization
type CostOptimizationPotential struct {
	TotalSavingsPotential    float64                 `json:"total_savings_potential"`
	QuickWinsSavings         float64                 `json:"quick_wins_savings"`
	LongTermSavings          float64                 `json:"long_term_savings"`
	OptimizationCategories   []*OptimizationCategory `json:"optimization_categories"`
	ImplementationComplexity string                  `json:"implementation_complexity"`
	TimeToRealizeSavings     string                  `json:"time_to_realize_savings"`
}

// OptimizationCategory represents a category of optimization
type OptimizationCategory struct {
	Category         string  `json:"category"`
	SavingsPotential float64 `json:"savings_potential"`
	Effort           string  `json:"effort"`
	Timeline         string  `json:"timeline"`
	RiskLevel        string  `json:"risk_level"`
}

// CostOptimizationRecommendations represents specific optimization recommendations
type CostOptimizationRecommendations struct {
	ID                    string                            `json:"id"`
	AnalysisID            string                            `json:"analysis_id"`
	TotalSavingsPotential float64                           `json:"total_savings_potential"`
	Recommendations       []*CostOptimizationRecommendation `json:"recommendations"`
	ImplementationPlan    *OptimizationImplementationPlan   `json:"implementation_plan"`
	RiskAssessment        *CostOptimizationRiskAssessment   `json:"risk_assessment"`
	ROIAnalysis           *OptimizationROIAnalysis          `json:"roi_analysis"`
	CreatedAt             time.Time                         `json:"created_at"`
}

// CostOptimizationRecommendation represents a specific optimization recommendation
type CostOptimizationRecommendation struct {
	ID                  string                    `json:"id"`
	Title               string                    `json:"title"`
	Category            string                    `json:"category"`
	Priority            string                    `json:"priority"` // "high", "medium", "low"
	SavingsPotential    float64                   `json:"savings_potential"`
	SavingsPercentage   float64                   `json:"savings_percentage"`
	ImplementationCost  float64                   `json:"implementation_cost"`
	PaybackPeriod       string                    `json:"payback_period"`
	Description         string                    `json:"description"`
	TechnicalDetails    string                    `json:"technical_details"`
	ImplementationSteps []*CostImplementationStep `json:"implementation_steps"`
	Prerequisites       []string                  `json:"prerequisites"`
	Risks               []string                  `json:"risks"`
	Benefits            []string                  `json:"benefits"`
	AffectedServices    []string                  `json:"affected_services"`
	Effort              string                    `json:"effort"`
	Timeline            string                    `json:"timeline"`
	Validation          *RecommendationValidation `json:"validation"`
}

// CostImplementationStep represents a step in implementing an optimization
type CostImplementationStep struct {
	StepNumber  int      `json:"step_number"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	Owner       string   `json:"owner"`
	Tools       []string `json:"tools"`
	Validation  []string `json:"validation"`
}

// RecommendationValidation represents validation criteria for recommendations
type RecommendationValidation struct {
	ValidationCriteria []string `json:"validation_criteria"`
	SuccessMetrics     []string `json:"success_metrics"`
	MonitoringPoints   []string `json:"monitoring_points"`
	RollbackPlan       []string `json:"rollback_plan"`
}

// OptimizationImplementationPlan represents the overall implementation plan
type OptimizationImplementationPlan struct {
	TotalDuration        string                            `json:"total_duration"`
	Phases               []*CostOptimizationPhase          `json:"phases"`
	ResourceRequirements *OptimizationResourceRequirements `json:"resource_requirements"`
	Dependencies         []string                          `json:"dependencies"`
	Milestones           []*OptimizationMilestone          `json:"milestones"`
	RiskMitigation       []string                          `json:"risk_mitigation"`
}

// CostOptimizationPhase represents a phase in the optimization implementation
type CostOptimizationPhase struct {
	PhaseName       string   `json:"phase_name"`
	Duration        string   `json:"duration"`
	Objectives      []string `json:"objectives"`
	Deliverables    []string `json:"deliverables"`
	Prerequisites   []string `json:"prerequisites"`
	SuccessCriteria []string `json:"success_criteria"`
}

// OptimizationResourceRequirements represents resource requirements for optimization
type OptimizationResourceRequirements struct {
	TeamSize       int      `json:"team_size"`
	SkillsRequired []string `json:"skills_required"`
	ToolsRequired  []string `json:"tools_required"`
	BudgetRequired float64  `json:"budget_required"`
	TimeCommitment string   `json:"time_commitment"`
}

// OptimizationMilestone represents a milestone in optimization implementation
type OptimizationMilestone struct {
	MilestoneName  string    `json:"milestone_name"`
	TargetDate     time.Time `json:"target_date"`
	Deliverables   []string  `json:"deliverables"`
	SuccessMetrics []string  `json:"success_metrics"`
}

// CostOptimizationRiskAssessment represents risk assessment for optimization
type CostOptimizationRiskAssessment struct {
	OverallRiskLevel string              `json:"overall_risk_level"`
	Risks            []*OptimizationRisk `json:"risks"`
	MitigationPlan   *RiskMitigationPlan `json:"mitigation_plan"`
	ContingencyPlan  *ContingencyPlan    `json:"contingency_plan"`
}

// OptimizationRisk represents a risk in optimization implementation
type OptimizationRisk struct {
	RiskID          string   `json:"risk_id"`
	RiskName        string   `json:"risk_name"`
	Category        string   `json:"category"`
	Impact          string   `json:"impact"`
	Probability     string   `json:"probability"`
	RiskScore       int      `json:"risk_score"`
	Description     string   `json:"description"`
	Triggers        []string `json:"triggers"`
	Consequences    []string `json:"consequences"`
	MitigationSteps []string `json:"mitigation_steps"`
}

// RiskMitigationPlan represents the risk mitigation plan
type RiskMitigationPlan struct {
	Strategy          string                  `json:"strategy"`
	MitigationActions []*RiskMitigationAction `json:"mitigation_actions"`
	MonitoringPlan    *RiskMonitoringPlan     `json:"monitoring_plan"`
	EscalationPlan    *RiskEscalationPlan     `json:"escalation_plan"`
}

// RiskMitigationAction represents a risk mitigation action
type RiskMitigationAction struct {
	ActionID  string   `json:"action_id"`
	RiskID    string   `json:"risk_id"`
	Action    string   `json:"action"`
	Owner     string   `json:"owner"`
	Timeline  string   `json:"timeline"`
	Resources []string `json:"resources"`
	Success   []string `json:"success_criteria"`
}

// RiskMonitoringPlan represents the risk monitoring plan
type RiskMonitoringPlan struct {
	MonitoringFrequency string   `json:"monitoring_frequency"`
	KeyIndicators       []string `json:"key_indicators"`
	ReportingSchedule   string   `json:"reporting_schedule"`
	EscalationTriggers  []string `json:"escalation_triggers"`
}

// RiskEscalationPlan represents the risk escalation plan
type RiskEscalationPlan struct {
	EscalationLevels []string `json:"escalation_levels"`
	ContactList      []string `json:"contact_list"`
	EscalationMatrix []string `json:"escalation_matrix"`
}

// ContingencyPlan represents the contingency plan
type ContingencyPlan struct {
	Scenarios         []*ContingencyScenario `json:"scenarios"`
	ResponsePlans     []*ResponsePlan        `json:"response_plans"`
	ResourceBackup    []string               `json:"resource_backup"`
	CommunicationPlan string                 `json:"communication_plan"`
}

// ContingencyScenario represents a contingency scenario
type ContingencyScenario struct {
	ScenarioID   string   `json:"scenario_id"`
	ScenarioName string   `json:"scenario_name"`
	Description  string   `json:"description"`
	Triggers     []string `json:"triggers"`
	Impact       string   `json:"impact"`
	Probability  string   `json:"probability"`
}

// ResponsePlan represents a response plan for contingency scenarios
type ResponsePlan struct {
	ScenarioID      string   `json:"scenario_id"`
	ResponseActions []string `json:"response_actions"`
	Timeline        string   `json:"timeline"`
	Resources       []string `json:"resources"`
	SuccessCriteria []string `json:"success_criteria"`
}

// OptimizationROIAnalysis represents ROI analysis for optimization
type OptimizationROIAnalysis struct {
	TotalInvestment     float64               `json:"total_investment"`
	TotalSavings        float64               `json:"total_savings"`
	NetBenefit          float64               `json:"net_benefit"`
	ROIPercentage       float64               `json:"roi_percentage"`
	PaybackPeriod       string                `json:"payback_period"`
	NPV                 float64               `json:"npv"`
	IRR                 float64               `json:"irr"`
	CashFlowProjection  []*CashFlowProjection `json:"cash_flow_projection"`
	SensitivityAnalysis *SensitivityAnalysis  `json:"sensitivity_analysis"`
	BreakEvenAnalysis   *BreakEvenAnalysis    `json:"break_even_analysis"`
}

// CashFlowProjection represents cash flow projection
type CashFlowProjection struct {
	Period      string  `json:"period"`
	Investment  float64 `json:"investment"`
	Savings     float64 `json:"savings"`
	NetCashFlow float64 `json:"net_cash_flow"`
	Cumulative  float64 `json:"cumulative"`
}

// SensitivityAnalysis represents sensitivity analysis
type SensitivityAnalysis struct {
	Variables   []*SensitivityVariable `json:"variables"`
	Scenarios   []*SensitivityScenario `json:"scenarios"`
	RiskFactors []string               `json:"risk_factors"`
}

// SensitivityVariable represents a variable in sensitivity analysis
type SensitivityVariable struct {
	VariableName string  `json:"variable_name"`
	BaseValue    float64 `json:"base_value"`
	LowValue     float64 `json:"low_value"`
	HighValue    float64 `json:"high_value"`
	Impact       string  `json:"impact"`
}

// SensitivityScenario represents a scenario in sensitivity analysis
type SensitivityScenario struct {
	ScenarioName  string  `json:"scenario_name"`
	ROI           float64 `json:"roi"`
	NPV           float64 `json:"npv"`
	PaybackPeriod string  `json:"payback_period"`
	Probability   string  `json:"probability"`
}

// BreakEvenAnalysis represents break-even analysis
type BreakEvenAnalysis struct {
	BreakEvenPoint   string  `json:"break_even_point"`
	BreakEvenSavings float64 `json:"break_even_savings"`
	TimeToBreakEven  string  `json:"time_to_break_even"`
	MarginOfSafety   float64 `json:"margin_of_safety"`
}

// Supporting types for cost analysis
type ServiceUsageMetrics struct {
	ComputeHours   float64 `json:"compute_hours"`
	StorageGB      float64 `json:"storage_gb"`
	DataTransferGB float64 `json:"data_transfer_gb"`
	RequestCount   int64   `json:"request_count"`
	Utilization    float64 `json:"utilization"`
	PeakUsage      float64 `json:"peak_usage"`
	AverageUsage   float64 `json:"average_usage"`
}

// CostComponent represents a component of service cost
type CostComponent struct {
	ComponentName string  `json:"component_name"`
	Cost          float64 `json:"cost"`
	Unit          string  `json:"unit"`
	Quantity      float64 `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
}

// OptimizationOption represents an optimization opportunity
type OptimizationOption struct {
	OptionName           string  `json:"option_name"`
	SavingsPotential     float64 `json:"savings_potential"`
	ImplementationEffort string  `json:"implementation_effort"`
	RiskLevel            string  `json:"risk_level"`
	Description          string  `json:"description"`
}

// CostArchitectureSpec represents architecture specification for cost analysis
type CostArchitectureSpec struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Components  []*CostArchitectureComponent `json:"components"`
	Services    []*ServiceSpec               `json:"services"`
	DataFlow    []*DataFlowSpec              `json:"data_flow"`
	Regions     []string                     `json:"regions"`
	Environment string                       `json:"environment"`
	Metadata    map[string]interface{}       `json:"metadata"`
}

// CostArchitectureComponent represents a component in the architecture
type CostArchitectureComponent struct {
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	Description    string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications"`
	Dependencies   []string               `json:"dependencies"`
	CostDrivers    []string               `json:"cost_drivers"`
}

// ServiceSpec represents a service specification
type ServiceSpec struct {
	ServiceName   string                 `json:"service_name"`
	Provider      string                 `json:"provider"`
	Region        string                 `json:"region"`
	Configuration map[string]interface{} `json:"configuration"`
	UsagePattern  *CostUsagePattern      `json:"usage_pattern"`
	PricingTier   string                 `json:"pricing_tier"`
}

// DataFlowSpec represents data flow specification
type DataFlowSpec struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	DataVolume  float64 `json:"data_volume"`
	Frequency   string  `json:"frequency"`
	Protocol    string  `json:"protocol"`
}

// CostUsagePattern represents usage patterns for services
type CostUsagePattern struct {
	Pattern     string  `json:"pattern"` // "steady", "burst", "seasonal", "unpredictable"
	BaseUsage   float64 `json:"base_usage"`
	PeakUsage   float64 `json:"peak_usage"`
	GrowthRate  float64 `json:"growth_rate"`
	Seasonality string  `json:"seasonality"`
	PeakHours   []int   `json:"peak_hours"`
	PeakDays    []int   `json:"peak_days"`
}

// Reserved Instance Analysis types

// UsageData represents usage data for reserved instance analysis
type UsageData struct {
	AccountID     string                  `json:"account_id"`
	Region        string                  `json:"region"`
	TimeRange     *TimeRange              `json:"time_range"`
	ComputeUsage  []*ComputeUsageData     `json:"compute_usage"`
	DatabaseUsage []*DatabaseUsageData    `json:"database_usage"`
	StorageUsage  []*StorageUsageData     `json:"storage_usage"`
	UsagePatterns []*UsagePatternAnalysis `json:"usage_patterns"`
	CostData      []*HistoricalCostData   `json:"cost_data"`
	Metadata      map[string]interface{}  `json:"metadata"`
}

// TimeRange represents a time range for analysis
type TimeRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Period    string    `json:"period"` // "daily", "weekly", "monthly"
}

// ComputeUsageData represents compute usage data
type ComputeUsageData struct {
	InstanceType     string    `json:"instance_type"`
	Region           string    `json:"region"`
	AvailabilityZone string    `json:"availability_zone"`
	Platform         string    `json:"platform"`
	Tenancy          string    `json:"tenancy"`
	UsageHours       float64   `json:"usage_hours"`
	OnDemandCost     float64   `json:"on_demand_cost"`
	Date             time.Time `json:"date"`
	Utilization      float64   `json:"utilization"`
}

// DatabaseUsageData represents database usage data
type DatabaseUsageData struct {
	Engine        string    `json:"engine"`
	InstanceClass string    `json:"instance_class"`
	Region        string    `json:"region"`
	MultiAZ       bool      `json:"multi_az"`
	UsageHours    float64   `json:"usage_hours"`
	OnDemandCost  float64   `json:"on_demand_cost"`
	Date          time.Time `json:"date"`
	Utilization   float64   `json:"utilization"`
}

// StorageUsageData represents storage usage data
type StorageUsageData struct {
	StorageType string    `json:"storage_type"`
	Region      string    `json:"region"`
	UsageGB     float64   `json:"usage_gb"`
	Cost        float64   `json:"cost"`
	Date        time.Time `json:"date"`
	GrowthRate  float64   `json:"growth_rate"`
}

// UsagePatternAnalysis represents usage pattern analysis
type UsagePatternAnalysis struct {
	ResourceType    string            `json:"resource_type"`
	Pattern         string            `json:"pattern"`
	Consistency     float64           `json:"consistency"`
	Predictability  float64           `json:"predictability"`
	SeasonalFactors []*SeasonalFactor `json:"seasonal_factors"`
	TrendAnalysis   *TrendAnalysis    `json:"trend_analysis"`
}

// SeasonalFactor represents seasonal factors in usage
type SeasonalFactor struct {
	Period      string  `json:"period"`
	Factor      float64 `json:"factor"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// TrendAnalysis represents trend analysis
type TrendAnalysis struct {
	Direction  string  `json:"direction"` // "increasing", "decreasing", "stable"
	Rate       float64 `json:"rate"`
	Confidence float64 `json:"confidence"`
	Projection string  `json:"projection"`
}

// HistoricalCostData represents historical cost data
type HistoricalCostData struct {
	Date         time.Time `json:"date"`
	Service      string    `json:"service"`
	ResourceType string    `json:"resource_type"`
	OnDemandCost float64   `json:"on_demand_cost"`
	ReservedCost float64   `json:"reserved_cost"`
	SavingsUsed  float64   `json:"savings_used"`
}

// ReservedInstanceAnalysis represents reserved instance analysis results
type ReservedInstanceAnalysis struct {
	ID                    string                `json:"id"`
	AnalysisDate          time.Time             `json:"analysis_date"`
	TotalSavingsPotential float64               `json:"total_savings_potential"`
	Recommendations       []*RIRecommendation   `json:"recommendations"`
	CurrentUtilization    *RICurrentUtilization `json:"current_utilization"`
	OptimalPortfolio      *RIOptimalPortfolio   `json:"optimal_portfolio"`
	PaybackAnalysis       *RIPaybackAnalysis    `json:"payback_analysis"`
	RiskAssessment        *RIRiskAssessment     `json:"risk_assessment"`
	ImplementationPlan    *RIImplementationPlan `json:"implementation_plan"`
	CreatedAt             time.Time             `json:"created_at"`
}

// RIRecommendation represents a reserved instance recommendation
type RIRecommendation struct {
	ID                   string                 `json:"id"`
	InstanceType         string                 `json:"instance_type"`
	Region               string                 `json:"region"`
	Platform             string                 `json:"platform"`
	Tenancy              string                 `json:"tenancy"`
	Term                 string                 `json:"term"`           // "1year", "3year"
	PaymentOption        string                 `json:"payment_option"` // "no_upfront", "partial_upfront", "all_upfront"
	RecommendedQuantity  int                    `json:"recommended_quantity"`
	CurrentOnDemandCost  float64                `json:"current_on_demand_cost"`
	ReservedInstanceCost float64                `json:"reserved_instance_cost"`
	AnnualSavings        float64                `json:"annual_savings"`
	SavingsPercentage    float64                `json:"savings_percentage"`
	UpfrontCost          float64                `json:"upfront_cost"`
	MonthlyCost          float64                `json:"monthly_cost"`
	BreakEvenMonths      int                    `json:"break_even_months"`
	UtilizationRequired  float64                `json:"utilization_required"`
	CurrentUtilization   float64                `json:"current_utilization"`
	ConfidenceLevel      string                 `json:"confidence_level"`
	RiskFactors          []string               `json:"risk_factors"`
	Justification        string                 `json:"justification"`
	AlternativeOptions   []*RIAlternativeOption `json:"alternative_options"`
}

// RIAlternativeOption represents alternative RI options
type RIAlternativeOption struct {
	Term              string  `json:"term"`
	PaymentOption     string  `json:"payment_option"`
	Quantity          int     `json:"quantity"`
	AnnualSavings     float64 `json:"annual_savings"`
	SavingsPercentage float64 `json:"savings_percentage"`
	UpfrontCost       float64 `json:"upfront_cost"`
	RiskLevel         string  `json:"risk_level"`
}

// RICurrentUtilization represents current RI utilization
type RICurrentUtilization struct {
	TotalReservedInstances    int                          `json:"total_reserved_instances"`
	UtilizedInstances         int                          `json:"utilized_instances"`
	UnderutilizedInstances    int                          `json:"underutilized_instances"`
	OverallUtilization        float64                      `json:"overall_utilization"`
	UtilizationByType         []*RIUtilizationByType       `json:"utilization_by_type"`
	WastedSpend               float64                      `json:"wasted_spend"`
	OptimizationOpportunities []*RIOptimizationOpportunity `json:"optimization_opportunities"`
}

// RIUtilizationByType represents RI utilization by instance type
type RIUtilizationByType struct {
	InstanceType   string  `json:"instance_type"`
	Reserved       int     `json:"reserved"`
	Utilized       int     `json:"utilized"`
	Utilization    float64 `json:"utilization"`
	WastedSpend    float64 `json:"wasted_spend"`
	Recommendation string  `json:"recommendation"`
}

// RIOptimizationOpportunity represents RI optimization opportunities
type RIOptimizationOpportunity struct {
	OpportunityType      string  `json:"opportunity_type"`
	Description          string  `json:"description"`
	SavingsPotential     float64 `json:"savings_potential"`
	ImplementationEffort string  `json:"implementation_effort"`
	RiskLevel            string  `json:"risk_level"`
	ActionRequired       string  `json:"action_required"`
}

// RIOptimalPortfolio represents the optimal RI portfolio
type RIOptimalPortfolio struct {
	TotalInvestment       float64            `json:"total_investment"`
	TotalAnnualSavings    float64            `json:"total_annual_savings"`
	OverallSavingsPercent float64            `json:"overall_savings_percent"`
	Portfolio             []*RIPortfolioItem `json:"portfolio"`
	RiskProfile           string             `json:"risk_profile"`
	DiversificationScore  float64            `json:"diversification_score"`
	FlexibilityScore      float64            `json:"flexibility_score"`
}

// RIPortfolioItem represents an item in the RI portfolio
type RIPortfolioItem struct {
	InstanceType     string  `json:"instance_type"`
	Region           string  `json:"region"`
	Term             string  `json:"term"`
	PaymentOption    string  `json:"payment_option"`
	Quantity         int     `json:"quantity"`
	Investment       float64 `json:"investment"`
	AnnualSavings    float64 `json:"annual_savings"`
	PortfolioWeight  float64 `json:"portfolio_weight"`
	RiskContribution float64 `json:"risk_contribution"`
}

// RIPaybackAnalysis represents payback analysis for RI investments
type RIPaybackAnalysis struct {
	AveragePaybackPeriod    string                       `json:"average_payback_period"`
	PaybackByRecommendation []*RIPaybackByRecommendation `json:"payback_by_recommendation"`
	CashFlowProjection      []*RICashFlowProjection      `json:"cash_flow_projection"`
	ROIAnalysis             *RIROIAnalysis               `json:"roi_analysis"`
}

// RIPaybackByRecommendation represents payback by recommendation
type RIPaybackByRecommendation struct {
	RecommendationID string  `json:"recommendation_id"`
	InstanceType     string  `json:"instance_type"`
	PaybackMonths    int     `json:"payback_months"`
	ROI              float64 `json:"roi"`
	NPV              float64 `json:"npv"`
}

// RICashFlowProjection represents RI cash flow projection
type RICashFlowProjection struct {
	Month             int     `json:"month"`
	Investment        float64 `json:"investment"`
	Savings           float64 `json:"savings"`
	NetCashFlow       float64 `json:"net_cash_flow"`
	CumulativeSavings float64 `json:"cumulative_savings"`
}

// RIROIAnalysis represents RI ROI analysis
type RIROIAnalysis struct {
	ThreeYearROI    float64 `json:"three_year_roi"`
	FiveYearROI     float64 `json:"five_year_roi"`
	NPV             float64 `json:"npv"`
	IRR             float64 `json:"irr"`
	PaybackPeriod   string  `json:"payback_period"`
	RiskAdjustedROI float64 `json:"risk_adjusted_roi"`
}

// RIRiskAssessment represents RI risk assessment
type RIRiskAssessment struct {
	OverallRiskLevel     string                  `json:"overall_risk_level"`
	RiskFactors          []*RIRiskFactor         `json:"risk_factors"`
	MitigationStrategies []*RIMitigationStrategy `json:"mitigation_strategies"`
	RiskScore            float64                 `json:"risk_score"`
	RecommendedActions   []string                `json:"recommended_actions"`
}

// RIRiskFactor represents an RI risk factor
type RIRiskFactor struct {
	RiskType    string  `json:"risk_type"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Probability string  `json:"probability"`
	RiskScore   float64 `json:"risk_score"`
	Mitigation  string  `json:"mitigation"`
}

// RIMitigationStrategy represents RI mitigation strategy
type RIMitigationStrategy struct {
	StrategyName        string   `json:"strategy_name"`
	Description         string   `json:"description"`
	ImplementationSteps []string `json:"implementation_steps"`
	Effectiveness       string   `json:"effectiveness"`
	Cost                float64  `json:"cost"`
}

// RIImplementationPlan represents RI implementation plan
type RIImplementationPlan struct {
	TotalDuration        string                    `json:"total_duration"`
	Phases               []*RIImplementationPhase  `json:"phases"`
	Prerequisites        []string                  `json:"prerequisites"`
	ResourceRequirements *RIResourceRequirements   `json:"resource_requirements"`
	Timeline             *RIImplementationTimeline `json:"timeline"`
	SuccessMetrics       []string                  `json:"success_metrics"`
}

// RIImplementationPhase represents an RI implementation phase
type RIImplementationPhase struct {
	PhaseName    string   `json:"phase_name"`
	Duration     string   `json:"duration"`
	Objectives   []string `json:"objectives"`
	Activities   []string `json:"activities"`
	Deliverables []string `json:"deliverables"`
	Dependencies []string `json:"dependencies"`
}

// RIResourceRequirements represents RI resource requirements
type RIResourceRequirements struct {
	TeamMembers    []string `json:"team_members"`
	SkillsRequired []string `json:"skills_required"`
	ToolsRequired  []string `json:"tools_required"`
	BudgetRequired float64  `json:"budget_required"`
	TimeCommitment string   `json:"time_commitment"`
}

// RIImplementationTimeline represents RI implementation timeline
type RIImplementationTimeline struct {
	StartDate    time.Time                    `json:"start_date"`
	EndDate      time.Time                    `json:"end_date"`
	Milestones   []*RIImplementationMilestone `json:"milestones"`
	CriticalPath []string                     `json:"critical_path"`
	Dependencies []string                     `json:"dependencies"`
}

// RIImplementationMilestone represents an RI implementation milestone
type RIImplementationMilestone struct {
	MilestoneName   string    `json:"milestone_name"`
	TargetDate      time.Time `json:"target_date"`
	Description     string    `json:"description"`
	Deliverables    []string  `json:"deliverables"`
	SuccessCriteria []string  `json:"success_criteria"`
}

// Savings Plans Analysis types

// SavingsPlansAnalysis represents savings plans analysis results
type SavingsPlansAnalysis struct {
	ID                    string                          `json:"id"`
	AnalysisDate          time.Time                       `json:"analysis_date"`
	TotalSavingsPotential float64                         `json:"total_savings_potential"`
	Recommendations       []*SavingsPlansRecommendation   `json:"recommendations"`
	CurrentCommitments    *CurrentSavingsPlansCommitments `json:"current_commitments"`
	OptimalPortfolio      *SavingsPlansOptimalPortfolio   `json:"optimal_portfolio"`
	PaybackAnalysis       *SavingsPlansPaybackAnalysis    `json:"payback_analysis"`
	RiskAssessment        *SavingsPlansRiskAssessment     `json:"risk_assessment"`
	ImplementationPlan    *SavingsPlansImplementationPlan `json:"implementation_plan"`
	CreatedAt             time.Time                       `json:"created_at"`
}

// SavingsPlansRecommendation represents a savings plans recommendation
type SavingsPlansRecommendation struct {
	ID                  string                           `json:"id"`
	PlanType            string                           `json:"plan_type"` // "compute", "ec2_instance", "sagemaker"
	Term                string                           `json:"term"`      // "1year", "3year"
	PaymentOption       string                           `json:"payment_option"`
	HourlyCommitment    float64                          `json:"hourly_commitment"`
	AnnualCommitment    float64                          `json:"annual_commitment"`
	CurrentOnDemandCost float64                          `json:"current_on_demand_cost"`
	SavingsPlansRate    float64                          `json:"savings_plans_rate"`
	AnnualSavings       float64                          `json:"annual_savings"`
	SavingsPercentage   float64                          `json:"savings_percentage"`
	UpfrontCost         float64                          `json:"upfront_cost"`
	MonthlyCost         float64                          `json:"monthly_cost"`
	BreakEvenMonths     int                              `json:"break_even_months"`
	UtilizationRequired float64                          `json:"utilization_required"`
	CurrentUtilization  float64                          `json:"current_utilization"`
	CoverageScope       []string                         `json:"coverage_scope"`
	ConfidenceLevel     string                           `json:"confidence_level"`
	RiskFactors         []string                         `json:"risk_factors"`
	Justification       string                           `json:"justification"`
	AlternativeOptions  []*SavingsPlansAlternativeOption `json:"alternative_options"`
}

// SavingsPlansAlternativeOption represents alternative savings plans options
type SavingsPlansAlternativeOption struct {
	Term              string  `json:"term"`
	PaymentOption     string  `json:"payment_option"`
	HourlyCommitment  float64 `json:"hourly_commitment"`
	AnnualSavings     float64 `json:"annual_savings"`
	SavingsPercentage float64 `json:"savings_percentage"`
	UpfrontCost       float64 `json:"upfront_cost"`
	RiskLevel         string  `json:"risk_level"`
	Flexibility       string  `json:"flexibility"`
}

// CurrentSavingsPlansCommitments represents current savings plans commitments
type CurrentSavingsPlansCommitments struct {
	TotalCommitments          float64                                `json:"total_commitments"`
	UtilizedCommitments       float64                                `json:"utilized_commitments"`
	UnderutilizedCommitments  float64                                `json:"underutilized_commitments"`
	OverallUtilization        float64                                `json:"overall_utilization"`
	CommitmentsByType         []*SavingsPlansCommitmentByType        `json:"commitments_by_type"`
	WastedSpend               float64                                `json:"wasted_spend"`
	OptimizationOpportunities []*SavingsPlansOptimizationOpportunity `json:"optimization_opportunities"`
}

// SavingsPlansCommitmentByType represents savings plans commitment by type
type SavingsPlansCommitmentByType struct {
	PlanType       string  `json:"plan_type"`
	Commitment     float64 `json:"commitment"`
	Utilized       float64 `json:"utilized"`
	Utilization    float64 `json:"utilization"`
	WastedSpend    float64 `json:"wasted_spend"`
	Recommendation string  `json:"recommendation"`
}

// SavingsPlansOptimizationOpportunity represents savings plans optimization opportunities
type SavingsPlansOptimizationOpportunity struct {
	OpportunityType      string  `json:"opportunity_type"`
	Description          string  `json:"description"`
	SavingsPotential     float64 `json:"savings_potential"`
	ImplementationEffort string  `json:"implementation_effort"`
	RiskLevel            string  `json:"risk_level"`
	ActionRequired       string  `json:"action_required"`
}

// SavingsPlansOptimalPortfolio represents the optimal savings plans portfolio
type SavingsPlansOptimalPortfolio struct {
	TotalCommitment       float64                      `json:"total_commitment"`
	TotalAnnualSavings    float64                      `json:"total_annual_savings"`
	OverallSavingsPercent float64                      `json:"overall_savings_percent"`
	Portfolio             []*SavingsPlansPortfolioItem `json:"portfolio"`
	RiskProfile           string                       `json:"risk_profile"`
	FlexibilityScore      float64                      `json:"flexibility_score"`
	CoverageScore         float64                      `json:"coverage_score"`
}

// SavingsPlansPortfolioItem represents an item in the savings plans portfolio
type SavingsPlansPortfolioItem struct {
	PlanType         string  `json:"plan_type"`
	Term             string  `json:"term"`
	PaymentOption    string  `json:"payment_option"`
	HourlyCommitment float64 `json:"hourly_commitment"`
	AnnualCommitment float64 `json:"annual_commitment"`
	AnnualSavings    float64 `json:"annual_savings"`
	PortfolioWeight  float64 `json:"portfolio_weight"`
	RiskContribution float64 `json:"risk_contribution"`
}

// SavingsPlansPaybackAnalysis represents payback analysis for savings plans
type SavingsPlansPaybackAnalysis struct {
	AveragePaybackPeriod    string                                 `json:"average_payback_period"`
	PaybackByRecommendation []*SavingsPlansPaybackByRecommendation `json:"payback_by_recommendation"`
	CashFlowProjection      []*SavingsPlansCashFlowProjection      `json:"cash_flow_projection"`
	ROIAnalysis             *SavingsPlansROIAnalysis               `json:"roi_analysis"`
}

// SavingsPlansPaybackByRecommendation represents payback by recommendation
type SavingsPlansPaybackByRecommendation struct {
	RecommendationID string  `json:"recommendation_id"`
	PlanType         string  `json:"plan_type"`
	PaybackMonths    int     `json:"payback_months"`
	ROI              float64 `json:"roi"`
	NPV              float64 `json:"npv"`
}

// SavingsPlansCashFlowProjection represents savings plans cash flow projection
type SavingsPlansCashFlowProjection struct {
	Month             int     `json:"month"`
	Commitment        float64 `json:"commitment"`
	Savings           float64 `json:"savings"`
	NetCashFlow       float64 `json:"net_cash_flow"`
	CumulativeSavings float64 `json:"cumulative_savings"`
}

// SavingsPlansROIAnalysis represents savings plans ROI analysis
type SavingsPlansROIAnalysis struct {
	ThreeYearROI    float64 `json:"three_year_roi"`
	FiveYearROI     float64 `json:"five_year_roi"`
	NPV             float64 `json:"npv"`
	IRR             float64 `json:"irr"`
	PaybackPeriod   string  `json:"payback_period"`
	RiskAdjustedROI float64 `json:"risk_adjusted_roi"`
}

// SavingsPlansRiskAssessment represents savings plans risk assessment
type SavingsPlansRiskAssessment struct {
	OverallRiskLevel     string                            `json:"overall_risk_level"`
	RiskFactors          []*SavingsPlansRiskFactor         `json:"risk_factors"`
	MitigationStrategies []*SavingsPlansMitigationStrategy `json:"mitigation_strategies"`
	RiskScore            float64                           `json:"risk_score"`
	RecommendedActions   []string                          `json:"recommended_actions"`
}

// SavingsPlansRiskFactor represents a savings plans risk factor
type SavingsPlansRiskFactor struct {
	RiskType    string  `json:"risk_type"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Probability string  `json:"probability"`
	RiskScore   float64 `json:"risk_score"`
	Mitigation  string  `json:"mitigation"`
}

// SavingsPlansMitigationStrategy represents savings plans mitigation strategy
type SavingsPlansMitigationStrategy struct {
	StrategyName        string   `json:"strategy_name"`
	Description         string   `json:"description"`
	ImplementationSteps []string `json:"implementation_steps"`
	Effectiveness       string   `json:"effectiveness"`
	Cost                float64  `json:"cost"`
}

// SavingsPlansImplementationPlan represents savings plans implementation plan
type SavingsPlansImplementationPlan struct {
	TotalDuration        string                              `json:"total_duration"`
	Phases               []*SavingsPlansImplementationPhase  `json:"phases"`
	Prerequisites        []string                            `json:"prerequisites"`
	ResourceRequirements *SavingsPlansResourceRequirements   `json:"resource_requirements"`
	Timeline             *SavingsPlansImplementationTimeline `json:"timeline"`
	SuccessMetrics       []string                            `json:"success_metrics"`
}

// SavingsPlansImplementationPhase represents a savings plans implementation phase
type SavingsPlansImplementationPhase struct {
	PhaseName    string   `json:"phase_name"`
	Duration     string   `json:"duration"`
	Objectives   []string `json:"objectives"`
	Activities   []string `json:"activities"`
	Deliverables []string `json:"deliverables"`
	Dependencies []string `json:"dependencies"`
}

// SavingsPlansResourceRequirements represents savings plans resource requirements
type SavingsPlansResourceRequirements struct {
	TeamMembers    []string `json:"team_members"`
	SkillsRequired []string `json:"skills_required"`
	ToolsRequired  []string `json:"tools_required"`
	BudgetRequired float64  `json:"budget_required"`
	TimeCommitment string   `json:"time_commitment"`
}

// SavingsPlansImplementationTimeline represents savings plans implementation timeline
type SavingsPlansImplementationTimeline struct {
	StartDate    time.Time                              `json:"start_date"`
	EndDate      time.Time                              `json:"end_date"`
	Milestones   []*SavingsPlansImplementationMilestone `json:"milestones"`
	CriticalPath []string                               `json:"critical_path"`
	Dependencies []string                               `json:"dependencies"`
}

// SavingsPlansImplementationMilestone represents a savings plans implementation milestone
type SavingsPlansImplementationMilestone struct {
	MilestoneName   string    `json:"milestone_name"`
	TargetDate      time.Time `json:"target_date"`
	Description     string    `json:"description"`
	Deliverables    []string  `json:"deliverables"`
	SuccessCriteria []string  `json:"success_criteria"`
}

// Purchase Recommendations types

// PurchaseRecommendations represents combined purchase recommendations
type PurchaseRecommendations struct {
	ID                      string                        `json:"id"`
	ReservedInstanceRecs    []*RIRecommendation           `json:"reserved_instance_recommendations"`
	SavingsPlansRecs        []*SavingsPlansRecommendation `json:"savings_plans_recommendations"`
	OptimalMix              *OptimalPurchaseMix           `json:"optimal_mix"`
	TotalSavingsPotential   float64                       `json:"total_savings_potential"`
	TotalInvestmentRequired float64                       `json:"total_investment_required"`
	OverallROI              float64                       `json:"overall_roi"`
	ImplementationPriority  []*PurchasePriority           `json:"implementation_priority"`
	CreatedAt               time.Time                     `json:"created_at"`
}

// OptimalPurchaseMix represents the optimal mix of RI and Savings Plans
type OptimalPurchaseMix struct {
	RIInvestment           float64 `json:"ri_investment"`
	SavingsPlansInvestment float64 `json:"savings_plans_investment"`
	RISavings              float64 `json:"ri_savings"`
	SavingsPlansSavings    float64 `json:"savings_plans_savings"`
	TotalSavings           float64 `json:"total_savings"`
	RiskProfile            string  `json:"risk_profile"`
	FlexibilityScore       float64 `json:"flexibility_score"`
	Justification          string  `json:"justification"`
}

// PurchasePriority represents purchase priority
type PurchasePriority struct {
	RecommendationID     string  `json:"recommendation_id"`
	RecommendationType   string  `json:"recommendation_type"`
	Priority             int     `json:"priority"`
	SavingsPotential     float64 `json:"savings_potential"`
	RiskLevel            string  `json:"risk_level"`
	ImplementationEffort string  `json:"implementation_effort"`
	Justification        string  `json:"justification"`
}

// SavingsCalculation represents exact savings calculation
type SavingsCalculation struct {
	ID                 string                      `json:"id"`
	RecommendationsID  string                      `json:"recommendations_id"`
	CalculationDate    time.Time                   `json:"calculation_date"`
	TotalCurrentCost   float64                     `json:"total_current_cost"`
	TotalOptimizedCost float64                     `json:"total_optimized_cost"`
	TotalSavings       float64                     `json:"total_savings"`
	SavingsPercentage  float64                     `json:"savings_percentage"`
	SavingsBreakdown   []*SavingsBreakdownItem     `json:"savings_breakdown"`
	MonthlyProjection  []*MonthlySavingsProjection `json:"monthly_projection"`
	AnnualProjection   []*AnnualSavingsProjection  `json:"annual_projection"`
	ConfidenceLevel    float64                     `json:"confidence_level"`
	Assumptions        []string                    `json:"assumptions"`
	RiskFactors        []string                    `json:"risk_factors"`
	CreatedAt          time.Time                   `json:"created_at"`
}

// SavingsBreakdownItem represents a savings breakdown item
type SavingsBreakdownItem struct {
	RecommendationID   string  `json:"recommendation_id"`
	RecommendationType string  `json:"recommendation_type"`
	Service            string  `json:"service"`
	CurrentCost        float64 `json:"current_cost"`
	OptimizedCost      float64 `json:"optimized_cost"`
	Savings            float64 `json:"savings"`
	SavingsPercentage  float64 `json:"savings_percentage"`
	ConfidenceLevel    float64 `json:"confidence_level"`
}

// MonthlySavingsProjection represents monthly savings projection
type MonthlySavingsProjection struct {
	Month         int     `json:"month"`
	CurrentCost   float64 `json:"current_cost"`
	OptimizedCost float64 `json:"optimized_cost"`
	Savings       float64 `json:"savings"`
	Cumulative    float64 `json:"cumulative_savings"`
}

// AnnualSavingsProjection represents annual savings projection
type AnnualSavingsProjection struct {
	Year          int     `json:"year"`
	CurrentCost   float64 `json:"current_cost"`
	OptimizedCost float64 `json:"optimized_cost"`
	Savings       float64 `json:"savings"`
	Cumulative    float64 `json:"cumulative_savings"`
}

// ImplementationPlan represents a plan for implementing cost optimizations
type ImplementationPlan struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Phases            []*ImplementationPhase `json:"phases"`
	EstimatedDuration string                 `json:"estimated_duration"`
	TotalCost         float64                `json:"total_cost"`
	Prerequisites     []string               `json:"prerequisites"`
	RiskFactors       []string               `json:"risk_factors"`
	SuccessMetrics    []string               `json:"success_metrics"`
	CreatedAt         time.Time              `json:"created_at"`
}

// UsageAnalysis represents AWS usage analysis
type UsageAnalysis struct {
	ID               string                 `json:"id"`
	AccountID        string                 `json:"account_id"`
	AnalysisPeriod   string                 `json:"analysis_period"`
	EC2Usage         map[string]interface{} `json:"ec2_usage"`
	RDSUsage         map[string]interface{} `json:"rds_usage"`
	ElastiCacheUsage map[string]interface{} `json:"elasticache_usage"`
	S3Usage          map[string]interface{} `json:"s3_usage"`
	LambdaUsage      map[string]interface{} `json:"lambda_usage"`
	TotalCost        float64                `json:"total_cost"`
	CreatedAt        time.Time              `json:"created_at"`
}

// RIUtilizationAnalysis represents RI utilization analysis
type RIUtilizationAnalysis struct {
	OverallUtilization        float64            `json:"overall_utilization"`
	UnderutilizedRIs          []*UnderutilizedRI `json:"underutilized_ris"`
	OptimizationOpportunities []string           `json:"optimization_opportunities"`
	PotentialSavings          float64            `json:"potential_savings"`
}

// UnderutilizedRI represents an underutilized reserved instance
type UnderutilizedRI struct {
	InstanceType    string  `json:"instance_type"`
	Quantity        int     `json:"quantity"`
	UtilizationRate float64 `json:"utilization_rate"`
	WastedCost      float64 `json:"wasted_cost"`
	Recommendation  string  `json:"recommendation"`
}

// OptimalRIPortfolio represents an optimal RI portfolio recommendation
type OptimalRIPortfolio struct {
	TotalInvestment   float64             `json:"total_investment"`
	AnnualSavings     float64             `json:"annual_savings"`
	PaybackPeriod     float64             `json:"payback_period"`
	RIRecommendations []*RIRecommendation `json:"ri_recommendations"`
	RiskAssessment    *RIRiskAssessment   `json:"risk_assessment"`
}
