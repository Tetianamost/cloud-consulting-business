package interfaces

import (
	"context"
	"time"
)

// ClientCommunicationService defines the interface for advanced client communication tools
type ClientCommunicationService interface {
	// Technical explanation generator
	GenerateTechnicalExplanation(ctx context.Context, req *TechnicalExplanationRequest) (*TechnicalExplanation, error)

	// Presentation slide generator
	GeneratePresentationSlides(ctx context.Context, req *PresentationRequest) (*PresentationSlides, error)

	// Email template generator
	GenerateEmailTemplate(ctx context.Context, req *EmailTemplateRequest) (*EmailTemplate, error)

	// Status report generator
	GenerateStatusReport(ctx context.Context, req *StatusReportRequest) (*StatusReport, error)
}

// CommunicationAudience represents the target audience for communications
type CommunicationAudience string

const (
	CommunicationAudienceExecutive   CommunicationAudience = "executive"
	CommunicationAudienceTechnical   CommunicationAudience = "technical"
	CommunicationAudienceBusiness    CommunicationAudience = "business"
	CommunicationAudienceMixed       CommunicationAudience = "mixed"
	CommunicationAudienceStakeholder CommunicationAudience = "stakeholder"
)

// ComplexityLevel represents the technical complexity level
type ComplexityLevel string

const (
	ComplexityBasic        ComplexityLevel = "basic"
	ComplexityIntermediate ComplexityLevel = "intermediate"
	ComplexityAdvanced     ComplexityLevel = "advanced"
	ComplexityExpert       ComplexityLevel = "expert"
)

// TechnicalExplanationRequest represents a request for technical explanation
type TechnicalExplanationRequest struct {
	TechnicalConcept string                 `json:"technical_concept" validate:"required"`
	AudienceType     CommunicationAudience  `json:"audience_type" validate:"required"`
	IndustryContext  string                 `json:"industry_context,omitempty"`
	ComplexityLevel  ComplexityLevel        `json:"complexity_level" validate:"required"`
	UseCase          string                 `json:"use_case,omitempty"`
	BusinessContext  string                 `json:"business_context,omitempty"`
	KeyTerms         []string               `json:"key_terms,omitempty"`
	Constraints      []string               `json:"constraints,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// TechnicalExplanation represents the generated technical explanation
type TechnicalExplanation struct {
	ID                string                 `json:"id"`
	Title             string                 `json:"title"`
	ExecutiveSummary  string                 `json:"executive_summary"`
	BusinessValue     string                 `json:"business_value"`
	TechnicalOverview string                 `json:"technical_overview"`
	KeyBenefits       []string               `json:"key_benefits"`
	Considerations    []string               `json:"considerations"`
	NextSteps         []string               `json:"next_steps"`
	Glossary          map[string]string      `json:"glossary"`
	VisualAids        []VisualAid            `json:"visual_aids"`
	References        []Reference            `json:"references"`
	AudienceType      CommunicationAudience  `json:"audience_type"`
	ComplexityLevel   ComplexityLevel        `json:"complexity_level"`
	GeneratedAt       time.Time              `json:"generated_at"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// VisualAid represents a visual aid for explanations
type VisualAid struct {
	Type        string `json:"type"` // "diagram", "chart", "image", "flowchart"
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"` // Mermaid diagram, chart data, etc.
	URL         string `json:"url,omitempty"`
}

// Reference represents a reference or documentation link
type Reference struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"` // "documentation", "best_practice", "case_study"
	Description string `json:"description,omitempty"`
}

// SlideType represents the type of presentation slide
type SlideType string

const (
	SlideTitle      SlideType = "title"
	SlideContent    SlideType = "content"
	SlideDiagram    SlideType = "diagram"
	SlideComparison SlideType = "comparison"
	SlideTimeline   SlideType = "timeline"
	SlideConclusion SlideType = "conclusion"
	SlideAppendix   SlideType = "appendix"
)

// PresentationRequest represents a request for presentation slides
type PresentationRequest struct {
	Topic           string                 `json:"topic" validate:"required"`
	AudienceType    CommunicationAudience  `json:"audience_type" validate:"required"`
	Duration        int                    `json:"duration"` // in minutes
	SlideCount      int                    `json:"slide_count,omitempty"`
	IncludeDiagrams bool                   `json:"include_diagrams"`
	IncludeCostInfo bool                   `json:"include_cost_info"`
	IndustryContext string                 `json:"industry_context,omitempty"`
	BusinessGoals   []string               `json:"business_goals,omitempty"`
	TechnicalScope  []string               `json:"technical_scope,omitempty"`
	ComplianceReqs  []string               `json:"compliance_requirements,omitempty"`
	BrandingStyle   string                 `json:"branding_style,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// PresentationSlides represents the generated presentation
type PresentationSlides struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Subtitle       string                 `json:"subtitle"`
	Slides         []Slide                `json:"slides"`
	SpeakerNotes   []string               `json:"speaker_notes"`
	AppendixSlides []Slide                `json:"appendix_slides"`
	References     []Reference            `json:"references"`
	EstimatedTime  int                    `json:"estimated_time"` // in minutes
	AudienceType   CommunicationAudience  `json:"audience_type"`
	GeneratedAt    time.Time              `json:"generated_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// Slide represents a single presentation slide
type Slide struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	SlideType    SlideType `json:"slide_type"`
	BulletPoints []string  `json:"bullet_points,omitempty"`
	Diagrams     []Diagram `json:"diagrams,omitempty"`
	Charts       []Chart   `json:"charts,omitempty"`
	Images       []Image   `json:"images,omitempty"`
	SpeakerNotes string    `json:"speaker_notes,omitempty"`
	Duration     int       `json:"duration"` // estimated time in minutes
	Order        int       `json:"order"`
}

