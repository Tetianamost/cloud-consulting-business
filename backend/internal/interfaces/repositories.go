package interfaces

import (
	"context"

	"github.com/cloud-consulting/backend/internal/domain"
)

// InquiryRepository defines the interface for inquiry data access
type InquiryRepository interface {
	Create(ctx context.Context, inquiry *domain.Inquiry) error
	GetByID(ctx context.Context, id string) (*domain.Inquiry, error)
	GetByEmail(ctx context.Context, email string) ([]*domain.Inquiry, error)
	List(ctx context.Context, filters *domain.InquiryFilters) ([]*domain.Inquiry, error)
	Update(ctx context.Context, inquiry *domain.Inquiry) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filters *domain.InquiryFilters) (int64, error)
	GetByStatus(ctx context.Context, status string) ([]*domain.Inquiry, error)
}

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	GetByID(ctx context.Context, id string) (*domain.Report, error)
	GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.Report, error)
	List(ctx context.Context, filters *ReportFilters) ([]*domain.Report, error)
	Update(ctx context.Context, report *domain.Report) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, status domain.ReportStatus) ([]*domain.Report, error)
	GetByType(ctx context.Context, reportType domain.ReportType) ([]*domain.Report, error)
}

// ActivityRepository defines the interface for activity log data access
type ActivityRepository interface {
	Create(ctx context.Context, activity *domain.Activity) error
	GetByID(ctx context.Context, id string) (*domain.Activity, error)
	GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.Activity, error)
	GetByType(ctx context.Context, activityType domain.ActivityType) ([]*domain.Activity, error)
	List(ctx context.Context, filters *ActivityFilters) ([]*domain.Activity, error)
	GetRecent(ctx context.Context, limit int) ([]*domain.Activity, error)
	Delete(ctx context.Context, id string) error
}

// EmailEventRepository defines the interface for email event data access
type EmailEventRepository interface {
	Create(ctx context.Context, event *domain.EmailEvent) error
	Update(ctx context.Context, event *domain.EmailEvent) error
	GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error)
	GetByMessageID(ctx context.Context, messageID string) (*domain.EmailEvent, error)
	GetMetrics(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error)
	List(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error)
}

// Supporting types for repositories

// ReportFilters represents filters for listing reports
type ReportFilters struct {
	InquiryID   *string              `json:"inquiry_id,omitempty"`
	Type        *domain.ReportType   `json:"type,omitempty"`
	Status      *domain.ReportStatus `json:"status,omitempty"`
	GeneratedBy *string              `json:"generated_by,omitempty"`
	DateFrom    *string              `json:"date_from,omitempty"`
	DateTo      *string              `json:"date_to,omitempty"`
	Limit       int                  `json:"limit,omitempty"`
	Offset      int                  `json:"offset,omitempty"`
}

// ActivityFilters represents filters for listing activities
type ActivityFilters struct {
	InquiryID *string              `json:"inquiry_id,omitempty"`
	Type      *domain.ActivityType `json:"type,omitempty"`
	Actor     *string              `json:"actor,omitempty"`
	DateFrom  *string              `json:"date_from,omitempty"`
	DateTo    *string              `json:"date_to,omitempty"`
	Limit     int                  `json:"limit,omitempty"`
	Offset    int                  `json:"offset,omitempty"`
}
