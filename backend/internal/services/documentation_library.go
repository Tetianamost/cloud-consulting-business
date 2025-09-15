package services

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// DocumentationLibraryService implements the DocumentationLibrary interface
type DocumentationLibraryService struct {
	links       map[string]*interfaces.DocumentationLink
	validations map[string]*interfaces.LinkValidation
	mutex       sync.RWMutex
	httpClient  *http.Client
	stats       *interfaces.DocumentationLibraryStats
}

// NewDocumentationLibraryService creates a new documentation library service
func NewDocumentationLibraryService() *DocumentationLibraryService {
	service := &DocumentationLibraryService{
		links:       make(map[string]*interfaces.DocumentationLink),
		validations: make(map[string]*interfaces.LinkValidation),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		stats: &interfaces.DocumentationLibraryStats{
			LinksByProvider: make(map[string]int),
			LinksByType:     make(map[string]int),
			LinksByCategory: make(map[string]int),
		},
	}

	// Initialize with default documentation links
	service.initializeDefaultLinks()
	
	return service
}

// GetDocumentationLinks retrieves documentation links for a specific provider and topic
func (d *DocumentationLibraryService) GetDocumentationLinks(ctx context.Context, provider, topic string) ([]*interfaces.DocumentationLink, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var results []*interfaces.DocumentationLink
	
	for _, link := range d.links {
		if (provider == "" || strings.EqualFold(link.Provider, provider)) &&
		   (topic == "" || d.matchesTopic(link, topic)) {
			results = append(results, link)
		}
	}

	// Sort by relevance and validation status
	sort.Slice(results, func(i, j int) bool {
		// Prioritize valid links
		if results[i].IsValid != results[j].IsValid {
			return results[i].IsValid
		}
		// Then by last validated (more recent first)
		return results[i].LastValidated.After(results[j].LastValidated)
	})

	return results, nil
}

// ValidateLinks validates a list of documentation links
func (d *DocumentationLibraryService) ValidateLinks(ctx context.Context, links []*interfaces.DocumentationLink) ([]*interfaces.LinkValidation, error) {
	var validations []*interfaces.LinkValidation
	var wg sync.WaitGroup
	validationChan := make(chan *interfaces.LinkValidation, len(links))

	// Validate links concurrently
	for _, link := range links {
		wg.Add(1)
		go func(l *interfaces.DocumentationLink) {
			defer wg.Done()
			validation := d.validateSingleLink(ctx, l)
			validationChan <- validation
		}(link)
	}

	// Wait for all validations to complete
	go func() {
		wg.Wait()
		close(validationChan)
	}()

	// Collect results
	for validation := range validationChan {
		validations = append(validations, validation)
		
		// Update the link's validation status
		d.mutex.Lock()
		if link, exists := d.links[validation.LinkID]; exists {
			link.IsValid = validation.IsValid
			link.LastValidated = validation.ValidatedAt
		}
		d.validations[validation.LinkID] = validation
		d.mutex.Unlock()
	}

	d.updateStats()
	return validations, nil
}

// UpdateDocumentationIndex updates the documentation index
func (d *DocumentationLibraryService) UpdateDocumentationIndex(ctx context.Context) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Get all links for validation
	var allLinks []*interfaces.DocumentationLink
	for _, link := range d.links {
		allLinks = append(allLinks, link)
	}

	// Validate all links
	_, err := d.ValidateLinks(ctx, allLinks)
	if err != nil {
		return fmt.Errorf("failed to validate links during index update: %w", err)
	}

	d.stats.LastValidationRun = time.Now()
	return nil
}

// SearchDocumentation searches for documentation links based on a query
func (d *DocumentationLibraryService) SearchDocumentation(ctx context.Context, query string, providers []string) ([]*interfaces.DocumentationLink, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var results []*interfaces.DocumentationLink
	queryLower := strings.ToLower(query)

	for _, link := range d.links {
		// Filter by providers if specified
		if len(providers) > 0 && !d.containsProvider(providers, link.Provider) {
			continue
		}

		// Check if query matches title, description, or tags
		if d.matchesQuery(link, queryLower) {
			results = append(results, link)
		}
	}

	// Sort by relevance
	sort.Slice(results, func(i, j int) bool {
		scoreI := d.calculateRelevanceScore(results[i], queryLower)
		scoreJ := d.calculateRelevanceScore(results[j], queryLower)
		return scoreI > scoreJ
	})

	return results, nil
}

