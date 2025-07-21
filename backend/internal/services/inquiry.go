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
	emailService    interfaces.EmailService
}

// NewInquiryService creates a new inquiry service instance
func NewInquiryService(storage *storage.InMemoryStorage, reportGenerator interfaces.ReportService, emailService interfaces.EmailService) interfaces.InquiryService {
	return &inquiryService{
		storage:         storage,
		reportGenerator: reportGenerator,
		emailService:    emailService,
	}
}

// CreateInquiry creates a new inquiry and generates a report
func (s *inquiryService) CreateInquiry(ctx context.Context, req *domain.CreateInquiryRequest) (*domain.Inquiry, error) {
	// Create the inquiry first
	inquiry, err := s.storage.CreateInquiry(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create inquiry: %w", err)
	}
	
	// We'll send customer confirmation email after report generation
	// This ensures we can include report links and attachments if available

	// Try to generate a report using Bedrock
	// This should not fail the inquiry creation if it fails
	if s.reportGenerator != nil {
		report, err := s.reportGenerator.GenerateReport(ctx, inquiry)
		if err != nil {
			// Log the error but don't fail the inquiry creation
			fmt.Printf("Warning: Failed to generate report for inquiry %s: %v\n", inquiry.ID, err)
			
			// If report generation fails, send basic inquiry notification to internal team
			if s.emailService != nil {
				if err := s.emailService.SendInquiryNotification(ctx, inquiry); err != nil {
					fmt.Printf("Warning: Failed to send inquiry notification email for inquiry %s: %v\n", inquiry.ID, err)
				}
				
				// Send basic customer confirmation since we don't have a report
				if err := s.emailService.SendCustomerConfirmation(ctx, inquiry); err != nil {
					fmt.Printf("Warning: Failed to send customer confirmation email for inquiry %s: %v\n", inquiry.ID, err)
				}
			}
		} else {
			// Store the report
			if err := s.storage.CreateReport(report); err != nil {
				fmt.Printf("Warning: Failed to store report for inquiry %s: %v\n", inquiry.ID, err)
				
				// If report storage fails, send basic inquiry notification
				if s.emailService != nil {
					if err := s.emailService.SendInquiryNotification(ctx, inquiry); err != nil {
						fmt.Printf("Warning: Failed to send inquiry notification email for inquiry %s: %v\n", inquiry.ID, err)
					}
					
					// Send basic customer confirmation since we couldn't store the report
					if err := s.emailService.SendCustomerConfirmation(ctx, inquiry); err != nil {
						fmt.Printf("Warning: Failed to send customer confirmation email for inquiry %s: %v\n", inquiry.ID, err)
					}
				}
			} else {
				// Send comprehensive internal email with the report (this is the ONLY internal email)
				if s.emailService != nil {
					// Try to generate PDF for the report
					var pdfData []byte
					if s.reportGenerator != nil {
						pdfBytes, pdfErr := s.reportGenerator.GeneratePDF(ctx, inquiry, report)
						if pdfErr != nil {
							fmt.Printf("Warning: Failed to generate PDF for inquiry %s: %v\n", inquiry.ID, pdfErr)
						} else {
							pdfData = pdfBytes
						}
					}
					
					// Send internal notification with PDF if available
					if pdfData != nil && len(pdfData) > 0 {
						if err := s.emailService.SendReportEmailWithPDF(ctx, inquiry, report, pdfData); err != nil {
							fmt.Printf("Warning: Failed to send report email with PDF for inquiry %s: %v\n", inquiry.ID, err)
						}
					} else {
						if err := s.emailService.SendReportEmail(ctx, inquiry, report); err != nil {
							fmt.Printf("Warning: Failed to send report email for inquiry %s: %v\n", inquiry.ID, err)
						}
					}
					
					// Send customer confirmation with PDF if available
					if pdfData != nil && len(pdfData) > 0 {
						if err := s.emailService.SendCustomerConfirmationWithPDF(ctx, inquiry, report, pdfData); err != nil {
							fmt.Printf("Warning: Failed to send customer confirmation with PDF for inquiry %s: %v\n", inquiry.ID, err)
						}
					} else {
						// The regular SendCustomerConfirmation will include report links if a report exists
						if err := s.emailService.SendCustomerConfirmation(ctx, inquiry); err != nil {
							fmt.Printf("Warning: Failed to send customer confirmation email for inquiry %s: %v\n", inquiry.ID, err)
						}
					}
				}
			}
		}
	} else {
		// If no report generator available, send basic inquiry notification
		if s.emailService != nil {
			if err := s.emailService.SendInquiryNotification(ctx, inquiry); err != nil {
				fmt.Printf("Warning: Failed to send inquiry notification email for inquiry %s: %v\n", inquiry.ID, err)
			}
			
			// Send basic customer confirmation since we don't have a report generator
			if err := s.emailService.SendCustomerConfirmation(ctx, inquiry); err != nil {
				fmt.Printf("Warning: Failed to send customer confirmation email for inquiry %s: %v\n", inquiry.ID, err)
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