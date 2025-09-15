package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/google/uuid"
)

// proposalGenerator implements the ProposalGenerator interface
type proposalGenerator struct {
	bedrockService     interfaces.BedrockService
	promptArchitect    interfaces.PromptArchitect
	knowledgeBase      interfaces.KnowledgeBase
	riskAssessor       interfaces.RiskAssessor
	multiCloudAnalyzer interfaces.MultiCloudAnalyzer
}

// NewProposalGenerator creates a new proposal generator instance
func NewProposalGenerator(
	bedrockService interfaces.BedrockService,
	promptArchitect interfaces.PromptArchitect,
	knowledgeBase interfaces.KnowledgeBase,
	riskAssessor interfaces.RiskAssessor,
	multiCloudAnalyzer interfaces.MultiCloudAnalyzer,
) interfaces.ProposalGenerator {
	return &proposalGenerator{
		bedrockService:     bedrockService,
		promptArchitect:    promptArchitect,
		knowledgeBase:      knowledgeBase,
		riskAssessor:       riskAssessor,
		multiCloudAnalyzer: multiCloudAnalyzer,
	}
}

// GenerateProposal generates a comprehensive proposal for the given inquiry
func (p *proposalGenerator) GenerateProposal(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.ProposalOptions) (*interfaces.Proposal, error) {
	// Create base proposal structure
	proposal := &interfaces.Proposal{
		ID:        uuid.New().String(),
		InquiryID: inquiry.ID,
		Title:     p.generateProposalTitle(inquiry),
		Status:    interfaces.ProposalStatusDraft,
		Version:   "1.0",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set expiration date (30 days from creation)
	expiresAt := time.Now().AddDate(0, 0, 30)
	proposal.ExpiresAt = &expiresAt

	// Generate executive summary
	executiveSummary, err := p.generateExecutiveSummary(ctx, inquiry, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate executive summary: %w", err)
	}
	proposal.ExecutiveSummary = executiveSummary

	// Generate problem statement
	problemStatement, err := p.generateProblemStatement(ctx, inquiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate problem statement: %w", err)
	}
	proposal.ProblemStatement = problemStatement

	// Generate proposed solution
	proposedSolution, err := p.generateProposedSolution(ctx, inquiry, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate proposed solution: %w", err)
	}
	proposal.ProposedSolution = proposedSolution

	// Generate project scope
	projectScope, err := p.generateProjectScope(ctx, inquiry, proposedSolution)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project scope: %w", err)
	}
	proposal.ProjectScope = projectScope

	// Generate timeline estimate if requested
	if options.IncludeTimeline {
		timelineEstimate, err := p.EstimateTimeline(ctx, inquiry, projectScope)
		if err == nil {
			proposal.TimelineEstimate = timelineEstimate
		}
	}

	// Generate resource estimate
	resourceEstimate, err := p.EstimateResources(ctx, inquiry, projectScope)
	if err == nil {
		proposal.ResourceEstimate = resourceEstimate
	}

	// Generate pricing recommendation if requested
	if options.IncludePricingBreakdown {
		pricingRecommendation, err := p.GeneratePricingRecommendation(ctx, inquiry, projectScope)
		if err == nil {
			proposal.PricingRecommendation = pricingRecommendation
		}
	}

	// Generate risk assessment if requested
	if options.IncludeRiskAssessment {
		riskAssessment, err := p.AssessProjectRisks(ctx, inquiry, projectScope)
		if err == nil {
			proposal.RiskAssessment = riskAssessment
		}
	}

	// Generate deliverables
	deliverables, err := p.generateDeliverables(ctx, inquiry, proposedSolution)
	if err == nil {
		proposal.Deliverables = deliverables
	}

	// Generate next steps
	proposal.NextSteps = p.generateNextSteps(inquiry, options)

	// Generate assumptions
	proposal.Assumptions = p.generateAssumptions(inquiry, proposedSolution)

	// Generate success metrics
	proposal.SuccessMetrics = p.generateSuccessMetrics(inquiry, proposedSolution)

	return proposal, nil
}

