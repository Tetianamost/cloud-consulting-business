package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// PerformanceAnalyticsService defines the interface for consultant performance analytics
type PerformanceAnalyticsService interface {
	// Engagement success tracking
	TrackEngagementOutcome(ctx context.Context, engagement *EngagementOutcome) error
	GetEngagementSuccessMetrics(ctx context.Context, filters *EngagementFilters) (*EngagementMetrics, error)
	AnalyzeEngagementPatterns(ctx context.Context, consultantID string, timeRange *TimeRange) (*EngagementPatterns, error)

	// Client satisfaction correlation
	RecordClientFeedback(ctx context.Context, feedback *ClientFeedback) error
	AnalyzeRecommendationEffectiveness(ctx context.Context, filters *RecommendationFilters) (*RecommendationAnalysis, error)
	GetSatisfactionCorrelations(ctx context.Context, consultantID string) (*SatisfactionCorrelations, error)

	// Consultant skill gap analysis
	AnalyzeConsultantSkills(ctx context.Context, consultantID string) (*SkillGapAnalysis, error)
	GetSkillDevelopmentRecommendations(ctx context.Context, consultantID string) ([]*SkillRecommendation, error)
	TrackSkillImprovement(ctx context.Context, consultantID string, skillArea string, improvement *SkillImprovement) error

	// Knowledge sharing system
	CaptureSuccessPattern(ctx context.Context, pattern *SuccessPattern) error
	GetSuccessPatterns(ctx context.Context, filters *PatternFilters) ([]*SuccessPattern, error)
	RecommendPatterns(ctx context.Context, inquiry *domain.Inquiry) ([]*PatternRecommendation, error)
	ShareKnowledge(ctx context.Context, knowledge *KnowledgeItem) error
	GetKnowledgeBase(ctx context.Context, filters *KnowledgeFilters) ([]*KnowledgeItem, error)

	// Analytics and reporting
	GeneratePerformanceReport(ctx context.Context, consultantID string, timeRange *TimeRange) (*PerformanceReport, error)
	GetTeamAnalytics(ctx context.Context, timeRange *TimeRange) (*TeamAnalytics, error)
	GetBenchmarkMetrics(ctx context.Context) (*BenchmarkMetrics, error)
}

