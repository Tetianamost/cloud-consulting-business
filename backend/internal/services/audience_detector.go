package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloud-consulting/backend/internal/domain"
)

// AudienceType represents the type of audience for content generation
type AudienceType string

const (
	AudienceTechnical AudienceType = "technical"
	AudienceBusiness  AudienceType = "business"
	AudienceMixed     AudienceType = "mixed"
	AudienceExecutive AudienceType = "executive"
)

// AudienceProfile represents the detected audience characteristics
type AudienceProfile struct {
	PrimaryType     AudienceType `json:"primary_type"`
	SecondaryType   AudienceType `json:"secondary_type,omitempty"`
	TechnicalDepth  int          `json:"technical_depth"`  // 1-5 scale
	BusinessFocus   int          `json:"business_focus"`   // 1-5 scale
	ExecutiveLevel  bool         `json:"executive_level"`
	Confidence      float64      `json:"confidence"`       // 0.0-1.0
	Indicators      []string     `json:"indicators"`       // What led to this classification
}

// ContentTemplate represents a template for audience-specific content
type ContentTemplate struct {
	AudienceType    AudienceType `json:"audience_type"`
	TechnicalDepth  int          `json:"technical_depth"`
	BusinessFocus   int          `json:"business_focus"`
	SectionTemplate string       `json:"section_template"`
	ToneGuidelines  []string     `json:"tone_guidelines"`
	ContentFocus    []string     `json:"content_focus"`
}

// AudienceDetector provides audience detection and content adaptation capabilities
type AudienceDetector interface {
	DetectAudience(ctx context.Context, inquiry *domain.Inquiry) (*AudienceProfile, error)
	GetContentTemplate(audienceType AudienceType, technicalDepth int) (*ContentTemplate, error)
	AdaptContentForAudience(content string, profile *AudienceProfile) (string, error)
	SeparateBusinessAndTechnical(content string) (*SeparatedContent, error)
}

// SeparatedContent represents content separated into business and technical sections
type SeparatedContent struct {
	BusinessJustification string   `json:"business_justification"`
	TechnicalExplanation  string   `json:"technical_explanation"`
	SharedContent         string   `json:"shared_content"`
	Recommendations       []string `json:"recommendations"`
}

// audienceDetector implements the AudienceDetector interface
type audienceDetector struct {
	templates map[string]*ContentTemplate
}

// NewAudienceDetector creates a new AudienceDetector instance
func NewAudienceDetector() AudienceDetector {
	ad := &audienceDetector{
		templates: make(map[string]*ContentTemplate),
	}
	
	// Initialize default content templates
	ad.initializeContentTemplates()
	
	return ad
}

// DetectAudience analyzes an inquiry to determine the target audience
func (ad *audienceDetector) DetectAudience(ctx context.Context, inquiry *domain.Inquiry) (*AudienceProfile, error) {
	profile := &AudienceProfile{
		TechnicalDepth: 3, // Default medium depth
		BusinessFocus:  3, // Default medium focus
		Indicators:     []string{},
	}

	// Analyze inquiry content for audience indicators
	content := strings.ToLower(inquiry.Message + " " + inquiry.Company + " " + strings.Join(inquiry.Services, " "))
	
	// Technical indicators
	technicalScore := ad.calculateTechnicalScore(content, profile)
	
	// Business indicators
	businessScore := ad.calculateBusinessScore(content, profile)
	
	// Executive indicators
	executiveScore := ad.calculateExecutiveScore(content, profile)
	
	// Determine primary audience type
	profile.PrimaryType = ad.determinePrimaryAudience(technicalScore, businessScore, executiveScore)
	
	// Set technical depth and business focus based on scores
	profile.TechnicalDepth = ad.mapScoreToDepth(technicalScore)
	profile.BusinessFocus = ad.mapScoreToFocus(businessScore)
	profile.ExecutiveLevel = executiveScore > 0.6
	
	// Calculate confidence based on score differences
	profile.Confidence = ad.calculateConfidence(technicalScore, businessScore, executiveScore)
	
	// Determine secondary audience if mixed
	if profile.Confidence < 0.7 {
		profile.SecondaryType = ad.determineSecondaryAudience(technicalScore, businessScore, executiveScore, profile.PrimaryType)
		if profile.SecondaryType != "" {
			profile.PrimaryType = AudienceMixed
		}
	}

	return profile, nil
}

