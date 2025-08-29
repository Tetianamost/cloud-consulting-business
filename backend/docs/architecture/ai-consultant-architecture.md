# AI Consultant Architecture Documentation

## Overview

The AI Consultant system provides an advanced, context-aware chat interface for cloud consulting conversations. It combines a React-based frontend with a Go backend, featuring real-time communication, session management, and intelligent response generation.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend Layer                           │
├─────────────────────────────────────────────────────────────────┤
│  AIConsultantPage.tsx                                          │
│  ├── Context Management (Client, Meeting Type)                 │
│  ├── Quick Actions (8 pre-defined consulting actions)          │
│  ├── Connection Management (WebSocket/Polling)                 │
│  ├── Message Interface (Input, Display, History)               │
│  └── Settings Panel (Configuration, Status)                    │
│                                                                 │
│  SimpleAIService.ts                                            │
│  ├── Session Management                                         │
│  ├── Message Handling                                           │
│  ├── Connection Health Monitoring                               │
│  └── Error Recovery                                             │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ HTTPS/WSS
                                │
┌─────────────────────────────────────────────────────────────────┐
│                        Backend Layer                            │
├─────────────────────────────────────────────────────────────────┤
│  SimpleChatHandler.go                                          │
│  ├── Message Processing                                         │
│  ├── Session Management                                         │
│  ├── AI Response Generation                                     │
│  └── In-Memory Storage                                          │
│                                                                 │
│  Authentication Middleware                                      │
│  ├── JWT Token Validation                                       │
│  ├── Admin Role Verification                                    │
│  └── Request Authorization                                       │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ Future Integration
                                │
┌─────────────────────────────────────────────────────────────────┐
│                      AI Services Layer                          │
├─────────────────────────────────────────────────────────────────┤
│  AWS Bedrock Integration (Planned)                             │
│  ├── Nova Model Integration                                     │
│  ├── Context-Aware Prompts                                     │
│  ├── Response Optimization                                      │
│  └── Cost Management                                            │
└─────────────────────────────────────────────────────────────────┘
```

## Component Architecture

### Frontend Components

#### AIConsultantPage Component

**Purpose**: Main interface for AI consulting conversations

**Key Features**:
- Context-aware conversations with client name and meeting type
- 8 pre-defined quick actions for common consulting scenarios
- Fullscreen mode for focused conversations
- Real-time connection status monitoring
- Debounced input for performance optimization

**State Management**:
```typescript
interface ComponentState {
  clientName: string;
  meetingContext: string;
  showSettings: boolean;
  isFullscreen: boolean;
}

// Redux Integration
const chatState = useSelector((state: RootState) => state.chat);
const connectionState = useSelector((state: RootState) => state.connection);
```

**Quick Actions System**:
```typescript
const QUICK_ACTIONS = [
  { id: 'cost-estimate', label: 'Cost Estimate', prompt: '...' },
  { id: 'security-review', label: 'Security Review', prompt: '...' },
  { id: 'best-practices', label: 'Best Practices', prompt: '...' },
  { id: 'alternatives', label: 'Alternatives', prompt: '...' },
  { id: 'next-steps', label: 'Next Steps', prompt: '...' },
  { id: 'migration-plan', label: 'Migration Plan', prompt: '...' },
  { id: 'compliance', label: 'Compliance', prompt: '...' },
  { id: 'performance', label: 'Performance', prompt: '...' },
];
```

#### SimpleAIService

**Purpose**: Service layer for AI chat functionality

**Key Responsibilities**:
- Session management with unique session IDs
- Message sending and retrieval
- Connection health monitoring
- Error handling and recovery

**Service Methods**:
```typescript
class SimpleAIService {
  async sendMessage(request: SimpleAIRequest): Promise<SimpleAIResponse>
  async getMessages(): Promise<SimpleChatMessage[]>
  async checkConnection(): Promise<boolean>
  isHealthy(): boolean
  getSessionId(): string
  resetSession(): void
}
```

### Backend Components

#### SimpleChatHandler

**Purpose**: HTTP handler for simple chat functionality

**Key Features**:
- In-memory message storage
- Immediate AI response generation
- Session-based message filtering
- RESTful API endpoints

**Handler Methods**:
```go
type SimpleChatHandler struct {
    logger   *logrus.Logger
    messages []SimpleChatMessage
}

func (h *SimpleChatHandler) SendMessage(c *gin.Context)
func (h *SimpleChatHandler) GetMessages(c *gin.Context)
func (h *SimpleChatHandler) generateSimpleAIResponse(userMessage string) string
```

**Message Structure**:
```go
type SimpleChatMessage struct {
    ID        string    `json:"id"`
    Content   string    `json:"content"`
    Role      string    `json:"role"`
    Timestamp time.Time `json:"timestamp"`
    SessionID string    `json:"session_id"`
}
```

## Data Flow

### Message Send Flow

```
1. User Input
   ├── Frontend: AIConsultantPage
   ├── Debounced Input Processing
   └── Redux State Update (Optimistic)

2. API Request
   ├── SimpleAIService.sendMessage()
   ├── JWT Authentication
   └── HTTP POST to /api/v1/admin/simple-chat/messages

3. Backend Processing
   ├── SimpleChatHandler.SendMessage()
   ├── Message Validation
   ├── User Message Storage
   ├── AI Response Generation
   └── AI Message Storage