// EngagementOutcome represents the outcome of a client engagement
type EngagementOutcome struct {
	ID                   string                 `json:"id"`
	InquiryID            string                 `json:"inquiry_id"`
	ConsultantID         string                 `json:"consultant_id"`
	ClientID             string                 `json:"client_id"`
	ProjectType          string                 `json:"project_type"`
	Industry             string                 `json:"industry"`
	RecommendationTypes  []string               `json:"recommendation_types"`
	ImplementationStatus string                 `json:"implementation_status"` // "not_started", "in_progress", "completed", "cancelled"
	SuccessMetrics       map[string]float64     `json:"success_metrics"`
	ClientSatisfaction   float64                `json:"client_satisfaction"` // 1-10 scale
	BusinessImpact       *BusinessImpact        `json:"business_impact"`
	TechnicalOutcomes    *TechnicalOutcomes     `json:"technical_outcomes"`
	LessonsLearned       []string               `json:"lessons_learned"`
	ChallengesFaced      []string               `json:"challenges_faced"`
	BestPracticesUsed    []string               `json:"best_practices_used"`
	TimeToValue          *time.Duration         `json:"time_to_value"`
	ProjectDuration      *time.Duration         `json:"project_duration"`
	BudgetVariance       float64                `json:"budget_variance"` // percentage over/under budget
	Metadata             map[string]interface{} `json:"metadata"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

// BusinessImpact represents the business impact of an engagement
type BusinessImpact struct {
	CostSavings           float64 `json:"cost_savings"`
	RevenueIncrease       float64 `json:"revenue_increase"`
	EfficiencyGains       float64 `json:"efficiency_gains"` // percentage improvement
	TimeToMarketReduction string  `json:"time_to_market_reduction"`
	RiskReduction         string  `json:"risk_reduction"`
	ComplianceImprovement string  `json:"compliance_improvement"`
}

// TechnicalOutcomes represents technical outcomes of an engagement
type TechnicalOutcomes struct {
	PerformanceImprovement float64  `json:"performance_improvement"` // percentage
	SecurityPosture        string   `json:"security_posture"`        // "improved", "maintained", "degraded"
	ScalabilityGains       string   `json:"scalability_gains"`
	ReliabilityMetrics     float64  `json:"reliability_metrics"` // uptime percentage
	TechnicalDebtReduction string   `json:"technical_debt_reduction"`
	AutomationLevel        float64  `json:"automation_level"` // percentage of processes automated
	ServicesImplemented    []string `json:"services_implemented"`
	ArchitectureChanges    []string `json:"architecture_changes"`
}

// EngagementFilters represents filters for engagement queries
type EngagementFilters struct {
	ConsultantID    string     `json:"consultant_id,omitempty"`
	Industry        string     `json:"industry,omitempty"`
	ProjectType     string     `json:"project_type,omitempty"`
	Status          string     `json:"status,omitempty"`
	TimeRange       *TimeRange `json:"time_range,omitempty"`
	MinSatisfaction float64    `json:"min_satisfaction,omitempty"`
	Limit           int        `json:"limit"`
	Offset          int        `json:"offset"`
}

// EngagementMetrics represents metrics for engagement success
type EngagementMetrics struct {
	TotalEngagements          int64                  `json:"total_engagements"`
	SuccessfulEngagements     int64                  `json:"successful_engagements"`
	SuccessRate               float64                `json:"success_rate"`
	AverageClientSatisfaction float64                `json:"average_client_satisfaction"`
	AverageProjectDuration    time.Duration          `json:"average_project_duration"`
	AverageTimeToValue        time.Duration          `json:"average_time_to_value"`
	TotalCostSavings          float64                `json:"total_cost_savings"`
	TotalRevenueImpact        float64                `json:"total_revenue_impact"`
	TopRecommendationTypes    []RecommendationMetric `json:"top_recommendation_types"`
	IndustryBreakdown         map[string]int64       `json:"industry_breakdown"`
	TrendData                 []MetricTrend          `json:"trend_data"`
}

// RecommendationMetric represents metrics for a specific recommendation type
type RecommendationMetric struct {
	Type               string  `json:"type"`
	Count              int64   `json:"count"`
	SuccessRate        float64 `json:"success_rate"`
	AverageRating      float64 `json:"average_rating"`
	AverageImpact      float64 `json:"average_impact"`
	ImplementationRate float64 `json:"implementation_rate"`
}

// MetricTrend represents trend data for metrics over time
type MetricTrend struct {
	Period string  `json:"period"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"` // percentage change from previous period
}

// EngagementPatterns represents patterns in consultant engagements
type EngagementPatterns struct {
	ConsultantID            string                   `json:"consultant_id"`
	MostSuccessfulTypes     []string                 `json:"most_successful_types"`
	PreferredIndustries     []string                 `json:"preferred_industries"`
	StrengthAreas           []string                 `json:"strength_areas"`
	CommonChallenges        []string                 `json:"common_challenges"`
	BestPracticePatterns    []string                 `json:"best_practice_patterns"`
	ClientSatisfactionTrend []MetricTrend            `json:"client_satisfaction_trend"`
	PerformanceIndicators   map[string]float64       `json:"performance_indicators"`
	RecommendationPatterns  map[string]PatternMetric `json:"recommendation_patterns"`
	SeasonalTrends          map[string]float64       `json:"seasonal_trends"`
}

// PatternMetric represents metrics for a specific pattern
type PatternMetric struct {
	Frequency   int64   `json:"frequency"`
	SuccessRate float64 `json:"success_rate"`
	Impact      float64 `json:"impact"`
	Confidence  float64 `json:"confidence"`
}

