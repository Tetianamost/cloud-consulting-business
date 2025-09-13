package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	storage             *storage.InMemoryStorage
	inquiryService      interfaces.InquiryService
	reportService       interfaces.ReportService
	emailService        interfaces.EmailService
	emailMetricsService interfaces.EmailMetricsService
	logger              *logrus.Logger
	errorHandler        *ErrorHandler
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	storage *storage.InMemoryStorage,
	inquiryService interfaces.InquiryService,
	reportService interfaces.ReportService,
	emailService interfaces.EmailService,
	emailMetricsService interfaces.EmailMetricsService,
	logger *logrus.Logger,
) *AdminHandler {
	return &AdminHandler{
		storage:             storage,
		inquiryService:      inquiryService,
		reportService:       reportService,
		emailService:        emailService,
		emailMetricsService: emailMetricsService,
		logger:              logger,
		errorHandler:        NewErrorHandler(logger),
	}
}

// AdminInquiryFilters represents filters for admin inquiry listing
type AdminInquiryFilters struct {
	Status   string `form:"status"`
	Priority string `form:"priority"`
	Service  string `form:"service"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
}

// AdminInquiriesResponse represents the response for admin inquiries listing
type AdminInquiriesResponse struct {
	Success bool              `json:"success"`
	Data    []*domain.Inquiry `json:"data"`
	Count   int               `json:"count"`
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	Pages   int               `json:"pages"`
}

// SystemMetrics represents system metrics for the admin dashboard
type SystemMetrics struct {
	TotalInquiries    int64   `json:"total_inquiries"`
	ReportsGenerated  int64   `json:"reports_generated"`
	EmailsSent        int64   `json:"emails_sent"`
	EmailDeliveryRate float64 `json:"email_delivery_rate"`
	AvgReportGenTime  float64 `json:"avg_report_gen_time_ms"`
	SystemUptime      string  `json:"system_uptime"`
	LastProcessedAt   string  `json:"last_processed_at,omitempty"`
}

// EmailStatus represents email delivery status
type EmailStatus struct {
	InquiryID       string    `json:"inquiry_id"`
	CustomerEmail   string    `json:"customer_email"`
	ConsultantEmail string    `json:"consultant_email"`
	Status          string    `json:"status"`
	SentAt          time.Time `json:"sent_at,omitempty"`
	DeliveredAt     time.Time `json:"delivered_at,omitempty"`
	ErrorMessage    string    `json:"error_message,omitempty"`
}

// ListInquiries handles GET /api/v1/admin/inquiries
func (h *AdminHandler) ListInquiries(c *gin.Context) {
	var filters AdminInquiryFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	// Convert to domain filters
	domainFilters := &domain.InquiryFilters{
		Status:   filters.Status,
		Priority: filters.Priority,
		Service:  filters.Service,
	}

	// Set default limit and offset if not provided
	if filters.Limit <= 0 {
		filters.Limit = 10
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Cap maximum limit
	}
	domainFilters.Limit = filters.Limit
	domainFilters.Offset = filters.Offset

	// Get total count for pagination
	total, err := h.inquiryService.GetInquiryCount(c.Request.Context(), domainFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to count inquiries",
		})
		return
	}

	// Get inquiries with pagination
	inquiries, err := h.inquiryService.ListInquiries(c.Request.Context(), domainFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiries",
		})
		return
	}

	// Calculate pagination info
	page := (filters.Offset / filters.Limit) + 1
	pages := int(total) / filters.Limit
	if int(total)%filters.Limit > 0 {
		pages++
	}

	response := AdminInquiriesResponse{
		Success: true,
		Data:    inquiries,
		Count:   len(inquiries),
		Total:   total,
		Page:    page,
		Pages:   pages,
	}

	c.JSON(http.StatusOK, response)
}

// AdminReportFilters represents filters for admin report listing
type AdminReportFilters struct {
	Status   string `form:"status"`
	Type     string `form:"type"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
}

