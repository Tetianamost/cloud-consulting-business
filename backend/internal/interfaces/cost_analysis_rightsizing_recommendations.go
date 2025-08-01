package interfaces

import (
	"time"
)

// RightSizingRecommendations represents right-sizing recommendations
type RightSizingRecommendations struct {
	ID                    string                         `json:"id"`
	AnalysisID            string                         `json:"analysis_id"`
	TotalSavingsPotential float64                        `json:"total_savings_potential"`
	Recommendations       []*RightSizingRecommendation   `json:"recommendations"`
	ImplementationPlan    *RightSizingImplementationPlan `json:"implementation_plan"`
	RiskAssessment        *RightSizingRiskAssessment     `json:"risk_assessment"`
	ValidationPlan        *RightSizingValidationPlan     `json:"validation_plan"`
	MonitoringPlan        *RightSizingMonitoringPlan     `json:"monitoring_plan"`
	CreatedAt             time.Time                      `json:"created_at"`
}

// RightSizingRecommendation represents a specific right-sizing recommendation
type RightSizingRecommendation struct {
	ID                       string                           `json:"id"`
	ResourceID               string                           `json:"resource_id"`
	ResourceType             string                           `json:"resource_type"`
	RecommendationType       string                           `json:"recommendation_type"` // "downsize", "upsize", "change_type", "terminate"
	Priority                 string                           `json:"priority"`
	CurrentConfiguration     *ResourceConfiguration           `json:"current_configuration"`
	RecommendedConfiguration *ResourceConfiguration           `json:"recommended_configuration"`
	SavingsPotential         float64                          `json:"savings_potential"`
	SavingsPercentage        float64                          `json:"savings_percentage"`
	PerformanceImpact        string                           `json:"performance_impact"`
	RiskLevel                string                           `json:"risk_level"`
	ConfidenceLevel          float64                          `json:"confidence_level"`
	ImplementationEffort     string                           `json:"implementation_effort"`
	ImplementationTime       string                           `json:"implementation_time"`
	Prerequisites            []string                         `json:"prerequisites"`
	ImplementationSteps      []*RightSizingImplementationStep `json:"implementation_steps"`
	ValidationCriteria       []string                         `json:"validation_criteria"`
	RollbackPlan             []string                         `json:"rollback_plan"`
	MonitoringMetrics        []string                         `json:"monitoring_metrics"`
	BusinessJustification    string                           `json:"business_justification"`
	TechnicalJustification   string                           `json:"technical_justification"`
	AlternativeOptions       []*RightSizingAlternativeOption  `json:"alternative_options"`
	Dependencies             []string                         `json:"dependencies"`
	Constraints              []string                         `json:"constraints"`
}

