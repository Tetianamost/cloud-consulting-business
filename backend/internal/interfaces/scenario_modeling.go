package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// ScenarioModelingEngine defines the interface for advanced scenario modeling
type ScenarioModelingEngine interface {
	// Comprehensive scenario analysis
	GenerateComprehensiveScenarioAnalysis(ctx context.Context, inquiry *domain.Inquiry) (*ComprehensiveScenarioAnalysis, error)
}

// ComprehensiveScenarioAnalysis represents comprehensive scenario analysis
type ComprehensiveScenarioAnalysis struct {
	ID                  string                      `json:"id"`
	InquiryID           string                      `json:"inquiry_id"`
	AnalysisDate        time.Time                   `json:"analysis_date"`
	BaseScenario        *ScenarioBaseScenario       `json:"base_scenario"`
	WhatIfScenarios     []*WhatIfScenario           `json:"what_if_scenarios"`
	MultiYearProjection *MultiYearProjection        `json:"multi_year_projection"`
	DRScenarios         []*DisasterRecoveryScenario `json:"dr_scenarios"`
	CapacityScenarios   []*CapacityScenario         `json:"capacity_scenarios"`
	IntegratedAnalysis  *IntegratedAnalysis         `json:"integrated_analysis"`
	ExecutiveSummary    *ExecutiveSummary           `json:"executive_summary"`
	KeyRecommendations  []*KeyRecommendation        `json:"key_recommendations"`
	NextSteps           []*ScenarioNextStep         `json:"next_steps"`
	CreatedAt           time.Time                   `json:"created_at"`
}

// ScenarioBaseScenario represents the baseline scenario for what-if analysis
type ScenarioBaseScenario struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// WhatIfScenario represents a what-if scenario analysis
type WhatIfScenario struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	ConfidenceLevel   float64            `json:"confidence_level"`
	RiskFactors       []string           `json:"risk_factors"`
	ProjectedOutcomes *ProjectedOutcomes `json:"projected_outcomes"`
	CreatedAt         time.Time          `json:"created_at"`
}

// ProjectedOutcomes represents projected outcomes of a scenario
type ProjectedOutcomes struct {
	CostProjection *CostProjection `json:"cost_projection"`
}

// CostProjection represents cost projections
type CostProjection struct {
	TotalCostChange   float64 `json:"total_cost_change"`
	CostChangePercent float64 `json:"cost_change_percent"`
}

// MultiYearProjection represents multi-year projection results
type MultiYearProjection struct {
	ID                string               `json:"id"`
	ProjectionPeriod  string               `json:"projection_period"`
	YearlyProjections []*YearlyProjection  `json:"yearly_projections"`
	KeyInsights       []*ProjectionInsight `json:"key_insights"`
	CreatedAt         time.Time            `json:"created_at"`
}

// YearlyProjection represents projection for a specific year
type YearlyProjection struct {
	Year          int     `json:"year"`
	ProjectedCost float64 `json:"projected_cost"`
}

// ProjectionInsight represents an insight from projections
type ProjectionInsight struct {
	InsightType string `json:"insight_type"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
}

// DisasterRecoveryScenario represents a disaster recovery scenario
type DisasterRecoveryScenario struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CapacityScenario represents a capacity planning scenario
type CapacityScenario struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// IntegratedAnalysis represents integrated analysis across all scenarios
type IntegratedAnalysis struct {
	OverallRiskLevel          string  `json:"overall_risk_level"`
	CostOptimizationPotential float64 `json:"cost_optimization_potential"`
	BusinessImpactSummary     string  `json:"business_impact_summary"`
	TechnicalComplexity       string  `json:"technical_complexity"`
	ImplementationTimeline    string  `json:"implementation_timeline"`
}

// ExecutiveSummary represents executive summary
type ExecutiveSummary struct {
	KeyFindings           []string `json:"key_findings"`
	CostImpactSummary     string   `json:"cost_impact_summary"`
	BusinessImpactSummary string   `json:"business_impact_summary"`
	RiskSummary           string   `json:"risk_summary"`
	RecommendedApproach   string   `json:"recommended_approach"`
	ExpectedOutcomes      []string `json:"expected_outcomes"`
	SuccessMetrics        []string `json:"success_metrics"`
}

// KeyRecommendation represents a key recommendation
type KeyRecommendation struct {
	RecommendationID   string  `json:"recommendation_id"`
	Title              string  `json:"title"`
	Priority           string  `json:"priority"`
	Category           string  `json:"category"`
	Description        string  `json:"description"`
	ExpectedBenefit    string  `json:"expected_benefit"`
	ImplementationCost float64 `json:"implementation_cost"`
	Timeline           string  `json:"timeline"`
	RiskLevel          string  `json:"risk_level"`
}

// ScenarioNextStep represents a next step
type ScenarioNextStep struct {
	StepID      string `json:"step_id"`
	StepName    string `json:"step_name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Timeline    string `json:"timeline"`
	Priority    string `json:"priority"`
}
