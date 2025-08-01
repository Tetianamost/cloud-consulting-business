package interfaces

import (
	"time"
)

// Right-sizing Analysis types

// ResourceUtilizationData represents resource utilization data for right-sizing analysis
type ResourceUtilizationData struct {
	AccountID         string                         `json:"account_id"`
	Region            string                         `json:"region"`
	TimeRange         *TimeRange                     `json:"time_range"`
	ComputeResources  []*ComputeResourceUtilization  `json:"compute_resources"`
	StorageResources  []*StorageResourceUtilization  `json:"storage_resources"`
	DatabaseResources []*DatabaseResourceUtilization `json:"database_resources"`
	NetworkResources  []*NetworkResourceUtilization  `json:"network_resources"`
	Metadata          map[string]interface{}         `json:"metadata"`
}

// ComputeResourceUtilization represents compute resource utilization
type ComputeResourceUtilization struct {
	ResourceID         string                 `json:"resource_id"`
	InstanceType       string                 `json:"instance_type"`
	Region             string                 `json:"region"`
	Platform           string                 `json:"platform"`
	CPUUtilization     *UtilizationMetrics    `json:"cpu_utilization"`
	MemoryUtilization  *UtilizationMetrics    `json:"memory_utilization"`
	NetworkUtilization *UtilizationMetrics    `json:"network_utilization"`
	StorageUtilization *UtilizationMetrics    `json:"storage_utilization"`
	CostData           *ResourceCostData      `json:"cost_data"`
	PerformanceData    *PerformanceData       `json:"performance_data"`
	UsagePatterns      *ResourceUsagePatterns `json:"usage_patterns"`
	Tags               map[string]string      `json:"tags"`
}

// StorageResourceUtilization represents storage resource utilization
type StorageResourceUtilization struct {
	ResourceID            string              `json:"resource_id"`
	StorageType           string              `json:"storage_type"`
	Region                string              `json:"region"`
	CapacityGB            float64             `json:"capacity_gb"`
	UsedGB                float64             `json:"used_gb"`
	UtilizationRate       float64             `json:"utilization_rate"`
	IOPSUtilization       *UtilizationMetrics `json:"iops_utilization"`
	ThroughputUtilization *UtilizationMetrics `json:"throughput_utilization"`
	CostData              *ResourceCostData   `json:"cost_data"`
	GrowthTrend           *GrowthTrend        `json:"growth_trend"`
	AccessPatterns        *AccessPatterns     `json:"access_patterns"`
	Tags                  map[string]string   `json:"tags"`
}

// DatabaseResourceUtilization represents database resource utilization
type DatabaseResourceUtilization struct {
	ResourceID            string                      `json:"resource_id"`
	Engine                string                      `json:"engine"`
	InstanceClass         string                      `json:"instance_class"`
	Region                string                      `json:"region"`
	CPUUtilization        *UtilizationMetrics         `json:"cpu_utilization"`
	MemoryUtilization     *UtilizationMetrics         `json:"memory_utilization"`
	StorageUtilization    *UtilizationMetrics         `json:"storage_utilization"`
	ConnectionUtilization *UtilizationMetrics         `json:"connection_utilization"`
	CostData              *ResourceCostData           `json:"cost_data"`
	PerformanceMetrics    *DatabasePerformanceMetrics `json:"performance_metrics"`
	UsagePatterns         *ResourceUsagePatterns      `json:"usage_patterns"`
	Tags                  map[string]string           `json:"tags"`
}

// NetworkResourceUtilization represents network resource utilization
type NetworkResourceUtilization struct {
	ResourceID            string              `json:"resource_id"`
	ResourceType          string              `json:"resource_type"`
	Region                string              `json:"region"`
	BandwidthUtilization  *UtilizationMetrics `json:"bandwidth_utilization"`
	PacketUtilization     *UtilizationMetrics `json:"packet_utilization"`
	ConnectionUtilization *UtilizationMetrics `json:"connection_utilization"`
	CostData              *ResourceCostData   `json:"cost_data"`
	TrafficPatterns       *TrafficPatterns    `json:"traffic_patterns"`
	Tags                  map[string]string   `json:"tags"`
}

