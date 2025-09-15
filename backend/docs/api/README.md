# Cloud Consulting Backend API Documentation

## Overview

The Cloud Consulting Backend provides RESTful APIs for managing service inquiries from potential clients. The system accepts inquiries, categorizes them by service type, and stores them for consultant review.

## Base URL

- **Local Development**: `http://localhost:8061`
- **Production**: TBD

## Authentication

Currently, the API does not require authentication for inquiry submission endpoints.

## API Endpoints

### Health Check

#### GET /health

Returns the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "service": "cloud-consulting-backend",
  "version": "1.0.0",
  "time": "2025-07-19T18:18:35Z"
}
```

### Service Configuration

#### GET /api/v1/config/services

Returns available service types and their descriptions.

**Response:**
```json
{
  "success": true,
  "data": {
    "services": [
      {
        "id": "assessment",
        "name": "Cloud Assessment",
        "description": "Comprehensive evaluation of your current cloud infrastructure, security posture, and optimization opportunities"
      },
      {
        "id": "migration",
        "name": "Cloud Migration",
        "description": "Strategic planning and execution support for migrating workloads to the cloud"
      },
      {
        "id": "optimization",
        "name": "Cloud Optimization",
        "description": "Performance tuning, cost optimization, and efficiency improvements for existing cloud deployments"
      },
      {
        "id": "architecture_review",
        "name": "Architecture Review",
        "description": "Expert review of cloud architecture designs for scalability, security, and best practices compliance"
      }
    ]
  }
}
```

### Inquiries

#### POST /api/v1/inquiries

Creates a new service inquiry.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "company": "Test Company",
  "phone": "+1-555-123-4567",
  "services": ["assessment", "migration"],
  "message": "I need help with cloud migration planning.",
  "source": "contact_form"
}
```

**Required Fields:**
- `name` (string): Client's full name
- `email` (string): Valid email address
- `services` (array): Array of service IDs (must be valid service types)
- `message` (string): Inquiry message

**Optional Fields:**
- `company` (string): Company name
- `phone` (string): Phone number
- `source` (string): Source of the inquiry (e.g., "contact_form", "quote_request")

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "a485bdeb-446f-4a45-b401-a369ce3a5318",
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Test Company",
    "phone": "+1-555-123-4567",
    "services": ["assessment", "migration"],
    "message": "I need help with cloud migration planning.",
    "status": "pending",
    "priority": "medium",
    "source": "contact_form",
    "created_at": "2025-07-19T18:19:28.009554-06:00",
    "updated_at": "2025-07-19T18:19:28.009554-06:00"
  },
  "message": "Inquiry created successfully"
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "At least one service must be selected"
}
```

#### GET /api/v1/inquiries/{id}

Retrieves a specific inquiry by ID.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "a485bdeb-446f-4a45-b401-a369ce3a5318",
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Test Company",
    "phone": "+1-555-123-4567",
    "services": ["assessment", "migration"],
    "message": "I need help with cloud migration planning.",
    "status": "pending",
    "priority": "medium",
    "source": "contact_form",
    "created_at": "2025-07-19T18:19:28.009554-06:00",
    "updated_at": "2025-07-19T18:19:28.009554-06:00"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "error": "Inquiry not found"
}
```

#### GET /api/v1/inquiries

Retrieves all inquiries.

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "a485bdeb-446f-4a45-b401-a369ce3a5318",
      "name": "John Doe",
      "email": "john@example.com",
      "company": "Test Company",
      "phone": "+1-555-123-4567",
      "services": ["assessment", "migration"],
      "message": "I need help with cloud migration planning.",
      "status": "pending",
      "priority": "medium",
      "source": "contact_form",
      "created_at": "2025-07-19T18:19:28.009554-06:00",
      "updated_at": "2025-07-19T18:19:28.009554-06:00"
    }
  ],
  "count": 1
}
```

## Service Types

The following service types are currently supported:

- **assessment**: Cloud Assessment
- **migration**: Cloud Migration  
- **optimization**: Cloud Optimization
- **architecture_review**: Architecture Review

## Error Handling

All API responses follow a consistent format:

**Success Response:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Optional success message"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error description"
}
```

## CORS Configuration

The API is configured to accept requests from:
- `http://localhost:3000` (React development server)
- `http://localhost:3001` (Alternative development port)

## Testing Examples

### Using curl

**Health Check:**
```bash
curl -X GET http://localhost:8061/health
```

**Get Service Configuration:**
```bash
curl -X GET http://localhost:8061/api/v1/config/services
```

**Create Inquiry:**
```bash
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "services": ["assessment"],
    "message": "I need a cloud assessment."
  }'
```

**Get All Inquiries:**
```bash
curl -X GET http://localhost:8061/api/v1/inquiries
```

## Frontend Integration

The frontend uses the `/src/services/api.ts` service to communicate with the backend:

```typescript
import { apiService } from '../services/api';

// Create an inquiry
const response = await apiService.createInquiry({
  name: 'John Doe',
  email: 'john@example.com',
  services: ['assessment'],
  message: 'I need help with cloud migration.'
});
```

## Status Codes

- **200 OK**: Successful GET requests
- **201 Created**: Successful POST requests (inquiry creation)
- **400 Bad Request**: Invalid request data or validation errors
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server-side errors

## Data Storage

Currently using in-memory storage for development. Data will be lost when the server restarts. Production deployment will use PostgreSQL database.