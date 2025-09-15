package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// RetentionPolicy defines email event retention rules
type RetentionPolicy struct {
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	RetentionDays      int           `json:"retention_days"`
	ArchiveDays        int           `json:"archive_days"`
	EmailTypes         []string      `json:"email_types,omitempty"`
	Statuses           []string      `json:"statuses,omitempty"`
	Enabled            bool          `json:"enabled"`
	LastExecuted       *time.Time    `json:"last_executed,omitempty"`
	NextExecution      *time.Time    `json:"next_execution,omitempty"`
	ExecutionFrequency time.Duration `json:"execution_frequency"`
}

// RetentionStats provides statistics about retention operations
type RetentionStats struct {
	PolicyName        string        `json:"policy_name"`
	ExecutionTime     time.Time     `json:"execution_time"`
	RecordsProcessed  int           `json:"records_processed"`
	RecordsArchived   int           `json:"records_archived"`
	RecordsDeleted    int           `json:"records_deleted"`
	ExecutionDuration time.Duration `json:"execution_duration"`
	ErrorMessage      string        `json:"error_message,omitempty"`
}

// EmailRetentionService manages email event retention policies
type EmailRetentionService struct {
	db       *sql.DB
	logger   *logrus.Logger
	policies map[string]*RetentionPolicy
	stats    []*RetentionStats

	// Configuration
	maxBatchSize     int
	maxExecutionTime time.Duration

	// State
	isRunning bool
	stopChan  chan struct{}
}

// NewEmailRetentionService creates a new email retention service
func NewEmailRetentionService(db *sql.DB, logger *logrus.Logger) *EmailRetentionService {
	service := &EmailRetentionService{
		db:               db,
		logger:           logger,
		policies:         make(map[string]*RetentionPolicy),
		stats:            make([]*RetentionStats, 0),
		maxBatchSize:     1000,
		maxExecutionTime: 30 * time.Minute,
		stopChan:         make(chan struct{}),
	}

	// Initialize default policies
	service.initializeDefaultPolicies()

	return service
}

// initializeDefaultPolicies sets up default retention policies
func (s *EmailRetentionService) initializeDefaultPolicies() {
	// Policy for customer confirmation emails (keep longer for audit)
	s.policies["customer_confirmations"] = &RetentionPolicy{
		Name:               "customer_confirmations",
		Description:        "Retention policy for customer confirmation emails",
		RetentionDays:      730, // 2 years
		ArchiveDays:        365, // Archive after 1 year
		EmailTypes:         []string{"customer_confirmation"},
		Enabled:            true,
		ExecutionFrequency: 24 * time.Hour, // Daily
	}

	// Policy for consultant notifications (shorter retention)
	s.policies["consultant_notifications"] = &RetentionPolicy{
		Name:               "consultant_notifications",
		Description:        "Retention policy for internal consultant notifications",
		RetentionDays:      365, // 1 year
		ArchiveDays:        90,  // Archive after 3 months
		EmailTypes:         []string{"consultant_notification"},
		Enabled:            true,
		ExecutionFrequency: 24 * time.Hour, // Daily
	}

	// Policy for failed emails (keep shorter for troubleshooting)
	s.policies["failed_emails"] = &RetentionPolicy{
		Name:               "failed_emails",
		Description:        "Retention policy for failed email events",
		RetentionDays:      180, // 6 months
		ArchiveDays:        30,  // Archive after 1 month
		Statuses:           []string{"failed", "bounced", "spam"},
		Enabled:            true,
		ExecutionFrequency: 7 * 24 * time.Hour, // Weekly
	}

	// Policy for successful deliveries (standard retention)
	s.policies["successful_deliveries"] = &RetentionPolicy{
		Name:               "successful_deliveries",
		Description:        "Retention policy for successfully delivered emails",
		RetentionDays:      365, // 1 year
		ArchiveDays:        180, // Archive after 6 months
		Statuses:           []string{"delivered"},
		Enabled:            true,
		ExecutionFrequency: 24 * time.Hour, // Daily
	}

	// Set initial next execution times
	now := time.Now()
	for _, policy := range s.policies {
		nextExecution := now.Add(policy.ExecutionFrequency)
		policy.NextExecution = &nextExecution
	}
}

