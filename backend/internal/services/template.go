package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// templateService implements the TemplateService interface
type templateService struct {
	templatesDir string
	templates    map[string]*template.Template
	logger       *logrus.Logger
}

// NewTemplateService creates a new template service instance
func NewTemplateService(templatesDir string, logger *logrus.Logger) interfaces.TemplateService {
	service := &templateService{
		templatesDir: templatesDir,
		templates:    make(map[string]*template.Template),
		logger:       logger,
	}
	
	// Load templates on initialization
	if err := service.loadTemplates(); err != nil {
		logger.WithError(err).Error("Failed to load email templates")
	}
	
	return service
}

// loadTemplates loads all email and report templates from the templates directory
func (t *templateService) loadTemplates() error {
	// Load email templates
	t.loadEmailTemplates()
	
	// Load report templates
	t.loadReportTemplates()
	
	return nil
}

// loadEmailTemplates loads email templates
func (t *templateService) loadEmailTemplates() {
	// Load customer confirmation template
	customerConfirmationPath := filepath.Join(t.templatesDir, "email", "customer_confirmation.html")
	customerTemplate, err := template.ParseFiles(customerConfirmationPath)
	if err != nil {
		t.logger.WithError(err).WithField("template", "customer_confirmation").Warn("Failed to load customer confirmation template")
	} else {
		t.templates["customer_confirmation"] = customerTemplate
		t.logger.WithField("template", "customer_confirmation").Info("Loaded customer confirmation template")
	}
	
	// Load consultant notification template
	consultantNotificationPath := filepath.Join(t.templatesDir, "email", "consultant_notification.html")
	consultantTemplate, err := template.ParseFiles(consultantNotificationPath)
	if err != nil {
		t.logger.WithError(err).WithField("template", "consultant_notification").Warn("Failed to load consultant notification template")
	} else {
		t.templates["consultant_notification"] = consultantTemplate
		t.logger.WithField("template", "consultant_notification").Info("Loaded consultant notification template")
	}
}

// loadReportTemplates loads report templates
func (t *templateService) loadReportTemplates() {
	reportTemplates := []string{"assessment", "migration", "optimization", "architecture"}
	
	for _, templateName := range reportTemplates {
		templatePath := filepath.Join(t.templatesDir, "reports", templateName+".html")
		reportTemplate, err := template.ParseFiles(templatePath)
		if err != nil {
			t.logger.WithError(err).WithField("template", templateName).Warn("Failed to load report template")
		} else {
			t.templates["report_"+templateName] = reportTemplate
			t.logger.WithField("template", "report_"+templateName).Info("Loaded report template")
		}
	}
}

// RenderEmailTemplate renders an email template with the provided data
func (t *templateService) RenderEmailTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	tmpl, exists := t.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.logger.WithError(err).WithFields(logrus.Fields{
			"template": templateName,
		}).Error("Failed to render email template")
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}
	
	return buf.String(), nil
}

// RenderReportTemplate renders a report template with the provided data
func (t *templateService) RenderReportTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	// Look for report template first
	reportTemplateName := "report_" + templateName
	tmpl, exists := t.templates[reportTemplateName]
	if !exists {
		return "", fmt.Errorf("report template %s not found", templateName)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.logger.WithError(err).WithFields(logrus.Fields{
			"template": reportTemplateName,
		}).Error("Failed to render report template")
		return "", fmt.Errorf("failed to render report template %s: %w", templateName, err)
	}
	
	return buf.String(), nil
}

// LoadTemplate loads a specific template by name
func (t *templateService) LoadTemplate(templateName string) (*template.Template, error) {
	tmpl, exists := t.templates[templateName]
	if !exists {
		return nil, fmt.Errorf("template %s not found", templateName)
	}
	return tmpl, nil
}

// ValidateTemplate validates a template content string
func (t *templateService) ValidateTemplate(templateContent string) error {
	_, err := template.New("validation").Parse(templateContent)
	if err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}
	return nil
}

