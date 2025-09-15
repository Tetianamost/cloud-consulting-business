package interfaces

import "time"

// Supporting types for client-specific solution engine that don't conflict with existing types

// ArchitecturalLayer represents a layer in the architecture
type ArchitecturalLayer struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Components   []string `json:"components"`
	Dependencies []string `json:"dependencies"`
	Technologies []string `json:"technologies"`
}

// DataFlowStep represents a step in data flow
type DataFlowStep struct {
	StepNumber  int      `json:"step_number"`
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	DataType    string   `json:"data_type"`
	Processing  []string `json:"processing"`
	Security    []string `json:"security"`
}

// TechnicalStack represents the technical stack for a solution
type TechnicalStack struct {
	CloudProvider    string            `json:"cloud_provider"`
	ComputeServices  []string          `json:"compute_services"`
	StorageServices  []string          `json:"storage_services"`
	DatabaseServices []string          `json:"database_services"`
	NetworkServices  []string          `json:"network_services"`
	SecurityServices []string          `json:"security_services"`
	MonitoringTools  []string          `json:"monitoring_tools"`
	DevOpsTools      []string          `json:"devops_tools"`
	Frameworks       []string          `json:"frameworks"`
	Languages        []string          `json:"languages"`
	Integrations     map[string]string `json:"integrations"`
}

// ScalabilityFactor represents factors affecting scalability
type ScalabilityFactor struct {
	Factor     string `json:"factor"`
	Impact     string `json:"impact"`
	Mitigation string `json:"mitigation"`
	Monitoring string `json:"monitoring"`
}

// DataTier represents a data tier in the architecture
type DataTier struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DataTypes   []string `json:"data_types"`
	Storage     []string `json:"storage"`
	Access      []string `json:"access"`
}

// SecurityZone represents a security zone in the architecture
type SecurityZone struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TrustLevel  string   `json:"trust_level"`
	Controls    []string `json:"controls"`
	Components  []string `json:"components"`
}

// IntegrationPoint represents an integration point
type IntegrationPoint struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Protocol   string   `json:"protocol"`
	Security   []string `json:"security"`
	DataFormat string   `json:"data_format"`
}

// ComplianceControl represents a compliance control
type ComplianceControl struct {
	ControlID      string   `json:"control_id"`
	Framework      string   `json:"framework"`
	Description    string   `json:"description"`
	Implementation []string `json:"implementation"`
	Testing        []string `json:"testing"`
}

// PerformanceTargets represents performance targets
type PerformanceTargets struct {
	ResponseTime string `json:"response_time"`
	Throughput   string `json:"throughput"`
	Availability string `json:"availability"`
	Latency      string `json:"latency"`
}

// AvailabilityTargets represents availability targets
type AvailabilityTargets struct {
	Uptime string `json:"uptime"`
	RTO    string `json:"rto"`
	RPO    string `json:"rpo"`
	MTTR   string `json:"mttr"`
}

// TechnicalControl represents a technical control
type TechnicalControl struct {
	ControlType    string   `json:"control_type"`
	Implementation []string `json:"implementation"`
	Validation     []string `json:"validation"`
	Monitoring     []string `json:"monitoring"`
}

// OperationalControl represents an operational control
type OperationalControl struct {
	ControlType      string   `json:"control_type"`
	Procedures       []string `json:"procedures"`
	Responsibilities []string `json:"responsibilities"`
	Documentation    []string `json:"documentation"`
}

// AuditRequirement represents an audit requirement
type AuditRequirement struct {
	RequirementType string   `json:"requirement_type"`
	Frequency       string   `json:"frequency"`
	Evidence        []string `json:"evidence"`
	Reporting       []string `json:"reporting"`
}

// CompliancePenalty represents a compliance penalty
type CompliancePenalty struct {
	ViolationType string `json:"violation_type"`
	Severity      string `json:"severity"`
	Penalty       string `json:"penalty"`
	Remediation   string `json:"remediation"`
}

