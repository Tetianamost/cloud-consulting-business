package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// InterviewPreparer defines the interface for interview preparation system
type InterviewPreparer interface {
	GenerateInterviewGuide(ctx context.Context, inquiry *domain.Inquiry) (*InterviewGuide, error)
	GenerateQuestionSet(ctx context.Context, category string, industry string) (*QuestionSet, error)
	GenerateDiscoveryChecklist(ctx context.Context, serviceType string) (*DiscoveryChecklist, error)
	GenerateFollowUpQuestions(ctx context.Context, responses []InterviewResponse) ([]*Question, error)

	// Task 7: Intelligent client meeting preparation system
	GeneratePreMeetingBriefing(ctx context.Context, inquiry *domain.Inquiry) (*PreMeetingBriefing, error)
	GenerateQuestionBank(ctx context.Context, industry string, challenges []string) (*QuestionBank, error)
	GenerateCompetitiveLandscapeAnalysis(ctx context.Context, industry string, currentSolutions []string) (*CompetitiveLandscapeAnalysis, error)
	GenerateFollowUpActionItems(ctx context.Context, meetingNotes string, clientResponses []InterviewResponse) (*FollowUpActionItems, error)
}

// InterviewGuide represents a comprehensive interview guide for client meetings
type InterviewGuide struct {
	ID                string             `json:"id"`
	InquiryID         string             `json:"inquiry_id"`
	Title             string             `json:"title"`
	Objective         string             `json:"objective"`
	EstimatedDuration string             `json:"estimated_duration"`
	PreparationNotes  []string           `json:"preparation_notes"`
	Sections          []InterviewSection `json:"sections"`
	FollowUpActions   []string           `json:"follow_up_actions"`
	CreatedAt         time.Time          `json:"created_at"`
}

// InterviewSection represents a section within an interview guide
type InterviewSection struct {
	Title            string      `json:"title"`
	Objective        string      `json:"objective"`
	Questions        []*Question `json:"questions"`
	ExpectedDuration string      `json:"expected_duration"`
	Notes            []string    `json:"notes"`
}

// Question represents an interview question
type Question struct {
	ID                 string   `json:"id"`
	Text               string   `json:"text"`
	Type               string   `json:"type"` // "open", "closed", "technical", "business"
	Category           string   `json:"category"`
	Priority           string   `json:"priority"` // "must-ask", "should-ask", "nice-to-ask"
	FollowUpQuestions  []string `json:"follow_up_questions"`
	ExpectedAnswerType string   `json:"expected_answer_type"`
	ValidationCriteria []string `json:"validation_criteria"`
}

// QuestionSet represents a set of questions for a specific category/industry
type QuestionSet struct {
	ID        string      `json:"id"`
	Category  string      `json:"category"`
	Industry  string      `json:"industry"`
	Questions []*Question `json:"questions"`
	CreatedAt time.Time   `json:"created_at"`
}

// DiscoveryChecklist represents a checklist for discovery sessions
type DiscoveryChecklist struct {
	ServiceType            string     `json:"service_type"`
	RequiredArtifacts      []Artifact `json:"required_artifacts"`
	TechnicalRequirements  []string   `json:"technical_requirements"`
	BusinessRequirements   []string   `json:"business_requirements"`
	ComplianceRequirements []string   `json:"compliance_requirements"`
	EnvironmentDetails     []string   `json:"environment_details"`
}

// Artifact represents a required artifact for discovery
type Artifact struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"` // "document", "diagram", "data", "access"
	Priority    string   `json:"priority"`
	Format      []string `json:"format"`
	Source      string   `json:"source"`
}

// InterviewResponse represents a client's response to an interview question
type InterviewResponse struct {
	QuestionID string `json:"question_id"`
	Response   string `json:"response"`
	Notes      string `json:"notes"`
	Confidence string `json:"confidence"` // "high", "medium", "low"
}

// Task 7: Intelligent client meeting preparation system data structures

