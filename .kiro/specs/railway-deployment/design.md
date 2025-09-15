# Railway Deployment Design

## Overview

This document provides a comprehensive design for deploying the Cloud Consulting Platform on Railway. Railway is a modern deployment platform that eliminates infrastructure complexity while providing production-ready features like automatic scaling, built-in databases, and seamless CI/CD.

## Architecture

### High-Level Architecture Diagram

```
                    Internet
                       │
                       ▼
              ┌─────────────────┐
              │   Custom Domain │
              │ (your-domain.com)│
              │   + Auto SSL    │
              └─────────┬───────┘
                        │
              ┌─────────▼───────┐
              │  Railway Edge   │
              │   Load Balancer │
              └─────────┬───────┘
                        │
    ┌───────────────────┼────────────────────┐
    │                   │                    │
    ▼                   ▼                    ▼
┌─────────┐      ┌─────────────┐      ┌─────────────┐
│Frontend │      │   Backend   │      │ PostgreSQL  │
│(React)  │      │    (Go)     │      │  Database   │
│Service  │      │   Service   │      │   Service   │
└─────────┘      └─────────────┘      └─────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │  External AWS   │
              │   Services      │
              │ • Bedrock (AI)  │
              │ • SES (Email)   │
              └─────────────────┘
```

### Railway Service Architecture

Railway automatically detects and deploys your application as separate services:

1. **Frontend Service**: React application served as static files
2. **Backend Service**: Go API server with automatic scaling
3. **Database Service**: Managed PostgreSQL with automatic backups
4. **Edge Network**: Global CDN and load balancing

## Components and Interfaces

### 1. Project Structure for Railway

Railway works best with a monorepo structure. Your current structure is perfect:

```
cloud-consulting-platform/
├── backend/
│   ├── cmd/server/main.go
│   ├── Dockerfile              # Railway auto-detects this
│   ├── go.mod
│   └── internal/
├── frontend/
│   ├── package.json           # Railway auto-detects this
│   ├── src/
│   └── public/
├── railway.json               # Railway configuration (optional)
└── README.md
```

### 2. Railway Configuration Files

#### Optional Railway Configuration (`railway.json`)
```json
{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "NIXPACKS"
  },
  "deploy": {
    "numReplicas": 1,
    "sleepApplication": false,
    "restartPolicyType": "ON_FAILURE"
  }
}
```

#### Backend Dockerfile (Railway auto-detects)
```dockerfile
# Railway will use this automatically
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates

EXPOSE 8061
CMD ["./server"]
```

#### Frontend Package.json (Railway auto-detects)
```json
{
  "name": "cloud-consulting-frontend",
  "scripts": {
    "build": "npm run build",
    "start": "serve -s build -l 3000"
  },
  "dependencies": {
    "serve": "^14.0.0"
  }
}
```

### 3. Environment Configuration

#### Backend Environment Variables
```bash
# Railway provides these automatically
PORT=8061
DATABASE_URL=postgresql://postgres:password@host:5432/railway

# Your existing AWS configuration (add these in Railway dashboard)
AWS_ACCESS_KEY_ID=your_existing_key
AWS_SECRET_ACCESS_KEY=your_existing_secret
AWS_BEARER_TOKEN_BEDROCK=your_existing_token
BEDROCK_REGION=us-east-1
BEDROCK_MODEL_ID=amazon.nova-lite-v1:0
AWS_SES_REGION=us-east-1
SES_SENDER_EMAIL=info@cloudpartner.pro
SES_REPLY_TO_EMAIL=info@cloudpartner.pro

# Application configuration
GIN_MODE=release
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=https://your-domain.com
JWT_SECRET=your_production_jwt_secret
CHAT_MODE=polling
CHAT_POLLING_INTERVAL=3000
ENABLE_EMAIL_EVENTS=true
```

#### Frontend Environment Variables
```bash
# Railway provides these automatically
REACT_APP_API_URL=https://your-backend-service.railway.app

# Or with custom domain
REACT_APP_API_URL=https://api.your-domain.com

# Chat configuration
REACT_APP_CHAT_MODE=polling
REACT_APP_CHAT_POLLING_INTERVAL=2000
REACT_APP_CHAT_MAX_RETRIES=3
```

