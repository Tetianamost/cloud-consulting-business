package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// RoadmapGeneratorService implements the RoadmapGenerator interface
type RoadmapGeneratorService struct {
	bedrockService interfaces.BedrockService
	templates      map[string]*interfaces.RoadmapTemplate
}

// NewRoadmapGeneratorService creates a new roadmap generator service
func NewRoadmapGeneratorService(bedrockService interfaces.BedrockService) *RoadmapGeneratorService {
	service := &RoadmapGeneratorService{
		bedrockService: bedrockService,
		templates:      make(map[string]*interfaces.RoadmapTemplate),
	}
	
	// Initialize default templates
	service.initializeDefaultTemplates()
	
	return service
}

// GenerateImplementationRoadmap generates a complete implementation roadmap
func (r *RoadmapGeneratorService) GenerateImplementationRoadmap(ctx context.Context, inquiry *domain.Inquiry, report *domain.Report) (*interfaces.ImplementationRoadmap, error) {
	if inquiry == nil {
		return nil, fmt.Errorf("inquiry cannot be nil")
	}
	
	// Determine project type and constraints from inquiry
	projectType := r.determineProjectType(inquiry)
	constraints := r.extractConstraints(inquiry)
	
	// Generate phases based on project requirements
	phases, err := r.GeneratePhases(ctx, r.extractRequirements(inquiry), constraints)
	if err != nil {
		return nil, fmt.Errorf("failed to generate phases: %w", err)
	}
	
	// Calculate dependencies between phases
	dependencies, err := r.CalculateDependencies(ctx, phases)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate dependencies: %w", err)
	}
	
	// Generate milestones for all phases
	_, err = r.GenerateMilestones(ctx, phases)
	if err != nil {
		return nil, fmt.Errorf("failed to generate milestones: %w", err)
	}
	
	// Estimate resources
	resourceEstimate, err := r.EstimateResources(ctx, phases)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate resources: %w", err)
	}
	
	// Create the roadmap
	roadmap := &interfaces.ImplementationRoadmap{
		ID:              uuid.New().String(),
		InquiryID:       inquiry.ID,
		Title:           r.generateTitle(inquiry, projectType),
		Overview:        r.generateOverview(inquiry, phases),
		TotalDuration:   r.calculateTotalDuration(phases),
		EstimatedCost:   resourceEstimate.TotalCost,
		Phases:          phases,
		Dependencies:    dependencies,
		Risks:           r.identifyRisks(phases, constraints),
		SuccessMetrics:  r.generateSuccessMetrics(inquiry, projectType),
		ProjectType:     projectType,
		CloudProviders:  r.extractCloudProviders(inquiry),
		IndustryContext: r.extractIndustryContext(inquiry),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	// Validate the roadmap
	validation, err := r.ValidateRoadmap(ctx, roadmap)
	if err != nil {
		return nil, fmt.Errorf("failed to validate roadmap: %w", err)
	}
	
	if !validation.IsValid {
		return nil, fmt.Errorf("generated roadmap is invalid: %v", validation.Errors)
	}
	
	return roadmap, nil
}

// GeneratePhases generates roadmap phases based on requirements and constraints
func (r *RoadmapGeneratorService) GeneratePhases(ctx context.Context, requirements []string, constraints *interfaces.ProjectConstraints) ([]interfaces.RoadmapPhase, error) {
	// Select appropriate template
	template := r.selectTemplate(requirements, constraints)
	
	var phases []interfaces.RoadmapPhase
	
	for i, phaseTemplate := range template.PhaseTemplates {
		phase := interfaces.RoadmapPhase{
			ID:          uuid.New().String(),
			Name:        phaseTemplate.Name,
			Description: phaseTemplate.Description,
			Duration:    phaseTemplate.EstimatedDuration,
			Priority:    r.calculatePhasePriority(i, len(template.PhaseTemplates)),
			RiskLevel:   r.assessPhaseRisk(phaseTemplate, constraints),
		}
		
		// Generate tasks for this phase
		tasks, err := r.generateTasks(ctx, phaseTemplate.TaskTemplates, requirements, constraints)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tasks for phase %s: %w", phase.Name, err)
		}
		phase.Tasks = tasks
		
		// Generate deliverables
		deliverables := r.generateDeliverables(phaseTemplate.DeliverableTemplates)
		phase.Deliverables = deliverables
		
		// Generate milestones for this phase
		milestones := r.generatePhaseMilestones(phaseTemplate.MilestoneTemplates, phase.ID)
		phase.Milestones = milestones
		
		// Set resource requirements
		phase.ResourceRequirements = r.adaptResourceRequirements(phaseTemplate.ResourceTemplate, constraints)
		
		// Calculate estimated cost for this phase
		phase.EstimatedCost = r.calculatePhaseCost(phase.ResourceRequirements, phase.Duration)
		
		// Set prerequisites
		if i > 0 {
			phase.Prerequisites = []string{phases[i-1].ID}
		}
		
		phases = append(phases, phase)
	}
	
	return phases, nil
}