// RightSizingImplementationStep represents an implementation step
type RightSizingImplementationStep struct {
	StepNumber    int      `json:"step_number"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Duration      string   `json:"duration"`
	Owner         string   `json:"owner"`
	Prerequisites []string `json:"prerequisites"`
	Tools         []string `json:"tools"`
	Commands      []string `json:"commands"`
	Validation    []string `json:"validation"`
	RollbackSteps []string `json:"rollback_steps"`
	RiskLevel     string   `json:"risk_level"`
}

// RightSizingAlternativeOption represents alternative right-sizing options
type RightSizingAlternativeOption struct {
	OptionName        string                 `json:"option_name"`
	Configuration     *ResourceConfiguration `json:"configuration"`
	SavingsPotential  float64                `json:"savings_potential"`
	PerformanceImpact string                 `json:"performance_impact"`
	RiskLevel         string                 `json:"risk_level"`
	Pros              []string               `json:"pros"`
	Cons              []string               `json:"cons"`
	Recommendation    string                 `json:"recommendation"`
}

// RightSizingImplementationPlan represents the implementation plan
type RightSizingImplementationPlan struct {
	TotalDuration        string                             `json:"total_duration"`
	TotalSavings         float64                            `json:"total_savings"`
	ImplementationPhases []*RightSizingImplementationPhase  `json:"implementation_phases"`
	ResourceRequirements *RightSizingResourceRequirements   `json:"resource_requirements"`
	Timeline             *RightSizingImplementationTimeline `json:"timeline"`
	Dependencies         []string                           `json:"dependencies"`
	CriticalPath         []string                           `json:"critical_path"`
	RiskMitigation       []string                           `json:"risk_mitigation"`
	SuccessMetrics       []string                           `json:"success_metrics"`
	CommunicationPlan    *RightSizingCommunicationPlan      `json:"communication_plan"`
}

// RightSizingImplementationPhase represents an implementation phase
type RightSizingImplementationPhase struct {
	PhaseName       string   `json:"phase_name"`
	PhaseNumber     int      `json:"phase_number"`
	Duration        string   `json:"duration"`
	Objectives      []string `json:"objectives"`
	Recommendations []string `json:"recommendations"` // IDs of recommendations in this phase
	Prerequisites   []string `json:"prerequisites"`
	Deliverables    []string `json:"deliverables"`
	SuccessCriteria []string `json:"success_criteria"`
	RiskLevel       string   `json:"risk_level"`
	ExpectedSavings float64  `json:"expected_savings"`
}

// RightSizingResourceRequirements represents resource requirements
type RightSizingResourceRequirements struct {
	TeamSize          int      `json:"team_size"`
	SkillsRequired    []string `json:"skills_required"`
	ToolsRequired     []string `json:"tools_required"`
	BudgetRequired    float64  `json:"budget_required"`
	TimeCommitment    string   `json:"time_commitment"`
	ExternalResources []string `json:"external_resources"`
}

// RightSizingImplementationTimeline represents implementation timeline
type RightSizingImplementationTimeline struct {
	StartDate     time.Time                             `json:"start_date"`
	EndDate       time.Time                             `json:"end_date"`
	Milestones    []*RightSizingImplementationMilestone `json:"milestones"`
	CriticalDates []*CriticalDate                       `json:"critical_dates"`
	BufferTime    string                                `json:"buffer_time"`
	ReviewPoints  []*ReviewPoint                        `json:"review_points"`
}

// RightSizingImplementationMilestone represents an implementation milestone
type RightSizingImplementationMilestone struct {
	MilestoneName   string    `json:"milestone_name"`
	TargetDate      time.Time `json:"target_date"`
	Description     string    `json:"description"`
	Deliverables    []string  `json:"deliverables"`
	SuccessCriteria []string  `json:"success_criteria"`
	Dependencies    []string  `json:"dependencies"`
	ExpectedSavings float64   `json:"expected_savings"`
}

// CriticalDate represents a critical date in the timeline
type CriticalDate struct {
	Date        time.Time `json:"date"`
	Event       string    `json:"event"`
	Impact      string    `json:"impact"`
	Preparation []string  `json:"preparation"`
}

// ReviewPoint represents a review point in the timeline
type ReviewPoint struct {
	Date           time.Time `json:"date"`
	ReviewType     string    `json:"review_type"`
	Participants   []string  `json:"participants"`
	ReviewCriteria []string  `json:"review_criteria"`
	DecisionPoints []string  `json:"decision_points"`
}

// RightSizingCommunicationPlan represents the communication plan
type RightSizingCommunicationPlan struct {
	Stakeholders          []*Stakeholder           `json:"stakeholders"`
	CommunicationChannels []string                 `json:"communication_channels"`
	ReportingSchedule     *ReportingSchedule       `json:"reporting_schedule"`
	EscalationMatrix      *EscalationMatrix        `json:"escalation_matrix"`
	Templates             []*CommunicationTemplate `json:"templates"`
}

// Stakeholder represents a stakeholder
type Stakeholder struct {
	Name          string   `json:"name"`
	Role          string   `json:"role"`
	Interests     []string `json:"interests"`
	Influence     string   `json:"influence"`
	Communication string   `json:"communication_preference"`
}

// ReportingSchedule represents the reporting schedule
type ReportingSchedule struct {
	Frequency  string   `json:"frequency"`
	Recipients []string `json:"recipients"`
	ReportType string   `json:"report_type"`
	Metrics    []string `json:"metrics"`
	Format     string   `json:"format"`
}

// EscalationMatrix represents the escalation matrix
type EscalationMatrix struct {
	Levels      []*EscalationLevel `json:"levels"`
	Triggers    []string           `json:"triggers"`
	Procedures  []string           `json:"procedures"`
	ContactList []string           `json:"contact_list"`
}

// EscalationLevel represents an escalation level
type EscalationLevel struct {
	Level     int      `json:"level"`
	Title     string   `json:"title"`
	Contacts  []string `json:"contacts"`
	Timeframe string   `json:"timeframe"`
	Authority []string `json:"authority"`
}

// CommunicationTemplate represents a communication template
type CommunicationTemplate struct {
	TemplateName string `json:"template_name"`
	Purpose      string `json:"purpose"`
	Audience     string `json:"audience"`
	Format       string `json:"format"`
	Content      string `json:"content"`
}

// RightSizingRiskAssessment represents risk assessment for right-sizing
type RightSizingRiskAssessment struct {
	OverallRiskLevel     string                           `json:"overall_risk_level"`
	RiskScore            float64                          `json:"risk_score"`
	Risks                []*RightSizingRisk               `json:"risks"`
	MitigationStrategies []*RightSizingMitigationStrategy `json:"mitigation_strategies"`
	ContingencyPlans     []*RightSizingContingencyPlan    `json:"contingency_plans"`
	RiskMonitoring       *RightSizingRiskMonitoring       `json:"risk_monitoring"`
	AcceptanceCriteria   []string                         `json:"acceptance_criteria"`
}

// RightSizingRisk represents a right-sizing risk
type RightSizingRisk struct {
	RiskID            string   `json:"risk_id"`
	RiskName          string   `json:"risk_name"`
	Category          string   `json:"category"`
	Description       string   `json:"description"`
	Impact            string   `json:"impact"`
	Probability       string   `json:"probability"`
	RiskScore         float64  `json:"risk_score"`
	AffectedResources []string `json:"affected_resources"`
	Triggers          []string `json:"triggers"`
	Consequences      []string `json:"consequences"`
	Owner             string   `json:"owner"`
	Status            string   `json:"status"`
}

// RightSizingMitigationStrategy represents a mitigation strategy
type RightSizingMitigationStrategy struct {
	StrategyID          string   `json:"strategy_id"`
	RiskID              string   `json:"risk_id"`
	StrategyName        string   `json:"strategy_name"`
	Description         string   `json:"description"`
	ImplementationSteps []string `json:"implementation_steps"`
	Resources           []string `json:"resources"`
	Timeline            string   `json:"timeline"`
	Effectiveness       string   `json:"effectiveness"`
	Cost                float64  `json:"cost"`
	Owner               string   `json:"owner"`
}

// RightSizingContingencyPlan represents a contingency plan
type RightSizingContingencyPlan struct {
	PlanID          string                       `json:"plan_id"`
	PlanName        string                       `json:"plan_name"`
	TriggerScenario string                       `json:"trigger_scenario"`
	ResponseActions []*ContingencyResponseAction `json:"response_actions"`
	Resources       []string                     `json:"resources"`
	Timeline        string                       `json:"timeline"`
	SuccessCriteria []string                     `json:"success_criteria"`
	Owner           string                       `json:"owner"`
}

// ContingencyResponseAction represents a contingency response action
type ContingencyResponseAction struct {
	ActionID     string   `json:"action_id"`
	Action       string   `json:"action"`
	Priority     int      `json:"priority"`
	Owner        string   `json:"owner"`
	Timeline     string   `json:"timeline"`
	Resources    []string `json:"resources"`
	Dependencies []string `json:"dependencies"`
}

// RightSizingRiskMonitoring represents risk monitoring
type RightSizingRiskMonitoring struct {
	MonitoringFrequency string           `json:"monitoring_frequency"`
	KeyIndicators       []string         `json:"key_indicators"`
	Thresholds          []*RiskThreshold `json:"thresholds"`
	AlertingRules       []*AlertingRule  `json:"alerting_rules"`
	ReportingSchedule   string           `json:"reporting_schedule"`
	ReviewSchedule      string           `json:"review_schedule"`
}

// RiskThreshold represents a risk threshold
type RiskThreshold struct {
	Metric    string  `json:"metric"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
	Severity  string  `json:"severity"`
	Action    string  `json:"action"`
}

