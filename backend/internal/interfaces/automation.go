package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// AutomationService defines the interface for advanced automation and integration
type AutomationService interface {
	// Client Environment Discovery
	DiscoverClientEnvironment(ctx context.Context, clientID string, credentials *ClientCredentials) (*EnvironmentDiscovery, error)
	AnalyzeEnvironmentChanges(ctx context.Context, clientID string, previousSnapshot, currentSnapshot *EnvironmentSnapshot) (*ChangeAnalysis, error)

	// Integration Management
	RegisterIntegration(ctx context.Context, integration *Integration) error
	GetIntegrations(ctx context.Context, clientID string) ([]*Integration, error)
	TestIntegration(ctx context.Context, integrationID string) (*IntegrationTestResult, error)

	// Automated Report Generation
	GenerateAutomatedReport(ctx context.Context, trigger *ReportTrigger) (*domain.Report, error)
	ScheduleReportGeneration(ctx context.Context, schedule *ReportSchedule) error

	// Proactive Recommendations
	GenerateProactiveRecommendations(ctx context.Context, clientID string, usagePatterns *UsagePatterns) ([]*ProactiveRecommendation, error)
	AnalyzeUsagePatterns(ctx context.Context, clientID string, timeRange AutomationTimeRange) (*UsagePatterns, error)
}

// EnvironmentDiscoveryService handles automated client environment discovery
type EnvironmentDiscoveryService interface {
	ScanAWSEnvironment(ctx context.Context, credentials *AWSCredentials) (*AWSEnvironmentSnapshot, error)
	ScanAzureEnvironment(ctx context.Context, credentials *AzureCredentials) (*AzureEnvironmentSnapshot, error)
	ScanGCPEnvironment(ctx context.Context, credentials *GCPCredentials) (*GCPEnvironmentSnapshot, error)

	CompareSnapshots(ctx context.Context, previous, current *EnvironmentSnapshot) (*ChangeAnalysis, error)
	GenerateDiscoveryReport(ctx context.Context, discovery *EnvironmentDiscovery) (*DiscoveryReport, error)
}

// IntegrationService handles third-party tool integrations
type IntegrationService interface {
	// Monitoring Tools
	IntegrateCloudWatch(ctx context.Context, config *CloudWatchConfig) (*Integration, error)
	IntegrateDatadog(ctx context.Context, config *DatadogConfig) (*Integration, error)
	IntegrateNewRelic(ctx context.Context, config *NewRelicConfig) (*Integration, error)

	// Ticketing Systems
	IntegrateJira(ctx context.Context, config *JiraConfig) (*Integration, error)
	IntegrateServiceNow(ctx context.Context, config *ServiceNowConfig) (*Integration, error)

	// Documentation Systems
	IntegrateConfluence(ctx context.Context, config *ConfluenceConfig) (*Integration, error)
	IntegrateNotion(ctx context.Context, config *NotionConfig) (*Integration, error)

	// Communication Tools
	IntegrateSlack(ctx context.Context, config *SlackConfig) (*Integration, error)
	IntegrateTeams(ctx context.Context, config *TeamsConfig) (*Integration, error)

	SyncData(ctx context.Context, integrationID string) error
	GetIntegrationData(ctx context.Context, integrationID string, dataType string) (interface{}, error)
}

// ProactiveRecommendationEngine generates recommendations based on usage patterns
type ProactiveRecommendationEngine interface {
	AnalyzeCostTrends(ctx context.Context, clientID string, timeRange AutomationTimeRange) (*AutomationCostTrendAnalysis, error)
	AnalyzePerformancePatterns(ctx context.Context, clientID string, timeRange AutomationTimeRange) (*AutomationPerformancePatternAnalysisUpdated, error)
	AnalyzeSecurityPosture(ctx context.Context, clientID string) (*AutomationSecurityPostureAnalysisUpdated, error)

	GenerateCostOptimizationRecommendations(ctx context.Context, analysis *AutomationCostTrendAnalysis) ([]*AutomationCostOptimizationRecommendationUpdated, error)
	GeneratePerformanceRecommendations(ctx context.Context, analysis *AutomationPerformancePatternAnalysisUpdated) ([]*AutomationPerformanceRecommendationUpdated, error)
	GenerateSecurityRecommendations(ctx context.Context, analysis *AutomationSecurityPostureAnalysisUpdated) ([]*AutomationSecurityRecommendationUpdated, error)
}

// Data Models

// ClientCredentials represents credentials for accessing client environments
type ClientCredentials struct {
	ClientID    string                 `json:"client_id"`
	Provider    string                 `json:"provider"` // "aws", "azure", "gcp"
	Credentials map[string]interface{} `json:"credentials"`
	Region      string                 `json:"region,omitempty"`
	AccountID   string                 `json:"account_id,omitempty"`
}

// AWSCredentials represents AWS-specific credentials
type AWSCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	SessionToken    string `json:"session_token,omitempty"`
	Region          string `json:"region"`
	RoleArn         string `json:"role_arn,omitempty"`
}

