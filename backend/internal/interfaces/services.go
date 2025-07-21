package interfaces

import (
	"context"
	"html/template"

	"github.com/cloud-consulting/backend/internal/domain"
)

// InquiryService defines the interface for inquiry management
type InquiryService interface {
	CreateInquiry(ctx context.Context, req *domain.CreateInquiryRequest) (*domain.Inquiry, error)
	GetInquiry(ctx context.Context, id string) (*domain.Inquiry, error)
	ListInquiries(ctx context.Context, filters *domain.InquiryFilters) ([]*domain.Inquiry, error)
	UpdateInquiryStatus(ctx context.Context, id string, status domain.InquiryStatus) error
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
	SendReportEmailWithPDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report, pdfData []byte) error
	SendInquiryNotification(ctx context.Context, inquiry *domain.Inquiry) error
	SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error
	SendCustomerConfirmationWithPDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report, pdfData []byte) error
	IsHealthy() bool
}

// SESService defines the interface for AWS SES integration
type SESService interface {
	SendEmail(ctx context.Context, email *EmailMessage) error
	VerifyEmailAddress(ctx context.Context, email string) error
	GetSendingQuota(ctx context.Context) (*SendingQuota, error)
}

// TemplateService defines the interface for email and report template management
type TemplateService interface {
	RenderEmailTemplate(ctx context.Context, templateName string, data interface{}) (string, error)
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
	NotificationID string    `json:"notification_id"`
	Status         string    `json:"status"`
	DeliveredAt    *string   `json:"delivered_at,omitempty"`
	ErrorMessage   *string   `json:"error_message,omitempty"`
	RetryCount     int       `json:"retry_count"`
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
	HookID      string                 `json:"hook_id"`
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Data        map[string]interface{} `json:"data,omitempty"`
	ExecutionTime int64                `json:"execution_time_ms"`
	Error       *string                `json:"error,omitempty"`
}

// HookStatus represents the status of a hook
type HookStatus struct {
	HookID        string `json:"hook_id"`
	Status        string `json:"status"`
	LastExecution *string `json:"last_execution,omitempty"`
	ExecutionCount int64  `json:"execution_count"`
	SuccessCount   int64  `json:"success_count"`
	ErrorCount     int64  `json:"error_count"`
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
	Content   string            `json:"content"`
	Usage     BedrockUsage      `json:"usage"`
	Metadata  map[string]string `json:"metadata"`
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

// PDFOptions represents options for PDF generation
type PDFOptions struct {
	PageSize        string            `json:"page_size"`        // A4, Letter, etc.
	Orientation     string            `json:"orientation"`      // Portrait, Landscape
	MarginTop       string            `json:"margin_top"`       // e.g., "1in", "2cm"
	MarginRight     string            `json:"margin_right"`
	MarginBottom    string            `json:"margin_bottom"`
	MarginLeft      string            `json:"margin_left"`
	HeaderHTML      string            `json:"header_html"`      // HTML for header
	FooterHTML      string            `json:"footer_html"`      // HTML for footer
	EnableJavaScript bool             `json:"enable_javascript"`
	LoadTimeout     int               `json:"load_timeout"`     // Timeout in seconds
	Quality         int               `json:"quality"`          // Image quality (0-100)
	CustomOptions   map[string]string `json:"custom_options"`   // Additional wkhtmltopdf options
}