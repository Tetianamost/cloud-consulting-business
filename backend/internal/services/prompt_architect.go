package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)



// promptArchitect implements the PromptArchitect interface
type promptArchitect struct {
	templates        map[string]*interfaces.PromptTemplate
	audienceDetector AudienceDetector
}

// NewPromptArchitect creates a new PromptArchitect instance
func NewPromptArchitect() interfaces.PromptArchitect {
	pa := &promptArchitect{
		templates:        make(map[string]*interfaces.PromptTemplate),
		audienceDetector: NewAudienceDetector(),
	}
	
	// Initialize with default templates
	pa.initializeDefaultTemplates()
	
	return pa
}

// BuildReportPrompt creates a sophisticated prompt for report generation
func (pa *promptArchitect) BuildReportPrompt(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.PromptOptions) (string, error) {
	if options == nil {
		options = &interfaces.PromptOptions{
			TargetAudience: "mixed",
			MaxTokens: 2000,
		}
	}

	// Detect audience if not explicitly specified or if set to "mixed"
	var audienceProfile *AudienceProfile
	var err error
	
	if options.TargetAudience == "mixed" || options.TargetAudience == "" {
		audienceProfile, err = pa.audienceDetector.DetectAudience(ctx, inquiry)
		if err != nil {
			// Fall back to mixed audience if detection fails
			audienceProfile = &AudienceProfile{
				PrimaryType:    AudienceMixed,
				TechnicalDepth: 3,
				BusinessFocus:  3,
				Confidence:     0.5,
			}
		}
		// Update options with detected audience
		options.TargetAudience = string(audienceProfile.PrimaryType)
	} else {
		// Create a profile based on explicitly specified audience
		audienceProfile = pa.createProfileFromAudience(options.TargetAudience)
	}

	// Select appropriate template based on detected/specified audience and options
	templateName := pa.selectReportTemplate(inquiry, options)
	template, err := pa.GetTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to get template %s: %w", templateName, err)
	}

	// Prepare template variables with audience information
	variables := pa.prepareReportVariables(inquiry, options)
	
	// Add audience-specific variables
	variables["AudienceProfile"] = audienceProfile
	variables["TechnicalDepth"] = audienceProfile.TechnicalDepth
	variables["BusinessFocus"] = audienceProfile.BusinessFocus
	variables["AudienceConfidence"] = audienceProfile.Confidence
	variables["AudienceIndicators"] = strings.Join(audienceProfile.Indicators, ", ")
	
	// Render the template
	prompt, err := pa.renderTemplate(template, variables)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	// Apply audience-specific content adaptation
	adaptedPrompt, err := pa.audienceDetector.AdaptContentForAudience(prompt, audienceProfile)
	if err != nil {
		// Use original prompt if adaptation fails
		adaptedPrompt = prompt
	}

	// Validate the generated prompt
	if err := pa.ValidatePrompt(adaptedPrompt); err != nil {
		return "", fmt.Errorf("prompt validation failed: %w", err)
	}

	return adaptedPrompt, nil
}

// BuildInterviewPrompt creates a prompt for interview guide generation
func (pa *promptArchitect) BuildInterviewPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	template, err := pa.GetTemplate("interview_guide")
	if err != nil {
		return "", fmt.Errorf("failed to get interview template: %w", err)
	}

	variables := pa.prepareInterviewVariables(inquiry)
	
	prompt, err := pa.renderTemplate(template, variables)
	if err != nil {
		return "", fmt.Errorf("failed to render interview template: %w", err)
	}

	return prompt, nil
}

// BuildRiskAssessmentPrompt creates a prompt for risk assessment generation
func (pa *promptArchitect) BuildRiskAssessmentPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	template, err := pa.GetTemplate("risk_assessment")
	if err != nil {
		return "", fmt.Errorf("failed to get risk assessment template: %w", err)
	}

	variables := pa.prepareRiskAssessmentVariables(inquiry)
	
	prompt, err := pa.renderTemplate(template, variables)
	if err != nil {
		return "", fmt.Errorf("failed to render risk assessment template: %w", err)
	}

	return prompt, nil
}

// BuildCompetitiveAnalysisPrompt creates a prompt for competitive analysis generation
func (pa *promptArchitect) BuildCompetitiveAnalysisPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	template, err := pa.GetTemplate("competitive_analysis")
	if err != nil {
		return "", fmt.Errorf("failed to get competitive analysis template: %w", err)
	}

	variables := pa.prepareCompetitiveAnalysisVariables(inquiry)
	
	prompt, err := pa.renderTemplate(template, variables)
	if err != nil {
		return "", fmt.Errorf("failed to render competitive analysis template: %w", err)
	}

	return prompt, nil
}