// calculateTechnicalScore analyzes content for technical indicators
func (ad *audienceDetector) calculateTechnicalScore(content string, profile *AudienceProfile) float64 {
	technicalKeywords := map[string]float64{
		// High technical indicators
		"architecture":     0.8,
		"infrastructure":   0.8,
		"api":             0.7,
		"microservices":   0.9,
		"kubernetes":      0.9,
		"docker":          0.8,
		"devops":          0.8,
		"ci/cd":           0.9,
		"database":        0.7,
		"performance":     0.6,
		"scalability":     0.7,
		"security":        0.6,
		"integration":     0.7,
		"deployment":      0.7,
		"monitoring":      0.6,
		"logging":         0.6,
		"networking":      0.8,
		"load balancer":   0.8,
		"auto scaling":    0.8,
		"serverless":      0.8,
		"lambda":          0.8,
		"containers":      0.8,
		"vpc":             0.9,
		"subnet":          0.9,
		"firewall":        0.7,
		"encryption":      0.7,
		"ssl":             0.7,
		"tls":             0.8,
		"oauth":           0.8,
		"saml":            0.8,
		"ldap":            0.7,
		"rest":            0.7,
		"graphql":         0.8,
		"json":            0.6,
		"xml":             0.6,
		"yaml":            0.7,
		"terraform":       0.9,
		"cloudformation":  0.9,
		"ansible":         0.8,
		"chef":            0.8,
		"puppet":          0.8,
		
		// Medium technical indicators
		"cloud":           0.4,
		"server":          0.4,
		"application":     0.3,
		"system":          0.3,
		"platform":        0.3,
		"service":         0.2,
		"solution":        0.2,
		"technology":      0.4,
		"technical":       0.5,
		"implementation":  0.4,
		"configuration":   0.5,
		"setup":           0.3,
		"install":         0.4,
		"upgrade":         0.4,
		"maintenance":     0.4,
		"backup":          0.5,
		"disaster recovery": 0.6,
		"high availability": 0.7,
		"redundancy":      0.6,
		"failover":        0.7,
	}

	score := 0.0
	matchCount := 0
	
	for keyword, weight := range technicalKeywords {
		if strings.Contains(content, keyword) {
			score += weight
			matchCount++
			profile.Indicators = append(profile.Indicators, "Technical: "+keyword)
		}
	}
	
	// Normalize score
	if matchCount > 0 {
		score = score / float64(matchCount)
	}
	
	// Boost score for multiple technical terms
	if matchCount > 5 {
		score = score * 1.2
	}
	
	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}
	
	return score
}