// GenerateSOW generates a detailed statement of work from a proposal
func (p *proposalGenerator) GenerateSOW(ctx context.Context, inquiry *domain.Inquiry, proposal *interfaces.Proposal) (*interfaces.StatementOfWork, error) {
	sow := &interfaces.StatementOfWork{
		ID:         uuid.New().String(),
		ProposalID: proposal.ID,
		InquiryID:  inquiry.ID,
		Title:      fmt.Sprintf("Statement of Work - %s", proposal.Title),
		Status:     interfaces.SOWStatusDraft,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Generate project overview
	projectOverview, err := p.generateProjectOverview(ctx, inquiry, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to generate project overview: %w", err)
	}
	sow.ProjectOverview = projectOverview

	// Generate detailed scope
	detailedScope, err := p.generateDetailedScope(ctx, inquiry, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to generate detailed scope: %w", err)
	}
	sow.Scope = detailedScope

	// Generate detailed deliverables
	detailedDeliverables, err := p.generateDetailedDeliverables(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to generate detailed deliverables: %w", err)
	}
	sow.Deliverables = detailedDeliverables

	// Generate detailed timeline
	if proposal.TimelineEstimate != nil {
		detailedTimeline, err := p.generateDetailedTimeline(ctx, proposal.TimelineEstimate)
		if err == nil {
			sow.Timeline = detailedTimeline
		}
	}

	// Generate detailed resource allocation
	if proposal.ResourceEstimate != nil {
		detailedResources, err := p.generateDetailedResources(ctx, proposal.ResourceEstimate)
		if err == nil {
			sow.ResourceAllocation = detailedResources
		}
	}

	// Generate payment schedule
	if proposal.PricingRecommendation != nil {
		paymentSchedule, err := p.generatePaymentSchedule(ctx, proposal.PricingRecommendation, proposal.TimelineEstimate)
		if err == nil {
			sow.PaymentSchedule = paymentSchedule
		}
	}

	// Generate acceptance criteria
	acceptanceCriteria, err := p.generateAcceptanceCriteria(ctx, proposal)
	if err == nil {
		sow.AcceptanceCriteria = acceptanceCriteria
	}

	return sow, nil
}

// EstimateTimeline estimates the project timeline based on similar projects
func (p *proposalGenerator) EstimateTimeline(ctx context.Context, inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) (*interfaces.TimelineEstimate, error) {
	// Get similar projects for comparison
	similarProjects, err := p.GetSimilarProjects(ctx, inquiry)
	if err != nil {
		similarProjects = []*interfaces.HistoricalProject{} // Continue with empty list
	}

	// Calculate base duration from similar projects
	baseDuration := p.calculateBaseDuration(inquiry, similarProjects)

	// Generate phases based on services requested
	phases, err := p.generateTimelinePhases(ctx, inquiry, projectScope, baseDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate timeline phases: %w", err)
	}

	// Generate milestones
	milestones := p.generateMilestones(phases)

	// Calculate total duration
	totalDuration := p.calculateTotalDuration(phases)

	// Determine confidence based on similar projects
	confidence := p.calculateTimelineConfidence(similarProjects, inquiry)

	// Generate critical path
	criticalPath := p.identifyCriticalPath(phases)

	// Calculate buffer time (15% of total duration)
	bufferTime := fmt.Sprintf("%.0f days", float64(totalDuration)*0.15)

	timeline := &interfaces.TimelineEstimate{
		TotalDuration:    fmt.Sprintf("%d days", totalDuration),
		Phases:           phases,
		Milestones:       milestones,
		CriticalPath:     criticalPath,
		BufferTime:       bufferTime,
		Confidence:       confidence,
		SimilarProjects:  p.convertToTimelineProjects(similarProjects),
		EstimationMethod: "Historical data analysis with complexity adjustments",
	}

	return timeline, nil
}

// EstimateResources estimates the resource requirements for the project
func (p *proposalGenerator) EstimateResources(ctx context.Context, inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) (*interfaces.ProposalResourceEstimate, error) {
	// Calculate total effort based on project complexity
	totalEffort := p.calculateTotalEffort(inquiry, projectScope)

	// Generate team composition based on services
	teamComposition := p.generateTeamComposition(inquiry, totalEffort)

	// Generate skill requirements
	skillRequirements := p.generateSkillRequirements(inquiry)

	// Generate external resource requirements
	externalResources := p.generateExternalResources(inquiry)

	// Generate tools and licenses requirements
	toolsAndLicenses := p.generateToolsAndLicenses(inquiry)

	// Generate training needs
	trainingNeeds := p.generateTrainingNeeds(inquiry, skillRequirements)

	// Calculate confidence based on project complexity
	confidence := p.calculateResourceConfidence(inquiry)

	resourceEstimate := &interfaces.ProposalResourceEstimate{
		TotalEffort:       totalEffort,
		TeamComposition:   teamComposition,
		SkillRequirements: skillRequirements,
		ExternalResources: externalResources,
		ToolsAndLicenses:  toolsAndLicenses,
		TrainingNeeds:     trainingNeeds,
		Confidence:        confidence,
	}

	return resourceEstimate, nil
}

// AssessProjectRisks assesses risks for the proposed project
func (p *proposalGenerator) AssessProjectRisks(ctx context.Context, inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) (*interfaces.ProjectRiskAssessment, error) {
	// Generate technical risks
	technicalRisks := p.generateTechnicalRisks(inquiry, projectScope)

	// Generate business risks
	businessRisks := p.generateBusinessRisks(inquiry, projectScope)

	// Generate resource risks
	resourceRisks := p.generateResourceRisks(inquiry, projectScope)

	// Generate timeline risks
	timelineRisks := p.generateTimelineRisks(inquiry, projectScope)

	// Generate budget risks
	budgetRisks := p.generateBudgetRisks(inquiry, projectScope)

	// Calculate overall risk level
	overallRiskLevel := p.calculateOverallRiskLevel(technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks)

	// Generate mitigation plan
	mitigationPlan := p.generateMitigationPlan(technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks)

	// Generate contingency planning
	contingencyPlanning := p.generateContingencyPlanning(inquiry, overallRiskLevel)

	// Generate risk monitoring indicators
	riskMonitoring := p.generateRiskMonitoring(technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks)

	riskAssessment := &interfaces.ProjectRiskAssessment{
		OverallRiskLevel:    overallRiskLevel,
		TechnicalRisks:      technicalRisks,
		BusinessRisks:       businessRisks,
		ResourceRisks:       resourceRisks,
		TimelineRisks:       timelineRisks,
		BudgetRisks:         budgetRisks,
		MitigationPlan:      mitigationPlan,
		ContingencyPlanning: contingencyPlanning,
		RiskMonitoring:      riskMonitoring,
	}

	return riskAssessment, nil
}

// GeneratePricingRecommendation generates pricing recommendations based on project complexity and market rates
func (p *proposalGenerator) GeneratePricingRecommendation(ctx context.Context, inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) (*interfaces.PricingRecommendation, error) {
	// Calculate base price from resource estimates
	basePrice := p.calculateBasePrice(inquiry, projectScope)

	// Generate price breakdown
	breakdown := p.generatePriceBreakdown(inquiry, projectScope, basePrice)

	// Calculate total price
	totalPrice := p.calculateTotalPrice(breakdown)

	// Generate market rate analysis
	marketRateAnalysis := p.generateMarketRateAnalysis(inquiry)

	// Generate competitive analysis
	competitiveAnalysis := p.generateCompetitivePricing(totalPrice, marketRateAnalysis)

	// Generate ROI projection for client
	roiProjection := p.generateROIProjection(totalPrice, inquiry)

	// Generate discounts if applicable
	discounts := p.generateDiscounts(totalPrice, inquiry)

	// Apply discounts to total price
	finalPrice := p.applyDiscounts(totalPrice, discounts)

	// Generate value proposition
	valueProposition := p.generateValueProposition(inquiry, finalPrice)

	pricingRecommendation := &interfaces.PricingRecommendation{
		TotalPrice:          finalPrice,
		Currency:            "USD",
		PricingModel:        p.determinePricingModel(inquiry),
		Breakdown:           breakdown,
		Discounts:           discounts,
		PaymentTerms:        p.generatePaymentTerms(finalPrice),
		ValidityPeriod:      "30 days",
		CompetitiveAnalysis: competitiveAnalysis,
		ValueProposition:    valueProposition,
		ROIProjection:       roiProjection,
		MarketRateAnalysis:  marketRateAnalysis,
	}

	return pricingRecommendation, nil
}

// GetSimilarProjects retrieves similar historical projects for comparison
func (p *proposalGenerator) GetSimilarProjects(ctx context.Context, inquiry *domain.Inquiry) ([]*interfaces.HistoricalProject, error) {
	// This would typically query a database of historical projects
	// For now, we'll return mock data based on the inquiry
	similarProjects := p.generateMockSimilarProjects(inquiry)

	// Calculate similarity scores
	for _, project := range similarProjects {
		project.SimilarityScore = p.calculateSimilarityScore(inquiry, project)
	}

	// Sort by similarity score (highest first)
	p.sortProjectsBySimilarity(similarProjects)

	// Return top 5 most similar projects
	if len(similarProjects) > 5 {
		similarProjects = similarProjects[:5]
	}

	return similarProjects, nil
}

// ValidateProposal validates a proposal for completeness and accuracy
func (p *proposalGenerator) ValidateProposal(proposal *interfaces.Proposal) error {
	if proposal == nil {
		return fmt.Errorf("proposal cannot be nil")
	}

	if proposal.ID == "" {
		return fmt.Errorf("proposal ID cannot be empty")
	}

	if proposal.InquiryID == "" {
		return fmt.Errorf("inquiry ID cannot be empty")
	}

	if proposal.Title == "" {
		return fmt.Errorf("proposal title cannot be empty")
	}

	if proposal.ExecutiveSummary == "" {
		return fmt.Errorf("executive summary cannot be empty")
	}

	if proposal.ProblemStatement == "" {
		return fmt.Errorf("problem statement cannot be empty")
	}

	if proposal.ProposedSolution == nil {
		return fmt.Errorf("proposed solution cannot be nil")
	}

	if proposal.ProjectScope == nil {
		return fmt.Errorf("project scope cannot be nil")
	}

	if len(proposal.NextSteps) == 0 {
		return fmt.Errorf("next steps cannot be empty")
	}

	return nil
}

// Helper methods for generating proposal components

// generateProposalTitle creates a title for the proposal
func (p *proposalGenerator) generateProposalTitle(inquiry *domain.Inquiry) string {
	company := inquiry.Company
	if company == "" {
		company = "Client Organization"
	}

	services := strings.Join(inquiry.Services, " & ")
	return fmt.Sprintf("%s - %s Proposal", company, strings.Title(services))
}

// generateExecutiveSummary generates the executive summary using AI
func (p *proposalGenerator) generateExecutiveSummary(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.ProposalOptions) (string, error) {
	prompt := p.buildExecutiveSummaryPrompt(inquiry, options)

	bedrockOptions := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   1000,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := p.bedrockService.GenerateText(ctx, prompt, bedrockOptions)
	if err != nil {
		return "", fmt.Errorf("failed to generate executive summary: %w", err)
	}

	return response.Content, nil
}

// generateProblemStatement generates the problem statement using AI
func (p *proposalGenerator) generateProblemStatement(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	prompt := p.buildProblemStatementPrompt(inquiry)

	bedrockOptions := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   800,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := p.bedrockService.GenerateText(ctx, prompt, bedrockOptions)
	if err != nil {
		return "", fmt.Errorf("failed to generate problem statement: %w", err)
	}

	return response.Content, nil
}

// generateProposedSolution generates the proposed solution using AI and multi-cloud analysis
func (p *proposalGenerator) generateProposedSolution(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.ProposalOptions) (*interfaces.ProposalSolution, error) {
	solution := &interfaces.ProposalSolution{
		ID:        uuid.New().String(),
		InquiryID: inquiry.ID,
	}

	// Generate solution overview using AI
	prompt := p.buildSolutionPrompt(inquiry, options)
	bedrockOptions := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   1500,
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := p.bedrockService.GenerateText(ctx, prompt, bedrockOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate solution overview: %w", err)
	}
	solution.SolutionOverview = response.Content

	// Determine cloud providers
	solution.CloudProviders = p.determineCloudProviders(inquiry, options)

	// Generate proposed services
	solution.Services = p.generateProposedServices(inquiry, solution.CloudProviders)

	// Generate architecture design
	solution.Architecture = p.generateArchitectureDesign(inquiry, solution.Services)

	// Generate technical approach
	solution.TechnicalApproach = p.generateTechnicalApproach(inquiry, solution.Services)

	// Generate security approach
	solution.SecurityApproach = p.generateSecurityApproach(inquiry, solution.Services)

	// Generate estimated cost and timeline
	solution.EstimatedCost = p.generateEstimatedCost(inquiry, solution.Services)
	solution.Timeline = p.generateEstimatedTimeline(inquiry, solution.Services)

	// Generate benefits
	solution.Benefits = p.generateSolutionBenefits(inquiry, solution.Services)

	return solution, nil
}

// generateProjectScope generates the project scope
func (p *proposalGenerator) generateProjectScope(ctx context.Context, inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) (*interfaces.ProjectScope, error) {
	scope := &interfaces.ProjectScope{
		InScope:      p.generateInScopeItems(inquiry, solution),
		OutOfScope:   p.generateOutOfScopeItems(inquiry, solution),
		Assumptions:  p.generateScopeAssumptions(inquiry, solution),
		Constraints:  p.generateScopeConstraints(inquiry, solution),
		Dependencies: p.generateScopeDependencies(inquiry, solution),
		Exclusions:   p.generateScopeExclusions(inquiry, solution),
	}

	return scope, nil
}

// Prompt building methods

// buildExecutiveSummaryPrompt builds a prompt for generating executive summary
func (p *proposalGenerator) buildExecutiveSummaryPrompt(inquiry *domain.Inquiry, options *interfaces.ProposalOptions) string {
	template := `Generate a compelling executive summary for a cloud consulting proposal.

CLIENT CONTEXT:
- Client: %s
- Company: %s
- Services Requested: %s
- Client Message: %s
- Target Audience: %s
- Proposal Type: %s

INSTRUCTIONS:
Create a concise executive summary (3-4 paragraphs) that:
1. Acknowledges the client's current challenges and objectives
2. Presents our understanding of their business needs
3. Highlights the key benefits of our proposed solution
4. Emphasizes the value proposition and expected outcomes
5. Maintains a professional, confident tone appropriate for %s audience

Focus on business value, competitive advantages, and measurable outcomes.
Avoid technical jargon unless the audience is technical.`

	return fmt.Sprintf(template,
		inquiry.Name,
		p.getCompanyOrDefault(inquiry.Company),
		strings.Join(inquiry.Services, ", "),
		inquiry.Message,
		options.TargetAudience,
		options.ProposalType,
		options.TargetAudience)
}

// buildProblemStatementPrompt builds a prompt for generating problem statement
func (p *proposalGenerator) buildProblemStatementPrompt(inquiry *domain.Inquiry) string {
	template := `Generate a clear problem statement for a cloud consulting proposal.

CLIENT CONTEXT:
- Client: %s
- Company: %s
- Services Requested: %s
- Client Message: %s

INSTRUCTIONS:
Create a problem statement (2-3 paragraphs) that:
1. Clearly articulates the client's current challenges
2. Identifies the root causes and pain points
3. Explains the business impact of not addressing these issues
4. Sets the context for why our solution is needed

Be specific and avoid generic statements. Focus on the client's actual situation.`

	return fmt.Sprintf(template,
		inquiry.Name,
		p.getCompanyOrDefault(inquiry.Company),
		strings.Join(inquiry.Services, ", "),
		inquiry.Message)
}

// buildSolutionPrompt builds a prompt for generating solution overview
func (p *proposalGenerator) buildSolutionPrompt(inquiry *domain.Inquiry, options *interfaces.ProposalOptions) string {
	template := `Generate a comprehensive solution overview for a cloud consulting proposal.

CLIENT CONTEXT:
- Client: %s
- Company: %s
- Services Requested: %s
- Client Message: %s
- Cloud Providers: %s
- Target Audience: %s

INSTRUCTIONS:
Create a solution overview (4-5 paragraphs) that:
1. Presents our recommended approach to address their challenges
2. Explains the high-level architecture and methodology
3. Highlights key technologies and cloud services to be used
4. Describes the implementation approach and phases
5. Emphasizes how this solution addresses their specific needs

Focus on the strategic approach rather than detailed technical specifications.
Tailor the technical depth to the %s audience.`

	cloudProviders := "AWS, Azure, GCP"
	if len(options.CloudProviders) > 0 {
		cloudProviders = strings.Join(options.CloudProviders, ", ")
	}

	return fmt.Sprintf(template,
		inquiry.Name,
		p.getCompanyOrDefault(inquiry.Company),
		strings.Join(inquiry.Services, ", "),
		inquiry.Message,
		cloudProviders,
		options.TargetAudience,
		options.TargetAudience)
}

// getCompanyOrDefault returns the company name or a default value
func (p *proposalGenerator) getCompanyOrDefault(company string) string {
	if company != "" {
		return company
	}
	return "Client Organization"
}

// Timeline and resource calculation methods

// calculateBaseDuration calculates base duration from similar projects
func (p *proposalGenerator) calculateBaseDuration(inquiry *domain.Inquiry, similarProjects []*interfaces.HistoricalProject) int {
	if len(similarProjects) == 0 {
		// Default durations based on service type
		return p.getDefaultDuration(inquiry.Services)
	}

	// Calculate weighted average duration from similar projects
	totalDuration := 0
	totalWeight := 0.0

	for _, project := range similarProjects {
		duration := p.parseDurationToDays(project.Duration)
		weight := project.SimilarityScore
		totalDuration += int(float64(duration) * weight)
		totalWeight += weight
	}

	if totalWeight > 0 {
		return int(float64(totalDuration) / totalWeight)
	}

	return p.getDefaultDuration(inquiry.Services)
}

// getDefaultDuration returns default duration based on services
func (p *proposalGenerator) getDefaultDuration(services []string) int {
	baseDuration := 30 // 30 days base

	for _, service := range services {
		switch strings.ToLower(service) {
		case "migration":
			baseDuration += 60 // Add 60 days for migration
		case "optimization":
			baseDuration += 30 // Add 30 days for optimization
		case "assessment":
			baseDuration += 15 // Add 15 days for assessment
		case "architecture":
			baseDuration += 45 // Add 45 days for architecture
		default:
			baseDuration += 20 // Add 20 days for other services
		}
	}

	return baseDuration
}

// parseDurationToDays parses duration string to days
func (p *proposalGenerator) parseDurationToDays(duration string) int {
	// Simple parsing - in real implementation, use proper duration parsing
	if strings.Contains(duration, "week") {
		// Extract number and multiply by 7
		return 84 // Default 12 weeks
	}
	if strings.Contains(duration, "month") {
		// Extract number and multiply by 30
		return 90 // Default 3 months
	}
	return 60 // Default 60 days
}

// generateTimelinePhases generates timeline phases based on services
func (p *proposalGenerator) generateTimelinePhases(ctx context.Context, inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope, totalDuration int) ([]interfaces.TimelinePhase, error) {
	var phases []interfaces.TimelinePhase

	// Discovery and Planning Phase (15% of total duration)
	discoveryDuration := int(float64(totalDuration) * 0.15)
	phases = append(phases, interfaces.TimelinePhase{
		ID:           uuid.New().String(),
		Name:         "Discovery and Planning",
		Description:  "Requirements gathering, current state assessment, and detailed planning",
		Duration:     fmt.Sprintf("%d days", discoveryDuration),
		Dependencies: []string{},
		Deliverables: []string{"Requirements Document", "Current State Assessment", "Project Plan"},
		Resources:    []string{"Solution Architect", "Business Analyst"},
		RiskLevel:    "Low",
	})

	// Implementation phases based on services
	remainingDuration := totalDuration - discoveryDuration
	implementationPhases := p.generateImplementationPhases(inquiry.Services, remainingDuration)
	phases = append(phases, implementationPhases...)

	// Testing and Validation Phase (10% of total duration)
	testingDuration := int(float64(totalDuration) * 0.10)
	phases = append(phases, interfaces.TimelinePhase{
		ID:           uuid.New().String(),
		Name:         "Testing and Validation",
		Description:  "System testing, user acceptance testing, and performance validation",
		Duration:     fmt.Sprintf("%d days", testingDuration),
		Dependencies: []string{phases[len(phases)-1].ID},
		Deliverables: []string{"Test Results", "Performance Report", "User Acceptance Sign-off"},
		Resources:    []string{"QA Engineer", "Performance Tester"},
		RiskLevel:    "Medium",
	})

	// Go-Live and Support Phase (5% of total duration)
	goLiveDuration := int(float64(totalDuration) * 0.05)
	phases = append(phases, interfaces.TimelinePhase{
		ID:           uuid.New().String(),
		Name:         "Go-Live and Support",
		Description:  "Production deployment, go-live support, and knowledge transfer",
		Duration:     fmt.Sprintf("%d days", goLiveDuration),
		Dependencies: []string{phases[len(phases)-1].ID},
		Deliverables: []string{"Production Deployment", "Support Documentation", "Knowledge Transfer"},
		Resources:    []string{"DevOps Engineer", "Support Specialist"},
		RiskLevel:    "High",
	})

	return phases, nil
}

// generateImplementationPhases generates implementation phases based on services
func (p *proposalGenerator) generateImplementationPhases(services []string, remainingDuration int) []interfaces.TimelinePhase {
	var phases []interfaces.TimelinePhase
	phaseDuration := remainingDuration / len(services)

	for i, service := range services {
		phase := interfaces.TimelinePhase{
			ID:          uuid.New().String(),
			Name:        fmt.Sprintf("%s Implementation", strings.Title(service)),
			Description: fmt.Sprintf("Implementation of %s services and components", service),
			Duration:    fmt.Sprintf("%d days", phaseDuration),
			RiskLevel:   p.getServiceRiskLevel(service),
		}

		// Set dependencies (each phase depends on the previous one)
		if i > 0 {
			phase.Dependencies = []string{phases[i-1].ID}
		}

		// Set deliverables based on service type
		phase.Deliverables = p.getServiceDeliverables(service)

		// Set resources based on service type
		phase.Resources = p.getServiceResources(service)

		phases = append(phases, phase)
	}

	return phases
}

// calculateTotalEffort calculates total effort in person-hours
func (p *proposalGenerator) calculateTotalEffort(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) float64 {
	baseEffort := 160.0 // 160 hours base (4 weeks * 40 hours)

	// Add effort based on services
	for _, service := range inquiry.Services {
		switch strings.ToLower(service) {
		case "migration":
			baseEffort += 480 // Add 480 hours (12 weeks)
		case "optimization":
			baseEffort += 240 // Add 240 hours (6 weeks)
		case "assessment":
			baseEffort += 120 // Add 120 hours (3 weeks)
		case "architecture":
			baseEffort += 320 // Add 320 hours (8 weeks)
		default:
			baseEffort += 160 // Add 160 hours (4 weeks)
		}
	}

	// Adjust based on project scope complexity
	complexityMultiplier := p.calculateComplexityMultiplier(projectScope)
	baseEffort *= complexityMultiplier

	return baseEffort
}

// calculateComplexityMultiplier calculates complexity multiplier based on project scope
func (p *proposalGenerator) calculateComplexityMultiplier(projectScope *interfaces.ProjectScope) float64 {
	multiplier := 1.0

	// Increase complexity based on scope size
	if len(projectScope.InScope) > 10 {
		multiplier += 0.3
	} else if len(projectScope.InScope) > 5 {
		multiplier += 0.2
	}

	// Increase complexity based on constraints
	if len(projectScope.Constraints) > 5 {
		multiplier += 0.2
	}

	// Increase complexity based on dependencies
	if len(projectScope.Dependencies) > 3 {
		multiplier += 0.3
	}

	return multiplier
}

// Additional helper methods

// generateMilestones generates milestones from phases
func (p *proposalGenerator) generateMilestones(phases []interfaces.TimelinePhase) []interfaces.ProposalMilestone {
	var milestones []interfaces.ProposalMilestone

	for _, phase := range phases {
		milestone := interfaces.ProposalMilestone{
			ID:          uuid.New().String(),
			Name:        fmt.Sprintf("%s Complete", phase.Name),
			Description: fmt.Sprintf("Completion of %s phase", phase.Name),
			Criteria:    phase.Deliverables,
			Critical:    phase.RiskLevel == "High",
		}
		milestones = append(milestones, milestone)
	}

	return milestones
}

// calculateTotalDuration calculates total duration from phases
func (p *proposalGenerator) calculateTotalDuration(phases []interfaces.TimelinePhase) int {
	totalDays := 0
	for _, phase := range phases {
		days := p.parseDurationToDays(phase.Duration)
		totalDays += days
	}
	return totalDays
}

// calculateTimelineConfidence calculates confidence based on similar projects
func (p *proposalGenerator) calculateTimelineConfidence(similarProjects []*interfaces.HistoricalProject, inquiry *domain.Inquiry) float64 {
	if len(similarProjects) == 0 {
		return 0.6 // Low confidence without historical data
	}

	// Higher confidence with more similar projects
	confidence := 0.5 + (float64(len(similarProjects)) * 0.1)
	if confidence > 0.9 {
		confidence = 0.9
	}

	return confidence
}

// identifyCriticalPath identifies the critical path through phases
func (p *proposalGenerator) identifyCriticalPath(phases []interfaces.TimelinePhase) []string {
	var criticalPath []string
	for _, phase := range phases {
		if phase.RiskLevel == "High" || phase.RiskLevel == "Medium" {
			criticalPath = append(criticalPath, phase.Name)
		}
	}
	return criticalPath
}

// convertToTimelineProjects converts historical projects for timeline
func (p *proposalGenerator) convertToTimelineProjects(projects []*interfaces.HistoricalProject) []interfaces.HistoricalProject {
	var timelineProjects []interfaces.HistoricalProject
	for _, project := range projects {
		timelineProjects = append(timelineProjects, *project)
	}
	return timelineProjects
}

// generateTeamComposition generates team composition based on services
func (p *proposalGenerator) generateTeamComposition(inquiry *domain.Inquiry, totalEffort float64) []interfaces.RoleRequirement {
	var roles []interfaces.RoleRequirement

	// Always include project manager
	roles = append(roles, interfaces.RoleRequirement{
		Role:        "Project Manager",
		Level:       "Senior",
		Allocation:  0.5, // 50% allocation
		Duration:    "Full project duration",
		Essential:   true,
		Description: "Overall project coordination and management",
		HourlyRate:  150.0,
	})

	// Add roles based on services
	for _, service := range inquiry.Services {
		serviceRoles := p.getServiceRoles(service, totalEffort)
		roles = append(roles, serviceRoles...)
	}

	return roles
}

// getServiceRoles returns roles required for a specific service
func (p *proposalGenerator) getServiceRoles(service string, totalEffort float64) []interfaces.RoleRequirement {
	var roles []interfaces.RoleRequirement

	switch strings.ToLower(service) {
	case "migration":
		roles = append(roles, interfaces.RoleRequirement{
			Role:        "Migration Specialist",
			Level:       "Senior",
			Allocation:  0.8,
			Duration:    "Migration phase",
			Essential:   true,
			Description: "Lead migration planning and execution",
			HourlyRate:  175.0,
		})
	case "architecture":
		roles = append(roles, interfaces.RoleRequirement{
			Role:        "Solution Architect",
			Level:       "Senior",
			Allocation:  0.7,
			Duration:    "Design and implementation phases",
			Essential:   true,
			Description: "Design and oversee architecture implementation",
			HourlyRate:  200.0,
		})
	case "optimization":
		roles = append(roles, interfaces.RoleRequirement{
			Role:        "Performance Engineer",
			Level:       "Mid",
			Allocation:  0.6,
			Duration:    "Optimization phase",
			Essential:   true,
			Description: "Analyze and optimize system performance",
			HourlyRate:  140.0,
		})
	}

	// Always add a cloud engineer
	roles = append(roles, interfaces.RoleRequirement{
		Role:        "Cloud Engineer",
		Level:       "Mid",
		Allocation:  0.8,
		Duration:    "Implementation phases",
		Essential:   true,
		Description: "Implement and configure cloud services",
		HourlyRate:  130.0,
	})

	return roles
}

// generateSkillRequirements generates skill requirements based on services
func (p *proposalGenerator) generateSkillRequirements(inquiry *domain.Inquiry) []interfaces.SkillRequirement {
	var skills []interfaces.SkillRequirement

	// Common skills
	skills = append(skills, interfaces.SkillRequirement{
		Skill:       "AWS",
		Level:       "Advanced",
		Essential:   true,
		Description: "Amazon Web Services platform expertise",
	})

	// Add service-specific skills
	for _, service := range inquiry.Services {
		serviceSkills := p.getServiceSkills(service)
		skills = append(skills, serviceSkills...)
	}

	return skills
}

// getServiceSkills returns skills required for a specific service
func (p *proposalGenerator) getServiceSkills(service string) []interfaces.SkillRequirement {
	var skills []interfaces.SkillRequirement

	switch strings.ToLower(service) {
	case "migration":
		skills = append(skills, interfaces.SkillRequirement{
			Skill:       "Database Migration",
			Level:       "Advanced",
			Essential:   true,
			Description: "Experience with database migration tools and strategies",
		})
	case "architecture":
		skills = append(skills, interfaces.SkillRequirement{
			Skill:       "System Design",
			Level:       "Expert",
			Essential:   true,
			Description: "Large-scale system architecture design",
		})
	}

	return skills
}

// generateMockSimilarProjects generates mock similar projects for demonstration
func (p *proposalGenerator) generateMockSimilarProjects(inquiry *domain.Inquiry) []*interfaces.HistoricalProject {
	var projects []*interfaces.HistoricalProject

	// Generate 3-5 mock projects based on services
	for i, service := range inquiry.Services {
		if i >= 5 { // Limit to 5 projects
			break
		}

		project := &interfaces.HistoricalProject{
			ID:         uuid.New().String(),
			Name:       fmt.Sprintf("%s Project %d", strings.Title(service), i+1),
			Industry:   p.inferIndustry(inquiry),
			Services:   []string{service},
			Duration:   p.getServiceDuration(service),
			TeamSize:   p.getServiceTeamSize(service),
			Budget:     p.getServiceBudget(service),
			Complexity: p.getServiceComplexity(service),
			SuccessMetrics: []string{
				"On-time delivery",
				"Budget adherence",
				"Client satisfaction > 4.5/5",
			},
			LessonsLearned: []string{
				"Early stakeholder engagement is critical",
				"Thorough testing prevents production issues",
			},
			CompletedAt: time.Now().AddDate(0, -6, 0), // 6 months ago
		}

		projects = append(projects, project)
	}

	return projects
}

// calculateSimilarityScore calculates similarity score between inquiry and project
func (p *proposalGenerator) calculateSimilarityScore(inquiry *domain.Inquiry, project *interfaces.HistoricalProject) float64 {
	score := 0.0

	// Service similarity (40% weight)
	serviceMatch := p.calculateServiceMatch(inquiry.Services, project.Services)
	score += serviceMatch * 0.4

	// Industry similarity (30% weight)
	industryMatch := p.calculateIndustryMatch(inquiry, project.Industry)
	score += industryMatch * 0.3

	// Complexity similarity (30% weight)
	complexityMatch := p.calculateComplexityMatch(inquiry, project.Complexity)
	score += complexityMatch * 0.3

	return score
}

// Helper methods for similarity calculation
func (p *proposalGenerator) calculateServiceMatch(inquiryServices, projectServices []string) float64 {
	matches := 0
	for _, iService := range inquiryServices {
		for _, pService := range projectServices {
			if strings.EqualFold(iService, pService) {
				matches++
				break
			}
		}
	}
	return float64(matches) / float64(len(inquiryServices))
}

func (p *proposalGenerator) calculateIndustryMatch(inquiry *domain.Inquiry, projectIndustry string) float64 {
	inferredIndustry := p.inferIndustry(inquiry)
	if strings.EqualFold(inferredIndustry, projectIndustry) {
		return 1.0
	}
	return 0.5 // Partial match for related industries
}

func (p *proposalGenerator) calculateComplexityMatch(inquiry *domain.Inquiry, projectComplexity string) float64 {
	inquiryComplexity := p.inferComplexity(inquiry)
	if strings.EqualFold(inquiryComplexity, projectComplexity) {
		return 1.0
	}
	return 0.7 // Partial match for similar complexity
}

// sortProjectsBySimilarity sorts projects by similarity score
func (p *proposalGenerator) sortProjectsBySimilarity(projects []*interfaces.HistoricalProject) {
	// Simple bubble sort by similarity score (descending)
	n := len(projects)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if projects[j].SimilarityScore < projects[j+1].SimilarityScore {
				projects[j], projects[j+1] = projects[j+1], projects[j]
			}
		}
	}
}

