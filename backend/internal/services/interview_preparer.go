package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// InterviewPreparerService implements the InterviewPreparer interface
type InterviewPreparerService struct {
	bedrockService interfaces.BedrockService
	promptArchitect interfaces.PromptArchitect
}

// NewInterviewPreparerService creates a new InterviewPreparerService
func NewInterviewPreparerService(bedrockService interfaces.BedrockService, promptArchitect interfaces.PromptArchitect) *InterviewPreparerService {
	return &InterviewPreparerService{
		bedrockService:  bedrockService,
		promptArchitect: promptArchitect,
	}
}

// GenerateInterviewGuide generates a comprehensive interview guide for a client inquiry
func (s *InterviewPreparerService) GenerateInterviewGuide(ctx context.Context, inquiry *domain.Inquiry) (*interfaces.InterviewGuide, error) {
	// Build prompt for interview guide generation
	prompt, err := s.buildInterviewGuidePrompt(inquiry)
	if err != nil {
		return nil, fmt.Errorf("failed to build interview guide prompt: %w", err)
	}

	// Generate content using Bedrock
	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate interview guide: %w", err)
	}

	// Parse the response and create structured interview guide
	guide := s.parseInterviewGuideResponse(response.Content, inquiry)
	
	return guide, nil
}

// GenerateQuestionSet generates a set of questions for a specific category and industry
func (s *InterviewPreparerService) GenerateQuestionSet(ctx context.Context, category string, industry string) (*interfaces.QuestionSet, error) {
	prompt := s.buildQuestionSetPrompt(category, industry)

	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   2000,
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate question set: %w", err)
	}

	questionSet := s.parseQuestionSetResponse(response.Content, category, industry)
	
	return questionSet, nil
}

// GenerateDiscoveryChecklist generates a discovery checklist for a specific service type
func (s *InterviewPreparerService) GenerateDiscoveryChecklist(ctx context.Context, serviceType string) (*interfaces.DiscoveryChecklist, error) {
	prompt := s.buildDiscoveryChecklistPrompt(serviceType)

	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   2000,
		Temperature: 0.2,
		TopP:        0.8,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate discovery checklist: %w", err)
	}

	checklist := s.parseDiscoveryChecklistResponse(response.Content, serviceType)
	
	return checklist, nil
}

// GenerateFollowUpQuestions generates follow-up questions based on client responses
func (s *InterviewPreparerService) GenerateFollowUpQuestions(ctx context.Context, responses []interfaces.InterviewResponse) ([]*interfaces.Question, error) {
	if len(responses) == 0 {
		return []*interfaces.Question{}, nil
	}

	prompt := s.buildFollowUpQuestionsPrompt(responses)

	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   1500,
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := s.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate follow-up questions: %w", err)
	}

	questions := s.parseFollowUpQuestionsResponse(response.Content)
	
	return questions, nil
}

// buildInterviewGuidePrompt creates a prompt for generating interview guides
func (s *InterviewPreparerService) buildInterviewGuidePrompt(inquiry *domain.Inquiry) (string, error) {
	services := strings.Join(inquiry.Services, ", ")
	industry := s.inferIndustryFromCompany(inquiry.Company)
	
	prompt := fmt.Sprintf(`You are an expert cloud consultant preparing for a client interview. Generate a comprehensive interview guide for the following client inquiry:

Client Information:
- Company: %s
- Industry: %s
- Services Requested: %s
- Initial Message: %s

Create a structured interview guide with the following format:

TITLE: [Descriptive title for the interview]

OBJECTIVE: [Clear objective for the interview session]

ESTIMATED DURATION: [Time estimate like "90 minutes"]

PREPARATION NOTES:
- [Key preparation point 1]
- [Key preparation point 2]
- [Key preparation point 3]

SECTION 1: BUSINESS CONTEXT AND OBJECTIVES
Objective: [Section objective]
Expected Duration: [Time estimate]
Questions:
1. [MUST-ASK] [Question text] (Type: business, Expected: [answer type])
2. [SHOULD-ASK] [Question text] (Type: open, Expected: [answer type])
3. [NICE-TO-ASK] [Question text] (Type: technical, Expected: [answer type])

SECTION 2: CURRENT INFRASTRUCTURE AND CHALLENGES
Objective: [Section objective]
Expected Duration: [Time estimate]
Questions:
[Continue with numbered questions...]

SECTION 3: TECHNICAL REQUIREMENTS
Objective: [Section objective]
Expected Duration: [Time estimate]
Questions:
[Continue with numbered questions...]

SECTION 4: COMPLIANCE AND SECURITY
Objective: [Section objective]
Expected Duration: [Time estimate]
Questions:
[Continue with numbered questions...]

SECTION 5: TIMELINE AND BUDGET
Objective: [Section objective]
Expected Duration: [Time estimate]
Questions:
[Continue with numbered questions...]

FOLLOW-UP ACTIONS:
- [Action item 1]
- [Action item 2]
- [Action item 3]

Focus on generating specific, actionable questions that will help understand the client's needs, current state, and desired outcomes. Prioritize questions using MUST-ASK, SHOULD-ASK, and NICE-TO-ASK labels.`, 
		inquiry.Company, industry, services, inquiry.Message)

	return prompt, nil
}

