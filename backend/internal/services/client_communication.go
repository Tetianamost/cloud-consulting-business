package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/google/uuid"
)

// ClientCommunicationService implements advanced client communication tools
type ClientCommunicationService struct {
	bedrockService  interfaces.BedrockService
	promptArchitect interfaces.PromptArchitect
}

// NewClientCommunicationService creates a new client communication service
func NewClientCommunicationService(
	bedrockService interfaces.BedrockService,
	promptArchitect interfaces.PromptArchitect,
) *ClientCommunicationService {
	return &ClientCommunicationService{
		bedrockService:  bedrockService,
		promptArchitect: promptArchitect,
	}
}

// GenerateTechnicalExplanation generates a technical explanation for business stakeholders
func (s *ClientCommunicationService) GenerateTechnicalExplanation(
	ctx context.Context,
	req *interfaces.TechnicalExplanationRequest,
) (*interfaces.TechnicalExplanation, error) {
	// Build the prompt for technical explanation
	prompt := s.buildTechnicalExplanationPrompt(req)

	// Generate content using Bedrock
	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate technical explanation: %w", err)
	}

	// Parse and structure the response
	explanation := s.parseTechnicalExplanation(response.Content, req)

	return explanation, nil
}

// GeneratePresentationSlides generates presentation slides with technical diagrams
func (s *ClientCommunicationService) GeneratePresentationSlides(
	ctx context.Context,
	req *interfaces.PresentationRequest,
) (*interfaces.PresentationSlides, error) {
	// Build the prompt for presentation generation
	prompt := s.buildPresentationPrompt(req)

	// Generate content using Bedrock
	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   6000,
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presentation slides: %w", err)
	}

	// Parse and structure the response
	slides := s.parsePresentationSlides(response.Content, req)

	return slides, nil
}

// GenerateEmailTemplate generates email templates for client communications
func (s *ClientCommunicationService) GenerateEmailTemplate(
	ctx context.Context,
	req *interfaces.EmailTemplateRequest,
) (*interfaces.EmailTemplate, error) {
	// Build the prompt for email template generation
	prompt := s.buildEmailTemplatePrompt(req)

	// Generate content using Bedrock
	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   2000,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate email template: %w", err)
	}

	// Parse and structure the response
	template := s.parseEmailTemplate(response.Content, req)

	return template, nil
}

// GenerateStatusReport generates status reports for ongoing engagements
func (s *ClientCommunicationService) GenerateStatusReport(
	ctx context.Context,
	req *interfaces.StatusReportRequest,
) (*interfaces.StatusReport, error) {
	// Build the prompt for status report generation
	prompt := s.buildStatusReportPrompt(req)

	// Generate content using Bedrock
	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   5000,
		Temperature: 0.2,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate status report: %w", err)
	}

	// Parse and structure the response
	report := s.parseStatusReport(response.Content, req)

	return report, nil
}

// buildTechnicalExplanationPrompt builds a prompt for technical explanation generation
func (s *ClientCommunicationService) buildTechnicalExplanationPrompt(req *interfaces.TechnicalExplanationRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are an expert AWS cloud consultant helping to explain complex technical concepts to business stakeholders.\n\n")

	prompt.WriteString(fmt.Sprintf("TASK: Create a comprehensive technical explanation for '%s'\n", req.TechnicalConcept))
	prompt.WriteString(fmt.Sprintf("AUDIENCE: %s stakeholders\n", req.AudienceType))
	prompt.WriteString(fmt.Sprintf("COMPLEXITY LEVEL: %s\n", req.ComplexityLevel))

	if req.IndustryContext != "" {
		prompt.WriteString(fmt.Sprintf("INDUSTRY CONTEXT: %s\n", req.IndustryContext))
	}

	if req.BusinessContext != "" {
		prompt.WriteString(fmt.Sprintf("BUSINESS CONTEXT: %s\n", req.BusinessContext))
	}

	if req.UseCase != "" {
		prompt.WriteString(fmt.Sprintf("USE CASE: %s\n", req.UseCase))
	}

	if len(req.KeyTerms) > 0 {
		prompt.WriteString(fmt.Sprintf("KEY TERMS TO INCLUDE: %s\n", strings.Join(req.KeyTerms, ", ")))
	}

	if len(req.Constraints) > 0 {
		prompt.WriteString(fmt.Sprintf("CONSTRAINTS: %s\n", strings.Join(req.Constraints, ", ")))
	}

	prompt.WriteString("\nPLEASE PROVIDE:\n")
	prompt.WriteString("1. TITLE: A clear, engaging title\n")
	prompt.WriteString("2. EXECUTIVE_SUMMARY: 2-3 sentences for executives\n")
	prompt.WriteString("3. BUSINESS_VALUE: Why this matters to the business\n")
	prompt.WriteString("4. TECHNICAL_OVERVIEW: Technical explanation appropriate for the audience\n")
	prompt.WriteString("5. KEY_BENEFITS: 3-5 specific benefits\n")
	prompt.WriteString("6. CONSIDERATIONS: Important factors to consider\n")
	prompt.WriteString("7. NEXT_STEPS: Recommended actions\n")
	prompt.WriteString("8. GLOSSARY: Key technical terms and definitions\n")
	prompt.WriteString("9. REFERENCES: Relevant AWS documentation links\n")

	prompt.WriteString("\nFormat your response with clear section headers. Be specific and actionable.")

	return prompt.String()
}

