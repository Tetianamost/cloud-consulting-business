# Cloud Consulting Business Platform

A full-stack application for managing cloud consulting services with AI-powered report generation using Amazon Bedrock.

## 🚀 Quick Start with Docker

### Prerequisites

- Docker and Docker Compose installed
- AWS account with Bedrock access
- Amazon Bedrock API key

### 1. Clone and Setup

```bash
git clone <your-repo>
cd cloud-consulting-business

# Copy environment configuration
cp .env.example .env
```

### 2. Configure AWS Credentials

Edit `.env` file with your AWS credentials:

```bash
# Required for Bedrock AI integration
AWS_BEARER_TOKEN_BEDROCK=your_bedrock_api_key_here
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=us-east-1

# AWS_SESSION_TOKEN is optional - only needed for temporary credentials
```

> 📖 **Need help with AWS setup?** See our detailed [AWS Setup Guide](docs/AWS_SETUP.md) for step-by-step instructions on getting your credentials.

### 3. Start the Application

```bash
# Start backend and frontend
docker-compose up -d backend frontend

# Or use the helper script
./scripts/docker-dev.sh up
```

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8061
- **Health Check**: http://localhost:8061/health

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React Frontend│    │   Go Backend    │    │  Amazon Bedrock │
│   (Port 3000)   │◄──►│   (Port 8061)   │◄──►│   AI Service    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Development Commands

### Using Docker Compose Helper Script

```bash
# Start services
./scripts/docker-dev.sh up              # Backend + Frontend
./scripts/docker-dev.sh up-backend      # Backend only
./scripts/docker-dev.sh up-full         # All services + DB + Cache

# View logs
./scripts/docker-dev.sh logs            # All services
./scripts/docker-dev.sh logs-backend    # Backend only
./scripts/docker-dev.sh logs-frontend   # Frontend only

# Stop and cleanup
./scripts/docker-dev.sh down            # Stop services
./scripts/docker-dev.sh clean           # Stop + remove volumes

# Test services
./scripts/docker-dev.sh test            # Health checks

# Build services
./scripts/docker-dev.sh build           # Rebuild containers
```

### Direct Docker Compose Commands

```bash
# Start all services
docker-compose up -d

# Start with optional services (database, cache, monitoring)
docker-compose --profile database --profile cache --profile monitoring up -d

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop services
docker-compose down

# Rebuild and start
docker-compose up --build -d

# Clean up everything
docker-compose down -v --remove-orphans
```

## 🧪 Testing the API

### Create an inquiry (triggers AI report generation):

```bash
curl -X POST http://localhost:8061/api/v1/inquiries \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Tech Corp",
    "services": ["assessment"],
    "message": "We need help assessing our current AWS infrastructure for cost optimization."
  }'
```

### Get the AI-generated report:

```bash
curl http://localhost:8061/api/v1/inquiries/{inquiry-id}/report
```

### List all inquiries:

```bash
curl http://localhost:8061/api/v1/inquiries
```

## 📁 Project Structure

```
cloud-consulting-business/
├── backend/                 # Go backend service
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Internal packages
│   │   ├── config/        # Configuration management
│   │   ├── domain/        # Domain models
│   │   ├── handlers/      # HTTP handlers
│   │   ├── services/      # Business logic (including Bedrock)
│   │   ├── storage/       # Data storage
│   │   └── server/        # Server setup
│   ├── Dockerfile         # Backend container
│   └── README.md          # Backend documentation
├── frontend/               # React frontend
│   ├── src/               # Source code
│   ├── Dockerfile         # Frontend container
│   └── nginx.conf         # Nginx configuration
├── scripts/               # Helper scripts
│   └── docker-dev.sh      # Docker development helper
├── docker-compose.yml     # Multi-service orchestration
├── .env.example          # Environment template
└── README.md             # This file
```

## 🔧 Features

### Backend (Go)
- ✅ RESTful API for inquiry management
- ✅ AI-powered report generation with Amazon Bedrock
- ✅ Graceful error handling and fallbacks
- ✅ Structured logging
- ✅ Health checks and monitoring
- ✅ CORS support for frontend integration

### Frontend (React)
- ✅ Modern React with TypeScript
- ✅ Responsive design with Tailwind CSS
- ✅ Form handling with validation
- ✅ API integration with backend
- ✅ Production-ready Nginx configuration

### AI Integration
- ✅ Amazon Bedrock Nova model integration
- ✅ Structured prompt engineering
- ✅ Professional report generation
- ✅ Error handling and fallbacks
- ✅ Configurable model parameters

## 🔒 Security

- Environment-based configuration
- Non-root container users
- HTTPS-ready setup
- Security headers in Nginx
- Input validation and sanitization

## 📊 Optional Services

Enable additional services with profiles:

```bash
# Database (PostgreSQL)
docker-compose --profile database up -d

# Cache (Redis)
docker-compose --profile cache up -d

# Monitoring (Prometheus + Grafana)
docker-compose --profile monitoring up -d
```

## 🚀 Production Deployment

For production deployment:

1. Use production environment variables
2. Enable HTTPS with SSL certificates
3. Use managed databases (RDS, ElastiCache)
4. Implement proper logging and monitoring
5. Set up CI/CD pipelines
6. Configure auto-scaling

## 📝 Environment Variables

### Required
- `AWS_BEARER_TOKEN_BEDROCK` - Your Bedrock API key
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key

### Optional
- `AWS_REGION` - AWS region (default: us-east-1)
- `BEDROCK_MODEL_ID` - Bedrock model (default: amazon.nova-lite-v1:0)
- `PORT` - Backend port (default: 8061)
- `LOG_LEVEL` - Logging level (default: 4)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with Docker Compose
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
1. Check the logs: `./scripts/docker-dev.sh logs`
2. Run health checks: `./scripts/docker-dev.sh test`
3. Review the documentation in `/backend/README.md`
4. Open an issue on GitHub