// CustomerConfirmationData represents the data structure for customer confirmation emails
type CustomerConfirmationData struct {
	Name     string
	Company  string
	Services string
	ID       string
}

// ConsultantNotificationData represents the data structure for consultant notification emails
type ConsultantNotificationData struct {
	Name           string
	Email          string
	Company        string
	Phone          string
	Services       string
	Message        string
	ID             string
	IsHighPriority bool
	Priority       string
	Report         *ReportData
}

// ReportData represents report information for email templates
type ReportData struct {
	ID          string
	HTMLContent template.HTML
	Content     string
}

// HTMLReportTemplateData represents the data structure for HTML report templates
type HTMLReportTemplateData struct {
	ID               string
	Title            string
	ClientName       string
	ClientEmail      string
	ClientCompany    string
	ClientPhone      string
	Services         string
	GeneratedDate    string
	IsHighPriority   bool
	FormattedContent template.HTML
}

// PrepareCustomerConfirmationData prepares data for customer confirmation email template
func (t *templateService) PrepareCustomerConfirmationData(inquiry *domain.Inquiry) *CustomerConfirmationData {
	return &CustomerConfirmationData{
		Name:     inquiry.Name,
		Company:  inquiry.Company,
		Services: strings.Join(inquiry.Services, ", "),
		ID:       inquiry.ID,
	}
}

// PrepareConsultantNotificationData prepares data for consultant notification email template
func (t *templateService) PrepareConsultantNotificationData(inquiry *domain.Inquiry, report *domain.Report, isHighPriority bool) interface{} {
	data := &ConsultantNotificationData{
		Name:           inquiry.Name,
		Email:          inquiry.Email,
		Company:        inquiry.Company,
		Phone:          inquiry.Phone,
		Services:       strings.Join(inquiry.Services, ", "),
		Message:        inquiry.Message,
		ID:             inquiry.ID,
		IsHighPriority: isHighPriority,
		Priority:       func() string {
			if isHighPriority {
				return "HIGH"
			}
			return "NORMAL"
		}(),
	}
	
	if report != nil {
		data.Report = &ReportData{
			ID:          report.ID,
			Content:     report.Content,
			HTMLContent: template.HTML(t.convertMarkdownToHTML(report.Content)),
		}
	}
	
	return data
}

// convertMarkdownToHTML converts Markdown text to HTML (simplified version)
func (t *templateService) convertMarkdownToHTML(markdown string) string {
	if markdown == "" {
		return "<p>No content available.</p>"
	}
	
	html := markdown
	
	// First, handle bold text properly
	html = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(html, "<strong>$1</strong>")
	
	// Handle horizontal rules
	html = strings.ReplaceAll(html, "---", "<hr>")
	
	// Split into sections by double line breaks
	sections := strings.Split(html, "\n\n")
	var formattedSections []string
	
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		
		// Check if this is a header (starts with ### or is all caps)
		if strings.HasPrefix(section, "### ") {
			headerText := strings.TrimPrefix(section, "### ")
			headerText = strings.TrimSpace(headerText)
			// Remove markdown formatting from headers
			headerText = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(headerText, "$1")
			formattedSections = append(formattedSections, fmt.Sprintf("<h3>%s</h3>", headerText))
		} else if strings.HasPrefix(section, "## ") {
			headerText := strings.TrimPrefix(section, "## ")
			headerText = strings.TrimSpace(headerText)
			headerText = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(headerText, "$1")
			formattedSections = append(formattedSections, fmt.Sprintf("<h2>%s</h2>", headerText))
		} else if strings.HasPrefix(section, "# ") {
			headerText := strings.TrimPrefix(section, "# ")
			headerText = strings.TrimSpace(headerText)
			headerText = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(headerText, "$1")
			formattedSections = append(formattedSections, fmt.Sprintf("<h1>%s</h1>", headerText))
		} else {
			// Regular content - process line by line
			lines := strings.Split(section, "\n")
			var processedLines []string
			inList := false
			
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				
				// Handle bullet points
				if strings.HasPrefix(line, "- ") {
					if !inList {
						processedLines = append(processedLines, "<ul>")
						inList = true
					}
					listItem := strings.TrimPrefix(line, "- ")
					processedLines = append(processedLines, fmt.Sprintf("  <li>%s</li>", listItem))
				} else if regexp.MustCompile(`^\d+\.\s`).MatchString(line) {
					// Handle numbered lists
					if !inList {
						processedLines = append(processedLines, "<ol>")
						inList = true
					}
					listItem := regexp.MustCompile(`^\d+\.\s`).ReplaceAllString(line, "")
					processedLines = append(processedLines, fmt.Sprintf("  <li>%s</li>", listItem))
				} else {
					if inList {
						// Close the list
						if strings.Contains(strings.Join(processedLines[len(processedLines)-3:], ""), "<ul>") {
							processedLines = append(processedLines, "</ul>")
						} else {
							processedLines = append(processedLines, "</ol>")
						}
						inList = false
					}
					processedLines = append(processedLines, line)
				}
			}
			
			if inList {
				// Close any open list
				if strings.Contains(strings.Join(processedLines, ""), "<ul>") {
					processedLines = append(processedLines, "</ul>")
				} else {
					processedLines = append(processedLines, "</ol>")
				}
			}
			
			// Join lines and wrap in paragraph if needed
			content := strings.Join(processedLines, "\n")
			if !strings.Contains(content, "<ul>") && !strings.Contains(content, "<ol>") && !strings.Contains(content, "<h") {
				content = fmt.Sprintf("<p>%s</p>", content)
			}
			
			formattedSections = append(formattedSections, content)
		}
	}
	
	return strings.Join(formattedSections, "\n\n")
}