// buildPresentationPrompt builds a prompt for presentation slide generation
func (s *ClientCommunicationService) buildPresentationPrompt(req *interfaces.PresentationRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are an expert AWS cloud consultant creating a professional presentation.\n\n")

	prompt.WriteString(fmt.Sprintf("TOPIC: %s\n", req.Topic))
	prompt.WriteString(fmt.Sprintf("AUDIENCE: %s\n", req.AudienceType))
	prompt.WriteString(fmt.Sprintf("DURATION: %d minutes\n", req.Duration))

	if req.SlideCount > 0 {
		prompt.WriteString(fmt.Sprintf("TARGET SLIDES: %d\n", req.SlideCount))
	}

	if req.IndustryContext != "" {
		prompt.WriteString(fmt.Sprintf("INDUSTRY: %s\n", req.IndustryContext))
	}

	if len(req.BusinessGoals) > 0 {
		prompt.WriteString(fmt.Sprintf("BUSINESS GOALS: %s\n", strings.Join(req.BusinessGoals, ", ")))
	}

	if len(req.TechnicalScope) > 0 {
		prompt.WriteString(fmt.Sprintf("TECHNICAL SCOPE: %s\n", strings.Join(req.TechnicalScope, ", ")))
	}

	if len(req.ComplianceReqs) > 0 {
		prompt.WriteString(fmt.Sprintf("COMPLIANCE REQUIREMENTS: %s\n", strings.Join(req.ComplianceReqs, ", ")))
	}

	prompt.WriteString("\nREQUIREMENTS:\n")
	if req.IncludeDiagrams {
		prompt.WriteString("- Include technical diagrams (provide Mermaid syntax)\n")
	}
	if req.IncludeCostInfo {
		prompt.WriteString("- Include cost analysis and optimization information\n")
	}

	prompt.WriteString("\nPLEASE PROVIDE:\n")
	prompt.WriteString("1. PRESENTATION_TITLE: Main title and subtitle\n")
	prompt.WriteString("2. SLIDE_OUTLINE: List of slides with titles and key points\n")
	prompt.WriteString("3. DETAILED_SLIDES: Full content for each slide\n")
	prompt.WriteString("4. SPEAKER_NOTES: Key talking points for each slide\n")
	prompt.WriteString("5. DIAGRAMS: Mermaid diagram syntax where applicable\n")
	prompt.WriteString("6. APPENDIX: Additional technical details\n")

	prompt.WriteString("\nFormat slides for professional business presentation. Include estimated timing for each slide.")

	return prompt.String()
}

