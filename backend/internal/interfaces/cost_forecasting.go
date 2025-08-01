package interfaces

import (
	"time"
)

// Cost Forecasting types

// ForecastParameters represents parameters for cost forecasting
type ForecastParameters struct {
	ForecastPeriod     string                 `json:"forecast_period"` // "1year", "3years", "5years"
	Granularity        string                 `json:"granularity"`     // "monthly", "quarterly", "yearly"
	GrowthAssumptions  *GrowthAssumptions     `json:"growth_assumptions"`
	ScenarioParameters []*ScenarioParameter   `json:"scenario_parameters"`
	ExternalFactors    []*ExternalFactor      `json:"external_factors"`
	BusinessDrivers    []*BusinessDriver      `json:"business_drivers"`
	ConstraintFactors  []*ConstraintFactor    `json:"constraint_factors"`
	ModelParameters    *ModelParameters       `json:"model_parameters"`
	Metadata           map[string]interface{} `json:"metadata"`
}

// GrowthAssumptions represents growth assumptions for forecasting
type GrowthAssumptions struct {
	UserGrowthRate        float64 `json:"user_growth_rate"`
	DataGrowthRate        float64 `json:"data_growth_rate"`
	TransactionGrowthRate float64 `json:"transaction_growth_rate"`
	ComputeGrowthRate     float64 `json:"compute_growth_rate"`
	StorageGrowthRate     float64 `json:"storage_growth_rate"`
	NetworkGrowthRate     float64 `json:"network_growth_rate"`
	SeasonalityFactor     float64 `json:"seasonality_factor"`
	InflationRate         float64 `json:"inflation_rate"`
}

// ScenarioParameter represents a scenario parameter
type ScenarioParameter struct {
	ScenarioName     string   `json:"scenario_name"`
	Probability      float64  `json:"probability"`
	GrowthMultiplier float64  `json:"growth_multiplier"`
	Description      string   `json:"description"`
	Assumptions      []string `json:"assumptions"`
}

// ExternalFactor represents an external factor affecting costs
type ExternalFactor struct {
	FactorName  string  `json:"factor_name"`
	Impact      string  `json:"impact"` // "positive", "negative", "neutral"
	Magnitude   float64 `json:"magnitude"`
	Probability float64 `json:"probability"`
	Description string  `json:"description"`
	TimeFrame   string  `json:"time_frame"`
}

// BusinessDriver represents a business driver affecting costs
type BusinessDriver struct {
	DriverName      string  `json:"driver_name"`
	DriverType      string  `json:"driver_type"`
	Impact          float64 `json:"impact"`
	Correlation     float64 `json:"correlation"`
	Description     string  `json:"description"`
	MeasurementUnit string  `json:"measurement_unit"`
}

// ConstraintFactor represents a constraint factor
type ConstraintFactor struct {
	ConstraintName string  `json:"constraint_name"`
	ConstraintType string  `json:"constraint_type"`
	LimitValue     float64 `json:"limit_value"`
	Impact         string  `json:"impact"`
	Description    string  `json:"description"`
	Mitigation     string  `json:"mitigation"`
}

// ModelParameters represents model parameters for forecasting
type ModelParameters struct {
	ModelType          string                 `json:"model_type"` // "linear", "exponential", "polynomial", "ml"
	ConfidenceLevel    float64                `json:"confidence_level"`
	SeasonalAdjustment bool                   `json:"seasonal_adjustment"`
	TrendAdjustment    bool                   `json:"trend_adjustment"`
	OutlierDetection   bool                   `json:"outlier_detection"`
	ValidationMethod   string                 `json:"validation_method"`
	CustomParameters   map[string]interface{} `json:"custom_parameters"`
}

