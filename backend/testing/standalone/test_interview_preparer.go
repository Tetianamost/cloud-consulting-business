package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockBedrockService for testing
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	// Simulate different responses based on prompt content
	var content string
	
	
	if contains(prompt, "interview guide") {
		content = `TITLE: Cloud Migration Discovery Interview for Healthcare Client

OBJECTIVE: Understand current infrastructure, compliance requirements, and migration goals for HIPAA-compliant cloud migration

ESTIMATED DURATION: 90 minutes

PREPARATION NOTES:
- Review HIPAA compliance requirements for cloud environments
- Prepare questions about current data handling processes
- Research healthcare-specific cloud solutions

SECTION 1: BUSINESS CONTEXT AND OBJECTIVES
Objective: Understand business drivers and success criteria
Expected Duration: 20 minutes
Questions:
1. [MUST-ASK] What are the primary business drivers for this cloud migration? (Type: business, Expected: strategic goals)
2. [SHOULD-ASK] What does success look like for this project? (Type: open, Expected: measurable outcomes)
3. [NICE-TO-ASK] How will this migration impact your patient care delivery? (Type: business, Expected: operational impact)

SECTION 2: CURRENT INFRASTRUCTURE AND CHALLENGES
Objective: Document existing environment and pain points
Expected Duration: 25 minutes
Questions:
1. [MUST-ASK] Can you describe your current IT infrastructure? (Type: technical, Expected: architecture overview)
2. [MUST-ASK] What are your biggest infrastructure challenges today? (Type: open, Expected: specific pain points)
3. [SHOULD-ASK] How do you currently handle data backups and disaster recovery? (Type: technical, Expected: process description)

SECTION 3: TECHNICAL REQUIREMENTS
Objective: Gather detailed technical specifications
Expected Duration: 25 minutes
Questions:
1. [MUST-ASK] What applications and systems need to be migrated? (Type: technical, Expected: application inventory)
2. [SHOULD-ASK] What are your performance requirements? (Type: technical, Expected: SLA specifications)
3. [SHOULD-ASK] Do you have any integration requirements with external systems? (Type: technical, Expected: integration points)

SECTION 4: COMPLIANCE AND SECURITY
Objective: Understand regulatory and security requirements
Expected Duration: 15 minutes
Questions:
1. [MUST-ASK] What compliance frameworks do you need to maintain? (Type: business, Expected: regulatory requirements)
2. [MUST-ASK] How do you currently handle PHI data security? (Type: technical, Expected: security controls)
3. [SHOULD-ASK] What audit requirements do you have? (Type: business, Expected: audit processes)

SECTION 5: TIMELINE AND BUDGET
Objective: Establish project constraints and expectations
Expected Duration: 5 minutes
Questions:
1. [MUST-ASK] What is your target timeline for this migration? (Type: business, Expected: project schedule)
2. [SHOULD-ASK] Do you have a budget range for this project? (Type: business, Expected: budget constraints)

FOLLOW-UP ACTIONS:
- Schedule technical deep-dive session with IT team
- Provide HIPAA compliance checklist for cloud migration
- Prepare detailed migration assessment and timeline`
	} else if contains(prompt, "Generate a focused set of interview questions") {
		content = `[MUST-ASK] What specific healthcare regulations must your cloud environment comply with? (Type: business, Expected: regulatory frameworks)
[MUST-ASK] How do you currently manage patient data access and audit trails? (Type: technical, Expected: access control processes)
[SHOULD-ASK] What are your data retention requirements for different types of healthcare records? (Type: business, Expected: retention policies)
[SHOULD-ASK] How do you handle data encryption for PHI both at rest and in transit? (Type: technical, Expected: encryption methods)
[SHOULD-ASK] What disaster recovery requirements do you have for patient-critical systems? (Type: technical, Expected: RTO/RPO specifications)
[NICE-TO-ASK] How do you currently handle software updates and patching in your environment? (Type: technical, Expected: maintenance processes)
[NICE-TO-ASK] What integration requirements do you have with electronic health record systems? (Type: technical, Expected: EHR integration needs)
[NICE-TO-ASK] How do you manage user access provisioning and deprovisioning? (Type: technical, Expected: identity management processes)`
	} else if contains(prompt, "discovery checklist") {
		content = `REQUIRED ARTIFACTS:
- Network Architecture Diagram: Current network topology and security zones (Type: diagram, Priority: high, Format: Visio/PDF, Source: IT team)
- Application Inventory: Complete list of applications and dependencies (Type: document, Priority: high, Format: Excel/CSV, Source: IT team)
- Data Classification Matrix: Types of data and sensitivity levels (Type: document, Priority: high, Format: Excel/Word, Source: compliance team)
- Current Security Policies: Existing security and compliance documentation (Type: document, Priority: medium, Format: PDF/Word, Source: security team)
- Backup and DR Procedures: Current backup and disaster recovery processes (Type: document, Priority: medium, Format: Word/PDF, Source: operations team)

TECHNICAL REQUIREMENTS TO GATHER:
- Current server specifications and utilization metrics
- Database sizes and performance requirements
- Network bandwidth and latency requirements
- Storage capacity and IOPS requirements
- Integration points and API dependencies
- Monitoring and alerting requirements

BUSINESS REQUIREMENTS TO GATHER:
- Compliance and regulatory requirements
- Business continuity requirements
- Budget constraints and approval processes
- Timeline and milestone expectations
- Success criteria and KPIs
- Stakeholder roles and responsibilities

COMPLIANCE REQUIREMENTS TO ASSESS:
- HIPAA compliance requirements and current controls
- State and federal healthcare regulations
- Data residency and sovereignty requirements
- Audit and reporting requirements
- Third-party vendor compliance requirements

ENVIRONMENT DETAILS TO DOCUMENT:
- Physical data center locations and specifications
- Current cloud services usage (if any)
- Software licensing and support agreements
- Change management and deployment processes
- User access patterns and peak usage times
- Geographic distribution of users and systems`
	} else if contains(prompt, "follow-up questions") {
		content = `[MUST-ASK] Can you provide more specific details about the compliance frameworks you mentioned? (Type: business, Expected: detailed regulatory requirements)
[SHOULD-ASK] What specific challenges are you facing with your current backup solution? (Type: technical, Expected: technical pain points)
[SHOULD-ASK] How many users would be affected by the migration, and what are their typical usage patterns? (Type: business, Expected: user impact analysis)
[NICE-TO-ASK] Are there any seasonal or cyclical patterns in your system usage that we should consider? (Type: technical, Expected: usage patterns)
[NICE-TO-ASK] What internal resources do you have available to support the migration project? (Type: business, Expected: resource availability)`
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
	return &interfaces.PromptTemplate{
		Name:     templateName,
		Template: "mock template",
	}, nil
}

