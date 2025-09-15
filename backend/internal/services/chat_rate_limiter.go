package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// RateLimitEntry represents a rate limit entry
type RateLimitEntry struct {
	Requests  []time.Time
	Limit     int
	Window    time.Duration
	LastReset time.Time
}

// ChatRateLimiter implements the ChatRateLimiter interface
type ChatRateLimiter struct {
	logger       *logrus.Logger
	limits       map[string]*RateLimitEntry
	customLimits map[string]*RateLimitEntry
	mutex        sync.RWMutex

	// Default limits
	messageLimit          int
	messageWindow         time.Duration
	connectionLimit       int
	connectionWindow      time.Duration
	sessionCreationLimit  int
	sessionCreationWindow time.Duration
}

// NewChatRateLimiter creates a new chat rate limiter
func NewChatRateLimiter(logger *logrus.Logger) *ChatRateLimiter {
	limiter := &ChatRateLimiter{
		logger:                logger,
		limits:                make(map[string]*RateLimitEntry),
		customLimits:          make(map[string]*RateLimitEntry),
		messageLimit:          60, // 60 messages per minute
		messageWindow:         time.Minute,
		connectionLimit:       10, // 10 connections per minute
		connectionWindow:      time.Minute,
		sessionCreationLimit:  5, // 5 sessions per hour
		sessionCreationWindow: time.Hour,
	}

	// Start cleanup routine
	go limiter.cleanupExpiredEntries()

	return limiter
}

// AllowMessage checks if a message is allowed for a user
func (r *ChatRateLimiter) AllowMessage(ctx context.Context, userID string) (*interfaces.RateLimitResult, error) {
	key := fmt.Sprintf("message:%s", userID)
	return r.checkAndUpdateLimit(key, r.messageLimit, r.messageWindow)
}

// AllowConnection checks if a connection is allowed for a user
func (r *ChatRateLimiter) AllowConnection(ctx context.Context, userID string) (*interfaces.RateLimitResult, error) {
	key := fmt.Sprintf("connection:%s", userID)
	return r.checkAndUpdateLimit(key, r.connectionLimit, r.connectionWindow)
}

// AllowSessionCreation checks if session creation is allowed for a user
func (r *ChatRateLimiter) AllowSessionCreation(ctx context.Context, userID string) (*interfaces.RateLimitResult, error) {
	key := fmt.Sprintf("session_creation:%s", userID)
	return r.checkAndUpdateLimit(key, r.sessionCreationLimit, r.sessionCreationWindow)
}

// CheckLimit checks a custom rate limit
func (r *ChatRateLimiter) CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (*interfaces.RateLimitResult, error) {
	return r.checkAndUpdateLimit(key, limit, window)
}

// SetCustomLimit sets a custom rate limit for a user and action
func (r *ChatRateLimiter) SetCustomLimit(ctx context.Context, userID string, action string, limit int, window time.Duration) error {
	key := fmt.Sprintf("custom:%s:%s", action, userID)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.customLimits[key] = &RateLimitEntry{
		Requests:  make([]time.Time, 0),
		Limit:     limit,
		Window:    window,
		LastReset: time.Now(),
	}

	r.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"action":  action,
		"limit":   limit,
		"window":  window,
	}).Info("Custom rate limit set")

	return nil
}

// GetLimitInfo gets information about a rate limit
func (r *ChatRateLimiter) GetLimitInfo(ctx context.Context, key string) (*interfaces.RateLimitInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Check custom limits first
	if entry, exists := r.customLimits[key]; exists {
		return r.createLimitInfo(key, entry), nil
	}

	// Check regular limits
	if entry, exists := r.limits[key]; exists {
		return r.createLimitInfo(key, entry), nil
	}

	// Return default info if not found
	return &interfaces.RateLimitInfo{
		Key:           key,
		Limit:         0,
		Remaining:     0,
		ResetTime:     time.Now(),
		WindowSize:    0,
		TotalRequests: 0,
	}, nil
}

// ResetUserLimits resets all rate limits for a user
func (r *ChatRateLimiter) ResetUserLimits(ctx context.Context, userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Reset regular limits
	for key, entry := range r.limits {
		if r.keyBelongsToUser(key, userID) {
			entry.Requests = make([]time.Time, 0)
			entry.LastReset = time.Now()
		}
	}

	// Reset custom limits
	for key, entry := range r.customLimits {
		if r.keyBelongsToUser(key, userID) {
			entry.Requests = make([]time.Time, 0)
			entry.LastReset = time.Now()
		}
	}

	r.logger.WithField("user_id", userID).Info("User rate limits reset")
	return nil
}

// GetUserLimitStatus gets the status of all limits for a user
func (r *ChatRateLimiter) GetUserLimitStatus(ctx context.Context, userID string) (map[string]*interfaces.RateLimitInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	status := make(map[string]*interfaces.RateLimitInfo)

	// Check regular limits
	for key, entry := range r.limits {
		if r.keyBelongsToUser(key, userID) {
			status[key] = r.createLimitInfo(key, entry)
		}
	}

	// Check custom limits
	for key, entry := range r.customLimits {
		if r.keyBelongsToUser(key, userID) {
			status[key] = r.createLimitInfo(key, entry)
		}
	}

	return status, nil
}

// Helper methods

