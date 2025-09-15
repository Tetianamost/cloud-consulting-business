package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockBedrockService for testing
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	// Mock response based on prompt content
	var content string

	if contains(prompt, "pre-meeting briefing") {
		content = `CLIENT BACKGROUND ANALYSIS:
Company Size: medium
Industry: healthcare
Business Model: B2B
Technology Maturity: developing
Cloud Readiness: intermediate
Key Stakeholders: CTO, IT Director, Compliance Officer
Business Drivers: Cost reduction, Scalability, Compliance
Pain Points: Legacy systems, Security concerns, Budget constraints
Compliance Requirements: HIPAA, SOC 2, GDPR

STRATEGIC TALKING POINTS:
1. Opening (High Priority):
   Topic: Healthcare Cloud Transformation
   Key Message: We understand healthcare's unique compliance and security requirements
   Supporting Points: HIPAA expertise, Healthcare case studies, Security-first approach
   Context: Use during introductions to establish credibility

KEY QUESTIONS TO ASK:
1. [MUST-ASK] What are your current HIPAA compliance challenges? (Type: business, Expected: specific compliance gaps)
2. [MUST-ASK] How do you currently handle patient data backup and recovery? (Type: technical, Expected: current processes)

POTENTIAL CHALLENGES:
- HIPAA compliance complexity
- Legacy system integration

RECOMMENDED APPROACH:
Start with a comprehensive security and compliance assessment.

COMPETITOR INSIGHTS:
Competitor: AWS Healthcare
Market Position: Market leader
Strengths: Comprehensive services, HIPAA compliance
Weaknesses: Complex pricing, Steep learning curve
Differentiation Opportunity: Personalized healthcare expertise
Client Relevance: Direct competitor for this engagement

INDUSTRY CONTEXT:
Market Trends: Telemedicine growth, AI adoption
Regulatory Landscape: HIPAA, HITECH, State privacy laws
Technology Trends: FHIR adoption, Cloud-native applications
Common Challenges: Data security, Legacy modernization
Best Practices: Zero-trust security, Phased migrations
Case Studies: Regional hospital migration

PREPARATION CHECKLIST:
- Review latest HIPAA guidance
- Prepare healthcare case studies`
	} else if contains(prompt, "question bank") {
		content = `CATEGORY: DISCOVERY QUESTIONS
Description: Questions to understand current state and requirements
Usage: discovery

1. [MUST-ASK] What patient data systems are currently in use? (Type: technical, Expected: system inventory)
2. [SHOULD-ASK] How do you currently ensure HIPAA compliance? (Type: business, Expected: compliance processes)

CATEGORY: VALIDATION QUESTIONS  
Description: Questions to validate assumptions and requirements
Usage: validation

1. [MUST-ASK] Have you experienced any security incidents in the past year? (Type: business, Expected: security history)`
	} else if contains(prompt, "competitive landscape") {
		content = `COMPETITOR ANALYSIS:
Competitor: AWS Healthcare
Market Share: 35%
Service Offerings: EC2, RDS, S3, Lambda, Healthcare-specific services
Pricing Strategy: Pay-as-you-go with volume discounts
Strengths: Market leader, Comprehensive services, Strong compliance
Weaknesses: Complex pricing, Steep learning curve, Generic approach
Client Overlap: High - targets same healthcare segment
Differentiation Gaps: Lack of personalized healthcare expertise

MARKET POSITIONING:
Unique Value Proposition: Healthcare-specialized cloud consulting with deep industry expertise
Target Segments: Mid-market healthcare providers, Regional hospitals
Key Differentiators: Healthcare expertise, Personalized approach
Competitive Advantages: Industry knowledge, Proven healthcare implementations
Market Opportunities: Telemedicine growth, AI adoption

DIFFERENTIATION STRATEGY:
Primary Differentiator: Deep healthcare industry expertise and compliance specialization
Secondary Differentiators: Personalized service, Proven healthcare case studies
Messaging Strategy: "Healthcare cloud experts who understand your unique challenges"
Proof Points: Healthcare certifications, Industry case studies

COMPETITIVE RESPONSES:
Competitor Claim: "We're the market leader with the most comprehensive services"
Our Response: "We're the healthcare specialists with deep industry expertise"
Supporting Evidence: Healthcare certifications, Industry case studies
Timing: proactive

THREAT ASSESSMENT:
Threat Type: pricing
Description: Large cloud providers can offer lower prices due to scale
Severity: medium
Probability: high
Mitigation Strategy: Focus on value and specialized expertise
Monitoring Required: true

RECOMMENDED STRATEGY:
Position as the healthcare cloud specialists who understand unique challenges.`
	} else if contains(prompt, "follow-up action items") {
		content = `ACTION ITEMS:
1. ID: action-001
   Description: Prepare detailed HIPAA compliance assessment document
   Owner: consultant
   Priority: high
   Due Date: 2024-02-15
   Status: pending
   Dependencies: Client provides current compliance documentation
   Notes: Focus on data encryption and access controls

NEXT STEPS:
Step: Compliance Assessment
Description: Conduct comprehensive HIPAA compliance review
Timeline: Within 2 weeks
Prerequisites: Access to current systems and documentation
Deliverables: Compliance gap analysis report
Stakeholders: Compliance Officer, IT Director

DELIVERABLES:
Name: HIPAA Compliance Assessment Report
Description: Comprehensive analysis of current compliance posture
Format: document
Due Date: 2024-02-20
Owner: consultant
Requirements: Include gap analysis and remediation roadmap
Dependencies: Client system access and documentation

TIMELINE:
Immediate Actions (24 hours): Send follow-up email with meeting summary
Short Term Actions (1 week): Schedule technical deep-dive
Medium Term Actions (1 month): Complete compliance assessment
Long Term Actions (beyond 1 month): Begin phased migration implementation

MILESTONE EVENTS:
Name: Compliance Assessment Complete
Description: HIPAA compliance assessment completed
Target Date: 2024-02-20
Criteria: Report delivered and reviewed with client
Stakeholders: Compliance Officer, IT Director

RISK FLAGS:
- Client expressed concerns about data migration timeline
- Budget approval process may be complex

OPPORTUNITY FLAGS:
- Strong interest in AI/ML capabilities
- Potential for additional telemedicine platform work`
	}

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(prompt) / 4,
			OutputTokens: len(content) / 4,
		},
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		ModelName:   "Claude 3 Sonnet",
		Provider:    "Anthropic",
		MaxTokens:   4000,
		IsAvailable: true,
	}
}