// PreMeetingBriefing represents a comprehensive briefing for client meetings
type PreMeetingBriefing struct {
	ID                   string              `json:"id"`
	InquiryID            string              `json:"inquiry_id"`
	ClientBackground     *ClientBackground   `json:"client_background"`
	TalkingPoints        []TalkingPoint      `json:"talking_points"`
	KeyQuestions         []*Question         `json:"key_questions"`
	PotentialChallenges  []string            `json:"potential_challenges"`
	RecommendedApproach  string              `json:"recommended_approach"`
	CompetitorInsights   []CompetitorInsight `json:"competitor_insights"`
	IndustryContext      *IndustryContext    `json:"industry_context"`
	PreparationChecklist []string            `json:"preparation_checklist"`
	CreatedAt            time.Time           `json:"created_at"`
}

// ClientBackground represents analyzed client background information
type ClientBackground struct {
	CompanySize            string   `json:"company_size"`
	Industry               string   `json:"industry"`
	BusinessModel          string   `json:"business_model"`
	TechnologyMaturity     string   `json:"technology_maturity"`
	CloudReadiness         string   `json:"cloud_readiness"`
	KeyStakeholders        []string `json:"key_stakeholders"`
	BusinessDrivers        []string `json:"business_drivers"`
	PainPoints             []string `json:"pain_points"`
	ComplianceRequirements []string `json:"compliance_requirements"`
}

// TalkingPoint represents a strategic talking point for the meeting
type TalkingPoint struct {
	Topic            string   `json:"topic"`
	KeyMessage       string   `json:"key_message"`
	SupportingPoints []string `json:"supporting_points"`
	Timing           string   `json:"timing"`   // "opening", "discovery", "solution", "closing"
	Priority         string   `json:"priority"` // "high", "medium", "low"
	Context          string   `json:"context"`
}

// CompetitorInsight represents insights about competitors in the client's space
type CompetitorInsight struct {
	CompetitorName             string   `json:"competitor_name"`
	MarketPosition             string   `json:"market_position"`
	Strengths                  []string `json:"strengths"`
	Weaknesses                 []string `json:"weaknesses"`
	DifferentiationOpportunity string   `json:"differentiation_opportunity"`
	ClientRelevance            string   `json:"client_relevance"`
}

// IndustryContext represents industry-specific context and insights
type IndustryContext struct {
	Industry            string   `json:"industry"`
	MarketTrends        []string `json:"market_trends"`
	RegulatoryLandscape []string `json:"regulatory_landscape"`
	TechnologyTrends    []string `json:"technology_trends"`
	CommonChallenges    []string `json:"common_challenges"`
	BestPractices       []string `json:"best_practices"`
	CaseStudies         []string `json:"case_studies"`
}

// QuestionBank represents a curated bank of questions for specific contexts
type QuestionBank struct {
	ID                 string             `json:"id"`
	Industry           string             `json:"industry"`
	Challenges         []string           `json:"challenges"`
	QuestionCategories []QuestionCategory `json:"question_categories"`
	CreatedAt          time.Time          `json:"created_at"`
}

// QuestionCategory represents a category of questions in the question bank
type QuestionCategory struct {
	Category    string      `json:"category"`
	Description string      `json:"description"`
	Questions   []*Question `json:"questions"`
	Usage       string      `json:"usage"` // "discovery", "validation", "objection_handling"
}

// CompetitiveLandscapeAnalysis represents analysis of the competitive landscape
type CompetitiveLandscapeAnalysis struct {
	ID                      string                   `json:"id"`
	Industry                string                   `json:"industry"`
	CurrentSolutions        []string                 `json:"current_solutions"`
	CompetitorAnalysis      []CompetitorAnalysis     `json:"competitor_analysis"`
	MarketPositioning       *MarketPositioning       `json:"market_positioning"`
	DifferentiationStrategy *DifferentiationStrategy `json:"differentiation_strategy"`
	CompetitiveAdvantages   []string                 `json:"competitive_advantages"`
	ThreatAssessment        []ThreatAssessment       `json:"threat_assessment"`
	RecommendedStrategy     string                   `json:"recommended_strategy"`
	CreatedAt               time.Time                `json:"created_at"`
}

// CompetitorAnalysis represents detailed analysis of a specific competitor
type CompetitorAnalysis struct {
	CompetitorName      string   `json:"competitor_name"`
	MarketShare         string   `json:"market_share"`
	ServiceOfferings    []string `json:"service_offerings"`
	PricingStrategy     string   `json:"pricing_strategy"`
	Strengths           []string `json:"strengths"`
	Weaknesses          []string `json:"weaknesses"`
	ClientOverlap       string   `json:"client_overlap"`
	DifferentiationGaps []string `json:"differentiation_gaps"`
}