// DataGovernanceFramework represents a data governance framework
type DataGovernanceFramework struct {
	Framework  string   `json:"framework"`
	Policies   []string `json:"policies"`
	Procedures []string `json:"procedures"`
	Controls   []string `json:"controls"`
	Monitoring []string `json:"monitoring"`
}

// RecommendedService represents a recommended service
type RecommendedService struct {
	ServiceName   string   `json:"service_name"`
	Provider      string   `json:"provider"`
	Category      string   `json:"category"`
	Justification string   `json:"justification"`
	Configuration []string `json:"configuration"`
}

// ImplementationPhase represents an implementation phase
type ImplementationPhase struct {
	PhaseNumber   int      `json:"phase_number"`
	PhaseName     string   `json:"phase_name"`
	Description   string   `json:"description"`
	Duration      string   `json:"duration"`
	Prerequisites []string `json:"prerequisites"`
	Deliverables  []string `json:"deliverables"`
}

// IndustryCostEstimation represents cost estimation for industry solutions
type IndustryCostEstimation struct {
	InitialCost float64            `json:"initial_cost"`
	MonthlyCost float64            `json:"monthly_cost"`
	AnnualCost  float64            `json:"annual_cost"`
	Breakdown   map[string]float64 `json:"breakdown"`
	Currency    string             `json:"currency"`
}

// IndustryRisk represents an industry-specific risk
type IndustryRisk struct {
	RiskType    string   `json:"risk_type"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	Probability string   `json:"probability"`
	Mitigation  []string `json:"mitigation"`
}

// IndustryCaseStudy represents an industry case study
type IndustryCaseStudy struct {
	Title      string   `json:"title"`
	Industry   string   `json:"industry"`
	Challenge  string   `json:"challenge"`
	Solution   string   `json:"solution"`
	Results    []string `json:"results"`
	Duration   string   `json:"duration"`
	Investment string   `json:"investment,omitempty"`
}

// Additional stub types for completeness
type WorkloadCharacteristic struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type ArchitecturePattern struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Components  []string `json:"components"`
	Benefits    []string `json:"benefits"`
	Limitations []string `json:"limitations"`
}

type WorkloadResourceRequirements struct {
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
	Network string `json:"network"`
}

type PerformanceProfile struct {
	Baseline    string   `json:"baseline"`
	Peak        string   `json:"peak"`
	Patterns    []string `json:"patterns"`
	Bottlenecks []string `json:"bottlenecks"`
}

type ScalingBehavior struct {
	TriggerMetrics  []string `json:"trigger_metrics"`
	ScaleUpPolicy   string   `json:"scale_up_policy"`
	ScaleDownPolicy string   `json:"scale_down_policy"`
	Limits          string   `json:"limits"`
}

type OptimizationTip struct {
	Category   string `json:"category"`
	Tip        string `json:"tip"`
	Impact     string `json:"impact"`
	Difficulty string `json:"difficulty"`
}

type ScalabilityNeeds struct {
	ExpectedGrowth string   `json:"expected_growth"`
	PeakLoads      string   `json:"peak_loads"`
	Seasonality    []string `json:"seasonality"`
	Constraints    []string `json:"constraints"`
}

// BudgetConstraints is defined in multicloud.go

type ExistingInfrastructure struct {
	OnPremise   []string `json:"on_premise"`
	Cloud       []string `json:"cloud"`
	Hybrid      []string `json:"hybrid"`
	Constraints []string `json:"constraints"`
}

type Application struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Criticality  string   `json:"criticality"`
	Dependencies []string `json:"dependencies"`
	DataAssets   []string `json:"data_assets"`
}

type DataAsset struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Size        string   `json:"size"`
	Sensitivity string   `json:"sensitivity"`
	Compliance  []string `json:"compliance"`
}

type MigrationStep struct {
	StepNumber  int      `json:"step_number"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	Resources   []string `json:"resources"`
}

// MigrationTool is defined in multicloud.go

type MigrationAssessment struct {
	Complexity   string   `json:"complexity"`
	Risks        []string `json:"risks"`
	Dependencies []string `json:"dependencies"`
	Readiness    string   `json:"readiness"`
}