// Utility methods for mock data generation
func (p *proposalGenerator) inferIndustry(inquiry *domain.Inquiry) string {
	company := strings.ToLower(inquiry.Company)
	message := strings.ToLower(inquiry.Message)

	if strings.Contains(company, "health") || strings.Contains(message, "hipaa") {
		return "Healthcare"
	}
	if strings.Contains(company, "bank") || strings.Contains(message, "financial") {
		return "Financial Services"
	}
	if strings.Contains(company, "retail") || strings.Contains(message, "ecommerce") {
		return "Retail"
	}
	return "Technology"
}

func (p *proposalGenerator) inferComplexity(inquiry *domain.Inquiry) string {
	if len(inquiry.Services) > 3 {
		return "High"
	}
	if len(inquiry.Services) > 1 {
		return "Medium"
	}
	return "Low"
}

func (p *proposalGenerator) getServiceDuration(service string) string {
	switch strings.ToLower(service) {
	case "migration":
		return "12 weeks"
	case "architecture":
		return "8 weeks"
	case "optimization":
		return "6 weeks"
	default:
		return "4 weeks"
	}
}

func (p *proposalGenerator) getServiceTeamSize(service string) int {
	switch strings.ToLower(service) {
	case "migration":
		return 6
	case "architecture":
		return 4
	case "optimization":
		return 3
	default:
		return 3
	}
}