// MarketPositioning represents our positioning in the market
type MarketPositioning struct {
	UniqueValueProposition string   `json:"unique_value_proposition"`
	TargetSegments         []string `json:"target_segments"`
	KeyDifferentiators     []string `json:"key_differentiators"`
	CompetitiveAdvantages  []string `json:"competitive_advantages"`
	MarketOpportunities    []string `json:"market_opportunities"`
}

// DifferentiationStrategy represents strategy for differentiating from competitors
type DifferentiationStrategy struct {
	PrimaryDifferentiator    string                `json:"primary_differentiator"`
	SecondaryDifferentiators []string              `json:"secondary_differentiators"`
	MessagingStrategy        string                `json:"messaging_strategy"`
	ProofPoints              []string              `json:"proof_points"`
	CompetitiveResponses     []CompetitiveResponse `json:"competitive_responses"`
}

// CompetitiveResponse represents how to respond to competitive challenges
type CompetitiveResponse struct {
	CompetitorClaim    string   `json:"competitor_claim"`
	OurResponse        string   `json:"our_response"`
	SupportingEvidence []string `json:"supporting_evidence"`
	Timing             string   `json:"timing"` // "proactive", "reactive"
}

// ThreatAssessment represents assessment of competitive threats
type ThreatAssessment struct {
	ThreatType         string `json:"threat_type"` // "pricing", "feature", "relationship", "timing"
	Description        string `json:"description"`
	Severity           string `json:"severity"`    // "high", "medium", "low"
	Probability        string `json:"probability"` // "high", "medium", "low"
	MitigationStrategy string `json:"mitigation_strategy"`
	MonitoringRequired bool   `json:"monitoring_required"`
}

// FollowUpActionItems represents action items generated from meeting notes
type FollowUpActionItems struct {
	ID               string               `json:"id"`
	MeetingDate      time.Time            `json:"meeting_date"`
	ClientCompany    string               `json:"client_company"`
	ActionItems      []ActionItem         `json:"action_items"`
	NextSteps        []NextStep           `json:"next_steps"`
	Deliverables     []MeetingDeliverable `json:"deliverables"`
	Timeline         *ActionTimeline      `json:"timeline"`
	RiskFlags        []string             `json:"risk_flags"`
	OpportunityFlags []string             `json:"opportunity_flags"`
	CreatedAt        time.Time            `json:"created_at"`
}

// ActionItem represents a specific action item from the meeting
type ActionItem struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	Owner        string    `json:"owner"`    // "consultant", "client", "both"
	Priority     string    `json:"priority"` // "high", "medium", "low"
	DueDate      time.Time `json:"due_date"`
	Status       string    `json:"status"` // "pending", "in_progress", "completed"
	Dependencies []string  `json:"dependencies"`
	Notes        string    `json:"notes"`
}

// NextStep represents a next step in the engagement process
type NextStep struct {
	Step          string   `json:"step"`
	Description   string   `json:"description"`
	Timeline      string   `json:"timeline"`
	Prerequisites []string `json:"prerequisites"`
	Deliverables  []string `json:"deliverables"`
	Stakeholders  []string `json:"stakeholders"`
}

// MeetingDeliverable represents a deliverable committed during the meeting
type MeetingDeliverable struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Format       string    `json:"format"` // "document", "presentation", "demo", "proposal"
	DueDate      time.Time `json:"due_date"`
	Owner        string    `json:"owner"`
	Requirements []string  `json:"requirements"`
	Dependencies []string  `json:"dependencies"`
}

// ActionTimeline represents the overall timeline for follow-up actions
type ActionTimeline struct {
	ImmediateActions  []string         `json:"immediate_actions"`   // Within 24 hours
	ShortTermActions  []string         `json:"short_term_actions"`  // Within 1 week
	MediumTermActions []string         `json:"medium_term_actions"` // Within 1 month
	LongTermActions   []string         `json:"long_term_actions"`   // Beyond 1 month
	MilestoneEvents   []MilestoneEvent `json:"milestone_events"`
}

// MilestoneEvent represents a key milestone in the engagement
type MilestoneEvent struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TargetDate   time.Time `json:"target_date"`
	Criteria     []string  `json:"criteria"`
	Stakeholders []string  `json:"stakeholders"`
}