// Diagram represents a technical diagram
type Diagram struct {
	Type        string `json:"type"` // "architecture", "flow", "network", "sequence"
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"` // Mermaid diagram syntax
	Notes       string `json:"notes,omitempty"`
}

// Chart represents a data chart
type Chart struct {
	Type   string                 `json:"type"` // "bar", "line", "pie", "scatter"
	Title  string                 `json:"title"`
	Data   map[string]interface{} `json:"data"`
	Config map[string]interface{} `json:"config"`
}

// Image represents an image reference
type Image struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Alt         string `json:"alt"`
	Description string `json:"description,omitempty"`
}

// EmailPurpose represents the purpose of the email
type EmailPurpose string

const (
	EmailFollowUp       EmailPurpose = "follow_up"
	EmailStatusUpdate   EmailPurpose = "status_update"
	EmailMeetingRequest EmailPurpose = "meeting_request"
	EmailProposal       EmailPurpose = "proposal"
	EmailDeliverable    EmailPurpose = "deliverable"
	EmailEscalation     EmailPurpose = "escalation"
	EmailIntroduction   EmailPurpose = "introduction"
	EmailClosing        EmailPurpose = "closing"
)

// RecipientType represents the type of email recipient
type RecipientType string

const (
	RecipientClient      RecipientType = "client"
	RecipientStakeholder RecipientType = "stakeholder"
	RecipientTeam        RecipientType = "team"
	RecipientVendor      RecipientType = "vendor"
	RecipientExecutive   RecipientType = "executive"
)

// ToneStyle represents the tone and style of communication
type ToneStyle string

const (
	ToneFormal        ToneStyle = "formal"
	ToneProfessional  ToneStyle = "professional"
	ToneFriendly      ToneStyle = "friendly"
	ToneConcise       ToneStyle = "concise"
	ToneUrgent        ToneStyle = "urgent"
	ToneCollaborative ToneStyle = "collaborative"
)

// UrgencyLevel represents the urgency level
type UrgencyLevel string

const (
	UrgencyLow      UrgencyLevel = "low"
	UrgencyMedium   UrgencyLevel = "medium"
	UrgencyHigh     UrgencyLevel = "high"
	UrgencyCritical UrgencyLevel = "critical"
)

