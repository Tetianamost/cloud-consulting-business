package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/storage"
)

// MockInquiryService is a mock implementation of the InquiryService interface
type MockInquiryService struct {
	mock.Mock
}

func (m *MockInquiryService) CreateInquiry(ctx context.Context, req *interfaces.CreateInquiryRequest) (*domain.Inquiry, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.Inquiry), args.Error(1)
}

func (m *MockInquiryService) GetInquiry(ctx context.Context, id string) (*domain.Inquiry, error) {
	args := m.Called(ctx, id)
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

func (m *MockInquiryService) UpdateInquiryStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockInquiryService) AssignConsultant(ctx context.Context, id string, consultantID string) error {
	args := m.Called(ctx, id, consultantID)
	return args.Error(0)
}

// MockReportService is a mock implementation of the ReportService interface
type MockReportService struct {
	mock.Mock
}

func (m *MockReportService) GenerateReport(ctx context.Context, inquiry *domain.Inquiry) (*domain.Report, error) {
	args := m.Called(ctx, inquiry)
	return args.Get(0).(*domain.Report), args.Error(1)
}

func (m *MockReportService) GeneratePDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) ([]byte, error) {
	args := m.Called(ctx, inquiry, report)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockReportService) GenerateHTML(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (string, error) {
	args := m.Called(ctx, inquiry, report)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockReportService) GetReport(ctx context.Context, id string) (*domain.Report, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Report), args.Error(1)
}

func (m *MockReportService) GetReportsByInquiry(ctx context.Context, inquiryID string) ([]*domain.Report, error) {
	args := m.Called(ctx, inquiryID)
	return args.Get(0).([]*domain.Report), args.Error(1)
}

func (m *MockReportService) UpdateReportStatus(ctx context.Context, id string, status domain.ReportStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockReportService) GetReportTemplate(serviceType domain.ServiceType) (*interfaces.ReportTemplate, error) {
	args := m.Called(serviceType)
	return args.Get(0).(*interfaces.ReportTemplate), args.Error(1)
}

func (m *MockReportService) ValidateReport(report *domain.Report) error {
	args := m.Called(report)
	return args.Error(0)
}

