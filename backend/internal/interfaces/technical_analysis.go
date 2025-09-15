package interfaces

import (
	"context"
	"time"
)

// TechnicalAnalysisService defines the interface for technical deep-dive analysis tools
type TechnicalAnalysisService interface {
	// Code review and architecture assessment
	AnalyzeCodebase(ctx context.Context, request *CodeAnalysisRequest) (*CodeAnalysisResult, error)
	AssessArchitecture(ctx context.Context, request *TechArchitectureAssessmentRequest) (*TechArchitectureAssessmentResult, error)

	// Security vulnerability assessment
	PerformSecurityAssessment(ctx context.Context, request *TechSecurityAssessmentRequest) (*TechSecurityAssessmentResult, error)
	GenerateSecurityRemediation(ctx context.Context, vulnerabilities []*TechSecurityVulnerability) ([]*RemediationRecommendation, error)

	// Performance benchmarking and optimization
	AnalyzePerformance(ctx context.Context, request *TechPerformanceAnalysisRequest) (*TechPerformanceAnalysisResult, error)
	GenerateOptimizationRecommendations(ctx context.Context, metrics *TechPerformanceMetrics) ([]*TechOptimizationRecommendation, error)

	// Compliance gap analysis
	AnalyzeCompliance(ctx context.Context, request *TechComplianceAnalysisRequest) (*TechComplianceAnalysisResult, error)
	GenerateComplianceRemediation(ctx context.Context, gaps []*TechComplianceGap) ([]*TechComplianceRemediation, error)

	// Comprehensive technical analysis
	PerformComprehensiveAnalysis(ctx context.Context, request *ComprehensiveAnalysisRequest) (*ComprehensiveAnalysisResult, error)
}

// CodeAnalysisRequest represents a request for code analysis
type CodeAnalysisRequest struct {
	InquiryID            string                 `json:"inquiry_id"`
	RepositoryURL        string                 `json:"repository_url,omitempty"`
	CodeSamples          []*CodeSample          `json:"code_samples,omitempty"`
	Languages            []string               `json:"languages"`
	AnalysisScope        []string               `json:"analysis_scope"` // "security", "performance", "maintainability", "architecture"
	CloudProvider        string                 `json:"cloud_provider"`
	ApplicationType      string                 `json:"application_type"` // "web", "api", "microservices", "batch", "data-pipeline"
	ComplianceFrameworks []string               `json:"compliance_frameworks,omitempty"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// CodeSample represents a code sample for analysis
type CodeSample struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Language    string `json:"language"`
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
	LineCount   int    `json:"line_count"`
}

// CodeAnalysisResult represents the result of code analysis
type CodeAnalysisResult struct {
	ID                      string                     `json:"id"`
	InquiryID               string                     `json:"inquiry_id"`
	OverallScore            float64                    `json:"overall_score"` // 0-100
	SecurityFindings        []*TechSecurityFinding     `json:"security_findings"`
	PerformanceFindings     []*TechPerformanceFinding  `json:"performance_findings"`
	ArchitectureFindings    []*TechArchitectureFinding `json:"architecture_findings"`
	MaintainabilityFindings []*MaintainabilityFinding  `json:"maintainability_findings"`
	BestPracticeViolations  []*BestPracticeViolation   `json:"best_practice_violations"`
	CloudOptimizations      []*CloudOptimization       `json:"cloud_optimizations"`
	Summary                 string                     `json:"summary"`
	Recommendations         []*CodeRecommendation      `json:"recommendations"`
	GeneratedAt             time.Time                  `json:"generated_at"`
}

// TechSecurityFinding represents a security-related finding in code
type TechSecurityFinding struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`     // "vulnerability", "misconfiguration", "exposure"
	Severity    string   `json:"severity"` // "critical", "high", "medium", "low"
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
	CWE         string   `json:"cwe,omitempty"`  // Common Weakness Enumeration
	CVSS        float64  `json:"cvss,omitempty"` // Common Vulnerability Scoring System
	Impact      string   `json:"impact"`
	Remediation string   `json:"remediation"`
	References  []string `json:"references"`
}

