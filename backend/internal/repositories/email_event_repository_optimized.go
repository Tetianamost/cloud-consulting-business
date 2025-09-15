package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// OptimizedEmailEventRepository provides performance-optimized email event operations
type OptimizedEmailEventRepository struct {
	db     *sql.DB
	logger *logrus.Logger

	// Prepared statements for common queries
	getByInquiryStmt   *sql.Stmt
	getByMessageIDStmt *sql.Stmt
	getMetricsStmt     *sql.Stmt
	getMetricsFastStmt *sql.Stmt
	createEventStmt    *sql.Stmt
	updateEventStmt    *sql.Stmt
}

// NewOptimizedEmailEventRepository creates a new optimized email event repository
func NewOptimizedEmailEventRepository(db *sql.DB, logger *logrus.Logger) (interfaces.EmailEventRepository, error) {
	repo := &OptimizedEmailEventRepository{
		db:     db,
		logger: logger,
	}

	// Prepare common statements for better performance
	if err := repo.prepareStatements(); err != nil {
		return nil, fmt.Errorf("failed to prepare statements: %w", err)
	}

	return repo, nil
}

// prepareStatements prepares commonly used SQL statements
func (r *OptimizedEmailEventRepository) prepareStatements() error {
	var err error

	// Prepare statement for getting events by inquiry ID
	r.getByInquiryStmt, err = r.db.Prepare(`
		SELECT id, inquiry_id, email_type, recipient_email, sender_email, subject, 
		       status, sent_at, delivered_at, error_message, bounce_type, 
		       ses_message_id, created_at, updated_at
		FROM email_events
		WHERE inquiry_id = $1
		ORDER BY sent_at DESC`)
	if err != nil {
		return fmt.Errorf("failed to prepare getByInquiry statement: %w", err)
	}

	// Prepare statement for getting events by message ID
	r.getByMessageIDStmt, err = r.db.Prepare(`
		SELECT id, inquiry_id, email_type, recipient_email, sender_email, subject, 
		       status, sent_at, delivered_at, error_message, bounce_type, 
		       ses_message_id, created_at, updated_at
		FROM email_events
		WHERE ses_message_id = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare getByMessageID statement: %w", err)
	}

	// Prepare statement for basic metrics calculation
	r.getMetricsStmt, err = r.db.Prepare(`
		SELECT 
			COUNT(*) as total_emails,
			COUNT(*) FILTER (WHERE status = 'delivered') as delivered_emails,
			COUNT(*) FILTER (WHERE status = 'failed') as failed_emails,
			COUNT(*) FILTER (WHERE status = 'bounced') as bounced_emails,
			COUNT(*) FILTER (WHERE status = 'spam') as spam_emails
		FROM email_events
		WHERE sent_at >= $1 AND sent_at <= $2`)
	if err != nil {
		return fmt.Errorf("failed to prepare getMetrics statement: %w", err)
	}

	// Prepare statement for fast metrics using materialized views
	r.getMetricsFastStmt, err = r.db.Prepare(`
		SELECT * FROM get_email_metrics_fast($1, $2, $3)`)
	if err != nil {
		// If the function doesn't exist, we'll fall back to regular metrics
		r.logger.Warn("Fast metrics function not available, falling back to regular metrics")
		r.getMetricsFastStmt = nil
	}

	// Prepare statement for creating events
	r.createEventStmt, err = r.db.Prepare(`
		INSERT INTO email_events (id, inquiry_id, email_type, recipient_email, sender_email, 
		                         subject, status, sent_at, delivered_at, error_message, 
		                         bounce_type, ses_message_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`)
	if err != nil {
		return fmt.Errorf("failed to prepare createEvent statement: %w", err)
	}

	// Prepare statement for updating events
	r.updateEventStmt, err = r.db.Prepare(`
		UPDATE email_events 
		SET inquiry_id = $2, email_type = $3, recipient_email = $4, sender_email = $5,
		    subject = $6, status = $7, sent_at = $8, delivered_at = $9, 
		    error_message = $10, bounce_type = $11, ses_message_id = $12, updated_at = $13
		WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare updateEvent statement: %w", err)
	}

	return nil
}

// Create creates a new email event using prepared statement
func (r *OptimizedEmailEventRepository) Create(ctx context.Context, event *domain.EmailEvent) error {
	// Set timestamps if not already set
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	if event.UpdatedAt.IsZero() {
		event.UpdatedAt = time.Now()
	}
	if event.SentAt.IsZero() {
		event.SentAt = time.Now()
	}

	_, err := r.createEventStmt.ExecContext(ctx,
		event.ID,
		event.InquiryID,
		event.EmailType,
		event.RecipientEmail,
		event.SenderEmail,
		event.Subject,
		event.Status,
		event.SentAt,
		event.DeliveredAt,
		event.ErrorMessage,
		event.BounceType,
		event.SESMessageID,
		event.CreatedAt,
		event.UpdatedAt,
	)

	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"event_id":   event.ID,
			"inquiry_id": event.InquiryID,
			"email_type": event.EmailType,
			"status":     event.Status,
		}).Error("Failed to create email event")
		return fmt.Errorf("failed to create email event: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"inquiry_id": event.InquiryID,
		"email_type": event.EmailType,
		"status":     event.Status,
	}).Debug("Email event created successfully")

	return nil
}

