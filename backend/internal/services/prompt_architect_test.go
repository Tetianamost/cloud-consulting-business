package services

import (
	"context"
	"strings"
	"testing"

	"github.com/cloud-consulting/backend/internal/domain"
)

func TestNewPromptArchitect(t *testing.T) {
	pa := NewPromptArchitect()
	
	if pa == nil {
		t.Fatal("NewPromptArchitect returned nil")
	}
	
	// Check that default templates are loaded
	templates := pa.ListTemplates()
	expectedTemplates := []string{
		"enhanced_report",
		"technical_report", 
		"business_report",
		"interview_guide",
		"risk_assessment",
		"competitive_analysis",
		"comprehensive_report",
	}
	
	if len(templates) != len(expectedTemplates) {
		t.Errorf("Expected %d templates, got %d", len(expectedTemplates), len(templates))
	}
	
	for _, expected := range expectedTemplates {
		found := false
		for _, actual := range templates {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template %s not found", expected)
		}
	}
}

func TestBuildReportPrompt(t *testing.T) {
	pa := NewPromptArchitect()
	
	inquiry := &domain.Inquiry{
		ID:       "test-123",
		Name:     "John Doe",
		Email:    "john@example.com",
		Company:  "Test Corp",
		Phone:    "555-1234",
		Services: []string{"assessment", "migration"},
		Message:  "We need help migrating our infrastructure to the cloud",
		Priority: "high",
	}
	
	options := &PromptOptions{
		TargetAudience:             "mixed",
		IncludeDocumentationLinks:  true,
		IncludeCompetitiveAnalysis: true,
		IncludeRiskAssessment:      true,
		CloudProviders:             []string{"AWS", "Azure"},
		MaxTokens:                  2000,
	}
	
	prompt, err := pa.BuildReportPrompt(context.Background(), inquiry, options)
	if err != nil {
		t.Fatalf("BuildReportPrompt failed: %v", err)
	}
	
	// Verify prompt contains expected content
	expectedContent := []string{
		"John Doe",
		"john@example.com", 
		"Test Corp",
		"555-1234",
		"assessment, migration",
		"We need help migrating our infrastructure to the cloud",
		"EXECUTIVE SUMMARY",
		"RECOMMENDATIONS",
		"NEXT STEPS",
	}
	
	for _, content := range expectedContent {
		if !strings.Contains(prompt, content) {
			t.Errorf("Prompt missing expected content: %s", content)
		}
	}
	
	// Verify multi-cloud content is included
	if !strings.Contains(prompt, "AWS, Azure") {
		t.Error("Prompt should include specified cloud providers")
	}
	
	// Verify enhanced features are included
	if !strings.Contains(prompt, "DOCUMENTATION REQUIREMENTS") {
		t.Error("Prompt should include documentation requirements when requested")
	}
}

func TestBuildReportPromptWithDefaults(t *testing.T) {
	pa := NewPromptArchitect()
	
	inquiry := &domain.Inquiry{
		ID:       "test-123",
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Services: []string{"optimization"},
		Message:  "Looking to optimize our cloud costs",
	}
	
	// Test with nil options (should use defaults)
	prompt, err := pa.BuildReportPrompt(context.Background(), inquiry, nil)
	if err != nil {
		t.Fatalf("BuildReportPrompt with nil options failed: %v", err)
	}
	
	// Should use default values
	if !strings.Contains(prompt, "Jane Smith") {
		t.Error("Prompt should contain client name")
	}
	
	if !strings.Contains(prompt, "Client Organization") {
		t.Error("Prompt should use default company name")
	}
	
	if !strings.Contains(prompt, "Not provided") {
		t.Error("Prompt should use default phone")
	}
}

func TestBuildInterviewPrompt(t *testing.T) {
	pa := NewPromptArchitect()
	
	inquiry := &domain.Inquiry{
		Name:     "Alice Johnson",
		Company:  "Tech Startup",
		Services: []string{"assessment"},
		Message:  "Need to assess our current cloud setup",
	}
	
	prompt, err := pa.BuildInterviewPrompt(context.Background(), inquiry)
	if err != nil {
		t.Fatalf("BuildInterviewPrompt failed: %v", err)
	}
	
	expectedContent := []string{
		"Alice Johnson",
		"Tech Startup",
		"assessment",
		"INTERVIEW OBJECTIVES",
		"BUSINESS CONTEXT QUESTIONS",
		"TECHNICAL DISCOVERY QUESTIONS",
	}
	
	for _, content := range expectedContent {
		if !strings.Contains(prompt, content) {
			t.Errorf("Interview prompt missing expected content: %s", content)
		}
	}
}

