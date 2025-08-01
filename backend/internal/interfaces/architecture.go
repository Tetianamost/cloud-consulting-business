package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// ArchitectureAnalyzer defines the interface for advanced AWS architecture analysis
type ArchitectureAnalyzer interface {
	AnalyzeArchitecture(ctx context.Context, inquiry *domain.Inquiry, architecture *ArchitectureSpec) (*ArchitectureAnalysis, error)
	IdentifyBottlenecks(ctx context.Context, architecture *ArchitectureSpec) ([]*PerformanceBottleneck, error)
	GenerateScalingRecommendations(ctx context.Context, architecture *ArchitectureSpec, workload *WorkloadProfile) ([]*ScalingRecommendation, error)
	ValidateArchitecture(ctx context.Context, architecture *ArchitectureSpec) (*ArchitectureValidation, error)
	CompareArchitectures(ctx context.Context, current, proposed *ArchitectureSpec) (*ArchitectureComparison, error)
}

// CostOptimizationEngine defines the interface for cost optimization analysis
type CostOptimizationEngine interface {
	AnalyzeCosts(ctx context.Context, architecture *ArchitectureSpec) (*CostAnalysis, error)
	IdentifySavingsOpportunities(ctx context.Context, architecture *ArchitectureSpec) ([]*SavingsOpportunity, error)
	GenerateCostOptimizationPlan(ctx context.Context, architecture *ArchitectureSpec) (*CostOptimizationPlan, error)
	EstimateROI(ctx context.Context, currentCosts, optimizedCosts *CostBreakdown) (*ROIAnalysis, error)
	GetRightsizingRecommendations(ctx context.Context, resources []*ResourceUsage) ([]*RightsizingRecommendation, error)
}

// SecurityAnalyzer defines the interface for security assessment
type SecurityAnalyzer interface {
	AssessSecurity(ctx context.Context, architecture *ArchitectureSpec) (*SecurityAssessment, error)
	IdentifyVulnerabilities(ctx context.Context, architecture *ArchitectureSpec) ([]*SecurityVulnerability, error)
	GenerateRemediationPlan(ctx context.Context, vulnerabilities []*SecurityVulnerability) (*SecurityRemediationPlan, error)
	ValidateCompliance(ctx context.Context, architecture *ArchitectureSpec, frameworks []string) (*ComplianceValidation, error)
	GenerateSecurityRecommendations(ctx context.Context, architecture *ArchitectureSpec) ([]*SecurityRecommendation, error)
}

// PerformanceAnalyzer defines the interface for performance analysis
type PerformanceAnalyzer interface {
	AnalyzePerformance(ctx context.Context, architecture *ArchitectureSpec, metrics *PerformanceMetrics) (*PerformanceAnalysis, error)
	IdentifyBottlenecks(ctx context.Context, architecture *ArchitectureSpec, metrics *PerformanceMetrics) ([]*PerformanceBottleneck, error)
	GenerateOptimizationRecommendations(ctx context.Context, bottlenecks []*PerformanceBottleneck) ([]*PerformanceOptimization, error)
	PredictScalingNeeds(ctx context.Context, architecture *ArchitectureSpec, growthProjection *GrowthProjection) (*ScalingPrediction, error)
}

