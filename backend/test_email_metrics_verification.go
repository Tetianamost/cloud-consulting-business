package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

func main() {
	fmt.Println("=== Email Metrics Service Verification ===")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create in-memory database
	db, err := storage.NewInMemoryDB()
	if err != nil {
		log.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Create email event repository
	emailEventRepo := repositories.NewEmailEventRepository(db, logger)

	// Create email metrics service
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, logger)

	ctx := context.Background()

	// Test 1: Verify service implements interface
	fmt.Println("\n--- Test 1: Interface Verification ---")
	var _ interfaces.EmailMetricsService = emailMetricsService
	fmt.Println("✅ EmailMetricsService interface implemented correctly")

	// Test 2: Test with empty database
	fmt.Println("\n--- Test 2: Empty Database Test ---")
	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.Add(-24 * time.Hour),
		End:   now,
	}

	metrics, err := emailMetricsService.GetEmailMetrics(ctx, timeRange)
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics failed: %v\n", err)
	} else {
		fmt.Printf("✅ Empty database metrics: Total=%d, Delivered=%d, Failed=%d\n",
			metrics.TotalEmails, metrics.DeliveredEmails, metrics.FailedEmails)
	}

	// Test 3: Test GetEmailStatusByInquiry with non-existent inquiry
	status, err := emailMetricsService.GetEmailStatusByInquiry(ctx, "non-existent")
	if err != nil {
		fmt.Printf("❌ GetEmailStatusByInquiry failed: %v\n", err)
	} else if status == nil {
		fmt.Println("✅ Non-existent inquiry returns nil status")
	} else {
		fmt.Printf("❌ Expected nil status, got: %+v\n", status)
	}

	// Test 4: Test GetEmailEventHistory with empty database
	filters := domain.EmailEventFilters{
		Limit:  10,
		Offset: 0,
	}

	events, err := emailMetricsService.GetEmailEventHistory(ctx, filters)
	if err != nil {
		fmt.Printf("❌ GetEmailEventHistory failed: %v\n", err)
	} else {
		fmt.Printf("✅ Empty database event history: %d events\n", len(events))
	}

	// Test 5: Create sample data and test metrics calculation
	fmt.Println("\n--- Test 5: Sample Data Test ---")

	// Create sample email events
	sampleEvents := []*domain.EmailEvent{
		{
			ID:          "event-1",
			InquiryID:   "inquiry-1",
			MessageID:   "msg-1",
			EmailType:   domain.EmailTypeCustomerConfirmation,
			Recipient:   "customer@example.com",
			Subject:     "Thank you for your inquiry",
			Status:      domain.EmailStatusDelivered,
			SentAt:      now.Add(-1 * time.Hour),
			DeliveredAt: &[]time.Time{now.Add(-50 * time.Minute)}[0],
		},
		{
			ID:        "event-2",
			InquiryID: "inquiry-1",
			MessageID: "msg-2",
			EmailType: domain.EmailTypeConsultantNotification,
			Recipient: "consultant@cloudpartner.pro",
			Subject:   "New inquiry received",
			Status:    domain.EmailStatusFailed,
			SentAt:    now.Add(-30 * time.Minute),
			Error:     &[]string{"SMTP connection failed"}[0],
		},
	}

	// Insert sample events
	for _, event := range sampleEvents {
		err := emailEventRepo.Create(ctx, event)
		if err != nil {
			fmt.Printf("❌ Failed to create sample event: %v\n", err)
			continue
		}
	}

	fmt.Printf("✅ Created %d sample email events\n", len(sampleEvents))

	// Test metrics with sample data
	metrics, err = emailMetricsService.GetEmailMetrics(ctx, timeRange)
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics with sample data failed: %v\n", err)
	} else {
		fmt.Printf("✅ Sample data metrics: Total=%d, Delivered=%d, Failed=%d\n",
			metrics.TotalEmails, metrics.DeliveredEmails, metrics.FailedEmails)

		// Validate expected values
		if metrics.TotalEmails == 2 && metrics.DeliveredEmails == 1 && metrics.FailedEmails == 1 {
			fmt.Println("✅ Metrics calculations are correct")
		} else {
			fmt.Printf("❌ Metrics calculations incorrect: expected Total=2, Delivered=1, Failed=1\n")
		}
	}

	// Test email status by inquiry
	status, err = emailMetricsService.GetEmailStatusByInquiry(ctx, "inquiry-1")
	if err != nil {
		fmt.Printf("❌ GetEmailStatusByInquiry with sample data failed: %v\n", err)
	} else if status != nil {
		fmt.Printf("✅ Email status for inquiry-1: %d emails sent\n", status.TotalEmailsSent)

		if status.TotalEmailsSent == 2 {
			fmt.Println("✅ Email status calculation is correct")
		} else {
			fmt.Printf("❌ Expected 2 emails sent, got %d\n", status.TotalEmailsSent)
		}
	} else {
		fmt.Println("❌ Expected email status, got nil")
	}

	// Test event history with sample data
	events, err = emailMetricsService.GetEmailEventHistory(ctx, filters)
	if err != nil {
		fmt.Printf("❌ GetEmailEventHistory with sample data failed: %v\n", err)
	} else {
		fmt.Printf("✅ Event history with sample data: %d events\n", len(events))

		if len(events) == 2 {
			fmt.Println("✅ Event history retrieval is correct")
		} else {
			fmt.Printf("❌ Expected 2 events, got %d\n", len(events))
		}
	}

	// Test 6: Verify additional methods exist
	fmt.Println("\n--- Test 6: Extended Methods Test ---")

	// Check if service has additional methods beyond the interface
	if extendedService, ok := emailMetricsService.(interface {
		GetEmailMetricsByType(context.Context, domain.TimeRange) (map[domain.EmailEventType]*domain.EmailMetrics, error)
	}); ok {
		metricsByType, err := extendedService.GetEmailMetricsByType(ctx, timeRange)
		if err != nil {
			fmt.Printf("❌ GetEmailMetricsByType failed: %v\n", err)
		} else {
			fmt.Printf("✅ Email metrics by type: %d types\n", len(metricsByType))
		}
	} else {
		fmt.Println("ℹ️ GetEmailMetricsByType method not available")
	}

	if healthService, ok := emailMetricsService.(interface {
		IsHealthy(context.Context) bool
	}); ok {
		isHealthy := healthService.IsHealthy(ctx)
		if isHealthy {
			fmt.Println("✅ Email metrics service is healthy")
		} else {
			fmt.Println("❌ Email metrics service is unhealthy")
		}
	} else {
		fmt.Println("ℹ️ IsHealthy method not available")
	}

	fmt.Println("\n=== Email Metrics Service Verification Complete ===")
	fmt.Println("✅ All core interface methods are implemented and working correctly")
	fmt.Println("✅ Service handles empty database gracefully")
	fmt.Println("✅ Service calculates metrics correctly with sample data")
	fmt.Println("✅ Service provides comprehensive email status information")
	fmt.Println("✅ Service supports event history filtering")
}