// buildQuestionSetPrompt creates a prompt for generating question sets
func (s *InterviewPreparerService) buildQuestionSetPrompt(category, industry string) string {
	return fmt.Sprintf(`Generate a focused set of interview questions for cloud consulting in the following context:

Category: %s
Industry: %s

Create 8-12 specific questions that are:
1. Relevant to the category and industry
2. Actionable and specific
3. Designed to uncover technical and business requirements

Format each question as:
[PRIORITY] Question text (Type: [question_type], Expected: [expected_answer_type])

Where:
- PRIORITY is MUST-ASK, SHOULD-ASK, or NICE-TO-ASK
- question_type is business, technical, open, or closed
- expected_answer_type describes what kind of answer is expected

Focus on questions that will help a cloud consultant provide better recommendations and solutions.`, category, industry)
}

// buildDiscoveryChecklistPrompt creates a prompt for generating discovery checklists
func (s *InterviewPreparerService) buildDiscoveryChecklistPrompt(serviceType string) string {
	return fmt.Sprintf(`Create a comprehensive discovery checklist for cloud consulting service type: %s

Generate a structured checklist with the following sections:

REQUIRED ARTIFACTS:
- [Artifact Name]: [Description] (Type: [document/diagram/data/access], Priority: [high/medium/low], Format: [formats], Source: [who provides it])

TECHNICAL REQUIREMENTS TO GATHER:
- [Requirement 1]
- [Requirement 2]
- [etc.]

BUSINESS REQUIREMENTS TO GATHER:
- [Requirement 1]
- [Requirement 2]
- [etc.]

COMPLIANCE REQUIREMENTS TO ASSESS:
- [Compliance area 1]
- [Compliance area 2]
- [etc.]

ENVIRONMENT DETAILS TO DOCUMENT:
- [Environment detail 1]
- [Environment detail 2]
- [etc.]

Focus on practical, actionable items that a cloud consultant would need to gather during discovery to provide accurate recommendations.`, serviceType)
}

// buildFollowUpQuestionsPrompt creates a prompt for generating follow-up questions
func (s *InterviewPreparerService) buildFollowUpQuestionsPrompt(responses []interfaces.InterviewResponse) string {
	var responseText strings.Builder
	responseText.WriteString("Based on the following client responses, generate 3-5 targeted follow-up questions:\n\n")
	
	for i, response := range responses {
		responseText.WriteString(fmt.Sprintf("Response %d:\n", i+1))
		responseText.WriteString(fmt.Sprintf("Question ID: %s\n", response.QuestionID))
		responseText.WriteString(fmt.Sprintf("Client Response: %s\n", response.Response))
		if response.Notes != "" {
			responseText.WriteString(fmt.Sprintf("Notes: %s\n", response.Notes))
		}
		responseText.WriteString(fmt.Sprintf("Confidence: %s\n\n", response.Confidence))
	}
	
	responseText.WriteString(`Generate follow-up questions that:
1. Clarify ambiguous or incomplete responses
2. Dig deeper into areas that need more detail
3. Identify potential risks or challenges
4. Uncover additional requirements

Format each question as:
[PRIORITY] Question text (Type: [question_type], Expected: [expected_answer_type])`)

	return responseText.String()
}