// AdminReportsResponse represents the response for admin reports listing
type AdminReportsResponse struct {
	Success bool                 `json:"success"`
	Data    []*ReportWithInquiry `json:"data"`
	Count   int                  `json:"count"`
	Total   int64                `json:"total"`
	Page    int                  `json:"page"`
	Pages   int                  `json:"pages"`
}

// ReportWithInquiry represents a report with its associated inquiry information
type ReportWithInquiry struct {
	*domain.Report
	Inquiry *domain.Inquiry `json:"inquiry"`
}

// ListReports handles GET /api/v1/admin/reports
func (h *AdminHandler) ListReports(c *gin.Context) {
	var filters AdminReportFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	// Set default limit and offset if not provided
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100 // Cap maximum limit
	}

	// Get all inquiries to extract reports
	inquiries, err := h.inquiryService.ListInquiries(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiries",
		})
		return
	}

	// Extract all reports with their inquiry information
	// Initialize with empty slice to ensure JSON marshals as [] not null
	allReports := make([]*ReportWithInquiry, 0)
	for _, inquiry := range inquiries {
		if len(inquiry.Reports) > 0 {
			for _, report := range inquiry.Reports {
				reportWithInquiry := &ReportWithInquiry{
					Report:  report,
					Inquiry: inquiry,
				}

				// Apply filters
				if filters.Status != "" && string(report.Status) != filters.Status {
					continue
				}
				if filters.Type != "" && string(report.Type) != filters.Type {
					continue
				}

				allReports = append(allReports, reportWithInquiry)
			}
		}
	}

	// Sort reports by creation date (newest first)
	for i := 0; i < len(allReports)-1; i++ {
		for j := i + 1; j < len(allReports); j++ {
			if allReports[i].CreatedAt.Before(allReports[j].CreatedAt) {
				allReports[i], allReports[j] = allReports[j], allReports[i]
			}
		}
	}

	// Apply pagination
	total := int64(len(allReports))
	start := filters.Offset
	end := start + filters.Limit

	if start > len(allReports) {
		start = len(allReports)
	}
	if end > len(allReports) {
		end = len(allReports)
	}

	paginatedReports := allReports[start:end]

	// Calculate pagination info
	page := (filters.Offset / filters.Limit) + 1
	pages := int(total) / filters.Limit
	if int(total)%filters.Limit > 0 {
		pages++
	}

	response := AdminReportsResponse{
		Success: true,
		Data:    paginatedReports,
		Count:   len(paginatedReports),
		Total:   total,
		Page:    page,
		Pages:   pages,
	}

	c.JSON(http.StatusOK, response)
}

// GetReport handles GET /api/v1/admin/reports/:reportId
func (h *AdminHandler) GetReport(c *gin.Context) {
	reportID := c.Param("reportId")

	// Get all inquiries to find the report
	inquiries, err := h.inquiryService.ListInquiries(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiries",
		})
		return
	}

	// Find the report by ID
	var foundReport *domain.Report
	var foundInquiry *domain.Inquiry

	for _, inquiry := range inquiries {
		for _, report := range inquiry.Reports {
			if report.ID == reportID {
				foundReport = report
				foundInquiry = inquiry
				break
			}
		}
		if foundReport != nil {
			break
		}
	}

	if foundReport == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Report not found",
		})
		return
	}

	reportWithInquiry := &ReportWithInquiry{
		Report:  foundReport,
		Inquiry: foundInquiry,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reportWithInquiry,
	})
}