func (p *proposalGenerator) getServiceBudget(service string) float64 {
	switch strings.ToLower(service) {
	case "migration":
		return 250000.0
	case "architecture":
		return 150000.0
	case "optimization":
		return 100000.0
	default:
		return 75000.0
	}
}

func (p *proposalGenerator) getServiceComplexity(service string) string {
	switch strings.ToLower(service) {
	case "migration":
		return "High"
	case "architecture":
		return "Medium"
	default:
		return "Low"
	}
}

// Pricing and risk assessment methods

// calculateBasePrice calculates base price from resource estimates
func (p *proposalGenerator) calculateBasePrice(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) float64 {
	basePrice := 50000.0 // Base price

	// Add price based on services
	for _, service := range inquiry.Services {
		basePrice += p.getServicePrice(service)
	}

	// Adjust based on complexity
	complexityMultiplier := p.calculateComplexityMultiplier(projectScope)
	basePrice *= complexityMultiplier

	return basePrice
}

// getServicePrice returns base price for a service
func (p *proposalGenerator) getServicePrice(service string) float64 {
	switch strings.ToLower(service) {
	case "migration":
		return 100000.0
	case "architecture":
		return 75000.0
	case "optimization":
		return 50000.0
	case "assessment":
		return 25000.0
	default:
		return 40000.0
	}
}