// AddPolicy adds a new retention policy
func (s *EmailRetentionService) AddPolicy(policy *RetentionPolicy) error {
	if policy.Name == "" {
		return fmt.Errorf("policy name cannot be empty")
	}

	if policy.RetentionDays <= 0 {
		return fmt.Errorf("retention days must be positive")
	}

	if policy.ArchiveDays < 0 {
		return fmt.Errorf("archive days cannot be negative")
	}

	if policy.ArchiveDays >= policy.RetentionDays {
		return fmt.Errorf("archive days must be less than retention days")
	}

	// Set default execution frequency if not specified
	if policy.ExecutionFrequency == 0 {
		policy.ExecutionFrequency = 24 * time.Hour
	}

	// Set next execution time
	if policy.NextExecution == nil {
		nextExecution := time.Now().Add(policy.ExecutionFrequency)
		policy.NextExecution = &nextExecution
	}

	s.policies[policy.Name] = policy

	s.logger.WithFields(logrus.Fields{
		"policy_name":    policy.Name,
		"retention_days": policy.RetentionDays,
		"archive_days":   policy.ArchiveDays,
		"execution_freq": policy.ExecutionFrequency,
	}).Info("Added email retention policy")

	return nil
}

// RemovePolicy removes a retention policy
func (s *EmailRetentionService) RemovePolicy(policyName string) error {
	if _, exists := s.policies[policyName]; !exists {
		return fmt.Errorf("policy %s not found", policyName)
	}

	delete(s.policies, policyName)

	s.logger.WithField("policy_name", policyName).Info("Removed email retention policy")
	return nil
}

// GetPolicies returns all retention policies
func (s *EmailRetentionService) GetPolicies() map[string]*RetentionPolicy {
	// Return a copy to prevent external modification
	policies := make(map[string]*RetentionPolicy)
	for name, policy := range s.policies {
		policyCopy := *policy
		policies[name] = &policyCopy
	}
	return policies
}

// ExecutePolicy executes a specific retention policy
func (s *EmailRetentionService) ExecutePolicy(ctx context.Context, policyName string) (*RetentionStats, error) {
	policy, exists := s.policies[policyName]
	if !exists {
		return nil, fmt.Errorf("policy %s not found", policyName)
	}

	if !policy.Enabled {
		return nil, fmt.Errorf("policy %s is disabled", policyName)
	}

	s.logger.WithField("policy_name", policyName).Info("Starting retention policy execution")

	stats := &RetentionStats{
		PolicyName:    policyName,
		ExecutionTime: time.Now(),
	}

	startTime := time.Now()

	// Execute archival first (if configured)
	if policy.ArchiveDays > 0 {
		archivedCount, err := s.archiveOldEvents(ctx, policy)
		if err != nil {
			stats.ErrorMessage = fmt.Sprintf("archival failed: %v", err)
			s.logger.WithError(err).WithField("policy_name", policyName).Error("Failed to archive old events")
			return stats, err
		}
		stats.RecordsArchived = archivedCount
	}

	// Execute deletion
	deletedCount, err := s.deleteOldEvents(ctx, policy)
	if err != nil {
		stats.ErrorMessage = fmt.Sprintf("deletion failed: %v", err)
		s.logger.WithError(err).WithField("policy_name", policyName).Error("Failed to delete old events")
		return stats, err
	}
	stats.RecordsDeleted = deletedCount

	stats.RecordsProcessed = stats.RecordsArchived + stats.RecordsDeleted
	stats.ExecutionDuration = time.Since(startTime)

	// Update policy execution time
	now := time.Now()
	policy.LastExecuted = &now
	nextExecution := now.Add(policy.ExecutionFrequency)
	policy.NextExecution = &nextExecution

	// Store stats
	s.stats = append(s.stats, stats)

	// Keep only last 100 stats entries
	if len(s.stats) > 100 {
		s.stats = s.stats[len(s.stats)-100:]
	}

	s.logger.WithFields(logrus.Fields{
		"policy_name":        policyName,
		"records_archived":   stats.RecordsArchived,
		"records_deleted":    stats.RecordsDeleted,
		"execution_duration": stats.ExecutionDuration,
	}).Info("Completed retention policy execution")

	return stats, nil
}

