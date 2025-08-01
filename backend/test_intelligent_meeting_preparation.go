package main

import (
	"context"
	"encoding/json"
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

2. Discovery (High Priority):
   Topic: Current Infrastructure Assessment
   Key Message: Let's understand your current state to build the right roadmap
   Supporting Points: Legacy system integration, Data migration strategies, Compliance mapping
   Context: During technical discovery phase

KEY QUESTIONS TO ASK:
1. [MUST-ASK] What are your current HIPAA compliance challenges? (Type: business, Expected: specific compliance gaps)
2. [MUST-ASK] How do you currently handle patient data backup and recovery? (Type: technical, Expected: current processes)
3. [SHOULD-ASK] What's driving the timeline for this cloud migration? (Type: open, Expected: business drivers)

POTENTIAL CHALLENGES:
- HIPAA compliance complexity
- Legacy system integration
- Staff training requirements

RECOMMENDED APPROACH:
Start with a comprehensive security and compliance assessment, then develop a phased migration plan that maintains HIPAA compliance throughout the process.

COMPETITOR INSIGHTS:
Competitor: AWS Healthcare
Market Position: Market leader
Strengths: Comprehensive services, HIPAA compliance
Weaknesses: Complex pricing, Steep learning curve
Differentiation Opportunity: Personalized healthcare expertise
Client Relevance: Direct competitor for this engagement

INDUSTRY CONTEXT:
Market Trends: Telemedicine growth, AI adoption, Interoperability focus
Regulatory Landscape: HIPAA, HITECH, State privacy laws
Technology Trends: FHIR adoption, Cloud-native applications, AI/ML integration
Common Challenges: Data security, Legacy modernization, Interoperability
Best Practices: Zero-trust security, Phased migrations, Staff training
Case Studies: Regional hospital migration, Telehealth platform deployment

PREPARATION CHECKLIST:
- Review latest HIPAA guidance
- Prepare healthcare case studies
- Research client's current technology stack`
	} else if contains(prompt, "question bank") {
		content = `CATEGORY: DISCOVERY QUESTIONS
Description: Questions to understand current state and requirements
Usage: discovery

1. [MUST-ASK] What patient data systems are currently in use? (Type: technical, Expected: system inventory)
2. [SHOULD-ASK] How do you currently ensure HIPAA compliance? (Type: business, Expected: compliance processes)
3. [NICE-TO-ASK] What are your disaster recovery procedures? (Type: technical, Expected: DR processes)

CATEGORY: VALIDATION QUESTIONS  
Description: Questions to validate assumptions and requirements
Usage: validation

1. [MUST-ASK] Have you experienced any security incidents in the past year? (Type: business, Expected: security history)
2. [SHOULD-ASK] What's your budget range for this cloud migration? (Type: business, Expected: budget constraints)

CATEGORY: OBJECTION HANDLING QUESTIONS
Description: Questions to address potential objections and concerns
Usage: objection_handling

1. [MUST-ASK] What are your main concerns about moving to the cloud? (Type: open, Expected: specific concerns)
2. [SHOULD-ASK] How do you evaluate the ROI of technology investments? (Type: business, Expected: ROI criteria)

CATEGORY: TECHNICAL DEEP-DIVE QUESTIONS
Description: Questions for technical stakeholders and detailed requirements
Usage: discovery

1. [MUST-ASK] What's your current network architecture? (Type: technical, Expected: network topology)
2. [SHOULD-ASK] How do you handle data encryption currently? (Type: technical, Expected: encryption methods)

CATEGORY: BUSINESS IMPACT QUESTIONS
Description: Questions to understand business impact and ROI expectations
Usage: validation

1. [MUST-ASK] What business outcomes are you hoping to achieve? (Type: business, Expected: success metrics)
2. [SHOULD-ASK] How will you measure the success of this migration? (Type: business, Expected: KPIs)`
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

Competitor: Microsoft Azure Healthcare
Market Share: 25%
Service Offerings: Virtual Machines, SQL Database, Healthcare Bot, FHIR services
Pricing Strategy: Hybrid licensing with enterprise discounts
Strengths: Enterprise integration, Office 365 synergy, AI capabilities
Weaknesses: Limited healthcare-specific features, Complex licensing
Client Overlap: Medium - focuses more on enterprise
Differentiation Gaps: Less healthcare industry specialization

MARKET POSITIONING:
Unique Value Proposition: Healthcare-specialized cloud consulting with deep industry expertise
Target Segments: Mid-market healthcare providers, Regional hospitals, Healthcare startups
Key Differentiators: Healthcare expertise, Personalized approach, Compliance specialization
Competitive Advantages: Industry knowledge, Proven healthcare implementations, Regulatory expertise
Market Opportunities: Telemedicine growth, AI adoption, Legacy modernization

DIFFERENTIATION STRATEGY:
Primary Differentiator: Deep healthcare industry expertise and compliance specialization
Secondary Differentiators: Personalized service, Proven healthcare case studies, Regulatory knowledge
Messaging Strategy: "Healthcare cloud experts who understand your unique challenges"
Proof Points: Healthcare certifications, Industry case studies, Compliance track record

COMPETITIVE RESPONSES:
Competitor Claim: "We're the market leader with the most comprehensive services"
Our Response: "We're the healthcare specialists with deep industry expertise and personalized service"
Supporting Evidence: Healthcare certifications, Industry case studies, Client testimonials
Timing: proactive

THREAT ASSESSMENT:
Threat Type: pricing
Description: Large cloud providers can offer lower prices due to scale
Severity: medium
Probability: high
Mitigation Strategy: Focus on value and specialized expertise rather than price competition
Monitoring Required: true

RECOMMENDED STRATEGY:
Position as the healthcare cloud specialists who understand the unique challenges of healthcare organizations. Focus on compliance expertise, industry knowledge, and personalized service rather than competing on price or breadth of services.`
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

2. ID: action-002
   Description: Schedule technical deep-dive session with IT team
   Owner: both
   Priority: high
   Due Date: 2024-02-10
   Status: pending
   Dependencies: IT team availability
   Notes: Include network architecture review

NEXT STEPS:
Step: Compliance Assessment
Description: Conduct comprehensive HIPAA compliance review
Timeline: Within 2 weeks
Prerequisites: Access to current systems and documentation
Deliverables: Compliance gap analysis report
Stakeholders: Compliance Officer, IT Director, Legal team

DELIVERABLES:
Name: HIPAA Compliance Assessment Report
Description: Comprehensive analysis of current compliance posture and recommendations
Format: document
Due Date: 2024-02-20
Owner: consultant
Requirements: Include gap analysis and remediation roadmap
Dependencies: Client system access and documentation

TIMELINE:
Immediate Actions (24 hours): Send follow-up email with meeting summary
Short Term Actions (1 week): Schedule technical deep-dive, Gather compliance documentation
Medium Term Actions (1 month): Complete compliance assessment, Develop migration roadmap
Long Term Actions (beyond 1 month): Begin phased migration implementation

MILESTONE EVENTS:
Name: Compliance Assessment Complete
Description: HIPAA compliance assessment and gap analysis completed
Target Date: 2024-02-20
Criteria: Report delivered and reviewed with client
Stakeholders: Compliance Officer, IT Director, Project team

RISK FLAGS:
- Client expressed concerns about data migration timeline
- Budget approval process may be complex
- Legacy system integration challenges identified

OPPORTUNITY FLAGS:
- Strong interest in AI/ML capabilities for patient analytics
- Potential for additional telemedicine platform work
- Client mentioned other facilities that might need similar services`
	}

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(prompt) / 4, // Rough estimate
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
	fmt.Println("Testing Intelligent Client Meeting Preparation System...")

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
		Message:  "We need to migrate our patient data systems to the cloud while maintaining HIPAA compliance. We're looking for expertise in healthcare cloud solutions.",
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

	briefingJSON, _ := json.MarshalIndent(briefing, "", "  ")
	fmt.Printf("Pre-Meeting Briefing:\n%s\n", briefingJSON)

	// Test 2: Generate Question Bank
	fmt.Println("\n=== Test 2: Question Bank ===")
	questionBank, err := interviewPreparer.GenerateQuestionBank(ctx, "healthcare", []string{"HIPAA compliance", "legacy system integration", "data security"})
	if err != nil {
		log.Fatalf("Failed to generate question bank: %v", err)
	}

	questionBankJSON, _ := json.MarshalIndent(questionBank, "", "  ")
	fmt.Printf("Question Bank:\n%s\n", questionBankJSON)

	// Test 3: Generate Competitive Landscape Analysis
	fmt.Println("\n=== Test 3: Competitive Landscape Analysis ===")
	competitiveAnalysis, err := interviewPreparer.GenerateCompetitiveLandscapeAnalysis(ctx, "healthcare", []string{"on-premises servers", "legacy EMR system", "basic backup solution"})
	if err != nil {
		log.Fatalf("Failed to generate competitive landscape analysis: %v", err)
	}

	competitiveJSON, _ := json.MarshalIndent(competitiveAnalysis, "", "  ")
	fmt.Printf("Competitive Landscape Analysis:\n%s\n", competitiveJSON)

	// Test 4: Generate Follow-Up Action Items
	fmt.Println("\n=== Test 4: Follow-Up Action Items ===")
	meetingNotes := "Client expressed strong interest in cloud migration but has concerns about HIPAA compliance. They currently use legacy EMR system and want to modernize. Budget is approved but timeline is aggressive. IT team is small but knowledgeable."

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

	actionItemsJSON, _ := json.MarshalIndent(actionItems, "", "  ")
	fmt.Printf("Follow-Up Action Items:\n%s\n", actionItemsJSON)

	fmt.Println("\n=== All Tests Completed Successfully! ===")

	// Summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Pre-Meeting Briefing: Generated with %d talking points and %d key questions\n",
		len(briefing.TalkingPoints), len(briefing.KeyQuestions))
	fmt.Printf("- Question Bank: Generated with %d categories\n", len(questionBank.QuestionCategories))
	fmt.Printf("- Competitive Analysis: Generated with %d competitors analyzed\n", len(competitiveAnalysis.CompetitorAnalysis))
	fmt.Printf("- Action Items: Generated with %d action items and %d deliverables\n",
		len(actionItems.ActionItems), len(actionItems.Deliverables))
}
