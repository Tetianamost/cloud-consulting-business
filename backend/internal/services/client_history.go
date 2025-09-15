package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ClientHistoryService manages client engagement history and integration
type ClientHistoryService struct {
	knowledgeBase interfaces.KnowledgeBase
}

// NewClientHistoryService creates a new client history service
func NewClientHistoryService(kb interfaces.KnowledgeBase) *ClientHistoryService {
	return &ClientHistoryService{
		knowledgeBase: kb,
	}
}

// RecordEngagement records a new client engagement from an inquiry
func (c *ClientHistoryService) RecordEngagement(ctx context.Context, inquiry *domain.Inquiry) error {
	if inquiry.Company == "" {
		return nil // Skip if no company information
	}

	// Create engagement record from inquiry
	engagement := &interfaces.ClientEngagement{
		ID:                 fmt.Sprintf("eng-%s", inquiry.ID),
		ClientName:         inquiry.Company,
		Industry:           c.inferIndustry(inquiry.Company),
		ProjectName:        fmt.Sprintf("Inquiry: %s", inquiry.Name),
		ServiceType:        c.inferServiceType(inquiry.Services),
		StartDate:          inquiry.CreatedAt,
		Status:             "inquiry",
		TeamMembers:        []string{}, // To be assigned later
		CloudProviders:     c.inferCloudProviders(inquiry.Message),
		TechnologiesUsed:   c.extractTechnologies(inquiry.Message),
		ProjectValue:       0.0, // To be determined
		ClientSatisfaction: 0.0, // To be measured
		KeyChallenges:      c.extractChallenges(inquiry.Message),
		SolutionsProvided:  []string{}, // To be filled during engagement
		LessonsLearned:     []string{}, // To be captured post-engagement
		Deliverables:       c.inferDeliverables(inquiry.Services),
		SuccessMetrics:     make(map[string]interface{}),
		ReferenceAllowed:   false, // Default to false, to be confirmed
		CaseStudyAllowed:   false, // Default to false, to be confirmed
		Metadata: map[string]interface{}{
			"inquiry_id": inquiry.ID,
			"source":     inquiry.Source,
			"priority":   inquiry.Priority,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store engagement (in a real implementation, this would go to a database)
	// For now, we'll just log it
	fmt.Printf("Recording client engagement: %+v\n", engagement)

	return nil
}

// GetClientInsights returns insights about a client based on history
func (c *ClientHistoryService) GetClientInsights(ctx context.Context, clientName string) (*ClientInsights, error) {
	history, err := c.knowledgeBase.GetClientHistory(ctx, clientName)
	if err != nil {
		return nil, err
	}

	insights := &ClientInsights{
		ClientName:          clientName,
		TotalEngagements:    len(history),
		EngagementHistory:   history,
		PreferredServices:   c.analyzePreferredServices(history),
		AverageSatisfaction: c.calculateAverageSatisfaction(history),
		TotalValue:          c.calculateTotalValue(history),
		KeyRelationships:    c.identifyKeyRelationships(history),
		RecommendedApproach: c.recommendApproach(history),
		RiskFactors:         c.identifyRiskFactors(history),
		Opportunities:       c.identifyOpportunities(history),
		LastEngagement:      c.getLastEngagement(history),
	}

	return insights, nil
}

// GetRecommendedServices returns service recommendations based on client history and company offerings
func (c *ClientHistoryService) GetRecommendedServices(ctx context.Context, clientName string, industry string) ([]*interfaces.ServiceOffering, error) {
	// Get client history to understand past preferences
	history, err := c.knowledgeBase.GetClientHistory(ctx, clientName)
	if err != nil {
		return nil, err
	}

	// Get all available service offerings
	allServices, err := c.knowledgeBase.GetServiceOfferings(ctx)
	if err != nil {
		return nil, err
	}

	// Filter and rank services based on client history and industry
	var recommendations []*interfaces.ServiceOffering

	for _, service := range allServices {
		// Check if service is suitable for client's industry
		if c.isServiceSuitableForIndustry(service, industry) {
			// Check if client has used this service before
			hasUsedBefore := c.hasClientUsedService(history, service.ServiceType)

			// Prioritize services the client hasn't used but are complementary
			if !hasUsedBefore && c.isComplementaryService(history, service.ServiceType) {
				recommendations = append(recommendations, service)
			}
		}
	}

	return recommendations, nil
}

// GetSimilarClientPatterns returns patterns from similar clients
func (c *ClientHistoryService) GetSimilarClientPatterns(ctx context.Context, clientName string, industry string) ([]*interfaces.ProjectPattern, error) {
	// Create a mock inquiry to use the existing pattern matching
	mockInquiry := &domain.Inquiry{
		Company:  clientName,
		Services: []string{}, // Will be filled based on history
	}

	// Get client history to understand their service preferences
	history, err := c.knowledgeBase.GetClientHistory(ctx, clientName)
	if err != nil {
		return nil, err
	}

	// Extract services from history
	serviceMap := make(map[string]bool)
	for _, engagement := range history {
		serviceMap[string(engagement.ServiceType)] = true
	}

	for service := range serviceMap {
		mockInquiry.Services = append(mockInquiry.Services, service)
	}

	// Get similar project patterns
	patterns, err := c.knowledgeBase.GetSimilarProjects(ctx, mockInquiry)
	if err != nil {
		return nil, err
	}

	// Filter patterns by industry if specified
	if industry != "" {
		var filteredPatterns []*interfaces.ProjectPattern
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(pattern.Industry), strings.ToLower(industry)) {
				filteredPatterns = append(filteredPatterns, pattern)
			}
		}
		return filteredPatterns, nil
	}

	return patterns, nil
}

