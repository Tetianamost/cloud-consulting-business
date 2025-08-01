package interfaces

import (
	"context"
	"time"
)

// DocumentationLink represents a link to official cloud provider documentation
type DocumentationLink struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Provider      string    `json:"provider"` // "AWS", "Azure", "GCP"
	Category      string    `json:"category"` // "compute", "storage", "networking", etc.
	Type          string    `json:"type"`     // "guide", "reference", "tutorial", "best-practice"
	Topic         string    `json:"topic"`
	Audience      string    `json:"audience"`
	Description   string    `json:"description"`
	Tags          []string  `json:"tags"`
	LastChecked   time.Time `json:"last_checked"`
	IsActive      bool      `json:"is_active"`
	IsValid       bool      `json:"is_valid"`
	LastValidated time.Time `json:"last_validated"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DocumentationLibrary manages and validates links to official cloud provider documentation
type DocumentationLibrary interface {
	GetDocumentationLinks(ctx context.Context, provider, topic string) ([]*DocumentationLink, error)
	ValidateLinks(ctx context.Context, links []*DocumentationLink) ([]*LinkValidation, error)
	UpdateDocumentationIndex(ctx context.Context) error
	SearchDocumentation(ctx context.Context, query string, providers []string) ([]*DocumentationLink, error)
	AddDocumentationLink(ctx context.Context, link *DocumentationLink) error
	RemoveDocumentationLink(ctx context.Context, linkID string) error
	GetLinksByCategory(ctx context.Context, category string) ([]*DocumentationLink, error)
	GetLinksByProvider(ctx context.Context, provider string) ([]*DocumentationLink, error)
	GetLinksByType(ctx context.Context, linkType string) ([]*DocumentationLink, error)
	GetLinkValidationStatus(ctx context.Context, linkID string) (*LinkValidation, error)
	IsHealthy() bool
	GetStats() *DocumentationLibraryStats
}

// LinkValidation represents the validation status of a documentation link
type LinkValidation struct {
	LinkID       string        `json:"link_id"`
	URL          string        `json:"url"`
	IsValid      bool          `json:"is_valid"`
	StatusCode   int           `json:"status_code"`
	Error        string        `json:"error,omitempty"`
	ValidatedAt  time.Time     `json:"validated_at"`
	ResponseTime time.Duration `json:"response_time"`
	ContentType  string        `json:"content_type,omitempty"`
	LastModified time.Time     `json:"last_modified,omitempty"`
}

// DocumentationLibraryStats represents statistics about the documentation library
type DocumentationLibraryStats struct {
	TotalLinks          int                   `json:"total_links"`
	ValidLinks          int                   `json:"valid_links"`
	InvalidLinks        int                   `json:"invalid_links"`
	UnvalidatedLinks    int                   `json:"unvalidated_links"`
	LinksByProvider     map[string]int        `json:"links_by_provider"`
	LinksByType         map[string]int        `json:"links_by_type"`
	LinksByCategory     map[string]int        `json:"links_by_category"`
	LastValidationRun   time.Time             `json:"last_validation_run"`
	AverageResponseTime time.Duration         `json:"average_response_time"`
	ValidationErrors    []LinkValidationError `json:"validation_errors"`
}

// LinkValidationError represents an error that occurred during link validation
type LinkValidationError struct {
	LinkID    string    `json:"link_id"`
	URL       string    `json:"url"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

// DocumentationSearchResult represents a search result from the documentation library
type DocumentationSearchResult struct {
	Link      *DocumentationLink `json:"link"`
	Relevance float64            `json:"relevance"`
	Matches   []string           `json:"matches"`
}

// DocumentationFilter represents filters for documentation queries
type DocumentationFilter struct {
	Providers    []string   `json:"providers,omitempty"`
	Types        []string   `json:"types,omitempty"`
	Categories   []string   `json:"categories,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
	ValidOnly    bool       `json:"valid_only"`
	UpdatedSince *time.Time `json:"updated_since,omitempty"`
}

// DocumentationLinkType represents the type of documentation link
type DocumentationLinkType string

const (
	LinkTypeGuide           DocumentationLinkType = "guide"
	LinkTypeReference       DocumentationLinkType = "reference"
	LinkTypeTutorial        DocumentationLinkType = "tutorial"
	LinkTypeBestPractice    DocumentationLinkType = "best-practice"
	LinkTypeAPI             DocumentationLinkType = "api"
	LinkTypePricing         DocumentationLinkType = "pricing"
	LinkTypeSecurity        DocumentationLinkType = "security"
	LinkTypeCompliance      DocumentationLinkType = "compliance"
	LinkTypeArchitecture    DocumentationLinkType = "architecture"
	LinkTypeTroubleshooting DocumentationLinkType = "troubleshooting"
	LinkTypeQuickStart      DocumentationLinkType = "quickstart"
	LinkTypeFAQ             DocumentationLinkType = "faq"
)

// DocumentationCategory represents categories for documentation organization
type DocumentationCategory string

const (
	CategoryGettingStarted  DocumentationCategory = "getting-started"
	CategoryServiceDocs     DocumentationCategory = "service-docs"
	CategoryBestPractices   DocumentationCategory = "best-practices"
	CategorySecurityDocs    DocumentationCategory = "security"
	CategoryCompliance      DocumentationCategory = "compliance"
	CategoryPricing         DocumentationCategory = "pricing"
	CategoryArchitecture    DocumentationCategory = "architecture"
	CategoryMigration       DocumentationCategory = "migration"
	CategoryMonitoringDocs  DocumentationCategory = "monitoring"
	CategoryTroubleshooting DocumentationCategory = "troubleshooting"
	CategorySDK             DocumentationCategory = "sdk"
	CategoryCLI             DocumentationCategory = "cli"
	CategoryWhitepapers     DocumentationCategory = "whitepapers"
	CategoryCaseStudies     DocumentationCategory = "case-studies"
)

// DocumentationAudience represents the target audience for documentation
type DocumentationAudience string

const (
	AudienceTechnical DocumentationAudience = "technical"
	AudienceBusiness  DocumentationAudience = "business"
	AudienceMixed     DocumentationAudience = "mixed"
	AudienceBeginner  DocumentationAudience = "beginner"
	AudienceAdvanced  DocumentationAudience = "advanced"
)
