package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// RiskAssessor defines the interface for risk assessment functionality
type RiskAssessor interface {
	AssessRisks(ctx context.Context, inquiry *domain.Inquiry, solution *ProposedSolution) (*RiskAssessment, error)
	IdentifySecurityRisks(ctx context.Context, architecture *Architecture) ([]*SecurityRisk, error)
	EvaluateComplianceRisks(ctx context.Context, industry string, solution *ProposedSolution) ([]*ComplianceRisk, error)
	GenerateMitigationStrategies(ctx context.Context, risks []*Risk) ([]*MitigationStrategy, error)
}

// RiskAssessment represents a comprehensive risk assessment
type RiskAssessment struct {
	ID                   string                `json:"id"`
	InquiryID            string                `json:"inquiry_id"`
	OverallRiskLevel     string                `json:"overall_risk_level"` // "low", "medium", "high", "critical"
	TechnicalRisks       []*TechnicalRisk      `json:"technical_risks"`
	SecurityRisks        []*SecurityRisk       `json:"security_risks"`
	ComplianceRisks      []*ComplianceRisk     `json:"compliance_risks"`
	BusinessRisks        []*BusinessRisk       `json:"business_risks"`
	MitigationStrategies []*MitigationStrategy `json:"mitigation_strategies"`
	RecommendedActions   []string              `json:"recommended_actions"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// Risk represents a base risk structure
type Risk struct {
	ID                 string   `json:"id"`
	Category           string   `json:"category"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Impact             string   `json:"impact"`      // "low", "medium", "high", "critical"
	Probability        string   `json:"probability"` // "low", "medium", "high"
	RiskScore          int      `json:"risk_score"`  // calculated from impact and probability
	AffectedComponents []string `json:"affected_components"`
	DocumentationURL   string   `json:"documentation_url,omitempty"`
}

// TechnicalRisk represents technical risks in cloud solutions
type TechnicalRisk struct {
	Risk
	ServiceType        string   `json:"service_type"`
	ArchitecturalLayer string   `json:"architectural_layer"`
	Dependencies       []string `json:"dependencies"`
	PerformanceImpact  string   `json:"performance_impact"`
	ScalabilityImpact  string   `json:"scalability_impact"`
}

// SecurityRisk represents security-related risks
type SecurityRisk struct {
	Risk
	ThreatType           string   `json:"threat_type"`
	AttackVectors        []string `json:"attack_vectors"`
	DataClassification   string   `json:"data_classification"`
	ComplianceFrameworks []string `json:"compliance_frameworks"`
	EncryptionRequired   bool     `json:"encryption_required"`
}

// ComplianceRisk represents compliance and regulatory risks
type ComplianceRisk struct {
	Risk
	Framework         string   `json:"framework"` // "HIPAA", "PCI-DSS", "SOX", "GDPR", etc.
	RequirementID     string   `json:"requirement_id"`
	Jurisdiction      string   `json:"jurisdiction"`
	PenaltyLevel      string   `json:"penalty_level"`
	AuditRequirements []string `json:"audit_requirements"`
}

// BusinessRisk represents business and operational risks
type BusinessRisk struct {
	Risk
	BusinessFunction  string `json:"business_function"`
	RevenueImpact     string `json:"revenue_impact"`
	CustomerImpact    string `json:"customer_impact"`
	OperationalImpact string `json:"operational_impact"`
	TimeToRecover     string `json:"time_to_recover"`
}

// MitigationStrategy represents a strategy to mitigate identified risks
type MitigationStrategy struct {
	ID                     string   `json:"id"`
	RiskID                 string   `json:"risk_id"`
	Strategy               string   `json:"strategy"`
	ImplementationSteps    []string `json:"implementation_steps"`
	EstimatedEffort        string   `json:"estimated_effort"`
	Cost                   string   `json:"cost"`
	Priority               string   `json:"priority"`
	Effectiveness          string   `json:"effectiveness"` // "low", "medium", "high"
	DocumentationURL       string   `json:"documentation_url,omitempty"`
	CloudProvider          string   `json:"cloud_provider,omitempty"`
	ServiceRecommendations []string `json:"service_recommendations,omitempty"`
}

// Note: ComplianceRequirement is defined in knowledge.go

// ProposedSolution represents a proposed cloud solution for risk assessment
type ProposedSolution struct {
	ID               string                  `json:"id"`
	InquiryID        string                  `json:"inquiry_id"`
	CloudProviders   []string                `json:"cloud_providers"`
	Services         []CloudService          `json:"services"`
	Architecture     *Architecture           `json:"architecture"`
	DataFlow         []DataFlowComponent     `json:"data_flow"`
	SecurityControls []SecurityControl       `json:"security_controls"`
	ComplianceNeeds  []ComplianceRequirement `json:"compliance_needs"`
	EstimatedCost    string                  `json:"estimated_cost"`
	Timeline         string                  `json:"timeline"`
}

// Architecture represents the technical architecture of a solution
type Architecture struct {
	ID               string                  `json:"id"`
	Type             string                  `json:"type"` // "microservices", "monolithic", "serverless", etc.
	Components       []ArchitectureComponent `json:"components"`
	NetworkTopology  NetworkTopology         `json:"network_topology"`
	DataStorage      []DataStorageComponent  `json:"data_storage"`
	SecurityLayers   []SecurityLayer         `json:"security_layers"`
	HighAvailability bool                    `json:"high_availability"`
	DisasterRecovery bool                    `json:"disaster_recovery"`
}

// CloudService represents a cloud service in the solution
type CloudService struct {
	Provider      string                 `json:"provider"`
	ServiceName   string                 `json:"service_name"`
	ServiceType   string                 `json:"service_type"`
	Configuration map[string]interface{} `json:"configuration"`
	Dependencies  []string               `json:"dependencies"`
	CriticalPath  bool                   `json:"critical_path"`
}

// Note: ArchitectureComponent and NetworkTopology are defined in architecture.go

// DataStorageComponent represents data storage configuration
type DataStorageComponent struct {
	Type             string                 `json:"type"`
	Provider         string                 `json:"provider"`
	ServiceName      string                 `json:"service_name"`
	DataType         string                 `json:"data_type"`
	SensitivityLevel string                 `json:"sensitivity_level"`
	BackupStrategy   string                 `json:"backup_strategy"`
	Configuration    map[string]interface{} `json:"configuration"`
}

// SecurityLayer represents a security layer in the architecture
type SecurityLayer struct {
	Layer    string            `json:"layer"`
	Controls []SecurityControl `json:"controls"`
	Tools    []string          `json:"tools"`
	Policies []string          `json:"policies"`
}

// Note: SecurityControl, SecurityGroup, and FirewallRule are defined in architecture.go

// Note: LoadBalancer and HealthCheck are defined in architecture.go

// DataFlowComponent represents a component in the data flow
type DataFlowComponent struct {
	Source        string                 `json:"source"`
	Destination   string                 `json:"destination"`
	DataType      string                 `json:"data_type"`
	Volume        string                 `json:"volume"`
	Frequency     string                 `json:"frequency"`
	Encryption    bool                   `json:"encryption"`
	Compliance    []string               `json:"compliance"`
	Configuration map[string]interface{} `json:"configuration"`
}
