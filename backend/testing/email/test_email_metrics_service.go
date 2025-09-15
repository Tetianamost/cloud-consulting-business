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
	fmt.Println("=== Email Metrics Service Test ===")

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

	// Test 1: Test with empty database
	fmt.Println("\n--- Test 1: Empty Database ---")
	testEmptyDatabase(ctx, emailMetricsService)

	// Test 2: Create sample email events
	fmt.Println("\n--- Test 2: Creating Sample Email Events ---")
	createSampleEmailEvents(ctx, emailEventRepo)

	// Test 3: Test GetEmailMetrics
	fmt.Println("\n--- Test 3: GetEmailMetrics ---")
	testGetEmailMetrics(ctx, emailMetricsService)

	// Test 4: Test GetEmailStatusByInquiry
	fmt.Println("\n--- Test 4: GetEmailStatusByInquiry ---")
	testGetEmailStatusByInquiry(ctx, emailMetricsService)

	// Test 5: Test GetEmailEventHistory
	fmt.Println("\n--- Test 5: GetEmailEventHistory ---")
	testGetEmailEventHistory(ctx, emailMetricsService)

	// Test 6: Test time range filtering
	fmt.Println("\n--- Test 6: Time Range Filtering ---")
	testTimeRangeFiltering(ctx, emailMetricsService)

	// Test 7: Test service health
	fmt.Println("\n--- Test 7: Service Health Check ---")
	testServiceHealth(ctx, emailMetricsService)

	fmt.Println("\n=== All Email Metrics Service Tests Completed ===")
}

func testEmptyDatabase(ctx context.Context, service interfaces.EmailMetricsService) {
	// Test metrics with empty database
	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.Add(-24 * time.Hour),
		End:   now,
	}

	metrics, err := service.GetEmailMetrics(ctx, timeRange)
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Empty database metrics: Total=%d, Delivered=%d, Failed=%d\n",
		metrics.TotalEmails, metrics.DeliveredEmails, metrics.FailedEmails)

	// Test email status for non-existent inquiry
	status, err := service.GetEmailStatusByInquiry(ctx, "non-existent")
	if err != nil {
		fmt.Printf("❌ GetEmailStatusByInquiry failed: %v\n", err)
		return
	}

	if status == nil {
		fmt.Println("✅ Non-existent inquiry returns nil status")
	} else {
		fmt.Printf("❌ Expected nil status, got: %+v\n", status)
	}
}

func createSampleEmailEvents(ctx context.Context, repo interfaces.EmailEventRepository) {
	now := time.Now()

	// Create sample email events for different inquiries and types
	events := []*domain.EmailEvent{
		{
			ID:          "event-1",
			InquiryID:   "inquiry-1",
			MessageID:   "msg-1",
			EmailType:   domain.EmailTypeCustomerConfirmation,
			Recipient:   "customer1@example.com",
			Subject:     "Thank you for your inquiry",
			Status:      domain.EmailStatusDelivered,
			SentAt:      now.Add(-2 * time.Hour),
			DeliveredAt: &[]time.Time{now.Add(-2*time.Hour + 30*time.Minute)}[0],
		},
		{
			ID:          "event-2",
			InquiryID:   "inquiry-1",
			MessageID:   "msg-2",
			EmailType:   domain.EmailTypeConsultantNotification,
			Recipient:   "consultant@cloudpartner.pro",
			Subject:     "New inquiry received",
			Status:      domain.EmailStatusDelivered,
			SentAt:      now.Add(-2 * time.Hour),
			DeliveredAt: &[]time.Time{now.Add(-2*time.Hour + 15*time.Minute)}[0],
		},
		{
			ID:        "event-3",
			InquiryID: "inquiry-2",
			MessageID: "msg-3",
			EmailType: domain.EmailTypeCustomerConfirmation,
			Recipient: "customer2@example.com",
			Subject:   "Thank you for your inquiry",
			Status:    domain.EmailStatusFailed,
			SentAt:    now.Add(-1 * time.Hour),
			Error:     &[]string{"SMTP connection failed"}[0],
		},
		{
			ID:        "event-4",
			InquiryID: "inquiry-3",
			MessageID: "msg-4",
			EmailType: domain.EmailTypeInquiryNotification,
			Recipient: "admin@cloudpartner.pro",
			Subject:   "High priority inquiry",
			Status:    domain.EmailStatusBounced,
			SentAt:    now.Add(-30 * time.Minute),
			BouncedAt: &[]time.Time{now.Add(-25 * time.Minute)}[0],
		},
		{
			ID:        "event-5",
			InquiryID: "inquiry-4",
			MessageID: "msg-5",
			EmailType: domain.EmailTypeConsultantNotification,
			Recipient: "consultant@cloudpartner.pro",
			Subject:   "New inquiry received",
			Status:    domain.EmailStatusSpam,
			SentAt:    now.Add(-10 * time.Minute),
		},
	}

	for _, event := range events {
		err := repo.Create(ctx, event)
		if err != nil {
			fmt.Printf("❌ Failed to create email event %s: %v\n", event.ID, err)
			return
		}
	}

	fmt.Printf("✅ Created %d sample email events\n", len(events))
}

