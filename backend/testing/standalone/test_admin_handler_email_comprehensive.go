package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/storage"
)

// MockEmailMetricsService provides a mock implementation for testing admin handler
type MockEmailMetricsService struct {
	mock.Mock
}

func (m *MockEmailMetricsService) GetEmailMetrics(ctx context.Context, timeRange domain.TimeRange) (*domain.EmailMetrics, error) {
	args := m.Called(ctx, timeRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmailMetrics), args.Error(1)
}

func (m *MockEmailMetricsService) GetEmailStatusByInquiry(ctx context.Context, inquiryID string) (*domain.EmailStatus, error) {
	args := m.Called(ctx, inquiryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EmailStatus), args.Error(1)
}

func (m *MockEmailMetricsService) GetEmailEventHistory(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.EmailEvent), args.Error(1)
}

func (m *MockEmailMetricsService) IsHealthy(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// MockInquiryService provides a mock implementation for testing
type MockInquiryService struct {
	mock.Mock
}

func (m *MockInquiryService) CreateInquiry(ctx context.Context, req *interfaces.CreateInquiryRequest) (*domain.Inquiry, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.Inquiry), args.Error(1)
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

func (m *MockInquiryService) UpdateInquiryStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockInquiryService) AssignConsultant(ctx context.Context, id string, consultantID string) error {
	args := m.Called(ctx, id, consultantID)
	return args.Error(0)
}

