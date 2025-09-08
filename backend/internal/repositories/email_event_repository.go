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

// EmailEventRepositoryImpl implements the EmailEventRepository interface
type EmailEventRepositoryImpl struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewEmailEventRepository creates a new email event repository
func NewEmailEventRepository(db *sql.DB, logger *logrus.Logger) interfaces.EmailEventRepository {
	return &EmailEventRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create creates a new email event in the database
func (r *EmailEventRepositoryImpl) Create(ctx context.Context, event *domain.EmailEvent) error {
	query := `
		INSERT INTO email_events (id, inquiry_id, email_type, recipient_email, sender_email, 
		                         subject, status, sent_at, delivered_at, error_message, 
		                         bounce_type, ses_message_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

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

	_, err := r.db.ExecContext(ctx, query,
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

// Update updates an existing email event
func (r *EmailEventRepositoryImpl) Update(ctx context.Context, event *domain.EmailEvent) error {
	query := `
		UPDATE email_events 
		SET inquiry_id = $2, email_type = $3, recipient_email = $4, sender_email = $5,
		    subject = $6, status = $7, sent_at = $8, delivered_at = $9, 
		    error_message = $10, bounce_type = $11, ses_message_id = $12, updated_at = $13
		WHERE id = $1`

	// Update the updated_at timestamp
	event.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
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

// GetByInquiryID retrieves all email events for a specific inquiry
func (r *EmailEventRepositoryImpl) GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	query := `
		SELECT id, inquiry_id, email_type, recipient_email, sender_email, subject, 
		       status, sent_at, delivered_at, error_message, bounce_type, 
		       ses_message_id, created_at, updated_at
		FROM email_events
		WHERE inquiry_id = $1
		ORDER BY sent_at DESC`

	events, err := r.queryEvents(ctx, query, inquiryID)
	if err != nil {
		r.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to get email events by inquiry ID")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"inquiry_id":  inquiryID,
		"event_count": len(events),
	}).Debug("Retrieved email events by inquiry ID")

	return events, nil
}

// GetByMessageID retrieves an email event by SES message ID
func (r *EmailEventRepositoryImpl) GetByMessageID(ctx context.Context, messageID string) (*domain.EmailEvent, error) {
	query := `
		SELECT id, inquiry_id, email_type, recipient_email, sender_email, subject, 
		       status, sent_at, delivered_at, error_message, bounce_type, 
		       ses_message_id, created_at, updated_at
		FROM email_events
		WHERE ses_message_id = $1`

	row := r.db.QueryRowContext(ctx, query, messageID)

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

// GetMetrics calculates aggregated email metrics based on filters
func (r *EmailEventRepositoryImpl) GetMetrics(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error) {
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
		r.logger.WithError(err).Error("Failed to calculate email metrics")
		return nil, fmt.Errorf("failed to calculate email metrics: %w", err)
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
		"time_range":       metrics.TimeRange,
	}).Debug("Calculated email metrics")

	return metrics, nil
}

// List retrieves email events based on filters with pagination
func (r *EmailEventRepositoryImpl) List(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	query := `
		SELECT id, inquiry_id, email_type, recipient_email, sender_email, subject, 
		       status, sent_at, delivered_at, error_message, bounce_type, 
		       ses_message_id, created_at, updated_at
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

	// Add ORDER BY
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

	events, err := r.queryEvents(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to list email events")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"event_count": len(events),
		"limit":       filters.Limit,
		"offset":      filters.Offset,
	}).Debug("Listed email events")

	return events, nil
}

// queryEvents is a helper method to execute email event queries and scan results
func (r *EmailEventRepositoryImpl) queryEvents(ctx context.Context, query string, args ...interface{}) ([]*domain.EmailEvent, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query email events: %w", err)
	}
	defer rows.Close()

	var events []*domain.EmailEvent

	for rows.Next() {
		event := &domain.EmailEvent{}
		err := r.scanEmailEvent(rows, event)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan email event row")
			return nil, fmt.Errorf("failed to scan email event row: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error iterating over email event rows")
		return nil, fmt.Errorf("error iterating over email event rows: %w", err)
	}

	return events, nil
}

// scanEmailEvent is a helper method to scan email event data from database rows
func (r *EmailEventRepositoryImpl) scanEmailEvent(scanner interface{}, event *domain.EmailEvent) error {
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