// calculateBusinessScore analyzes content for business indicators
func (ad *audienceDetector) calculateBusinessScore(content string, profile *AudienceProfile) float64 {
	businessKeywords := map[string]float64{
		// High business indicators
		"roi":             0.9,
		"return on investment": 0.9,
		"cost savings":    0.8,
		"budget":          0.7,
		"business case":   0.9,
		"competitive advantage": 0.8,
		"market":          0.6,
		"revenue":         0.8,
		"profit":          0.8,
		"efficiency":      0.7,
		"productivity":    0.7,
		"growth":          0.6,
		"strategy":        0.7,
		"strategic":       0.7,
		"business value":  0.9,
		"value proposition": 0.8,
		"stakeholder":     0.7,
		"executive":       0.8,
		"board":           0.9,
		"c-level":         0.9,
		"ceo":             0.9,
		"cto":             0.8,
		"cfo":             0.9,
		"cio":             0.8,
		"decision maker":  0.8,
		"business owner":  0.8,
		"investment":      0.7,
		"funding":         0.7,
		"approval":        0.6,
		"business impact": 0.8,
		"operational efficiency": 0.7,
		"cost optimization": 0.7,
		"business continuity": 0.7,
		"risk management": 0.6,
		"compliance":      0.6,
		"governance":      0.6,
		"audit":           0.6,
		"regulatory":      0.6,
		
		// Medium business indicators
		"business":        0.4,
		"organization":    0.3,
		"company":         0.2,
		"enterprise":      0.4,
		"corporate":       0.4,
		"management":      0.4,
		"operations":      0.3,
		"process":         0.3,
		"workflow":        0.4,
		"customer":        0.4,
		"client":          0.3,
		"user":            0.2,
		"team":            0.2,
		"department":      0.3,
		"division":        0.3,
		"project":         0.2,
		"initiative":      0.3,
		"program":         0.3,
		"timeline":        0.4,
		"deadline":        0.5,
		"milestone":       0.4,
		"deliverable":     0.4,
		"outcome":         0.4,
		"result":          0.3,
		"benefit":         0.5,
		"advantage":       0.4,
		"opportunity":     0.4,
		"challenge":       0.3,
		"problem":         0.3,
		"solution":        0.2,
		"requirement":     0.3,
		"objective":       0.4,
		"goal":            0.4,
		"target":          0.3,
		"kpi":             0.6,
		"metric":          0.5,
		"measurement":     0.4,
		"performance":     0.3,
		"quality":         0.3,
		"standard":        0.3,
		"policy":          0.4,
		"procedure":       0.3,
		"guideline":       0.3,
	}

	score := 0.0
	matchCount := 0
	
	for keyword, weight := range businessKeywords {
		if strings.Contains(content, keyword) {
			score += weight
			matchCount++
			profile.Indicators = append(profile.Indicators, "Business: "+keyword)
		}
	}
	
	// Normalize score
	if matchCount > 0 {
		score = score / float64(matchCount)
	}
	
	// Boost score for multiple business terms
	if matchCount > 5 {
		score = score * 1.2
	}
	
	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}
	
	return score
}

// calculateExecutiveScore analyzes content for executive-level indicators
func (ad *audienceDetector) calculateExecutiveScore(content string, profile *AudienceProfile) float64 {
	executiveKeywords := map[string]float64{
		"ceo":             1.0,
		"cto":             1.0,
		"cfo":             1.0,
		"cio":             1.0,
		"president":       0.9,
		"vice president":  0.9,
		"vp":              0.9,
		"director":        0.7,
		"executive":       0.8,
		"c-level":         1.0,
		"c-suite":         1.0,
		"board":           0.9,
		"board of directors": 0.9,
		"senior management": 0.8,
		"leadership":      0.7,
		"strategic":       0.6,
		"vision":          0.6,
		"transformation":  0.7,
		"digital transformation": 0.8,
		"business transformation": 0.8,
		"enterprise":      0.5,
		"corporate":       0.5,
		"organization":    0.3,
		"high level":      0.6,
		"overview":        0.4,
		"summary":         0.4,
		"executive summary": 0.8,
		"business case":   0.7,
		"investment":      0.6,
		"funding":         0.6,
		"approval":        0.6,
		"decision":        0.5,
		"stakeholder":     0.5,
		"governance":      0.6,
		"compliance":      0.5,
		"risk":            0.5,
		"competitive":     0.6,
		"market":          0.5,
		"industry":        0.4,
	}

	score := 0.0
	matchCount := 0
	
	for keyword, weight := range executiveKeywords {
		if strings.Contains(content, keyword) {
			score += weight
			matchCount++
			profile.Indicators = append(profile.Indicators, "Executive: "+keyword)
		}
	}
	
	// Normalize score
	if matchCount > 0 {
		score = score / float64(matchCount)
	}
	
	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}
	
	return score
}