// CostForecast represents cost forecast results
type CostForecast struct {
	ID                 string                 `json:"id"`
	ArchitectureID     string                 `json:"architecture_id"`
	ForecastDate       time.Time              `json:"forecast_date"`
	ForecastPeriod     string                 `json:"forecast_period"`
	Granularity        string                 `json:"granularity"`
	BaselineCost       float64                `json:"baseline_cost"`
	ForecastedCost     float64                `json:"forecasted_cost"`
	Currency           string                 `json:"currency"`
	ForecastData       []*ForecastDataPoint   `json:"forecast_data"`
	ServiceForecasts   []*ServiceForecast     `json:"service_forecasts"`
	ScenarioForecasts  []*ScenarioForecast    `json:"scenario_forecasts"`
	CostDriverAnalysis *CostDriverAnalysis    `json:"cost_driver_analysis"`
	ForecastAccuracy   *ForecastAccuracy      `json:"forecast_accuracy"`
	ModelMetadata      *ForecastModelMetadata `json:"model_metadata"`
	Assumptions        []string               `json:"assumptions"`
	Limitations        []string               `json:"limitations"`
	CreatedAt          time.Time              `json:"created_at"`
}

// ForecastDataPoint represents a data point in the forecast
type ForecastDataPoint struct {
	Date           time.Time `json:"date"`
	Period         string    `json:"period"`
	ForecastedCost float64   `json:"forecasted_cost"`
	LowerBound     float64   `json:"lower_bound"`
	UpperBound     float64   `json:"upper_bound"`
	Confidence     float64   `json:"confidence"`
	GrowthRate     float64   `json:"growth_rate"`
	SeasonalFactor float64   `json:"seasonal_factor"`
	TrendComponent float64   `json:"trend_component"`
}

// ServiceForecast represents forecast for a specific service
type ServiceForecast struct {
	ServiceName    string               `json:"service_name"`
	Provider       string               `json:"provider"`
	Category       string               `json:"category"`
	BaselineCost   float64              `json:"baseline_cost"`
	ForecastedCost float64              `json:"forecasted_cost"`
	GrowthRate     float64              `json:"growth_rate"`
	ForecastData   []*ForecastDataPoint `json:"forecast_data"`
	CostDrivers    []*ServiceCostDriver `json:"cost_drivers"`
	Assumptions    []string             `json:"assumptions"`
	RiskFactors    []string             `json:"risk_factors"`
}

// ServiceCostDriver represents a cost driver for a service
type ServiceCostDriver struct {
	DriverName      string  `json:"driver_name"`
	CurrentValue    float64 `json:"current_value"`
	ForecastedValue float64 `json:"forecasted_value"`
	Impact          float64 `json:"impact"`
	Confidence      float64 `json:"confidence"`
	Unit            string  `json:"unit"`
}

// ScenarioForecast represents forecast for different scenarios
type ScenarioForecast struct {
	ScenarioName   string               `json:"scenario_name"`
	Probability    float64              `json:"probability"`
	Description    string               `json:"description"`
	ForecastedCost float64              `json:"forecasted_cost"`
	CostVariance   float64              `json:"cost_variance"`
	ForecastData   []*ForecastDataPoint `json:"forecast_data"`
	KeyAssumptions []string             `json:"key_assumptions"`
	RiskFactors    []string             `json:"risk_factors"`
	Mitigation     []string             `json:"mitigation"`
}

// CostDriverAnalysis represents analysis of cost drivers
type CostDriverAnalysis struct {
	PrimaryCostDrivers   []*CostDriverImpact        `json:"primary_cost_drivers"`
	SecondaryCostDrivers []*CostDriverImpact        `json:"secondary_cost_drivers"`
	DriverCorrelations   []*DriverCorrelation       `json:"driver_correlations"`
	SensitivityAnalysis  *DriverSensitivityAnalysis `json:"sensitivity_analysis"`
	DriverTrends         []*DriverTrend             `json:"driver_trends"`
}

// CostDriverImpact represents the impact of a cost driver
type CostDriverImpact struct {
	DriverName      string  `json:"driver_name"`
	DriverType      string  `json:"driver_type"`
	ImpactMagnitude float64 `json:"impact_magnitude"`
	ImpactDirection string  `json:"impact_direction"`
	Confidence      float64 `json:"confidence"`
	Description     string  `json:"description"`
	Controllability string  `json:"controllability"`
}

// DriverCorrelation represents correlation between cost drivers
type DriverCorrelation struct {
	Driver1          string  `json:"driver1"`
	Driver2          string  `json:"driver2"`
	CorrelationCoeff float64 `json:"correlation_coefficient"`
	Significance     string  `json:"significance"`
	Relationship     string  `json:"relationship"`
}