// UtilizationMetrics represents utilization metrics
type UtilizationMetrics struct {
	Average    float64   `json:"average"`
	Maximum    float64   `json:"maximum"`
	Minimum    float64   `json:"minimum"`
	P95        float64   `json:"p95"`
	P99        float64   `json:"p99"`
	Samples    int       `json:"samples"`
	DataPoints []float64 `json:"data_points"`
	Trend      string    `json:"trend"`
}

// ResourceCostData represents cost data for a resource
type ResourceCostData struct {
	MonthlyCost    float64 `json:"monthly_cost"`
	DailyCost      float64 `json:"daily_cost"`
	HourlyCost     float64 `json:"hourly_cost"`
	CostTrend      string  `json:"cost_trend"`
	CostEfficiency float64 `json:"cost_efficiency"`
	WastedSpend    float64 `json:"wasted_spend"`
}

// PerformanceData represents performance data
type PerformanceData struct {
	ResponseTime     float64 `json:"response_time"`
	Throughput       float64 `json:"throughput"`
	ErrorRate        float64 `json:"error_rate"`
	Availability     float64 `json:"availability"`
	PerformanceScore float64 `json:"performance_score"`
}

// ResourceUsagePatterns represents usage patterns for a resource
type ResourceUsagePatterns struct {
	Pattern         string            `json:"pattern"`
	PeakHours       []int             `json:"peak_hours"`
	PeakDays        []int             `json:"peak_days"`
	SeasonalFactors []*SeasonalFactor `json:"seasonal_factors"`
	Predictability  float64           `json:"predictability"`
	Variability     float64           `json:"variability"`
}

// GrowthTrend represents growth trend data
type GrowthTrend struct {
	GrowthRate      float64 `json:"growth_rate"`
	Direction       string  `json:"direction"`
	Confidence      float64 `json:"confidence"`
	ProjectedGrowth float64 `json:"projected_growth"`
}

// AccessPatterns represents access patterns for storage
type AccessPatterns struct {
	ReadFrequency   float64 `json:"read_frequency"`
	WriteFrequency  float64 `json:"write_frequency"`
	AccessPattern   string  `json:"access_pattern"`
	HotDataPercent  float64 `json:"hot_data_percent"`
	ColdDataPercent float64 `json:"cold_data_percent"`
}

// DatabasePerformanceMetrics represents database performance metrics
type DatabasePerformanceMetrics struct {
	QueryLatency    float64 `json:"query_latency"`
	TransactionRate float64 `json:"transaction_rate"`
	LockWaitTime    float64 `json:"lock_wait_time"`
	BufferHitRatio  float64 `json:"buffer_hit_ratio"`
	IndexEfficiency float64 `json:"index_efficiency"`
	ReplicationLag  float64 `json:"replication_lag"`
}

// TrafficPatterns represents network traffic patterns
type TrafficPatterns struct {
	InboundTraffic   float64 `json:"inbound_traffic"`
	OutboundTraffic  float64 `json:"outbound_traffic"`
	PeakTraffic      float64 `json:"peak_traffic"`
	TrafficPattern   string  `json:"traffic_pattern"`
	DataTransferCost float64 `json:"data_transfer_cost"`
}

// RightSizingAnalysis represents right-sizing analysis results
type RightSizingAnalysis struct {
	ID                            string                         `json:"id"`
	AnalysisDate                  time.Time                      `json:"analysis_date"`
	TotalResourcesAnalyzed        int                            `json:"total_resources_analyzed"`
	OverProvisionedResources      int                            `json:"over_provisioned_resources"`
	UnderProvisionedResources     int                            `json:"under_provisioned_resources"`
	OptimallyProvisionedResources int                            `json:"optimally_provisioned_resources"`
	TotalSavingsPotential         float64                        `json:"total_savings_potential"`
	ResourceAnalysis              []*ResourceRightSizingAnalysis `json:"resource_analysis"`
	UtilizationSummary            *UtilizationSummary            `json:"utilization_summary"`
	CostImpactAnalysis            *CostImpactAnalysis            `json:"cost_impact_analysis"`
	PerformanceImpactAnalysis     *PerformanceImpactAnalysis     `json:"performance_impact_analysis"`
	CreatedAt                     time.Time                      `json:"created_at"`
}