// checkAndUpdateLimit checks and updates a rate limit
func (r *ChatRateLimiter) checkAndUpdateLimit(key string, limit int, window time.Duration) (*interfaces.RateLimitResult, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	// Get or create entry
	entry, exists := r.limits[key]
	if !exists {
		entry = &RateLimitEntry{
			Requests:  make([]time.Time, 0),
			Limit:     limit,
			Window:    window,
			LastReset: now,
		}
		r.limits[key] = entry
	}

	// Clean up old requests
	cutoff := now.Add(-window)
	validRequests := make([]time.Time, 0)
	for _, reqTime := range entry.Requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	entry.Requests = validRequests

	// Check if limit is exceeded
	allowed := len(entry.Requests) < limit
	remaining := limit - len(entry.Requests)
	if remaining < 0 {
		remaining = 0
	}

	// Calculate reset time
	resetTime := now.Add(window)
	if len(entry.Requests) > 0 {
		oldestRequest := entry.Requests[0]
		resetTime = oldestRequest.Add(window)
	}

	// Calculate retry after
	var retryAfter time.Duration
	if !allowed {
		retryAfter = time.Until(resetTime)
		if retryAfter < 0 {
			retryAfter = 0
		}
	}

	// Add current request if allowed
	if allowed {
		entry.Requests = append(entry.Requests, now)
		remaining--
	}

	result := &interfaces.RateLimitResult{
		Allowed:       allowed,
		Remaining:     remaining,
		ResetTime:     resetTime,
		RetryAfter:    retryAfter,
		TotalRequests: len(entry.Requests),
	}

	// Log rate limit check
	r.logger.WithFields(logrus.Fields{
		"key":       key,
		"allowed":   allowed,
		"remaining": remaining,
		"total":     len(entry.Requests),
		"limit":     limit,
	}).Debug("Rate limit checked")

	return result, nil
}

// createLimitInfo creates a RateLimitInfo from an entry
func (r *ChatRateLimiter) createLimitInfo(key string, entry *RateLimitEntry) *interfaces.RateLimitInfo {
	now := time.Now()
	cutoff := now.Add(-entry.Window)

	// Count valid requests
	validRequests := 0
	for _, reqTime := range entry.Requests {
		if reqTime.After(cutoff) {
			validRequests++
		}
	}

	remaining := entry.Limit - validRequests
	if remaining < 0 {
		remaining = 0
	}

	// Calculate reset time
	resetTime := now.Add(entry.Window)
	if len(entry.Requests) > 0 {
		oldestRequest := entry.Requests[0]
		resetTime = oldestRequest.Add(entry.Window)
	}

	return &interfaces.RateLimitInfo{
		Key:           key,
		Limit:         entry.Limit,
		Remaining:     remaining,
		ResetTime:     resetTime,
		WindowSize:    entry.Window,
		TotalRequests: validRequests,
	}
}

// keyBelongsToUser checks if a rate limit key belongs to a user
func (r *ChatRateLimiter) keyBelongsToUser(key, userID string) bool {
	return fmt.Sprintf(":%s", userID) == key[len(key)-len(userID)-1:] ||
		fmt.Sprintf(":%s:", userID) != "" &&
			(key == fmt.Sprintf("message:%s", userID) ||
				key == fmt.Sprintf("connection:%s", userID) ||
				key == fmt.Sprintf("session_creation:%s", userID) ||
				(len(key) > len(userID)+1 && key[len(key)-len(userID)-1:] == fmt.Sprintf(":%s", userID)))
}

// cleanupExpiredEntries periodically removes expired rate limit entries
func (r *ChatRateLimiter) cleanupExpiredEntries() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mutex.Lock()
			now := time.Now()

			// Clean up regular limits
			for key, entry := range r.limits {
				// Remove entries that haven't been used in 2x their window
				if now.Sub(entry.LastReset) > 2*entry.Window {
					delete(r.limits, key)
					continue
				}

				// Clean up old requests within entries
				cutoff := now.Add(-entry.Window)
				validRequests := make([]time.Time, 0)
				for _, reqTime := range entry.Requests {
					if reqTime.After(cutoff) {
						validRequests = append(validRequests, reqTime)
					}
				}
				entry.Requests = validRequests
			}

			// Clean up custom limits
			for key, entry := range r.customLimits {
				if now.Sub(entry.LastReset) > 2*entry.Window {
					delete(r.customLimits, key)
					continue
				}

				cutoff := now.Add(-entry.Window)
				validRequests := make([]time.Time, 0)
				for _, reqTime := range entry.Requests {
					if reqTime.After(cutoff) {
						validRequests = append(validRequests, reqTime)
					}
				}
				entry.Requests = validRequests
			}

			r.mutex.Unlock()

			r.logger.Debug("Rate limiter cleanup completed")
		}
	}
}

// SetMessageLimit sets the default message rate limit
func (r *ChatRateLimiter) SetMessageLimit(limit int, window time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.messageLimit = limit
	r.messageWindow = window

	r.logger.WithFields(logrus.Fields{
		"limit":  limit,
		"window": window,
	}).Info("Message rate limit updated")
}

// SetConnectionLimit sets the default connection rate limit
func (r *ChatRateLimiter) SetConnectionLimit(limit int, window time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.connectionLimit = limit
	r.connectionWindow = window

	r.logger.WithFields(logrus.Fields{
		"limit":  limit,
		"window": window,
	}).Info("Connection rate limit updated")
}

// SetSessionCreationLimit sets the default session creation rate limit
func (r *ChatRateLimiter) SetSessionCreationLimit(limit int, window time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.sessionCreationLimit = limit
	r.sessionCreationWindow = window

	r.logger.WithFields(logrus.Fields{
		"limit":  limit,
		"window": window,
	}).Info("Session creation rate limit updated")
}