// ArchitectureSpec represents a detailed architecture specification
type ArchitectureSpec struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	CloudProvider    string                   `json:"cloud_provider"`
	Region           string                   `json:"region"`
	Components       []*ArchitectureComponent `json:"components"`
	NetworkTopology  *NetworkTopology         `json:"network_topology"`
	DataFlow         []*DataFlowPath          `json:"data_flow"`
	SecurityControls []*SecurityControl       `json:"security_controls"`
	Monitoring       *MonitoringConfiguration `json:"monitoring"`
	BackupStrategy   *BackupConfiguration     `json:"backup_strategy"`
	DisasterRecovery *DisasterRecoveryConfig  `json:"disaster_recovery"`
	Compliance       []string                 `json:"compliance"`
	Tags             map[string]string        `json:"tags"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

// ArchitectureComponent represents a component in the architecture
type ArchitectureComponent struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`         // "compute", "storage", "database", "network", "security"
	Layer            string                 `json:"layer"`        // "presentation", "application", "data", "infrastructure"
	ServiceName      string                 `json:"service_name"` // e.g., "EC2", "RDS", "S3"
	Configuration    map[string]interface{} `json:"configuration"`
	Dependencies     []string               `json:"dependencies"`
	Criticality      string                 `json:"criticality"` // "low", "medium", "high", "critical"
	AvailabilityZone string                 `json:"availability_zone"`
	ResourceTags     map[string]string      `json:"resource_tags"`
	CostCenter       string                 `json:"cost_center"`
	Owner            string                 `json:"owner"`
}

// ArchitectureAnalysis represents the result of architecture analysis
type ArchitectureAnalysis struct {
	ID                    string                        `json:"id"`
	ArchitectureID        string                        `json:"architecture_id"`
	OverallScore          float64                       `json:"overall_score"` // 0-100
	Strengths             []string                      `json:"strengths"`
	Weaknesses            []string                      `json:"weaknesses"`
	Recommendations       []*ArchitectureRecommendation `json:"recommendations"`
	CostAnalysis          *CostAnalysis                 `json:"cost_analysis"`
	SecurityAssessment    *SecurityAssessment           `json:"security_assessment"`
	PerformanceAnalysis   *PerformanceAnalysis          `json:"performance_analysis"`
	ComplianceStatus      *ComplianceValidation         `json:"compliance_status"`
	RiskAssessment        *RiskAssessment               `json:"risk_assessment"`
	OptimizationPotential *OptimizationPotential        `json:"optimization_potential"`
	CreatedAt             time.Time                     `json:"created_at"`
}

// ArchitectureRecommendation represents a recommendation for architecture improvement
type ArchitectureRecommendation struct {
	ID                  string    `json:"id"`
	Category            string    `json:"category"` // "cost", "security", "performance", "reliability"
	Priority            string    `json:"priority"` // "low", "medium", "high", "critical"
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	Impact              string    `json:"impact"`
	ImplementationSteps []string  `json:"implementation_steps"`
	EstimatedEffort     string    `json:"estimated_effort"`
	EstimatedCost       string    `json:"estimated_cost"`
	EstimatedSavings    string    `json:"estimated_savings"`
	Timeline            string    `json:"timeline"`
	Dependencies        []string  `json:"dependencies"`
	RiskLevel           string    `json:"risk_level"`
	Documentation       []string  `json:"documentation"`
	CreatedAt           time.Time `json:"created_at"`
}

// CostAnalysis represents detailed cost analysis
type CostAnalysis struct {
	ID                    string                  `json:"id"`
	TotalMonthlyCost      float64                 `json:"total_monthly_cost"`
	TotalAnnualCost       float64                 `json:"total_annual_cost"`
	CostBreakdown         *CostBreakdown          `json:"cost_breakdown"`
	SavingsOpportunities  []*SavingsOpportunity   `json:"savings_opportunities"`
	CostTrends            *CostTrends             `json:"cost_trends"`
	BudgetRecommendations []*BudgetRecommendation `json:"budget_recommendations"`
	OptimizationPlan      *CostOptimizationPlan   `json:"optimization_plan"`
	ROIAnalysis           *ROIAnalysis            `json:"roi_analysis"`
	CreatedAt             time.Time               `json:"created_at"`
}