// buildEmailTemplatePrompt builds a prompt for email template generation
func (s *ClientCommunicationService) buildEmailTemplatePrompt(req *interfaces.EmailTemplateRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are an expert AWS cloud consultant creating professional email templates.\n\n")

	prompt.WriteString(fmt.Sprintf("PURPOSE: %s\n", req.Purpose))
	prompt.WriteString(fmt.Sprintf("RECIPIENT: %s\n", req.RecipientType))
	prompt.WriteString(fmt.Sprintf("TONE: %s\n", req.ToneStyle))
	prompt.WriteString(fmt.Sprintf("URGENCY: %s\n", req.Urgency))
	prompt.WriteString(fmt.Sprintf("CONTEXT: %s\n", req.Context))

	if req.ProjectPhase != "" {
		prompt.WriteString(fmt.Sprintf("PROJECT PHASE: %s\n", req.ProjectPhase))
	}

	if len(req.KeyPoints) > 0 {
		prompt.WriteString(fmt.Sprintf("KEY POINTS TO INCLUDE: %s\n", strings.Join(req.KeyPoints, ", ")))
	}

	if req.CallToAction != "" {
		prompt.WriteString(fmt.Sprintf("CALL TO ACTION: %s\n", req.CallToAction))
	}

	if req.Deadline != nil {
		prompt.WriteString(fmt.Sprintf("DEADLINE: %s\n", req.Deadline.Format("January 2, 2006")))
	}

	prompt.WriteString("\nREQUIREMENTS:\n")
	if req.IncludeAttachment {
		prompt.WriteString("- Mention attachment in email\n")
	}
	if req.ActionRequired {
		prompt.WriteString("- Clear action items required from recipient\n")
	}

	prompt.WriteString("\nPLEASE PROVIDE:\n")
	prompt.WriteString("1. SUBJECT: Professional subject line\n")
	prompt.WriteString("2. EMAIL_BODY: Complete email content\n")
	prompt.WriteString("3. ALTERNATIVE_SUBJECTS: 2-3 alternative subject lines\n")
	prompt.WriteString("4. FOLLOW_UP_ACTIONS: Suggested follow-up steps\n")
	prompt.WriteString("5. TIMING_SUGGESTION: Best time to send this email\n")

	prompt.WriteString("\nCreate a professional, clear, and actionable email appropriate for cloud consulting business.")

	return prompt.String()
}

// buildStatusReportPrompt builds a prompt for status report generation
func (s *ClientCommunicationService) buildStatusReportPrompt(req *interfaces.StatusReportRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are an expert AWS cloud consultant creating a professional project status report.\n\n")

	prompt.WriteString(fmt.Sprintf("PROJECT ID: %s\n", req.ProjectID))
	prompt.WriteString(fmt.Sprintf("REPORTING PERIOD: %s\n", req.ReportingPeriod))
	prompt.WriteString(fmt.Sprintf("AUDIENCE: %s\n", req.AudienceType))

	if req.ProjectPhase != "" {
		prompt.WriteString(fmt.Sprintf("PROJECT PHASE: %s\n", req.ProjectPhase))
	}

	// Include milestones if provided
	if len(req.Milestones) > 0 {
		prompt.WriteString("\nMILESTONES:\n")
		for _, milestone := range req.Milestones {
			prompt.WriteString(fmt.Sprintf("- %s: %s (Due: %s, Status: %s, Progress: %d%%)\n",
				milestone.Name, milestone.Description, milestone.DueDate.Format("Jan 2"), milestone.Status, milestone.Progress))
		}
	}

	// Include issues if provided
	if len(req.Issues) > 0 {
		prompt.WriteString("\nISSUES:\n")
		for _, issue := range req.Issues {
			prompt.WriteString(fmt.Sprintf("- %s: %s (Severity: %s, Status: %s)\n",
				issue.Title, issue.Description, issue.Severity, issue.Status))
		}
	}

	// Include achievements if provided
	if len(req.Achievements) > 0 {
		prompt.WriteString("\nACHIEVEMENTS:\n")
		for _, achievement := range req.Achievements {
			prompt.WriteString(fmt.Sprintf("- %s: %s (Impact: %s)\n",
				achievement.Title, achievement.Description, achievement.Impact))
		}
	}

	// Include budget information if provided
	if req.Budget != nil {
		prompt.WriteString(fmt.Sprintf("\nBUDGET STATUS:\n"))
		prompt.WriteString(fmt.Sprintf("- Total Budget: $%.2f\n", req.Budget.TotalBudget))
		prompt.WriteString(fmt.Sprintf("- Spent: $%.2f\n", req.Budget.SpentAmount))
		prompt.WriteString(fmt.Sprintf("- Remaining: $%.2f\n", req.Budget.RemainingAmount))
		prompt.WriteString(fmt.Sprintf("- Status: %s\n", req.Budget.Status))
	}

	// Include timeline information if provided
	if req.Timeline != nil {
		prompt.WriteString(fmt.Sprintf("\nTIMELINE STATUS:\n"))
		prompt.WriteString(fmt.Sprintf("- Progress: %d%%\n", req.Timeline.ProgressPercent))
		prompt.WriteString(fmt.Sprintf("- Days Remaining: %d\n", req.Timeline.DaysRemaining))
		prompt.WriteString(fmt.Sprintf("- Status: %s\n", req.Timeline.Status))
	}

	prompt.WriteString("\nREQUIREMENTS:\n")
	if req.IncludeMetrics {
		prompt.WriteString("- Include relevant project metrics\n")
	}
	if req.IncludeRisks {
		prompt.WriteString("- Include risk assessment\n")
	}
	if req.IncludeNextSteps {
		prompt.WriteString("- Include clear next steps\n")
	}

	prompt.WriteString("\nPLEASE PROVIDE:\n")
	prompt.WriteString("1. REPORT_TITLE: Professional report title\n")
	prompt.WriteString("2. EXECUTIVE_SUMMARY: High-level status summary\n")
	prompt.WriteString("3. OVERALL_STATUS: Project health (on_track, at_risk, delayed, blocked)\n")
	prompt.WriteString("4. PROGRESS_SUMMARY: Detailed progress update\n")
	prompt.WriteString("5. KEY_ACCOMPLISHMENTS: Major achievements this period\n")
	prompt.WriteString("6. CURRENT_ACTIVITIES: What's happening now\n")
	prompt.WriteString("7. UPCOMING_MILESTONES: Next major milestones\n")
	prompt.WriteString("8. RISKS_AND_ISSUES: Current risks and issues\n")
	prompt.WriteString("9. NEXT_STEPS: Recommended actions\n")
	prompt.WriteString("10. RECOMMENDATIONS: Strategic recommendations\n")

	prompt.WriteString("\nCreate a professional, comprehensive status report appropriate for cloud consulting projects.")

	return prompt.String()
}