// generatePriceBreakdown generates detailed price breakdown
func (p *proposalGenerator) generatePriceBreakdown(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope, basePrice float64) []interfaces.PriceComponent {
	var breakdown []interfaces.PriceComponent

	// Professional services (70% of total)
	professionalServices := basePrice * 0.7
	breakdown = append(breakdown, interfaces.PriceComponent{
		Category:    "Professional Services",
		Description: "Consulting, implementation, and project management services",
		Quantity:    1,
		UnitPrice:   professionalServices,
		TotalPrice:  professionalServices,
	})

	// Tools and licenses (20% of total)
	toolsLicenses := basePrice * 0.2
	breakdown = append(breakdown, interfaces.PriceComponent{
		Category:    "Tools and Licenses",
		Description: "Required software tools and cloud service licenses",
		Quantity:    1,
		UnitPrice:   toolsLicenses,
		TotalPrice:  toolsLicenses,
	})

	// Training and knowledge transfer (10% of total)
	training := basePrice * 0.1
	breakdown = append(breakdown, interfaces.PriceComponent{
		Category:    "Training and Knowledge Transfer",
		Description: "Team training and documentation",
		Quantity:    1,
		UnitPrice:   training,
		TotalPrice:  training,
	})

	return breakdown
}

// calculateTotalPrice calculates total price from breakdown
func (p *proposalGenerator) calculateTotalPrice(breakdown []interfaces.PriceComponent) float64 {
	total := 0.0
	for _, component := range breakdown {
		total += component.TotalPrice
	}
	return total
}

// generateTechnicalRisks generates technical risks for the project
func (p *proposalGenerator) generateTechnicalRisks(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) []interfaces.ProjectRisk {
	var risks []interfaces.ProjectRisk

	// Common technical risks
	risks = append(risks, interfaces.ProjectRisk{
		ID:          uuid.New().String(),
		Category:    "Technical",
		Title:       "Integration Complexity",
		Description: "Challenges integrating with existing systems",
		Impact:      "Medium",
		Probability: "Medium",
		RiskScore:   6,
		Mitigation:  "Thorough integration testing and phased rollout",
		Owner:       "Technical Lead",
		Status:      "Identified",
	})

	// Service-specific risks
	for _, service := range inquiry.Services {
		serviceRisks := p.getServiceRisks(service)
		risks = append(risks, serviceRisks...)
	}

	return risks
}

