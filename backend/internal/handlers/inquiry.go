package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// InquiryHandler handles inquiry-related HTTP requests
type InquiryHandler struct {
	inquiryService  interfaces.InquiryService
	reportService   interfaces.ReportService
}

// NewInquiryHandler creates a new inquiry handler
func NewInquiryHandler(inquiryService interfaces.InquiryService, reportService interfaces.ReportService) *InquiryHandler {
	return &InquiryHandler{
		inquiryService: inquiryService,
		reportService:  reportService,
	}
}

// CreateInquiry handles POST /api/v1/inquiries
func (h *InquiryHandler) CreateInquiry(c *gin.Context) {
	var req interfaces.CreateInquiryRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Basic validation for services
	if len(req.Services) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "At least one service must be selected",
		})
		return
	}

	// Validate service types
	for _, service := range req.Services {
		valid := false
		for _, validService := range domain.ValidServiceTypes {
			if service == validService {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid service type: " + service,
			})
			return
		}
	}

	inquiry, err := h.inquiryService.CreateInquiry(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create inquiry",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    inquiry,
		"message": "Inquiry created successfully",
	})
}

// GetInquiry handles GET /api/v1/inquiries/:id
func (h *InquiryHandler) GetInquiry(c *gin.Context) {
	id := c.Param("id")
	
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), id)
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    inquiry,
	})
}

// ListInquiries handles GET /api/v1/inquiries
func (h *InquiryHandler) ListInquiries(c *gin.Context) {
	inquiries, err := h.inquiryService.ListInquiries(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve inquiries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    inquiries,
		"count":   len(inquiries),
	})
}

// GetInquiryReportHTML handles GET /api/v1/inquiries/{id}/report/html
func (h *InquiryHandler) GetInquiryReportHTML(c *gin.Context) {
	id := c.Param("id")
	
	// Get the inquiry
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), id)
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

	// Generate HTML version of the report
	if h.reportService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Report service not available",
		})
		return
	}

	htmlContent, err := h.reportService.GenerateHTML(c.Request.Context(), inquiry, report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate HTML report",
		})
		return
	}

	// Return HTML content
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}

// GetInquiryReportPDF handles GET /api/v1/inquiries/{id}/report/pdf
func (h *InquiryHandler) GetInquiryReportPDF(c *gin.Context) {
	id := c.Param("id")
	
	// Get the inquiry
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), id)
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

	// Generate PDF version of the report
	if h.reportService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Report service not available",
		})
		return
	}

	pdfBytes, err := h.reportService.GeneratePDF(c.Request.Context(), inquiry, report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to generate PDF report: " + err.Error(),
		})
		return
	}

	// Set appropriate headers for PDF
	filename := generatePDFFilename(inquiry, report)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	
	// Return PDF content
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// DownloadInquiryReport handles GET /api/v1/inquiries/{id}/report/download
func (h *InquiryHandler) DownloadInquiryReport(c *gin.Context) {
	id := c.Param("id")
	format := c.DefaultQuery("format", "pdf") // Default to PDF
	
	// Get the inquiry
	inquiry, err := h.inquiryService.GetInquiry(c.Request.Context(), id)
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

	if h.reportService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Report service not available",
		})
		return
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
		filename := generatePDFFilename(inquiry, report)
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
		filename := generateHTMLFilename(inquiry, report)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Length", fmt.Sprintf("%d", len(htmlContent)))
		
		c.String(http.StatusOK, htmlContent)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid format. Supported formats: pdf, html",
		})
	}
}

// generatePDFFilename creates a filename for PDF downloads
func generatePDFFilename(inquiry *domain.Inquiry, report *domain.Report) string {
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
	
	return sanitized + "_" + string(report.Type) + "_report.pdf"
}

// generateHTMLFilename creates a filename for HTML downloads
func generateHTMLFilename(inquiry *domain.Inquiry, report *domain.Report) string {
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
	
	return sanitized + "_" + string(report.Type) + "_report.html"
}