// GetRecommendedTeam returns team member recommendations based on client needs and past success
func (c *ClientHistoryService) GetRecommendedTeam(ctx context.Context, clientName string, serviceType string) ([]*interfaces.TeamExpertise, error) {
	// Get client history to see who worked successfully with this client before
	history, err := c.knowledgeBase.GetClientHistory(ctx, clientName)
	if err != nil {
		return nil, err
	}

	// Get all team expertise
	allExpertise, err := c.knowledgeBase.GetTeamExpertise(ctx)
	if err != nil {
		return nil, err
	}

	var recommendations []*interfaces.TeamExpertise

	// First, prioritize team members who worked successfully with this client before
	for _, expertise := range allExpertise {
		for _, engagement := range history {
			for _, teamMember := range engagement.TeamMembers {
				if teamMember == expertise.ConsultantID && engagement.ClientSatisfaction >= 8.0 {
					// Check if they have relevant expertise for the service type
					if c.hasRelevantExpertise(expertise, serviceType) {
						recommendations = append(recommendations, expertise)
						break
					}
				}
			}
		}
	}

	// If we don't have enough recommendations, add other qualified team members
	if len(recommendations) < 3 {
		for _, expertise := range allExpertise {
			// Skip if already recommended
			alreadyRecommended := false
			for _, rec := range recommendations {
				if rec.ConsultantID == expertise.ConsultantID {
					alreadyRecommended = true
					break
				}
			}

			if !alreadyRecommended && c.hasRelevantExpertise(expertise, serviceType) {
				recommendations = append(recommendations, expertise)
				if len(recommendations) >= 5 { // Limit to top 5 recommendations
					break
				}
			}
		}
	}

	return recommendations, nil
}

// UpdateEngagementStatus updates the status of an engagement
func (c *ClientHistoryService) UpdateEngagementStatus(ctx context.Context, engagementID string, status string) error {
	// In a real implementation, this would update the database
	fmt.Printf("Updating engagement %s status to: %s\n", engagementID, status)
	return nil
}

// AddEngagementFeedback adds feedback to an engagement
func (c *ClientHistoryService) AddEngagementFeedback(ctx context.Context, engagementID string, feedback *EngagementFeedback) error {
	// In a real implementation, this would update the database
	fmt.Printf("Adding feedback to engagement %s: %+v\n", engagementID, feedback)
	return nil
}

// Helper methods for data extraction and analysis

func (c *ClientHistoryService) inferIndustry(company string) string {
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

func (c *ClientHistoryService) inferServiceType(services []string) domain.ServiceType {
	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "migration") {
			return domain.ServiceTypeMigration
		}
		if strings.Contains(serviceLower, "architecture") {
			return domain.ServiceTypeArchitectureReview
		}
		if strings.Contains(serviceLower, "assessment") {
			return domain.ServiceTypeAssessment
		}
		if strings.Contains(serviceLower, "optimization") {
			return domain.ServiceTypeOptimization
		}
	}
	return domain.ServiceTypeGeneral
}

func (c *ClientHistoryService) inferCloudProviders(message string) []string {
	providers := []string{}
	messageLower := strings.ToLower(message)

	if strings.Contains(messageLower, "aws") || strings.Contains(messageLower, "amazon") {
		providers = append(providers, "AWS")
	}
	if strings.Contains(messageLower, "azure") || strings.Contains(messageLower, "microsoft") {
		providers = append(providers, "Azure")
	}
	if strings.Contains(messageLower, "gcp") || strings.Contains(messageLower, "google cloud") {
		providers = append(providers, "Google Cloud")
	}

	return providers
}

