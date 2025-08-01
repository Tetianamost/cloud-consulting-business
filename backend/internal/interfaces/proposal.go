package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// ProposalGenerator defines the interface for proposal and SOW generation
type ProposalGenerator interface {
	GenerateProposal(ctx context.Context, inquiry *domain.Inquiry, options *ProposalOptions) (*Proposal, error)
	GenerateSOW(ctx context.Context, inquiry *domain.Inquiry, proposal *Proposal) (*StatementOfWork, error)
	EstimateTimeline(ctx context.Context, inquiry *domain.Inquiry, projectScope *ProjectScope) (*TimelineEstimate, error)
	EstimateResources(ctx context.Context, inquiry *domain.Inquiry, projectScope *ProjectScope) (*ProposalResourceEstimate, error)
	AssessProjectRisks(ctx context.Context, inquiry *domain.Inquiry, projectScope *ProjectScope) (*ProjectRiskAssessment, error)
	GeneratePricingRecommendation(ctx context.Context, inquiry *domain.Inquiry, projectScope *ProjectScope) (*PricingRecommendation, error)
	GetSimilarProjects(ctx context.Context, inquiry *domain.Inquiry) ([]*HistoricalProject, error)
	ValidateProposal(proposal *Proposal) error
}

// ProposalOptions defines options for proposal generation
type ProposalOptions struct {
	IncludeDetailedSOW      bool     `json:"include_detailed_sow"`
	IncludeRiskAssessment   bool     `json:"include_risk_assessment"`
	IncludePricingBreakdown bool     `json:"include_pricing_breakdown"`
	IncludeTimeline         bool     `json:"include_timeline"`
	TargetAudience          string   `json:"target_audience"` // "technical", "business", "executive"
	ProposalType            string   `json:"proposal_type"`   // "rfp", "informal", "follow_up"
	CloudProviders          []string `json:"cloud_providers"`
	MaxBudget               float64  `json:"max_budget,omitempty"`
	PreferredTimeline       string   `json:"preferred_timeline,omitempty"`
}

// Proposal represents a generated proposal document
type Proposal struct {
	ID                    string                    `json:"id"`
	InquiryID             string                    `json:"inquiry_id"`
	Title                 string                    `json:"title"`
	ExecutiveSummary      string                    `json:"executive_summary"`
	ProblemStatement      string                    `json:"problem_statement"`
	ProposedSolution      *ProposalSolution         `json:"proposed_solution"`
	ProjectScope          *ProjectScope             `json:"project_scope"`
	TimelineEstimate      *TimelineEstimate         `json:"timeline_estimate"`
	ResourceEstimate      *ProposalResourceEstimate `json:"resource_estimate"`
	PricingRecommendation *PricingRecommendation    `json:"pricing_recommendation"`
	RiskAssessment        *ProjectRiskAssessment    `json:"risk_assessment"`
	NextSteps             []string                  `json:"next_steps"`
	Assumptions           []string                  `json:"assumptions"`
	Deliverables          []ProposalDeliverable     `json:"deliverables"`
	SuccessMetrics        []string                  `json:"success_metrics"`
	Status                ProposalStatus            `json:"status"`
	Version               string                    `json:"version"`
	CreatedAt             time.Time                 `json:"created_at"`
	UpdatedAt             time.Time                 `json:"updated_at"`
	ExpiresAt             *time.Time                `json:"expires_at,omitempty"`
}