func (m *MockInquiryService) GetInquiryCount(ctx context.Context, filters *domain.InquiryFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

// MockReportService provides a mock implementation for testing
type MockReportService struct {
	mock.Mock
}

func (m *MockReportService) GenerateReport(ctx context.Context, inquiry *domain.Inquiry) (*domain.Report, error) {
	args := m.Called(ctx, inquiry)
	return args.Get(0).(*domain.Report), args.Error(1)
}

func (m *MockReportService) GenerateHTML(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (string, error) {
	args := m.Called(ctx, inquiry, report)
	return args.String(0), args.Error(1)
}

func (m *MockReportService) GeneratePDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) ([]byte, error) {
	args := m.Called(ctx, inquiry, report)
	return args.Get(0).([]byte), args.Error(1)
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

// MockEmailService provides a mock implementation for testing
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendReportEmail(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) error {
	args := m.Called(ctx, inquiry, report)
	return args.Error(0)
}

func (m *MockEmailService) SendInquiryNotification(ctx context.Context, inquiry *domain.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) SendCustomerConfirmation(ctx context.Context, inquiry *domain.Inquiry) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// TestAdminHandlerEmailEndpoints provides comprehensive tests for admin handler email-related endpoints
func TestAdminHandlerEmailEndpoints(t *testing.T) {
	t.Run("GetSystemMetrics", func(t *testing.T) {
		testGetSystemMetrics(t)
	})

	t.Run("GetEmailStatus", func(t *testing.T) {
		testGetEmailStatus(t)
	})

	t.Run("GetEmailEventHistory", func(t *testing.T) {
		testGetEmailEventHistory(t)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t)
	})
}

func setupTestAdminHandler(t *testing.T) (*handlers.AdminHandler, *MockEmailMetricsService, *MockInquiryService, *MockReportService, *MockEmailService) {
	// Create mocks
	mockEmailMetrics := &MockEmailMetricsService{}
	mockInquiry := &MockInquiryService{}
	mockReport := &MockReportService{}
	mockEmail := &MockEmailService{}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Create in-memory storage
	storage := storage.NewInMemoryStorage()

	// Create admin handler
	handler := handlers.NewAdminHandler(
		storage,
		mockInquiry,
		mockReport,
		mockEmail,
		mockEmailMetrics,
		logger,
	)

	require.NotNil(t, handler)

	return handler, mockEmailMetrics, mockInquiry, mockReport, mockEmail
}

func testGetSystemMetrics(t *testing.T) {
	handler, mockEmailMetrics, mockInquiry, _, _ := setupTestAdminHandler(t)

	t.Run("SuccessWithRealEmailMetrics", func(t *testing.T) {
		// Setup mocks
		mockInquiry.On("GetInquiryCount", mock.Anything, mock.Anything).Return(int64(25), nil).Once()
		mockInquiry.On("ListInquiries", mock.Anything, mock.Anything).Return([]*domain.Inquiry{
			{
				ID:        "inquiry-1",
				Reports:   []*domain.Report{{ID: "report-1"}, {ID: "report-2"}},
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
			{
				ID:        "inquiry-2",
				Reports:   []*domain.Report{{ID: "report-3"}},
				UpdatedAt: time.Now().Add(-2 * time.Hour),
			},
		}, nil).Once()

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailMetrics", mock.Anything, mock.MatchedBy(func(tr domain.TimeRange) bool {
			return !tr.Start.IsZero() && !tr.End.IsZero()
		})).Return(&domain.EmailMetrics{
			TotalEmails:     50,
			DeliveredEmails: 45,
			FailedEmails:    3,
			BouncedEmails:   1,
			SpamEmails:      1,
			DeliveryRate:    90.0,
			BounceRate:      2.0,
			SpamRate:        2.0,
			TimeRange:       "30d",
		}, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics?time_range=30d", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetSystemMetrics(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response, "data")

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(25), data["total_inquiries"])
		assert.Equal(t, float64(3), data["reports_generated"])
		assert.Equal(t, float64(50), data["emails_sent"])
		assert.Equal(t, 90.0, data["email_delivery_rate"])

		meta := response["meta"].(map[string]interface{})
		assert.True(t, meta["email_metrics_available"].(bool))

		// Verify no warnings
		assert.NotContains(t, response, "warnings")

		mockEmailMetrics.AssertExpectations(t)
		mockInquiry.AssertExpectations(t)
	})

	t.Run("SuccessWithUnhealthyEmailService", func(t *testing.T) {
		// Setup mocks
		mockInquiry.On("GetInquiryCount", mock.Anything, mock.Anything).Return(int64(10), nil).Once()
		mockInquiry.On("ListInquiries", mock.Anything, mock.Anything).Return([]*domain.Inquiry{}, nil).Once()

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(false).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetSystemMetrics(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(0), data["emails_sent"])
		assert.Equal(t, 0.0, data["email_delivery_rate"])

		meta := response["meta"].(map[string]interface{})
		assert.False(t, meta["email_metrics_available"].(bool))

		// Should have warnings
		assert.Contains(t, response, "warnings")
		warnings := response["warnings"].([]interface{})
		assert.Len(t, warnings, 1)
		assert.Contains(t, warnings[0].(string), "Email monitoring system is currently unavailable")

		mockEmailMetrics.AssertExpectations(t)
		mockInquiry.AssertExpectations(t)
	})

	t.Run("SuccessWithNoEmailMetricsService", func(t *testing.T) {
		// Create handler without email metrics service
		logger := logrus.New()
		logger.SetLevel(logrus.ErrorLevel)
		storage := storage.NewInMemoryStorage()

		handlerNoEmail := handlers.NewAdminHandler(
			storage,
			mockInquiry,
			&MockReportService{},
			&MockEmailService{},
			nil, // No email metrics service
			logger,
		)

		mockInquiry.On("GetInquiryCount", mock.Anything, mock.Anything).Return(int64(5), nil).Once()
		mockInquiry.On("ListInquiries", mock.Anything, mock.Anything).Return([]*domain.Inquiry{}, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handlerNoEmail.GetSystemMetrics(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(0), data["emails_sent"])
		assert.Equal(t, 0.0, data["email_delivery_rate"])

		meta := response["meta"].(map[string]interface{})
		assert.False(t, meta["email_metrics_available"].(bool))

		// Should have warnings
		assert.Contains(t, response, "warnings")
		warnings := response["warnings"].([]interface{})
		assert.Len(t, warnings, 1)
		assert.Contains(t, warnings[0].(string), "Email monitoring is not configured")

		mockInquiry.AssertExpectations(t)
	})

	t.Run("InvalidTimeRange", func(t *testing.T) {
		// Create request with invalid time range
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics?time_range=invalid", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetSystemMetrics(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Invalid time range parameter")
		assert.Equal(t, "INVALID_TIME_RANGE", response["code"])
	})

	t.Run("InquiryServiceError", func(t *testing.T) {
		mockInquiry.On("GetInquiryCount", mock.Anything, mock.Anything).Return(int64(0), fmt.Errorf("database error")).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetSystemMetrics(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Unable to retrieve system metrics")
		assert.Equal(t, "INQUIRY_COUNT_ERROR", response["code"])

		mockInquiry.AssertExpectations(t)
	})
}

