package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// EnhancedBedrockService wraps the base Bedrock service with knowledge base integration
type EnhancedBedrockService struct {
	bedrockService        interfaces.BedrockService
	knowledgeBase         interfaces.KnowledgeBase
	companyKnowledgeInteg *CompanyKnowledgeIntegrationService
}

// NewEnhancedBedrockService creates a new enhanced Bedrock service
func NewEnhancedBedrockService(bedrock interfaces.BedrockService, kb interfaces.KnowledgeBase, companyInteg *CompanyKnowledgeIntegrationService) *EnhancedBedrockService {
	return &EnhancedBedrockService{
		bedrockService:        bedrock,
		knowledgeBase:         kb,
		companyKnowledgeInteg: companyInteg,
	}
}

// GenerateEnhancedResponse generates AI responses enhanced with company knowledge
func (e *EnhancedBedrockService) GenerateEnhancedResponse(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	// Build base prompt
	basePrompt := fmt.Sprintf("Generate a professional consulting response for the following client inquiry:\n\nClient: %s (%s)\nServices: %s\nMessage: %s\n",
		inquiry.Name, inquiry.Company, strings.Join(inquiry.Services, ", "), inquiry.Message)

	// Use company knowledge integration to create contextual prompt
	enhancedPrompt, err := e.companyKnowledgeInteg.GenerateContextualPrompt(ctx, inquiry, basePrompt)
	if err != nil {
		// Fall back to base prompt if enhancement fails
		enhancedPrompt = basePrompt
	}

	// Generate response using enhanced prompt
	return e.bedrockService.GenerateText(ctx, enhancedPrompt, options)
}

// GenerateEnhancedResponseWithRecommendations generates AI responses with specific recommendations
func (e *EnhancedBedrockService) GenerateEnhancedResponseWithRecommendations(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.BedrockOptions) (*EnhancedBedrockResponse, error) {
	// Generate the enhanced response
	response, err := e.GenerateEnhancedResponse(ctx, inquiry, options)
	if err != nil {
		return nil, err
	}

	// Get specific recommendations
	recommendations, err := e.companyKnowledgeInteg.GetRecommendationsForInquiry(ctx, inquiry)
	if err != nil {
		// Continue without recommendations if they fail
		recommendations = &InquiryRecommendations{
			InquiryID:   inquiry.ID,
			GeneratedAt: time.Now(),
		}
	}

	return &EnhancedBedrockResponse{
		BedrockResponse: *response,
		Recommendations: *recommendations,
		CompanyContext:  nil, // Will be populated if needed
	}, nil
}

// buildEnhancedPrompt creates a prompt enriched with company-specific knowledge
func (e *EnhancedBedrockService) buildEnhancedPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	var promptBuilder strings.Builder

	// Start with company context
	promptBuilder.WriteString("You are an AI assistant for a cloud consulting company. ")
	promptBuilder.WriteString("Use the following company knowledge to provide accurate, relevant responses that align with our business model and expertise.\n\n")

	// Add service offerings context
	if err := e.addServiceOfferingsContext(ctx, &promptBuilder, inquiry); err != nil {
		return "", err
	}

	// Add team expertise context
	if err := e.addTeamExpertiseContext(ctx, &promptBuilder, inquiry); err != nil {
		return "", err
	}

	// Add past solutions context
	if err := e.addPastSolutionsContext(ctx, &promptBuilder, inquiry); err != nil {
		return "", err
	}

	// Add methodology context
	if err := e.addMethodologyContext(ctx, &promptBuilder, inquiry); err != nil {
		return "", err
	}

	// Add client history context if available
	if inquiry.Company != "" {
		if err := e.addClientHistoryContext(ctx, &promptBuilder, inquiry); err != nil {
			return "", err
		}
	}

	// Add the actual inquiry
	promptBuilder.WriteString("\n## CLIENT INQUIRY\n")
	promptBuilder.WriteString(fmt.Sprintf("Client: %s (%s)\n", inquiry.Name, inquiry.Company))
	promptBuilder.WriteString(fmt.Sprintf("Services Requested: %s\n", strings.Join(inquiry.Services, ", ")))
	promptBuilder.WriteString(fmt.Sprintf("Message: %s\n", inquiry.Message))

	// Add response instructions
	promptBuilder.WriteString("\n## RESPONSE INSTRUCTIONS\n")
	promptBuilder.WriteString("Based on the company knowledge above, provide a professional response that:\n")
	promptBuilder.WriteString("1. Addresses the client's specific needs using our service offerings\n")
	promptBuilder.WriteString("2. References relevant past experience and solutions\n")
	promptBuilder.WriteString("3. Suggests appropriate team members based on expertise\n")
	promptBuilder.WriteString("4. Follows our consulting methodology and approach\n")
	promptBuilder.WriteString("5. Provides realistic timelines and expectations\n")
	promptBuilder.WriteString("6. Maintains a professional, consultative tone\n")
	promptBuilder.WriteString("7. Includes next steps for engagement\n\n")

	promptBuilder.WriteString("Generate a comprehensive response:")

	return promptBuilder.String(), nil
}

// addServiceOfferingsContext adds relevant service offerings to the prompt
func (e *EnhancedBedrockService) addServiceOfferingsContext(ctx context.Context, builder *strings.Builder, inquiry *domain.Inquiry) error {
	offerings, err := e.knowledgeBase.GetServiceOfferings(ctx)
	if err != nil {
		return err
	}

	builder.WriteString("## OUR SERVICE OFFERINGS\n")

	// Filter relevant offerings based on inquiry services
	relevantOfferings := e.filterRelevantOfferings(offerings, inquiry.Services)

	for _, offering := range relevantOfferings {
		builder.WriteString(fmt.Sprintf("### %s\n", offering.Name))
		builder.WriteString(fmt.Sprintf("- Description: %s\n", offering.Description))
		builder.WriteString(fmt.Sprintf("- Duration: %s\n", offering.TypicalDuration))
		builder.WriteString(fmt.Sprintf("- Team Size: %s\n", offering.TeamSize))
		builder.WriteString("- Key Benefits:\n")
		for _, benefit := range offering.KeyBenefits {
			builder.WriteString(fmt.Sprintf("  • %s\n", benefit))
		}
		builder.WriteString("- Deliverables:\n")
		for _, deliverable := range offering.Deliverables {
			builder.WriteString(fmt.Sprintf("  • %s\n", deliverable))
		}
		builder.WriteString("\n")
	}

	return nil
}

