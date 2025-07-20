package interfaces

import (
	"context"

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