// AzureCredentials represents Azure-specific credentials
type AzureCredentials struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TenantID       string `json:"tenant_id"`
	SubscriptionID string `json:"subscription_id"`
}

// GCPCredentials represents GCP-specific credentials
type GCPCredentials struct {
	ProjectID           string `json:"project_id"`
	ServiceAccountKey   string `json:"service_account_key"`
	ServiceAccountEmail string `json:"service_account_email"`
}

// EnvironmentDiscovery represents the result of environment discovery
type EnvironmentDiscovery struct {
	ClientID         string                     `json:"client_id"`
	Provider         string                     `json:"provider"`
	DiscoveredAt     time.Time                  `json:"discovered_at"`
	Resources        []*DiscoveredResource      `json:"resources"`
	Services         []*DiscoveredService       `json:"services"`
	Configurations   map[string]interface{}     `json:"configurations"`
	CostEstimate     *CostEstimate              `json:"cost_estimate"`
	SecurityFindings []*SecurityFinding         `json:"security_findings"`
	Recommendations  []*AutomatedRecommendation `json:"recommendations"`
}

// EnvironmentSnapshot represents a snapshot of client environment at a point in time
type EnvironmentSnapshot struct {
	ID           string                 `json:"id"`
	ClientID     string                 `json:"client_id"`
	Provider     string                 `json:"provider"`
	SnapshotTime time.Time              `json:"snapshot_time"`
	Resources    []*ResourceSnapshot    `json:"resources"`
	Metrics      map[string]interface{} `json:"metrics"`
	Costs        *CostSnapshot          `json:"costs"`
	Checksum     string                 `json:"checksum"`
}

// AWSEnvironmentSnapshot represents AWS-specific environment snapshot
type AWSEnvironmentSnapshot struct {
	*EnvironmentSnapshot
	EC2Instances     []*EC2Instance     `json:"ec2_instances"`
	RDSInstances     []*RDSInstance     `json:"rds_instances"`
	S3Buckets        []*S3Bucket        `json:"s3_buckets"`
	LambdaFunctions  []*LambdaFunction  `json:"lambda_functions"`
	VPCs             []*VPC             `json:"vpcs"`
	AWSLoadBalancers []*AWSLoadBalancer `json:"load_balancers"`
}

// AzureEnvironmentSnapshot represents Azure-specific environment snapshot
type AzureEnvironmentSnapshot struct {
	*EnvironmentSnapshot
	VirtualMachines []*AzureVM        `json:"virtual_machines"`
	StorageAccounts []*StorageAccount `json:"storage_accounts"`
	SQLDatabases    []*SQLDatabase    `json:"sql_databases"`
	AppServices     []*AppService     `json:"app_services"`
	VirtualNetworks []*VirtualNetwork `json:"virtual_networks"`
}

// GCPEnvironmentSnapshot represents GCP-specific environment snapshot
type GCPEnvironmentSnapshot struct {
	*EnvironmentSnapshot
	ComputeInstances []*ComputeInstance `json:"compute_instances"`
	CloudStorage     []*CloudStorage    `json:"cloud_storage"`
	CloudSQL         []*CloudSQL        `json:"cloud_sql"`
	CloudFunctions   []*CloudFunction   `json:"cloud_functions"`
	VPCNetworks      []*VPCNetwork      `json:"vpc_networks"`
}

// ChangeAnalysis represents analysis of environment changes
type ChangeAnalysis struct {
	ClientID          string                        `json:"client_id"`
	AnalysisTime      time.Time                     `json:"analysis_time"`
	TimeRange         TimeRange                     `json:"time_range"`
	AddedResources    []*ResourceChange             `json:"added_resources"`
	ModifiedResources []*ResourceChange             `json:"modified_resources"`
	DeletedResources  []*ResourceChange             `json:"deleted_resources"`
	CostImpact        *AutomationCostImpactAnalysis `json:"cost_impact"`
	SecurityImpact    *SecurityImpactAnalysis       `json:"security_impact"`
	Recommendations   []*ChangeRecommendation       `json:"recommendations"`
}

