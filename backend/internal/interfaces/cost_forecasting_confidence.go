package interfaces

import (
	"time"
)

// ForecastConfidenceIntervals represents confidence intervals for cost forecasts
type ForecastConfidenceIntervals struct {
	ID                  string                        `json:"id"`
	ForecastID          string                        `json:"forecast_id"`
	CalculationDate     time.Time                     `json:"calculation_date"`
	ConfidenceLevel     float64                       `json:"confidence_level"`
	OverallConfidence   *OverallConfidenceInterval    `json:"overall_confidence"`
	ServiceConfidence   []*ServiceConfidenceInterval  `json:"service_confidence"`
	PeriodConfidence    []*PeriodConfidenceInterval   `json:"period_confidence"`
	ScenarioConfidence  []*ScenarioConfidenceInterval `json:"scenario_confidence"`
	UncertaintyAnalysis *UncertaintyAnalysis          `json:"uncertainty_analysis"`
	ConfidenceFactors   []*ConfidenceFactor           `json:"confidence_factors"`
	RiskAdjustments     []*RiskAdjustment             `json:"risk_adjustments"`
	CreatedAt           time.Time                     `json:"created_at"`
}

// OverallConfidenceInterval represents overall confidence interval
type OverallConfidenceInterval struct {
	ForecastValue   float64  `json:"forecast_value"`
	LowerBound      float64  `json:"lower_bound"`
	UpperBound      float64  `json:"upper_bound"`
	ConfidenceWidth float64  `json:"confidence_width"`
	RelativeWidth   float64  `json:"relative_width"`
	ConfidenceLevel float64  `json:"confidence_level"`
	Methodology     string   `json:"methodology"`
	Assumptions     []string `json:"assumptions"`
}

// ServiceConfidenceInterval represents confidence interval for a service
type ServiceConfidenceInterval struct {
	ServiceName       string   `json:"service_name"`
	Provider          string   `json:"provider"`
	ForecastValue     float64  `json:"forecast_value"`
	LowerBound        float64  `json:"lower_bound"`
	UpperBound        float64  `json:"upper_bound"`
	ConfidenceWidth   float64  `json:"confidence_width"`
	RelativeWidth     float64  `json:"relative_width"`
	ConfidenceLevel   float64  `json:"confidence_level"`
	UncertaintySource []string `json:"uncertainty_source"`
	DataQuality       string   `json:"data_quality"`
}

// PeriodConfidenceInterval represents confidence interval for a time period
type PeriodConfidenceInterval struct {
	Period             string    `json:"period"`
	Date               time.Time `json:"date"`
	ForecastValue      float64   `json:"forecast_value"`
	LowerBound         float64   `json:"lower_bound"`
	UpperBound         float64   `json:"upper_bound"`
	ConfidenceWidth    float64   `json:"confidence_width"`
	RelativeWidth      float64   `json:"relative_width"`
	ConfidenceLevel    float64   `json:"confidence_level"`
	SeasonalAdjustment float64   `json:"seasonal_adjustment"`
	TrendConfidence    float64   `json:"trend_confidence"`
}

// ScenarioConfidenceInterval represents confidence interval for scenarios
type ScenarioConfidenceInterval struct {
	ScenarioName     string   `json:"scenario_name"`
	Probability      float64  `json:"probability"`
	ForecastValue    float64  `json:"forecast_value"`
	LowerBound       float64  `json:"lower_bound"`
	UpperBound       float64  `json:"upper_bound"`
	ConfidenceWidth  float64  `json:"confidence_width"`
	RelativeWidth    float64  `json:"relative_width"`
	ConfidenceLevel  float64  `json:"confidence_level"`
	ScenarioRisk     string   `json:"scenario_risk"`
	KeyUncertainties []string `json:"key_uncertainties"`
}

