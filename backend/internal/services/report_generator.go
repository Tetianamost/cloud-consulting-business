package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/google/uuid"
)

// reportGenerator implements enhanced report generation using Bedrock with AI assistant capabilities
type reportGenerator struct {
	bedrockService     interfaces.BedrockService
	templateService    interfaces.TemplateService
	pdfService         interfaces.PDFService
	promptArchitect    interfaces.PromptArchitect
	knowledgeBase      interfaces.KnowledgeBase
	multiCloudAnalyzer interfaces.MultiCloudAnalyzer
	riskAssessor       interfaces.RiskAssessor
	docLibrary         interfaces.DocumentationLibrary
	audienceDetector   AudienceDetector
}

// NewReportGenerator creates a new enhanced report generator instance
func NewReportGenerator(
	bedrockService interfaces.BedrockService,
	templateService interfaces.TemplateService,
	pdfService interfaces.PDFService,
	promptArchitect interfaces.PromptArchitect,
	knowledgeBase interfaces.KnowledgeBase,
	multiCloudAnalyzer interfaces.MultiCloudAnalyzer,
	riskAssessor interfaces.RiskAssessor,
	docLibrary interfaces.DocumentationLibrary,
) interfaces.ReportService {
	return &reportGenerator{
		bedrockService:     bedrockService,
		templateService:    templateService,
		pdfService:         pdfService,
		promptArchitect:    promptArchitect,
		knowledgeBase:      knowledgeBase,
		multiCloudAnalyzer: multiCloudAnalyzer,
		riskAssessor:       riskAssessor,
		docLibrary:         docLibrary,
		audienceDetector:   NewAudienceDetector(),
	}
}

// NewBasicReportGenerator creates a basic report generator for backward compatibility
func NewBasicReportGenerator(bedrockService interfaces.BedrockService, templateService interfaces.TemplateService, pdfService interfaces.PDFService) interfaces.ReportService {
	return &reportGenerator{
		bedrockService:  bedrockService,
		templateService: templateService,
		pdfService:      pdfService,
	}
}

// GenerateReport generates an enhanced report for the given inquiry using AI assistant capabilities
func (r *reportGenerator) GenerateReport(ctx context.Context, inquiry *domain.Inquiry) (*domain.Report, error) {
	// Use enhanced prompt generation if PromptArchitect is available
	var prompt string
	var err error

	if r.promptArchitect != nil {
		prompt, err = r.buildEnhancedPrompt(ctx, inquiry)
		if err != nil {
			// Fall back to basic prompt if enhanced fails
			prompt = r.buildPrompt(inquiry)
		}
	} else {
		// Use basic prompt for backward compatibility
		prompt = r.buildPrompt(inquiry)
	}

	// Set Bedrock options with higher token limit for enhanced reports
	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   4000, // Increased for enhanced content
		Temperature: 0.7,
		TopP:        0.9,
	}

	// Generate content using Bedrock
	response, err := r.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate report with Bedrock: %w", err)
	}

	// Create the enhanced report
	report := &domain.Report{
		ID:          uuid.New().String(),
		InquiryID:   inquiry.ID,
		Type:        r.getReportType(inquiry.Services),
		Title:       r.generateTitle(inquiry),
		Content:     response.Content,
		Status:      domain.ReportStatusDraft,
		GeneratedBy: "enhanced-bedrock-ai",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return report, nil
}

