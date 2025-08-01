package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// CompetitiveIntelligenceService defines the interface for competitive intelligence analysis
type CompetitiveIntelligenceService interface {
	// Competitor analysis
	AnalyzeCompetitors(ctx context.Context, inquiry *domain.Inquiry) (*CompetitiveAnalysis, error)
	GetCompetitorProfile(ctx context.Context, competitorName string) (*CompetitorProfile, error)
	UpdateCompetitorData(ctx context.Context, competitor *CompetitorProfile) error

	// Pricing intelligence
	GetPricingIntelligence(ctx context.Context, serviceType string, region string) (*PricingIntelligence, error)
	ComparePricing(ctx context.Context, services []string, competitors []string) (*PricingComparison, error)

	// Technology trends
	GetTechnologyTrends(ctx context.Context, category string) (*TechnologyTrends, error)
	AnalyzeTrendImpact(ctx context.Context, inquiry *domain.Inquiry) (*TrendImpact, error)

	// Differentiation strategies
	GenerateDifferentiationStrategy(ctx context.Context, inquiry *domain.Inquiry, competitors []string) (*CompetitiveDifferentiationStrategy, error)
	IdentifyCompetitorWeaknesses(ctx context.Context, competitors []string) (*CompetitorWeaknesses, error)
}

// CompetitiveAnalysis represents a comprehensive competitor analysis
type CompetitiveAnalysis struct {
	ID                string                     `json:"id"`
	InquiryID         string                     `json:"inquiry_id"`
	Competitors       []CompetitorProfile        `json:"competitors"`
	MarketPositioning *CompetitiveMarketPosition `json:"market_positioning"`
	CompetitiveMatrix *CompetitiveMatrix         `json:"competitive_matrix"`
	Recommendations   []string                   `json:"recommendations"`
	GeneratedAt       time.Time                  `json:"generated_at"`
}

// CompetitorProfile represents detailed information about a competitor
type CompetitorProfile struct {
	Name             string                      `json:"name"`
	Type             CompetitorType              `json:"type"` // "direct", "indirect", "substitute"
	MarketShare      float64                     `json:"market_share"`
	Strengths        []string                    `json:"strengths"`
	Weaknesses       []string                    `json:"weaknesses"`
	ServiceOfferings []CompetitorServiceOffering `json:"service_offerings"`
	PricingModel     CompetitorPricingModel      `json:"pricing_model"`
	TargetMarkets    []string                    `json:"target_markets"`
	RecentNews       []NewsItem                  `json:"recent_news"`
	FinancialMetrics *FinancialMetrics           `json:"financial_metrics,omitempty"`
	TechnologyStack  []string                    `json:"technology_stack"`
	Partnerships     []Partnership               `json:"partnerships"`
	LastUpdated      time.Time                   `json:"last_updated"`
}

// CompetitorType defines the type of competitor
type CompetitorType string

const (
	CompetitorTypeDirect     CompetitorType = "direct"
	CompetitorTypeIndirect   CompetitorType = "indirect"
	CompetitorTypeSubstitute CompetitorType = "substitute"
)

// CompetitorServiceOffering represents a service offered by a competitor
type CompetitorServiceOffering struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	Features     []string `json:"features"`
	PriceRange   string   `json:"price_range"`
	Availability []string `json:"availability"` // regions/markets
}

// CompetitorPricingModel represents a competitor's pricing approach
type CompetitorPricingModel struct {
	Type        string            `json:"type"` // "fixed", "hourly", "project", "retainer"
	Structure   string            `json:"structure"`
	Ranges      map[string]string `json:"ranges"` // service -> price range
	Discounts   []string          `json:"discounts"`
	LastUpdated time.Time         `json:"last_updated"`
}