func testGetEmailMetrics(ctx context.Context, service interfaces.EmailMetricsService) {
	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.Add(-3 * time.Hour),
		End:   now,
	}

	metrics, err := service.GetEmailMetrics(ctx, timeRange)
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Email Metrics:\n")
	fmt.Printf("   Total Emails: %d\n", metrics.TotalEmails)
	fmt.Printf("   Delivered: %d\n", metrics.DeliveredEmails)
	fmt.Printf("   Failed: %d\n", metrics.FailedEmails)
	fmt.Printf("   Bounced: %d\n", metrics.BouncedEmails)
	fmt.Printf("   Spam: %d\n", metrics.SpamEmails)
	fmt.Printf("   Delivery Rate: %.2f%%\n", metrics.DeliveryRate)
	fmt.Printf("   Bounce Rate: %.2f%%\n", metrics.BounceRate)
	fmt.Printf("   Spam Rate: %.2f%%\n", metrics.SpamRate)

	// Validate calculations
	expectedTotal := int64(5)
	expectedDelivered := int64(2)
	expectedFailed := int64(1)
	expectedBounced := int64(1)
	expectedSpam := int64(1)

	if metrics.TotalEmails != expectedTotal {
		fmt.Printf("❌ Expected total emails %d, got %d\n", expectedTotal, metrics.TotalEmails)
	}
	if metrics.DeliveredEmails != expectedDelivered {
		fmt.Printf("❌ Expected delivered emails %d, got %d\n", expectedDelivered, metrics.DeliveredEmails)
	}
	if metrics.FailedEmails != expectedFailed {
		fmt.Printf("❌ Expected failed emails %d, got %d\n", expectedFailed, metrics.FailedEmails)
	}
	if metrics.BouncedEmails != expectedBounced {
		fmt.Printf("❌ Expected bounced emails %d, got %d\n", expectedBounced, metrics.BouncedEmails)
	}
	if metrics.SpamEmails != expectedSpam {
		fmt.Printf("❌ Expected spam emails %d, got %d\n", expectedSpam, metrics.SpamEmails)
	}

	// Check delivery rate calculation (delivered / total * 100)
	expectedDeliveryRate := float64(expectedDelivered) / float64(expectedTotal) * 100
	if metrics.DeliveryRate != expectedDeliveryRate {
		fmt.Printf("❌ Expected delivery rate %.2f%%, got %.2f%%\n", expectedDeliveryRate, metrics.DeliveryRate)
	}
}

func testGetEmailStatusByInquiry(ctx context.Context, service interfaces.EmailMetricsService) {
	// Test inquiry with multiple email types
	status, err := service.GetEmailStatusByInquiry(ctx, "inquiry-1")
	if err != nil {
		fmt.Printf("❌ GetEmailStatusByInquiry failed: %v\n", err)
		return
	}

	if status == nil {
		fmt.Println("❌ Expected email status, got nil")
		return
	}

	fmt.Printf("✅ Email Status for inquiry-1:\n")
	fmt.Printf("   Total Emails Sent: %d\n", status.TotalEmailsSent)
	fmt.Printf("   Has Customer Email: %t\n", status.CustomerEmail != nil)
	fmt.Printf("   Has Consultant Email: %t\n", status.ConsultantEmail != nil)
	fmt.Printf("   Has Inquiry Notification: %t\n", status.InquiryNotification != nil)
	if status.LastEmailSent != nil {
		fmt.Printf("   Last Email Sent: %s\n", status.LastEmailSent.Format(time.RFC3339))
	}

	// Validate expected values
	if status.TotalEmailsSent != 2 {
		fmt.Printf("❌ Expected 2 emails sent, got %d\n", status.TotalEmailsSent)
	}
	if status.CustomerEmail == nil {
		fmt.Println("❌ Expected customer email, got nil")
	}
	if status.ConsultantEmail == nil {
		fmt.Println("❌ Expected consultant email, got nil")
	}
	if status.InquiryNotification != nil {
		fmt.Println("❌ Expected no inquiry notification, got one")
	}

	// Test inquiry with failed email
	status2, err := service.GetEmailStatusByInquiry(ctx, "inquiry-2")
	if err != nil {
		fmt.Printf("❌ GetEmailStatusByInquiry for inquiry-2 failed: %v\n", err)
		return
	}

	if status2 == nil {
		fmt.Println("❌ Expected email status for inquiry-2, got nil")
		return
	}

	fmt.Printf("✅ Email Status for inquiry-2 (failed email):\n")
	fmt.Printf("   Total Emails Sent: %d\n", status2.TotalEmailsSent)
	if status2.CustomerEmail != nil {
		fmt.Printf("   Customer Email Status: %s\n", status2.CustomerEmail.Status)
	}
}

