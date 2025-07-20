package storage

import (
	"sync"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/google/uuid"
)

// InMemoryStorage provides in-memory storage for inquiries
type InMemoryStorage struct {
	inquiries map[string]*domain.Inquiry
	mutex     sync.RWMutex
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		inquiries: make(map[string]*domain.Inquiry),
	}
}

// CreateInquiry stores a new inquiry
func (s *InMemoryStorage) CreateInquiry(req *domain.CreateInquiryRequest) (*domain.Inquiry, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	inquiry := &domain.Inquiry{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Company:   req.Company,
		Phone:     req.Phone,
		Services:  req.Services,
		Message:   req.Message,
		Status:    "pending",
		Priority:  "medium",
		Source:    req.Source,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.inquiries[inquiry.ID] = inquiry
	return inquiry, nil
}

// GetInquiry retrieves an inquiry by ID
func (s *InMemoryStorage) GetInquiry(id string) (*domain.Inquiry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	inquiry, exists := s.inquiries[id]
	if !exists {
		return nil, nil // Not found
	}

	return inquiry, nil
}

// ListInquiries returns all inquiries
func (s *InMemoryStorage) ListInquiries() ([]*domain.Inquiry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	inquiries := make([]*domain.Inquiry, 0, len(s.inquiries))
	for _, inquiry := range s.inquiries {
		inquiries = append(inquiries, inquiry)
	}

	return inquiries, nil
}