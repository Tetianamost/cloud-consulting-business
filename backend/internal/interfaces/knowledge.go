package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// KnowledgeBase defines the interface for company-specific knowledge management
type KnowledgeBase interface {
	// Service offerings and pricing
	GetServiceOfferings(ctx context.Context) ([]*ServiceOffering, error)
	GetServiceOffering(ctx context.Context, id string) (*ServiceOffering, error)
	GetPricingModels(ctx context.Context, serviceType string) ([]*PricingModel, error)

	// Team expertise and specializations
	GetTeamExpertise(ctx context.Context) ([]*TeamExpertise, error)
	GetConsultantSpecializations(ctx context.Context, consultantID string) ([]*Specialization, error)
	GetExpertiseByArea(ctx context.Context, area string) ([]*TeamExpertise, error)

	// Client history and past engagements
	GetClientHistory(ctx context.Context, clientName string) ([]*ClientEngagement, error)
	GetPastSolutions(ctx context.Context, serviceType string, industry string) ([]*PastSolution, error)
	GetSimilarProjects(ctx context.Context, inquiry *domain.Inquiry) ([]*ProjectPattern, error)

	// Methodology templates
	GetMethodologyTemplates(ctx context.Context, serviceType string) ([]*MethodologyTemplate, error)
	GetConsultingApproach(ctx context.Context, serviceType string) (*ConsultingApproach, error)
	GetDeliverableTemplates(ctx context.Context, serviceType string) ([]*DeliverableTemplate, error)

	// Knowledge base management
	UpdateKnowledgeBase(ctx context.Context) error
	SearchKnowledge(ctx context.Context, query string, category string) ([]*KnowledgeItem, error)
	GetKnowledgeStats(ctx context.Context) (*KnowledgeStats, error)

	// Additional methods for report generation
	GetBestPractices(ctx context.Context, category string) ([]*BestPractice, error)
	GetComplianceRequirements(ctx context.Context, framework string) ([]*ComplianceRequirement, error)
}

