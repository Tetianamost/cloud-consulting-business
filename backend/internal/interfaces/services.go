package interfaces

import (
	"context"
	"html/template"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// CreateInquiryRequest represents the request to create an inquiry
type CreateInquiryRequest struct {
	Name     string   `json:"name" validate:"required,min=2,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Company  string   `json:"company,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	Services []string `json:"services" validate:"required,min=1"`
	Message  string   `json:"message" validate:"required,min=1"`
	Source   string   `json:"source,omitempty"`
}

// InquiryService defines the interface for inquiry management
type InquiryService interface {
	CreateInquiry(ctx context.Context, req *CreateInquiryRequest) (*domain.Inquiry, error)
	GetInquiry(ctx context.Context, id string) (*domain.Inquiry, error)
	ListInquiries(ctx context.Context, filters *domain.InquiryFilters) ([]*domain.Inquiry, error)
	UpdateInquiryStatus(ctx context.Context, id string, status string) error
	AssignConsultant(ctx context.Context, id string, consultantID string) error
	GetInquiryCount(ctx context.Context, filters *domain.InquiryFilters) (int64, error)
}

// ReportService defines the interface for report management
type ReportService interface {
	GenerateReport(ctx context.Context, inquiry *domain.Inquiry) (*domain.Report, error)
	GenerateHTML(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (string, error)
	GeneratePDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) ([]byte, error)
	GetReport(ctx context.Context, id string) (*domain.Report, error)
	GetReportsByInquiry(ctx context.Context, inquiryID string) ([]*domain.Report, error)
	UpdateReportStatus(ctx context.Context, id string, status domain.ReportStatus) error
	GetReportTemplate(serviceType domain.ServiceType) (*ReportTemplate, error)
	ValidateReport(report *domain.Report) error
}

// NotificationService defines the interface for notification management
type NotificationService interface {
	SendNotification(ctx context.Context, notification *Notification) error
	RegisterChannel(channel NotificationChannel) error
	GetDeliveryStatus(notificationID string) (*DeliveryStatus, error)
	SendInquiryNotification(ctx context.Context, inquiry *domain.Inquiry) error
	SendReportNotification(ctx context.Context, report *domain.Report) error
}

// AgentHooksService defines the interface for agent hooks management
type AgentHooksService interface {
	RegisterHook(hookType HookType, handler HookHandler) error
	TriggerHook(ctx context.Context, hookType HookType, payload interface{}) error
	ListActiveHooks() []HookInfo
	ExecuteHook(ctx context.Context, hookID string, payload interface{}) (*HookResult, error)
	GetHookStatus(hookID string) (*HookStatus, error)
}

// ActivityService defines the interface for activity logging
type ActivityService interface {
	LogActivity(ctx context.Context, activity *domain.Activity) error
	GetActivitiesByInquiry(ctx context.Context, inquiryID string) ([]*domain.Activity, error)
	GetActivitiesByType(ctx context.Context, activityType domain.ActivityType) ([]*domain.Activity, error)
	GetRecentActivities(ctx context.Context, limit int) ([]*domain.Activity, error)
}

// BedrockService defines the interface for Amazon Bedrock AI service
type BedrockService interface {
	GenerateText(ctx context.Context, prompt string, options *BedrockOptions) (*BedrockResponse, error)
	GetModelInfo() BedrockModelInfo
	IsHealthy() bool
}

// EmailService defines the interface for email notifications
type EmailService interface {
	SendReportEmail(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) error
	SendInquiryNotification(ctx context.Context, inquiry *domain.Inquiry) error
	SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error
	// SendCustomerConfirmationWithPDF has been REMOVED - reports should never be sent to customers
	IsHealthy() bool
}

// SESService defines the interface for AWS SES integration
type SESService interface {
	SendEmail(ctx context.Context, email *EmailMessage) error
	VerifyEmailAddress(ctx context.Context, email string) error
	GetSendingQuota(ctx context.Context) (*SendingQuota, error)
	GetDeliveryStatus(ctx context.Context, messageID string) (*EmailDeliveryStatus, error)
	ProcessSESNotification(ctx context.Context, notification *SESNotification) (*SESNotificationResult, error)
	CategorizeError(errorType string, errorMessage string) *EmailErrorCategory
}

// TemplateService defines the interface for email and report template management
// CustomerConfirmationData represents the data structure for customer confirmation emails
type CustomerConfirmationData struct {
	Name     string
	Company  string
	Services string
	ID       string
}

type TemplateService interface {
	RenderEmailTemplate(ctx context.Context, templateName string, data interface{}) (string, error)
	PrepareCustomerConfirmationData(inquiry *domain.Inquiry) *CustomerConfirmationData
	RenderReportTemplate(ctx context.Context, templateName string, data interface{}) (string, error)
	LoadTemplate(templateName string) (*template.Template, error)
	ValidateTemplate(templateContent string) error
	GetAvailableTemplates() []string
	ReloadTemplates() error
	PrepareReportTemplateData(inquiry *domain.Inquiry, report *domain.Report) interface{}
	PrepareConsultantNotificationData(inquiry *domain.Inquiry, report *domain.Report, isHighPriority bool) interface{}
}

// PDFService defines the interface for PDF generation
type PDFService interface {
	GeneratePDF(ctx context.Context, htmlContent string, options *PDFOptions) ([]byte, error)
	GeneratePDFFromURL(ctx context.Context, url string, options *PDFOptions) ([]byte, error)
	IsHealthy() bool
	GetVersion() string
}

// PromptArchitect defines the interface for sophisticated prompt generation
type PromptArchitect interface {
	BuildReportPrompt(ctx context.Context, inquiry *domain.Inquiry, options *PromptOptions) (string, error)
	BuildInterviewPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error)
	BuildRiskAssessmentPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error)
	BuildCompetitiveAnalysisPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error)
	ValidatePrompt(prompt string) error
	GetTemplate(templateName string) (*PromptTemplate, error)
	RegisterTemplate(template *PromptTemplate) error
	ListTemplates() []string
}

