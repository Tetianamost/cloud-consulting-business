package interfaces

import (
	"context"
	"time"
)

// QualityAssuranceService defines the interface for AI recommendation quality assurance
type QualityAssuranceService interface {
	// Recommendation accuracy tracking
	TrackRecommendationAccuracy(ctx context.Context, recommendation *RecommendationTracking) error
	GetAccuracyMetrics(ctx context.Context, filters *AccuracyFilters) (*QualityAccuracyMetrics, error)
	UpdateRecommendationOutcome(ctx context.Context, recommendationID string, outcome *RecommendationOutcome) error

	// Peer review system
	SubmitForPeerReview(ctx context.Context, review *PeerReviewRequest) error
	GetPendingReviews(ctx context.Context, reviewerID string) ([]*PeerReview, error)
	SubmitPeerReview(ctx context.Context, reviewID string, feedback *PeerReviewFeedback) error
	GetReviewHistory(ctx context.Context, filters *ReviewFilters) ([]*PeerReview, error)

	// Client outcome tracking
	RecordClientOutcome(ctx context.Context, outcome *ClientOutcome) error
	GetOutcomeAnalytics(ctx context.Context, filters *OutcomeFilters) (*OutcomeAnalytics, error)
	ValidateRecommendationEffectiveness(ctx context.Context, recommendationID string) (*EffectivenessReport, error)

	// Continuous improvement system
	GenerateImprovementInsights(ctx context.Context, timeRange *TimeRange) (*ImprovementInsights, error)
	UpdateQualityThresholds(ctx context.Context, thresholds *QualityThresholds) error
	GetQualityTrends(ctx context.Context, filters *TrendFilters) (*QualityTrends, error)
	TriggerQualityAlert(ctx context.Context, alert *QualityAlert) error

	// Quality control validation
	ValidateRecommendationQuality(ctx context.Context, recommendation *AIRecommendation) (*QualityValidation, error)
	GetQualityScore(ctx context.Context, recommendationID string) (*QualityScore, error)
	SetQualityStandards(ctx context.Context, standards *QualityStandards) error
}