func testGetEmailStatus(t *testing.T) {
	handler, mockEmailMetrics, mockInquiry, _, _ := setupTestAdminHandler(t)

	t.Run("SuccessWithEmailEvents", func(t *testing.T) {
		inquiryID := "inquiry-" + uuid.New().String()

		// Setup mocks
		mockInquiry.On("GetInquiry", mock.Anything, inquiryID).Return(&domain.Inquiry{
			ID:      inquiryID,
			Name:    "John Doe",
			Email:   "john@example.com",
			Company: "Example Corp",
		}, nil).Once()

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailStatusByInquiry", mock.Anything, inquiryID).Return(&domain.EmailStatus{
			InquiryID:       inquiryID,
			TotalEmailsSent: 2,
			CustomerEmail: &domain.EmailEvent{
				ID:             "event-1",
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "john@example.com",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now().Add(-1 * time.Hour),
				DeliveredAt:    timePtr(time.Now().Add(-50 * time.Minute)),
			},
			ConsultantEmail: &domain.EmailEvent{
				ID:             "event-2",
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeConsultantNotification,
				RecipientEmail: "consultant@cloudpartner.pro",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now().Add(-2 * time.Hour),
				DeliveredAt:    timePtr(time.Now().Add(-110 * time.Minute)),
			},
			LastEmailSent: timePtr(time.Now().Add(-1 * time.Hour)),
		}, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/"+inquiryID, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: inquiryID}}

		// Execute
		handler.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response, "data")

		data := response["data"].(map[string]interface{})
		assert.Equal(t, inquiryID, data["inquiry_id"])
		assert.Equal(t, "john@example.com", data["customer_email"])
		assert.Equal(t, "info@cloudpartner.pro", data["consultant_email"])
		assert.Equal(t, "delivered", data["status"])

		meta := response["meta"].(map[string]interface{})
		assert.Equal(t, float64(2), meta["total_emails_sent"])

		mockInquiry.AssertExpectations(t)
		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("InquiryNotFound", func(t *testing.T) {
		inquiryID := "non-existent-" + uuid.New().String()

		mockInquiry.On("GetInquiry", mock.Anything, inquiryID).Return(nil, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/"+inquiryID, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: inquiryID}}

		// Execute
		handler.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Inquiry not found")
		assert.Equal(t, "INQUIRY_NOT_FOUND", response["code"])
		assert.Equal(t, inquiryID, response["inquiry_id"])

		mockInquiry.AssertExpectations(t)
	})

	t.Run("NoEmailEvents", func(t *testing.T) {
		inquiryID := "inquiry-no-emails-" + uuid.New().String()

		mockInquiry.On("GetInquiry", mock.Anything, inquiryID).Return(&domain.Inquiry{
			ID:      inquiryID,
			Name:    "Jane Doe",
			Email:   "jane@example.com",
			Company: "Test Corp",
		}, nil).Once()

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailStatusByInquiry", mock.Anything, inquiryID).Return(nil, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/"+inquiryID, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: inquiryID}}

		// Execute
		handler.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "No email events found")
		assert.Equal(t, "NO_EMAIL_EVENTS", response["code"])

		mockInquiry.AssertExpectations(t)
		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("EmailMetricsServiceUnavailable", func(t *testing.T) {
		inquiryID := "inquiry-unavailable-" + uuid.New().String()

		mockInquiry.On("GetInquiry", mock.Anything, inquiryID).Return(&domain.Inquiry{
			ID: inquiryID,
		}, nil).Once()

		// Create handler without email metrics service
		logger := logrus.New()
		logger.SetLevel(logrus.ErrorLevel)
		storage := storage.NewInMemoryStorage()

		handlerNoEmail := handlers.NewAdminHandler(
			storage,
			mockInquiry,
			&MockReportService{},
			&MockEmailService{},
			nil, // No email metrics service
			logger,
		)

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/"+inquiryID, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: inquiryID}}

		// Execute
		handlerNoEmail.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Email monitoring is not configured")
		assert.Equal(t, "EMAIL_MONITORING_UNAVAILABLE", response["code"])

		mockInquiry.AssertExpectations(t)
	})

	t.Run("EmailMetricsServiceUnhealthy", func(t *testing.T) {
		inquiryID := "inquiry-unhealthy-" + uuid.New().String()

		mockInquiry.On("GetInquiry", mock.Anything, inquiryID).Return(&domain.Inquiry{
			ID: inquiryID,
		}, nil).Once()

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(false).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/"+inquiryID, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: inquiryID}}

		// Execute
		handler.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Email monitoring system is currently unavailable")
		assert.Equal(t, "EMAIL_MONITORING_UNHEALTHY", response["code"])

		mockInquiry.AssertExpectations(t)
		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("MissingInquiryID", func(t *testing.T) {
		// Create request without inquiry ID
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: ""}}

		// Execute
		handler.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Inquiry ID is required")
		assert.Equal(t, "MISSING_INQUIRY_ID", response["code"])
	})
}