// AddDocumentationLink adds a new documentation link
func (d *DocumentationLibraryService) AddDocumentationLink(ctx context.Context, link *interfaces.DocumentationLink) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if link.ID == "" {
		link.ID = uuid.New().String()
	}

	// Validate the link before adding
	validation := d.validateSingleLink(ctx, link)
	link.IsValid = validation.IsValid
	link.LastValidated = validation.ValidatedAt

	d.links[link.ID] = link
	d.validations[link.ID] = validation
	
	d.updateStats()
	return nil
}

// RemoveDocumentationLink removes a documentation link
func (d *DocumentationLibraryService) RemoveDocumentationLink(ctx context.Context, linkID string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.links, linkID)
	delete(d.validations, linkID)
	
	d.updateStats()
	return nil
}

// GetLinksByCategory retrieves links by category
func (d *DocumentationLibraryService) GetLinksByCategory(ctx context.Context, category string) ([]*interfaces.DocumentationLink, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var results []*interfaces.DocumentationLink
	for _, link := range d.links {
		if strings.EqualFold(link.Category, category) {
			results = append(results, link)
		}
	}

	return results, nil
}

// GetLinksByProvider retrieves links by provider
func (d *DocumentationLibraryService) GetLinksByProvider(ctx context.Context, provider string) ([]*interfaces.DocumentationLink, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var results []*interfaces.DocumentationLink
	for _, link := range d.links {
		if strings.EqualFold(link.Provider, provider) {
			results = append(results, link)
		}
	}

	return results, nil
}

// GetLinksByType retrieves links by type
func (d *DocumentationLibraryService) GetLinksByType(ctx context.Context, linkType string) ([]*interfaces.DocumentationLink, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var results []*interfaces.DocumentationLink
	for _, link := range d.links {
		if strings.EqualFold(link.Type, linkType) {
			results = append(results, link)
		}
	}

	return results, nil
}

// GetLinkValidationStatus retrieves the validation status of a specific link
func (d *DocumentationLibraryService) GetLinkValidationStatus(ctx context.Context, linkID string) (*interfaces.LinkValidation, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	validation, exists := d.validations[linkID]
	if !exists {
		return nil, fmt.Errorf("validation not found for link ID: %s", linkID)
	}

	return validation, nil
}

// IsHealthy checks if the documentation library is healthy
func (d *DocumentationLibraryService) IsHealthy() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Consider healthy if we have links and most are valid
	totalLinks := len(d.links)
	if totalLinks == 0 {
		return false
	}

	validLinks := 0
	for _, link := range d.links {
		if link.IsValid {
			validLinks++
		}
	}

	// Consider healthy if at least 70% of links are valid
	healthThreshold := 0.7
	return float64(validLinks)/float64(totalLinks) >= healthThreshold
}

// GetStats returns statistics about the documentation library
func (d *DocumentationLibraryService) GetStats() *interfaces.DocumentationLibraryStats {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Return a copy of the stats
	statsCopy := *d.stats
	return &statsCopy
}

// Helper methods

func (d *DocumentationLibraryService) validateSingleLink(ctx context.Context, link *interfaces.DocumentationLink) *interfaces.LinkValidation {
	validation := &interfaces.LinkValidation{
		LinkID:      link.ID,
		URL:         link.URL,
		ValidatedAt: time.Now(),
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "HEAD", link.URL, nil)
	if err != nil {
		validation.IsValid = false
		validation.Error = fmt.Sprintf("failed to create request: %v", err)
		return validation
	}

	resp, err := d.httpClient.Do(req)
	validation.ResponseTime = time.Since(start)
	
	if err != nil {
		validation.IsValid = false
		validation.Error = fmt.Sprintf("request failed: %v", err)
		return validation
	}
	defer resp.Body.Close()

	validation.StatusCode = resp.StatusCode
	validation.ContentType = resp.Header.Get("Content-Type")
	
	if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		if parsed, err := time.Parse(time.RFC1123, lastModified); err == nil {
			validation.LastModified = parsed
		}
	}

	// Consider 2xx and 3xx status codes as valid
	validation.IsValid = resp.StatusCode >= 200 && resp.StatusCode < 400
	if !validation.IsValid {
		validation.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return validation
}