// determinePrimaryAudience determines the primary audience type based on scores
func (ad *audienceDetector) determinePrimaryAudience(technicalScore, businessScore, executiveScore float64) AudienceType {
	// Executive takes precedence if score is high enough
	if executiveScore > 0.6 {
		return AudienceExecutive
	}
	
	// If scores are close, it's mixed
	if abs(technicalScore-businessScore) < 0.2 {
		return AudienceMixed
	}
	
	// Otherwise, choose the higher score
	if technicalScore > businessScore {
		return AudienceTechnical
	}
	
	return AudienceBusiness
}

// determineSecondaryAudience determines secondary audience for mixed audiences
func (ad *audienceDetector) determineSecondaryAudience(technicalScore, businessScore, executiveScore float64, primaryType AudienceType) AudienceType {
	if primaryType == AudienceExecutive {
		if technicalScore > businessScore {
			return AudienceTechnical
		}
		return AudienceBusiness
	}
	
	if primaryType == AudienceTechnical && businessScore > 0.3 {
		return AudienceBusiness
	}
	
	if primaryType == AudienceBusiness && technicalScore > 0.3 {
		return AudienceTechnical
	}
	
	return ""
}

// mapScoreToDepth maps technical score to depth level (1-5)
func (ad *audienceDetector) mapScoreToDepth(score float64) int {
	if score >= 0.8 {
		return 5 // Very high technical depth
	} else if score >= 0.6 {
		return 4 // High technical depth
	} else if score >= 0.4 {
		return 3 // Medium technical depth
	} else if score >= 0.2 {
		return 2 // Low technical depth
	}
	return 1 // Very low technical depth
}

// mapScoreToFocus maps business score to focus level (1-5)
func (ad *audienceDetector) mapScoreToFocus(score float64) int {
	if score >= 0.8 {
		return 5 // Very high business focus
	} else if score >= 0.6 {
		return 4 // High business focus
	} else if score >= 0.4 {
		return 3 // Medium business focus
	} else if score >= 0.2 {
		return 2 // Low business focus
	}
	return 1 // Very low business focus
}

// calculateConfidence calculates confidence in audience detection
func (ad *audienceDetector) calculateConfidence(technicalScore, businessScore, executiveScore float64) float64 {
	maxScore := max(technicalScore, businessScore, executiveScore)
	minScore := min(technicalScore, businessScore, executiveScore)
	
	// Higher difference between scores means higher confidence
	difference := maxScore - minScore
	
	// Base confidence on the difference and the absolute values
	confidence := difference * 0.7 + maxScore * 0.3
	
	// Ensure confidence is between 0.0 and 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	return confidence
}

// GetContentTemplate retrieves a content template for the specified audience
func (ad *audienceDetector) GetContentTemplate(audienceType AudienceType, technicalDepth int) (*ContentTemplate, error) {
	key := string(audienceType) + "_" + string(rune(technicalDepth+'0'))
	
	contentTemplate, exists := ad.templates[key]
	if !exists {
		// Fall back to default template for the audience type
		contentTemplate, exists = ad.templates[string(audienceType)]
		if !exists {
			return nil, fmt.Errorf("no template found for audience type: %s", audienceType)
		}
	}
	
	return contentTemplate, nil
}