// Supporting types for services

// ReportTemplate represents a template for report generation
type ReportTemplate struct {
	ID          string                 `json:"id"`
	ServiceType domain.ServiceType     `json:"service_type"`
	Name        string                 `json:"name"`
	Template    string                 `json:"template"`
	Variables   map[string]interface{} `json:"variables"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// Notification represents a notification to be sent
type Notification struct {
	ID          string                 `json:"id"`
	Type        NotificationType       `json:"type"`
	Recipient   string                 `json:"recipient"`
	Subject     string                 `json:"subject"`
	Message     string                 `json:"message"`
	Channel     ChannelType            `json:"channel"`
	Priority    domain.Priority        `json:"priority"`
	Metadata    map[string]interface{} `json:"metadata"`
	ScheduledAt *string                `json:"scheduled_at,omitempty"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationInquiryReceived NotificationType = "inquiry_received"
	NotificationReportGenerated NotificationType = "report_generated"
	NotificationStatusChanged   NotificationType = "status_changed"
	NotificationReminder        NotificationType = "reminder"
)

// ChannelType represents the notification channel
type ChannelType string

const (
	ChannelEmail ChannelType = "email"
	ChannelSlack ChannelType = "slack"
	ChannelSMS   ChannelType = "sms"
)

// DeliveryStatus represents the delivery status of a notification
type DeliveryStatus struct {
	NotificationID string  `json:"notification_id"`
	Status         string  `json:"status"`
	DeliveredAt    *string `json:"delivered_at,omitempty"`
	ErrorMessage   *string `json:"error_message,omitempty"`
	RetryCount     int     `json:"retry_count"`
}

// HookType represents the type of agent hook
type HookType string

const (
	HookInquiryCreated     HookType = "inquiry_created"
	HookReportGenerated    HookType = "report_generated"
	HookStatusChanged      HookType = "status_changed"
	HookNotificationSent   HookType = "notification_sent"
	HookConsultantAssigned HookType = "consultant_assigned"
)

// HookInfo represents information about a registered hook
type HookInfo struct {
	ID          string                 `json:"id"`
	Type        HookType               `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"created_at"`
}

// HookResult represents the result of hook execution
type HookResult struct {
	HookID        string                 `json:"hook_id"`
	Success       bool                   `json:"success"`
	Message       string                 `json:"message"`
	Data          map[string]interface{} `json:"data,omitempty"`
	ExecutionTime int64                  `json:"execution_time_ms"`
	Error         *string                `json:"error,omitempty"`
}

// HookStatus represents the status of a hook
type HookStatus struct {
	HookID         string  `json:"hook_id"`
	Status         string  `json:"status"`
	LastExecution  *string `json:"last_execution,omitempty"`
	ExecutionCount int64   `json:"execution_count"`
	SuccessCount   int64   `json:"success_count"`
	ErrorCount     int64   `json:"error_count"`
}

// BedrockOptions represents options for Bedrock API calls
type BedrockOptions struct {
	ModelID     string  `json:"modelId"`
	MaxTokens   int     `json:"maxTokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"topP"`
}

// BedrockResponse represents the response from Bedrock API
type BedrockResponse struct {
	Content  string            `json:"content"`
	Usage    BedrockUsage      `json:"usage"`
	Metadata map[string]string `json:"metadata"`
}

// BedrockUsage represents token usage information from Bedrock
type BedrockUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
}

