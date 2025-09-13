package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/storage"
)

// MockInquiryService for testing
type MockInquiryService struct {
	mock.Mock
}

func (m *MockInquiryService) CreateInquiry(ctx context.Context, inquiry *domain.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockInquiryService) GetInquiry(ctx context.Context, id string) (*domain.Inquiry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Inquiry), args.Error(1)
}

func (m *MockInquiryService) ListInquiries(ctx context.Context, filters *domain.InquiryFilters) ([]*domain.Inquiry, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*domain.Inquiry), args.Error(1)
}

func (m *MockInquiryService) GetInquiryCount(ctx context.Context, filters *domain.InquiryFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

// MockReportService for testing
type MockReportService struct {
	mock.Mock
}

func (m *MockReportService) GenerateReport(ctx context.Context, inquiry *domain.Inquiry) (*domain.Report, error) {
	args := m.Called(ctx, inquiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Report), args.Error(1)
}

func (m *MockReportService) GeneratePDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) ([]byte, error) {
	args := m.Called(ctx, inquiry, report)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockReportService) GenerateHTML(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (string, error) {
	args := m.Called(ctx, inquiry, report)
	return args.Get(0).(string), args.Error(1)
}

// MockEmailService for testing
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendReportEmail(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) error {
	args := m.Called(ctx, inquiry, report)
	return args.Error(0)
}

func (m *MockEmailService) SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) SendInquiryNotification(ctx context.Context, inquiry *domain.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestEnhancedDownloadErrorHandling(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create logger with buffer to capture logs
	var logBuffer bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&logBuffer)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create mocks
	mockInquiryService := &MockInquiryService{}
	mockReportService := &MockReportService{}
	mockEmailService := &MockEmailService{}
	memStorage := storage.NewInMemoryStorage()

	// Create handler
	adminHandler := handlers.NewAdminHandler(
		memStorage,
		mockInquiryService,
		mockReportService,
		mockEmailService,
		logger,
	)

	// Setup router
	router := gin.New()
	router.GET("/api/v1/admin/reports/:inquiryId/download/:format", adminHandler.DownloadReport)

	tests := []struct {
		name           string
		inquiryID      string
		format         string
		setupMocks     func()
		expectedStatus int
		expectedCode   string
		expectedMsg    string
		checkLogs      bool
	}{
		{
			name:           "invalid format parameter",
			inquiryID:      "test-123",
			format:         "xml",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_FORMAT",
			expectedMsg:    "Invalid format parameter",
			checkLogs:      true,
		},
		{
			name:           "empty inquiry ID",
			inquiryID:      "",
			format:         "pdf",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
			expectedMsg:    "inquiry ID is required",
			checkLogs:      true,
		},
		{
			name:      "inquiry not found",
			inquiryID: "missing-123",
			format:    "pdf",
			setupMocks: func() {
				mockInquiryService.On("GetInquiry", mock.Anything, "missing-123").Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "INQUIRY_NOT_FOUND",
			expectedMsg:    "Inquiry not found",
			checkLogs:      true,
		},
		{
			name:      "inquiry service error",
			inquiryID: "error-123",
			format:    "pdf",
			setupMocks: func() {
				mockInquiryService.On("GetInquiry", mock.Anything, "error-123").Return(nil, fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
			expectedMsg:    "Internal server error",
			checkLogs:      true,
		},
		{
			name:      "no reports available",
			inquiryID: "empty-123",
			format:    "html",
			setupMocks: func() {
				inquiry := &domain.Inquiry{
					ID:      "empty-123",
					Reports: []*domain.Report{}, // Empty reports
				}
				mockInquiryService.On("GetInquiry", mock.Anything, "empty-123").Return(inquiry, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   "NO_REPORTS",
			expectedMsg:    "No reports available",
			checkLogs:      true,
		},
		{
			name:      "PDF generation failure",
			inquiryID: "pdf-fail-123",
			format:    "pdf",
			setupMocks: func() {
				report := &domain.Report{
					ID:        "report-1",
					Type:      domain.ReportTypeAssessment,
					CreatedAt: time.Now(),
				}
				inquiry := &domain.Inquiry{
					ID:      "pdf-fail-123",
					Reports: []*domain.Report{report},
				}
				mockInquiryService.On("GetInquiry", mock.Anything, "pdf-fail-123").Return(inquiry, nil)
				mockReportService.On("GeneratePDF", mock.Anything, inquiry, report).Return(nil, fmt.Errorf("PDF generation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "REPORT_GENERATION_ERROR",
			expectedMsg:    "Report generation failed",
			checkLogs:      true,
		},
		{
			name:      "HTML generation failure",
			inquiryID: "html-fail-123",
			format:    "html",
			setupMocks: func() {
				report := &domain.Report{
					ID:        "report-2",
					Type:      domain.ReportTypeMigration,
					CreatedAt: time.Now(),
				}
				inquiry := &domain.Inquiry{
					ID:      "html-fail-123",
					Reports: []*domain.Report{report},
				}
				mockInquiryService.On("GetInquiry", mock.Anything, "html-fail-123").Return(inquiry, nil)
				mockReportService.On("GenerateHTML", mock.Anything, inquiry, report).Return("", fmt.Errorf("HTML generation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "REPORT_GENERATION_ERROR",
			expectedMsg:    "Report generation failed",
			checkLogs:      true,
		},
		{
			name:      "successful PDF download",
			inquiryID: "success-123",
			format:    "pdf",
			setupMocks: func() {
				report := &domain.Report{
					ID:        "report-3",
					Type:      domain.ReportTypeOptimization,
					CreatedAt: time.Now(),
				}
				inquiry := &domain.Inquiry{
					ID:      "success-123",
					Company: "Test Company",
					Reports: []*domain.Report{report},
				}
				pdfData := []byte("fake PDF content")
				mockInquiryService.On("GetInquiry", mock.Anything, "success-123").Return(inquiry, nil)
				mockReportService.On("GeneratePDF", mock.Anything, inquiry, report).Return(pdfData, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			expectedMsg:    "",
			checkLogs:      true,
		},
		{
			name:      "successful HTML download",
			inquiryID: "success-html-123",
			format:    "html",
			setupMocks: func() {
				report := &domain.Report{
					ID:        "report-4",
					Type:      domain.ReportTypeArchitecture,
					CreatedAt: time.Now(),
				}
				inquiry := &domain.Inquiry{
					ID:      "success-html-123",
					Company: "HTML Test Company",
					Reports: []*domain.Report{report},
				}
				htmlContent := "<html><body>Test Report</body></html>"
				mockInquiryService.On("GetInquiry", mock.Anything, "success-html-123").Return(inquiry, nil)
				mockReportService.On("GenerateHTML", mock.Anything, inquiry, report).Return(htmlContent, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			expectedMsg:    "",
			checkLogs:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log buffer and reset mocks
			logBuffer.Reset()
			mockInquiryService.ExpectedCalls = nil
			mockReportService.ExpectedCalls = nil

			// Setup mocks
			tt.setupMocks()

			// Create request
			url := fmt.Sprintf("/api/v1/admin/reports/%s/download/%s", tt.inquiryID, tt.format)
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-Trace-ID", fmt.Sprintf("test-trace-%s", tt.name))

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code, "Status code mismatch for test: %s", tt.name)

			if tt.expectedStatus == http.StatusOK {
				// For successful downloads, check headers
				if tt.format == "pdf" {
					assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
					assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
					assert.Contains(t, w.Header().Get("Content-Disposition"), ".pdf")
				} else if tt.format == "html" {
					assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
					assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
					assert.Contains(t, w.Header().Get("Content-Disposition"), ".html")
				}
			} else {
				// For error responses, check JSON structure
				var response handlers.APIErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Failed to parse error response for test: %s", tt.name)

				assert.False(t, response.Success)
				assert.Equal(t, tt.expectedCode, string(response.Code))
				assert.Contains(t, response.Error, tt.expectedMsg)
				assert.NotEmpty(t, response.Timestamp)
				assert.NotEmpty(t, response.TraceID)
				assert.Contains(t, response.TraceID, "test-trace")

				// Check context information
				if tt.inquiryID != "" && response.Context != nil {
					assert.Equal(t, tt.inquiryID, response.Context["inquiry_id"])
				}
				if tt.format != "" && response.Context != nil {
					assert.Equal(t, tt.format, response.Context["format"])
				}
			}

			// Check logs if required
			if tt.checkLogs {
				logOutput := logBuffer.String()
				assert.NotEmpty(t, logOutput, "Expected log output for test: %s", tt.name)

				if tt.expectedStatus == http.StatusOK {
					// Check for success log
					assert.Contains(t, logOutput, "success", "Expected success log for test: %s", tt.name)
					assert.Contains(t, logOutput, tt.inquiryID, "Expected inquiry ID in log for test: %s", tt.name)
				} else {
					// Check for error log
					assert.Contains(t, logOutput, tt.expectedCode, "Expected error code in log for test: %s", tt.name)
					if tt.inquiryID != "" {
						assert.Contains(t, logOutput, tt.inquiryID, "Expected inquiry ID in log for test: %s", tt.name)
					}
				}
			}

			// Verify mock expectations
			mockInquiryService.AssertExpectations(t)
			mockReportService.AssertExpectations(t)
		})
	}
}

func main() {
	// Run the test
	testing.Main(func(pat, str string) (bool, error) { return true, nil }, []testing.InternalTest{
		{
			Name: "TestEnhancedDownloadErrorHandling",
			F:    TestEnhancedDownloadErrorHandling,
		},
	}, nil, nil)

	fmt.Println("Enhanced download error handling tests completed successfully!")
}