// TechPerformanceFinding represents a performance-related finding
type TechPerformanceFinding struct {
	ID                   string   `json:"id"`
	Type                 string   `json:"type"` // "bottleneck", "inefficiency", "resource-waste"
	Severity             string   `json:"severity"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	Location             Location `json:"location"`
	Impact               string   `json:"impact"`
	Suggestion           string   `json:"suggestion"`
	EstimatedImprovement string   `json:"estimated_improvement"`
}

// TechArchitectureFinding represents an architecture-related finding
type TechArchitectureFinding struct {
	ID                string   `json:"id"`
	Type              string   `json:"type"` // "anti-pattern", "coupling", "cohesion", "scalability"
	Severity          string   `json:"severity"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Location          Location `json:"location"`
	Impact            string   `json:"impact"`
	Recommendation    string   `json:"recommendation"`
	RefactoringEffort string   `json:"refactoring_effort"` // "low", "medium", "high"
}

// MaintainabilityFinding represents a maintainability-related finding
type MaintainabilityFinding struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"` // "complexity", "duplication", "documentation", "testing"
	Severity    string             `json:"severity"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Location    Location           `json:"location"`
	Metrics     map[string]float64 `json:"metrics"` // cyclomatic complexity, etc.
	Suggestion  string             `json:"suggestion"`
}

// BestPracticeViolation represents a violation of best practices
type BestPracticeViolation struct {
	ID          string   `json:"id"`
	Practice    string   `json:"practice"`
	Category    string   `json:"category"` // "cloud", "security", "performance", "maintainability"
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	Location    Location `json:"location"`
	Correction  string   `json:"correction"`
	References  []string `json:"references"`
}

// CloudOptimization represents cloud-specific optimization opportunities
type CloudOptimization struct {
	ID                   string  `json:"id"`
	Type                 string  `json:"type"` // "cost", "performance", "scalability", "reliability"
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	CurrentState         string  `json:"current_state"`
	RecommendedState     string  `json:"recommended_state"`
	EstimatedSavings     float64 `json:"estimated_savings,omitempty"`
	ImplementationEffort string  `json:"implementation_effort"`
	Priority             string  `json:"priority"`
}

// CodeRecommendation represents a code improvement recommendation
type CodeRecommendation struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`     // "refactor", "optimize", "secure", "modernize"
	Priority     string   `json:"priority"` // "critical", "high", "medium", "low"
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Effort       string   `json:"effort"` // "low", "medium", "high"
	Timeline     string   `json:"timeline"`
	Dependencies []string `json:"dependencies"`
	References   []string `json:"references"`
}

// Location represents a location in code
type Location struct {
	File     string `json:"file"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Function string `json:"function,omitempty"`
	Class    string `json:"class,omitempty"`
	Module   string `json:"module,omitempty"`
}

// TechArchitectureAssessmentRequest represents a request for architecture assessment
type TechArchitectureAssessmentRequest struct {
	InquiryID                   string                       `json:"inquiry_id"`
	ArchitectureDiagrams        []*ArchitectureDiagram       `json:"architecture_diagrams,omitempty"`
	SystemDescription           string                       `json:"system_description"`
	CloudProvider               string                       `json:"cloud_provider"`
	Services                    []string                     `json:"services"` // AWS services, Azure services, etc.
	SystemTrafficPatterns       *SystemTrafficPatterns       `json:"traffic_patterns,omitempty"`
	DataFlow                    *DataFlow                    `json:"data_flow,omitempty"`
	SecurityRequirements        []string                     `json:"security_requirements"`
	ComplianceFrameworks        []string                     `json:"compliance_frameworks"`
	TechPerformanceRequirements *TechPerformanceRequirements `json:"performance_requirements,omitempty"`
	ScalabilityRequirements     *ScalabilityRequirements     `json:"scalability_requirements,omitempty"`
	Metadata                    map[string]interface{}       `json:"metadata"`
}

// ArchitectureDiagram represents an architecture diagram
type ArchitectureDiagram struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`   // "high-level", "detailed", "network", "security", "data-flow"
	Format      string `json:"format"` // "image", "text", "json", "yaml"
	Content     string `json:"content"`
	Description string `json:"description"`
}

// SystemTrafficPatterns represents traffic patterns in the system
type SystemTrafficPatterns struct {
	PeakTraffic            int64              `json:"peak_traffic"` // requests per second
	AverageTraffic         int64              `json:"average_traffic"`
	TrafficSources         []string           `json:"traffic_sources"`
	GeographicDistribution []string           `json:"geographic_distribution"`
	SeasonalPatterns       map[string]float64 `json:"seasonal_patterns"`
	GrowthProjection       float64            `json:"growth_projection"` // percentage per year
}

// DataFlow represents data flow in the system
type DataFlow struct {
	DataSources            []string          `json:"data_sources"`
	DataDestinations       []string          `json:"data_destinations"`
	DataVolume             int64             `json:"data_volume"` // bytes per day
	DataTypes              []string          `json:"data_types"`
	ProcessingRequirements []string          `json:"processing_requirements"`
	RetentionRequirements  map[string]string `json:"retention_requirements"`
}