// NewsItem represents recent news about a competitor
type NewsItem struct {
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Source      string    `json:"source"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	Impact      string    `json:"impact"` // "positive", "negative", "neutral"
}

// FinancialMetrics represents financial information about a competitor
type FinancialMetrics struct {
	Revenue      string    `json:"revenue"`
	Growth       string    `json:"growth"`
	Funding      string    `json:"funding"`
	Valuation    string    `json:"valuation"`
	LastReported time.Time `json:"last_reported"`
}

// Partnership represents a strategic partnership
type Partnership struct {
	PartnerName string `json:"partner_name"`
	Type        string `json:"type"` // "technology", "channel", "strategic"
	Description string `json:"description"`
	Impact      string `json:"impact"`
}

// CompetitiveMarketPosition represents market positioning analysis
type CompetitiveMarketPosition struct {
	OurPosition     Position            `json:"our_position"`
	CompetitorMap   map[string]Position `json:"competitor_map"`
	MarketSegments  []MarketSegment     `json:"market_segments"`
	OpportunityGaps []string            `json:"opportunity_gaps"`
}

// Position represents a market position
type Position struct {
	X           float64  `json:"x"` // e.g., price axis
	Y           float64  `json:"y"` // e.g., quality axis
	Quadrant    string   `json:"quadrant"`
	Description string   `json:"description"`
	Advantages  []string `json:"advantages"`
}

// MarketSegment represents a market segment
type MarketSegment struct {
	Name        string   `json:"name"`
	Size        string   `json:"size"`
	Growth      string   `json:"growth"`
	Leaders     []string `json:"leaders"`
	Opportunity string   `json:"opportunity"`
}

// CompetitiveMatrix represents a feature/capability comparison matrix
type CompetitiveMatrix struct {
	Criteria    []ComparisonCriteria          `json:"criteria"`
	Scores      map[string]map[string]float64 `json:"scores"`       // competitor -> criteria -> score
	Weights     map[string]float64            `json:"weights"`      // criteria -> weight
	TotalScores map[string]float64            `json:"total_scores"` // competitor -> weighted total
}

// ComparisonCriteria represents criteria for competitive comparison
type ComparisonCriteria struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Type        string  `json:"type"` // "capability", "price", "quality", "service"
}

// PricingIntelligence represents pricing intelligence data
type PricingIntelligence struct {
	ServiceType       string                           `json:"service_type"`
	Region            string                           `json:"region"`
	MarketRates       map[string]CompetitivePriceRange `json:"market_rates"` // competitor -> price range
	AverageRate       float64                          `json:"average_rate"`
	MedianRate        float64                          `json:"median_rate"`
	PriceDistribution []PriceBucket                    `json:"price_distribution"`
	Trends            *PricingTrends                   `json:"trends"`
	LastUpdated       time.Time                        `json:"last_updated"`
}

// CompetitivePriceRange represents a price range for a service
type CompetitivePriceRange struct {
	Min         float64   `json:"min"`
	Max         float64   `json:"max"`
	Currency    string    `json:"currency"`
	Unit        string    `json:"unit"`       // "hour", "project", "month"
	Confidence  float64   `json:"confidence"` // 0-1
	Source      string    `json:"source"`
	LastUpdated time.Time `json:"last_updated"`
}

// PriceBucket represents a price distribution bucket
type PriceBucket struct {
	Range      string  `json:"range"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// PricingTrends represents pricing trend analysis
type PricingTrends struct {
	Direction   string    `json:"direction"` // "increasing", "decreasing", "stable"
	Rate        float64   `json:"rate"`      // percentage change
	Factors     []string  `json:"factors"`   // factors influencing trends
	Forecast    string    `json:"forecast"`
	LastUpdated time.Time `json:"last_updated"`
}

// PricingComparison represents a pricing comparison across competitors
type PricingComparison struct {
	Services        []string                                    `json:"services"`
	Competitors     []string                                    `json:"competitors"`
	Matrix          map[string]map[string]CompetitivePriceRange `json:"matrix"` // service -> competitor -> price
	Analysis        *PricingAnalysis                            `json:"analysis"`
	Recommendations []CompetitivePricingRecommendation          `json:"recommendations"`
	GeneratedAt     time.Time                                   `json:"generated_at"`
}

// PricingAnalysis represents analysis of pricing comparison
type PricingAnalysis struct {
	LowestPrices  map[string]string `json:"lowest_prices"`  // service -> competitor
	HighestPrices map[string]string `json:"highest_prices"` // service -> competitor
	OurPosition   map[string]string `json:"our_position"`   // service -> position description
	Opportunities []string          `json:"opportunities"`
	Risks         []string          `json:"risks"`
}

// CompetitivePricingRecommendation represents a pricing recommendation
type CompetitivePricingRecommendation struct {
	Service    string  `json:"service"`
	Action     string  `json:"action"` // "increase", "decrease", "maintain"
	Rationale  string  `json:"rationale"`
	Impact     string  `json:"impact"`
	Confidence float64 `json:"confidence"`
}

// TechnologyTrends represents technology trend analysis
type TechnologyTrends struct {
	Category        string                `json:"category"`
	EmergingTrends  []TechnologyTrend     `json:"emerging_trends"`
	DecliningTrends []TechnologyTrend     `json:"declining_trends"`
	MarketImpact    *MarketImpact         `json:"market_impact"`
	Recommendations []TrendRecommendation `json:"recommendations"`
	LastUpdated     time.Time             `json:"last_updated"`
}

// TechnologyTrend represents a specific technology trend
type TechnologyTrend struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Stage            string   `json:"stage"` // "emerging", "growing", "mature", "declining"
	AdoptionRate     float64  `json:"adoption_rate"`
	KeyPlayers       []string `json:"key_players"`
	BusinessImpact   string   `json:"business_impact"`
	TimeToMainstream string   `json:"time_to_mainstream"`
	Relevance        float64  `json:"relevance"` // 0-1 relevance to our business
}

// MarketImpact represents the impact of trends on the market
type MarketImpact struct {
	DisruptionLevel string   `json:"disruption_level"` // "low", "medium", "high"
	AffectedSectors []string `json:"affected_sectors"`
	Opportunities   []string `json:"opportunities"`
	Threats         []string `json:"threats"`
	Timeline        string   `json:"timeline"`
}

