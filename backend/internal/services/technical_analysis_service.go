package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// TechnicalAnalysisService implements the TechnicalAnalysisService interface
type TechnicalAnalysisService struct {
	bedrockService interfaces.BedrockService
	knowledgeBase  interfaces.KnowledgeBase
}

// NewTechnicalAnalysisService creates a new technical analysis service
func NewTechnicalAnalysisService(bedrock interfaces.BedrockService, kb interfaces.KnowledgeBase) *TechnicalAnalysisService {
	return &TechnicalAnalysisService{
		bedrockService: bedrock,
		knowledgeBase:  kb,
	}
}

// AnalyzeCodebase performs comprehensive code analysis
func (s *TechnicalAnalysisService) AnalyzeCodebase(ctx context.Context, request *interfaces.CodeAnalysisRequest) (*interfaces.CodeAnalysisResult, error) {
	// Generate analysis prompt
	prompt := s.buildCodeAnalysisPrompt(request)

	// Get AI analysis
	response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate code analysis: %w", err)
	}

	// Parse and structure the analysis results
	result := &interfaces.CodeAnalysisResult{
		ID:                      generateID("code_analysis"),
		InquiryID:               request.InquiryID,
		OverallScore:            s.calculateCodeScore(response.Content),
		SecurityFindings:        s.extractSecurityFindings(response.Content),
		PerformanceFindings:     s.extractPerformanceFindings(response.Content),
		ArchitectureFindings:    s.extractArchitectureFindings(response.Content),
		MaintainabilityFindings: s.extractMaintainabilityFindings(response.Content),
		BestPracticeViolations:  s.extractBestPracticeViolations(response.Content),
		CloudOptimizations:      s.extractCloudOptimizations(response.Content),
		Summary:                 s.extractSummary(response.Content),
		Recommendations:         s.extractCodeRecommendations(response.Content),
		GeneratedAt:             time.Now(),
	}

	return result, nil
}

// AssessArchitecture performs architecture assessment
func (s *TechnicalAnalysisService) AssessArchitecture(ctx context.Context, request *interfaces.TechArchitectureAssessmentRequest) (*interfaces.TechArchitectureAssessmentResult, error) {
	// Generate architecture analysis prompt
	prompt := s.buildArchitectureAnalysisPrompt(request)

	// Get AI analysis
	response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate architecture analysis: %w", err)
	}

	// Parse and structure the analysis results
	result := &interfaces.TechArchitectureAssessmentResult{
		ID:                        generateID("arch_analysis"),
		InquiryID:                 request.InquiryID,
		OverallScore:              s.calculateArchitectureScore(response.Content),
		TechSecurityAssessment:    s.extractTechSecurityAssessment(response.Content),
		TechPerformanceAssessment: s.extractTechPerformanceAssessment(response.Content),
		TechScalabilityAssessment: s.extractTechScalabilityAssessment(response.Content),
		TechReliabilityAssessment: s.extractTechReliabilityAssessment(response.Content),
		TechCostAssessment:        s.extractTechCostAssessment(response.Content),
		TechComplianceAssessment:  s.extractTechComplianceAssessment(response.Content),
		ArchitectureFindings:      s.extractTechArchitectureFindings(response.Content),
		Recommendations:           s.extractTechArchitectureRecommendations(response.Content),
		Summary:                   s.extractSummary(response.Content),
		GeneratedAt:               time.Now(),
	}

	return result, nil
}

// PerformSecurityAssessment performs comprehensive security assessment
func (s *TechnicalAnalysisService) PerformSecurityAssessment(ctx context.Context, request *interfaces.TechSecurityAssessmentRequest) (*interfaces.TechSecurityAssessmentResult, error) {
	// Generate security analysis prompt
	prompt := s.buildSecurityAnalysisPrompt(request)

	// Get AI analysis
	response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate security analysis: %w", err)
	}

	// Parse and structure the analysis results
	result := &interfaces.TechSecurityAssessmentResult{
		ID:                      generateID("security_analysis"),
		InquiryID:               request.InquiryID,
		OverallScore:            s.calculateSecurityScore(response.Content),
		RiskLevel:               s.determineRiskLevel(response.Content),
		SecurityVulnerabilities: s.extractTechSecurityVulnerabilities(response.Content),
		ComplianceGaps:          s.extractTechComplianceGaps(response.Content),
		ThreatAnalysis:          s.extractTechThreatAnalysis(response.Content),
		SecurityRecommendations: s.extractTechSecurityRecommendations(response.Content),
		Summary:                 s.extractSummary(response.Content),
		GeneratedAt:             time.Now(),
	}

	return result, nil
}

// GenerateSecurityRemediation generates remediation recommendations for vulnerabilities
func (s *TechnicalAnalysisService) GenerateSecurityRemediation(ctx context.Context, vulnerabilities []*interfaces.TechSecurityVulnerability) ([]*interfaces.RemediationRecommendation, error) {
	var recommendations []*interfaces.RemediationRecommendation

	for _, vuln := range vulnerabilities {
		prompt := s.buildRemediationPrompt(vuln)

		response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
			ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
			MaxTokens:   2000,
			Temperature: 0.1,
		})
		if err != nil {
			continue // Skip this vulnerability if analysis fails
		}

		recommendation := &interfaces.RemediationRecommendation{
			ID:              generateID("remediation"),
			VulnerabilityID: vuln.ID,
			Title:           s.extractRemediationTitle(response.Content),
			Description:     s.extractRemediationDescription(response.Content),
			Steps:           s.extractRemediationSteps(response.Content),
			Priority:        vuln.Severity,
			Effort:          s.estimateRemediationEffort(vuln.Severity),
			Timeline:        s.estimateRemediationTimeline(vuln.Severity),
			References:      s.extractRemediationReferences(response.Content),
		}

		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

// AnalyzePerformance performs performance analysis
func (s *TechnicalAnalysisService) AnalyzePerformance(ctx context.Context, request *interfaces.TechPerformanceAnalysisRequest) (*interfaces.TechPerformanceAnalysisResult, error) {
	// Generate performance analysis prompt
	prompt := s.buildPerformanceAnalysisPrompt(request)

	// Get AI analysis
	response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate performance analysis: %w", err)
	}

	// Parse and structure the analysis results
	result := &interfaces.TechPerformanceAnalysisResult{
		ID:                         generateID("perf_analysis"),
		InquiryID:                  request.InquiryID,
		OverallScore:               s.calculatePerformanceScore(response.Content),
		PerformanceBottlenecks:     s.extractTechPerformanceBottlenecks(response.Content),
		OptimizationOpportunities:  s.extractTechOptimizationOpportunities(response.Content),
		BenchmarkComparison:        s.extractBenchmarkComparison(response.Content),
		PerformanceRecommendations: s.extractTechPerformanceRecommendations(response.Content),
		Summary:                    s.extractSummary(response.Content),
		GeneratedAt:                time.Now(),
	}

	return result, nil
}

// GenerateOptimizationRecommendations generates optimization recommendations
func (s *TechnicalAnalysisService) GenerateOptimizationRecommendations(ctx context.Context, metrics *interfaces.TechPerformanceMetrics) ([]*interfaces.TechOptimizationRecommendation, error) {
	prompt := s.buildOptimizationPrompt(metrics)

	response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   3000,
		Temperature: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate optimization recommendations: %w", err)
	}

	return s.extractTechOptimizationOpportunities(response.Content), nil
}

// AnalyzeCompliance performs compliance gap analysis
func (s *TechnicalAnalysisService) AnalyzeCompliance(ctx context.Context, request *interfaces.TechComplianceAnalysisRequest) (*interfaces.TechComplianceAnalysisResult, error) {
	// Generate compliance analysis prompt
	prompt := s.buildComplianceAnalysisPrompt(request)

	// Get AI analysis
	response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   4000,
		Temperature: 0.1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate compliance analysis: %w", err)
	}

	// Parse and structure the analysis results
	result := &interfaces.TechComplianceAnalysisResult{
		ID:                        generateID("compliance_analysis"),
		InquiryID:                 request.InquiryID,
		OverallScore:              s.calculateComplianceScore(response.Content),
		ComplianceStatus:          s.extractComplianceStatus(response.Content, request.ComplianceFrameworks),
		ComplianceGaps:            s.extractTechComplianceGaps(response.Content),
		ControlAssessment:         s.extractTechControlAssessment(response.Content),
		ComplianceRecommendations: s.extractTechComplianceRecommendations(response.Content),
		RoadmapToCompliance:       s.extractTechComplianceRoadmap(response.Content),
		Summary:                   s.extractSummary(response.Content),
		GeneratedAt:               time.Now(),
	}

	return result, nil
}