func testGetEmailEventHistory(t *testing.T) {
	handler, mockEmailMetrics, _, _, _ := setupTestAdminHandler(t)

	t.Run("SuccessWithValidFilters", func(t *testing.T) {
		expectedEvents := []*domain.EmailEvent{
			{
				ID:             "event-1",
				InquiryID:      "inquiry-1",
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now().Add(-1 * time.Hour),
			},
			{
				ID:             "event-2",
				InquiryID:      "inquiry-2",
				EmailType:      domain.EmailTypeConsultantNotification,
				RecipientEmail: "consultant@cloudpartner.pro",
				Status:         domain.EmailStatusSent,
				SentAt:         time.Now().Add(-2 * time.Hour),
			},
		}

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailEventHistory", mock.Anything, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.Limit == 50 && filters.Offset == 0
		})).Return(expectedEvents, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?limit=50&offset=0", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))
		assert.Contains(t, response, "data")

		data := response["data"].([]interface{})
		assert.Len(t, data, 2)

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(50), pagination["limit"])
		assert.Equal(t, float64(0), pagination["offset"])

		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("SuccessWithEmailTypeFilter", func(t *testing.T) {
		expectedEvents := []*domain.EmailEvent{
			{
				ID:             "event-1",
				InquiryID:      "inquiry-1",
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now().Add(-1 * time.Hour),
			},
		}

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailEventHistory", mock.Anything, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.EmailType != nil && *filters.EmailType == domain.EmailTypeCustomerConfirmation
		})).Return(expectedEvents, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?email_type=customer_confirmation", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 1)

		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("SuccessWithStatusFilter", func(t *testing.T) {
		expectedEvents := []*domain.EmailEvent{
			{
				ID:             "event-1",
				InquiryID:      "inquiry-1",
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now().Add(-1 * time.Hour),
			},
		}

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailEventHistory", mock.Anything, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.Status != nil && *filters.Status == domain.EmailStatusDelivered
		})).Return(expectedEvents, nil).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?status=delivered", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 1)

		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("InvalidTimeRange", func(t *testing.T) {
		// Create request with invalid time range
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?time_range=invalid", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Invalid time range parameter")
		assert.Equal(t, "INVALID_TIME_RANGE", response["code"])
	})

	t.Run("InvalidEmailType", func(t *testing.T) {
		// Create request with invalid email type
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?email_type=invalid_type", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Invalid email type parameter")
		assert.Equal(t, "INVALID_EMAIL_TYPE", response["code"])
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		// Create request with invalid status
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?status=invalid_status", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Invalid status parameter")
		assert.Equal(t, "INVALID_STATUS", response["code"])
	})

	t.Run("EmailMetricsServiceUnavailable", func(t *testing.T) {
		// Create handler without email metrics service
		logger := logrus.New()
		logger.SetLevel(logrus.ErrorLevel)
		storage := storage.NewInMemoryStorage()

		handlerNoEmail := handlers.NewAdminHandler(
			storage,
			&MockInquiryService{},
			&MockReportService{},
			&MockEmailService{},
			nil, // No email metrics service
			logger,
		)

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handlerNoEmail.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Email monitoring is not configured")
		assert.Equal(t, "EMAIL_MONITORING_UNAVAILABLE", response["code"])
	})

	t.Run("EmailMetricsServiceUnhealthy", func(t *testing.T) {
		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(false).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Email monitoring system is currently unavailable")
		assert.Equal(t, "EMAIL_MONITORING_UNHEALTHY", response["code"])

		mockEmailMetrics.AssertExpectations(t)
	})
}