// buildPrompt creates a structured prompt for Bedrock based on the inquiry
func (r *reportGenerator) buildPrompt(inquiry *domain.Inquiry) string {
	template := `Generate a professional consulting report draft for the following client inquiry:

Client: %s (%s)
Company: %s
Phone: %s
Services Requested: %s
Message: %s

PRIORITY ANALYSIS: Carefully analyze the client's message for urgency indicators and meeting requests. Look for:
- Time-sensitive language: "urgent", "ASAP", "immediately", "today", "tomorrow", "this week"
- Meeting requests: "schedule", "meeting", "call", "discuss", "talk", "available"
- Specific dates/times mentioned
- Business impact language: "critical", "blocking", "emergency", "deadline"

Please provide a structured report with the following sections:

1. EXECUTIVE SUMMARY
   - Brief overview of the client's needs
   - Key recommendations summary
   - **PRIORITY LEVEL**: If urgent language or immediate meeting requests are detected, clearly mark as "HIGH PRIORITY" and explain why

2. CURRENT STATE ASSESSMENT
   - Analysis of the client's current situation based on their inquiry
   - Identified challenges and opportunities

3. RECOMMENDATIONS
   - Specific actionable recommendations for each requested service
   - Prioritized implementation approach

4. NEXT STEPS
   - Immediate actions the client should consider
   - Proposed engagement timeline
   - **MEETING SCHEDULING**: If the client requested a meeting or mentioned specific dates/times, highlight this prominently

5. URGENCY ASSESSMENT
   - If any urgency indicators were found, create a special section highlighting:
     * Specific urgent language detected
     * Requested meeting timeframes
     * Recommended response timeline
     * Any specific dates/times mentioned by the client

IMPORTANT: At the end of the report, include a "Contact Information" section with the client's actual details for follow-up:
- Use the client's name: %s
- Use the client's email: %s
- Use the client's company: %s
- Use the client's phone: %s

Do NOT use placeholder text like [Your Name] or [Your Contact Information]. Use the actual client information provided above.

Keep the tone professional and focus on actionable insights. The report should be comprehensive but concise, suitable for a business audience.`

	return fmt.Sprintf(template,
		inquiry.Name,
		inquiry.Email,
		r.getCompanyOrDefault(inquiry.Company),
		r.getPhoneOrDefault(inquiry.Phone),
		strings.Join(inquiry.Services, ", "),
		inquiry.Message,
		inquiry.Name,
		inquiry.Email,
		r.getCompanyOrDefault(inquiry.Company),
		r.getPhoneOrDefault(inquiry.Phone),
	)
}

// generateTitle creates a title for the report based on the inquiry
func (r *reportGenerator) generateTitle(inquiry *domain.Inquiry) string {
	serviceType := r.getReportType(inquiry.Services)
	companyName := r.getCompanyOrDefault(inquiry.Company)

	return fmt.Sprintf("%s Report for %s",
		strings.Title(string(serviceType)),
		companyName)
}

// getReportType determines the primary report type based on services requested
func (r *reportGenerator) getReportType(services []string) domain.ReportType {
	if len(services) == 0 {
		return domain.ReportTypeGeneral
	}

	// Use the first service as the primary type
	switch strings.ToLower(services[0]) {
	case domain.ServiceTypeAssessment:
		return domain.ReportTypeAssessment
	case domain.ServiceTypeMigration:
		return domain.ReportTypeMigration
	case domain.ServiceTypeOptimization:
		return domain.ReportTypeOptimization
	case domain.ServiceTypeArchitectureReview:
		return domain.ReportTypeArchitecture
	default:
		return domain.ReportTypeGeneral
	}
}

// getCompanyOrDefault returns the company name or a default value
func (r *reportGenerator) getCompanyOrDefault(company string) string {
	if company != "" {
		return company
	}
	return "Client Organization"
}

// getPhoneOrDefault returns the phone number or a default value
func (r *reportGenerator) getPhoneOrDefault(phone string) string {
	if phone != "" {
		return phone
	}
	return "Not provided"
}