// TechPerformanceRequirements represents performance requirements
type TechPerformanceRequirements struct {
	ResponseTime int     `json:"response_time"` // milliseconds
	Throughput   int64   `json:"throughput"`    // requests per second
	Availability float64 `json:"availability"`  // percentage
	Latency      int     `json:"latency"`       // milliseconds
	Concurrency  int     `json:"concurrency"`   // concurrent users
}

// ScalabilityRequirements represents scalability requirements
type ScalabilityRequirements struct {
	HorizontalScaling bool     `json:"horizontal_scaling"`
	VerticalScaling   bool     `json:"vertical_scaling"`
	AutoScaling       bool     `json:"auto_scaling"`
	MaxInstances      int      `json:"max_instances"`
	ScalingTriggers   []string `json:"scaling_triggers"`
	LoadBalancing     string   `json:"load_balancing"`
}

// TechArchitectureAssessmentResult represents the result of architecture assessment
type TechArchitectureAssessmentResult struct {
	ID                        string                            `json:"id"`
	InquiryID                 string                            `json:"inquiry_id"`
	OverallScore              float64                           `json:"overall_score"` // 0-100
	TechSecurityAssessment    *TechSecurityAssessment           `json:"security_assessment"`
	TechPerformanceAssessment *TechPerformanceAssessment        `json:"performance_assessment"`
	TechScalabilityAssessment *TechScalabilityAssessment        `json:"scalability_assessment"`
	TechReliabilityAssessment *TechReliabilityAssessment        `json:"reliability_assessment"`
	TechCostAssessment        *TechCostAssessment               `json:"cost_assessment"`
	TechComplianceAssessment  *TechComplianceAssessment         `json:"compliance_assessment"`
	ArchitectureFindings      []*TechArchitectureFinding        `json:"architecture_findings"`
	Recommendations           []*TechArchitectureRecommendation `json:"recommendations"`
	Summary                   string                            `json:"summary"`
	GeneratedAt               time.Time                         `json:"generated_at"`
}

// TechSecurityAssessment represents security assessment results
type TechSecurityAssessment struct {
	Score            float64                `json:"score"` // 0-100
	SecurityFindings []*TechSecurityFinding `json:"security_findings"`
	ComplianceGaps   []*TechComplianceGap   `json:"compliance_gaps"`
	Recommendations  []string               `json:"recommendations"`
}

// TechPerformanceAssessment represents performance assessment results
type TechPerformanceAssessment struct {
	Score                     float64                           `json:"score"` // 0-100
	BottleneckAnalysis        []*TechPerformanceBottleneck      `json:"bottleneck_analysis"`
	OptimizationOpportunities []*TechOptimizationRecommendation `json:"optimization_opportunities"`
	Recommendations           []string                          `json:"recommendations"`
}

// TechScalabilityAssessment represents scalability assessment results
type TechScalabilityAssessment struct {
	Score           float64  `json:"score"` // 0-100
	ScalabilityGaps []string `json:"scalability_gaps"`
	Recommendations []string `json:"recommendations"`
}

// TechReliabilityAssessment represents reliability assessment results
type TechReliabilityAssessment struct {
	Score                 float64  `json:"score"` // 0-100
	SinglePointsOfFailure []string `json:"single_points_of_failure"`
	DisasterRecoveryGaps  []string `json:"disaster_recovery_gaps"`
	Recommendations       []string `json:"recommendations"`
}

// TechCostAssessment represents cost assessment results
type TechCostAssessment struct {
	EstimatedMonthlyCost float64                 `json:"estimated_monthly_cost"`
	CostOptimizations    []*TechCostOptimization `json:"cost_optimizations"`
	Recommendations      []string                `json:"recommendations"`
}

// TechComplianceAssessment represents compliance assessment results
type TechComplianceAssessment struct {
	Score           float64              `json:"score"` // 0-100
	ComplianceGaps  []*TechComplianceGap `json:"compliance_gaps"`
	Recommendations []string             `json:"recommendations"`
}