// AlertingRule represents an alerting rule
type AlertingRule struct {
	RuleID     string   `json:"rule_id"`
	RuleName   string   `json:"rule_name"`
	Condition  string   `json:"condition"`
	Severity   string   `json:"severity"`
	Recipients []string `json:"recipients"`
	Message    string   `json:"message"`
	Frequency  string   `json:"frequency"`
}

// RightSizingValidationPlan represents the validation plan
type RightSizingValidationPlan struct {
	ValidationStrategy  string                        `json:"validation_strategy"`
	ValidationPhases    []*RightSizingValidationPhase `json:"validation_phases"`
	TestingPlan         *RightSizingTestingPlan       `json:"testing_plan"`
	PerformanceBaseline *PerformanceBaseline          `json:"performance_baseline"`
	ValidationMetrics   []string                      `json:"validation_metrics"`
	AcceptanceCriteria  []string                      `json:"acceptance_criteria"`
	ValidationTimeline  string                        `json:"validation_timeline"`
}

// RightSizingValidationPhase represents a validation phase
type RightSizingValidationPhase struct {
	PhaseName       string   `json:"phase_name"`
	PhaseNumber     int      `json:"phase_number"`
	Duration        string   `json:"duration"`
	Objectives      []string `json:"objectives"`
	TestScenarios   []string `json:"test_scenarios"`
	SuccessCriteria []string `json:"success_criteria"`
	Tools           []string `json:"tools"`
	Deliverables    []string `json:"deliverables"`
}

