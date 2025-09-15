package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ConsultantSessionLoadBalancer manages load balancing for multiple concurrent consultant sessions
type ConsultantSessionLoadBalancer struct {
	logger      *logrus.Logger
	sessions    map[string]*SessionInfo
	consultants map[string]*ConsultantInfo
	mu          sync.RWMutex

	// Load balancing configuration
	maxSessionsPerConsultant int
	sessionTimeout           time.Duration
	loadBalancingStrategy    LoadBalancingStrategy

	// Metrics
	totalSessions    int64
	activeSessions   int64
	balancedSessions int64
	rejectedSessions int64
}

// SessionInfo represents information about a consultant session
type SessionInfo struct {
	SessionID    string                 `json:"session_id"`
	ConsultantID string                 `json:"consultant_id"`
	ClientName   string                 `json:"client_name"`
	StartTime    time.Time              `json:"start_time"`
	LastActivity time.Time              `json:"last_activity"`
	RequestCount int                    `json:"request_count"`
	Priority     SessionPriority        `json:"priority"`
	Metadata     map[string]interface{} `json:"metadata"`
	Status       SessionStatus          `json:"status"`
}

// ConsultantInfo represents information about a consultant
type ConsultantInfo struct {
	ConsultantID    string                 `json:"consultant_id"`
	Name            string                 `json:"name"`
	ActiveSessions  []string               `json:"active_sessions"`
	MaxSessions     int                    `json:"max_sessions"`
	CurrentLoad     float64                `json:"current_load"`
	Specializations []string               `json:"specializations"`
	Status          ConsultantStatus       `json:"status"`
	LastActivity    time.Time              `json:"last_activity"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// LoadBalancingStrategy defines the strategy for load balancing
type LoadBalancingStrategy string

const (
	StrategyRoundRobin     LoadBalancingStrategy = "round_robin"
	StrategyLeastLoaded    LoadBalancingStrategy = "least_loaded"
	StrategySpecialization LoadBalancingStrategy = "specialization"
	StrategyPriority       LoadBalancingStrategy = "priority"
)

// SessionPriority defines session priority levels
type SessionPriority string

const (
	PriorityLow      SessionPriority = "low"
	PriorityNormal   SessionPriority = "normal"
	PriorityHigh     SessionPriority = "high"
	PriorityCritical SessionPriority = "critical"
)

// SessionStatus defines session status
type SessionStatus string

const (
	StatusActive    SessionStatus = "active"
	StatusIdle      SessionStatus = "idle"
	StatusCompleted SessionStatus = "completed"
	StatusTimeout   SessionStatus = "timeout"
)

// ConsultantStatus defines consultant status
type ConsultantStatus string

const (
	ConsultantAvailable ConsultantStatus = "available"
	ConsultantBusy      ConsultantStatus = "busy"
	ConsultantOffline   ConsultantStatus = "offline"
)

// NewConsultantSessionLoadBalancer creates a new consultant session load balancer
func NewConsultantSessionLoadBalancer(logger *logrus.Logger) *ConsultantSessionLoadBalancer {
	lb := &ConsultantSessionLoadBalancer{
		logger:                   logger,
		sessions:                 make(map[string]*SessionInfo),
		consultants:              make(map[string]*ConsultantInfo),
		maxSessionsPerConsultant: 5,
		sessionTimeout:           30 * time.Minute,
		loadBalancingStrategy:    StrategyLeastLoaded,
	}

	// Initialize with default consultants (in production, this would come from a database)
	lb.initializeDefaultConsultants()

	return lb
}

// AssignSession assigns a session to the most appropriate consultant
func (lb *ConsultantSessionLoadBalancer) AssignSession(sessionID, preferredConsultantID string) string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.totalSessions++

	// Check if session already exists
	if existingSession, exists := lb.sessions[sessionID]; exists {
		existingSession.LastActivity = time.Now()
		existingSession.RequestCount++
		return existingSession.ConsultantID
	}

	// Find the best consultant for this session
	consultantID := lb.findBestConsultant(preferredConsultantID)
	if consultantID == "" {
		lb.rejectedSessions++
		lb.logger.WithField("session_id", sessionID).Warn("No available consultant for session")
		return ""
	}

	// Create new session
	session := &SessionInfo{
		SessionID:    sessionID,
		ConsultantID: consultantID,
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		RequestCount: 1,
		Priority:     PriorityNormal,
		Status:       StatusActive,
		Metadata:     make(map[string]interface{}),
	}

	// Assign session to consultant
	consultant := lb.consultants[consultantID]
	consultant.ActiveSessions = append(consultant.ActiveSessions, sessionID)
	consultant.CurrentLoad = float64(len(consultant.ActiveSessions)) / float64(consultant.MaxSessions)
	consultant.LastActivity = time.Now()

	// Update consultant status based on load
	if consultant.CurrentLoad >= 1.0 {
		consultant.Status = ConsultantBusy
	} else if consultant.CurrentLoad > 0 {
		consultant.Status = ConsultantAvailable
	}

	// Store session
	lb.sessions[sessionID] = session
	lb.activeSessions++
	lb.balancedSessions++

	lb.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"consultant_id": consultantID,
		"load":          consultant.CurrentLoad,
		"strategy":      lb.loadBalancingStrategy,
	}).Info("Session assigned to consultant")

	return consultantID
}

// ReleaseSession releases a session from a consultant
func (lb *ConsultantSessionLoadBalancer) ReleaseSession(sessionID string) error {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	session, exists := lb.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Update session status
	session.Status = StatusCompleted

	// Remove session from consultant
	consultant := lb.consultants[session.ConsultantID]
	consultant.ActiveSessions = lb.removeSessionFromSlice(consultant.ActiveSessions, sessionID)
	consultant.CurrentLoad = float64(len(consultant.ActiveSessions)) / float64(consultant.MaxSessions)

	// Update consultant status
	if len(consultant.ActiveSessions) == 0 {
		consultant.Status = ConsultantAvailable
	} else if consultant.CurrentLoad < 1.0 {
		consultant.Status = ConsultantAvailable
	}

	// Remove from active sessions
	delete(lb.sessions, sessionID)
	lb.activeSessions--

	lb.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"consultant_id": session.ConsultantID,
		"duration":      time.Since(session.StartTime).String(),
		"requests":      session.RequestCount,
	}).Info("Session released")

	return nil
}

// GetActiveSessionCount returns the number of active sessions
func (lb *ConsultantSessionLoadBalancer) GetActiveSessionCount() int {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	return len(lb.sessions)
}

// GetLoadBalancingMetrics returns load balancing metrics
func (lb *ConsultantSessionLoadBalancer) GetLoadBalancingMetrics() *LoadBalancingMetrics {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// Calculate consultant utilization
	var totalLoad float64
	availableConsultants := 0
	busyConsultants := 0

	for _, consultant := range lb.consultants {
		totalLoad += consultant.CurrentLoad
		switch consultant.Status {
		case ConsultantAvailable:
			availableConsultants++
		case ConsultantBusy:
			busyConsultants++
		}
	}

	avgLoad := float64(0)
	if len(lb.consultants) > 0 {
		avgLoad = totalLoad / float64(len(lb.consultants))
	}

	return &LoadBalancingMetrics{
		TotalSessions:        lb.totalSessions,
		ActiveSessions:       lb.activeSessions,
		BalancedSessions:     lb.balancedSessions,
		RejectedSessions:     lb.rejectedSessions,
		TotalConsultants:     int64(len(lb.consultants)),
		AvailableConsultants: int64(availableConsultants),
		BusyConsultants:      int64(busyConsultants),
		AverageLoad:          avgLoad,
		Strategy:             string(lb.loadBalancingStrategy),
		Timestamp:            time.Now(),
	}
}

// OptimizeForResponseTime optimizes load balancer settings for response time
func (lb *ConsultantSessionLoadBalancer) OptimizeForResponseTime() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Reduce max sessions per consultant for better response times
	if lb.maxSessionsPerConsultant > 3 {
		lb.maxSessionsPerConsultant = 3

		// Update consultant max sessions
		for _, consultant := range lb.consultants {
			consultant.MaxSessions = lb.maxSessionsPerConsultant
			consultant.CurrentLoad = float64(len(consultant.ActiveSessions)) / float64(consultant.MaxSessions)
		}

		lb.logger.WithField("max_sessions", lb.maxSessionsPerConsultant).Info("Optimized load balancer for response time")
	}

	// Switch to least loaded strategy for better distribution
	if lb.loadBalancingStrategy != StrategyLeastLoaded {
		lb.loadBalancingStrategy = StrategyLeastLoaded
		lb.logger.Info("Switched to least loaded strategy for response time optimization")
	}
}

// CleanupExpiredSessions removes expired sessions
func (lb *ConsultantSessionLoadBalancer) CleanupExpiredSessions(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lb.performCleanup()
		}
	}
}

// Private helper methods

func (lb *ConsultantSessionLoadBalancer) findBestConsultant(preferredConsultantID string) string {
	// Check if preferred consultant is available
	if preferredConsultantID != "" {
		if consultant, exists := lb.consultants[preferredConsultantID]; exists {
			if consultant.Status == ConsultantAvailable && len(consultant.ActiveSessions) < consultant.MaxSessions {
				return preferredConsultantID
			}
		}
	}

	// Apply load balancing strategy
	switch lb.loadBalancingStrategy {
	case StrategyLeastLoaded:
		return lb.findLeastLoadedConsultant()
	case StrategyRoundRobin:
		return lb.findRoundRobinConsultant()
	case StrategySpecialization:
		return lb.findSpecializedConsultant("")
	case StrategyPriority:
		return lb.findPriorityConsultant()
	default:
		return lb.findLeastLoadedConsultant()
	}
}

func (lb *ConsultantSessionLoadBalancer) findLeastLoadedConsultant() string {
	var bestConsultant string
	lowestLoad := float64(2.0) // Higher than max possible load

	for id, consultant := range lb.consultants {
		if consultant.Status == ConsultantAvailable && consultant.CurrentLoad < lowestLoad {
			lowestLoad = consultant.CurrentLoad
			bestConsultant = id
		}
	}

	return bestConsultant
}

func (lb *ConsultantSessionLoadBalancer) findRoundRobinConsultant() string {
	// Simple round-robin implementation
	availableConsultants := make([]string, 0)

	for id, consultant := range lb.consultants {
		if consultant.Status == ConsultantAvailable && len(consultant.ActiveSessions) < consultant.MaxSessions {
			availableConsultants = append(availableConsultants, id)
		}
	}

	if len(availableConsultants) == 0 {
		return ""
	}

	// Use total sessions as round-robin counter
	index := int(lb.totalSessions) % len(availableConsultants)
	return availableConsultants[index]
}

func (lb *ConsultantSessionLoadBalancer) findSpecializedConsultant(specialization string) string {
	// For now, fall back to least loaded
	// In production, this would match specializations
	return lb.findLeastLoadedConsultant()
}

func (lb *ConsultantSessionLoadBalancer) findPriorityConsultant() string {
	// For now, fall back to least loaded
	// In production, this would consider consultant priorities
	return lb.findLeastLoadedConsultant()
}

func (lb *ConsultantSessionLoadBalancer) removeSessionFromSlice(sessions []string, sessionID string) []string {
	for i, id := range sessions {
		if id == sessionID {
			return append(sessions[:i], sessions[i+1:]...)
		}
	}
	return sessions
}

func (lb *ConsultantSessionLoadBalancer) performCleanup() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	expiredSessions := make([]string, 0)

	// Find expired sessions
	for sessionID, session := range lb.sessions {
		if now.Sub(session.LastActivity) > lb.sessionTimeout {
			expiredSessions = append(expiredSessions, sessionID)
		}
	}

	// Remove expired sessions
	for _, sessionID := range expiredSessions {
		session := lb.sessions[sessionID]
		session.Status = StatusTimeout

		// Remove from consultant
		consultant := lb.consultants[session.ConsultantID]
		consultant.ActiveSessions = lb.removeSessionFromSlice(consultant.ActiveSessions, sessionID)
		consultant.CurrentLoad = float64(len(consultant.ActiveSessions)) / float64(consultant.MaxSessions)

		// Update consultant status
		if len(consultant.ActiveSessions) == 0 {
			consultant.Status = ConsultantAvailable
		}

		delete(lb.sessions, sessionID)
		lb.activeSessions--
	}

	if len(expiredSessions) > 0 {
		lb.logger.WithField("expired_count", len(expiredSessions)).Info("Cleaned up expired sessions")
	}
}

func (lb *ConsultantSessionLoadBalancer) initializeDefaultConsultants() {
	// Initialize with default consultants
	defaultConsultants := []struct {
		id              string
		name            string
		specializations []string
	}{
		{"consultant-1", "Senior AWS Consultant", []string{"aws", "architecture", "migration"}},
		{"consultant-2", "Cloud Security Expert", []string{"security", "compliance", "aws", "azure"}},
		{"consultant-3", "Multi-Cloud Architect", []string{"aws", "azure", "gcp", "architecture"}},
		{"consultant-4", "DevOps Specialist", []string{"devops", "kubernetes", "ci/cd", "automation"}},
		{"consultant-5", "Cost Optimization Expert", []string{"cost", "optimization", "finops", "aws"}},
	}

	for _, consultant := range defaultConsultants {
		lb.consultants[consultant.id] = &ConsultantInfo{
			ConsultantID:    consultant.id,
			Name:            consultant.name,
			ActiveSessions:  make([]string, 0),
			MaxSessions:     lb.maxSessionsPerConsultant,
			CurrentLoad:     0.0,
			Specializations: consultant.specializations,
			Status:          ConsultantAvailable,
			LastActivity:    time.Now(),
			Metadata:        make(map[string]interface{}),
		}
	}

	lb.logger.WithField("consultant_count", len(lb.consultants)).Info("Initialized default consultants")
}

// LoadBalancingMetrics represents load balancing metrics
type LoadBalancingMetrics struct {
	TotalSessions        int64     `json:"total_sessions"`
	ActiveSessions       int64     `json:"active_sessions"`
	BalancedSessions     int64     `json:"balanced_sessions"`
	RejectedSessions     int64     `json:"rejected_sessions"`
	TotalConsultants     int64     `json:"total_consultants"`
	AvailableConsultants int64     `json:"available_consultants"`
	BusyConsultants      int64     `json:"busy_consultants"`
	AverageLoad          float64   `json:"average_load"`
	Strategy             string    `json:"strategy"`
	Timestamp            time.Time `json:"timestamp"`
}