// GetSystemMetrics handles GET /api/v1/admin/metrics
func (h *AdminHandler) GetSystemMetrics(c *gin.Context) {
	// Get time range from query params (default to last 30 days)
	timeRangeParam := c.DefaultQuery("time_range", "30d")
	timeRange, err := h.parseTimeRange(timeRangeParam)
	if err != nil {
		h.logger.WithError(err).WithField("time_range", timeRangeParam).Error("Invalid time range parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid time range parameter: " + err.Error(),
			"code":    "INVALID_TIME_RANGE",
		})
		return
	}

	// Get total inquiries with error handling
	totalInquiries, err := h.inquiryService.GetInquiryCount(c.Request.Context(), nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve inquiry count")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to retrieve system metrics at this time",
			"code":    "INQUIRY_COUNT_ERROR",
			"details": "Failed to access inquiry data",
		})
		return
	}

	// Calculate reports generated with error handling
	inquiries, err := h.inquiryService.ListInquiries(c.Request.Context(), nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve inquiries for metrics calculation")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to calculate report metrics",
			"code":    "INQUIRY_LIST_ERROR",
			"details": "Failed to access inquiry data for report calculation",
		})
		return
	}

	var reportsGenerated int64 = 0
	var lastProcessed time.Time

	for _, inquiry := range inquiries {
		if len(inquiry.Reports) > 0 {
			reportsGenerated += int64(len(inquiry.Reports))
		}
		if inquiry.UpdatedAt.After(lastProcessed) {
			lastProcessed = inquiry.UpdatedAt
		}
	}

	// Get real email metrics with comprehensive error handling
	var emailsSent int64 = 0
	var emailDeliveryRate float64 = 0.0
	var emailMetricsAvailable = false
	var emailMetricsError string

	if h.emailMetricsService != nil {
		// Check if email metrics service is healthy first
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if h.emailMetricsService.IsHealthy(ctx) {
			emailMetrics, err := h.emailMetricsService.GetEmailMetrics(ctx, *timeRange)
			if err != nil {
				h.logger.WithError(err).WithFields(logrus.Fields{
					"time_range_start": timeRange.Start,
					"time_range_end":   timeRange.End,
				}).Error("Failed to get email metrics from healthy service")
				emailMetricsError = "Email metrics temporarily unavailable"
			} else {
				emailsSent = emailMetrics.TotalEmails
				emailDeliveryRate = emailMetrics.DeliveryRate
				emailMetricsAvailable = true
				h.logger.WithFields(logrus.Fields{
					"emails_sent":         emailsSent,
					"email_delivery_rate": emailDeliveryRate,
				}).Debug("Successfully retrieved email metrics")
			}
		} else {
			h.logger.Warn("Email metrics service is unhealthy")
			emailMetricsError = "Email monitoring system is currently unavailable"
		}
	} else {
		h.logger.Warn("Email metrics service not configured")
		emailMetricsError = "Email monitoring is not configured"
	}

	// Calculate average report generation time (for demo, use a fixed value)
	avgReportGenTime := 1250.0 // 1.25 seconds in ms

	// Calculate system uptime (for demo, use a fixed value)
	systemUptime := "3d 7h 22m"

	metrics := SystemMetrics{
		TotalInquiries:    totalInquiries,
		ReportsGenerated:  reportsGenerated,
		EmailsSent:        emailsSent,
		EmailDeliveryRate: emailDeliveryRate,
		AvgReportGenTime:  avgReportGenTime,
		SystemUptime:      systemUptime,
	}

	if !lastProcessed.IsZero() {
		metrics.LastProcessedAt = lastProcessed.Format(time.RFC3339)
	}

	// Prepare response with data availability information
	response := gin.H{
		"success": true,
		"data":    metrics,
		"meta": gin.H{
			"email_metrics_available": emailMetricsAvailable,
			"time_range":              timeRangeParam,
		},
	}

	// Add warning if email metrics are not available
	if !emailMetricsAvailable {
		response["warnings"] = []string{emailMetricsError}
		h.logger.WithField("email_metrics_error", emailMetricsError).Info("Returning system metrics with email metrics warning")
	}

	c.JSON(http.StatusOK, response)
}