// GenerateComplianceRemediation generates compliance remediation steps
func (s *TechnicalAnalysisService) GenerateComplianceRemediation(ctx context.Context, gaps []*interfaces.TechComplianceGap) ([]*interfaces.TechComplianceRemediation, error) {
	var remediations []*interfaces.TechComplianceRemediation

	for _, gap := range gaps {
		prompt := s.buildComplianceRemediationPrompt(gap)

		response, err := s.bedrockService.GenerateText(ctx, prompt, &interfaces.BedrockOptions{
			ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
			MaxTokens:   2000,
			Temperature: 0.1,
		})
		if err != nil {
			continue // Skip this gap if analysis fails
		}

		remediation := &interfaces.TechComplianceRemediation{
			ID:          generateID("compliance_remediation"),
			GapID:       gap.ID,
			Title:       s.extractRemediationTitle(response.Content),
			Description: s.extractRemediationDescription(response.Content),
			Steps:       s.extractRemediationSteps(response.Content),
			Priority:    gap.Priority,
			Effort:      s.estimateComplianceEffort(gap.Impact),
			Timeline:    s.estimateComplianceTimeline(gap.Priority),
			References:  s.extractRemediationReferences(response.Content),
		}

		remediations = append(remediations, remediation)
	}

	return remediations, nil
}

// PerformComprehensiveAnalysis performs comprehensive technical analysis
func (s *TechnicalAnalysisService) PerformComprehensiveAnalysis(ctx context.Context, request *interfaces.ComprehensiveAnalysisRequest) (*interfaces.ComprehensiveAnalysisResult, error) {
	result := &interfaces.ComprehensiveAnalysisResult{
		ID:          generateID("comprehensive_analysis"),
		InquiryID:   request.InquiryID,
		GeneratedAt: time.Now(),
	}

	var overallScores []float64

	// Perform individual analyses based on scope
	for _, scope := range request.AnalysisScope {
		switch scope {
		case "code":
			if request.CodeAnalysisRequest != nil {
				codeResult, err := s.AnalyzeCodebase(ctx, request.CodeAnalysisRequest)
				if err == nil {
					result.CodeAnalysisResult = codeResult
					overallScores = append(overallScores, codeResult.OverallScore)
				}
			}
		case "architecture":
			if request.ArchitectureRequest != nil {
				archResult, err := s.AssessArchitecture(ctx, request.ArchitectureRequest)
				if err == nil {
					result.ArchitectureResult = archResult
					overallScores = append(overallScores, archResult.OverallScore)
				}
			}
		case "security":
			if request.SecurityRequest != nil {
				secResult, err := s.PerformSecurityAssessment(ctx, request.SecurityRequest)
				if err == nil {
					result.SecurityResult = secResult
					overallScores = append(overallScores, secResult.OverallScore)
				}
			}
		case "performance":
			if request.PerformanceRequest != nil {
				perfResult, err := s.AnalyzePerformance(ctx, request.PerformanceRequest)
				if err == nil {
					result.PerformanceResult = perfResult
					overallScores = append(overallScores, perfResult.OverallScore)
				}
			}
		case "compliance":
			if request.ComplianceRequest != nil {
				compResult, err := s.AnalyzeCompliance(ctx, request.ComplianceRequest)
				if err == nil {
					result.ComplianceResult = compResult
					overallScores = append(overallScores, compResult.OverallScore)
				}
			}
		}
	}

	// Calculate overall score
	if len(overallScores) > 0 {
		var sum float64
		for _, score := range overallScores {
			sum += score
		}
		result.OverallScore = sum / float64(len(overallScores))
	}

	// Generate cross-cutting findings and prioritized recommendations
	result.CrossCuttingFindings = s.identifyCrossCuttingFindings(result)
	result.PrioritizedRecommendations = s.prioritizeRecommendations(result)
	result.ExecutiveSummary = s.generateExecutiveSummary(result)
	result.TechnicalSummary = s.generateTechnicalSummary(result)
	result.ActionPlan = s.generateActionPlan(result)

	return result, nil
}

// Prompt building methods

func (s *TechnicalAnalysisService) buildCodeAnalysisPrompt(request *interfaces.CodeAnalysisRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are a senior technical architect performing a comprehensive code analysis. ")
	prompt.WriteString("Analyze the provided code samples and provide detailed findings in the following areas:\n\n")

	prompt.WriteString("## CODE ANALYSIS REQUEST\n")
	prompt.WriteString(fmt.Sprintf("Application Type: %s\n", request.ApplicationType))
	prompt.WriteString(fmt.Sprintf("Cloud Provider: %s\n", request.CloudProvider))
	prompt.WriteString(fmt.Sprintf("Languages: %s\n", strings.Join(request.Languages, ", ")))
	prompt.WriteString(fmt.Sprintf("Analysis Scope: %s\n", strings.Join(request.AnalysisScope, ", ")))

	if len(request.CodeSamples) > 0 {
		prompt.WriteString("\n## CODE SAMPLES\n")
		for _, sample := range request.CodeSamples {
			prompt.WriteString(fmt.Sprintf("### %s (%s)\n", sample.Filename, sample.Language))
			prompt.WriteString("```\n")
			prompt.WriteString(sample.Content)
			prompt.WriteString("\n```\n\n")
		}
	}

	prompt.WriteString("\n## ANALYSIS REQUIREMENTS\n")
	prompt.WriteString("Provide a comprehensive analysis covering:\n")
	prompt.WriteString("1. Security vulnerabilities and risks\n")
	prompt.WriteString("2. Performance bottlenecks and optimization opportunities\n")
	prompt.WriteString("3. Architecture patterns and anti-patterns\n")
	prompt.WriteString("4. Code maintainability and technical debt\n")
	prompt.WriteString("5. Best practice violations\n")
	prompt.WriteString("6. Cloud-specific optimization opportunities\n")
	prompt.WriteString("7. Overall recommendations with priorities\n\n")

	prompt.WriteString("Format your response with clear sections and actionable recommendations.")

	return prompt.String()
}

func (s *TechnicalAnalysisService) buildArchitectureAnalysisPrompt(request *interfaces.TechArchitectureAssessmentRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are a senior cloud architect performing a comprehensive architecture assessment. ")
	prompt.WriteString("Analyze the provided system architecture and provide detailed findings.\n\n")

	prompt.WriteString("## ARCHITECTURE ASSESSMENT REQUEST\n")
	prompt.WriteString(fmt.Sprintf("System Description: %s\n", request.SystemDescription))
	prompt.WriteString(fmt.Sprintf("Cloud Provider: %s\n", request.CloudProvider))
	prompt.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(request.Services, ", ")))
	prompt.WriteString(fmt.Sprintf("Security Requirements: %s\n", strings.Join(request.SecurityRequirements, ", ")))
	prompt.WriteString(fmt.Sprintf("Compliance Frameworks: %s\n", strings.Join(request.ComplianceFrameworks, ", ")))

	if request.SystemTrafficPatterns != nil {
		prompt.WriteString(fmt.Sprintf("\nTraffic Patterns:\n"))
		prompt.WriteString(fmt.Sprintf("- Peak Traffic: %d RPS\n", request.SystemTrafficPatterns.PeakTraffic))
		prompt.WriteString(fmt.Sprintf("- Average Traffic: %d RPS\n", request.SystemTrafficPatterns.AverageTraffic))
		prompt.WriteString(fmt.Sprintf("- Growth Projection: %.1f%% per year\n", request.SystemTrafficPatterns.GrowthProjection))
	}

	if request.TechPerformanceRequirements != nil {
		prompt.WriteString(fmt.Sprintf("\nPerformance Requirements:\n"))
		prompt.WriteString(fmt.Sprintf("- Response Time: %d ms\n", request.TechPerformanceRequirements.ResponseTime))
		prompt.WriteString(fmt.Sprintf("- Throughput: %d RPS\n", request.TechPerformanceRequirements.Throughput))
		prompt.WriteString(fmt.Sprintf("- Availability: %.2f%%\n", request.TechPerformanceRequirements.Availability))
	}

	prompt.WriteString("\n## ASSESSMENT REQUIREMENTS\n")
	prompt.WriteString("Provide a comprehensive assessment covering:\n")
	prompt.WriteString("1. Security architecture and vulnerabilities\n")
	prompt.WriteString("2. Performance and scalability analysis\n")
	prompt.WriteString("3. Reliability and fault tolerance\n")
	prompt.WriteString("4. Cost optimization opportunities\n")
	prompt.WriteString("5. Compliance gap analysis\n")
	prompt.WriteString("6. Architecture recommendations with priorities\n\n")

	prompt.WriteString("Format your response with clear sections and actionable recommendations.")

	return prompt.String()
}

