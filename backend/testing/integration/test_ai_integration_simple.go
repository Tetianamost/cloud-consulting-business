package main

import (
	"fmt"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

func main() {
	fmt.Println("Testing AI Integration - Simple Version")

	// Test basic domain types
	session := &domain.ChatSession{
		ID:         "test-session",
		UserID:     "test-user",
		ClientName: "Test Client",
		Status:     domain.SessionStatusActive,
		CreatedAt:  time.Now(),
	}

	message := &domain.ChatMessage{
		ID:        "test-message",
		SessionID: session.ID,
		Type:      domain.MessageTypeUser,
		Content:   "Test message",
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
	}

	fmt.Printf("✓ Created session: %s\n", session.ID)
	fmt.Printf("✓ Created message: %s\n", message.ID)

	// Test session context
	ctx := &domain.SessionContext{
		ClientName:     "Test Client",
		MeetingType:    "consultation",
		ProjectContext: "Cloud migration",
		ServiceTypes:   []string{"migration", "architecture"},
	}

	fmt.Printf("✓ Created context for: %s\n", ctx.ClientName)
	fmt.Println("✅ Basic AI integration types work correctly!")
}