// getServiceRisks returns risks specific to a service
func (p *proposalGenerator) getServiceRisks(service string) []interfaces.ProjectRisk {
	var risks []interfaces.ProjectRisk

	switch strings.ToLower(service) {
	case "migration":
		risks = append(risks, interfaces.ProjectRisk{
			ID:          uuid.New().String(),
			Category:    "Technical",
			Title:       "Data Migration Issues",
			Description: "Potential data loss or corruption during migration",
			Impact:      "High",
			Probability: "Low",
			RiskScore:   6,
			Mitigation:  "Comprehensive backup and validation procedures",
			Owner:       "Migration Specialist",
			Status:      "Identified",
		})
	case "architecture":
		risks = append(risks, interfaces.ProjectRisk{
			ID:          uuid.New().String(),
			Category:    "Technical",
			Title:       "Scalability Concerns",
			Description: "Architecture may not scale as expected",
			Impact:      "Medium",
			Probability: "Medium",
			RiskScore:   6,
			Mitigation:  "Load testing and performance monitoring",
			Owner:       "Solution Architect",
			Status:      "Identified",
		})
	}

	return risks
}

// generateBusinessRisks generates business risks
func (p *proposalGenerator) generateBusinessRisks(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) []interfaces.ProjectRisk {
	var risks []interfaces.ProjectRisk

	risks = append(risks, interfaces.ProjectRisk{
		ID:          uuid.New().String(),
		Category:    "Business",
		Title:       "Stakeholder Alignment",
		Description: "Risk of misaligned expectations between stakeholders",
		Impact:      "Medium",
		Probability: "Medium",
		RiskScore:   6,
		Mitigation:  "Regular stakeholder meetings and clear communication",
		Owner:       "Project Manager",
		Status:      "Identified",
	})

	return risks
}

// generateResourceRisks generates resource-related risks
func (p *proposalGenerator) generateResourceRisks(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) []interfaces.ProjectRisk {
	var risks []interfaces.ProjectRisk

	risks = append(risks, interfaces.ProjectRisk{
		ID:          uuid.New().String(),
		Category:    "Resource",
		Title:       "Key Personnel Availability",
		Description: "Risk of key team members being unavailable",
		Impact:      "High",
		Probability: "Low",
		RiskScore:   6,
		Mitigation:  "Cross-training and backup resource identification",
		Owner:       "Project Manager",
		Status:      "Identified",
	})

	return risks
}

// generateTimelineRisks generates timeline-related risks
func (p *proposalGenerator) generateTimelineRisks(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) []interfaces.ProjectRisk {
	var risks []interfaces.ProjectRisk

	risks = append(risks, interfaces.ProjectRisk{
		ID:          uuid.New().String(),
		Category:    "Timeline",
		Title:       "Scope Creep",
		Description: "Additional requirements extending project timeline",
		Impact:      "Medium",
		Probability: "High",
		RiskScore:   8,
		Mitigation:  "Clear scope definition and change control process",
		Owner:       "Project Manager",
		Status:      "Identified",
	})

	return risks
}

// generateBudgetRisks generates budget-related risks
func (p *proposalGenerator) generateBudgetRisks(inquiry *domain.Inquiry, projectScope *interfaces.ProjectScope) []interfaces.ProjectRisk {
	var risks []interfaces.ProjectRisk

	risks = append(risks, interfaces.ProjectRisk{
		ID:          uuid.New().String(),
		Category:    "Budget",
		Title:       "Cost Overrun",
		Description: "Project costs exceeding approved budget",
		Impact:      "High",
		Probability: "Medium",
		RiskScore:   8,
		Mitigation:  "Regular budget monitoring and contingency planning",
		Owner:       "Project Manager",
		Status:      "Identified",
	})

	return risks
}

// calculateOverallRiskLevel calculates overall risk level
func (p *proposalGenerator) calculateOverallRiskLevel(technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks []interfaces.ProjectRisk) string {
	totalScore := 0
	riskCount := 0

	allRisks := [][]interfaces.ProjectRisk{technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks}
	for _, riskCategory := range allRisks {
		for _, risk := range riskCategory {
			totalScore += risk.RiskScore
			riskCount++
		}
	}

	if riskCount == 0 {
		return "Low"
	}

	averageScore := float64(totalScore) / float64(riskCount)
	if averageScore >= 8 {
		return "High"
	} else if averageScore >= 6 {
		return "Medium"
	}
	return "Low"
}

// generateMitigationPlan generates a mitigation plan for risks
func (p *proposalGenerator) generateMitigationPlan(technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks []interfaces.ProjectRisk) *interfaces.MitigationPlan {
	var strategies []interfaces.ProposalMitigationStrategy

	// Generate strategies for high-impact risks
	allRisks := [][]interfaces.ProjectRisk{technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks}
	for _, riskCategory := range allRisks {
		for _, risk := range riskCategory {
			if risk.Impact == "High" || risk.RiskScore >= 8 {
				strategy := interfaces.ProposalMitigationStrategy{
					RiskID:        risk.ID,
					Strategy:      risk.Mitigation,
					Actions:       []string{"Monitor regularly", "Implement preventive measures"},
					Timeline:      "Throughout project",
					Cost:          5000.0, // Estimated mitigation cost
					Owner:         risk.Owner,
					Effectiveness: "High",
				}
				strategies = append(strategies, strategy)
			}
		}
	}

	return &interfaces.MitigationPlan{
		Strategies:      strategies,
		ContingencyFund: 25000.0, // 10% of typical project budget
		EscalationPlan:  "Escalate high-risk issues to project sponsor within 24 hours",
		ReviewSchedule:  "Weekly risk review meetings",
	}
}

// Additional utility methods for completing the implementation

// generateNextSteps generates next steps for the proposal
func (p *proposalGenerator) generateNextSteps(inquiry *domain.Inquiry, options *interfaces.ProposalOptions) []string {
	return []string{
		"Review and approve this proposal",
		"Schedule kick-off meeting with key stakeholders",
		"Finalize project timeline and resource allocation",
		"Execute statement of work and contracts",
		"Begin discovery and planning phase",
	}
}

// generateAssumptions generates project assumptions
func (p *proposalGenerator) generateAssumptions(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Client will provide timely access to required systems and data",
		"Key stakeholders will be available for regular meetings",
		"Existing infrastructure meets minimum requirements",
		"No major organizational changes during project execution",
		"Required approvals and sign-offs will be obtained promptly",
	}
}

// generateSuccessMetrics generates success metrics
func (p *proposalGenerator) generateSuccessMetrics(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Project delivered on time and within budget",
		"All functional requirements met and tested",
		"Client satisfaction score of 4.5/5 or higher",
		"Zero critical defects in production",
		"Successful knowledge transfer to client team",
	}
}

// generateDeliverables generates project deliverables
func (p *proposalGenerator) generateDeliverables(ctx context.Context, inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) ([]interfaces.ProposalDeliverable, error) {
	var deliverables []interfaces.ProposalDeliverable

	// Common deliverables
	deliverables = append(deliverables, interfaces.ProposalDeliverable{
		ID:          uuid.New().String(),
		Name:        "Project Plan",
		Description: "Detailed project plan with timeline and milestones",
		Type:        "Document",
		DueDate:     time.Now().AddDate(0, 0, 7), // 1 week
		Owner:       "Project Manager",
		Status:      "Planned",
	})

	// Service-specific deliverables
	for _, service := range inquiry.Services {
		serviceDeliverables := p.getServiceDeliverables(service)
		for _, deliverableName := range serviceDeliverables {
			deliverable := interfaces.ProposalDeliverable{
				ID:          uuid.New().String(),
				Name:        deliverableName,
				Description: fmt.Sprintf("%s deliverable for %s service", deliverableName, service),
				Type:        "Technical",
				DueDate:     time.Now().AddDate(0, 1, 0), // 1 month
				Owner:       "Technical Lead",
				Status:      "Planned",
			}
			deliverables = append(deliverables, deliverable)
		}
	}

	return deliverables, nil
}

// getServiceDeliverables returns deliverables for a service
func (p *proposalGenerator) getServiceDeliverables(service string) []string {
	switch strings.ToLower(service) {
	case "migration":
		return []string{"Migration Plan", "Data Migration Scripts", "Cutover Procedures"}
	case "architecture":
		return []string{"Architecture Design", "Technical Specifications", "Implementation Guide"}
	case "optimization":
		return []string{"Performance Analysis", "Optimization Recommendations", "Monitoring Setup"}
	case "assessment":
		return []string{"Current State Assessment", "Gap Analysis", "Recommendations Report"}
	default:
		return []string{"Technical Documentation", "Implementation Guide"}
	}
}

// getServiceResources returns resources needed for a service
func (p *proposalGenerator) getServiceResources(service string) []string {
	switch strings.ToLower(service) {
	case "migration":
		return []string{"Migration Specialist", "Database Administrator", "Cloud Engineer"}
	case "architecture":
		return []string{"Solution Architect", "Cloud Engineer", "Security Specialist"}
	case "optimization":
		return []string{"Performance Engineer", "Cloud Engineer", "Monitoring Specialist"}
	default:
		return []string{"Cloud Engineer", "Technical Consultant"}
	}
}

