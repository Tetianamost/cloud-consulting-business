---
inclusion: fileMatch
fileMatchPattern: '*chat*'
---

# Chat System Development Patterns

## Overview

This document provides specific guidance for developing chat-related features in the Cloud Consulting Platform, including WebSocket and polling-based implementations, AI integration patterns, and real-time communication best practices.

## Chat Architecture Patterns

### Dual Communication Strategy
The platform supports both WebSocket and polling-based chat systems to ensure reliability:

1. **WebSocket (Primary)**: For real-time, low-latency communication
2. **Polling (Fallback)**: For environments with WebSocket connectivity issues

### Implementation Pattern
```go
// Backend: Unified chat service interface
type ChatService interface {
    SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
    GetMessages(ctx context.Context, sessionID string, since time.Time) ([]*ChatMessage, error)
    CreateSession(ctx context.Context, userID string) (*ChatSession, error)
}

// Support both WebSocket and HTTP handlers
func (h *ChatHandler) HandleWebSocket(c *gin.Context) { /* WebSocket logic */ }
func (h *ChatHandler) HandlePolling(c *gin.Context) { /* HTTP polling logic */ }
```

## Session Management Patterns

### Session Lifecycle
1. **Creation**: Generate secure session ID with user authentication
2. **Validation**: Verify session on each request/connection
3. **Persistence**: Store session state in Redis with TTL
4. **Cleanup**: Automatic session expiration and cleanup

```go
type ChatSession struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    ExpiresAt time.Time `json:"expires_at"`
    IsActive  bool      `json:"is_active"`
}
```

### Session Security
- Use cryptographically secure session IDs
- Implement session timeout (30 minutes default)
- Validate session on every request
- Log session creation/destruction for audit

## Message Handling Patterns

### Message Structure
```go
type ChatMessage struct {
    ID        string    `json:"id"`
    SessionID string    `json:"session_id"`
    Type      string    `json:"type"` // "user", "ai", "system"
    Content   string    `json:"content"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    Status    string    `json:"status"` // "sending", "sent", "delivered", "failed"
}
```

### Message Processing Pipeline
1. **Validation**: Sanitize and validate message content
2. **Persistence**: Store message immediately
3. **AI Processing**: Send to Bedrock for AI response (async)
4. **Response Handling**: Process AI response and store
5. **Delivery**: Send to client via WebSocket or polling

## AI Integration Patterns

### Bedrock Integration
```go
type BedrockService struct {
    client *bedrock.Client
    config *BedrockConfig
}

func (s *BedrockService) GenerateResponse(ctx context.Context, prompt string) (*AIResponse, error) {
    // Implement with proper error handling and retries
    // Include context from previous messages
    // Use structured prompts for consistency
}
```

### Prompt Engineering Standards
- Include conversation context (last 5-10 messages)
- Use system prompts to define AI role and behavior
- Implement prompt templates for different scenarios
- Version control prompts for A/B testing

### AI Response Handling
- Stream responses when possible for better UX
- Implement proper error handling and fallbacks
- Cache responses for similar queries
- Monitor AI service usage and costs

## Real-time Communication Patterns

### WebSocket Implementation
```go
type WebSocketManager struct {
    connections map[string]*websocket.Conn
    broadcast   chan []byte
    register    chan *Client
    unregister  chan *Client
}

func (m *WebSocketManager) HandleConnection(conn *websocket.Conn, sessionID string) {
    // Implement connection lifecycle management
    // Handle ping/pong for connection health
    // Implement graceful disconnection
}
```

### Polling Implementation
```go
type PollingHandler struct {
    chatService ChatService
    pollInterval time.Duration
}