// parseInterviewGuideResponse parses the AI response into a structured InterviewGuide
func (s *InterviewPreparerService) parseInterviewGuideResponse(content string, inquiry *domain.Inquiry) *interfaces.InterviewGuide {
	guide := &interfaces.InterviewGuide{
		ID:        uuid.New().String(),
		InquiryID: inquiry.ID,
		CreatedAt: time.Now(),
	}

	lines := strings.Split(content, "\n")
	var currentSection *interfaces.InterviewSection
	var inSection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse title
		if strings.HasPrefix(line, "TITLE:") {
			guide.Title = strings.TrimSpace(strings.TrimPrefix(line, "TITLE:"))
		}

		// Parse objective
		if strings.HasPrefix(line, "OBJECTIVE:") {
			guide.Objective = strings.TrimSpace(strings.TrimPrefix(line, "OBJECTIVE:"))
		}

		// Parse estimated duration
		if strings.HasPrefix(line, "ESTIMATED DURATION:") {
			guide.EstimatedDuration = strings.TrimSpace(strings.TrimPrefix(line, "ESTIMATED DURATION:"))
		}

		// Parse preparation notes
		if strings.HasPrefix(line, "PREPARATION NOTES:") {
			continue
		}
		if strings.HasPrefix(line, "- ") && !inSection {
			note := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			guide.PreparationNotes = append(guide.PreparationNotes, note)
		}

		// Parse sections
		if strings.HasPrefix(line, "SECTION ") && strings.Contains(line, ":") {
			if currentSection != nil {
				guide.Sections = append(guide.Sections, *currentSection)
			}
			
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				currentSection = &interfaces.InterviewSection{
					Title: strings.TrimSpace(parts[1]),
				}
				inSection = true
			}
		}

		// Parse section details
		if inSection && currentSection != nil {
			if strings.HasPrefix(line, "Objective:") {
				currentSection.Objective = strings.TrimSpace(strings.TrimPrefix(line, "Objective:"))
			}
			if strings.HasPrefix(line, "Expected Duration:") {
				currentSection.ExpectedDuration = strings.TrimSpace(strings.TrimPrefix(line, "Expected Duration:"))
			}

			// Parse questions
			if strings.Contains(line, "[MUST-ASK]") || strings.Contains(line, "[SHOULD-ASK]") || strings.Contains(line, "[NICE-TO-ASK]") {
				question := s.parseQuestionFromLine(line)
				if question != nil {
					currentSection.Questions = append(currentSection.Questions, question)
				}
			}
		}

		// Parse follow-up actions
		if strings.HasPrefix(line, "FOLLOW-UP ACTIONS:") {
			inSection = false
			if currentSection != nil {
				guide.Sections = append(guide.Sections, *currentSection)
				currentSection = nil
			}
			continue
		}
		if strings.HasPrefix(line, "- ") && !inSection {
			action := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			guide.FollowUpActions = append(guide.FollowUpActions, action)
		}
	}

	// Add the last section if it exists
	if currentSection != nil {
		guide.Sections = append(guide.Sections, *currentSection)
	}

	// Set defaults if not parsed
	if guide.Title == "" {
		guide.Title = fmt.Sprintf("Interview Guide for %s", inquiry.Company)
	}
	if guide.EstimatedDuration == "" {
		guide.EstimatedDuration = "90 minutes"
	}

	return guide
}

// parseQuestionSetResponse parses the AI response into a QuestionSet
func (s *InterviewPreparerService) parseQuestionSetResponse(content string, category, industry string) *interfaces.QuestionSet {
	questionSet := &interfaces.QuestionSet{
		ID:        uuid.New().String(),
		Category:  category,
		Industry:  industry,
		CreatedAt: time.Now(),
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Look for lines that start with priority markers
		if strings.HasPrefix(line, "[MUST-ASK]") || strings.HasPrefix(line, "[SHOULD-ASK]") || strings.HasPrefix(line, "[NICE-TO-ASK]") {
			question := s.parseQuestionFromLine(line)
			if question != nil {
				questionSet.Questions = append(questionSet.Questions, question)
			}
		}
	}

	return questionSet
}

