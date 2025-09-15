package interfaces

import "context"

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	CheckHealth(ctx context.Context) *HealthStatus
	GetComponentName() string
}

// HealthService defines the interface for overall health management
type HealthService interface {
	RegisterChecker(checker HealthChecker)
	GetOverallHealth(ctx context.Context) *OverallHealthStatus
	GetComponentHealth(ctx context.Context, component string) *HealthStatus
}

// HealthStatus represents the health status of a component
type HealthStatus struct {
	Status    HealthStatusType       `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp string                 `json:"timestamp"`
	Duration  int64                  `json:"duration_ms"`
}

// OverallHealthStatus represents the overall system health
type OverallHealthStatus struct {
	Status     HealthStatusType         `json:"status"`
	Components map[string]*HealthStatus `json:"components"`
	Timestamp  string                   `json:"timestamp"`
	Version    string                   `json:"version"`
	Uptime     int64                    `json:"uptime_seconds"`
}

// HealthStatusType represents the health status type
type HealthStatusType string

const (
	HealthStatusHealthy   HealthStatusType = "healthy"
	HealthStatusDegraded  HealthStatusType = "degraded"
	HealthStatusUnhealthy HealthStatusType = "unhealthy"
	HealthStatusUnknown   HealthStatusType = "unknown"
)