// StatementOfWork represents a detailed statement of work
type StatementOfWork struct {
	ID                 string                `json:"id"`
	ProposalID         string                `json:"proposal_id"`
	InquiryID          string                `json:"inquiry_id"`
	Title              string                `json:"title"`
	ProjectOverview    string                `json:"project_overview"`
	Scope              *DetailedScope        `json:"scope"`
	Deliverables       []DetailedDeliverable `json:"deliverables"`
	Timeline           *DetailedTimeline     `json:"timeline"`
	ResourceAllocation *DetailedResources    `json:"resource_allocation"`
	PaymentSchedule    *PaymentSchedule      `json:"payment_schedule"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptance_criteria"`
	Status             SOWStatus             `json:"status"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

// ProposalSolution represents the proposed technical solution
type ProposalSolution struct {
	ID                string              `json:"id"`
	InquiryID         string              `json:"inquiry_id"`
	SolutionOverview  string              `json:"solution_overview"`
	CloudProviders    []string            `json:"cloud_providers"`
	Services          []ProposedService   `json:"services"`
	Architecture      *ArchitectureDesign `json:"architecture"`
	TechnicalApproach string              `json:"technical_approach"`
	SecurityApproach  string              `json:"security_approach"`
	EstimatedCost     string              `json:"estimated_cost"`
	Timeline          string              `json:"timeline"`
	Benefits          []string            `json:"benefits"`
}

// ProjectScope defines the scope of the project
type ProjectScope struct {
	InScope      []string `json:"in_scope"`
	OutOfScope   []string `json:"out_of_scope"`
	Assumptions  []string `json:"assumptions"`
	Constraints  []string `json:"constraints"`
	Dependencies []string `json:"dependencies"`
	Exclusions   []string `json:"exclusions"`
}

// ProposedService represents a proposed service
type ProposedService struct {
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	Description string   `json:"description"`
	Purpose     string   `json:"purpose"`
	Benefits    []string `json:"benefits"`
	Cost        string   `json:"cost"`
}

// ArchitectureDesign represents the architecture design
type ArchitectureDesign struct {
	Overview    string           `json:"overview"`
	Components  []ArchComponent  `json:"components"`
	Connections []ArchConnection `json:"connections"`
	Patterns    []string         `json:"patterns"`
	Principles  []string         `json:"principles"`
}

// TimelineEstimate represents project timeline estimation
type TimelineEstimate struct {
	TotalDuration    string              `json:"total_duration"`
	StartDate        *time.Time          `json:"start_date,omitempty"`
	EndDate          *time.Time          `json:"end_date,omitempty"`
	Phases           []TimelinePhase     `json:"phases"`
	Milestones       []ProposalMilestone `json:"milestones"`
	CriticalPath     []string            `json:"critical_path"`
	BufferTime       string              `json:"buffer_time"`
	Confidence       float64             `json:"confidence"` // 0.0 to 1.0
	SimilarProjects  []HistoricalProject `json:"similar_projects"`
	EstimationMethod string              `json:"estimation_method"`
}

// ProposalResourceEstimate represents resource estimation for the project
type ProposalResourceEstimate struct {
	TotalEffort       float64            `json:"total_effort"` // in person-hours
	TeamComposition   []RoleRequirement  `json:"team_composition"`
	SkillRequirements []SkillRequirement `json:"skill_requirements"`
	ExternalResources []ExternalResource `json:"external_resources"`
	ToolsAndLicenses  []ToolRequirement  `json:"tools_and_licenses"`
	TrainingNeeds     []TrainingNeed     `json:"training_needs"`
	Confidence        float64            `json:"confidence"` // 0.0 to 1.0
}

// TimelinePhase represents a phase in the project timeline
type TimelinePhase struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Duration     string     `json:"duration"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Dependencies []string   `json:"dependencies"`
	Deliverables []string   `json:"deliverables"`
	Resources    []string   `json:"resources"`
	RiskLevel    string     `json:"risk_level"`
}

// ProposalMilestone represents a project milestone
type ProposalMilestone struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Criteria     []string   `json:"criteria"`
	Dependencies []string   `json:"dependencies"`
	Critical     bool       `json:"critical"`
}

// RoleRequirement represents a role requirement for the project
type RoleRequirement struct {
	Role        string  `json:"role"`
	Level       string  `json:"level"`      // "junior", "mid", "senior", "lead"
	Allocation  float64 `json:"allocation"` // percentage of time
	Duration    string  `json:"duration"`
	Essential   bool    `json:"essential"`
	Description string  `json:"description"`
	HourlyRate  float64 `json:"hourly_rate,omitempty"`
}

