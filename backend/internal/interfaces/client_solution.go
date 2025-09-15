package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// ClientSpecificSolutionEngine defines the interface for generating client-specific solutions
type ClientSpecificSolutionEngine interface {
	// Industry-specific solution patterns
	GenerateIndustrySolution(ctx context.Context, inquiry *domain.Inquiry, industry string) (*IndustrySolution, error)
	GetIndustryPatterns(ctx context.Context, industry string) ([]*IndustryPattern, error)
	GetComplianceRequirements(ctx context.Context, industry string) ([]*IndustryComplianceRequirement, error)

	// Workload-specific optimization recommendations
	GenerateWorkloadOptimization(ctx context.Context, workloadType string, requirements *WorkloadRequirements) (*WorkloadOptimization, error)
	GetWorkloadPatterns(ctx context.Context, workloadType string) ([]*WorkloadPattern, error)
	AnalyzeWorkloadPerformance(ctx context.Context, workload *WorkloadSpec) (*WorkloadPerformanceAnalysis, error)

	// Migration strategy generator
	GenerateMigrationStrategy(ctx context.Context, migrationRequest *MigrationRequest) (*MigrationStrategy, error)
	GetMigrationPatterns(ctx context.Context, sourceType, targetType string) ([]*MigrationPattern, error)
	EstimateMigrationComplexity(ctx context.Context, migrationRequest *MigrationRequest) (*MigrationComplexityAssessment, error)

	// Disaster recovery and business continuity
	GenerateDisasterRecoveryPlan(ctx context.Context, requirements *DRRequirements) (*DisasterRecoveryPlan, error)
	GenerateBusinessContinuityPlan(ctx context.Context, requirements *BCPRequirements) (*BusinessContinuityPlan, error)
	AssessRPORTO(ctx context.Context, architecture *Architecture) (*RPORTOAssessment, error)
}

// IndustrySolution represents a solution tailored to a specific industry
type IndustrySolution struct {
	ID                   string                           `json:"id"`
	Industry             string                           `json:"industry"`
	InquiryID            string                           `json:"inquiry_id"`
	SolutionName         string                           `json:"solution_name"`
	Description          string                           `json:"description"`
	ArchitecturePattern  *IndustryArchitecturePattern     `json:"architecture_pattern"`
	ComplianceFrameworks []*IndustryComplianceRequirement `json:"compliance_frameworks"`
	SecurityControls     []*IndustrySecurityControl       `json:"security_controls"`
	DataGovernance       *DataGovernanceFramework         `json:"data_governance"`
	RecommendedServices  []*RecommendedService            `json:"recommended_services"`
	ImplementationPhases []*ImplementationPhase           `json:"implementation_phases"`
	CostEstimation       *IndustryCostEstimation          `json:"cost_estimation"`
	RiskConsiderations   []*IndustryRisk                  `json:"risk_considerations"`
	BestPractices        []string                         `json:"best_practices"`
	CaseStudies          []*IndustryCaseStudy             `json:"case_studies"`
	CreatedAt            time.Time                        `json:"created_at"`
}

// IndustryPattern represents common patterns for specific industries
type IndustryPattern struct {
	ID                  string                 `json:"id"`
	Industry            string                 `json:"industry"`
	PatternName         string                 `json:"pattern_name"`
	Description         string                 `json:"description"`
	UseCases            []string               `json:"use_cases"`
	ArchitecturalLayers []ArchitecturalLayer   `json:"architectural_layers"`
	DataFlow            []DataFlowStep         `json:"data_flow"`
	SecurityLayers      []string               `json:"security_layers"`    // Simplified to avoid conflicts
	ComplianceMapping   map[string][]string    `json:"compliance_mapping"` // framework -> requirements
	TechnicalStack      *TechnicalStack        `json:"technical_stack"`
	ScalabilityFactors  []ScalabilityFactor    `json:"scalability_factors"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// IndustryComplianceRequirement represents compliance requirements for specific industries
type IndustryComplianceRequirement struct {
	Framework           string                 `json:"framework"`
	Industry            string                 `json:"industry"`
	RequirementID       string                 `json:"requirement_id"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	Category            string                 `json:"category"`
	Severity            string                 `json:"severity"`
	TechnicalControls   []TechnicalControl     `json:"technical_controls"`
	OperationalControls []OperationalControl   `json:"operational_controls"`
	AuditRequirements   []AuditRequirement     `json:"audit_requirements"`
	Penalties           []CompliancePenalty    `json:"penalties"`
	CloudMappings       map[string][]string    `json:"cloud_mappings"` // provider -> services
	Metadata            map[string]interface{} `json:"metadata"`
}

