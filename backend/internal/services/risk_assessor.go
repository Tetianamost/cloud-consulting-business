package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// RiskAssessorService implements the RiskAssessor interface
type RiskAssessorService struct {
	knowledgeBase interfaces.KnowledgeBase
	docLibrary    interfaces.DocumentationLibrary
}

// NewRiskAssessorService creates a new risk assessor service
func NewRiskAssessorService(kb interfaces.KnowledgeBase, docLib interfaces.DocumentationLibrary) interfaces.RiskAssessor {
	return &RiskAssessorService{
		knowledgeBase: kb,
		docLibrary:    docLib,
	}
}

// AssessRisks performs a comprehensive risk assessment for a given inquiry and solution
func (r *RiskAssessorService) AssessRisks(ctx context.Context, inquiry *domain.Inquiry, solution *interfaces.ProposedSolution) (*interfaces.RiskAssessment, error) {
	assessment := &interfaces.RiskAssessment{
		ID:        uuid.New().String(),
		InquiryID: inquiry.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Identify technical risks
	technicalRisks, err := r.identifyTechnicalRisks(ctx, inquiry, solution)
	if err != nil {
		return nil, fmt.Errorf("failed to identify technical risks: %w", err)
	}
	assessment.TechnicalRisks = technicalRisks

	// Identify security risks
	securityRisks, err := r.IdentifySecurityRisks(ctx, solution.Architecture)
	if err != nil {
		return nil, fmt.Errorf("failed to identify security risks: %w", err)
	}
	assessment.SecurityRisks = securityRisks

	// Identify compliance risks
	industry := r.extractIndustryFromInquiry(inquiry)
	complianceRisks, err := r.EvaluateComplianceRisks(ctx, industry, solution)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate compliance risks: %w", err)
	}
	assessment.ComplianceRisks = complianceRisks

	// Identify business risks
	businessRisks, err := r.identifyBusinessRisks(ctx, inquiry, solution)
	if err != nil {
		return nil, fmt.Errorf("failed to identify business risks: %w", err)
	}
	assessment.BusinessRisks = businessRisks

	// Collect all risks for mitigation strategy generation
	allRisks := r.collectAllRisks(assessment)

	// Generate mitigation strategies
	mitigationStrategies, err := r.GenerateMitigationStrategies(ctx, allRisks)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mitigation strategies: %w", err)
	}
	assessment.MitigationStrategies = mitigationStrategies

	// Calculate overall risk level
	assessment.OverallRiskLevel = r.calculateOverallRiskLevel(allRisks)

	// Generate recommended actions
	assessment.RecommendedActions = r.generateRecommendedActions(assessment)

	return assessment, nil
}

// IdentifySecurityRisks identifies security risks in the proposed architecture
func (r *RiskAssessorService) IdentifySecurityRisks(ctx context.Context, architecture *interfaces.Architecture) ([]*interfaces.SecurityRisk, error) {
	var risks []*interfaces.SecurityRisk

	if architecture == nil {
		return risks, nil
	}

	// Check for common security risks
	risks = append(risks, r.assessDataEncryptionRisks(architecture)...)
	risks = append(risks, r.assessNetworkSecurityRisks(architecture)...)
	risks = append(risks, r.assessAccessControlRisks(architecture)...)
	risks = append(risks, r.assessDataStorageSecurityRisks(architecture)...)

	return risks, nil
}

// EvaluateComplianceRisks evaluates compliance risks based on industry and solution
func (r *RiskAssessorService) EvaluateComplianceRisks(ctx context.Context, industry string, solution *interfaces.ProposedSolution) ([]*interfaces.ComplianceRisk, error) {
	var risks []*interfaces.ComplianceRisk

	// Get compliance requirements for the industry
	complianceFrameworks := r.getComplianceFrameworksForIndustry(industry)

	for _, framework := range complianceFrameworks {
		frameworkRisks := r.assessComplianceFrameworkRisks(framework, solution)
		risks = append(risks, frameworkRisks...)
	}

	return risks, nil
}

// GenerateMitigationStrategies generates mitigation strategies for identified risks
func (r *RiskAssessorService) GenerateMitigationStrategies(ctx context.Context, risks []*interfaces.Risk) ([]*interfaces.MitigationStrategy, error) {
	var strategies []*interfaces.MitigationStrategy

	for _, risk := range risks {
		strategy := r.generateMitigationStrategy(risk)
		if strategy != nil {
			strategies = append(strategies, strategy)
		}
	}

	return strategies, nil
}