// addTeamExpertiseContext adds relevant team expertise to the prompt
func (e *EnhancedBedrockService) addTeamExpertiseContext(ctx context.Context, builder *strings.Builder, inquiry *domain.Inquiry) error {
	expertise, err := e.knowledgeBase.GetTeamExpertise(ctx)
	if err != nil {
		return err
	}

	builder.WriteString("## OUR TEAM EXPERTISE\n")

	// Filter relevant expertise based on inquiry
	relevantExpertise := e.filterRelevantExpertise(expertise, inquiry.Services)

	for _, exp := range relevantExpertise {
		builder.WriteString(fmt.Sprintf("### %s - %s\n", exp.ConsultantName, exp.Role))
		builder.WriteString(fmt.Sprintf("- Experience: %d years\n", exp.ExperienceYears))
		builder.WriteString("- Expertise Areas:\n")
		for _, area := range exp.ExpertiseAreas {
			builder.WriteString(fmt.Sprintf("  • %s\n", area))
		}
		builder.WriteString("- Cloud Providers:\n")
		for _, provider := range exp.CloudProviders {
			builder.WriteString(fmt.Sprintf("  • %s\n", provider))
		}
		builder.WriteString("\n")
	}

	return nil
}

// addPastSolutionsContext adds relevant past solutions to the prompt
func (e *EnhancedBedrockService) addPastSolutionsContext(ctx context.Context, builder *strings.Builder, inquiry *domain.Inquiry) error {
	// Infer industry from company name or use general search
	industry := e.inferIndustry(inquiry.Company)
	serviceType := e.inferServiceType(inquiry.Services)

	solutions, err := e.knowledgeBase.GetPastSolutions(ctx, serviceType, industry)
	if err != nil {
		return err
	}

	if len(solutions) > 0 {
		builder.WriteString("## RELEVANT PAST SOLUTIONS\n")

		// Limit to top 3 most relevant solutions
		limit := 3
		if len(solutions) < limit {
			limit = len(solutions)
		}

		for i := 0; i < limit; i++ {
			solution := solutions[i]
			builder.WriteString(fmt.Sprintf("### %s\n", solution.Title))
			builder.WriteString(fmt.Sprintf("- Industry: %s\n", solution.Industry))
			builder.WriteString(fmt.Sprintf("- Problem: %s\n", solution.ProblemStatement))
			builder.WriteString(fmt.Sprintf("- Solution: %s\n", solution.SolutionApproach))
			builder.WriteString(fmt.Sprintf("- Technologies: %s\n", strings.Join(solution.Technologies, ", ")))
			builder.WriteString(fmt.Sprintf("- Time to Value: %s\n", solution.TimeToValue))
			if solution.CostSavings > 0 {
				builder.WriteString(fmt.Sprintf("- Cost Savings: $%.0f\n", solution.CostSavings))
			}
			builder.WriteString("\n")
		}
	}

	return nil
}

// addMethodologyContext adds relevant methodology information to the prompt
func (e *EnhancedBedrockService) addMethodologyContext(ctx context.Context, builder *strings.Builder, inquiry *domain.Inquiry) error {
	serviceType := e.inferServiceType(inquiry.Services)

	methodology, err := e.knowledgeBase.GetConsultingApproach(ctx, serviceType)
	if err != nil {
		// If no specific methodology found, continue without error
		return nil
	}

	builder.WriteString("## OUR CONSULTING APPROACH\n")
	builder.WriteString(fmt.Sprintf("### %s\n", methodology.Name))
	builder.WriteString(fmt.Sprintf("- Philosophy: %s\n", methodology.Philosophy))
	builder.WriteString(fmt.Sprintf("- Engagement Model: %s\n", methodology.EngagementModel))
	builder.WriteString("- Key Principles:\n")
	for _, principle := range methodology.KeyPrinciples {
		builder.WriteString(fmt.Sprintf("  • %s\n", principle))
	}
	builder.WriteString(fmt.Sprintf("- Client Involvement: %s\n", methodology.ClientInvolvement))
	builder.WriteString(fmt.Sprintf("- Knowledge Transfer: %s\n", methodology.KnowledgeTransfer))
	builder.WriteString("\n")

	return nil
}

// addClientHistoryContext adds client history if available
func (e *EnhancedBedrockService) addClientHistoryContext(ctx context.Context, builder *strings.Builder, inquiry *domain.Inquiry) error {
	history, err := e.knowledgeBase.GetClientHistory(ctx, inquiry.Company)
	if err != nil || len(history) == 0 {
		return nil // No history available, continue without error
	}

	builder.WriteString("## CLIENT HISTORY\n")
	builder.WriteString(fmt.Sprintf("Previous engagements with %s:\n", inquiry.Company))

	for _, engagement := range history {
		builder.WriteString(fmt.Sprintf("- %s (%s): %s\n",
			engagement.ProjectName,
			engagement.StartDate.Format("2006"),
			engagement.Status))
		if engagement.ClientSatisfaction > 0 {
			builder.WriteString(fmt.Sprintf("  Satisfaction: %.1f/10\n", engagement.ClientSatisfaction))
		}
	}
	builder.WriteString("\n")

	return nil
}

// Helper methods for filtering and inference

func (e *EnhancedBedrockService) filterRelevantOfferings(offerings []*interfaces.ServiceOffering, services []string) []*interfaces.ServiceOffering {
	var relevant []*interfaces.ServiceOffering

	for _, offering := range offerings {
		for _, service := range services {
			if e.isServiceMatch(offering, service) {
				relevant = append(relevant, offering)
				break
			}
		}
	}

	// If no specific matches, return all offerings (up to 3)
	if len(relevant) == 0 {
		limit := 3
		if len(offerings) < limit {
			limit = len(offerings)
		}
		relevant = offerings[:limit]
	}

	return relevant
}