// GetReport retrieves a report by ID (placeholder implementation)
func (r *reportGenerator) GetReport(ctx context.Context, id string) (*domain.Report, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetReportsByInquiry retrieves reports for a specific inquiry (placeholder implementation)
func (r *reportGenerator) GetReportsByInquiry(ctx context.Context, inquiryID string) ([]*domain.Report, error) {
	return nil, fmt.Errorf("not implemented")
}

// UpdateReportStatus updates the status of a report (placeholder implementation)
func (r *reportGenerator) UpdateReportStatus(ctx context.Context, id string, status domain.ReportStatus) error {
	return fmt.Errorf("not implemented")
}

// GetReportTemplate retrieves a report template (placeholder implementation)
func (r *reportGenerator) GetReportTemplate(serviceType domain.ServiceType) (*interfaces.ReportTemplate, error) {
	return nil, fmt.Errorf("not implemented")
}

// ValidateReport validates a report (placeholder implementation)
func (r *reportGenerator) ValidateReport(report *domain.Report) error {
	if report.Content == "" {
		return fmt.Errorf("report content cannot be empty")
	}
	return nil
}

// GenerateHTML generates HTML version of a report
func (r *reportGenerator) GenerateHTML(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (string, error) {
	if r.templateService == nil {
		return "", fmt.Errorf("template service not available")
	}

	// Determine template name based on report type
	templateName := r.getTemplateName(report.Type)

	// Prepare template data
	templateData := r.prepareTemplateData(inquiry, report)

	// Render the template
	htmlContent, err := r.templateService.RenderReportTemplate(ctx, templateName, templateData)
	if err != nil {
		return "", fmt.Errorf("failed to render HTML report: %w", err)
	}

	return htmlContent, nil
}

// getTemplateName returns the template name for a report type
func (r *reportGenerator) getTemplateName(reportType domain.ReportType) string {
	switch reportType {
	case domain.ReportTypeAssessment:
		return "assessment"
	case domain.ReportTypeMigration:
		return "migration"
	case domain.ReportTypeOptimization:
		return "optimization"
	case domain.ReportTypeArchitecture:
		return "architecture"
	default:
		return "assessment" // Default fallback
	}
}

// prepareTemplateData prepares data for template rendering
func (r *reportGenerator) prepareTemplateData(inquiry *domain.Inquiry, report *domain.Report) interface{} {
	// Use the template service to prepare data if available
	if r.templateService != nil {
		return r.templateService.PrepareReportTemplateData(inquiry, report)
	}

	// Fallback to basic data structure
	return map[string]interface{}{
		"ID":               report.ID,
		"Title":            report.Title,
		"ClientName":       inquiry.Name,
		"ClientEmail":      inquiry.Email,
		"ClientCompany":    r.getCompanyOrDefault(inquiry.Company),
		"ClientPhone":      r.getPhoneOrDefault(inquiry.Phone),
		"Services":         strings.Join(inquiry.Services, ", "),
		"GeneratedDate":    report.CreatedAt.Format("January 2, 2006"),
		"IsHighPriority":   r.detectHighPriority(report.Content),
		"FormattedContent": r.formatContentForHTML(report.Content),
	}
}

// detectHighPriority analyzes report content for priority indicators
func (r *reportGenerator) detectHighPriority(content string) bool {
	priorityKeywords := []string{
		"HIGH PRIORITY", "URGENT", "CRITICAL", "IMMEDIATE", "ASAP",
		"urgent", "critical", "emergency", "deadline", "time-sensitive",
		"meeting", "schedule", "call", "discuss", "today", "tomorrow",
	}

	contentLower := strings.ToLower(content)
	for _, keyword := range priorityKeywords {
		if strings.Contains(contentLower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// formatContentForHTML converts plain text content to structured HTML
func (r *reportGenerator) formatContentForHTML(content string) string {
	if content == "" {
		return "<p>No content available.</p>"
	}

	// Enhanced HTML formatting with better structure detection
	sections := strings.Split(content, "\n\n")
	var htmlSections []string

	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}

		// Format the section based on its content type
		formattedSection := r.formatSection(section)
		htmlSections = append(htmlSections, formattedSection)
	}

	return strings.Join(htmlSections, "\n\n")
}

// formatSection formats a single section of content
func (r *reportGenerator) formatSection(section string) string {
	// Check if this is a main header (numbered or all caps)
	if r.isMainHeader(section) {
		return r.formatMainHeader(section)
	}

	// Check if this is a sub-header
	if r.isSubHeader(section) {
		return r.formatSubHeader(section)
	}

	// Check if this contains a list
	if r.containsList(section) {
		return r.formatListSection(section)
	}

	// Format as regular paragraph content
	return r.formatParagraphSection(section)
}

// isMainHeader checks if a section is a main header
func (r *reportGenerator) isMainHeader(text string) bool {
	text = strings.TrimSpace(text)

	// Don't treat multi-line content as headers
	if strings.Count(text, "\n") > 1 {
		return false
	}

	// Check for numbered headers (1., 2., etc.)
	if regexp.MustCompile(`^\d+\.\s+[A-Z]`).MatchString(text) {
		return true
	}

	// Check for main section headers
	mainHeaders := []string{
		"EXECUTIVE SUMMARY", "CURRENT STATE ASSESSMENT", "RECOMMENDATIONS",
		"NEXT STEPS", "URGENCY ASSESSMENT", "CONTACT INFORMATION",
		"PRIORITY LEVEL", "MEETING SCHEDULING",
	}

	textUpper := strings.ToUpper(text)
	for _, header := range mainHeaders {
		if strings.Contains(textUpper, header) {
			return true
		}
	}

	return false
}

// isSubHeader checks if a section is a sub-header
func (r *reportGenerator) isSubHeader(text string) bool {
	text = strings.TrimSpace(text)

	// Single line ending with colon
	if !strings.Contains(text, "\n") && strings.HasSuffix(text, ":") && len(text) < 100 {
		return true
	}

	// Bold text that looks like a header
	if strings.HasPrefix(text, "**") && strings.HasSuffix(text, "**") && !strings.Contains(text, "\n") {
		return true
	}

	return false
}

// containsList checks if a section contains list items
func (r *reportGenerator) containsList(text string) bool {
	lines := strings.Split(text, "\n")
	listCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") ||
			regexp.MustCompile(`^\d+\.\s`).MatchString(line) {
			listCount++
		}
	}

	return listCount >= 2 // At least 2 list items
}

// formatMainHeader formats a main header
func (r *reportGenerator) formatMainHeader(text string) string {
	text = strings.TrimSpace(text)

	// Remove numbered prefixes for cleaner headers
	text = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(text, "")

	// Convert to title case if all caps
	if strings.ToUpper(text) == text {
		text = r.toTitleCase(text)
	}

	// Remove trailing colons
	text = strings.TrimSuffix(text, ":")

	return fmt.Sprintf("<h2 class=\"section-header\">%s</h2>", text)
}

// formatSubHeader formats a sub-header
func (r *reportGenerator) formatSubHeader(text string) string {
	text = strings.TrimSpace(text)

	// Remove bold markdown
	text = strings.Trim(text, "*")

	// Remove trailing colons
	text = strings.TrimSuffix(text, ":")

	return fmt.Sprintf("<h3 class=\"subsection-header\">%s</h3>", text)
}

// formatListSection formats a section containing lists
func (r *reportGenerator) formatListSection(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	var currentParagraph []string
	inList := false
	listType := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a list item
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") {
			// Close any open paragraph
			if len(currentParagraph) > 0 {
				result = append(result, r.formatParagraph(strings.Join(currentParagraph, " ")))
				currentParagraph = nil
			}

			// Start or continue unordered list
			if !inList || listType != "ul" {
				if inList && listType == "ol" {
					result = append(result, "</ol>")
				}
				if !inList {
					result = append(result, "<ul class=\"report-list\">")
				}
				inList = true
				listType = "ul"
			}

			item := strings.TrimPrefix(line, "- ")
			item = strings.TrimPrefix(item, "• ")
			result = append(result, fmt.Sprintf("  <li>%s</li>", r.formatInlineText(item)))

		} else if regexp.MustCompile(`^\d+\.\s`).MatchString(line) {
			// Close any open paragraph
			if len(currentParagraph) > 0 {
				result = append(result, r.formatParagraph(strings.Join(currentParagraph, " ")))
				currentParagraph = nil
			}

			// Start or continue ordered list
			if !inList || listType != "ol" {
				if inList && listType == "ul" {
					result = append(result, "</ul>")
				}
				if !inList {
					result = append(result, "<ol class=\"report-list\">")
				}
				inList = true
				listType = "ol"
			}

			item := regexp.MustCompile(`^\d+\.\s`).ReplaceAllString(line, "")
			result = append(result, fmt.Sprintf("  <li>%s</li>", r.formatInlineText(item)))

		} else {
			// Regular text - add to current paragraph
			currentParagraph = append(currentParagraph, line)
		}
	}

	// Close any open list
	if inList {
		if listType == "ul" {
			result = append(result, "</ul>")
		} else {
			result = append(result, "</ol>")
		}
	}

	// Add any remaining paragraph
	if len(currentParagraph) > 0 {
		result = append(result, r.formatParagraph(strings.Join(currentParagraph, " ")))
	}

	return strings.Join(result, "\n")
}

