package domain

// Service type constants
const (
	ServiceTypeAssessment         = "assessment"
	ServiceTypeMigration          = "migration"
	ServiceTypeOptimization       = "optimization"
	ServiceTypeArchitectureReview = "architecture_review"
	ServiceTypeGeneral            = "general"
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
	ServiceTypeGeneral,
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
	ServiceTypeAssessment:         "Comprehensive evaluation of your current cloud infrastructure, security posture, and optimization opportunities",
	ServiceTypeMigration:          "Strategic planning and execution support for migrating workloads to the cloud",
	ServiceTypeOptimization:       "Performance tuning, cost optimization, and efficiency improvements for existing cloud deployments",
	ServiceTypeArchitectureReview: "Expert review of cloud architecture designs for scalability, security, and best practices compliance",
	ServiceTypeGeneral:            "General cloud consulting and advisory services",
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

// Chat-related constants

// Chat session status constants
const (
	ChatSessionStatusActive     = "active"
	ChatSessionStatusInactive   = "inactive"
	ChatSessionStatusExpired    = "expired"
	ChatSessionStatusTerminated = "terminated"
)

// Chat message type constants
const (
	ChatMessageTypeUser      = "user"
	ChatMessageTypeAssistant = "assistant"
	ChatMessageTypeSystem    = "system"
	ChatMessageTypeError     = "error"
)

// Chat message status constants
const (
	ChatMessageStatusSent      = "sent"
	ChatMessageStatusDelivered = "delivered"
	ChatMessageStatusRead      = "read"
	ChatMessageStatusFailed    = "failed"
)

// Chat session configuration constants
const (
	DefaultSessionExpirationHours = 24
	MaxSessionsPerUser            = 10
	MaxMessageLength              = 10000
	MaxMessagesPerSession         = 1000
	DefaultMessagePageSize        = 50
	MaxMessagePageSize            = 1000
)

// Valid chat session statuses slice
var ValidChatSessionStatuses = []string{
	ChatSessionStatusActive,
	ChatSessionStatusInactive,
	ChatSessionStatusExpired,
	ChatSessionStatusTerminated,
}

// Valid chat message types slice
var ValidChatMessageTypes = []string{
	ChatMessageTypeUser,
	ChatMessageTypeAssistant,
	ChatMessageTypeSystem,
	ChatMessageTypeError,
}

// Valid chat message statuses slice
var ValidChatMessageStatuses = []string{
	ChatMessageStatusSent,
	ChatMessageStatusDelivered,
	ChatMessageStatusRead,
	ChatMessageStatusFailed,
}

// Chat quick action constants
const (
	QuickActionCostEstimate     = "cost_estimate"
	QuickActionArchitectureHelp = "architecture_help"
	QuickActionBestPractices    = "best_practices"
	QuickActionTroubleshooting  = "troubleshooting"
	QuickActionServiceInfo      = "service_info"
)

// Valid quick actions slice
var ValidQuickActions = []string{
	QuickActionCostEstimate,
	QuickActionArchitectureHelp,
	QuickActionBestPractices,
	QuickActionTroubleshooting,
	QuickActionServiceInfo,
}

// Quick action descriptions
var QuickActionDescriptions = map[string]string{
	QuickActionCostEstimate:     "Get quick cost estimates for AWS services and architectures",
	QuickActionArchitectureHelp: "Get help with cloud architecture design and best practices",
	QuickActionBestPractices:    "Learn about AWS best practices for security, performance, and cost optimization",
	QuickActionTroubleshooting:  "Get help troubleshooting common AWS issues and problems",
	QuickActionServiceInfo:      "Get information about specific AWS services and their capabilities",
}

// Validation functions for chat-related constants

// IsValidChatSessionStatus checks if a chat session status is valid
func IsValidChatSessionStatus(status string) bool {
	for _, valid := range ValidChatSessionStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// IsValidChatMessageType checks if a chat message type is valid
func IsValidChatMessageType(messageType string) bool {
	for _, valid := range ValidChatMessageTypes {
		if messageType == valid {
			return true
		}
	}
	return false
}

// IsValidChatMessageStatus checks if a chat message status is valid
func IsValidChatMessageStatus(status string) bool {
	for _, valid := range ValidChatMessageStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// IsValidQuickAction checks if a quick action is valid
func IsValidQuickAction(action string) bool {
	for _, valid := range ValidQuickActions {
		if action == valid {
			return true
		}
	}
	return false
}
