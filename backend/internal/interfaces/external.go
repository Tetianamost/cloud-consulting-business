package interfaces

import (
	"context"
	"io"
)

// LLMService defines the interface for Large Language Model integration
type LLMService interface {
	GenerateText(ctx context.Context, prompt string, options *GenerationOptions) (*GenerationResult, error)
	GenerateReport(ctx context.Context, inquiry interface{}, template string) (*GenerationResult, error)
	ValidateResponse(response string) error
	GetModelInfo() ModelInfo
	IsHealthy(ctx context.Context) bool
}

// StorageService defines the interface for file storage operations
type StorageService interface {
	Upload(ctx context.Context, key string, data io.Reader, contentType string) (*UploadResult, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GeneratePresignedURL(ctx context.Context, key string, expiration int64) (string, error)
	ListFiles(ctx context.Context, prefix string) ([]FileInfo, error)
	GetFileInfo(ctx context.Context, key string) (*FileInfo, error)
	IsHealthy(ctx context.Context) bool
}

// MessageQueueService defines the interface for message queue operations
type MessageQueueService interface {
	Publish(ctx context.Context, topic string, message *QueueMessage) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	CreateTopic(ctx context.Context, topic string) error
	DeleteTopic(ctx context.Context, topic string) error
	GetQueueInfo(ctx context.Context, topic string) (*QueueInfo, error)
	IsHealthy(ctx context.Context) bool
}

// NotificationChannel defines the interface for notification channels
type NotificationChannel interface {
	Send(ctx context.Context, message *Message) error
	GetChannelType() ChannelType
	IsHealthy() bool
	GetConfiguration() map[string]interface{}
	ValidateConfiguration() error
}

// HookHandler defines the interface for agent hook handlers
type HookHandler interface {
	Execute(ctx context.Context, payload interface{}) (*HookResult, error)
	GetMetadata() HookMetadata
	Validate(payload interface{}) error
	GetRequiredPermissions() []string
}

// CacheService defines the interface for caching operations
type CacheService interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration int64) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context, pattern string) error
	GetTTL(ctx context.Context, key string) (int64, error)
	IsHealthy(ctx context.Context) bool
}

// MetricsService defines the interface for metrics collection
type MetricsService interface {
	IncrementCounter(name string, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
	RecordDuration(name string, duration int64, labels map[string]string)
	GetMetrics() map[string]interface{}
}

// Supporting types for external services

// GenerationOptions represents options for LLM text generation
type GenerationOptions struct {
	Model       string                 `json:"model,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	TopP        float64                `json:"top_p,omitempty"`
	Stop        []string               `json:"stop,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// GenerationResult represents the result of LLM text generation
type GenerationResult struct {
	Text         string                 `json:"text"`
	TokensUsed   int                    `json:"tokens_used"`
	Model        string                 `json:"model"`
	FinishReason string                 `json:"finish_reason"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Error        *string                `json:"error,omitempty"`
}

// ModelInfo represents information about an LLM model
type ModelInfo struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	MaxTokens    int      `json:"max_tokens"`
	Capabilities []string `json:"capabilities"`
	Provider     string   `json:"provider"`
}

// UploadResult represents the result of a file upload
type UploadResult struct {
	Key         string `json:"key"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	ETag        string `json:"etag"`
	UploadedAt  string `json:"uploaded_at"`
}

// FileInfo represents information about a stored file
type FileInfo struct {
	Key          string `json:"key"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	LastModified string `json:"last_modified"`
	ETag         string `json:"etag"`
	URL          string `json:"url,omitempty"`
}

// QueueMessage represents a message in the queue
type QueueMessage struct {
	ID          string                 `json:"id"`
	Body        string                 `json:"body"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	DelaySeconds int                   `json:"delay_seconds,omitempty"`
	Priority    int                    `json:"priority,omitempty"`
}

// MessageHandler defines the function signature for message handlers
type MessageHandler func(ctx context.Context, message *QueueMessage) error

// QueueInfo represents information about a message queue
type QueueInfo struct {
	Name                string `json:"name"`
	MessageCount        int64  `json:"message_count"`
	VisibleMessages     int64  `json:"visible_messages"`
	InFlightMessages    int64  `json:"in_flight_messages"`
	DelayedMessages     int64  `json:"delayed_messages"`
	CreatedAt           string `json:"created_at"`
	LastModified        string `json:"last_modified"`
}

// Message represents a notification message
type Message struct {
	To          string                 `json:"to"`
	Subject     string                 `json:"subject"`
	Body        string                 `json:"body"`
	ContentType string                 `json:"content_type"`
	Attachments []Attachment           `json:"attachments,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Attachment represents a message attachment
type Attachment struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Data        []byte `json:"data"`
	URL         string `json:"url,omitempty"`
}

// HookMetadata represents metadata about a hook handler
type HookMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Async       bool     `json:"async"`
	Timeout     int      `json:"timeout_seconds"`
}