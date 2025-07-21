package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/google/uuid"
)

// reportGenerator implements report generation using Bedrock
type reportGenerator struct {
	bedrockService  interfaces.BedrockService
	templateService interfaces.TemplateService
	pdfService      interfaces.PDFService
}

// NewReportGenerator creates a new report generator instance
func NewReportGenerator(bedrockService interfaces.BedrockService, templateService interfaces.TemplateService, pdfService interfaces.PDFService) interfaces.ReportService {
	return &reportGenerator{
		bedrockService:  bedrockService,
		templateService: templateService,
		pdfService:      pdfService,
	}
}

// GenerateReport generates a report for the given inquiry using Bedrock
func (r *reportGenerator) GenerateReport(ctx context.Context, inquiry *domain.Inquiry) (*domain.Report, error) {
	// Build the prompt based on inquiry details
	prompt := r.buildPrompt(inquiry)
	
	// Set Bedrock options
	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   2000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	// Generate content using Bedrock
	response, err := r.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		// Return error but don't fail the entire process
		return nil, fmt.Errorf("failed to generate report with Bedrock: %w", err)
	}

	// Create the report
	report := &domain.Report{
		ID:          uuid.New().String(),
		InquiryID:   inquiry.ID,
		Type:        r.getReportType(inquiry.Services),
		Title:       r.generateTitle(inquiry),
		Content:     response.Content,
		Status:      domain.ReportStatusDraft,
		GeneratedBy: "bedrock-ai",
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

// formatContentForHTML converts plain text content to HTML
func (r *reportGenerator) formatContentForHTML(content string) string {
	if content == "" {
		return "<p>No content available.</p>"
	}
	
	// Basic HTML formatting
	// Split into paragraphs
	paragraphs := strings.Split(content, "\n\n")
	var htmlParagraphs []string
	
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		
		// Convert line breaks within paragraphs
		p = strings.ReplaceAll(p, "\n", "<br>")
		
		// Convert **bold** to <strong>
		p = strings.ReplaceAll(p, "**", "<strong>")
		
		// Simple header detection
		if strings.HasSuffix(p, ":") && len(p) < 100 {
			htmlParagraphs = append(htmlParagraphs, fmt.Sprintf("<h3>%s</h3>", p))
		} else {
			htmlParagraphs = append(htmlParagraphs, fmt.Sprintf("<p>%s</p>", p))
		}
	}
	
	return strings.Join(htmlParagraphs, "\n")
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
	options := getReportPDFOptions()
	
	// Generate PDF from HTML
	pdfBytes, err := r.pdfService.GeneratePDF(ctx, htmlContent, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	
	return pdfBytes, nil
}