### 4. Database Configuration

Railway automatically provisions PostgreSQL with:

```yaml
Database: PostgreSQL 15
Storage: 1GB (expandable)
Connections: 100 concurrent
Backups: Automatic daily backups
SSL: Enabled by default
Connection URL: Provided as DATABASE_URL environment variable
```

#### Database Migration Strategy
```go
// Add to your main.go or init function
func runMigrations(db *sql.DB) error {
    // Railway supports running migrations on startup
    migrationFiles := []string{
        "scripts/init.sql",
        "scripts/chat_migration.sql", 
        "scripts/email_events_migration.sql",
    }
    
    for _, file := range migrationFiles {
        content, err := os.ReadFile(file)
        if err != nil {
            return err
        }
        
        if _, err := db.Exec(string(content)); err != nil {
            return err
        }
    }
    
    return nil
}
```

### 5. Service Communication

#### Internal Service URLs
Railway provides internal networking between services:

```go
// Backend can connect to database using Railway-provided DATABASE_URL
// Frontend connects to backend using public URL or internal service name

// Example backend configuration
type Config struct {
    Port        string `env:"PORT" envDefault:"8061"`
    DatabaseURL string `env:"DATABASE_URL"`
    // ... other config
}
```

#### External Service Integration
Your existing AWS services work without changes:

```go
// AWS SES configuration (unchanged)
sesConfig := &aws.Config{
    Region:      aws.String(os.Getenv("AWS_SES_REGION")),
    Credentials: credentials.NewStaticCredentials(
        os.Getenv("AWS_ACCESS_KEY_ID"),
        os.Getenv("AWS_SECRET_ACCESS_KEY"),
        "",
    ),
}

// AWS Bedrock configuration (unchanged)
bedrockConfig := &aws.Config{
    Region: aws.String(os.Getenv("BEDROCK_REGION")),
    // ... existing configuration
}
```

## Data Models

### Railway Service Configuration

#### Service Definitions
```yaml
# Railway automatically creates these services
services:
  backend:
    type: web
    source: ./backend
    build:
      dockerfile: ./backend/Dockerfile
    environment:
      - PORT=8061
      - DATABASE_URL=${{Postgres.DATABASE_URL}}
    
  frontend:
    type: web  
    source: ./frontend
    build:
      buildCommand: npm run build
      startCommand: serve -s build -l 3000
    environment:
      - REACT_APP_API_URL=${{backend.url}}
      
  database:
    type: postgresql
    version: "15"
```

### Environment Variable Management

#### Secure Configuration
```bash
# Public variables (safe to expose)
REACT_APP_API_URL=https://api.your-domain.com
REACT_APP_CHAT_MODE=polling

# Private variables (encrypted by Railway)
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...
DATABASE_URL=postgresql://...
JWT_SECRET=...
```

## Error Handling

### Railway-Specific Error Handling

#### Service Health Checks
```go
// Add health check endpoint for Railway monitoring
func healthHandler(c *gin.Context) {
    // Check database connection
    if err := db.Ping(); err != nil {
        c.JSON(500, gin.H{
            "status": "unhealthy",
            "database": "disconnected",
            "error": err.Error(),
        })
        return
    }
    
    // Check external services
    awsHealthy := checkAWSServices()
    
    c.JSON(200, gin.H{
        "status": "healthy",
        "database": "connected",
        "aws_services": awsHealthy,
        "timestamp": time.Now(),
    })
}
```

#### Graceful Shutdown
```go
// Handle Railway's shutdown signals
func main() {
    // ... setup code
    
    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        log.Println("Shutting down gracefully...")
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := server.Shutdown(ctx); err != nil {
            log.Fatal("Server forced to shutdown:", err)
        }
    }()
    
    // Start server
    server.ListenAndServe()
}
```

### Deployment Error Recovery

#### Automatic Rollback
```yaml
# Railway automatically handles failed deployments
deployment:
  strategy: rolling
  healthCheck:
    path: /health
    timeout: 30s
  rollback:
    automatic: true
    onFailure: true
```