func (e *EnhancedBedrockService) filterRelevantExpertise(expertise []*interfaces.TeamExpertise, services []string) []*interfaces.TeamExpertise {
	var relevant []*interfaces.TeamExpertise

	for _, exp := range expertise {
		for _, service := range services {
			if e.isExpertiseMatch(exp, service) {
				relevant = append(relevant, exp)
				break
			}
		}
	}

	// Limit to top 3 consultants
	limit := 3
	if len(relevant) < limit {
		limit = len(relevant)
	}
	if len(relevant) > 0 {
		relevant = relevant[:limit]
	}

	return relevant
}

func (e *EnhancedBedrockService) isServiceMatch(offering *interfaces.ServiceOffering, service string) bool {
	serviceLower := strings.ToLower(service)
	return strings.Contains(strings.ToLower(offering.Name), serviceLower) ||
		strings.Contains(strings.ToLower(offering.Description), serviceLower) ||
		strings.Contains(strings.ToLower(offering.Category), serviceLower)
}

func (e *EnhancedBedrockService) isExpertiseMatch(expertise *interfaces.TeamExpertise, service string) bool {
	serviceLower := strings.ToLower(service)

	for _, area := range expertise.ExpertiseAreas {
		if strings.Contains(strings.ToLower(area), serviceLower) {
			return true
		}
	}

	return false
}

func (e *EnhancedBedrockService) inferIndustry(company string) string {
	companyLower := strings.ToLower(company)

	// Simple industry inference based on company name
	if strings.Contains(companyLower, "bank") || strings.Contains(companyLower, "financial") {
		return "Financial Services"
	}
	if strings.Contains(companyLower, "health") || strings.Contains(companyLower, "medical") {
		return "Healthcare"
	}
	if strings.Contains(companyLower, "tech") || strings.Contains(companyLower, "software") {
		return "Technology"
	}
	if strings.Contains(companyLower, "retail") || strings.Contains(companyLower, "store") {
		return "Retail"
	}

	return "" // Return empty for general search
}

func (e *EnhancedBedrockService) inferServiceType(services []string) string {
	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "migration") {
			return string(domain.ReportTypeMigration)
		}
		if strings.Contains(serviceLower, "architecture") {
			return string(domain.ReportTypeArchitectureReview)
		}
		if strings.Contains(serviceLower, "assessment") {
			return string(domain.ReportTypeAssessment)
		}
		if strings.Contains(serviceLower, "optimization") {
			return string(domain.ReportTypeOptimization)
		}
	}

	return string(domain.ReportTypeGeneral)
}

// Delegate methods to base service

func (e *EnhancedBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	return e.bedrockService.GenerateText(ctx, prompt, options)
}

func (e *EnhancedBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return e.bedrockService.GetModelInfo()
}

func (e *EnhancedBedrockService) IsHealthy() bool {
	return e.bedrockService.IsHealthy()
}

// EnhancedBedrockResponse represents an enhanced response with company-specific recommendations
type EnhancedBedrockResponse struct {
	interfaces.BedrockResponse
	Recommendations InquiryRecommendations `json:"recommendations"`
	CompanyContext  *CompanyContext        `json:"company_context,omitempty"`
}

// ChatAwareBedrockService extends EnhancedBedrockService with chat-specific capabilities
type ChatAwareBedrockService struct {
	*EnhancedBedrockService
	promptArchitect interfaces.PromptArchitect
	responseCache   map[string]*CachedResponse
	quickActions    map[string]QuickActionHandler
}

// CachedResponse represents a cached AI response
type CachedResponse struct {
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  time.Time              `json:"created_at"`
	ExpiresAt  time.Time              `json:"expires_at"`
	TokensUsed int                    `json:"tokens_used"`
}

// QuickActionHandler defines a handler for quick actions
type QuickActionHandler func(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error)

// NewChatAwareBedrockService creates a new chat-aware Bedrock service
func NewChatAwareBedrockService(
	enhanced *EnhancedBedrockService,
	promptArchitect interfaces.PromptArchitect,
) *ChatAwareBedrockService {
	service := &ChatAwareBedrockService{
		EnhancedBedrockService: enhanced,
		promptArchitect:        promptArchitect,
		responseCache:          make(map[string]*CachedResponse),
		quickActions:           make(map[string]QuickActionHandler),
	}

	// Initialize quick action handlers
	service.initializeQuickActions()

	return service
}