// formatParagraphSection formats a regular paragraph section
func (r *reportGenerator) formatParagraphSection(text string) string {
	// Split into individual lines and rejoin as a paragraph
	lines := strings.Split(text, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	if len(cleanLines) == 0 {
		return ""
	}

	paragraph := strings.Join(cleanLines, " ")
	return r.formatParagraph(paragraph)
}

// formatParagraph formats a single paragraph with inline formatting
func (r *reportGenerator) formatParagraph(text string) string {
	if text == "" {
		return ""
	}

	formatted := r.formatInlineText(text)
	return fmt.Sprintf("<p class=\"report-paragraph\">%s</p>", formatted)
}

// formatInlineText applies inline formatting (bold, italic, etc.)
func (r *reportGenerator) formatInlineText(text string) string {
	// Convert **bold** to <strong>
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "<strong>$1</strong>")

	// Convert *italic* to <em>
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "<em>$1</em>")

	// Convert URLs to links (basic implementation)
	text = regexp.MustCompile(`https?://[^\s]+`).ReplaceAllStringFunc(text, func(url string) string {
		return fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, url)
	})

	return text
}

// toTitleCase converts text to title case
func (r *reportGenerator) toTitleCase(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		if len(word) > 0 {
			// Keep certain words lowercase (articles, prepositions)
			lowercaseWords := map[string]bool{
				"a": true, "an": true, "the": true, "and": true, "or": true, "but": true,
				"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
				"with": true, "by": true,
			}

			if i == 0 || !lowercaseWords[strings.ToLower(word)] {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			} else {
				words[i] = strings.ToLower(word)
			}
		}
	}
	return strings.Join(words, " ")
}