func (m *MockPromptArchitect) RegisterTemplate(template *interfaces.PromptTemplate) error {
	return nil
}

func (m *MockPromptArchitect) ListTemplates() []string {
	return []string{"mock-template"}
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func main() {
	fmt.Println("Testing Interview Preparer Service...")

	// Create mock services
	bedrockService := &MockBedrockService{}
	promptArchitect := &MockPromptArchitect{}

	// Create the interview preparer service
	interviewPreparer := services.NewInterviewPreparerService(bedrockService, promptArchitect)

	// Create a test inquiry
	testInquiry := &domain.Inquiry{
		ID:       "test-inquiry-1",
		Name:     "John Smith",
		Email:    "john.smith@regionalhospital.com",
		Company:  "Regional Hospital",
		Services: []string{"migration", "assessment"},
		Message:  "We need to migrate our patient data systems to the cloud while maintaining HIPAA compliance. Our current infrastructure is aging and we're experiencing performance issues.",
		Status:   "new",
		Priority: "high",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	// Test 1: Generate Interview Guide
	fmt.Println("\n=== Test 1: Generate Interview Guide ===")
	guide, err := interviewPreparer.GenerateInterviewGuide(ctx, testInquiry)
	if err != nil {
		log.Fatalf("Failed to generate interview guide: %v", err)
	}

	fmt.Printf("Generated Interview Guide:\n")
	fmt.Printf("ID: %s\n", guide.ID)
	fmt.Printf("Title: %s\n", guide.Title)
	fmt.Printf("Objective: %s\n", guide.Objective)
	fmt.Printf("Estimated Duration: %s\n", guide.EstimatedDuration)
	fmt.Printf("Number of Sections: %d\n", len(guide.Sections))
	fmt.Printf("Preparation Notes: %d items\n", len(guide.PreparationNotes))
	fmt.Printf("Follow-up Actions: %d items\n", len(guide.FollowUpActions))

	// Print first section details
	if len(guide.Sections) > 0 {
		section := guide.Sections[0]
		fmt.Printf("\nFirst Section: %s\n", section.Title)
		fmt.Printf("  Objective: %s\n", section.Objective)
		fmt.Printf("  Duration: %s\n", section.ExpectedDuration)
		fmt.Printf("  Questions: %d\n", len(section.Questions))
		
		if len(section.Questions) > 0 {
			q := section.Questions[0]
			fmt.Printf("  First Question: %s\n", q.Text)
			fmt.Printf("    Priority: %s\n", q.Priority)
			fmt.Printf("    Type: %s\n", q.Type)
			fmt.Printf("    Category: %s\n", q.Category)
		}
	}

	// Test 2: Generate Question Set
	fmt.Println("\n=== Test 2: Generate Question Set ===")
	questionSet, err := interviewPreparer.GenerateQuestionSet(ctx, "migration", "healthcare")
	if err != nil {
		log.Fatalf("Failed to generate question set: %v", err)
	}

	fmt.Printf("Generated Question Set:\n")
	fmt.Printf("ID: %s\n", questionSet.ID)
	fmt.Printf("Category: %s\n", questionSet.Category)
	fmt.Printf("Industry: %s\n", questionSet.Industry)
	fmt.Printf("Number of Questions: %d\n", len(questionSet.Questions))

	for i, q := range questionSet.Questions {
		fmt.Printf("  %d. [%s] %s (Type: %s)\n", i+1, q.Priority, q.Text, q.Type)
	}

	// Test 3: Generate Discovery Checklist
	fmt.Println("\n=== Test 3: Generate Discovery Checklist ===")
	checklist, err := interviewPreparer.GenerateDiscoveryChecklist(ctx, "migration")
	if err != nil {
		log.Fatalf("Failed to generate discovery checklist: %v", err)
	}

	fmt.Printf("Generated Discovery Checklist:\n")
	fmt.Printf("Service Type: %s\n", checklist.ServiceType)
	fmt.Printf("Required Artifacts: %d\n", len(checklist.RequiredArtifacts))
	fmt.Printf("Technical Requirements: %d\n", len(checklist.TechnicalRequirements))
	fmt.Printf("Business Requirements: %d\n", len(checklist.BusinessRequirements))
	fmt.Printf("Compliance Requirements: %d\n", len(checklist.ComplianceRequirements))
	fmt.Printf("Environment Details: %d\n", len(checklist.EnvironmentDetails))

	// Print some artifacts
	fmt.Println("\nRequired Artifacts:")
	for i, artifact := range checklist.RequiredArtifacts {
		if i >= 3 { // Limit output
			break
		}
		fmt.Printf("  %d. %s: %s\n", i+1, artifact.Name, artifact.Description)
		fmt.Printf("     Type: %s, Priority: %s, Source: %s\n", artifact.Type, artifact.Priority, artifact.Source)
	}

	// Test 4: Generate Follow-up Questions
	fmt.Println("\n=== Test 4: Generate Follow-up Questions ===")
	responses := []interfaces.InterviewResponse{
		{
			QuestionID: "q1",
			Response:   "We need to comply with HIPAA and some state regulations, but I'm not sure about all the specific requirements.",
			Notes:      "Client seems uncertain about full compliance scope",
			Confidence: "medium",
		},
		{
			QuestionID: "q2",
			Response:   "Our current backup takes too long and sometimes fails, especially for our large database files.",
			Notes:      "Backup issues are a major pain point",
			Confidence: "high",
		},
	}

	followUpQuestions, err := interviewPreparer.GenerateFollowUpQuestions(ctx, responses)
	if err != nil {
		log.Fatalf("Failed to generate follow-up questions: %v", err)
	}

	fmt.Printf("Generated Follow-up Questions: %d\n", len(followUpQuestions))
	for i, q := range followUpQuestions {
		fmt.Printf("  %d. [%s] %s\n", i+1, q.Priority, q.Text)
		fmt.Printf("     Type: %s, Expected: %s\n", q.Type, q.ExpectedAnswerType)
	}

	// Test 5: JSON Serialization
	fmt.Println("\n=== Test 5: JSON Serialization ===")
	guideJSON, err := json.MarshalIndent(guide, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize interview guide: %v", err)
	}
	fmt.Printf("Interview Guide JSON (first 500 chars):\n%s...\n", string(guideJSON)[:min(500, len(guideJSON))])

	checklistJSON, err := json.MarshalIndent(checklist, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize discovery checklist: %v", err)
	}
	fmt.Printf("Discovery Checklist JSON (first 500 chars):\n%s...\n", string(checklistJSON)[:min(500, len(checklistJSON))])

	fmt.Println("\n=== All Tests Completed Successfully! ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}