// GenerateChatResponse generates an AI response for chat scenarios with context awareness
func (c *ChatAwareBedrockService) GenerateChatResponse(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	// Check cache first for similar queries
	if cachedResponse := c.getCachedResponse(request.Content); cachedResponse != nil {
		return &ChatResponse{
			Content:     cachedResponse.Content,
			Metadata:    cachedResponse.Metadata,
			TokensUsed:  cachedResponse.TokensUsed,
			ProcessTime: 0, // Cached response is instant
		}, nil
	}

	// Handle quick actions
	if request.QuickAction != "" {
		if handler, exists := c.quickActions[request.QuickAction]; exists {
			content, err := handler(ctx, request.Session, request.Context)
			if err != nil {
				return nil, fmt.Errorf("quick action failed: %w", err)
			}

			return &ChatResponse{
				Content:     content,
				Metadata:    map[string]interface{}{"quick_action": request.QuickAction},
				TokensUsed:  0, // Quick actions don't use tokens
				ProcessTime: 0,
			}, nil
		}
	}

	// Generate context-aware prompt
	prompt, err := c.GetContextualPrompt(ctx, request.Session, request.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate contextual prompt: %w", err)
	}

	// Optimize prompt to reduce token usage
	optimizedPrompt := c.optimizePrompt(prompt)

	// Generate response using enhanced Bedrock service
	// Use the configured model ID from the base service
	modelInfo := c.EnhancedBedrockService.GetModelInfo()
	options := &interfaces.BedrockOptions{
		ModelID:     modelInfo.ModelID,
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	startTime := time.Now()
	response, err := c.EnhancedBedrockService.GenerateText(ctx, optimizedPrompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}
	processingTime := time.Since(startTime)

	// Validate response quality
	if err := c.validateResponseQuality(response.Content); err != nil {
		// Try fallback response
		fallbackContent := c.getFallbackResponse(request.Content)
		if fallbackContent != "" {
			response.Content = fallbackContent
		} else {
			return nil, fmt.Errorf("response quality validation failed: %w", err)
		}
	}

	chatResponse := &ChatResponse{
		Content: response.Content,
		Metadata: map[string]interface{}{
			"model_id":          options.ModelID,
			"processing_time":   processingTime.Milliseconds(),
			"prompt_tokens":     response.Usage.InputTokens,
			"completion_tokens": response.Usage.OutputTokens,
		},
		TokensUsed:  response.Usage.OutputTokens,
		ProcessTime: float64(processingTime.Milliseconds()),
	}

	// Cache the response for similar future queries
	c.cacheResponse(request.Content, chatResponse)

	return chatResponse, nil
}

// GetContextualPrompt generates a context-aware prompt for chat scenarios
func (c *ChatAwareBedrockService) GetContextualPrompt(ctx context.Context, session *domain.ChatSession, message string) (string, error) {
	var promptBuilder strings.Builder

	// Start with chat-specific system prompt
	promptBuilder.WriteString("You are an expert AWS cloud consultant assistant providing real-time support during client meetings. ")
	promptBuilder.WriteString("Your role is to help consultants provide immediate, authoritative, and specific answers to client questions.\n\n")

	// Add consultant guidelines
	promptBuilder.WriteString("CONSULTANT GUIDELINES:\n")
	promptBuilder.WriteString("- Always provide specific AWS service recommendations\n")
	promptBuilder.WriteString("- Include approximate costs when relevant\n")
	promptBuilder.WriteString("- Mention compliance considerations for regulated industries\n")
	promptBuilder.WriteString("- Suggest next steps or follow-up actions\n")
	promptBuilder.WriteString("- Reference AWS documentation or whitepapers when helpful\n")
	promptBuilder.WriteString("- Keep responses concise but comprehensive (2-3 sentences max unless technical details are needed)\n\n")

	// Add session context
	if session.ClientName != "" {
		promptBuilder.WriteString(fmt.Sprintf("CLIENT: %s\n", session.ClientName))
	}

	if session.Context != "" {
		promptBuilder.WriteString(fmt.Sprintf("MEETING CONTEXT: %s\n", session.Context))
	}

	// Add session metadata context
	if session.Metadata != nil {
		if meetingType, ok := session.Metadata["meeting_type"].(string); ok && meetingType != "" {
			promptBuilder.WriteString(fmt.Sprintf("MEETING TYPE: %s\n", meetingType))
		}

		if serviceTypes, ok := session.Metadata["service_types"].([]string); ok && len(serviceTypes) > 0 {
			promptBuilder.WriteString(fmt.Sprintf("SERVICES OF INTEREST: %s\n", strings.Join(serviceTypes, ", ")))
		}

		if cloudProviders, ok := session.Metadata["cloud_providers"].([]string); ok && len(cloudProviders) > 0 {
			promptBuilder.WriteString(fmt.Sprintf("CLOUD PROVIDERS: %s\n", strings.Join(cloudProviders, ", ")))
		}
	}

	// Add relevant company knowledge based on session context
	if err := c.addChatContextualKnowledge(ctx, &promptBuilder, session, message); err != nil {
		// Log error but continue without enhanced context
		fmt.Printf("Warning: Failed to add contextual knowledge: %v\n", err)
	}

	// Add the current question
	promptBuilder.WriteString(fmt.Sprintf("\nCURRENT QUESTION: %s\n", message))
	promptBuilder.WriteString("\nProvide a direct, expert-level response that the consultant can immediately use with their client:")

	return promptBuilder.String(), nil
}

// ProcessQuickAction processes a quick action request
func (c *ChatAwareBedrockService) ProcessQuickAction(ctx context.Context, action string, context *domain.SessionContext) (string, error) {
	if handler, exists := c.quickActions[action]; exists {
		// Create a temporary session for quick action processing
		session := &domain.ChatSession{
			ClientName: context.ClientName,
			Context:    context.ProjectContext,
			Metadata: map[string]interface{}{
				"meeting_type":    context.MeetingType,
				"service_types":   context.ServiceTypes,
				"cloud_providers": context.CloudProviders,
				"custom_fields":   context.CustomFields,
			},
		}

		return handler(ctx, session, context)
	}

	return "", fmt.Errorf("unknown quick action: %s", action)
}

// ValidateResponse validates the quality of an AI response
func (c *ChatAwareBedrockService) ValidateResponse(response string) error {
	return c.validateResponseQuality(response)
}

// Private helper methods

// initializeQuickActions sets up handlers for common quick actions
func (c *ChatAwareBedrockService) initializeQuickActions() {
	c.quickActions["cost_estimate"] = c.handleCostEstimate
	c.quickActions["architecture_review"] = c.handleArchitectureReview
	c.quickActions["migration_plan"] = c.handleMigrationPlan
	c.quickActions["security_assessment"] = c.handleSecurityAssessment
	c.quickActions["compliance_check"] = c.handleComplianceCheck
	c.quickActions["best_practices"] = c.handleBestPractices
}

// handleCostEstimate provides quick cost estimation guidance
func (c *ChatAwareBedrockService) handleCostEstimate(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error) {
	response := "For accurate cost estimation, I recommend using the AWS Pricing Calculator (calculator.aws). "
	response += "Key factors to consider: compute hours, storage requirements, data transfer, and any managed services. "
	response += "Would you like me to help estimate costs for specific services or workloads?"

	return response, nil
}

// handleArchitectureReview provides architecture review guidance
func (c *ChatAwareBedrockService) handleArchitectureReview(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error) {
	response := "For architecture reviews, I follow the AWS Well-Architected Framework's 6 pillars: "
	response += "Operational Excellence, Security, Reliability, Performance Efficiency, Cost Optimization, and Sustainability. "
	response += "What specific aspect of the architecture would you like me to focus on?"

	return response, nil
}

// handleMigrationPlan provides migration planning guidance
func (c *ChatAwareBedrockService) handleMigrationPlan(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error) {
	response := "AWS migration typically follows the 6 R's strategy: Rehost, Replatform, Refactor, Repurchase, Retain, or Retire. "
	response += "I recommend starting with AWS Application Discovery Service and Migration Hub for assessment. "
	response += "What type of workloads are you looking to migrate?"

	return response, nil
}

// handleSecurityAssessment provides security assessment guidance
func (c *ChatAwareBedrockService) handleSecurityAssessment(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error) {
	response := "For security assessments, I recommend AWS Security Hub, GuardDuty, and Config for continuous monitoring. "
	response += "Key areas to review: IAM policies, network security groups, encryption at rest/transit, and compliance frameworks. "
	response += "Are there specific security concerns or compliance requirements we should address?"

	return response, nil
}

// handleComplianceCheck provides compliance guidance
func (c *ChatAwareBedrockService) handleComplianceCheck(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error) {
	response := "AWS offers compliance programs for SOC, PCI DSS, HIPAA, FedRAMP, and more. "
	response += "I recommend AWS Artifact for compliance reports and AWS Config for compliance monitoring. "
	response += "What specific compliance frameworks do you need to meet?"

	return response, nil
}

// handleBestPractices provides best practices guidance
func (c *ChatAwareBedrockService) handleBestPractices(ctx context.Context, session *domain.ChatSession, context *domain.SessionContext) (string, error) {
	response := "AWS best practices include: using IAM roles instead of keys, enabling CloudTrail logging, "
	response += "implementing least privilege access, using multiple AZs for high availability, and regular backups. "
	response += "What specific area would you like best practice recommendations for?"

	return response, nil
}

// addChatContextualKnowledge adds relevant company knowledge based on chat context
func (c *ChatAwareBedrockService) addChatContextualKnowledge(ctx context.Context, builder *strings.Builder, session *domain.ChatSession, message string) error {
	// Infer service type from message content and session context
	serviceType := c.inferServiceTypeFromChat(message, session)

	// Add relevant service offerings
	if offerings, err := c.knowledgeBase.GetServiceOfferings(ctx); err == nil {
		relevantOfferings := c.filterOfferingsForChat(offerings, serviceType, message)
		if len(relevantOfferings) > 0 {
			builder.WriteString("\nRELEVANT SERVICES:\n")
			for _, offering := range relevantOfferings {
				fmt.Fprintf(builder, "- %s: %s\n", offering.Name, offering.Description)
			}
		}
	}

	// Add relevant past solutions
	industry := c.inferIndustryFromSession(session)
	if solutions, err := c.knowledgeBase.GetPastSolutions(ctx, serviceType, industry); err == nil && len(solutions) > 0 {
		builder.WriteString("\nRELEVANT EXPERIENCE:\n")
		// Limit to 2 most relevant solutions for chat context
		limit := 2
		if len(solutions) < limit {
			limit = len(solutions)
		}
		for i := 0; i < limit; i++ {
			solution := solutions[i]
			fmt.Fprintf(builder, "- %s: %s (Industry: %s)\n",
				solution.Title, solution.SolutionApproach, solution.Industry)
		}
	}

	return nil
}

// inferServiceTypeFromChat infers service type from chat message and session
func (c *ChatAwareBedrockService) inferServiceTypeFromChat(message string, session *domain.ChatSession) string {
	messageLower := strings.ToLower(message)

	// Check session metadata first
	if session.Metadata != nil {
		if serviceTypes, ok := session.Metadata["service_types"].([]string); ok && len(serviceTypes) > 0 {
			// Use the first service type as primary
			return strings.ToLower(serviceTypes[0])
		}
	}

	// Infer from message content
	if strings.Contains(messageLower, "migrate") || strings.Contains(messageLower, "migration") {
		return string(domain.ReportTypeMigration)
	}
	if strings.Contains(messageLower, "architecture") || strings.Contains(messageLower, "design") {
		return string(domain.ReportTypeArchitectureReview)
	}
	if strings.Contains(messageLower, "assess") || strings.Contains(messageLower, "evaluation") {
		return string(domain.ReportTypeAssessment)
	}
	if strings.Contains(messageLower, "optimize") || strings.Contains(messageLower, "cost") {
		return string(domain.ReportTypeOptimization)
	}

	return string(domain.ReportTypeGeneral)
}

// inferIndustryFromSession infers industry from session context
func (c *ChatAwareBedrockService) inferIndustryFromSession(session *domain.ChatSession) string {
	if session.ClientName != "" {
		return c.inferIndustry(session.ClientName)
	}

	// Check for industry hints in context
	if session.Context != "" {
		contextLower := strings.ToLower(session.Context)
		if strings.Contains(contextLower, "healthcare") || strings.Contains(contextLower, "medical") {
			return "Healthcare"
		}
		if strings.Contains(contextLower, "financial") || strings.Contains(contextLower, "banking") {
			return "Financial Services"
		}
		if strings.Contains(contextLower, "retail") || strings.Contains(contextLower, "ecommerce") {
			return "Retail"
		}
	}

	return ""
}

// filterOfferingsForChat filters service offerings relevant to chat context
func (c *ChatAwareBedrockService) filterOfferingsForChat(offerings []*interfaces.ServiceOffering, serviceType, message string) []*interfaces.ServiceOffering {
	var relevant []*interfaces.ServiceOffering
	messageLower := strings.ToLower(message)

	for _, offering := range offerings {
		// Check if offering matches service type or message content
		if c.isServiceMatch(offering, serviceType) ||
			strings.Contains(strings.ToLower(offering.Name), messageLower) ||
			strings.Contains(strings.ToLower(offering.Description), messageLower) {
			relevant = append(relevant, offering)
		}
	}

	// Limit to top 3 for chat context
	if len(relevant) > 3 {
		relevant = relevant[:3]
	}

	return relevant
}

// getCachedResponse retrieves a cached response for similar content
func (c *ChatAwareBedrockService) getCachedResponse(content string) *CachedResponse {
	// Simple cache key based on content hash (in production, use proper hashing)
	cacheKey := fmt.Sprintf("%x", len(content)) // Simplified for demo

	if cached, exists := c.responseCache[cacheKey]; exists {
		if time.Now().Before(cached.ExpiresAt) {
			return cached
		}
		// Remove expired cache entry
		delete(c.responseCache, cacheKey)
	}

	return nil
}

// cacheResponse caches a response for future use
func (c *ChatAwareBedrockService) cacheResponse(content string, response *ChatResponse) {
	cacheKey := fmt.Sprintf("%x", len(content)) // Simplified for demo

	cached := &CachedResponse{
		Content:    response.Content,
		Metadata:   response.Metadata,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(1 * time.Hour), // Cache for 1 hour
		TokensUsed: response.TokensUsed,
	}

	c.responseCache[cacheKey] = cached
}

// optimizePrompt optimizes the prompt to reduce token usage
func (c *ChatAwareBedrockService) optimizePrompt(prompt string) string {
	// Remove excessive whitespace
	optimized := strings.TrimSpace(prompt)

	// Replace multiple consecutive spaces/newlines with single ones
	spaceRegex := regexp.MustCompile(`\s+`)
	optimized = spaceRegex.ReplaceAllString(optimized, " ")

	newlineRegex := regexp.MustCompile(`\n\s*\n\s*\n+`)
	optimized = newlineRegex.ReplaceAllString(optimized, "\n\n")

	// Truncate if too long (keep within reasonable token limits)
	maxLength := 4000 // Approximate token limit
	if len(optimized) > maxLength {
		optimized = optimized[:maxLength] + "..."
	}

	return optimized
}

// validateResponseQuality validates the quality of an AI response
func (c *ChatAwareBedrockService) validateResponseQuality(response string) error {
	if len(strings.TrimSpace(response)) < 10 {
		return fmt.Errorf("response too short")
	}

	// Check for common error patterns
	responseLower := strings.ToLower(response)
	errorPatterns := []string{
		"i cannot",
		"i don't know",
		"i'm not sure",
		"error occurred",
		"something went wrong",
	}

	for _, pattern := range errorPatterns {
		if strings.Contains(responseLower, pattern) {
			return fmt.Errorf("response contains error pattern: %s", pattern)
		}
	}

	return nil
}

// getFallbackResponse provides a fallback response for common scenarios
func (c *ChatAwareBedrockService) getFallbackResponse(content string) string {
	contentLower := strings.ToLower(content)

	if strings.Contains(contentLower, "cost") {
		return "For accurate cost estimates, I recommend using the AWS Pricing Calculator. I can help you identify the specific services and configurations needed for your use case."
	}

	if strings.Contains(contentLower, "security") {
		return "AWS provides comprehensive security services including IAM, GuardDuty, and Security Hub. I can help you design a security architecture that meets your specific requirements."
	}

	if strings.Contains(contentLower, "migration") {
		return "AWS offers several migration tools and services including Application Migration Service and Database Migration Service. I can help you plan a migration strategy based on your current infrastructure."
	}

	return "I'd be happy to help you with your AWS question. Could you provide more specific details about your requirements or use case?"
}

// Supporting types for chat-aware Bedrock service

// ChatRequest represents a request for chat-aware AI response
type ChatRequest struct {
	Content     string                 `json:"content"`
	Session     *domain.ChatSession    `json:"session"`
	Context     *domain.SessionContext `json:"context"`
	QuickAction string                 `json:"quick_action,omitempty"`
	History     []*domain.ChatMessage  `json:"history,omitempty"`
}

// ChatResponse represents a chat-aware AI response
type ChatResponse struct {
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	TokensUsed  int                    `json:"tokens_used"`
	ProcessTime float64                `json:"process_time"`
}

// ResponseOptimizer handles AI response optimization
type ResponseOptimizer struct {
	cache           *ResponseCache
	promptOptimizer *PromptOptimizer
	qualityFilter   *ResponseQualityFilter
}

// ResponseCache manages caching of AI responses
type ResponseCache struct {
	cache     map[string]*CachedResponse
	maxSize   int
	ttl       time.Duration
	hitCount  int64
	missCount int64
}

// PromptOptimizer optimizes prompts to reduce token usage
type PromptOptimizer struct {
	maxTokens      int
	compressionMap map[string]string
}

// ResponseQualityFilter validates and filters AI responses
type ResponseQualityFilter struct {
	minLength        int
	maxLength        int
	bannedPatterns   []string
	requiredPatterns []string
}

// NewResponseOptimizer creates a new response optimizer
func NewResponseOptimizer() *ResponseOptimizer {
	return &ResponseOptimizer{
		cache: &ResponseCache{
			cache:   make(map[string]*CachedResponse),
			maxSize: 1000,
			ttl:     2 * time.Hour,
		},
		promptOptimizer: &PromptOptimizer{
			maxTokens: 4000,
			compressionMap: map[string]string{
				"Amazon Web Services":            "AWS",
				"Elastic Compute Cloud":          "EC2",
				"Simple Storage Service":         "S3",
				"Relational Database Service":    "RDS",
				"Virtual Private Cloud":          "VPC",
				"Identity and Access Management": "IAM",
				"CloudFormation":                 "CFN",
				"Application Load Balancer":      "ALB",
				"Network Load Balancer":          "NLB",
				"Auto Scaling Group":             "ASG",
			},
		},
		qualityFilter: &ResponseQualityFilter{
			minLength: 20,
			maxLength: 2000,
			bannedPatterns: []string{
				"I cannot",
				"I don't know",
				"I'm not sure",
				"I apologize",
				"I'm sorry",
				"error occurred",
				"something went wrong",
				"unable to process",
			},
			requiredPatterns: []string{
				// At least one AWS service should be mentioned
				"AWS|Amazon|EC2|S3|RDS|Lambda|CloudFormation|VPC|IAM",
			},
		},
	}
}

// OptimizeResponse optimizes an AI response using caching, prompt optimization, and quality filtering
func (c *ChatAwareBedrockService) OptimizeResponse(ctx context.Context, request *ChatRequest) (*ChatResponse, error) {
	optimizer := NewResponseOptimizer()

	// Step 1: Check cache for similar queries
	cacheKey := optimizer.generateCacheKey(request)
	if cachedResponse := optimizer.cache.Get(cacheKey); cachedResponse != nil {
		return &ChatResponse{
			Content:     cachedResponse.Content,
			Metadata:    cachedResponse.Metadata,
			TokensUsed:  cachedResponse.TokensUsed,
			ProcessTime: 0, // Cached response is instant
		}, nil
	}

	// Step 2: Optimize prompt to reduce token usage
	optimizedPrompt, err := c.GetContextualPrompt(ctx, request.Session, request.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate contextual prompt: %w", err)
	}

	optimizedPrompt = optimizer.promptOptimizer.Optimize(optimizedPrompt)

	// Step 3: Generate response with optimized settings
	// Use the configured model ID from the base service
	modelInfo := c.EnhancedBedrockService.GetModelInfo()
	options := &interfaces.BedrockOptions{
		ModelID:     modelInfo.ModelID,
		MaxTokens:   800, // Reduced from 1000 for optimization
		Temperature: 0.7,
		TopP:        0.9,
	}

	startTime := time.Now()
	response, err := c.EnhancedBedrockService.GenerateText(ctx, optimizedPrompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}
	processingTime := time.Since(startTime)

	// Step 4: Validate and filter response quality
	if err := optimizer.qualityFilter.Validate(response.Content); err != nil {
		// Try fallback response
		fallbackContent := c.getFallbackResponse(request.Content)
		if fallbackContent != "" {
			response.Content = fallbackContent
		} else {
			return nil, fmt.Errorf("response quality validation failed: %w", err)
		}
	}

	chatResponse := &ChatResponse{
		Content: response.Content,
		Metadata: map[string]interface{}{
			"model_id":          options.ModelID,
			"processing_time":   processingTime.Milliseconds(),
			"prompt_tokens":     response.Usage.InputTokens,
			"completion_tokens": response.Usage.OutputTokens,
			"optimized":         true,
			"cache_hit":         false,
		},
		TokensUsed:  response.Usage.OutputTokens,
		ProcessTime: float64(processingTime.Milliseconds()),
	}

	// Step 5: Cache the optimized response
	optimizer.cache.Set(cacheKey, &CachedResponse{
		Content:    chatResponse.Content,
		Metadata:   chatResponse.Metadata,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(optimizer.cache.ttl),
		TokensUsed: chatResponse.TokensUsed,
	})

	return chatResponse, nil
}

// ResponseCache methods

// generateCacheKey generates a cache key based on request content and context
func (ro *ResponseOptimizer) generateCacheKey(request *ChatRequest) string {
	// Create a hash-like key based on content and key context elements
	key := fmt.Sprintf("%s|%s|%s",
		request.Content,
		request.Session.ClientName,
		request.QuickAction,
	)

	// Add service types if available
	if request.Context != nil && len(request.Context.ServiceTypes) > 0 {
		key += "|" + strings.Join(request.Context.ServiceTypes, ",")
	}

	// Simple hash simulation (in production, use proper hashing like SHA256)
	return fmt.Sprintf("chat_%x", len(key)%10000)
}

// Get retrieves a cached response
func (rc *ResponseCache) Get(key string) *CachedResponse {
	if cached, exists := rc.cache[key]; exists {
		if time.Now().Before(cached.ExpiresAt) {
			rc.hitCount++
			return cached
		}
		// Remove expired entry
		delete(rc.cache, key)
	}

	rc.missCount++
	return nil
}

// Set stores a response in cache
func (rc *ResponseCache) Set(key string, response *CachedResponse) {
	// Implement LRU eviction if cache is full
	if len(rc.cache) >= rc.maxSize {
		rc.evictOldest()
	}

	rc.cache[key] = response
}

// evictOldest removes the oldest cache entry
func (rc *ResponseCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range rc.cache {
		if oldestKey == "" || cached.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.CreatedAt
		}
	}

	if oldestKey != "" {
		delete(rc.cache, oldestKey)
	}
}