// GeneratePDF generates a PDF version of a report
func (r *reportGenerator) GeneratePDF(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) ([]byte, error) {
	if r.pdfService == nil {
		return nil, fmt.Errorf("PDF service not available")
	}

	// First generate the HTML version
	htmlContent, err := r.GenerateHTML(ctx, inquiry, report)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML for PDF: %w", err)
	}

	// Get optimized PDF options for reports
	options := getEnhancedReportPDFOptions()

	// Generate PDF from HTML
	pdfBytes, err := r.pdfService.GeneratePDF(ctx, htmlContent, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBytes, nil
}

// buildEnhancedPrompt creates an enhanced prompt using AI assistant capabilities
func (r *reportGenerator) buildEnhancedPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	// Detect audience for the inquiry
	var audienceProfile *AudienceProfile
	var err error

	if r.audienceDetector != nil {
		audienceProfile, err = r.audienceDetector.DetectAudience(ctx, inquiry)
		if err != nil {
			// Fall back to mixed audience if detection fails
			audienceProfile = &AudienceProfile{
				PrimaryType:    AudienceMixed,
				TechnicalDepth: 3,
				BusinessFocus:  3,
				Confidence:     0.5,
			}
		}
	} else {
		// Default audience profile if detector not available
		audienceProfile = &AudienceProfile{
			PrimaryType:    AudienceMixed,
			TechnicalDepth: 3,
			BusinessFocus:  3,
			Confidence:     0.5,
		}
	}

	// Determine prompt options based on detected audience and available services
	options := &interfaces.PromptOptions{
		TargetAudience:             string(audienceProfile.PrimaryType),
		MaxTokens:                  4000,
		IncludeDocumentationLinks:  r.docLibrary != nil,
		IncludeCompetitiveAnalysis: r.multiCloudAnalyzer != nil,
		IncludeRiskAssessment:      r.riskAssessor != nil,
		IncludeImplementationSteps: true,
		CloudProviders:             []string{"AWS", "Azure", "GCP"},
	}

	// Extract industry context if available
	if r.knowledgeBase != nil {
		options.IndustryContext = r.extractIndustryContext(inquiry)
	}

	// Build the enhanced prompt using PromptArchitect (which will use audience detection internally)
	prompt, err := r.promptArchitect.BuildReportPrompt(ctx, inquiry, options)
	if err != nil {
		return "", fmt.Errorf("failed to build enhanced prompt: %w", err)
	}

	// Enhance the prompt with additional context from available services
	enhancedPrompt, err := r.enrichPromptWithContext(ctx, prompt, inquiry, options)
	if err != nil {
		// Log error but continue with base prompt
		return prompt, nil
	}

	return enhancedPrompt, nil
}

