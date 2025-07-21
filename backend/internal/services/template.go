package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
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

// loadTemplates loads all email templates from the templates directory
func (t *templateService) loadTemplates() error {
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
	
	return nil
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
	// For now, we'll use the same mechanism as email templates
	// In the future, this could be extended for specific report templates
	return t.RenderEmailTemplate(ctx, templateName, data)
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
	HTMLContent string
	Content     string
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
func (t *templateService) PrepareConsultantNotificationData(inquiry *domain.Inquiry, report *domain.Report, isHighPriority bool) *ConsultantNotificationData {
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
			HTMLContent: t.convertMarkdownToHTML(report.Content),
		}
	}
	
	return data
}

// convertMarkdownToHTML converts Markdown text to HTML (simplified version)
func (t *templateService) convertMarkdownToHTML(markdown string) string {
	// This is a simplified conversion - in a real implementation,
	// you might want to use a proper markdown library like blackfriday
	html := markdown
	
	// Convert headers
	html = strings.ReplaceAll(html, "# ", "<h1>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "### ", "<h3>")
	html = strings.ReplaceAll(html, "#### ", "<h4>")
	
	// Convert bold text
	html = strings.ReplaceAll(html, "**", "<strong>")
	html = strings.ReplaceAll(html, "</strong>", "</strong>")
	
	// Convert line breaks to paragraphs
	paragraphs := strings.Split(html, "\n\n")
	for i, p := range paragraphs {
		if strings.TrimSpace(p) != "" && !strings.HasPrefix(p, "<h") {
			paragraphs[i] = "<p>" + strings.ReplaceAll(p, "\n", "<br>") + "</p>"
		}
	}
	
	return strings.Join(paragraphs, "\n")
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