func TestBuildRiskAssessmentPrompt(t *testing.T) {
	pa := NewPromptArchitect()
	
	inquiry := &domain.Inquiry{
		Name:     "Bob Wilson",
		Company:  "Healthcare Corp",
		Services: []string{"migration"},
		Message:  "Migrating patient data to cloud with HIPAA compliance",
	}
	
	prompt, err := pa.BuildRiskAssessmentPrompt(context.Background(), inquiry)
	if err != nil {
		t.Fatalf("BuildRiskAssessmentPrompt failed: %v", err)
	}
	
	expectedContent := []string{
		"Bob Wilson",
		"Healthcare Corp",
		"TECHNICAL RISKS",
		"SECURITY RISKS",
		"OPERATIONAL RISKS",
		"BUSINESS RISKS",
		"COMPLIANCE RISKS",
		"Healthcare", // Should detect healthcare industry
	}
	
	for _, content := range expectedContent {
		if !strings.Contains(prompt, content) {
			t.Errorf("Risk assessment prompt missing expected content: %s", content)
		}
	}
}

func TestBuildCompetitiveAnalysisPrompt(t *testing.T) {
	pa := NewPromptArchitect()
	
	inquiry := &domain.Inquiry{
		Name:     "Carol Davis",
		Company:  "E-commerce Inc",
		Services: []string{"optimization"},
		Message:  "Need to optimize our e-commerce platform performance",
	}
	
	prompt, err := pa.BuildCompetitiveAnalysisPrompt(context.Background(), inquiry)
	if err != nil {
		t.Fatalf("BuildCompetitiveAnalysisPrompt failed: %v", err)
	}
	
	expectedContent := []string{
		"Carol Davis",
		"E-commerce Inc",
		"AWS ANALYSIS",
		"MICROSOFT AZURE ANALYSIS", 
		"GOOGLE CLOUD PLATFORM ANALYSIS",
		"SCENARIO-BASED RECOMMENDATIONS",
		"FINAL RECOMMENDATION",
	}
	
	for _, content := range expectedContent {
		if !strings.Contains(prompt, content) {
			t.Errorf("Competitive analysis prompt missing expected content: %s", content)
		}
	}
}

func TestValidatePrompt(t *testing.T) {
	pa := NewPromptArchitect()
	
	tests := []struct {
		name    string
		prompt  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid prompt",
			prompt:  strings.Repeat("This is a valid prompt with sufficient content. ", 10) + "EXECUTIVE SUMMARY and RECOMMENDATIONS and NEXT STEPS",
			wantErr: false,
		},
		{
			name:    "empty prompt",
			prompt:  "",
			wantErr: true,
			errMsg:  "prompt cannot be empty",
		},
		{
			name:    "too short",
			prompt:  "Short",
			wantErr: true,
			errMsg:  "prompt too short",
		},
		{
			name:    "too long",
			prompt:  strings.Repeat("x", 50001),
			wantErr: true,
			errMsg:  "prompt too long",
		},
		{
			name:    "missing required sections",
			prompt:  "Generate a professional consulting report " + strings.Repeat("This prompt is long enough but missing required sections. ", 10),
			wantErr: true,
			errMsg:  "prompt missing required section",
		},
		{
			name:    "unsubstituted variables",
			prompt:  strings.Repeat("This prompt has unsubstituted variables like {{.ClientName}}. ", 10) + "EXECUTIVE SUMMARY and RECOMMENDATIONS and NEXT STEPS",
			wantErr: true,
			errMsg:  "prompt contains unsubstituted variables",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pa.ValidatePrompt(tt.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidatePrompt() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	pa := NewPromptArchitect()
	
	// Test getting existing template
	template, err := pa.GetTemplate("enhanced_report")
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}
	
	if template.Name != "enhanced_report" {
		t.Errorf("Expected template name 'enhanced_report', got %s", template.Name)
	}
	
	if template.Category != "report" {
		t.Errorf("Expected template category 'report', got %s", template.Category)
	}
	
	// Test getting non-existent template
	_, err = pa.GetTemplate("non_existent")
	if err == nil {
		t.Error("Expected error for non-existent template")
	}
}

func TestRegisterTemplate(t *testing.T) {
	pa := NewPromptArchitect()
	
	template := &PromptTemplate{
		Name:        "test_template",
		Category:    "test",
		Description: "Test template",
		Template:    "Hello {{.ClientName}}, this is a test template with {{.Message}}",
		RequiredVariables: []string{"ClientName", "Message"},
		OptionalVariables: []string{},
		ValidationRules: []ValidationRule{
			{Name: "HasHello", Pattern: "Hello", ErrorMessage: "Must contain Hello", Required: true},
		},
	}
	
	err := pa.RegisterTemplate(template)
	if err != nil {
		t.Fatalf("RegisterTemplate failed: %v", err)
	}
	
	// Verify template was registered
	retrieved, err := pa.GetTemplate("test_template")
	if err != nil {
		t.Fatalf("Failed to retrieve registered template: %v", err)
	}
	
	if retrieved.Name != template.Name {
		t.Errorf("Expected name %s, got %s", template.Name, retrieved.Name)
	}
	
	// Test registering template with empty name
	invalidTemplate := &PromptTemplate{
		Name:     "",
		Template: "Invalid template",
	}
	
	err = pa.RegisterTemplate(invalidTemplate)
	if err == nil {
		t.Error("Expected error for template with empty name")
	}
	
	// Test registering template with empty content
	invalidTemplate2 := &PromptTemplate{
		Name:     "invalid2",
		Template: "",
	}
	
	err = pa.RegisterTemplate(invalidTemplate2)
	if err == nil {
		t.Error("Expected error for template with empty content")
	}
}

func TestExtractIndustryHints(t *testing.T) {
	pa := &promptArchitect{}
	
	tests := []struct {
		message  string
		company  string
		expected string
	}{
		{
			message:  "We need HIPAA compliant cloud storage for patient data",
			company:  "Regional Hospital",
			expected: "Healthcare",
		},
		{
			message:  "Banking application needs PCI compliance",
			company:  "First National Bank",
			expected: "Financial",
		},
		{
			message:  "E-commerce platform optimization",
			company:  "Online Retail Store",
			expected: "Retail",
		},
		{
			message:  "Manufacturing supply chain management",
			company:  "Industrial Corp",
			expected: "Manufacturing",
		},
		{
			message:  "Student information system migration",
			company:  "State University",
			expected: "Education",
		},
		{
			message:  "General cloud consulting",
			company:  "Generic Corp",
			expected: "General",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := pa.extractIndustryHints(tt.message, tt.company)
			if result != tt.expected {
				t.Errorf("extractIndustryHints(%q, %q) = %q, want %q", tt.message, tt.company, result, tt.expected)
			}
		})
	}
}