// parseDiscoveryChecklistResponse parses the AI response into a DiscoveryChecklist
func (s *InterviewPreparerService) parseDiscoveryChecklistResponse(content string, serviceType string) *interfaces.DiscoveryChecklist {
	checklist := &interfaces.DiscoveryChecklist{
		ServiceType: serviceType,
	}

	lines := strings.Split(content, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Identify sections
		if strings.HasPrefix(line, "REQUIRED ARTIFACTS:") {
			currentSection = "artifacts"
			continue
		}
		if strings.HasPrefix(line, "TECHNICAL REQUIREMENTS") {
			currentSection = "technical"
			continue
		}
		if strings.HasPrefix(line, "BUSINESS REQUIREMENTS") {
			currentSection = "business"
			continue
		}
		if strings.HasPrefix(line, "COMPLIANCE REQUIREMENTS") {
			currentSection = "compliance"
			continue
		}
		if strings.HasPrefix(line, "ENVIRONMENT DETAILS") {
			currentSection = "environment"
			continue
		}

		// Parse items based on current section
		if strings.HasPrefix(line, "- ") {
			item := strings.TrimSpace(strings.TrimPrefix(line, "- "))
			
			switch currentSection {
			case "artifacts":
				artifact := s.parseArtifactFromLine(item)
				if artifact != nil {
					checklist.RequiredArtifacts = append(checklist.RequiredArtifacts, *artifact)
				}
			case "technical":
				checklist.TechnicalRequirements = append(checklist.TechnicalRequirements, item)
			case "business":
				checklist.BusinessRequirements = append(checklist.BusinessRequirements, item)
			case "compliance":
				checklist.ComplianceRequirements = append(checklist.ComplianceRequirements, item)
			case "environment":
				checklist.EnvironmentDetails = append(checklist.EnvironmentDetails, item)
			}
		}
	}

	return checklist
}

// parseFollowUpQuestionsResponse parses the AI response into follow-up questions
func (s *InterviewPreparerService) parseFollowUpQuestionsResponse(content string) []*interfaces.Question {
	var questions []*interfaces.Question

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "[MUST-ASK]") || strings.Contains(line, "[SHOULD-ASK]") || strings.Contains(line, "[NICE-TO-ASK]") {
			question := s.parseQuestionFromLine(line)
			if question != nil {
				questions = append(questions, question)
			}
		}
	}

	return questions
}

// parseQuestionFromLine parses a question from a formatted line
func (s *InterviewPreparerService) parseQuestionFromLine(line string) *interfaces.Question {
	question := &interfaces.Question{
		ID: uuid.New().String(),
	}

	// Extract priority
	if strings.Contains(line, "[MUST-ASK]") {
		question.Priority = "must-ask"
		line = strings.Replace(line, "[MUST-ASK]", "", 1)
	} else if strings.Contains(line, "[SHOULD-ASK]") {
		question.Priority = "should-ask"
		line = strings.Replace(line, "[SHOULD-ASK]", "", 1)
	} else if strings.Contains(line, "[NICE-TO-ASK]") {
		question.Priority = "nice-to-ask"
		line = strings.Replace(line, "[NICE-TO-ASK]", "", 1)
	}

	// Extract question text and metadata
	if strings.Contains(line, "(Type:") {
		parts := strings.Split(line, "(Type:")
		question.Text = strings.TrimSpace(parts[0])
		
		if len(parts) > 1 {
			metadata := parts[1]
			
			// Extract type
			if strings.Contains(metadata, ",") {
				typePart := strings.Split(metadata, ",")[0]
				question.Type = strings.TrimSpace(typePart)
			}
			
			// Extract expected answer type
			if strings.Contains(metadata, "Expected:") {
				expectedPart := strings.Split(metadata, "Expected:")[1]
				expectedPart = strings.TrimSuffix(expectedPart, ")")
				question.ExpectedAnswerType = strings.TrimSpace(expectedPart)
			}
		}
	} else {
		question.Text = strings.TrimSpace(line)
	}

	// Set defaults
	if question.Type == "" {
		question.Type = "open"
	}
	if question.Priority == "" {
		question.Priority = "should-ask"
	}

	// Infer category from question content
	question.Category = s.inferQuestionCategory(question.Text)

	return question
}