// Private helper methods

func (r *RiskAssessorService) identifyTechnicalRisks(ctx context.Context, inquiry *domain.Inquiry, solution *interfaces.ProposedSolution) ([]*interfaces.TechnicalRisk, error) {
	var risks []*interfaces.TechnicalRisk

	// Assess single points of failure
	risks = append(risks, r.assessSinglePointsOfFailure(solution)...)

	// Assess scalability risks
	risks = append(risks, r.assessScalabilityRisks(solution)...)

	// Assess performance risks
	risks = append(risks, r.assessPerformanceRisks(solution)...)

	// Assess dependency risks
	risks = append(risks, r.assessDependencyRisks(solution)...)

	// Assess migration risks if it's a migration project
	if r.isMigrationProject(inquiry) {
		risks = append(risks, r.assessMigrationRisks(solution)...)
	}

	return risks, nil
}

func (r *RiskAssessorService) identifyBusinessRisks(ctx context.Context, inquiry *domain.Inquiry, solution *interfaces.ProposedSolution) ([]*interfaces.BusinessRisk, error) {
	var risks []*interfaces.BusinessRisk

	// Assess cost overrun risks
	risks = append(risks, r.assessCostRisks(solution)...)

	// Assess timeline risks
	risks = append(risks, r.assessTimelineRisks(solution)...)

	// Assess vendor lock-in risks
	risks = append(risks, r.assessVendorLockInRisks(solution)...)

	// Assess operational risks
	risks = append(risks, r.assessOperationalRisks(solution)...)

	return risks, nil
}