func testGetEmailEventHistory(ctx context.Context, service interfaces.EmailMetricsService) {
	// Test with no filters
	filters := domain.EmailEventFilters{
		Limit:  10,
		Offset: 0,
	}

	events, err := service.GetEmailEventHistory(ctx, filters)
	if err != nil {
		fmt.Printf("❌ GetEmailEventHistory failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Retrieved %d email events (no filters)\n", len(events))

	// Test with email type filter
	customerEmailType := domain.EmailTypeCustomerConfirmation
	filters.EmailType = &customerEmailType

	customerEvents, err := service.GetEmailEventHistory(ctx, filters)
	if err != nil {
		fmt.Printf("❌ GetEmailEventHistory with email type filter failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Retrieved %d customer confirmation events\n", len(customerEvents))

	// Validate all events are customer confirmation type
	for _, event := range customerEvents {
		if event.EmailType != domain.EmailTypeCustomerConfirmation {
			fmt.Printf("❌ Expected customer confirmation email, got %s\n", event.EmailType)
		}
	}

	// Test with status filter
	filters.EmailType = nil
	deliveredStatus := domain.EmailStatusDelivered
	filters.Status = &deliveredStatus

	deliveredEvents, err := service.GetEmailEventHistory(ctx, filters)
	if err != nil {
		fmt.Printf("❌ GetEmailEventHistory with status filter failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Retrieved %d delivered events\n", len(deliveredEvents))

	// Validate all events are delivered
	for _, event := range deliveredEvents {
		if event.Status != domain.EmailStatusDelivered {
			fmt.Printf("❌ Expected delivered status, got %s\n", event.Status)
		}
	}
}

func testTimeRangeFiltering(ctx context.Context, service interfaces.EmailMetricsService) {
	now := time.Now()

	// Test recent time range (last hour)
	recentTimeRange := domain.TimeRange{
		Start: now.Add(-1 * time.Hour),
		End:   now,
	}

	recentMetrics, err := service.GetEmailMetrics(ctx, recentTimeRange)
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics for recent time range failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Recent metrics (last hour): %d total emails\n", recentMetrics.TotalEmails)

	// Test older time range (2-3 hours ago)
	olderTimeRange := domain.TimeRange{
		Start: now.Add(-3 * time.Hour),
		End:   now.Add(-2 * time.Hour),
	}

	olderMetrics, err := service.GetEmailMetrics(ctx, olderTimeRange)
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics for older time range failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Older metrics (2-3 hours ago): %d total emails\n", olderMetrics.TotalEmails)

	// Validate that recent + older should be less than or equal to total
	totalMetrics, err := service.GetEmailMetrics(ctx, domain.TimeRange{
		Start: now.Add(-3 * time.Hour),
		End:   now,
	})
	if err != nil {
		fmt.Printf("❌ GetEmailMetrics for total time range failed: %v\n", err)
		return
	}

	if recentMetrics.TotalEmails+olderMetrics.TotalEmails > totalMetrics.TotalEmails {
		fmt.Printf("❌ Time range filtering inconsistent: recent(%d) + older(%d) > total(%d)\n",
			recentMetrics.TotalEmails, olderMetrics.TotalEmails, totalMetrics.TotalEmails)
	} else {
		fmt.Println("✅ Time range filtering is consistent")
	}
}

func testServiceHealth(ctx context.Context, service interfaces.EmailMetricsService) {
	// Test if service implements health check
	if healthyService, ok := service.(interface{ IsHealthy(context.Context) bool }); ok {
		isHealthy := healthyService.IsHealthy(ctx)
		if isHealthy {
			fmt.Println("✅ Email metrics service is healthy")
		} else {
			fmt.Println("❌ Email metrics service is unhealthy")
		}
	} else {
		fmt.Println("ℹ️ Email metrics service does not implement health check")
	}

	// Test additional methods if they exist
	if extendedService, ok := service.(interface {
		GetEmailMetricsByType(context.Context, domain.TimeRange) (map[domain.EmailEventType]*domain.EmailMetrics, error)
	}); ok {
		now := time.Now()
		timeRange := domain.TimeRange{
			Start: now.Add(-3 * time.Hour),
			End:   now,
		}

		metricsByType, err := extendedService.GetEmailMetricsByType(ctx, timeRange)
		if err != nil {
			fmt.Printf("❌ GetEmailMetricsByType failed: %v\n", err)
		} else {
			fmt.Printf("✅ Email metrics by type:\n")
			for emailType, metrics := range metricsByType {
				fmt.Printf("   %s: %d emails\n", emailType, metrics.TotalEmails)
			}
		}
	}

	if recentActivityService, ok := service.(interface {
		GetRecentEmailActivity(context.Context, int) ([]*domain.EmailEvent, error)
	}); ok {
		recentEvents, err := recentActivityService.GetRecentEmailActivity(ctx, 24)
		if err != nil {
			fmt.Printf("❌ GetRecentEmailActivity failed: %v\n", err)
		} else {
			fmt.Printf("✅ Recent email activity (24h): %d events\n", len(recentEvents))
		}
	}
}
