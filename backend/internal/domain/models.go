package domain

import (
	"time"
)

// Inquiry represents a client inquiry
type Inquiry struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Company     string    `json:"company"`
	Phone       string    `json:"phone"`
	Services    []string  `json:"services"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Reports     []*Report `json:"reports,omitempty"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
}

// Report represents a generated report
type Report struct {
	ID          string       `json:"id"`
	InquiryID   string       `json:"inquiry_id"`
	Type        ReportType   `json:"type"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Status      ReportStatus `json:"status"`
	GeneratedBy string       `json:"generated_by"`
	ReviewedBy  *string      `json:"reviewed_by"`
	S3Key       string       `json:"s3_key"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	// Enhanced consultant-grade fields
	ExecutiveSummary   *ExecutiveSummary     `json:"executive_summary,omitempty"`
	BusinessCase       *BusinessCase         `json:"business_case,omitempty"`
	TechnicalSpecs     *TechnicalSpecs       `json:"technical_specs,omitempty"`
	ImplementationPlan *ImplementationPlan   `json:"implementation_plan,omitempty"`
	RiskAssessment     *RiskAssessment       `json:"risk_assessment,omitempty"`
	ROIAnalysis        *ROIAnalysis          `json:"roi_analysis,omitempty"`
	QualityMetrics     *ReportQualityMetrics `json:"quality_metrics,omitempty"`
}

// Priority represents the priority level of an inquiry
type Priority string

// ReportType represents the type of report
type ReportType string

const (
	ReportTypeAssessment         ReportType = "assessment"
	ReportTypeMigration          ReportType = "migration"
	ReportTypeOptimization       ReportType = "optimization"
	ReportTypeArchitectureReview ReportType = "architecture_review"
	ReportTypeGeneral            ReportType = "general"
	ReportTypeArchitecture       ReportType = "architecture"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusDraft     ReportStatus = "draft"
	ReportStatusGenerated ReportStatus = "generated"
	ReportStatusReviewed  ReportStatus = "reviewed"
	ReportStatusSent      ReportStatus = "sent"
)

// ServiceType represents the type of service
type ServiceType string

// Activity represents an activity log entry
type Activity struct {
	ID          string       `json:"id"`
	InquiryID   string       `json:"inquiry_id"`
	Type        ActivityType `json:"type"`
	Description string       `json:"description"`
	Actor       string       `json:"actor"`
	Metadata    string       `json:"metadata"` // JSON string
	CreatedAt   time.Time    `json:"created_at"`
}

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityTypeInquiryCreated     ActivityType = "inquiry_created"
	ActivityTypeReportGenerated    ActivityType = "report_generated"
	ActivityTypeStatusChanged      ActivityType = "status_changed"
	ActivityTypeEmailSent          ActivityType = "email_sent"
	ActivityTypeConsultantAssigned ActivityType = "consultant_assigned"
)

// InquiryFilters represents filters for querying inquiries
type InquiryFilters struct {
	Status   string `json:"status,omitempty"`
	Priority string `json:"priority,omitempty"`
	Service  string `json:"service,omitempty"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

// Enhanced consultant-grade report components

// ExecutiveSummary represents the executive summary of a report
type ExecutiveSummary struct {
	KeyFindings       []string `json:"key_findings"`
	Recommendations   []string `json:"recommendations"`
	BusinessImpact    string   `json:"business_impact"`
	InvestmentSummary string   `json:"investment_summary"`
	Timeline          string   `json:"timeline"`
}

// BusinessCase represents the business case for recommendations
type BusinessCase struct {
	ProblemStatement    string   `json:"problem_statement"`
	ProposedSolution    string   `json:"proposed_solution"`
	ExpectedBenefits    []string `json:"expected_benefits"`
	CostBenefitAnalysis string   `json:"cost_benefit_analysis"`
	ROI                 float64  `json:"roi"`
	PaybackPeriod       string   `json:"payback_period"`
}

// TechnicalSpecs represents technical specifications
type TechnicalSpecs struct {
	ArchitectureOverview   string                 `json:"architecture_overview"`
	TechnicalRequirements  []string               `json:"technical_requirements"`
	ComponentDetails       map[string]interface{} `json:"component_details"`
	IntegrationPoints      []string               `json:"integration_points"`
	SecurityConsiderations []string               `json:"security_considerations"`
}

// ImplementationPlan represents the implementation plan
type ImplementationPlan struct {
	Phases               []ImplementationPhase `json:"phases"`
	Timeline             string                `json:"timeline"`
	ResourceRequirements []string              `json:"resource_requirements"`
	RiskMitigation       []string              `json:"risk_mitigation"`
	SuccessMetrics       []string              `json:"success_metrics"`
}

// ImplementationPhase represents a phase in the implementation plan
type ImplementationPhase struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Duration     string   `json:"duration"`
	Deliverables []string `json:"deliverables"`
	Dependencies []string `json:"dependencies"`
}

// RiskAssessment represents risk assessment for the report
type RiskAssessment struct {
	OverallRiskLevel     string   `json:"overall_risk_level"`
	TechnicalRisks       []Risk   `json:"technical_risks"`
	BusinessRisks        []Risk   `json:"business_risks"`
	MitigationStrategies []string `json:"mitigation_strategies"`
}

// Risk represents a risk item
type Risk struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Probability string `json:"probability"`
	Mitigation  string `json:"mitigation"`
}

// ROIAnalysis represents return on investment analysis
type ROIAnalysis struct {
	InitialInvestment float64 `json:"initial_investment"`
	AnnualSavings     float64 `json:"annual_savings"`
	PaybackPeriod     string  `json:"payback_period"`
	NPV               float64 `json:"npv"`
	IRR               float64 `json:"irr"`
	ROIPercentage     float64 `json:"roi_percentage"`
}

// ReportQualityMetrics represents quality metrics for the report
type ReportQualityMetrics struct {
	CompletenessScore   float64  `json:"completeness_score"`
	AccuracyScore       float64  `json:"accuracy_score"`
	RelevanceScore      float64  `json:"relevance_score"`
	ActionabilityScore  float64  `json:"actionability_score"`
	OverallQualityScore float64  `json:"overall_quality_score"`
	ReviewComments      []string `json:"review_comments"`
}
