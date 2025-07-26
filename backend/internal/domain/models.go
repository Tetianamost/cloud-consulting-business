package domain

import (
	"time"
)

// Inquiry represents a client inquiry
type Inquiry struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Company   string          `json:"company"`
	Phone     string          `json:"phone"`
	Services  []string        `json:"services"`
	Message   string          `json:"message"`
	Status    string          `json:"status"`
	Priority  string          `json:"priority"`
	Source    string          `json:"source"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Reports   []*Report       `json:"reports,omitempty"`
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
}



// Priority represents the priority level of an inquiry
type Priority string

// ReportType represents the type of report
type ReportType string

const (
	ReportTypeAssessment        ReportType = "assessment"
	ReportTypeMigration         ReportType = "migration"
	ReportTypeOptimization      ReportType = "optimization"
	ReportTypeArchitectureReview ReportType = "architecture_review"
	ReportTypeGeneral           ReportType = "general"
	ReportTypeArchitecture      ReportType = "architecture"
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