func (h *PollingHandler) GetNewMessages(c *gin.Context) {
    // Implement efficient polling with since timestamp
    // Return empty response for no new messages
    // Include proper caching headers
}
```

## Frontend Chat Patterns

### React Chat Components
```tsx
// Main chat container
export const ChatContainer: React.FC = () => {
  const { messages, sendMessage, connectionStatus } = useChat();
  
  return (
    <div className="chat-container">
      <MessageList messages={messages} />
      <MessageInput onSend={sendMessage} disabled={connectionStatus !== 'connected'} />
      <ConnectionStatus status={connectionStatus} />
    </div>
  );
};
```

### Connection Management
```tsx
// Custom hook for chat connection
export const useChat = () => {
  const [connectionType, setConnectionType] = useState<'websocket' | 'polling'>('websocket');
  
  // Implement fallback logic
  useEffect(() => {
    if (websocketFailed) {
      setConnectionType('polling');
    }
  }, [websocketFailed]);
  
  return { messages, sendMessage, connectionStatus };
};
```

### State Management
```tsx
// Redux slice for chat state
const chatSlice = createSlice({
  name: 'chat',
  initialState: {
    messages: [],
    sessions: {},
    connectionStatus: 'disconnected',
    currentSessionId: null,
  },
  reducers: {
    messageReceived: (state, action) => {
      state.messages.push(action.payload);
    },
    messageSent: (state, action) => {
      // Optimistic update
      state.messages.push({ ...action.payload, status: 'sending' });
    },
  },
});
```

## Error Handling Patterns

### Backend Error Handling
```go
func (h *ChatHandler) SendMessage(c *gin.Context) {
    // Validate input
    if err := validateMessage(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
        return
    }
    
    // Handle service errors gracefully
    if err := h.chatService.SendMessage(ctx, req); err != nil {
        log.WithError(err).Error("Failed to send message")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Message delivery failed"})
        return
    }
}
```

### Frontend Error Handling
```tsx
const MessageInput: React.FC<Props> = ({ onSend }) => {
  const [error, setError] = useState<string | null>(null);
  
  const handleSend = async (message: string) => {
    try {
      await onSend(message);
      setError(null);
    } catch (err) {
      setError('Failed to send message. Please try again.');
      // Implement retry logic
    }
  };
  
  return (
    <div>
      {error && <ErrorMessage message={error} onRetry={() => handleSend(lastMessage)} />}
      {/* Input component */}
    </div>
  );
};
```

## Performance Optimization Patterns

### Message Pagination
- Implement virtual scrolling for large message lists
- Load messages in chunks (50-100 messages)
- Use infinite scrolling for chat history
- Cache recent messages in local storage

### Connection Optimization
- Implement connection pooling for WebSockets
- Use efficient polling intervals (3-5 seconds)
- Implement exponential backoff for reconnection
- Reduce polling frequency when user is inactive

### Caching Strategies
- Cache AI responses for similar queries
- Store recent chat sessions in Redis
- Implement client-side message caching
- Use CDN for static chat assets

## Testing Patterns

### Backend Testing
```go
func TestChatHandler_SendMessage(t *testing.T) {
    // Setup
    mockService := &MockChatService{}
    handler := NewChatHandler(mockService)
    
    // Test cases
    tests := []struct {
        name           string
        request        *SendMessageRequest
        expectedStatus int
        expectedError  string
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Frontend Testing
```tsx
describe('ChatContainer', () => {
  it('should send message when user submits', async () => {
    const mockSendMessage = jest.fn();
    render(<ChatContainer onSendMessage={mockSendMessage} />);
    
    const input = screen.getByPlaceholderText('Type a message...');
    const sendButton = screen.getByRole('button', { name: 'Send' });
    
    fireEvent.change(input, { target: { value: 'Hello' } });
    fireEvent.click(sendButton);
    
    expect(mockSendMessage).toHaveBeenCalledWith('Hello');
  });
});
```

## Security Considerations

### Message Validation
- Sanitize all user input
- Implement message length limits
- Validate message types and formats
- Filter potentially harmful content

### Authentication Integration
- Validate JWT tokens on WebSocket connections
- Implement session-based authentication for polling
- Use secure session storage
- Implement proper logout handling

### Data Privacy
- Encrypt sensitive message content
- Implement message retention policies
- Provide user data deletion capabilities
- Comply with privacy regulations (GDPR, CCPA)

## Monitoring and Observability

### Metrics to Track
- Message delivery success rate
- WebSocket connection stability
- AI response times and success rates
- User engagement metrics
- Error rates by component

### Logging Standards
```go
log.WithFields(log.Fields{
    "session_id": sessionID,
    "user_id": userID,
    "message_type": messageType,
    "ai_model": "bedrock-nova",
    "response_time_ms": responseTime,
}).Info("Message processed successfully")
```

### Alerting
- Set up alerts for high error rates
- Monitor AI service availability
- Track WebSocket connection failures
- Alert on unusual message patterns