// SkillRequirement represents a skill requirement
type SkillRequirement struct {
	Skill       string `json:"skill"`
	Level       string `json:"level"` // "basic", "intermediate", "advanced", "expert"
	Essential   bool   `json:"essential"`
	Description string `json:"description"`
}

// ExternalResource represents an external resource requirement
type ExternalResource struct {
	Type        string  `json:"type"` // "consultant", "vendor", "service"
	Description string  `json:"description"`
	Duration    string  `json:"duration"`
	Cost        float64 `json:"cost"`
	Essential   bool    `json:"essential"`
}

// ToolRequirement represents a tool or license requirement
type ToolRequirement struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "software", "license", "hardware"
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	Duration    string  `json:"duration"`
	Essential   bool    `json:"essential"`
}

// TrainingNeed represents a training requirement
type TrainingNeed struct {
	Topic       string   `json:"topic"`
	Audience    []string `json:"audience"`
	Duration    string   `json:"duration"`
	Cost        float64  `json:"cost"`
	Essential   bool     `json:"essential"`
	Description string   `json:"description"`
}

// ProjectRiskAssessment represents risk assessment for the project
type ProjectRiskAssessment struct {
	OverallRiskLevel    string          `json:"overall_risk_level"`
	TechnicalRisks      []ProjectRisk   `json:"technical_risks"`
	BusinessRisks       []ProjectRisk   `json:"business_risks"`
	ResourceRisks       []ProjectRisk   `json:"resource_risks"`
	TimelineRisks       []ProjectRisk   `json:"timeline_risks"`
	BudgetRisks         []ProjectRisk   `json:"budget_risks"`
	MitigationPlan      *MitigationPlan `json:"mitigation_plan"`
	ContingencyPlanning string          `json:"contingency_planning"`
	RiskMonitoring      []RiskIndicator `json:"risk_monitoring"`
}