// GetAvailableTemplates returns a list of available template names
func (t *templateService) GetAvailableTemplates() []string {
	templates := make([]string, 0, len(t.templates))
	for name := range t.templates {
		templates = append(templates, name)
	}
	return templates
}

// ReloadTemplates reloads all templates from disk
func (t *templateService) ReloadTemplates() error {
	t.templates = make(map[string]*template.Template)
	return t.loadTemplates()
}

// PrepareReportTemplateData prepares data for report template rendering
func (t *templateService) PrepareReportTemplateData(inquiry *domain.Inquiry, report *domain.Report) interface{} {
	// Detect high priority based on report content
	isHighPriority := t.detectHighPriority(report.Content)
	
	return &HTMLReportTemplateData{
		ID:               report.ID,
		Title:            report.Title,
		ClientName:       inquiry.Name,
		ClientEmail:      inquiry.Email,
		ClientCompany:    t.getCompanyOrDefault(inquiry.Company),
		ClientPhone:      t.getPhoneOrDefault(inquiry.Phone),
		Services:         strings.Join(inquiry.Services, ", "),
		GeneratedDate:    report.CreatedAt.Format("January 2, 2006"),
		IsHighPriority:   isHighPriority,
		FormattedContent: template.HTML(t.formatReportContent(report.Content)),
	}
}

// detectHighPriority analyzes report content for priority indicators
func (t *templateService) detectHighPriority(content string) bool {
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

// formatReportContent converts plain text report content to structured HTML
func (t *templateService) formatReportContent(content string) string {
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
		formattedSection := t.formatSection(section)
		htmlSections = append(htmlSections, formattedSection)
	}
	
	return strings.Join(htmlSections, "\n\n")
}

// formatSection formats a single section of content
func (t *templateService) formatSection(section string) string {
	// Check if this is a main header (numbered or all caps)
	if t.isMainHeader(section) {
		return t.formatMainHeader(section)
	}
	
	// Check if this is a sub-header
	if t.isSubHeader(section) {
		return t.formatSubHeader(section)
	}
	
	// Check if this contains a list
	if t.containsList(section) {
		return t.formatListSection(section)
	}
	
	// Format as regular paragraph content
	return t.formatParagraphSection(section)
}