// parseTechnicalExplanation parses the AI response into a structured technical explanation
func (s *ClientCommunicationService) parseTechnicalExplanation(content string, req *interfaces.TechnicalExplanationRequest) *interfaces.TechnicalExplanation {
	explanation := &interfaces.TechnicalExplanation{
		ID:              uuid.New().String(),
		AudienceType:    req.AudienceType,
		ComplexityLevel: req.ComplexityLevel,
		GeneratedAt:     time.Now(),
		Metadata:        req.Metadata,
	}

	// Parse sections from the content
	sections := s.parseSections(content)

	if title, exists := sections["TITLE"]; exists {
		explanation.Title = title
	} else {
		explanation.Title = fmt.Sprintf("Technical Explanation: %s", req.TechnicalConcept)
	}

	if summary, exists := sections["EXECUTIVE_SUMMARY"]; exists {
		explanation.ExecutiveSummary = summary
	}

	if value, exists := sections["BUSINESS_VALUE"]; exists {
		explanation.BusinessValue = value
	}

	if overview, exists := sections["TECHNICAL_OVERVIEW"]; exists {
		explanation.TechnicalOverview = overview
	}

	if benefits, exists := sections["KEY_BENEFITS"]; exists {
		explanation.KeyBenefits = s.parseList(benefits)
	}

	if considerations, exists := sections["CONSIDERATIONS"]; exists {
		explanation.Considerations = s.parseList(considerations)
	}

	if steps, exists := sections["NEXT_STEPS"]; exists {
		explanation.NextSteps = s.parseList(steps)
	}

	if glossary, exists := sections["GLOSSARY"]; exists {
		explanation.Glossary = s.parseGlossary(glossary)
	}

	if references, exists := sections["REFERENCES"]; exists {
		explanation.References = s.parseReferences(references)
	}

	return explanation
}

// parsePresentationSlides parses the AI response into structured presentation slides
func (s *ClientCommunicationService) parsePresentationSlides(content string, req *interfaces.PresentationRequest) *interfaces.PresentationSlides {
	slides := &interfaces.PresentationSlides{
		ID:           uuid.New().String(),
		AudienceType: req.AudienceType,
		GeneratedAt:  time.Now(),
		Metadata:     req.Metadata,
	}

	// Parse sections from the content
	sections := s.parseSections(content)

	if title, exists := sections["PRESENTATION_TITLE"]; exists {
		parts := strings.SplitN(title, "\n", 2)
		slides.Title = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			slides.Subtitle = strings.TrimSpace(parts[1])
		}
	} else {
		slides.Title = req.Topic
	}

	if slideContent, exists := sections["DETAILED_SLIDES"]; exists {
		slides.Slides = s.parseSlides(slideContent)
	}

	if notes, exists := sections["SPEAKER_NOTES"]; exists {
		slides.SpeakerNotes = s.parseList(notes)
	}

	if references, exists := sections["REFERENCES"]; exists {
		slides.References = s.parseReferences(references)
	}

	// Estimate total time
	slides.EstimatedTime = req.Duration
	if slides.EstimatedTime == 0 && len(slides.Slides) > 0 {
		slides.EstimatedTime = len(slides.Slides) * 2 // 2 minutes per slide default
	}

	return slides
}