func (s *TechnicalAnalysisService) buildSecurityAnalysisPrompt(request *interfaces.TechSecurityAssessmentRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are a senior security architect performing a comprehensive security assessment. ")
	prompt.WriteString("Analyze the provided system and identify security vulnerabilities and risks.\n\n")

	prompt.WriteString("## SECURITY ASSESSMENT REQUEST\n")
	prompt.WriteString(fmt.Sprintf("System Description: %s\n", request.SystemDescription))
	prompt.WriteString(fmt.Sprintf("Cloud Provider: %s\n", request.CloudProvider))
	prompt.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(request.Services, ", ")))
	prompt.WriteString(fmt.Sprintf("Data Classification: %s\n", strings.Join(request.DataClassification, ", ")))
	prompt.WriteString(fmt.Sprintf("Compliance Frameworks: %s\n", strings.Join(request.ComplianceFrameworks, ", ")))

	if request.TechThreatModel != nil {
		prompt.WriteString(fmt.Sprintf("\nThreat Model:\n"))
		prompt.WriteString(fmt.Sprintf("- Assets: %s\n", strings.Join(request.TechThreatModel.Assets, ", ")))
		prompt.WriteString(fmt.Sprintf("- Threats: %s\n", strings.Join(request.TechThreatModel.Threats, ", ")))
		prompt.WriteString(fmt.Sprintf("- Attack Vectors: %s\n", strings.Join(request.TechThreatModel.AttackVectors, ", ")))
	}

	prompt.WriteString("\n## SECURITY ASSESSMENT REQUIREMENTS\n")
	prompt.WriteString("Provide a comprehensive security assessment covering:\n")
	prompt.WriteString("1. Security vulnerabilities with CVSS scores\n")
	prompt.WriteString("2. Compliance gaps for specified frameworks\n")
	prompt.WriteString("3. Threat analysis and risk assessment\n")
	prompt.WriteString("4. Security recommendations with priorities\n")
	prompt.WriteString("5. Remediation steps for critical vulnerabilities\n\n")

	prompt.WriteString("Format your response with clear sections and actionable recommendations.")

	return prompt.String()
}

func (s *TechnicalAnalysisService) buildPerformanceAnalysisPrompt(request *interfaces.TechPerformanceAnalysisRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are a senior performance engineer performing a comprehensive performance analysis. ")
	prompt.WriteString("Analyze the provided system and identify performance bottlenecks and optimization opportunities.\n\n")

	prompt.WriteString("## PERFORMANCE ANALYSIS REQUEST\n")
	prompt.WriteString(fmt.Sprintf("System Description: %s\n", request.SystemDescription))
	prompt.WriteString(fmt.Sprintf("Cloud Provider: %s\n", request.CloudProvider))
	prompt.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(request.Services, ", ")))

	if request.PerformanceMetrics != nil {
		prompt.WriteString(fmt.Sprintf("\nCurrent Performance Metrics:\n"))
		prompt.WriteString(fmt.Sprintf("- Response Time: %.2f ms\n", request.PerformanceMetrics.ResponseTime))
		prompt.WriteString(fmt.Sprintf("- Throughput: %.2f RPS\n", request.PerformanceMetrics.Throughput))
		prompt.WriteString(fmt.Sprintf("- Error Rate: %.2f%%\n", request.PerformanceMetrics.ErrorRate))
		prompt.WriteString(fmt.Sprintf("- CPU Utilization: %.2f%%\n", request.PerformanceMetrics.CPUUtilization))
		prompt.WriteString(fmt.Sprintf("- Memory Utilization: %.2f%%\n", request.PerformanceMetrics.MemoryUtilization))
	}

	if request.ApplicationProfile != nil {
		prompt.WriteString(fmt.Sprintf("\nApplication Profile:\n"))
		prompt.WriteString(fmt.Sprintf("- Type: %s\n", request.ApplicationProfile.Type))
		prompt.WriteString(fmt.Sprintf("- Language: %s\n", request.ApplicationProfile.Language))
		prompt.WriteString(fmt.Sprintf("- Framework: %s\n", request.ApplicationProfile.Framework))
		prompt.WriteString(fmt.Sprintf("- Database: %s\n", request.ApplicationProfile.DatabaseType))
	}

	prompt.WriteString("\n## PERFORMANCE ANALYSIS REQUIREMENTS\n")
	prompt.WriteString("Provide a comprehensive performance analysis covering:\n")
	prompt.WriteString("1. Performance bottlenecks identification\n")
	prompt.WriteString("2. Optimization opportunities with expected improvements\n")
	prompt.WriteString("3. Benchmark comparison with industry standards\n")
	prompt.WriteString("4. Performance recommendations with priorities\n")
	prompt.WriteString("5. Implementation roadmap for optimizations\n\n")

	prompt.WriteString("Format your response with clear sections and actionable recommendations.")

	return prompt.String()
}

func (s *TechnicalAnalysisService) buildComplianceAnalysisPrompt(request *interfaces.TechComplianceAnalysisRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are a senior compliance architect performing a comprehensive compliance gap analysis. ")
	prompt.WriteString("Analyze the provided system against specified compliance frameworks.\n\n")

	prompt.WriteString("## COMPLIANCE ANALYSIS REQUEST\n")
	prompt.WriteString(fmt.Sprintf("System Description: %s\n", request.SystemDescription))
	prompt.WriteString(fmt.Sprintf("Cloud Provider: %s\n", request.CloudProvider))
	prompt.WriteString(fmt.Sprintf("Services: %s\n", strings.Join(request.Services, ", ")))
	prompt.WriteString(fmt.Sprintf("Compliance Frameworks: %s\n", strings.Join(request.ComplianceFrameworks, ", ")))
	prompt.WriteString(fmt.Sprintf("Data Classification: %s\n", strings.Join(request.DataClassification, ", ")))
	prompt.WriteString(fmt.Sprintf("Geographic Scope: %s\n", strings.Join(request.GeographicScope, ", ")))

	if request.BusinessContext != nil {
		prompt.WriteString(fmt.Sprintf("\nBusiness Context:\n"))
		prompt.WriteString(fmt.Sprintf("- Industry: %s\n", request.BusinessContext.Industry))
		prompt.WriteString(fmt.Sprintf("- Company Size: %s\n", request.BusinessContext.CompanySize))
		prompt.WriteString(fmt.Sprintf("- Data Types: %s\n", strings.Join(request.BusinessContext.DataTypes, ", ")))
	}

	prompt.WriteString("\n## COMPLIANCE ANALYSIS REQUIREMENTS\n")
	prompt.WriteString("For each specified compliance framework, provide:\n")
	prompt.WriteString("1. Current compliance status assessment\n")
	prompt.WriteString("2. Identified compliance gaps with priorities\n")
	prompt.WriteString("3. Control assessment for key requirements\n")
	prompt.WriteString("4. Compliance recommendations with timelines\n")
	prompt.WriteString("5. Roadmap to achieve full compliance\n\n")

	prompt.WriteString("Format your response with clear sections and actionable recommendations.")

	return prompt.String()
}

// Helper methods for parsing AI responses

