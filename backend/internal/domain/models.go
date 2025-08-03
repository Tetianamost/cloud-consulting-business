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

// Chat-related models for AI consultant live chat feature

// SessionStatus represents the status of a chat session
type SessionStatus string

const (
	SessionStatusActive     SessionStatus = "active"
	SessionStatusInactive   SessionStatus = "inactive"
	SessionStatusExpired    SessionStatus = "expired"
	SessionStatusTerminated SessionStatus = "terminated"
)

// MessageType represents the type of a chat message
type MessageType string

const (
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
	MessageTypeSystem    MessageType = "system"
	MessageTypeError     MessageType = "error"
)

// MessageStatus represents the status of a chat message
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// ChatSession represents a chat session between a user and the AI assistant
type ChatSession struct {
	ID           string                 `json:"id" db:"id" validate:"required"`
	UserID       string                 `json:"user_id" db:"user_id" validate:"required,max=100"`
	ClientName   string                 `json:"client_name" db:"client_name" validate:"max=255"`
	Context      string                 `json:"context" db:"context"`
	Status       SessionStatus          `json:"status" db:"status" validate:"required"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	LastActivity time.Time              `json:"last_activity" db:"last_activity"`
	ExpiresAt    *time.Time             `json:"expires_at" db:"expires_at"`
}

// ChatMessage represents a single message within a chat session
type ChatMessage struct {
	ID        string                 `json:"id" db:"id" validate:"required"`
	SessionID string                 `json:"session_id" db:"session_id" validate:"required"`
	Type      MessageType            `json:"type" db:"type" validate:"required"`
	Content   string                 `json:"content" db:"content" validate:"required"`
	Metadata  map[string]interface{} `json:"metadata" db:"metadata"`
	Status    MessageStatus          `json:"status" db:"status" validate:"required"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

// SessionContext represents the context information for a chat session
type SessionContext struct {
	ClientName     string            `json:"client_name" validate:"max=255"`
	MeetingType    string            `json:"meeting_type" validate:"max=100"`
	ProjectContext string            `json:"project_context"`
	ServiceTypes   []string          `json:"service_types"`
	CloudProviders []string          `json:"cloud_providers"`
	CustomFields   map[string]string `json:"custom_fields"`
}

// ChatRequest represents a request to send a chat message
type ChatRequest struct {
	SessionID   string                 `json:"session_id" validate:"required"`
	Content     string                 `json:"content" validate:"required,max=10000"`
	Type        MessageType            `json:"type" validate:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
	QuickAction string                 `json:"quick_action,omitempty"`
}

// ChatResponse represents a response from the AI assistant
type ChatResponse struct {
	MessageID   string                 `json:"message_id"`
	SessionID   string                 `json:"session_id"`
	Content     string                 `json:"content"`
	Type        MessageType            `json:"type"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	Status      MessageStatus          `json:"status"`
	TokensUsed  int                    `json:"tokens_used,omitempty"`
	ProcessTime float64                `json:"process_time,omitempty"`
}

// SessionMetadata represents metadata for a chat session
type SessionMetadata struct {
	ClientName     string            `json:"client_name,omitempty"`
	MeetingType    string            `json:"meeting_type,omitempty"`
	ProjectContext string            `json:"project_context,omitempty"`
	ServiceTypes   []string          `json:"service_types,omitempty"`
	CloudProviders []string          `json:"cloud_providers,omitempty"`
	CustomFields   map[string]string `json:"custom_fields,omitempty"`
	UserAgent      string            `json:"user_agent,omitempty"`
	IPAddress      string            `json:"ip_address,omitempty"`
	SessionVersion string            `json:"session_version,omitempty"`
}

// ChatSessionFilters represents filters for querying chat sessions
type ChatSessionFilters struct {
	UserID     string        `json:"user_id,omitempty"`
	Status     SessionStatus `json:"status,omitempty"`
	ClientName string        `json:"client_name,omitempty"`
	FromDate   *time.Time    `json:"from_date,omitempty"`
	ToDate     *time.Time    `json:"to_date,omitempty"`
	Limit      int           `json:"limit" validate:"min=1,max=100"`
	Offset     int           `json:"offset" validate:"min=0"`
}

// ChatMessageFilters represents filters for querying chat messages
type ChatMessageFilters struct {
	SessionID string        `json:"session_id,omitempty"`
	Type      MessageType   `json:"type,omitempty"`
	Status    MessageStatus `json:"status,omitempty"`
	FromDate  *time.Time    `json:"from_date,omitempty"`
	ToDate    *time.Time    `json:"to_date,omitempty"`
	Limit     int           `json:"limit" validate:"min=1,max=1000"`
	Offset    int           `json:"offset" validate:"min=0"`
}

// Validation methods for ChatSession
func (cs *ChatSession) Validate() error {
	if cs.UserID == "" {
		return NewValidationError("user_id", "User ID is required")
	}
	if len(cs.UserID) > 100 {
		return NewValidationError("user_id", "User ID must be 100 characters or less")
	}
	if len(cs.ClientName) > 255 {
		return NewValidationError("client_name", "Client name must be 255 characters or less")
	}
	if cs.Status == "" {
		return NewValidationError("status", "Status is required")
	}
	return nil
}

// Validation methods for ChatMessage
func (cm *ChatMessage) Validate() error {
	if cm.SessionID == "" {
		return NewValidationError("session_id", "Session ID is required")
	}
	if cm.Content == "" {
		return NewValidationError("content", "Content is required")
	}
	if len(cm.Content) > 10000 {
		return NewValidationError("content", "Content must be 10000 characters or less")
	}
	if cm.Type == "" {
		return NewValidationError("type", "Message type is required")
	}
	return nil
}

// IsExpired checks if a chat session has expired
func (cs *ChatSession) IsExpired() bool {
	if cs.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*cs.ExpiresAt)
}

// IsActive checks if a chat session is active and not expired
func (cs *ChatSession) IsActive() bool {
	return cs.Status == SessionStatusActive && !cs.IsExpired()
}

// SetExpiration sets the expiration time for a chat session
func (cs *ChatSession) SetExpiration(duration time.Duration) {
	expiresAt := time.Now().Add(duration)
	cs.ExpiresAt = &expiresAt
}

// UpdateActivity updates the last activity timestamp
func (cs *ChatSession) UpdateActivity() {
	cs.LastActivity = time.Now()
	cs.UpdatedAt = time.Now()
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}