// GetEmailStatus handles GET /api/v1/admin/email-status/:inquiryId
func (h *AdminHandler) GetEmailStatus(c *gin.Context) {
	inquiryID := c.Param("inquiryId")

	// Validate inquiry ID
	if inquiryID == "" {
		h.logger.Error("Email status request missing inquiry ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Inquiry ID is required",
			"code":    "MISSING_INQUIRY_ID",
		})
		return
	}

	// Get the inquiry to verify it exists with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	inquiry, err := h.inquiryService.GetInquiry(ctx, inquiryID)
	if err != nil {
		h.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to retrieve inquiry for email status")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to retrieve inquiry information",
			"code":    "INQUIRY_RETRIEVAL_ERROR",
			"details": "Failed to access inquiry data",
		})
		return
	}

	if inquiry == nil {
		h.logger.WithField("inquiry_id", inquiryID).Warn("Inquiry not found for email status request")
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error":      "Inquiry not found",
			"code":       "INQUIRY_NOT_FOUND",
			"inquiry_id": inquiryID,
		})
		return
	}

	// Check if email metrics service is available and healthy
	if h.emailMetricsService == nil {
		h.logger.WithField("inquiry_id", inquiryID).Error("Email metrics service not configured")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email monitoring is not configured",
			"code":    "EMAIL_MONITORING_UNAVAILABLE",
			"details": "Email status tracking is not available",
		})
		return
	}

	// Check service health with timeout
	healthCtx, healthCancel := context.WithTimeout(ctx, 3*time.Second)
	defer healthCancel()

	if !h.emailMetricsService.IsHealthy(healthCtx) {
		h.logger.WithField("inquiry_id", inquiryID).Error("Email metrics service is unhealthy")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email monitoring system is currently unavailable",
			"code":    "EMAIL_MONITORING_UNHEALTHY",
			"details": "Email status service is experiencing issues",
		})
		return
	}

	// Get real email status with timeout
	statusCtx, statusCancel := context.WithTimeout(ctx, 5*time.Second)
	defer statusCancel()

	emailStatus, err := h.emailMetricsService.GetEmailStatusByInquiry(statusCtx, inquiryID)
	if err != nil {
		h.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to get email status from metrics service")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to retrieve email status at this time",
			"code":    "EMAIL_STATUS_RETRIEVAL_ERROR",
			"details": "Failed to access email event data",
		})
		return
	}

	if emailStatus == nil {
		h.logger.WithField("inquiry_id", inquiryID).Info("No email events found for inquiry")
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error":      "No email events found for this inquiry",
			"code":       "NO_EMAIL_EVENTS",
			"inquiry_id": inquiryID,
			"details":    "This inquiry may not have triggered any email notifications yet",
		})
		return
	}

	// Convert domain.EmailStatus to API response format
	response := h.convertEmailStatusToResponse(emailStatus, inquiry)

	h.logger.WithFields(logrus.Fields{
		"inquiry_id":  inquiryID,
		"email_count": emailStatus.TotalEmailsSent,
		"status":      response.Status,
	}).Debug("Successfully retrieved email status")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"meta": gin.H{
			"total_emails_sent": emailStatus.TotalEmailsSent,
			"last_email_sent":   emailStatus.LastEmailSent,
		},
	})
}