type MigrationResourceRequirements struct {
	Personnel []string `json:"personnel"`
	Tools     []string `json:"tools"`
	Budget    string   `json:"budget"`
	Timeline  string   `json:"timeline"`
}

type MigrationTestingStrategy struct {
	TestTypes  []string `json:"test_types"`
	TestData   []string `json:"test_data"`
	Validation []string `json:"validation"`
	Rollback   []string `json:"rollback"`
}

type SuccessMetric struct {
	Metric      string `json:"metric"`
	Target      string `json:"target"`
	Measurement string `json:"measurement"`
	Frequency   string `json:"frequency"`
}

type PostMigrationTask struct {
	Task         string   `json:"task"`
	Description  string   `json:"description"`
	Owner        string   `json:"owner"`
	Timeline     string   `json:"timeline"`
	Dependencies []string `json:"dependencies"`
}

type MigrationComplexityAssessment struct {
	OverallComplexity string   `json:"overall_complexity"`
	TechnicalFactors  []string `json:"technical_factors"`
	BusinessFactors   []string `json:"business_factors"`
	RiskFactors       []string `json:"risk_factors"`
	Recommendations   []string `json:"recommendations"`
}

type DisasterScenario struct {
	ScenarioName string   `json:"scenario_name"`
	Description  string   `json:"description"`
	Impact       string   `json:"impact"`
	Probability  string   `json:"probability"`
	Response     []string `json:"response"`
}

type RecoveryStrategy struct {
	StrategyName  string   `json:"strategy_name"`
	Description   string   `json:"description"`
	Applicability []string `json:"applicability"`
	Steps         []string `json:"steps"`
	Resources     []string `json:"resources"`
}

type ReplicationStrategy struct {
	Type       string   `json:"type"`
	Frequency  string   `json:"frequency"`
	Targets    []string `json:"targets"`
	Validation []string `json:"validation"`
}

type FailoverProcedure struct {
	ProcedureName string   `json:"procedure_name"`
	TriggerEvents []string `json:"trigger_events"`
	Steps         []string `json:"steps"`
	Validation    []string `json:"validation"`
	Rollback      []string `json:"rollback"`
}

type DRTestingPlan struct {
	TestTypes       []string `json:"test_types"`
	Frequency       string   `json:"frequency"`
	Scenarios       []string `json:"scenarios"`
	SuccessCriteria []string `json:"success_criteria"`
}

type DRResourceRequirements struct {
	Personnel []string `json:"personnel"`
	Systems   []string `json:"systems"`
	Budget    string   `json:"budget"`
	Vendors   []string `json:"vendors"`
}

type DRCostAnalysis struct {
	InitialCost float64 `json:"initial_cost"`
	OngoingCost float64 `json:"ongoing_cost"`
	TestingCost float64 `json:"testing_cost"`
	Currency    string  `json:"currency"`
}

type DRMaintenancePlan struct {
	UpdateFrequency  string   `json:"update_frequency"`
	ReviewCycle      string   `json:"review_cycle"`
	TestSchedule     []string `json:"test_schedule"`
	Responsibilities []string `json:"responsibilities"`
}

type BusinessImpactAnalysis struct {
	CriticalFunctions  []string `json:"critical_functions"`
	ImpactAssessment   []string `json:"impact_assessment"`
	RecoveryPriorities []string `json:"recovery_priorities"`
	Dependencies       []string `json:"dependencies"`
}

type CriticalBusinessFunction struct {
	FunctionName string   `json:"function_name"`
	Description  string   `json:"description"`
	Criticality  string   `json:"criticality"`
	Dependencies []string `json:"dependencies"`
	RecoveryTime string   `json:"recovery_time"`
}

type BCPRecoveryStrategy struct {
	StrategyName string   `json:"strategy_name"`
	Description  string   `json:"description"`
	Functions    []string `json:"functions"`
	Resources    []string `json:"resources"`
	Timeline     string   `json:"timeline"`
}

type AlternateWorkArrangements struct {
	RemoteWork     []string `json:"remote_work"`
	AlternateSites []string `json:"alternate_sites"`
	Equipment      []string `json:"equipment"`
	Communications []string `json:"communications"`
}