func testErrorHandling(t *testing.T) {
	handler, mockEmailMetrics, mockInquiry, _, _ := setupTestAdminHandler(t)

	t.Run("GetSystemMetricsWithEmailMetricsError", func(t *testing.T) {
		mockInquiry.On("GetInquiryCount", mock.Anything, mock.Anything).Return(int64(10), nil).Once()
		mockInquiry.On("ListInquiries", mock.Anything, mock.Anything).Return([]*domain.Inquiry{}, nil).Once()
		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailMetrics", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("database timeout")).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetSystemMetrics(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code) // Should still return OK but with warnings

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(0), data["emails_sent"])
		assert.Equal(t, 0.0, data["email_delivery_rate"])

		meta := response["meta"].(map[string]interface{})
		assert.False(t, meta["email_metrics_available"].(bool))

		// Should have warnings
		assert.Contains(t, response, "warnings")
		warnings := response["warnings"].([]interface{})
		assert.Len(t, warnings, 1)
		assert.Contains(t, warnings[0].(string), "Email metrics temporarily unavailable")

		mockInquiry.AssertExpectations(t)
		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("GetEmailStatusWithServiceError", func(t *testing.T) {
		inquiryID := "inquiry-error-" + uuid.New().String()

		mockInquiry.On("GetInquiry", mock.Anything, inquiryID).Return(&domain.Inquiry{
			ID: inquiryID,
		}, nil).Once()

		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailStatusByInquiry", mock.Anything, inquiryID).Return(nil, fmt.Errorf("service error")).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-status/"+inquiryID, nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "inquiryId", Value: inquiryID}}

		// Execute
		handler.GetEmailStatus(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Unable to retrieve email status")
		assert.Equal(t, "EMAIL_STATUS_RETRIEVAL_ERROR", response["code"])

		mockInquiry.AssertExpectations(t)
		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("GetEmailEventHistoryWithServiceError", func(t *testing.T) {
		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailEventHistory", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("service error")).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(string), "Unable to retrieve email event history")
		assert.Equal(t, "EMAIL_HISTORY_RETRIEVAL_ERROR", response["code"])

		mockEmailMetrics.AssertExpectations(t)
	})
}

func testEdgeCases(t *testing.T) {
	handler, mockEmailMetrics, mockInquiry, _, _ := setupTestAdminHandler(t)

	t.Run("GetSystemMetricsWithTimeout", func(t *testing.T) {
		mockInquiry.On("GetInquiryCount", mock.Anything, mock.Anything).Return(int64(5), nil).Once()
		mockInquiry.On("ListInquiries", mock.Anything, mock.Anything).Return([]*domain.Inquiry{}, nil).Once()

		// Simulate timeout by making health check slow
		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Run(func(args mock.Arguments) {
			time.Sleep(6 * time.Second) // Longer than timeout
		}).Once()

		// Create request
		req := httptest.NewRequest("GET", "/api/v1/admin/metrics", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetSystemMetrics(c)

		// Assert - should handle timeout gracefully
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		mockInquiry.AssertExpectations(t)
		// Note: mockEmailMetrics expectations may not be fully met due to timeout
	})

	t.Run("GetEmailEventHistoryWithLargeOffset", func(t *testing.T) {
		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailEventHistory", mock.Anything, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.Offset == 10000
		})).Return([]*domain.EmailEvent{}, nil).Once()

		// Create request with large offset
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?offset=10000", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 0) // No events at large offset

		mockEmailMetrics.AssertExpectations(t)
	})

	t.Run("GetEmailEventHistoryWithMaxLimit", func(t *testing.T) {
		mockEmailMetrics.On("IsHealthy", mock.Anything).Return(true).Once()
		mockEmailMetrics.On("GetEmailEventHistory", mock.Anything, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.Limit == 1000 // Should be capped at 1000
		})).Return([]*domain.EmailEvent{}, nil).Once()

		// Create request with limit exceeding maximum
		req := httptest.NewRequest("GET", "/api/v1/admin/email-events?limit=5000", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Execute
		handler.GetEmailEventHistory(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.True(t, response["success"].(bool))

		pagination := response["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1000), pagination["limit"]) // Should be capped

		mockEmailMetrics.AssertExpectations(t)
	})
}

// Helper functions

func timePtr(t time.Time) *time.Time {
	return &t
}

// Main function for running tests standalone
func main() {
	fmt.Println("=== Comprehensive Admin Handler Email Tests ===")

	// Note: This would normally be run with `go test` command
	// This main function is for demonstration purposes

	fmt.Println("Run with: go test -v ./test_admin_handler_email_comprehensive.go")
	fmt.Println("Or integrate into your test suite")
}
