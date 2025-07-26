package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cloud-consulting/backend/internal/storage"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	storage *storage.InMemoryStorage
}

// NewReportHandler creates a new report handler
func NewReportHandler(storage *storage.InMemoryStorage) *ReportHandler {
	return &ReportHandler{
		storage: storage,
	}
}

// GetInquiryReport handles GET /api/v1/inquiries/:id/report
func (h *ReportHandler) GetInquiryReport(c *gin.Context) {
	inquiryID := c.Param("id")
	
	// Get reports for the inquiry
	reports, err := h.storage.GetReportsByInquiry(c.Request.Context(), inquiryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve reports",
		})
		return
	}

	if len(reports) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "No reports found for this inquiry",
		})
		return
	}

	// Return the most recent report (for now, just return the first one)
	report := reports[0]
	if len(reports) > 1 {
		// Find the most recent report
		for _, r := range reports {
			if r.CreatedAt.After(report.CreatedAt) {
				report = r
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
	})
}