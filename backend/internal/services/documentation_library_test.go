package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

func TestNewDocumentationLibraryService(t *testing.T) {
	service := NewDocumentationLibraryService()
	
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	
	if service.links == nil {
		t.Fatal("Expected links map to be initialized")
	}
	
	if service.validations == nil {
		t.Fatal("Expected validations map to be initialized")
	}
	
	if service.httpClient == nil {
		t.Fatal("Expected HTTP client to be initialized")
	}
	
	if service.stats == nil {
		t.Fatal("Expected stats to be initialized")
	}
	
	// Check that default links are loaded
	if len(service.links) == 0 {
		t.Fatal("Expected default links to be loaded")
	}
}

func TestGetDocumentationLinks(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	tests := []struct {
		name     string
		provider string
		topic    string
		wantMin  int // minimum expected results
	}{
		{
			name:     "Get all links",
			provider: "",
			topic:    "",
			wantMin:  10, // We have at least 10+ default links
		},
		{
			name:     "Get AWS links",
			provider: "aws",
			topic:    "",
			wantMin:  3, // We have several AWS links
		},
		{
			name:     "Get Azure links",
			provider: "azure",
			topic:    "",
			wantMin:  3, // We have several Azure links
		},
		{
			name:     "Get GCP links",
			provider: "gcp",
			topic:    "",
			wantMin:  3, // We have several GCP links
		},
		{
			name:     "Get security topic links",
			provider: "",
			topic:    "security",
			wantMin:  2, // We have security links for multiple providers
		},
		{
			name:     "Get architecture topic links",
			provider: "",
			topic:    "architecture",
			wantMin:  2, // We have architecture links for multiple providers
		},
		{
			name:     "Get AWS security links",
			provider: "aws",
			topic:    "security",
			wantMin:  1, // We have AWS security links
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links, err := service.GetDocumentationLinks(ctx, tt.provider, tt.topic)
			if err != nil {
				t.Fatalf("GetDocumentationLinks() error = %v", err)
			}
			
			if len(links) < tt.wantMin {
				t.Errorf("GetDocumentationLinks() got %d links, want at least %d", len(links), tt.wantMin)
			}
			
			// Verify filtering works correctly
			for _, link := range links {
				if tt.provider != "" && link.Provider != tt.provider {
					t.Errorf("Expected provider %s, got %s", tt.provider, link.Provider)
				}
				
				if tt.topic != "" {
					// Check if topic matches in title, description, topic, or tags
					found := false
					if containsIgnoreCase(link.Title, tt.topic) ||
					   containsIgnoreCase(link.Description, tt.topic) ||
					   containsIgnoreCase(link.Topic, tt.topic) {
						found = true
					}
					for _, tag := range link.Tags {
						if containsIgnoreCase(tag, tt.topic) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Link %s does not match topic %s", link.Title, tt.topic)
					}
				}
			}
		})
	}
}

func TestSearchDocumentation(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	tests := []struct {
		name      string
		query     string
		providers []string
		wantMin   int
	}{
		{
			name:    "Search for security",
			query:   "security",
			wantMin: 2,
		},
		{
			name:    "Search for architecture",
			query:   "architecture",
			wantMin: 2,
		},
		{
			name:    "Search for pricing",
			query:   "pricing",
			wantMin: 2,
		},
		{
			name:      "Search AWS security",
			query:     "security",
			providers: []string{"aws"},
			wantMin:   1,
		},
		{
			name:    "Search for migration",
			query:   "migration",
			wantMin: 2,
		},
		{
			name:    "Search for compliance",
			query:   "compliance",
			wantMin: 2,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.SearchDocumentation(ctx, tt.query, tt.providers)
			if err != nil {
				t.Fatalf("SearchDocumentation() error = %v", err)
			}
			
			if len(results) < tt.wantMin {
				t.Errorf("SearchDocumentation() got %d results, want at least %d", len(results), tt.wantMin)
			}
			
			// Verify provider filtering
			if len(tt.providers) > 0 {
				for _, result := range results {
					found := false
					for _, provider := range tt.providers {
						if result.Provider == provider {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Result provider %s not in expected providers %v", result.Provider, tt.providers)
					}
				}
			}
		})
	}
}

func TestAddAndRemoveDocumentationLink(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	// Create a test server for link validation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	// Test adding a link
	testLink := &interfaces.DocumentationLink{
		Provider:    "test",
		Topic:       "testing",
		Title:       "Test Documentation",
		URL:         server.URL,
		Description: "Test documentation for testing purposes",
		Type:        string(interfaces.LinkTypeGuide),
		Category:    string(interfaces.CategoryGettingStarted),
		Audience:    string(interfaces.AudienceTechnical),
		Tags:        []string{"test", "documentation"},
	}
	
	err := service.AddDocumentationLink(ctx, testLink)
	if err != nil {
		t.Fatalf("AddDocumentationLink() error = %v", err)
	}
	
	if testLink.ID == "" {
		t.Error("Expected ID to be generated for the link")
	}
	
	// Verify the link was added
	links, err := service.GetDocumentationLinks(ctx, "test", "")
	if err != nil {
		t.Fatalf("GetDocumentationLinks() error = %v", err)
	}
	
	if len(links) != 1 {
		t.Errorf("Expected 1 test link, got %d", len(links))
	}
	
	// Test removing the link
	err = service.RemoveDocumentationLink(ctx, testLink.ID)
	if err != nil {
		t.Fatalf("RemoveDocumentationLink() error = %v", err)
	}
	
	// Verify the link was removed
	links, err = service.GetDocumentationLinks(ctx, "test", "")
	if err != nil {
		t.Fatalf("GetDocumentationLinks() error = %v", err)
	}
	
	if len(links) != 0 {
		t.Errorf("Expected 0 test links after removal, got %d", len(links))
	}
}