// DriverSensitivityAnalysis represents sensitivity analysis of drivers
type DriverSensitivityAnalysis struct {
	SensitivityResults []*DriverSensitivityResult `json:"sensitivity_results"`
	TornadoChart       *TornadoChartData          `json:"tornado_chart"`
	KeyInsights        []string                   `json:"key_insights"`
}

// DriverSensitivityResult represents sensitivity result for a driver
type DriverSensitivityResult struct {
	DriverName       string  `json:"driver_name"`
	BaseValue        float64 `json:"base_value"`
	LowValue         float64 `json:"low_value"`
	HighValue        float64 `json:"high_value"`
	LowImpact        float64 `json:"low_impact"`
	HighImpact       float64 `json:"high_impact"`
	SensitivityIndex float64 `json:"sensitivity_index"`
}

// TornadoChartData represents data for tornado chart
type TornadoChartData struct {
	ChartData  []*TornadoChartItem `json:"chart_data"`
	BaseValue  float64             `json:"base_value"`
	Title      string              `json:"title"`
	XAxisLabel string              `json:"x_axis_label"`
	YAxisLabel string              `json:"y_axis_label"`
}

// TornadoChartItem represents an item in tornado chart
type TornadoChartItem struct {
	DriverName string  `json:"driver_name"`
	LowImpact  float64 `json:"low_impact"`
	HighImpact float64 `json:"high_impact"`
	Range      float64 `json:"range"`
	Rank       int     `json:"rank"`
}

// DriverTrend represents trend for a cost driver
type DriverTrend struct {
	DriverName     string                 `json:"driver_name"`
	TrendDirection string                 `json:"trend_direction"`
	TrendMagnitude float64                `json:"trend_magnitude"`
	TrendData      []*DriverTrendPoint    `json:"trend_data"`
	Seasonality    *DriverSeasonality     `json:"seasonality"`
	Forecast       []*DriverForecastPoint `json:"forecast"`
}

// DriverTrendPoint represents a trend point for a driver
type DriverTrendPoint struct {
	Date     time.Time `json:"date"`
	Value    float64   `json:"value"`
	Trend    float64   `json:"trend"`
	Seasonal float64   `json:"seasonal"`
}

// DriverSeasonality represents seasonality for a driver
type DriverSeasonality struct {
	HasSeasonality   bool                    `json:"has_seasonality"`
	SeasonalStrength float64                 `json:"seasonal_strength"`
	SeasonalPattern  []*SeasonalPatternPoint `json:"seasonal_pattern"`
	PeakPeriods      []string                `json:"peak_periods"`
	LowPeriods       []string                `json:"low_periods"`
}

// SeasonalPatternPoint represents a seasonal pattern point
type SeasonalPatternPoint struct {
	Period   string  `json:"period"`
	Factor   float64 `json:"factor"`
	Strength float64 `json:"strength"`
}

// DriverForecastPoint represents a forecast point for a driver
type DriverForecastPoint struct {
	Date          time.Time `json:"date"`
	ForecastValue float64   `json:"forecast_value"`
	LowerBound    float64   `json:"lower_bound"`
	UpperBound    float64   `json:"upper_bound"`
	Confidence    float64   `json:"confidence"`
}

// ForecastAccuracy represents forecast accuracy metrics
type ForecastAccuracy struct {
	HistoricalAccuracy *HistoricalAccuracy `json:"historical_accuracy"`
	ValidationResults  *ValidationResults  `json:"validation_results"`
	AccuracyMetrics    *AccuracyMetrics    `json:"accuracy_metrics"`
	ModelPerformance   *ModelPerformance   `json:"model_performance"`
	AccuracyTrends     []*AccuracyTrend    `json:"accuracy_trends"`
}

// HistoricalAccuracy represents historical accuracy
type HistoricalAccuracy struct {
	PeriodAnalyzed    string                 `json:"period_analyzed"`
	AccuracyScore     float64                `json:"accuracy_score"`
	AccuracyByPeriod  []*AccuracyByPeriod    `json:"accuracy_by_period"`
	AccuracyByService []*AccuracyByService   `json:"accuracy_by_service"`
	ImprovementTrends []*AccuracyImprovement `json:"improvement_trends"`
}