// Update updates an existing email event using prepared statement
func (r *OptimizedEmailEventRepository) Update(ctx context.Context, event *domain.EmailEvent) error {
	// Update the updated_at timestamp
	event.UpdatedAt = time.Now()

	result, err := r.updateEventStmt.ExecContext(ctx,
		event.ID,
		event.InquiryID,
		event.EmailType,
		event.RecipientEmail,
		event.SenderEmail,
		event.Subject,
		event.Status,
		event.SentAt,
		event.DeliveredAt,
		event.ErrorMessage,
		event.BounceType,
		event.SESMessageID,
		event.UpdatedAt,
	)

	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"event_id":   event.ID,
			"inquiry_id": event.InquiryID,
			"status":     event.Status,
		}).Error("Failed to update email event")
		return fmt.Errorf("failed to update email event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email event not found")
	}

	r.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"inquiry_id": event.InquiryID,
		"status":     event.Status,
	}).Debug("Email event updated successfully")

	return nil
}

// GetByInquiryID retrieves all email events for a specific inquiry using prepared statement
func (r *OptimizedEmailEventRepository) GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	rows, err := r.getByInquiryStmt.QueryContext(ctx, inquiryID)
	if err != nil {
		r.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to query email events by inquiry ID")
		return nil, fmt.Errorf("failed to query email events by inquiry ID: %w", err)
	}
	defer rows.Close()

	events, err := r.scanEmailEvents(rows)
	if err != nil {
		r.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to scan email events")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"inquiry_id":  inquiryID,
		"event_count": len(events),
	}).Debug("Retrieved email events by inquiry ID")

	return events, nil
}

// GetByMessageID retrieves an email event by SES message ID using prepared statement
func (r *OptimizedEmailEventRepository) GetByMessageID(ctx context.Context, messageID string) (*domain.EmailEvent, error) {
	row := r.getByMessageIDStmt.QueryRowContext(ctx, messageID)

	event := &domain.EmailEvent{}
	err := r.scanEmailEvent(row, event)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WithField("message_id", messageID).Debug("Email event not found by message ID")
			return nil, nil
		}
		r.logger.WithError(err).WithField("message_id", messageID).Error("Failed to get email event by message ID")
		return nil, fmt.Errorf("failed to get email event by message ID: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"message_id": messageID,
		"status":     event.Status,
	}).Debug("Retrieved email event by message ID")

	return event, nil
}

// GetMetrics calculates aggregated email metrics with optimization
func (r *OptimizedEmailEventRepository) GetMetrics(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error) {
	// Try to use fast metrics function if available and filters are simple
	if r.getMetricsFastStmt != nil && r.canUseFastMetrics(filters) {
		return r.getMetricsFast(ctx, filters)
	}

	// Fall back to regular metrics calculation
	return r.getMetricsRegular(ctx, filters)
}

// canUseFastMetrics determines if we can use the optimized materialized view function
func (r *OptimizedEmailEventRepository) canUseFastMetrics(filters domain.EmailEventFilters) bool {
	// Fast metrics can be used if:
	// 1. We have a time range
	// 2. No specific inquiry ID filter
	// 3. No complex status filters
	return filters.TimeRange != nil &&
		filters.InquiryID == nil &&
		(filters.Status == nil || *filters.Status != domain.EmailStatusSent)
}

// getMetricsFast uses the optimized database function with materialized views
func (r *OptimizedEmailEventRepository) getMetricsFast(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error) {
	var emailTypeParam interface{}
	if filters.EmailType != nil {
		emailTypeParam = *filters.EmailType
	}

	row := r.getMetricsFastStmt.QueryRowContext(ctx,
		filters.TimeRange.Start,
		filters.TimeRange.End,
		emailTypeParam)

	metrics := &domain.EmailMetrics{}
	err := row.Scan(
		&metrics.TotalEmails,
		&metrics.DeliveredEmails,
		&metrics.FailedEmails,
		&metrics.BouncedEmails,
		&metrics.SpamEmails,
		&metrics.DeliveryRate,
		&metrics.BounceRate,
		&metrics.SpamRate,
		&metrics.TimeRange,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to calculate fast email metrics")
		return nil, fmt.Errorf("failed to calculate fast email metrics: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"total_emails":     metrics.TotalEmails,
		"delivered_emails": metrics.DeliveredEmails,
		"delivery_rate":    metrics.DeliveryRate,
		"method":           "fast_materialized_view",
	}).Debug("Calculated email metrics using fast method")

	return metrics, nil
}