// enrichPromptWithContext adds additional context from knowledge base and other services
func (r *reportGenerator) enrichPromptWithContext(ctx context.Context, basePrompt string, inquiry *domain.Inquiry, options *interfaces.PromptOptions) (string, error) {
	var contextSections []string

	// Add knowledge base context
	if r.knowledgeBase != nil && options.IncludeDocumentationLinks {
		kbContext, err := r.buildKnowledgeBaseContext(ctx, inquiry)
		if err == nil && kbContext != "" {
			contextSections = append(contextSections, kbContext)
		}
	}

	// Add multi-cloud analysis context
	if r.multiCloudAnalyzer != nil && options.IncludeCompetitiveAnalysis {
		mcContext, err := r.buildMultiCloudContext(ctx, inquiry)
		if err == nil && mcContext != "" {
			contextSections = append(contextSections, mcContext)
		}
	}

	// Add risk assessment context
	if r.riskAssessor != nil && options.IncludeRiskAssessment {
		riskContext, err := r.buildRiskAssessmentContext(ctx, inquiry)
		if err == nil && riskContext != "" {
			contextSections = append(contextSections, riskContext)
		}
	}

	// Add documentation links context
	if r.docLibrary != nil && options.IncludeDocumentationLinks {
		docContext, err := r.buildDocumentationContext(ctx, inquiry)
		if err == nil && docContext != "" {
			contextSections = append(contextSections, docContext)
		}
	}

	// Combine base prompt with additional context
	if len(contextSections) > 0 {
		contextHeader := "\n\nADDITIONAL CONTEXT FOR ENHANCED RECOMMENDATIONS:\n"
		contextContent := strings.Join(contextSections, "\n\n")
		return basePrompt + contextHeader + contextContent, nil
	}

	return basePrompt, nil
}

// buildKnowledgeBaseContext builds context from the knowledge base
func (r *reportGenerator) buildKnowledgeBaseContext(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	var contextParts []string

	// Get best practices for the requested services
	for _, service := range inquiry.Services {
		bestPractices, err := r.knowledgeBase.GetBestPractices(ctx, service)
		if err == nil && len(bestPractices) > 0 {
			var practices []string
			for i, bp := range bestPractices {
				if i >= 3 { // Limit to top 3 best practices per service
					break
				}
				practices = append(practices, fmt.Sprintf("- %s: %s", bp.Title, bp.Description))
			}
			if len(practices) > 0 {
				contextParts = append(contextParts, fmt.Sprintf("BEST PRACTICES FOR %s:\n%s",
					strings.ToUpper(service), strings.Join(practices, "\n")))
			}
		}
	}

	// Get industry-specific compliance requirements
	industry := r.extractIndustryContext(inquiry)
	if industry != "" {
		complianceReqs, err := r.knowledgeBase.GetComplianceRequirements(ctx, industry)
		if err == nil && len(complianceReqs) > 0 {
			var requirements []string
			for i, req := range complianceReqs {
				if i >= 3 { // Limit to top 3 compliance requirements
					break
				}
				requirements = append(requirements, fmt.Sprintf("- %s (%s): %s",
					req.Framework, req.Severity, req.Description))
			}
			if len(requirements) > 0 {
				contextParts = append(contextParts, fmt.Sprintf("COMPLIANCE REQUIREMENTS FOR %s INDUSTRY:\n%s",
					strings.ToUpper(industry), strings.Join(requirements, "\n")))
			}
		}
	}

	if len(contextParts) > 0 {
		return "KNOWLEDGE BASE CONTEXT:\n" + strings.Join(contextParts, "\n\n"), nil
	}

	return "", nil
}

// buildMultiCloudContext builds context from multi-cloud analysis
func (r *reportGenerator) buildMultiCloudContext(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	// Get provider recommendation based on inquiry
	recommendation, err := r.multiCloudAnalyzer.GetProviderRecommendation(ctx, inquiry)
	if err != nil {
		return "", err
	}

	var contextParts []string

	// Add recommended provider information
	contextParts = append(contextParts, fmt.Sprintf("RECOMMENDED CLOUD PROVIDER: %s",
		strings.ToUpper(recommendation.RecommendedProvider)))

	if len(recommendation.Reasoning) > 0 {
		contextParts = append(contextParts, fmt.Sprintf("REASONING:\n%s",
			strings.Join(recommendation.Reasoning, "\n- ")))
	}

	// Add alternative options
	if len(recommendation.AlternativeOptions) > 0 {
		contextParts = append(contextParts, fmt.Sprintf("ALTERNATIVE OPTIONS: %s",
			strings.Join(recommendation.AlternativeOptions, ", ")))
	}

	// Add cost implications
	if recommendation.CostImplications != "" {
		contextParts = append(contextParts, fmt.Sprintf("COST IMPLICATIONS: %s",
			recommendation.CostImplications))
	}

	if len(contextParts) > 0 {
		return "MULTI-CLOUD ANALYSIS:\n" + strings.Join(contextParts, "\n\n"), nil
	}

	return "", nil
}