func TestExtractUseCase(t *testing.T) {
	pa := &promptArchitect{}
	
	tests := []struct {
		services []string
		message  string
		expected string
	}{
		{
			services: []string{"migration"},
			message:  "Need to migrate to cloud",
			expected: "Migration",
		},
		{
			services: []string{"assessment", "optimization"},
			message:  "Assess and optimize our setup",
			expected: "Assessment",
		},
		{
			services: []string{},
			message:  "We need to optimize our cloud costs",
			expected: "Optimization",
		},
		{
			services: []string{},
			message:  "Architecture review needed",
			expected: "Architecture",
		},
		{
			services: []string{},
			message:  "General consulting help",
			expected: "General Consulting",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := pa.extractUseCase(tt.services, tt.message)
			if result != tt.expected {
				t.Errorf("extractUseCase(%v, %q) = %q, want %q", tt.services, tt.message, result, tt.expected)
			}
		})
	}
}

func TestTemplateRendering(t *testing.T) {
	pa := NewPromptArchitect()
	
	// Test that templates can be rendered without errors
	inquiry := &domain.Inquiry{
		Name:     "Test Client",
		Email:    "test@example.com",
		Company:  "Test Company",
		Phone:    "555-0123",
		Services: []string{"assessment"},
		Message:  "Test message for template rendering",
		Priority: "medium",
	}
	
	templateNames := []string{
		"enhanced_report",
		"technical_report",
		"business_report",
		"interview_guide",
		"risk_assessment",
		"competitive_analysis",
		"comprehensive_report",
	}
	
	for _, templateName := range templateNames {
		t.Run(templateName, func(t *testing.T) {
			template, err := pa.GetTemplate(templateName)
			if err != nil {
				t.Fatalf("Failed to get template %s: %v", templateName, err)
			}
			
			var variables map[string]interface{}
			
			switch templateName {
			case "interview_guide":
				variables = pa.(*promptArchitect).prepareInterviewVariables(inquiry)
			case "risk_assessment":
				variables = pa.(*promptArchitect).prepareRiskAssessmentVariables(inquiry)
			case "competitive_analysis":
				variables = pa.(*promptArchitect).prepareCompetitiveAnalysisVariables(inquiry)
			default:
				options := &PromptOptions{
					TargetAudience: "mixed",
					CloudProviders: []string{"AWS", "Azure"},
				}
				variables = pa.(*promptArchitect).prepareReportVariables(inquiry, options)
			}
			
			rendered, err := pa.(*promptArchitect).renderTemplate(template, variables)
			if err != nil {
				t.Fatalf("Failed to render template %s: %v", templateName, err)
			}
			
			if len(rendered) == 0 {
				t.Errorf("Template %s rendered to empty string", templateName)
			}
			
			// Verify no unsubstituted variables remain
			if strings.Contains(rendered, "{{") || strings.Contains(rendered, "}}") {
				t.Errorf("Template %s contains unsubstituted variables", templateName)
			}
		})
	}
}