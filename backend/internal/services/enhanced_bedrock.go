package services

import (
	"context"
	"fmt"
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
