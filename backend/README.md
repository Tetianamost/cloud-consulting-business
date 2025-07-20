# Cloud Consulting Backend

A Go-based REST API for managing cloud consulting service inquiries.

## Features

- RESTful API for inquiry management
- Service type configuration
- In-memory storage (development)
- CORS support for frontend integration
- Structured logging
- Health check endpoint

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository and navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Copy the environment file:
```bash
cp .env.example .env
```

4. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8061 by default.

### Verify Installation

Test the health endpoint:
```bash
curl http://localhost:8061/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "cloud-consulting-backend",
  "version": "1.0.0",
  "time": "2025-07-19T18:18:35Z"
}
```

## Configuration

The application uses environment variables for configuration. Key settings:

- `PORT`: Server port (default: 8061)
- `LOG_LEVEL`: Logging level (default: 4 - Info)
- `GIN_MODE`: Gin framework mode (debug/release)
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins

## API Documentation

See [API Documentation](docs/api/README.md) for detailed endpoint information.

### Available Endpoints

- `GET /health` - Health check
- `GET /api/v1/config/services` - Get available service types
- `POST /api/v1/inquiries` - Create new inquiry
- `GET /api/v1/inquiries` - List all inquiries
- `GET /api/v1/inquiries/{id}` - Get specific inquiry

## Development

### Project Structure

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── domain/         # Domain models and constants
│   ├── handlers/       # HTTP request handlers
│   ├── server/         # Server setup and routing
│   └── storage/        # Data storage layer
├── docs/               # Documentation
└── scripts/            # Database scripts
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o server cmd/server/main.go
```

## Frontend Integration

The backend is designed to work with the React frontend. The frontend uses the API service located at `/frontend/src/services/api.ts` to communicate with these endpoints.

### CORS Configuration

The server is configured to accept requests from:
- `http://localhost:3000` (React development server)
- `http://localhost:3001` (Alternative development port)

## Data Storage

Currently using in-memory storage for development purposes. Data will be lost when the server restarts.

**Note**: Production deployment will use PostgreSQL database as specified in the design document.

## Logging

The application uses structured JSON logging with the following levels:
- Error (1)
- Warn (2) 
- Info (4)
- Debug (5)

Logs include request details such as:
- HTTP method and path
- Response status
- Request latency
- Client IP
- User agent

## Next Steps

This is a minimal working implementation. Future enhancements include:

- Database integration (PostgreSQL)
- Authentication and authorization
- Input validation and sanitization
- Rate limiting
- Comprehensive testing
- Docker containerization
- Production deployment configuration
- AI report generation
- Email notifications
- Monitoring and metrics

## Support

For questions or issues, please refer to the main project documentation or contact the development team.