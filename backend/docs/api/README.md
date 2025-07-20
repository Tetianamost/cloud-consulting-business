# API Documentation

This directory contains API documentation for the Cloud Consulting Backend.

## Overview

The API follows RESTful conventions and returns JSON responses. All endpoints use standard HTTP status codes and include comprehensive error handling.

## Base URL

- Development: `http://localhost:8080`
- Production: `https://api.your-domain.com`

## Authentication

Authentication is optional and can be enabled via the `ENABLE_AUTHENTICATION` environment variable. When enabled, endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <jwt-token>
```

## Rate Limiting

API requests are rate-limited to prevent abuse:
- Default: 100 requests per second per IP
- Burst: 200 requests
- Configurable via `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST`

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "timestamp": "2024-01-15T10:30:00Z",
  "trace_id": "abc123"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": "Email is required",
    "trace_id": "abc123",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "trace_id": "abc123"
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `NOT_FOUND` | Resource not found |
| `UNAUTHORIZED` | Authentication required |
| `FORBIDDEN` | Insufficient permissions |
| `CONFLICT` | Resource conflict |
| `RATE_LIMIT_EXCEEDED` | Rate limit exceeded |
| `INTERNAL_ERROR` | Internal server error |
| `SERVICE_UNAVAILABLE` | Service temporarily unavailable |
| `TIMEOUT` | Request timeout |
| `BAD_REQUEST` | Invalid request |

## Endpoints

### Health Check
- [GET /health](endpoints/health.md) - System health status

### Inquiries
- [POST /api/v1/inquiries](endpoints/inquiries.md#create-inquiry) - Create new inquiry
- [GET /api/v1/inquiries/{id}](endpoints/inquiries.md#get-inquiry) - Get inquiry details
- [GET /api/v1/inquiries](endpoints/inquiries.md#list-inquiries) - List inquiries
- [PUT /api/v1/inquiries/{id}/status](endpoints/inquiries.md#update-status) - Update inquiry status
- [GET /api/v1/inquiries/{id}/report](endpoints/inquiries.md#get-report) - Get generated report

### System Management
- [GET /api/v1/metrics](endpoints/system.md#metrics) - Prometheus metrics
- [POST /api/v1/hooks/trigger](endpoints/system.md#trigger-hook) - Manual hook trigger
- [GET /api/v1/hooks](endpoints/system.md#list-hooks) - List active hooks
- [GET /api/v1/config/services](endpoints/system.md#service-config) - Get service configuration

## Data Models

### Inquiry
```json
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "company": "Acme Corp",
  "phone": "+1-555-0123",
  "services": ["assessment", "migration"],
  "message": "We need help with our cloud migration",
  "status": "pending",
  "priority": "medium",
  "source": "website",
  "utm_params": {
    "utm_source": "google",
    "utm_medium": "cpc"
  },
  "assigned_to": "consultant-id",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Report
```json
{
  "id": "uuid",
  "inquiry_id": "uuid",
  "type": "draft",
  "title": "Cloud Assessment Report",
  "content": "Report content...",
  "status": "draft",
  "generated_by": "ai-agent",
  "reviewed_by": "consultant-id",
  "s3_key": "reports/uuid.pdf",
  "metadata": {
    "word_count": 1500,
    "generation_time": 30
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Activity
```json
{
  "id": "uuid",
  "inquiry_id": "uuid",
  "type": "inquiry_created",
  "description": "New inquiry received from John Doe",
  "actor": "system",
  "metadata": {
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0..."
  },
  "created_at": "2024-01-15T10:30:00Z"
}
```

## Testing

Use the provided Postman collection or curl commands to test the API endpoints. Examples are provided in each endpoint documentation file.