// EstimateResources estimates resource requirements for all phases
func (r *RoadmapGeneratorService) EstimateResources(ctx context.Context, phases []interfaces.RoadmapPhase) (*interfaces.ResourceEstimate, error) {
	estimate := &interfaces.ResourceEstimate{
		RoleBreakdown:       make(map[string]int),
		PhaseBreakdown:      make(map[string]int),
		CostBreakdown:       make(map[string]string),
		ResourceUtilization: make(map[string]float64),
	}
	
	totalHours := 0
	totalCostFloat := 0.0
	
	for _, phase := range phases {
		phaseHours := 0
		
		// Calculate hours from tasks
		for _, task := range phase.Tasks {
			phaseHours += task.EstimatedHours
			totalHours += task.EstimatedHours
			
			// Track role breakdown
			for _, skill := range task.SkillsRequired {
				estimate.RoleBreakdown[skill] += task.EstimatedHours
			}
		}
		
		estimate.PhaseBreakdown[phase.Name] = phaseHours
		
		// Parse and add phase cost
		if phaseCost := r.parseCostString(phase.EstimatedCost); phaseCost > 0 {
			totalCostFloat += phaseCost
			estimate.CostBreakdown[phase.Name] = phase.EstimatedCost
		}
	}
	
	estimate.TotalHours = totalHours
	estimate.TotalCost = fmt.Sprintf("$%.2f", totalCostFloat)
	
	// Calculate resource utilization (simplified)
	for role, hours := range estimate.RoleBreakdown {
		utilization := float64(hours) / float64(totalHours) * 100
		estimate.ResourceUtilization[role] = utilization
	}
	
	return estimate, nil
}

// CalculateDependencies calculates dependencies between phases
func (r *RoadmapGeneratorService) CalculateDependencies(ctx context.Context, phases []interfaces.RoadmapPhase) ([]interfaces.Dependency, error) {
	var dependencies []interfaces.Dependency
	
	// Create sequential dependencies between phases
	for i := 1; i < len(phases); i++ {
		dependency := interfaces.Dependency{
			ID:          uuid.New().String(),
			FromID:      phases[i-1].ID,
			ToID:        phases[i].ID,
			Type:        "finish_to_start",
			Description: fmt.Sprintf("Phase '%s' must complete before '%s' can start", phases[i-1].Name, phases[i].Name),
			IsCritical:  true,
		}
		dependencies = append(dependencies, dependency)
	}
	
	// Add task-level dependencies within phases
	for _, phase := range phases {
		for _, task := range phase.Tasks {
			for _, depTaskID := range task.Dependencies {
				dependency := interfaces.Dependency{
					ID:          uuid.New().String(),
					FromID:      depTaskID,
					ToID:        task.ID,
					Type:        "finish_to_start",
					Description: fmt.Sprintf("Task dependency within phase '%s'", phase.Name),
					IsCritical:  task.Priority == "high" || task.Priority == "critical",
				}
				dependencies = append(dependencies, dependency)
			}
		}
	}
	
	return dependencies, nil
}

