package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/sirupsen/logrus"
)

// QualityAssuranceService implements the quality assurance interface
type QualityAssuranceService struct {
	db     interfaces.DatabaseService
	cache  interfaces.CacheService
	logger *logrus.Logger

	// Quality thresholds and standards
	qualityThresholds *interfaces.QualityThresholds
	qualityStandards  *interfaces.QualityStandards

	// Metrics and tracking
	metricsCollector interfaces.MetricsService
}

// NewQualityAssuranceService creates a new quality assurance service
func NewQualityAssuranceService(
	db interfaces.DatabaseService,
	cache interfaces.CacheService,
	logger *logrus.Logger,
	metricsCollector interfaces.MetricsService,
) *QualityAssuranceService {
	service := &QualityAssuranceService{
		db:               db,
		cache:            cache,
		logger:           logger,
		metricsCollector: metricsCollector,
	}

	// Initialize default quality thresholds and standards
	service.initializeDefaults()

	return service
}

// initializeDefaults sets up default quality thresholds and standards
func (s *QualityAssuranceService) initializeDefaults() {
	s.qualityThresholds = &interfaces.QualityThresholds{
		MinAccuracy:           0.8,
		MinClientSatisfaction: 7,
		MinPeerReviewScore:    0.75,
		MaxResponseTime:       5000, // 5 seconds
		MinImplementationRate: 0.7,
		AlertThresholds: map[string]float64{
			"accuracy_drop":     0.1,
			"satisfaction_drop": 2.0,
			"review_score_drop": 0.15,
		},
		UpdatedAt: time.Now(),
	}

	s.qualityStandards = &interfaces.QualityStandards{
		MinOverallScore: 0.75,
		ComponentWeights: map[string]float64{
			"technical_accuracy": 0.25,
			"business_relevance": 0.20,
			"completeness":       0.15,
			"clarity":            0.15,
			"actionability":      0.15,
			"innovation":         0.10,
		},
	}
}

// TrackRecommendationAccuracy tracks the accuracy of AI recommendations
func (s *QualityAssuranceService) TrackRecommendationAccuracy(ctx context.Context, recommendation *interfaces.RecommendationTracking) error {
	s.logger.WithFields(logrus.Fields{
		"recommendation_id": recommendation.RecommendationID,
		"type":              recommendation.RecommendationType,
		"confidence":        recommendation.Confidence,
	}).Info("Tracking recommendation accuracy")

	// Store recommendation tracking data
	query := `
		INSERT INTO recommendation_tracking (
			id, recommendation_id, inquiry_id, consultant_id, recommendation_type,
			content, confidence, generated_at, model_version, prompt_version,
			context, tags, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	contextJSON, _ := json.Marshal(recommendation.Context)
	tagsJSON, _ := json.Marshal(recommendation.Tags)

	_, err := s.db.Exec(ctx, query,
		recommendation.ID,
		recommendation.RecommendationID,
		recommendation.InquiryID,
		recommendation.ConsultantID,
		recommendation.RecommendationType,
		recommendation.Content,
		recommendation.Confidence,
		recommendation.GeneratedAt,
		recommendation.ModelVersion,
		recommendation.PromptVersion,
		string(contextJSON),
		string(tagsJSON),
		recommendation.Status,
		recommendation.CreatedAt,
		recommendation.UpdatedAt,
	)

	if err != nil {
		s.logger.WithError(err).Error("Failed to track recommendation accuracy")
		return fmt.Errorf("failed to track recommendation accuracy: %w", err)
	}

	// Update metrics
	s.metricsCollector.IncrementCounter("recommendations_tracked", map[string]string{
		"type":       recommendation.RecommendationType,
		"consultant": recommendation.ConsultantID,
	})

	// Cache recent tracking data for quick access
	cacheKey := fmt.Sprintf("recommendation_tracking:%s", recommendation.RecommendationID)
	trackingData, _ := json.Marshal(recommendation)
	s.cache.Set(ctx, cacheKey, trackingData, 3600) // 1 hour

	return nil
}

// GetAccuracyMetrics retrieves accuracy metrics based on filters
func (s *QualityAssuranceService) GetAccuracyMetrics(ctx context.Context, filters *interfaces.AccuracyFilters) (*interfaces.QualityAccuracyMetrics, error) {
	s.logger.WithFields(logrus.Fields{
		"consultant_id": filters.ConsultantID,
		"type":          filters.RecommendationType,
	}).Info("Getting accuracy metrics")

	// Build query with filters
	query := `
		SELECT 
			rt.recommendation_type,
			rt.consultant_id,
			COUNT(*) as total_recommendations,
			COUNT(CASE WHEN ro.accuracy >= 0.8 THEN 1 END) as high_accuracy_count,
			COUNT(CASE WHEN ro.outcome_type = 'accepted' THEN 1 END) as accepted_count,
			COUNT(CASE WHEN ro.outcome_type = 'rejected' THEN 1 END) as rejected_count,
			AVG(rt.confidence) as avg_confidence,
			AVG(ro.accuracy) as avg_accuracy,
			DATE_TRUNC('day', rt.created_at) as date
		FROM recommendation_tracking rt
		LEFT JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filters.ConsultantID != "" {
		query += fmt.Sprintf(" AND rt.consultant_id = $%d", argIndex)
		args = append(args, filters.ConsultantID)
		argIndex++
	}

	if filters.RecommendationType != "" {
		query += fmt.Sprintf(" AND rt.recommendation_type = $%d", argIndex)
		args = append(args, filters.RecommendationType)
		argIndex++
	}

	if filters.TimeRange != nil {
		query += fmt.Sprintf(" AND rt.created_at >= $%d AND rt.created_at <= $%d", argIndex, argIndex+1)
		args = append(args, filters.TimeRange.StartDate, filters.TimeRange.EndDate)
		argIndex += 2
	}

	query += " GROUP BY rt.recommendation_type, rt.consultant_id, DATE_TRUNC('day', rt.created_at) ORDER BY date DESC"

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get accuracy metrics")
		return nil, fmt.Errorf("failed to get accuracy metrics: %w", err)
	}
	defer rows.Close()

	metrics := &interfaces.QualityAccuracyMetrics{
		AccuracyByType:       make(map[string]float64),
		AccuracyByConsultant: make(map[string]float64),
		QualityDistribution:  make(map[string]int64),
		TrendData:            []interfaces.AccuracyTrendPoint{},
		Metadata:             make(map[string]interface{}),
		GeneratedAt:          time.Now(),
	}

	var totalRecommendations, validatedRecommendations, acceptedRecommendations, rejectedRecommendations int64
	var totalConfidence, totalAccuracy float64
	var dataPointCount int

	for rows.Next() {
		var recType, consultantID string
		var total, highAccuracy, accepted, rejected int64
		var avgConfidence, avgAccuracy float64
		var date time.Time

		err := rows.Scan(&recType, &consultantID, &total, &highAccuracy, &accepted, &rejected, &avgConfidence, &avgAccuracy, &date)
		if err != nil {
			continue
		}

		totalRecommendations += total
		acceptedRecommendations += accepted
		rejectedRecommendations += rejected
		totalConfidence += avgConfidence
		totalAccuracy += avgAccuracy
		dataPointCount++

		// Aggregate by type
		if _, exists := metrics.AccuracyByType[recType]; !exists {
			metrics.AccuracyByType[recType] = 0
		}
		metrics.AccuracyByType[recType] += avgAccuracy

		// Aggregate by consultant
		if _, exists := metrics.AccuracyByConsultant[consultantID]; !exists {
			metrics.AccuracyByConsultant[consultantID] = 0
		}
		metrics.AccuracyByConsultant[consultantID] += avgAccuracy

		// Add trend data point
		metrics.TrendData = append(metrics.TrendData, interfaces.AccuracyTrendPoint{
			Date:     date,
			Accuracy: avgAccuracy,
			Count:    total,
		})
	}

	// Calculate overall metrics
	if dataPointCount > 0 {
		metrics.OverallAccuracy = totalAccuracy / float64(dataPointCount)
		metrics.AverageConfidence = totalConfidence / float64(dataPointCount)
	}

	metrics.TotalRecommendations = totalRecommendations
	metrics.ValidatedRecommendations = validatedRecommendations
	metrics.AcceptedRecommendations = acceptedRecommendations
	metrics.RejectedRecommendations = rejectedRecommendations

	// Cache metrics for quick access
	cacheKey := fmt.Sprintf("accuracy_metrics:%s:%s", filters.ConsultantID, filters.RecommendationType)
	metricsData, _ := json.Marshal(metrics)
	s.cache.Set(ctx, cacheKey, metricsData, 1800) // 30 minutes

	return metrics, nil
}