// TechArchitectureRecommendation represents an architecture improvement recommendation
type TechArchitectureRecommendation struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`     // "security", "performance", "cost", "scalability", "reliability"
	Priority     string   `json:"priority"` // "critical", "high", "medium", "low"
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Effort       string   `json:"effort"` // "low", "medium", "high"
	Timeline     string   `json:"timeline"`
	Cost         float64  `json:"cost,omitempty"`
	Dependencies []string `json:"dependencies"`
	References   []string `json:"references"`
}

// TechSecurityAssessmentRequest represents a request for security assessment
type TechSecurityAssessmentRequest struct {
	InquiryID            string                 `json:"inquiry_id"`
	SystemDescription    string                 `json:"system_description"`
	CloudProvider        string                 `json:"cloud_provider"`
	Services             []string               `json:"services"`
	DataClassification   []string               `json:"data_classification"` // "public", "internal", "confidential", "restricted"
	ComplianceFrameworks []string               `json:"compliance_frameworks"`
	TechThreatModel      *TechThreatModel       `json:"threat_model,omitempty"`
	SecurityControls     []*TechSecurityControl `json:"security_controls,omitempty"`
	NetworkArchitecture  *NetworkArchitecture   `json:"network_architecture,omitempty"`
	AccessControls       *AccessControls        `json:"access_controls,omitempty"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// TechThreatModel represents a threat model for the system
type TechThreatModel struct {
	Assets        []string `json:"assets"`
	Threats       []string `json:"threats"`
	Attackers     []string `json:"attackers"`
	AttackVectors []string `json:"attack_vectors"`
	Mitigations   []string `json:"mitigations"`
}

// TechSecurityControl represents a security control
type TechSecurityControl struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`     // "preventive", "detective", "corrective"
	Category      string `json:"category"` // "access", "network", "data", "application"
	Description   string `json:"description"`
	Status        string `json:"status"`        // "implemented", "partial", "missing"
	Effectiveness string `json:"effectiveness"` // "high", "medium", "low"
}

// NetworkArchitecture represents network architecture details
type NetworkArchitecture struct {
	VPCs              []string `json:"vpcs"`
	Subnets           []string `json:"subnets"`
	SecurityGroups    []string `json:"security_groups"`
	NACLs             []string `json:"nacls"`
	LoadBalancers     []string `json:"load_balancers"`
	Firewalls         []string `json:"firewalls"`
	VPNConnections    []string `json:"vpn_connections"`
	DirectConnections []string `json:"direct_connections"`
}

// AccessControls represents access control configuration
type AccessControls struct {
	IdentityProvider      string             `json:"identity_provider"`
	AuthenticationMethods []string           `json:"authentication_methods"`
	AuthorizationModel    string             `json:"authorization_model"` // "RBAC", "ABAC", "ACL"
	MFAEnabled            bool               `json:"mfa_enabled"`
	PasswordPolicy        *PasswordPolicy    `json:"password_policy,omitempty"`
	SessionManagement     *SessionManagement `json:"session_management,omitempty"`
}

// PasswordPolicy represents password policy configuration
type PasswordPolicy struct {
	MinLength        int  `json:"min_length"`
	RequireUppercase bool `json:"require_uppercase"`
	RequireLowercase bool `json:"require_lowercase"`
	RequireNumbers   bool `json:"require_numbers"`
	RequireSymbols   bool `json:"require_symbols"`
	MaxAge           int  `json:"max_age"` // days
	HistoryCount     int  `json:"history_count"`
}

// SessionManagement represents session management configuration
type SessionManagement struct {
	SessionTimeout     int  `json:"session_timeout"` // minutes
	IdleTimeout        int  `json:"idle_timeout"`    // minutes
	ConcurrentSessions int  `json:"concurrent_sessions"`
	SecureCookies      bool `json:"secure_cookies"`
}

// TechSecurityAssessmentResult represents the result of security assessment
type TechSecurityAssessmentResult struct {
	ID                      string                        `json:"id"`
	InquiryID               string                        `json:"inquiry_id"`
	OverallScore            float64                       `json:"overall_score"` // 0-100
	RiskLevel               string                        `json:"risk_level"`    // "critical", "high", "medium", "low"
	SecurityVulnerabilities []*TechSecurityVulnerability  `json:"security_vulnerabilities"`
	ComplianceGaps          []*TechComplianceGap          `json:"compliance_gaps"`
	ThreatAnalysis          *TechThreatAnalysis           `json:"threat_analysis"`
	SecurityRecommendations []*TechSecurityRecommendation `json:"security_recommendations"`
	Summary                 string                        `json:"summary"`
	GeneratedAt             time.Time                     `json:"generated_at"`
}

// TechSecurityVulnerability represents a security vulnerability
type TechSecurityVulnerability struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`     // "configuration", "access", "network", "data", "application"
	Severity           string   `json:"severity"` // "critical", "high", "medium", "low"
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Impact             string   `json:"impact"`
	Likelihood         string   `json:"likelihood"` // "high", "medium", "low"
	CVSS               float64  `json:"cvss,omitempty"`
	CWE                string   `json:"cwe,omitempty"`
	References         []string `json:"references"`
	AffectedComponents []string `json:"affected_components"`
}