// RightSizingTestingPlan represents the testing plan
type RightSizingTestingPlan struct {
	TestingStrategy   string                 `json:"testing_strategy"`
	TestTypes         []string               `json:"test_types"`
	TestEnvironments  []string               `json:"test_environments"`
	TestScenarios     []*TestScenario        `json:"test_scenarios"`
	LoadTestingPlan   *LoadTestingPlan       `json:"load_testing_plan"`
	FailoverTesting   *FailoverTestingPlan   `json:"failover_testing"`
	RegressionTesting *RegressionTestingPlan `json:"regression_testing"`
}

// TestScenario represents a test scenario
type TestScenario struct {
	ScenarioID      string   `json:"scenario_id"`
	ScenarioName    string   `json:"scenario_name"`
	Description     string   `json:"description"`
	TestSteps       []string `json:"test_steps"`
	ExpectedResults []string `json:"expected_results"`
	PassCriteria    []string `json:"pass_criteria"`
	Tools           []string `json:"tools"`
	Duration        string   `json:"duration"`
}

// LoadTestingPlan represents load testing plan
type LoadTestingPlan struct {
	Strategy     string   `json:"strategy"`
	LoadProfiles []string `json:"load_profiles"`
	TestDuration string   `json:"test_duration"`
	Metrics      []string `json:"metrics"`
	Tools        []string `json:"tools"`
	PassCriteria []string `json:"pass_criteria"`
}

// FailoverTestingPlan represents failover testing plan
type FailoverTestingPlan struct {
	Strategy      string   `json:"strategy"`
	FailoverTypes []string `json:"failover_types"`
	TestScenarios []string `json:"test_scenarios"`
	RecoveryTime  string   `json:"recovery_time"`
	PassCriteria  []string `json:"pass_criteria"`
}

// RegressionTestingPlan represents regression testing plan
type RegressionTestingPlan struct {
	Strategy     string   `json:"strategy"`
	TestSuite    []string `json:"test_suite"`
	Automation   string   `json:"automation"`
	Frequency    string   `json:"frequency"`
	PassCriteria []string `json:"pass_criteria"`
}

// PerformanceBaseline represents performance baseline
type PerformanceBaseline struct {
	BaselineDate       time.Time         `json:"baseline_date"`
	BaselineMetrics    []*BaselineMetric `json:"baseline_metrics"`
	BaselineConditions string            `json:"baseline_conditions"`
	MeasurementPeriod  string            `json:"measurement_period"`
	Tools              []string          `json:"tools"`
}