## Testing Strategy

### Railway-Specific Testing

#### Preview Deployments
```yaml
# Railway creates preview deployments for PRs automatically
preview:
  enabled: true
  environment:
    - DATABASE_URL=${{preview-database.url}}
    - REACT_APP_API_URL=${{preview-backend.url}}
```

#### Integration Testing
```go
// Test with Railway environment
func TestRailwayIntegration(t *testing.T) {
    if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
        t.Skip("Skipping Railway integration test")
    }
    
    // Test database connection
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    require.NoError(t, err)
    defer db.Close()
    
    err = db.Ping()
    assert.NoError(t, err)
}
```

### Load Testing
```javascript
// K6 load test for Railway deployment
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 10 },
    { duration: '5m', target: 50 },
    { duration: '2m', target: 0 },
  ],
};

export default function() {
  let response = http.get('https://your-app.railway.app/health');
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
```

## Security Considerations

### Railway Security Features

#### Built-in Security
- **Automatic SSL**: Free SSL certificates with auto-renewal
- **Environment Encryption**: All environment variables encrypted at rest
- **Network Isolation**: Services isolated by default
- **DDoS Protection**: Built-in DDoS protection
- **Security Headers**: Automatic security headers

#### Additional Security Measures
```go
// Add security middleware
func securityMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        // Security headers
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        c.Next()
    })
}
```

### AWS Integration Security
```go
// Secure AWS credential management
func loadAWSConfig() (*aws.Config, error) {
    // Railway encrypts these environment variables
    accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
    secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    
    if accessKey == "" || secretKey == "" {
        return nil, errors.New("AWS credentials not configured")
    }
    
    return &aws.Config{
        Region: aws.String("us-east-1"),
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
    }, nil
}
```

## Monitoring and Observability

### Railway Built-in Monitoring

#### Metrics Dashboard
Railway provides:
- **Response Times**: P50, P95, P99 percentiles
- **Error Rates**: 4xx and 5xx error tracking
- **Throughput**: Requests per second
- **Resource Usage**: CPU and memory utilization
- **Database Metrics**: Connection count, query performance

#### Custom Metrics
```go
// Add custom metrics for Railway monitoring
func metricsMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        // Log metrics in structured format for Railway
        log.Printf("REQUEST method=%s path=%s status=%d duration=%v",
            c.Request.Method,
            c.Request.URL.Path,
            c.Writer.Status(),
            duration,
        )
    })
}
```

### Log Management
```go
// Structured logging for Railway
func setupLogging() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    
    if os.Getenv("RAILWAY_ENVIRONMENT") == "production" {
        logrus.SetLevel(logrus.InfoLevel)
    } else {
        logrus.SetLevel(logrus.DebugLevel)
    }
}
```

## Cost Optimization

### Railway Pricing Model

#### Usage-Based Pricing
```yaml
# Railway pricing (as of 2024)
starter_plan:
  cost: $5/month
  includes:
    - 512MB RAM per service
    - 1GB disk per service
    - 100GB bandwidth
    - Custom domains
    - SSL certificates

pro_plan:
  cost: $20/month
  includes:
    - 8GB RAM per service
    - 100GB disk per service
    - 1TB bandwidth
    - Priority support
    - Advanced metrics
```

#### Cost Optimization Strategies
```go
// Optimize resource usage
func optimizeForRailway() {
    // Use connection pooling
    db.SetMaxOpenConns(10)  // Railway starter: keep connections low
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)
    
    // Implement caching to reduce database load
    cache := make(map[string]interface{})
    
    // Use efficient JSON serialization
    gin.SetMode(gin.ReleaseMode)
}
```

### Monitoring Costs
```bash
# Railway provides cost tracking in dashboard
# Set up alerts for usage thresholds
RAILWAY_USAGE_ALERT_THRESHOLD=80  # Alert at 80% of plan limits
```

This design provides a comprehensive approach to deploying your Cloud Consulting Platform on Railway with minimal configuration while maintaining production-ready features and your existing AWS integrations.