// DownloadReport handles GET /api/v1/admin/reports/:inquiryId/download/:format
func (h *AdminHandler) DownloadReport(c *gin.Context) {
	inquiryID := c.Param("inquiryId")
	format := c.Param("format")
	userID := GetUserIDFromContext(c)
	requestID := GetRequestIDFromContext(c)

	// Create error context for logging
	errorContext := &ErrorContext{
		InquiryID: inquiryID,
		Format:    format,
		UserID:    userID,
		Action:    "download_report",
		RequestID: requestID,
	}

	// Validate inquiry ID - return 404 for empty ID to match routing behavior
	if inquiryID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Not found",
		})
		return
	}

	// Validate format parameter
	if err := h.errorHandler.ValidateDownloadFormat(format); err != nil {
		h.errorHandler.HandleDownloadError(c, err, ErrCodeInvalidFormat, errorContext)
		return
	}

	// Get the inquiry
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), inquiryID)
	if err != nil {
		h.errorHandler.HandleDownloadError(c, err, ErrCodeInternalError, errorContext)
		return
	}

	if inquiry == nil {
		h.errorHandler.HandleDownloadError(c, fmt.Errorf("inquiry with ID %s not found", inquiryID), ErrCodeInquiryNotFound, errorContext)
		return
	}

	// Check if inquiry has reports
	if len(inquiry.Reports) == 0 {
		h.errorHandler.HandleDownloadError(c, fmt.Errorf("no reports available for inquiry %s", inquiryID), ErrCodeNoReports, errorContext)
		return
	}

	// Get the most recent report
	report := inquiry.Reports[0]
	for _, r := range inquiry.Reports {
		if r.CreatedAt.After(report.CreatedAt) {
			report = r
		}
	}

	// Log download attempt
	h.logger.WithFields(logrus.Fields{
		"inquiry_id":  inquiryID,
		"format":      format,
		"user_id":     userID,
		"request_id":  requestID,
		"report_id":   report.ID,
		"report_type": string(report.Type),
		"action":      "download_report_start",
	}).Info("Starting report download")

	switch format {
	case "pdf":
		// Generate PDF
		pdfBytes, err := h.reportService.GeneratePDF(c.Request.Context(), inquiry, report)
		if err != nil {
			h.errorHandler.HandleDownloadError(c, fmt.Errorf("PDF generation failed: %w", err), ErrCodeReportGeneration, errorContext)
			return
		}

		// Set download headers
		filename := generateAdminPDFFilename(inquiry, report)
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

		// Log successful download
		h.errorHandler.LogSuccessfulDownload(c, errorContext, int64(len(pdfBytes)))

		c.Data(http.StatusOK, "application/pdf", pdfBytes)

	case "html":
		// Generate HTML
		htmlContent, err := h.reportService.GenerateHTML(c.Request.Context(), inquiry, report)
		if err != nil {
			h.errorHandler.HandleDownloadError(c, fmt.Errorf("HTML generation failed: %w", err), ErrCodeReportGeneration, errorContext)
			return
		}

		// Set download headers
		filename := generateAdminHTMLFilename(inquiry, report)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Length", fmt.Sprintf("%d", len(htmlContent)))

		// Log successful download
		h.errorHandler.LogSuccessfulDownload(c, errorContext, int64(len(htmlContent)))

		c.String(http.StatusOK, htmlContent)
	}
}

// generateAdminPDFFilename creates a filename for admin PDF downloads
func generateAdminPDFFilename(inquiry *domain.Inquiry, report *domain.Report) string {
	// Sanitize company name for filename
	companyName := inquiry.Company
	if companyName == "" {
		companyName = "Client"
	}

	// Remove special characters and spaces
	sanitized := ""
	for _, r := range companyName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sanitized += string(r)
		} else if r == ' ' || r == '-' || r == '_' {
			sanitized += "_"
		}
	}

	if sanitized == "" {
		sanitized = "Client"
	}

	timestamp := time.Now().Format("20060102")
	return "admin_" + sanitized + "_" + string(report.Type) + "_" + timestamp + ".pdf"
}

// generateAdminHTMLFilename creates a filename for admin HTML downloads
func generateAdminHTMLFilename(inquiry *domain.Inquiry, report *domain.Report) string {
	// Sanitize company name for filename
	companyName := inquiry.Company
	if companyName == "" {
		companyName = "Client"
	}

	// Remove special characters and spaces
	sanitized := ""
	for _, r := range companyName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sanitized += string(r)
		} else if r == ' ' || r == '-' || r == '_' {
			sanitized += "_"
		}
	}

	if sanitized == "" {
		sanitized = "Client"
	}

	timestamp := time.Now().Format("20060102")
	return "admin_" + sanitized + "_" + string(report.Type) + "_" + timestamp + ".html"
}

// parseTimeRange parses time range parameter into domain.TimeRange
func (h *AdminHandler) parseTimeRange(timeRangeParam string) (*domain.TimeRange, error) {
	now := time.Now()
	var start time.Time

	switch timeRangeParam {
	case "1h":
		start = now.Add(-1 * time.Hour)
	case "24h", "1d":
		start = now.Add(-24 * time.Hour)
	case "7d":
		start = now.Add(-7 * 24 * time.Hour)
	case "30d":
		start = now.Add(-30 * 24 * time.Hour)
	case "90d":
		start = now.Add(-90 * 24 * time.Hour)
	default:
		return nil, fmt.Errorf("unsupported time range: %s (supported: 1h, 1d, 7d, 30d, 90d)", timeRangeParam)
	}

	return &domain.TimeRange{
		Start: start,
		End:   now,
	}, nil
}