// ValidatePrompt validates a generated prompt against various criteria
func (pa *promptArchitect) ValidatePrompt(prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	// Check minimum length
	if len(prompt) < 100 {
		return fmt.Errorf("prompt too short: minimum 100 characters required")
	}

	// Check maximum length (reasonable limit for AI models)
	if len(prompt) > 50000 {
		return fmt.Errorf("prompt too long: maximum 50000 characters allowed")
	}

	// Check for required sections in report prompts
	if strings.Contains(prompt, "Generate a professional consulting report") || 
	   strings.Contains(prompt, "consulting report") {
		requiredSections := []string{"EXECUTIVE SUMMARY", "RECOMMENDATIONS", "NEXT STEPS"}
		for _, section := range requiredSections {
			if !strings.Contains(prompt, section) {
				return fmt.Errorf("prompt missing required section: %s", section)
			}
		}
	}

	// Check for placeholder variables that weren't substituted
	placeholderPattern := regexp.MustCompile(`\{\{[^}]+\}\}`)
	if placeholderPattern.MatchString(prompt) {
		matches := placeholderPattern.FindAllString(prompt, -1)
		return fmt.Errorf("prompt contains unsubstituted variables: %v", matches)
	}

	return nil
}

// GetTemplate retrieves a template by name
func (pa *promptArchitect) GetTemplate(templateName string) (*interfaces.PromptTemplate, error) {
	template, exists := pa.templates[templateName]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}
	return template, nil
}

// RegisterTemplate registers a new template
func (pa *promptArchitect) RegisterTemplate(template *interfaces.PromptTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if template.Template == "" {
		return fmt.Errorf("template content cannot be empty")
	}

	// Validate template syntax
	if err := pa.validateTemplateSyntax(template); err != nil {
		return fmt.Errorf("template syntax validation failed: %w", err)
	}

	template.UpdatedAt = time.Now().Format(time.RFC3339)
	if template.CreatedAt == "" {
		template.CreatedAt = time.Now().Format(time.RFC3339)
	}

	pa.templates[template.Name] = template
	return nil
}

// ListTemplates returns a list of available template names
func (pa *promptArchitect) ListTemplates() []string {
	names := make([]string, 0, len(pa.templates))
	for name := range pa.templates {
		names = append(names, name)
	}
	return names
}

// selectReportTemplate selects the appropriate template based on inquiry and options
func (pa *promptArchitect) selectReportTemplate(inquiry *domain.Inquiry, options *interfaces.PromptOptions) string {
	// Enhanced template selection logic
	if options.IncludeCompetitiveAnalysis && options.IncludeRiskAssessment {
		return "comprehensive_report"
	}
	
	if options.TargetAudience == "technical" {
		return "technical_report"
	}
	
	if options.TargetAudience == "business" {
		return "business_report"
	}
	
	// Default to enhanced report template
	return "enhanced_report"
}

// prepareReportVariables prepares variables for report template rendering
func (pa *promptArchitect) prepareReportVariables(inquiry *domain.Inquiry, options *interfaces.PromptOptions) map[string]interface{} {
	variables := map[string]interface{}{
		"ClientName":     inquiry.Name,
		"ClientEmail":    inquiry.Email,
		"ClientCompany":  pa.getCompanyOrDefault(inquiry.Company),
		"ClientPhone":    pa.getPhoneOrDefault(inquiry.Phone),
		"Services":       strings.Join(inquiry.Services, ", "),
		"Message":        inquiry.Message,
		"Priority":       inquiry.Priority,
		"GeneratedDate":  time.Now().Format("January 2, 2006"),
		"TargetAudience": options.TargetAudience,
		"IndustryContext": options.IndustryContext,
	}

	// Add cloud provider specific variables
	if len(options.CloudProviders) > 0 {
		variables["CloudProviders"] = strings.Join(options.CloudProviders, ", ")
		variables["MultiCloud"] = len(options.CloudProviders) > 1
	} else {
		variables["CloudProviders"] = "AWS, Azure, GCP"
		variables["MultiCloud"] = true
	}

	// Add feature flags
	variables["IncludeDocumentationLinks"] = options.IncludeDocumentationLinks
	variables["IncludeCompetitiveAnalysis"] = options.IncludeCompetitiveAnalysis
	variables["IncludeRiskAssessment"] = options.IncludeRiskAssessment
	variables["IncludeImplementationSteps"] = options.IncludeImplementationSteps

	return variables
}

