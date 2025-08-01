package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// CompanyKnowledgeIntegrationService integrates company-specific knowledge with AI responses
type CompanyKnowledgeIntegrationService struct {
	knowledgeBase interfaces.KnowledgeBase
	clientHistory *ClientHistoryService
}

// NewCompanyKnowledgeIntegrationService creates a new company knowledge integration service
func NewCompanyKnowledgeIntegrationService(kb interfaces.KnowledgeBase, ch *ClientHistoryService) *CompanyKnowledgeIntegrationService {
	return &CompanyKnowledgeIntegrationService{
		knowledgeBase: kb,
		clientHistory: ch,
	}
}

// CompanyContext represents company-specific context for AI responses
type CompanyContext struct {
	ServiceOfferings     []*interfaces.ServiceOffering     `json:"service_offerings"`
	TeamExpertise        []*interfaces.TeamExpertise       `json:"team_expertise"`
	PastSolutions        []*interfaces.PastSolution        `json:"past_solutions"`
	MethodologyTemplates []*interfaces.MethodologyTemplate `json:"methodology_templates"`
	ClientHistory        []*interfaces.ClientEngagement    `json:"client_history"`
	ProjectPatterns      []*interfaces.ProjectPattern      `json:"project_patterns"`
	ConsultingApproach   *interfaces.ConsultingApproach    `json:"consulting_approach"`
	PricingModels        []*interfaces.PricingModel        `json:"pricing_models"`
	RecommendedTeam      []*interfaces.TeamExpertise       `json:"recommended_team"`
	SimilarProjects      []*interfaces.ProjectPattern      `json:"similar_projects"`
}

// GetCompanyContextForInquiry returns company-specific context for an inquiry
func (c *CompanyKnowledgeIntegrationService) GetCompanyContextForInquiry(ctx context.Context, inquiry *domain.Inquiry) (*CompanyContext, error) {
	context := &CompanyContext{}

	// Get relevant service offerings based on inquiry services
	if len(inquiry.Services) > 0 {
		for _, service := range inquiry.Services {
			serviceType := c.mapServiceStringToType(service)
			offerings, err := c.getServiceOfferingsByType(ctx, serviceType)
			if err == nil {
				context.ServiceOfferings = append(context.ServiceOfferings, offerings...)
			}
		}
	} else {
		// Get all service offerings if no specific services mentioned
		offerings, err := c.knowledgeBase.GetServiceOfferings(ctx)
		if err == nil {
			context.ServiceOfferings = offerings
		}
	}

	// Get team expertise relevant to the inquiry
	if len(inquiry.Services) > 0 {
		for _, service := range inquiry.Services {
			expertise, err := c.knowledgeBase.GetExpertiseByArea(ctx, service)
			if err == nil {
				context.TeamExpertise = append(context.TeamExpertise, expertise...)
			}
		}
	}

	// Get client history if company name is provided
	if inquiry.Company != "" {
		history, err := c.knowledgeBase.GetClientHistory(ctx, inquiry.Company)
		if err == nil {
			context.ClientHistory = history
		}

		// Get recommended team based on client history
		if len(inquiry.Services) > 0 {
			recommendedTeam, err := c.clientHistory.GetRecommendedTeam(ctx, inquiry.Company, inquiry.Services[0])
			if err == nil {
				context.RecommendedTeam = recommendedTeam
			}
		}
	}

	// Get past solutions relevant to the inquiry
	industry := c.inferIndustryFromCompany(inquiry.Company)
	if len(inquiry.Services) > 0 {
		serviceType := c.mapServiceStringToType(inquiry.Services[0])
		solutions, err := c.knowledgeBase.GetPastSolutions(ctx, string(serviceType), industry)
		if err == nil {
			context.PastSolutions = solutions
		}
	}

	// Get similar project patterns
	patterns, err := c.knowledgeBase.GetSimilarProjects(ctx, inquiry)
	if err == nil {
		context.ProjectPatterns = patterns
		context.SimilarProjects = patterns
	}

	// Get methodology templates
	if len(inquiry.Services) > 0 {
		serviceType := c.mapServiceStringToType(inquiry.Services[0])
		templates, err := c.knowledgeBase.GetMethodologyTemplates(ctx, string(serviceType))
		if err == nil {
			context.MethodologyTemplates = templates
		}

		// Get consulting approach
		approach, err := c.knowledgeBase.GetConsultingApproach(ctx, string(serviceType))
		if err == nil {
			context.ConsultingApproach = approach
		}

		// Get pricing models
		pricing, err := c.knowledgeBase.GetPricingModels(ctx, string(serviceType))
		if err == nil {
			context.PricingModels = pricing
		}
	}

	return context, nil
}