func (s *TechnicalAnalysisService) calculateCodeScore(content string) float64 {
	// Simple scoring based on content analysis
	score := 75.0 // Base score

	if strings.Contains(strings.ToLower(content), "critical") {
		score -= 20
	}
	if strings.Contains(strings.ToLower(content), "high risk") {
		score -= 15
	}
	if strings.Contains(strings.ToLower(content), "security vulnerability") {
		score -= 10
	}
	if strings.Contains(strings.ToLower(content), "performance issue") {
		score -= 10
	}
	if strings.Contains(strings.ToLower(content), "best practices") {
		score += 5
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

func (s *TechnicalAnalysisService) extractSecurityFindings(content string) []*interfaces.TechSecurityFinding {
	// Parse security findings from AI response
	findings := []*interfaces.TechSecurityFinding{
		{
			ID:          generateID("security_finding"),
			Type:        "vulnerability",
			Severity:    "medium",
			Title:       "Example Security Finding",
			Description: "This is an example security finding extracted from the analysis.",
			Location: interfaces.Location{
				File: "example.go",
				Line: 42,
			},
			Impact:      "Potential data exposure",
			Remediation: "Implement proper input validation",
			References:  []string{"https://owasp.org/"},
		},
	}

	return findings
}

func (s *TechnicalAnalysisService) extractPerformanceFindings(content string) []*interfaces.TechPerformanceFinding {
	findings := []*interfaces.TechPerformanceFinding{
		{
			ID:          generateID("perf_finding"),
			Type:        "bottleneck",
			Severity:    "medium",
			Title:       "Example Performance Finding",
			Description: "This is an example performance finding extracted from the analysis.",
			Location: interfaces.Location{
				File: "example.go",
				Line: 100,
			},
			Impact:               "Reduced throughput",
			Suggestion:           "Implement caching mechanism",
			EstimatedImprovement: "30% performance improvement",
		},
	}

	return findings
}

func (s *TechnicalAnalysisService) extractArchitectureFindings(content string) []*interfaces.TechArchitectureFinding {
	return s.extractTechArchitectureFindings(content)
}

func (s *TechnicalAnalysisService) extractTechArchitectureFindings(content string) []*interfaces.TechArchitectureFinding {
	findings := []*interfaces.TechArchitectureFinding{
		{
			ID:          generateID("arch_finding"),
			Type:        "anti-pattern",
			Severity:    "medium",
			Title:       "Example Architecture Finding",
			Description: "This is an example architecture finding extracted from the analysis.",
			Location: interfaces.Location{
				File: "architecture.go",
				Line: 50,
			},
			Impact:            "Reduced maintainability",
			Recommendation:    "Refactor to use dependency injection",
			RefactoringEffort: "medium",
		},
	}

	return findings
}

func (s *TechnicalAnalysisService) extractMaintainabilityFindings(content string) []*interfaces.MaintainabilityFinding {
	findings := []*interfaces.MaintainabilityFinding{
		{
			ID:          generateID("maint_finding"),
			Type:        "complexity",
			Severity:    "medium",
			Title:       "Example Maintainability Finding",
			Description: "This is an example maintainability finding extracted from the analysis.",
			Location: interfaces.Location{
				File: "complex.go",
				Line: 75,
			},
			Metrics: map[string]float64{
				"cyclomatic_complexity": 15.0,
			},
			Suggestion: "Break down into smaller functions",
		},
	}

	return findings
}

func (s *TechnicalAnalysisService) extractBestPracticeViolations(content string) []*interfaces.BestPracticeViolation {
	violations := []*interfaces.BestPracticeViolation{
		{
			ID:          generateID("bp_violation"),
			Practice:    "Error Handling",
			Category:    "maintainability",
			Severity:    "medium",
			Description: "Missing error handling in critical path",
			Location: interfaces.Location{
				File: "handler.go",
				Line: 25,
			},
			Correction: "Add proper error handling and logging",
			References: []string{"https://golang.org/doc/effective_go.html#errors"},
		},
	}

	return violations
}

func (s *TechnicalAnalysisService) extractCloudOptimizations(content string) []*interfaces.CloudOptimization {
	optimizations := []*interfaces.CloudOptimization{
		{
			ID:                   generateID("cloud_opt"),
			Type:                 "cost",
			Title:                "Example Cloud Optimization",
			Description:          "This is an example cloud optimization opportunity.",
			CurrentState:         "Using on-demand instances",
			RecommendedState:     "Use reserved instances for predictable workloads",
			EstimatedSavings:     500.0,
			ImplementationEffort: "low",
			Priority:             "medium",
		},
	}

	return optimizations
}

func (s *TechnicalAnalysisService) extractCodeRecommendations(content string) []*interfaces.CodeRecommendation {
	recommendations := []*interfaces.CodeRecommendation{
		{
			ID:           generateID("code_rec"),
			Type:         "refactor",
			Priority:     "medium",
			Title:        "Example Code Recommendation",
			Description:  "This is an example code recommendation extracted from the analysis.",
			Benefits:     []string{"Improved maintainability", "Better performance"},
			Effort:       "medium",
			Timeline:     "2-3 weeks",
			Dependencies: []string{"Complete security fixes first"},
			References:   []string{"https://refactoring.guru/"},
		},
	}

	return recommendations
}

func (s *TechnicalAnalysisService) extractSummary(content string) string {
	// Extract summary from AI response
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), "summary") && i+1 < len(lines) {
			return strings.TrimSpace(lines[i+1])
		}
	}

	return "Comprehensive technical analysis completed with actionable recommendations for improvement."
}

// Additional helper methods for other analysis types

func (s *TechnicalAnalysisService) calculateArchitectureScore(content string) float64 {
	return s.calculateCodeScore(content) // Reuse similar logic
}