// parseArtifactFromLine parses an artifact from a formatted line
func (s *InterviewPreparerService) parseArtifactFromLine(line string) *interfaces.Artifact {
	// Expected format: [Name]: [Description] (Type: [type], Priority: [priority], Format: [formats], Source: [source])
	if !strings.Contains(line, ":") {
		return &interfaces.Artifact{
			Name:        line,
			Description: line,
			Type:        "document",
			Priority:    "medium",
			Format:      []string{"any"},
			Source:      "client",
		}
	}

	parts := strings.SplitN(line, ":", 2)
	artifact := &interfaces.Artifact{
		Name: strings.TrimSpace(parts[0]),
	}

	if len(parts) > 1 {
		remaining := parts[1]
		
		// Extract description (everything before the first parenthesis)
		if strings.Contains(remaining, "(") {
			descParts := strings.Split(remaining, "(")
			artifact.Description = strings.TrimSpace(descParts[0])
			
			// Parse metadata from parentheses
			if len(descParts) > 1 {
				metadata := descParts[1]
				metadata = strings.TrimSuffix(metadata, ")")
				
				// Parse type
				if strings.Contains(metadata, "Type:") {
					typePart := s.extractMetadataValue(metadata, "Type:")
					artifact.Type = typePart
				}
				
				// Parse priority
				if strings.Contains(metadata, "Priority:") {
					priorityPart := s.extractMetadataValue(metadata, "Priority:")
					artifact.Priority = priorityPart
				}
				
				// Parse format
				if strings.Contains(metadata, "Format:") {
					formatPart := s.extractMetadataValue(metadata, "Format:")
					artifact.Format = strings.Split(formatPart, "/")
				}
				
				// Parse source
				if strings.Contains(metadata, "Source:") {
					sourcePart := s.extractMetadataValue(metadata, "Source:")
					artifact.Source = sourcePart
				}
			}
		} else {
			artifact.Description = strings.TrimSpace(remaining)
		}
	}

	// Set defaults
	if artifact.Type == "" {
		artifact.Type = "document"
	}
	if artifact.Priority == "" {
		artifact.Priority = "medium"
	}
	if len(artifact.Format) == 0 {
		artifact.Format = []string{"any"}
	}
	if artifact.Source == "" {
		artifact.Source = "client"
	}

	return artifact
}

// extractMetadataValue extracts a value from metadata string
func (s *InterviewPreparerService) extractMetadataValue(metadata, key string) string {
	if !strings.Contains(metadata, key) {
		return ""
	}
	
	parts := strings.Split(metadata, key)
	if len(parts) < 2 {
		return ""
	}
	
	value := parts[1]
	if strings.Contains(value, ",") {
		value = strings.Split(value, ",")[0]
	}
	
	return strings.TrimSpace(value)
}

// inferIndustryFromCompany attempts to infer industry from company name
func (s *InterviewPreparerService) inferIndustryFromCompany(company string) string {
	company = strings.ToLower(company)
	
	// Healthcare indicators
	if strings.Contains(company, "hospital") || strings.Contains(company, "medical") || 
	   strings.Contains(company, "health") || strings.Contains(company, "clinic") {
		return "healthcare"
	}
	
	// Financial indicators
	if strings.Contains(company, "bank") || strings.Contains(company, "financial") || 
	   strings.Contains(company, "credit") || strings.Contains(company, "insurance") {
		return "financial"
	}
	
	// Technology indicators
	if strings.Contains(company, "tech") || strings.Contains(company, "software") || 
	   strings.Contains(company, "systems") || strings.Contains(company, "solutions") {
		return "technology"
	}
	
	// Manufacturing indicators
	if strings.Contains(company, "manufacturing") || strings.Contains(company, "industrial") || 
	   strings.Contains(company, "factory") || strings.Contains(company, "production") {
		return "manufacturing"
	}
	
	// Retail indicators
	if strings.Contains(company, "retail") || strings.Contains(company, "store") || 
	   strings.Contains(company, "commerce") || strings.Contains(company, "shopping") {
		return "retail"
	}
	
	return "general"
}

// inferQuestionCategory infers the category of a question from its content
func (s *InterviewPreparerService) inferQuestionCategory(questionText string) string {
	text := strings.ToLower(questionText)
	
	// Business category indicators
	if strings.Contains(text, "budget") || strings.Contains(text, "cost") || 
	   strings.Contains(text, "business") || strings.Contains(text, "roi") ||
	   strings.Contains(text, "timeline") || strings.Contains(text, "stakeholder") {
		return "business"
	}
	
	// Security category indicators
	if strings.Contains(text, "security") || strings.Contains(text, "compliance") || 
	   strings.Contains(text, "audit") || strings.Contains(text, "access") ||
	   strings.Contains(text, "authentication") || strings.Contains(text, "encryption") {
		return "security"
	}
	
	// Infrastructure category indicators
	if strings.Contains(text, "infrastructure") || strings.Contains(text, "server") || 
	   strings.Contains(text, "network") || strings.Contains(text, "storage") ||
	   strings.Contains(text, "database") || strings.Contains(text, "architecture") {
		return "infrastructure"
	}
	
	// Migration category indicators
	if strings.Contains(text, "migration") || strings.Contains(text, "migrate") || 
	   strings.Contains(text, "move") || strings.Contains(text, "transfer") || 
	   strings.Contains(text, "legacy") {
		return "migration"
	}
	
	return "general"
}