// convertEmailStatusToResponse converts domain.EmailStatus to API response format
func (h *AdminHandler) convertEmailStatusToResponse(emailStatus *domain.EmailStatus, inquiry *domain.Inquiry) EmailStatus {
	response := EmailStatus{
		InquiryID:       emailStatus.InquiryID,
		CustomerEmail:   inquiry.Email,
		ConsultantEmail: "info@cloudpartner.pro",
		Status:          "no_emails",
	}

	// Determine overall status based on email events
	if emailStatus.CustomerEmail != nil || emailStatus.ConsultantEmail != nil || emailStatus.InquiryNotification != nil {
		// Find the most recent email event to determine overall status
		var mostRecentEvent *domain.EmailEvent
		var mostRecentTime time.Time

		for _, event := range []*domain.EmailEvent{emailStatus.CustomerEmail, emailStatus.ConsultantEmail, emailStatus.InquiryNotification} {
			if event != nil && event.SentAt.After(mostRecentTime) {
				mostRecentEvent = event
				mostRecentTime = event.SentAt
			}
		}

		if mostRecentEvent != nil {
			response.Status = string(mostRecentEvent.Status)
			response.SentAt = mostRecentEvent.SentAt
			if mostRecentEvent.DeliveredAt != nil {
				response.DeliveredAt = *mostRecentEvent.DeliveredAt
			}
			if mostRecentEvent.ErrorMessage != "" {
				response.ErrorMessage = mostRecentEvent.ErrorMessage
			}
		}
	}

	return response
}