// WorkloadRequirements represents requirements for workload optimization
type WorkloadRequirements struct {
	WorkloadType           string                  `json:"workload_type"`
	PerformanceTargets     *PerformanceTargets     `json:"performance_targets"`
	ScalabilityNeeds       *ScalabilityNeeds       `json:"scalability_needs"`
	BudgetConstraints      *BudgetConstraints      `json:"budget_constraints"`
	ComplianceNeeds        []string                `json:"compliance_needs"`
	ExistingInfrastructure *ExistingInfrastructure `json:"existing_infrastructure"`
	BusinessRequirements   []string                `json:"business_requirements"`
	TechnicalConstraints   []string                `json:"technical_constraints"`
	Metadata               map[string]interface{}  `json:"metadata"`
}

// WorkloadOptimization represents optimization recommendations for specific workloads
type WorkloadOptimization struct {
	ID                      string                           `json:"id"`
	WorkloadType            string                           `json:"workload_type"`
	OptimizationStrategy    string                           `json:"optimization_strategy"`
	PerformanceOptimization *WorkloadPerformanceOptimization `json:"performance_optimization"`
	CostOptimization        *WorkloadCostOptimization        `json:"cost_optimization"`
	ScalabilityOptimization *WorkloadScalabilityOptimization `json:"scalability_optimization"`
	SecurityOptimization    *WorkloadSecurityOptimization    `json:"security_optimization"`
	RecommendedArchitecture *WorkloadOptimizedArchitecture   `json:"recommended_architecture"`
	ServiceRecommendations  []*WorkloadServiceRecommendation `json:"service_recommendations"`
	ConfigurationTuning     []*WorkloadConfigurationTuning   `json:"configuration_tuning"`
	MonitoringStrategy      *WorkloadMonitoringStrategy      `json:"monitoring_strategy"`
	ImplementationPlan      []*WorkloadOptimizationPhase     `json:"implementation_plan"`
	ExpectedBenefits        *WorkloadOptimizationBenefits    `json:"expected_benefits"`
	CreatedAt               time.Time                        `json:"created_at"`
}

// WorkloadPattern represents common patterns for specific workload types
type WorkloadPattern struct {
	ID                   string                        `json:"id"`
	WorkloadType         string                        `json:"workload_type"`
	PatternName          string                        `json:"pattern_name"`
	Description          string                        `json:"description"`
	Characteristics      []WorkloadCharacteristic      `json:"characteristics"`
	ArchitecturePattern  *ArchitecturePattern          `json:"architecture_pattern"`
	ResourceRequirements *WorkloadResourceRequirements `json:"resource_requirements"`
	PerformanceProfile   *PerformanceProfile           `json:"performance_profile"`
	ScalingBehavior      *ScalingBehavior              `json:"scaling_behavior"`
	OptimizationTips     []OptimizationTip             `json:"optimization_tips"`
	CommonChallenges     []string                      `json:"common_challenges"`
	BestPractices        []string                      `json:"best_practices"`
	Metadata             map[string]interface{}        `json:"metadata"`
}