// SavingsOpportunity represents a cost savings opportunity
type SavingsOpportunity struct {
	ID                  string    `json:"id"`
	Type                string    `json:"type"` // "rightsizing", "reserved_instances", "spot_instances", "storage_optimization"
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	CurrentCost         float64   `json:"current_cost"`
	OptimizedCost       float64   `json:"optimized_cost"`
	MonthlySavings      float64   `json:"monthly_savings"`
	AnnualSavings       float64   `json:"annual_savings"`
	SavingsPercentage   float64   `json:"savings_percentage"`
	ImplementationSteps []string  `json:"implementation_steps"`
	Effort              string    `json:"effort"` // "low", "medium", "high"
	Risk                string    `json:"risk"`   // "low", "medium", "high"
	Timeline            string    `json:"timeline"`
	Prerequisites       []string  `json:"prerequisites"`
	AffectedResources   []string  `json:"affected_resources"`
	CreatedAt           time.Time `json:"created_at"`
}

// SecurityAssessment represents security assessment results
type SecurityAssessment struct {
	ID                      string                      `json:"id"`
	OverallSecurityScore    float64                     `json:"overall_security_score"` // 0-100
	SecurityPosture         string                      `json:"security_posture"`       // "poor", "fair", "good", "excellent"
	Vulnerabilities         []*SecurityVulnerability    `json:"vulnerabilities"`
	ComplianceGaps          []*ComplianceGap            `json:"compliance_gaps"`
	SecurityRecommendations []*SecurityRecommendation   `json:"security_recommendations"`
	RemediationPlan         *SecurityRemediationPlan    `json:"remediation_plan"`
	ThreatModel             *ThreatModel                `json:"threat_model"`
	SecurityControls        *SecurityControlsAssessment `json:"security_controls"`
	CreatedAt               time.Time                   `json:"created_at"`
}

// SecurityVulnerability represents a security vulnerability
type SecurityVulnerability struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Severity          string    `json:"severity"` // "low", "medium", "high", "critical"
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	AffectedResources []string  `json:"affected_resources"`
	ThreatLevel       string    `json:"threat_level"`
	ExploitComplexity string    `json:"exploit_complexity"`
	Impact            string    `json:"impact"`
	CVSS              float64   `json:"cvss"`
	CWE               string    `json:"cwe"`
	References        []string  `json:"references"`
	CreatedAt         time.Time `json:"created_at"`
}

// PerformanceAnalysis represents performance analysis results
type PerformanceAnalysis struct {
	ID                          string                     `json:"id"`
	OverallPerformanceScore     float64                    `json:"overall_performance_score"` // 0-100
	Bottlenecks                 []*PerformanceBottleneck   `json:"bottlenecks"`
	OptimizationRecommendations []*PerformanceOptimization `json:"optimization_recommendations"`
	ScalingRecommendations      []*ScalingRecommendation   `json:"scaling_recommendations"`
	CapacityPlanning            *CapacityPlan              `json:"capacity_planning"`
	PerformanceMetrics          *PerformanceMetrics        `json:"performance_metrics"`
	BenchmarkComparison         *BenchmarkComparison       `json:"benchmark_comparison"`
	CreatedAt                   time.Time                  `json:"created_at"`
}

// PerformanceBottleneck represents a performance bottleneck
type PerformanceBottleneck struct {
	ID                 string             `json:"id"`
	Type               string             `json:"type"`     // "cpu", "memory", "network", "storage", "database"
	Severity           string             `json:"severity"` // "low", "medium", "high", "critical"
	Component          string             `json:"component"`
	Description        string             `json:"description"`
	Impact             string             `json:"impact"`
	CurrentMetrics     map[string]float64 `json:"current_metrics"`
	ThresholdMetrics   map[string]float64 `json:"threshold_metrics"`
	AffectedOperations []string           `json:"affected_operations"`
	RootCause          string             `json:"root_cause"`
	CreatedAt          time.Time          `json:"created_at"`
}

// Supporting types (simplified for this implementation)
type NetworkTopology struct {
	VPCConfiguration map[string]interface{} `json:"vpc_configuration"`
	SubnetStrategy   string                 `json:"subnet_strategy"`
	SecurityGroups   []SecurityGroup        `json:"security_groups"`
	LoadBalancers    []LoadBalancer         `json:"load_balancers"`
	CDNConfiguration map[string]interface{} `json:"cdn_configuration"`
}