// GenerateMilestones generates milestones for all phases
func (r *RoadmapGeneratorService) GenerateMilestones(ctx context.Context, phases []interfaces.RoadmapPhase) ([]interfaces.Milestone, error) {
	var allMilestones []interfaces.Milestone
	
	for _, phase := range phases {
		// Add phase completion milestone
		milestone := interfaces.Milestone{
			ID:           uuid.New().String(),
			Name:         fmt.Sprintf("%s Completion", phase.Name),
			Description:  fmt.Sprintf("All tasks and deliverables for %s are completed", phase.Name),
			Type:         "phase_completion",
			Criteria:     []string{"All tasks completed", "All deliverables reviewed and approved"},
			Status:       "pending",
			Stakeholders: []string{"Project Manager", "Technical Lead"},
			Importance:   "high",
		}
		allMilestones = append(allMilestones, milestone)
		
		// Add existing phase milestones
		allMilestones = append(allMilestones, phase.Milestones...)
	}
	
	// Add project-level milestones
	projectMilestones := []interfaces.Milestone{
		{
			ID:           uuid.New().String(),
			Name:         "Project Kickoff",
			Description:  "Project officially starts with all stakeholders aligned",
			Type:         "checkpoint",
			Criteria:     []string{"Team assembled", "Requirements confirmed", "Timeline approved"},
			Status:       "pending",
			Stakeholders: []string{"Project Sponsor", "Project Manager", "Technical Lead"},
			Importance:   "critical",
		},
		{
			ID:           uuid.New().String(),
			Name:         "Go-Live",
			Description:  "Solution is deployed and operational in production",
			Type:         "go_live",
			Criteria:     []string{"Production deployment successful", "User acceptance testing passed", "Monitoring in place"},
			Status:       "pending",
			Stakeholders: []string{"Business Owner", "Operations Team", "End Users"},
			Importance:   "critical",
		},
	}
	
	allMilestones = append(allMilestones, projectMilestones...)
	
	return allMilestones, nil
}

// ValidateRoadmap validates the generated roadmap
func (r *RoadmapGeneratorService) ValidateRoadmap(ctx context.Context, roadmap *interfaces.ImplementationRoadmap) (*interfaces.ValidationResult, error) {
	result := &interfaces.ValidationResult{
		IsValid: true,
	}
	
	// Validate basic structure
	if roadmap.Title == "" {
		result.Errors = append(result.Errors, "Roadmap title is required")
		result.IsValid = false
	}
	
	if len(roadmap.Phases) == 0 {
		result.Errors = append(result.Errors, "Roadmap must have at least one phase")
		result.IsValid = false
	}
	
	// Validate phases
	for i, phase := range roadmap.Phases {
		if phase.Name == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Phase %d is missing a name", i+1))
			result.IsValid = false
		}
		
		if len(phase.Tasks) == 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Phase '%s' has no tasks", phase.Name))
		}
		
		if phase.Duration == "" {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Phase '%s' has no duration estimate", phase.Name))
		}
	}
	
	// Validate dependencies
	phaseIDs := make(map[string]bool)
	for _, phase := range roadmap.Phases {
		phaseIDs[phase.ID] = true
	}
	
	for _, dep := range roadmap.Dependencies {
		if !phaseIDs[dep.FromID] && !r.isTaskID(dep.FromID, roadmap.Phases) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Dependency references unknown ID: %s", dep.FromID))
		}
		if !phaseIDs[dep.ToID] && !r.isTaskID(dep.ToID, roadmap.Phases) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Dependency references unknown ID: %s", dep.ToID))
		}
	}
	
	// Calculate quality scores
	result.QualityScore = r.calculateQualityScore(roadmap)
	result.CompletenessScore = r.calculateCompletenessScore(roadmap)
	
	// Add suggestions
	if result.QualityScore < 0.7 {
		result.Suggestions = append(result.Suggestions, "Consider adding more detailed task descriptions and acceptance criteria")
	}
	
	if result.CompletenessScore < 0.8 {
		result.Suggestions = append(result.Suggestions, "Consider adding more milestones and deliverables to track progress")
	}
	
	return result, nil
}

// Helper methods