// AdaptContentForAudience adapts content based on the audience profile
func (ad *audienceDetector) AdaptContentForAudience(content string, profile *AudienceProfile) (string, error) {
	_, err := ad.GetContentTemplate(profile.PrimaryType, profile.TechnicalDepth)
	if err != nil {
		return content, err // Return original content if template not found
	}
	
	// Apply audience-specific adaptations
	adaptedContent := content
	
	switch profile.PrimaryType {
	case AudienceTechnical:
		adaptedContent = ad.enhanceTechnicalContent(content, profile.TechnicalDepth)
	case AudienceBusiness:
		adaptedContent = ad.enhanceBusinessContent(content, profile.BusinessFocus)
	case AudienceExecutive:
		adaptedContent = ad.enhanceExecutiveContent(content)
	case AudienceMixed:
		adaptedContent = ad.enhanceMixedContent(content, profile)
	}
	
	return adaptedContent, nil
}

// enhanceTechnicalContent enhances content for technical audiences
func (ad *audienceDetector) enhanceTechnicalContent(content string, technicalDepth int) string {
	// Add technical depth based on the level
	enhancements := []string{}
	
	if technicalDepth >= 3 {
		enhancements = append(enhancements, 
			"\n\nTECHNICAL IMPLEMENTATION DETAILS:",
			"- Include specific service configurations and parameters",
			"- Provide architecture diagrams and technical specifications",
			"- Detail integration patterns and API specifications",
		)
	}
	
	if technicalDepth >= 4 {
		enhancements = append(enhancements,
			"- Include code examples and configuration snippets",
			"- Provide performance benchmarks and optimization strategies",
			"- Detail monitoring and alerting configurations",
		)
	}
	
	if technicalDepth >= 5 {
		enhancements = append(enhancements,
			"- Include advanced architectural patterns and design principles",
			"- Provide detailed security configurations and best practices",
			"- Include automation scripts and infrastructure as code examples",
		)
	}
	
	if len(enhancements) > 0 {
		content += strings.Join(enhancements, "\n")
	}
	
	return content
}

// enhanceBusinessContent enhances content for business audiences
func (ad *audienceDetector) enhanceBusinessContent(content string, businessFocus int) string {
	enhancements := []string{}
	
	if businessFocus >= 3 {
		enhancements = append(enhancements,
			"\n\nBUSINESS VALUE ANALYSIS:",
			"- Include ROI calculations and cost-benefit analysis",
			"- Provide business impact assessment and success metrics",
			"- Detail competitive advantages and market positioning",
		)
	}
	
	if businessFocus >= 4 {
		enhancements = append(enhancements,
			"- Include financial projections and budget implications",
			"- Provide risk assessment from business perspective",
			"- Detail change management and organizational impact",
		)
	}
	
	if businessFocus >= 5 {
		enhancements = append(enhancements,
			"- Include strategic alignment and long-term vision",
			"- Provide market analysis and competitive intelligence",
			"- Detail governance and compliance implications",
		)
	}
	
	if len(enhancements) > 0 {
		content += strings.Join(enhancements, "\n")
	}
	
	return content
}

// enhanceExecutiveContent enhances content for executive audiences
func (ad *audienceDetector) enhanceExecutiveContent(content string) string {
	enhancements := []string{
		"\n\nEXECUTIVE SUMMARY ENHANCEMENTS:",
		"- Focus on strategic implications and business transformation",
		"- Highlight competitive advantages and market opportunities",
		"- Emphasize ROI, cost savings, and business value creation",
		"- Include risk mitigation and governance considerations",
		"- Provide high-level timeline and resource requirements",
		"- Detail stakeholder impact and change management needs",
	}
	
	return content + strings.Join(enhancements, "\n")
}

// enhanceMixedContent enhances content for mixed audiences
func (ad *audienceDetector) enhanceMixedContent(content string, profile *AudienceProfile) string {
	// Combine both technical and business enhancements
	technicalContent := ad.enhanceTechnicalContent("", profile.TechnicalDepth)
	businessContent := ad.enhanceBusinessContent("", profile.BusinessFocus)
	
	enhancements := []string{
		"\n\nMIXED AUDIENCE CONSIDERATIONS:",
		"- Provide both technical details and business justification",
		"- Include executive summary for leadership review",
		"- Detail implementation approach with business context",
	}
	
	if technicalContent != "" {
		enhancements = append(enhancements, technicalContent)
	}
	
	if businessContent != "" {
		enhancements = append(enhancements, businessContent)
	}
	
	return content + strings.Join(enhancements, "\n")
}