// GenerateContextualPrompt creates a prompt that includes company-specific context
func (c *CompanyKnowledgeIntegrationService) GenerateContextualPrompt(ctx context.Context, inquiry *domain.Inquiry, basePrompt string) (string, error) {
	companyContext, err := c.GetCompanyContextForInquiry(ctx, inquiry)
	if err != nil {
		return basePrompt, err // Return base prompt if context retrieval fails
	}

	contextualPrompt := basePrompt + "\n\n"
	contextualPrompt += "## Company-Specific Context\n\n"

	// Add service offerings context
	if len(companyContext.ServiceOfferings) > 0 {
		contextualPrompt += "### Our Service Offerings:\n"
		for _, offering := range companyContext.ServiceOfferings {
			contextualPrompt += fmt.Sprintf("- **%s**: %s\n", offering.Name, offering.Description)
			contextualPrompt += fmt.Sprintf("  - Duration: %s\n", offering.TypicalDuration)
			contextualPrompt += fmt.Sprintf("  - Team Size: %s\n", offering.TeamSize)
			if len(offering.KeyBenefits) > 0 {
				contextualPrompt += fmt.Sprintf("  - Key Benefits: %s\n", strings.Join(offering.KeyBenefits, ", "))
			}
		}
		contextualPrompt += "\n"
	}

	// Add team expertise context
	if len(companyContext.TeamExpertise) > 0 {
		contextualPrompt += "### Our Team Expertise:\n"
		for _, expert := range companyContext.TeamExpertise {
			contextualPrompt += fmt.Sprintf("- **%s** (%s): %s\n", expert.ConsultantName, expert.Role, strings.Join(expert.ExpertiseAreas, ", "))
		}
		contextualPrompt += "\n"
	}

	// Add client history context
	if len(companyContext.ClientHistory) > 0 {
		contextualPrompt += "### Previous Engagements with This Client:\n"
		for _, engagement := range companyContext.ClientHistory {
			contextualPrompt += fmt.Sprintf("- **%s** (%s): %s\n", engagement.ProjectName, engagement.Status, strings.Join(engagement.SolutionsProvided, ", "))
			if engagement.ClientSatisfaction > 0 {
				contextualPrompt += fmt.Sprintf("  - Client Satisfaction: %.1f/10\n", engagement.ClientSatisfaction)
			}
		}
		contextualPrompt += "\n"
	}

	// Add past solutions context
	if len(companyContext.PastSolutions) > 0 {
		contextualPrompt += "### Relevant Past Solutions:\n"
		for _, solution := range companyContext.PastSolutions {
			contextualPrompt += fmt.Sprintf("- **%s**: %s\n", solution.Title, solution.Description)
			contextualPrompt += fmt.Sprintf("  - Industry: %s\n", solution.Industry)
			contextualPrompt += fmt.Sprintf("  - Technologies: %s\n", strings.Join(solution.Technologies, ", "))
			if solution.CostSavings > 0 {
				contextualPrompt += fmt.Sprintf("  - Cost Savings: $%.0f\n", solution.CostSavings)
			}
		}
		contextualPrompt += "\n"
	}

	// Add methodology context
	if companyContext.ConsultingApproach != nil {
		contextualPrompt += "### Our Consulting Approach:\n"
		contextualPrompt += fmt.Sprintf("- **Philosophy**: %s\n", companyContext.ConsultingApproach.Philosophy)
		contextualPrompt += fmt.Sprintf("- **Key Principles**: %s\n", strings.Join(companyContext.ConsultingApproach.KeyPrinciples, ", "))
		contextualPrompt += fmt.Sprintf("- **Engagement Model**: %s\n", companyContext.ConsultingApproach.EngagementModel)
		contextualPrompt += "\n"
	}

	// Add pricing context
	if len(companyContext.PricingModels) > 0 {
		contextualPrompt += "### Our Pricing Models:\n"
		for _, pricing := range companyContext.PricingModels {
			contextualPrompt += fmt.Sprintf("- **%s**: %s pricing starting at $%.0f\n", pricing.Name, pricing.PricingType, pricing.BasePrice)
		}
		contextualPrompt += "\n"
	}

	contextualPrompt += "## Instructions\n"
	contextualPrompt += "Please use the above company-specific context to provide responses that:\n"
	contextualPrompt += "1. Reference our specific service offerings and capabilities\n"
	contextualPrompt += "2. Mention relevant team members and their expertise\n"
	contextualPrompt += "3. Draw from our past successful solutions and client experiences\n"
	contextualPrompt += "4. Align with our consulting methodology and approach\n"
	contextualPrompt += "5. Provide realistic pricing estimates based on our models\n"
	contextualPrompt += "6. Demonstrate deep understanding of our company's unique value proposition\n\n"

	return contextualPrompt, nil
}

