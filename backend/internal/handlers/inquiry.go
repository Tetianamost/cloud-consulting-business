package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/storage"
)

// InquiryHandler handles inquiry-related HTTP requests
type InquiryHandler struct {
	storage *storage.InMemoryStorage
}

// NewInquiryHandler creates a new inquiry handler
func NewInquiryHandler(storage *storage.InMemoryStorage) *InquiryHandler {
	return &InquiryHandler{
		storage: storage,
	}
}

// CreateInquiry handles POST /api/v1/inquiries
func (h *InquiryHandler) CreateInquiry(c *gin.Context) {
	var req domain.CreateInquiryRequest
	
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

	inquiry, err := h.storage.CreateInquiry(&req)
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
	
	inquiry, err := h.storage.GetInquiry(id)
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
	inquiries, err := h.storage.ListInquiries()
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