// UncertaintyAnalysis represents uncertainty analysis
type UncertaintyAnalysis struct {
	TotalUncertainty     float64                 `json:"total_uncertainty"`
	UncertaintyBreakdown []*UncertaintyComponent `json:"uncertainty_breakdown"`
	UncertaintySources   []*UncertaintySource    `json:"uncertainty_sources"`
	SensitivityAnalysis  *UncertaintySensitivity `json:"sensitivity_analysis"`
	MonteCarloResults    *MonteCarloAnalysis     `json:"monte_carlo_results"`
	UncertaintyTrends    []*UncertaintyTrend     `json:"uncertainty_trends"`
}

// UncertaintyComponent represents a component of uncertainty
type UncertaintyComponent struct {
	ComponentName       string  `json:"component_name"`
	ComponentType       string  `json:"component_type"`
	UncertaintyValue    float64 `json:"uncertainty_value"`
	Contribution        float64 `json:"contribution"`
	ContributionPercent float64 `json:"contribution_percent"`
	Controllability     string  `json:"controllability"`
	Description         string  `json:"description"`
}

// UncertaintySource represents a source of uncertainty
type UncertaintySource struct {
	SourceName         string  `json:"source_name"`
	SourceType         string  `json:"source_type"`
	Impact             string  `json:"impact"`
	Probability        float64 `json:"probability"`
	UncertaintyLevel   string  `json:"uncertainty_level"`
	Mitigation         string  `json:"mitigation"`
	MonitoringRequired bool    `json:"monitoring_required"`
	Description        string  `json:"description"`
}

// UncertaintySensitivity represents uncertainty sensitivity analysis
type UncertaintySensitivity struct {
	SensitivityResults []*UncertaintySensitivityResult `json:"sensitivity_results"`
	KeyDrivers         []string                        `json:"key_drivers"`
	InteractionEffects []*InteractionEffect            `json:"interaction_effects"`
	Recommendations    []string                        `json:"recommendations"`
}

// UncertaintySensitivityResult represents sensitivity result for uncertainty
type UncertaintySensitivityResult struct {
	ParameterName    string  `json:"parameter_name"`
	BaseUncertainty  float64 `json:"base_uncertainty"`
	LowUncertainty   float64 `json:"low_uncertainty"`
	HighUncertainty  float64 `json:"high_uncertainty"`
	SensitivityIndex float64 `json:"sensitivity_index"`
	Impact           string  `json:"impact"`
}

// InteractionEffect represents interaction effects between parameters
type InteractionEffect struct {
	Parameter1          string  `json:"parameter1"`
	Parameter2          string  `json:"parameter2"`
	InteractionStrength float64 `json:"interaction_strength"`
	Effect              string  `json:"effect"`
	Significance        string  `json:"significance"`
}

// MonteCarloAnalysis represents Monte Carlo analysis results
type MonteCarloAnalysis struct {
	SimulationCount   int                     `json:"simulation_count"`
	ConvergenceStatus string                  `json:"convergence_status"`
	Results           *MonteCarloResults      `json:"results"`
	Distribution      *MonteCarloDistribution `json:"distribution"`
	Percentiles       []*MonteCarloPercentile `json:"percentiles"`
	Statistics        *MonteCarloStatistics   `json:"statistics"`
}

// MonteCarloResults represents Monte Carlo results
type MonteCarloResults struct {
	Mean              float64 `json:"mean"`
	Median            float64 `json:"median"`
	StandardDeviation float64 `json:"standard_deviation"`
	Variance          float64 `json:"variance"`
	Minimum           float64 `json:"minimum"`
	Maximum           float64 `json:"maximum"`
	Range             float64 `json:"range"`
	Skewness          float64 `json:"skewness"`
	Kurtosis          float64 `json:"kurtosis"`
}

// MonteCarloDistribution represents Monte Carlo distribution
type MonteCarloDistribution struct {
	DistributionType string             `json:"distribution_type"`
	Parameters       map[string]float64 `json:"parameters"`
	GoodnessOfFit    *GoodnessOfFitTest `json:"goodness_of_fit"`
	HistogramData    []*HistogramBin    `json:"histogram_data"`
	DensityFunction  []*DensityPoint    `json:"density_function"`
}

