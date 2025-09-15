# Design Document

## Overview

The polling-based chat system replaces the problematic WebSocket implementation with a reliable HTTP-based approach. Instead of maintaining persistent connections, the system uses regular HTTP requests to send messages and poll for updates. This eliminates connection lifecycle issues while providing a near real-time chat experience.

## Architecture

The system uses a simple client-server polling architecture:

1. **Frontend Polling Service** - Manages HTTP requests for sending/receiving messages
2. **Backend REST API** - Handles chat messages via standard HTTP endpoints
3. **Message Queue** - Stores messages temporarily for reliable delivery
4. **Polling Coordinator** - Manages polling intervals and connection state

## Components and Interfaces

### 1. Polling Chat Service

```typescript
interface PollingChatService {
  startPolling(): void;
  stopPolling(): void;
  sendMessage(message: ChatMessage): Promise<MessageResponse>;
  getMessages(sessionId: string, lastMessageId?: string): Promise<ChatMessage[]>;
  setPollingInterval(interval: number): void;
  getConnectionStatus(): ConnectionStatus;
}

interface ChatMessage {
  id: string;
  type: 'user' | 'assistant';
  content: string;
  timestamp: string;
  session_id: string;
  status: 'sending' | 'sent' | 'delivered' | 'failed';
}

interface MessageResponse {
  success: boolean;
  message_id: string;
  error?: string;
}
```

### 2. Backend REST Endpoints

```go
// POST /api/v1/admin/chat/messages - Send a new message
type SendMessageRequest struct {
    Content   string `json:"content"`
    SessionID string `json:"session_id"`
    ClientName string `json:"client_name,omitempty"`
}

type SendMessageResponse struct {
    Success   bool   `json:"success"`
    MessageID string `json:"message_id"`
    Error     string `json:"error,omitempty"`
}

// GET /api/v1/admin/chat/messages?session_id=X&since=Y - Poll for new messages
type GetMessagesResponse struct {
    Success  bool          `json:"success"`
    Messages []ChatMessage `json:"messages"`
    HasMore  bool          `json:"has_more"`
}
```

### 3. Polling Strategy

```typescript
interface PollingStrategy {
  baseInterval: number;        // 3 seconds default
  maxInterval: number;         // 30 seconds max
  backoffMultiplier: number;   // 1.5x on errors
  activeInterval: number;      // 2 seconds when actively chatting
  inactiveInterval: number;    // 10 seconds when idle
}
```

## Data Models

### Message Storage

```typescript
interface StoredMessage {
  id: string;
  session_id: string;
  type: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  user_id: string;
  metadata?: Record<string, any>;
  delivered: boolean;
  read: boolean;
}
```

### Session Management

```typescript
interface ChatSession {
  id: string;
  client_name?: string;
  created_at: Date;
  last_activity: Date;
  status: 'active' | 'inactive' | 'closed';
  message_count: number;
  last_message_id?: string;
}
```

## Error Handling

### 1. Network Error Handling

```typescript
interface ErrorHandler {
  handleNetworkError(error: NetworkError): void;
  handleServerError(error: ServerError): void;
  handleTimeoutError(error: TimeoutError): void;
  retryWithBackoff(operation: () => Promise<any>): Promise<any>;
}
```

### 2. Message Delivery Guarantees

- **Optimistic Updates**: Messages appear immediately with "sending" status
- **Retry Logic**: Failed messages retry with exponential backoff (1s, 2s, 4s, 8s)
- **Delivery Confirmation**: Server confirms message receipt with unique ID
- **Duplicate Prevention**: Client tracks sent messages to prevent duplicates

### 3. Connection State Management

```typescript
type ConnectionState = 'connected' | 'polling' | 'error' | 'offline';

interface ConnectionManager {
  state: ConnectionState;
  lastSuccessfulPoll: Date;
  errorCount: number;
  retryAttempts: number;
  isPollingActive: boolean;
}
```

## Testing Strategy

### 1. Polling Reliability Tests

```typescript
describe('Polling Chat Service', () => {
  test('should poll for messages at regular intervals');
  test('should handle network failures gracefully');
  test('should retry failed message sends');
  test('should adjust polling frequency based on activity');
  test('should prevent duplicate message delivery');
});
```

### 2. Backend API Tests

```go
func TestChatAPI(t *testing.T) {
    // Test message sending endpoint
    // Test message retrieval endpoint
    // Test session management
    // Test concurrent access
    // Test rate limiting
}
```

### 3. Integration Tests

- End-to-end message flow testing
- Multiple user concurrent chat testing
- Network interruption simulation
- Server restart recovery testing

## Implementation Approach

### Phase 1: Core Polling Infrastructure

1. **HTTP Chat Service**: Create polling-based service to replace WebSocket service
2. **Backend Endpoints**: Implement REST API for message send/receive
3. **Basic Polling**: Implement simple polling mechanism with fixed intervals

### Phase 2: Enhanced Reliability

1. **Smart Polling**: Implement adaptive polling intervals based on activity
2. **Error Handling**: Add comprehensive error handling and retry logic
3. **Message Queue**: Implement client-side message queuing for offline scenarios

### Phase 3: Performance Optimization

1. **Efficient Polling**: Optimize polling to reduce server load
2. **Caching**: Implement message caching to reduce redundant requests
3. **Compression**: Add response compression for large message payloads

### Phase 4: User Experience

1. **Status Indicators**: Show connection status and message delivery status
2. **Offline Support**: Handle offline scenarios gracefully
3. **Performance Metrics**: Add monitoring for polling performance

## Advantages Over WebSocket

1. **Reliability**: No connection lifecycle issues or unexpected disconnections
2. **Simplicity**: Standard HTTP requests are easier to debug and maintain
3. **Firewall Friendly**: Works through corporate firewalls and proxies
4. **Stateless**: Each request is independent, reducing server complexity
5. **Caching**: Can leverage HTTP caching mechanisms
6. **Load Balancing**: Works seamlessly with standard HTTP load balancers

## Performance Considerations

### Polling Efficiency

- **Smart Intervals**: Reduce polling when inactive (10s) vs active (2s)
- **Long Polling**: Consider long polling for better efficiency
- **Conditional Requests**: Use ETags or timestamps to avoid unnecessary data transfer
- **Batch Operations**: Send multiple messages in single request when possible

### Server Load

- **Rate Limiting**: Implement per-user rate limiting for polling requests
- **Caching**: Cache recent messages to reduce database queries
- **Connection Pooling**: Reuse HTTP connections for polling requests
- **Compression**: Use gzip compression for message payloads

## Migration Strategy

1. **Parallel Implementation**: Build polling system alongside existing WebSocket
2. **Feature Flag**: Use feature flag to switch between WebSocket and polling
3. **Gradual Rollout**: Test with subset of users before full deployment
4. **Fallback Mechanism**: Automatically fall back to polling if WebSocket fails
5. **Performance Monitoring**: Monitor both systems during transition period