// BedrockModelInfo represents information about the Bedrock model
type BedrockModelInfo struct {
	ModelID     string `json:"modelId"`
	ModelName   string `json:"modelName"`
	Provider    string `json:"provider"`
	MaxTokens   int    `json:"maxTokens"`
	IsAvailable bool   `json:"isAvailable"`
}

// EmailMessage represents an email message to be sent
type EmailMessage struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	TextBody    string            `json:"text_body"`
	HTMLBody    string            `json:"html_body,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	MessageID   string            `json:"message_id,omitempty"` // SES Message ID returned after sending
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data"`
}

// SendingQuota represents AWS SES sending quota information
type SendingQuota struct {
	Max24HourSend   float64 `json:"max_24_hour_send"`
	MaxSendRate     float64 `json:"max_send_rate"`
	SentLast24Hours float64 `json:"sent_last_24_hours"`
}

// EmailDeliveryStatus represents the delivery status of an email from SES
type EmailDeliveryStatus struct {
	MessageID     string    `json:"message_id"`
	Status        string    `json:"status"` // "sent", "delivered", "bounced", "complained"
	Timestamp     time.Time `json:"timestamp"`
	Destination   string    `json:"destination"`
	BounceType    string    `json:"bounce_type,omitempty"` // "Permanent", "Temporary"
	BounceReason  string    `json:"bounce_reason,omitempty"`
	ComplaintType string    `json:"complaint_type,omitempty"`
}

// SESNotification represents an SES notification (bounce, complaint, delivery)
type SESNotification struct {
	NotificationType string                 `json:"notificationType"` // "Bounce", "Complaint", "Delivery"
	MessageID        string                 `json:"messageId"`
	Timestamp        time.Time              `json:"timestamp"`
	Source           string                 `json:"source"`
	Destination      []string               `json:"destination"`
	Bounce           *SESBounceInfo         `json:"bounce,omitempty"`
	Complaint        *SESComplaintInfo      `json:"complaint,omitempty"`
	Delivery         *SESDeliveryInfo       `json:"delivery,omitempty"`
	RawMessage       map[string]interface{} `json:"rawMessage,omitempty"`
}

// SESBounceInfo represents bounce information from SES
type SESBounceInfo struct {
	BounceType        string                `json:"bounceType"`    // "Permanent", "Temporary"
	BounceSubType     string                `json:"bounceSubType"` // "General", "NoEmail", "Suppressed", etc.
	BouncedRecipients []SESBouncedRecipient `json:"bouncedRecipients"`
	Timestamp         time.Time             `json:"timestamp"`
	FeedbackID        string                `json:"feedbackId"`
	ReportingMTA      string                `json:"reportingMTA,omitempty"`
}

// SESBouncedRecipient represents a bounced recipient
type SESBouncedRecipient struct {
	EmailAddress   string `json:"emailAddress"`
	Action         string `json:"action,omitempty"`
	Status         string `json:"status,omitempty"`
	DiagnosticCode string `json:"diagnosticCode,omitempty"`
}

// SESComplaintInfo represents complaint information from SES
type SESComplaintInfo struct {
	ComplainedRecipients  []SESComplainedRecipient `json:"complainedRecipients"`
	Timestamp             time.Time                `json:"timestamp"`
	FeedbackID            string                   `json:"feedbackId"`
	UserAgent             string                   `json:"userAgent,omitempty"`
	ComplaintFeedbackType string                   `json:"complaintFeedbackType,omitempty"`
}

// SESComplainedRecipient represents a recipient who complained
type SESComplainedRecipient struct {
	EmailAddress string `json:"emailAddress"`
}

// SESDeliveryInfo represents delivery information from SES
type SESDeliveryInfo struct {
	Timestamp            time.Time `json:"timestamp"`
	ProcessingTimeMillis int64     `json:"processingTimeMillis"`
	Recipients           []string  `json:"recipients"`
	SMTPResponse         string    `json:"smtpResponse,omitempty"`
	ReportingMTA         string    `json:"reportingMTA,omitempty"`
}

// SESNotificationResult represents the result of processing an SES notification
type SESNotificationResult struct {
	MessageID        string    `json:"message_id"`
	NotificationType string    `json:"notification_type"`
	Status           string    `json:"status"`
	ProcessedAt      time.Time `json:"processed_at"`
	UpdatedEvents    int       `json:"updated_events"`
	Error            string    `json:"error,omitempty"`
}