type SupplierContinuity struct {
	CriticalSuppliers  []string `json:"critical_suppliers"`
	AlternateSuppliers []string `json:"alternate_suppliers"`
	ContractTerms      []string `json:"contract_terms"`
	MonitoringProcess  []string `json:"monitoring_process"`
}

type BCPCommunicationStrategy struct {
	InternalComms []string `json:"internal_comms"`
	ExternalComms []string `json:"external_comms"`
	MediaStrategy []string `json:"media_strategy"`
	Stakeholders  []string `json:"stakeholders"`
}

type TrainingAndAwareness struct {
	TrainingPrograms    []string `json:"training_programs"`
	AwarenessActivities []string `json:"awareness_activities"`
	Frequency           string   `json:"frequency"`
	Audience            []string `json:"audience"`
}

type BCPTestingPlan struct {
	TestTypes     []string `json:"test_types"`
	TestFrequency string   `json:"test_frequency"`
	Scenarios     []string `json:"scenarios"`
	Participants  []string `json:"participants"`
}

type BCPMaintenancePlan struct {
	ReviewCycle      string   `json:"review_cycle"`
	UpdateTriggers   []string `json:"update_triggers"`
	Responsibilities []string `json:"responsibilities"`
	Documentation    []string `json:"documentation"`
}

type BCPResourceRequirements struct {
	Personnel  []string `json:"personnel"`
	Facilities []string `json:"facilities"`
	Technology []string `json:"technology"`
	Budget     string   `json:"budget"`
}

type BCPCostAnalysis struct {
	PreventionCost float64 `json:"prevention_cost"`
	ResponseCost   float64 `json:"response_cost"`
	RecoveryCost   float64 `json:"recovery_cost"`
	Currency       string  `json:"currency"`
}

// WorkloadPerformanceOptimization represents performance optimization for workloads
type WorkloadPerformanceOptimization struct {
	Strategy            string   `json:"strategy"`
	CPUOptimization     []string `json:"cpu_optimization"`
	MemoryOptimization  []string `json:"memory_optimization"`
	StorageOptimization []string `json:"storage_optimization"`
	NetworkOptimization []string `json:"network_optimization"`
	ExpectedImprovement string   `json:"expected_improvement"`
}

// WorkloadCostOptimization represents cost optimization for workloads
type WorkloadCostOptimization struct {
	Strategy            string   `json:"strategy"`
	RightSizing         []string `json:"right_sizing"`
	ReservedInstances   []string `json:"reserved_instances"`
	SpotInstances       []string `json:"spot_instances"`
	StorageOptimization []string `json:"storage_optimization"`
	ExpectedSavings     string   `json:"expected_savings"`
}

// WorkloadScalabilityOptimization represents scalability optimization for workloads
type WorkloadScalabilityOptimization struct {
	Strategy              string   `json:"strategy"`
	HorizontalScaling     []string `json:"horizontal_scaling"`
	VerticalScaling       []string `json:"vertical_scaling"`
	AutoScalingPolicies   []string `json:"auto_scaling_policies"`
	LoadBalancingStrategy []string `json:"load_balancing_strategy"`
	ExpectedCapacity      string   `json:"expected_capacity"`
}

// WorkloadSecurityOptimization represents security optimization for workloads
type WorkloadSecurityOptimization struct {
	Strategy           string   `json:"strategy"`
	AccessControls     []string `json:"access_controls"`
	EncryptionStrategy []string `json:"encryption_strategy"`
	NetworkSecurity    []string `json:"network_security"`
	MonitoringControls []string `json:"monitoring_controls"`
	ComplianceControls []string `json:"compliance_controls"`
}

// WorkloadOptimizedArchitecture represents optimized architecture for workloads
type WorkloadOptimizedArchitecture struct {
	ArchitectureType    string                   `json:"architecture_type"`
	Components          []*ArchitectureComponent `json:"components"`
	DataFlow            []*DataFlowStep          `json:"data_flow"`
	SecurityLayers      []string                 `json:"security_layers"`
	MonitoringStrategy  []string                 `json:"monitoring_strategy"`
	ScalabilityFeatures []string                 `json:"scalability_features"`
}