// prepareInterviewVariables prepares variables for interview template rendering
func (pa *promptArchitect) prepareInterviewVariables(inquiry *domain.Inquiry) map[string]interface{} {
	return map[string]interface{}{
		"ClientName":    inquiry.Name,
		"ClientCompany": pa.getCompanyOrDefault(inquiry.Company),
		"Services":      strings.Join(inquiry.Services, ", "),
		"Message":       inquiry.Message,
		"Priority":      inquiry.Priority,
	}
}

// prepareRiskAssessmentVariables prepares variables for risk assessment template rendering
func (pa *promptArchitect) prepareRiskAssessmentVariables(inquiry *domain.Inquiry) map[string]interface{} {
	return map[string]interface{}{
		"ClientName":    inquiry.Name,
		"ClientCompany": pa.getCompanyOrDefault(inquiry.Company),
		"Services":      strings.Join(inquiry.Services, ", "),
		"Message":       inquiry.Message,
		"IndustryHints": pa.extractIndustryHints(inquiry.Message, inquiry.Company),
	}
}

// prepareCompetitiveAnalysisVariables prepares variables for competitive analysis template rendering
func (pa *promptArchitect) prepareCompetitiveAnalysisVariables(inquiry *domain.Inquiry) map[string]interface{} {
	return map[string]interface{}{
		"ClientName":    inquiry.Name,
		"ClientCompany": pa.getCompanyOrDefault(inquiry.Company),
		"Services":      strings.Join(inquiry.Services, ", "),
		"Message":       inquiry.Message,
		"UseCase":       pa.extractUseCase(inquiry.Services, inquiry.Message),
	}
}

// renderTemplate renders a template with the provided variables
func (pa *promptArchitect) renderTemplate(promptTemplate *interfaces.PromptTemplate, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New(promptTemplate.Name).Parse(promptTemplate.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// validateTemplateSyntax validates the syntax of a template
func (pa *promptArchitect) validateTemplateSyntax(promptTemplate *interfaces.PromptTemplate) error {
	_, err := template.New("validation").Parse(promptTemplate.Template)
	if err != nil {
		return fmt.Errorf("invalid template syntax: %w", err)
	}

	// Check for required variables
	for _, variable := range promptTemplate.RequiredVariables {
		if !strings.Contains(promptTemplate.Template, "{{."+variable+"}}") &&
		   !strings.Contains(promptTemplate.Template, "{{ ."+variable+" }}") {
			return fmt.Errorf("template missing required variable: %s", variable)
		}
	}

	return nil
}

// Helper functions

func (pa *promptArchitect) getCompanyOrDefault(company string) string {
	if company != "" {
		return company
	}
	return "Client Organization"
}

func (pa *promptArchitect) getPhoneOrDefault(phone string) string {
	if phone != "" {
		return phone
	}
	return "Not provided"
}

func (pa *promptArchitect) extractIndustryHints(message, company string) string {
	text := strings.ToLower(message + " " + company)
	
	industryKeywords := map[string][]string{
		"Healthcare": {"healthcare", "hospital", "medical", "patient", "hipaa", "health"},
		"Financial": {"bank", "financial", "finance", "payment", "pci", "trading"},
		"Retail": {"retail", "ecommerce", "store", "shopping", "customer"},
		"Manufacturing": {"manufacturing", "factory", "production", "supply chain"},
		"Education": {"education", "school", "university", "student", "learning"},
		"Government": {"government", "public", "federal", "state", "municipal"},
	}
	
	for industry, keywords := range industryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				return industry
			}
		}
	}
	
	return "General"
}

func (pa *promptArchitect) extractUseCase(services []string, message string) string {
	text := strings.ToLower(message)
	
	if len(services) > 0 {
		return strings.Title(services[0])
	}
	
	useCases := map[string][]string{
		"Migration": {"migrate", "migration", "move", "transfer"},
		"Optimization": {"optimize", "cost", "performance", "efficiency"},
		"Assessment": {"assess", "review", "evaluate", "audit"},
		"Architecture": {"architecture", "design", "structure", "blueprint"},
	}
	
	for useCase, keywords := range useCases {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				return useCase
			}
		}
	}
	
	return "General Consulting"
}