// TrendRecommendation represents a recommendation based on trends
type TrendRecommendation struct {
	Trend      string  `json:"trend"`
	Action     string  `json:"action"`
	Priority   string  `json:"priority"` // "high", "medium", "low"
	Timeline   string  `json:"timeline"`
	Investment string  `json:"investment"`
	Rationale  string  `json:"rationale"`
	Confidence float64 `json:"confidence"`
}

// TrendImpact represents the impact of trends on a specific inquiry
type TrendImpact struct {
	InquiryID       string                 `json:"inquiry_id"`
	RelevantTrends  []TechnologyTrend      `json:"relevant_trends"`
	ImpactAnalysis  *ImpactAnalysis        `json:"impact_analysis"`
	Recommendations []ImpactRecommendation `json:"recommendations"`
	GeneratedAt     time.Time              `json:"generated_at"`
}

// ImpactAnalysis represents analysis of trend impact
type ImpactAnalysis struct {
	PositiveImpacts []string `json:"positive_impacts"`
	NegativeImpacts []string `json:"negative_impacts"`
	Opportunities   []string `json:"opportunities"`
	Risks           []string `json:"risks"`
	OverallImpact   string   `json:"overall_impact"` // "positive", "negative", "neutral"
}

// ImpactRecommendation represents a recommendation based on trend impact
type ImpactRecommendation struct {
	Trend      string  `json:"trend"`
	Action     string  `json:"action"`
	Rationale  string  `json:"rationale"`
	Priority   string  `json:"priority"`
	Timeline   string  `json:"timeline"`
	Confidence float64 `json:"confidence"`
}

// CompetitiveDifferentiationStrategy represents a competitive differentiation strategy
type CompetitiveDifferentiationStrategy struct {
	InquiryID            string                `json:"inquiry_id"`
	Competitors          []string              `json:"competitors"`
	OurStrengths         []Strength            `json:"our_strengths"`
	CompetitorWeaknesses []CompetitorWeakness  `json:"competitor_weaknesses"`
	DifferentiationAreas []DifferentiationArea `json:"differentiation_areas"`
	PositioningStrategy  *PositioningStrategy  `json:"positioning_strategy"`
	MessageFramework     *MessageFramework     `json:"message_framework"`
	GeneratedAt          time.Time             `json:"generated_at"`
}

// Strength represents our competitive strength
type Strength struct {
	Area        string   `json:"area"`
	Description string   `json:"description"`
	Evidence    []string `json:"evidence"`
	Impact      string   `json:"impact"`
	Uniqueness  float64  `json:"uniqueness"` // 0-1 how unique this strength is
}

// CompetitorWeakness represents a competitor's weakness
type CompetitorWeakness struct {
	Competitor  string   `json:"competitor"`
	Area        string   `json:"area"`
	Description string   `json:"description"`
	Evidence    []string `json:"evidence"`
	Severity    string   `json:"severity"` // "low", "medium", "high"
	Exploitable bool     `json:"exploitable"`
}

// DifferentiationArea represents an area where we can differentiate
type DifferentiationArea struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	OurAdvantage         string   `json:"our_advantage"`
	CompetitorGaps       []string `json:"competitor_gaps"`
	BusinessValue        string   `json:"business_value"`
	ImplementationEffort string   `json:"implementation_effort"`
	Priority             string   `json:"priority"`
}

// PositioningStrategy represents our market positioning strategy
type PositioningStrategy struct {
	PrimaryPosition   string   `json:"primary_position"`
	SecondaryPosition string   `json:"secondary_position"`
	TargetSegments    []string `json:"target_segments"`
	ValueProposition  string   `json:"value_proposition"`
	KeyMessages       []string `json:"key_messages"`
	ProofPoints       []string `json:"proof_points"`
}

// MessageFramework represents messaging framework for competitive positioning
type MessageFramework struct {
	CoreMessage         string            `json:"core_message"`
	AudienceMessages    map[string]string `json:"audience_messages"`    // audience -> message
	CompetitorResponses map[string]string `json:"competitor_responses"` // competitor -> response
	SupportingEvidence  []string          `json:"supporting_evidence"`
	CallToAction        string            `json:"call_to_action"`
}

// CompetitorWeaknesses represents analysis of competitor weaknesses
type CompetitorWeaknesses struct {
	Competitors []string                        `json:"competitors"`
	Analysis    map[string][]CompetitorWeakness `json:"analysis"` // competitor -> weaknesses
	Summary     *WeaknessAnalysisSummary        `json:"summary"`
	GeneratedAt time.Time                       `json:"generated_at"`
}

// WeaknessAnalysisSummary represents summary of weakness analysis
type WeaknessAnalysisSummary struct {
	CommonWeaknesses    []string `json:"common_weaknesses"`
	ExploitableGaps     []string `json:"exploitable_gaps"`
	MarketOpportunities []string `json:"market_opportunities"`
	RecommendedActions  []string `json:"recommended_actions"`
}
