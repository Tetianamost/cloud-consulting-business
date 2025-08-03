# Cloud Consulting Backend

A comprehensive Go backend system for processing and categorizing cloud consulting service inquiries with AI-powered report generation and automated notifications.

## Overview

This backend system handles four main service types:
- **Assessment**: Comprehensive evaluation of current cloud infrastructure
- **Migration**: Strategic planning and execution support for cloud migrations  
- **Optimization**: Performance tuning and cost optimization for existing deployments
- **Architecture Review**: Expert review of cloud architecture designs

## Features

- RESTful API with comprehensive validation
- AI-powered draft report generation using LLM integration
- Automated consultant notifications via email and Slack
- Agent hooks system for extensible workflow automation
- Comprehensive logging and monitoring with Prometheus metrics
- Secure AWS deployment with horizontal scaling support
- PostgreSQL database with GORM ORM
- Redis caching and session management
- Circuit breaker and retry patterns for resilience

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Redis 6+
- Docker and Docker Compose (for local development)

### Local Development

1. Clone the repository and navigate to the backend directory
2. Copy the environment file: `cp .env.example .env`
3. Update the `.env` file with your configuration
4. Start dependencies: `docker-compose up -d db redis`
5. Run database migrations: `make migrate`
6. Start the server: `make run`

The API will be available at `http://localhost:8080`

### Docker Development

```bash
# Start all services including the API
docker-compose up

# Or start in detached mode
docker-compose up -d
```

## API Documentation

### Health Check
- `GET /health` - System health status

### Inquiries
- `POST /api/v1/inquiries` - Create new inquiry
- `GET /api/v1/inquiries/{id}` - Get inquiry details
- `GET /api/v1/inquiries` - List inquiries with filters
- `PUT /api/v1/inquiries/{id}/status` - Update inquiry status
- `GET /api/v1/inquiries/{id}/report` - Get generated report

### AI Consultant Live Chat
- `GET /api/v1/admin/chat/ws` - WebSocket connection for real-time chat
- `POST /api/v1/admin/chat/sessions` - Create new chat session
- `GET /api/v1/admin/chat/sessions` - List chat sessions
- `GET /api/v1/admin/chat/sessions/{id}` - Get specific chat session
- `PUT /api/v1/admin/chat/sessions/{id}` - Update chat session
- `DELETE /api/v1/admin/chat/sessions/{id}` - Delete chat session
- `GET /api/v1/admin/chat/sessions/{id}/history` - Get session message history
- `GET /api/v1/admin/chat/metrics` - Chat system metrics
- `GET /api/v1/admin/chat/health` - Chat system health status

### System Management
- `GET /api/v1/metrics` - Prometheus metrics
- `POST /api/v1/hooks/trigger` - Manual hook trigger
- `GET /api/v1/hooks` - List active hooks
- `GET /api/v1/config/services` - Get available service types

## Architecture

The system follows a layered architecture pattern:

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and constants
│   ├── handlers/        # HTTP request handlers
│   ├── interfaces/      # Interface definitions
│   ├── middleware/      # HTTP middleware
│   ├── repositories/    # Data access layer
│   ├── server/          # HTTP server setup
│   ├── services/        # Business logic layer
│   └── utils/           # Utility functions
├── pkg/                 # Shared packages
│   └── logger/          # Logging utilities
├── docs/                # Documentation
└── scripts/             # Database scripts and utilities
```

## Configuration

The application uses environment variables for configuration. See `.env.example` for all available options.

Key configuration areas:
- **Database**: PostgreSQL connection and pool settings
- **Redis**: Cache and session configuration  
- **AWS**: S3 storage and SQS messaging
- **LLM**: AI service integration for report generation
- **Notifications**: Email and Slack integration
- **Security**: JWT, CORS, and rate limiting
- **Monitoring**: Metrics, logging, and health checks

## Development

### Make Commands

- `make run` - Start the development server
- `make build` - Build the application binary
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make lint` - Run code linting
- `make migrate` - Run database migrations
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container

### Testing

The project includes comprehensive testing:
- Unit tests for all services and utilities
- Integration tests for database operations
- API endpoint testing
- Mock implementations for external services

Run tests with: `make test`

### Code Quality

- Go fmt for code formatting
- golangci-lint for comprehensive linting
- Minimum 80% test coverage requirement
- Pre-commit hooks for quality checks

## Deployment

### AWS Deployment

The application is designed for AWS deployment using:
- **ECS Fargate** or **EKS** for container orchestration
- **RDS PostgreSQL** for the database
- **ElastiCache Redis** for caching
- **S3** for file storage
- **SQS** for message queuing
- **CloudWatch** for logging and monitoring
- **Application Load Balancer** for traffic distribution

### Environment Setup

1. **Development**: Local Docker Compose setup
2. **Staging**: AWS ECS with shared RDS instance
3. **Production**: AWS ECS with dedicated RDS and Redis clusters

## Monitoring

### Health Checks
- Database connectivity
- Redis connectivity  
- External service availability
- Memory and CPU usage

### Metrics
- Request count and latency
- Error rates by endpoint
- Database query performance
- Cache hit/miss ratios
- Business metrics (inquiries, reports, notifications)

### Logging
- Structured JSON logging
- Request/response logging with correlation IDs
- Error logging with stack traces
- Audit logging for all business operations

## Security

- JWT-based authentication (optional)
- CORS configuration for frontend integration
- Rate limiting per IP and user
- Input validation and sanitization
- SQL injection prevention via GORM
- Secrets management via AWS Secrets Manager
- TLS encryption in transit
- Database encryption at rest

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run the test suite and linting
5. Submit a pull request

## License

This project is proprietary software for cloud consulting services.

## Support

For questions or issues, please contact the development team or create an issue in the project repository.