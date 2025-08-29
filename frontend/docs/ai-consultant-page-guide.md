# AI Consultant Page User Guide

## Overview

The AI Consultant Page is an advanced, full-featured interface for interacting with the AI-powered consulting assistant. It provides a comprehensive set of tools for conducting professional consulting conversations with context awareness, quick actions, and advanced connection management.

## Key Features

### 1. Quick Actions System

The AI Consultant Page includes eight pre-defined quick actions for common consulting scenarios:

#### Available Quick Actions

| Action | Purpose | Example Use Case |
|--------|---------|------------------|
| **Cost Estimate** | Get pricing analysis for solutions | "How much would it cost to migrate 100 VMs to AWS?" |
| **Security Review** | Analyze security considerations | "What security measures are needed for HIPAA compliance?" |
| **Best Practices** | Get AWS best practice recommendations | "What are the best practices for RDS deployment?" |
| **Alternatives** | Explore alternative approaches | "What alternatives exist to EC2 for this workload?" |
| **Next Steps** | Get actionable next steps | "What should we do first in this migration?" |
| **Migration Plan** | Get migration strategy recommendations | "How should we approach migrating this legacy system?" |
| **Compliance** | Address compliance requirements | "What compliance considerations apply to financial data?" |
| **Performance** | Get performance optimization advice | "How can we optimize this architecture for performance?" |

#### Using Quick Actions

1. Click any quick action button to send a pre-defined prompt
2. The AI will respond with context-aware advice
3. Follow up with specific questions to dive deeper
4. Quick actions work with your configured client context

### 2. Context Management

#### Client Context Configuration

**Client Name**
- Set the client name for personalized responses
- AI will reference the client by name in responses
- Helps maintain professional, personalized communication

**Meeting Context**
- Specify the type of meeting or engagement
- Examples: "Migration planning", "Cost optimization", "Security assessment"
- Provides situational awareness to the AI assistant

#### Session Context

The AI Consultant Page automatically maintains:
- Conversation history within the session
- Client and meeting context across messages
- Session persistence across page reloads
- Message status and delivery confirmation

### 3. Connection Management

#### Connection Modes

The AI Consultant Page supports three connection modes:

**WebSocket Mode (Default)**
- Real-time, low-latency communication
- Instant message delivery
- Live typing indicators
- Best for active conversations

**Polling Mode**
- HTTP-based message polling
- Reliable fallback for network issues
- Works in restrictive network environments
- Automatic retry on failures

**Auto Mode**
- Automatically switches between WebSocket and polling
- Starts with WebSocket, falls back to polling if needed
- Provides best user experience with maximum reliability

#### Connection Status Indicators

- **Green Dot**: Connected and healthy
- **Red Dot**: Connection issues detected
- **Status Message**: Descriptive connection status
- **Health Check**: Real-time connection monitoring

### 4. User Interface Features

#### Fullscreen Mode

**Benefits**:
- Distraction-free conversation experience
- Larger message display area
- Better focus for complex discussions
- All features remain available

**Toggle**: Click the maximize/minimize button in the header

#### Settings Panel

**Access**: Click the settings gear icon in the header

**Configuration Options**:
- Client name input field
- Meeting context input field
- Connection mode selection (WebSocket/Polling/Auto)
- Real-time connection status display

#### Message Interface

**Message Display**:
- Clear distinction between user and AI messages
- Professional styling with proper spacing
- Timestamp display for all messages
- Message status indicators (sending, sent, delivered)

**Input Features**:
- Debounced input to prevent excessive API calls
- Visual feedback during message processing
- Disabled state during connection issues
- Enter key to send (Shift+Enter for new line)

**Loading States**:
- Animated typing indicators while AI is responding
- "AI is thinking..." message during processing
- Smooth message appearance animations

## Usage Patterns

### 1. Starting a Consultation Session

1. **Navigate** to the AI Consultant page from the admin dashboard
2. **Configure Context** (optional but recommended):
   - Enter client name in settings panel
   - Add meeting context (e.g., "AWS migration planning")
3. **Begin Conversation**:
   - Use a quick action for common scenarios
   - Or type a custom question
4. **Engage** in back-and-forth conversation
5. **Switch to Fullscreen** for focused discussions

### 2. Using Quick Actions Effectively

**For New Conversations**:
- Start with "Best Practices" to get foundational advice
- Follow with "Cost Estimate" for budget planning
- Use "Security Review" for compliance discussions

**For Ongoing Projects**:
- Use "Next Steps" to get actionable recommendations
- Use "Alternatives" to explore different approaches
- Use "Performance" for optimization discussions

**For Specific Scenarios**:
- Use "Migration Plan" for legacy system modernization
- Use "Compliance" for regulatory requirements
- Use "Cost Estimate" for budget planning

### 3. Managing Connection Issues

**If Connection Fails**:
1. Check the connection status indicator
2. Try switching to polling mode in settings
3. Use the refresh button to restart the connection
4. Check network connectivity if issues persist

**Best Practices**:
- Use auto mode for automatic failover
- Monitor connection status during important conversations
- Keep settings panel accessible for quick mode switching

## Advanced Features

### 1. Session Persistence

- Conversations are automatically saved
- Messages persist across page reloads
- Session context is maintained
- No data loss during connection issues

### 2. Error Handling

**Connection Errors**:
- Automatic retry for failed messages
- Visual error indicators
- Graceful degradation to polling mode
- Clear error messages with suggested actions

**Message Errors**:
- Retry buttons for failed messages
- Status indicators for message delivery
- Automatic resend on connection recovery