// AccuracyByPeriod represents accuracy by time period
type AccuracyByPeriod struct {
	Period          string  `json:"period"`
	ActualCost      float64 `json:"actual_cost"`
	ForecastedCost  float64 `json:"forecasted_cost"`
	AccuracyPercent float64 `json:"accuracy_percent"`
	AbsoluteError   float64 `json:"absolute_error"`
	RelativeError   float64 `json:"relative_error"`
}

// AccuracyByService represents accuracy by service
type AccuracyByService struct {
	ServiceName    string  `json:"service_name"`
	AccuracyScore  float64 `json:"accuracy_score"`
	ErrorMagnitude float64 `json:"error_magnitude"`
	ErrorDirection string  `json:"error_direction"`
	Reliability    string  `json:"reliability"`
}

// AccuracyImprovement represents accuracy improvement over time
type AccuracyImprovement struct {
	Period             string   `json:"period"`
	AccuracyScore      float64  `json:"accuracy_score"`
	ImprovementPercent float64  `json:"improvement_percent"`
	ImprovementFactors []string `json:"improvement_factors"`
}

// ValidationResults represents validation results
type ValidationResults struct {
	ValidationMethod   string                  `json:"validation_method"`
	ValidationPeriod   string                  `json:"validation_period"`
	ValidationScore    float64                 `json:"validation_score"`
	CrossValidation    *CrossValidationResults `json:"cross_validation"`
	OutOfSampleTesting *OutOfSampleResults     `json:"out_of_sample_testing"`
	ValidationInsights []string                `json:"validation_insights"`
}

// CrossValidationResults represents cross-validation results
type CrossValidationResults struct {
	FoldCount     int                    `json:"fold_count"`
	AverageScore  float64                `json:"average_score"`
	ScoreVariance float64                `json:"score_variance"`
	FoldResults   []*CrossValidationFold `json:"fold_results"`
	Stability     string                 `json:"stability"`
}

// CrossValidationFold represents a cross-validation fold
type CrossValidationFold struct {
	FoldNumber  int     `json:"fold_number"`
	Score       float64 `json:"score"`
	TrainSize   int     `json:"train_size"`
	TestSize    int     `json:"test_size"`
	Performance string  `json:"performance"`
}

// OutOfSampleResults represents out-of-sample testing results
type OutOfSampleResults struct {
	TestPeriod      string   `json:"test_period"`
	TestScore       float64  `json:"test_score"`
	Overfitting     bool     `json:"overfitting"`
	Generalization  string   `json:"generalization"`
	Recommendations []string `json:"recommendations"`
}

// AccuracyMetrics represents various accuracy metrics
type AccuracyMetrics struct {
	MAE        float64 `json:"mae"`      // Mean Absolute Error
	MAPE       float64 `json:"mape"`     // Mean Absolute Percentage Error
	RMSE       float64 `json:"rmse"`     // Root Mean Square Error
	R2Score    float64 `json:"r2_score"` // R-squared
	AdjustedR2 float64 `json:"adjusted_r2"`
	AIC        float64 `json:"aic"`     // Akaike Information Criterion
	BIC        float64 `json:"bic"`     // Bayesian Information Criterion
	TheilU     float64 `json:"theil_u"` // Theil's U statistic
}

// ModelPerformance represents model performance metrics
type ModelPerformance struct {
	ModelType         string               `json:"model_type"`
	TrainingTime      string               `json:"training_time"`
	PredictionTime    string               `json:"prediction_time"`
	ModelComplexity   string               `json:"model_complexity"`
	FeatureImportance []*FeatureImportance `json:"feature_importance"`
	ModelDiagnostics  *ModelDiagnostics    `json:"model_diagnostics"`
	PerformanceScore  float64              `json:"performance_score"`
}

// FeatureImportance represents feature importance in the model
type FeatureImportance struct {
	FeatureName string  `json:"feature_name"`
	Importance  float64 `json:"importance"`
	Rank        int     `json:"rank"`
	Coefficient float64 `json:"coefficient"`
	PValue      float64 `json:"p_value"`
}

// ModelDiagnostics represents model diagnostics
type ModelDiagnostics struct {
	Residuals          *ResidualAnalysis       `json:"residuals"`
	Autocorrelation    *AutocorrelationTest    `json:"autocorrelation"`
	Heteroscedasticity *HeteroscedasticityTest `json:"heteroscedasticity"`
	Normality          *NormalityTest          `json:"normality"`
	Multicollinearity  *MulticollinearityTest  `json:"multicollinearity"`
}