// ProjectRisk represents a project risk
type ProjectRisk struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`      // "low", "medium", "high", "critical"
	Probability string `json:"probability"` // "low", "medium", "high"
	RiskScore   int    `json:"risk_score"`
	Mitigation  string `json:"mitigation"`
	Owner       string `json:"owner"`
	Status      string `json:"status"`
}

// MitigationPlan represents a risk mitigation plan
type MitigationPlan struct {
	Strategies      []ProposalMitigationStrategy `json:"strategies"`
	ContingencyFund float64                      `json:"contingency_fund"`
	EscalationPlan  string                       `json:"escalation_plan"`
	ReviewSchedule  string                       `json:"review_schedule"`
}

// ProposalMitigationStrategy represents a risk mitigation strategy
type ProposalMitigationStrategy struct {
	RiskID        string   `json:"risk_id"`
	Strategy      string   `json:"strategy"`
	Actions       []string `json:"actions"`
	Timeline      string   `json:"timeline"`
	Cost          float64  `json:"cost"`
	Owner         string   `json:"owner"`
	Effectiveness string   `json:"effectiveness"`
}

// RiskIndicator represents a risk monitoring indicator
type RiskIndicator struct {
	RiskID    string `json:"risk_id"`
	Indicator string `json:"indicator"`
	Threshold string `json:"threshold"`
	Frequency string `json:"frequency"`
	Owner     string `json:"owner"`
}

// PricingRecommendation represents pricing recommendations
type PricingRecommendation struct {
	TotalPrice          float64                `json:"total_price"`
	Currency            string                 `json:"currency"`
	PricingModel        string                 `json:"pricing_model"` // "fixed", "time_materials", "milestone"
	Breakdown           []PriceComponent       `json:"breakdown"`
	Discounts           []Discount             `json:"discounts"`
	PaymentTerms        string                 `json:"payment_terms"`
	ValidityPeriod      string                 `json:"validity_period"`
	CompetitiveAnalysis *CompetitivePricing    `json:"competitive_analysis,omitempty"`
	ValueProposition    string                 `json:"value_proposition"`
	ROIProjection       *ProposalROIProjection `json:"roi_projection,omitempty"`
	MarketRateAnalysis  *MarketRateAnalysis    `json:"market_rate_analysis,omitempty"`
}

// PriceComponent represents a component of the total price
type PriceComponent struct {
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	Notes       string  `json:"notes,omitempty"`
}

// Discount represents a pricing discount
type Discount struct {
	Type        string   `json:"type"` // "percentage", "fixed", "volume"
	Description string   `json:"description"`
	Amount      float64  `json:"amount"`
	Conditions  []string `json:"conditions"`
}

// CompetitivePricing represents competitive pricing analysis
type CompetitivePricing struct {
	MarketRange     PriceRange `json:"market_range"`
	OurPosition     string     `json:"our_position"` // "below", "at", "above"
	Justification   string     `json:"justification"`
	Differentiators []string   `json:"differentiators"`
}

// PriceRange represents a price range
type PriceRange struct {
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
	Average float64 `json:"average"`
}

// ProposalROIProjection represents ROI projection for the client
type ProposalROIProjection struct {
	InitialInvestment float64  `json:"initial_investment"`
	AnnualSavings     float64  `json:"annual_savings"`
	PaybackPeriod     string   `json:"payback_period"`
	ThreeYearROI      float64  `json:"three_year_roi"`
	FiveYearROI       float64  `json:"five_year_roi"`
	Assumptions       []string `json:"assumptions"`
}

// MarketRateAnalysis represents market rate analysis
type MarketRateAnalysis struct {
	Region          string                `json:"region"`
	ServiceCategory string                `json:"service_category"`
	RateRanges      map[string]PriceRange `json:"rate_ranges"` // role -> price range
	DataSource      string                `json:"data_source"`
	LastUpdated     time.Time             `json:"last_updated"`
}

// HistoricalProject represents a historical project for comparison
type HistoricalProject struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Industry        string    `json:"industry"`
	Services        []string  `json:"services"`
	Duration        string    `json:"duration"`
	TeamSize        int       `json:"team_size"`
	Budget          float64   `json:"budget"`
	Complexity      string    `json:"complexity"`
	SuccessMetrics  []string  `json:"success_metrics"`
	LessonsLearned  []string  `json:"lessons_learned"`
	SimilarityScore float64   `json:"similarity_score"`
	CompletedAt     time.Time `json:"completed_at"`
}

// DetailedScope represents detailed project scope
type DetailedScope struct {
	ProjectScope
	WorkBreakdownStructure []WorkPackage           `json:"work_breakdown_structure"`
	TechnicalRequirements  []TechnicalRequirement  `json:"technical_requirements"`
	FunctionalRequirements []FunctionalRequirement `json:"functional_requirements"`
}

// DetailedDeliverable represents a detailed deliverable
type DetailedDeliverable struct {
	ProposalDeliverable
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptance_criteria"`
	QualityStandards   []QualityStandard     `json:"quality_standards"`
	Dependencies       []string              `json:"dependencies"`
	ReviewProcess      string                `json:"review_process"`
}

// DetailedTimeline represents detailed project timeline
type DetailedTimeline struct {
	TimelineEstimate
	WorkSchedule     []WorkScheduleItem     `json:"work_schedule"`
	ResourceCalendar []ResourceCalendarItem `json:"resource_calendar"`
	BufferAllocation []BufferAllocation     `json:"buffer_allocation"`
}

// DetailedResources represents detailed resource allocation
type DetailedResources struct {
	ResourceEstimate
	ResourcePlan    []ProposalResourcePlanItem `json:"resource_plan"`
	SkillMatrix     []SkillMatrixItem          `json:"skill_matrix"`
	OnboardingPlan  string                     `json:"onboarding_plan"`
	OffboardingPlan string                     `json:"offboarding_plan"`
}

// Supporting types for detailed SOW

