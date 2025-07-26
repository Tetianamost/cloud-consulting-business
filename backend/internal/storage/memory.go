package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// InMemoryStorage provides in-memory storage for inquiries and reports
type InMemoryStorage struct {
	inquiries map[string]*domain.Inquiry
	reports   map[string]*domain.Report
	mutex     sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		inquiries: make(map[string]*domain.Inquiry),
		reports:   make(map[string]*domain.Report),
	}
}

// CreateInquiry creates a new inquiry from a request
func (s *InMemoryStorage) CreateInquiry(req *interfaces.CreateInquiryRequest) (*domain.Inquiry, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Generate a unique ID
	id := fmt.Sprintf("inq_%d", time.Now().UnixNano())
	
	inquiry := &domain.Inquiry{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		Company:   req.Company,
		Phone:     req.Phone,
		Services:  req.Services,
		Message:   req.Message,
		Status:    domain.InquiryStatusPending,
		Priority:  domain.PriorityMedium,
		Source:    req.Source,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	s.inquiries[inquiry.ID] = inquiry
	return inquiry, nil
}

// StoreInquiry stores an inquiry object directly
func (s *InMemoryStorage) StoreInquiry(ctx context.Context, inquiry *domain.Inquiry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if inquiry.CreatedAt.IsZero() {
		inquiry.CreatedAt = time.Now()
	}
	inquiry.UpdatedAt = time.Now()
	s.inquiries[inquiry.ID] = inquiry
	return nil
}

// GetInquiry retrieves an inquiry by ID
func (s *InMemoryStorage) GetInquiry(ctx context.Context, id string) (*domain.Inquiry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	inquiry, exists := s.inquiries[id]
	if !exists {
		return nil, fmt.Errorf("inquiry not found")
	}
	
	// Load associated reports
	var reports []*domain.Report
	for _, report := range s.reports {
		if report.InquiryID == id {
			reports = append(reports, report)
		}
	}
	inquiry.Reports = reports
	
	return inquiry, nil
}

// ListInquiries retrieves all inquiries with optional filtering
func (s *InMemoryStorage) ListInquiries(ctx context.Context, filters *domain.InquiryFilters) ([]*domain.Inquiry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var inquiries []*domain.Inquiry
	for _, inquiry := range s.inquiries {
		// Apply filters
		if filters != nil {
			if filters.Status != "" && string(inquiry.Status) != filters.Status {
				continue
			}
			if filters.Priority != "" && string(inquiry.Priority) != filters.Priority {
				continue
			}
			if filters.Service != "" {
				found := false
				for _, service := range inquiry.Services {
					if service == filters.Service {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}
		
		// Load associated reports
		var reports []*domain.Report
		for _, report := range s.reports {
			if report.InquiryID == inquiry.ID {
				reports = append(reports, report)
			}
		}
		inquiryCopy := *inquiry
		inquiryCopy.Reports = reports
		inquiries = append(inquiries, &inquiryCopy)
	}
	
	// Apply pagination
	if filters != nil {
		start := filters.Offset
		end := start + filters.Limit
		
		if start > len(inquiries) {
			return []*domain.Inquiry{}, nil
		}
		if end > len(inquiries) {
			end = len(inquiries)
		}
		if filters.Limit > 0 {
			inquiries = inquiries[start:end]
		}
	}
	
	return inquiries, nil
}

// GetInquiryCount returns the total count of inquiries matching the filters
func (s *InMemoryStorage) GetInquiryCount(ctx context.Context, filters *domain.InquiryFilters) (int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	count := int64(0)
	for _, inquiry := range s.inquiries {
		// Apply filters
		if filters != nil {
			if filters.Status != "" && string(inquiry.Status) != filters.Status {
				continue
			}
			if filters.Priority != "" && string(inquiry.Priority) != filters.Priority {
				continue
			}
			if filters.Service != "" {
				found := false
				for _, service := range inquiry.Services {
					if service == filters.Service {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}
		count++
	}
	
	return count, nil
}

// UpdateInquiry updates an existing inquiry
func (s *InMemoryStorage) UpdateInquiry(ctx context.Context, inquiry *domain.Inquiry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.inquiries[inquiry.ID]; !exists {
		return fmt.Errorf("inquiry not found")
	}
	
	inquiry.UpdatedAt = time.Now()
	s.inquiries[inquiry.ID] = inquiry
	return nil
}

// CreateReport stores a new report
func (s *InMemoryStorage) CreateReport(report *domain.Report) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if report.CreatedAt.IsZero() {
		report.CreatedAt = time.Now()
	}
	report.UpdatedAt = time.Now()
	s.reports[report.ID] = report
	return nil
}

// CreateReportWithContext stores a new report with context
func (s *InMemoryStorage) CreateReportWithContext(ctx context.Context, report *domain.Report) error {
	return s.CreateReport(report)
}

// GetReport retrieves a report by ID
func (s *InMemoryStorage) GetReport(ctx context.Context, id string) (*domain.Report, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	report, exists := s.reports[id]
	if !exists {
		return nil, fmt.Errorf("report not found")
	}
	
	return report, nil
}

// GetReportsByInquiryID retrieves all reports for a specific inquiry
func (s *InMemoryStorage) GetReportsByInquiryID(ctx context.Context, inquiryID string) ([]*domain.Report, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var reports []*domain.Report
	for _, report := range s.reports {
		if report.InquiryID == inquiryID {
			reports = append(reports, report)
		}
	}
	
	return reports, nil
}

// GetReportsByInquiry is an alias for GetReportsByInquiryID
func (s *InMemoryStorage) GetReportsByInquiry(ctx context.Context, inquiryID string) ([]*domain.Report, error) {
	return s.GetReportsByInquiryID(ctx, inquiryID)
}