// SeparateBusinessAndTechnical separates content into business and technical sections
func (ad *audienceDetector) SeparateBusinessAndTechnical(content string) (*SeparatedContent, error) {
	separated := &SeparatedContent{
		Recommendations: []string{},
	}
	
	// Split content into sections
	sections := strings.Split(content, "\n\n")
	
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		
		// Classify section as business, technical, or shared
		classification := ad.classifySection(section)
		
		switch classification {
		case "business":
			if separated.BusinessJustification != "" {
				separated.BusinessJustification += "\n\n"
			}
			separated.BusinessJustification += section
		case "technical":
			if separated.TechnicalExplanation != "" {
				separated.TechnicalExplanation += "\n\n"
			}
			separated.TechnicalExplanation += section
		case "recommendation":
			separated.Recommendations = append(separated.Recommendations, section)
		default:
			if separated.SharedContent != "" {
				separated.SharedContent += "\n\n"
			}
			separated.SharedContent += section
		}
	}
	
	return separated, nil
}

// classifySection classifies a content section as business, technical, or shared
func (ad *audienceDetector) classifySection(section string) string {
	sectionLower := strings.ToLower(section)
	
	// Business indicators
	businessIndicators := []string{
		"roi", "return on investment", "cost savings", "budget", "business case",
		"competitive advantage", "market", "revenue", "profit", "efficiency",
		"business value", "stakeholder", "executive", "investment", "funding",
		"business impact", "operational efficiency", "cost optimization",
	}
	
	// Technical indicators
	technicalIndicators := []string{
		"architecture", "infrastructure", "api", "microservices", "kubernetes",
		"docker", "devops", "ci/cd", "database", "performance", "scalability",
		"security", "integration", "deployment", "monitoring", "networking",
		"load balancer", "auto scaling", "serverless", "containers", "vpc",
	}
	
	// Recommendation indicators
	recommendationIndicators := []string{
		"recommend", "suggest", "propose", "should", "must", "need to",
		"next steps", "action items", "implementation", "approach",
	}
	
	businessScore := 0
	technicalScore := 0
	recommendationScore := 0
	
	// Count indicators
	for _, indicator := range businessIndicators {
		if strings.Contains(sectionLower, indicator) {
			businessScore++
		}
	}
	
	for _, indicator := range technicalIndicators {
		if strings.Contains(sectionLower, indicator) {
			technicalScore++
		}
	}
	
	for _, indicator := range recommendationIndicators {
		if strings.Contains(sectionLower, indicator) {
			recommendationScore++
		}
	}
	
	// Classify based on highest score
	if recommendationScore > 0 && (recommendationScore >= businessScore || recommendationScore >= technicalScore) {
		return "recommendation"
	} else if businessScore > technicalScore && businessScore > 0 {
		return "business"
	} else if technicalScore > businessScore && technicalScore > 0 {
		return "technical"
	}
	
	return "shared"
}