// getMetricsRegular uses the regular metrics calculation for complex filters
func (r *OptimizedEmailEventRepository) getMetricsRegular(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error) {
	query := `
		SELECT 
			COUNT(*) as total_emails,
			COUNT(*) FILTER (WHERE status = 'delivered') as delivered_emails,
			COUNT(*) FILTER (WHERE status = 'failed') as failed_emails,
			COUNT(*) FILTER (WHERE status = 'bounced') as bounced_emails,
			COUNT(*) FILTER (WHERE status = 'spam') as spam_emails
		FROM email_events`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if filters.TimeRange != nil {
		conditions = append(conditions, fmt.Sprintf("sent_at >= $%d AND sent_at <= $%d", argIndex, argIndex+1))
		args = append(args, filters.TimeRange.Start, filters.TimeRange.End)
		argIndex += 2
	}

	if filters.EmailType != nil {
		conditions = append(conditions, fmt.Sprintf("email_type = $%d", argIndex))
		args = append(args, *filters.EmailType)
		argIndex++
	}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.InquiryID != nil {
		conditions = append(conditions, fmt.Sprintf("inquiry_id = $%d", argIndex))
		args = append(args, *filters.InquiryID)
		argIndex++
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	metrics := &domain.EmailMetrics{}
	err := row.Scan(
		&metrics.TotalEmails,
		&metrics.DeliveredEmails,
		&metrics.FailedEmails,
		&metrics.BouncedEmails,
		&metrics.SpamEmails,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to calculate regular email metrics")
		return nil, fmt.Errorf("failed to calculate regular email metrics: %w", err)
	}

	// Calculate rates
	if metrics.TotalEmails > 0 {
		metrics.DeliveryRate = float64(metrics.DeliveredEmails) / float64(metrics.TotalEmails) * 100
		metrics.BounceRate = float64(metrics.BouncedEmails) / float64(metrics.TotalEmails) * 100
		metrics.SpamRate = float64(metrics.SpamEmails) / float64(metrics.TotalEmails) * 100
	}

	// Set time range description
	if filters.TimeRange != nil {
		metrics.TimeRange = fmt.Sprintf("%s to %s",
			filters.TimeRange.Start.Format("2006-01-02"),
			filters.TimeRange.End.Format("2006-01-02"))
	} else {
		metrics.TimeRange = "All time"
	}

	r.logger.WithFields(logrus.Fields{
		"total_emails":     metrics.TotalEmails,
		"delivered_emails": metrics.DeliveredEmails,
		"delivery_rate":    metrics.DeliveryRate,
		"method":           "regular_calculation",
	}).Debug("Calculated email metrics using regular method")

	return metrics, nil
}

// List retrieves email events with optimized pagination
func (r *OptimizedEmailEventRepository) List(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	// Use covering index hint for better performance
	query := `
		SELECT /*+ INDEX(email_events idx_email_events_metrics_covering) */
		       id, inquiry_id, email_type, recipient_email, sender_email, subject, 
		       status, sent_at, delivered_at, error_message, bounce_type, 
		       ses_message_id, created_at, updated_at
		FROM email_events`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions with optimal order for index usage
	if filters.TimeRange != nil {
		conditions = append(conditions, fmt.Sprintf("sent_at >= $%d AND sent_at <= $%d", argIndex, argIndex+1))
		args = append(args, filters.TimeRange.Start, filters.TimeRange.End)
		argIndex += 2
	}

	if filters.EmailType != nil {
		conditions = append(conditions, fmt.Sprintf("email_type = $%d", argIndex))
		args = append(args, *filters.EmailType)
		argIndex++
	}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.InquiryID != nil {
		conditions = append(conditions, fmt.Sprintf("inquiry_id = $%d", argIndex))
		args = append(args, *filters.InquiryID)
		argIndex++
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ORDER BY with index-friendly ordering
	query += " ORDER BY sent_at DESC"

	// Add LIMIT and OFFSET
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to list email events")
		return nil, fmt.Errorf("failed to list email events: %w", err)
	}
	defer rows.Close()

	events, err := r.scanEmailEvents(rows)
	if err != nil {
		r.logger.WithError(err).Error("Failed to scan email events")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"event_count": len(events),
		"limit":       filters.Limit,
		"offset":      filters.Offset,
	}).Debug("Listed email events")

	return events, nil
}