// MigrationRequest represents a request for migration strategy
type MigrationRequest struct {
	SourceType           string                 `json:"source_type"` // "on-premises", "cloud", "hybrid"
	TargetType           string                 `json:"target_type"` // "aws", "azure", "gcp", "multi-cloud"
	ApplicationPortfolio []Application          `json:"application_portfolio"`
	DataAssets           []DataAsset            `json:"data_assets"`
	BusinessDrivers      []string               `json:"business_drivers"`
	Constraints          []string               `json:"constraints"`
	Timeline             string                 `json:"timeline"`
	Budget               float64                `json:"budget"`
	RiskTolerance        string                 `json:"risk_tolerance"`
	ComplianceNeeds      []string               `json:"compliance_needs"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// MigrationStrategy represents a comprehensive migration strategy
type MigrationStrategy struct {
	ID                   string                         `json:"id"`
	MigrationName        string                         `json:"migration_name"`
	SourceEnvironment    *EnvironmentSpec               `json:"source_environment"`
	TargetEnvironment    *EnvironmentSpec               `json:"target_environment"`
	MigrationType        string                         `json:"migration_type"`     // "lift-and-shift", "re-platform", "re-architect", "hybrid"
	MigrationApproach    string                         `json:"migration_approach"` // "big-bang", "phased", "parallel-run"
	AssessmentResults    *MigrationAssessment           `json:"assessment_results"`
	MigrationPhases      []*MigrationPhase              `json:"migration_phases"`
	RiskMitigation       []*MigrationRisk               `json:"risk_mitigation"`
	ResourceRequirements *MigrationResourceRequirements `json:"resource_requirements"`
	Timeline             *MigrationTimeline             `json:"timeline"`
	CostEstimation       *MigrationCostEstimation       `json:"cost_estimation"`
	TestingStrategy      *MigrationTestingStrategy      `json:"testing_strategy"`
	RollbackPlan         *RollbackPlan                  `json:"rollback_plan"`
	SuccessMetrics       []SuccessMetric                `json:"success_metrics"`
	PostMigrationTasks   []PostMigrationTask            `json:"post_migration_tasks"`
	CreatedAt            time.Time                      `json:"created_at"`
}

// MigrationPattern represents common migration patterns
type MigrationPattern struct {
	ID               string                 `json:"id"`
	PatternName      string                 `json:"pattern_name"`
	SourceType       string                 `json:"source_type"`
	TargetType       string                 `json:"target_type"`
	Description      string                 `json:"description"`
	Applicability    []string               `json:"applicability"`
	Steps            []MigrationStep        `json:"steps"`
	Tools            []MigrationTool        `json:"tools"`
	Duration         string                 `json:"duration"`
	Complexity       string                 `json:"complexity"`
	RiskLevel        string                 `json:"risk_level"`
	SuccessFactors   []string               `json:"success_factors"`
	CommonPitfalls   []string               `json:"common_pitfalls"`
	CostImplications []string               `json:"cost_implications"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// DRRequirements represents disaster recovery requirements
type DRRequirements struct {
	BusinessCriticality string                 `json:"business_criticality"`
	RPORequirements     map[string]string      `json:"rpo_requirements"`
	RTORequirements     map[string]string      `json:"rto_requirements"`
	ComplianceNeeds     []string               `json:"compliance_needs"`
	BudgetConstraints   *BudgetConstraints     `json:"budget_constraints"`
	GeographicScope     []string               `json:"geographic_scope"`
	DataSensitivity     string                 `json:"data_sensitivity"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// DisasterRecoveryPlan represents a comprehensive disaster recovery plan
type DisasterRecoveryPlan struct {
	ID                   string                  `json:"id"`
	PlanName             string                  `json:"plan_name"`
	Scope                string                  `json:"scope"`
	RPOTargets           map[string]string       `json:"rpo_targets"` // service -> RPO
	RTOTargets           map[string]string       `json:"rto_targets"` // service -> RTO
	DisasterScenarios    []*DisasterScenario     `json:"disaster_scenarios"`
	RecoveryStrategies   []*RecoveryStrategy     `json:"recovery_strategies"`
	BackupStrategy       *BackupStrategy         `json:"backup_strategy"`
	ReplicationStrategy  *ReplicationStrategy    `json:"replication_strategy"`
	FailoverProcedures   []*FailoverProcedure    `json:"failover_procedures"`
	RecoveryProcedures   []*RecoveryProcedure    `json:"recovery_procedures"`
	TestingPlan          *DRTestingPlan          `json:"testing_plan"`
	CommunicationPlan    *CommunicationPlan      `json:"communication_plan"`
	ResourceRequirements *DRResourceRequirements `json:"resource_requirements"`
	CostAnalysis         *DRCostAnalysis         `json:"cost_analysis"`
	ComplianceMapping    map[string][]string     `json:"compliance_mapping"`
	MaintenancePlan      *DRMaintenancePlan      `json:"maintenance_plan"`
	CreatedAt            time.Time               `json:"created_at"`
}

// BCPRequirements represents business continuity planning requirements
type BCPRequirements struct {
	BusinessFunctions  []string               `json:"business_functions"`
	CriticalityLevels  map[string]string      `json:"criticality_levels"`
	RecoveryObjectives map[string]string      `json:"recovery_objectives"`
	StakeholderGroups  []string               `json:"stakeholder_groups"`
	RegulatoryNeeds    []string               `json:"regulatory_needs"`
	BudgetConstraints  *BudgetConstraints     `json:"budget_constraints"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// BusinessContinuityPlan represents a comprehensive business continuity plan
type BusinessContinuityPlan struct {
	ID                        string                      `json:"id"`
	PlanName                  string                      `json:"plan_name"`
	BusinessImpactAnalysis    *BusinessImpactAnalysis     `json:"business_impact_analysis"`
	CriticalBusinessFunctions []*CriticalBusinessFunction `json:"critical_business_functions"`
	RecoveryStrategies        []*BCPRecoveryStrategy      `json:"recovery_strategies"`
	AlternateWorkArrangements *AlternateWorkArrangements  `json:"alternate_work_arrangements"`
	SupplierContinuity        *SupplierContinuity         `json:"supplier_continuity"`
	CommunicationStrategy     *BCPCommunicationStrategy   `json:"communication_strategy"`
	TrainingAndAwareness      *TrainingAndAwareness       `json:"training_and_awareness"`
	TestingAndExercises       *BCPTestingPlan             `json:"testing_and_exercises"`
	MaintenanceAndReview      *BCPMaintenancePlan         `json:"maintenance_and_review"`
	ResourceRequirements      *BCPResourceRequirements    `json:"resource_requirements"`
	CostAnalysis              *BCPCostAnalysis            `json:"cost_analysis"`
	CreatedAt                 time.Time                   `json:"created_at"`
}

// RPORTOAssessment represents RPO/RTO assessment results
type RPORTOAssessment struct {
	ServiceAssessments []ServiceRPORTO `json:"service_assessments"`
	OverallRPO         string          `json:"overall_rpo"`
	OverallRTO         string          `json:"overall_rto"`
	Recommendations    []string        `json:"recommendations"`
	CostImplications   []string        `json:"cost_implications"`
	RiskFactors        []string        `json:"risk_factors"`
}

// ServiceRPORTO represents RPO/RTO for a specific service
type ServiceRPORTO struct {
	ServiceName     string   `json:"service_name"`
	CurrentRPO      string   `json:"current_rpo"`
	CurrentRTO      string   `json:"current_rto"`
	TargetRPO       string   `json:"target_rpo"`
	TargetRTO       string   `json:"target_rto"`
	GapAnalysis     string   `json:"gap_analysis"`
	Recommendations []string `json:"recommendations"`
}