func (c *ClientHistoryService) extractTechnologies(message string) []string {
	technologies := []string{}
	messageLower := strings.ToLower(message)

	// Common technology keywords
	techKeywords := []string{
		"kubernetes", "docker", "terraform", "ansible", "jenkins",
		"lambda", "ec2", "s3", "rds", "dynamodb",
		"microservices", "serverless", "containers",
		"python", "java", "nodejs", "react", "angular",
	}

	for _, tech := range techKeywords {
		if strings.Contains(messageLower, tech) {
			technologies = append(technologies, tech)
		}
	}

	return technologies
}

func (c *ClientHistoryService) extractChallenges(message string) []string {
	challenges := []string{}
	messageLower := strings.ToLower(message)

	// Common challenge keywords
	challengeKeywords := map[string]string{
		"cost":        "Cost optimization",
		"performance": "Performance issues",
		"security":    "Security concerns",
		"scalability": "Scalability challenges",
		"compliance":  "Compliance requirements",
		"downtime":    "Downtime concerns",
		"legacy":      "Legacy system integration",
	}

	for keyword, challenge := range challengeKeywords {
		if strings.Contains(messageLower, keyword) {
			challenges = append(challenges, challenge)
		}
	}

	return challenges
}

func (c *ClientHistoryService) inferDeliverables(services []string) []string {
	deliverables := []string{}

	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "migration") {
			deliverables = append(deliverables, "Migration Strategy", "Implementation Plan", "Risk Assessment")
		}
		if strings.Contains(serviceLower, "architecture") {
			deliverables = append(deliverables, "Architecture Review", "Recommendations Report", "Best Practices Guide")
		}
		if strings.Contains(serviceLower, "assessment") {
			deliverables = append(deliverables, "Current State Assessment", "Gap Analysis", "Roadmap")
		}
	}

	return deliverables
}

// Analysis methods for client insights

func (c *ClientHistoryService) analyzePreferredServices(history []*interfaces.ClientEngagement) []string {
	serviceCount := make(map[string]int)

	for _, engagement := range history {
		serviceType := string(engagement.ServiceType)
		serviceCount[serviceType]++
	}

	// Return services sorted by frequency
	var preferred []string
	for service := range serviceCount {
		preferred = append(preferred, service)
	}

	return preferred
}