// ClientFeedback represents feedback from clients
type ClientFeedback struct {
	ID                    string                 `json:"id"`
	InquiryID             string                 `json:"inquiry_id"`
	ConsultantID          string                 `json:"consultant_id"`
	OverallSatisfaction   float64                `json:"overall_satisfaction"` // 1-10 scale
	CommunicationRating   float64                `json:"communication_rating"`
	TechnicalExpertise    float64                `json:"technical_expertise"`
	ResponseTime          float64                `json:"response_time"`
	SolutionQuality       float64                `json:"solution_quality"`
	ValueForMoney         float64                `json:"value_for_money"`
	WouldRecommend        bool                   `json:"would_recommend"`
	SpecificFeedback      string                 `json:"specific_feedback"`
	AreasForImprovement   []string               `json:"areas_for_improvement"`
	StrengthsHighlighted  []string               `json:"strengths_highlighted"`
	RecommendationRatings map[string]float64     `json:"recommendation_ratings"` // recommendation type -> rating
	ImplementationSupport float64                `json:"implementation_support"`
	FollowUpSatisfaction  float64                `json:"follow_up_satisfaction"`
	Metadata              map[string]interface{} `json:"metadata"`
	CreatedAt             time.Time              `json:"created_at"`
}

// RecommendationAnalysis represents analysis of recommendation effectiveness
type RecommendationAnalysis struct {
	RecommendationType    string                 `json:"recommendation_type"`
	TotalRecommendations  int64                  `json:"total_recommendations"`
	ImplementationRate    float64                `json:"implementation_rate"`
	AverageClientRating   float64                `json:"average_client_rating"`
	AverageBusinessImpact float64                `json:"average_business_impact"`
	SuccessFactors        []string               `json:"success_factors"`
	FailureReasons        []string               `json:"failure_reasons"`
	IndustryEffectiveness map[string]float64     `json:"industry_effectiveness"`
	ConsultantPerformance map[string]float64     `json:"consultant_performance"`
	TrendAnalysis         []MetricTrend          `json:"trend_analysis"`
	Recommendations       []string               `json:"recommendations"`
	Metadata              map[string]interface{} `json:"metadata"`
}

// SatisfactionCorrelations represents correlations between various factors and client satisfaction
type SatisfactionCorrelations struct {
	ConsultantID               string                 `json:"consultant_id"`
	OverallSatisfactionScore   float64                `json:"overall_satisfaction_score"`
	RecommendationCorrelations map[string]float64     `json:"recommendation_correlations"`
	IndustryCorrelations       map[string]float64     `json:"industry_correlations"`
	ProjectSizeCorrelations    map[string]float64     `json:"project_size_correlations"`
	TimelineCorrelations       map[string]float64     `json:"timeline_correlations"`
	CommunicationCorrelations  map[string]float64     `json:"communication_correlations"`
	TechnicalSkillCorrelations map[string]float64     `json:"technical_skill_correlations"`
	StrongestPositiveFactors   []CorrelationFactor    `json:"strongest_positive_factors"`
	StrongestNegativeFactors   []CorrelationFactor    `json:"strongest_negative_factors"`
	ImprovementOpportunities   []string               `json:"improvement_opportunities"`
	Metadata                   map[string]interface{} `json:"metadata"`
}

// CorrelationFactor represents a factor that correlates with satisfaction
type CorrelationFactor struct {
	Factor      string  `json:"factor"`
	Correlation float64 `json:"correlation"`
	Confidence  float64 `json:"confidence"`
	SampleSize  int64   `json:"sample_size"`
}

// RecommendationFilters represents filters for recommendation analysis
type RecommendationFilters struct {
	RecommendationType string     `json:"recommendation_type,omitempty"`
	ConsultantID       string     `json:"consultant_id,omitempty"`
	Industry           string     `json:"industry,omitempty"`
	TimeRange          *TimeRange `json:"time_range,omitempty"`
	MinRating          float64    `json:"min_rating,omitempty"`
	Limit              int        `json:"limit"`
	Offset             int        `json:"offset"`
}

