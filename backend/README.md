# Cloud Consulting Backend

A Go-based REST API for managing cloud consulting service inquiries.

## Features

- RESTful API for inquiry management
- AI-powered draft report generation using Amazon Bedrock
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

4. Configure Amazon Bedrock API key:
   - Generate an API key from the Amazon Bedrock console
   - Get your AWS access keys from IAM console
   - Update the `.env` file with your credentials:
   ```bash
   AWS_ACCESS_KEY_ID=your_aws_access_key
   AWS_SECRET_ACCESS_KEY=your_aws_secret_key
   AWS_BEARER_TOKEN_BEDROCK=your_bedrock_api_key_here
   # AWS_SESSION_TOKEN is optional - only needed for temporary credentials
   ```

5. Start the server:
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

### Basic Configuration
- `PORT`: Server port (default: 8061)
- `LOG_LEVEL`: Logging level (default: 4 - Info)
- `GIN_MODE`: Gin framework mode (debug/release)
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed origins

### Amazon Bedrock Configuration
- `AWS_BEARER_TOKEN_BEDROCK`: Your Bedrock API key (required for AI report generation)
- `BEDROCK_REGION`: AWS region for Bedrock (default: us-east-1)
- `BEDROCK_MODEL_ID`: Bedrock model to use (default: amazon.nova-lite-v1:0)
- `BEDROCK_BASE_URL`: Bedrock API base URL (default: https://bedrock-runtime.us-east-1.amazonaws.com)
- `BEDROCK_TIMEOUT_SECONDS`: Request timeout in seconds (default: 30)

## API Documentation

See [API Documentation](docs/api/README.md) for detailed endpoint information.

### Available Endpoints

- `GET /health` - Health check
- `GET /api/v1/config/services` - Get available service types
- `POST /api/v1/inquiries` - Create new inquiry (automatically generates AI report)
- `GET /api/v1/inquiries` - List all inquiries
- `GET /api/v1/inquiries/{id}` - Get specific inquiry
- `GET /api/v1/inquiries/{id}/report` - Get AI-generated report for inquiry

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

## AI Report Generation

The backend automatically generates draft reports for new inquiries using Amazon Bedrock's Nova model.

### How it works:
1. When a new inquiry is created via `POST /api/v1/inquiries`, the system automatically triggers report generation
2. The inquiry details are sent to Amazon Bedrock with a structured prompt
3. Bedrock generates a professional consulting report draft
4. The report is stored and linked to the inquiry
5. Reports can be retrieved via `GET /api/v1/inquiries/{id}/report`

### Error Handling:
- If Bedrock API fails, the inquiry is still created successfully
- Report generation failures are logged but don't block inquiry processing
- The system gracefully degrades when Bedrock is unavailable

### Example Report Structure:
- Executive Summary
- Current State Assessment  
- Recommendations
- Next Steps

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

## Docker Development Setup

### Quick Start with Docker Compose

1. **Copy environment configuration:**
   ```bash
   cp .env.example .env
   ```

2. **Edit .env file with your AWS credentials:**
   ```bash
   # Required for Bedrock integration
   AWS_BEARER_TOKEN_BEDROCK=your_bedrock_api_key_here
   AWS_ACCESS_KEY_ID=your_aws_access_key
   AWS_SECRET_ACCESS_KEY=your_aws_secret_key
   ```

3. **Start the services:**
   ```bash
   # Start backend and frontend
   docker-compose up -d backend frontend
   
   # Or use the helper script
   ./scripts/docker-dev.sh up
   ```

4. **Access the application:**
   - Backend API: http://localhost:8061
   - Frontend: http://localhost:3000
   - Health Check: http://localhost:8061/health

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Start with database and cache
docker-compose --profile database --profile cache up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down

# Rebuild services
docker-compose build

# Clean up everything
docker-compose down -v --remove-orphans
```

## Testing the API

### Create an inquiry with AI report generation:
```bash
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Tech Corp",
    "services": ["assessment"],
    "message": "We need help assessing our current AWS infrastructure for cost optimization and security improvements."
  }'
```

### Get the generated report:
```bash
curl http://localhost:8061/api/v1/inquiries/{inquiry-id}/report
```

### List all inquiries (with reports):
```bash
curl http://localhost:8061/api/v1/inquiries
```

## Next Steps

This is a minimal working implementation. Future enhancements include:

- Database integration (PostgreSQL)
- Authentication and authorization
- Input validation and sanitization
- Rate limiting
- Comprehensive testing
- Docker containerization
- Production deployment configuration
- Email notifications
- Monitoring and metrics

## Support

For questions or issues, please refer to the main project documentation or contact the development team.