func (c *ClientHistoryService) calculateAverageSatisfaction(history []*interfaces.ClientEngagement) float64 {
	if len(history) == 0 {
		return 0.0
	}

	total := 0.0
	count := 0

	for _, engagement := range history {
		if engagement.ClientSatisfaction > 0 {
			total += engagement.ClientSatisfaction
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return total / float64(count)
}

func (c *ClientHistoryService) calculateTotalValue(history []*interfaces.ClientEngagement) float64 {
	total := 0.0
	for _, engagement := range history {
		total += engagement.ProjectValue
	}
	return total
}

func (c *ClientHistoryService) identifyKeyRelationships(history []*interfaces.ClientEngagement) []string {
	teamMembers := make(map[string]int)

	for _, engagement := range history {
		for _, member := range engagement.TeamMembers {
			teamMembers[member]++
		}
	}

	// Return team members who worked on multiple projects
	var keyRelationships []string
	for member, count := range teamMembers {
		if count > 1 {
			keyRelationships = append(keyRelationships, member)
		}
	}

	return keyRelationships
}

func (c *ClientHistoryService) recommendApproach(history []*interfaces.ClientEngagement) string {
	if len(history) == 0 {
		return "Standard engagement approach"
	}

	// Analyze past satisfaction and adjust approach
	avgSatisfaction := c.calculateAverageSatisfaction(history)

	if avgSatisfaction >= 8.0 {
		return "Continue with proven successful approach"
	} else if avgSatisfaction >= 6.0 {
		return "Standard approach with enhanced communication"
	} else {
		return "Careful approach with frequent check-ins and stakeholder alignment"
	}
}

func (c *ClientHistoryService) identifyRiskFactors(history []*interfaces.ClientEngagement) []string {
	risks := []string{}

	// Analyze past challenges
	challengeCount := make(map[string]int)
	for _, engagement := range history {
		for _, challenge := range engagement.KeyChallenges {
			challengeCount[challenge]++
		}
	}

	// Identify recurring challenges as risk factors
	for challenge, count := range challengeCount {
		if count > 1 {
			risks = append(risks, fmt.Sprintf("Recurring challenge: %s", challenge))
		}
	}

	return risks
}

func (c *ClientHistoryService) identifyOpportunities(history []*interfaces.ClientEngagement) []string {
	opportunities := []string{}

	if len(history) > 0 {
		// Look for expansion opportunities
		lastEngagement := c.getLastEngagement(history)
		if lastEngagement != nil && lastEngagement.ClientSatisfaction >= 8.0 {
			opportunities = append(opportunities, "High satisfaction - opportunity for expanded engagement")
		}

		// Look for service expansion
		serviceTypes := make(map[string]bool)
		for _, engagement := range history {
			serviceTypes[string(engagement.ServiceType)] = true
		}

		if len(serviceTypes) == 1 {
			opportunities = append(opportunities, "Single service type - opportunity to introduce additional services")
		}
	}

	return opportunities
}

func (c *ClientHistoryService) getLastEngagement(history []*interfaces.ClientEngagement) *interfaces.ClientEngagement {
	if len(history) == 0 {
		return nil
	}

	var latest *interfaces.ClientEngagement
	for _, engagement := range history {
		if latest == nil || engagement.StartDate.After(latest.StartDate) {
			latest = engagement
		}
	}

	return latest
}

// Helper methods for service recommendations

func (c *ClientHistoryService) isServiceSuitableForIndustry(service *interfaces.ServiceOffering, industry string) bool {
	if industry == "" {
		return true
	}

	for _, targetIndustry := range service.TargetIndustries {
		if strings.Contains(strings.ToLower(targetIndustry), strings.ToLower(industry)) {
			return true
		}
	}

	return false
}

func (c *ClientHistoryService) hasClientUsedService(history []*interfaces.ClientEngagement, serviceType domain.ServiceType) bool {
	for _, engagement := range history {
		if engagement.ServiceType == serviceType {
			return true
		}
	}
	return false
}

func (c *ClientHistoryService) isComplementaryService(history []*interfaces.ClientEngagement, serviceType domain.ServiceType) bool {
	// Define complementary service relationships
	complementaryServices := map[domain.ServiceType][]domain.ServiceType{
		domain.ServiceTypeMigration:          {domain.ServiceTypeOptimization, domain.ServiceTypeAssessment},
		domain.ServiceTypeAssessment:         {domain.ServiceTypeMigration, domain.ServiceTypeOptimization},
		domain.ServiceTypeOptimization:       {domain.ServiceTypeArchitectureReview, domain.ServiceTypeAssessment},
		domain.ServiceTypeArchitectureReview: {domain.ServiceTypeOptimization, domain.ServiceTypeMigration},
	}

	// Check if client has used services that are complementary to the proposed service
	for _, engagement := range history {
		if complementary, exists := complementaryServices[engagement.ServiceType]; exists {
			for _, comp := range complementary {
				if comp == serviceType {
					return true
				}
			}
		}
	}

	return false
}

func (c *ClientHistoryService) hasRelevantExpertise(expertise *interfaces.TeamExpertise, serviceType string) bool {
	// Check if consultant has relevant expertise areas
	serviceKeywords := map[string][]string{
		"migration":           {"migration", "aws", "azure", "cloud"},
		"assessment":          {"security", "assessment", "compliance", "architecture"},
		"optimization":        {"optimization", "cost", "performance", "devops"},
		"architecture_review": {"architecture", "design", "aws", "azure", "security"},
	}

	keywords, exists := serviceKeywords[serviceType]
	if !exists {
		return true // Default to true for unknown service types
	}

	// Check expertise areas
	for _, area := range expertise.ExpertiseAreas {
		areaLower := strings.ToLower(area)
		for _, keyword := range keywords {
			if strings.Contains(areaLower, keyword) {
				return true
			}
		}
	}

	// Check specializations
	for _, spec := range expertise.Specializations {
		specLower := strings.ToLower(spec.Area)
		for _, keyword := range keywords {
			if strings.Contains(specLower, keyword) {
				return true
			}
		}
	}

	return false
}

// Supporting types

// ClientInsights represents insights about a client based on engagement history
type ClientInsights struct {
	ClientName          string                         `json:"client_name"`
	TotalEngagements    int                            `json:"total_engagements"`
	EngagementHistory   []*interfaces.ClientEngagement `json:"engagement_history"`
	PreferredServices   []string                       `json:"preferred_services"`
	AverageSatisfaction float64                        `json:"average_satisfaction"`
	TotalValue          float64                        `json:"total_value"`
	KeyRelationships    []string                       `json:"key_relationships"`
	RecommendedApproach string                         `json:"recommended_approach"`
	RiskFactors         []string                       `json:"risk_factors"`
	Opportunities       []string                       `json:"opportunities"`
	LastEngagement      *interfaces.ClientEngagement   `json:"last_engagement"`
}

// EngagementFeedback represents feedback for an engagement
type EngagementFeedback struct {
	EngagementID       string                 `json:"engagement_id"`
	ClientSatisfaction float64                `json:"client_satisfaction"`
	Feedback           string                 `json:"feedback"`
	LessonsLearned     []string               `json:"lessons_learned"`
	SuccessMetrics     map[string]interface{} `json:"success_metrics"`
	ReferenceAllowed   bool                   `json:"reference_allowed"`
	CaseStudyAllowed   bool                   `json:"case_study_allowed"`
	Testimonial        string                 `json:"testimonial,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
}
