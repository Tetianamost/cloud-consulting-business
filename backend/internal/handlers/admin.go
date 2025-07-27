package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/storage"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	storage        *storage.InMemoryStorage
	inquiryService interfaces.InquiryService
	reportService  interfaces.ReportService
	emailService   interfaces.EmailService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	storage *storage.InMemoryStorage,
	inquiryService interfaces.InquiryService,
	reportService interfaces.ReportService,
	emailService interfaces.EmailService,
) *AdminHandler {
	return &AdminHandler{
		storage:       storage,
		inquiryService: inquiryService,
		reportService:  reportService,
		emailService:   emailService,
	}
}

// AdminInquiryFilters represents filters for admin inquiry listing
type AdminInquiryFilters struct {
	Status     string     `form:"status"`
	Priority   string     `form:"priority"`
	Service    string     `form:"service"`
	DateFrom   string     `form:"date_from"`
	DateTo     string     `form:"date_to"`
	Limit      int        `form:"limit"`
	Offset     int        `form:"offset"`
}

// AdminInquiriesResponse represents the response for admin inquiries listing
type AdminInquiriesResponse struct {
	Success bool                `json:"success"`
	Data    []*domain.Inquiry   `json:"data"`
	Count   int                 `json:"count"`
	Total   int64               `json:"total"`
	Page    int                 `json:"page"`
	Pages   int                 `json:"pages"`
}

// SystemMetrics represents system metrics for the admin dashboard
type SystemMetrics struct {
	TotalInquiries      int64   `json:"total_inquiries"`
	ReportsGenerated    int64   `json:"reports_generated"`
	EmailsSent          int64   `json:"emails_sent"`
	EmailDeliveryRate   float64 `json:"email_delivery_rate"`
	AvgReportGenTime    float64 `json:"avg_report_gen_time_ms"`
	SystemUptime        string  `json:"system_uptime"`
	LastProcessedAt     string  `json:"last_processed_at,omitempty"`
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
	Status     string `form:"status"`
	Type       string `form:"type"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
}

// AdminReportsResponse represents the response for admin reports listing
type AdminReportsResponse struct {
	Success bool              `json:"success"`
	Data    []*ReportWithInquiry `json:"data"`
	Count   int               `json:"count"`
	Total   int64             `json:"total"`
	Page    int               `json:"page"`
	Pages   int               `json:"pages"`
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
	// Get total inquiries
	totalInquiries, err := h.inquiryService.GetInquiryCount(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiry count",
		})
		return
	}

	// For demo purposes, we'll calculate some metrics based on the inquiries
	// In a real system, these would be tracked and stored
	inquiries, _ := h.inquiryService.ListInquiries(c.Request.Context(), nil)
	
	var reportsGenerated int64 = 0
	var emailsSent int64 = 0
	var lastProcessed time.Time
	
	for _, inquiry := range inquiries {
		if len(inquiry.Reports) > 0 {
			reportsGenerated += int64(len(inquiry.Reports))
			emailsSent += 2 // Assume 1 customer email + 1 consultant email per report
		}
		
		if inquiry.UpdatedAt.After(lastProcessed) {
			lastProcessed = inquiry.UpdatedAt
		}
	}
	
	// Calculate email delivery rate (for demo, assume 95% success)
	emailDeliveryRate := 95.0
	if emailsSent == 0 {
		emailDeliveryRate = 100.0
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetEmailStatus handles GET /api/v1/admin/email-status/:inquiryId
func (h *AdminHandler) GetEmailStatus(c *gin.Context) {
	inquiryID := c.Param("inquiryId")
	
	// Get the inquiry
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), inquiryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiry",
		})
		return
	}

	if inquiry == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Inquiry not found",
		})
		return
	}

	// For demo purposes, we'll create a mock email status
	// In a real system, this would be tracked and stored
	status := EmailStatus{
		InquiryID:       inquiry.ID,
		CustomerEmail:   inquiry.Email,
		ConsultantEmail: "info@cloudpartner.pro",
		Status:          "delivered",
		SentAt:          inquiry.CreatedAt.Add(time.Minute * 1),
		DeliveredAt:     inquiry.CreatedAt.Add(time.Minute * 1).Add(time.Second * 3),
	}
	
	// If the inquiry was created less than 1 minute ago, show as "sending"
	if time.Since(inquiry.CreatedAt) < time.Minute {
		status.Status = "sending"
		status.DeliveredAt = time.Time{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// DownloadReport handles GET /api/v1/admin/reports/:inquiryId/download/:format
func (h *AdminHandler) DownloadReport(c *gin.Context) {
	inquiryID := c.Param("inquiryId")
	format := c.Param("format")
	
	// Validate format
	if format != "pdf" && format != "html" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid format. Supported formats: pdf, html",
		})
		return
	}
	
	// Get the inquiry
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), inquiryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiry",
		})
		return
	}

	if inquiry == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Inquiry not found",
		})
		return
	}

	// Check if inquiry has reports
	if len(inquiry.Reports) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "No reports found for this inquiry",
		})
		return
	}

	// Get the first (latest) report
	report := inquiry.Reports[0]
	
	// Find the most recent report if there are multiple
	for _, r := range inquiry.Reports {
		if r.CreatedAt.After(report.CreatedAt) {
			report = r
		}
	}

	switch format {
	case "pdf":
		// Generate PDF
		pdfBytes, err := h.reportService.GeneratePDF(c.Request.Context(), inquiry, report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to generate PDF report: " + err.Error(),
			})
			return
		}

		// Set download headers
		filename := generateAdminPDFFilename(inquiry, report)
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
		
		c.Data(http.StatusOK, "application/pdf", pdfBytes)

	case "html":
		// Generate HTML
		htmlContent, err := h.reportService.GenerateHTML(c.Request.Context(), inquiry, report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to generate HTML report",
			})
			return
		}

		// Set download headers
		filename := generateAdminHTMLFilename(inquiry, report)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Length", fmt.Sprintf("%d", len(htmlContent)))
		
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