// ExecuteAllPolicies executes all enabled retention policies
func (s *EmailRetentionService) ExecuteAllPolicies(ctx context.Context) ([]*RetentionStats, error) {
	var allStats []*RetentionStats
	var errors []error

	for policyName, policy := range s.policies {
		if !policy.Enabled {
			continue
		}

		// Check if policy is due for execution
		if policy.NextExecution != nil && time.Now().Before(*policy.NextExecution) {
			continue
		}

		stats, err := s.ExecutePolicy(ctx, policyName)
		if err != nil {
			errors = append(errors, fmt.Errorf("policy %s failed: %w", policyName, err))
			continue
		}

		allStats = append(allStats, stats)
	}

	if len(errors) > 0 {
		return allStats, fmt.Errorf("some policies failed: %v", errors)
	}

	return allStats, nil
}

// archiveOldEvents moves old events to archive table
func (s *EmailRetentionService) archiveOldEvents(ctx context.Context, policy *RetentionPolicy) (int, error) {
	cutoffDate := time.Now().AddDate(0, 0, -policy.ArchiveDays)

	query := s.buildArchiveQuery(policy, cutoffDate)

	totalArchived := 0
	batchSize := s.maxBatchSize

	for {
		// Execute in batches to avoid long-running transactions
		result, err := s.db.ExecContext(ctx, query, batchSize)
		if err != nil {
			return totalArchived, fmt.Errorf("failed to archive events: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return totalArchived, fmt.Errorf("failed to get rows affected: %w", err)
		}

		totalArchived += int(rowsAffected)

		// If we processed fewer rows than batch size, we're done
		if int(rowsAffected) < batchSize {
			break
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return totalArchived, ctx.Err()
		default:
		}

		// Brief pause between batches
		time.Sleep(100 * time.Millisecond)
	}

	return totalArchived, nil
}

// deleteOldEvents removes old events from the main table
func (s *EmailRetentionService) deleteOldEvents(ctx context.Context, policy *RetentionPolicy) (int, error) {
	cutoffDate := time.Now().AddDate(0, 0, -policy.RetentionDays)

	query := s.buildDeleteQuery(policy, cutoffDate)

	totalDeleted := 0
	batchSize := s.maxBatchSize

	for {
		// Execute in batches
		result, err := s.db.ExecContext(ctx, query, batchSize)
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to delete events: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to get rows affected: %w", err)
		}

		totalDeleted += int(rowsAffected)

		// If we processed fewer rows than batch size, we're done
		if int(rowsAffected) < batchSize {
			break
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return totalDeleted, ctx.Err()
		default:
		}

		// Brief pause between batches
		time.Sleep(100 * time.Millisecond)
	}

	return totalDeleted, nil
}

// buildArchiveQuery builds the SQL query for archiving events
func (s *EmailRetentionService) buildArchiveQuery(policy *RetentionPolicy, cutoffDate time.Time) string {
	query := `
		WITH archived_events AS (
			DELETE FROM email_events 
			WHERE id IN (
				SELECT id FROM email_events 
				WHERE sent_at < $1`

	// Add email type filter
	if len(policy.EmailTypes) > 0 {
		query += " AND email_type IN ("
		for i, emailType := range policy.EmailTypes {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("'%s'", emailType)
		}
		query += ")"
	}

	// Add status filter
	if len(policy.Statuses) > 0 {
		query += " AND status IN ("
		for i, status := range policy.Statuses {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("'%s'", status)
		}
		query += ")"
	}

	query += `
				ORDER BY sent_at 
				LIMIT $2
			)
			RETURNING *
		)
		INSERT INTO email_events_archive 
		SELECT * FROM archived_events`

	return query
}