// BaselineMetric represents a baseline metric
type BaselineMetric struct {
	MetricName  string  `json:"metric_name"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Tolerance   float64 `json:"tolerance"`
	Threshold   float64 `json:"threshold"`
	Criticality string  `json:"criticality"`
}

// RightSizingMonitoringPlan represents the monitoring plan
type RightSizingMonitoringPlan struct {
	MonitoringStrategy     string                        `json:"monitoring_strategy"`
	MonitoringPhases       []*RightSizingMonitoringPhase `json:"monitoring_phases"`
	MonitoringMetrics      []*MonitoringMetric           `json:"monitoring_metrics"`
	AlertingConfiguration  *AlertingConfiguration        `json:"alerting_configuration"`
	DashboardConfiguration *DashboardConfiguration       `json:"dashboard_configuration"`
	ReportingConfiguration *ReportingConfiguration       `json:"reporting_configuration"`
	MonitoringTools        []string                      `json:"monitoring_tools"`
}

// RightSizingMonitoringPhase represents a monitoring phase
type RightSizingMonitoringPhase struct {
	PhaseName       string   `json:"phase_name"`
	Duration        string   `json:"duration"`
	Objectives      []string `json:"objectives"`
	MonitoringLevel string   `json:"monitoring_level"`
	Metrics         []string `json:"metrics"`
	Frequency       string   `json:"frequency"`
	Thresholds      []string `json:"thresholds"`
}

// MonitoringMetric represents a monitoring metric
type MonitoringMetric struct {
	MetricName  string  `json:"metric_name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Frequency   string  `json:"frequency"`
	Threshold   float64 `json:"threshold"`
	AlertLevel  string  `json:"alert_level"`
	Criticality string  `json:"criticality"`
}

// AlertingConfiguration represents alerting configuration
type AlertingConfiguration struct {
	AlertingRules        []*AlertingRule   `json:"alerting_rules"`
	NotificationChannels []string          `json:"notification_channels"`
	EscalationRules      []*EscalationRule `json:"escalation_rules"`
	AlertSuppression     *AlertSuppression `json:"alert_suppression"`
}

// EscalationRule represents an escalation rule
type EscalationRule struct {
	RuleID          string   `json:"rule_id"`
	Condition       string   `json:"condition"`
	EscalationLevel int      `json:"escalation_level"`
	Delay           string   `json:"delay"`
	Recipients      []string `json:"recipients"`
}

// AlertSuppression represents alert suppression configuration
type AlertSuppression struct {
	SuppressionRules   []*SuppressionRule `json:"suppression_rules"`
	MaintenanceWindows []string           `json:"maintenance_windows"`
}

// SuppressionRule represents a suppression rule
type SuppressionRule struct {
	RuleID    string `json:"rule_id"`
	Condition string `json:"condition"`
	Duration  string `json:"duration"`
	Reason    string `json:"reason"`
}

// DashboardConfiguration represents dashboard configuration
type DashboardConfiguration struct {
	DashboardName   string             `json:"dashboard_name"`
	Widgets         []*DashboardWidget `json:"widgets"`
	RefreshInterval string             `json:"refresh_interval"`
	AccessControl   []string           `json:"access_control"`
}

// DashboardWidget represents a dashboard widget
type DashboardWidget struct {
	WidgetID   string   `json:"widget_id"`
	WidgetType string   `json:"widget_type"`
	Title      string   `json:"title"`
	Metrics    []string `json:"metrics"`
	TimeRange  string   `json:"time_range"`
	Filters    []string `json:"filters"`
}

// ReportingConfiguration represents reporting configuration
type ReportingConfiguration struct {
	ReportTypes []string         `json:"report_types"`
	Schedule    string           `json:"schedule"`
	Recipients  []string         `json:"recipients"`
	Format      string           `json:"format"`
	Content     []*ReportContent `json:"content"`
}

// ReportContent represents report content
type ReportContent struct {
	Section  string   `json:"section"`
	Metrics  []string `json:"metrics"`
	Charts   []string `json:"charts"`
	Analysis string   `json:"analysis"`
}