// SkillGapAnalysis represents analysis of consultant skill gaps
type SkillGapAnalysis struct {
	ConsultantID        string                 `json:"consultant_id"`
	OverallSkillLevel   float64                `json:"overall_skill_level"`
	SkillAreas          map[string]SkillMetric `json:"skill_areas"`
	StrengthAreas       []string               `json:"strength_areas"`
	ImprovementAreas    []string               `json:"improvement_areas"`
	CriticalGaps        []SkillGap             `json:"critical_gaps"`
	BenchmarkComparison map[string]float64     `json:"benchmark_comparison"`
	CareerProgression   *CareerProgression     `json:"career_progression"`
	TrainingNeeds       []SkillTrainingNeed    `json:"training_needs"`
	Metadata            map[string]interface{} `json:"metadata"`
	AnalyzedAt          time.Time              `json:"analyzed_at"`
}

// SkillMetric represents metrics for a specific skill area
type SkillMetric struct {
	CurrentLevel    float64       `json:"current_level"` // 1-10 scale
	RequiredLevel   float64       `json:"required_level"`
	Gap             float64       `json:"gap"`
	Trend           []MetricTrend `json:"trend"`
	ClientFeedback  float64       `json:"client_feedback"`
	PeerAssessment  float64       `json:"peer_assessment"`
	SelfAssessment  float64       `json:"self_assessment"`
	ProjectOutcomes float64       `json:"project_outcomes"`
}

// SkillGap represents a specific skill gap
type SkillGap struct {
	SkillArea       string  `json:"skill_area"`
	CurrentLevel    float64 `json:"current_level"`
	RequiredLevel   float64 `json:"required_level"`
	Gap             float64 `json:"gap"`
	Priority        string  `json:"priority"` // "critical", "high", "medium", "low"
	BusinessImpact  string  `json:"business_impact"`
	DevelopmentPath string  `json:"development_path"`
}

// CareerProgression represents career progression analysis
type CareerProgression struct {
	CurrentLevel   string    `json:"current_level"`
	NextLevel      string    `json:"next_level"`
	ProgressScore  float64   `json:"progress_score"`
	RequiredSkills []string  `json:"required_skills"`
	EstimatedTime  string    `json:"estimated_time"`
	Milestones     []string  `json:"milestones"`
	LastPromotion  time.Time `json:"last_promotion"`
}

// SkillTrainingNeed represents a specific training need for skill development
type SkillTrainingNeed struct {
	SkillArea             string   `json:"skill_area"`
	Priority              string   `json:"priority"`
	TrainingType          string   `json:"training_type"` // "course", "certification", "mentoring", "project"
	EstimatedHours        int      `json:"estimated_hours"`
	RecommendedPath       []string `json:"recommended_path"`
	BusinessJustification string   `json:"business_justification"`
}

// SkillRecommendation represents a recommendation for skill development
type SkillRecommendation struct {
	ID             string    `json:"id"`
	ConsultantID   string    `json:"consultant_id"`
	SkillArea      string    `json:"skill_area"`
	Recommendation string    `json:"recommendation"`
	Priority       string    `json:"priority"`
	ExpectedImpact string    `json:"expected_impact"`
	Resources      []string  `json:"resources"`
	Timeline       string    `json:"timeline"`
	SuccessMetrics []string  `json:"success_metrics"`
	CreatedAt      time.Time `json:"created_at"`
}

// SkillImprovement represents tracked improvement in a skill area
type SkillImprovement struct {
	SkillArea      string    `json:"skill_area"`
	PreviousLevel  float64   `json:"previous_level"`
	CurrentLevel   float64   `json:"current_level"`
	Improvement    float64   `json:"improvement"`
	TrainingTaken  []string  `json:"training_taken"`
	ProjectsWorked []string  `json:"projects_worked"`
	MentorFeedback string    `json:"mentor_feedback"`
	ClientFeedback string    `json:"client_feedback"`
	AssessmentDate time.Time `json:"assessment_date"`
}

