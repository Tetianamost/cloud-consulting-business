# Project Standards and Guidelines

## Overview

This document outlines the coding standards, architectural patterns, and development practices for the Cloud Consulting Platform - a comprehensive Go backend with React frontend for managing cloud consulting service inquiries with AI-powered report generation.

## Architecture Overview

### Backend (Go)
- **Framework**: Gin web framework for HTTP routing
- **Language**: Go 1.24+ with modern toolchain
- **Database**: PostgreSQL with GORM ORM (in-memory for development)
- **AI Integration**: AWS Bedrock (Nova model) for report generation
- **Email**: AWS SES for notifications
- **Authentication**: JWT-based admin authentication
- **Caching**: Redis for performance optimization
- **Real-time**: WebSocket and polling-based chat systems

### Frontend (React)
- **Framework**: React 19+ with TypeScript
- **State Management**: Redux Toolkit
- **UI Components**: Radix UI primitives with Tailwind CSS
- **Routing**: React Router v6
- **Forms**: Formik with Yup validation
- **HTTP Client**: Axios
- **Testing**: Jest with React Testing Library

## Code Organization

### Backend Structure
```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── domain/         # Domain models and business logic
│   ├── handlers/       # HTTP request handlers
│   ├── interfaces/     # Service interfaces and contracts
│   ├── repositories/   # Data access layer
│   ├── services/       # Business logic services
│   ├── server/         # Server setup and middleware
│   └── storage/        # Storage implementations
├── docs/               # API and system documentation
└── scripts/            # Database migrations and utilities
```

### Frontend Structure
```
frontend/src/
├── components/         # Reusable UI components
│   ├── admin/         # Admin-specific components
│   └── ui/            # Base UI components
├── hooks/             # Custom React hooks
├── services/          # API clients and external services
├── store/             # Redux store and slices
├── types/             # TypeScript type definitions
└── utils/             # Utility functions
```

## Development Standards

### Go Backend Standards

#### Error Handling
- Use structured error responses with appropriate HTTP status codes
- Log errors with context using structured logging (logrus)
- Implement graceful degradation for external service failures
- Never expose internal errors to clients

```go
// Good
if err != nil {
    log.WithFields(log.Fields{
        "error": err.Error(),
        "context": "bedrock_api_call",
    }).Error("Failed to generate AI report")
    return c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Report generation temporarily unavailable"
    })
}
```

#### Service Layer Pattern
- Implement business logic in service layer, not handlers
- Use interfaces for service contracts
- Inject dependencies through constructors
- Keep handlers thin - they should only handle HTTP concerns

```go
type ChatService interface {
    SendMessage(ctx context.Context, sessionID string, message string) (*ChatMessage, error)
    GetHistory(ctx context.Context, sessionID string) ([]*ChatMessage, error)
}
```

#### Database Patterns
- Use repository pattern for data access
- Implement proper transaction handling
- Use GORM for ORM operations
- Include database migrations in scripts/

#### API Design
- Follow RESTful conventions
- Use consistent response formats
- Implement proper CORS configuration
- Include comprehensive API documentation

### React Frontend Standards

#### Component Organization
- Use functional components with hooks
- Implement proper TypeScript typing
- Follow single responsibility principle
- Use composition over inheritance

```tsx
interface ChatMessageProps {
  message: ChatMessage;
  isOwn: boolean;
  onRetry?: () => void;
}

export const ChatMessage: React.FC<ChatMessageProps> = ({ message, isOwn, onRetry }) => {
  // Component implementation
};
```

#### State Management
- Use Redux Toolkit for global state
- Keep local state in components when appropriate
- Implement proper action creators and reducers
- Use RTK Query for API state management

#### Styling Standards
- Use Tailwind CSS for styling
- Implement responsive design patterns
- Follow accessibility guidelines (WCAG 2.1)
- Use Radix UI for complex interactive components

## Testing Standards

### Backend Testing
- Minimum 90% test coverage for new features
- Unit tests for all service layer functions
- Integration tests for API endpoints
- End-to-end tests for critical user flows

```go
func TestChatService_SendMessage(t *testing.T) {
    // Arrange
    service := NewChatService(mockRepo, mockAI)
    
    // Act
    result, err := service.SendMessage(ctx, "session-123", "Hello")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Frontend Testing
- Component unit tests with React Testing Library
- Integration tests for user workflows
- Accessibility testing with jest-axe
- Visual regression testing for UI components

## Security Standards

### Authentication & Authorization
- Use JWT tokens with appropriate expiration
- Implement proper session management
- Validate all user inputs
- Use HTTPS in production

### Data Protection
- Encrypt sensitive data at rest
- Use TLS for all communications
- Implement proper CORS policies
- Follow OWASP security guidelines

## Performance Standards

### Backend Performance
- Implement caching strategies (Redis)
- Use connection pooling for databases
- Optimize database queries
- Implement rate limiting

### Frontend Performance
- Use React.memo for expensive components
- Implement virtual scrolling for large lists
- Optimize bundle size with code splitting
- Use proper loading states and error boundaries

## AI Integration Standards

### AWS Bedrock Integration
- Use structured prompts for consistent outputs
- Implement proper error handling and fallbacks
- Cache AI responses when appropriate
- Monitor API usage and costs

### Prompt Engineering
- Use clear, specific prompts
- Include relevant context and constraints
- Implement prompt versioning
- Test prompts thoroughly before deployment

## Deployment Standards

### Docker Configuration
- Use multi-stage builds for optimization
- Implement proper health checks
- Use environment variables for configuration
- Follow security best practices in containers

### Environment Management
- Separate configurations for dev/staging/prod
- Use environment variables for secrets
- Implement proper logging and monitoring
- Use infrastructure as code (Kubernetes manifests)

## Documentation Standards

### Code Documentation
- Document all public interfaces
- Include usage examples
- Maintain up-to-date README files
- Use clear, descriptive comments

### API Documentation
- Maintain OpenAPI/Swagger specifications
- Include request/response examples
- Document error conditions
- Provide integration guides

## Quality Assurance

### Code Review Process
- All code must be reviewed before merging
- Use automated linting and formatting
- Run full test suite on CI/CD
- Check for security vulnerabilities

### Continuous Integration
- Automated testing on all pull requests
- Code quality checks (linting, formatting)
- Security scanning
- Performance regression testing

## Monitoring and Observability

### Logging
- Use structured logging (JSON format)
- Include correlation IDs for request tracing
- Log at appropriate levels (ERROR, WARN, INFO, DEBUG)
- Implement log aggregation and analysis

### Metrics
- Monitor application performance metrics
- Track business metrics (inquiries, reports generated)
- Monitor external service dependencies
- Implement alerting for critical issues

### Health Checks
- Implement comprehensive health check endpoints
- Monitor database connectivity
- Check external service availability
- Include dependency status in health responses