// TechComplianceGap represents a compliance gap
type TechComplianceGap struct {
	ID            string `json:"id"`
	Framework     string `json:"framework"` // "SOC2", "HIPAA", "PCI-DSS", "GDPR"
	Control       string `json:"control"`
	Requirement   string `json:"requirement"`
	CurrentState  string `json:"current_state"`
	RequiredState string `json:"required_state"`
	Gap           string `json:"gap"`
	Impact        string `json:"impact"`
	Priority      string `json:"priority"`
}

// TechThreatAnalysis represents threat analysis results
type TechThreatAnalysis struct {
	IdentifiedThreats []*TechIdentifiedThreat `json:"identified_threats"`
	AttackVectors     []*TechAttackVector     `json:"attack_vectors"`
	RiskMatrix        *TechRiskMatrix         `json:"risk_matrix"`
}

// TechIdentifiedThreat represents an identified threat
type TechIdentifiedThreat struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Likelihood  string   `json:"likelihood"`
	Impact      string   `json:"impact"`
	RiskScore   float64  `json:"risk_score"`
	Mitigations []string `json:"mitigations"`
}

// TechAttackVector represents an attack vector
type TechAttackVector struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Complexity    string   `json:"complexity"` // "low", "medium", "high"
	Prerequisites []string `json:"prerequisites"`
	Mitigations   []string `json:"mitigations"`
}

// TechRiskMatrix represents a risk assessment matrix
type TechRiskMatrix struct {
	Risks  [][]float64 `json:"risks"` // 2D matrix of risk scores
	Labels struct {
		Impact     []string `json:"impact"`
		Likelihood []string `json:"likelihood"`
	} `json:"labels"`
}

// TechSecurityRecommendation represents a security recommendation
type TechSecurityRecommendation struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`     // "immediate", "short-term", "long-term"
	Priority     string   `json:"priority"` // "critical", "high", "medium", "low"
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Effort       string   `json:"effort"` // "low", "medium", "high"
	Timeline     string   `json:"timeline"`
	Cost         float64  `json:"cost,omitempty"`
	Dependencies []string `json:"dependencies"`
	References   []string `json:"references"`
}

// RemediationRecommendation represents a remediation recommendation
type RemediationRecommendation struct {
	ID              string   `json:"id"`
	VulnerabilityID string   `json:"vulnerability_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Steps           []string `json:"steps"`
	Priority        string   `json:"priority"`
	Effort          string   `json:"effort"`
	Timeline        string   `json:"timeline"`
	Cost            float64  `json:"cost,omitempty"`
	References      []string `json:"references"`
}

// TechPerformanceAnalysisRequest represents a request for performance analysis
type TechPerformanceAnalysisRequest struct {
	InquiryID           string                  `json:"inquiry_id"`
	SystemDescription   string                  `json:"system_description"`
	CloudProvider       string                  `json:"cloud_provider"`
	Services            []string                `json:"services"`
	PerformanceMetrics  *TechPerformanceMetrics `json:"performance_metrics,omitempty"`
	LoadPatterns        *LoadPatterns           `json:"load_patterns,omitempty"`
	ResourceUtilization *ResourceUtilization    `json:"resource_utilization,omitempty"`
	ApplicationProfile  *ApplicationProfile     `json:"application_profile,omitempty"`
	Metadata            map[string]interface{}  `json:"metadata"`
}

// TechPerformanceMetrics represents performance metrics
type TechPerformanceMetrics struct {
	ResponseTime       float64 `json:"response_time"`       // milliseconds
	Throughput         float64 `json:"throughput"`          // requests per second
	ErrorRate          float64 `json:"error_rate"`          // percentage
	Availability       float64 `json:"availability"`        // percentage
	Latency            float64 `json:"latency"`             // milliseconds
	CPUUtilization     float64 `json:"cpu_utilization"`     // percentage
	MemoryUtilization  float64 `json:"memory_utilization"`  // percentage
	DiskUtilization    float64 `json:"disk_utilization"`    // percentage
	NetworkUtilization float64 `json:"network_utilization"` // percentage
}