// getServiceRiskLevel returns risk level for a service
func (p *proposalGenerator) getServiceRiskLevel(service string) string {
	switch strings.ToLower(service) {
	case "migration":
		return "High"
	case "architecture":
		return "Medium"
	default:
		return "Low"
	}
}

// Placeholder methods for detailed SOW generation (to be implemented)
func (p *proposalGenerator) generateProjectOverview(ctx context.Context, inquiry *domain.Inquiry, proposal *interfaces.Proposal) (string, error) {
	return fmt.Sprintf("This project will %s for %s, delivering comprehensive cloud solutions.",
		strings.Join(inquiry.Services, " and "),
		p.getCompanyOrDefault(inquiry.Company)), nil
}

func (p *proposalGenerator) generateDetailedScope(ctx context.Context, inquiry *domain.Inquiry, proposal *interfaces.Proposal) (*interfaces.DetailedScope, error) {
	return &interfaces.DetailedScope{
		ProjectScope:           *proposal.ProjectScope,
		WorkBreakdownStructure: []interfaces.WorkPackage{},
		TechnicalRequirements:  []interfaces.TechnicalRequirement{},
		FunctionalRequirements: []interfaces.FunctionalRequirement{},
	}, nil
}

func (p *proposalGenerator) generateDetailedDeliverables(ctx context.Context, proposal *interfaces.Proposal) ([]interfaces.DetailedDeliverable, error) {
	var detailed []interfaces.DetailedDeliverable
	for _, deliverable := range proposal.Deliverables {
		detailed = append(detailed, interfaces.DetailedDeliverable{
			ProposalDeliverable: deliverable,
			AcceptanceCriteria:  []interfaces.AcceptanceCriterion{},
			QualityStandards:    []interfaces.QualityStandard{},
		})
	}
	return detailed, nil
}

func (p *proposalGenerator) generateDetailedTimeline(ctx context.Context, timeline *interfaces.TimelineEstimate) (*interfaces.DetailedTimeline, error) {
	return &interfaces.DetailedTimeline{
		TimelineEstimate: *timeline,
		WorkSchedule:     []interfaces.WorkScheduleItem{},
		ResourceCalendar: []interfaces.ResourceCalendarItem{},
		BufferAllocation: []interfaces.BufferAllocation{},
	}, nil
}

func (p *proposalGenerator) generateDetailedResources(ctx context.Context, resources *interfaces.ProposalResourceEstimate) (*interfaces.DetailedResources, error) {
	return &interfaces.DetailedResources{
		ResourcePlan:    []interfaces.ProposalResourcePlanItem{},
		SkillMatrix:     []interfaces.SkillMatrixItem{},
		OnboardingPlan:  "Standard onboarding process for new team members",
		OffboardingPlan: "Knowledge transfer and documentation handover",
	}, nil
}

func (p *proposalGenerator) generatePaymentSchedule(ctx context.Context, pricing *interfaces.PricingRecommendation, timeline *interfaces.TimelineEstimate) (*interfaces.PaymentSchedule, error) {
	milestones := []interfaces.PaymentMilestone{
		{
			ID:         uuid.New().String(),
			Name:       "Project Initiation",
			Amount:     pricing.TotalPrice * 0.3,
			Percentage: 30.0,
			DueDate:    time.Now().AddDate(0, 0, 7),
			Criteria:   []string{"Signed contract", "Project kickoff completed"},
			Status:     "Pending",
		},
		{
			ID:         uuid.New().String(),
			Name:       "Mid-project Milestone",
			Amount:     pricing.TotalPrice * 0.4,
			Percentage: 40.0,
			DueDate:    time.Now().AddDate(0, 1, 0),
			Criteria:   []string{"50% of deliverables completed", "Client approval received"},
			Status:     "Pending",
		},
		{
			ID:         uuid.New().String(),
			Name:       "Project Completion",
			Amount:     pricing.TotalPrice * 0.3,
			Percentage: 30.0,
			DueDate:    time.Now().AddDate(0, 2, 0),
			Criteria:   []string{"All deliverables completed", "Final acceptance received"},
			Status:     "Pending",
		},
	}

	return &interfaces.PaymentSchedule{
		TotalAmount: pricing.TotalPrice,
		Currency:    pricing.Currency,
		Milestones:  milestones,
		Terms:       "Net 30 days",
		LateFees:    "1.5% per month on overdue amounts",
	}, nil
}

func (p *proposalGenerator) generateAcceptanceCriteria(ctx context.Context, proposal *interfaces.Proposal) ([]interfaces.AcceptanceCriterion, error) {
	return []interfaces.AcceptanceCriterion{
		{
			ID:           uuid.New().String(),
			Description:  "All functional requirements implemented and tested",
			TestMethod:   "User acceptance testing",
			PassCriteria: "100% of test cases pass",
		},
		{
			ID:           uuid.New().String(),
			Description:  "Performance meets specified requirements",
			TestMethod:   "Load testing",
			PassCriteria: "Response time < 2 seconds under normal load",
		},
	}, nil
}

// Additional utility methods for pricing
func (p *proposalGenerator) generateMarketRateAnalysis(inquiry *domain.Inquiry) *interfaces.MarketRateAnalysis {
	return &interfaces.MarketRateAnalysis{
		Region:          "North America",
		ServiceCategory: strings.Join(inquiry.Services, ", "),
		RateRanges: map[string]interfaces.PriceRange{
			"Senior Consultant": {Low: 150, High: 250, Average: 200},
			"Cloud Engineer":    {Low: 100, High: 180, Average: 140},
			"Project Manager":   {Low: 120, High: 200, Average: 160},
		},
		DataSource:  "Industry salary surveys and market research",
		LastUpdated: time.Now(),
	}
}

func (p *proposalGenerator) generateCompetitivePricing(totalPrice float64, marketAnalysis *interfaces.MarketRateAnalysis) *interfaces.CompetitivePricing {
	marketAverage := totalPrice * 1.1 // Assume market is 10% higher
	return &interfaces.CompetitivePricing{
		MarketRange: interfaces.PriceRange{
			Low:     totalPrice * 0.8,
			High:    totalPrice * 1.3,
			Average: marketAverage,
		},
		OurPosition:   "below",
		Justification: "Competitive pricing while maintaining high quality standards",
		Differentiators: []string{
			"Proven track record",
			"Dedicated project management",
			"Post-implementation support",
		},
	}
}

func (p *proposalGenerator) generateROIProjection(totalPrice float64, inquiry *domain.Inquiry) *interfaces.ProposalROIProjection {
	annualSavings := totalPrice * 0.4 // Assume 40% annual savings
	return &interfaces.ProposalROIProjection{
		InitialInvestment: totalPrice,
		AnnualSavings:     annualSavings,
		PaybackPeriod:     "2.5 years",
		ThreeYearROI:      (annualSavings*3 - totalPrice) / totalPrice * 100,
		FiveYearROI:       (annualSavings*5 - totalPrice) / totalPrice * 100,
		Assumptions: []string{
			"Current operational costs remain stable",
			"Projected efficiency gains are realized",
			"No major technology changes required",
		},
	}
}

func (p *proposalGenerator) generateDiscounts(totalPrice float64, inquiry *domain.Inquiry) []interfaces.Discount {
	var discounts []interfaces.Discount

	// Volume discount for large projects
	if totalPrice > 200000 {
		discounts = append(discounts, interfaces.Discount{
			Type:        "percentage",
			Description: "Large project volume discount",
			Amount:      5.0,
			Conditions:  []string{"Project value > $200,000"},
		})
	}

	return discounts
}

func (p *proposalGenerator) applyDiscounts(totalPrice float64, discounts []interfaces.Discount) float64 {
	finalPrice := totalPrice
	for _, discount := range discounts {
		if discount.Type == "percentage" {
			finalPrice *= (1.0 - discount.Amount/100.0)
		} else if discount.Type == "fixed" {
			finalPrice -= discount.Amount
		}
	}
	return finalPrice
}

func (p *proposalGenerator) generateValueProposition(inquiry *domain.Inquiry, finalPrice float64) string {
	return fmt.Sprintf("Our comprehensive %s solution delivers measurable business value through improved efficiency, reduced operational costs, and enhanced scalability, with an estimated ROI of 120%% over 3 years.",
		strings.Join(inquiry.Services, " and "))
}