func (d *DocumentationLibraryService) matchesTopic(link *interfaces.DocumentationLink, topic string) bool {
	topicLower := strings.ToLower(topic)
	
	// Check title, description, and tags
	if strings.Contains(strings.ToLower(link.Title), topicLower) ||
	   strings.Contains(strings.ToLower(link.Description), topicLower) ||
	   strings.Contains(strings.ToLower(link.Topic), topicLower) {
		return true
	}

	// Check tags
	for _, tag := range link.Tags {
		if strings.Contains(strings.ToLower(tag), topicLower) {
			return true
		}
	}

	return false
}

func (d *DocumentationLibraryService) containsProvider(providers []string, provider string) bool {
	for _, p := range providers {
		if strings.EqualFold(p, provider) {
			return true
		}
	}
	return false
}

func (d *DocumentationLibraryService) matchesQuery(link *interfaces.DocumentationLink, query string) bool {
	// Check title, description, topic, and tags
	searchableText := strings.ToLower(fmt.Sprintf("%s %s %s %s", 
		link.Title, link.Description, link.Topic, strings.Join(link.Tags, " ")))
	
	return strings.Contains(searchableText, query)
}

func (d *DocumentationLibraryService) calculateRelevanceScore(link *interfaces.DocumentationLink, query string) float64 {
	score := 0.0
	
	// Title match gets highest score
	if strings.Contains(strings.ToLower(link.Title), query) {
		score += 10.0
	}
	
	// Topic match gets high score
	if strings.Contains(strings.ToLower(link.Topic), query) {
		score += 8.0
	}
	
	// Description match gets medium score
	if strings.Contains(strings.ToLower(link.Description), query) {
		score += 5.0
	}
	
	// Tag match gets lower score
	for _, tag := range link.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			score += 2.0
		}
	}
	
	// Boost score for valid links
	if link.IsValid {
		score *= 1.5
	}
	
	// Boost score for recently validated links
	daysSinceValidation := time.Since(link.LastValidated).Hours() / 24
	if daysSinceValidation < 7 {
		score *= 1.2
	}
	
	return score
}

func (d *DocumentationLibraryService) updateStats() {
	d.stats.TotalLinks = len(d.links)
	d.stats.ValidLinks = 0
	d.stats.InvalidLinks = 0
	d.stats.UnvalidatedLinks = 0
	
	// Reset counters
	d.stats.LinksByProvider = make(map[string]int)
	d.stats.LinksByType = make(map[string]int)
	d.stats.LinksByCategory = make(map[string]int)
	
	var totalResponseTime time.Duration
	validationCount := 0
	
	for _, link := range d.links {
		// Count by validation status
		if link.LastValidated.IsZero() {
			d.stats.UnvalidatedLinks++
		} else if link.IsValid {
			d.stats.ValidLinks++
		} else {
			d.stats.InvalidLinks++
		}
		
		// Count by provider
		d.stats.LinksByProvider[link.Provider]++
		
		// Count by type
		d.stats.LinksByType[link.Type]++
		
		// Count by category
		d.stats.LinksByCategory[link.Category]++
		
		// Calculate average response time
		if validation, exists := d.validations[link.ID]; exists {
			totalResponseTime += validation.ResponseTime
			validationCount++
		}
	}
	
	if validationCount > 0 {
		d.stats.AverageResponseTime = totalResponseTime / time.Duration(validationCount)
	}
}

