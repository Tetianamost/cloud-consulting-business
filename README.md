# Cloud Consulting Backend

A comprehensive Go backend system for processing and categorizing cloud consulting service inquiries with AI-powered report generation and admin dashboard.

## Features

- 🚀 **Service Inquiry Processing**: Handle four main service types (Assessment, Migration, Optimization, Architecture Review)
- 🤖 **AI Report Generation**: Automatically generate draft reports using Amazon Bedrock AI
- 📧 **Email Notifications**: Professional email notifications using AWS SES
- 👨‍💼 **Admin Dashboard**: React-based admin interface for monitoring and management
- 🔐 **Authentication**: JWT-based admin authentication
- 📊 **Metrics & Monitoring**: System metrics and email delivery tracking
- 🐳 **Docker Support**: Full containerization with Docker Compose

## Quick Start

### Prerequisites

Choose one of the following setups based on your environment:

**Option 1: Native Development**
- Go 1.21+
- Node.js 18+
- AWS credentials (optional for AI and email features)

**Option 2: Docker (Local Builds)**
- Docker and Docker Compose
- No internet required after initial setup

**Option 3: Docker (Registry Images)**
- Docker and Docker Compose
- Internet connection required

### 1. Clone and Setup

```bash
git clone <repository-url>
cd cloud-consulting-backend
```

### 2. Configure Environment

The `.env` file is already configured with demo values. Update with your AWS credentials if you want AI and email features:

```bash
# Edit .env with your AWS credentials if you want AI and email features
nano .env
```

### 3. Choose Your Development Method

#### Option A: Interactive Setup (Recommended)
```bash
./start-local-options.sh
```

This will give you a menu to choose between:
1. Docker with local builds (no registry pulls)
2. Native development (Go + Node.js directly)
3. Docker with registry images (requires internet)

#### Option B: Direct Commands

**Native Development (No Docker):**
```bash
./start-local-dev.sh
```

**Docker with Local Builds:**
```bash
docker-compose -f docker-compose.local.yml up --build
```

**Docker with Registry Images:**
```bash
./start-local.sh
```

**Completely Offline Docker:**
```bash
./build-offline.sh
docker-compose -f docker-compose.offline.yml up
```

### 4. Access the Application

All methods will make the application available at:
- **Frontend**: http://localhost:3006
- **Backend API**: http://localhost:8061
- **Health Check**: http://localhost:8061/health
- **Admin Login**: http://localhost:3006/admin/login

### 4. Access the Application

- **Frontend**: http://localhost:3006
- **Backend API**: http://localhost:8061
- **Health Check**: http://localhost:8061/health
- **Admin Login**: http://localhost:3006/admin/login

**Admin Credentials:**
- Username: `admin`
- Password: `cloudadmin`

## API Endpoints

### Public Endpoints
- `POST /api/v1/inquiries` - Create new inquiry
- `GET /api/v1/inquiries/{id}` - Get inquiry details
- `GET /api/v1/config/services` - Get available service types
- `GET /health` - Health check

### Admin Endpoints (Protected)
- `POST /api/v1/auth/login` - Admin login
- `GET /api/v1/admin/inquiries` - List all inquiries
- `GET /api/v1/admin/metrics` - System metrics
- `GET /api/v1/admin/email-status/{id}` - Email delivery status
- `GET /api/v1/admin/reports/{id}/download/{format}` - Download reports

## Development

### Manual Docker Commands

```bash
# Build and start services
docker-compose up --build -d

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop services
docker-compose down
```

### Environment Variables

Key environment variables in `.env`:

```bash
# Backend Configuration
PORT=8061
JWT_SECRET=cloud-consulting-demo-secret

# AWS Bedrock (for AI report generation)
AWS_BEARER_TOKEN_BEDROCK=your-bedrock-token
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key

# AWS SES (for email notifications)
SES_SENDER_EMAIL=noreply@yourdomain.com
```

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: React with TypeScript
- **Storage**: In-memory (for demo) with plans for PostgreSQL
- **AI**: Amazon Bedrock Nova model
- **Email**: AWS SES
- **Authentication**: JWT tokens
- **Containerization**: Docker with multi-stage builds

## Service Types

1. **Assessment** - Cloud readiness assessment and migration planning
2. **Migration** - End-to-end cloud migration services
3. **Optimization** - Cloud cost optimization and performance tuning
4. **Architecture Review** - Cloud architecture review and best practices

## Troubleshooting

### Common Issues

1. **Port conflicts**: Make sure ports 3006 and 8061 are available
2. **Docker issues**: Ensure Docker is running and you have sufficient resources
3. **AWS credentials**: AI and email features require valid AWS credentials
4. **Network/Registry issues**: If you get "failed to resolve source metadata" errors:
   - Use `./start-local-options.sh` and choose option 1 (local builds)
   - Or use native development with option 2
   - Or build completely offline with `./build-offline.sh`
5. **DNS resolution issues**: If Docker can't reach registries:
   - Check your DNS settings
   - Try using local builds instead of registry pulls
   - Use the offline build option

### Logs

```bash
# Backend logs
docker-compose logs backend

# Frontend logs
docker-compose logs frontend

# All logs
docker-compose logs
```

### Health Checks

```bash
# Backend health
curl http://localhost:8061/health

# Frontend health
curl http://localhost:3006
```

## Stopping the Application

```bash
./stop-local.sh
# or
docker-compose down
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally with `./start-local.sh`
5. Submit a pull request

## License

This project is licensed under the MIT License.