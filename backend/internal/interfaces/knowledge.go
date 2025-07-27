package interfaces

import (
	"context"
	"time"
)

// KnowledgeBase defines the interface for cloud knowledge management
type KnowledgeBase interface {
	GetCloudServiceInfo(provider, service string) (*CloudServiceInfo, error)
	GetBestPractices(category, provider string) ([]*BestPractice, error)
	GetComplianceRequirements(industry string) ([]*ComplianceRequirement, error)
	GetArchitecturalPatterns(useCase, provider string) ([]*ArchitecturalPattern, error)
	GetDocumentationLinks(provider, topic string) ([]*DocumentationLink, error)
	UpdateKnowledgeBase(ctx context.Context) error
	SearchServices(ctx context.Context, query string, providers []string) ([]*CloudServiceInfo, error)
	GetServiceAlternatives(provider, service string) (map[string]string, error)
	IsHealthy() bool
}

// CloudServiceInfo represents information about a cloud service
type CloudServiceInfo struct {
	Provider         string            `json:"provider"`
	ServiceName      string            `json:"service_name"`
	Category         string            `json:"category"`
	Description      string            `json:"description"`
	UseCases         []string          `json:"use_cases"`
	PricingModel     string            `json:"pricing_model"`
	DocumentationURL string            `json:"documentation_url"`
	BestPracticesURL string            `json:"best_practices_url"`
	Alternatives     map[string]string `json:"alternatives"` // provider -> service mapping
	Features         []string          `json:"features"`
	Limitations      []string          `json:"limitations"`
	LastUpdated      time.Time         `json:"last_updated"`
}

// BestPractice represents a cloud best practice recommendation
type BestPractice struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	Provider         string    `json:"provider"`
	Industry         string    `json:"industry,omitempty"`
	DocumentationURL string    `json:"documentation_url"`
	Priority         string    `json:"priority"` // "high", "medium", "low"
	Tags             []string  `json:"tags"`
	Implementation   string    `json:"implementation"`
	Benefits         []string  `json:"benefits"`
	Considerations   []string  `json:"considerations"`
	LastUpdated      time.Time `json:"last_updated"`
}

// ComplianceRequirement represents industry-specific compliance requirements
type ComplianceRequirement struct {
	Framework        string              `json:"framework"` // "HIPAA", "PCI-DSS", "SOX", etc.
	Industry         string              `json:"industry"`
	Requirement      string              `json:"requirement"`
	Description      string              `json:"description"`
	CloudControls    map[string][]string `json:"cloud_controls"` // provider -> controls
	DocumentationURL string              `json:"documentation_url"`
	Severity         string              `json:"severity"`
	Implementation   []string            `json:"implementation"`
	ValidationSteps  []string            `json:"validation_steps"`
	LastUpdated      time.Time           `json:"last_updated"`
}

// ArchitecturalPattern represents a cloud architectural pattern
type ArchitecturalPattern struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	UseCase          string            `json:"use_case"`
	Provider         string            `json:"provider"`
	Components       []string          `json:"components"`
	Benefits         []string          `json:"benefits"`
	Drawbacks        []string          `json:"drawbacks"`
	Implementation   string            `json:"implementation"`
	DocumentationURL string            `json:"documentation_url"`
	DiagramURL       string            `json:"diagram_url,omitempty"`
	EstimatedCost    string            `json:"estimated_cost"`
	Complexity       string            `json:"complexity"` // "low", "medium", "high"
	Tags             []string          `json:"tags"`
	Alternatives     map[string]string `json:"alternatives"` // provider -> pattern mapping
	LastUpdated      time.Time         `json:"last_updated"`
}

// DocumentationLink represents a link to official cloud provider documentation
type DocumentationLink struct {
	ID           string    `json:"id"`
	Provider     string    `json:"provider"`
	Topic        string    `json:"topic"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Description  string    `json:"description"`
	Type         string    `json:"type"` // "guide", "reference", "tutorial", "best-practice"
	LastValidated time.Time `json:"last_validated"`
	IsValid      bool      `json:"is_valid"`
	Tags         []string  `json:"tags"`
	Category     string    `json:"category"`
	Audience     string    `json:"audience"` // "technical", "business", "mixed"
}

// KnowledgeBaseStats represents statistics about the knowledge base
type KnowledgeBaseStats struct {
	TotalServices           int       `json:"total_services"`
	TotalBestPractices      int       `json:"total_best_practices"`
	TotalComplianceRules    int       `json:"total_compliance_rules"`
	TotalArchitecturalPatterns int    `json:"total_architectural_patterns"`
	TotalDocumentationLinks int       `json:"total_documentation_links"`
	ProvidersSupported      []string  `json:"providers_supported"`
	IndustriesSupported     []string  `json:"industries_supported"`
	LastUpdated             time.Time `json:"last_updated"`
}

// ServiceCategory represents categories of cloud services
type ServiceCategory string

const (
	CategoryCompute     ServiceCategory = "compute"
	CategoryStorage     ServiceCategory = "storage"
	CategoryDatabase    ServiceCategory = "database"
	CategoryNetworking  ServiceCategory = "networking"
	CategorySecurity    ServiceCategory = "security"
	CategoryAnalytics   ServiceCategory = "analytics"
	CategoryAI          ServiceCategory = "ai"
	CategoryDevOps      ServiceCategory = "devops"
	CategoryMonitoring  ServiceCategory = "monitoring"
	CategoryIntegration ServiceCategory = "integration"
	CategoryIoT         ServiceCategory = "iot"
	CategoryMobile      ServiceCategory = "mobile"
	CategoryWeb         ServiceCategory = "web"
	CategoryContainers  ServiceCategory = "containers"
	CategoryServerless  ServiceCategory = "serverless"
)

// CloudProvider represents supported cloud providers
type CloudProvider string

const (
	ProviderAWS     CloudProvider = "aws"
	ProviderAzure   CloudProvider = "azure"
	ProviderGCP     CloudProvider = "gcp"
	ProviderAlibaba CloudProvider = "alibaba"
	ProviderOracle  CloudProvider = "oracle"
	ProviderIBM     CloudProvider = "ibm"
)

// ComplianceFramework represents supported compliance frameworks
type ComplianceFramework string

const (
	FrameworkHIPAA    ComplianceFramework = "HIPAA"
	FrameworkPCIDSS   ComplianceFramework = "PCI-DSS"
	FrameworkSOX      ComplianceFramework = "SOX"
	FrameworkGDPR     ComplianceFramework = "GDPR"
	FrameworkSOC2     ComplianceFramework = "SOC2"
	FrameworkISO27001 ComplianceFramework = "ISO27001"
	FrameworkFedRAMP  ComplianceFramework = "FedRAMP"
	FrameworkNIST     ComplianceFramework = "NIST"
)

// Industry represents supported industries
type Industry string

const (
	IndustryHealthcare    Industry = "healthcare"
	IndustryFinancial     Industry = "financial"
	IndustryRetail        Industry = "retail"
	IndustryManufacturing Industry = "manufacturing"
	IndustryEducation     Industry = "education"
	IndustryGovernment    Industry = "government"
	IndustryTechnology    Industry = "technology"
	IndustryMedia         Industry = "media"
	IndustryEnergy        Industry = "energy"
	IndustryTelecom       Industry = "telecom"
)