func (r *RiskAssessorService) assessSinglePointsOfFailure(solution *interfaces.ProposedSolution) []*interfaces.TechnicalRisk {
	var risks []*interfaces.TechnicalRisk

	if solution.Architecture == nil {
		return risks
	}

	// Check for single points of failure in architecture
	if !solution.Architecture.HighAvailability {
		risk := &interfaces.TechnicalRisk{
			Risk: interfaces.Risk{
				ID:          uuid.New().String(),
				Category:    "technical",
				Title:       "Single Point of Failure - No High Availability",
				Description: "The proposed architecture lacks high availability configuration, creating single points of failure",
				Impact:      "high",
				Probability: "medium",
				RiskScore:   r.calculateRiskScore("high", "medium"),
				AffectedComponents: []string{"entire_system"},
			},
			ServiceType:        "infrastructure",
			ArchitecturalLayer: "infrastructure",
			PerformanceImpact:  "high",
			ScalabilityImpact:  "medium",
		}
		risks = append(risks, risk)
	}

	// Check for critical path services without redundancy
	for _, service := range solution.Services {
		if service.CriticalPath && len(service.Dependencies) == 0 {
			risk := &interfaces.TechnicalRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "technical",
					Title:       fmt.Sprintf("Critical Service Without Redundancy - %s", service.ServiceName),
					Description: fmt.Sprintf("Critical service %s lacks redundancy and backup mechanisms", service.ServiceName),
					Impact:      "high",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("high", "medium"),
					AffectedComponents: []string{service.ServiceName},
				},
				ServiceType:        service.ServiceType,
				ArchitecturalLayer: "application",
				PerformanceImpact:  "high",
				ScalabilityImpact:  "medium",
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessScalabilityRisks(solution *interfaces.ProposedSolution) []*interfaces.TechnicalRisk {
	var risks []*interfaces.TechnicalRisk

	// Check for potential scalability bottlenecks
	for _, service := range solution.Services {
		if service.ServiceType == "database" && service.Provider != "managed" {
			risk := &interfaces.TechnicalRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "technical",
					Title:       fmt.Sprintf("Database Scalability Risk - %s", service.ServiceName),
					Description: "Self-managed database may face scalability challenges under high load",
					Impact:      "medium",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("medium", "medium"),
					AffectedComponents: []string{service.ServiceName},
				},
				ServiceType:        service.ServiceType,
				ArchitecturalLayer: "data",
				PerformanceImpact:  "high",
				ScalabilityImpact:  "high",
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessPerformanceRisks(solution *interfaces.ProposedSolution) []*interfaces.TechnicalRisk {
	var risks []*interfaces.TechnicalRisk

	// Check for performance risks in data flow
	for _, dataFlow := range solution.DataFlow {
		if dataFlow.Volume == "high" && !dataFlow.Encryption {
			risk := &interfaces.TechnicalRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "technical",
					Title:       "Performance Risk - High Volume Unencrypted Data",
					Description: fmt.Sprintf("High volume data flow from %s to %s without encryption may impact performance", dataFlow.Source, dataFlow.Destination),
					Impact:      "medium",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("medium", "medium"),
					AffectedComponents: []string{dataFlow.Source, dataFlow.Destination},
				},
				ServiceType:        "data_transfer",
				ArchitecturalLayer: "network",
				PerformanceImpact:  "medium",
				ScalabilityImpact:  "low",
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessDependencyRisks(solution *interfaces.ProposedSolution) []*interfaces.TechnicalRisk {
	var risks []*interfaces.TechnicalRisk

	// Check for complex dependency chains
	dependencyMap := make(map[string][]string)
	for _, service := range solution.Services {
		dependencyMap[service.ServiceName] = service.Dependencies
	}

	// Find services with many dependencies
	for serviceName, deps := range dependencyMap {
		if len(deps) > 3 {
			risk := &interfaces.TechnicalRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "technical",
					Title:       fmt.Sprintf("Complex Dependency Chain - %s", serviceName),
					Description: fmt.Sprintf("Service %s has %d dependencies, creating complexity and potential failure points", serviceName, len(deps)),
					Impact:      "medium",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("medium", "medium"),
					AffectedComponents: append([]string{serviceName}, deps...),
				},
				ServiceType:        "integration",
				ArchitecturalLayer: "application",
				Dependencies:       deps,
				PerformanceImpact:  "medium",
				ScalabilityImpact:  "medium",
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessMigrationRisks(solution *interfaces.ProposedSolution) []*interfaces.TechnicalRisk {
	var risks []*interfaces.TechnicalRisk

	// Common migration risks
	risk := &interfaces.TechnicalRisk{
		Risk: interfaces.Risk{
			ID:          uuid.New().String(),
			Category:    "technical",
			Title:       "Data Migration Risk",
			Description: "Risk of data loss, corruption, or downtime during migration process",
			Impact:      "high",
			Probability: "medium",
			RiskScore:   r.calculateRiskScore("high", "medium"),
			AffectedComponents: []string{"data_migration"},
		},
		ServiceType:        "migration",
		ArchitecturalLayer: "data",
		PerformanceImpact:  "high",
		ScalabilityImpact:  "low",
	}
	risks = append(risks, risk)

	return risks
}

func (r *RiskAssessorService) assessDataEncryptionRisks(architecture *interfaces.Architecture) []*interfaces.SecurityRisk {
	var risks []*interfaces.SecurityRisk

	// Check data storage encryption
	for _, storage := range architecture.DataStorage {
		if storage.SensitivityLevel == "high" || storage.SensitivityLevel == "critical" {
			// Assume encryption is not configured if not explicitly mentioned
			risk := &interfaces.SecurityRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "security",
					Title:       fmt.Sprintf("Data Encryption Risk - %s", storage.ServiceName),
					Description: fmt.Sprintf("Sensitive data in %s may not be properly encrypted at rest and in transit", storage.ServiceName),
					Impact:      "high",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("high", "medium"),
					AffectedComponents: []string{storage.ServiceName},
				},
				ThreatType:         "data_exposure",
				AttackVectors:      []string{"data_breach", "unauthorized_access"},
				DataClassification: storage.SensitivityLevel,
				EncryptionRequired: true,
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessNetworkSecurityRisks(architecture *interfaces.Architecture) []*interfaces.SecurityRisk {
	var risks []*interfaces.SecurityRisk

	// Check for overly permissive security groups
	for _, sg := range architecture.NetworkTopology.SecurityGroups {
		for _, rule := range sg.Rules {
			if rule.Source == "0.0.0.0/0" && rule.Action == "allow" {
				risk := &interfaces.SecurityRisk{
					Risk: interfaces.Risk{
						ID:          uuid.New().String(),
						Category:    "security",
						Title:       fmt.Sprintf("Overly Permissive Security Group - %s", sg.Name),
						Description: fmt.Sprintf("Security group %s allows traffic from anywhere (0.0.0.0/0)", sg.Name),
						Impact:      "high",
						Probability: "high",
						RiskScore:   r.calculateRiskScore("high", "high"),
						AffectedComponents: []string{sg.Name},
					},
					ThreatType:    "network_intrusion",
					AttackVectors: []string{"port_scanning", "brute_force", "ddos"},
				}
				risks = append(risks, risk)
			}
		}
	}

	return risks
}

func (r *RiskAssessorService) assessAccessControlRisks(architecture *interfaces.Architecture) []*interfaces.SecurityRisk {
	var risks []*interfaces.SecurityRisk

	// Check for missing access controls
	hasAccessControls := false
	for _, layer := range architecture.SecurityLayers {
		for _, control := range layer.Controls {
			if control.Type == "access_control" || control.Type == "identity_management" {
				hasAccessControls = true
				break
			}
		}
	}

	if !hasAccessControls {
		risk := &interfaces.SecurityRisk{
			Risk: interfaces.Risk{
				ID:          uuid.New().String(),
				Category:    "security",
				Title:       "Missing Access Control",
				Description: "No explicit access control or identity management mechanisms identified in the architecture",
				Impact:      "high",
				Probability: "high",
				RiskScore:   r.calculateRiskScore("high", "high"),
				AffectedComponents: []string{"entire_system"},
			},
			ThreatType:    "unauthorized_access",
			AttackVectors: []string{"privilege_escalation", "insider_threat"},
		}
		risks = append(risks, risk)
	}

	return risks
}

func (r *RiskAssessorService) assessDataStorageSecurityRisks(architecture *interfaces.Architecture) []*interfaces.SecurityRisk {
	var risks []*interfaces.SecurityRisk

	for _, storage := range architecture.DataStorage {
		if storage.BackupStrategy == "" || storage.BackupStrategy == "none" {
			risk := &interfaces.SecurityRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "security",
					Title:       fmt.Sprintf("Missing Backup Strategy - %s", storage.ServiceName),
					Description: fmt.Sprintf("Data storage %s lacks a proper backup and recovery strategy", storage.ServiceName),
					Impact:      "high",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("high", "medium"),
					AffectedComponents: []string{storage.ServiceName},
				},
				ThreatType:         "data_loss",
				AttackVectors:      []string{"ransomware", "hardware_failure", "human_error"},
				DataClassification: storage.SensitivityLevel,
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessCostRisks(solution *interfaces.ProposedSolution) []*interfaces.BusinessRisk {
	var risks []*interfaces.BusinessRisk

	// Check for potential cost overruns
	if solution.EstimatedCost == "" || solution.EstimatedCost == "unknown" {
		risk := &interfaces.BusinessRisk{
			Risk: interfaces.Risk{
				ID:          uuid.New().String(),
				Category:    "business",
				Title:       "Cost Estimation Risk",
				Description: "Lack of detailed cost estimation may lead to budget overruns",
				Impact:      "medium",
				Probability: "high",
				RiskScore:   r.calculateRiskScore("medium", "high"),
				AffectedComponents: []string{"budget"},
			},
			BusinessFunction: "finance",
			RevenueImpact:    "medium",
			CustomerImpact:   "low",
			OperationalImpact: "medium",
		}
		risks = append(risks, risk)
	}

	return risks
}

func (r *RiskAssessorService) assessTimelineRisks(solution *interfaces.ProposedSolution) []*interfaces.BusinessRisk {
	var risks []*interfaces.BusinessRisk

	// Check for unrealistic timelines
	if solution.Timeline == "" || solution.Timeline == "unknown" {
		risk := &interfaces.BusinessRisk{
			Risk: interfaces.Risk{
				ID:          uuid.New().String(),
				Category:    "business",
				Title:       "Timeline Estimation Risk",
				Description: "Lack of detailed timeline estimation may lead to project delays",
				Impact:      "medium",
				Probability: "high",
				RiskScore:   r.calculateRiskScore("medium", "high"),
				AffectedComponents: []string{"project_timeline"},
			},
			BusinessFunction:  "project_management",
			RevenueImpact:     "medium",
			CustomerImpact:    "medium",
			OperationalImpact: "high",
		}
		risks = append(risks, risk)
	}

	return risks
}

func (r *RiskAssessorService) assessVendorLockInRisks(solution *interfaces.ProposedSolution) []*interfaces.BusinessRisk {
	var risks []*interfaces.BusinessRisk

	// Check for single cloud provider dependency
	if len(solution.CloudProviders) == 1 {
		risk := &interfaces.BusinessRisk{
			Risk: interfaces.Risk{
				ID:          uuid.New().String(),
				Category:    "business",
				Title:       "Vendor Lock-in Risk",
				Description: fmt.Sprintf("Solution depends entirely on %s, creating vendor lock-in risk", solution.CloudProviders[0]),
				Impact:      "medium",
				Probability: "medium",
				RiskScore:   r.calculateRiskScore("medium", "medium"),
				AffectedComponents: []string{"cloud_provider"},
			},
			BusinessFunction:  "strategy",
			RevenueImpact:     "medium",
			CustomerImpact:    "low",
			OperationalImpact: "medium",
		}
		risks = append(risks, risk)
	}

	return risks
}

func (r *RiskAssessorService) assessOperationalRisks(solution *interfaces.ProposedSolution) []*interfaces.BusinessRisk {
	var risks []*interfaces.BusinessRisk

	// Check for operational complexity
	if len(solution.Services) > 10 {
		risk := &interfaces.BusinessRisk{
			Risk: interfaces.Risk{
				ID:          uuid.New().String(),
				Category:    "business",
				Title:       "Operational Complexity Risk",
				Description: fmt.Sprintf("Solution involves %d services, creating operational complexity", len(solution.Services)),
				Impact:      "medium",
				Probability: "medium",
				RiskScore:   r.calculateRiskScore("medium", "medium"),
				AffectedComponents: []string{"operations"},
			},
			BusinessFunction:  "operations",
			RevenueImpact:     "low",
			CustomerImpact:    "medium",
			OperationalImpact: "high",
			TimeToRecover:     "hours",
		}
		risks = append(risks, risk)
	}

	return risks
}

func (r *RiskAssessorService) getComplianceFrameworksForIndustry(industry string) []string {
	industryLower := strings.ToLower(industry)
	
	switch {
	case strings.Contains(industryLower, "healthcare") || strings.Contains(industryLower, "medical"):
		return []string{"HIPAA", "HITECH"}
	case strings.Contains(industryLower, "financial") || strings.Contains(industryLower, "bank"):
		return []string{"PCI-DSS", "SOX", "GLBA"}
	case strings.Contains(industryLower, "retail") || strings.Contains(industryLower, "ecommerce"):
		return []string{"PCI-DSS", "GDPR"}
	case strings.Contains(industryLower, "government"):
		return []string{"FedRAMP", "FISMA"}
	default:
		return []string{"GDPR", "SOC2"}
	}
}

func (r *RiskAssessorService) assessComplianceFrameworkRisks(framework string, solution *interfaces.ProposedSolution) []*interfaces.ComplianceRisk {
	var risks []*interfaces.ComplianceRisk

	switch framework {
	case "HIPAA":
		risks = append(risks, r.assessHIPAARisks(solution)...)
	case "PCI-DSS":
		risks = append(risks, r.assessPCIDSSRisks(solution)...)
	case "GDPR":
		risks = append(risks, r.assessGDPRRisks(solution)...)
	case "SOX":
		risks = append(risks, r.assessSOXRisks(solution)...)
	}

	return risks
}

func (r *RiskAssessorService) assessHIPAARisks(solution *interfaces.ProposedSolution) []*interfaces.ComplianceRisk {
	var risks []*interfaces.ComplianceRisk

	// Check for PHI data handling
	for _, storage := range solution.Architecture.DataStorage {
		if storage.DataType == "personal" || storage.DataType == "health" {
			risk := &interfaces.ComplianceRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "compliance",
					Title:       "HIPAA PHI Protection Risk",
					Description: "Personal health information may not be adequately protected according to HIPAA requirements",
					Impact:      "critical",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("critical", "medium"),
					AffectedComponents: []string{storage.ServiceName},
				},
				Framework:     "HIPAA",
				RequirementID: "164.312",
				Jurisdiction:  "US",
				PenaltyLevel:  "high",
				AuditRequirements: []string{"access_logs", "encryption_verification", "backup_procedures"},
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessPCIDSSRisks(solution *interfaces.ProposedSolution) []*interfaces.ComplianceRisk {
	var risks []*interfaces.ComplianceRisk

	// Check for payment data handling
	for _, dataFlow := range solution.DataFlow {
		if strings.Contains(strings.ToLower(dataFlow.DataType), "payment") || 
		   strings.Contains(strings.ToLower(dataFlow.DataType), "card") {
			if !dataFlow.Encryption {
				risk := &interfaces.ComplianceRisk{
					Risk: interfaces.Risk{
						ID:          uuid.New().String(),
						Category:    "compliance",
						Title:       "PCI-DSS Data Encryption Risk",
						Description: "Payment card data transmission lacks proper encryption",
						Impact:      "critical",
						Probability: "high",
						RiskScore:   r.calculateRiskScore("critical", "high"),
						AffectedComponents: []string{dataFlow.Source, dataFlow.Destination},
					},
					Framework:     "PCI-DSS",
					RequirementID: "4.1",
					Jurisdiction:  "Global",
					PenaltyLevel:  "critical",
					AuditRequirements: []string{"encryption_verification", "key_management", "network_security"},
				}
				risks = append(risks, risk)
			}
		}
	}

	return risks
}

func (r *RiskAssessorService) assessGDPRRisks(solution *interfaces.ProposedSolution) []*interfaces.ComplianceRisk {
	var risks []*interfaces.ComplianceRisk

	// Check for personal data handling
	for _, storage := range solution.Architecture.DataStorage {
		if storage.DataType == "personal" {
			risk := &interfaces.ComplianceRisk{
				Risk: interfaces.Risk{
					ID:          uuid.New().String(),
					Category:    "compliance",
					Title:       "GDPR Personal Data Protection Risk",
					Description: "Personal data processing may not comply with GDPR requirements",
					Impact:      "high",
					Probability: "medium",
					RiskScore:   r.calculateRiskScore("high", "medium"),
					AffectedComponents: []string{storage.ServiceName},
				},
				Framework:     "GDPR",
				RequirementID: "Article 32",
				Jurisdiction:  "EU",
				PenaltyLevel:  "high",
				AuditRequirements: []string{"consent_management", "data_portability", "right_to_erasure"},
			}
			risks = append(risks, risk)
		}
	}

	return risks
}

func (r *RiskAssessorService) assessSOXRisks(solution *interfaces.ProposedSolution) []*interfaces.ComplianceRisk {
	var risks []*interfaces.ComplianceRisk

	// Check for financial data controls
	risk := &interfaces.ComplianceRisk{
		Risk: interfaces.Risk{
			ID:          uuid.New().String(),
			Category:    "compliance",
			Title:       "SOX Internal Controls Risk",
			Description: "Financial reporting systems may lack adequate internal controls",
			Impact:      "high",
			Probability: "medium",
			RiskScore:   r.calculateRiskScore("high", "medium"),
			AffectedComponents: []string{"financial_systems"},
		},
		Framework:     "SOX",
		RequirementID: "Section 404",
		Jurisdiction:  "US",
		PenaltyLevel:  "high",
		AuditRequirements: []string{"access_controls", "change_management", "audit_trails"},
	}
	risks = append(risks, risk)

	return risks
}

func (r *RiskAssessorService) generateMitigationStrategy(risk *interfaces.Risk) *interfaces.MitigationStrategy {
	strategy := &interfaces.MitigationStrategy{
		ID:     uuid.New().String(),
		RiskID: risk.ID,
	}

	// Generate mitigation strategies based on risk category and type
	switch risk.Category {
	case "technical":
		r.generateTechnicalMitigationStrategy(strategy, risk)
	case "security":
		r.generateSecurityMitigationStrategy(strategy, risk)
	case "compliance":
		r.generateComplianceMitigationStrategy(strategy, risk)
	case "business":
		r.generateBusinessMitigationStrategy(strategy, risk)
	}

	// Set priority based on risk score
	if risk.RiskScore >= 12 {
		strategy.Priority = "critical"
	} else if risk.RiskScore >= 8 {
		strategy.Priority = "high"
	} else if risk.RiskScore >= 4 {
		strategy.Priority = "medium"
	} else {
		strategy.Priority = "low"
	}

	return strategy
}

func (r *RiskAssessorService) generateTechnicalMitigationStrategy(strategy *interfaces.MitigationStrategy, risk *interfaces.Risk) {
	if strings.Contains(risk.Title, "Single Point of Failure") {
		strategy.Strategy = "Implement high availability and redundancy"
		strategy.ImplementationSteps = []string{
			"Configure multi-AZ deployment",
			"Implement load balancing",
			"Set up automated failover",
			"Create backup and recovery procedures",
		}
		strategy.EstimatedEffort = "2-4 weeks"
		strategy.Cost = "medium"
		strategy.Effectiveness = "high"
	} else if strings.Contains(risk.Title, "Scalability") {
		strategy.Strategy = "Implement auto-scaling and managed services"
		strategy.ImplementationSteps = []string{
			"Configure auto-scaling groups",
			"Migrate to managed database services",
			"Implement caching layers",
			"Optimize database queries",
		}
		strategy.EstimatedEffort = "3-6 weeks"
		strategy.Cost = "medium"
		strategy.Effectiveness = "high"
	} else if strings.Contains(risk.Title, "Performance") {
		strategy.Strategy = "Optimize performance and implement monitoring"
		strategy.ImplementationSteps = []string{
			"Implement performance monitoring",
			"Optimize data transfer protocols",
			"Add content delivery network",
			"Implement caching strategies",
		}
		strategy.EstimatedEffort = "2-3 weeks"
		strategy.Cost = "low"
		strategy.Effectiveness = "medium"
	}
}

func (r *RiskAssessorService) generateSecurityMitigationStrategy(strategy *interfaces.MitigationStrategy, risk *interfaces.Risk) {
	if strings.Contains(risk.Title, "Encryption") {
		strategy.Strategy = "Implement comprehensive encryption"
		strategy.ImplementationSteps = []string{
			"Enable encryption at rest for all data stores",
			"Implement TLS for data in transit",
			"Set up key management service",
			"Configure encryption for backups",
		}
		strategy.EstimatedEffort = "1-2 weeks"
		strategy.Cost = "low"
		strategy.Effectiveness = "high"
	} else if strings.Contains(risk.Title, "Security Group") {
		strategy.Strategy = "Implement principle of least privilege"
		strategy.ImplementationSteps = []string{
			"Review and restrict security group rules",
			"Implement network segmentation",
			"Set up VPN for administrative access",
			"Configure network access control lists",
		}
		strategy.EstimatedEffort = "1 week"
		strategy.Cost = "low"
		strategy.Effectiveness = "high"
	} else if strings.Contains(risk.Title, "Access Control") {
		strategy.Strategy = "Implement identity and access management"
		strategy.ImplementationSteps = []string{
			"Set up identity provider integration",
			"Implement role-based access control",
			"Configure multi-factor authentication",
			"Set up access logging and monitoring",
		}
		strategy.EstimatedEffort = "2-3 weeks"
		strategy.Cost = "medium"
		strategy.Effectiveness = "high"
	}
}

func (r *RiskAssessorService) generateComplianceMitigationStrategy(strategy *interfaces.MitigationStrategy, risk *interfaces.Risk) {
	strategy.Strategy = "Implement compliance controls and monitoring"
	strategy.ImplementationSteps = []string{
		"Conduct compliance gap analysis",
		"Implement required security controls",
		"Set up audit logging and monitoring",
		"Create compliance documentation",
		"Establish regular compliance reviews",
	}
	strategy.EstimatedEffort = "4-8 weeks"
	strategy.Cost = "high"
	strategy.Effectiveness = "high"
}

func (r *RiskAssessorService) generateBusinessMitigationStrategy(strategy *interfaces.MitigationStrategy, risk *interfaces.Risk) {
	if strings.Contains(risk.Title, "Cost") {
		strategy.Strategy = "Implement cost monitoring and optimization"
		strategy.ImplementationSteps = []string{
			"Set up detailed cost monitoring",
			"Implement budget alerts",
			"Create cost optimization plan",
			"Regular cost reviews and adjustments",
		}
		strategy.EstimatedEffort = "1-2 weeks"
		strategy.Cost = "low"
		strategy.Effectiveness = "medium"
	} else if strings.Contains(risk.Title, "Timeline") {
		strategy.Strategy = "Implement project management best practices"
		strategy.ImplementationSteps = []string{
			"Create detailed project plan",
			"Implement milestone tracking",
			"Set up regular progress reviews",
			"Identify and mitigate dependencies",
		}
		strategy.EstimatedEffort = "1 week"
		strategy.Cost = "low"
		strategy.Effectiveness = "medium"
	} else if strings.Contains(risk.Title, "Vendor Lock-in") {
		strategy.Strategy = "Implement multi-cloud strategy"
		strategy.ImplementationSteps = []string{
			"Design cloud-agnostic architecture",
			"Use containerization and orchestration",
			"Implement infrastructure as code",
			"Create migration procedures",
		}
		strategy.EstimatedEffort = "4-6 weeks"
		strategy.Cost = "high"
		strategy.Effectiveness = "medium"
	}
}

func (r *RiskAssessorService) collectAllRisks(assessment *interfaces.RiskAssessment) []*interfaces.Risk {
	var allRisks []*interfaces.Risk

	for _, risk := range assessment.TechnicalRisks {
		allRisks = append(allRisks, &risk.Risk)
	}
	for _, risk := range assessment.SecurityRisks {
		allRisks = append(allRisks, &risk.Risk)
	}
	for _, risk := range assessment.ComplianceRisks {
		allRisks = append(allRisks, &risk.Risk)
	}
	for _, risk := range assessment.BusinessRisks {
		allRisks = append(allRisks, &risk.Risk)
	}

	return allRisks
}

func (r *RiskAssessorService) calculateOverallRiskLevel(risks []*interfaces.Risk) string {
	if len(risks) == 0 {
		return "low"
	}

	totalScore := 0
	criticalCount := 0
	highCount := 0

	for _, risk := range risks {
		totalScore += risk.RiskScore
		if risk.Impact == "critical" {
			criticalCount++
		} else if risk.Impact == "high" {
			highCount++
		}
	}

	avgScore := float64(totalScore) / float64(len(risks))

	// Determine overall risk level
	if criticalCount > 0 || avgScore >= 12 {
		return "critical"
	} else if highCount > 2 || avgScore >= 8 {
		return "high"
	} else if avgScore >= 4 {
		return "medium"
	}

	return "low"
}

func (r *RiskAssessorService) generateRecommendedActions(assessment *interfaces.RiskAssessment) []string {
	var actions []string

	// Generate actions based on overall risk level
	switch assessment.OverallRiskLevel {
	case "critical":
		actions = append(actions, "Immediate action required - address critical risks before proceeding")
		actions = append(actions, "Conduct emergency risk review meeting")
		actions = append(actions, "Consider alternative solutions or additional safeguards")
	case "high":
		actions = append(actions, "Address high-priority risks before implementation")
		actions = append(actions, "Implement additional monitoring and controls")
		actions = append(actions, "Schedule regular risk review meetings")
	case "medium":
		actions = append(actions, "Monitor identified risks during implementation")
		actions = append(actions, "Implement recommended mitigation strategies")
		actions = append(actions, "Schedule periodic risk assessments")
	case "low":
		actions = append(actions, "Proceed with standard risk monitoring")
		actions = append(actions, "Implement basic mitigation strategies")
	}

	// Add specific actions based on risk types
	if len(assessment.SecurityRisks) > 0 {
		actions = append(actions, "Conduct security review with cybersecurity team")
	}
	if len(assessment.ComplianceRisks) > 0 {
		actions = append(actions, "Engage compliance team for regulatory review")
	}

	return actions
}

func (r *RiskAssessorService) calculateRiskScore(impact, probability string) int {
	impactScore := r.getImpactScore(impact)
	probabilityScore := r.getProbabilityScore(probability)
	return impactScore * probabilityScore
}

func (r *RiskAssessorService) getImpactScore(impact string) int {
	switch impact {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 1
	}
}

func (r *RiskAssessorService) getProbabilityScore(probability string) int {
	switch probability {
	case "high":
		return 4
	case "medium":
		return 3
	case "low":
		return 2
	default:
		return 1
	}
}

func (r *RiskAssessorService) extractIndustryFromInquiry(inquiry *domain.Inquiry) string {
	// Try to extract industry from company name or message
	message := strings.ToLower(inquiry.Message)
	company := strings.ToLower(inquiry.Company)
	
	industries := map[string][]string{
		"healthcare": {"health", "medical", "hospital", "clinic", "pharma"},
		"financial": {"bank", "finance", "financial", "credit", "insurance"},
		"retail": {"retail", "ecommerce", "store", "shop", "commerce"},
		"government": {"government", "gov", "public", "federal", "state"},
		"education": {"school", "university", "education", "academic"},
		"manufacturing": {"manufacturing", "factory", "production", "industrial"},
	}

	for industry, keywords := range industries {
		for _, keyword := range keywords {
			if strings.Contains(message, keyword) || strings.Contains(company, keyword) {
				return industry
			}
		}
	}

	return "general"
}

func (r *RiskAssessorService) isMigrationProject(inquiry *domain.Inquiry) bool {
	for _, service := range inquiry.Services {
		if service == "migration" {
			return true
		}
	}
	return strings.Contains(strings.ToLower(inquiry.Message), "migrat")
}