// WorkloadServiceRecommendation represents service recommendations for workloads
type WorkloadServiceRecommendation struct {
	ServiceName       string   `json:"service_name"`
	ServiceType       string   `json:"service_type"`
	Justification     string   `json:"justification"`
	Configuration     []string `json:"configuration"`
	CostImpact        string   `json:"cost_impact"`
	PerformanceImpact string   `json:"performance_impact"`
}

// WorkloadConfigurationTuning represents configuration tuning for workloads
type WorkloadConfigurationTuning struct {
	Component         string   `json:"component"`
	CurrentConfig     []string `json:"current_config"`
	RecommendedConfig []string `json:"recommended_config"`
	Justification     string   `json:"justification"`
	Impact            string   `json:"impact"`
	RiskLevel         string   `json:"risk_level"`
}

// WorkloadMonitoringStrategy represents monitoring strategy for workloads
type WorkloadMonitoringStrategy struct {
	Strategy        string   `json:"strategy"`
	KeyMetrics      []string `json:"key_metrics"`
	AlertingRules   []string `json:"alerting_rules"`
	DashboardSetup  []string `json:"dashboard_setup"`
	LoggingStrategy []string `json:"logging_strategy"`
	Tools           []string `json:"tools"`
}

// WorkloadOptimizationPhase represents a phase in workload optimization
type WorkloadOptimizationPhase struct {
	PhaseNumber   int      `json:"phase_number"`
	PhaseName     string   `json:"phase_name"`
	Description   string   `json:"description"`
	Duration      string   `json:"duration"`
	Prerequisites []string `json:"prerequisites"`
	Deliverables  []string `json:"deliverables"`
	RiskLevel     string   `json:"risk_level"`
}

// WorkloadOptimizationBenefits represents expected benefits from workload optimization
type WorkloadOptimizationBenefits struct {
	PerformanceGain string `json:"performance_gain"`
	CostSavings     string `json:"cost_savings"`
	ScalabilityGain string `json:"scalability_gain"`
	SecurityGain    string `json:"security_gain"`
	OperationalGain string `json:"operational_gain"`
}

// WorkloadSpec is defined in multicloud.go

// WorkloadPerformanceAnalysis represents performance analysis results
type WorkloadPerformanceAnalysis struct {
	WorkloadName     string                 `json:"workload_name"`
	AnalysisDate     time.Time              `json:"analysis_date"`
	PerformanceScore float64                `json:"performance_score"`
	Bottlenecks      []string               `json:"bottlenecks"`
	Recommendations  []string               `json:"recommendations"`
	MetricAnalysis   map[string]interface{} `json:"metric_analysis"`
	ImprovementAreas []string               `json:"improvement_areas"`
}

// Architecture is defined in risk.go

// ArchitectureComponent is defined in architecture.go

// SecurityLayer is defined in risk.go

// EnvironmentSpec represents an environment specification
type EnvironmentSpec struct {
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Provider      string                 `json:"provider"`
	Region        string                 `json:"region"`
	Resources     []string               `json:"resources"`
	Configuration map[string]interface{} `json:"configuration"`
	Constraints   []string               `json:"constraints"`
}

// MigrationPhase is defined in multicloud.go

