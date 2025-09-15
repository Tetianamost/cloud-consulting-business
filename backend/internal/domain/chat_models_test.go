package domain

import (
	"testing"
	"time"
)

func TestChatSession_Validate(t *testing.T) {
	tests := []struct {
		name    string
		session ChatSession
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid session",
			session: ChatSession{
				ID:           "test-id",
				UserID:       "user123",
				ClientName:   "Test Client",
				Context:      "Test context",
				Status:       SessionStatusActive,
				Metadata:     map[string]interface{}{"key": "value"},
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				LastActivity: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			session: ChatSession{
				ID:     "test-id",
				Status: SessionStatusActive,
			},
			wantErr: true,
			errMsg:  "User ID is required",
		},
		{
			name: "user_id too long",
			session: ChatSession{
				ID:     "test-id",
				UserID: string(make([]byte, 101)), // 101 characters
				Status: SessionStatusActive,
			},
			wantErr: true,
			errMsg:  "User ID must be 100 characters or less",
		},
		{
			name: "client_name too long",
			session: ChatSession{
				ID:         "test-id",
				UserID:     "user123",
				ClientName: string(make([]byte, 256)), // 256 characters
				Status:     SessionStatusActive,
			},
			wantErr: true,
			errMsg:  "Client name must be 255 characters or less",
		},
		{
			name: "missing status",
			session: ChatSession{
				ID:     "test-id",
				UserID: "user123",
			},
			wantErr: true,
			errMsg:  "Status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.session.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatSession.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ChatSession.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestChatMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message ChatMessage
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid message",
			message: ChatMessage{
				ID:        "msg-id",
				SessionID: "session-id",
				Type:      MessageTypeUser,
				Content:   "Hello, world!",
				Status:    MessageStatusSent,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing session_id",
			message: ChatMessage{
				ID:      "msg-id",
				Type:    MessageTypeUser,
				Content: "Hello, world!",
				Status:  MessageStatusSent,
			},
			wantErr: true,
			errMsg:  "Session ID is required",
		},
		{
			name: "missing content",
			message: ChatMessage{
				ID:        "msg-id",
				SessionID: "session-id",
				Type:      MessageTypeUser,
				Status:    MessageStatusSent,
			},
			wantErr: true,
			errMsg:  "Content is required",
		},
		{
			name: "content too long",
			message: ChatMessage{
				ID:        "msg-id",
				SessionID: "session-id",
				Type:      MessageTypeUser,
				Content:   string(make([]byte, 10001)), // 10001 characters
				Status:    MessageStatusSent,
			},
			wantErr: true,
			errMsg:  "Content must be 10000 characters or less",
		},
		{
			name: "missing type",
			message: ChatMessage{
				ID:        "msg-id",
				SessionID: "session-id",
				Content:   "Hello, world!",
				Status:    MessageStatusSent,
			},
			wantErr: true,
			errMsg:  "Message type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatMessage.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ChatMessage.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestChatSession_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		session ChatSession
		want    bool
	}{
		{
			name: "no expiration set",
			session: ChatSession{
				ExpiresAt: nil,
			},
			want: false,
		},
		{
			name: "not expired",
			session: ChatSession{
				ExpiresAt: &[]time.Time{now.Add(time.Hour)}[0],
			},
			want: false,
		},
		{
			name: "expired",
			session: ChatSession{
				ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.IsExpired(); got != tt.want {
				t.Errorf("ChatSession.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatSession_IsActive(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		session ChatSession
		want    bool
	}{
		{
			name: "active and not expired",
			session: ChatSession{
				Status:    SessionStatusActive,
				ExpiresAt: &[]time.Time{now.Add(time.Hour)}[0],
			},
			want: true,
		},
		{
			name: "active but expired",
			session: ChatSession{
				Status:    SessionStatusActive,
				ExpiresAt: &[]time.Time{now.Add(-time.Hour)}[0],
			},
			want: false,
		},
		{
			name: "inactive",
			session: ChatSession{
				Status:    SessionStatusInactive,
				ExpiresAt: &[]time.Time{now.Add(time.Hour)}[0],
			},
			want: false,
		},
		{
			name: "active with no expiration",
			session: ChatSession{
				Status:    SessionStatusActive,
				ExpiresAt: nil,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.IsActive(); got != tt.want {
				t.Errorf("ChatSession.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatSession_SetExpiration(t *testing.T) {
	session := &ChatSession{}
	duration := time.Hour * 24

	session.SetExpiration(duration)

	if session.ExpiresAt == nil {
		t.Error("ChatSession.SetExpiration() did not set ExpiresAt")
		return
	}

	expected := time.Now().Add(duration)
	if session.ExpiresAt.Sub(expected) > time.Second {
		t.Errorf("ChatSession.SetExpiration() set incorrect expiration time")
	}
}

func TestChatSession_UpdateActivity(t *testing.T) {
	session := &ChatSession{
		LastActivity: time.Now().Add(-time.Hour),
		UpdatedAt:    time.Now().Add(-time.Hour),
	}

	oldActivity := session.LastActivity
	oldUpdated := session.UpdatedAt

	session.UpdateActivity()

	if !session.LastActivity.After(oldActivity) {
		t.Error("ChatSession.UpdateActivity() did not update LastActivity")
	}

	if !session.UpdatedAt.After(oldUpdated) {
		t.Error("ChatSession.UpdateActivity() did not update UpdatedAt")
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("test_field", "test message")

	if err.Field != "test_field" {
		t.Errorf("ValidationError.Field = %v, want %v", err.Field, "test_field")
	}

	if err.Message != "test message" {
		t.Errorf("ValidationError.Message = %v, want %v", err.Message, "test message")
	}

	if err.Error() != "test message" {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), "test message")
	}
}

func TestSessionStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   SessionStatus
		expected string
	}{
		{"active", SessionStatusActive, "active"},
		{"inactive", SessionStatusInactive, "inactive"},
		{"expired", SessionStatusExpired, "expired"},
		{"terminated", SessionStatusTerminated, "terminated"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("SessionStatus constant %v = %v, want %v", tt.name, string(tt.status), tt.expected)
			}
		})
	}
}

func TestMessageType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		msgType  MessageType
		expected string
	}{
		{"user", MessageTypeUser, "user"},
		{"assistant", MessageTypeAssistant, "assistant"},
		{"system", MessageTypeSystem, "system"},
		{"error", MessageTypeError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.msgType) != tt.expected {
				t.Errorf("MessageType constant %v = %v, want %v", tt.name, string(tt.msgType), tt.expected)
			}
		})
	}
}

func TestMessageStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   MessageStatus
		expected string
	}{
		{"sent", MessageStatusSent, "sent"},
		{"delivered", MessageStatusDelivered, "delivered"},
		{"read", MessageStatusRead, "read"},
		{"failed", MessageStatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("MessageStatus constant %v = %v, want %v", tt.name, string(tt.status), tt.expected)
			}
		})
	}
}