// ServiceOffering represents a company service offering
type ServiceOffering struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Category         string                 `json:"category"`
	ServiceType      domain.ServiceType     `json:"service_type"`
	KeyBenefits      []string               `json:"key_benefits"`
	Deliverables     []string               `json:"deliverables"`
	TypicalDuration  string                 `json:"typical_duration"`
	Prerequisites    []string               `json:"prerequisites"`
	TargetIndustries []string               `json:"target_industries"`
	CloudProviders   []string               `json:"cloud_providers"`
	ComplexityLevel  string                 `json:"complexity_level"` // "basic", "intermediate", "advanced"
	TeamSize         string                 `json:"team_size"`
	SuccessMetrics   []string               `json:"success_metrics"`
	RiskFactors      []string               `json:"risk_factors"`
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// PricingModel represents pricing information for services
type PricingModel struct {
	ID              string                 `json:"id"`
	ServiceID       string                 `json:"service_id"`
	Name            string                 `json:"name"`
	PricingType     string                 `json:"pricing_type"` // "fixed", "hourly", "value-based", "retainer"
	BasePrice       float64                `json:"base_price"`
	Currency        string                 `json:"currency"`
	BillingCycle    string                 `json:"billing_cycle"` // "one-time", "monthly", "quarterly", "project"
	PriceFactors    []PriceFactor          `json:"price_factors"`
	DiscountTiers   []DiscountTier         `json:"discount_tiers"`
	AdditionalCosts []AdditionalCost       `json:"additional_costs"`
	EstimationRules []EstimationRule       `json:"estimation_rules"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// PriceFactor represents factors that affect pricing
type PriceFactor struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "multiplier", "additive", "percentage"
	Value       float64 `json:"value"`
	Condition   string  `json:"condition"`
	Description string  `json:"description"`
}

// DiscountTier represents volume or loyalty discounts
type DiscountTier struct {
	Name        string  `json:"name"`
	Threshold   float64 `json:"threshold"`
	Discount    float64 `json:"discount"` // percentage or fixed amount
	Type        string  `json:"type"`     // "percentage", "fixed"
	Description string  `json:"description"`
}

// AdditionalCost represents additional costs that may apply
type AdditionalCost struct {
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Type        string  `json:"type"` // "fixed", "percentage", "per-unit"
	Condition   string  `json:"condition"`
	Description string  `json:"description"`
}

// EstimationRule represents rules for estimating project costs
type EstimationRule struct {
	Name        string                 `json:"name"`
	Condition   string                 `json:"condition"`
	Formula     string                 `json:"formula"`
	Variables   map[string]interface{} `json:"variables"`
	Description string                 `json:"description"`
}

// TeamExpertise represents team member expertise areas
type TeamExpertise struct {
	ID                 string            `json:"id"`
	ConsultantID       string            `json:"consultant_id"`
	ConsultantName     string            `json:"consultant_name"`
	Role               string            `json:"role"`
	ExpertiseAreas     []string          `json:"expertise_areas"`
	Specializations    []*Specialization `json:"specializations"`
	Certifications     []Certification   `json:"certifications"`
	ExperienceYears    int               `json:"experience_years"`
	IndustryFocus      []string          `json:"industry_focus"`
	CloudProviders     []string          `json:"cloud_providers"`
	TechnicalSkills    []TechnicalSkill  `json:"technical_skills"`
	ProjectHistory     []string          `json:"project_history"` // Project IDs
	AvailabilityStatus string            `json:"availability_status"`
	HourlyRate         float64           `json:"hourly_rate"`
	PreferredProjects  []string          `json:"preferred_projects"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// Specialization represents a specific area of specialization
type Specialization struct {
	Area            string    `json:"area"`
	Level           string    `json:"level"` // "beginner", "intermediate", "expert", "thought-leader"
	YearsExperience int       `json:"years_experience"`
	KeyProjects     []string  `json:"key_projects"`
	Certifications  []string  `json:"certifications"`
	LastUpdated     time.Time `json:"last_updated"`
}

// Certification represents professional certifications
type Certification struct {
	Name         string     `json:"name"`
	Provider     string     `json:"provider"`
	Level        string     `json:"level"`
	ObtainedDate time.Time  `json:"obtained_date"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	CertID       string     `json:"cert_id"`
	VerifyURL    string     `json:"verify_url,omitempty"`
}

// TechnicalSkill represents specific technical skills
type TechnicalSkill struct {
	Technology      string    `json:"technology"`
	Proficiency     string    `json:"proficiency"` // "basic", "intermediate", "advanced", "expert"
	YearsExperience int       `json:"years_experience"`
	LastUsed        time.Time `json:"last_used"`
	ProjectsUsed    []string  `json:"projects_used"`
}

// ClientEngagement represents past client engagements
type ClientEngagement struct {
	ID                 string                 `json:"id"`
	ClientName         string                 `json:"client_name"`
	Industry           string                 `json:"industry"`
	ProjectName        string                 `json:"project_name"`
	ServiceType        domain.ServiceType     `json:"service_type"`
	StartDate          time.Time              `json:"start_date"`
	EndDate            *time.Time             `json:"end_date,omitempty"`
	Status             string                 `json:"status"`
	TeamMembers        []string               `json:"team_members"`
	CloudProviders     []string               `json:"cloud_providers"`
	TechnologiesUsed   []string               `json:"technologies_used"`
	ProjectValue       float64                `json:"project_value"`
	ClientSatisfaction float64                `json:"client_satisfaction"` // 1-10 scale
	KeyChallenges      []string               `json:"key_challenges"`
	SolutionsProvided  []string               `json:"solutions_provided"`
	LessonsLearned     []string               `json:"lessons_learned"`
	Deliverables       []string               `json:"deliverables"`
	SuccessMetrics     map[string]interface{} `json:"success_metrics"`
	ReferenceAllowed   bool                   `json:"reference_allowed"`
	CaseStudyAllowed   bool                   `json:"case_study_allowed"`
	Metadata           map[string]interface{} `json:"metadata"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// PastSolution represents solutions from previous projects
type PastSolution struct {
	ID               string                 `json:"id"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	ServiceType      domain.ServiceType     `json:"service_type"`
	Industry         string                 `json:"industry"`
	ProblemStatement string                 `json:"problem_statement"`
	SolutionApproach string                 `json:"solution_approach"`
	TechnicalDetails string                 `json:"technical_details"`
	CloudProviders   []string               `json:"cloud_providers"`
	Technologies     []string               `json:"technologies"`
	Architecture     string                 `json:"architecture"`
	Implementation   []ImplementationStep   `json:"implementation"`
	Results          []Result               `json:"results"`
	CostSavings      float64                `json:"cost_savings"`
	TimeToValue      string                 `json:"time_to_value"`
	Reusability      string                 `json:"reusability"` // "high", "medium", "low"
	Complexity       string                 `json:"complexity"`
	RiskLevel        string                 `json:"risk_level"`
	ClientFeedback   string                 `json:"client_feedback"`
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// ImplementationStep represents a step in solution implementation
type ImplementationStep struct {
	StepNumber   int      `json:"step_number"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Duration     string   `json:"duration"`
	Dependencies []string `json:"dependencies"`
	Resources    []string `json:"resources"`
	Deliverables []string `json:"deliverables"`
	RiskFactors  []string `json:"risk_factors"`
}

// Result represents project results and outcomes
type Result struct {
	Metric      string      `json:"metric"`
	Value       interface{} `json:"value"`
	Unit        string      `json:"unit"`
	Improvement string      `json:"improvement"`
	Timeframe   string      `json:"timeframe"`
}

// ProjectPattern represents patterns from similar projects
type ProjectPattern struct {
	ID                string                 `json:"id"`
	PatternName       string                 `json:"pattern_name"`
	Description       string                 `json:"description"`
	ServiceType       domain.ServiceType     `json:"service_type"`
	Industry          string                 `json:"industry"`
	CompanySize       string                 `json:"company_size"`
	Complexity        string                 `json:"complexity"`
	CommonChallenges  []string               `json:"common_challenges"`
	TypicalSolutions  []string               `json:"typical_solutions"`
	BestPractices     []string               `json:"best_practices"`
	PitfallsToAvoid   []string               `json:"pitfalls_to_avoid"`
	EstimatedDuration string                 `json:"estimated_duration"`
	EstimatedCost     string                 `json:"estimated_cost"`
	SuccessFactors    []string               `json:"success_factors"`
	KPIs              []string               `json:"kpis"`
	SimilarProjects   []string               `json:"similar_projects"` // Project IDs
	Confidence        float64                `json:"confidence"`       // 0-1 scale
	Metadata          map[string]interface{} `json:"metadata"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// MethodologyTemplate represents consulting methodology templates
type MethodologyTemplate struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ServiceType     domain.ServiceType     `json:"service_type"`
	Phases          []MethodologyPhase     `json:"phases"`
	Prerequisites   []string               `json:"prerequisites"`
	Deliverables    []string               `json:"deliverables"`
	Tools           []string               `json:"tools"`
	BestPractices   []string               `json:"best_practices"`
	QualityGates    []QualityGate          `json:"quality_gates"`
	RiskMitigation  []string               `json:"risk_mitigation"`
	SuccessMetrics  []string               `json:"success_metrics"`
	AdaptationNotes []string               `json:"adaptation_notes"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// MethodologyPhase represents a phase in the methodology
type MethodologyPhase struct {
	PhaseNumber     int                    `json:"phase_number"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Objectives      []string               `json:"objectives"`
	Activities      []Activity             `json:"activities"`
	Duration        string                 `json:"duration"`
	Prerequisites   []string               `json:"prerequisites"`
	Deliverables    []string               `json:"deliverables"`
	QualityCriteria []string               `json:"quality_criteria"`
	RiskFactors     []string               `json:"risk_factors"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// Activity represents an activity within a methodology phase
type Activity struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Duration     string   `json:"duration"`
	Resources    []string `json:"resources"`
	Skills       []string `json:"skills"`
	Tools        []string `json:"tools"`
	Deliverables []string `json:"deliverables"`
	Dependencies []string `json:"dependencies"`
}

// QualityGate represents quality checkpoints in methodology
type QualityGate struct {
	Name          string   `json:"name"`
	Phase         string   `json:"phase"`
	Criteria      []string `json:"criteria"`
	Approver      string   `json:"approver"`
	Documentation []string `json:"documentation"`
	ExitCriteria  []string `json:"exit_criteria"`
}

// ConsultingApproach represents the overall consulting approach
type ConsultingApproach struct {
	ID                 string                 `json:"id"`
	ServiceType        domain.ServiceType     `json:"service_type"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	Philosophy         string                 `json:"philosophy"`
	KeyPrinciples      []string               `json:"key_principles"`
	EngagementModel    string                 `json:"engagement_model"`
	CommunicationStyle string                 `json:"communication_style"`
	DecisionFramework  string                 `json:"decision_framework"`
	QualityAssurance   []string               `json:"quality_assurance"`
	ClientInvolvement  string                 `json:"client_involvement"`
	ChangeManagement   string                 `json:"change_management"`
	KnowledgeTransfer  string                 `json:"knowledge_transfer"`
	Metadata           map[string]interface{} `json:"metadata"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// DeliverableTemplate represents templates for project deliverables
type DeliverableTemplate struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	ServiceType      domain.ServiceType     `json:"service_type"`
	DeliverableType  string                 `json:"deliverable_type"` // "document", "presentation", "code", "configuration"
	Template         string                 `json:"template"`
	Sections         []DeliverableSection   `json:"sections"`
	RequiredFields   []string               `json:"required_fields"`
	OptionalFields   []string               `json:"optional_fields"`
	QualityChecklist []string               `json:"quality_checklist"`
	ReviewProcess    []string               `json:"review_process"`
	ApprovalCriteria []string               `json:"approval_criteria"`
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// DeliverableSection represents a section within a deliverable template
type DeliverableSection struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Required     bool     `json:"required"`
	Order        int      `json:"order"`
	Content      string   `json:"content"`
	Variables    []string `json:"variables"`
	Instructions []string `json:"instructions"`
}

// KnowledgeItem represents a searchable knowledge item
type KnowledgeItem struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Category  string                 `json:"category"`
	Type      string                 `json:"type"` // "service", "expertise", "solution", "methodology"
	Tags      []string               `json:"tags"`
	Relevance float64                `json:"relevance"` // Search relevance score
	Source    string                 `json:"source"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// KnowledgeStats represents statistics about the knowledge base
type KnowledgeStats struct {
	TotalServices      int                    `json:"total_services"`
	TotalExpertise     int                    `json:"total_expertise"`
	TotalEngagements   int                    `json:"total_engagements"`
	TotalSolutions     int                    `json:"total_solutions"`
	TotalMethodologies int                    `json:"total_methodologies"`
	LastUpdated        time.Time              `json:"last_updated"`
	Categories         map[string]int         `json:"categories"`
	Industries         map[string]int         `json:"industries"`
	CloudProviders     map[string]int         `json:"cloud_providers"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// BestPractice represents a best practice recommendation
type BestPractice struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Provider    string   `json:"provider"`
	Tags        []string `json:"tags"`
	References  []string `json:"references"`
}

// ComplianceRequirement represents a compliance requirement
type ComplianceRequirement struct {
	ID          string   `json:"id"`
	Framework   string   `json:"framework"`
	Requirement string   `json:"requirement"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Controls    []string `json:"controls"`
}

// CloudServiceInfo represents information about cloud services
type CloudServiceInfo struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Provider      string                 `json:"provider"`
	Category      string                 `json:"category"`
	Description   string                 `json:"description"`
	Features      []string               `json:"features"`
	UseCases      []string               `json:"use_cases"`
	PricingModel  string                 `json:"pricing_model"`
	Regions       []string               `json:"regions"`
	Integrations  []string               `json:"integrations"`
	Limitations   []string               `json:"limitations"`
	BestPractices []string               `json:"best_practices"`
	Documentation []string               `json:"documentation"`
	Metadata      map[string]interface{} `json:"metadata"`
	LastUpdated   time.Time              `json:"last_updated"`
}

// ArchitecturalPattern represents architectural patterns and best practices
type ArchitecturalPattern struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Category         string                 `json:"category"`
	Provider         string                 `json:"provider"`
	UseCases         []string               `json:"use_cases"`
	Components       []string               `json:"components"`
	Benefits         []string               `json:"benefits"`
	Drawbacks        []string               `json:"drawbacks"`
	Implementation   []string               `json:"implementation"`
	BestPractices    []string               `json:"best_practices"`
	AntiPatterns     []string               `json:"anti_patterns"`
	CostImplications []string               `json:"cost_implications"`
	SecurityNotes    []string               `json:"security_notes"`
	Examples         []string               `json:"examples"`
	References       []string               `json:"references"`
	Metadata         map[string]interface{} `json:"metadata"`
	LastUpdated      time.Time              `json:"last_updated"`
}