// GetStats returns cache statistics
func (rc *ResponseCache) GetStats() map[string]interface{} {
	total := rc.hitCount + rc.missCount
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(rc.hitCount) / float64(total)
	}

	return map[string]interface{}{
		"size":       len(rc.cache),
		"max_size":   rc.maxSize,
		"hit_count":  rc.hitCount,
		"miss_count": rc.missCount,
		"hit_rate":   hitRate,
	}
}

// PromptOptimizer methods

// Optimize optimizes a prompt to reduce token usage
func (po *PromptOptimizer) Optimize(prompt string) string {
	optimized := prompt

	// Step 1: Apply compression mappings
	for longForm, shortForm := range po.compressionMap {
		optimized = strings.ReplaceAll(optimized, longForm, shortForm)
	}

	// Step 2: Remove redundant phrases
	redundantPhrases := []string{
		"As an AI assistant, ",
		"I understand that ",
		"Let me help you with ",
		"Based on my knowledge, ",
		"In my experience, ",
	}

	for _, phrase := range redundantPhrases {
		optimized = strings.ReplaceAll(optimized, phrase, "")
	}

	// Step 3: Compress whitespace
	spaceRegex := regexp.MustCompile(`\s+`)
	optimized = spaceRegex.ReplaceAllString(optimized, " ")

	// Step 4: Remove excessive newlines
	newlineRegex := regexp.MustCompile(`\n\s*\n\s*\n+`)
	optimized = newlineRegex.ReplaceAllString(optimized, "\n\n")

	// Step 5: Truncate if still too long
	if len(optimized) > po.maxTokens {
		// Find a good breaking point (end of sentence)
		truncateAt := po.maxTokens
		for i := po.maxTokens - 100; i < po.maxTokens && i < len(optimized); i++ {
			if optimized[i] == '.' || optimized[i] == '!' || optimized[i] == '?' {
				truncateAt = i + 1
				break
			}
		}
		optimized = optimized[:truncateAt]
	}

	return strings.TrimSpace(optimized)
}

