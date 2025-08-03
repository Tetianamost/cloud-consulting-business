package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// Performance test configuration
const (
	defaultConcurrency = 50
	defaultDuration    = 30 * time.Second
	defaultRampUp      = 5 * time.Second
)

// Performance test metrics
type PerformanceMetrics struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	AverageLatency     time.Duration
	MinLatency         time.Duration
	MaxLatency         time.Duration
	P95Latency         time.Duration
	P99Latency         time.Duration
	RequestsPerSecond  float64
	ErrorRate          float64
	Latencies          []time.Duration
	mutex              sync.RWMutex
}

func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		MinLatency: time.Hour, // Initialize with high value
		Latencies:  make([]time.Duration, 0),
	}
}

func (pm *PerformanceMetrics) RecordRequest(latency time.Duration, success bool) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalRequests++
	pm.Latencies = append(pm.Latencies, latency)

	if success {
		pm.SuccessfulRequests++
	} else {
		pm.FailedRequests++
	}

	if latency < pm.MinLatency {
		pm.MinLatency = latency
	}
	if latency > pm.MaxLatency {
		pm.MaxLatency = latency
	}
}

func (pm *PerformanceMetrics) Calculate(duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if pm.TotalRequests == 0 {
		return
	}

	// Calculate average latency
	var totalLatency time.Duration
	for _, latency := range pm.Latencies {
		totalLatency += latency
	}
	pm.AverageLatency = totalLatency / time.Duration(pm.TotalRequests)

	// Calculate percentiles
	if len(pm.Latencies) > 0 {
		// Sort latencies for percentile calculation
		latencies := make([]time.Duration, len(pm.Latencies))
		copy(latencies, pm.Latencies)

		// Simple bubble sort for small datasets
		for i := 0; i < len(latencies); i++ {
			for j := 0; j < len(latencies)-1-i; j++ {
				if latencies[j] > latencies[j+1] {
					latencies[j], latencies[j+1] = latencies[j+1], latencies[j]
				}
			}
		}

		p95Index := int(float64(len(latencies)) * 0.95)
		p99Index := int(float64(len(latencies)) * 0.99)

		if p95Index < len(latencies) {
			pm.P95Latency = latencies[p95Index]
		}
		if p99Index < len(latencies) {
			pm.P99Latency = latencies[p99Index]
		}
	}

	// Calculate requests per second
	pm.RequestsPerSecond = float64(pm.TotalRequests) / duration.Seconds()

	// Calculate error rate
	pm.ErrorRate = float64(pm.FailedRequests) / float64(pm.TotalRequests) * 100
}

func (pm *PerformanceMetrics) String() string {
	return fmt.Sprintf(`Performance Metrics:
  Total Requests: %d
  Successful: %d
  Failed: %d
  Requests/Second: %.2f
  Error Rate: %.2f%%
  Average Latency: %v
  Min Latency: %v
  Max Latency: %v
  P95 Latency: %v
  P99 Latency: %v`,
		pm.TotalRequests,
		pm.SuccessfulRequests,
		pm.FailedRequests,
		pm.RequestsPerSecond,
		pm.ErrorRate,
		pm.AverageLatency,
		pm.MinLatency,
		pm.MaxLatency,
		pm.P95Latency,
		pm.P99Latency,
	)
}

// Load test suite
type LoadTestSuite struct {
	server      *httptest.Server
	chatHandler *handlers.ChatHandler
	metrics     *PerformanceMetrics
}

func setupLoadTestSuite(t *testing.T) *LoadTestSuite {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Use in-memory repositories for performance testing
	sessionRepo := storage.NewMemoryChatSessionRepository()
	messageRepo := storage.NewMemoryChatMessageRepository()

	sessionSvc := services.NewSessionService(sessionRepo, logger)
	chatSvc := services.NewChatService(messageRepo, sessionSvc, nil, logger)

	bedrockSvc := NewMockBedrockIntegration()
	knowledgeBase := NewMockKnowledgeBaseIntegration()

	chatHandler := handlers.NewChatHandler(
		logger,
		bedrockSvc,
		knowledgeBase,
		sessionSvc,
		chatSvc,
		nil,
		testJWTSecret,
	)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add routes
	api := router.Group("/api/v1")
	{
		admin := api.Group("/admin")
		{
			chat := admin.Group("/chat")
			{
				chat.POST("/sessions", chatHandler.CreateChatSession)
				chat.GET("/sessions", chatHandler.GetChatSessions)
				chat.GET("/sessions/:sessionId", chatHandler.GetChatSession)
				chat.POST("/sessions/:sessionId/messages", chatHandler.SendMessage)
				chat.GET("/sessions/:sessionId/history", chatHandler.GetSessionHistory)
				chat.GET("/ws", chatHandler.HandleWebSocket)
			}
		}
	}

	server := httptest.NewServer(router)

	return &LoadTestSuite{
		server:      server,
		chatHandler: chatHandler,
		metrics:     NewPerformanceMetrics(),
	}
}