4. Response Handling
   ├── Success Response to Frontend
   ├── Message History Retrieval
   ├── Redux State Update
   └── UI Update with New Messages
```

### Session Management Flow

```
1. Session Creation
   ├── Frontend: Generate unique session ID
   ├── Format: session-{timestamp}-{random}
   └── Store in SimpleAIService instance

2. Session Persistence
   ├── Session ID included in all requests
   ├── Backend filters messages by session
   └── Frontend maintains session across page reloads

3. Session Reset
   ├── User-initiated via clear chat button
   ├── Generate new session ID
   └── Clear local message history
```

## API Design

### Endpoints

#### Send Message
```http
POST /api/v1/admin/simple-chat/messages
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "content": "How much would it cost to migrate to AWS?",
  "session_id": "session-1234567890-abc123"
}
```

#### Get Messages
```http
GET /api/v1/admin/simple-chat/messages?session_id=session-1234567890-abc123
Authorization: Bearer <jwt_token>
```

### Response Format

**Success Response**:
```json
{
  "success": true,
  "message_id": "msg-1234567890-user"
}
```

**Message List Response**:
```json
{
  "success": true,
  "messages": [
    {
      "id": "msg-1234567890-user",
      "content": "User message content",
      "role": "user",
      "timestamp": "2024-01-15T10:30:00Z",
      "session_id": "session-1234567890-abc123"
    }
  ]
}
```

## Security Architecture

### Authentication Flow

```
1. Admin Login
   ├── Username/Password Validation
   ├── JWT Token Generation
   └── Token Storage in localStorage

2. API Request Authentication
   ├── JWT Token from localStorage
   ├── Authorization Header
   └── Backend Token Validation

3. Session Security
   ├── Session IDs are not predictable
   ├── Messages isolated by session
   └── No cross-session data leakage
```

### Security Measures

- **JWT Authentication**: All API requests require valid JWT tokens
- **HTTPS/TLS**: All communications encrypted in transit
- **Session Isolation**: Messages are isolated by session ID
- **Input Validation**: All user inputs are validated and sanitized
- **Admin-Only Access**: Only admin users can access AI consultant features

## Performance Considerations

### Frontend Optimizations

- **Debounced Input**: 300ms delay prevents excessive API calls
- **Memoized Components**: React.useMemo for expensive computations
- **Efficient Rendering**: Minimal re-renders with proper dependencies
- **Auto-scroll Optimization**: Smooth scrolling without performance impact

### Backend Optimizations

- **In-Memory Storage**: Fast message retrieval for demo purposes
- **Simple Response Generation**: Immediate responses without external API calls
- **Efficient Filtering**: O(n) message filtering by session ID
- **Minimal Processing**: Lightweight request handling

### Scalability Considerations

**Current Limitations**:
- In-memory storage limits scalability
- No message persistence across server restarts
- Single-server architecture

**Future Improvements**:
- Database storage for message persistence
- Redis caching for session management
- Load balancing for multiple server instances
- Message queuing for high-volume scenarios

## Error Handling

### Frontend Error Handling

```typescript
// Service-level error handling
try {
  const response = await simpleAIService.sendMessage(request);
  dispatch(addMessage(response));
} catch (error) {
  dispatch(setError('Failed to send message'));
  // Add error message to UI
  dispatch(addMessage(errorMessage));
}
```

### Backend Error Handling

```go
// Handler-level error handling
if err := c.ShouldBindJSON(&req); err != nil {
    h.logger.WithError(err).Error("Failed to bind request")
    c.JSON(http.StatusBadRequest, SimpleSendMessageResponse{
        Success: false,
        Error:   "Invalid request format",
    })
    return
}
```

### Error Recovery Strategies

- **Automatic Retry**: Failed requests are automatically retried
- **Graceful Degradation**: System continues to function with reduced features
- **User Feedback**: Clear error messages with suggested actions
- **Logging**: Comprehensive error logging for debugging

## Monitoring and Observability

### Logging Strategy

**Frontend Logging**:
- Connection status changes
- API request/response cycles
- User interaction events
- Error conditions

**Backend Logging**:
- Message send/receive operations
- Session management events
- Authentication attempts
- Performance metrics

### Metrics Collection

**Key Metrics**:
- Messages per session
- Average response time
- Error rates by endpoint
- Active session count
- User engagement metrics

## Future Architecture Enhancements

### Planned Improvements

#### AI Integration
- **AWS Bedrock Integration**: Replace simple responses with AI-generated content
- **Context-Aware Prompts**: Use client and meeting context for better responses
- **Response Caching**: Cache similar responses for performance
- **Cost Optimization**: Monitor and optimize AI API usage

#### Data Persistence
- **Database Integration**: PostgreSQL for message persistence
- **Redis Caching**: Session and response caching
- **Message History**: Long-term conversation storage
- **Analytics Database**: Separate analytics data store

#### Real-time Features
- **WebSocket Integration**: Real-time message delivery
- **Typing Indicators**: Show when AI is generating responses
- **Presence Management**: Online/offline status
- **Push Notifications**: Browser notifications for new messages

#### Scalability Enhancements
- **Microservices Architecture**: Separate services for different concerns
- **Load Balancing**: Multiple server instances
- **Message Queuing**: Asynchronous message processing
- **CDN Integration**: Static asset delivery optimization

This architecture provides a solid foundation for AI-powered consulting conversations while maintaining simplicity and performance. The modular design allows for incremental enhancements and scaling as requirements evolve.