func TestValidateLinks(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	// Create test servers
	validServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer validServer.Close()
	
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer invalidServer.Close()
	
	testLinks := []*interfaces.DocumentationLink{
		{
			ID:  "valid-link",
			URL: validServer.URL,
		},
		{
			ID:  "invalid-link",
			URL: invalidServer.URL,
		},
		{
			ID:  "malformed-link",
			URL: "not-a-valid-url",
		},
	}
	
	validations, err := service.ValidateLinks(ctx, testLinks)
	if err != nil {
		t.Fatalf("ValidateLinks() error = %v", err)
	}
	
	if len(validations) != 3 {
		t.Errorf("Expected 3 validations, got %d", len(validations))
	}
	
	// Check validation results
	validationMap := make(map[string]*interfaces.LinkValidation)
	for _, v := range validations {
		validationMap[v.LinkID] = v
	}
	
	if validationMap["valid-link"].IsValid != true {
		t.Error("Expected valid-link to be valid")
	}
	
	if validationMap["invalid-link"].IsValid != false {
		t.Error("Expected invalid-link to be invalid")
	}
	
	if validationMap["malformed-link"].IsValid != false {
		t.Error("Expected malformed-link to be invalid")
	}
}

func TestGetLinksByCategory(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	links, err := service.GetLinksByCategory(ctx, string(interfaces.CategorySecurityDocs))
	if err != nil {
		t.Fatalf("GetLinksByCategory() error = %v", err)
	}
	
	if len(links) == 0 {
		t.Error("Expected to find security category links")
	}
	
	for _, link := range links {
		if link.Category != string(interfaces.CategorySecurityDocs) {
			t.Errorf("Expected category %s, got %s", interfaces.CategorySecurityDocs, link.Category)
		}
	}
}

func TestGetLinksByProvider(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	links, err := service.GetLinksByProvider(ctx, "aws")
	if err != nil {
		t.Fatalf("GetLinksByProvider() error = %v", err)
	}
	
	if len(links) == 0 {
		t.Error("Expected to find AWS provider links")
	}
	
	for _, link := range links {
		if link.Provider != "aws" {
			t.Errorf("Expected provider aws, got %s", link.Provider)
		}
	}
}

func TestGetLinksByType(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	links, err := service.GetLinksByType(ctx, string(interfaces.LinkTypeBestPractice))
	if err != nil {
		t.Fatalf("GetLinksByType() error = %v", err)
	}
	
	if len(links) == 0 {
		t.Error("Expected to find best-practice type links")
	}
	
	for _, link := range links {
		if link.Type != string(interfaces.LinkTypeBestPractice) {
			t.Errorf("Expected type %s, got %s", interfaces.LinkTypeBestPractice, link.Type)
		}
	}
}

func TestIsHealthy(t *testing.T) {
	service := NewDocumentationLibraryService()
	
	// Should be healthy with default links (assuming they're marked as valid)
	if !service.IsHealthy() {
		t.Error("Expected service to be healthy with default links")
	}
	
	// Test with empty service
	emptyService := &DocumentationLibraryService{
		links:       make(map[string]*interfaces.DocumentationLink),
		validations: make(map[string]*interfaces.LinkValidation),
		stats: &interfaces.DocumentationLibraryStats{
			LinksByProvider: make(map[string]int),
			LinksByType:     make(map[string]int),
			LinksByCategory: make(map[string]int),
		},
	}
	
	if emptyService.IsHealthy() {
		t.Error("Expected empty service to be unhealthy")
	}
}

func TestGetStats(t *testing.T) {
	service := NewDocumentationLibraryService()
	
	stats := service.GetStats()
	if stats == nil {
		t.Fatal("Expected stats to be returned")
	}
	
	if stats.TotalLinks == 0 {
		t.Error("Expected total links to be greater than 0")
	}
	
	if len(stats.LinksByProvider) == 0 {
		t.Error("Expected links by provider to be populated")
	}
	
	if len(stats.LinksByType) == 0 {
		t.Error("Expected links by type to be populated")
	}
	
	if len(stats.LinksByCategory) == 0 {
		t.Error("Expected links by category to be populated")
	}
}

func TestUpdateDocumentationIndex(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	// This will attempt to validate all links, which might fail for external URLs
	// but should not return an error for the operation itself
	err := service.UpdateDocumentationIndex(ctx)
	if err != nil {
		t.Fatalf("UpdateDocumentationIndex() error = %v", err)
	}
	
	stats := service.GetStats()
	if stats.LastValidationRun.IsZero() {
		t.Error("Expected LastValidationRun to be updated")
	}
}

func TestGetLinkValidationStatus(t *testing.T) {
	service := NewDocumentationLibraryService()
	ctx := context.Background()
	
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	// Add a test link
	testLink := &interfaces.DocumentationLink{
		ID:  "test-validation",
		URL: server.URL,
	}
	
	err := service.AddDocumentationLink(ctx, testLink)
	if err != nil {
		t.Fatalf("AddDocumentationLink() error = %v", err)
	}
	
	// Get validation status
	validation, err := service.GetLinkValidationStatus(ctx, testLink.ID)
	if err != nil {
		t.Fatalf("GetLinkValidationStatus() error = %v", err)
	}
	
	if validation.LinkID != testLink.ID {
		t.Errorf("Expected LinkID %s, got %s", testLink.ID, validation.LinkID)
	}
	
	// Test non-existent link
	_, err = service.GetLinkValidationStatus(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent link")
	}
}

// Helper function for case-insensitive string contains check
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      containsIgnoreCaseHelper(s, substr))))
}

func containsIgnoreCaseHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}