// LoadPatterns represents load patterns
type LoadPatterns struct {
	PeakLoad         float64            `json:"peak_load"`
	AverageLoad      float64            `json:"average_load"`
	LoadDistribution map[string]float64 `json:"load_distribution"` // time -> load
	ConcurrentUsers  int                `json:"concurrent_users"`
	SessionDuration  float64            `json:"session_duration"` // minutes
}

// ResourceUtilization represents resource utilization
type ResourceUtilization struct {
	CPU     *ResourceMetric `json:"cpu"`
	Memory  *ResourceMetric `json:"memory"`
	Disk    *ResourceMetric `json:"disk"`
	Network *ResourceMetric `json:"network"`
}

// ResourceMetric represents a resource metric
type ResourceMetric struct {
	Current   float64 `json:"current"`   // percentage
	Average   float64 `json:"average"`   // percentage
	Peak      float64 `json:"peak"`      // percentage
	Allocated float64 `json:"allocated"` // units
	Used      float64 `json:"used"`      // units
}

// ApplicationProfile represents application profile
type ApplicationProfile struct {
	Type            string   `json:"type"` // "web", "api", "batch", "streaming"
	Language        string   `json:"language"`
	Framework       string   `json:"framework"`
	DatabaseType    string   `json:"database_type"`
	CachingStrategy string   `json:"caching_strategy"`
	Dependencies    []string `json:"dependencies"`
}

// TechPerformanceAnalysisResult represents the result of performance analysis
type TechPerformanceAnalysisResult struct {
	ID                         string                            `json:"id"`
	InquiryID                  string                            `json:"inquiry_id"`
	OverallScore               float64                           `json:"overall_score"` // 0-100
	PerformanceBottlenecks     []*TechPerformanceBottleneck      `json:"performance_bottlenecks"`
	OptimizationOpportunities  []*TechOptimizationRecommendation `json:"optimization_opportunities"`
	BenchmarkComparison        *TechBenchmarkComparison          `json:"benchmark_comparison"`
	PerformanceRecommendations []*TechPerformanceRecommendation  `json:"performance_recommendations"`
	Summary                    string                            `json:"summary"`
	GeneratedAt                time.Time                         `json:"generated_at"`
}

// TechPerformanceBottleneck represents a performance bottleneck
type TechPerformanceBottleneck struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"`     // "cpu", "memory", "disk", "network", "database", "application"
	Severity    string             `json:"severity"` // "critical", "high", "medium", "low"
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Impact      string             `json:"impact"`
	Location    string             `json:"location"`
	Metrics     map[string]float64 `json:"metrics"`
}

// TechOptimizationRecommendation represents an optimization recommendation
type TechOptimizationRecommendation struct {
	ID                  string   `json:"id"`
	Type                string   `json:"type"`     // "scaling", "caching", "database", "code", "infrastructure"
	Priority            string   `json:"priority"` // "critical", "high", "medium", "low"
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	ExpectedImprovement string   `json:"expected_improvement"`
	Effort              string   `json:"effort"` // "low", "medium", "high"
	Timeline            string   `json:"timeline"`
	Cost                float64  `json:"cost,omitempty"`
	Dependencies        []string `json:"dependencies"`
	References          []string `json:"references"`
}

// TechBenchmarkComparison represents benchmark comparison results
type TechBenchmarkComparison struct {
	Industry   string             `json:"industry"`
	Benchmarks map[string]float64 `json:"benchmarks"` // metric -> benchmark value
	Current    map[string]float64 `json:"current"`    // metric -> current value
	Comparison map[string]string  `json:"comparison"` // metric -> "above", "below", "at"
	Percentile map[string]float64 `json:"percentile"` // metric -> percentile
}

// TechPerformanceRecommendation represents a performance recommendation
type TechPerformanceRecommendation struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`     // "immediate", "short-term", "long-term"
	Priority     string   `json:"priority"` // "critical", "high", "medium", "low"
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Effort       string   `json:"effort"` // "low", "medium", "high"
	Timeline     string   `json:"timeline"`
	Cost         float64  `json:"cost,omitempty"`
	Dependencies []string `json:"dependencies"`
	References   []string `json:"references"`
}

// TechComplianceAnalysisRequest represents a request for compliance analysis
type TechComplianceAnalysisRequest struct {
	InquiryID            string                 `json:"inquiry_id"`
	SystemDescription    string                 `json:"system_description"`
	CloudProvider        string                 `json:"cloud_provider"`
	Services             []string               `json:"services"`
	ComplianceFrameworks []string               `json:"compliance_frameworks"` // "SOC2", "HIPAA", "PCI-DSS", "GDPR"
	DataClassification   []string               `json:"data_classification"`
	GeographicScope      []string               `json:"geographic_scope"`
	BusinessContext      *BusinessContext       `json:"business_context,omitempty"`
	CurrentControls      []*TechSecurityControl `json:"current_controls,omitempty"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// BusinessContext represents business context for compliance