// ResourceRightSizingAnalysis represents right-sizing analysis for a specific resource
type ResourceRightSizingAnalysis struct {
	ResourceID           string                       `json:"resource_id"`
	ResourceType         string                       `json:"resource_type"`
	CurrentConfiguration *ResourceConfiguration       `json:"current_configuration"`
	UtilizationAnalysis  *ResourceUtilizationAnalysis `json:"utilization_analysis"`
	RightSizingStatus    string                       `json:"right_sizing_status"`
	RecommendedAction    string                       `json:"recommended_action"`
	SavingsPotential     float64                      `json:"savings_potential"`
	PerformanceImpact    string                       `json:"performance_impact"`
	RiskLevel            string                       `json:"risk_level"`
	ConfidenceLevel      float64                      `json:"confidence_level"`
	Justification        string                       `json:"justification"`
}

// ResourceConfiguration represents resource configuration
type ResourceConfiguration struct {
	InstanceType    string                 `json:"instance_type"`
	CPU             int                    `json:"cpu"`
	Memory          float64                `json:"memory"`
	Storage         float64                `json:"storage"`
	NetworkCapacity float64                `json:"network_capacity"`
	Specifications  map[string]interface{} `json:"specifications"`
	MonthlyCost     float64                `json:"monthly_cost"`
}

// ResourceUtilizationAnalysis represents utilization analysis for a resource
type ResourceUtilizationAnalysis struct {
	CPUAnalysis        *UtilizationAnalysisDetail `json:"cpu_analysis"`
	MemoryAnalysis     *UtilizationAnalysisDetail `json:"memory_analysis"`
	StorageAnalysis    *UtilizationAnalysisDetail `json:"storage_analysis"`
	NetworkAnalysis    *UtilizationAnalysisDetail `json:"network_analysis"`
	OverallStatus      string                     `json:"overall_status"`
	BottleneckAnalysis *BottleneckAnalysis        `json:"bottleneck_analysis"`
}

// UtilizationAnalysisDetail represents detailed utilization analysis
type UtilizationAnalysisDetail struct {
	AverageUtilization float64 `json:"average_utilization"`
	PeakUtilization    float64 `json:"peak_utilization"`
	UtilizationTrend   string  `json:"utilization_trend"`
	Status             string  `json:"status"` // "over_provisioned", "under_provisioned", "optimal"
	Recommendation     string  `json:"recommendation"`
	ConfidenceLevel    float64 `json:"confidence_level"`
}

// BottleneckAnalysis represents bottleneck analysis
type BottleneckAnalysis struct {
	PrimaryBottleneck    string   `json:"primary_bottleneck"`
	SecondaryBottlenecks []string `json:"secondary_bottlenecks"`
	BottleneckImpact     string   `json:"bottleneck_impact"`
	ResolutionPriority   string   `json:"resolution_priority"`
}

// UtilizationSummary represents utilization summary
type UtilizationSummary struct {
	OverallUtilization float64              `json:"overall_utilization"`
	UtilizationByType  []*UtilizationByType `json:"utilization_by_type"`
	UtilizationTrends  []*UtilizationTrend  `json:"utilization_trends"`
	WasteAnalysis      *WasteAnalysis       `json:"waste_analysis"`
	EfficiencyMetrics  *EfficiencyMetrics   `json:"efficiency_metrics"`
}

// UtilizationByType represents utilization by resource type
type UtilizationByType struct {
	ResourceType       string  `json:"resource_type"`
	AverageUtilization float64 `json:"average_utilization"`
	ResourceCount      int     `json:"resource_count"`
	WastedCapacity     float64 `json:"wasted_capacity"`
	SavingsPotential   float64 `json:"savings_potential"`
}