// UpdateRecommendationOutcome updates the outcome of a recommendation
func (s *QualityAssuranceService) UpdateRecommendationOutcome(ctx context.Context, recommendationID string, outcome *interfaces.RecommendationOutcome) error {
	s.logger.WithFields(logrus.Fields{
		"recommendation_id": recommendationID,
		"outcome_type":      outcome.OutcomeType,
		"accuracy":          outcome.Accuracy,
	}).Info("Updating recommendation outcome")

	query := `
		INSERT INTO recommendation_outcomes (
			recommendation_id, outcome_type, actual_result, client_feedback,
			consultant_notes, implementation_cost, time_to_implement,
			business_impact, technical_impact, lessons_learned, accuracy,
			effectiveness, metadata, recorded_at, recorded_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (recommendation_id) DO UPDATE SET
			outcome_type = EXCLUDED.outcome_type,
			actual_result = EXCLUDED.actual_result,
			client_feedback = EXCLUDED.client_feedback,
			consultant_notes = EXCLUDED.consultant_notes,
			implementation_cost = EXCLUDED.implementation_cost,
			time_to_implement = EXCLUDED.time_to_implement,
			business_impact = EXCLUDED.business_impact,
			technical_impact = EXCLUDED.technical_impact,
			lessons_learned = EXCLUDED.lessons_learned,
			accuracy = EXCLUDED.accuracy,
			effectiveness = EXCLUDED.effectiveness,
			metadata = EXCLUDED.metadata,
			recorded_at = EXCLUDED.recorded_at,
			recorded_by = EXCLUDED.recorded_by`

	lessonsJSON, _ := json.Marshal(outcome.LessonsLearned)
	metadataJSON, _ := json.Marshal(outcome.Metadata)

	_, err := s.db.Exec(ctx, query,
		outcome.RecommendationID,
		outcome.OutcomeType,
		outcome.ActualResult,
		outcome.ClientFeedback,
		outcome.ConsultantNotes,
		outcome.ImplementationCost,
		outcome.TimeToImplement,
		outcome.BusinessImpact,
		outcome.TechnicalImpact,
		string(lessonsJSON),
		outcome.Accuracy,
		outcome.Effectiveness,
		string(metadataJSON),
		outcome.RecordedAt,
		outcome.RecordedBy,
	)

	if err != nil {
		s.logger.WithError(err).Error("Failed to update recommendation outcome")
		return fmt.Errorf("failed to update recommendation outcome: %w", err)
	}

	// Update metrics
	s.metricsCollector.IncrementCounter("recommendation_outcomes_updated", map[string]string{
		"outcome_type": outcome.OutcomeType,
	})

	// Check if outcome triggers quality alerts
	s.checkQualityAlerts(ctx, outcome)

	return nil
}