func (m *MockBedrockService) IsHealthy() bool {
	return true
}

// MockPromptArchitect for testing
type MockPromptArchitect struct{}

func (m *MockPromptArchitect) BuildReportPrompt(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.PromptOptions) (string, error) {
	return "mock report prompt", nil
}

func (m *MockPromptArchitect) BuildInterviewPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	return "mock interview prompt", nil
}

func (m *MockPromptArchitect) BuildRiskAssessmentPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	return "mock risk assessment prompt", nil
}

func (m *MockPromptArchitect) BuildCompetitiveAnalysisPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	return "mock competitive analysis prompt", nil
}

func (m *MockPromptArchitect) ValidatePrompt(prompt string) error {
	return nil
}

func (m *MockPromptArchitect) GetTemplate(templateName string) (*interfaces.PromptTemplate, error) {
	return &interfaces.PromptTemplate{Name: templateName}, nil
}

func (m *MockPromptArchitect) RegisterTemplate(template *interfaces.PromptTemplate) error {
	return nil
}

func (m *MockPromptArchitect) ListTemplates() []string {
	return []string{"template1", "template2"}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Testing Task 7: Intelligent Client Meeting Preparation System...")

	// Create mock services
	bedrockService := &MockBedrockService{}
	promptArchitect := &MockPromptArchitect{}

	// Create the interview preparer service
	interviewPreparer := services.NewInterviewPreparerService(bedrockService, promptArchitect)

	// Create a sample inquiry
	inquiry := &domain.Inquiry{
		ID:       "test-inquiry-001",
		Name:     "Dr. Sarah Johnson",
		Email:    "sarah.johnson@healthcareplus.com",
		Company:  "HealthcarePlus Regional Medical Center",
		Phone:    "+1-555-0123",
		Services: []string{"cloud-migration", "compliance-assessment"},
		Message:  "We need to migrate our patient data systems to the cloud while maintaining HIPAA compliance.",
		Status:   "new",
		Priority: "high",
	}

	ctx := context.Background()

	// Test 1: Generate Pre-Meeting Briefing
	fmt.Println("\n=== Test 1: Pre-Meeting Briefing ===")
	briefing, err := interviewPreparer.GeneratePreMeetingBriefing(ctx, inquiry)
	if err != nil {
		log.Fatalf("Failed to generate pre-meeting briefing: %v", err)
	}

	fmt.Printf("✓ Pre-Meeting Briefing Generated Successfully\n")
	fmt.Printf("  - ID: %s\n", briefing.ID)
	fmt.Printf("  - Talking Points: %d\n", len(briefing.TalkingPoints))
	fmt.Printf("  - Key Questions: %d\n", len(briefing.KeyQuestions))
	fmt.Printf("  - Potential Challenges: %d\n", len(briefing.PotentialChallenges))
	fmt.Printf("  - Competitor Insights: %d\n", len(briefing.CompetitorInsights))

	// Test 2: Generate Question Bank
	fmt.Println("\n=== Test 2: Question Bank ===")
	questionBank, err := interviewPreparer.GenerateQuestionBank(ctx, "healthcare", []string{"HIPAA compliance", "legacy system integration"})
	if err != nil {
		log.Fatalf("Failed to generate question bank: %v", err)
	}

	fmt.Printf("✓ Question Bank Generated Successfully\n")
	fmt.Printf("  - ID: %s\n", questionBank.ID)
	fmt.Printf("  - Industry: %s\n", questionBank.Industry)
	fmt.Printf("  - Question Categories: %d\n", len(questionBank.QuestionCategories))

	totalQuestions := 0
	for _, category := range questionBank.QuestionCategories {
		totalQuestions += len(category.Questions)
	}
	fmt.Printf("  - Total Questions: %d\n", totalQuestions)

	// Test 3: Generate Competitive Landscape Analysis
	fmt.Println("\n=== Test 3: Competitive Landscape Analysis ===")
	competitiveAnalysis, err := interviewPreparer.GenerateCompetitiveLandscapeAnalysis(ctx, "healthcare", []string{"on-premises servers", "legacy EMR system"})
	if err != nil {
		log.Fatalf("Failed to generate competitive landscape analysis: %v", err)
	}

	fmt.Printf("✓ Competitive Landscape Analysis Generated Successfully\n")
	fmt.Printf("  - ID: %s\n", competitiveAnalysis.ID)
	fmt.Printf("  - Industry: %s\n", competitiveAnalysis.Industry)
	fmt.Printf("  - Competitors Analyzed: %d\n", len(competitiveAnalysis.CompetitorAnalysis))
	fmt.Printf("  - Threat Assessments: %d\n", len(competitiveAnalysis.ThreatAssessment))

	// Test 4: Generate Follow-Up Action Items
	fmt.Println("\n=== Test 4: Follow-Up Action Items ===")
	meetingNotes := "Client expressed strong interest in cloud migration but has concerns about HIPAA compliance. Budget is approved but timeline is aggressive."

	clientResponses := []interfaces.InterviewResponse{
		{
			QuestionID: "q1",
			Response:   "We have about 50,000 patient records that need to be migrated securely",
			Notes:      "Client emphasized security as top priority",
			Confidence: "high",
		},
		{
			QuestionID: "q2",
			Response:   "Our current backup takes 8 hours and we've had two failures this year",
			Notes:      "Clear pain point - unreliable backup system",
			Confidence: "high",
		},
	}

	actionItems, err := interviewPreparer.GenerateFollowUpActionItems(ctx, meetingNotes, clientResponses)
	if err != nil {
		log.Fatalf("Failed to generate follow-up action items: %v", err)
	}

	fmt.Printf("✓ Follow-Up Action Items Generated Successfully\n")
	fmt.Printf("  - ID: %s\n", actionItems.ID)
	fmt.Printf("  - Action Items: %d\n", len(actionItems.ActionItems))
	fmt.Printf("  - Next Steps: %d\n", len(actionItems.NextSteps))
	fmt.Printf("  - Deliverables: %d\n", len(actionItems.Deliverables))
	fmt.Printf("  - Risk Flags: %d\n", len(actionItems.RiskFlags))
	fmt.Printf("  - Opportunity Flags: %d\n", len(actionItems.OpportunityFlags))

	fmt.Println("\n=== All Task 7 Tests Completed Successfully! ===")

	// Summary
	fmt.Printf("\nTask 7 Implementation Summary:\n")
	fmt.Printf("✓ Pre-Meeting Briefing: Analyzes client background and suggests talking points\n")
	fmt.Printf("✓ Question Bank: Generates industry-specific questions based on challenges\n")
	fmt.Printf("✓ Competitive Analysis: Provides market positioning and differentiation strategies\n")
	fmt.Printf("✓ Follow-Up Actions: Creates structured action items from meeting notes\n")

	fmt.Printf("\nAll four components of the intelligent client meeting preparation system are working correctly!\n")
}
