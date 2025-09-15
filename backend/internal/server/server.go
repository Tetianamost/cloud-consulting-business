package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// Server represents the HTTP server
type Server struct {
	config             *config.Config
	logger             *logrus.Logger
	router             *gin.Engine
	inquiryHandler     *handlers.InquiryHandler
	reportHandler      *handlers.ReportHandler
	adminHandler       *handlers.AdminHandler
	authHandler        *handlers.AuthHandler
	chatHandler        *handlers.ChatHandler
	pollingChatHandler *handlers.PollingChatHandler
	chatConfigHandler  *handlers.ChatConfigHandler
	healthHandler      *handlers.HealthHandler
	simpleChatHandler  *handlers.SimpleChatHandler
	bedrockService     interfaces.BedrockService
	dbConnection       *storage.DatabaseConnection
}

// New creates a new server instance
func New(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	fmt.Printf("[SERVER DEBUG] cfg.CORSAllowedOrigins: %v\n", cfg.CORSAllowedOrigins)

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg.CORSAllowedOrigins))
	router.Use(loggingMiddleware(logger))

	// Initialize storage
	memStorage := storage.NewInMemoryStorage()

	// Initialize database connection if configured
	var dbConnection *storage.DatabaseConnection
	var emailEventRecorder interfaces.EmailEventRecorder
	var emailMetricsService interfaces.EmailMetricsService

	if cfg.IsEmailEventTrackingEnabled() {
		var err error
		dbConnection, err = storage.NewDatabaseConnection(&cfg.Database, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize database connection, email event tracking will be disabled")
		} else {
			// Run email events migration
			migrationSQL := getEmailEventsMigrationSQL()
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := dbConnection.RunMigration(ctx, migrationSQL); err != nil {
				logger.WithError(err).Warn("Failed to run email events migration, email event tracking will be disabled")
				dbConnection.Close()
				dbConnection = nil
			} else {
				// Initialize email event repository and services
				emailEventRepo := repositories.NewEmailEventRepository(dbConnection.GetDB(), logger)
				emailEventRecorder = services.NewEmailEventRecorder(emailEventRepo, logger)
				emailMetricsService = services.NewEmailMetricsService(emailEventRepo, logger)
				logger.Info("Email event tracking initialized successfully")
			}
		}
	} else {
		logger.Info("Database not configured or email events disabled, using in-memory storage only")
	}

	// Initialize template service first
	logger.Info("Initializing template service with path: /usr/local/bin/templates")
	templateService := services.NewTemplateService("/usr/local/bin/templates", logger)
	logger.Info("Template service initialized successfully")

	// Initialize PDF service
	pdfService := services.NewPDFService(logger)

	// Initialize services
	bedrockService := services.NewBedrockService(&cfg.Bedrock)
	promptArchitect := services.NewPromptArchitect()
	knowledgeBase := services.NewInMemoryKnowledgeBase()
	documentationLibrary := services.NewDocumentationLibraryService()
	riskAssessor := services.NewRiskAssessorService(knowledgeBase, documentationLibrary)
	multiCloudAnalyzer := services.NewMultiCloudAnalyzerService(knowledgeBase, documentationLibrary)
	reportGenerator := services.NewReportGenerator(
		bedrockService,
		templateService,
		pdfService,
		promptArchitect,
		knowledgeBase,
		multiCloudAnalyzer,
		riskAssessor,
		documentationLibrary,
	)

	// Initialize email services using the primary factory method
	var emailService interfaces.EmailService
	if cfg.SES.AccessKeyID != "" && cfg.SES.SecretAccessKey != "" && cfg.SES.SenderEmail != "" {
		var err error
		// Use the primary factory method which automatically handles event recording
		emailService, err = services.NewEmailServiceFactory(cfg.SES, emailEventRecorder, logger)
		if err != nil {
			logger.WithError(err).Warn("Failed to initialize email service, email notifications will be disabled")
			emailService = nil
		} else {
			eventRecordingStatus := "disabled"
			if emailEventRecorder != nil {
				eventRecordingStatus = "enabled"
			}
			logger.WithField("event_recording", eventRecordingStatus).Info("Email service initialized successfully with branded templates")
		}
	} else {
		logger.Warn("SES configuration incomplete, email notifications will be disabled")
	}

	inquiryService := services.NewInquiryService(memStorage, reportGenerator, emailService)

	// Initialize chat repositories and services
	chatSessionRepository := storage.NewInMemoryChatSessionRepository(logger)
	chatMessageRepository := storage.NewInMemoryChatMessageRepository(logger)
	sessionService := services.NewSessionService(chatSessionRepository, logger)
	chatService := services.NewChatService(chatMessageRepository, sessionService, bedrockService, logger)

	// Initialize chat monitoring services
	chatPerformanceMonitor := services.NewChatPerformanceMonitor(logger)
	// Create a simple cache monitor for in-memory usage (pass nil for Redis cache)
	cacheMonitor := services.NewCacheMonitor(nil, logger)
	chatMetricsCollector := services.NewChatMetricsCollector(chatPerformanceMonitor, cacheMonitor, logger)

	// Email metrics service is already initialized above if database is available

	// Initialize handlers
	inquiryHandler := handlers.NewInquiryHandler(inquiryService, reportGenerator)
	reportHandler := handlers.NewReportHandler(memStorage)
	adminHandler := handlers.NewAdminHandler(memStorage, inquiryService, reportGenerator, emailService, emailMetricsService, logger)
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)
	chatHandler := handlers.NewChatHandler(logger, bedrockService, knowledgeBase, sessionService, chatService, authHandler, cfg.JWTSecret, cfg.CORSAllowedOrigins, chatMetricsCollector, chatPerformanceMonitor)

	// Initialize security services for polling chat
	chatAuthService := services.NewChatAuthService(cfg.JWTSecret, logger)
	chatRateLimiter := services.NewChatRateLimiter(logger)
	chatAuditLogger := services.NewChatAuditLogger(logger)
	chatSecurityService := services.NewChatSecurityService(logger, cfg.JWTSecret, chatRateLimiter, chatAuditLogger)

	// Initialize polling chat handler
	pollingChatHandler := handlers.NewPollingChatHandler(logger, chatService, sessionService, authHandler, chatAuthService, chatSecurityService, chatRateLimiter, chatAuditLogger)

	// Initialize chat configuration handler
	chatConfigHandler := handlers.NewChatConfigHandler(cfg)

	// Initialize health handler with database connection and email event services
	var db *sql.DB
	if dbConnection != nil {
		db = dbConnection.GetDB()
	}
	healthHandler := handlers.NewHealthHandler(logger, db, chatHandler.GetConnectionPool(), emailEventRecorder, emailMetricsService)

	// Initialize simple chat handler
	simpleChatHandler := handlers.NewSimpleChatHandler(logger, bedrockService)

	server := &Server{
		config:             cfg,
		logger:             logger,
		router:             router,
		inquiryHandler:     inquiryHandler,
		reportHandler:      reportHandler,
		adminHandler:       adminHandler,
		authHandler:        authHandler,
		chatHandler:        chatHandler,
		pollingChatHandler: pollingChatHandler,
		chatConfigHandler:  chatConfigHandler,
		healthHandler:      healthHandler,
		simpleChatHandler:  simpleChatHandler,
		bedrockService:     bedrockService,
		dbConnection:       dbConnection,
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
	// Serve static files (CSS, etc.)
	s.router.Static("/static", "./static")

	// Health check endpoints
	s.router.GET("/health", s.healthCheck)

	// Enhanced health check endpoints (if health handler is available)
	if s.healthHandler != nil {
		health := s.router.Group("/health")
		{
			health.GET("/detailed", s.healthHandler.DetailedHealthCheck)
			health.GET("/ready", s.healthHandler.ReadinessCheck)
			health.GET("/live", s.healthHandler.LivenessCheck)
			health.GET("/websocket", s.healthHandler.WebSocketHealthCheck)
			health.GET("/websocket/config", s.healthHandler.ConfigurationValidationCheck)
		}
	}

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
			inquiries.GET("/:id/report/html", s.inquiryHandler.GetInquiryReportHTML)
			inquiries.GET("/:id/report/pdf", s.inquiryHandler.GetInquiryReportPDF)
			inquiries.GET("/:id/report/download", s.inquiryHandler.DownloadInquiryReport)
		}

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", s.authHandler.Login)
		}

		// Admin routes - protected by auth middleware
		admin := v1.Group("/admin", s.authHandler.AuthMiddleware())
		{
			admin.GET("/inquiries", s.adminHandler.ListInquiries)
			admin.GET("/reports", s.adminHandler.ListReports)
			admin.GET("/reports/download/:inquiryId/:format", s.adminHandler.DownloadReport)
			admin.GET("/reports/view/:reportId", s.adminHandler.GetReport)
			admin.GET("/metrics", s.adminHandler.GetSystemMetrics)
			admin.GET("/email-status/:inquiryId", s.adminHandler.GetEmailStatus)
			admin.GET("/email-events", s.adminHandler.GetEmailEventHistory)

			// Chat management routes - organized to avoid conflicts
			chatMgmt := admin.Group("/chat-mgmt")
			{
				chatMgmt.POST("/sessions", s.chatHandler.CreateChatSession)
				chatMgmt.GET("/sessions", s.chatHandler.ListChatSessions)
				chatMgmt.GET("/sessions/:id", s.chatHandler.GetChatSessionByID)
				chatMgmt.PUT("/sessions/:id", s.chatHandler.UpdateChatSession)
				chatMgmt.DELETE("/sessions/:id", s.chatHandler.DeleteChatSession)
				chatMgmt.GET("/sessions/:id/history", s.chatHandler.GetChatSessionHistory)
			}

			// Legacy endpoints for backward compatibility
			admin.GET("/chat/sessions-legacy", s.chatHandler.GetChatSessions)
			admin.GET("/chat/sessions-legacy/:sessionId", s.chatHandler.GetChatSession)

			// Simple working chat endpoints (bypass complex polling)
			simpleChat := admin.Group("/simple-chat")
			{
				simpleChat.POST("/messages", s.simpleChatHandler.SendMessage)
				simpleChat.GET("/messages", s.simpleChatHandler.GetMessages)
			}

			// Polling-based chat endpoints
			chatPolling := admin.Group("/chat-polling")
			{
				chatPolling.POST("/messages", s.pollingChatHandler.SendMessage)
				chatPolling.GET("/messages", s.pollingChatHandler.GetMessages)
			}

			// Chat configuration endpoints
			admin.GET("/chat/config", s.chatConfigHandler.GetChatConfig)
			admin.PUT("/chat/config", s.chatConfigHandler.UpdateChatConfig)
		}

		// WebSocket routes - use WebSocket-specific auth middleware that supports query parameters
		adminWS := v1.Group("/admin")
		{
			// Chat WebSocket route with WebSocket-specific authentication
			adminWS.GET("/chat/ws", s.chatHandler.HandleWebSocket)
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
	Success bool              `json:"success"`
	Data    ServiceConfigData `json:"data"`
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
		domain.ServiceTypeAssessment:         "Cloud Assessment",
		domain.ServiceTypeMigration:          "Cloud Migration",
		domain.ServiceTypeOptimization:       "Cloud Optimization",
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
	fmt.Printf("[CORS DEBUG] Allowed origins: %v\n", allowedOrigins)
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		fmt.Printf("[CORS DEBUG] Incoming request Origin: %s\n", origin)

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			trimmedOrigin := strings.TrimSpace(allowedOrigin)
			trimmedRequestOrigin := strings.TrimSpace(origin)
			if trimmedRequestOrigin == trimmedOrigin || trimmedOrigin == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "null")
		}
		// Log CORS diagnostics for debugging
		c.Header("X-Debug-CORS-Origin", origin)
		c.Header("X-Debug-CORS-Allowed", fmt.Sprintf("%v", allowed))
		fmt.Printf("[CORS DEBUG] Response Access-Control-Allow-Origin: %s\n", c.Writer.Header().Get("Access-Control-Allow-Origin"))

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, Upgrade, Connection, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Protocol")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle WebSocket upgrade requests
		if c.Request.Header.Get("Upgrade") == "websocket" {
			// Allow WebSocket upgrade for allowed origins
			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		fmt.Printf("[CORS DEBUG] End of middleware for %s, Access-Control-Allow-Origin: %s\n", c.Request.URL.Path, c.Writer.Header().Get("Access-Control-Allow-Origin"))
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

// getEmailEventsMigrationSQL returns the SQL for creating email events table and related objects
func getEmailEventsMigrationSQL() string {
	return `
-- Email events tracking database migration
-- This migration adds email_events table for real email monitoring data
-- Replaces mock data with actual email delivery tracking

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for better type safety and performance
DO $$ BEGIN
    CREATE TYPE email_event_type AS ENUM ('customer_confirmation', 'consultant_notification', 'inquiry_notification');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE email_event_status AS ENUM ('sent', 'delivered', 'failed', 'bounced', 'spam');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE bounce_type AS ENUM ('permanent', 'temporary', 'complaint');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Create email_events table
CREATE TABLE IF NOT EXISTS email_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inquiry_id VARCHAR(255) NOT NULL,
    email_type email_event_type NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    sender_email VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    status email_event_status NOT NULL DEFAULT 'sent',
    sent_at TIMESTAMP WITH TIME ZONE NOT NULL,
    delivered_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    bounce_type bounce_type,
    ses_message_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for optimal query performance
CREATE INDEX IF NOT EXISTS idx_email_events_inquiry_id ON email_events(inquiry_id);
CREATE INDEX IF NOT EXISTS idx_email_events_status ON email_events(status);
CREATE INDEX IF NOT EXISTS idx_email_events_sent_at ON email_events(sent_at);
CREATE INDEX IF NOT EXISTS idx_email_events_email_type ON email_events(email_type);

-- Additional composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_email_events_inquiry_type ON email_events(inquiry_id, email_type);
CREATE INDEX IF NOT EXISTS idx_email_events_status_sent ON email_events(status, sent_at DESC);
CREATE INDEX IF NOT EXISTS idx_email_events_type_status ON email_events(email_type, status);
CREATE INDEX IF NOT EXISTS idx_email_events_ses_message_id ON email_events(ses_message_id) WHERE ses_message_id IS NOT NULL;

-- Partial indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_email_events_failed_status ON email_events(sent_at DESC, error_message) WHERE status IN ('failed', 'bounced', 'spam');
CREATE INDEX IF NOT EXISTS idx_email_events_recent_events ON email_events(sent_at DESC, status) WHERE sent_at > NOW() - INTERVAL '30 days';

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_email_events_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_email_events_updated_at ON email_events;
CREATE TRIGGER update_email_events_updated_at 
    BEFORE UPDATE ON email_events 
    FOR EACH ROW EXECUTE FUNCTION update_email_events_updated_at();

-- Add table constraints for data integrity
ALTER TABLE email_events 
DROP CONSTRAINT IF EXISTS chk_email_events_dates;
ALTER TABLE email_events 
ADD CONSTRAINT chk_email_events_dates 
CHECK (sent_at <= COALESCE(delivered_at, sent_at) AND created_at <= updated_at);

ALTER TABLE email_events 
DROP CONSTRAINT IF EXISTS chk_email_events_recipient_email;
ALTER TABLE email_events 
ADD CONSTRAINT chk_email_events_recipient_email 
CHECK (recipient_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE email_events 
DROP CONSTRAINT IF EXISTS chk_email_events_sender_email;
ALTER TABLE email_events 
ADD CONSTRAINT chk_email_events_sender_email 
CHECK (sender_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE email_events 
DROP CONSTRAINT IF EXISTS chk_email_events_subject_length;
ALTER TABLE email_events 
ADD CONSTRAINT chk_email_events_subject_length 
CHECK (LENGTH(TRIM(COALESCE(subject, ''))) <= 500);

ALTER TABLE email_events 
DROP CONSTRAINT IF EXISTS chk_email_events_error_message_length;
ALTER TABLE email_events 
ADD CONSTRAINT chk_email_events_error_message_length 
CHECK (LENGTH(COALESCE(error_message, '')) <= 10000);
`
}