func (suite *LoadTestSuite) tearDown() {
	suite.server.Close()
}

// Load test for session creation
func TestLoadTest_SessionCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	suite := setupLoadTestSuite(t)
	defer suite.tearDown()

	concurrency := 100
	requestsPerWorker := 50

	var wg sync.WaitGroup
	metrics := NewPerformanceMetrics()

	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			for j := 0; j < requestsPerWorker; j++ {
				sessionData := map[string]interface{}{
					"client_name": fmt.Sprintf("Load Test Client %d-%d", workerID, j),
					"context":     fmt.Sprintf("Load test session %d-%d", workerID, j),
				}

				jsonData, _ := json.Marshal(sessionData)

				reqStart := time.Now()
				req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/chat/sessions", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token")
				req = req.WithContext(context.WithValue(req.Context(), "user_id", fmt.Sprintf("load-test-user-%d", workerID)))

				resp, err := client.Do(req)
				latency := time.Since(reqStart)

				success := err == nil && resp != nil && resp.StatusCode == http.StatusCreated
				if resp != nil {
					resp.Body.Close()
				}

				metrics.RecordRequest(latency, success)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	metrics.Calculate(duration)

	t.Logf("Session Creation Load Test Results:\n%s", metrics.String())

	// Performance assertions
	assert.Greater(t, metrics.RequestsPerSecond, 50.0, "Should handle at least 50 requests/second")
	assert.Less(t, metrics.ErrorRate, 5.0, "Error rate should be less than 5%")
	assert.Less(t, metrics.AverageLatency, 200*time.Millisecond, "Average latency should be less than 200ms")
	assert.Less(t, metrics.P95Latency, 500*time.Millisecond, "P95 latency should be less than 500ms")
}