func (r *RoadmapGeneratorService) initializeDefaultTemplates() {
	// Migration template
	migrationTemplate := &interfaces.RoadmapTemplate{
		ID:              "migration-template",
		Name:            "Cloud Migration Template",
		Description:     "Standard template for cloud migration projects",
		ProjectType:     "migration",
		IndustryContext: "general",
		PhaseTemplates: []interfaces.PhaseTemplate{
			{
				Name:              "Assessment & Planning",
				Description:       "Assess current infrastructure and plan migration strategy",
				EstimatedDuration: "2-4 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Infrastructure Assessment",
						Description:        "Analyze current infrastructure and dependencies",
						EstimatedHours:     40,
						SkillsRequired:     []string{"Cloud Architect", "Infrastructure Engineer"},
						Priority:           "high",
						CompletionCriteria: []string{"Current state documented", "Dependencies mapped"},
					},
					{
						Name:               "Migration Strategy Design",
						Description:        "Design migration approach and timeline",
						EstimatedHours:     32,
						SkillsRequired:     []string{"Cloud Architect", "Project Manager"},
						Priority:           "high",
						CompletionCriteria: []string{"Migration strategy approved", "Risk assessment completed"},
					},
				},
				DeliverableTemplates: []interfaces.DeliverableTemplate{
					{
						Name:        "Assessment Report",
						Description: "Comprehensive analysis of current infrastructure",
						Type:        "document",
						Format:      "PDF",
						Reviewers:   []string{"Technical Lead", "Business Stakeholder"},
					},
				},
				MilestoneTemplates: []interfaces.MilestoneTemplate{
					{
						Name:         "Assessment Complete",
						Description:  "Infrastructure assessment and migration plan approved",
						Type:         "checkpoint",
						Criteria:     []string{"Assessment report approved", "Migration strategy signed off"},
						Stakeholders: []string{"Project Sponsor", "Technical Lead"},
						Importance:   "high",
					},
				},
			},
			{
				Name:              "Environment Setup",
				Description:       "Set up target cloud environment and tools",
				EstimatedDuration: "1-2 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Cloud Environment Provisioning",
						Description:        "Provision target cloud infrastructure",
						EstimatedHours:     24,
						SkillsRequired:     []string{"Cloud Engineer", "DevOps Engineer"},
						Priority:           "high",
						CompletionCriteria: []string{"Infrastructure provisioned", "Security configured"},
					},
				},
			},
			{
				Name:              "Migration Execution",
				Description:       "Execute the migration according to the plan",
				EstimatedDuration: "3-6 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Data Migration",
						Description:        "Migrate data to target environment",
						EstimatedHours:     60,
						SkillsRequired:     []string{"Data Engineer", "Database Administrator"},
						Priority:           "critical",
						CompletionCriteria: []string{"Data migrated successfully", "Data integrity verified"},
					},
				},
			},
			{
				Name:              "Testing & Validation",
				Description:       "Test migrated systems and validate functionality",
				EstimatedDuration: "1-2 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "System Testing",
						Description:        "Comprehensive testing of migrated systems",
						EstimatedHours:     40,
						SkillsRequired:     []string{"QA Engineer", "System Administrator"},
						Priority:           "high",
						CompletionCriteria: []string{"All tests passed", "Performance validated"},
					},
				},
			},
			{
				Name:              "Go-Live & Support",
				Description:       "Deploy to production and provide initial support",
				EstimatedDuration: "1 week",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Production Deployment",
						Description:        "Deploy to production environment",
						EstimatedHours:     16,
						SkillsRequired:     []string{"DevOps Engineer", "System Administrator"},
						Priority:           "critical",
						CompletionCriteria: []string{"Production deployment successful", "Monitoring active"},
					},
				},
			},
		},
	}
	
	r.templates["migration"] = migrationTemplate
	
	// Add more templates for other project types
	r.templates["optimization"] = r.createOptimizationTemplate()
	r.templates["assessment"] = r.createAssessmentTemplate()
	r.templates["architecture"] = r.createArchitectureTemplate()
}

func (r *RoadmapGeneratorService) createOptimizationTemplate() *interfaces.RoadmapTemplate {
	return &interfaces.RoadmapTemplate{
		ID:              "optimization-template",
		Name:            "Cloud Optimization Template",
		Description:     "Template for cloud cost and performance optimization projects",
		ProjectType:     "optimization",
		IndustryContext: "general",
		PhaseTemplates: []interfaces.PhaseTemplate{
			{
				Name:              "Current State Analysis",
				Description:       "Analyze current cloud usage and costs",
				EstimatedDuration: "1-2 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Cost Analysis",
						Description:        "Analyze current cloud spending patterns",
						EstimatedHours:     24,
						SkillsRequired:     []string{"Cloud Architect", "FinOps Specialist"},
						Priority:           "high",
						CompletionCriteria: []string{"Cost breakdown completed", "Optimization opportunities identified"},
					},
				},
			},
			{
				Name:              "Optimization Implementation",
				Description:       "Implement identified optimizations",
				EstimatedDuration: "2-4 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Resource Right-sizing",
						Description:        "Optimize resource allocation and sizing",
						EstimatedHours:     32,
						SkillsRequired:     []string{"Cloud Engineer", "Performance Engineer"},
						Priority:           "high",
						CompletionCriteria: []string{"Resources optimized", "Performance maintained"},
					},
				},
			},
		},
	}
}