// initializeDefaultLinks populates the library with default cloud provider documentation links
func (d *DocumentationLibraryService) initializeDefaultLinks() {
	defaultLinks := []*interfaces.DocumentationLink{
		// AWS Documentation
		{
			ID:          "aws-well-architected",
			Provider:    "aws",
			Topic:       "architecture",
			Title:       "AWS Well-Architected Framework",
			URL:         "https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html",
			Description: "The AWS Well-Architected Framework helps you understand the pros and cons of decisions you make while building systems on AWS.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategoryArchitecture),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"architecture", "best-practices", "framework", "design"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour), // Assume validated yesterday
		},
		{
			ID:          "aws-security-best-practices",
			Provider:    "aws",
			Topic:       "security",
			Title:       "AWS Security Best Practices",
			URL:         "https://docs.aws.amazon.com/security/",
			Description: "Security best practices and guidelines for AWS services and infrastructure.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategorySecurityDocs),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"security", "best-practices", "compliance", "iam"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "aws-pricing-calculator",
			Provider:    "aws",
			Topic:       "pricing",
			Title:       "AWS Pricing Calculator",
			URL:         "https://calculator.aws/",
			Description: "Estimate the cost for your architecture solution with AWS Pricing Calculator.",
			Type:        string(interfaces.LinkTypePricing),
			Category:    string(interfaces.CategoryPricing),
			Audience:    string(interfaces.AudienceMixed),
			Tags:        []string{"pricing", "calculator", "cost", "estimation"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "aws-migration-hub",
			Provider:    "aws",
			Topic:       "migration",
			Title:       "AWS Migration Hub Documentation",
			URL:         "https://docs.aws.amazon.com/migrationhub/",
			Description: "AWS Migration Hub provides a single location to track the progress of application migrations across multiple AWS and partner solutions.",
			Type:        string(interfaces.LinkTypeGuide),
			Category:    string(interfaces.CategoryMigration),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"migration", "hub", "tracking", "applications"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "aws-compliance-center",
			Provider:    "aws",
			Topic:       "compliance",
			Title:       "AWS Compliance Center",
			URL:         "https://aws.amazon.com/compliance/",
			Description: "AWS compliance programs and certifications to help you meet regulatory requirements.",
			Type:        string(interfaces.LinkTypeCompliance),
			Category:    string(interfaces.CategoryCompliance),
			Audience:    string(interfaces.AudienceBusiness),
			Tags:        []string{"compliance", "certifications", "regulatory", "governance"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},

		// Azure Documentation
		{
			ID:          "azure-architecture-center",
			Provider:    "azure",
			Topic:       "architecture",
			Title:       "Azure Architecture Center",
			URL:         "https://docs.microsoft.com/en-us/azure/architecture/",
			Description: "Architecture guidance for Azure including reference architectures, design patterns, and best practices.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategoryArchitecture),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"architecture", "patterns", "reference", "design"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "azure-security-center",
			Provider:    "azure",
			Topic:       "security",
			Title:       "Azure Security Documentation",
			URL:         "https://docs.microsoft.com/en-us/azure/security/",
			Description: "Security best practices, guidelines, and services for Azure cloud platform.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategorySecurityDocs),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"security", "best-practices", "defender", "sentinel"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "azure-pricing-calculator",
			Provider:    "azure",
			Topic:       "pricing",
			Title:       "Azure Pricing Calculator",
			URL:         "https://azure.microsoft.com/en-us/pricing/calculator/",
			Description: "Estimate your expected monthly costs for using Azure services.",
			Type:        string(interfaces.LinkTypePricing),
			Category:    string(interfaces.CategoryPricing),
			Audience:    string(interfaces.AudienceMixed),
			Tags:        []string{"pricing", "calculator", "cost", "estimation"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "azure-migrate",
			Provider:    "azure",
			Topic:       "migration",
			Title:       "Azure Migrate Documentation",
			URL:         "https://docs.microsoft.com/en-us/azure/migrate/",
			Description: "Azure Migrate provides a centralized hub to assess and migrate to Azure.",
			Type:        string(interfaces.LinkTypeGuide),
			Category:    string(interfaces.CategoryMigration),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"migration", "assessment", "discovery", "modernization"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "azure-compliance",
			Provider:    "azure",
			Topic:       "compliance",
			Title:       "Azure Compliance Documentation",
			URL:         "https://docs.microsoft.com/en-us/azure/compliance/",
			Description: "Azure compliance offerings and regulatory compliance information.",
			Type:        string(interfaces.LinkTypeCompliance),
			Category:    string(interfaces.CategoryCompliance),
			Audience:    string(interfaces.AudienceBusiness),
			Tags:        []string{"compliance", "regulatory", "certifications", "governance"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},

		// Google Cloud Documentation
		{
			ID:          "gcp-architecture-center",
			Provider:    "gcp",
			Topic:       "architecture",
			Title:       "Google Cloud Architecture Center",
			URL:         "https://cloud.google.com/architecture",
			Description: "Reference architectures, diagrams, design patterns, guidance, and best practices for Google Cloud.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategoryArchitecture),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"architecture", "patterns", "reference", "best-practices"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "gcp-security-best-practices",
			Provider:    "gcp",
			Topic:       "security",
			Title:       "Google Cloud Security Best Practices",
			URL:         "https://cloud.google.com/security/best-practices",
			Description: "Security best practices and recommendations for Google Cloud Platform.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategorySecurityDocs),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"security", "best-practices", "iam", "encryption"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "gcp-pricing-calculator",
			Provider:    "gcp",
			Topic:       "pricing",
			Title:       "Google Cloud Pricing Calculator",
			URL:         "https://cloud.google.com/products/calculator",
			Description: "Estimate your Google Cloud costs with the pricing calculator.",
			Type:        string(interfaces.LinkTypePricing),
			Category:    string(interfaces.CategoryPricing),
			Audience:    string(interfaces.AudienceMixed),
			Tags:        []string{"pricing", "calculator", "cost", "estimation"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "gcp-migrate-for-compute-engine",
			Provider:    "gcp",
			Topic:       "migration",
			Title:       "Migrate for Compute Engine Documentation",
			URL:         "https://cloud.google.com/migrate/compute-engine/docs",
			Description: "Migrate for Compute Engine enables you to migrate VMs from on-premises or other clouds to Google Cloud.",
			Type:        string(interfaces.LinkTypeGuide),
			Category:    string(interfaces.CategoryMigration),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"migration", "compute-engine", "vm", "lift-and-shift"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "gcp-compliance",
			Provider:    "gcp",
			Topic:       "compliance",
			Title:       "Google Cloud Compliance Resource Center",
			URL:         "https://cloud.google.com/security/compliance",
			Description: "Google Cloud compliance certifications, attestations, and regulatory information.",
			Type:        string(interfaces.LinkTypeCompliance),
			Category:    string(interfaces.CategoryCompliance),
			Audience:    string(interfaces.AudienceBusiness),
			Tags:        []string{"compliance", "certifications", "regulatory", "attestations"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},

		// Multi-Cloud and General Documentation
		{
			ID:          "nist-cloud-computing",
			Provider:    "nist",
			Topic:       "standards",
			Title:       "NIST Cloud Computing Standards",
			URL:         "https://www.nist.gov/programs-projects/nist-cloud-computing-program-nccp",
			Description: "NIST cloud computing standards and guidelines for secure cloud adoption.",
			Type:        string(interfaces.LinkTypeCompliance),
			Category:    string(interfaces.CategoryCompliance),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"nist", "standards", "security", "compliance"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "cis-controls",
			Provider:    "cis",
			Topic:       "security",
			Title:       "CIS Controls for Cloud Security",
			URL:         "https://www.cisecurity.org/controls/",
			Description: "Center for Internet Security (CIS) Controls for effective cyber defense in cloud environments.",
			Type:        string(interfaces.LinkTypeBestPractice),
			Category:    string(interfaces.CategorySecurityDocs),
			Audience:    string(interfaces.AudienceTechnical),
			Tags:        []string{"cis", "controls", "security", "cyber-defense"},
			IsValid:     true,
			LastValidated: time.Now().Add(-24 * time.Hour),
		},
	}

	// Add all default links
	for _, link := range defaultLinks {
		d.links[link.ID] = link
	}

	// Update stats after initialization
	d.updateStats()
}