// buildRiskAssessmentContext builds context from risk assessment
func (r *reportGenerator) buildRiskAssessmentContext(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	// Create a basic proposed solution for risk assessment
	solution := &interfaces.ProposedSolution{
		ID:             uuid.New().String(),
		InquiryID:      inquiry.ID,
		CloudProviders: []string{"aws"}, // Default for assessment
		Services:       r.buildBasicServices(inquiry),
		Architecture:   r.buildBasicArchitecture(inquiry),
		EstimatedCost:  "TBD",
		Timeline:       "TBD",
	}

	// Perform risk assessment
	riskAssessment, err := r.riskAssessor.AssessRisks(ctx, inquiry, solution)
	if err != nil {
		return "", err
	}

	var contextParts []string

	// Add overall risk level
	contextParts = append(contextParts, fmt.Sprintf("OVERALL RISK LEVEL: %s",
		strings.ToUpper(riskAssessment.OverallRiskLevel)))

	// Add top risks by category
	if len(riskAssessment.TechnicalRisks) > 0 {
		var risks []string
		for i, risk := range riskAssessment.TechnicalRisks {
			if i >= 2 { // Limit to top 2 risks per category
				break
			}
			risks = append(risks, fmt.Sprintf("- %s (%s impact)", risk.Title, risk.Impact))
		}
		contextParts = append(contextParts, fmt.Sprintf("KEY TECHNICAL RISKS:\n%s",
			strings.Join(risks, "\n")))
	}

	if len(riskAssessment.SecurityRisks) > 0 {
		var risks []string
		for i, risk := range riskAssessment.SecurityRisks {
			if i >= 2 { // Limit to top 2 risks per category
				break
			}
			risks = append(risks, fmt.Sprintf("- %s (%s impact)", risk.Title, risk.Impact))
		}
		contextParts = append(contextParts, fmt.Sprintf("KEY SECURITY RISKS:\n%s",
			strings.Join(risks, "\n")))
	}

	// Add recommended actions
	if len(riskAssessment.RecommendedActions) > 0 {
		actions := riskAssessment.RecommendedActions
		if len(actions) > 3 {
			actions = actions[:3] // Limit to top 3 actions
		}
		contextParts = append(contextParts, fmt.Sprintf("RECOMMENDED RISK MITIGATION ACTIONS:\n- %s",
			strings.Join(actions, "\n- ")))
	}

	if len(contextParts) > 0 {
		return "RISK ASSESSMENT CONTEXT:\n" + strings.Join(contextParts, "\n\n"), nil
	}

	return "", nil
}

// buildDocumentationContext builds context from documentation library
func (r *reportGenerator) buildDocumentationContext(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	var contextParts []string

	// Get relevant documentation links for each service
	for _, service := range inquiry.Services {
		docLinks, err := r.docLibrary.GetDocumentationLinks(ctx, "", service)
		if err == nil && len(docLinks) > 0 {
			var links []string
			for i, link := range docLinks {
				if i >= 3 { // Limit to top 3 links per service
					break
				}
				if link.IsValid {
					links = append(links, fmt.Sprintf("- %s: %s (%s)",
						link.Title, link.URL, link.Provider))
				}
			}
			if len(links) > 0 {
				contextParts = append(contextParts, fmt.Sprintf("DOCUMENTATION FOR %s:\n%s",
					strings.ToUpper(service), strings.Join(links, "\n")))
			}
		}
	}

	// Get general best practices documentation
	bestPracticeLinks, err := r.docLibrary.GetLinksByType(ctx, "best-practice")
	if err == nil && len(bestPracticeLinks) > 0 {
		var links []string
		for i, link := range bestPracticeLinks {
			if i >= 3 { // Limit to top 3 best practice links
				break
			}
			if link.IsValid {
				links = append(links, fmt.Sprintf("- %s: %s (%s)",
					link.Title, link.URL, link.Provider))
			}
		}
		if len(links) > 0 {
			contextParts = append(contextParts, fmt.Sprintf("BEST PRACTICES DOCUMENTATION:\n%s",
				strings.Join(links, "\n")))
		}
	}

	if len(contextParts) > 0 {
		return "DOCUMENTATION REFERENCES:\n" + strings.Join(contextParts, "\n\n"), nil
	}

	return "", nil
}