// Integration represents a third-party tool integration
type Integration struct {
	ID            string                 `json:"id"`
	ClientID      string                 `json:"client_id"`
	Type          IntegrationType        `json:"type"`
	Name          string                 `json:"name"`
	Configuration map[string]interface{} `json:"configuration"`
	Status        IntegrationStatus      `json:"status"`
	LastSync      *time.Time             `json:"last_sync,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// IntegrationType represents the type of integration
type IntegrationType string

const (
	IntegrationTypeMonitoring       IntegrationType = "monitoring"
	IntegrationTypeTicketing        IntegrationType = "ticketing"
	IntegrationTypeDocumentation    IntegrationType = "documentation"
	IntegrationTypeCommunication    IntegrationType = "communication"
	IntegrationTypeCI_CD            IntegrationType = "ci_cd"
	IntegrationTypeSecurityScanning IntegrationType = "security_scanning"
)

// IntegrationStatus represents the status of an integration
type IntegrationStatus string

const (
	IntegrationStatusActive   IntegrationStatus = "active"
	IntegrationStatusInactive IntegrationStatus = "inactive"
	IntegrationStatusError    IntegrationStatus = "error"
	IntegrationStatusPending  IntegrationStatus = "pending"
)

// ReportTrigger represents a trigger for automated report generation
type ReportTrigger struct {
	ID          string                 `json:"id"`
	ClientID    string                 `json:"client_id"`
	TriggerType TriggerType            `json:"trigger_type"`
	Conditions  map[string]interface{} `json:"conditions"`
	ReportType  domain.ReportType      `json:"report_type"`
	Recipients  []string               `json:"recipients"`
	CreatedAt   time.Time              `json:"created_at"`
}

// TriggerType represents the type of report trigger
type TriggerType string

const (
	TriggerTypeScheduled         TriggerType = "scheduled"
	TriggerTypeThreshold         TriggerType = "threshold"
	TriggerTypeEnvironmentChange TriggerType = "environment_change"
	TriggerTypeCostAnomaly       TriggerType = "cost_anomaly"
	TriggerTypeSecurityAlert     TriggerType = "security_alert"
)

// ReportSchedule represents a schedule for automated report generation
type ReportSchedule struct {
	ID             string            `json:"id"`
	ClientID       string            `json:"client_id"`
	ReportType     domain.ReportType `json:"report_type"`
	CronExpression string            `json:"cron_expression"`
	Recipients     []string          `json:"recipients"`
	Enabled        bool              `json:"enabled"`
	NextRun        time.Time         `json:"next_run"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// ProactiveRecommendation represents a proactive recommendation
type ProactiveRecommendation struct {
	ID               string                 `json:"id"`
	ClientID         string                 `json:"client_id"`
	Type             RecommendationType     `json:"type"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Priority         RecommendationPriority `json:"priority"`
	Impact           string                 `json:"impact"`
	Effort           string                 `json:"effort"`
	PotentialSavings float64                `json:"potential_savings,omitempty"`
	ActionItems      []string               `json:"action_items"`
	Resources        []string               `json:"resources"`
	CreatedAt        time.Time              `json:"created_at"`
	ExpiresAt        *time.Time             `json:"expires_at,omitempty"`
}

// RecommendationType represents the type of recommendation
type RecommendationType string

const (
	RecommendationTypeCostOptimization RecommendationType = "cost_optimization"
	RecommendationTypePerformance      RecommendationType = "performance"
	RecommendationTypeSecurity         RecommendationType = "security"
	RecommendationTypeCompliance       RecommendationType = "compliance"
	RecommendationTypeArchitecture     RecommendationType = "architecture"
	RecommendationTypeOperational      RecommendationType = "operational"
)

// RecommendationPriority represents the priority of a recommendation
type RecommendationPriority string

const (
	RecommendationPriorityLow      RecommendationPriority = "low"
	RecommendationPriorityMedium   RecommendationPriority = "medium"
	RecommendationPriorityHigh     RecommendationPriority = "high"
	RecommendationPriorityCritical RecommendationPriority = "critical"
)

// UsagePatterns represents usage patterns for a client
type UsagePatterns struct {
	ClientID                 string                         `json:"client_id"`
	TimeRange                AutomationTimeRange            `json:"time_range"`
	CostTrends               *AutomationCostTrends          `json:"cost_trends"`
	ResourceUtilization      *AutomationResourceUtilization `json:"resource_utilization"`
	PerformanceMetrics       *AutomationPerformanceMetrics  `json:"performance_metrics"`
	AutomationSecurityEvents []*AutomationSecurityEvent     `json:"security_events"`
	AnomaliesDetected        []*Anomaly                     `json:"anomalies_detected"`
}

// AutomationTimeRange represents a time range for automation (different from existing TimeRange)
type AutomationTimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// Supporting data structures for various components

// DiscoveredResource represents a discovered cloud resource
type DiscoveredResource struct {
	ID            string                         `json:"id"`
	Type          string                         `json:"type"`
	Name          string                         `json:"name"`
	Region        string                         `json:"region"`
	Tags          map[string]string              `json:"tags"`
	Configuration map[string]interface{}         `json:"configuration"`
	Cost          *ResourceCost                  `json:"cost,omitempty"`
	Utilization   *AutomationResourceUtilization `json:"utilization,omitempty"`
}

// DiscoveredService represents a discovered cloud service
type DiscoveredService struct {
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	Version       string                 `json:"version,omitempty"`
	Configuration map[string]interface{} `json:"configuration"`
	Dependencies  []string               `json:"dependencies"`
	Endpoints     []string               `json:"endpoints"`
}

// CostEstimate represents cost estimation for discovered resources
type CostEstimate struct {
	MonthlyCost float64            `json:"monthly_cost"`
	AnnualCost  float64            `json:"annual_cost"`
	Currency    string             `json:"currency"`
	Breakdown   map[string]float64 `json:"breakdown"`
	LastUpdated time.Time          `json:"last_updated"`
}

// SecurityFinding represents a security finding from environment discovery
type SecurityFinding struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Remediation string                 `json:"remediation"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AutomatedRecommendation represents an automated recommendation
type AutomatedRecommendation struct {
	Type             string   `json:"type"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Priority         string   `json:"priority"`
	ActionItems      []string `json:"action_items"`
	EstimatedSavings float64  `json:"estimated_savings,omitempty"`
}

// Integration configuration types
type CloudWatchConfig struct {
	Region          string   `json:"region"`
	AccessKeyID     string   `json:"access_key_id"`
	SecretAccessKey string   `json:"secret_access_key"`
	LogGroups       []string `json:"log_groups"`
	Metrics         []string `json:"metrics"`
}

type DatadogConfig struct {
	APIKey string   `json:"api_key"`
	AppKey string   `json:"app_key"`
	Site   string   `json:"site"`
	Tags   []string `json:"tags"`
}

type NewRelicConfig struct {
	APIKey    string `json:"api_key"`
	AccountID string `json:"account_id"`
	Region    string `json:"region"`
}

type JiraConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	APIToken string `json:"api_token"`
	Project  string `json:"project"`
}

type ServiceNowConfig struct {
	Instance string `json:"instance"`
	Username string `json:"username"`
	Password string `json:"password"`
	Table    string `json:"table"`
}

type ConfluenceConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	APIToken string `json:"api_token"`
	Space    string `json:"space"`
}

type NotionConfig struct {
	APIToken   string `json:"api_token"`
	DatabaseID string `json:"database_id"`
}

type SlackConfig struct {
	BotToken  string `json:"bot_token"`
	Channel   string `json:"channel"`
	Workspace string `json:"workspace"`
}

type TeamsConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel"`
}

// IntegrationTestResult represents the result of testing an integration
type IntegrationTestResult struct {
	IntegrationID string                 `json:"integration_id"`
	Success       bool                   `json:"success"`
	Message       string                 `json:"message"`
	TestedAt      time.Time              `json:"tested_at"`
	ResponseTime  int64                  `json:"response_time_ms"`
	Details       map[string]interface{} `json:"details"`
}

// Additional supporting types for comprehensive automation

// ResourceSnapshot represents a snapshot of a resource
type ResourceSnapshot struct {
	ResourceID    string                 `json:"resource_id"`
	Type          string                 `json:"type"`
	State         string                 `json:"state"`
	Configuration map[string]interface{} `json:"configuration"`
	Metrics       map[string]float64     `json:"metrics"`
	Tags          map[string]string      `json:"tags"`
}

// CostSnapshot represents cost information at a point in time
type CostSnapshot struct {
	TotalCost    float64            `json:"total_cost"`
	Currency     string             `json:"currency"`
	Breakdown    map[string]float64 `json:"breakdown"`
	Period       string             `json:"period"`
	SnapshotTime time.Time          `json:"snapshot_time"`
}

// ResourceChange represents a change in a resource
type ResourceChange struct {
	ResourceID string                 `json:"resource_id"`
	ChangeType string                 `json:"change_type"` // "added", "modified", "deleted"
	OldState   map[string]interface{} `json:"old_state,omitempty"`
	NewState   map[string]interface{} `json:"new_state,omitempty"`
	Impact     string                 `json:"impact"`
	Timestamp  time.Time              `json:"timestamp"`
}

// AutomationCostImpactAnalysis represents the cost impact of changes (prefixed to avoid conflicts)
type AutomationCostImpactAnalysis struct {
	TotalImpact     float64            `json:"total_impact"`
	Currency        string             `json:"currency"`
	ImpactBreakdown map[string]float64 `json:"impact_breakdown"`
	Trend           string             `json:"trend"` // "increasing", "decreasing", "stable"
	Recommendations []string           `json:"recommendations"`
}

// SecurityImpactAnalysis represents the security impact of changes
type SecurityImpactAnalysis struct {
	RiskLevel          string   `json:"risk_level"`
	NewVulnerabilities []string `json:"new_vulnerabilities"`
	ResolvedIssues     []string `json:"resolved_issues"`
	Recommendations    []string `json:"recommendations"`
}

// ChangeRecommendation represents a recommendation based on environment changes
type ChangeRecommendation struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Actions     []string `json:"actions"`
	Impact      string   `json:"impact"`
}

// AWS-specific resource types
type EC2Instance struct {
	InstanceID   string            `json:"instance_id"`
	InstanceType string            `json:"instance_type"`
	State        string            `json:"state"`
	Region       string            `json:"region"`
	Tags         map[string]string `json:"tags"`
	LaunchTime   time.Time         `json:"launch_time"`
	Utilization  *CPUUtilization   `json:"utilization,omitempty"`
}

type RDSInstance struct {
	DBInstanceIdentifier string            `json:"db_instance_identifier"`
	DBInstanceClass      string            `json:"db_instance_class"`
	Engine               string            `json:"engine"`
	EngineVersion        string            `json:"engine_version"`
	Status               string            `json:"status"`
	Tags                 map[string]string `json:"tags"`
}

type S3Bucket struct {
	Name         string            `json:"name"`
	Region       string            `json:"region"`
	CreationDate time.Time         `json:"creation_date"`
	Size         int64             `json:"size"`
	ObjectCount  int64             `json:"object_count"`
	Tags         map[string]string `json:"tags"`
}

type LambdaFunction struct {
	FunctionName string            `json:"function_name"`
	Runtime      string            `json:"runtime"`
	Handler      string            `json:"handler"`
	CodeSize     int64             `json:"code_size"`
	Timeout      int               `json:"timeout"`
	MemorySize   int               `json:"memory_size"`
	Tags         map[string]string `json:"tags"`
}

type VPC struct {
	VpcID     string            `json:"vpc_id"`
	CidrBlock string            `json:"cidr_block"`
	State     string            `json:"state"`
	Region    string            `json:"region"`
	Tags      map[string]string `json:"tags"`
}

// AWSLoadBalancer represents AWS load balancer (prefixed to avoid conflicts)
type AWSLoadBalancer struct {
	LoadBalancerArn  string            `json:"load_balancer_arn"`
	LoadBalancerName string            `json:"load_balancer_name"`
	Type             string            `json:"type"`
	Scheme           string            `json:"scheme"`
	State            string            `json:"state"`
	Tags             map[string]string `json:"tags"`
}

// Azure-specific resource types
type AzureVM struct {
	Name          string            `json:"name"`
	Size          string            `json:"size"`
	Location      string            `json:"location"`
	Status        string            `json:"status"`
	ResourceGroup string            `json:"resource_group"`
	Tags          map[string]string `json:"tags"`
}

type StorageAccount struct {
	Name          string            `json:"name"`
	Kind          string            `json:"kind"`
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resource_group"`
	Tags          map[string]string `json:"tags"`
}

type SQLDatabase struct {
	Name          string            `json:"name"`
	ServerName    string            `json:"server_name"`
	Edition       string            `json:"edition"`
	ServiceTier   string            `json:"service_tier"`
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resource_group"`
	Tags          map[string]string `json:"tags"`
}

type AppService struct {
	Name          string            `json:"name"`
	Kind          string            `json:"kind"`
	Location      string            `json:"location"`
	State         string            `json:"state"`
	ResourceGroup string            `json:"resource_group"`
	Tags          map[string]string `json:"tags"`
}

type VirtualNetwork struct {
	Name          string            `json:"name"`
	AddressSpace  []string          `json:"address_space"`
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resource_group"`
	Tags          map[string]string `json:"tags"`
}

// GCP-specific resource types
type ComputeInstance struct {
	Name        string            `json:"name"`
	MachineType string            `json:"machine_type"`
	Zone        string            `json:"zone"`
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
}

type CloudStorage struct {
	Name         string            `json:"name"`
	Location     string            `json:"location"`
	StorageClass string            `json:"storage_class"`
	Labels       map[string]string `json:"labels"`
}

type CloudSQL struct {
	Name            string            `json:"name"`
	DatabaseVersion string            `json:"database_version"`
	Tier            string            `json:"tier"`
	Region          string            `json:"region"`
	State           string            `json:"state"`
	Labels          map[string]string `json:"labels"`
}

type CloudFunction struct {
	Name       string            `json:"name"`
	Runtime    string            `json:"runtime"`
	EntryPoint string            `json:"entry_point"`
	Region     string            `json:"region"`
	Status     string            `json:"status"`
	Labels     map[string]string `json:"labels"`
}

type VPCNetwork struct {
	Name                  string   `json:"name"`
	AutoCreateSubnetworks bool     `json:"auto_create_subnetworks"`
	Subnetworks           []string `json:"subnetworks"`
	RoutingMode           string   `json:"routing_mode"`
}

// Analysis and metrics types (prefixed to avoid conflicts)
type AutomationCostTrends struct {
	TotalCost        float64            `json:"total_cost"`
	PreviousCost     float64            `json:"previous_cost"`
	PercentChange    float64            `json:"percent_change"`
	Trend            string             `json:"trend"`
	ServiceBreakdown map[string]float64 `json:"service_breakdown"`
	DailyTrends      []DailyCost        `json:"daily_trends"`
}

type DailyCost struct {
	Date time.Time `json:"date"`
	Cost float64   `json:"cost"`
}

type AutomationResourceUtilization struct {
	CPU     *CPUUtilization     `json:"cpu,omitempty"`
	Memory  *MemoryUtilization  `json:"memory,omitempty"`
	Storage *StorageUtilization `json:"storage,omitempty"`
	Network *NetworkUtilization `json:"network,omitempty"`
}

type CPUUtilization struct {
	Average float64 `json:"average"`
	Maximum float64 `json:"maximum"`
	Minimum float64 `json:"minimum"`
}

type MemoryUtilization struct {
	Average float64 `json:"average"`
	Maximum float64 `json:"maximum"`
	Minimum float64 `json:"minimum"`
}

type StorageUtilization struct {
	Used    int64   `json:"used"`
	Total   int64   `json:"total"`
	Percent float64 `json:"percent"`
}

type NetworkUtilization struct {
	InboundMbps  float64 `json:"inbound_mbps"`
	OutboundMbps float64 `json:"outbound_mbps"`
}

type AutomationPerformanceMetrics struct {
	ResponseTime *ResponseTimeMetrics `json:"response_time,omitempty"`
	Throughput   *ThroughputMetrics   `json:"throughput,omitempty"`
	ErrorRate    float64              `json:"error_rate"`
	Availability float64              `json:"availability"`
}

type ResponseTimeMetrics struct {
	Average float64 `json:"average"`
	P95     float64 `json:"p95"`
	P99     float64 `json:"p99"`
}

type ThroughputMetrics struct {
	RequestsPerSecond     float64 `json:"requests_per_second"`
	TransactionsPerSecond float64 `json:"transactions_per_second"`
}

type AutomationSecurityEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
	Source      string                 `json:"source"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type Anomaly struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
	Description string                 `json:"description"`
	Value       float64                `json:"value"`
	Expected    float64                `json:"expected"`
	Deviation   float64                `json:"deviation"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Analysis result types (prefixed to avoid conflicts)
type AutomationCostTrendAnalysis struct {
	ClientID     string                                   `json:"client_id"`
	TimeRange    AutomationTimeRange                      `json:"time_range"`
	TotalCost    float64                                  `json:"total_cost"`
	CostTrend    string                                   `json:"cost_trend"`
	Anomalies    []*CostAnomaly                           `json:"anomalies"`
	Forecasts    []*AutomationCostForecast                `json:"forecasts"`
	Optimization []*AutomationCostOptimizationOpportunity `json:"optimization"`
}

type CostAnomaly struct {
	Date         time.Time `json:"date"`
	ExpectedCost float64   `json:"expected_cost"`
	ActualCost   float64   `json:"actual_cost"`
	Deviation    float64   `json:"deviation"`
	Service      string    `json:"service"`
	Reason       string    `json:"reason"`
}

type AutomationCostForecast struct {
	Date         time.Time `json:"date"`
	ForecastCost float64   `json:"forecast_cost"`
	Confidence   float64   `json:"confidence"`
}

type AutomationCostOptimizationOpportunity struct {
	Type             string   `json:"type"`
	Description      string   `json:"description"`
	PotentialSavings float64  `json:"potential_savings"`
	Effort           string   `json:"effort"`
	Resources        []string `json:"resources"`
}

type AutomationPerformancePatternAnalysis struct {
	ClientID        string                                 `json:"client_id"`
	TimeRange       AutomationTimeRange                    `json:"time_range"`
	OverallHealth   string                                 `json:"overall_health"`
	Bottlenecks     []*AutomationPerformanceBottleneck     `json:"bottlenecks"`
	Trends          []*AutomationPerformanceTrend          `json:"trends"`
	Recommendations []*AutomationPerformanceRecommendation `json:"recommendations"`
}

type AutomationPerformanceBottleneck struct {
	Resource    string `json:"resource"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Suggestion  string `json:"suggestion"`
}

type AutomationPerformanceTrend struct {
	Metric    string    `json:"metric"`
	Trend     string    `json:"trend"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Change    float64   `json:"change"`
}

type AutomationPerformanceRecommendation struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Actions     []string `json:"actions"`
	Impact      string   `json:"impact"`
	Resources   []string `json:"resources"`
}

