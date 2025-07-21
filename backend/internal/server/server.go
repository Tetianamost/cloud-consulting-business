package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// Server represents the HTTP server
type Server struct {
	config         *config.Config
	logger         *logrus.Logger
	router         *gin.Engine
	inquiryHandler *handlers.InquiryHandler
	reportHandler  *handlers.ReportHandler
}

// New creates a new server instance
func New(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	router := gin.New()
	
	// Add middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg.CORSAllowedOrigins))
	router.Use(loggingMiddleware(logger))

	// Initialize storage
	memStorage := storage.NewInMemoryStorage()
	
	// Initialize services
	bedrockService := services.NewBedrockService(&cfg.Bedrock)
	reportGenerator := services.NewReportGenerator(bedrockService)
	
	// Initialize template service
	templateService := services.NewTemplateService("templates", logger)
	
	// Initialize email services (with graceful degradation if SES config is missing)
	var emailService interfaces.EmailService
	if cfg.SES.AccessKeyID != "" && cfg.SES.SecretAccessKey != "" && cfg.SES.SenderEmail != "" {
		sesService, err := services.NewSESService(cfg.SES, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize SES service, email notifications will be disabled")
		} else {
			emailService = services.NewEmailService(sesService, templateService, cfg.SES, logger)
			logger.Info("Email service initialized successfully with branded templates")
		}
	} else {
		logger.Warn("SES configuration incomplete, email notifications will be disabled")
	}
	
	inquiryService := services.NewInquiryService(memStorage, reportGenerator, emailService)
	
	// Initialize handlers
	inquiryHandler := handlers.NewInquiryHandler(inquiryService)
	reportHandler := handlers.NewReportHandler(memStorage)

	server := &Server{
		config:         cfg,
		logger:         logger,
		router:         router,
		inquiryHandler: inquiryHandler,
		reportHandler:  reportHandler,
	}

	// Setup routes
	server.setupRoutes()

	return server, nil
}

// Handler returns the HTTP handler
func (s *Server) Handler() http.Handler {
	return s.router
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)
	
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Inquiry routes
		inquiries := v1.Group("/inquiries")
		{
			inquiries.POST("", s.inquiryHandler.CreateInquiry)
			inquiries.GET("/:id", s.inquiryHandler.GetInquiry)
			inquiries.GET("", s.inquiryHandler.ListInquiries)
			inquiries.GET("/:id/report", s.reportHandler.GetInquiryReport)
		}

		// System management routes
		v1.GET("/config/services", s.getServiceConfig)
	}
}

// Health check handler
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "cloud-consulting-backend",
		"version": "1.0.0",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

// ServiceConfigResponse represents the service configuration response
type ServiceConfigResponse struct {
	Success bool                   `json:"success"`
	Data    ServiceConfigData      `json:"data"`
}

// ServiceConfigData represents the service configuration data
type ServiceConfigData struct {
	Services []ServiceInfo `json:"services"`
}

// ServiceInfo represents information about a service type
type ServiceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Service configuration handler
func (s *Server) getServiceConfig(c *gin.Context) {
	// Build service information from domain constants
	services := make([]ServiceInfo, 0, len(domain.ValidServiceTypes))
	
	// Service name mapping for better frontend display
	serviceNames := map[string]string{
		domain.ServiceTypeAssessment:        "Cloud Assessment",
		domain.ServiceTypeMigration:         "Cloud Migration",
		domain.ServiceTypeOptimization:      "Cloud Optimization",
		domain.ServiceTypeArchitectureReview: "Architecture Review",
	}
	
	for _, serviceType := range domain.ValidServiceTypes {
		services = append(services, ServiceInfo{
			ID:          serviceType,
			Name:        serviceNames[serviceType],
			Description: domain.ServiceDescriptions[serviceType],
		})
	}
	
	response := ServiceConfigResponse{
		Success: true,
		Data: ServiceConfigData{
			Services: services,
		},
	}
	
	c.JSON(http.StatusOK, response)
}
// Middleware functions

// corsMiddleware handles CORS headers
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Log request
		latency := time.Since(start)
		
		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP Request")
	}
}