// GoodnessOfFitTest represents goodness of fit test
type GoodnessOfFitTest struct {
	TestName      string  `json:"test_name"`
	TestStatistic float64 `json:"test_statistic"`
	PValue        float64 `json:"p_value"`
	Conclusion    string  `json:"conclusion"`
	FitQuality    string  `json:"fit_quality"`
}

// HistogramBin represents a histogram bin
type HistogramBin struct {
	BinStart    float64 `json:"bin_start"`
	BinEnd      float64 `json:"bin_end"`
	BinCenter   float64 `json:"bin_center"`
	Frequency   int     `json:"frequency"`
	Probability float64 `json:"probability"`
	Density     float64 `json:"density"`
}

// DensityPoint represents a density function point
type DensityPoint struct {
	Value   float64 `json:"value"`
	Density float64 `json:"density"`
}

// MonteCarloPercentile represents Monte Carlo percentiles
type MonteCarloPercentile struct {
	Percentile float64 `json:"percentile"`
	Value      float64 `json:"value"`
	Label      string  `json:"label"`
}

// MonteCarloStatistics represents Monte Carlo statistics
type MonteCarloStatistics struct {
	ConfidenceIntervals []*MonteCarloConfidenceInterval `json:"confidence_intervals"`
	RiskMetrics         *MonteCarloRiskMetrics          `json:"risk_metrics"`
	TailAnalysis        *TailAnalysis                   `json:"tail_analysis"`
	ExtremeValues       *ExtremeValueAnalysis           `json:"extreme_values"`
}

// MonteCarloConfidenceInterval represents Monte Carlo confidence interval
type MonteCarloConfidenceInterval struct {
	ConfidenceLevel float64 `json:"confidence_level"`
	LowerBound      float64 `json:"lower_bound"`
	UpperBound      float64 `json:"upper_bound"`
	Width           float64 `json:"width"`
	RelativeWidth   float64 `json:"relative_width"`
}

// MonteCarloRiskMetrics represents Monte Carlo risk metrics
type MonteCarloRiskMetrics struct {
	ValueAtRisk       float64 `json:"value_at_risk"`
	ConditionalVaR    float64 `json:"conditional_var"`
	DownsideDeviation float64 `json:"downside_deviation"`
	UpsidePotential   float64 `json:"upside_potential"`
	ProbabilityOfLoss float64 `json:"probability_of_loss"`
	ExpectedShortfall float64 `json:"expected_shortfall"`
}

// TailAnalysis represents tail analysis
type TailAnalysis struct {
	LeftTail       *TailStatistics `json:"left_tail"`
	RightTail      *TailStatistics `json:"right_tail"`
	TailDependence float64         `json:"tail_dependence"`
	TailRisk       string          `json:"tail_risk"`
	TailEvents     []*TailEvent    `json:"tail_events"`
}

// TailStatistics represents tail statistics
type TailStatistics struct {
	TailProbability float64 `json:"tail_probability"`
	TailMean        float64 `json:"tail_mean"`
	TailVariance    float64 `json:"tail_variance"`
	TailSkewness    float64 `json:"tail_skewness"`
	TailKurtosis    float64 `json:"tail_kurtosis"`
}

// TailEvent represents a tail event
type TailEvent struct {
	EventType     string  `json:"event_type"`
	Threshold     float64 `json:"threshold"`
	Probability   float64 `json:"probability"`
	ExpectedValue float64 `json:"expected_value"`
	Impact        string  `json:"impact"`
}

// ExtremeValueAnalysis represents extreme value analysis
type ExtremeValueAnalysis struct {
	ExtremeValues         []*ExtremeValue    `json:"extreme_values"`
	ExtremeStatistics     *ExtremeStatistics `json:"extreme_statistics"`
	ReturnLevels          []*ReturnLevel     `json:"return_levels"`
	ExceedanceProbability float64            `json:"exceedance_probability"`
}

