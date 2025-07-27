package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// RoadmapGenerator defines the interface for generating implementation roadmaps
type RoadmapGenerator interface {
	GenerateImplementationRoadmap(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (*ImplementationRoadmap, error)
	GeneratePhases(ctx context.Context, requirements []string, constraints *ProjectConstraints) ([]RoadmapPhase, error)
	EstimateResources(ctx context.Context, phases []RoadmapPhase) (*ResourceEstimate, error)
	CalculateDependencies(ctx context.Context, phases []RoadmapPhase) ([]Dependency, error)
	GenerateMilestones(ctx context.Context, phases []RoadmapPhase) ([]Milestone, error)
	ValidateRoadmap(ctx context.Context, roadmap *ImplementationRoadmap) (*ValidationResult, error)
}

// ImplementationRoadmap represents a complete implementation roadmap
type ImplementationRoadmap struct {
	ID               string            `json:"id"`
	InquiryID        string            `json:"inquiry_id"`
	Title            string            `json:"title"`
	Overview         string            `json:"overview"`
	TotalDuration    string            `json:"total_duration"`
	EstimatedCost    string            `json:"estimated_cost"`
	Phases           []RoadmapPhase    `json:"phases"`
	Dependencies     []Dependency      `json:"dependencies"`
	Risks            []string          `json:"risks"`
	SuccessMetrics   []string          `json:"success_metrics"`
	ProjectType      string            `json:"project_type"`
	CloudProviders   []string          `json:"cloud_providers"`
	IndustryContext  string            `json:"industry_context"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

// RoadmapPhase represents a phase in the implementation roadmap
type RoadmapPhase struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	Duration             string               `json:"duration"`
	EstimatedCost        string               `json:"estimated_cost"`
	Prerequisites        []string             `json:"prerequisites"`
	Deliverables         []Deliverable        `json:"deliverables"`
	Tasks                []Task               `json:"tasks"`
	Milestones           []Milestone          `json:"milestones"`
	ResourceRequirements ResourceRequirements `json:"resource_requirements"`
	RiskLevel            string               `json:"risk_level"`
	Priority             string               `json:"priority"`
	StartDate            *time.Time           `json:"start_date,omitempty"`
	EndDate              *time.Time           `json:"end_date,omitempty"`
}

// Task represents a specific task within a phase
type Task struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	EstimatedHours      int       `json:"estimated_hours"`
	SkillsRequired      []string  `json:"skills_required"`
	Dependencies        []string  `json:"dependencies"`
	Priority            string    `json:"priority"`
	DocumentationLinks  []string  `json:"documentation_links"`
	Status              string    `json:"status"`
	AssignedTo          string    `json:"assigned_to,omitempty"`
	CompletionCriteria  []string  `json:"completion_criteria"`
	EstimatedStartDate  *time.Time `json:"estimated_start_date,omitempty"`
	EstimatedEndDate    *time.Time `json:"estimated_end_date,omitempty"`
}

// Deliverable represents a deliverable within a phase
type Deliverable struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"` // "document", "code", "infrastructure", "training"
	Format       string    `json:"format"`
	Owner        string    `json:"owner"`
	Reviewers    []string  `json:"reviewers"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Status       string    `json:"status"`
	Dependencies []string  `json:"dependencies"`
	Artifacts    []string  `json:"artifacts"`
}

// Milestone represents a milestone within a phase
type Milestone struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Type         string     `json:"type"` // "phase_completion", "go_live", "checkpoint", "approval"
	Criteria     []string   `json:"criteria"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Status       string     `json:"status"`
	Dependencies []string   `json:"dependencies"`
	Stakeholders []string   `json:"stakeholders"`
	Importance   string     `json:"importance"` // "critical", "high", "medium", "low"
}

// Dependency represents a dependency between phases, tasks, or milestones
type Dependency struct {
	ID           string `json:"id"`
	FromID       string `json:"from_id"`
	ToID         string `json:"to_id"`
	Type         string `json:"type"` // "finish_to_start", "start_to_start", "finish_to_finish", "start_to_finish"
	Description  string `json:"description"`
	IsCritical   bool   `json:"is_critical"`
	LeadTime     string `json:"lead_time,omitempty"`
}

// ResourceRequirements represents resource requirements for a phase
type ResourceRequirements struct {
	TechnicalRoles   []Role   `json:"technical_roles"`
	BusinessRoles    []Role   `json:"business_roles"`
	ExternalServices []string `json:"external_services"`
	Tools            []string `json:"tools"`
	Budget           string   `json:"budget"`
	Infrastructure   []string `json:"infrastructure"`
	Licenses         []string `json:"licenses"`
}

// Role represents a role requirement
type Role struct {
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	SkillsRequired   []string `json:"skills_required"`
	ExperienceLevel  string   `json:"experience_level"` // "junior", "mid", "senior", "expert"
	TimeCommitment   string   `json:"time_commitment"`  // "full_time", "part_time", "consultant"
	Duration         string   `json:"duration"`
	EstimatedCost    string   `json:"estimated_cost"`
	IsOptional       bool     `json:"is_optional"`
}

// ProjectConstraints represents constraints for roadmap generation
type ProjectConstraints struct {
	Budget           string    `json:"budget"`
	Timeline         string    `json:"timeline"`
	TeamSize         int       `json:"team_size"`
	SkillsAvailable  []string  `json:"skills_available"`
	PreferredStart   *time.Time `json:"preferred_start,omitempty"`
	MustCompleteBy   *time.Time `json:"must_complete_by,omitempty"`
	RiskTolerance    string    `json:"risk_tolerance"` // "low", "medium", "high"
	ComplianceReqs   []string  `json:"compliance_requirements"`
	TechnologyStack  []string  `json:"technology_stack"`
	CloudProviders   []string  `json:"cloud_providers"`
}

// ResourceEstimate represents resource estimation results
type ResourceEstimate struct {
	TotalHours       int               `json:"total_hours"`
	TotalCost        string            `json:"total_cost"`
	RoleBreakdown    map[string]int    `json:"role_breakdown"`    // role -> hours
	PhaseBreakdown   map[string]int    `json:"phase_breakdown"`   // phase -> hours
	CostBreakdown    map[string]string `json:"cost_breakdown"`    // category -> cost
	ResourceUtilization map[string]float64 `json:"resource_utilization"` // role -> utilization %
}

// ValidationResult represents roadmap validation results
type ValidationResult struct {
	IsValid          bool     `json:"is_valid"`
	Errors           []string `json:"errors"`
	Warnings         []string `json:"warnings"`
	Suggestions      []string `json:"suggestions"`
	QualityScore     float64  `json:"quality_score"`
	CompletenessScore float64 `json:"completeness_score"`
}

// RoadmapTemplate represents a template for generating roadmaps
type RoadmapTemplate struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ProjectType     string                 `json:"project_type"`
	IndustryContext string                 `json:"industry_context"`
	PhaseTemplates  []PhaseTemplate        `json:"phase_templates"`
	DefaultTasks    []TaskTemplate         `json:"default_tasks"`
	Variables       map[string]interface{} `json:"variables"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// PhaseTemplate represents a template for a roadmap phase
type PhaseTemplate struct {
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	EstimatedDuration    string               `json:"estimated_duration"`
	TaskTemplates        []TaskTemplate       `json:"task_templates"`
	DeliverableTemplates []DeliverableTemplate `json:"deliverable_templates"`
	MilestoneTemplates   []MilestoneTemplate  `json:"milestone_templates"`
	ResourceTemplate     ResourceRequirements `json:"resource_template"`
}

// TaskTemplate represents a template for a task
type TaskTemplate struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	EstimatedHours     int      `json:"estimated_hours"`
	SkillsRequired     []string `json:"skills_required"`
	Priority           string   `json:"priority"`
	DocumentationLinks []string `json:"documentation_links"`
	CompletionCriteria []string `json:"completion_criteria"`
}

// DeliverableTemplate represents a template for a deliverable
type DeliverableTemplate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Format      string   `json:"format"`
	Reviewers   []string `json:"reviewers"`
}

// MilestoneTemplate represents a template for a milestone
type MilestoneTemplate struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Criteria     []string `json:"criteria"`
	Stakeholders []string `json:"stakeholders"`
	Importance   string   `json:"importance"`
}