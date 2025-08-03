# AI Consultant Live Chat API Documentation

## Overview

The AI Consultant Live Chat API provides real-time chat capabilities for admin users to interact with an AI assistant powered by AWS Bedrock. The system supports both WebSocket connections for real-time communication and REST endpoints for session management.

## Base URL

- **Local Development**: `http://localhost:8080`
- **Production**: TBD

## Authentication

All chat endpoints require admin authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

For WebSocket connections, the token can be provided as:
- Query parameter: `?token=<your-jwt-token>`
- Authorization header: `Authorization: Bearer <your-jwt-token>`

## WebSocket Connection

### Endpoint
```
GET /api/v1/admin/chat/ws
```

### Connection Protocol

The WebSocket connection follows a structured message protocol for real-time communication.

#### Message Types

- `message` - Chat messages between user and AI
- `typing` - Typing indicators
- `status` - Connection status updates
- `error` - Error notifications
- `presence` - User presence updates
- `ack` - Message acknowledgments
- `heartbeat` - Connection health checks

#### WebSocket Message Format

```json
{
  "type": "message|typing|status|error|presence|ack|heartbeat",
  "session_id": "string",
  "message_id": "string",
  "content": "string",
  "metadata": {
    "client_name": "string",
    "context": "string",
    "quick_action": "string"
  },
  "timestamp": "2025-02-08T10:30:00Z"
}
```

#### Example Messages

**Send Chat Message:**
```json
{
  "type": "message",
  "session_id": "session-123",
  "content": "How do I optimize my AWS costs?",
  "metadata": {
    "client_name": "Acme Corp",
    "context": "cost_optimization"
  },
  "timestamp": "2025-02-08T10:30:00Z"
}
```

**Typing Indicator:**
```json
{
  "type": "typing",
  "session_id": "session-123",
  "content": "true",
  "timestamp": "2025-02-08T10:30:00Z"
}
```

**AI Response:**
```json
{
  "type": "message",
  "session_id": "session-123",
  "message_id": "msg-456",
  "content": "Here are several strategies to optimize your AWS costs...",
  "metadata": {
    "message_type": "assistant",
    "tokens_used": 150
  },
  "timestamp": "2025-02-08T10:30:05Z"
}
```

## REST API Endpoints

### Chat Sessions

#### Create Chat Session

**POST** `/api/v1/admin/chat/sessions`

Creates a new chat session for the authenticated user.

**Request Body:**
```json
{
  "client_name": "Acme Corporation",
  "context": "Architecture review for e-commerce platform",
  "metadata": {
    "meeting_type": "initial_consultation",
    "service_types": ["architecture_review", "optimization"]
  }
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "session-123",
    "user_id": "user-456",
    "client_name": "Acme Corporation",
    "context": "Architecture review for e-commerce platform",
    "status": "active",
    "metadata": {
      "meeting_type": "initial_consultation",
      "service_types": ["architecture_review", "optimization"]
    },
    "created_at": "2025-02-08T10:30:00Z",
    "updated_at": "2025-02-08T10:30:00Z",
    "last_activity": "2025-02-08T10:30:00Z"
  }
}
```

#### List Chat Sessions

**GET** `/api/v1/admin/chat/sessions`

Retrieves all chat sessions for the authenticated user.

**Query Parameters:**
- `limit` (optional): Number of sessions to return (default: 50)
- `offset` (optional): Number of sessions to skip (default: 0)
- `status` (optional): Filter by session status (`active`, `closed`, `expired`)

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "session-123",
      "user_id": "user-456",
      "client_name": "Acme Corporation",
      "context": "Architecture review",
      "status": "active",
      "created_at": "2025-02-08T10:30:00Z",
      "last_activity": "2025-02-08T10:35:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "limit": 50,
    "offset": 0
  }
}
```

#### Get Chat Session

**GET** `/api/v1/admin/chat/sessions/{id}`

Retrieves a specific chat session by ID.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "session-123",
    "user_id": "user-456",
    "client_name": "Acme Corporation",
    "context": "Architecture review for e-commerce platform",
    "status": "active",
    "metadata": {
      "meeting_type": "initial_consultation",
      "service_types": ["architecture_review", "optimization"]
    },
    "created_at": "2025-02-08T10:30:00Z",
    "updated_at": "2025-02-08T10:30:00Z",
    "last_activity": "2025-02-08T10:35:00Z"
  }
}
```

#### Update Chat Session

**PUT** `/api/v1/admin/chat/sessions/{id}`

Updates session context or metadata.

**Request Body:**
```json
{
  "client_name": "Acme Corporation Ltd",
  "context": "Updated context for architecture review",
  "metadata": {
    "meeting_type": "follow_up",
    "priority": "high"
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "session-123",
    "user_id": "user-456",
    "client_name": "Acme Corporation Ltd",
    "context": "Updated context for architecture review",
    "status": "active",
    "metadata": {
      "meeting_type": "follow_up",
      "priority": "high"
    },
    "updated_at": "2025-02-08T10:40:00Z"
  }
}
```

#### Delete Chat Session

**DELETE** `/api/v1/admin/chat/sessions/{id}`