// ExtremeValue represents an extreme value
type ExtremeValue struct {
	Value       float64 `json:"value"`
	Probability float64 `json:"probability"`
	Rank        int     `json:"rank"`
	Type        string  `json:"type"`
	Impact      string  `json:"impact"`
}

// ExtremeStatistics represents extreme statistics
type ExtremeStatistics struct {
	MaximumValue   float64 `json:"maximum_value"`
	MinimumValue   float64 `json:"minimum_value"`
	ExtremeRange   float64 `json:"extreme_range"`
	ExtremeRatio   float64 `json:"extreme_ratio"`
	OutlierCount   int     `json:"outlier_count"`
	OutlierPercent float64 `json:"outlier_percent"`
}

// ReturnLevel represents return level analysis
type ReturnLevel struct {
	ReturnPeriod       float64        `json:"return_period"`
	ReturnLevel        float64        `json:"return_level"`
	ConfidenceInterval *ReturnLevelCI `json:"confidence_interval"`
	Probability        float64        `json:"probability"`
}

// ReturnLevelCI represents return level confidence interval
type ReturnLevelCI struct {
	LowerBound float64 `json:"lower_bound"`
	UpperBound float64 `json:"upper_bound"`
	Width      float64 `json:"width"`
}

// UncertaintyTrend represents uncertainty trend over time
type UncertaintyTrend struct {
	Period           string   `json:"period"`
	UncertaintyLevel float64  `json:"uncertainty_level"`
	Trend            string   `json:"trend"`
	ChangeRate       float64  `json:"change_rate"`
	Factors          []string `json:"factors"`
}

// ConfidenceFactor represents factors affecting confidence
type ConfidenceFactor struct {
	FactorName      string  `json:"factor_name"`
	FactorType      string  `json:"factor_type"`
	Impact          string  `json:"impact"`
	Magnitude       float64 `json:"magnitude"`
	Controllability string  `json:"controllability"`
	Description     string  `json:"description"`
	Mitigation      string  `json:"mitigation"`
}

// RiskAdjustment represents risk adjustments to confidence
type RiskAdjustment struct {
	AdjustmentName  string  `json:"adjustment_name"`
	AdjustmentType  string  `json:"adjustment_type"`
	RiskFactor      string  `json:"risk_factor"`
	AdjustmentValue float64 `json:"adjustment_value"`
	Justification   string  `json:"justification"`
	Impact          string  `json:"impact"`
}

// ComprehensiveCostAnalysis represents comprehensive cost analysis
type ComprehensiveCostAnalysis struct {
	ID                          string                           `json:"id"`
	InquiryID                   string                           `json:"inquiry_id"`
	AnalysisDate                time.Time                        `json:"analysis_date"`
	CostBreakdown               *CostBreakdownAnalysis           `json:"cost_breakdown"`
	OptimizationRecommendations *CostOptimizationRecommendations `json:"optimization_recommendations"`
	ReservedInstanceAnalysis    *ReservedInstanceAnalysis        `json:"reserved_instance_analysis"`
	SavingsPlansAnalysis        *SavingsPlansAnalysis            `json:"savings_plans_analysis"`
	RightSizingAnalysis         *RightSizingAnalysis             `json:"right_sizing_analysis"`
	RightSizingRecommendations  *RightSizingRecommendations      `json:"right_sizing_recommendations"`
	CostForecast                *CostForecast                    `json:"cost_forecast"`
	ConfidenceIntervals         *ForecastConfidenceIntervals     `json:"confidence_intervals"`
	ExecutiveSummary            *CostAnalysisExecutiveSummary    `json:"executive_summary"`
	ActionPlan                  *CostOptimizationActionPlan      `json:"action_plan"`
	CreatedAt                   time.Time                        `json:"created_at"`
}