// Helper methods for building basic structures for risk assessment

func (r *reportGenerator) buildBasicServices(inquiry *domain.Inquiry) []interfaces.CloudService {
	var services []interfaces.CloudService

	for _, serviceType := range inquiry.Services {
		service := interfaces.CloudService{
			Provider:      "aws", // Default provider
			ServiceName:   r.mapServiceTypeToAWSService(serviceType),
			ServiceType:   serviceType,
			Configuration: make(map[string]interface{}),
			Dependencies:  []string{},
			CriticalPath:  true,
		}
		services = append(services, service)
	}

	return services
}

func (r *reportGenerator) buildBasicArchitecture(inquiry *domain.Inquiry) *interfaces.Architecture {
	return &interfaces.Architecture{
		ID:   uuid.New().String(),
		Type: "cloud-native",
		Components: []interfaces.ArchitectureComponent{
			{
				Name:          "application-tier",
				Type:          "compute",
				Layer:         "application",
				Dependencies:  []string{},
				Criticality:   "high",
				Configuration: make(map[string]interface{}),
			},
		},
		NetworkTopology: interfaces.NetworkTopology{
			VPCConfiguration: make(map[string]interface{}),
			SubnetStrategy:   "multi-az",
			SecurityGroups:   []interfaces.SecurityGroup{},
			LoadBalancers:    []interfaces.LoadBalancer{},
			CDNConfiguration: make(map[string]interface{}),
		},
		DataStorage: []interfaces.DataStorageComponent{
			{
				Type:             "database",
				Provider:         "aws",
				ServiceName:      "RDS",
				DataType:         "application",
				SensitivityLevel: "medium",
				BackupStrategy:   "automated",
				Configuration:    make(map[string]interface{}),
			},
		},
		SecurityLayers:   []interfaces.SecurityLayer{},
		HighAvailability: true,
		DisasterRecovery: false,
	}
}

func (r *reportGenerator) mapServiceTypeToAWSService(serviceType string) string {
	serviceMap := map[string]string{
		"assessment":          "Well-Architected Review",
		"migration":           "Migration Hub",
		"optimization":        "Cost Explorer",
		"architecture-review": "Well-Architected Tool",
		"security":            "Security Hub",
		"compliance":          "Config",
		"monitoring":          "CloudWatch",
		"backup":              "Backup",
	}

	if awsService, exists := serviceMap[serviceType]; exists {
		return awsService
	}
	return "EC2" // Default fallback
}

func (r *reportGenerator) extractIndustryContext(inquiry *domain.Inquiry) string {
	// Extract industry hints from company name and message
	text := strings.ToLower(inquiry.Company + " " + inquiry.Message)

	industryKeywords := map[string][]string{
		"healthcare":    {"healthcare", "hospital", "medical", "patient", "hipaa", "health", "clinic"},
		"financial":     {"bank", "financial", "finance", "payment", "pci", "trading", "credit", "loan"},
		"retail":        {"retail", "ecommerce", "store", "shopping", "customer", "sales", "commerce"},
		"manufacturing": {"manufacturing", "factory", "production", "supply chain", "industrial", "plant"},
		"education":     {"education", "school", "university", "student", "learning", "academic", "campus"},
		"government":    {"government", "public", "federal", "state", "municipal", "agency", "civic"},
		"technology":    {"software", "tech", "saas", "platform", "development", "startup", "app"},
		"media":         {"media", "entertainment", "content", "streaming", "publishing", "broadcast"},
	}

	for industry, keywords := range industryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				return industry
			}
		}
	}

	return ""
}

// getEnhancedReportPDFOptions returns optimized PDF options for enhanced reports
func getEnhancedReportPDFOptions() *interfaces.PDFOptions {
	return &interfaces.PDFOptions{
		PageSize:     "A4",
		Orientation:  "Portrait",
		MarginTop:    "1in",
		MarginRight:  "0.75in",
		MarginBottom: "1in",
		MarginLeft:   "0.75in",
		Quality:      90,
		LoadTimeout:  30,
		CustomOptions: map[string]string{
			"enable-local-file-access": "true",
			"print-media-type":         "true",
		},
	}
}
