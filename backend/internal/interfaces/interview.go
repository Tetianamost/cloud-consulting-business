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
	Title             string      `json:"title"`
	Objective         string      `json:"objective"`
	Questions         []*Question `json:"questions"`
	ExpectedDuration  string      `json:"expected_duration"`
	Notes             []string    `json:"notes"`
}

// Question represents an interview question
type Question struct {
	ID                   string   `json:"id"`
	Text                 string   `json:"text"`
	Type                 string   `json:"type"` // "open", "closed", "technical", "business"
	Category             string   `json:"category"`
	Priority             string   `json:"priority"` // "must-ask", "should-ask", "nice-to-ask"
	FollowUpQuestions    []string `json:"follow_up_questions"`
	ExpectedAnswerType   string   `json:"expected_answer_type"`
	ValidationCriteria   []string `json:"validation_criteria"`
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