func (r *RoadmapGeneratorService) createAssessmentTemplate() *interfaces.RoadmapTemplate {
	return &interfaces.RoadmapTemplate{
		ID:              "assessment-template",
		Name:            "Cloud Assessment Template",
		Description:     "Template for comprehensive cloud readiness assessments",
		ProjectType:     "assessment",
		IndustryContext: "general",
		PhaseTemplates: []interfaces.PhaseTemplate{
			{
				Name:              "Discovery & Analysis",
				Description:       "Discover and analyze current infrastructure",
				EstimatedDuration: "2-3 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Infrastructure Discovery",
						Description:        "Document current infrastructure and applications",
						EstimatedHours:     40,
						SkillsRequired:     []string{"Cloud Architect", "System Analyst"},
						Priority:           "high",
						CompletionCriteria: []string{"Infrastructure documented", "Dependencies mapped"},
					},
				},
			},
		},
	}
}

func (r *RoadmapGeneratorService) createArchitectureTemplate() *interfaces.RoadmapTemplate {
	return &interfaces.RoadmapTemplate{
		ID:              "architecture-template",
		Name:            "Architecture Review Template",
		Description:     "Template for cloud architecture review and design projects",
		ProjectType:     "architecture",
		IndustryContext: "general",
		PhaseTemplates: []interfaces.PhaseTemplate{
			{
				Name:              "Architecture Review",
				Description:       "Review current architecture and identify improvements",
				EstimatedDuration: "1-2 weeks",
				TaskTemplates: []interfaces.TaskTemplate{
					{
						Name:               "Architecture Analysis",
						Description:        "Analyze current architecture against best practices",
						EstimatedHours:     32,
						SkillsRequired:     []string{"Solution Architect", "Cloud Architect"},
						Priority:           "high",
						CompletionCriteria: []string{"Architecture reviewed", "Recommendations documented"},
					},
				},
			},
		},
	}
}

func (r *RoadmapGeneratorService) determineProjectType(inquiry *domain.Inquiry) string {
	services := strings.Join(inquiry.Services, " ")
	message := strings.ToLower(inquiry.Message)
	
	if strings.Contains(services, "migration") || strings.Contains(message, "migrate") {
		return "migration"
	}
	if strings.Contains(services, "optimization") || strings.Contains(message, "optimize") {
		return "optimization"
	}
	if strings.Contains(services, "assessment") || strings.Contains(message, "assess") {
		return "assessment"
	}
	if strings.Contains(services, "architecture") || strings.Contains(message, "architecture") {
		return "architecture"
	}
	
	return "general"
}

func (r *RoadmapGeneratorService) extractConstraints(inquiry *domain.Inquiry) *interfaces.ProjectConstraints {
	// Extract constraints from inquiry message and context
	// This is a simplified implementation - in practice, this would use NLP or structured input
	
	constraints := &interfaces.ProjectConstraints{
		RiskTolerance:   "medium",
		CloudProviders:  []string{"AWS"}, // Default
		TechnologyStack: []string{},
	}
	
	message := strings.ToLower(inquiry.Message)
	
	// Extract budget constraints
	if strings.Contains(message, "budget") || strings.Contains(message, "cost") {
		constraints.Budget = "medium"
	}
	
	// Extract timeline constraints
	if strings.Contains(message, "urgent") || strings.Contains(message, "asap") {
		constraints.Timeline = "aggressive"
	} else if strings.Contains(message, "flexible") {
		constraints.Timeline = "flexible"
	} else {
		constraints.Timeline = "standard"
	}
	
	// Extract cloud provider preferences
	if strings.Contains(message, "azure") {
		constraints.CloudProviders = append(constraints.CloudProviders, "Azure")
	}
	if strings.Contains(message, "gcp") || strings.Contains(message, "google") {
		constraints.CloudProviders = append(constraints.CloudProviders, "GCP")
	}
	
	return constraints
}

