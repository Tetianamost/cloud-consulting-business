-- Quality Assurance System Database Migration
-- This script creates tables for the quality assurance and validation system

-- Recommendation tracking table
CREATE TABLE IF NOT EXISTS recommendation_tracking (
    id VARCHAR(255) PRIMARY KEY,
    recommendation_id VARCHAR(255) NOT NULL UNIQUE,
    inquiry_id VARCHAR(255) NOT NULL,
    consultant_id VARCHAR(255) NOT NULL,
    recommendation_type VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    confidence DECIMAL(3,2) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    generated_at TIMESTAMP NOT NULL,
    model_version VARCHAR(100),
    prompt_version VARCHAR(100),
    context JSONB,
    tags JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Recommendation outcomes table
CREATE TABLE IF NOT EXISTS recommendation_outcomes (
    recommendation_id VARCHAR(255) PRIMARY KEY,
    outcome_type VARCHAR(50) NOT NULL,
    actual_result TEXT,
    client_feedback TEXT,
    consultant_notes TEXT,
    implementation_cost DECIMAL(12,2),
    time_to_implement VARCHAR(100),
    business_impact TEXT,
    technical_impact TEXT,
    lessons_learned JSONB,
    accuracy DECIMAL(3,2) CHECK (accuracy >= 0 AND accuracy <= 1),
    effectiveness DECIMAL(3,2) CHECK (effectiveness >= 0 AND effectiveness <= 1),
    metadata JSONB,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    recorded_by VARCHAR(255),
    FOREIGN KEY (recommendation_id) REFERENCES recommendation_tracking(recommendation_id)
);

-- Peer review requests table
CREATE TABLE IF NOT EXISTS peer_review_requests (
    id VARCHAR(255) PRIMARY KEY,
    recommendation_id VARCHAR(255) NOT NULL,
    requested_by VARCHAR(255) NOT NULL,
    assigned_to VARCHAR(255) NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium',
    review_type VARCHAR(50) NOT NULL,
    context TEXT,
    specific_areas JSONB,
    due_date TIMESTAMP,
    instructions TEXT,
    metadata JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recommendation_id) REFERENCES recommendation_tracking(recommendation_id)
);

-- Peer reviews table
CREATE TABLE IF NOT EXISTS peer_reviews (
    id VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    reviewer_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    time_spent INTEGER DEFAULT 0, -- minutes
    feedback JSONB,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES peer_review_requests(id)
);

-- Client outcomes table
CREATE TABLE IF NOT EXISTS client_outcomes (
    id VARCHAR(255) PRIMARY KEY,
    inquiry_id VARCHAR(255) NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    recommendation_ids JSONB,
    engagement_type VARCHAR(50) NOT NULL,
    outcome_type VARCHAR(50) NOT NULL,
    client_satisfaction INTEGER CHECK (client_satisfaction >= 1 AND client_satisfaction <= 10),
    business_value DECIMAL(12,2),
    cost_savings DECIMAL(12,2),
    time_to_value INTEGER, -- days
    implementation_rate DECIMAL(3,2) CHECK (implementation_rate >= 0 AND implementation_rate <= 1),
    client_testimonial TEXT,
    challenges_faced JSONB,
    success_factors JSONB,
    lessons_learned JSONB,
    follow_up_opportunities JSONB,
    reference_permission BOOLEAN DEFAULT FALSE,
    case_study_permission BOOLEAN DEFAULT FALSE,
    net_promoter_score INTEGER CHECK (net_promoter_score >= -100 AND net_promoter_score <= 100),
    metadata JSONB,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    recorded_by VARCHAR(255)
);

-- Quality validations table
CREATE TABLE IF NOT EXISTS quality_validations (
    id VARCHAR(255) PRIMARY KEY,
    recommendation_id VARCHAR(255) NOT NULL,
    overall_score DECIMAL(3,2) NOT NULL CHECK (overall_score >= 0 AND overall_score <= 1),
    validation_checks JSONB NOT NULL,
    passed_checks INTEGER NOT NULL,
    total_checks INTEGER NOT NULL,
    pass_rate DECIMAL(3,2) NOT NULL CHECK (pass_rate >= 0 AND pass_rate <= 1),
    issues JSONB,
    recommendations JSONB,
    requires_review BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    validated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    validated_by VARCHAR(255),
    FOREIGN KEY (recommendation_id) REFERENCES recommendation_tracking(recommendation_id)
);