// EmailErrorCategory represents categorized email error information
type EmailErrorCategory struct {
	Category   string `json:"category"`              // "bounce", "complaint", "delivery_delay", "unknown"
	Severity   string `json:"severity"`              // "permanent", "temporary", "warning"
	Reason     string `json:"reason"`                // Human-readable reason
	Actionable bool   `json:"actionable"`            // Whether the error can be acted upon
	RetryAfter *int   `json:"retry_after,omitempty"` // Seconds to wait before retry (for temporary errors)
}

// PDFOptions represents options for PDF generation
type PDFOptions struct {
	PageSize         string            `json:"page_size"`   // A4, Letter, etc.
	Orientation      string            `json:"orientation"` // Portrait, Landscape
	MarginTop        string            `json:"margin_top"`  // e.g., "1in", "2cm"
	MarginRight      string            `json:"margin_right"`
	MarginBottom     string            `json:"margin_bottom"`
	MarginLeft       string            `json:"margin_left"`
	HeaderHTML       string            `json:"header_html"` // HTML for header
	FooterHTML       string            `json:"footer_html"` // HTML for footer
	EnableJavaScript bool              `json:"enable_javascript"`
	LoadTimeout      int               `json:"load_timeout"`   // Timeout in seconds
	Quality          int               `json:"quality"`        // Image quality (0-100)
	CustomOptions    map[string]string `json:"custom_options"` // Additional wkhtmltopdf options
}

// PromptOptions defines options for prompt generation
type PromptOptions struct {
	IncludeDocumentationLinks  bool     `json:"include_documentation_links"`
	IncludeCompetitiveAnalysis bool     `json:"include_competitive_analysis"`
	IncludeRiskAssessment      bool     `json:"include_risk_assessment"`
	IncludeImplementationSteps bool     `json:"include_implementation_steps"`
	TargetAudience             string   `json:"target_audience"` // "technical", "business", "mixed"
	IndustryContext            string   `json:"industry_context"`
	CloudProviders             []string `json:"cloud_providers"`
	MaxTokens                  int      `json:"max_tokens"`
}

// PromptTemplate represents a structured prompt template
type PromptTemplate struct {
	Name              string           `json:"name"`
	Template          string           `json:"template"`
	RequiredVariables []string         `json:"required_variables"`
	OptionalVariables []string         `json:"optional_variables"`
	ValidationRules   []ValidationRule `json:"validation_rules"`
	Category          string           `json:"category"`
	Description       string           `json:"description"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
}

// ValidationRule defines a rule for prompt validation
type ValidationRule struct {
	Name         string `json:"name"`
	Pattern      string `json:"pattern"`
	ErrorMessage string `json:"error_message"`
	Required     bool   `json:"required"`
}

// EmailEventRecorder defines the interface for recording email events
type EmailEventRecorder interface {
	RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error
	UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error
	GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error)
	IsHealthy() bool
	IsHealthyWithContext(ctx context.Context) bool
	GetMetrics() *EmailEventRecorderMetrics
}

// EmailMetricsService defines the interface for email metrics calculation and retrieval
type EmailMetricsService interface {
	GetEmailMetrics(ctx context.Context, timeRange domain.TimeRange) (*domain.EmailMetrics, error)
	GetEmailStatusByInquiry(ctx context.Context, inquiryID string) (*domain.EmailStatus, error)
	GetEmailEventHistory(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error)
	GetRecentEmailActivity(ctx context.Context, hours int) ([]*domain.EmailEvent, error)
	IsHealthy(ctx context.Context) bool
	GetMetrics() *EmailMetricsServiceMetrics
}

// EmailEventRecorderMetrics represents metrics for email event recording
type EmailEventRecorderMetrics struct {
	TotalRecordingAttempts int64     `json:"total_recording_attempts"`
	SuccessfulRecordings   int64     `json:"successful_recordings"`
	FailedRecordings       int64     `json:"failed_recordings"`
	SuccessRate            float64   `json:"success_rate"`
	AverageRecordingTime   float64   `json:"average_recording_time_ms"`
	LastRecordingTime      time.Time `json:"last_recording_time"`
	RetryAttempts          int64     `json:"retry_attempts"`
	HealthCheckFailures    int64     `json:"health_check_failures"`
}

// EmailMetricsServiceMetrics represents metrics for email metrics service
type EmailMetricsServiceMetrics struct {
	TotalMetricsRequests int64     `json:"total_metrics_requests"`
	SuccessfulRequests   int64     `json:"successful_requests"`
	FailedRequests       int64     `json:"failed_requests"`
	SuccessRate          float64   `json:"success_rate"`
	AverageResponseTime  float64   `json:"average_response_time_ms"`
	LastRequestTime      time.Time `json:"last_request_time"`
	CacheHits            int64     `json:"cache_hits"`
	CacheMisses          int64     `json:"cache_misses"`
	HealthCheckFailures  int64     `json:"health_check_failures"`
}