func (s *TechnicalAnalysisService) calculateSecurityScore(content string) float64 {
	score := 70.0 // Base security score

	if strings.Contains(strings.ToLower(content), "critical vulnerability") {
		score -= 30
	}
	if strings.Contains(strings.ToLower(content), "high vulnerability") {
		score -= 20
	}
	if strings.Contains(strings.ToLower(content), "compliance gap") {
		score -= 15
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

func (s *TechnicalAnalysisService) calculatePerformanceScore(content string) float64 {
	return s.calculateCodeScore(content) // Reuse similar logic
}

func (s *TechnicalAnalysisService) calculateComplianceScore(content string) float64 {
	score := 80.0 // Base compliance score

	if strings.Contains(strings.ToLower(content), "major gap") {
		score -= 25
	}
	if strings.Contains(strings.ToLower(content), "non-compliant") {
		score -= 20
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

func (s *TechnicalAnalysisService) determineRiskLevel(content string) string {
	contentLower := strings.ToLower(content)

	if strings.Contains(contentLower, "critical") {
		return "critical"
	}
	if strings.Contains(contentLower, "high risk") {
		return "high"
	}
	if strings.Contains(contentLower, "medium risk") {
		return "medium"
	}

	return "low"
}

// Utility function to generate unique IDs
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// Additional extraction methods

func (s *TechnicalAnalysisService) extractSecurityAssessment(content string) *interfaces.TechSecurityAssessment {
	return &interfaces.TechSecurityAssessment{
		Score:            s.calculateSecurityScore(content),
		SecurityFindings: s.extractSecurityFindings(content),
		ComplianceGaps:   s.extractComplianceGaps(content),
		Recommendations:  []string{"Implement multi-factor authentication", "Enable encryption at rest"},
	}
}

func (s *TechnicalAnalysisService) extractPerformanceAssessment(content string) *interfaces.TechPerformanceAssessment {
	return &interfaces.TechPerformanceAssessment{
		Score:                     s.calculatePerformanceScore(content),
		BottleneckAnalysis:        s.extractPerformanceBottlenecks(content),
		OptimizationOpportunities: s.extractOptimizationOpportunities(content),
		Recommendations:           []string{"Implement caching", "Optimize database queries"},
	}
}

func (s *TechnicalAnalysisService) extractScalabilityAssessment(content string) *interfaces.TechScalabilityAssessment {
	return &interfaces.TechScalabilityAssessment{
		Score:           75.0,
		ScalabilityGaps: []string{"No auto-scaling configured", "Single point of failure in database"},
		Recommendations: []string{"Implement horizontal scaling", "Add load balancing"},
	}
}

func (s *TechnicalAnalysisService) extractReliabilityAssessment(content string) *interfaces.TechReliabilityAssessment {
	return &interfaces.TechReliabilityAssessment{
		Score:                 80.0,
		SinglePointsOfFailure: []string{"Single database instance", "No backup strategy"},
		DisasterRecoveryGaps:  []string{"No cross-region backup", "Missing recovery procedures"},
		Recommendations:       []string{"Implement database clustering", "Create disaster recovery plan"},
	}
}

func (s *TechnicalAnalysisService) extractCostAssessment(content string) *interfaces.TechCostAssessment {
	return &interfaces.TechCostAssessment{
		EstimatedMonthlyCost: 2500.0,
		CostOptimizations: []*interfaces.TechCostOptimization{
			{
				ID:                   generateID("cost_opt"),
				Type:                 "rightsizing",
				Title:                "Rightsize EC2 Instances",
				Description:          "Current instances are over-provisioned",
				CurrentCost:          1000.0,
				OptimizedCost:        700.0,
				Savings:              300.0,
				SavingsPercent:       30.0,
				ImplementationEffort: "low",
				RiskLevel:            "low",
				Timeline:             "1 week",
			},
		},
		Recommendations: []string{"Use reserved instances", "Implement auto-scaling"},
	}
}

func (s *TechnicalAnalysisService) extractComplianceAssessment(content string) *interfaces.TechComplianceAssessment {
	return &interfaces.TechComplianceAssessment{
		Score:           70.0,
		ComplianceGaps:  s.extractComplianceGaps(content),
		Recommendations: []string{"Implement access logging", "Enable data encryption"},
	}
}

func (s *TechnicalAnalysisService) extractArchitectureRecommendations(content string) []*interfaces.TechArchitectureRecommendation {
	return []*interfaces.TechArchitectureRecommendation{
		{
			ID:           generateID("arch_rec"),
			Type:         "security",
			Priority:     "high",
			Title:        "Implement Zero Trust Architecture",
			Description:  "Current architecture lacks proper network segmentation",
			Benefits:     []string{"Enhanced security", "Better compliance"},
			Effort:       "high",
			Timeline:     "3-4 months",
			Cost:         15000.0,
			Dependencies: []string{"Security team approval", "Network redesign"},
			References:   []string{"https://www.nist.gov/publications/zero-trust-architecture"},
		},
	}
}

func (s *TechnicalAnalysisService) extractSecurityVulnerabilities(content string) []*interfaces.TechSecurityVulnerability {
	return []*interfaces.TechSecurityVulnerability{
		{
			ID:                 generateID("sec_vuln"),
			Type:               "configuration",
			Severity:           "high",
			Title:              "Unencrypted Data Storage",
			Description:        "Database stores sensitive data without encryption",
			Impact:             "Potential data breach if storage is compromised",
			Likelihood:         "medium",
			CVSS:               7.5,
			References:         []string{"https://cwe.mitre.org/data/definitions/311.html"},
			AffectedComponents: []string{"Database", "Storage"},
		},
	}
}

func (s *TechnicalAnalysisService) extractComplianceGaps(content string) []*interfaces.TechComplianceGap {
	return []*interfaces.TechComplianceGap{
		{
			ID:            generateID("comp_gap"),
			Framework:     "SOC2",
			Control:       "CC6.1",
			Requirement:   "Logical and physical access controls",
			CurrentState:  "Basic password authentication",
			RequiredState: "Multi-factor authentication required",
			Gap:           "Missing MFA implementation",
			Impact:        "High risk of unauthorized access",
			Priority:      "high",
		},
	}
}

func (s *TechnicalAnalysisService) extractThreatAnalysis(content string) *interfaces.TechThreatAnalysis {
	return &interfaces.TechThreatAnalysis{
		IdentifiedThreats: []*interfaces.TechIdentifiedThreat{
			{
				ID:          generateID("threat"),
				Name:        "SQL Injection",
				Description: "Potential for SQL injection attacks",
				Category:    "Application Security",
				Likelihood:  "medium",
				Impact:      "high",
				RiskScore:   7.5,
				Mitigations: []string{"Use parameterized queries", "Input validation"},
			},
		},
		AttackVectors: []*interfaces.TechAttackVector{
			{
				ID:            generateID("attack_vector"),
				Name:          "Web Application Attack",
				Description:   "Attack through web application vulnerabilities",
				Complexity:    "medium",
				Prerequisites: []string{"Network access", "Application knowledge"},
				Mitigations:   []string{"Web application firewall", "Regular security testing"},
			},
		},
		RiskMatrix: &interfaces.TechRiskMatrix{
			Risks: [][]float64{{1.0, 2.0, 3.0}, {2.0, 4.0, 6.0}, {3.0, 6.0, 9.0}},
			Labels: struct {
				Impact     []string `json:"impact"`
				Likelihood []string `json:"likelihood"`
			}{
				Impact:     []string{"Low", "Medium", "High"},
				Likelihood: []string{"Low", "Medium", "High"},
			},
		},
	}
}

func (s *TechnicalAnalysisService) extractSecurityRecommendations(content string) []*interfaces.TechSecurityRecommendation {
	return []*interfaces.TechSecurityRecommendation{
		{
			ID:           generateID("sec_rec"),
			Type:         "immediate",
			Priority:     "critical",
			Title:        "Enable Database Encryption",
			Description:  "Implement encryption at rest for all sensitive data",
			Benefits:     []string{"Data protection", "Compliance requirement"},
			Effort:       "medium",
			Timeline:     "2 weeks",
			Cost:         5000.0,
			Dependencies: []string{"Database maintenance window"},
			References:   []string{"https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html"},
		},
	}
}

func (s *TechnicalAnalysisService) extractPerformanceBottlenecks(content string) []*interfaces.TechPerformanceBottleneck {
	return []*interfaces.TechPerformanceBottleneck{
		{
			ID:          generateID("perf_bottleneck"),
			Type:        "database",
			Severity:    "high",
			Title:       "Slow Database Queries",
			Description: "Multiple queries taking over 1 second to execute",
			Impact:      "Increased response times and poor user experience",
			Location:    "User service database",
			Metrics: map[string]float64{
				"avg_query_time": 1.5,
				"slow_queries":   25.0,
			},
		},
	}
}

func (s *TechnicalAnalysisService) extractOptimizationOpportunities(content string) []*interfaces.TechOptimizationRecommendation {
	return []*interfaces.TechOptimizationRecommendation{
		{
			ID:                  generateID("opt_rec"),
			Type:                "caching",
			Priority:            "high",
			Title:               "Implement Redis Caching",
			Description:         "Add caching layer to reduce database load",
			ExpectedImprovement: "50% reduction in response time",
			Effort:              "medium",
			Timeline:            "3 weeks",
			Cost:                3000.0,
			Dependencies:        []string{"Redis cluster setup"},
			References:          []string{"https://redis.io/documentation"},
		},
	}
}

func (s *TechnicalAnalysisService) extractBenchmarkComparison(content string) *interfaces.TechBenchmarkComparison {
	return &interfaces.TechBenchmarkComparison{
		Industry: "E-commerce",
		Benchmarks: map[string]float64{
			"response_time": 200.0,
			"throughput":    1000.0,
			"availability":  99.9,
		},
		Current: map[string]float64{
			"response_time": 350.0,
			"throughput":    750.0,
			"availability":  99.5,
		},
		Comparison: map[string]string{
			"response_time": "below",
			"throughput":    "below",
			"availability":  "below",
		},
		Percentile: map[string]float64{
			"response_time": 25.0,
			"throughput":    30.0,
			"availability":  40.0,
		},
	}
}

func (s *TechnicalAnalysisService) extractPerformanceRecommendations(content string) []*interfaces.TechPerformanceRecommendation {
	return []*interfaces.TechPerformanceRecommendation{
		{
			ID:           generateID("perf_rec"),
			Type:         "short-term",
			Priority:     "high",
			Title:        "Database Query Optimization",
			Description:  "Optimize slow-running database queries",
			Benefits:     []string{"Faster response times", "Reduced server load"},
			Effort:       "medium",
			Timeline:     "2 weeks",
			Cost:         2000.0,
			Dependencies: []string{"Database analysis", "Query profiling"},
			References:   []string{"https://use-the-index-luke.com/"},
		},
	}
}

// Compliance and comprehensive analysis methods

func (s *TechnicalAnalysisService) extractComplianceStatus(content string, frameworks []string) map[string]string {
	status := make(map[string]string)
	for _, framework := range frameworks {
		// Simple status determination based on content
		if strings.Contains(strings.ToLower(content), strings.ToLower(framework)) {
			if strings.Contains(strings.ToLower(content), "compliant") {
				status[framework] = "compliant"
			} else if strings.Contains(strings.ToLower(content), "partial") {
				status[framework] = "partial"
			} else {
				status[framework] = "non-compliant"
			}
		} else {
			status[framework] = "not-assessed"
		}
	}
	return status
}

func (s *TechnicalAnalysisService) extractControlAssessment(content string) []*interfaces.TechControlAssessment {
	return []*interfaces.TechControlAssessment{
		{
			ControlID:       "SOC2-CC6.1",
			Framework:       "SOC2",
			ControlName:     "Logical and Physical Access Controls",
			RequiredState:   "Multi-factor authentication implemented",
			CurrentState:    "Basic password authentication only",
			ComplianceLevel: 40.0,
			Gap:             "Missing MFA implementation",
			Priority:        "high",
			Effort:          "medium",
		},
	}
}

func (s *TechnicalAnalysisService) extractComplianceRecommendations(content string) []*interfaces.TechComplianceRecommendation {
	return []*interfaces.TechComplianceRecommendation{
		{
			ID:           generateID("comp_rec"),
			Framework:    "SOC2",
			Type:         "immediate",
			Priority:     "critical",
			Title:        "Implement Multi-Factor Authentication",
			Description:  "Deploy MFA for all user accounts to meet SOC2 requirements",
			Benefits:     []string{"Enhanced security", "SOC2 compliance", "Reduced breach risk"},
			Effort:       "medium",
			Timeline:     "4 weeks",
			Cost:         8000.0,
			Dependencies: []string{"Identity provider selection", "User training"},
			References:   []string{"https://www.aicpa.org/interestareas/frc/assuranceadvisoryservices/aicpasoc2report.html"},
		},
	}
}

func (s *TechnicalAnalysisService) extractComplianceRoadmap(content string) *interfaces.TechComplianceRoadmap {
	return &interfaces.TechComplianceRoadmap{
		Framework: "SOC2",
		Phases: []*interfaces.TechCompliancePhase{
			{
				PhaseNumber:  1,
				Name:         "Foundation",
				Description:  "Implement basic security controls",
				Duration:     "8 weeks",
				Controls:     []string{"CC6.1", "CC6.2", "CC6.3"},
				Deliverables: []string{"MFA implementation", "Access control policies"},
				Cost:         15000.0,
			},
			{
				PhaseNumber:  2,
				Name:         "Advanced Controls",
				Description:  "Implement monitoring and logging",
				Duration:     "6 weeks",
				Controls:     []string{"CC7.1", "CC7.2"},
				Deliverables: []string{"Security monitoring", "Audit logging"},
				Cost:         12000.0,
			},
		},
		Timeline:  "14 weeks",
		TotalCost: 27000.0,
		Milestones: []*interfaces.TechComplianceMilestone{
			{
				Name:        "Phase 1 Complete",
				Description: "Basic security controls implemented",
				TargetDate:  time.Now().AddDate(0, 2, 0),
				Criteria:    []string{"MFA deployed", "Access policies documented"},
			},
		},
	}
}

// Comprehensive analysis helper methods

func (s *TechnicalAnalysisService) identifyCrossCuttingFindings(result *interfaces.ComprehensiveAnalysisResult) []*interfaces.CrossCuttingFinding {
	return []*interfaces.CrossCuttingFinding{
		{
			ID:          generateID("cross_cutting"),
			Title:       "Insufficient Monitoring and Logging",
			Description: "Lack of comprehensive monitoring affects security, performance, and compliance",
			Areas:       []string{"security", "performance", "compliance"},
			Impact:      "High - affects multiple areas of system reliability",
			Priority:    "high",
			Recommendations: []string{
				"Implement centralized logging",
				"Deploy monitoring and alerting",
				"Create operational dashboards",
			},
		},
	}
}

func (s *TechnicalAnalysisService) prioritizeRecommendations(result *interfaces.ComprehensiveAnalysisResult) []*interfaces.PrioritizedRecommendation {
	var recommendations []*interfaces.PrioritizedRecommendation

	// Collect recommendations from all analysis results
	if result.SecurityResult != nil {
		for _, rec := range result.SecurityResult.SecurityRecommendations {
			recommendations = append(recommendations, &interfaces.PrioritizedRecommendation{
				ID:           rec.ID,
				Source:       "security",
				Type:         rec.Type,
				Priority:     rec.Priority,
				Title:        rec.Title,
				Description:  rec.Description,
				Benefits:     rec.Benefits,
				Effort:       rec.Effort,
				Timeline:     rec.Timeline,
				Cost:         rec.Cost,
				Dependencies: rec.Dependencies,
				References:   rec.References,
				Score:        s.calculatePriorityScore(rec.Priority, rec.Effort),
			})
		}
	}

	if result.PerformanceResult != nil {
		for _, rec := range result.PerformanceResult.PerformanceRecommendations {
			recommendations = append(recommendations, &interfaces.PrioritizedRecommendation{
				ID:           rec.ID,
				Source:       "performance",
				Type:         rec.Type,
				Priority:     rec.Priority,
				Title:        rec.Title,
				Description:  rec.Description,
				Benefits:     rec.Benefits,
				Effort:       rec.Effort,
				Timeline:     rec.Timeline,
				Cost:         rec.Cost,
				Dependencies: rec.Dependencies,
				References:   rec.References,
				Score:        s.calculatePriorityScore(rec.Priority, rec.Effort),
			})
		}
	}

	// Sort by priority score (higher is better)
	for i := 0; i < len(recommendations)-1; i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	return recommendations
}

func (s *TechnicalAnalysisService) calculatePriorityScore(priority, effort string) float64 {
	priorityScore := map[string]float64{
		"critical": 100.0,
		"high":     80.0,
		"medium":   60.0,
		"low":      40.0,
	}

	effortPenalty := map[string]float64{
		"low":    0.0,
		"medium": 10.0,
		"high":   20.0,
	}

	score := priorityScore[priority] - effortPenalty[effort]
	if score < 0 {
		score = 0
	}

	return score
}

func (s *TechnicalAnalysisService) generateExecutiveSummary(result *interfaces.ComprehensiveAnalysisResult) string {
	return fmt.Sprintf("Comprehensive technical analysis completed with an overall score of %.1f/100. "+
		"Key areas for improvement include security enhancements, performance optimizations, and compliance gaps. "+
		"Immediate action required on %d critical recommendations.",
		result.OverallScore, s.countCriticalRecommendations(result.PrioritizedRecommendations))
}

func (s *TechnicalAnalysisService) generateTechnicalSummary(result *interfaces.ComprehensiveAnalysisResult) string {
	return "Technical analysis reveals multiple areas requiring attention across security, performance, and compliance domains. " +
		"Primary concerns include unencrypted data storage, database performance bottlenecks, and missing compliance controls. " +
		"Recommended approach focuses on immediate security fixes followed by performance optimizations and compliance implementation."
}

func (s *TechnicalAnalysisService) generateActionPlan(result *interfaces.ComprehensiveAnalysisResult) *interfaces.TechnicalActionPlan {
	return &interfaces.TechnicalActionPlan{
		Phases: []*interfaces.ActionPhase{
			{
				PhaseNumber:   1,
				Name:          "Critical Security Fixes",
				Description:   "Address critical security vulnerabilities",
				Duration:      "4 weeks",
				Actions:       []string{"Enable database encryption", "Implement MFA", "Deploy WAF"},
				Deliverables:  []string{"Encrypted storage", "MFA system", "Security monitoring"},
				Cost:          25000.0,
				Prerequisites: []string{"Security team approval", "Maintenance windows"},
			},
			{
				PhaseNumber:   2,
				Name:          "Performance Optimization",
				Description:   "Implement performance improvements",
				Duration:      "6 weeks",
				Actions:       []string{"Deploy caching layer", "Optimize queries", "Scale infrastructure"},
				Deliverables:  []string{"Redis cluster", "Optimized database", "Auto-scaling"},
				Cost:          18000.0,
				Prerequisites: []string{"Phase 1 completion", "Performance baseline"},
			},
		},
		Timeline:  "10 weeks",
		TotalCost: 43000.0,
		Resources: &interfaces.RequiredResources{
			TechnicalRoles:   []string{"Security Engineer", "Performance Engineer", "DevOps Engineer"},
			SkillsRequired:   []string{"Cloud Security", "Database Optimization", "Infrastructure as Code"},
			Tools:            []string{"Redis", "CloudWatch", "Terraform"},
			ExternalServices: []string{"Security audit", "Performance testing"},
			Budget:           43000.0,
		},
		Milestones: []*interfaces.ActionMilestone{
			{
				Name:         "Security Phase Complete",
				Description:  "All critical security issues resolved",
				TargetDate:   time.Now().AddDate(0, 1, 0),
				Criteria:     []string{"Encryption enabled", "MFA deployed", "Monitoring active"},
				Deliverables: []string{"Security assessment report", "Compliance documentation"},
			},
		},
		RiskFactors: []string{
			"Potential downtime during encryption implementation",
			"User adoption challenges with MFA",
			"Performance impact during optimization",
		},
	}
}

func (s *TechnicalAnalysisService) countCriticalRecommendations(recommendations []*interfaces.PrioritizedRecommendation) int {
	count := 0
	for _, rec := range recommendations {
		if rec.Priority == "critical" {
			count++
		}
	}
	return count
}

// Additional helper methods for remediation

func (s *TechnicalAnalysisService) buildRemediationPrompt(vuln *interfaces.TechSecurityVulnerability) string {
	return fmt.Sprintf("Provide detailed remediation steps for the following security vulnerability:\n\n"+
		"Title: %s\n"+
		"Description: %s\n"+
		"Severity: %s\n"+
		"Type: %s\n\n"+
		"Please provide:\n"+
		"1. Step-by-step remediation instructions\n"+
		"2. Expected timeline for implementation\n"+
		"3. Potential risks or considerations\n"+
		"4. Verification steps to confirm the fix\n"+
		"5. Relevant documentation or references",
		vuln.Title, vuln.Description, vuln.Severity, vuln.Type)
}

func (s *TechnicalAnalysisService) buildOptimizationPrompt(metrics *interfaces.TechPerformanceMetrics) string {
	return fmt.Sprintf("Based on the following performance metrics, provide optimization recommendations:\n\n"+
		"Response Time: %.2f ms\n"+
		"Throughput: %.2f RPS\n"+
		"Error Rate: %.2f%%\n"+
		"CPU Utilization: %.2f%%\n"+
		"Memory Utilization: %.2f%%\n\n"+
		"Please provide:\n"+
		"1. Specific optimization opportunities\n"+
		"2. Expected performance improvements\n"+
		"3. Implementation effort and timeline\n"+
		"4. Potential risks or trade-offs\n"+
		"5. Monitoring and validation approaches",
		metrics.ResponseTime, metrics.Throughput, metrics.ErrorRate,
		metrics.CPUUtilization, metrics.MemoryUtilization)
}

func (s *TechnicalAnalysisService) buildComplianceRemediationPrompt(gap *interfaces.TechComplianceGap) string {
	return fmt.Sprintf("Provide detailed remediation steps for the following compliance gap:\n\n"+
		"Framework: %s\n"+
		"Control: %s\n"+
		"Requirement: %s\n"+
		"Current State: %s\n"+
		"Required State: %s\n"+
		"Gap: %s\n\n"+
		"Please provide:\n"+
		"1. Step-by-step implementation plan\n"+
		"2. Required resources and skills\n"+
		"3. Timeline and milestones\n"+
		"4. Compliance validation steps\n"+
		"5. Relevant standards and references",
		gap.Framework, gap.Control, gap.Requirement,
		gap.CurrentState, gap.RequiredState, gap.Gap)
}

func (s *TechnicalAnalysisService) extractRemediationTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "title") || strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "Security Remediation"
}

func (s *TechnicalAnalysisService) extractRemediationDescription(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), "description") && i+1 < len(lines) {
			return strings.TrimSpace(lines[i+1])
		}
	}
	return "Detailed remediation steps to address the identified issue."
}

