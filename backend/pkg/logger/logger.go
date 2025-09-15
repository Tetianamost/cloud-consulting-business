package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus with additional functionality
type Logger struct {
	*logrus.Logger
}

// New creates a new logger instance
func New(level logrus.Level, format string) *Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	
	// Set formatter based on format
	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}
	
	return &Logger{Logger: logger}
}

// WithContext adds context information to log entries
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.WithFields(logrus.Fields{})
	
	// Add trace ID if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		entry = entry.WithField("trace_id", traceID)
	}
	
	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		entry = entry.WithField("user_id", userID)
	}
	
	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}
	
	return entry
}

// WithInquiry adds inquiry context to log entries
func (l *Logger) WithInquiry(inquiryID string) *logrus.Entry {
	return l.WithField("inquiry_id", inquiryID)
}

// WithReport adds report context to log entries
func (l *Logger) WithReport(reportID string) *logrus.Entry {
	return l.WithField("report_id", reportID)
}

// WithService adds service context to log entries
func (l *Logger) WithService(serviceName string) *logrus.Entry {
	return l.WithField("service", serviceName)
}

// LogError logs an error with additional context
func (l *Logger) LogError(ctx context.Context, err error, message string, fields map[string]interface{}) {
	entry := l.WithContext(ctx).WithError(err)
	
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	
	entry.Error(message)
}

// LogInfo logs an info message with context
func (l *Logger) LogInfo(ctx context.Context, message string, fields map[string]interface{}) {
	entry := l.WithContext(ctx)
	
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	
	entry.Info(message)
}

// LogDebug logs a debug message with context
func (l *Logger) LogDebug(ctx context.Context, message string, fields map[string]interface{}) {
	entry := l.WithContext(ctx)
	
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	
	entry.Debug(message)
}

// LogWarn logs a warning message with context
func (l *Logger) LogWarn(ctx context.Context, message string, fields map[string]interface{}) {
	entry := l.WithContext(ctx)
	
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	
	entry.Warn(message)
}