// UtilizationTrend represents utilization trend
type UtilizationTrend struct {
	Period      string  `json:"period"`
	Utilization float64 `json:"utilization"`
	Trend       string  `json:"trend"`
	Forecast    float64 `json:"forecast"`
}

// WasteAnalysis represents waste analysis
type WasteAnalysis struct {
	TotalWastedCapacity float64            `json:"total_wasted_capacity"`
	TotalWastedCost     float64            `json:"total_wasted_cost"`
	WasteByCategory     []*WasteByCategory `json:"waste_by_category"`
	WasteReduction      *WasteReduction    `json:"waste_reduction"`
}

// WasteByCategory represents waste by category
type WasteByCategory struct {
	Category      string  `json:"category"`
	WastedCost    float64 `json:"wasted_cost"`
	WastedPercent float64 `json:"wasted_percent"`
	Opportunity   string  `json:"opportunity"`
}

// WasteReduction represents waste reduction opportunities
type WasteReduction struct {
	ImmediateReduction float64 `json:"immediate_reduction"`
	ShortTermReduction float64 `json:"short_term_reduction"`
	LongTermReduction  float64 `json:"long_term_reduction"`
	TotalReduction     float64 `json:"total_reduction"`
}

// EfficiencyMetrics represents efficiency metrics
type EfficiencyMetrics struct {
	CostEfficiency        float64 `json:"cost_efficiency"`
	ResourceEfficiency    float64 `json:"resource_efficiency"`
	PerformanceEfficiency float64 `json:"performance_efficiency"`
	OverallEfficiency     float64 `json:"overall_efficiency"`
	EfficiencyScore       string  `json:"efficiency_score"`
}

// CostImpactAnalysis represents cost impact analysis
type CostImpactAnalysis struct {
	CurrentMonthlyCost   float64                 `json:"current_monthly_cost"`
	OptimizedMonthlyCost float64                 `json:"optimized_monthly_cost"`
	MonthlySavings       float64                 `json:"monthly_savings"`
	AnnualSavings        float64                 `json:"annual_savings"`
	SavingsPercentage    float64                 `json:"savings_percentage"`
	CostImpactByResource []*CostImpactByResource `json:"cost_impact_by_resource"`
	ImplementationCost   float64                 `json:"implementation_cost"`
	PaybackPeriod        string                  `json:"payback_period"`
	ROI                  float64                 `json:"roi"`
}

// CostImpactByResource represents cost impact by resource
type CostImpactByResource struct {
	ResourceID         string  `json:"resource_id"`
	ResourceType       string  `json:"resource_type"`
	CurrentCost        float64 `json:"current_cost"`
	OptimizedCost      float64 `json:"optimized_cost"`
	Savings            float64 `json:"savings"`
	SavingsPercent     float64 `json:"savings_percent"`
	ImplementationCost float64 `json:"implementation_cost"`
}

// PerformanceImpactAnalysis represents performance impact analysis
type PerformanceImpactAnalysis struct {
	OverallPerformanceImpact    string                         `json:"overall_performance_impact"`
	PerformanceRisk             string                         `json:"performance_risk"`
	PerformanceImpactByResource []*PerformanceImpactByResource `json:"performance_impact_by_resource"`
	MitigationStrategies        []string                       `json:"mitigation_strategies"`
	MonitoringRecommendations   []string                       `json:"monitoring_recommendations"`
}

// PerformanceImpactByResource represents performance impact by resource
type PerformanceImpactByResource struct {
	ResourceID           string   `json:"resource_id"`
	ResourceType         string   `json:"resource_type"`
	CurrentPerformance   float64  `json:"current_performance"`
	ProjectedPerformance float64  `json:"projected_performance"`
	PerformanceChange    float64  `json:"performance_change"`
	ImpactLevel          string   `json:"impact_level"`
	RiskFactors          []string `json:"risk_factors"`
}