// EmailTemplateRequest represents a request for email template generation
type EmailTemplateRequest struct {
	Purpose           EmailPurpose           `json:"purpose" validate:"required"`
	RecipientType     RecipientType          `json:"recipient_type" validate:"required"`
	Context           string                 `json:"context" validate:"required"`
	ToneStyle         ToneStyle              `json:"tone_style" validate:"required"`
	IncludeAttachment bool                   `json:"include_attachment"`
	ActionRequired    bool                   `json:"action_required"`
	Urgency           UrgencyLevel           `json:"urgency"`
	ProjectPhase      string                 `json:"project_phase,omitempty"`
	KeyPoints         []string               `json:"key_points,omitempty"`
	CallToAction      string                 `json:"call_to_action,omitempty"`
	Deadline          *time.Time             `json:"deadline,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// EmailTemplate represents the generated email template
type EmailTemplate struct {
	ID              string                 `json:"id"`
	Subject         string                 `json:"subject"`
	Body            string                 `json:"body"`
	HTMLBody        string                 `json:"html_body,omitempty"`
	Purpose         EmailPurpose           `json:"purpose"`
	RecipientType   RecipientType          `json:"recipient_type"`
	ToneStyle       ToneStyle              `json:"tone_style"`
	SuggestedTiming string                 `json:"suggested_timing"`
	FollowUpActions []string               `json:"follow_up_actions"`
	Alternatives    []string               `json:"alternatives"` // alternative subject lines or approaches
	GeneratedAt     time.Time              `json:"generated_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ReportingPeriod represents the reporting period
type ReportingPeriod string

const (
	PeriodWeekly    ReportingPeriod = "weekly"
	PeriodBiweekly  ReportingPeriod = "biweekly"
	PeriodMonthly   ReportingPeriod = "monthly"
	PeriodQuarterly ReportingPeriod = "quarterly"
	PeriodMilestone ReportingPeriod = "milestone"
	PeriodAdHoc     ReportingPeriod = "ad_hoc"
)

// ProjectStatus represents the overall project status
type ProjectStatus string

const (
	StatusOnTrack   ProjectStatus = "on_track"
	StatusAtRisk    ProjectStatus = "at_risk"
	StatusDelayed   ProjectStatus = "delayed"
	StatusBlocked   ProjectStatus = "blocked"
	StatusCompleted ProjectStatus = "completed"
	StatusCancelled ProjectStatus = "cancelled"
)

// StatusReportRequest represents a request for status report generation
type StatusReportRequest struct {
	ProjectID        string                 `json:"project_id" validate:"required"`
	ReportingPeriod  ReportingPeriod        `json:"reporting_period" validate:"required"`
	AudienceType     CommunicationAudience  `json:"audience_type" validate:"required"`
	IncludeMetrics   bool                   `json:"include_metrics"`
	IncludeRisks     bool                   `json:"include_risks"`
	IncludeNextSteps bool                   `json:"include_next_steps"`
	ProjectPhase     string                 `json:"project_phase,omitempty"`
	Milestones       []ProjectMilestone     `json:"milestones,omitempty"`
	Issues           []Issue                `json:"issues,omitempty"`
	Achievements     []Achievement          `json:"achievements,omitempty"`
	Budget           *BudgetStatus          `json:"budget,omitempty"`
	Timeline         *TimelineStatus        `json:"timeline,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// StatusReport represents the generated status report
type StatusReport struct {
	ID                 string                 `json:"id"`
	ProjectID          string                 `json:"project_id"`
	Title              string                 `json:"title"`
	ExecutiveSummary   string                 `json:"executive_summary"`
	OverallStatus      ProjectStatus          `json:"overall_status"`
	ProgressSummary    string                 `json:"progress_summary"`
	KeyAccomplishments []string               `json:"key_accomplishments"`
	CurrentActivities  []string               `json:"current_activities"`
	UpcomingMilestones []ProjectMilestone     `json:"upcoming_milestones"`
	RisksAndIssues     []RiskIssue            `json:"risks_and_issues"`
	BudgetStatus       *BudgetStatus          `json:"budget_status,omitempty"`
	TimelineStatus     *TimelineStatus        `json:"timeline_status,omitempty"`
	Metrics            []Metric               `json:"metrics,omitempty"`
	NextSteps          []string               `json:"next_steps"`
	Recommendations    []string               `json:"recommendations"`
	ReportingPeriod    ReportingPeriod        `json:"reporting_period"`
	AudienceType       CommunicationAudience  `json:"audience_type"`
	GeneratedAt        time.Time              `json:"generated_at"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// ProjectMilestone represents a project milestone
type ProjectMilestone struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DueDate      time.Time `json:"due_date"`
	Status       string    `json:"status"`
	Progress     int       `json:"progress"` // percentage
	Dependencies []string  `json:"dependencies,omitempty"`
}

// Issue represents a project issue
type Issue struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    string     `json:"severity"`
	Status      string     `json:"status"`
	AssignedTo  string     `json:"assigned_to,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// Achievement represents a project achievement
type Achievement struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Impact      string    `json:"impact"`
	Category    string    `json:"category"`
}

// BudgetStatus represents budget status information
type BudgetStatus struct {
	TotalBudget     float64 `json:"total_budget"`
	SpentAmount     float64 `json:"spent_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	BurnRate        float64 `json:"burn_rate"`
	ProjectedSpend  float64 `json:"projected_spend"`
	Status          string  `json:"status"` // "on_budget", "over_budget", "under_budget"
	Variance        float64 `json:"variance"`
}

// TimelineStatus represents timeline status information
type TimelineStatus struct {
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	CurrentDate     time.Time `json:"current_date"`
	ProgressPercent int       `json:"progress_percent"`
	DaysRemaining   int       `json:"days_remaining"`
	Status          string    `json:"status"` // "on_schedule", "ahead", "behind"
	CriticalPath    []string  `json:"critical_path"`
	DelayedTasks    []string  `json:"delayed_tasks,omitempty"`
}

// RiskIssue represents a risk or issue in the status report
type RiskIssue struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"` // "risk", "issue"
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Impact      string     `json:"impact"`
	Probability string     `json:"probability,omitempty"` // for risks
	Status      string     `json:"status"`
	Mitigation  string     `json:"mitigation"`
	Owner       string     `json:"owner"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// Metric represents a project metric
type Metric struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Unit        string      `json:"unit,omitempty"`
	Target      interface{} `json:"target,omitempty"`
	Trend       string      `json:"trend"` // "up", "down", "stable"
	Description string      `json:"description,omitempty"`
}