// RecommendationTracking represents tracking data for AI recommendations
type RecommendationTracking struct {
	ID                 string                 `json:"id"`
	RecommendationID   string                 `json:"recommendation_id"`
	InquiryID          string                 `json:"inquiry_id"`
	ConsultantID       string                 `json:"consultant_id"`
	RecommendationType string                 `json:"recommendation_type"` // "architecture", "cost", "security", "migration"
	Content            string                 `json:"content"`
	Confidence         float64                `json:"confidence"` // 0-1 scale
	GeneratedAt        time.Time              `json:"generated_at"`
	ModelVersion       string                 `json:"model_version"`
	PromptVersion      string                 `json:"prompt_version"`
	Context            map[string]interface{} `json:"context"`
	Tags               []string               `json:"tags"`
	Status             string                 `json:"status"` // "pending", "reviewed", "validated", "rejected"
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// AccuracyFilters defines filters for accuracy metrics queries
type AccuracyFilters struct {
	ConsultantID       string            `json:"consultant_id,omitempty"`
	RecommendationType string            `json:"recommendation_type,omitempty"`
	TimeRange          *QualityTimeRange `json:"time_range,omitempty"`
	MinConfidence      float64           `json:"min_confidence,omitempty"`
	Status             string            `json:"status,omitempty"`
	Tags               []string          `json:"tags,omitempty"`
}

// QualityAccuracyMetrics represents accuracy metrics for recommendations
type QualityAccuracyMetrics struct {
	OverallAccuracy          float64                `json:"overall_accuracy"`
	AccuracyByType           map[string]float64     `json:"accuracy_by_type"`
	AccuracyByConsultant     map[string]float64     `json:"accuracy_by_consultant"`
	TotalRecommendations     int64                  `json:"total_recommendations"`
	ValidatedRecommendations int64                  `json:"validated_recommendations"`
	AcceptedRecommendations  int64                  `json:"accepted_recommendations"`
	RejectedRecommendations  int64                  `json:"rejected_recommendations"`
	AverageConfidence        float64                `json:"average_confidence"`
	TrendData                []AccuracyTrendPoint   `json:"trend_data"`
	QualityDistribution      map[string]int64       `json:"quality_distribution"`
	Metadata                 map[string]interface{} `json:"metadata"`
	GeneratedAt              time.Time              `json:"generated_at"`
}

// AccuracyTrendPoint represents a point in accuracy trend data
type AccuracyTrendPoint struct {
	Date     time.Time `json:"date"`
	Accuracy float64   `json:"accuracy"`
	Count    int64     `json:"count"`
}

// RecommendationOutcome represents the outcome of a recommendation
type RecommendationOutcome struct {
	RecommendationID   string                 `json:"recommendation_id"`
	OutcomeType        string                 `json:"outcome_type"` // "accepted", "rejected", "modified", "implemented"
	ActualResult       string                 `json:"actual_result"`
	ClientFeedback     string                 `json:"client_feedback"`
	ConsultantNotes    string                 `json:"consultant_notes"`
	ImplementationCost float64                `json:"implementation_cost"`
	TimeToImplement    string                 `json:"time_to_implement"`
	BusinessImpact     string                 `json:"business_impact"`
	TechnicalImpact    string                 `json:"technical_impact"`
	LessonsLearned     []string               `json:"lessons_learned"`
	Accuracy           float64                `json:"accuracy"`      // 0-1 scale
	Effectiveness      float64                `json:"effectiveness"` // 0-1 scale
	Metadata           map[string]interface{} `json:"metadata"`
	RecordedAt         time.Time              `json:"recorded_at"`
	RecordedBy         string                 `json:"recorded_by"`
}

// PeerReviewRequest represents a request for peer review
type PeerReviewRequest struct {
	ID               string                 `json:"id"`
	RecommendationID string                 `json:"recommendation_id"`
	RequestedBy      string                 `json:"requested_by"`
	AssignedTo       string                 `json:"assigned_to"`
	Priority         string                 `json:"priority"`    // "low", "medium", "high", "urgent"
	ReviewType       string                 `json:"review_type"` // "technical", "business", "compliance", "comprehensive"
	Context          string                 `json:"context"`
	SpecificAreas    []string               `json:"specific_areas"` // Areas to focus review on
	DueDate          *time.Time             `json:"due_date,omitempty"`
	Instructions     string                 `json:"instructions"`
	Metadata         map[string]interface{} `json:"metadata"`
	Status           string                 `json:"status"` // "pending", "in_progress", "completed", "cancelled"
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// PeerReview represents a peer review record
type PeerReview struct {
	ID               string                 `json:"id"`
	RequestID        string                 `json:"request_id"`
	RecommendationID string                 `json:"recommendation_id"`
	ReviewerID       string                 `json:"reviewer_id"`
	ReviewerName     string                 `json:"reviewer_name"`
	ReviewType       string                 `json:"review_type"`
	Status           string                 `json:"status"`
	StartedAt        *time.Time             `json:"started_at,omitempty"`
	CompletedAt      *time.Time             `json:"completed_at,omitempty"`
	TimeSpent        int64                  `json:"time_spent"` // minutes
	Feedback         *PeerReviewFeedback    `json:"feedback,omitempty"`
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// PeerReviewFeedback represents feedback from a peer review
type PeerReviewFeedback struct {
	OverallRating     int                    `json:"overall_rating"`     // 1-5 scale
	TechnicalAccuracy int                    `json:"technical_accuracy"` // 1-5 scale
	BusinessRelevance int                    `json:"business_relevance"` // 1-5 scale
	Completeness      int                    `json:"completeness"`       // 1-5 scale
	Clarity           int                    `json:"clarity"`            // 1-5 scale
	Actionability     int                    `json:"actionability"`      // 1-5 scale
	Comments          string                 `json:"comments"`
	Strengths         []string               `json:"strengths"`
	Weaknesses        []string               `json:"weaknesses"`
	Suggestions       []string               `json:"suggestions"`
	RequiredChanges   []string               `json:"required_changes"`
	OptionalChanges   []string               `json:"optional_changes"`
	Approved          bool                   `json:"approved"`
	ApprovalLevel     string                 `json:"approval_level"` // "conditional", "full", "rejected"
	FollowUpRequired  bool                   `json:"follow_up_required"`
	FollowUpNotes     string                 `json:"follow_up_notes"`
	Metadata          map[string]interface{} `json:"metadata"`
	SubmittedAt       time.Time              `json:"submitted_at"`
}

// ReviewFilters defines filters for review history queries
type ReviewFilters struct {
	ReviewerID    string            `json:"reviewer_id,omitempty"`
	RequestedBy   string            `json:"requested_by,omitempty"`
	ReviewType    string            `json:"review_type,omitempty"`
	Status        string            `json:"status,omitempty"`
	Priority      string            `json:"priority,omitempty"`
	TimeRange     *QualityTimeRange `json:"time_range,omitempty"`
	MinRating     int               `json:"min_rating,omitempty"`
	ApprovalLevel string            `json:"approval_level,omitempty"`
}

// ClientOutcome represents client engagement outcomes
type ClientOutcome struct {
	ID                    string                 `json:"id"`
	InquiryID             string                 `json:"inquiry_id"`
	ClientName            string                 `json:"client_name"`
	RecommendationIDs     []string               `json:"recommendation_ids"`
	EngagementType        string                 `json:"engagement_type"`     // "consultation", "implementation", "ongoing"
	OutcomeType           string                 `json:"outcome_type"`        // "successful", "partially_successful", "unsuccessful", "cancelled"
	ClientSatisfaction    int                    `json:"client_satisfaction"` // 1-10 scale
	BusinessValue         float64                `json:"business_value"`      // Monetary value delivered
	CostSavings           float64                `json:"cost_savings"`
	TimeToValue           int64                  `json:"time_to_value"`       // Days
	ImplementationRate    float64                `json:"implementation_rate"` // Percentage of recommendations implemented
	ClientTestimonial     string                 `json:"client_testimonial"`
	ChallengesFaced       []string               `json:"challenges_faced"`
	SuccessFactors        []string               `json:"success_factors"`
	LessonsLearned        []string               `json:"lessons_learned"`
	FollowUpOpportunities []string               `json:"follow_up_opportunities"`
	ReferencePermission   bool                   `json:"reference_permission"`
	CaseStudyPermission   bool                   `json:"case_study_permission"`
	NetPromoterScore      int                    `json:"net_promoter_score"` // -100 to 100
	Metadata              map[string]interface{} `json:"metadata"`
	RecordedAt            time.Time              `json:"recorded_at"`
	RecordedBy            string                 `json:"recorded_by"`
}

// OutcomeFilters defines filters for outcome analytics queries
type OutcomeFilters struct {
	ClientName         string            `json:"client_name,omitempty"`
	EngagementType     string            `json:"engagement_type,omitempty"`
	OutcomeType        string            `json:"outcome_type,omitempty"`
	MinSatisfaction    int               `json:"min_satisfaction,omitempty"`
	MinBusinessValue   float64           `json:"min_business_value,omitempty"`
	TimeRange          *QualityTimeRange `json:"time_range,omitempty"`
	ConsultantID       string            `json:"consultant_id,omitempty"`
	RecommendationType string            `json:"recommendation_type,omitempty"`
}

// OutcomeAnalytics represents analytics for client outcomes
type OutcomeAnalytics struct {
	TotalEngagements          int64                     `json:"total_engagements"`
	SuccessfulEngagements     int64                     `json:"successful_engagements"`
	SuccessRate               float64                   `json:"success_rate"`
	AverageClientSatisfaction float64                   `json:"average_client_satisfaction"`
	TotalBusinessValue        float64                   `json:"total_business_value"`
	TotalCostSavings          float64                   `json:"total_cost_savings"`
	AverageTimeToValue        float64                   `json:"average_time_to_value"`
	AverageImplementationRate float64                   `json:"average_implementation_rate"`
	AverageNetPromoterScore   float64                   `json:"average_net_promoter_score"`
	OutcomesByType            map[string]int64          `json:"outcomes_by_type"`
	SatisfactionDistribution  map[string]int64          `json:"satisfaction_distribution"`
	BusinessValueTrend        []BusinessValueTrendPoint `json:"business_value_trend"`
	TopSuccessFactors         []SuccessFactor           `json:"top_success_factors"`
	CommonChallenges          []Challenge               `json:"common_challenges"`
	Metadata                  map[string]interface{}    `json:"metadata"`
	GeneratedAt               time.Time                 `json:"generated_at"`
}

// BusinessValueTrendPoint represents a point in business value trend
type BusinessValueTrendPoint struct {
	Date          time.Time `json:"date"`
	BusinessValue float64   `json:"business_value"`
	CostSavings   float64   `json:"cost_savings"`
	Engagements   int64     `json:"engagements"`
}

// SuccessFactor represents a success factor with frequency
type SuccessFactor struct {
	Factor    string  `json:"factor"`
	Frequency int64   `json:"frequency"`
	Impact    float64 `json:"impact"`
}

// Challenge represents a common challenge with frequency
type Challenge struct {
	Challenge string  `json:"challenge"`
	Frequency int64   `json:"frequency"`
	Impact    float64 `json:"impact"`
}

// EffectivenessReport represents a recommendation effectiveness report
type EffectivenessReport struct {
	RecommendationID      string                 `json:"recommendation_id"`
	OverallEffectiveness  float64                `json:"overall_effectiveness"` // 0-1 scale
	TechnicalAccuracy     float64                `json:"technical_accuracy"`
	BusinessRelevance     float64                `json:"business_relevance"`
	ImplementationSuccess float64                `json:"implementation_success"`
	ClientSatisfaction    float64                `json:"client_satisfaction"`
	BusinessImpact        float64                `json:"business_impact"`
	CostEffectiveness     float64                `json:"cost_effectiveness"`
	TimeEffectiveness     float64                `json:"time_effectiveness"`
	PeerReviewScore       float64                `json:"peer_review_score"`
	ClientOutcomeScore    float64                `json:"client_outcome_score"`
	Strengths             []string               `json:"strengths"`
	Weaknesses            []string               `json:"weaknesses"`
	ImprovementAreas      []string               `json:"improvement_areas"`
	Recommendations       []string               `json:"recommendations"`
	Metadata              map[string]interface{} `json:"metadata"`
	GeneratedAt           time.Time              `json:"generated_at"`
}

// ImprovementInsights represents insights for continuous improvement
type ImprovementInsights struct {
	TimeRange                     *QualityTimeRange           `json:"time_range"`
	OverallQualityTrend           string                      `json:"overall_quality_trend"` // "improving", "stable", "declining"
	KeyImprovementAreas           []ImprovementArea           `json:"key_improvement_areas"`
	SuccessPatterns               []Pattern                   `json:"success_patterns"`
	FailurePatterns               []Pattern                   `json:"failure_patterns"`
	ModelPerformanceInsights      []ModelInsight              `json:"model_performance_insights"`
	ConsultantPerformanceInsights []ConsultantInsight         `json:"consultant_performance_insights"`
	RecommendedActions            []RecommendedAction         `json:"recommended_actions"`
	QualityMetricsTrend           []QualityMetricTrendPoint   `json:"quality_metrics_trend"`
	BenchmarkComparison           *QualityBenchmarkComparison `json:"benchmark_comparison"`
	Metadata                      map[string]interface{}      `json:"metadata"`
	GeneratedAt                   time.Time                   `json:"generated_at"`
}

// ImprovementArea represents an area for improvement
type ImprovementArea struct {
	Area            string   `json:"area"`
	CurrentScore    float64  `json:"current_score"`
	TargetScore     float64  `json:"target_score"`
	Priority        string   `json:"priority"` // "high", "medium", "low"
	Impact          string   `json:"impact"`   // "high", "medium", "low"
	Effort          string   `json:"effort"`   // "high", "medium", "low"
	Description     string   `json:"description"`
	Recommendations []string `json:"recommendations"`
}

// Pattern represents a success or failure pattern
type Pattern struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Frequency   int64    `json:"frequency"`
	Confidence  float64  `json:"confidence"`
	Conditions  []string `json:"conditions"`
	Outcomes    []string `json:"outcomes"`
	Examples    []string `json:"examples"`
}

// ModelInsight represents insights about model performance
type ModelInsight struct {
	ModelVersion     string                 `json:"model_version"`
	PerformanceScore float64                `json:"performance_score"`
	Strengths        []string               `json:"strengths"`
	Weaknesses       []string               `json:"weaknesses"`
	Recommendations  []string               `json:"recommendations"`
	UsageStats       map[string]interface{} `json:"usage_stats"`
}

// ConsultantInsight represents insights about consultant performance
type ConsultantInsight struct {
	ConsultantID     string   `json:"consultant_id"`
	ConsultantName   string   `json:"consultant_name"`
	PerformanceScore float64  `json:"performance_score"`
	Strengths        []string `json:"strengths"`
	ImprovementAreas []string `json:"improvement_areas"`
	TrainingNeeds    []string `json:"training_needs"`
	BestPractices    []string `json:"best_practices"`
}

// RecommendedAction represents a recommended improvement action
type RecommendedAction struct {
	Action      string   `json:"action"`
	Priority    string   `json:"priority"`
	Impact      string   `json:"impact"`
	Effort      string   `json:"effort"`
	Timeline    string   `json:"timeline"`
	Owner       string   `json:"owner"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

// QualityMetricTrendPoint represents a point in quality metrics trend
type QualityMetricTrendPoint struct {
	Date               time.Time `json:"date"`
	OverallQuality     float64   `json:"overall_quality"`
	TechnicalAccuracy  float64   `json:"technical_accuracy"`
	BusinessRelevance  float64   `json:"business_relevance"`
	ClientSatisfaction float64   `json:"client_satisfaction"`
	ImplementationRate float64   `json:"implementation_rate"`
	PeerReviewScore    float64   `json:"peer_review_score"`
}

// QualityBenchmarkComparison represents comparison against benchmarks
type QualityBenchmarkComparison struct {
	IndustryBenchmark  float64  `json:"industry_benchmark"`
	CompanyBenchmark   float64  `json:"company_benchmark"`
	CurrentPerformance float64  `json:"current_performance"`
	PerformanceGap     float64  `json:"performance_gap"`
	RankingPercentile  float64  `json:"ranking_percentile"`
	ComparisonInsights []string `json:"comparison_insights"`
}

// QualityThresholds represents quality thresholds for alerts
type QualityThresholds struct {
	MinAccuracy           float64            `json:"min_accuracy"`
	MinClientSatisfaction int                `json:"min_client_satisfaction"`
	MinPeerReviewScore    float64            `json:"min_peer_review_score"`
	MaxResponseTime       int64              `json:"max_response_time"` // milliseconds
	MinImplementationRate float64            `json:"min_implementation_rate"`
	AlertThresholds       map[string]float64 `json:"alert_thresholds"`
	UpdatedBy             string             `json:"updated_by"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

// TrendFilters defines filters for quality trend queries
type TrendFilters struct {
	MetricType     string            `json:"metric_type,omitempty"`
	ConsultantID   string            `json:"consultant_id,omitempty"`
	TimeRange      *QualityTimeRange `json:"time_range,omitempty"`
	Granularity    string            `json:"granularity,omitempty"`     // "daily", "weekly", "monthly"
	ComparisonType string            `json:"comparison_type,omitempty"` // "period", "benchmark"
}

// QualityTrends represents quality trends over time
type QualityTrends struct {
	MetricType     string                    `json:"metric_type"`
	TrendDirection string                    `json:"trend_direction"` // "improving", "stable", "declining"
	TrendStrength  float64                   `json:"trend_strength"`  // 0-1 scale
	DataPoints     []QualityMetricTrendPoint `json:"data_points"`
	Seasonality    *SeasonalityAnalysis      `json:"seasonality,omitempty"`
	Anomalies      []QualityAnomaly          `json:"anomalies"`
	Forecasts      []ForecastPoint           `json:"forecasts"`
	Insights       []string                  `json:"insights"`
	Metadata       map[string]interface{}    `json:"metadata"`
	GeneratedAt    time.Time                 `json:"generated_at"`
}

// SeasonalityAnalysis represents seasonality analysis
type SeasonalityAnalysis struct {
	HasSeasonality  bool               `json:"has_seasonality"`
	SeasonalPeriod  string             `json:"seasonal_period"`
	SeasonalFactors map[string]float64 `json:"seasonal_factors"`
	Confidence      float64            `json:"confidence"`
}

// QualityAnomaly represents an anomaly in quality metrics
type QualityAnomaly struct {
	Date           time.Time `json:"date"`
	MetricValue    float64   `json:"metric_value"`
	ExpectedValue  float64   `json:"expected_value"`
	Deviation      float64   `json:"deviation"`
	Severity       string    `json:"severity"` // "low", "medium", "high"
	Description    string    `json:"description"`
	PossibleCauses []string  `json:"possible_causes"`
}

// ForecastPoint represents a forecast point
type ForecastPoint struct {
	Date               time.Time `json:"date"`
	ForecastValue      float64   `json:"forecast_value"`
	ConfidenceInterval struct {
		Lower float64 `json:"lower"`
		Upper float64 `json:"upper"`
	} `json:"confidence_interval"`
	Confidence float64 `json:"confidence"`
}

// QualityAlert represents a quality alert
type QualityAlert struct {
	ID          string                 `json:"id"`
	AlertType   string                 `json:"alert_type"` // "threshold", "trend", "anomaly"
	Severity    string                 `json:"severity"`   // "low", "medium", "high", "critical"
	MetricType  string                 `json:"metric_type"`
	MetricValue float64                `json:"metric_value"`
	Threshold   float64                `json:"threshold"`
	Description string                 `json:"description"`
	Context     map[string]interface{} `json:"context"`
	Recipients  []string               `json:"recipients"`
	Actions     []string               `json:"actions"`
	Status      string                 `json:"status"` // "active", "acknowledged", "resolved"
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// AIRecommendation represents an AI-generated recommendation for validation
type AIRecommendation struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Content       string                 `json:"content"`
	Confidence    float64                `json:"confidence"`
	Context       map[string]interface{} `json:"context"`
	GeneratedBy   string                 `json:"generated_by"`
	ModelVersion  string                 `json:"model_version"`
	PromptVersion string                 `json:"prompt_version"`
	GeneratedAt   time.Time              `json:"generated_at"`
}

// QualityValidation represents the result of quality validation
type QualityValidation struct {
	RecommendationID string                 `json:"recommendation_id"`
	OverallScore     float64                `json:"overall_score"` // 0-1 scale
	ValidationChecks []ValidationCheck      `json:"validation_checks"`
	PassedChecks     int                    `json:"passed_checks"`
	TotalChecks      int                    `json:"total_checks"`
	PassRate         float64                `json:"pass_rate"`
	Issues           []QualityIssue         `json:"issues"`
	Recommendations  []string               `json:"recommendations"`
	RequiresReview   bool                   `json:"requires_review"`
	Metadata         map[string]interface{} `json:"metadata"`
	ValidatedAt      time.Time              `json:"validated_at"`
	ValidatedBy      string                 `json:"validated_by"`
}

// ValidationCheck represents a single validation check
type ValidationCheck struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Passed      bool     `json:"passed"`
	Score       float64  `json:"score"`
	Weight      float64  `json:"weight"`
	Details     string   `json:"details"`
	Suggestions []string `json:"suggestions"`
}

// QualityIssue represents a quality issue found during validation
type QualityIssue struct {
	Type        string   `json:"type"`     // "accuracy", "completeness", "clarity", "relevance"
	Severity    string   `json:"severity"` // "low", "medium", "high", "critical"
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Suggestions []string `json:"suggestions"`
	AutoFixable bool     `json:"auto_fixable"`
}

// QualityScore represents a comprehensive quality score
type QualityScore struct {
	RecommendationID     string                 `json:"recommendation_id"`
	OverallScore         float64                `json:"overall_score"` // 0-1 scale
	ComponentScores      map[string]float64     `json:"component_scores"`
	WeightedScore        float64                `json:"weighted_score"`
	HistoricalComparison float64                `json:"historical_comparison"`
	BenchmarkComparison  float64                `json:"benchmark_comparison"`
	ScoreBreakdown       *ScoreBreakdown        `json:"score_breakdown"`
	QualityGrade         string                 `json:"quality_grade"` // "A", "B", "C", "D", "F"
	Metadata             map[string]interface{} `json:"metadata"`
	CalculatedAt         time.Time              `json:"calculated_at"`
}

// ScoreBreakdown represents detailed score breakdown
type ScoreBreakdown struct {
	TechnicalAccuracy float64 `json:"technical_accuracy"`
	BusinessRelevance float64 `json:"business_relevance"`
	Completeness      float64 `json:"completeness"`
	Clarity           float64 `json:"clarity"`
	Actionability     float64 `json:"actionability"`
	Innovation        float64 `json:"innovation"`
	RiskAssessment    float64 `json:"risk_assessment"`
	CostEffectiveness float64 `json:"cost_effectiveness"`
}

// QualityStandards represents quality standards configuration
type QualityStandards struct {
	MinOverallScore    float64                `json:"min_overall_score"`
	ComponentWeights   map[string]float64     `json:"component_weights"`
	ValidationRules    []ValidationRule       `json:"validation_rules"`
	ReviewThresholds   map[string]float64     `json:"review_thresholds"`
	AlertThresholds    map[string]float64     `json:"alert_thresholds"`
	QualityGradeRanges map[string]ScoreRange  `json:"quality_grade_ranges"`
	AutoReviewTriggers []AutoReviewTrigger    `json:"auto_review_triggers"`
	Metadata           map[string]interface{} `json:"metadata"`
	UpdatedBy          string                 `json:"updated_by"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// ScoreRange represents a score range for quality grades
type ScoreRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// AutoReviewTrigger represents conditions that trigger automatic review
type AutoReviewTrigger struct {
	Name        string                 `json:"name"`
	Conditions  map[string]interface{} `json:"conditions"`
	ReviewType  string                 `json:"review_type"`
	Priority    string                 `json:"priority"`
	Assignee    string                 `json:"assignee"`
	Description string                 `json:"description"`
}

// QualityTimeRange represents a time range for queries
type QualityTimeRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