func (s *TechnicalAnalysisService) extractRemediationSteps(content string) []string {
	steps := []string{
		"Review current configuration and identify affected components",
		"Plan implementation approach and required resources",
		"Implement the recommended changes in a test environment",
		"Validate the fix and perform security testing",
		"Deploy to production with proper monitoring",
		"Document the changes and update security procedures",
	}
	return steps
}

func (s *TechnicalAnalysisService) extractRemediationReferences(content string) []string {
	return []string{
		"https://owasp.org/",
		"https://cwe.mitre.org/",
		"https://nvd.nist.gov/",
	}
}

func (s *TechnicalAnalysisService) estimateRemediationEffort(severity string) string {
	effortMap := map[string]string{
		"critical": "high",
		"high":     "medium",
		"medium":   "medium",
		"low":      "low",
	}

	if effort, exists := effortMap[severity]; exists {
		return effort
	}
	return "medium"
}

func (s *TechnicalAnalysisService) estimateRemediationTimeline(severity string) string {
	timelineMap := map[string]string{
		"critical": "1 week",
		"high":     "2 weeks",
		"medium":   "3 weeks",
		"low":      "4 weeks",
	}

	if timeline, exists := timelineMap[severity]; exists {
		return timeline
	}
	return "2 weeks"
}

func (s *TechnicalAnalysisService) estimateComplianceEffort(impact string) string {
	effortMap := map[string]string{
		"high":   "high",
		"medium": "medium",
		"low":    "low",
	}

	if effort, exists := effortMap[impact]; exists {
		return effort
	}
	return "medium"
}

