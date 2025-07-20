package domain

// Service type constants
const (
	ServiceTypeAssessment        = "assessment"
	ServiceTypeMigration         = "migration"
	ServiceTypeOptimization      = "optimization"
	ServiceTypeArchitectureReview = "architecture_review"
)

// Inquiry status constants
const (
	InquiryStatusPending    = "pending"
	InquiryStatusProcessing = "processing"
	InquiryStatusReviewed   = "reviewed"
	InquiryStatusResponded  = "responded"
	InquiryStatusClosed     = "closed"
)

// Priority constants
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

// Default values
const (
	DefaultInquiryStatus = InquiryStatusPending
	DefaultPriority      = PriorityMedium
)

// Valid service types slice
var ValidServiceTypes = []string{
	ServiceTypeAssessment,
	ServiceTypeMigration,
	ServiceTypeOptimization,
	ServiceTypeArchitectureReview,
}

// Valid inquiry statuses slice
var ValidInquiryStatuses = []string{
	InquiryStatusPending,
	InquiryStatusProcessing,
	InquiryStatusReviewed,
	InquiryStatusResponded,
	InquiryStatusClosed,
}

// Valid priorities slice
var ValidPriorities = []string{
	PriorityLow,
	PriorityMedium,
	PriorityHigh,
	PriorityUrgent,
}

// Service descriptions for UI and documentation
var ServiceDescriptions = map[string]string{
	ServiceTypeAssessment: "Comprehensive evaluation of your current cloud infrastructure, security posture, and optimization opportunities",
	ServiceTypeMigration:  "Strategic planning and execution support for migrating workloads to the cloud",
	ServiceTypeOptimization: "Performance tuning, cost optimization, and efficiency improvements for existing cloud deployments",
	ServiceTypeArchitectureReview: "Expert review of cloud architecture designs for scalability, security, and best practices compliance",
}

// IsValidServiceType checks if a service type is valid
func IsValidServiceType(serviceType string) bool {
	for _, valid := range ValidServiceTypes {
		if serviceType == valid {
			return true
		}
	}
	return false
}

// IsValidInquiryStatus checks if an inquiry status is valid
func IsValidInquiryStatus(status string) bool {
	for _, valid := range ValidInquiryStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// IsValidPriority checks if a priority is valid
func IsValidPriority(priority string) bool {
	for _, valid := range ValidPriorities {
		if priority == valid {
			return true
		}
	}
	return false
}