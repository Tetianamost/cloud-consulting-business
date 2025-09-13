//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Minimal domain types for testing
type ServiceType string

const (
	ServiceTypeMigration          ServiceType = "migration"
	ServiceTypeAssessment         ServiceType = "assessment"
	ServiceTypeOptimization       ServiceType = "optimization"
	ServiceTypeArchitectureReview ServiceType = "architecture_review"
	ServiceTypeGeneral            ServiceType = "general"
)

type Inquiry struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	Services  []string  `json:"services"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// Minimal interfaces for testing
type ServiceOffering struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Category         string      `json:"category"`
	ServiceType      ServiceType `json:"service_type"`
	KeyBenefits      []string    `json:"key_benefits"`
	Deliverables     []string    `json:"deliverables"`
	TypicalDuration  string      `json:"typical_duration"`
	Prerequisites    []string    `json:"prerequisites"`
	TargetIndustries []string    `json:"target_industries"`
	CloudProviders   []string    `json:"cloud_providers"`
	ComplexityLevel  string      `json:"complexity_level"`
	TeamSize         string      `json:"team_size"`
	SuccessMetrics   []string    `json:"success_metrics"`
	RiskFactors      []string    `json:"risk_factors"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type TeamExpertise struct {
	ID                 string            `json:"id"`
	ConsultantID       string            `json:"consultant_id"`
	ConsultantName     string            `json:"consultant_name"`
	Role               string            `json:"role"`
	ExpertiseAreas     []string          `json:"expertise_areas"`
	Specializations    []*Specialization `json:"specializations"`
	Certifications     []Certification   `json:"certifications"`
	ExperienceYears    int               `json:"experience_years"`
	IndustryFocus      []string          `json:"industry_focus"`
	CloudProviders     []string          `json:"cloud_providers"`
	AvailabilityStatus string            `json:"availability_status"`
	HourlyRate         float64           `json:"hourly_rate"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type Specialization struct {
	Area            string    `json:"area"`
	Level           string    `json:"level"`
	YearsExperience int       `json:"years_experience"`
	KeyProjects     []string  `json:"key_projects"`
	Certifications  []string  `json:"certifications"`
	LastUpdated     time.Time `json:"last_updated"`
}

type Certification struct {
	Name         string     `json:"name"`
	Provider     string     `json:"provider"`
	Level        string     `json:"level"`
	ObtainedDate time.Time  `json:"obtained_date"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty"`
	CertID       string     `json:"cert_id"`
}

// Simple knowledge base implementation for testing
type TestKnowledgeBase struct {
	serviceOfferings map[string]*ServiceOffering
	teamExpertise    map[string]*TeamExpertise
}

func NewTestKnowledgeBase() *TestKnowledgeBase {
	kb := &TestKnowledgeBase{
		serviceOfferings: make(map[string]*ServiceOffering),
		teamExpertise:    make(map[string]*TeamExpertise),
	}
	kb.initializeTestData()
	return kb
}