// ResidualAnalysis represents residual analysis
type ResidualAnalysis struct {
	MeanResidual         float64           `json:"mean_residual"`
	ResidualVariance     float64           `json:"residual_variance"`
	ResidualDistribution string            `json:"residual_distribution"`
	OutlierCount         int               `json:"outlier_count"`
	ResidualPlot         *ResidualPlotData `json:"residual_plot"`
}

// ResidualPlotData represents residual plot data
type ResidualPlotData struct {
	PlotType   string           `json:"plot_type"`
	DataPoints []*ResidualPoint `json:"data_points"`
	TrendLine  *TrendLineData   `json:"trend_line"`
	Outliers   []*OutlierPoint  `json:"outliers"`
}

// ResidualPoint represents a residual point
type ResidualPoint struct {
	Predicted    float64 `json:"predicted"`
	Residual     float64 `json:"residual"`
	Standardized float64 `json:"standardized"`
}

// TrendLineData represents trend line data
type TrendLineData struct {
	Slope        float64 `json:"slope"`
	Intercept    float64 `json:"intercept"`
	RSquared     float64 `json:"r_squared"`
	Significance string  `json:"significance"`
}

// OutlierPoint represents an outlier point
type OutlierPoint struct {
	Index     int     `json:"index"`
	Value     float64 `json:"value"`
	Residual  float64 `json:"residual"`
	Influence float64 `json:"influence"`
	Leverage  float64 `json:"leverage"`
}

// AutocorrelationTest represents autocorrelation test results
type AutocorrelationTest struct {
	TestName        string    `json:"test_name"`
	TestStatistic   float64   `json:"test_statistic"`
	PValue          float64   `json:"p_value"`
	CriticalValue   float64   `json:"critical_value"`
	Conclusion      string    `json:"conclusion"`
	LagCorrelations []float64 `json:"lag_correlations"`
}

// HeteroscedasticityTest represents heteroscedasticity test results
type HeteroscedasticityTest struct {
	TestName       string  `json:"test_name"`
	TestStatistic  float64 `json:"test_statistic"`
	PValue         float64 `json:"p_value"`
	Conclusion     string  `json:"conclusion"`
	Recommendation string  `json:"recommendation"`
}

// NormalityTest represents normality test results
type NormalityTest struct {
	TestName      string  `json:"test_name"`
	TestStatistic float64 `json:"test_statistic"`
	PValue        float64 `json:"p_value"`
	Conclusion    string  `json:"conclusion"`
	Skewness      float64 `json:"skewness"`
	Kurtosis      float64 `json:"kurtosis"`
}

// MulticollinearityTest represents multicollinearity test results
type MulticollinearityTest struct {
	VIFScores       []*VIFScore `json:"vif_scores"`
	ConditionNumber float64     `json:"condition_number"`
	Conclusion      string      `json:"conclusion"`
	Recommendations []string    `json:"recommendations"`
}

// VIFScore represents Variance Inflation Factor score
type VIFScore struct {
	FeatureName string  `json:"feature_name"`
	VIFValue    float64 `json:"vif_value"`
	Status      string  `json:"status"`
}

// AccuracyTrend represents accuracy trend over time
type AccuracyTrend struct {
	Period        string   `json:"period"`
	AccuracyScore float64  `json:"accuracy_score"`
	Trend         string   `json:"trend"`
	Improvement   float64  `json:"improvement"`
	Factors       []string `json:"factors"`
}

// ForecastModelMetadata represents metadata about the forecast model
type ForecastModelMetadata struct {
	ModelVersion          string                 `json:"model_version"`
	ModelType             string                 `json:"model_type"`
	TrainingDataPeriod    string                 `json:"training_data_period"`
	FeatureCount          int                    `json:"feature_count"`
	TrainingDataPoints    int                    `json:"training_data_points"`
	ModelParameters       map[string]interface{} `json:"model_parameters"`
	LastTrainingDate      time.Time              `json:"last_training_date"`
	ModelPerformanceScore float64                `json:"model_performance_score"`
	DataQualityScore      float64                `json:"data_quality_score"`
	ModelReliability      string                 `json:"model_reliability"`
}