func (r *RoadmapGeneratorService) extractRequirements(inquiry *domain.Inquiry) []string {
	requirements := []string{}
	
	// Extract requirements from services
	for _, service := range inquiry.Services {
		requirements = append(requirements, service)
	}
	
	// Extract additional requirements from message
	message := strings.ToLower(inquiry.Message)
	
	if strings.Contains(message, "compliance") || strings.Contains(message, "hipaa") || strings.Contains(message, "pci") {
		requirements = append(requirements, "compliance")
	}
	if strings.Contains(message, "security") {
		requirements = append(requirements, "security")
	}
	if strings.Contains(message, "performance") {
		requirements = append(requirements, "performance")
	}
	if strings.Contains(message, "scalability") || strings.Contains(message, "scale") {
		requirements = append(requirements, "scalability")
	}
	
	return requirements
}

func (r *RoadmapGeneratorService) selectTemplate(requirements []string, constraints *interfaces.ProjectConstraints) *interfaces.RoadmapTemplate {
	// Determine the best template based on requirements
	for _, req := range requirements {
		if template, exists := r.templates[req]; exists {
			return template
		}
	}
	
	// Default to migration template
	return r.templates["migration"]
}

func (r *RoadmapGeneratorService) generateTasks(ctx context.Context, taskTemplates []interfaces.TaskTemplate, requirements []string, constraints *interfaces.ProjectConstraints) ([]interfaces.Task, error) {
	var tasks []interfaces.Task
	
	for _, template := range taskTemplates {
		task := interfaces.Task{
			ID:                 uuid.New().String(),
			Name:               template.Name,
			Description:        template.Description,
			EstimatedHours:     template.EstimatedHours,
			SkillsRequired:     template.SkillsRequired,
			Priority:           template.Priority,
			DocumentationLinks: template.DocumentationLinks,
			CompletionCriteria: template.CompletionCriteria,
			Status:             "not_started",
		}
		
		// Adjust task based on constraints
		if constraints.Timeline == "aggressive" {
			task.EstimatedHours = int(float64(task.EstimatedHours) * 0.8) // Reduce time estimate for aggressive timeline
		}
		
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *RoadmapGeneratorService) generateDeliverables(templates []interfaces.DeliverableTemplate) []interfaces.Deliverable {
	var deliverables []interfaces.Deliverable
	
	for _, template := range templates {
		deliverable := interfaces.Deliverable{
			ID:          uuid.New().String(),
			Name:        template.Name,
			Description: template.Description,
			Type:        template.Type,
			Format:      template.Format,
			Reviewers:   template.Reviewers,
			Status:      "not_started",
		}
		deliverables = append(deliverables, deliverable)
	}
	
	return deliverables
}

func (r *RoadmapGeneratorService) generatePhaseMilestones(templates []interfaces.MilestoneTemplate, phaseID string) []interfaces.Milestone {
	var milestones []interfaces.Milestone
	
	for _, template := range templates {
		milestone := interfaces.Milestone{
			ID:           uuid.New().String(),
			Name:         template.Name,
			Description:  template.Description,
			Type:         template.Type,
			Criteria:     template.Criteria,
			Status:       "pending",
			Stakeholders: template.Stakeholders,
			Importance:   template.Importance,
		}
		milestones = append(milestones, milestone)
	}
	
	return milestones
}

func (r *RoadmapGeneratorService) adaptResourceRequirements(template interfaces.ResourceRequirements, constraints *interfaces.ProjectConstraints) interfaces.ResourceRequirements {
	adapted := template
	
	// Adjust based on constraints
	if constraints.TeamSize > 0 && constraints.TeamSize < 5 {
		// For small teams, combine roles
		adapted.TechnicalRoles = r.consolidateRoles(adapted.TechnicalRoles)
	}
	
	return adapted
}

func (r *RoadmapGeneratorService) consolidateRoles(roles []interfaces.Role) []interfaces.Role {
	// Simplified role consolidation for small teams
	if len(roles) <= 2 {
		return roles
	}
	
	// Combine similar roles
	consolidated := []interfaces.Role{}
	for i, role := range roles {
		if i < 2 {
			consolidated = append(consolidated, role)
		} else {
			// Merge with existing roles
			consolidated[i%2].SkillsRequired = append(consolidated[i%2].SkillsRequired, role.SkillsRequired...)
		}
	}
	
	return consolidated
}

func (r *RoadmapGeneratorService) calculatePhaseCost(requirements interfaces.ResourceRequirements, duration string) string {
	// Simplified cost calculation
	baseCost := 10000.0 // Base cost per phase
	
	// Adjust based on roles
	roleCost := float64(len(requirements.TechnicalRoles)+len(requirements.BusinessRoles)) * 5000.0
	
	// Adjust based on duration (extract weeks from duration string)
	weeks := r.extractWeeksFromDuration(duration)
	durationMultiplier := float64(weeks) * 0.5
	
	totalCost := baseCost + roleCost + (baseCost * durationMultiplier)
	
	return fmt.Sprintf("$%.2f", totalCost)
}

func (r *RoadmapGeneratorService) extractWeeksFromDuration(duration string) int {
	// Simple extraction of weeks from duration string like "2-4 weeks"
	if strings.Contains(duration, "week") {
		parts := strings.Fields(duration)
		for _, part := range parts {
			if strings.Contains(part, "-") {
				// Take the higher number from range
				rangeParts := strings.Split(part, "-")
				if len(rangeParts) == 2 {
					if weeks, err := strconv.Atoi(rangeParts[1]); err == nil {
						return weeks
					}
				}
			} else if weeks, err := strconv.Atoi(part); err == nil {
				return weeks
			}
		}
	}
	return 2 // Default
}

func (r *RoadmapGeneratorService) calculatePhasePriority(index, total int) string {
	if index == 0 {
		return "critical"
	}
	if index < total/2 {
		return "high"
	}
	return "medium"
}

func (r *RoadmapGeneratorService) assessPhaseRisk(template interfaces.PhaseTemplate, constraints *interfaces.ProjectConstraints) string {
	// Simplified risk assessment
	if constraints.Timeline == "aggressive" {
		return "high"
	}
	if len(template.TaskTemplates) > 5 {
		return "medium"
	}
	return "low"
}

func (r *RoadmapGeneratorService) generateTitle(inquiry *domain.Inquiry, projectType string) string {
	return fmt.Sprintf("%s Implementation Roadmap - %s", 
		strings.Title(projectType), 
		inquiry.Company)
}

func (r *RoadmapGeneratorService) generateOverview(inquiry *domain.Inquiry, phases []interfaces.RoadmapPhase) string {
	return fmt.Sprintf("This implementation roadmap outlines the %d-phase approach for %s's cloud %s project. The roadmap includes detailed tasks, milestones, and resource requirements to ensure successful project delivery.",
		len(phases),
		inquiry.Company,
		strings.Join(inquiry.Services, " and "))
}

func (r *RoadmapGeneratorService) calculateTotalDuration(phases []interfaces.RoadmapPhase) string {
	totalWeeks := 0
	
	for _, phase := range phases {
		weeks := r.extractWeeksFromDuration(phase.Duration)
		totalWeeks += weeks
	}
	
	if totalWeeks <= 4 {
		return fmt.Sprintf("%d weeks", totalWeeks)
	} else {
		months := totalWeeks / 4
		remainingWeeks := totalWeeks % 4
		if remainingWeeks == 0 {
			return fmt.Sprintf("%d months", months)
		}
		return fmt.Sprintf("%d months %d weeks", months, remainingWeeks)
	}
}

func (r *RoadmapGeneratorService) identifyRisks(phases []interfaces.RoadmapPhase, constraints *interfaces.ProjectConstraints) []string {
	risks := []string{}
	
	// Common risks
	risks = append(risks, "Resource availability constraints")
	risks = append(risks, "Technical complexity underestimation")
	risks = append(risks, "Stakeholder alignment challenges")
	
	// Timeline-based risks
	if constraints.Timeline == "aggressive" {
		risks = append(risks, "Aggressive timeline may impact quality")
		risks = append(risks, "Limited time for thorough testing")
	}
	
	// Budget-based risks
	if constraints.Budget == "tight" {
		risks = append(risks, "Budget constraints may limit scope")
	}
	
	// Phase-specific risks
	for _, phase := range phases {
		if phase.RiskLevel == "high" {
			risks = append(risks, fmt.Sprintf("High risk in %s phase", phase.Name))
		}
	}
	
	return risks
}

func (r *RoadmapGeneratorService) generateSuccessMetrics(inquiry *domain.Inquiry, projectType string) []string {
	metrics := []string{
		"Project delivered on time and within budget",
		"All stakeholder requirements met",
		"Zero critical issues in production",
		"User acceptance criteria satisfied",
	}
	
	// Project-type specific metrics
	switch projectType {
	case "migration":
		metrics = append(metrics, "Zero data loss during migration")
		metrics = append(metrics, "Performance maintained or improved")
	case "optimization":
		metrics = append(metrics, "Cost reduction targets achieved")
		metrics = append(metrics, "Performance improvements measured")
	case "assessment":
		metrics = append(metrics, "Comprehensive assessment report delivered")
		metrics = append(metrics, "Actionable recommendations provided")
	}
	
	return metrics
}

func (r *RoadmapGeneratorService) extractCloudProviders(inquiry *domain.Inquiry) []string {
	providers := []string{"AWS"} // Default
	
	message := strings.ToLower(inquiry.Message)
	if strings.Contains(message, "azure") {
		providers = append(providers, "Azure")
	}
	if strings.Contains(message, "gcp") || strings.Contains(message, "google") {
		providers = append(providers, "GCP")
	}
	
	return providers
}

func (r *RoadmapGeneratorService) extractIndustryContext(inquiry *domain.Inquiry) string {
	message := strings.ToLower(inquiry.Message)
	
	if strings.Contains(message, "healthcare") || strings.Contains(message, "hipaa") {
		return "healthcare"
	}
	if strings.Contains(message, "financial") || strings.Contains(message, "banking") {
		return "financial"
	}
	if strings.Contains(message, "retail") || strings.Contains(message, "ecommerce") {
		return "retail"
	}
	
	return "general"
}

func (r *RoadmapGeneratorService) parseCostString(costStr string) float64 {
	// Remove $ and parse
	cleaned := strings.ReplaceAll(costStr, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	
	if cost, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return cost
	}
	
	return 0.0
}

func (r *RoadmapGeneratorService) isTaskID(id string, phases []interfaces.RoadmapPhase) bool {
	for _, phase := range phases {
		for _, task := range phase.Tasks {
			if task.ID == id {
				return true
			}
		}
	}
	return false
}

func (r *RoadmapGeneratorService) calculateQualityScore(roadmap *interfaces.ImplementationRoadmap) float64 {
	score := 0.0
	maxScore := 0.0
	
	// Check roadmap completeness
	if roadmap.Title != "" {
		score += 1.0
	}
	maxScore += 1.0
	
	if roadmap.Overview != "" {
		score += 1.0
	}
	maxScore += 1.0
	
	if len(roadmap.Phases) > 0 {
		score += 1.0
	}
	maxScore += 1.0
	
	// Check phase quality
	for _, phase := range roadmap.Phases {
		if phase.Description != "" {
			score += 0.5
		}
		maxScore += 0.5
		
		if len(phase.Tasks) > 0 {
			score += 0.5
		}
		maxScore += 0.5
		
		if len(phase.Deliverables) > 0 {
			score += 0.3
		}
		maxScore += 0.3
	}
	
	if maxScore == 0 {
		return 0.0
	}
	
	return score / maxScore
}

func (r *RoadmapGeneratorService) calculateCompletenessScore(roadmap *interfaces.ImplementationRoadmap) float64 {
	score := 0.0
	maxScore := 0.0
	
	// Check for essential elements
	elements := []bool{
		roadmap.TotalDuration != "",
		roadmap.EstimatedCost != "",
		len(roadmap.Dependencies) > 0,
		len(roadmap.Risks) > 0,
		len(roadmap.SuccessMetrics) > 0,
	}
	
	for _, hasElement := range elements {
		if hasElement {
			score += 1.0
		}
		maxScore += 1.0
	}
	
	// Check phase completeness
	for _, phase := range roadmap.Phases {
		phaseElements := []bool{
			phase.Duration != "",
			phase.EstimatedCost != "",
			len(phase.Tasks) > 0,
			len(phase.Deliverables) > 0,
			len(phase.Milestones) > 0,
		}
		
		for _, hasElement := range phaseElements {
			if hasElement {
				score += 0.2
			}
			maxScore += 0.2
		}
	}
	
	if maxScore == 0 {
		return 0.0
	}
	
	return score / maxScore
}