// createProfileFromAudience creates an audience profile from explicitly specified audience type
func (pa *promptArchitect) createProfileFromAudience(audienceType string) *AudienceProfile {
	switch strings.ToLower(audienceType) {
	case "technical":
		return &AudienceProfile{
			PrimaryType:    AudienceTechnical,
			TechnicalDepth: 4,
			BusinessFocus:  2,
			Confidence:     0.9,
			Indicators:     []string{"Explicitly specified as technical"},
		}
	case "business":
		return &AudienceProfile{
			PrimaryType:    AudienceBusiness,
			TechnicalDepth: 2,
			BusinessFocus:  4,
			Confidence:     0.9,
			Indicators:     []string{"Explicitly specified as business"},
		}
	case "executive":
		return &AudienceProfile{
			PrimaryType:    AudienceExecutive,
			TechnicalDepth: 1,
			BusinessFocus:  5,
			ExecutiveLevel: true,
			Confidence:     0.9,
			Indicators:     []string{"Explicitly specified as executive"},
		}
	default: // "mixed" or unknown
		return &AudienceProfile{
			PrimaryType:    AudienceMixed,
			TechnicalDepth: 3,
			BusinessFocus:  3,
			Confidence:     0.7,
			Indicators:     []string{"Explicitly specified as mixed or default"},
		}
	}
}