type BusinessContext struct {
	Industry        string   `json:"industry"`
	CompanySize     string   `json:"company_size"` // "startup", "small", "medium", "large", "enterprise"
	Revenue         float64  `json:"revenue,omitempty"`
	CustomerBase    string   `json:"customer_base"`
	DataTypes       []string `json:"data_types"`
	RegulatoryScope []string `json:"regulatory_scope"`
}

// TechComplianceAnalysisResult represents the result of compliance analysis
type TechComplianceAnalysisResult struct {
	ID                        string                          `json:"id"`
	InquiryID                 string                          `json:"inquiry_id"`
	OverallScore              float64                         `json:"overall_score"`     // 0-100
	ComplianceStatus          map[string]string               `json:"compliance_status"` // framework -> status
	ComplianceGaps            []*TechComplianceGap            `json:"compliance_gaps"`
	ControlAssessment         []*TechControlAssessment        `json:"control_assessment"`
	ComplianceRecommendations []*TechComplianceRecommendation `json:"compliance_recommendations"`
	RoadmapToCompliance       *TechComplianceRoadmap          `json:"roadmap_to_compliance"`
	Summary                   string                          `json:"summary"`
	GeneratedAt               time.Time                       `json:"generated_at"`
}

// TechControlAssessment represents assessment of a specific control
type TechControlAssessment struct {
	ControlID       string  `json:"control_id"`
	Framework       string  `json:"framework"`
	ControlName     string  `json:"control_name"`
	RequiredState   string  `json:"required_state"`
	CurrentState    string  `json:"current_state"`
	ComplianceLevel float64 `json:"compliance_level"` // 0-100
	Gap             string  `json:"gap"`
	Priority        string  `json:"priority"`
	Effort          string  `json:"effort"`
}

// TechComplianceRecommendation represents a compliance recommendation
type TechComplianceRecommendation struct {
	ID           string   `json:"id"`
	Framework    string   `json:"framework"`
	Type         string   `json:"type"`     // "immediate", "short-term", "long-term"
	Priority     string   `json:"priority"` // "critical", "high", "medium", "low"
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Effort       string   `json:"effort"` // "low", "medium", "high"
	Timeline     string   `json:"timeline"`
	Cost         float64  `json:"cost,omitempty"`
	Dependencies []string `json:"dependencies"`
	References   []string `json:"references"`
}

// TechComplianceRoadmap represents a roadmap to achieve compliance
type TechComplianceRoadmap struct {
	Framework  string                     `json:"framework"`
	Phases     []*TechCompliancePhase     `json:"phases"`
	Timeline   string                     `json:"timeline"`
	TotalCost  float64                    `json:"total_cost"`
	Milestones []*TechComplianceMilestone `json:"milestones"`
}

// TechCompliancePhase represents a phase in compliance roadmap
type TechCompliancePhase struct {
	PhaseNumber  int      `json:"phase_number"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Duration     string   `json:"duration"`
	Controls     []string `json:"controls"`
	Deliverables []string `json:"deliverables"`
	Cost         float64  `json:"cost"`
}

// TechComplianceMilestone represents a milestone in compliance roadmap
type TechComplianceMilestone struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TargetDate  time.Time `json:"target_date"`
	Criteria    []string  `json:"criteria"`
}

// TechComplianceRemediation represents compliance remediation steps
type TechComplianceRemediation struct {
	ID          string   `json:"id"`
	GapID       string   `json:"gap_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Priority    string   `json:"priority"`
	Effort      string   `json:"effort"`
	Timeline    string   `json:"timeline"`
	Cost        float64  `json:"cost,omitempty"`
	References  []string `json:"references"`
}

// TechCostOptimization represents a cost optimization opportunity
type TechCostOptimization struct {
	ID                   string  `json:"id"`
	Type                 string  `json:"type"` // "rightsizing", "reserved-instances", "spot-instances", "storage-optimization"
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	CurrentCost          float64 `json:"current_cost"`
	OptimizedCost        float64 `json:"optimized_cost"`
	Savings              float64 `json:"savings"`
	SavingsPercent       float64 `json:"savings_percent"`
	ImplementationEffort string  `json:"implementation_effort"`
	RiskLevel            string  `json:"risk_level"`
	Timeline             string  `json:"timeline"`
}