// isMainHeader checks if a section is a main header
func (t *templateService) isMainHeader(text string) bool {
	text = strings.TrimSpace(text)
	
	// Don't treat multi-line content as headers
	if strings.Count(text, "\n") > 1 {
		return false
	}
	
	// Don't treat long sentences as headers
	if len(text) > 150 {
		return false
	}
	
	// Check for numbered headers (1., 2., etc.)
	if regexp.MustCompile(`^\d+\.\s+[A-Z]`).MatchString(text) {
		return true
	}
	
	// Check for main section headers (must be exact matches or start with the header)
	mainHeaders := []string{
		"EXECUTIVE SUMMARY", "CURRENT STATE ASSESSMENT", "RECOMMENDATIONS",
		"NEXT STEPS", "URGENCY ASSESSMENT", "CONTACT INFORMATION",
		"PRIORITY LEVEL", "MEETING SCHEDULING",
	}
	
	textUpper := strings.ToUpper(text)
	for _, header := range mainHeaders {
		// Exact match or starts with header followed by colon or space
		if textUpper == header || 
		   strings.HasPrefix(textUpper, header+":") ||
		   (strings.HasPrefix(textUpper, header+" ") && len(text) < 100) {
			return true
		}
	}
	
	// Check if it's all caps and short (likely a header)
	if strings.ToUpper(text) == text && len(text) < 80 && !strings.Contains(text, ".") {
		return true
	}
	
	return false
}