-- Quality alerts table
CREATE TABLE IF NOT EXISTS quality_alerts (
    id VARCHAR(255) PRIMARY KEY,
    alert_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    metric_value DECIMAL(10,4) NOT NULL,
    threshold_value DECIMAL(10,4),
    description TEXT NOT NULL,
    context JSONB,
    recipients JSONB,
    actions JSONB,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quality thresholds configuration table
CREATE TABLE IF NOT EXISTS quality_thresholds (
    id VARCHAR(255) PRIMARY KEY DEFAULT 'default',
    min_accuracy DECIMAL(3,2) NOT NULL,
    min_client_satisfaction INTEGER NOT NULL,
    min_peer_review_score DECIMAL(3,2) NOT NULL,
    max_response_time INTEGER NOT NULL, -- milliseconds
    min_implementation_rate DECIMAL(3,2) NOT NULL,
    alert_thresholds JSONB,
    updated_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quality standards configuration table
CREATE TABLE IF NOT EXISTS quality_standards (
    id VARCHAR(255) PRIMARY KEY DEFAULT 'default',
    min_overall_score DECIMAL(3,2) NOT NULL,
    component_weights JSONB NOT NULL,
    validation_rules JSONB,
    review_thresholds JSONB,
    alert_thresholds JSONB,
    quality_grade_ranges JSONB,
    auto_review_triggers JSONB,
    metadata JSONB,
    updated_by VARCHAR(255),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_recommendation_tracking_inquiry_id ON recommendation_tracking(inquiry_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_tracking_consultant_id ON recommendation_tracking(consultant_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_tracking_type ON recommendation_tracking(recommendation_type);
CREATE INDEX IF NOT EXISTS idx_recommendation_tracking_status ON recommendation_tracking(status);
CREATE INDEX IF NOT EXISTS idx_recommendation_tracking_created_at ON recommendation_tracking(created_at);

CREATE INDEX IF NOT EXISTS idx_recommendation_outcomes_outcome_type ON recommendation_outcomes(outcome_type);
CREATE INDEX IF NOT EXISTS idx_recommendation_outcomes_recorded_at ON recommendation_outcomes(recorded_at);

CREATE INDEX IF NOT EXISTS idx_peer_review_requests_assigned_to ON peer_review_requests(assigned_to);
CREATE INDEX IF NOT EXISTS idx_peer_review_requests_status ON peer_review_requests(status);
CREATE INDEX IF NOT EXISTS idx_peer_review_requests_priority ON peer_review_requests(priority);
CREATE INDEX IF NOT EXISTS idx_peer_review_requests_due_date ON peer_review_requests(due_date);

CREATE INDEX IF NOT EXISTS idx_peer_reviews_reviewer_id ON peer_reviews(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_peer_reviews_status ON peer_reviews(status);
CREATE INDEX IF NOT EXISTS idx_peer_reviews_completed_at ON peer_reviews(completed_at);

CREATE INDEX IF NOT EXISTS idx_client_outcomes_inquiry_id ON client_outcomes(inquiry_id);
CREATE INDEX IF NOT EXISTS idx_client_outcomes_client_name ON client_outcomes(client_name);
CREATE INDEX IF NOT EXISTS idx_client_outcomes_outcome_type ON client_outcomes(outcome_type);
CREATE INDEX IF NOT EXISTS idx_client_outcomes_engagement_type ON client_outcomes(engagement_type);
CREATE INDEX IF NOT EXISTS idx_client_outcomes_recorded_at ON client_outcomes(recorded_at);

CREATE INDEX IF NOT EXISTS idx_quality_validations_recommendation_id ON quality_validations(recommendation_id);
CREATE INDEX IF NOT EXISTS idx_quality_validations_overall_score ON quality_validations(overall_score);
CREATE INDEX IF NOT EXISTS idx_quality_validations_requires_review ON quality_validations(requires_review);
CREATE INDEX IF NOT EXISTS idx_quality_validations_validated_at ON quality_validations(validated_at);

CREATE INDEX IF NOT EXISTS idx_quality_alerts_status ON quality_alerts(status);
CREATE INDEX IF NOT EXISTS idx_quality_alerts_severity ON quality_alerts(severity);
CREATE INDEX IF NOT EXISTS idx_quality_alerts_metric_type ON quality_alerts(metric_type);
CREATE INDEX IF NOT EXISTS idx_quality_alerts_created_at ON quality_alerts(created_at);

-- Insert default quality thresholds
INSERT INTO quality_thresholds (
    id, min_accuracy, min_client_satisfaction, min_peer_review_score,
    max_response_time, min_implementation_rate, alert_thresholds,
    updated_by, updated_at
) VALUES (
    'default', 0.80, 7, 0.75, 5000, 0.70,
    '{"accuracy_drop": 0.1, "satisfaction_drop": 2.0, "review_score_drop": 0.15}',
    'system', CURRENT_TIMESTAMP
) ON CONFLICT (id) DO NOTHING;

-- Insert default quality standards
INSERT INTO quality_standards (
    id, min_overall_score, component_weights, validation_rules,
    review_thresholds, alert_thresholds, quality_grade_ranges,
    auto_review_triggers, updated_by, updated_at
) VALUES (
    'default', 0.75,
    '{"technical_accuracy": 0.25, "business_relevance": 0.20, "completeness": 0.15, "clarity": 0.15, "actionability": 0.15, "innovation": 0.10}',
    '[]',
    '{"low_score": 0.6, "peer_review_required": 0.7}',
    '{"critical": 0.5, "high": 0.6, "medium": 0.7}',
    '{"A": {"min": 0.9, "max": 1.0}, "B": {"min": 0.8, "max": 0.89}, "C": {"min": 0.7, "max": 0.79}, "D": {"min": 0.6, "max": 0.69}, "F": {"min": 0.0, "max": 0.59}}',
    '[]',
    'system', CURRENT_TIMESTAMP
) ON CONFLICT (id) DO NOTHING;

-- Add triggers to update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_recommendation_tracking_updated_at BEFORE UPDATE ON recommendation_tracking FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_peer_review_requests_updated_at BEFORE UPDATE ON peer_review_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_peer_reviews_updated_at BEFORE UPDATE ON peer_reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_quality_alerts_updated_at BEFORE UPDATE ON quality_alerts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_quality_thresholds_updated_at BEFORE UPDATE ON quality_thresholds FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_quality_standards_updated_at BEFORE UPDATE ON quality_standards FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();