type SecurityGroup struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Rules       []FirewallRule `json:"rules"`
	Scope       string         `json:"scope"`
}

type FirewallRule struct {
	Direction   string `json:"direction"` // "inbound", "outbound"
	Protocol    string `json:"protocol"`
	Port        string `json:"port"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Action      string `json:"action"` // "allow", "deny"
}

type DataFlowPath struct {
	ID            string                 `json:"id"`
	Source        string                 `json:"source"`
	Destination   string                 `json:"destination"`
	Protocol      string                 `json:"protocol"`
	Port          int                    `json:"port"`
	DataType      string                 `json:"data_type"`
	Volume        string                 `json:"volume"`
	Frequency     string                 `json:"frequency"`
	Encryption    bool                   `json:"encryption"`
	Compression   bool                   `json:"compression"`
	Configuration map[string]interface{} `json:"configuration"`
}

type SecurityControl struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Description    string   `json:"description"`
	Implementation string   `json:"implementation"`
	Effectiveness  string   `json:"effectiveness"`
	Coverage       []string `json:"coverage"`
}

type MonitoringConfiguration struct {
	Metrics    []*MetricConfiguration `json:"metrics"`
	Alarms     []*AlarmConfiguration  `json:"alarms"`
	Dashboards []string               `json:"dashboards"`
	LogGroups  []string               `json:"log_groups"`
	Tracing    bool                   `json:"tracing"`
	APMEnabled bool                   `json:"apm_enabled"`
}

type MetricConfiguration struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Dimensions    map[string]string      `json:"dimensions"`
	Statistic     string                 `json:"statistic"`
	Period        int                    `json:"period"`
	Unit          string                 `json:"unit"`
	Threshold     float64                `json:"threshold"`
	Configuration map[string]interface{} `json:"configuration"`
}

type AlarmConfiguration struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	MetricName         string   `json:"metric_name"`
	Threshold          float64  `json:"threshold"`
	ComparisonOperator string   `json:"comparison_operator"`
	EvaluationPeriods  int      `json:"evaluation_periods"`
	Actions            []string `json:"actions"`
}

type BackupConfiguration struct {
	Strategy        string                 `json:"strategy"`
	Frequency       string                 `json:"frequency"`
	RetentionPeriod string                 `json:"retention_period"`
	CrossRegion     bool                   `json:"cross_region"`
	Encryption      bool                   `json:"encryption"`
	Configuration   map[string]interface{} `json:"configuration"`
}

type DisasterRecoveryConfig struct {
	Strategy        string                 `json:"strategy"`
	RPO             string                 `json:"rpo"`
	RTO             string                 `json:"rto"`
	SecondaryRegion string                 `json:"secondary_region"`
	Automation      bool                   `json:"automation"`
	Configuration   map[string]interface{} `json:"configuration"`
}

// Additional supporting types
type CostBreakdown struct {
	ByService     map[string]float64 `json:"by_service"`
	ByCategory    map[string]float64 `json:"by_category"`
	ByRegion      map[string]float64 `json:"by_region"`
	ByCostCenter  map[string]float64 `json:"by_cost_center"`
	ByEnvironment map[string]float64 `json:"by_environment"`
}

type CostTrends struct {
	MonthlyTrend string             `json:"monthly_trend"`
	GrowthRate   float64            `json:"growth_rate"`
	Seasonality  bool               `json:"seasonality"`
	Forecast     map[string]float64 `json:"forecast"`
}

type BudgetRecommendation struct {
	Category          string  `json:"category"`
	RecommendedBudget float64 `json:"recommended_budget"`
	Justification     string  `json:"justification"`
}

type CostOptimizationPlan struct {
	ID                    string                      `json:"id"`
	TotalPotentialSavings float64                     `json:"total_potential_savings"`
	Phases                []*OptimizationPhase        `json:"phases"`
	QuickWins             []*SavingsOpportunity       `json:"quick_wins"`
	LongTermInitiatives   []*SavingsOpportunity       `json:"long_term_initiatives"`
	RiskAssessment        *OptimizationRiskAssessment `json:"risk_assessment"`
	Timeline              string                      `json:"timeline"`
	ROIProjection         *ROIProjection              `json:"roi_projection"`
	CreatedAt             time.Time                   `json:"created_at"`
}

type OptimizationPhase struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	Duration        string                `json:"duration"`
	Opportunities   []*SavingsOpportunity `json:"opportunities"`
	ExpectedSavings float64               `json:"expected_savings"`
	Prerequisites   []string              `json:"prerequisites"`
	SuccessMetrics  []string              `json:"success_metrics"`
	RiskMitigation  []string              `json:"risk_mitigation"`
}

type ROIAnalysis struct {
	InvestmentCost float64 `json:"investment_cost"`
	AnnualSavings  float64 `json:"annual_savings"`
	PaybackPeriod  string  `json:"payback_period"`
	ROIPercentage  float64 `json:"roi_percentage"`
	NPV            float64 `json:"npv"`
}

type ROIProjection struct {
	Year1Savings float64 `json:"year1_savings"`
	Year2Savings float64 `json:"year2_savings"`
	Year3Savings float64 `json:"year3_savings"`
	TotalROI     float64 `json:"total_roi"`
}

type OptimizationRiskAssessment struct {
	OverallRisk     string   `json:"overall_risk"`
	RiskFactors     []string `json:"risk_factors"`
	MitigationSteps []string `json:"mitigation_steps"`
}

type SecurityRecommendation struct {
	ID                   string    `json:"id"`
	Category             string    `json:"category"`
	Priority             string    `json:"priority"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	RemediationSteps     []string  `json:"remediation_steps"`
	EstimatedEffort      string    `json:"estimated_effort"`
	SecurityImprovement  string    `json:"security_improvement"`
	ComplianceFrameworks []string  `json:"compliance_frameworks"`
	AffectedComponents   []string  `json:"affected_components"`
	Documentation        []string  `json:"documentation"`
	CreatedAt            time.Time `json:"created_at"`
}