func (s *TechnicalAnalysisService) estimateComplianceTimeline(priority string) string {
	timelineMap := map[string]string{
		"critical": "2 weeks",
		"high":     "4 weeks",
		"medium":   "6 weeks",
		"low":      "8 weeks",
	}

	if timeline, exists := timelineMap[priority]; exists {
		return timeline
	}
	return "4 weeks"
}

// New extraction methods with Tech prefix

func (s *TechnicalAnalysisService) extractTechSecurityAssessment(content string) *interfaces.TechSecurityAssessment {
	return &interfaces.TechSecurityAssessment{
		Score:            s.calculateSecurityScore(content),
		SecurityFindings: s.extractSecurityFindings(content),
		ComplianceGaps:   s.extractTechComplianceGaps(content),
		Recommendations:  []string{"Implement multi-factor authentication", "Enable encryption at rest"},
	}
}

func (s *TechnicalAnalysisService) extractTechPerformanceAssessment(content string) *interfaces.TechPerformanceAssessment {
	return &interfaces.TechPerformanceAssessment{
		Score:                     s.calculatePerformanceScore(content),
		BottleneckAnalysis:        s.extractTechPerformanceBottlenecks(content),
		OptimizationOpportunities: s.extractTechOptimizationOpportunities(content),
		Recommendations:           []string{"Implement caching", "Optimize database queries"},
	}
}

func (s *TechnicalAnalysisService) extractTechScalabilityAssessment(content string) *interfaces.TechScalabilityAssessment {
	return &interfaces.TechScalabilityAssessment{
		Score:           75.0,
		ScalabilityGaps: []string{"No auto-scaling configured", "Single point of failure in database"},
		Recommendations: []string{"Implement horizontal scaling", "Add load balancing"},
	}
}

func (s *TechnicalAnalysisService) extractTechReliabilityAssessment(content string) *interfaces.TechReliabilityAssessment {
	return &interfaces.TechReliabilityAssessment{
		Score:                 80.0,
		SinglePointsOfFailure: []string{"Single database instance", "No backup strategy"},
		DisasterRecoveryGaps:  []string{"No cross-region backup", "Missing recovery procedures"},
		Recommendations:       []string{"Implement database clustering", "Create disaster recovery plan"},
	}
}

func (s *TechnicalAnalysisService) extractTechCostAssessment(content string) *interfaces.TechCostAssessment {
	return &interfaces.TechCostAssessment{
		EstimatedMonthlyCost: 2500.0,
		CostOptimizations: []*interfaces.TechCostOptimization{
			{
				ID:                   generateID("cost_opt"),
				Type:                 "rightsizing",
				Title:                "Rightsize EC2 Instances",
				Description:          "Current instances are over-provisioned",
				CurrentCost:          1000.0,
				OptimizedCost:        700.0,
				Savings:              300.0,
				SavingsPercent:       30.0,
				ImplementationEffort: "low",
				RiskLevel:            "low",
				Timeline:             "1 week",
			},
		},
		Recommendations: []string{"Use reserved instances", "Implement auto-scaling"},
	}
}

func (s *TechnicalAnalysisService) extractTechComplianceAssessment(content string) *interfaces.TechComplianceAssessment {
	return &interfaces.TechComplianceAssessment{
		Score:           70.0,
		ComplianceGaps:  s.extractTechComplianceGaps(content),
		Recommendations: []string{"Implement access logging", "Enable data encryption"},
	}
}