// BatchCreate creates multiple email events in a single transaction for better performance
func (r *OptimizedEmailEventRepository) BatchCreate(ctx context.Context, events []*domain.EmailEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statement within transaction
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO email_events (id, inquiry_id, email_type, recipient_email, sender_email, 
		                         subject, status, sent_at, delivered_at, error_message, 
		                         bounce_type, ses_message_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`)
	if err != nil {
		return fmt.Errorf("failed to prepare batch insert statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now()
	for _, event := range events {
		// Set timestamps if not already set
		if event.CreatedAt.IsZero() {
			event.CreatedAt = now
		}
		if event.UpdatedAt.IsZero() {
			event.UpdatedAt = now
		}
		if event.SentAt.IsZero() {
			event.SentAt = now
		}

		_, err := stmt.ExecContext(ctx,
			event.ID,
			event.InquiryID,
			event.EmailType,
			event.RecipientEmail,
			event.SenderEmail,
			event.Subject,
			event.Status,
			event.SentAt,
			event.DeliveredAt,
			event.ErrorMessage,
			event.BounceType,
			event.SESMessageID,
			event.CreatedAt,
			event.UpdatedAt,
		)

		if err != nil {
			r.logger.WithError(err).WithField("event_id", event.ID).Error("Failed to insert event in batch")
			return fmt.Errorf("failed to insert event %s in batch: %w", event.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch insert: %w", err)
	}

	r.logger.WithField("event_count", len(events)).Info("Successfully batch created email events")
	return nil
}

// scanEmailEvents scans multiple email events from database rows
func (r *OptimizedEmailEventRepository) scanEmailEvents(rows *sql.Rows) ([]*domain.EmailEvent, error) {
	var events []*domain.EmailEvent

	for rows.Next() {
		event := &domain.EmailEvent{}
		err := r.scanEmailEvent(rows, event)
		if err != nil {
			return nil, fmt.Errorf("failed to scan email event row: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over email event rows: %w", err)
	}

	return events, nil
}

// scanEmailEvent scans a single email event from database row
func (r *OptimizedEmailEventRepository) scanEmailEvent(scanner interface{}, event *domain.EmailEvent) error {
	var deliveredAt sql.NullTime
	var errorMessage sql.NullString
	var bounceType sql.NullString
	var sesMessageID sql.NullString
	var subject sql.NullString

	var err error
	switch s := scanner.(type) {
	case *sql.Row:
		err = s.Scan(
			&event.ID,
			&event.InquiryID,
			&event.EmailType,
			&event.RecipientEmail,
			&event.SenderEmail,
			&subject,
			&event.Status,
			&event.SentAt,
			&deliveredAt,
			&errorMessage,
			&bounceType,
			&sesMessageID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
	case *sql.Rows:
		err = s.Scan(
			&event.ID,
			&event.InquiryID,
			&event.EmailType,
			&event.RecipientEmail,
			&event.SenderEmail,
			&subject,
			&event.Status,
			&event.SentAt,
			&deliveredAt,
			&errorMessage,
			&bounceType,
			&sesMessageID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
	default:
		return fmt.Errorf("unsupported scanner type")
	}

	if err != nil {
		return err
	}

	// Handle nullable fields
	if deliveredAt.Valid {
		event.DeliveredAt = &deliveredAt.Time
	}
	if errorMessage.Valid {
		event.ErrorMessage = errorMessage.String
	}
	if bounceType.Valid {
		event.BounceType = bounceType.String
	}
	if sesMessageID.Valid {
		event.SESMessageID = sesMessageID.String
	}
	if subject.Valid {
		event.Subject = subject.String
	}

	return nil
}

// Close closes prepared statements
func (r *OptimizedEmailEventRepository) Close() error {
	var errors []error

	if r.getByInquiryStmt != nil {
		if err := r.getByInquiryStmt.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if r.getByMessageIDStmt != nil {
		if err := r.getByMessageIDStmt.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if r.getMetricsStmt != nil {
		if err := r.getMetricsStmt.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if r.getMetricsFastStmt != nil {
		if err := r.getMetricsFastStmt.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if r.createEventStmt != nil {
		if err := r.createEventStmt.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if r.updateEventStmt != nil {
		if err := r.updateEventStmt.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close prepared statements: %v", errors)
	}

	return nil
}