// GetCompressionStats returns statistics about prompt compression
func (po *PromptOptimizer) GetCompressionStats(original, optimized string) map[string]interface{} {
	originalLen := len(original)
	optimizedLen := len(optimized)
	compressionRatio := float64(optimizedLen) / float64(originalLen)

	return map[string]interface{}{
		"original_length":   originalLen,
		"optimized_length":  optimizedLen,
		"compression_ratio": compressionRatio,
		"tokens_saved":      originalLen - optimizedLen,
	}
}

// ResponseQualityFilter methods

// Validate validates the quality of a response
func (rqf *ResponseQualityFilter) Validate(response string) error {
	trimmed := strings.TrimSpace(response)

	// Check length constraints
	if len(trimmed) < rqf.minLength {
		return fmt.Errorf("response too short: %d characters (minimum: %d)", len(trimmed), rqf.minLength)
	}

	if len(trimmed) > rqf.maxLength {
		return fmt.Errorf("response too long: %d characters (maximum: %d)", len(trimmed), rqf.maxLength)
	}

	// Check for banned patterns
	responseLower := strings.ToLower(response)
	for _, pattern := range rqf.bannedPatterns {
		if strings.Contains(responseLower, strings.ToLower(pattern)) {
			return fmt.Errorf("response contains banned pattern: %s", pattern)
		}
	}

	// Check for required patterns
	for _, pattern := range rqf.requiredPatterns {
		matched, err := regexp.MatchString(`(?i)`+pattern, response)
		if err != nil {
			continue // Skip invalid regex
		}
		if !matched {
			return fmt.Errorf("response missing required pattern: %s", pattern)
		}
	}

	return nil
}