// SuccessPattern represents a successful solution pattern
type SuccessPattern struct {
	ID                  string                 `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Category            string                 `json:"category"`
	Industry            string                 `json:"industry"`
	ProblemType         string                 `json:"problem_type"`
	SolutionApproach    string                 `json:"solution_approach"`
	TechnicalComponents []string               `json:"technical_components"`
	ImplementationSteps []string               `json:"implementation_steps"`
	SuccessMetrics      map[string]float64     `json:"success_metrics"`
	Prerequisites       []string               `json:"prerequisites"`
	RiskFactors         []string               `json:"risk_factors"`
	BestPractices       []string               `json:"best_practices"`
	LessonsLearned      []string               `json:"lessons_learned"`
	ApplicableScenarios []string               `json:"applicable_scenarios"`
	ConsultantID        string                 `json:"consultant_id"`
	ClientFeedback      string                 `json:"client_feedback"`
	BusinessImpact      *BusinessImpact        `json:"business_impact"`
	UsageCount          int64                  `json:"usage_count"`
	SuccessRate         float64                `json:"success_rate"`
	Tags                []string               `json:"tags"`
	Metadata            map[string]interface{} `json:"metadata"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// PatternRecommendation represents a recommended pattern for an inquiry
type PatternRecommendation struct {
	PatternID       string   `json:"pattern_id"`
	PatternName     string   `json:"pattern_name"`
	RelevanceScore  float64  `json:"relevance_score"`
	ConfidenceLevel float64  `json:"confidence_level"`
	Reasoning       string   `json:"reasoning"`
	Adaptations     []string `json:"adaptations"`
	ExpectedOutcome string   `json:"expected_outcome"`
}

// PatternFilters represents filters for success pattern queries
type PatternFilters struct {
	Category    string     `json:"category,omitempty"`
	Industry    string     `json:"industry,omitempty"`
	ProblemType string     `json:"problem_type,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	MinRating   float64    `json:"min_rating,omitempty"`
	TimeRange   *TimeRange `json:"time_range,omitempty"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

// KnowledgeFilters represents filters for knowledge base queries
type KnowledgeFilters struct {
	Category      string     `json:"category,omitempty"`
	Type          string     `json:"type,omitempty"`
	Industry      string     `json:"industry,omitempty"`
	TechnicalArea string     `json:"technical_area,omitempty"`
	Difficulty    string     `json:"difficulty,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	Author        string     `json:"author,omitempty"`
	MinRating     float64    `json:"min_rating,omitempty"`
	TimeRange     *TimeRange `json:"time_range,omitempty"`
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
}

// PerformanceReport represents a comprehensive performance report
type PerformanceReport struct {
	ConsultantID        string                          `json:"consultant_id"`
	ReportPeriod        *TimeRange                      `json:"report_period"`
	OverallPerformance  *OverallPerformance             `json:"overall_performance"`
	EngagementMetrics   *EngagementMetrics              `json:"engagement_metrics"`
	ClientSatisfaction  *ClientSatisfactionSummary      `json:"client_satisfaction"`
	SkillAssessment     *SkillGapAnalysis               `json:"skill_assessment"`
	SuccessPatterns     []*SuccessPattern               `json:"success_patterns"`
	ImprovementAreas    []string                        `json:"improvement_areas"`
	Achievements        []string                        `json:"achievements"`
	Goals               []PerformanceGoal               `json:"goals"`
	BenchmarkComparison *PerformanceBenchmarkComparison `json:"benchmark_comparison"`
	CareerDevelopment   *CareerProgression              `json:"career_development"`
	GeneratedAt         time.Time                       `json:"generated_at"`
}

// OverallPerformance represents overall performance metrics
type OverallPerformance struct {
	PerformanceScore float64  `json:"performance_score"` // 1-100 scale
	Ranking          int      `json:"ranking"`           // among peers
	TotalPeers       int      `json:"total_peers"`
	PerformanceTrend string   `json:"performance_trend"` // "improving", "stable", "declining"
	KeyStrengths     []string `json:"key_strengths"`
	DevelopmentAreas []string `json:"development_areas"`
}

// ClientSatisfactionSummary represents client satisfaction summary
type ClientSatisfactionSummary struct {
	AverageRating       float64          `json:"average_rating"`
	TotalFeedbacks      int64            `json:"total_feedbacks"`
	RatingDistribution  map[string]int64 `json:"rating_distribution"`
	TopStrengths        []string         `json:"top_strengths"`
	AreasForImprovement []string         `json:"areas_for_improvement"`
	TrendAnalysis       []MetricTrend    `json:"trend_analysis"`
	ComparisonToPeers   float64          `json:"comparison_to_peers"`
}

// PerformanceGoal represents a performance goal
type PerformanceGoal struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	TargetValue  float64   `json:"target_value"`
	CurrentValue float64   `json:"current_value"`
	Progress     float64   `json:"progress"` // percentage
	Deadline     time.Time `json:"deadline"`
	Status       string    `json:"status"` // "on_track", "at_risk", "behind", "completed"
}

// PerformanceBenchmarkComparison represents comparison to benchmarks (renamed to avoid conflict)
type PerformanceBenchmarkComparison struct {
	ClientSatisfaction float64 `json:"client_satisfaction"`
	ProjectSuccessRate float64 `json:"project_success_rate"`
	AverageProjectTime float64 `json:"average_project_time"`
	RevenuePerProject  float64 `json:"revenue_per_project"`
	SkillLevel         float64 `json:"skill_level"`
	PerformanceRanking int     `json:"performance_ranking"`
}

// TeamAnalytics represents team-wide analytics
type TeamAnalytics struct {
	TeamSize                  int64                        `json:"team_size"`
	TotalEngagements          int64                        `json:"total_engagements"`
	AverageClientSatisfaction float64                      `json:"average_client_satisfaction"`
	TopPerformers             []ConsultantRanking          `json:"top_performers"`
	SkillDistribution         map[string]SkillDistribution `json:"skill_distribution"`
	KnowledgeSharing          *KnowledgeSharingMetrics     `json:"knowledge_sharing"`
	TeamTrends                []MetricTrend                `json:"team_trends"`
	BenchmarkMetrics          *BenchmarkMetrics            `json:"benchmark_metrics"`
	ReportPeriod              *TimeRange                   `json:"report_period"`
	GeneratedAt               time.Time                    `json:"generated_at"`
}

// ConsultantRanking represents consultant ranking
type ConsultantRanking struct {
	ConsultantID       string   `json:"consultant_id"`
	Name               string   `json:"name"`
	PerformanceScore   float64  `json:"performance_score"`
	ClientSatisfaction float64  `json:"client_satisfaction"`
	ProjectSuccessRate float64  `json:"project_success_rate"`
	Specializations    []string `json:"specializations"`
}

// SkillDistribution represents skill distribution across the team
type SkillDistribution struct {
	SkillArea        string  `json:"skill_area"`
	AverageLevel     float64 `json:"average_level"`
	HighPerformers   int     `json:"high_performers"`
	NeedsDevelopment int     `json:"needs_development"`
	CriticalGaps     int     `json:"critical_gaps"`
}

// KnowledgeSharingMetrics represents knowledge sharing metrics
type KnowledgeSharingMetrics struct {
	TotalKnowledgeItems  int64    `json:"total_knowledge_items"`
	ActiveContributors   int64    `json:"active_contributors"`
	AverageRating        float64  `json:"average_rating"`
	MostViewedItems      []string `json:"most_viewed_items"`
	RecentContributions  int64    `json:"recent_contributions"`
	KnowledgeUtilization float64  `json:"knowledge_utilization"`
}

// BenchmarkMetrics represents benchmark metrics
type BenchmarkMetrics struct {
	IndustryAverages      map[string]float64 `json:"industry_averages"`
	CompanyTargets        map[string]float64 `json:"company_targets"`
	BestInClassMetrics    map[string]float64 `json:"best_in_class_metrics"`
	PerformanceThresholds map[string]float64 `json:"performance_thresholds"`
	LastUpdated           time.Time          `json:"last_updated"`
}