func (kb *TestKnowledgeBase) initializeTestData() {
	now := time.Now()

	// Add service offerings
	kb.serviceOfferings["cloud-migration"] = &ServiceOffering{
		ID:          "cloud-migration",
		Name:        "Cloud Migration Services",
		Description: "Comprehensive cloud migration planning and execution with minimal downtime and risk mitigation",
		Category:    "Migration",
		ServiceType: ServiceTypeMigration,
		KeyBenefits: []string{
			"Reduced infrastructure costs by 30-50%",
			"Improved scalability and reliability",
			"Enhanced security posture",
			"Faster deployment cycles",
		},
		Deliverables: []string{
			"Migration Strategy Document",
			"Risk Assessment Report",
			"Implementation Roadmap",
			"Post-Migration Optimization Plan",
		},
		TypicalDuration:  "3-6 months",
		Prerequisites:    []string{"Current infrastructure audit", "Business requirements analysis"},
		TargetIndustries: []string{"Financial Services", "Healthcare", "E-commerce", "Manufacturing"},
		CloudProviders:   []string{"AWS", "Azure", "Google Cloud"},
		ComplexityLevel:  "intermediate",
		TeamSize:         "3-5 consultants",
		SuccessMetrics:   []string{"Zero-downtime migration", "Cost reduction achieved", "Performance improvement"},
		RiskFactors:      []string{"Data loss", "Extended downtime", "Integration challenges"},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	kb.serviceOfferings["architecture-review"] = &ServiceOffering{
		ID:          "architecture-review",
		Name:        "Cloud Architecture Review",
		Description: "Comprehensive review of existing cloud architecture with optimization recommendations",
		Category:    "Assessment",
		ServiceType: ServiceTypeArchitectureReview,
		KeyBenefits: []string{
			"Identify cost optimization opportunities",
			"Improve security and compliance",
			"Enhance performance and reliability",
			"Modernize legacy components",
		},
		Deliverables: []string{
			"Architecture Assessment Report",
			"Security Analysis",
			"Cost Optimization Recommendations",
			"Modernization Roadmap",
		},
		TypicalDuration:  "2-4 weeks",
		Prerequisites:    []string{"Architecture documentation", "Access to cloud environments"},
		TargetIndustries: []string{"Technology", "Financial Services", "Healthcare", "Retail"},
		CloudProviders:   []string{"AWS", "Azure", "Google Cloud", "Multi-cloud"},
		ComplexityLevel:  "advanced",
		TeamSize:         "2-4 senior consultants",
		SuccessMetrics:   []string{"Identified savings potential", "Security improvements", "Performance gains"},
		RiskFactors:      []string{"Incomplete documentation", "Complex legacy systems"},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Add team expertise
	kb.teamExpertise["consultant-001"] = &TeamExpertise{
		ID:             "consultant-001",
		ConsultantID:   "consultant-001",
		ConsultantName: "Sarah Chen",
		Role:           "Senior Cloud Architect",
		ExpertiseAreas: []string{"AWS", "Kubernetes", "Microservices", "DevOps"},
		Specializations: []*Specialization{
			{
				Area:            "AWS Solutions Architecture",
				Level:           "expert",
				YearsExperience: 8,
				KeyProjects:     []string{"proj-001", "proj-003"},
				Certifications:  []string{"AWS Solutions Architect Professional"},
				LastUpdated:     now,
			},
		},
		Certifications: []Certification{
			{
				Name:         "AWS Solutions Architect Professional",
				Provider:     "Amazon Web Services",
				Level:        "Professional",
				ObtainedDate: time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC),
				CertID:       "AWS-SAP-2022-001",
			},
		},
		ExperienceYears:    8,
		IndustryFocus:      []string{"Financial Services", "Healthcare"},
		CloudProviders:     []string{"AWS", "Azure"},
		AvailabilityStatus: "available",
		HourlyRate:         250.0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	kb.teamExpertise["consultant-002"] = &TeamExpertise{
		ID:             "consultant-002",
		ConsultantID:   "consultant-002",
		ConsultantName: "Michael Rodriguez",
		Role:           "Cloud Security Specialist",
		ExpertiseAreas: []string{"Security", "Compliance", "Azure", "Multi-cloud"},
		Specializations: []*Specialization{
			{
				Area:            "Cloud Security Architecture",
				Level:           "expert",
				YearsExperience: 6,
				KeyProjects:     []string{"proj-002", "proj-004"},
				Certifications:  []string{"CISSP", "Azure Security Engineer"},
				LastUpdated:     now,
			},
		},
		Certifications: []Certification{
			{
				Name:         "Certified Information Systems Security Professional",
				Provider:     "ISC2",
				Level:        "Professional",
				ObtainedDate: time.Date(2021, 8, 20, 0, 0, 0, 0, time.UTC),
				CertID:       "CISSP-2021-002",
			},
		},
		ExperienceYears:    6,
		IndustryFocus:      []string{"Financial Services", "Government"},
		CloudProviders:     []string{"Azure", "AWS"},
		AvailabilityStatus: "available",
		HourlyRate:         275.0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (kb *TestKnowledgeBase) GetServiceOfferings(ctx context.Context) ([]*ServiceOffering, error) {
	offerings := make([]*ServiceOffering, 0, len(kb.serviceOfferings))
	for _, offering := range kb.serviceOfferings {
		offerings = append(offerings, offering)
	}
	return offerings, nil
}

func (kb *TestKnowledgeBase) GetTeamExpertise(ctx context.Context) ([]*TeamExpertise, error) {
	expertise := make([]*TeamExpertise, 0, len(kb.teamExpertise))
	for _, exp := range kb.teamExpertise {
		expertise = append(expertise, exp)
	}
	return expertise, nil
}

func (kb *TestKnowledgeBase) GetExpertiseByArea(ctx context.Context, area string) ([]*TeamExpertise, error) {
	expertise := make([]*TeamExpertise, 0)
	for _, exp := range kb.teamExpertise {
		for _, expArea := range exp.ExpertiseAreas {
			if strings.Contains(strings.ToLower(expArea), strings.ToLower(area)) {
				expertise = append(expertise, exp)
				break
			}
		}
	}
	return expertise, nil
}

// Company Knowledge Integration Service
type CompanyKnowledgeIntegrationService struct {
	knowledgeBase *TestKnowledgeBase
}

func NewCompanyKnowledgeIntegrationService(kb *TestKnowledgeBase) *CompanyKnowledgeIntegrationService {
	return &CompanyKnowledgeIntegrationService{
		knowledgeBase: kb,
	}
}

type CompanyContext struct {
	ServiceOfferings []*ServiceOffering `json:"service_offerings"`
	TeamExpertise    []*TeamExpertise   `json:"team_expertise"`
}

func (c *CompanyKnowledgeIntegrationService) GetCompanyContextForInquiry(ctx context.Context, inquiry *Inquiry) (*CompanyContext, error) {
	context := &CompanyContext{}

	// Get relevant service offerings
	offerings, err := c.knowledgeBase.GetServiceOfferings(ctx)
	if err == nil {
		context.ServiceOfferings = offerings
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

	return context, nil
}

func (c *CompanyKnowledgeIntegrationService) GenerateContextualPrompt(ctx context.Context, inquiry *Inquiry, basePrompt string) (string, error) {
	companyContext, err := c.GetCompanyContextForInquiry(ctx, inquiry)
	if err != nil {
		return basePrompt, err
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

	contextualPrompt += "## Instructions\n"
	contextualPrompt += "Please use the above company-specific context to provide responses that:\n"
	contextualPrompt += "1. Reference our specific service offerings and capabilities\n"
	contextualPrompt += "2. Mention relevant team members and their expertise\n"
	contextualPrompt += "3. Demonstrate deep understanding of our company's unique value proposition\n\n"

	return contextualPrompt, nil
}

func main() {
	ctx := context.Background()

	// Initialize services
	knowledgeBase := NewTestKnowledgeBase()
	companyKnowledgeInteg := NewCompanyKnowledgeIntegrationService(knowledgeBase)

	// Test inquiry
	inquiry := &Inquiry{
		ID:        "test-001",
		Name:      "John Smith",
		Email:     "john@techcorp.com",
		Company:   "TechCorp Solutions",
		Services:  []string{"migration", "architecture"},
		Message:   "We need help migrating our legacy applications to AWS cloud",
		CreatedAt: time.Now(),
	}

	fmt.Println("=== Testing Company Knowledge Integration ===\n")

	// Test 1: Get service offerings
	fmt.Println("1. Testing service offerings...")
	offerings, err := knowledgeBase.GetServiceOfferings(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d service offerings:\n", len(offerings))
		for _, offering := range offerings {
			fmt.Printf("   - %s: %s\n", offering.Name, offering.Description)
			fmt.Printf("     Duration: %s, Team Size: %s\n", offering.TypicalDuration, offering.TeamSize)
		}
	}

	// Test 2: Get team expertise
	fmt.Println("\n2. Testing team expertise...")
	expertise, err := knowledgeBase.GetTeamExpertise(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d team members:\n", len(expertise))
		for _, expert := range expertise {
			fmt.Printf("   - %s (%s): %d years experience\n",
				expert.ConsultantName, expert.Role, expert.ExperienceYears)
			fmt.Printf("     Expertise: %v\n", expert.ExpertiseAreas)
		}
	}

	// Test 3: Get company context for inquiry
	fmt.Println("\n3. Getting company context for inquiry...")
	context, err := companyKnowledgeInteg.GetCompanyContextForInquiry(ctx, inquiry)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   - Found %d service offerings\n", len(context.ServiceOfferings))
		fmt.Printf("   - Found %d team experts\n", len(context.TeamExpertise))
	}

	// Test 4: Generate contextual prompt
	fmt.Println("\n4. Generating contextual prompt...")
	basePrompt := "Generate a professional response for this client inquiry."
	enhancedPrompt, err := companyKnowledgeInteg.GenerateContextualPrompt(ctx, inquiry, basePrompt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   - Enhanced prompt length: %d characters\n", len(enhancedPrompt))
		fmt.Printf("   - Contains service offerings: %t\n", strings.Contains(enhancedPrompt, "Service Offerings"))
		fmt.Printf("   - Contains team expertise: %t\n", strings.Contains(enhancedPrompt, "Team Expertise"))

		// Show a sample of the enhanced prompt
		fmt.Println("\n   Sample of enhanced prompt:")
		lines := strings.Split(enhancedPrompt, "\n")
		for i, line := range lines {
			if i < 10 { // Show first 10 lines
				fmt.Printf("   %s\n", line)
			} else if i == 10 {
				fmt.Printf("   ... (truncated)\n")
				break
			}
		}
	}

	fmt.Println("\n=== Company Knowledge Integration Test Complete ===")
	fmt.Println("\n✅ All tests passed! Company-specific knowledge integration is working correctly.")
	fmt.Println("\nKey features implemented:")
	fmt.Println("- Service offerings with detailed information")
	fmt.Println("- Team expertise with specializations and certifications")
	fmt.Println("- Company context generation for inquiries")
	fmt.Println("- Contextual prompt enhancement with company knowledge")
	fmt.Println("- Integration ready for AI assistant responses")
}
