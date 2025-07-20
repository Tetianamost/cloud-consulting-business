package domain

import (
	"time"
)

// Inquiry represents a service inquiry from a client
type Inquiry struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	Phone     string    `json:"phone"`
	Services  []string  `json:"services"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Optional relationship to reports
	Reports []*Report `json:"reports,omitempty"`
}

// CreateInquiryRequest represents the request to create a new inquiry
type CreateInquiryRequest struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Company  string   `json:"company"`
	Phone    string   `json:"phone"`
	Services []string `json:"services" binding:"required"`
	Message  string   `json:"message" binding:"required"`
	Source   string   `json:"source"`
}

// Report represents a generated report for an inquiry
type Report struct {
	ID          string       `json:"id"`
	InquiryID   string       `json:"inquiry_id"`
	Type        ReportType   `json:"type"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Status      ReportStatus `json:"status"`
	GeneratedBy string       `json:"generated_by"`
	ReviewedBy  *string      `json:"reviewed_by,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ReportType represents the type of report
type ReportType string

const (
	ReportTypeAssessment   ReportType = "assessment"
	ReportTypeMigration    ReportType = "migration"
	ReportTypeOptimization ReportType = "optimization"
	ReportTypeArchitecture ReportType = "architecture_review"
	ReportTypeGeneral      ReportType = "general"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusDraft     ReportStatus = "draft"
	ReportStatusReviewed  ReportStatus = "reviewed"
	ReportStatusApproved  ReportStatus = "approved"
	ReportStatusSent      ReportStatus = "sent"
)

// Activity represents an activity log entry
type Activity struct {
	ID          string       `json:"id"`
	InquiryID   string       `json:"inquiry_id"`
	Type        ActivityType `json:"type"`
	Description string       `json:"description"`
	Actor       string       `json:"actor"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityInquiryCreated     ActivityType = "inquiry_created"
	ActivityReportGenerated    ActivityType = "report_generated"
	ActivityStatusChanged      ActivityType = "status_changed"
	ActivityConsultantAssigned ActivityType = "consultant_assigned"
	ActivityNotificationSent   ActivityType = "notification_sent"
)

// InquiryFilters represents filters for inquiry queries
type InquiryFilters struct {
	Status   string `json:"status,omitempty"`
	Priority string `json:"priority,omitempty"`
	Service  string `json:"service,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

// InquiryStatus represents the status of an inquiry
type InquiryStatus string

// Priority represents the priority level
type Priority string

// ServiceType represents the type of consulting service
type ServiceType string

