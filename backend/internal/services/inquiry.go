package services

import (
	"context"
	"fmt"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/storage"
)

// inquiryService implements the InquiryService interface
type inquiryService struct {
	storage         *storage.InMemoryStorage
	reportGenerator interfaces.ReportService
}

// NewInquiryService creates a new inquiry service instance
func NewInquiryService(storage *storage.InMemoryStorage, reportGenerator interfaces.ReportService) interfaces.InquiryService {
	return &inquiryService{
		storage:         storage,
		reportGenerator: reportGenerator,
	}
}

// CreateInquiry creates a new inquiry and generates a report
func (s *inquiryService) CreateInquiry(ctx context.Context, req *domain.CreateInquiryRequest) (*domain.Inquiry, error) {
	// Create the inquiry first
	inquiry, err := s.storage.CreateInquiry(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create inquiry: %w", err)
	}

	// Try to generate a report using Bedrock
	// This should not fail the inquiry creation if it fails
	if s.reportGenerator != nil {
		report, err := s.reportGenerator.GenerateReport(ctx, inquiry)
		if err != nil {
			// Log the error but don't fail the inquiry creation
			fmt.Printf("Warning: Failed to generate report for inquiry %s: %v\n", inquiry.ID, err)
		} else {
			// Store the report
			if err := s.storage.CreateReport(report); err != nil {
				fmt.Printf("Warning: Failed to store report for inquiry %s: %v\n", inquiry.ID, err)
			}
		}
	}

	// Return the inquiry (potentially with reports if generation succeeded)
	return s.storage.GetInquiry(inquiry.ID)
}

// GetInquiry retrieves an inquiry by ID
func (s *inquiryService) GetInquiry(ctx context.Context, id string) (*domain.Inquiry, error) {
	return s.storage.GetInquiry(id)
}

// ListInquiries returns all inquiries with optional filters
func (s *inquiryService) ListInquiries(ctx context.Context, filters *domain.InquiryFilters) ([]*domain.Inquiry, error) {
	// For now, just return all inquiries (filtering can be added later)
	return s.storage.ListInquiries()
}

// UpdateInquiryStatus updates the status of an inquiry
func (s *inquiryService) UpdateInquiryStatus(ctx context.Context, id string, status domain.InquiryStatus) error {
	return fmt.Errorf("not implemented")
}

// AssignConsultant assigns a consultant to an inquiry
func (s *inquiryService) AssignConsultant(ctx context.Context, id string, consultantID string) error {
	return fmt.Errorf("not implemented")
}

// GetInquiryCount returns the count of inquiries with optional filters
func (s *inquiryService) GetInquiryCount(ctx context.Context, filters *domain.InquiryFilters) (int64, error) {
	inquiries, err := s.storage.ListInquiries()
	if err != nil {
		return 0, err
	}
	return int64(len(inquiries)), nil
}