// ComprehensiveAnalysisRequest represents a request for comprehensive analysis
type ComprehensiveAnalysisRequest struct {
	InquiryID           string                             `json:"inquiry_id"`
	CodeAnalysisRequest *CodeAnalysisRequest               `json:"code_analysis_request,omitempty"`
	ArchitectureRequest *TechArchitectureAssessmentRequest `json:"architecture_request,omitempty"`
	SecurityRequest     *TechSecurityAssessmentRequest     `json:"security_request,omitempty"`
	PerformanceRequest  *TechPerformanceAnalysisRequest    `json:"performance_request,omitempty"`
	ComplianceRequest   *TechComplianceAnalysisRequest     `json:"compliance_request,omitempty"`
	AnalysisScope       []string                           `json:"analysis_scope"` // "code", "architecture", "security", "performance", "compliance"
	PriorityAreas       []string                           `json:"priority_areas"`
	Metadata            map[string]interface{}             `json:"metadata"`
}

// ComprehensiveAnalysisResult represents the result of comprehensive analysis
type ComprehensiveAnalysisResult struct {
	ID                         string                            `json:"id"`
	InquiryID                  string                            `json:"inquiry_id"`
	OverallScore               float64                           `json:"overall_score"` // 0-100
	CodeAnalysisResult         *CodeAnalysisResult               `json:"code_analysis_result,omitempty"`
	ArchitectureResult         *TechArchitectureAssessmentResult `json:"architecture_result,omitempty"`
	SecurityResult             *TechSecurityAssessmentResult     `json:"security_result,omitempty"`
	PerformanceResult          *TechPerformanceAnalysisResult    `json:"performance_result,omitempty"`
	ComplianceResult           *TechComplianceAnalysisResult     `json:"compliance_result,omitempty"`
	CrossCuttingFindings       []*CrossCuttingFinding            `json:"cross_cutting_findings"`
	PrioritizedRecommendations []*PrioritizedRecommendation      `json:"prioritized_recommendations"`
	ExecutiveSummary           string                            `json:"executive_summary"`
	TechnicalSummary           string                            `json:"technical_summary"`
	ActionPlan                 *TechnicalActionPlan              `json:"action_plan"`
	GeneratedAt                time.Time                         `json:"generated_at"`
}

// CrossCuttingFinding represents findings that span multiple analysis areas
type CrossCuttingFinding struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Areas           []string `json:"areas"` // analysis areas this finding affects
	Impact          string   `json:"impact"`
	Priority        string   `json:"priority"`
	Recommendations []string `json:"recommendations"`
}

// PrioritizedRecommendation represents a prioritized recommendation
type PrioritizedRecommendation struct {
	ID           string   `json:"id"`
	Source       string   `json:"source"` // "code", "architecture", "security", "performance", "compliance"
	Type         string   `json:"type"`
	Priority     string   `json:"priority"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Benefits     []string `json:"benefits"`
	Effort       string   `json:"effort"`
	Timeline     string   `json:"timeline"`
	Cost         float64  `json:"cost,omitempty"`
	Dependencies []string `json:"dependencies"`
	References   []string `json:"references"`
	Score        float64  `json:"score"` // prioritization score
}

// TechnicalActionPlan represents a technical action plan
type TechnicalActionPlan struct {
	Phases      []*ActionPhase     `json:"phases"`
	Timeline    string             `json:"timeline"`
	TotalCost   float64            `json:"total_cost"`
	Resources   *RequiredResources `json:"resources"`
	Milestones  []*ActionMilestone `json:"milestones"`
	RiskFactors []string           `json:"risk_factors"`
}

// ActionPhase represents a phase in the action plan
type ActionPhase struct {
	PhaseNumber   int      `json:"phase_number"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Duration      string   `json:"duration"`
	Actions       []string `json:"actions"`
	Deliverables  []string `json:"deliverables"`
	Cost          float64  `json:"cost"`
	Prerequisites []string `json:"prerequisites"`
}

// RequiredResources represents required resources for implementation
type RequiredResources struct {
	TechnicalRoles   []string `json:"technical_roles"`
	SkillsRequired   []string `json:"skills_required"`
	Tools            []string `json:"tools"`
	ExternalServices []string `json:"external_services"`
	Budget           float64  `json:"budget"`
}

// ActionMilestone represents a milestone in the action plan
type ActionMilestone struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TargetDate   time.Time `json:"target_date"`
	Criteria     []string  `json:"criteria"`
	Deliverables []string  `json:"deliverables"`
}