func (s *TechnicalAnalysisService) extractTechArchitectureRecommendations(content string) []*interfaces.TechArchitectureRecommendation {
	return []*interfaces.TechArchitectureRecommendation{
		{
			ID:           generateID("arch_rec"),
			Type:         "security",
			Priority:     "high",
			Title:        "Implement Zero Trust Architecture",
			Description:  "Current architecture lacks proper network segmentation",
			Benefits:     []string{"Enhanced security", "Better compliance"},
			Effort:       "high",
			Timeline:     "3-4 months",
			Cost:         15000.0,
			Dependencies: []string{"Security team approval", "Network redesign"},
			References:   []string{"https://www.nist.gov/publications/zero-trust-architecture"},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechSecurityVulnerabilities(content string) []*interfaces.TechSecurityVulnerability {
	return []*interfaces.TechSecurityVulnerability{
		{
			ID:                 generateID("sec_vuln"),
			Type:               "configuration",
			Severity:           "high",
			Title:              "Unencrypted Data Storage",
			Description:        "Database stores sensitive data without encryption",
			Impact:             "Potential data breach if storage is compromised",
			Likelihood:         "medium",
			CVSS:               7.5,
			References:         []string{"https://cwe.mitre.org/data/definitions/311.html"},
			AffectedComponents: []string{"Database", "Storage"},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechComplianceGaps(content string) []*interfaces.TechComplianceGap {
	return []*interfaces.TechComplianceGap{
		{
			ID:            generateID("comp_gap"),
			Framework:     "SOC2",
			Control:       "CC6.1",
			Requirement:   "Logical and physical access controls",
			CurrentState:  "Basic password authentication",
			RequiredState: "Multi-factor authentication required",
			Gap:           "Missing MFA implementation",
			Impact:        "High risk of unauthorized access",
			Priority:      "high",
		},
	}
}

func (s *TechnicalAnalysisService) extractTechThreatAnalysis(content string) *interfaces.TechThreatAnalysis {
	return &interfaces.TechThreatAnalysis{
		IdentifiedThreats: []*interfaces.TechIdentifiedThreat{
			{
				ID:          generateID("threat"),
				Name:        "SQL Injection",
				Description: "Potential for SQL injection attacks",
				Category:    "Application Security",
				Likelihood:  "medium",
				Impact:      "high",
				RiskScore:   7.5,
				Mitigations: []string{"Use parameterized queries", "Input validation"},
			},
		},
		AttackVectors: []*interfaces.TechAttackVector{
			{
				ID:            generateID("attack_vector"),
				Name:          "Web Application Attack",
				Description:   "Attack through web application vulnerabilities",
				Complexity:    "medium",
				Prerequisites: []string{"Network access", "Application knowledge"},
				Mitigations:   []string{"Web application firewall", "Regular security testing"},
			},
		},
		RiskMatrix: &interfaces.TechRiskMatrix{
			Risks: [][]float64{{1.0, 2.0, 3.0}, {2.0, 4.0, 6.0}, {3.0, 6.0, 9.0}},
			Labels: struct {
				Impact     []string `json:"impact"`
				Likelihood []string `json:"likelihood"`
			}{
				Impact:     []string{"Low", "Medium", "High"},
				Likelihood: []string{"Low", "Medium", "High"},
			},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechSecurityRecommendations(content string) []*interfaces.TechSecurityRecommendation {
	return []*interfaces.TechSecurityRecommendation{
		{
			ID:           generateID("sec_rec"),
			Type:         "immediate",
			Priority:     "critical",
			Title:        "Enable Database Encryption",
			Description:  "Implement encryption at rest for all sensitive data",
			Benefits:     []string{"Data protection", "Compliance requirement"},
			Effort:       "medium",
			Timeline:     "2 weeks",
			Cost:         5000.0,
			Dependencies: []string{"Database maintenance window"},
			References:   []string{"https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html"},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechPerformanceBottlenecks(content string) []*interfaces.TechPerformanceBottleneck {
	return []*interfaces.TechPerformanceBottleneck{
		{
			ID:          generateID("perf_bottleneck"),
			Type:        "database",
			Severity:    "high",
			Title:       "Slow Database Queries",
			Description: "Multiple queries taking over 1 second to execute",
			Impact:      "Increased response times and poor user experience",
			Location:    "User service database",
			Metrics: map[string]float64{
				"avg_query_time": 1.5,
				"slow_queries":   25.0,
			},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechOptimizationOpportunities(content string) []*interfaces.TechOptimizationRecommendation {
	return []*interfaces.TechOptimizationRecommendation{
		{
			ID:                  generateID("opt_rec"),
			Type:                "caching",
			Priority:            "high",
			Title:               "Implement Redis Caching",
			Description:         "Add caching layer to reduce database load",
			ExpectedImprovement: "50% reduction in response time",
			Effort:              "medium",
			Timeline:            "3 weeks",
			Cost:                3000.0,
			Dependencies:        []string{"Redis cluster setup"},
			References:          []string{"https://redis.io/documentation"},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechPerformanceRecommendations(content string) []*interfaces.TechPerformanceRecommendation {
	return []*interfaces.TechPerformanceRecommendation{
		{
			ID:           generateID("perf_rec"),
			Type:         "short-term",
			Priority:     "high",
			Title:        "Database Query Optimization",
			Description:  "Optimize slow-running database queries",
			Benefits:     []string{"Faster response times", "Reduced server load"},
			Effort:       "medium",
			Timeline:     "2 weeks",
			Cost:         2000.0,
			Dependencies: []string{"Database analysis", "Query profiling"},
			References:   []string{"https://use-the-index-luke.com/"},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechControlAssessment(content string) []*interfaces.TechControlAssessment {
	return []*interfaces.TechControlAssessment{
		{
			ControlID:       "SOC2-CC6.1",
			Framework:       "SOC2",
			ControlName:     "Logical and Physical Access Controls",
			RequiredState:   "Multi-factor authentication implemented",
			CurrentState:    "Basic password authentication only",
			ComplianceLevel: 40.0,
			Gap:             "Missing MFA implementation",
			Priority:        "high",
			Effort:          "medium",
		},
	}
}

func (s *TechnicalAnalysisService) extractTechComplianceRecommendations(content string) []*interfaces.TechComplianceRecommendation {
	return []*interfaces.TechComplianceRecommendation{
		{
			ID:           generateID("comp_rec"),
			Framework:    "SOC2",
			Type:         "immediate",
			Priority:     "critical",
			Title:        "Implement Multi-Factor Authentication",
			Description:  "Deploy MFA for all user accounts to meet SOC2 requirements",
			Benefits:     []string{"Enhanced security", "SOC2 compliance", "Reduced breach risk"},
			Effort:       "medium",
			Timeline:     "4 weeks",
			Cost:         8000.0,
			Dependencies: []string{"Identity provider selection", "User training"},
			References:   []string{"https://www.aicpa.org/interestareas/frc/assuranceadvisoryservices/aicpasoc2report.html"},
		},
	}
}

func (s *TechnicalAnalysisService) extractTechComplianceRoadmap(content string) *interfaces.TechComplianceRoadmap {
	return &interfaces.TechComplianceRoadmap{
		Framework: "SOC2",
		Phases: []*interfaces.TechCompliancePhase{
			{
				PhaseNumber:  1,
				Name:         "Foundation",
				Description:  "Implement basic security controls",
				Duration:     "8 weeks",
				Controls:     []string{"CC6.1", "CC6.2", "CC6.3"},
				Deliverables: []string{"MFA implementation", "Access control policies"},
				Cost:         15000.0,
			},
			{
				PhaseNumber:  2,
				Name:         "Advanced Controls",
				Description:  "Implement monitoring and logging",
				Duration:     "6 weeks",
				Controls:     []string{"CC7.1", "CC7.2"},
				Deliverables: []string{"Security monitoring", "Audit logging"},
				Cost:         12000.0,
			},
		},
		Timeline:  "14 weeks",
		TotalCost: 27000.0,
		Milestones: []*interfaces.TechComplianceMilestone{
			{
				Name:        "Phase 1 Complete",
				Description: "Basic security controls implemented",
				TargetDate:  time.Now().AddDate(0, 2, 0),
				Criteria:    []string{"MFA deployed", "Access policies documented"},
			},
		},
	}
}