Deletes a chat session and all associated messages.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Chat session deleted successfully"
}
```

#### Get Chat Session History

**GET** `/api/v1/admin/chat/sessions/{id}/history`

Retrieves message history for a specific session.

**Query Parameters:**
- `limit` (optional): Number of messages to return (default: 100)
- `offset` (optional): Number of messages to skip (default: 0)
- `before` (optional): Get messages before this timestamp (ISO 8601)
- `after` (optional): Get messages after this timestamp (ISO 8601)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "session_id": "session-123",
    "messages": [
      {
        "id": "msg-001",
        "session_id": "session-123",
        "type": "user",
        "content": "How do I optimize my AWS costs?",
        "metadata": {
          "client_name": "Acme Corp"
        },
        "created_at": "2025-02-08T10:30:00Z",
        "status": "delivered"
      },
      {
        "id": "msg-002",
        "session_id": "session-123",
        "type": "assistant",
        "content": "Here are several strategies to optimize your AWS costs...",
        "metadata": {
          "tokens_used": 150,
          "response_time_ms": 1200
        },
        "created_at": "2025-02-08T10:30:05Z",
        "status": "delivered"
      }
    ],
    "pagination": {
      "total": 2,
      "limit": 100,
      "offset": 0
    }
  }
}
```

### Chat Metrics

#### Get Chat Metrics

**GET** `/api/v1/admin/chat/metrics`

Retrieves comprehensive chat system metrics.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "connections": {
      "active": 5,
      "total_today": 25,
      "peak_concurrent": 8
    },
    "messages": {
      "sent_today": 150,
      "received_today": 145,
      "average_response_time_ms": 1200
    },
    "ai_usage": {
      "requests_today": 145,
      "tokens_used_today": 25000,
      "average_tokens_per_request": 172
    },
    "errors": {
      "total_today": 3,
      "rate_percent": 2.1
    }
  }
}
```

#### Get Connection Metrics

**GET** `/api/v1/admin/chat/metrics/connections`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "active_connections": 5,
    "total_connections_today": 25,
    "peak_concurrent_connections": 8,
    "average_connection_duration_minutes": 15.5,
    "connection_success_rate": 98.5
  }
}
```

#### Get Health Status

**GET** `/api/v1/admin/chat/health`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "overall_status": "healthy",
    "components": {
      "websocket_server": "healthy",
      "database": "healthy",
      "redis_cache": "healthy",
      "ai_service": "healthy"
    },
    "performance": {
      "average_response_time_ms": 1200,
      "error_rate_percent": 2.1,
      "uptime_percent": 99.9
    }
  }
}
```

## Error Handling

### Error Response Format

All API errors follow a consistent format:

```json
{
  "success": false,
  "error": {
    "code": "CHAT_ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "specific error details"
    },
    "retry_after": 30,
    "correlation_id": "uuid"
  }
}
```

### Common Error Codes

- `CHAT_SESSION_NOT_FOUND` - Session does not exist
- `CHAT_UNAUTHORIZED` - Invalid or expired authentication token
- `CHAT_RATE_LIMITED` - Too many requests, rate limit exceeded
- `CHAT_VALIDATION_ERROR` - Invalid request data
- `CHAT_AI_SERVICE_ERROR` - AI service unavailable or error
- `CHAT_DATABASE_ERROR` - Database connection or query error
- `CHAT_WEBSOCKET_ERROR` - WebSocket connection error

### HTTP Status Codes

- `200 OK` - Successful GET requests
- `201 Created` - Successful POST requests (session creation)
- `400 Bad Request` - Invalid request data or validation errors
- `401 Unauthorized` - Authentication required or invalid token
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server-side errors
- `503 Service Unavailable` - AI service or database unavailable

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **WebSocket Messages**: 60 messages per minute per user
- **REST API Calls**: 100 requests per minute per user
- **Session Creation**: 10 sessions per hour per user

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1644307200
```

## Testing Examples

### Using curl

**Create Chat Session:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/chat/sessions \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "client_name": "Test Client",
    "context": "Testing chat functionality"
  }'
```

**Get Session History:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/chat/sessions/session-123/history \
  -H "Authorization: Bearer your-jwt-token"
```

**Get Chat Metrics:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/chat/metrics \
  -H "Authorization: Bearer your-jwt-token"
```

### WebSocket Testing

Use a WebSocket client to test real-time functionality:

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/admin/chat/ws?token=your-jwt-token');

ws.onopen = function() {
  // Send a chat message
  ws.send(JSON.stringify({
    type: 'message',
    session_id: 'session-123',
    content: 'Hello, AI assistant!',
    timestamp: new Date().toISOString()
  }));
};

ws.onmessage = function(event) {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

## Security Considerations

- All communications use TLS encryption in production
- JWT tokens have configurable expiration times
- Rate limiting prevents abuse and DoS attacks
- Input validation and sanitization prevent XSS attacks
- Session tokens are validated on every request
- Audit logging tracks all chat activities

## Performance Optimization

- WebSocket connections are pooled and reused
- Message history is paginated to reduce load times
- Redis caching improves session retrieval performance
- Database queries are optimized with proper indexing
- AI responses are cached for similar queries

## Monitoring and Observability

The chat system provides comprehensive monitoring through:
- Prometheus metrics for system health
- Structured logging with correlation IDs
- Real-time performance monitoring
- Error tracking and alerting
- Business metrics for usage analytics