### 3. Performance Optimization

**Debounced Input**:
- Prevents excessive API calls during typing
- Visual feedback during debounce period
- Optimizes server resources

**Efficient Rendering**:
- Virtual scrolling for large message histories
- Optimized re-renders with React.memo
- Smooth animations without performance impact

## Integration with Other Systems

### 1. Redux State Management

The AI Consultant Page integrates with the global Redux store:

```typescript
// State structure
interface ChatState {
  currentSession: ChatSession | null;
  messages: ChatMessage[];
  isLoading: boolean;
  error: string | null;
}

interface ConnectionState {
  status: 'disconnected' | 'connecting' | 'connected' | 'error';
  mode: 'websocket' | 'polling';
  lastConnected: Date | null;
}
```

### 2. Service Integration

**Chat Mode Manager**:
- Handles connection mode switching
- Provides health monitoring
- Manages service initialization

**Polling Chat Service**:
- Reliable HTTP-based messaging
- Automatic retry logic
- Error handling and recovery

**WebSocket Service**:
- Real-time communication
- Connection health monitoring
- Automatic reconnection

## Troubleshooting

### Common Issues

#### Connection Problems

**Symptom**: Red connection indicator, messages not sending
**Solutions**:
1. Switch to polling mode in settings
2. Check network connectivity
3. Refresh the page
4. Contact system administrator if issues persist

#### Slow Response Times

**Symptom**: Long delays between messages
**Solutions**:
1. Check connection mode (WebSocket is faster)
2. Verify network speed
3. Monitor system load in dashboard
4. Consider switching to a less busy time

#### Context Not Working

**Symptom**: AI doesn't reference client name or context
**Solutions**:
1. Verify client name is set in settings
2. Check meeting context configuration
3. Start a new session if context was added mid-conversation
4. Ensure settings are saved properly

### Debug Information

**Browser Console**:
- Check for JavaScript errors
- Monitor network requests
- Review WebSocket connection status

**Network Tab**:
- Verify API requests are successful
- Check response times
- Monitor WebSocket connections

## Best Practices

### 1. Professional Usage

- Always set client context for personalized responses
- Use appropriate meeting context for situational awareness
- Start with quick actions for common scenarios
- Follow up with specific questions for detailed advice

### 2. Performance

- Use fullscreen mode for long conversations
- Monitor connection status during important discussions
- Switch to polling mode in unreliable network conditions
- Clear old sessions periodically for optimal performance

### 3. Security

- Don't include sensitive information in messages
- Use secure connections (HTTPS/WSS)
- Log out properly when finished
- Follow organizational data handling policies

## API Integration

### Backend Endpoints

The AI Consultant Page integrates with the following backend endpoints:

- `POST /api/v1/admin/simple-chat/messages` - Send messages to AI assistant
- `GET /api/v1/admin/simple-chat/messages` - Retrieve chat history
- `GET /health` - Health check for connection monitoring

### Service Architecture

```typescript
// Service layer architecture
SimpleAIService
├── Connection Management
├── Session Management  
├── Message Handling
├── Error Recovery
└── Health Monitoring
```

### Authentication Flow

1. JWT token retrieved from localStorage
2. Token included in all API requests
3. Automatic token refresh handling
4. Graceful degradation on auth failures

## Technical Implementation

### Component Architecture

```typescript
AIConsultantPage
├── Context Management (Client name, meeting type)
├── Connection Status (WebSocket/Polling indicators)
├── Quick Actions (8 pre-defined actions)
├── Message Interface (Input, display, history)
├── Settings Panel (Configuration options)
└── Fullscreen Mode (Distraction-free interface)
```

### State Management

The component uses Redux for state management:

```typescript
// Chat state structure
interface ChatState {
  currentSession: ChatSession | null;
  messages: ChatMessage[];
  isLoading: boolean;
  error: string | null;
}

// Connection state structure  
interface ConnectionState {
  status: 'connected' | 'disconnected' | 'error';
  mode: 'websocket' | 'polling';
  lastConnected: Date | null;
}
```

### Performance Optimizations

- **Debounced Input**: 300ms delay to prevent excessive API calls
- **Memoized Messages**: React.useMemo for message list optimization
- **Efficient Rendering**: Minimal re-renders with proper dependency arrays
- **Auto-scroll**: Smooth scrolling to new messages

## Future Enhancements

### Planned Features

- **File Upload**: Share documents and diagrams with AI
- **Voice Input**: Speech-to-text for hands-free operation
- **Message Export**: Save conversations as PDF or text
- **Templates**: Custom quick action templates
- **Multi-language**: Support for multiple languages
- **Collaboration**: Multi-user chat sessions
- **Message Search**: Full-text search across conversation history
- **Conversation Branching**: Create alternative conversation paths

### Performance Improvements

- **Intelligent Caching**: Response caching for similar queries
- **Message Compression**: Reduce bandwidth usage
- **Offline Mode**: Limited functionality without connection
- **Mobile App**: Dedicated mobile application
- **Virtual Scrolling**: Handle large message histories efficiently
- **Background Sync**: Sync messages when connection restored

### Integration Enhancements

- **CRM Integration**: Pull client data from external systems
- **Calendar Integration**: Schedule follow-up meetings
- **Document Management**: Integration with document repositories
- **Analytics Dashboard**: Conversation analytics and insights
- **Webhook Support**: Real-time notifications to external systems

This AI Consultant Page provides a professional, feature-rich interface for conducting AI-powered consulting conversations with advanced context management, reliable connectivity, and an intuitive user experience optimized for enterprise consulting workflows.