// buildDeleteQuery builds the SQL query for deleting events
func (s *EmailRetentionService) buildDeleteQuery(policy *RetentionPolicy, cutoffDate time.Time) string {
	query := `
		DELETE FROM email_events 
		WHERE id IN (
			SELECT id FROM email_events 
			WHERE sent_at < $1`

	// Add email type filter
	if len(policy.EmailTypes) > 0 {
		query += " AND email_type IN ("
		for i, emailType := range policy.EmailTypes {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("'%s'", emailType)
		}
		query += ")"
	}

	// Add status filter
	if len(policy.Statuses) > 0 {
		query += " AND status IN ("
		for i, status := range policy.Statuses {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("'%s'", status)
		}
		query += ")"
	}

	query += `
			ORDER BY sent_at 
			LIMIT $2
		)`

	return query
}

// GetStats returns retention execution statistics
func (s *EmailRetentionService) GetStats() []*RetentionStats {
	// Return a copy to prevent external modification
	stats := make([]*RetentionStats, len(s.stats))
	copy(stats, s.stats)
	return stats
}

// StartScheduler starts the background scheduler for retention policies
func (s *EmailRetentionService) StartScheduler(ctx context.Context) {
	if s.isRunning {
		s.logger.Warn("Retention scheduler is already running")
		return
	}

	s.isRunning = true
	s.logger.Info("Starting email retention scheduler")

	go s.schedulerLoop(ctx)
}

// StopScheduler stops the background scheduler
func (s *EmailRetentionService) StopScheduler() {
	if !s.isRunning {
		return
	}

	s.logger.Info("Stopping email retention scheduler")
	close(s.stopChan)
	s.isRunning = false
}

// schedulerLoop runs the background scheduler
func (s *EmailRetentionService) schedulerLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Retention scheduler stopped due to context cancellation")
			return
		case <-s.stopChan:
			s.logger.Info("Retention scheduler stopped")
			return
		case <-ticker.C:
			s.checkAndExecutePolicies(ctx)
		}
	}
}

// checkAndExecutePolicies checks if any policies are due for execution
func (s *EmailRetentionService) checkAndExecutePolicies(ctx context.Context) {
	now := time.Now()

	for policyName, policy := range s.policies {
		if !policy.Enabled {
			continue
		}

		if policy.NextExecution == nil || now.Before(*policy.NextExecution) {
			continue
		}

		s.logger.WithField("policy_name", policyName).Info("Executing scheduled retention policy")

		_, err := s.ExecutePolicy(ctx, policyName)
		if err != nil {
			s.logger.WithError(err).WithField("policy_name", policyName).Error("Scheduled retention policy execution failed")
		}
	}
}

// GetRetentionSummary returns a summary of retention status
func (s *EmailRetentionService) GetRetentionSummary(ctx context.Context) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Get total record counts
	var totalRecords, recentRecords, oldRecords int64

	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM email_events").Scan(&totalRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to get total records: %w", err)
	}

	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM email_events WHERE sent_at > $1",
		time.Now().AddDate(0, 0, -30)).Scan(&recentRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent records: %w", err)
	}

	err = s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM email_events WHERE sent_at < $1",
		time.Now().AddDate(0, 0, -365)).Scan(&oldRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to get old records: %w", err)
	}

	summary["total_records"] = totalRecords
	summary["recent_records"] = recentRecords
	summary["old_records"] = oldRecords
	summary["policies_count"] = len(s.policies)
	summary["enabled_policies"] = s.countEnabledPolicies()
	summary["scheduler_running"] = s.isRunning

	// Get archive table size if it exists
	var archiveRecords int64
	err = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM email_events_archive").Scan(&archiveRecords)
	if err == nil {
		summary["archived_records"] = archiveRecords
	}

	return summary, nil
}

// countEnabledPolicies counts the number of enabled policies
func (s *EmailRetentionService) countEnabledPolicies() int {
	count := 0
	for _, policy := range s.policies {
		if policy.Enabled {
			count++
		}
	}
	return count
}
