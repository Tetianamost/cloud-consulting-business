package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// InMemoryKnowledgeBase implements the KnowledgeBase interface with in-memory storage
type InMemoryKnowledgeBase struct {
	serviceOfferings       []*interfaces.ServiceOffering
	pricingModels          []*interfaces.PricingModel
	teamExpertise          []*interfaces.TeamExpertise
	clientHistory          []*interfaces.ClientEngagement
	pastSolutions          []*interfaces.PastSolution
	methodologyTemplates   []*interfaces.MethodologyTemplate
	consultingApproaches   []*interfaces.ConsultingApproach
	deliverableTemplates   []*interfaces.DeliverableTemplate
	bestPractices          []*interfaces.BestPractice
	complianceRequirements []*interfaces.ComplianceRequirement
	knowledgeItems         []*interfaces.KnowledgeItem
	stats                  *interfaces.KnowledgeStats
	mu                     sync.RWMutex
	lastUpdated            time.Time
}

// NewInMemoryKnowledgeBase creates a new in-memory knowledge base with initial data
func NewInMemoryKnowledgeBase() *InMemoryKnowledgeBase {
	return &InMemoryKnowledgeBase{
		serviceOfferings:       []*interfaces.ServiceOffering{},
		pricingModels:          []*interfaces.PricingModel{},
		teamExpertise:          []*interfaces.TeamExpertise{},
		clientHistory:          []*interfaces.ClientEngagement{},
		pastSolutions:          []*interfaces.PastSolution{},
		methodologyTemplates:   []*interfaces.MethodologyTemplate{},
		consultingApproaches:   []*interfaces.ConsultingApproach{},
		deliverableTemplates:   []*interfaces.DeliverableTemplate{},
		bestPractices:          []*interfaces.BestPractice{},
		complianceRequirements: []*interfaces.ComplianceRequirement{},
		knowledgeItems:         []*interfaces.KnowledgeItem{},
		stats:                  &interfaces.KnowledgeStats{},
		lastUpdated:            time.Now(),
	}
}

/*
Implement the KnowledgeBase interface from interfaces/knowledge.go.
All methods below are minimal stubs for demonstration and error-free compilation.
You should fill in real logic as needed for your application.
*/

// Service offerings and pricing
func (kb *InMemoryKnowledgeBase) GetServiceOfferings(ctx context.Context) ([]*interfaces.ServiceOffering, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.serviceOfferings, nil
}

func (kb *InMemoryKnowledgeBase) GetServiceOffering(ctx context.Context, id string) (*interfaces.ServiceOffering, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	for _, s := range kb.serviceOfferings {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("service offering not found")
}

func (kb *InMemoryKnowledgeBase) GetPricingModels(ctx context.Context, serviceType string) ([]*interfaces.PricingModel, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.pricingModels, nil
}

// Team expertise and specializations
func (kb *InMemoryKnowledgeBase) GetTeamExpertise(ctx context.Context) ([]*interfaces.TeamExpertise, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.teamExpertise, nil
}

func (kb *InMemoryKnowledgeBase) GetConsultantSpecializations(ctx context.Context, consultantID string) ([]*interfaces.Specialization, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	for _, t := range kb.teamExpertise {
		if t.ConsultantID == consultantID {
			return t.Specializations, nil
		}
	}
	return nil, errors.New("consultant not found")
}

func (kb *InMemoryKnowledgeBase) GetExpertiseByArea(ctx context.Context, area string) ([]*interfaces.TeamExpertise, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	var result []*interfaces.TeamExpertise
	for _, t := range kb.teamExpertise {
		for _, a := range t.ExpertiseAreas {
			if a == area {
				result = append(result, t)
				break
			}
		}
	}
	return result, nil
}

// Client history and past engagements
func (kb *InMemoryKnowledgeBase) GetClientHistory(ctx context.Context, clientName string) ([]*interfaces.ClientEngagement, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	var result []*interfaces.ClientEngagement
	for _, c := range kb.clientHistory {
		if c.ClientName == clientName {
			result = append(result, c)
		}
	}
	return result, nil
}

func (kb *InMemoryKnowledgeBase) GetPastSolutions(ctx context.Context, serviceType string, industry string) ([]*interfaces.PastSolution, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.pastSolutions, nil
}

func (kb *InMemoryKnowledgeBase) GetSimilarProjects(ctx context.Context, inquiry *domain.Inquiry) ([]*interfaces.ProjectPattern, error) {
	return nil, nil // Not implemented
}

// Methodology templates
func (kb *InMemoryKnowledgeBase) GetMethodologyTemplates(ctx context.Context, serviceType string) ([]*interfaces.MethodologyTemplate, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.methodologyTemplates, nil
}

func (kb *InMemoryKnowledgeBase) GetConsultingApproach(ctx context.Context, serviceType string) (*interfaces.ConsultingApproach, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	if len(kb.consultingApproaches) > 0 {
		return kb.consultingApproaches[0], nil
	}
	return nil, errors.New("not found")
}

func (kb *InMemoryKnowledgeBase) GetDeliverableTemplates(ctx context.Context, serviceType string) ([]*interfaces.DeliverableTemplate, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.deliverableTemplates, nil
}

// Knowledge base management
func (kb *InMemoryKnowledgeBase) UpdateKnowledgeBase(ctx context.Context) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()
	kb.lastUpdated = time.Now()
	return nil
}

func (kb *InMemoryKnowledgeBase) SearchKnowledge(ctx context.Context, query string, category string) ([]*interfaces.KnowledgeItem, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.knowledgeItems, nil
}

func (kb *InMemoryKnowledgeBase) GetKnowledgeStats(ctx context.Context) (*interfaces.KnowledgeStats, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.stats, nil
}

// Additional methods for report generation
func (kb *InMemoryKnowledgeBase) GetBestPractices(ctx context.Context, category string) ([]*interfaces.BestPractice, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.bestPractices, nil
}

func (kb *InMemoryKnowledgeBase) GetComplianceRequirements(ctx context.Context, framework string) ([]*interfaces.ComplianceRequirement, error) {
	kb.mu.RLock()
	defer kb.mu.RUnlock()
	return kb.complianceRequirements, nil
}