// CostAnalysisExecutiveSummary represents executive summary of cost analysis
type CostAnalysisExecutiveSummary struct {
	CurrentMonthlyCost       float64  `json:"current_monthly_cost"`
	OptimizedMonthlyCost     float64  `json:"optimized_monthly_cost"`
	TotalSavingsPotential    float64  `json:"total_savings_potential"`
	SavingsPercentage        float64  `json:"savings_percentage"`
	PaybackPeriod            string   `json:"payback_period"`
	ROI                      float64  `json:"roi"`
	KeyFindings              []string `json:"key_findings"`
	TopRecommendations       []string `json:"top_recommendations"`
	RiskAssessment           string   `json:"risk_assessment"`
	ImplementationComplexity string   `json:"implementation_complexity"`
	TimeToRealizeSavings     string   `json:"time_to_realize_savings"`
	ConfidenceLevel          string   `json:"confidence_level"`
}

// CostOptimizationActionPlan represents action plan for cost optimization
type CostOptimizationActionPlan struct {
	TotalDuration        string                        `json:"total_duration"`
	TotalSavings         float64                       `json:"total_savings"`
	TotalInvestment      float64                       `json:"total_investment"`
	ActionItems          []*CostOptimizationActionItem `json:"action_items"`
	ImplementationPhases []*ForecastOptimizationPhase  `json:"implementation_phases"`
	ResourceRequirements *CostOptimizationResources    `json:"resource_requirements"`
	Timeline             *CostOptimizationTimeline     `json:"timeline"`
	RiskMitigation       []string                      `json:"risk_mitigation"`
	SuccessMetrics       []string                      `json:"success_metrics"`
	GovernanceStructure  *CostOptimizationGovernance   `json:"governance_structure"`
}

// CostOptimizationActionItem represents an action item
type CostOptimizationActionItem struct {
	ItemID             string   `json:"item_id"`
	Title              string   `json:"title"`
	Category           string   `json:"category"`
	Priority           string   `json:"priority"`
	SavingsPotential   float64  `json:"savings_potential"`
	ImplementationCost float64  `json:"implementation_cost"`
	Timeline           string   `json:"timeline"`
	Owner              string   `json:"owner"`
	Prerequisites      []string `json:"prerequisites"`
	Dependencies       []string `json:"dependencies"`
	RiskLevel          string   `json:"risk_level"`
	Status             string   `json:"status"`
}

// ForecastOptimizationPhase represents an optimization phase
type ForecastOptimizationPhase struct {
	PhaseName       string   `json:"phase_name"`
	Duration        string   `json:"duration"`
	Objectives      []string `json:"objectives"`
	ActionItems     []string `json:"action_items"`
	ExpectedSavings float64  `json:"expected_savings"`
	RiskLevel       string   `json:"risk_level"`
	SuccessCriteria []string `json:"success_criteria"`
}

// CostOptimizationResources represents resource requirements
type CostOptimizationResources struct {
	TeamSize       int      `json:"team_size"`
	SkillsRequired []string `json:"skills_required"`
	ToolsRequired  []string `json:"tools_required"`
	BudgetRequired float64  `json:"budget_required"`
	TimeCommitment string   `json:"time_commitment"`
}

// CostOptimizationTimeline represents optimization timeline
type CostOptimizationTimeline struct {
	StartDate    time.Time                    `json:"start_date"`
	EndDate      time.Time                    `json:"end_date"`
	Milestones   []*CostOptimizationMilestone `json:"milestones"`
	CriticalPath []string                     `json:"critical_path"`
}

// CostOptimizationMilestone represents an optimization milestone
type CostOptimizationMilestone struct {
	MilestoneName   string    `json:"milestone_name"`
	TargetDate      time.Time `json:"target_date"`
	Deliverables    []string  `json:"deliverables"`
	ExpectedSavings float64   `json:"expected_savings"`
}

// CostOptimizationGovernance represents governance structure
type CostOptimizationGovernance struct {
	SteeringCommittee  []string `json:"steering_committee"`
	ProjectManager     string   `json:"project_manager"`
	TechnicalLeads     []string `json:"technical_leads"`
	BusinessOwners     []string `json:"business_owners"`
	ReportingStructure string   `json:"reporting_structure"`
	DecisionAuthority  string   `json:"decision_authority"`
	EscalationProcess  string   `json:"escalation_process"`
}