// MigrationRisk represents a migration risk
type MigrationRisk struct {
	RiskType    string   `json:"risk_type"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	Probability string   `json:"probability"`
	Mitigation  []string `json:"mitigation"`
	Owner       string   `json:"owner"`
}

// MigrationCostEstimation represents cost estimation for migration
type MigrationCostEstimation struct {
	AssessmentCost float64            `json:"assessment_cost"`
	MigrationCost  float64            `json:"migration_cost"`
	ValidationCost float64            `json:"validation_cost"`
	OngoingCost    float64            `json:"ongoing_cost"`
	Breakdown      map[string]float64 `json:"breakdown"`
	Currency       string             `json:"currency"`
}

// MigrationTimeline represents migration timeline
type MigrationTimeline struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Milestones   []string  `json:"milestones"`
	Dependencies []string  `json:"dependencies"`
	CriticalPath []string  `json:"critical_path"`
}

// RollbackPlan represents a rollback plan
type RollbackPlan struct {
	TriggerConditions []string `json:"trigger_conditions"`
	RollbackSteps     []string `json:"rollback_steps"`
	ValidationSteps   []string `json:"validation_steps"`
	Timeline          string   `json:"timeline"`
	Resources         []string `json:"resources"`
}

// BackupStrategy represents a backup strategy
type BackupStrategy struct {
	BackupTypes []string `json:"backup_types"`
	Frequency   string   `json:"frequency"`
	Retention   string   `json:"retention"`
	Storage     []string `json:"storage"`
	Validation  []string `json:"validation"`
}

// RecoveryProcedure represents a recovery procedure
type RecoveryProcedure struct {
	ProcedureName string   `json:"procedure_name"`
	Description   string   `json:"description"`
	Steps         []string `json:"steps"`
	Prerequisites []string `json:"prerequisites"`
	Validation    []string `json:"validation"`
	Timeline      string   `json:"timeline"`
}

// CommunicationPlan represents a communication plan
type CommunicationPlan struct {
	Stakeholders   []string `json:"stakeholders"`
	Channels       []string `json:"channels"`
	Frequency      string   `json:"frequency"`
	Templates      []string `json:"templates"`
	EscalationPath []string `json:"escalation_path"`
}

// IndustryArchitecturePattern represents architecture patterns specific to industries
type IndustryArchitecturePattern struct {
	PatternType         string                  `json:"pattern_type"`
	Components          []ArchitectureComponent `json:"components"`
	DataTiers           []DataTier              `json:"data_tiers"`
	SecurityZones       []SecurityZone          `json:"security_zones"`
	IntegrationPoints   []IntegrationPoint      `json:"integration_points"`
	ComplianceControls  []ComplianceControl     `json:"compliance_controls"`
	PerformanceTargets  *PerformanceTargets     `json:"performance_targets"`
	AvailabilityTargets *AvailabilityTargets    `json:"availability_targets"`
}

// IndustrySecurityControl represents security controls specific to industries
type IndustrySecurityControl struct {
	ControlID       string   `json:"control_id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Implementation  string   `json:"implementation"`
	CloudServices   []string `json:"cloud_services"`
	ConfigExamples  []string `json:"config_examples"`
	ValidationSteps []string `json:"validation_steps"`
	Frameworks      []string `json:"frameworks"`
}

// WorkloadComputeAnalysis represents compute analysis for workloads
type WorkloadComputeAnalysis struct {
	CPUUtilization    float64  `json:"cpu_utilization"`
	MemoryUtilization float64  `json:"memory_utilization"`
	InstanceTypes     []string `json:"instance_types"`
	Recommendations   []string `json:"recommendations"`
}

// StorageAnalysis represents storage analysis
type StorageAnalysis struct {
	StorageType     string   `json:"storage_type"`
	Capacity        int64    `json:"capacity"`
	Utilization     float64  `json:"utilization"`
	Performance     string   `json:"performance"`
	Recommendations []string `json:"recommendations"`
}

// WorkloadNetworkAnalysis represents network analysis for workloads
type WorkloadNetworkAnalysis struct {
	Bandwidth       float64  `json:"bandwidth"`
	Latency         float64  `json:"latency"`
	Throughput      float64  `json:"throughput"`
	SecurityGroups  []string `json:"security_groups"`
	Recommendations []string `json:"recommendations"`
}

// WorkloadScalingAnalysis represents scaling analysis for workloads
type WorkloadScalingAnalysis struct {
	CurrentCapacity int      `json:"current_capacity"`
	MaxCapacity     int      `json:"max_capacity"`
	ScalingTriggers []string `json:"scaling_triggers"`
	ScalingPolicies []string `json:"scaling_policies"`
	Recommendations []string `json:"recommendations"`
}