func (p *proposalGenerator) determinePricingModel(inquiry *domain.Inquiry) string {
	if len(inquiry.Services) > 2 {
		return "milestone" // Complex projects use milestone-based pricing
	}
	return "fixed" // Simple projects use fixed pricing
}

func (p *proposalGenerator) generatePaymentTerms(finalPrice float64) string {
	if finalPrice > 100000 {
		return "30% upfront, 40% at mid-project milestone, 30% upon completion"
	}
	return "50% upfront, 50% upon completion"
}

func (p *proposalGenerator) generateExternalResources(inquiry *domain.Inquiry) []interfaces.ExternalResource {
	return []interfaces.ExternalResource{
		{
			Type:        "service",
			Description: "Third-party security audit",
			Duration:    "1 week",
			Cost:        15000.0,
			Essential:   false,
		},
	}
}

func (p *proposalGenerator) generateToolsAndLicenses(inquiry *domain.Inquiry) []interfaces.ToolRequirement {
	return []interfaces.ToolRequirement{
		{
			Name:        "Project Management Software",
			Type:        "software",
			Description: "Collaborative project management platform",
			Cost:        2000.0,
			Duration:    "Project duration",
			Essential:   true,
		},
	}
}

func (p *proposalGenerator) generateTrainingNeeds(inquiry *domain.Inquiry, skills []interfaces.SkillRequirement) []interfaces.TrainingNeed {
	return []interfaces.TrainingNeed{
		{
			Topic:       "Cloud Best Practices",
			Audience:    []string{"Client technical team"},
			Duration:    "2 days",
			Cost:        5000.0,
			Essential:   true,
			Description: "Training on cloud architecture and operational best practices",
		},
	}
}

func (p *proposalGenerator) calculateResourceConfidence(inquiry *domain.Inquiry) float64 {
	// Higher confidence for simpler projects
	if len(inquiry.Services) == 1 {
		return 0.9
	} else if len(inquiry.Services) <= 3 {
		return 0.8
	}
	return 0.7
}

func (p *proposalGenerator) generateContingencyPlanning(inquiry *domain.Inquiry, overallRiskLevel string) string {
	switch overallRiskLevel {
	case "High":
		return "Comprehensive contingency plan with 20% budget buffer and alternative implementation approaches"
	case "Medium":
		return "Standard contingency plan with 15% budget buffer and backup resource allocation"
	default:
		return "Basic contingency plan with 10% budget buffer and risk monitoring procedures"
	}
}

func (p *proposalGenerator) generateRiskMonitoring(technicalRisks, businessRisks, resourceRisks, timelineRisks, budgetRisks []interfaces.ProjectRisk) []interfaces.RiskIndicator {
	return []interfaces.RiskIndicator{
		{
			RiskID:    "general",
			Indicator: "Project milestone delays",
			Threshold: "> 1 week behind schedule",
			Frequency: "Weekly",
			Owner:     "Project Manager",
		},
		{
			RiskID:    "general",
			Indicator: "Budget variance",
			Threshold: "> 10% over budget",
			Frequency: "Monthly",
			Owner:     "Project Manager",
		},
	}
}

// Utility methods for scope generation
func (p *proposalGenerator) generateInScopeItems(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	items := []string{
		"Requirements analysis and documentation",
		"Solution design and architecture",
		"Implementation and configuration",
		"Testing and validation",
		"Documentation and knowledge transfer",
	}

	// Add service-specific items
	for _, service := range inquiry.Services {
		switch strings.ToLower(service) {
		case "migration":
			items = append(items, "Data migration and validation", "Application migration")
		case "optimization":
			items = append(items, "Performance analysis", "Cost optimization recommendations")
		case "architecture":
			items = append(items, "Architecture review", "Design documentation")
		}
	}

	return items
}

func (p *proposalGenerator) generateOutOfScopeItems(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Hardware procurement and installation",
		"Third-party software licensing costs",
		"Ongoing operational support beyond go-live",
		"Training for end users (beyond technical team)",
		"Data center or network infrastructure changes",
	}
}

func (p *proposalGenerator) generateScopeAssumptions(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Client will provide necessary access to systems and data",
		"Key stakeholders will be available for regular meetings",
		"Existing infrastructure meets minimum requirements",
		"No major organizational changes during project",
		"Required approvals will be obtained in timely manner",
	}
}

func (p *proposalGenerator) generateScopeConstraints(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Project must be completed within specified timeline",
		"Budget constraints as outlined in the proposal",
		"Compliance with existing security policies",
		"Minimal disruption to ongoing operations",
		"Use of approved vendors and technologies only",
	}
}

func (p *proposalGenerator) generateScopeDependencies(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Client approval of project plan and timeline",
		"Access to required systems and environments",
		"Availability of key client personnel",
		"Completion of any prerequisite projects",
		"Procurement of necessary licenses and tools",
	}
}

func (p *proposalGenerator) generateScopeExclusions(inquiry *domain.Inquiry, solution *interfaces.ProposalSolution) []string {
	return []string{
		"Legacy system decommissioning",
		"End-user training programs",
		"Ongoing maintenance and support",
		"Hardware and infrastructure costs",
		"Third-party integration beyond specified scope",
	}
}

// Additional utility methods for architecture and services
func (p *proposalGenerator) determineCloudProviders(inquiry *domain.Inquiry, options *interfaces.ProposalOptions) []string {
	if len(options.CloudProviders) > 0 {
		return options.CloudProviders
	}

	// Default to AWS for most services
	providers := []string{"AWS"}

	// Add additional providers based on services
	for _, service := range inquiry.Services {
		if strings.ToLower(service) == "migration" {
			providers = append(providers, "Azure") // Multi-cloud migration
			break
		}
	}

	return providers
}

func (p *proposalGenerator) generateProposedServices(inquiry *domain.Inquiry, cloudProviders []string) []interfaces.ProposedService {
	var services []interfaces.ProposedService

	for _, service := range inquiry.Services {
		proposedService := interfaces.ProposedService{
			Name:        strings.Title(service),
			Provider:    cloudProviders[0], // Use primary provider
			Description: fmt.Sprintf("Comprehensive %s services", service),
			Purpose:     fmt.Sprintf("To address %s requirements", service),
			Benefits:    []string{"Improved efficiency", "Cost optimization", "Enhanced scalability"},
			Cost:        "Included in total project cost",
		}
		services = append(services, proposedService)
	}

	return services
}

func (p *proposalGenerator) generateArchitectureDesign(inquiry *domain.Inquiry, services []interfaces.ProposedService) *interfaces.ArchitectureDesign {
	return &interfaces.ArchitectureDesign{
		Overview:    "Modern cloud-native architecture designed for scalability and reliability",
		Components:  []interfaces.ArchComponent{},
		Connections: []interfaces.ArchConnection{},
		Patterns:    []string{"Microservices", "Event-driven", "Serverless"},
		Principles:  []string{"Security by design", "High availability", "Cost optimization"},
	}
}

func (p *proposalGenerator) generateTechnicalApproach(inquiry *domain.Inquiry, services []interfaces.ProposedService) string {
	return "Our technical approach follows industry best practices with a focus on security, scalability, and maintainability. We will implement a phased approach to minimize risk and ensure smooth delivery."
}

func (p *proposalGenerator) generateSecurityApproach(inquiry *domain.Inquiry, services []interfaces.ProposedService) string {
	return "Security is integrated throughout the solution with defense-in-depth principles, encryption at rest and in transit, identity and access management, and continuous monitoring."
}

func (p *proposalGenerator) generateEstimatedCost(inquiry *domain.Inquiry, services []interfaces.ProposedService) string {
	basePrice := p.calculateBasePrice(inquiry, &interfaces.ProjectScope{})
	return fmt.Sprintf("$%.0f - $%.0f", basePrice*0.9, basePrice*1.1)
}

func (p *proposalGenerator) generateEstimatedTimeline(inquiry *domain.Inquiry, services []interfaces.ProposedService) string {
	duration := p.getDefaultDuration(inquiry.Services)
	return fmt.Sprintf("%d - %d weeks", duration/7-2, duration/7+2)
}

func (p *proposalGenerator) generateSolutionBenefits(inquiry *domain.Inquiry, services []interfaces.ProposedService) []string {
	return []string{
		"Improved operational efficiency",
		"Reduced infrastructure costs",
		"Enhanced scalability and flexibility",
		"Better security and compliance",
		"Faster time to market for new features",
	}
}

// ptrTime returns a pointer to a time.Time value
func ptrTime(t time.Time) *time.Time {
	return &t
}

// convertProposalResourceEstimateToResourceEstimate maps ProposalResourceEstimate to ResourceEstimate for embedding in DetailedResources.