// parseEmailTemplate parses the AI response into a structured email template
func (s *ClientCommunicationService) parseEmailTemplate(content string, req *interfaces.EmailTemplateRequest) *interfaces.EmailTemplate {
	template := &interfaces.EmailTemplate{
		ID:            uuid.New().String(),
		Purpose:       req.Purpose,
		RecipientType: req.RecipientType,
		ToneStyle:     req.ToneStyle,
		GeneratedAt:   time.Now(),
		Metadata:      req.Metadata,
	}

	// Parse sections from the content
	sections := s.parseSections(content)

	if subject, exists := sections["SUBJECT"]; exists {
		template.Subject = strings.TrimSpace(subject)
	}

	if body, exists := sections["EMAIL_BODY"]; exists {
		template.Body = strings.TrimSpace(body)
	}

	if alternatives, exists := sections["ALTERNATIVE_SUBJECTS"]; exists {
		template.Alternatives = s.parseList(alternatives)
	}

	if followUp, exists := sections["FOLLOW_UP_ACTIONS"]; exists {
		template.FollowUpActions = s.parseList(followUp)
	}

	if timing, exists := sections["TIMING_SUGGESTION"]; exists {
		template.SuggestedTiming = strings.TrimSpace(timing)
	}

	return template
}

// parseStatusReport parses the AI response into a structured status report
func (s *ClientCommunicationService) parseStatusReport(content string, req *interfaces.StatusReportRequest) *interfaces.StatusReport {
	report := &interfaces.StatusReport{
		ID:              uuid.New().String(),
		ProjectID:       req.ProjectID,
		ReportingPeriod: req.ReportingPeriod,
		AudienceType:    req.AudienceType,
		GeneratedAt:     time.Now(),
		Metadata:        req.Metadata,
	}

	// Parse sections from the content
	sections := s.parseSections(content)

	if title, exists := sections["REPORT_TITLE"]; exists {
		report.Title = strings.TrimSpace(title)
	} else {
		report.Title = fmt.Sprintf("Status Report - %s", req.ProjectID)
	}

	if summary, exists := sections["EXECUTIVE_SUMMARY"]; exists {
		report.ExecutiveSummary = strings.TrimSpace(summary)
	}

	if status, exists := sections["OVERALL_STATUS"]; exists {
		statusStr := strings.ToLower(strings.TrimSpace(status))
		switch statusStr {
		case "on_track", "on track":
			report.OverallStatus = interfaces.StatusOnTrack
		case "at_risk", "at risk":
			report.OverallStatus = interfaces.StatusAtRisk
		case "delayed":
			report.OverallStatus = interfaces.StatusDelayed
		case "blocked":
			report.OverallStatus = interfaces.StatusBlocked
		case "completed":
			report.OverallStatus = interfaces.StatusCompleted
		default:
			report.OverallStatus = interfaces.StatusOnTrack
		}
	}

	if progress, exists := sections["PROGRESS_SUMMARY"]; exists {
		report.ProgressSummary = strings.TrimSpace(progress)
	}

	if accomplishments, exists := sections["KEY_ACCOMPLISHMENTS"]; exists {
		report.KeyAccomplishments = s.parseList(accomplishments)
	}

	if activities, exists := sections["CURRENT_ACTIVITIES"]; exists {
		report.CurrentActivities = s.parseList(activities)
	}

	if steps, exists := sections["NEXT_STEPS"]; exists {
		report.NextSteps = s.parseList(steps)
	}

	if recommendations, exists := sections["RECOMMENDATIONS"]; exists {
		report.Recommendations = s.parseList(recommendations)
	}

	// Copy provided data
	if req.Budget != nil {
		report.BudgetStatus = req.Budget
	}

	if req.Timeline != nil {
		report.TimelineStatus = req.Timeline
	}

	// Convert milestones
	for _, milestone := range req.Milestones {
		report.UpcomingMilestones = append(report.UpcomingMilestones, milestone)
	}

	return report
}