// WorkPackage represents a work package in WBS
type WorkPackage struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ParentID     string   `json:"parent_id,omitempty"`
	Level        int      `json:"level"`
	Effort       float64  `json:"effort"`
	Duration     string   `json:"duration"`
	Dependencies []string `json:"dependencies"`
	Resources    []string `json:"resources"`
	Deliverables []string `json:"deliverables"`
}

// TechnicalRequirement represents a technical requirement
type TechnicalRequirement struct {
	ID           string   `json:"id"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Priority     string   `json:"priority"`
	Rationale    string   `json:"rationale"`
	TestCriteria []string `json:"test_criteria"`
}

// FunctionalRequirement represents a functional requirement
type FunctionalRequirement struct {
	ID           string   `json:"id"`
	Feature      string   `json:"feature"`
	Description  string   `json:"description"`
	Priority     string   `json:"priority"`
	UserStory    string   `json:"user_story"`
	TestCriteria []string `json:"test_criteria"`
}

// AcceptanceCriterion represents an acceptance criterion
type AcceptanceCriterion struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	TestMethod   string `json:"test_method"`
	PassCriteria string `json:"pass_criteria"`
}

// QualityStandard represents a quality standard
type QualityStandard struct {
	Standard    string `json:"standard"`
	Description string `json:"description"`
	Metric      string `json:"metric"`
	Target      string `json:"target"`
}

// ProposalDeliverable represents a project deliverable
type ProposalDeliverable struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	DueDate     time.Time `json:"due_date"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
}

// PaymentSchedule represents the payment schedule
type PaymentSchedule struct {
	TotalAmount float64            `json:"total_amount"`
	Currency    string             `json:"currency"`
	Milestones  []PaymentMilestone `json:"milestones"`
	Terms       string             `json:"terms"`
	LateFees    string             `json:"late_fees"`
}

// PaymentMilestone represents a payment milestone
type PaymentMilestone struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Amount     float64   `json:"amount"`
	Percentage float64   `json:"percentage"`
	DueDate    time.Time `json:"due_date"`
	Criteria   []string  `json:"criteria"`
	Status     string    `json:"status"`
}

// Additional supporting types
type ArchComponent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
}

type ArchConnection struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type WorkScheduleItem struct {
	TaskID    string    `json:"task_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Duration  string    `json:"duration"`
	Resources []string  `json:"resources"`
}

type ResourceCalendarItem struct {
	ResourceID   string    `json:"resource_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Availability float64   `json:"availability"`
	Notes        string    `json:"notes"`
}

type BufferAllocation struct {
	PhaseID    string  `json:"phase_id"`
	BufferType string  `json:"buffer_type"`
	Amount     float64 `json:"amount"`
	Rationale  string  `json:"rationale"`
}

type ProposalResourcePlanItem struct {
	ResourceID string    `json:"resource_id"`
	Role       string    `json:"role"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Allocation float64   `json:"allocation"`
	Cost       float64   `json:"cost"`
}

type SkillMatrixItem struct {
	ResourceID string            `json:"resource_id"`
	Skills     map[string]string `json:"skills"` // skill -> level
	Gaps       []string          `json:"gaps"`
	Training   []string          `json:"training"`
}

// Enums and constants

// ProposalStatus represents the status of a proposal
type ProposalStatus string

const (
	ProposalStatusDraft    ProposalStatus = "draft"
	ProposalStatusReview   ProposalStatus = "review"
	ProposalStatusSent     ProposalStatus = "sent"
	ProposalStatusAccepted ProposalStatus = "accepted"
	ProposalStatusRejected ProposalStatus = "rejected"
	ProposalStatusExpired  ProposalStatus = "expired"
)

// SOWStatus represents the status of a statement of work
type SOWStatus string

const (
	SOWStatusDraft     SOWStatus = "draft"
	SOWStatusReview    SOWStatus = "review"
	SOWStatusApproved  SOWStatus = "approved"
	SOWStatusSigned    SOWStatus = "signed"
	SOWStatusActive    SOWStatus = "active"
	SOWStatusCompleted SOWStatus = "completed"
	SOWStatusCancelled SOWStatus = "cancelled"
)