// GetRecommendationsForInquiry returns specific recommendations based on company knowledge
func (c *CompanyKnowledgeIntegrationService) GetRecommendationsForInquiry(ctx context.Context, inquiry *domain.Inquiry) (*InquiryRecommendations, error) {
	recommendations := &InquiryRecommendations{
		InquiryID:   inquiry.ID,
		GeneratedAt: time.Now(),
	}

	// Get service recommendations
	if inquiry.Company != "" {
		industry := c.inferIndustryFromCompany(inquiry.Company)
		services, err := c.clientHistory.GetRecommendedServices(ctx, inquiry.Company, industry)
		if err == nil {
			recommendations.RecommendedServices = services
		}

		// Get team recommendations
		if len(inquiry.Services) > 0 {
			team, err := c.clientHistory.GetRecommendedTeam(ctx, inquiry.Company, inquiry.Services[0])
			if err == nil {
				recommendations.RecommendedTeam = team
			}
		}

		// Get similar client patterns
		patterns, err := c.clientHistory.GetSimilarClientPatterns(ctx, inquiry.Company, industry)
		if err == nil {
			recommendations.SimilarPatterns = patterns
		}
	}

	// Get methodology recommendations
	if len(inquiry.Services) > 0 {
		serviceType := c.mapServiceStringToType(inquiry.Services[0])
		templates, err := c.knowledgeBase.GetMethodologyTemplates(ctx, string(serviceType))
		if err == nil {
			recommendations.MethodologyTemplates = templates
		}
	}

	return recommendations, nil
}

// Helper methods

func (c *CompanyKnowledgeIntegrationService) mapServiceStringToType(service string) domain.ServiceType {
	serviceLower := strings.ToLower(service)

	if strings.Contains(serviceLower, "migration") {
		return domain.ServiceTypeMigration
	}
	if strings.Contains(serviceLower, "assessment") || strings.Contains(serviceLower, "security") {
		return domain.ServiceTypeAssessment
	}
	if strings.Contains(serviceLower, "optimization") || strings.Contains(serviceLower, "cost") {
		return domain.ServiceTypeOptimization
	}
	if strings.Contains(serviceLower, "architecture") || strings.Contains(serviceLower, "review") {
		return domain.ServiceTypeArchitectureReview
	}

	return domain.ServiceTypeGeneral
}

func (c *CompanyKnowledgeIntegrationService) getServiceOfferingsByType(ctx context.Context, serviceType domain.ServiceType) ([]*interfaces.ServiceOffering, error) {
	allOfferings, err := c.knowledgeBase.GetServiceOfferings(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*interfaces.ServiceOffering
	for _, offering := range allOfferings {
		if offering.ServiceType == serviceType {
			filtered = append(filtered, offering)
		}
	}

	return filtered, nil
}

func (c *CompanyKnowledgeIntegrationService) inferIndustryFromCompany(company string) string {
	if company == "" {
		return ""
	}

	companyLower := strings.ToLower(company)

	// Industry keywords mapping
	industryKeywords := map[string][]string{
		"Financial Services": {"bank", "financial", "finance", "credit", "investment", "insurance"},
		"Healthcare":         {"health", "medical", "hospital", "clinic", "pharma", "biotech"},
		"Technology":         {"tech", "software", "digital", "data", "ai", "ml", "saas"},
		"Retail":             {"retail", "store", "shop", "commerce", "marketplace"},
		"Manufacturing":      {"manufacturing", "factory", "industrial", "production"},
		"Education":          {"education", "school", "university", "learning", "academic"},
		"Government":         {"gov", "government", "public", "municipal", "federal"},
		"Energy":             {"energy", "oil", "gas", "renewable", "solar", "wind"},
	}

	for industry, keywords := range industryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(companyLower, keyword) {
				return industry
			}
		}
	}

	return "Other"
}

// Supporting types

// InquiryRecommendations represents recommendations for an inquiry
type InquiryRecommendations struct {
	InquiryID            string                            `json:"inquiry_id"`
	RecommendedServices  []*interfaces.ServiceOffering     `json:"recommended_services"`
	RecommendedTeam      []*interfaces.TeamExpertise       `json:"recommended_team"`
	SimilarPatterns      []*interfaces.ProjectPattern      `json:"similar_patterns"`
	MethodologyTemplates []*interfaces.MethodologyTemplate `json:"methodology_templates"`
	EstimatedTimeline    string                            `json:"estimated_timeline"`
	EstimatedCost        string                            `json:"estimated_cost"`
	RiskFactors          []string                          `json:"risk_factors"`
	SuccessFactors       []string                          `json:"success_factors"`
	NextSteps            []string                          `json:"next_steps"`
	GeneratedAt          time.Time                         `json:"generated_at"`
}