// GetEmailEventHistory handles GET /api/v1/admin/email-events
func (h *AdminHandler) GetEmailEventHistory(c *gin.Context) {
	// Parse query parameters with validation
	var filters domain.EmailEventFilters

	// Parse time range with validation
	timeRangeParam := c.DefaultQuery("time_range", "7d")
	timeRange, err := h.parseTimeRange(timeRangeParam)
	if err != nil {
		h.logger.WithError(err).WithField("time_range", timeRangeParam).Error("Invalid time range parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"error":        "Invalid time range parameter: " + err.Error(),
			"code":         "INVALID_TIME_RANGE",
			"valid_ranges": []string{"1h", "1d", "7d", "30d", "90d"},
		})
		return
	}
	filters.TimeRange = timeRange

	// Parse and validate email type filter
	if emailTypeParam := c.Query("email_type"); emailTypeParam != "" {
		emailType := domain.EmailEventType(emailTypeParam)
		// Validate email type
		validTypes := []domain.EmailEventType{
			domain.EmailTypeCustomerConfirmation,
			domain.EmailTypeConsultantNotification,
			domain.EmailTypeInquiryNotification,
		}
		isValid := false
		for _, validType := range validTypes {
			if emailType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			h.logger.WithField("email_type", emailTypeParam).Error("Invalid email type parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid email type parameter",
				"code":    "INVALID_EMAIL_TYPE",
				"valid_types": []string{
					string(domain.EmailTypeCustomerConfirmation),
					string(domain.EmailTypeConsultantNotification),
					string(domain.EmailTypeInquiryNotification),
				},
			})
			return
		}
		filters.EmailType = &emailType
	}

	// Parse and validate status filter
	if statusParam := c.Query("status"); statusParam != "" {
		status := domain.EmailEventStatus(statusParam)
		// Validate status
		validStatuses := []domain.EmailEventStatus{
			domain.EmailStatusSent,
			domain.EmailStatusDelivered,
			domain.EmailStatusFailed,
			domain.EmailStatusBounced,
			domain.EmailStatusSpam,
		}
		isValid := false
		for _, validStatus := range validStatuses {
			if status == validStatus {
				isValid = true
				break
			}
		}
		if !isValid {
			h.logger.WithField("status", statusParam).Error("Invalid status parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid status parameter",
				"code":    "INVALID_STATUS",
				"valid_statuses": []string{
					string(domain.EmailStatusSent),
					string(domain.EmailStatusDelivered),
					string(domain.EmailStatusFailed),
					string(domain.EmailStatusBounced),
					string(domain.EmailStatusSpam),
				},
			})
			return
		}
		filters.Status = &status
	}

	// Parse inquiry ID filter with validation
	if inquiryIDParam := c.Query("inquiry_id"); inquiryIDParam != "" {
		if len(inquiryIDParam) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Inquiry ID cannot be empty",
				"code":    "EMPTY_INQUIRY_ID",
			})
			return
		}
		filters.InquiryID = &inquiryIDParam
	}

	// Parse and validate pagination parameters
	limit := 50 // Default limit
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil || parsedLimit != 1 {
			h.logger.WithField("limit", limitParam).Error("Invalid limit parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid limit parameter - must be a number",
				"code":    "INVALID_LIMIT",
			})
			return
		}
		if limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Limit must be greater than 0",
				"code":    "INVALID_LIMIT_RANGE",
			})
			return
		}
		if limit > 1000 {
			limit = 1000 // Cap maximum limit
		}
	}
	filters.Limit = limit

	offset := 0
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if parsedOffset, err := fmt.Sscanf(offsetParam, "%d", &offset); err != nil || parsedOffset != 1 {
			h.logger.WithField("offset", offsetParam).Error("Invalid offset parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid offset parameter - must be a number",
				"code":    "INVALID_OFFSET",
			})
			return
		}
		if offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Offset cannot be negative",
				"code":    "INVALID_OFFSET_RANGE",
			})
			return
		}
	}
	filters.Offset = offset

	// Check if email metrics service is available and healthy
	if h.emailMetricsService == nil {
		h.logger.Error("Email metrics service not configured for event history request")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email monitoring is not configured",
			"code":    "EMAIL_MONITORING_UNAVAILABLE",
			"details": "Email event history is not available",
		})
		return
	}

	// Check service health with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	healthCtx, healthCancel := context.WithTimeout(ctx, 3*time.Second)
	defer healthCancel()

	if !h.emailMetricsService.IsHealthy(healthCtx) {
		h.logger.Error("Email metrics service is unhealthy for event history request")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "Email monitoring system is currently unavailable",
			"code":    "EMAIL_MONITORING_UNHEALTHY",
			"details": "Email event history service is experiencing issues",
		})
		return
	}

	// Get email event history with timeout
	events, err := h.emailMetricsService.GetEmailEventHistory(ctx, filters)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"time_range": timeRangeParam,
			"email_type": filters.EmailType,
			"status":     filters.Status,
			"inquiry_id": filters.InquiryID,
			"limit":      limit,
			"offset":     offset,
		}).Error("Failed to get email event history")

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Unable to retrieve email event history at this time",
			"code":    "EMAIL_HISTORY_RETRIEVAL_ERROR",
			"details": "Failed to access email event data",
		})
		return
	}

	// Calculate pagination info
	total := len(events)
	page := (offset / limit) + 1
	pages := total / limit
	if total%limit > 0 {
		pages++
	}

	h.logger.WithFields(logrus.Fields{
		"event_count": len(events),
		"total":       total,
		"page":        page,
		"pages":       pages,
		"time_range":  timeRangeParam,
	}).Debug("Successfully retrieved email event history")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    events,
		"count":   len(events),
		"total":   total,
		"page":    page,
		"pages":   pages,
		"filters": map[string]interface{}{
			"time_range": timeRangeParam,
			"email_type": filters.EmailType,
			"status":     filters.Status,
			"inquiry_id": filters.InquiryID,
			"limit":      limit,
			"offset":     offset,
		},
		"meta": gin.H{
			"query_duration_ms": time.Since(time.Now()).Milliseconds(),
			"service_healthy":   true,
		},
	})
}