// isSubHeader checks if a section is a sub-header
func (t *templateService) isSubHeader(text string) bool {
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
func (t *templateService) containsList(text string) bool {
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
func (t *templateService) formatMainHeader(text string) string {
	text = strings.TrimSpace(text)
	
	// Remove numbered prefixes for cleaner headers
	text = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(text, "")
	
	// Convert to title case if all caps
	if strings.ToUpper(text) == text {
		text = t.toTitleCase(text)
	}
	
	// Remove trailing colons
	text = strings.TrimSuffix(text, ":")
	
	return fmt.Sprintf("<h2 class=\"section-header\">%s</h2>", text)
}

// formatSubHeader formats a sub-header
func (t *templateService) formatSubHeader(text string) string {
	text = strings.TrimSpace(text)
	
	// Remove bold markdown
	text = strings.Trim(text, "*")
	
	// Remove trailing colons
	text = strings.TrimSuffix(text, ":")
	
	return fmt.Sprintf("<h3 class=\"subsection-header\">%s</h3>", text)
}

// formatListSection formats a section containing lists
func (t *templateService) formatListSection(text string) string {
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
				result = append(result, t.formatParagraphContent(strings.Join(currentParagraph, " ")))
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
			result = append(result, fmt.Sprintf("  <li>%s</li>", t.formatInlineText(item)))
			
		} else if regexp.MustCompile(`^\d+\.\s`).MatchString(line) {
			// Close any open paragraph
			if len(currentParagraph) > 0 {
				result = append(result, t.formatParagraphContent(strings.Join(currentParagraph, " ")))
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
			result = append(result, fmt.Sprintf("  <li>%s</li>", t.formatInlineText(item)))
			
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
		result = append(result, t.formatParagraphContent(strings.Join(currentParagraph, " ")))
	}
	
	return strings.Join(result, "\n")
}

// formatParagraphSection formats a regular paragraph section
func (t *templateService) formatParagraphSection(text string) string {
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
	return t.formatParagraphContent(paragraph)
}

// formatParagraphContent formats a single paragraph with inline formatting
func (t *templateService) formatParagraphContent(text string) string {
	if text == "" {
		return ""
	}
	
	formatted := t.formatInlineText(text)
	return fmt.Sprintf("<p class=\"report-paragraph\">%s</p>", formatted)
}

// formatInlineText applies inline formatting (bold, italic, etc.)
func (t *templateService) formatInlineText(text string) string {
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
func (t *templateService) toTitleCase(text string) string {
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

// isHeader determines if a section is a header
func (t *templateService) isHeader(text string) bool {
	text = strings.TrimSpace(text)
	
	// Don't treat multi-line content as headers
	if strings.Contains(text, "\n") && len(strings.Split(text, "\n")) > 2 {
		return false
	}
	
	// Check for numbered headers (1., 2., etc.)
	if matched, _ := regexp.MatchString(`^\d+\.`, text); matched {
		return true
	}
	
	// Check for all caps headers (but not too long)
	if strings.ToUpper(text) == text && len(text) < 100 && !strings.Contains(text, "\n") {
		return true
	}
	
	// Check for headers with specific keywords (single line only)
	headerKeywords := []string{
		"EXECUTIVE SUMMARY", "CURRENT STATE", "RECOMMENDATIONS", "NEXT STEPS",
		"ASSESSMENT", "MIGRATION", "OPTIMIZATION", "ARCHITECTURE",
		"PRIORITY LEVEL", "URGENCY ASSESSMENT", "CONTACT INFORMATION",
	}
	
	textUpper := strings.ToUpper(text)
	for _, keyword := range headerKeywords {
		if textUpper == keyword || strings.HasPrefix(textUpper, keyword+":") {
			return true
		}
	}
	
	return false
}

// getHeaderLevel determines the header level (1-4)
func (t *templateService) getHeaderLevel(text string) int {
	text = strings.TrimSpace(text)
	
	// Main sections (h1)
	mainSections := []string{
		"EXECUTIVE SUMMARY", "CURRENT STATE", "RECOMMENDATIONS", "NEXT STEPS",
		"URGENCY ASSESSMENT", "CONTACT INFORMATION",
	}
	
	textUpper := strings.ToUpper(text)
	for _, section := range mainSections {
		if strings.Contains(textUpper, section) {
			return 1
		}
	}
	
	// Numbered sections (h2)
	if matched, _ := regexp.MatchString(`^\d+\.`, text); matched {
		return 2
	}
	
	// Sub-sections (h3)
	if strings.Contains(textUpper, "PRIORITY") || strings.Contains(textUpper, "MEETING") {
		return 3
	}
	
	// Default to h2
	return 2
}

// cleanHeaderText removes formatting artifacts from header text
func (t *templateService) cleanHeaderText(text string) string {
	// Remove numbered prefixes
	text = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(text, "")
	
	// Convert to title case if all caps
	if strings.ToUpper(text) == text {
		words := strings.Fields(text)
		for i, word := range words {
			if len(word) > 0 {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
		text = strings.Join(words, " ")
	}
	
	return strings.TrimSpace(text)
}

// formatParagraph formats a paragraph with proper HTML
func (t *templateService) formatParagraph(text string) string {
	// Convert **bold** to <strong> first
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "<strong>$1</strong>")
	
	// Convert *italic* to <em>
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "<em>$1</em>")
	
	// Split into lines for processing
	lines := strings.Split(text, "\n")
	var formattedLines []string
	inList := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check for bullet points
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") {
			if !inList {
				formattedLines = append(formattedLines, "<ul>")
				inList = true
			}
			listItem := strings.TrimPrefix(line, "- ")
			listItem = strings.TrimPrefix(listItem, "• ")
			formattedLines = append(formattedLines, fmt.Sprintf("  <li>%s</li>", listItem))
		} else {
			if inList {
				formattedLines = append(formattedLines, "</ul>")
				inList = false
			}
			formattedLines = append(formattedLines, line)
		}
	}
	
	if inList {
		formattedLines = append(formattedLines, "</ul>")
	}
	
	result := strings.Join(formattedLines, "<br>\n")
	
	// Wrap in paragraph if it doesn't contain block elements
	if !strings.Contains(result, "<ul>") && !strings.Contains(result, "<h") {
		result = fmt.Sprintf("<p>%s</p>", result)
	}
	
	return result
}

// getCompanyOrDefault returns the company name or a default value
func (t *templateService) getCompanyOrDefault(company string) string {
	if company != "" {
		return company
	}
	return "Client Organization"
}

// getPhoneOrDefault returns the phone number or a default value
func (t *templateService) getPhoneOrDefault(phone string) string {
	if phone != "" {
		return phone
	}
	return "Not provided"
}