type SecurityRemediationPlan struct {
	ID                   string                    `json:"id"`
	TotalVulnerabilities int                       `json:"total_vulnerabilities"`
	CriticalCount        int                       `json:"critical_count"`
	HighCount            int                       `json:"high_count"`
	MediumCount          int                       `json:"medium_count"`
	LowCount             int                       `json:"low_count"`
	Phases               []*RemediationPhase       `json:"phases"`
	QuickFixes           []*SecurityRecommendation `json:"quick_fixes"`
	LongTermActions      []*SecurityRecommendation `json:"long_term_actions"`
	Timeline             string                    `json:"timeline"`
	EstimatedCost        string                    `json:"estimated_cost"`
	RiskReduction        string                    `json:"risk_reduction"`
	CreatedAt            time.Time                 `json:"created_at"`
}

type RemediationPhase struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Duration        string                    `json:"duration"`
	Vulnerabilities []*SecurityVulnerability  `json:"vulnerabilities"`
	Recommendations []*SecurityRecommendation `json:"recommendations"`
	Priority        string                    `json:"priority"`
	Dependencies    []string                  `json:"dependencies"`
	SuccessMetrics  []string                  `json:"success_metrics"`
}

type ComplianceValidation struct {
	OverallCompliance float64                     `json:"overall_compliance"`
	FrameworkStatus   map[string]ComplianceStatus `json:"framework_status"`
	Gaps              []*ComplianceGap            `json:"gaps"`
}

type ComplianceStatus struct {
	Framework   string    `json:"framework"`
	Compliance  float64   `json:"compliance"`
	Status      string    `json:"status"`
	LastChecked time.Time `json:"last_checked"`
}