// GetQualityScore calculates a quality score for a response
func (rqf *ResponseQualityFilter) GetQualityScore(response string) float64 {
	score := 1.0
	trimmed := strings.TrimSpace(response)

	// Length score (optimal range: 100-500 characters)
	length := len(trimmed)
	if length < 100 {
		score *= float64(length) / 100.0
	} else if length > 500 {
		score *= 500.0 / float64(length)
	}

	// AWS service mention score
	awsServices := []string{"AWS", "Amazon", "EC2", "S3", "RDS", "Lambda", "VPC", "IAM"}
	mentionCount := 0
	responseLower := strings.ToLower(response)
	for _, service := range awsServices {
		if strings.Contains(responseLower, strings.ToLower(service)) {
			mentionCount++
		}
	}

	if mentionCount > 0 {
		score *= 1.0 + (float64(mentionCount) * 0.1) // Bonus for AWS service mentions
	} else {
		score *= 0.5 // Penalty for no AWS mentions
	}

	// Specificity score (presence of numbers, specific terms)
	specificityPatterns := []string{`\$\d+`, `\d+%`, `\d+\s*(GB|TB|MB)`, `\d+\s*(hours?|days?|months?)`}
	for _, pattern := range specificityPatterns {
		if matched, _ := regexp.MatchString(pattern, response); matched {
			score *= 1.1 // Bonus for specific information
		}
	}

	// Cap the score at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// Enhanced fallback responses with categorization
func (c *ChatAwareBedrockService) getEnhancedFallbackResponse(content string, context *domain.SessionContext) string {
	contentLower := strings.ToLower(content)

	// Cost-related queries
	if strings.Contains(contentLower, "cost") || strings.Contains(contentLower, "price") || strings.Contains(contentLower, "budget") {
		response := "For accurate cost estimates, I recommend using the AWS Pricing Calculator (calculator.aws). "
		response += "Key cost factors include: compute instances, storage volume, data transfer, and managed services. "

		if context != nil && len(context.ServiceTypes) > 0 {
			response += fmt.Sprintf("For %s services, ", strings.Join(context.ServiceTypes, " and "))
		}

		response += "I can help you estimate costs for specific workloads or services."
		return response
	}

	// Security-related queries
	if strings.Contains(contentLower, "security") || strings.Contains(contentLower, "secure") || strings.Contains(contentLower, "compliance") {
		response := "AWS provides comprehensive security through IAM, GuardDuty, Security Hub, and Config. "
		response += "Key security practices include: least privilege access, encryption at rest/transit, network segmentation, and continuous monitoring. "

		if context != nil && context.MeetingType == "compliance_review" {
			response += "For compliance requirements, AWS offers SOC, PCI DSS, HIPAA, and FedRAMP certifications. "
		}

		response += "What specific security concerns should we address?"
		return response
	}

	// Migration-related queries
	if strings.Contains(contentLower, "migrate") || strings.Contains(contentLower, "migration") || strings.Contains(contentLower, "move") {
		response := "AWS migration follows the 6 R's strategy: Rehost, Replatform, Refactor, Repurchase, Retain, Retire. "
		response += "Key tools include Application Migration Service, Database Migration Service, and Migration Hub. "
		response += "I recommend starting with discovery and assessment using AWS Application Discovery Service. "
		response += "What type of workloads are you planning to migrate?"
		return response
	}

	// Architecture-related queries
	if strings.Contains(contentLower, "architecture") || strings.Contains(contentLower, "design") || strings.Contains(contentLower, "pattern") {
		response := "AWS Well-Architected Framework provides guidance across 6 pillars: Operational Excellence, Security, Reliability, "
		response += "Performance Efficiency, Cost Optimization, and Sustainability. "
		response += "Common patterns include multi-tier architectures, microservices, serverless, and event-driven designs. "
		response += "What specific architectural challenge are you addressing?"
		return response
	}

	// Performance-related queries
	if strings.Contains(contentLower, "performance") || strings.Contains(contentLower, "optimize") || strings.Contains(contentLower, "slow") {
		response := "AWS performance optimization involves right-sizing instances, using CloudFront CDN, implementing caching strategies, "
		response += "and leveraging auto-scaling. Key services include ElastiCache, CloudWatch, and X-Ray for monitoring. "
		response += "What specific performance issues are you experiencing?"
		return response
	}

	// Default fallback
	response := "I'm here to help with your AWS questions. "
	if context != nil && context.ClientName != "" {
		response += fmt.Sprintf("For %s, ", context.ClientName)
	}
	response += "I can provide guidance on architecture, migration, security, cost optimization, and best practices. "
	response += "Could you provide more specific details about your requirements?"

	return response
}