// Load test for message sending
func TestLoadTest_MessageSending(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	suite := setupLoadTestSuite(t)
	defer suite.tearDown()

	// Pre-create sessions for testing
	ctx := context.Background()
	sessionCount := 20
	sessions := make([]*domain.ChatSession, sessionCount)

	for i := 0; i < sessionCount; i++ {
		session := &domain.ChatSession{
			ID:           fmt.Sprintf("load-test-session-%d", i),
			UserID:       fmt.Sprintf("load-test-user-%d", i),
			ClientName:   fmt.Sprintf("Load Test Client %d", i),
			Status:       domain.SessionStatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LastActivity: time.Now(),
			ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
		}
		sessions[i] = session

		// Create session through service
		err := suite.chatHandler.(*handlers.ChatHandler).GetSessionService().CreateSession(ctx, session)
		require.NoError(t, err)
	}

	concurrency := 50
	messagesPerWorker := 20
	metrics := NewPerformanceMetrics()

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 15 * time.Second,
			}

			sessionIndex := workerID % sessionCount
			session := sessions[sessionIndex]

			for j := 0; j < messagesPerWorker; j++ {
				messageData := map[string]interface{}{
					"message":    fmt.Sprintf("Load test message %d from worker %d", j, workerID),
					"session_id": session.ID,
				}

				jsonData, _ := json.Marshal(messageData)

				reqStart := time.Now()
				req, _ := http.NewRequest("POST",
					fmt.Sprintf("%s/api/v1/admin/chat/sessions/%s/messages", suite.server.URL, session.ID),
					bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token")
				req = req.WithContext(context.WithValue(req.Context(), "user_id", session.UserID))

				resp, err := client.Do(req)
				latency := time.Since(reqStart)

				success := err == nil && resp != nil && resp.StatusCode == http.StatusOK
				if resp != nil {
					resp.Body.Close()
				}

				metrics.RecordRequest(latency, success)

				// Small delay to simulate realistic usage
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	metrics.Calculate(duration)

	t.Logf("Message Sending Load Test Results:\n%s", metrics.String())

	// Performance assertions
	assert.Greater(t, metrics.RequestsPerSecond, 30.0, "Should handle at least 30 messages/second")
	assert.Less(t, metrics.ErrorRate, 5.0, "Error rate should be less than 5%")
	assert.Less(t, metrics.AverageLatency, 500*time.Millisecond, "Average latency should be less than 500ms")
	assert.Less(t, metrics.P95Latency, 1*time.Second, "P95 latency should be less than 1 second")
}

// WebSocket load test
func TestLoadTest_WebSocketConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	suite := setupLoadTestSuite(t)
	defer suite.tearDown()

	concurrency := 25
	messagesPerConnection := 10
	metrics := NewPerformanceMetrics()

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			wsURL := fmt.Sprintf("ws%s/api/v1/admin/chat/ws?token=test-token",
				suite.server.URL[4:]) // Remove "http" prefix

			dialer := websocket.Dialer{
				HandshakeTimeout: 10 * time.Second,
			}

			conn, _, err := dialer.Dial(wsURL, nil)
			if err != nil {
				metrics.RecordRequest(0, false)
				return
			}
			defer conn.Close()

			for j := 0; j < messagesPerConnection; j++ {
				message := map[string]interface{}{
					"message":     fmt.Sprintf("WebSocket load test message %d from worker %d", j, workerID),
					"client_name": fmt.Sprintf("WS Load Test Client %d", workerID),
				}

				reqStart := time.Now()

				err := conn.WriteJSON(message)
				if err != nil {
					metrics.RecordRequest(time.Since(reqStart), false)
					continue
				}

				// Read response
				var response map[string]interface{}
				err = conn.ReadJSON(&response)
				latency := time.Since(reqStart)

				success := err == nil && response["success"] == true
				metrics.RecordRequest(latency, success)

				// Small delay between messages
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	metrics.Calculate(duration)

	t.Logf("WebSocket Load Test Results:\n%s", metrics.String())

	// Performance assertions
	assert.Greater(t, metrics.RequestsPerSecond, 20.0, "Should handle at least 20 WebSocket messages/second")
	assert.Less(t, metrics.ErrorRate, 10.0, "Error rate should be less than 10%")
	assert.Less(t, metrics.AverageLatency, 1*time.Second, "Average latency should be less than 1 second")
}

// Stress test with gradual ramp-up
func TestStressTest_GradualRampUp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	suite := setupLoadTestSuite(t)
	defer suite.tearDown()

	maxConcurrency := 100
	rampUpDuration := 30 * time.Second
	testDuration := 60 * time.Second

	metrics := NewPerformanceMetrics()
	var activeWorkers int64
	var wg sync.WaitGroup
	stopChan := make(chan bool)

	startTime := time.Now()

	// Gradual ramp-up
	go func() {
		rampUpInterval := rampUpDuration / time.Duration(maxConcurrency)

		for i := 0; i < maxConcurrency; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				for {
					select {
					case <-stopChan:
						return
					default:
						sessionData := map[string]interface{}{
							"client_name": fmt.Sprintf("Stress Test Client %d", workerID),
							"context":     fmt.Sprintf("Stress test session %d", workerID),
						}

						jsonData, _ := json.Marshal(sessionData)

						reqStart := time.Now()
						req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/chat/sessions", bytes.NewBuffer(jsonData))
						req.Header.Set("Content-Type", "application/json")
						req.Header.Set("Authorization", "Bearer test-token")
						req = req.WithContext(context.WithValue(req.Context(), "user_id", fmt.Sprintf("stress-test-user-%d", workerID)))

						resp, err := client.Do(req)
						latency := time.Since(reqStart)

						success := err == nil && resp != nil && resp.StatusCode == http.StatusCreated
						if resp != nil {
							resp.Body.Close()
						}

						metrics.RecordRequest(latency, success)

						// Delay between requests
						time.Sleep(100 * time.Millisecond)
					}
				}
			}(i)

			time.Sleep(rampUpInterval)
		}
	}()

	// Stop test after duration
	time.Sleep(testDuration)
	close(stopChan)
	wg.Wait()

	duration := time.Since(startTime)
	metrics.Calculate(duration)

	t.Logf("Stress Test Results:\n%s", metrics.String())

	// Stress test assertions (more lenient than load tests)
	assert.Greater(t, metrics.RequestsPerSecond, 10.0, "Should handle at least 10 requests/second under stress")
	assert.Less(t, metrics.ErrorRate, 20.0, "Error rate should be less than 20% under stress")
	assert.Less(t, metrics.P99Latency, 5*time.Second, "P99 latency should be less than 5 seconds under stress")
}