type ComplianceGap struct {
	Framework        string   `json:"framework"`
	Requirement      string   `json:"requirement"`
	CurrentState     string   `json:"current_state"`
	RequiredState    string   `json:"required_state"`
	Severity         string   `json:"severity"`
	RemediationSteps []string `json:"remediation_steps"`
}

type ThreatModel struct {
	Threats       []*ThreatScenario `json:"threats"`
	AttackVectors []*AttackVector   `json:"attack_vectors"`
	RiskMatrix    map[string]string `json:"risk_matrix"`
}

type ThreatScenario struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Likelihood  string   `json:"likelihood"`
	Impact      string   `json:"impact"`
	Mitigations []string `json:"mitigations"`
}

type AttackVector struct {
	Vector      string   `json:"vector"`
	Description string   `json:"description"`
	Complexity  string   `json:"complexity"`
	Mitigations []string `json:"mitigations"`
}

type SecurityControlsAssessment struct {
	ImplementedControls []string `json:"implemented_controls"`
	MissingControls     []string `json:"missing_controls"`
	WeakControls        []string `json:"weak_controls"`
	Recommendations     []string `json:"recommendations"`
}

type PerformanceOptimization struct {
	ID                  string    `json:"id"`
	BottleneckID        string    `json:"bottleneck_id"`
	Type                string    `json:"type"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	ExpectedImprovement string    `json:"expected_improvement"`
	ImplementationSteps []string  `json:"implementation_steps"`
	EstimatedEffort     string    `json:"estimated_effort"`
	EstimatedCost       string    `json:"estimated_cost"`
	RiskLevel           string    `json:"risk_level"`
	Prerequisites       []string  `json:"prerequisites"`
	ValidationMetrics   []string  `json:"validation_metrics"`
	CreatedAt           time.Time `json:"created_at"`
}

type ScalingRecommendation struct {
	ID                  string                 `json:"id"`
	Component           string                 `json:"component"`
	ScalingType         string                 `json:"scaling_type"`
	Trigger             string                 `json:"trigger"`
	CurrentCapacity     map[string]interface{} `json:"current_capacity"`
	RecommendedCapacity map[string]interface{} `json:"recommended_capacity"`
	Justification       string                 `json:"justification"`
	ExpectedBenefit     string                 `json:"expected_benefit"`
	ImplementationSteps []string               `json:"implementation_steps"`
	CostImpact          string                 `json:"cost_impact"`
	Timeline            string                 `json:"timeline"`
	CreatedAt           time.Time              `json:"created_at"`
}

type WorkloadProfile struct {
	Type            string                 `json:"type"`
	TrafficPattern  string                 `json:"traffic_pattern"`
	PeakLoad        map[string]interface{} `json:"peak_load"`
	AverageLoad     map[string]interface{} `json:"average_load"`
	GrowthRate      float64                `json:"growth_rate"`
	SeasonalFactors map[string]float64     `json:"seasonal_factors"`
	SLARequirements *SLARequirements       `json:"sla_requirements"`
}

type SLARequirements struct {
	Availability float64 `json:"availability"`
	ResponseTime string  `json:"response_time"`
	Throughput   string  `json:"throughput"`
	ErrorRate    float64 `json:"error_rate"`
	RecoveryTime string  `json:"recovery_time"`
}

type PerformanceMetrics struct {
	CPU                *MetricData `json:"cpu"`
	Memory             *MetricData `json:"memory"`
	Network            *MetricData `json:"network"`
	Storage            *MetricData `json:"storage"`
	DatabaseMetrics    *MetricData `json:"database_metrics"`
	ApplicationMetrics *MetricData `json:"application_metrics"`
	UserExperience     *MetricData `json:"user_experience"`
}

type MetricData struct {
	Current   float64                `json:"current"`
	Average   float64                `json:"average"`
	Peak      float64                `json:"peak"`
	Trend     string                 `json:"trend"`
	Threshold float64                `json:"threshold"`
	Unit      string                 `json:"unit"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type CapacityPlan struct {
	CurrentCapacity   map[string]interface{} `json:"current_capacity"`
	ProjectedCapacity map[string]interface{} `json:"projected_capacity"`
	ScalingMilestones []CapacityMilestone    `json:"scaling_milestones"`
	ResourcePlan      []ResourcePlanItem     `json:"resource_plan"`
}

type CapacityMilestone struct {
	Date          time.Time              `json:"date"`
	Trigger       string                 `json:"trigger"`
	Capacity      map[string]interface{} `json:"capacity"`
	Actions       []string               `json:"actions"`
	EstimatedCost float64                `json:"estimated_cost"`
}

type ResourcePlanItem struct {
	ResourceType  string                 `json:"resource_type"`
	Timeline      string                 `json:"timeline"`
	Quantity      int                    `json:"quantity"`
	Configuration map[string]interface{} `json:"configuration"`
	Cost          float64                `json:"cost"`
}

type BenchmarkComparison struct {
	Industry      string             `json:"industry"`
	Percentile    float64            `json:"percentile"`
	Metrics       map[string]float64 `json:"metrics"`
	Gaps          []string           `json:"gaps"`
	Opportunities []string           `json:"opportunities"`
}

type GrowthProjection struct {
	TimeHorizon     string                 `json:"time_horizon"`
	GrowthRate      float64                `json:"growth_rate"`
	SeasonalFactors map[string]float64     `json:"seasonal_factors"`
	BusinessDrivers []string               `json:"business_drivers"`
	Assumptions     map[string]interface{} `json:"assumptions"`
}

type ScalingPrediction struct {
	TimeToScale     string                   `json:"time_to_scale"`
	ScalingTriggers []string                 `json:"scaling_triggers"`
	Recommendations []*ScalingRecommendation `json:"recommendations"`
	CostProjection  map[string]float64       `json:"cost_projection"`
}

type ResourceUsage struct {
	ResourceID   string                 `json:"resource_id"`
	ResourceType string                 `json:"resource_type"`
	CurrentUsage map[string]interface{} `json:"current_usage"`
	Capacity     map[string]interface{} `json:"capacity"`
	Utilization  float64                `json:"utilization"`
}

type RightsizingRecommendation struct {
	ResourceID               string                 `json:"resource_id"`
	CurrentConfiguration     map[string]interface{} `json:"current_configuration"`
	RecommendedConfiguration map[string]interface{} `json:"recommended_configuration"`
	MonthlySavings           float64                `json:"monthly_savings"`
	Justification            string                 `json:"justification"`
}

type OptimizationPotential struct {
	CostSavings            float64 `json:"cost_savings"`
	PerformanceGain        float64 `json:"performance_gain"`
	SecurityImprovement    float64 `json:"security_improvement"`
	ReliabilityImprovement float64 `json:"reliability_improvement"`
}

type ArchitectureValidation struct {
	IsValid     bool     `json:"is_valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	Suggestions []string `json:"suggestions"`
}

type ArchitectureComparison struct {
	Improvements []string `json:"improvements"`
	Regressions  []string `json:"regressions"`
	CostDelta    float64  `json:"cost_delta"`
	ScoreDelta   float64  `json:"score_delta"`
}

// LoadBalancer represents a load balancer configuration
type LoadBalancer struct {
	Type          string                 `json:"type"`
	Provider      string                 `json:"provider"`
	Configuration map[string]interface{} `json:"configuration"`
	HealthChecks  []HealthCheck          `json:"health_checks"`
}

// HealthCheck represents a health check configuration
type HealthCheck struct {
	Type               string `json:"type"`
	Endpoint           string `json:"endpoint"`
	Interval           string `json:"interval"`
	Timeout            string `json:"timeout"`
	HealthyThreshold   int    `json:"healthy_threshold"`
	UnhealthyThreshold int    `json:"unhealthy_threshold"`
}
