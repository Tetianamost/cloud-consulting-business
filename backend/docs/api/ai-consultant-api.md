# AI Consultant API Documentation

## Overview

The AI Consultant API provides endpoints for interacting with the AI-powered consulting assistant. It supports both simple chat functionality and advanced context-aware conversations with quick actions and session management.

## Authentication

All AI Consultant endpoints require JWT authentication:

```http
Authorization: Bearer <jwt_token>
```

## Endpoints

### Simple Chat API

#### Send Message

**Endpoint:** `POST /api/v1/admin/simple-chat/messages`

**Description:** Send a message to the AI assistant and receive an immediate response.

**Request Body:**
```json
{
  "content": "How much would it cost to migrate 100 VMs to AWS?",
  "session_id": "session-1234567890-abc123"
}
```

**Response:**
```json
{
  "success": true,
  "message_id": "msg-1234567890-user"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Invalid request format"
}
```

#### Get Messages

**Endpoint:** `GET /api/v1/admin/simple-chat/messages`

**Description:** Retrieve all messages for a specific chat session.

**Query Parameters:**
- `session_id` (required): The session ID to retrieve messages for

**Example Request:**
```http
GET /api/v1/admin/simple-chat/messages?session_id=session-1234567890-abc123
```

**Response:**
```json
{
  "success": true,
  "messages": [
    {
      "id": "msg-1234567890-user",
      "content": "How much would it cost to migrate 100 VMs to AWS?",
      "role": "user",
      "timestamp": "2024-01-15T10:30:00Z",
      "session_id": "session-1234567890-abc123"
    },
    {
      "id": "msg-1234567890-ai",
      "content": "Based on your requirements, I recommend implementing a cloud-first architecture...",
      "role": "assistant",
      "timestamp": "2024-01-15T10:30:01Z",
      "session_id": "session-1234567890-abc123"
    }
  ]
}
```

## Quick Actions

The AI Consultant supports 8 pre-defined quick actions that can be triggered from the frontend:

### Available Quick Actions

| Action ID | Label | Purpose | Example Prompt |
|-----------|-------|---------|----------------|
| `cost-estimate` | Cost Estimate | Get pricing analysis | "Provide a cost estimate for this solution" |
| `security-review` | Security Review | Analyze security considerations | "What are the security considerations for this approach?" |
| `best-practices` | Best Practices | Get AWS best practices | "What are the AWS best practices for this scenario?" |
| `alternatives` | Alternatives | Explore alternative approaches | "What are alternative approaches to consider?" |
| `next-steps` | Next Steps | Get actionable next steps | "What are the recommended next steps?" |
| `migration-plan` | Migration Plan | Get migration strategy | "What is the recommended migration approach?" |
| `compliance` | Compliance | Address compliance requirements | "What compliance considerations should we address?" |
| `performance` | Performance | Get performance optimization | "How can we optimize performance for this solution?" |

## Context Management

### Client Context

The AI assistant supports context-aware conversations by maintaining:

- **Client Name**: Personalizes responses with client-specific references
- **Meeting Type**: Provides situational awareness (e.g., "Migration planning", "Cost optimization")
- **Session History**: Maintains conversation context within a session

### Session Management

Sessions are automatically created and managed:

- **Session ID Format**: `session-{timestamp}-{random}`
- **Session Persistence**: Messages are stored in memory during the session
- **Session Reset**: Can be reset to start fresh conversations

## Response Format

### AI Response Structure

AI responses follow a consistent format:

```json
{
  "id": "msg-{timestamp}-ai",
  "content": "Professional consulting response with actionable recommendations...",
  "role": "assistant",
  "timestamp": "2024-01-15T10:30:01Z",
  "session_id": "session-1234567890-abc123"
}
```

### Response Characteristics

- **Professional Tone**: Responses maintain a professional consulting tone
- **Actionable Advice**: Includes specific, implementable recommendations
- **Context Awareness**: References client name and meeting context when available
- **AWS Focus**: Specialized in AWS cloud consulting scenarios

## Error Handling

### Common Error Codes

| Status Code | Error Type | Description |
|-------------|------------|-------------|
| 400 | Bad Request | Invalid request format or missing required fields |
| 401 | Unauthorized | Missing or invalid JWT token |
| 500 | Internal Server Error | Server-side processing error |

### Error Response Format

```json
{
  "success": false,
  "error": "Descriptive error message"
}
```

## Rate Limiting

- **Default Limit**: 60 requests per minute per user
- **Burst Limit**: 10 requests per 10 seconds
- **Headers**: Rate limit information included in response headers

## Integration Examples

### Frontend Integration (React/TypeScript)

```typescript
import simpleAIService from '../services/simpleAIService';

// Send a message
const response = await simpleAIService.sendMessage({
  message: "How can I optimize AWS costs?",
  context: {
    clientName: "Acme Corp",
    meetingType: "Cost optimization"
  }
});

// Get message history
const messages = await simpleAIService.getMessages();
```

### cURL Examples

**Send Message:**
```bash
curl -X POST http://localhost:8061/api/v1/admin/simple-chat/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "content": "What are the best practices for AWS security?",
    "session_id": "session-1234567890-abc123"
  }'
```

**Get Messages:**
```bash
curl -X GET "http://localhost:8061/api/v1/admin/simple-chat/messages?session_id=session-1234567890-abc123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Performance Considerations

### Response Times

- **Target Response Time**: < 200ms for simple responses
- **Complex Queries**: May take up to 2 seconds for detailed analysis
- **Caching**: Responses are not cached to ensure fresh, contextual advice

### Scalability

- **In-Memory Storage**: Current implementation uses in-memory message storage
- **Session Limits**: No hard limits on session duration or message count
- **Concurrent Users**: Supports multiple concurrent chat sessions

## Security

### Data Protection

- **Message Encryption**: All API communications use HTTPS/TLS
- **Session Isolation**: Messages are isolated by session ID
- **No Persistent Storage**: Messages are not permanently stored (demo mode)

### Authentication

- **JWT Validation**: All requests require valid JWT tokens
- **Token Expiration**: Tokens expire after 24 hours by default
- **Admin Only**: Only admin users can access AI consultant features

## Monitoring and Logging

### Logging

The API logs the following events:

- Message send/receive operations
- Session creation and management
- Error conditions and recovery
- Performance metrics

### Metrics

Key metrics tracked:

- Messages per session
- Response times
- Error rates
- Active sessions

## Future Enhancements

### Planned Features

- **Persistent Storage**: Database storage for message history
- **Advanced Context**: Integration with client CRM data
- **File Uploads**: Support for document analysis
- **Voice Integration**: Speech-to-text capabilities
- **Multi-language**: Support for multiple languages

### API Versioning

- **Current Version**: v1
- **Backward Compatibility**: Maintained for at least 6 months
- **Deprecation Notice**: 30-day notice for breaking changes

This API provides a robust foundation for AI-powered consulting conversations with professional-grade features and enterprise-ready architecture.