// Memory usage test
func TestMemoryUsage_UnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory usage test in short mode")
	}

	suite := setupLoadTestSuite(t)
	defer suite.tearDown()

	// This test would require runtime memory profiling
	// For now, we'll create many sessions and messages to test memory handling

	ctx := context.Background()
	sessionCount := 1000
	messagesPerSession := 50

	startTime := time.Now()

	// Create many sessions and messages
	for i := 0; i < sessionCount; i++ {
		session := &domain.ChatSession{
			ID:           fmt.Sprintf("memory-test-session-%d", i),
			UserID:       fmt.Sprintf("memory-test-user-%d", i),
			ClientName:   fmt.Sprintf("Memory Test Client %d", i),
			Context:      fmt.Sprintf("Memory test context with some longer text to use more memory %d", i),
			Status:       domain.SessionStatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LastActivity: time.Now(),
			ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
			Metadata: map[string]interface{}{
				"test_data": fmt.Sprintf("Additional metadata for session %d", i),
				"numbers":   []int{1, 2, 3, 4, 5, i},
			},
		}

		// Create session
		sessionSvc := suite.chatHandler.(*handlers.ChatHandler).GetSessionService()
		err := sessionSvc.CreateSession(ctx, session)
		require.NoError(t, err)

		// Create messages for this session
		for j := 0; j < messagesPerSession; j++ {
			message := &domain.ChatMessage{
				ID:        fmt.Sprintf("memory-test-message-%d-%d", i, j),
				SessionID: session.ID,
				Type:      domain.MessageTypeUser,
				Content:   fmt.Sprintf("Memory test message %d in session %d with some additional content to use more memory", j, i),
				Status:    domain.MessageStatusSent,
				CreatedAt: time.Now(),
				Metadata: map[string]interface{}{
					"message_index": j,
					"session_index": i,
				},
			}

			chatSvc := suite.chatHandler.(*handlers.ChatHandler).GetChatService()
			_, err := chatSvc.SendMessage(ctx, &domain.ChatRequest{
				SessionID: session.ID,
				Content:   message.Content,
				Type:      message.Type,
			})
			require.NoError(t, err)
		}

		// Log progress every 100 sessions
		if (i+1)%100 == 0 {
			t.Logf("Created %d sessions with %d messages each", i+1, messagesPerSession)
		}
	}

	duration := time.Since(startTime)
	totalMessages := sessionCount * messagesPerSession

	t.Logf("Memory Usage Test Completed:")
	t.Logf("  Created %d sessions", sessionCount)
	t.Logf("  Created %d messages", totalMessages)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Sessions/second: %.2f", float64(sessionCount)/duration.Seconds())
	t.Logf("  Messages/second: %.2f", float64(totalMessages)/duration.Seconds())

	// Test that we can still query data efficiently
	queryStart := time.Now()

	// Query some sessions
	for i := 0; i < 10; i++ {
		sessionID := fmt.Sprintf("memory-test-session-%d", i*100)
		sessionSvc := suite.chatHandler.(*handlers.ChatHandler).GetSessionService()
		session, err := sessionSvc.GetSession(ctx, sessionID)
		require.NoError(t, err)
		assert.Equal(t, sessionID, session.ID)
	}

	queryDuration := time.Since(queryStart)
	t.Logf("Query performance: %v for 10 session lookups", queryDuration)

	// Query should still be fast even with lots of data
	assert.Less(t, queryDuration, 100*time.Millisecond, "Queries should remain fast even with lots of data")
}

// Benchmark tests
func BenchmarkChatAPI_SessionCreation(b *testing.B) {
	suite := setupLoadTestSuite(&testing.T{})
	defer suite.tearDown()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		workerID := 0
		for pb.Next() {
			sessionData := map[string]interface{}{
				"client_name": fmt.Sprintf("Benchmark Client %d", workerID),
				"context":     fmt.Sprintf("Benchmark session %d", workerID),
			}

			jsonData, _ := json.Marshal(sessionData)
			req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/chat/sessions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			req = req.WithContext(context.WithValue(req.Context(), "user_id", fmt.Sprintf("benchmark-user-%d", workerID)))

			resp, err := client.Do(req)
			if err != nil {
				b.Error(err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				b.Errorf("Expected status 201, got %d", resp.StatusCode)
			}

			workerID++
		}
	})
}

func BenchmarkChatAPI_MessageSending(b *testing.B) {
	suite := setupLoadTestSuite(&testing.T{})
	defer suite.tearDown()

	// Pre-create a session
	ctx := context.Background()
	session := &domain.ChatSession{
		ID:           "benchmark-session",
		UserID:       "benchmark-user",
		ClientName:   "Benchmark Client",
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
	}

	sessionSvc := suite.chatHandler.(*handlers.ChatHandler).GetSessionService()
	err := sessionSvc.CreateSession(ctx, session)
	if err != nil {
		b.Fatal(err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		messageID := 0
		for pb.Next() {
			messageData := map[string]interface{}{
				"message":    fmt.Sprintf("Benchmark message %d", messageID),
				"session_id": session.ID,
			}

			jsonData, _ := json.Marshal(messageData)
			req, _ := http.NewRequest("POST",
				fmt.Sprintf("%s/api/v1/admin/chat/sessions/%s/messages", suite.server.URL, session.ID),
				bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")
			req = req.WithContext(context.WithValue(req.Context(), "user_id", session.UserID))

			resp, err := client.Do(req)
			if err != nil {
				b.Error(err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			messageID++
		}
	})
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}
