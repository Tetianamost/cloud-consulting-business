package interfaces

import "context"

// ErrorHandler defines the interface for error handling
type ErrorHandler interface {
	HandleError(ctx context.Context, err error) *APIError
	LogError(ctx context.Context, err error, metadata map[string]interface{})
	ShouldRetry(err error) bool
}

// CircuitBreaker defines the interface for circuit breaker pattern
type CircuitBreaker interface {
	Execute(ctx context.Context, fn func() error) error
	GetState() CircuitState
	GetStats() CircuitStats
	Reset()
}

// RetryService defines the interface for retry logic
type RetryService interface {
	ExecuteWithRetry(ctx context.Context, fn func() error, config *RetryConfig) error
	GetRetryDelay(attempt int, config *RetryConfig) int64
}

// APIError represents a structured API error
type APIError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	TraceID    string                 `json:"trace_id,omitempty"`
	Timestamp  string                 `json:"timestamp"`
	StatusCode int                    `json:"-"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// CircuitState represents the state of a circuit breaker
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"
	CircuitStateOpen     CircuitState = "open"
	CircuitStateHalfOpen CircuitState = "half_open"
)

// CircuitStats represents circuit breaker statistics
type CircuitStats struct {
	State           CircuitState `json:"state"`
	FailureCount    int64        `json:"failure_count"`
	SuccessCount    int64        `json:"success_count"`
	RequestCount    int64        `json:"request_count"`
	LastFailureTime string       `json:"last_failure_time,omitempty"`
	NextRetryTime   string       `json:"next_retry_time,omitempty"`
}

// RetryConfig represents configuration for retry logic
type RetryConfig struct {
	MaxAttempts   int                    `json:"max_attempts"`
	InitialDelay  int64                  `json:"initial_delay_ms"`
	MaxDelay      int64                  `json:"max_delay_ms"`
	BackoffFactor float64                `json:"backoff_factor"`
	RetryIf       func(error) bool       `json:"-"`
	OnRetry       func(int, error)       `json:"-"`
	Context       map[string]interface{} `json:"context,omitempty"`
}

// Standard error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeRateLimit      = "RATE_LIMIT_EXCEEDED"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeServiceUnavail = "SERVICE_UNAVAILABLE"
	ErrCodeTimeout        = "TIMEOUT"
	ErrCodeBadRequest     = "BAD_REQUEST"
)