// initializeDefaultTemplates sets up the default prompt templates
func (pa *promptArchitect) initializeDefaultTemplates() {
	// Enhanced Report Template
	enhancedReportTemplate := &interfaces.PromptTemplate{
		Name:        "enhanced_report",
		Category:    "report",
		Description: "Enhanced report template with documentation links and multi-cloud support",
		Template: `Generate a professional consulting report for the following client inquiry:

Client: {{.ClientName}} ({{.ClientEmail}})
Company: {{.ClientCompany}}
Phone: {{.ClientPhone}}
Services Requested: {{.Services}}
Message: {{.Message}}

{{if .MultiCloud}}
MULTI-CLOUD ANALYSIS: This report should consider solutions across {{.CloudProviders}} and provide comparative analysis where relevant.
{{end}}

{{if .IncludeDocumentationLinks}}
DOCUMENTATION REQUIREMENTS: Include specific links to official cloud provider documentation for all recommendations. Reference AWS Well-Architected Framework, Azure Architecture Center, and GCP best practices guides where applicable.
{{end}}

PRIORITY ANALYSIS: Carefully analyze the client's message for urgency indicators and meeting requests. Look for:
- Time-sensitive language: "urgent", "ASAP", "immediately", "today", "tomorrow", "this week"
- Meeting requests: "schedule", "meeting", "call", "discuss", "talk", "available"
- Specific dates/times mentioned
- Business impact language: "critical", "blocking", "emergency", "deadline"

Please provide a structured report with the following sections:

1. EXECUTIVE SUMMARY
   - Brief overview of the client's needs
   - Key recommendations summary
   {{if .IncludeRiskAssessment}}- Risk level assessment{{end}}
   - **PRIORITY LEVEL**: If urgent language or immediate meeting requests are detected, clearly mark as "HIGH PRIORITY" and explain why

2. CURRENT STATE ASSESSMENT
   - Analysis of the client's current situation based on their inquiry
   - Identified challenges and opportunities
   {{if .IndustryContext}}- Industry-specific considerations for {{.IndustryContext}}{{end}}

3. RECOMMENDATIONS
   - Specific actionable recommendations for each requested service
   {{if .MultiCloud}}- Multi-cloud provider comparison and recommendations{{end}}
   {{if .IncludeDocumentationLinks}}- Include specific documentation links for each recommendation{{end}}
   - Prioritized implementation approach

{{if .IncludeRiskAssessment}}
4. RISK ASSESSMENT
   - Technical risks and mitigation strategies
   - Security considerations
   - Compliance requirements
   - Business continuity factors
{{end}}

{{if .IncludeCompetitiveAnalysis}}
5. COMPETITIVE ANALYSIS
   - Cloud provider comparison for this specific use case
   - Pros and cons of each platform
   - Cost comparison where relevant
   - Recommendation rationale
{{end}}

{{if .IncludeImplementationSteps}}
6. IMPLEMENTATION ROADMAP
   - Phase-based implementation plan
   - Timeline estimates
   - Resource requirements
   - Key milestones and deliverables
{{end}}

7. NEXT STEPS
   - Immediate actions the client should consider
   - Proposed engagement timeline
   - **MEETING SCHEDULING**: If the client requested a meeting or mentioned specific dates/times, highlight this prominently

8. URGENCY ASSESSMENT
   - If any urgency indicators were found, create a special section highlighting:
     * Specific urgent language detected
     * Requested meeting timeframes
     * Recommended response timeline
     * Any specific dates/times mentioned by the client

IMPORTANT: At the end of the report, include a "Contact Information" section with the client's actual details for follow-up:
- Use the client's name: {{.ClientName}}
- Use the client's email: {{.ClientEmail}}
- Use the client's company: {{.ClientCompany}}
- Use the client's phone: {{.ClientPhone}}

Do NOT use placeholder text like [Your Name] or [Your Contact Information]. Use the actual client information provided above.

AUDIENCE ADAPTATION:
{{if .AudienceProfile}}
- Detected audience: {{.AudienceProfile.PrimaryType}} (confidence: {{.AudienceConfidence}})
- Technical depth level: {{.TechnicalDepth}}/5
- Business focus level: {{.BusinessFocus}}/5
{{if .AudienceIndicators}}- Detection indicators: {{.AudienceIndicators}}{{end}}
{{end}}

{{if eq .TargetAudience "technical"}}
TECHNICAL FOCUS: This report should emphasize technical architecture, implementation details, performance specifications, and operational considerations. Include specific service configurations, integration patterns, and technical best practices.
{{else if eq .TargetAudience "business"}}
BUSINESS FOCUS: This report should emphasize business value, ROI, cost-benefit analysis, and strategic implications. Minimize technical jargon and focus on business outcomes and competitive advantages.
{{else if eq .TargetAudience "executive"}}
EXECUTIVE FOCUS: This report should be high-level and strategic, focusing on business transformation, competitive positioning, investment requirements, and strategic alignment. Keep technical details minimal and focus on business impact.
{{else}}
MIXED AUDIENCE: This report should balance technical details with business justification. Provide both technical implementation guidance and business value explanation. Separate technical and business sections clearly.
{{end}}

Keep the tone professional and focus on actionable insights. The report should be comprehensive but concise, suitable for a {{.TargetAudience}} audience.`,
		RequiredVariables: []string{"ClientName", "ClientEmail", "ClientCompany", "ClientPhone", "Services", "Message"},
		OptionalVariables: []string{"MultiCloud", "CloudProviders", "IncludeDocumentationLinks", "IncludeRiskAssessment", "IncludeCompetitiveAnalysis", "IncludeImplementationSteps", "IndustryContext", "TargetAudience"},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasClientInfo", Pattern: `{{\.ClientName}}`, ErrorMessage: "Template must include client name", Required: true},
			{Name: "HasRecommendations", Pattern: `RECOMMENDATIONS`, ErrorMessage: "Template must include recommendations section", Required: true},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Technical Report Template
	technicalReportTemplate := &interfaces.PromptTemplate{
		Name:        "technical_report",
		Category:    "report",
		Description: "Technical report template focused on technical stakeholders",
		Template: `Generate a technical consulting report for the following client inquiry:

Client: {{.ClientName}} ({{.ClientEmail}})
Company: {{.ClientCompany}}
Services Requested: {{.Services}}
Technical Requirements: {{.Message}}

This report is targeted at technical stakeholders and should include detailed technical specifications, architecture diagrams (described in text), and implementation details.

Please provide a structured technical report with the following sections:

1. TECHNICAL EXECUTIVE SUMMARY
   - Technical overview of requirements
   - Recommended technical architecture
   - Key technical decisions and rationale

2. TECHNICAL REQUIREMENTS ANALYSIS
   - Detailed analysis of technical requirements
   - Performance, scalability, and reliability considerations
   - Security and compliance technical requirements

3. ARCHITECTURE RECOMMENDATIONS
   - Detailed technical architecture description
   - Service selection with technical justification
   - Integration patterns and data flow
   - Network architecture and security boundaries

4. IMPLEMENTATION SPECIFICATIONS
   - Detailed implementation steps with technical specifics
   - Configuration requirements
   - Code examples or pseudocode where relevant
   - Testing and validation approaches

5. TECHNICAL RISKS AND MITIGATION
   - Technical risks assessment
   - Performance bottlenecks and solutions
   - Security vulnerabilities and controls
   - Operational considerations

6. TECHNICAL NEXT STEPS
   - Immediate technical actions
   - Development timeline with technical milestones
   - Required technical resources and skills

Contact Information:
- Client: {{.ClientName}} ({{.ClientEmail}})
- Company: {{.ClientCompany}}
- Phone: {{.ClientPhone}}`,
		RequiredVariables: []string{"ClientName", "ClientEmail", "ClientCompany", "ClientPhone", "Services", "Message"},
		OptionalVariables: []string{},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasTechnicalFocus", Pattern: `TECHNICAL`, ErrorMessage: "Template must have technical focus", Required: true},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Business Report Template
	businessReportTemplate := &interfaces.PromptTemplate{
		Name:        "business_report",
		Category:    "report",
		Description: "Business report template focused on business stakeholders",
		Template: `Generate a business-focused consulting report for the following client inquiry:

Client: {{.ClientName}} ({{.ClientEmail}})
Company: {{.ClientCompany}}
Business Requirements: {{.Message}}

This report is targeted at business stakeholders and should focus on business value, ROI, and strategic implications.

Please provide a structured business report with the following sections:

1. BUSINESS EXECUTIVE SUMMARY
   - Business overview of the opportunity
   - Expected business value and ROI
   - Strategic recommendations

2. BUSINESS CASE ANALYSIS
   - Current business challenges
   - Proposed solution benefits
   - Cost-benefit analysis
   - Risk assessment from business perspective

3. STRATEGIC RECOMMENDATIONS
   - Business-aligned cloud strategy
   - Competitive advantages
   - Market positioning implications
   - Scalability for business growth

4. FINANCIAL ANALYSIS
   - Investment requirements
   - Expected cost savings
   - ROI projections
   - Budget considerations

5. BUSINESS RISKS AND MITIGATION
   - Business continuity risks
   - Market and competitive risks
   - Operational risks
   - Change management considerations

6. BUSINESS NEXT STEPS
   - Business decision points
   - Stakeholder engagement plan
   - Timeline for business value realization
   - Success metrics and KPIs

Contact Information:
- Client: {{.ClientName}} ({{.ClientEmail}})
- Company: {{.ClientCompany}}
- Phone: {{.ClientPhone}}`,
		RequiredVariables: []string{"ClientName", "ClientEmail", "ClientCompany", "ClientPhone", "Services", "Message"},
		OptionalVariables: []string{},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasBusinessFocus", Pattern: `BUSINESS`, ErrorMessage: "Template must have business focus", Required: true},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Interview Guide Template
	interviewGuideTemplate := &interfaces.PromptTemplate{
		Name:        "interview_guide",
		Category:    "interview",
		Description: "Template for generating interview guides and discovery questions",
		Template: `Generate a comprehensive interview guide for the following client engagement:

Client: {{.ClientName}}
Company: {{.ClientCompany}}
Services Requested: {{.Services}}
Initial Message: {{.Message}}

Create a structured interview guide that will help gather detailed requirements and provide better recommendations. The guide should be organized into sections with specific questions designed to uncover technical, business, and operational requirements.

Please provide an interview guide with the following structure:

1. INTERVIEW OBJECTIVES
   - Primary goals for this discovery session
   - Key information to gather
   - Expected outcomes

2. PRE-INTERVIEW PREPARATION
   - Research to conduct beforehand
   - Documents to request from client
   - Technical artifacts needed

3. BUSINESS CONTEXT QUESTIONS
   - Company background and industry
   - Current business challenges
   - Strategic objectives and goals
   - Success criteria and metrics

4. TECHNICAL DISCOVERY QUESTIONS
   - Current infrastructure and architecture
   - Existing cloud usage and experience
   - Technical constraints and requirements
   - Performance and scalability needs

5. OPERATIONAL QUESTIONS
   - Current operational processes
   - Team structure and capabilities
   - Compliance and security requirements
   - Budget and timeline constraints

6. FOLLOW-UP QUESTIONS
   - Probing questions based on initial responses
   - Clarification questions for ambiguous areas
   - Questions to identify hidden requirements

7. NEXT STEPS DISCUSSION
   - Proposed engagement approach
   - Timeline and deliverables
   - Resource requirements
   - Decision-making process

8. POST-INTERVIEW ACTIONS
   - Information to validate
   - Additional stakeholders to engage
   - Documents to review
   - Follow-up meeting schedule

Focus on open-ended questions that encourage detailed responses and help build a comprehensive understanding of the client's needs.`,
		RequiredVariables: []string{"ClientName", "ClientCompany", "Services", "Message"},
		OptionalVariables: []string{"Priority"},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasQuestions", Pattern: `QUESTIONS`, ErrorMessage: "Template must include question sections", Required: true},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Risk Assessment Template
	riskAssessmentTemplate := &interfaces.PromptTemplate{
		Name:        "risk_assessment",
		Category:    "analysis",
		Description: "Template for generating comprehensive risk assessments",
		Template: `Generate a comprehensive risk assessment for the following client scenario:

Client: {{.ClientName}}
Company: {{.ClientCompany}}
Services: {{.Services}}
Requirements: {{.Message}}
{{if .IndustryHints}}Industry Context: {{.IndustryHints}}{{end}}

Analyze the potential risks associated with the proposed cloud initiative and provide detailed mitigation strategies.

Please provide a structured risk assessment with the following sections:

1. RISK ASSESSMENT OVERVIEW
   - Scope of risk analysis
   - Assessment methodology
   - Risk rating criteria

2. TECHNICAL RISKS
   - Architecture and design risks
   - Performance and scalability risks
   - Integration and compatibility risks
   - Data migration and synchronization risks
   - Mitigation strategies for each technical risk

3. SECURITY RISKS
   - Data security and privacy risks
   - Access control and identity management risks
   - Network security risks
   - Compliance and regulatory risks
   - Security mitigation strategies

4. OPERATIONAL RISKS
   - Service availability and reliability risks
   - Operational complexity risks
   - Skills and knowledge gaps
   - Change management risks
   - Operational mitigation strategies

5. BUSINESS RISKS
   - Cost overrun and budget risks
   - Timeline and delivery risks
   - Vendor lock-in risks
   - Business continuity risks
   - Strategic alignment risks
   - Business mitigation strategies

6. COMPLIANCE RISKS
   {{if .IndustryHints}}- Industry-specific compliance requirements for {{.IndustryHints}}{{end}}
   - Data governance and privacy regulations
   - Audit and reporting requirements
   - Compliance mitigation strategies

7. RISK PRIORITIZATION
   - High-priority risks requiring immediate attention
   - Medium-priority risks for ongoing monitoring
   - Low-priority risks for periodic review
   - Risk interdependencies and cascading effects

8. RISK MONITORING AND GOVERNANCE
   - Risk monitoring framework
   - Key risk indicators (KRIs)
   - Escalation procedures
   - Regular review and update processes

Provide specific, actionable mitigation strategies for each identified risk, including responsible parties, timelines, and success criteria.`,
		RequiredVariables: []string{"ClientName", "ClientCompany", "Services", "Message"},
		OptionalVariables: []string{"IndustryHints"},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasRiskSections", Pattern: `RISKS`, ErrorMessage: "Template must include risk sections", Required: true},
		},
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	// Competitive Analysis Template
	competitiveAnalysisTemplate := &interfaces.PromptTemplate{
		Name:        "competitive_analysis",
		Category:    "analysis",
		Description: "Template for generating competitive analysis between cloud providers",
		Template: `Generate a comprehensive competitive analysis for the following client scenario:

Client: {{.ClientName}}
Company: {{.ClientCompany}}
Use Case: {{.UseCase}}
Requirements: {{.Message}}

Provide an objective comparison of major cloud providers (AWS, Azure, GCP) for this specific use case, helping the client make an informed decision.

Please provide a structured competitive analysis with the following sections:

1. ANALYSIS OVERVIEW
   - Scope of competitive analysis
   - Evaluation criteria and methodology
   - Key decision factors

2. PROVIDER COMPARISON MATRIX
   - Service capabilities comparison
   - Feature and functionality analysis
   - Performance characteristics
   - Pricing model comparison

3. AWS ANALYSIS
   - Strengths for this use case
   - Relevant services and features
   - Pricing considerations
   - Implementation complexity
   - Documentation and support quality

4. MICROSOFT AZURE ANALYSIS
   - Strengths for this use case
   - Relevant services and features
   - Pricing considerations
   - Implementation complexity
   - Documentation and support quality

5. GOOGLE CLOUD PLATFORM ANALYSIS
   - Strengths for this use case
   - Relevant services and features
   - Pricing considerations
   - Implementation complexity
   - Documentation and support quality

6. SCENARIO-BASED RECOMMENDATIONS
   - Best choice for cost optimization
   - Best choice for performance
   - Best choice for specific technical requirements
   - Best choice for compliance and security
   - Best choice for ease of implementation

7. MULTI-CLOUD CONSIDERATIONS
   - Benefits of multi-cloud approach
   - Complexity and management overhead
   - Data portability and vendor lock-in
   - Integration and interoperability

8. FINAL RECOMMENDATION
   - Primary recommended provider with rationale
   - Alternative options and scenarios
   - Risk factors and mitigation strategies
   - Next steps for decision making

Include specific service names, pricing estimates where possible, and links to relevant documentation for each provider.`,
		RequiredVariables: []string{"ClientName", "ClientCompany", "UseCase", "Message"},
		OptionalVariables: []string{},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasProviderAnalysis", Pattern: `AWS|AZURE|GCP`, ErrorMessage: "Template must include provider analysis", Required: true},
		},
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	// Comprehensive Report Template
	comprehensiveReportTemplate := &interfaces.PromptTemplate{
		Name:        "comprehensive_report",
		Category:    "report",
		Description: "Comprehensive report template with all enhanced features",
		Template: `Generate a comprehensive consulting report for the following client inquiry:

Client: {{.ClientName}} ({{.ClientEmail}})
Company: {{.ClientCompany}}
Phone: {{.ClientPhone}}
Services Requested: {{.Services}}
Message: {{.Message}}

This is a comprehensive analysis that should include detailed recommendations, risk assessment, competitive analysis, and implementation planning.

MULTI-CLOUD ANALYSIS: Consider solutions across AWS, Azure, and GCP, providing comparative analysis and recommendations.

DOCUMENTATION REQUIREMENTS: Include specific links to official cloud provider documentation for all recommendations.

Please provide a comprehensive structured report with the following sections:

1. EXECUTIVE SUMMARY
   - Overview of client needs and proposed solutions
   - Key recommendations and expected outcomes
   - Risk level assessment and priority indicators
   - Investment requirements and expected ROI

2. CURRENT STATE ASSESSMENT
   - Analysis of current infrastructure and processes
   - Identified challenges and pain points
   - Opportunities for improvement
   - Baseline metrics and benchmarks

3. DETAILED RECOMMENDATIONS
   - Specific actionable recommendations for each service area
   - Multi-cloud provider comparison and selection rationale
   - Architecture and design recommendations
   - Technology stack and service selections
   - Include specific documentation links for each recommendation

4. COMPETITIVE ANALYSIS
   - Detailed comparison of AWS, Azure, and GCP for this use case
   - Strengths and weaknesses of each platform
   - Cost comparison and pricing analysis
   - Feature and capability comparison
   - Final provider recommendation with rationale

5. COMPREHENSIVE RISK ASSESSMENT
   - Technical risks and mitigation strategies
   - Security and compliance risks
   - Operational and business risks
   - Risk prioritization and monitoring approach
   - Contingency planning recommendations

6. IMPLEMENTATION ROADMAP
   - Phase-based implementation plan with timelines
   - Resource requirements and skill sets needed
   - Key milestones and deliverables
   - Dependencies and critical path analysis
   - Success criteria and measurement approach

7. FINANCIAL ANALYSIS
   - Detailed cost breakdown and estimates
   - ROI analysis and payback period
   - Cost optimization opportunities
   - Budget planning and cash flow considerations

8. NEXT STEPS AND ENGAGEMENT PLAN
   - Immediate actions and quick wins
   - Proposed engagement approach and timeline
   - Stakeholder engagement and communication plan
   - Decision points and approval processes

9. APPENDICES
   - Technical specifications and requirements
   - Vendor comparison matrices
   - Risk register and mitigation plans
   - Implementation checklists and templates

IMPORTANT: Include actual client contact information:
- Client: {{.ClientName}} ({{.ClientEmail}})
- Company: {{.ClientCompany}}
- Phone: {{.ClientPhone}}

Maintain a professional tone while ensuring all recommendations are specific, actionable, and well-documented with official cloud provider references.`,
		RequiredVariables: []string{"ClientName", "ClientEmail", "ClientCompany", "ClientPhone", "Services", "Message"},
		OptionalVariables: []string{},
		ValidationRules: []interfaces.ValidationRule{
			{Name: "HasComprehensiveSections", Pattern: `COMPETITIVE ANALYSIS.*RISK ASSESSMENT.*IMPLEMENTATION ROADMAP`, ErrorMessage: "Template must include all comprehensive sections", Required: true},
		},
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	// Register all templates
	templates := []*interfaces.PromptTemplate{
		enhancedReportTemplate,
		technicalReportTemplate,
		businessReportTemplate,
		interviewGuideTemplate,
		riskAssessmentTemplate,
		competitiveAnalysisTemplate,
		comprehensiveReportTemplate,
	}

	for _, template := range templates {
		pa.templates[template.Name] = template
	}
}