type AutomationSecurityPostureAnalysis struct {
	ClientID         string                                 `json:"client_id"`
	AnalysisTime     time.Time                              `json:"analysis_time"`
	OverallScore     float64                                `json:"overall_score"`
	RiskLevel        string                                 `json:"risk_level"`
	Vulnerabilities  []*Vulnerability                       `json:"vulnerabilities"`
	ComplianceStatus []*AutomationComplianceFrameworkStatus `json:"compliance_status"`
	Recommendations  []*AutomationSecurityRecommendation    `json:"recommendations"`
}

type Vulnerability struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	CVSS        float64   `json:"cvss"`
	FirstSeen   time.Time `json:"first_seen"`
	Status      string    `json:"status"`
}

type AutomationComplianceFrameworkStatus struct {
	Framework     string  `json:"framework"`
	Score         float64 `json:"score"`
	Status        string  `json:"status"`
	Controls      int     `json:"controls"`
	Passed        int     `json:"passed"`
	Failed        int     `json:"failed"`
	NotApplicable int     `json:"not_applicable"`
}

type AutomationSecurityRecommendation struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Actions     []string `json:"actions"`
	Impact      string   `json:"impact"`
	Resources   []string `json:"resources"`
	Compliance  []string `json:"compliance"`
}

type AutomationCostOptimizationRecommendation struct {
	Type                 string   `json:"type"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	Priority             string   `json:"priority"`
	PotentialSavings     float64  `json:"potential_savings"`
	ImplementationEffort string   `json:"implementation_effort"`
	Actions              []string `json:"actions"`
	Resources            []string `json:"resources"`
	RiskLevel            string   `json:"risk_level"`
}

// DiscoveryReport represents a comprehensive discovery report
type DiscoveryReport struct {
	ID               string                     `json:"id"`
	ClientID         string                     `json:"client_id"`
	Title            string                     `json:"title"`
	ExecutiveSummary string                     `json:"executive_summary"`
	Discovery        *EnvironmentDiscovery      `json:"discovery"`
	Analysis         *EnvironmentAnalysis       `json:"analysis"`
	Recommendations  []*AutomatedRecommendation `json:"recommendations"`
	NextSteps        []string                   `json:"next_steps"`
	GeneratedAt      time.Time                  `json:"generated_at"`
}

type EnvironmentAnalysis struct {
	ResourceCount        int      `json:"resource_count"`
	ServiceCount         int      `json:"service_count"`
	EstimatedMonthlyCost float64  `json:"estimated_monthly_cost"`
	SecurityScore        float64  `json:"security_score"`
	ComplianceScore      float64  `json:"compliance_score"`
	OptimizationScore    float64  `json:"optimization_score"`
	KeyFindings          []string `json:"key_findings"`
	RiskAreas            []string `json:"risk_areas"`
}

// ResourceCost represents cost information for a resource
type ResourceCost struct {
	MonthlyCost  float64   `json:"monthly_cost"`
	Currency     string    `json:"currency"`
	PricingModel string    `json:"pricing_model"`
	LastUpdated  time.Time `json:"last_updated"`
}

// Additional interfaces for proactive recommendation engine

type ResponseTimePattern struct {
	AverageResponseTime float64 `json:"average_response_time"`
	PeakResponseTime    float64 `json:"peak_response_time"`
	TrendDirection      string  `json:"trend_direction"`
	PeakHours           []int   `json:"peak_hours"`
}

type ThroughputPattern struct {
	AverageThroughput float64 `json:"average_throughput"`
	PeakThroughput    float64 `json:"peak_throughput"`
	TrendDirection    string  `json:"trend_direction"`
	PeakHours         []int   `json:"peak_hours"`
}

type ErrorPattern struct {
	AverageErrorRate float64  `json:"average_error_rate"`
	PeakErrorRate    float64  `json:"peak_error_rate"`
	TrendDirection   string   `json:"trend_direction"`
	CommonErrors     []string `json:"common_errors"`
}

type ResourceUtilizationPattern struct {
	CPUUtilization    *UtilizationPattern `json:"cpu_utilization"`
	MemoryUtilization *UtilizationPattern `json:"memory_utilization"`
	DiskUtilization   *UtilizationPattern `json:"disk_utilization"`
}

type UtilizationPattern struct {
	Average        float64 `json:"average"`
	Peak           float64 `json:"peak"`
	TrendDirection string  `json:"trend_direction"`
	PeakHours      []int   `json:"peak_hours"`
}

type PerformanceAnomaly struct {
	ID                 string    `json:"id"`
	Type               string    `json:"type"`
	Timestamp          time.Time `json:"timestamp"`
	Severity           string    `json:"severity"`
	Description        string    `json:"description"`
	Value              float64   `json:"value"`
	Expected           float64   `json:"expected"`
	Deviation          float64   `json:"deviation"`
	AffectedComponents []string  `json:"affected_components"`
}

type SecurityMetrics struct {
	VulnerabilityCount   *VulnerabilityCount   `json:"vulnerability_count"`
	ComplianceScore      float64               `json:"compliance_score"`
	SecurityEvents       *SecurityEventMetrics `json:"security_events"`
	AccessControlMetrics *AccessControlMetrics `json:"access_control_metrics"`
}

type VulnerabilityCount struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

type SecurityEventMetrics struct {
	TotalEvents           int     `json:"total_events"`
	HighSeverityEvents    int     `json:"high_severity_events"`
	ResolvedEvents        int     `json:"resolved_events"`
	AverageResolutionTime float64 `json:"average_resolution_time"`
}

type AccessControlMetrics struct {
	TotalUsers        int     `json:"total_users"`
	PrivilegedUsers   int     `json:"privileged_users"`
	InactiveUsers     int     `json:"inactive_users"`
	MFAEnabledUsers   int     `json:"mfa_enabled_users"`
	MFAComplianceRate float64 `json:"mfa_compliance_rate"`
}

type ThreatLandscape struct {
	ActiveThreats     []string          `json:"active_threats"`
	ThreatTrends      map[string]string `json:"threat_trends"`
	GeographicThreats map[string]int    `json:"geographic_threats"`
}

type AutomationComplianceStatus struct {
	Frameworks map[string]*ComplianceFrameworkStatus `json:"frameworks"`
}

type ComplianceFrameworkStatus struct {
	OverallScore    float64   `json:"overall_score"`
	ControlsPassed  int       `json:"controls_passed"`
	ControlsFailed  int       `json:"controls_failed"`
	ControlsPartial int       `json:"controls_partial"`
	LastAssessment  time.Time `json:"last_assessment"`
	NextAssessment  time.Time `json:"next_assessment"`
}

type CostSavings struct {
	MonthlySavings float64 `json:"monthly_savings"`
	AnnualSavings  float64 `json:"annual_savings"`
	Currency       string  `json:"currency"`
	Confidence     float64 `json:"confidence"`
}

type PerformanceImprovement struct {
	ResponseTimeReduction float64 `json:"response_time_reduction"`
	ThroughputIncrease    float64 `json:"throughput_increase"`
	ErrorRateReduction    float64 `json:"error_rate_reduction"`
}

type AutomationCostForecastDetailed struct {
	Period          string   `json:"period"`
	ForecastedCost  float64  `json:"forecasted_cost"`
	ConfidenceLevel float64  `json:"confidence_level"`
	Currency        string   `json:"currency"`
	ForecastFactors []string `json:"forecast_factors"`
}

// Update the existing AutomationPerformancePatternAnalysis to match the implementation
type AutomationPerformancePatternAnalysisUpdated struct {
	ClientID                    string                      `json:"client_id"`
	TimeRange                   AutomationTimeRange         `json:"time_range"`
	ResponseTimePatterns        *ResponseTimePattern        `json:"response_time_patterns"`
	ThroughputPatterns          *ThroughputPattern          `json:"throughput_patterns"`
	ErrorPatterns               *ErrorPattern               `json:"error_patterns"`
	ResourceUtilizationPatterns *ResourceUtilizationPattern `json:"resource_utilization_patterns"`
	PerformanceAnomalies        []*PerformanceAnomaly       `json:"performance_anomalies"`
	Recommendations             []string                    `json:"recommendations"`
}

// Update the existing AutomationSecurityPostureAnalysis to match the implementation
type AutomationSecurityPostureAnalysisUpdated struct {
	ClientID                string                      `json:"client_id"`
	AnalysisTime            time.Time                   `json:"analysis_time"`
	OverallRiskScore        float64                     `json:"overall_risk_score"`
	SecurityMetrics         *SecurityMetrics            `json:"security_metrics"`
	ThreatLandscape         *ThreatLandscape            `json:"threat_landscape"`
	ComplianceStatus        *AutomationComplianceStatus `json:"compliance_status"`
	SecurityRecommendations []string                    `json:"security_recommendations"`
}

// Update the existing AutomationCostOptimizationRecommendation to match the implementation
type AutomationCostOptimizationRecommendationUpdated struct {
	ID                  string                 `json:"id"`
	ClientID            string                 `json:"client_id"`
	Type                string                 `json:"type"`
	Title               string                 `json:"title"`
	Description         string                 `json:"description"`
	Priority            RecommendationPriority `json:"priority"`
	PotentialSavings    *CostSavings           `json:"potential_savings"`
	AffectedResources   []string               `json:"affected_resources"`
	ImplementationSteps []string               `json:"implementation_steps"`
	EstimatedEffort     string                 `json:"estimated_effort"`
	RiskLevel           string                 `json:"risk_level"`
	CreatedAt           time.Time              `json:"created_at"`
	ExpiresAt           *time.Time             `json:"expires_at,omitempty"`
}

// Update the existing AutomationPerformanceRecommendation to match the implementation
type AutomationPerformanceRecommendationUpdated struct {
	ID                  string                  `json:"id"`
	ClientID            string                  `json:"client_id"`
	Type                string                  `json:"type"`
	Title               string                  `json:"title"`
	Description         string                  `json:"description"`
	Priority            RecommendationPriority  `json:"priority"`
	Impact              string                  `json:"impact"`
	AffectedComponents  []string                `json:"affected_components"`
	OptimizationSteps   []string                `json:"optimization_steps"`
	ExpectedImprovement *PerformanceImprovement `json:"expected_improvement"`
	EstimatedEffort     string                  `json:"estimated_effort"`
	CreatedAt           time.Time               `json:"created_at"`
}

// Update the existing AutomationSecurityRecommendation to match the implementation
type AutomationSecurityRecommendationUpdated struct {
	ID               string                 `json:"id"`
	ClientID         string                 `json:"client_id"`
	Type             string                 `json:"type"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Priority         RecommendationPriority `json:"priority"`
	RiskLevel        string                 `json:"risk_level"`
	Impact           string                 `json:"impact"`
	AffectedSystems  []string               `json:"affected_systems"`
	RemediationSteps []string               `json:"remediation_steps"`
	EstimatedEffort  string                 `json:"estimated_effort"`
	Deadline         *time.Time             `json:"deadline,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
}