// SubmitForPeerReview submits a recommendation for peer review
func (s *QualityAssuranceService) SubmitForPeerReview(ctx context.Context, review *interfaces.PeerReviewRequest) error {
	s.logger.WithFields(logrus.Fields{
		"recommendation_id": review.RecommendationID,
		"requested_by":      review.RequestedBy,
		"assigned_to":       review.AssignedTo,
		"priority":          review.Priority,
	}).Info("Submitting recommendation for peer review")

	query := `
		INSERT INTO peer_review_requests (
			id, recommendation_id, requested_by, assigned_to, priority,
			review_type, context, specific_areas, due_date, instructions,
			metadata, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	areasJSON, _ := json.Marshal(review.SpecificAreas)
	metadataJSON, _ := json.Marshal(review.Metadata)

	_, err := s.db.Exec(ctx, query,
		review.ID,
		review.RecommendationID,
		review.RequestedBy,
		review.AssignedTo,
		review.Priority,
		review.ReviewType,
		review.Context,
		string(areasJSON),
		review.DueDate,
		review.Instructions,
		string(metadataJSON),
		review.Status,
		review.CreatedAt,
		review.UpdatedAt,
	)

	if err != nil {
		s.logger.WithError(err).Error("Failed to submit peer review request")
		return fmt.Errorf("failed to submit peer review request: %w", err)
	}

	// Update metrics
	s.metricsCollector.IncrementCounter("peer_reviews_requested", map[string]string{
		"priority":    review.Priority,
		"review_type": review.ReviewType,
	})

	return nil
}

// GetPendingReviews retrieves pending reviews for a reviewer
func (s *QualityAssuranceService) GetPendingReviews(ctx context.Context, reviewerID string) ([]*interfaces.PeerReview, error) {
	s.logger.WithField("reviewer_id", reviewerID).Info("Getting pending reviews")

	query := `
		SELECT 
			pr.id, prr.id as request_id, prr.recommendation_id, pr.reviewer_id,
			pr.reviewer_name, prr.review_type, pr.status, pr.started_at,
			pr.completed_at, pr.time_spent, pr.metadata, pr.created_at, pr.updated_at
		FROM peer_reviews pr
		JOIN peer_review_requests prr ON pr.request_id = prr.id
		WHERE pr.reviewer_id = $1 AND pr.status IN ('pending', 'in_progress')
		ORDER BY prr.priority DESC, prr.due_date ASC`

	rows, err := s.db.Query(ctx, query, reviewerID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get pending reviews")
		return nil, fmt.Errorf("failed to get pending reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*interfaces.PeerReview
	for rows.Next() {
		review := &interfaces.PeerReview{}
		var metadataJSON string

		err := rows.Scan(
			&review.ID,
			&review.RequestID,
			&review.RecommendationID,
			&review.ReviewerID,
			&review.ReviewerName,
			&review.ReviewType,
			&review.Status,
			&review.StartedAt,
			&review.CompletedAt,
			&review.TimeSpent,
			&metadataJSON,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(metadataJSON), &review.Metadata)
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// SubmitPeerReview submits peer review feedback
func (s *QualityAssuranceService) SubmitPeerReview(ctx context.Context, reviewID string, feedback *interfaces.PeerReviewFeedback) error {
	s.logger.WithFields(logrus.Fields{
		"review_id":      reviewID,
		"overall_rating": feedback.OverallRating,
		"approved":       feedback.Approved,
	}).Info("Submitting peer review feedback")

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update peer review with feedback
	updateQuery := `
		UPDATE peer_reviews SET
			status = 'completed',
			completed_at = $2,
			feedback = $3,
			updated_at = $4
		WHERE id = $1`

	feedbackJSON, _ := json.Marshal(feedback)
	_, err = tx.Exec(ctx, updateQuery, reviewID, time.Now(), string(feedbackJSON), time.Now())
	if err != nil {
		s.logger.WithError(err).Error("Failed to update peer review")
		return fmt.Errorf("failed to update peer review: %w", err)
	}

	// Update recommendation status based on review outcome
	if feedback.Approved {
		statusQuery := `
			UPDATE recommendation_tracking SET
				status = CASE 
					WHEN $2 = 'full' THEN 'validated'
					WHEN $2 = 'conditional' THEN 'conditionally_approved'
					ELSE 'reviewed'
				END,
				updated_at = $3
			WHERE recommendation_id = (
				SELECT prr.recommendation_id 
				FROM peer_reviews pr 
				JOIN peer_review_requests prr ON pr.request_id = prr.id 
				WHERE pr.id = $1
			)`

		_, err = tx.Exec(ctx, statusQuery, reviewID, feedback.ApprovalLevel, time.Now())
		if err != nil {
			s.logger.WithError(err).Error("Failed to update recommendation status")
			return fmt.Errorf("failed to update recommendation status: %w", err)
		}
	} else {
		statusQuery := `
			UPDATE recommendation_tracking SET
				status = 'rejected',
				updated_at = $2
			WHERE recommendation_id = (
				SELECT prr.recommendation_id 
				FROM peer_reviews pr 
				JOIN peer_review_requests prr ON pr.request_id = prr.id 
				WHERE pr.id = $1
			)`

		_, err = tx.Exec(ctx, statusQuery, reviewID, time.Now())
		if err != nil {
			s.logger.WithError(err).Error("Failed to update recommendation status")
			return fmt.Errorf("failed to update recommendation status: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update metrics
	s.metricsCollector.IncrementCounter("peer_reviews_completed", map[string]string{
		"approved":       fmt.Sprintf("%t", feedback.Approved),
		"approval_level": feedback.ApprovalLevel,
	})

	s.metricsCollector.RecordHistogram("peer_review_rating", float64(feedback.OverallRating), map[string]string{
		"approved": fmt.Sprintf("%t", feedback.Approved),
	})

	return nil
}

// GetReviewHistory retrieves review history based on filters
func (s *QualityAssuranceService) GetReviewHistory(ctx context.Context, filters *interfaces.ReviewFilters) ([]*interfaces.PeerReview, error) {
	s.logger.WithFields(logrus.Fields{
		"reviewer_id": filters.ReviewerID,
		"status":      filters.Status,
	}).Info("Getting review history")

	query := `
		SELECT 
			pr.id, prr.id as request_id, prr.recommendation_id, pr.reviewer_id,
			pr.reviewer_name, prr.review_type, pr.status, pr.started_at,
			pr.completed_at, pr.time_spent, pr.feedback, pr.metadata,
			pr.created_at, pr.updated_at
		FROM peer_reviews pr
		JOIN peer_review_requests prr ON pr.request_id = prr.id
		WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filters.ReviewerID != "" {
		query += fmt.Sprintf(" AND pr.reviewer_id = $%d", argIndex)
		args = append(args, filters.ReviewerID)
		argIndex++
	}

	if filters.Status != "" {
		query += fmt.Sprintf(" AND pr.status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.ReviewType != "" {
		query += fmt.Sprintf(" AND prr.review_type = $%d", argIndex)
		args = append(args, filters.ReviewType)
		argIndex++
	}

	if filters.TimeRange != nil {
		query += fmt.Sprintf(" AND pr.created_at >= $%d AND pr.created_at <= $%d", argIndex, argIndex+1)
		args = append(args, filters.TimeRange.StartDate, filters.TimeRange.EndDate)
		argIndex += 2
	}

	query += " ORDER BY pr.created_at DESC"

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get review history")
		return nil, fmt.Errorf("failed to get review history: %w", err)
	}
	defer rows.Close()

	var reviews []*interfaces.PeerReview
	for rows.Next() {
		review := &interfaces.PeerReview{}
		var feedbackJSON, metadataJSON string

		err := rows.Scan(
			&review.ID,
			&review.RequestID,
			&review.RecommendationID,
			&review.ReviewerID,
			&review.ReviewerName,
			&review.ReviewType,
			&review.Status,
			&review.StartedAt,
			&review.CompletedAt,
			&review.TimeSpent,
			&feedbackJSON,
			&metadataJSON,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if feedbackJSON != "" {
			json.Unmarshal([]byte(feedbackJSON), &review.Feedback)
		}
		json.Unmarshal([]byte(metadataJSON), &review.Metadata)

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// RecordClientOutcome records client engagement outcomes
func (s *QualityAssuranceService) RecordClientOutcome(ctx context.Context, outcome *interfaces.ClientOutcome) error {
	s.logger.WithFields(logrus.Fields{
		"inquiry_id":          outcome.InquiryID,
		"client_name":         outcome.ClientName,
		"outcome_type":        outcome.OutcomeType,
		"client_satisfaction": outcome.ClientSatisfaction,
	}).Info("Recording client outcome")

	query := `
		INSERT INTO client_outcomes (
			id, inquiry_id, client_name, recommendation_ids, engagement_type,
			outcome_type, client_satisfaction, business_value, cost_savings,
			time_to_value, implementation_rate, client_testimonial,
			challenges_faced, success_factors, lessons_learned,
			follow_up_opportunities, reference_permission, case_study_permission,
			net_promoter_score, metadata, recorded_at, recorded_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`

	recommendationIDsJSON, _ := json.Marshal(outcome.RecommendationIDs)
	challengesJSON, _ := json.Marshal(outcome.ChallengesFaced)
	successFactorsJSON, _ := json.Marshal(outcome.SuccessFactors)
	lessonsJSON, _ := json.Marshal(outcome.LessonsLearned)
	followUpJSON, _ := json.Marshal(outcome.FollowUpOpportunities)
	metadataJSON, _ := json.Marshal(outcome.Metadata)

	_, err := s.db.Exec(ctx, query,
		outcome.ID,
		outcome.InquiryID,
		outcome.ClientName,
		string(recommendationIDsJSON),
		outcome.EngagementType,
		outcome.OutcomeType,
		outcome.ClientSatisfaction,
		outcome.BusinessValue,
		outcome.CostSavings,
		outcome.TimeToValue,
		outcome.ImplementationRate,
		outcome.ClientTestimonial,
		string(challengesJSON),
		string(successFactorsJSON),
		string(lessonsJSON),
		string(followUpJSON),
		outcome.ReferencePermission,
		outcome.CaseStudyPermission,
		outcome.NetPromoterScore,
		string(metadataJSON),
		outcome.RecordedAt,
		outcome.RecordedBy,
	)

	if err != nil {
		s.logger.WithError(err).Error("Failed to record client outcome")
		return fmt.Errorf("failed to record client outcome: %w", err)
	}

	// Update metrics
	s.metricsCollector.IncrementCounter("client_outcomes_recorded", map[string]string{
		"outcome_type":    outcome.OutcomeType,
		"engagement_type": outcome.EngagementType,
	})

	s.metricsCollector.RecordHistogram("client_satisfaction", float64(outcome.ClientSatisfaction), map[string]string{
		"outcome_type": outcome.OutcomeType,
	})

	s.metricsCollector.RecordHistogram("business_value", outcome.BusinessValue, map[string]string{
		"outcome_type": outcome.OutcomeType,
	})

	return nil
}

// GetOutcomeAnalytics retrieves outcome analytics based on filters
func (s *QualityAssuranceService) GetOutcomeAnalytics(ctx context.Context, filters *interfaces.OutcomeFilters) (*interfaces.OutcomeAnalytics, error) {
	s.logger.WithFields(logrus.Fields{
		"client_name":     filters.ClientName,
		"engagement_type": filters.EngagementType,
	}).Info("Getting outcome analytics")

	// Build base query
	query := `
		SELECT 
			COUNT(*) as total_engagements,
			COUNT(CASE WHEN outcome_type = 'successful' THEN 1 END) as successful_engagements,
			AVG(client_satisfaction) as avg_satisfaction,
			SUM(business_value) as total_business_value,
			SUM(cost_savings) as total_cost_savings,
			AVG(time_to_value) as avg_time_to_value,
			AVG(implementation_rate) as avg_implementation_rate,
			AVG(net_promoter_score) as avg_nps,
			outcome_type,
			DATE_TRUNC('month', recorded_at) as month
		FROM client_outcomes
		WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filters.ClientName != "" {
		query += fmt.Sprintf(" AND client_name = $%d", argIndex)
		args = append(args, filters.ClientName)
		argIndex++
	}

	if filters.EngagementType != "" {
		query += fmt.Sprintf(" AND engagement_type = $%d", argIndex)
		args = append(args, filters.EngagementType)
		argIndex++
	}

	if filters.TimeRange != nil {
		query += fmt.Sprintf(" AND recorded_at >= $%d AND recorded_at <= $%d", argIndex, argIndex+1)
		args = append(args, filters.TimeRange.StartDate, filters.TimeRange.EndDate)
		argIndex += 2
	}

	query += " GROUP BY outcome_type, DATE_TRUNC('month', recorded_at) ORDER BY month DESC"

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get outcome analytics")
		return nil, fmt.Errorf("failed to get outcome analytics: %w", err)
	}
	defer rows.Close()

	analytics := &interfaces.OutcomeAnalytics{
		OutcomesByType:           make(map[string]int64),
		SatisfactionDistribution: make(map[string]int64),
		BusinessValueTrend:       []interfaces.BusinessValueTrendPoint{},
		TopSuccessFactors:        []interfaces.SuccessFactor{},
		CommonChallenges:         []interfaces.Challenge{},
		Metadata:                 make(map[string]interface{}),
		GeneratedAt:              time.Now(),
	}

	var totalEngagements, successfulEngagements int64
	var totalBusinessValue, totalCostSavings, avgTimeToValue, avgImplementationRate, avgNPS float64

	for rows.Next() {
		var total, successful int64
		var avgSatisfaction, businessValue, costSavings, timeToValue, implementationRate, nps float64
		var outcomeType string
		var month time.Time

		err := rows.Scan(&total, &successful, &avgSatisfaction, &businessValue, &costSavings,
			&timeToValue, &implementationRate, &nps, &outcomeType, &month)
		if err != nil {
			continue
		}

		totalEngagements += total
		successfulEngagements += successful
		totalBusinessValue += businessValue
		totalCostSavings += costSavings
		avgTimeToValue += timeToValue
		avgImplementationRate += implementationRate
		avgNPS += nps

		analytics.OutcomesByType[outcomeType] = total

		// Add business value trend point
		analytics.BusinessValueTrend = append(analytics.BusinessValueTrend, interfaces.BusinessValueTrendPoint{
			Date:          month,
			BusinessValue: businessValue,
			CostSavings:   costSavings,
			Engagements:   total,
		})
	}

	// Calculate overall metrics
	analytics.TotalEngagements = totalEngagements
	analytics.SuccessfulEngagements = successfulEngagements
	if totalEngagements > 0 {
		analytics.SuccessRate = float64(successfulEngagements) / float64(totalEngagements)
	}
	analytics.TotalBusinessValue = totalBusinessValue
	analytics.TotalCostSavings = totalCostSavings
	analytics.AverageTimeToValue = avgTimeToValue
	analytics.AverageImplementationRate = avgImplementationRate
	analytics.AverageNetPromoterScore = avgNPS

	// Get success factors and challenges
	analytics.TopSuccessFactors = s.getTopSuccessFactors(ctx, filters)
	analytics.CommonChallenges = s.getCommonChallenges(ctx, filters)

	return analytics, nil
}

// ValidateRecommendationEffectiveness validates the effectiveness of a recommendation
func (s *QualityAssuranceService) ValidateRecommendationEffectiveness(ctx context.Context, recommendationID string) (*interfaces.EffectivenessReport, error) {
	s.logger.WithField("recommendation_id", recommendationID).Info("Validating recommendation effectiveness")

	// Get recommendation tracking data
	trackingQuery := `
		SELECT rt.*, ro.accuracy, ro.effectiveness, ro.client_feedback,
			   ro.business_impact, ro.implementation_cost
		FROM recommendation_tracking rt
		LEFT JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		WHERE rt.recommendation_id = $1`

	var tracking interfaces.RecommendationTracking
	var outcome interfaces.RecommendationOutcome

	row := s.db.QueryRow(ctx, trackingQuery, recommendationID)
	var contextJSON, tagsJSON string

	err := row.Scan(
		&tracking.ID, &tracking.RecommendationID, &tracking.InquiryID,
		&tracking.ConsultantID, &tracking.RecommendationType, &tracking.Content,
		&tracking.Confidence, &tracking.GeneratedAt, &tracking.ModelVersion,
		&tracking.PromptVersion, &contextJSON, &tagsJSON, &tracking.Status,
		&tracking.CreatedAt, &tracking.UpdatedAt,
		&outcome.Accuracy, &outcome.Effectiveness, &outcome.ClientFeedback,
		&outcome.BusinessImpact, &outcome.ImplementationCost,
	)

	if err != nil {
		s.logger.WithError(err).Error("Failed to get recommendation data")
		return nil, fmt.Errorf("failed to get recommendation data: %w", err)
	}

	// Get peer review scores
	peerReviewQuery := `
		SELECT AVG(
			(COALESCE((feedback->>'overall_rating')::int, 0) +
			 COALESCE((feedback->>'technical_accuracy')::int, 0) +
			 COALESCE((feedback->>'business_relevance')::int, 0) +
			 COALESCE((feedback->>'completeness')::int, 0) +
			 COALESCE((feedback->>'clarity')::int, 0) +
			 COALESCE((feedback->>'actionability')::int, 0)) / 6.0
		) as avg_peer_score
		FROM peer_reviews pr
		JOIN peer_review_requests prr ON pr.request_id = prr.id
		WHERE prr.recommendation_id = $1 AND pr.status = 'completed'`

	var peerReviewScore float64
	err = s.db.QueryRow(ctx, peerReviewQuery, recommendationID).Scan(&peerReviewScore)
	if err != nil {
		peerReviewScore = 0 // No peer reviews available
	}

	// Calculate effectiveness scores
	report := &interfaces.EffectivenessReport{
		RecommendationID:      recommendationID,
		TechnicalAccuracy:     outcome.Accuracy,
		BusinessRelevance:     s.calculateBusinessRelevance(outcome.BusinessImpact),
		ImplementationSuccess: outcome.Effectiveness,
		PeerReviewScore:       peerReviewScore / 5.0, // Convert from 1-5 to 0-1 scale
		ClientOutcomeScore:    s.calculateClientOutcomeScore(ctx, tracking.InquiryID),
		GeneratedAt:           time.Now(),
	}

	// Calculate overall effectiveness
	weights := map[string]float64{
		"technical_accuracy":     0.25,
		"business_relevance":     0.20,
		"implementation_success": 0.20,
		"peer_review_score":      0.20,
		"client_outcome_score":   0.15,
	}

	report.OverallEffectiveness =
		report.TechnicalAccuracy*weights["technical_accuracy"] +
			report.BusinessRelevance*weights["business_relevance"] +
			report.ImplementationSuccess*weights["implementation_success"] +
			report.PeerReviewScore*weights["peer_review_score"] +
			report.ClientOutcomeScore*weights["client_outcome_score"]

	// Generate insights
	report.Strengths = s.identifyStrengths(report)
	report.Weaknesses = s.identifyWeaknesses(report)
	report.ImprovementAreas = s.identifyEffectivenessImprovementAreas(report)
	report.Recommendations = s.generateEffectivenessRecommendations(report)

	return report, nil
}

// GenerateImprovementInsights generates insights for continuous improvement
func (s *QualityAssuranceService) GenerateImprovementInsights(ctx context.Context, timeRange *interfaces.QualityTimeRange) (*interfaces.ImprovementInsights, error) {
	s.logger.WithFields(logrus.Fields{
		"start_date": timeRange.StartDate,
		"end_date":   timeRange.EndDate,
	}).Info("Generating improvement insights")

	insights := &interfaces.ImprovementInsights{
		TimeRange:                     timeRange,
		KeyImprovementAreas:           []interfaces.ImprovementArea{},
		SuccessPatterns:               []interfaces.Pattern{},
		FailurePatterns:               []interfaces.Pattern{},
		ModelPerformanceInsights:      []interfaces.ModelInsight{},
		ConsultantPerformanceInsights: []interfaces.ConsultantInsight{},
		RecommendedActions:            []interfaces.RecommendedAction{},
		QualityMetricsTrend:           []interfaces.QualityMetricTrendPoint{},
		GeneratedAt:                   time.Now(),
	}

	// Analyze quality trends
	trendQuery := `
		SELECT 
			DATE_TRUNC('week', rt.created_at) as week,
			AVG(ro.accuracy) as avg_accuracy,
			AVG(ro.effectiveness) as avg_effectiveness,
			AVG(co.client_satisfaction) as avg_satisfaction,
			AVG(co.implementation_rate) as avg_implementation_rate
		FROM recommendation_tracking rt
		LEFT JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		LEFT JOIN client_outcomes co ON rt.inquiry_id = co.inquiry_id
		WHERE rt.created_at >= $1 AND rt.created_at <= $2
		GROUP BY DATE_TRUNC('week', rt.created_at)
		ORDER BY week`

	rows, err := s.db.Query(ctx, trendQuery, timeRange.StartDate, timeRange.EndDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to query quality trends")
		return nil, fmt.Errorf("failed to query quality trends: %w", err)
	}
	defer rows.Close()

	var trendPoints []interfaces.QualityMetricTrendPoint
	var overallTrend string = "stable"
	var trendCount int

	for rows.Next() {
		var week time.Time
		var avgAccuracy, avgEffectiveness, avgSatisfaction, avgImplementationRate float64

		err := rows.Scan(&week, &avgAccuracy, &avgEffectiveness, &avgSatisfaction, &avgImplementationRate)
		if err != nil {
			continue
		}

		trendPoint := interfaces.QualityMetricTrendPoint{
			Date:               week,
			OverallQuality:     (avgAccuracy + avgEffectiveness + (avgSatisfaction / 10) + avgImplementationRate) / 4,
			TechnicalAccuracy:  avgAccuracy,
			BusinessRelevance:  avgEffectiveness,
			ClientSatisfaction: avgSatisfaction,
			ImplementationRate: avgImplementationRate,
		}

		trendPoints = append(trendPoints, trendPoint)
		trendCount++
	}

	insights.QualityMetricsTrend = trendPoints

	// Determine overall trend direction
	if len(trendPoints) >= 2 {
		firstHalf := trendPoints[:len(trendPoints)/2]
		secondHalf := trendPoints[len(trendPoints)/2:]

		var firstAvg, secondAvg float64
		for _, point := range firstHalf {
			firstAvg += point.OverallQuality
		}
		firstAvg /= float64(len(firstHalf))

		for _, point := range secondHalf {
			secondAvg += point.OverallQuality
		}
		secondAvg /= float64(len(secondHalf))

		if secondAvg > firstAvg+0.05 {
			overallTrend = "improving"
		} else if secondAvg < firstAvg-0.05 {
			overallTrend = "declining"
		}
	}

	insights.OverallQualityTrend = overallTrend

	// Generate improvement areas
	insights.KeyImprovementAreas = s.identifyImprovementAreas(ctx, timeRange, trendPoints)

	// Generate success and failure patterns
	insights.SuccessPatterns = s.identifySuccessPatterns(ctx, timeRange)
	insights.FailurePatterns = s.identifyFailurePatterns(ctx, timeRange)

	// Generate model performance insights
	insights.ModelPerformanceInsights = s.analyzeModelPerformance(ctx, timeRange)

	// Generate consultant performance insights
	insights.ConsultantPerformanceInsights = s.analyzeConsultantPerformance(ctx, timeRange)

	// Generate recommended actions
	insights.RecommendedActions = s.generateRecommendedActions(insights)

	// Generate benchmark comparison
	insights.BenchmarkComparison = s.generateBenchmarkComparison(trendPoints)

	return insights, nil
}

// identifyImprovementAreas identifies key areas for improvement
func (s *QualityAssuranceService) identifyImprovementAreas(ctx context.Context, timeRange *interfaces.QualityTimeRange, trendPoints []interfaces.QualityMetricTrendPoint) []interfaces.ImprovementArea {
	var areas []interfaces.ImprovementArea

	if len(trendPoints) == 0 {
		return areas
	}

	// Calculate average scores
	var avgAccuracy, avgRelevance, avgSatisfaction, avgImplementation float64
	for _, point := range trendPoints {
		avgAccuracy += point.TechnicalAccuracy
		avgRelevance += point.BusinessRelevance
		avgSatisfaction += point.ClientSatisfaction
		avgImplementation += point.ImplementationRate
	}

	count := float64(len(trendPoints))
	avgAccuracy /= count
	avgRelevance /= count
	avgSatisfaction /= count
	avgImplementation /= count

	// Identify areas below thresholds
	if avgAccuracy < 0.8 {
		areas = append(areas, interfaces.ImprovementArea{
			Area:         "Technical Accuracy",
			CurrentScore: avgAccuracy,
			TargetScore:  0.85,
			Priority:     "high",
			Impact:       "high",
			Effort:       "medium",
			Description:  "Technical accuracy of recommendations needs improvement",
			Recommendations: []string{
				"Enhance AI model training with more technical examples",
				"Implement additional technical validation checks",
				"Increase peer review focus on technical aspects",
			},
		})
	}

	if avgRelevance < 0.75 {
		areas = append(areas, interfaces.ImprovementArea{
			Area:         "Business Relevance",
			CurrentScore: avgRelevance,
			TargetScore:  0.8,
			Priority:     "medium",
			Impact:       "high",
			Effort:       "medium",
			Description:  "Business relevance of recommendations could be improved",
			Recommendations: []string{
				"Better understand client business context",
				"Include more business impact analysis",
				"Train consultants on business value articulation",
			},
		})
	}

	if avgSatisfaction < 7.5 {
		areas = append(areas, interfaces.ImprovementArea{
			Area:         "Client Satisfaction",
			CurrentScore: avgSatisfaction / 10, // Convert to 0-1 scale
			TargetScore:  0.85,
			Priority:     "high",
			Impact:       "high",
			Effort:       "high",
			Description:  "Client satisfaction scores need improvement",
			Recommendations: []string{
				"Improve communication and follow-up processes",
				"Enhance recommendation presentation quality",
				"Implement client feedback collection system",
			},
		})
	}

	if avgImplementation < 0.7 {
		areas = append(areas, interfaces.ImprovementArea{
			Area:         "Implementation Success",
			CurrentScore: avgImplementation,
			TargetScore:  0.8,
			Priority:     "medium",
			Impact:       "medium",
			Effort:       "high",
			Description:  "Implementation success rate needs improvement",
			Recommendations: []string{
				"Provide more detailed implementation guidance",
				"Offer implementation support services",
				"Create implementation templates and tools",
			},
		})
	}

	return areas
}

// identifySuccessPatterns identifies patterns in successful recommendations
func (s *QualityAssuranceService) identifySuccessPatterns(ctx context.Context, timeRange *interfaces.QualityTimeRange) []interfaces.Pattern {
	query := `
		SELECT 
			rt.recommendation_type,
			rt.consultant_id,
			COUNT(*) as count,
			AVG(ro.accuracy) as avg_accuracy,
			AVG(co.client_satisfaction) as avg_satisfaction
		FROM recommendation_tracking rt
		JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		JOIN client_outcomes co ON rt.inquiry_id = co.inquiry_id
		WHERE rt.created_at >= $1 AND rt.created_at <= $2
		AND ro.accuracy >= 0.8 AND co.client_satisfaction >= 8
		GROUP BY rt.recommendation_type, rt.consultant_id
		HAVING COUNT(*) >= 3
		ORDER BY avg_accuracy DESC, avg_satisfaction DESC`

	rows, err := s.db.Query(ctx, query, timeRange.StartDate, timeRange.EndDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to identify success patterns")
		return []interfaces.Pattern{}
	}
	defer rows.Close()

	var patterns []interfaces.Pattern
	for rows.Next() {
		var recType, consultantID string
		var count int64
		var avgAccuracy, avgSatisfaction float64

		err := rows.Scan(&recType, &consultantID, &count, &avgAccuracy, &avgSatisfaction)
		if err != nil {
			continue
		}

		pattern := interfaces.Pattern{
			Name:        fmt.Sprintf("High-Quality %s Recommendations", strings.Title(recType)),
			Description: fmt.Sprintf("Consultant %s consistently delivers high-quality %s recommendations", consultantID, recType),
			Frequency:   count,
			Confidence:  (avgAccuracy + avgSatisfaction/10) / 2,
			Conditions: []string{
				fmt.Sprintf("Recommendation type: %s", recType),
				fmt.Sprintf("Consultant: %s", consultantID),
				fmt.Sprintf("Average accuracy: %.2f", avgAccuracy),
				fmt.Sprintf("Average satisfaction: %.1f", avgSatisfaction),
			},
			Outcomes: []string{
				"High client satisfaction",
				"Successful implementation",
				"Positive business impact",
			},
			Examples: []string{
				fmt.Sprintf("%d successful %s recommendations", count, recType),
			},
		}

		patterns = append(patterns, pattern)
	}

	return patterns
}

// identifyFailurePatterns identifies patterns in failed recommendations
func (s *QualityAssuranceService) identifyFailurePatterns(ctx context.Context, timeRange *interfaces.QualityTimeRange) []interfaces.Pattern {
	query := `
		SELECT 
			rt.recommendation_type,
			COUNT(*) as count,
			AVG(ro.accuracy) as avg_accuracy,
			AVG(co.client_satisfaction) as avg_satisfaction,
			STRING_AGG(DISTINCT ro.consultant_notes, '; ') as common_issues
		FROM recommendation_tracking rt
		JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		JOIN client_outcomes co ON rt.inquiry_id = co.inquiry_id
		WHERE rt.created_at >= $1 AND rt.created_at <= $2
		AND (ro.accuracy < 0.6 OR co.client_satisfaction < 6)
		GROUP BY rt.recommendation_type
		HAVING COUNT(*) >= 2
		ORDER BY count DESC`

	rows, err := s.db.Query(ctx, query, timeRange.StartDate, timeRange.EndDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to identify failure patterns")
		return []interfaces.Pattern{}
	}
	defer rows.Close()

	var patterns []interfaces.Pattern
	for rows.Next() {
		var recType, commonIssues string
		var count int64
		var avgAccuracy, avgSatisfaction float64

		err := rows.Scan(&recType, &count, &avgAccuracy, &avgSatisfaction, &commonIssues)
		if err != nil {
			continue
		}

		pattern := interfaces.Pattern{
			Name:        fmt.Sprintf("Low-Quality %s Recommendations", strings.Title(recType)),
			Description: fmt.Sprintf("Pattern of poor performance in %s recommendations", recType),
			Frequency:   count,
			Confidence:  1.0 - (avgAccuracy+avgSatisfaction/10)/2, // Inverse confidence for failure patterns
			Conditions: []string{
				fmt.Sprintf("Recommendation type: %s", recType),
				fmt.Sprintf("Low accuracy: %.2f", avgAccuracy),
				fmt.Sprintf("Low satisfaction: %.1f", avgSatisfaction),
			},
			Outcomes: []string{
				"Poor client satisfaction",
				"Implementation challenges",
				"Negative business impact",
			},
			Examples: []string{
				fmt.Sprintf("%d failed %s recommendations", count, recType),
				fmt.Sprintf("Common issues: %s", commonIssues),
			},
		}

		patterns = append(patterns, pattern)
	}

	return patterns
}

// analyzeModelPerformance analyzes AI model performance
func (s *QualityAssuranceService) analyzeModelPerformance(ctx context.Context, timeRange *interfaces.QualityTimeRange) []interfaces.ModelInsight {
	query := `
		SELECT 
			rt.model_version,
			COUNT(*) as total_recommendations,
			AVG(rt.confidence) as avg_confidence,
			AVG(ro.accuracy) as avg_accuracy,
			AVG(ro.effectiveness) as avg_effectiveness
		FROM recommendation_tracking rt
		LEFT JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		WHERE rt.created_at >= $1 AND rt.created_at <= $2
		GROUP BY rt.model_version
		ORDER BY avg_accuracy DESC`

	rows, err := s.db.Query(ctx, query, timeRange.StartDate, timeRange.EndDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to analyze model performance")
		return []interfaces.ModelInsight{}
	}
	defer rows.Close()

	var insights []interfaces.ModelInsight
	for rows.Next() {
		var modelVersion string
		var totalRecs int64
		var avgConfidence, avgAccuracy, avgEffectiveness float64

		err := rows.Scan(&modelVersion, &totalRecs, &avgConfidence, &avgAccuracy, &avgEffectiveness)
		if err != nil {
			continue
		}

		performanceScore := (avgAccuracy + avgEffectiveness) / 2

		insight := interfaces.ModelInsight{
			ModelVersion:     modelVersion,
			PerformanceScore: performanceScore,
			UsageStats: map[string]interface{}{
				"total_recommendations": totalRecs,
				"average_confidence":    avgConfidence,
				"average_accuracy":      avgAccuracy,
				"average_effectiveness": avgEffectiveness,
			},
		}

		// Generate strengths and weaknesses based on performance
		if performanceScore >= 0.8 {
			insight.Strengths = []string{
				"High overall performance",
				"Consistent accuracy",
				"Good effectiveness ratings",
			}
		} else {
			insight.Weaknesses = []string{
				"Below target performance",
				"Inconsistent results",
				"Needs improvement",
			}
		}

		// Generate recommendations
		if avgAccuracy < 0.75 {
			insight.Recommendations = append(insight.Recommendations, "Improve training data quality")
		}
		if avgConfidence > avgAccuracy+0.1 {
			insight.Recommendations = append(insight.Recommendations, "Calibrate confidence scoring")
		}
		if totalRecs < 10 {
			insight.Recommendations = append(insight.Recommendations, "Increase model usage for better insights")
		}

		insights = append(insights, insight)
	}

	return insights
}

// analyzeConsultantPerformance analyzes consultant performance
func (s *QualityAssuranceService) analyzeConsultantPerformance(ctx context.Context, timeRange *interfaces.QualityTimeRange) []interfaces.ConsultantInsight {
	query := `
		SELECT 
			rt.consultant_id,
			COUNT(*) as total_recommendations,
			AVG(ro.accuracy) as avg_accuracy,
			AVG(co.client_satisfaction) as avg_satisfaction,
			AVG(co.implementation_rate) as avg_implementation_rate,
			COUNT(CASE WHEN pr.feedback->>'approved' = 'true' THEN 1 END) as peer_approvals,
			COUNT(pr.id) as total_peer_reviews
		FROM recommendation_tracking rt
		LEFT JOIN recommendation_outcomes ro ON rt.recommendation_id = ro.recommendation_id
		LEFT JOIN client_outcomes co ON rt.inquiry_id = co.inquiry_id
		LEFT JOIN peer_review_requests prr ON rt.recommendation_id = prr.recommendation_id
		LEFT JOIN peer_reviews pr ON prr.id = pr.request_id
		WHERE rt.created_at >= $1 AND rt.created_at <= $2
		GROUP BY rt.consultant_id
		HAVING COUNT(*) >= 3
		ORDER BY avg_accuracy DESC, avg_satisfaction DESC`

	rows, err := s.db.Query(ctx, query, timeRange.StartDate, timeRange.EndDate)
	if err != nil {
		s.logger.WithError(err).Error("Failed to analyze consultant performance")
		return []interfaces.ConsultantInsight{}
	}
	defer rows.Close()

	var insights []interfaces.ConsultantInsight
	for rows.Next() {
		var consultantID string
		var totalRecs, peerApprovals, totalPeerReviews int64
		var avgAccuracy, avgSatisfaction, avgImplementationRate float64

		err := rows.Scan(&consultantID, &totalRecs, &avgAccuracy, &avgSatisfaction, &avgImplementationRate, &peerApprovals, &totalPeerReviews)
		if err != nil {
			continue
		}

		performanceScore := (avgAccuracy + avgSatisfaction/10 + avgImplementationRate) / 3

		insight := interfaces.ConsultantInsight{
			ConsultantID:     consultantID,
			ConsultantName:   fmt.Sprintf("Consultant %s", consultantID), // In real implementation, would lookup name
			PerformanceScore: performanceScore,
		}

		// Generate strengths and improvement areas
		if avgAccuracy >= 0.8 {
			insight.Strengths = append(insight.Strengths, "High technical accuracy")
		} else {
			insight.ImprovementAreas = append(insight.ImprovementAreas, "Technical accuracy needs improvement")
		}

		if avgSatisfaction >= 8.0 {
			insight.Strengths = append(insight.Strengths, "Excellent client satisfaction")
		} else {
			insight.ImprovementAreas = append(insight.ImprovementAreas, "Client communication and satisfaction")
		}

		if avgImplementationRate >= 0.8 {
			insight.Strengths = append(insight.Strengths, "High implementation success rate")
		} else {
			insight.ImprovementAreas = append(insight.ImprovementAreas, "Implementation guidance and support")
		}

		// Generate training needs
		if avgAccuracy < 0.75 {
			insight.TrainingNeeds = append(insight.TrainingNeeds, "Technical skills enhancement")
		}
		if avgSatisfaction < 7.0 {
			insight.TrainingNeeds = append(insight.TrainingNeeds, "Client communication skills")
		}
		if totalPeerReviews > 0 && float64(peerApprovals)/float64(totalPeerReviews) < 0.8 {
			insight.TrainingNeeds = append(insight.TrainingNeeds, "Peer review and quality standards")
		}

		// Generate best practices
		if performanceScore >= 0.8 {
			insight.BestPractices = []string{
				"Consistent high-quality deliverables",
				"Strong client relationship management",
				"Effective implementation guidance",
			}
		}

		insights = append(insights, insight)
	}

	return insights
}

// generateRecommendedActions generates recommended improvement actions
func (s *QualityAssuranceService) generateRecommendedActions(insights *interfaces.ImprovementInsights) []interfaces.RecommendedAction {
	var actions []interfaces.RecommendedAction

	// Generate actions based on improvement areas
	for _, area := range insights.KeyImprovementAreas {
		action := interfaces.RecommendedAction{
			Action:      fmt.Sprintf("Improve %s", area.Area),
			Priority:    area.Priority,
			Impact:      area.Impact,
			Effort:      area.Effort,
			Timeline:    s.getTimelineForPriority(area.Priority),
			Owner:       "Quality Assurance Team",
			Description: area.Description,
			Steps:       area.Recommendations,
		}
		actions = append(actions, action)
	}

	// Generate actions based on failure patterns
	for _, pattern := range insights.FailurePatterns {
		if pattern.Frequency >= 3 {
			action := interfaces.RecommendedAction{
				Action:      fmt.Sprintf("Address %s", pattern.Name),
				Priority:    "high",
				Impact:      "high",
				Effort:      "medium",
				Timeline:    "2-4 weeks",
				Owner:       "Quality Assurance Team",
				Description: fmt.Sprintf("Address recurring pattern: %s", pattern.Description),
				Steps: []string{
					"Analyze root causes of the pattern",
					"Develop targeted interventions",
					"Implement process improvements",
					"Monitor for pattern recurrence",
				},
			}
			actions = append(actions, action)
		}
	}

	// Generate actions based on model performance
	for _, modelInsight := range insights.ModelPerformanceInsights {
		if modelInsight.PerformanceScore < 0.75 {
			action := interfaces.RecommendedAction{
				Action:      fmt.Sprintf("Improve %s Performance", modelInsight.ModelVersion),
				Priority:    "medium",
				Impact:      "medium",
				Effort:      "high",
				Timeline:    "4-8 weeks",
				Owner:       "AI/ML Team",
				Description: fmt.Sprintf("Improve performance of %s model", modelInsight.ModelVersion),
				Steps:       modelInsight.Recommendations,
			}
			actions = append(actions, action)
		}
	}

	return actions
}

// generateBenchmarkComparison generates benchmark comparison
func (s *QualityAssuranceService) generateBenchmarkComparison(trendPoints []interfaces.QualityMetricTrendPoint) *interfaces.QualityBenchmarkComparison {
	if len(trendPoints) == 0 {
		return &interfaces.QualityBenchmarkComparison{
			IndustryBenchmark:  0.75,
			CompanyBenchmark:   0.80,
			CurrentPerformance: 0.0,
			PerformanceGap:     -0.75,
			RankingPercentile:  0.0,
			ComparisonInsights: []string{"Insufficient data for comparison"},
		}
	}

	// Calculate current performance as average of recent trend points
	var currentPerformance float64
	recentPoints := trendPoints
	if len(trendPoints) > 4 {
		recentPoints = trendPoints[len(trendPoints)-4:] // Last 4 weeks
	}

	for _, point := range recentPoints {
		currentPerformance += point.OverallQuality
	}
	currentPerformance /= float64(len(recentPoints))

	industryBenchmark := 0.75 // Industry average
	companyBenchmark := 0.80  // Company target
	performanceGap := currentPerformance - companyBenchmark

	// Calculate ranking percentile (simplified)
	var rankingPercentile float64
	if currentPerformance >= 0.9 {
		rankingPercentile = 95.0
	} else if currentPerformance >= 0.85 {
		rankingPercentile = 80.0
	} else if currentPerformance >= 0.8 {
		rankingPercentile = 65.0
	} else if currentPerformance >= 0.75 {
		rankingPercentile = 50.0
	} else {
		rankingPercentile = 25.0
	}

	var insights []string
	if currentPerformance > companyBenchmark {
		insights = append(insights, "Performance exceeds company benchmark")
	} else {
		insights = append(insights, "Performance below company benchmark")
	}

	if currentPerformance > industryBenchmark {
		insights = append(insights, "Performance above industry average")
	} else {
		insights = append(insights, "Performance below industry average")
	}

	return &interfaces.QualityBenchmarkComparison{
		IndustryBenchmark:  industryBenchmark,
		CompanyBenchmark:   companyBenchmark,
		CurrentPerformance: currentPerformance,
		PerformanceGap:     performanceGap,
		RankingPercentile:  rankingPercentile,
		ComparisonInsights: insights,
	}
}

// getTimelineForPriority returns appropriate timeline based on priority
func (s *QualityAssuranceService) getTimelineForPriority(priority string) string {
	switch priority {
	case "high":
		return "1-2 weeks"
	case "medium":
		return "2-4 weeks"
	case "low":
		return "4-8 weeks"
	default:
		return "2-4 weeks"
	}
}

// ValidateRecommendationQuality validates the quality of an AI recommendation
func (s *QualityAssuranceService) ValidateRecommendationQuality(ctx context.Context, recommendation *interfaces.AIRecommendation) (*interfaces.QualityValidation, error) {
	s.logger.WithFields(logrus.Fields{
		"recommendation_id": recommendation.ID,
		"type":              recommendation.Type,
		"confidence":        recommendation.Confidence,
	}).Info("Validating recommendation quality")

	validation := &interfaces.QualityValidation{
		RecommendationID: recommendation.ID,
		ValidationChecks: []interfaces.ValidationCheck{},
		Issues:           []interfaces.QualityIssue{},
		Recommendations:  []string{},
		ValidatedAt:      time.Now(),
		ValidatedBy:      "quality_assurance_service",
	}

	// Perform validation checks
	checks := []struct {
		name        string
		description string
		weight      float64
		checkFunc   func(*interfaces.AIRecommendation) (bool, float64, string, []string)
	}{
		{
			"content_completeness",
			"Check if recommendation content is complete and comprehensive",
			0.20,
			s.checkContentCompleteness,
		},
		{
			"technical_accuracy",
			"Validate technical accuracy of the recommendation",
			0.25,
			s.checkTechnicalAccuracy,
		},
		{
			"business_relevance",
			"Assess business relevance and value proposition",
			0.20,
			s.checkBusinessRelevance,
		},
		{
			"clarity_readability",
			"Evaluate clarity and readability of the content",
			0.15,
			s.checkClarityReadability,
		},
		{
			"actionability",
			"Check if recommendations are actionable and specific",
			0.20,
			s.checkActionability,
		},
	}

	var totalScore, totalWeight float64
	passedChecks := 0

	for _, check := range checks {
		passed, score, details, suggestions := check.checkFunc(recommendation)

		validationCheck := interfaces.ValidationCheck{
			Name:        check.name,
			Description: check.description,
			Passed:      passed,
			Score:       score,
			Weight:      check.weight,
			Details:     details,
			Suggestions: suggestions,
		}

		validation.ValidationChecks = append(validation.ValidationChecks, validationCheck)

		if passed {
			passedChecks++
		}

		totalScore += score * check.weight
		totalWeight += check.weight

		// Add issues for failed checks
		if !passed {
			severity := "medium"
			if score < 0.3 {
				severity = "high"
			} else if score < 0.6 {
				severity = "medium"
			} else {
				severity = "low"
			}

			issue := interfaces.QualityIssue{
				Type:        check.name,
				Severity:    severity,
				Description: fmt.Sprintf("Quality check failed: %s", check.description),
				Location:    "recommendation_content",
				Suggestions: suggestions,
				AutoFixable: false,
			}
			validation.Issues = append(validation.Issues, issue)
		}
	}

	validation.PassedChecks = passedChecks
	validation.TotalChecks = len(checks)
	validation.PassRate = float64(passedChecks) / float64(len(checks))
	validation.OverallScore = totalScore / totalWeight

	// Determine if review is required
	validation.RequiresReview = validation.OverallScore < s.qualityStandards.MinOverallScore ||
		validation.PassRate < 0.8 ||
		len(validation.Issues) > 2

	// Generate improvement recommendations
	validation.Recommendations = s.generateQualityRecommendations(validation)

	// Update metrics
	s.metricsCollector.RecordHistogram("quality_validation_score", validation.OverallScore, map[string]string{
		"type": recommendation.Type,
	})

	s.metricsCollector.RecordHistogram("quality_validation_pass_rate", validation.PassRate, map[string]string{
		"type": recommendation.Type,
	})

	return validation, nil
}

// GetQualityScore retrieves the quality score for a recommendation
func (s *QualityAssuranceService) GetQualityScore(ctx context.Context, recommendationID string) (*interfaces.QualityScore, error) {
	s.logger.WithField("recommendation_id", recommendationID).Info("Getting quality score")

	// Check cache first
	cacheKey := fmt.Sprintf("quality_score:%s", recommendationID)
	if cachedData, err := s.cache.Get(ctx, cacheKey); err == nil {
		var score interfaces.QualityScore
		if json.Unmarshal(cachedData, &score) == nil {
			return &score, nil
		}
	}

	// Calculate quality score from various sources
	score := &interfaces.QualityScore{
		RecommendationID: recommendationID,
		ComponentScores:  make(map[string]float64),
		Metadata:         make(map[string]interface{}),
		CalculatedAt:     time.Now(),
	}

	// Get validation score
	validationQuery := `
		SELECT overall_score, validation_checks, pass_rate
		FROM quality_validations
		WHERE recommendation_id = $1
		ORDER BY validated_at DESC
		LIMIT 1`

	var validationScore, passRate float64
	var validationChecksJSON string
	err := s.db.QueryRow(ctx, validationQuery, recommendationID).Scan(&validationScore, &validationChecksJSON, &passRate)
	if err == nil {
		score.ComponentScores["validation"] = validationScore
	}

	// Get peer review score
	peerReviewQuery := `
		SELECT AVG(
			(COALESCE((feedback->>'overall_rating')::int, 0) +
			 COALESCE((feedback->>'technical_accuracy')::int, 0) +
			 COALESCE((feedback->>'business_relevance')::int, 0) +
			 COALESCE((feedback->>'completeness')::int, 0) +
			 COALESCE((feedback->>'clarity')::int, 0) +
			 COALESCE((feedback->>'actionability')::int, 0)) / 6.0 / 5.0
		) as avg_peer_score
		FROM peer_reviews pr
		JOIN peer_review_requests prr ON pr.request_id = prr.id
		WHERE prr.recommendation_id = $1 AND pr.status = 'completed'`

	var peerReviewScore float64
	err = s.db.QueryRow(ctx, peerReviewQuery, recommendationID).Scan(&peerReviewScore)
	if err == nil {
		score.ComponentScores["peer_review"] = peerReviewScore
	}

	// Get outcome score
	outcomeQuery := `
		SELECT accuracy, effectiveness
		FROM recommendation_outcomes
		WHERE recommendation_id = $1`

	var accuracy, effectiveness float64
	err = s.db.QueryRow(ctx, outcomeQuery, recommendationID).Scan(&accuracy, &effectiveness)
	if err == nil {
		score.ComponentScores["accuracy"] = accuracy
		score.ComponentScores["effectiveness"] = effectiveness
	}

	// Calculate weighted overall score
	weights := s.qualityStandards.ComponentWeights
	var totalScore, totalWeight float64

	for component, componentScore := range score.ComponentScores {
		if weight, exists := weights[component]; exists {
			totalScore += componentScore * weight
			totalWeight += weight
		}
	}

	if totalWeight > 0 {
		score.OverallScore = totalScore / totalWeight
		score.WeightedScore = totalScore
	}

	// Determine quality grade
	score.QualityGrade = s.calculateQualityGrade(score.OverallScore)

	// Calculate score breakdown
	score.ScoreBreakdown = &interfaces.ScoreBreakdown{
		TechnicalAccuracy: score.ComponentScores["accuracy"],
		BusinessRelevance: score.ComponentScores["business_relevance"],
		Completeness:      score.ComponentScores["completeness"],
		Clarity:           score.ComponentScores["clarity"],
		Actionability:     score.ComponentScores["actionability"],
	}

	// Cache the score
	scoreData, _ := json.Marshal(score)
	s.cache.Set(ctx, cacheKey, scoreData, 1800) // 30 minutes

	return score, nil
}

// Helper methods for quality validation checks

func (s *QualityAssuranceService) checkContentCompleteness(recommendation *interfaces.AIRecommendation) (bool, float64, string, []string) {
	content := strings.TrimSpace(recommendation.Content)

	// Check minimum content length
	if len(content) < 100 {
		return false, 0.2, "Content is too short and lacks detail", []string{
			"Expand the recommendation with more detailed explanations",
			"Include specific implementation steps",
			"Add context and background information",
		}
	}

	// Check for key sections
	score := 0.0
	suggestions := []string{}

	if strings.Contains(strings.ToLower(content), "recommendation") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Include clear recommendation statements")
	}

	if strings.Contains(strings.ToLower(content), "implementation") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Add implementation guidance")
	}

	if strings.Contains(strings.ToLower(content), "benefit") || strings.Contains(strings.ToLower(content), "advantage") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Explain benefits and advantages")
	}

	if strings.Contains(strings.ToLower(content), "cost") || strings.Contains(strings.ToLower(content), "price") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Include cost considerations")
	}

	if strings.Contains(strings.ToLower(content), "risk") || strings.Contains(strings.ToLower(content), "consideration") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Address potential risks and considerations")
	}

	passed := score >= 0.8
	details := fmt.Sprintf("Content completeness score: %.2f", score)

	return passed, score, details, suggestions
}

func (s *QualityAssuranceService) checkTechnicalAccuracy(recommendation *interfaces.AIRecommendation) (bool, float64, string, []string) {
	content := strings.ToLower(recommendation.Content)
	score := 0.8 // Default score, would be enhanced with actual technical validation
	suggestions := []string{}

	// Check for technical terms and specificity
	technicalTerms := []string{"aws", "azure", "gcp", "cloud", "service", "architecture", "configuration"}
	termCount := 0
	for _, term := range technicalTerms {
		if strings.Contains(content, term) {
			termCount++
		}
	}

	if termCount < 3 {
		score -= 0.2
		suggestions = append(suggestions, "Include more specific technical details")
	}

	// Check for specific service names or configurations
	if !strings.Contains(content, "ec2") && !strings.Contains(content, "s3") &&
		!strings.Contains(content, "lambda") && !strings.Contains(content, "rds") {
		score -= 0.1
		suggestions = append(suggestions, "Reference specific cloud services where applicable")
	}

	passed := score >= 0.7
	details := fmt.Sprintf("Technical accuracy assessment score: %.2f", score)

	return passed, score, details, suggestions
}

func (s *QualityAssuranceService) checkBusinessRelevance(recommendation *interfaces.AIRecommendation) (bool, float64, string, []string) {
	content := strings.ToLower(recommendation.Content)
	score := 0.0
	suggestions := []string{}

	// Check for business value indicators
	businessTerms := []string{"roi", "cost saving", "efficiency", "productivity", "revenue", "business value"}
	for _, term := range businessTerms {
		if strings.Contains(content, term) {
			score += 0.2
		}
	}

	// Check for business context
	if strings.Contains(content, "business") || strings.Contains(content, "organization") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Include business context and impact")
	}

	// Check for stakeholder considerations
	if strings.Contains(content, "stakeholder") || strings.Contains(content, "team") || strings.Contains(content, "user") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Consider stakeholder impact")
	}

	if score < 0.4 {
		suggestions = append(suggestions, "Strengthen business value proposition")
	}

	passed := score >= 0.6
	details := fmt.Sprintf("Business relevance score: %.2f", score)

	return passed, score, details, suggestions
}

func (s *QualityAssuranceService) checkClarityReadability(recommendation *interfaces.AIRecommendation) (bool, float64, string, []string) {
	content := recommendation.Content
	score := 0.8 // Base score
	suggestions := []string{}

	// Check sentence length (simple heuristic)
	sentences := strings.Split(content, ".")
	longSentences := 0
	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		if len(words) > 25 {
			longSentences++
		}
	}

	if longSentences > len(sentences)/3 {
		score -= 0.2
		suggestions = append(suggestions, "Break down long sentences for better readability")
	}

	// Check for structure indicators
	if strings.Contains(content, "1.") || strings.Contains(content, "•") || strings.Contains(content, "-") {
		score += 0.1
	} else {
		suggestions = append(suggestions, "Use bullet points or numbered lists for better structure")
	}

	passed := score >= 0.7
	details := fmt.Sprintf("Clarity and readability score: %.2f", score)

	return passed, score, details, suggestions
}

func (s *QualityAssuranceService) checkActionability(recommendation *interfaces.AIRecommendation) (bool, float64, string, []string) {
	content := strings.ToLower(recommendation.Content)
	score := 0.0
	suggestions := []string{}

	// Check for action verbs
	actionVerbs := []string{"implement", "configure", "deploy", "setup", "create", "establish", "migrate"}
	verbCount := 0
	for _, verb := range actionVerbs {
		if strings.Contains(content, verb) {
			verbCount++
			score += 0.15
		}
	}

	if verbCount == 0 {
		suggestions = append(suggestions, "Include specific action items and implementation steps")
	}

	// Check for step-by-step guidance
	if strings.Contains(content, "step") || strings.Contains(content, "first") || strings.Contains(content, "next") {
		score += 0.2
	} else {
		suggestions = append(suggestions, "Provide step-by-step implementation guidance")
	}

	// Check for specific timelines or priorities
	if strings.Contains(content, "immediate") || strings.Contains(content, "priority") || strings.Contains(content, "phase") {
		score += 0.15
	} else {
		suggestions = append(suggestions, "Include timeline and priority information")
	}

	passed := score >= 0.6
	details := fmt.Sprintf("Actionability score: %.2f", score)

	return passed, score, details, suggestions
}

// Helper methods for analytics and insights

func (s *QualityAssuranceService) getTopSuccessFactors(ctx context.Context, filters *interfaces.OutcomeFilters) []interfaces.SuccessFactor {
	// This would analyze success factors from client outcomes
	// For now, returning mock data
	return []interfaces.SuccessFactor{
		{Factor: "Clear communication", Frequency: 15, Impact: 0.8},
		{Factor: "Stakeholder engagement", Frequency: 12, Impact: 0.75},
		{Factor: "Proper planning", Frequency: 10, Impact: 0.7},
	}
}

func (s *QualityAssuranceService) getCommonChallenges(ctx context.Context, filters *interfaces.OutcomeFilters) []interfaces.Challenge {
	// This would analyze challenges from client outcomes
	// For now, returning mock data
	return []interfaces.Challenge{
		{Challenge: "Resource constraints", Frequency: 8, Impact: 0.6},
		{Challenge: "Technical complexity", Frequency: 6, Impact: 0.7},
		{Challenge: "Timeline pressure", Frequency: 5, Impact: 0.5},
	}
}

func (s *QualityAssuranceService) calculateBusinessRelevance(businessImpact string) float64 {
	// Simple heuristic based on business impact description
	impact := strings.ToLower(businessImpact)
	if strings.Contains(impact, "significant") || strings.Contains(impact, "major") {
		return 0.9
	} else if strings.Contains(impact, "moderate") || strings.Contains(impact, "good") {
		return 0.7
	} else if strings.Contains(impact, "minor") || strings.Contains(impact, "small") {
		return 0.5
	}
	return 0.6 // Default
}

func (s *QualityAssuranceService) calculateClientOutcomeScore(ctx context.Context, inquiryID string) float64 {
	query := `
		SELECT AVG(client_satisfaction) / 10.0 as normalized_satisfaction
		FROM client_outcomes
		WHERE inquiry_id = $1`

	var score float64
	err := s.db.QueryRow(ctx, query, inquiryID).Scan(&score)
	if err != nil {
		return 0.5 // Default score if no data
	}
	return score
}

func (s *QualityAssuranceService) calculateQualityGrade(score float64) string {
	if score >= 0.9 {
		return "A"
	} else if score >= 0.8 {
		return "B"
	} else if score >= 0.7 {
		return "C"
	} else if score >= 0.6 {
		return "D"
	}
	return "F"
}

// Additional helper methods

func (s *QualityAssuranceService) identifyStrengths(report *interfaces.EffectivenessReport) []string {
	strengths := []string{}

	if report.TechnicalAccuracy >= 0.8 {
		strengths = append(strengths, "High technical accuracy")
	}
	if report.BusinessRelevance >= 0.8 {
		strengths = append(strengths, "Strong business relevance")
	}
	if report.ImplementationSuccess >= 0.8 {
		strengths = append(strengths, "Successful implementation")
	}
	if report.PeerReviewScore >= 0.8 {
		strengths = append(strengths, "Positive peer feedback")
	}
	if report.ClientOutcomeScore >= 0.8 {
		strengths = append(strengths, "Excellent client outcomes")
	}

	return strengths
}

func (s *QualityAssuranceService) identifyWeaknesses(report *interfaces.EffectivenessReport) []string {
	weaknesses := []string{}

	if report.TechnicalAccuracy < 0.6 {
		weaknesses = append(weaknesses, "Low technical accuracy")
	}
	if report.BusinessRelevance < 0.6 {
		weaknesses = append(weaknesses, "Limited business relevance")
	}
	if report.ImplementationSuccess < 0.6 {
		weaknesses = append(weaknesses, "Implementation challenges")
	}
	if report.PeerReviewScore < 0.6 {
		weaknesses = append(weaknesses, "Peer review concerns")
	}
	if report.ClientOutcomeScore < 0.6 {
		weaknesses = append(weaknesses, "Poor client outcomes")
	}

	return weaknesses
}

func (s *QualityAssuranceService) generateEffectivenessRecommendations(report *interfaces.EffectivenessReport) []string {
	recommendations := []string{}

	if report.TechnicalAccuracy < 0.7 {
		recommendations = append(recommendations, "Enhance technical validation processes")
	}
	if report.BusinessRelevance < 0.7 {
		recommendations = append(recommendations, "Strengthen business context in recommendations")
	}
	if report.ImplementationSuccess < 0.7 {
		recommendations = append(recommendations, "Improve implementation guidance and support")
	}
	if report.PeerReviewScore < 0.7 {
		recommendations = append(recommendations, "Increase peer review participation and quality")
	}
	if report.ClientOutcomeScore < 0.7 {
		recommendations = append(recommendations, "Focus on client satisfaction and outcome tracking")
	}

	return recommendations
}

func (s *QualityAssuranceService) calculateAverageQuality(points []interfaces.QualityMetricTrendPoint) float64 {
	if len(points) == 0 {
		return 0
	}

	total := 0.0
	for _, point := range points {
		total += point.OverallQuality
	}

	return total / float64(len(points))
}

// checkQualityAlerts checks if outcome triggers any quality alerts
func (s *QualityAssuranceService) checkQualityAlerts(ctx context.Context, outcome *interfaces.RecommendationOutcome) {
	// Check if outcome triggers any quality alerts
	if outcome.Accuracy < s.qualityThresholds.MinAccuracy {
		alert := &interfaces.QualityAlert{
			ID:          fmt.Sprintf("accuracy_alert_%d", time.Now().Unix()),
			AlertType:   "threshold",
			Severity:    "high",
			MetricType:  "accuracy",
			MetricValue: outcome.Accuracy,
			Threshold:   s.qualityThresholds.MinAccuracy,
			Description: fmt.Sprintf("Recommendation accuracy (%.2f) below threshold (%.2f)", outcome.Accuracy, s.qualityThresholds.MinAccuracy),
			Status:      "active",
			CreatedAt:   time.Now(),
		}

		s.TriggerQualityAlert(ctx, alert)
	}
}

// generateQualityRecommendations generates quality improvement recommendations
func (s *QualityAssuranceService) generateQualityRecommendations(validation *interfaces.QualityValidation) []string {
	recommendations := []string{}

	if validation.OverallScore < 0.6 {
		recommendations = append(recommendations, "Comprehensive review and revision required")
	} else if validation.OverallScore < 0.8 {
		recommendations = append(recommendations, "Minor improvements needed before approval")
	}

	if validation.PassRate < 0.8 {
		recommendations = append(recommendations, "Address failed validation checks")
	}

	for _, issue := range validation.Issues {
		if issue.Severity == "high" {
			recommendations = append(recommendations, fmt.Sprintf("High priority: %s", issue.Description))
		}
	}

	return recommendations
}

// identifyEffectivenessImprovementAreas identifies improvement areas for effectiveness reports
func (s *QualityAssuranceService) identifyEffectivenessImprovementAreas(report *interfaces.EffectivenessReport) []string {
	areas := []string{}

	scores := map[string]float64{
		"Technical accuracy":     report.TechnicalAccuracy,
		"Business relevance":     report.BusinessRelevance,
		"Implementation success": report.ImplementationSuccess,
		"Peer review quality":    report.PeerReviewScore,
		"Client satisfaction":    report.ClientOutcomeScore,
	}

	// Sort by score to identify lowest performing areas
	type scoreItem struct {
		area  string
		score float64
	}

	var items []scoreItem
	for area, score := range scores {
		items = append(items, scoreItem{area, score})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].score < items[j].score
	})

	// Take the lowest 2-3 areas for improvement
	for i := 0; i < len(items) && i < 3; i++ {
		if items[i].score < 0.8 {
			areas = append(areas, items[i].area)
		}
	}

	return areas
}

// TriggerQualityAlert triggers a quality alert
func (s *QualityAssuranceService) TriggerQualityAlert(ctx context.Context, alert *interfaces.QualityAlert) error {
	s.logger.WithFields(logrus.Fields{
		"alert_id": alert.ID,
		"severity": alert.Severity,
		"metric":   alert.MetricType,
		"value":    alert.MetricValue,
	}).Warn("Quality alert triggered")

	// Store alert in database
	query := `
		INSERT INTO quality_alerts (
			id, alert_type, severity, metric_type, metric_value,
			threshold_value, description, context, recipients,
			actions, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	contextJSON, _ := json.Marshal(alert.Context)
	recipientsJSON, _ := json.Marshal(alert.Recipients)
	actionsJSON, _ := json.Marshal(alert.Actions)

	_, err := s.db.Exec(ctx, query,
		alert.ID,
		alert.AlertType,
		alert.Severity,
		alert.MetricType,
		alert.MetricValue,
		alert.Threshold,
		alert.Description,
		string(contextJSON),
		string(recipientsJSON),
		string(actionsJSON),
		alert.Status,
		alert.CreatedAt,
		alert.UpdatedAt,
	)

	if err != nil {
		s.logger.WithError(err).Error("Failed to store quality alert")
		return err
	}

	// Update metrics
	s.metricsCollector.IncrementCounter("quality_alerts_triggered", map[string]string{
		"severity":    alert.Severity,
		"metric_type": alert.MetricType,
		"alert_type":  alert.AlertType,
	})

	return nil
}

// GetQualityTrends retrieves quality trends based on filters
func (s *QualityAssuranceService) GetQualityTrends(ctx context.Context, filters *interfaces.TrendFilters) (*interfaces.QualityTrends, error) {
	s.logger.WithFields(logrus.Fields{
		"metric_type":   filters.MetricType,
		"consultant_id": filters.ConsultantID,
		"granularity":   filters.Granularity,
	}).Info("Getting quality trends")

	// Set default granularity if not specified
	granularity := filters.Granularity
	if granularity == "" {
		granularity = "weekly"
	}

	// For this implementation, return mock trend data
	trends := &interfaces.QualityTrends{
		MetricType:     filters.MetricType,
		TrendDirection: "improving",
		TrendStrength:  0.7,
		DataPoints: []interfaces.QualityMetricTrendPoint{
			{
				Date:               time.Now().AddDate(0, 0, -21),
				OverallQuality:     0.75,
				TechnicalAccuracy:  0.78,
				BusinessRelevance:  0.72,
				ClientSatisfaction: 7.5,
				ImplementationRate: 0.73,
			},
			{
				Date:               time.Now().AddDate(0, 0, -14),
				OverallQuality:     0.78,
				TechnicalAccuracy:  0.80,
				BusinessRelevance:  0.75,
				ClientSatisfaction: 7.8,
				ImplementationRate: 0.76,
			},
			{
				Date:               time.Now().AddDate(0, 0, -7),
				OverallQuality:     0.82,
				TechnicalAccuracy:  0.83,
				BusinessRelevance:  0.78,
				ClientSatisfaction: 8.2,
				ImplementationRate: 0.80,
			},
		},
		Anomalies: []interfaces.QualityAnomaly{},
		Forecasts: []interfaces.ForecastPoint{
			{
				Date:          time.Now().AddDate(0, 0, 7),
				ForecastValue: 0.85,
				ConfidenceInterval: struct {
					Lower float64 `json:"lower"`
					Upper float64 `json:"upper"`
				}{
					Lower: 0.80,
					Upper: 0.90,
				},
				Confidence: 0.85,
			},
		},
		Insights: []string{
			"Quality metrics show consistent improvement over the analysis period",
			"Technical accuracy has improved by 6.4% in the last 3 weeks",
			"Client satisfaction is trending upward",
		},
		Seasonality: &interfaces.SeasonalityAnalysis{
			HasSeasonality: false,
			Confidence:     0.3,
		},
		Metadata: map[string]interface{}{
			"granularity":   granularity,
			"data_points":   3,
			"consultant_id": filters.ConsultantID,
		},
		GeneratedAt: time.Now(),
	}

	return trends, nil
}