// MockEmailService is a mock implementation of the EmailService interface
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) SendReportEmail(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) error {
	args := m.Called(ctx, inquiry, report)
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

func TestDownloadReport_RouteParameterExtraction(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockInquiryService := &MockInquiryService{}
	mockReportService := &MockReportService{}
	mockEmailService := &MockEmailService{}
	memStorage := storage.NewInMemoryStorage()

	// Create admin handler
	adminHandler := NewAdminHandler(memStorage, mockInquiryService, mockReportService, mockEmailService)

	// Create test inquiry with report
	testInquiry := &domain.Inquiry{
		ID:      "test-inquiry-123",
		Name:    "Test User",
		Email:   "test@example.com",
		Company: "Test Company",
		Reports: []*domain.Report{
			{
				ID:        "test-report-456",
				Type:      domain.ReportTypeAssessment,
				Status:    domain.ReportStatusGenerated,
				Content:   "Test report content",
				CreatedAt: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		inquiryID      string
		format         string
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful PDF download with correct parameters",
			inquiryID: "test-inquiry-123",
			format:    "pdf",
			setupMocks: func() {
				mockInquiryService.On("GetInquiry", mock.Anything, "test-inquiry-123").Return(testInquiry, nil)
				mockReportService.On("GeneratePDF", mock.Anything, testInquiry, testInquiry.Reports[0]).Return([]byte("fake pdf content"), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "successful HTML download with correct parameters",
			inquiryID: "test-inquiry-123",
			format:    "html",
			setupMocks: func() {
				mockInquiryService.On("GetInquiry", mock.Anything, "test-inquiry-123").Return(testInquiry, nil)
				mockReportService.On("GenerateHTML", mock.Anything, testInquiry, testInquiry.Reports[0]).Return("<html>fake html content</html>", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid format parameter",
			inquiryID:      "test-inquiry-123",
			format:         "invalid",
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid format",
		},
		{
			name:      "inquiry not found",
			inquiryID: "nonexistent-inquiry",
			format:    "pdf",
			setupMocks: func() {
				mockInquiryService.On("GetInquiry", mock.Anything, "nonexistent-inquiry").Return((*domain.Inquiry)(nil), nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Inquiry not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockInquiryService.ExpectedCalls = nil
			mockReportService.ExpectedCalls = nil

			// Setup mocks for this test
			tt.setupMocks()

			// Create router and register route
			router := gin.New()
			router.GET("/api/v1/admin/reports/:inquiryId/download/:format", adminHandler.DownloadReport)

			// Create request
			url := "/api/v1/admin/reports/" + tt.inquiryID + "/download/" + tt.format
			req, err := http.NewRequest("GET", url, nil)
			assert.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			// Verify mocks were called as expected
			mockInquiryService.AssertExpectations(t)
			mockReportService.AssertExpectations(t)
		})
	}
}

func TestDownloadReport_RouteRegistration(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockInquiryService := &MockInquiryService{}
	mockReportService := &MockReportService{}
	mockEmailService := &MockEmailService{}
	memStorage := storage.NewInMemoryStorage()

	// Create admin handler
	adminHandler := NewAdminHandler(memStorage, mockInquiryService, mockReportService, mockEmailService)

	// Create router and register the new route
	router := gin.New()
	router.GET("/api/v1/admin/reports/:inquiryId/download/:format", adminHandler.DownloadReport)

	// Test that the route is registered correctly
	routes := router.Routes()

	// Find our route
	var foundRoute gin.RouteInfo
	for _, route := range routes {
		if route.Path == "/api/v1/admin/reports/:inquiryId/download/:format" && route.Method == "GET" {
			foundRoute = route
			break
		}
	}

	// Assert route was found
	assert.NotEmpty(t, foundRoute.Path, "Route should be registered")
	assert.Equal(t, "GET", foundRoute.Method, "Route should be GET method")
	assert.Equal(t, "/api/v1/admin/reports/:inquiryId/download/:format", foundRoute.Path, "Route path should match expected pattern")
}

func TestDownloadReport_ParameterValidation(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockInquiryService := &MockInquiryService{}
	mockReportService := &MockReportService{}
	mockEmailService := &MockEmailService{}
	memStorage := storage.NewInMemoryStorage()

	// Create admin handler
	adminHandler := NewAdminHandler(memStorage, mockInquiryService, mockReportService, mockEmailService)

	// Create router and register route
	router := gin.New()
	router.GET("/api/v1/admin/reports/:inquiryId/download/:format", adminHandler.DownloadReport)

	// Test various parameter combinations
	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		shouldContain  string
	}{
		{
			name:           "valid parameters",
			url:            "/api/v1/admin/reports/test-123/download/pdf",
			expectedStatus: http.StatusNotFound, // Will fail at inquiry lookup, but parameters are extracted
			shouldContain:  "Inquiry not found",
		},
		{
			name:           "invalid format parameter",
			url:            "/api/v1/admin/reports/test-123/download/invalid",
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Invalid format",
		},
		{
			name:           "empty inquiry ID",
			url:            "/api/v1/admin/reports//download/pdf",
			expectedStatus: http.StatusNotFound, // Gin routing will not match
		},
		{
			name:           "empty format",
			url:            "/api/v1/admin/reports/test-123/download/",
			expectedStatus: http.StatusNotFound, // Gin routing will not match
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mocks
			mockInquiryService.ExpectedCalls = nil

			// Setup mock for inquiry lookup (will be called for valid parameter cases)
			if tc.expectedStatus != http.StatusNotFound || tc.shouldContain == "Inquiry not found" {
				mockInquiryService.On("GetInquiry", mock.Anything, mock.AnythingOfType("string")).Return((*domain.Inquiry)(nil), nil)
			}

			// Create request
			req, err := http.NewRequest("GET", tc.url, nil)
			assert.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert response contains expected content
			if tc.shouldContain != "" {
				assert.Contains(t, w.Body.String(), tc.shouldContain)
			}
		})
	}
}