// initializeContentTemplates initializes default content templates
func (ad *audienceDetector) initializeContentTemplates() {
	// Technical audience template
	ad.templates["technical"] = &ContentTemplate{
		AudienceType:   AudienceTechnical,
		TechnicalDepth: 4,
		BusinessFocus:  2,
		SectionTemplate: `
TECHNICAL ANALYSIS:
- Detailed architecture and implementation specifications
- Performance, scalability, and reliability considerations
- Security configurations and best practices
- Integration patterns and API specifications
- Monitoring, logging, and operational considerations

IMPLEMENTATION DETAILS:
- Step-by-step technical implementation guide
- Configuration examples and code snippets
- Testing and validation procedures
- Deployment and rollback strategies
`,
		ToneGuidelines: []string{
			"Use technical terminology and industry jargon appropriately",
			"Provide specific implementation details and examples",
			"Focus on how rather than why",
			"Include performance metrics and benchmarks",
			"Reference technical documentation and standards",
		},
		ContentFocus: []string{
			"Architecture and design patterns",
			"Implementation specifics",
			"Performance and scalability",
			"Security and compliance",
			"Operational considerations",
		},
	}
	
	// Business audience template
	ad.templates["business"] = &ContentTemplate{
		AudienceType:   AudienceBusiness,
		TechnicalDepth: 2,
		BusinessFocus:  4,
		SectionTemplate: `
BUSINESS VALUE PROPOSITION:
- Return on investment and cost-benefit analysis
- Competitive advantages and market positioning
- Operational efficiency improvements
- Risk mitigation and business continuity
- Strategic alignment and growth enablement

FINANCIAL ANALYSIS:
- Investment requirements and budget implications
- Cost savings and revenue opportunities
- ROI projections and payback period
- Total cost of ownership considerations
`,
		ToneGuidelines: []string{
			"Focus on business value and outcomes",
			"Use business terminology and avoid technical jargon",
			"Emphasize ROI and cost-benefit analysis",
			"Highlight competitive advantages",
			"Address risk and compliance concerns",
		},
		ContentFocus: []string{
			"Business value and ROI",
			"Cost savings and efficiency",
			"Competitive positioning",
			"Risk management",
			"Strategic alignment",
		},
	}
	
	// Executive audience template
	ad.templates["executive"] = &ContentTemplate{
		AudienceType:   AudienceExecutive,
		TechnicalDepth: 1,
		BusinessFocus:  5,
		SectionTemplate: `
EXECUTIVE SUMMARY:
- Strategic implications and business transformation
- High-level investment and resource requirements
- Key risks and mitigation strategies
- Competitive advantages and market opportunities
- Timeline and critical success factors

STRATEGIC RECOMMENDATIONS:
- Long-term vision and roadmap alignment
- Organizational impact and change management
- Governance and oversight requirements
- Success metrics and KPIs
`,
		ToneGuidelines: []string{
			"Keep content high-level and strategic",
			"Focus on business transformation and competitive advantage",
			"Emphasize strategic alignment and vision",
			"Address governance and risk management",
			"Provide clear recommendations and next steps",
		},
		ContentFocus: []string{
			"Strategic transformation",
			"Competitive positioning",
			"Investment and ROI",
			"Risk and governance",
			"Organizational impact",
		},
	}
	
	// Mixed audience template
	ad.templates["mixed"] = &ContentTemplate{
		AudienceType:   AudienceMixed,
		TechnicalDepth: 3,
		BusinessFocus:  3,
		SectionTemplate: `
EXECUTIVE SUMMARY:
- High-level business value and strategic implications
- Key technical approach and architecture overview
- Investment requirements and expected ROI
- Risk assessment and mitigation strategies

BUSINESS JUSTIFICATION:
- Cost-benefit analysis and financial projections
- Competitive advantages and market positioning
- Operational efficiency and process improvements
- Compliance and risk management benefits

TECHNICAL APPROACH:
- Architecture overview and key design decisions
- Implementation approach and methodology
- Technology stack and integration considerations
- Security, performance, and scalability factors
`,
		ToneGuidelines: []string{
			"Balance business value with technical details",
			"Provide both strategic and tactical perspectives",
			"Use clear, accessible language for mixed audiences",
			"Separate business justification from technical implementation",
			"Include both high-level and detailed recommendations",
		},
		ContentFocus: []string{
			"Balanced business and technical content",
			"Strategic and tactical recommendations",
			"Clear separation of concerns",
			"Accessible language",
			"Comprehensive coverage",
		},
	}
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b, c float64) float64 {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}

func min(a, b, c float64) float64 {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}