// parseSections parses content into sections based on headers
func (s *ClientCommunicationService) parseSections(content string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(content, "\n")

	var currentSection string
	var currentContent strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if this is a section header
		if strings.HasSuffix(line, ":") && len(line) > 1 {
			// Save previous section
			if currentSection != "" {
				sections[currentSection] = strings.TrimSpace(currentContent.String())
			}

			// Start new section
			currentSection = strings.TrimSuffix(strings.ToUpper(line), ":")
			currentContent.Reset()
		} else if currentSection != "" && line != "" {
			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
		}
	}

	// Save last section
	if currentSection != "" {
		sections[currentSection] = strings.TrimSpace(currentContent.String())
	}

	return sections
}

// parseList parses a text block into a list of items
func (s *ClientCommunicationService) parseList(content string) []string {
	var items []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove bullet points and numbering
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")

		// Remove numbered list prefixes (1. 2. etc.)
		if len(line) > 2 && line[1] == '.' && line[0] >= '0' && line[0] <= '9' {
			line = strings.TrimSpace(line[2:])
		}

		if line != "" {
			items = append(items, line)
		}
	}

	return items
}

// parseGlossary parses glossary content into a map
func (s *ClientCommunicationService) parseGlossary(content string) map[string]string {
	glossary := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove bullet points
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")

		// Split on colon or dash
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			parts = strings.SplitN(line, " - ", 2)
		}

		if len(parts) == 2 {
			term := strings.TrimSpace(parts[0])
			definition := strings.TrimSpace(parts[1])
			if term != "" && definition != "" {
				glossary[term] = definition
			}
		}
	}

	return glossary
}

// parseReferences parses reference content into a list of references
func (s *ClientCommunicationService) parseReferences(content string) []interfaces.Reference {
	var references []interfaces.Reference
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove bullet points
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "* ")

		// Try to extract URL and title
		if strings.Contains(line, "http") {
			parts := strings.Split(line, "http")
			if len(parts) >= 2 {
				title := strings.TrimSpace(parts[0])
				url := "http" + strings.TrimSpace(parts[1])

				// Clean up title (remove trailing colons, dashes)
				title = strings.TrimSuffix(title, ":")
				title = strings.TrimSuffix(title, " -")
				title = strings.TrimSpace(title)

				if title == "" {
					title = "AWS Documentation"
				}

				references = append(references, interfaces.Reference{
					Title: title,
					URL:   url,
					Type:  "documentation",
				})
			}
		} else if line != "" {
			// Just a title without URL
			references = append(references, interfaces.Reference{
				Title: line,
				Type:  "reference",
			})
		}
	}

	return references
}

// parseSlides parses slide content into structured slides
func (s *ClientCommunicationService) parseSlides(content string) []interfaces.Slide {
	var slides []interfaces.Slide
	sections := strings.Split(content, "\n\n")

	for i, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}

		lines := strings.Split(section, "\n")
		if len(lines) == 0 {
			continue
		}

		slide := interfaces.Slide{
			ID:       uuid.New().String(),
			Order:    i + 1,
			Duration: 2, // Default 2 minutes per slide
		}

		// First line is typically the title
		slide.Title = strings.TrimSpace(lines[0])

		// Determine slide type based on title
		titleLower := strings.ToLower(slide.Title)
		switch {
		case strings.Contains(titleLower, "introduction") || strings.Contains(titleLower, "overview"):
			slide.SlideType = interfaces.SlideTitle
		case strings.Contains(titleLower, "architecture") || strings.Contains(titleLower, "diagram"):
			slide.SlideType = interfaces.SlideDiagram
		case strings.Contains(titleLower, "comparison") || strings.Contains(titleLower, "vs"):
			slide.SlideType = interfaces.SlideComparison
		case strings.Contains(titleLower, "timeline") || strings.Contains(titleLower, "roadmap"):
			slide.SlideType = interfaces.SlideTimeline
		case strings.Contains(titleLower, "conclusion") || strings.Contains(titleLower, "summary"):
			slide.SlideType = interfaces.SlideConclusion
		default:
			slide.SlideType = interfaces.SlideContent
		}

		// Rest of the lines are content
		if len(lines) > 1 {
			contentLines := lines[1:]
			var bulletPoints []string
			var contentBuilder strings.Builder

			for _, line := range contentLines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// Check if it's a bullet point
				if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
					bulletPoints = append(bulletPoints, strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "))
				} else {
					if contentBuilder.Len() > 0 {
						contentBuilder.WriteString("\n")
					}
					contentBuilder.WriteString(line)
				}
			}

			slide.Content = contentBuilder.String()
			slide.BulletPoints = bulletPoints
		}

		slides = append(slides, slide)
	}

	return slides
}
