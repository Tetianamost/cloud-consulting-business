# Cloud Consulting Backend

A comprehensive Go backend system for processing and categorizing cloud consulting service inquiries with AI-powered report generation and admin dashboard.

## Features

- 🚀 **Service Inquiry Processing**: Handle four main service types (Assessment, Migration, Optimization, Architecture Review)
- 🤖 **AI Report Generation**: Automatically generate draft reports using Amazon Bedrock AI
- 💬 **Real-time Chat System**: Polling-based chat with AI assistant for reliable communication
- 📧 **Email Notifications**: Professional email notifications using AWS SES
- 👨‍💼 **Admin Dashboard**: React-based admin interface with comprehensive management tools
- 🔐 **Secure Authentication**: JWT-based admin authentication with session management
- 📊 **Advanced Analytics**: System metrics, performance monitoring, and quality assurance
- 🔧 **Automation & Integration**: Proactive recommendations and third-party integrations
- 🎯 **Meeting Preparation**: AI-powered client meeting preparation and competitive analysis
- 📈 **Performance Optimization**: Intelligent caching, load balancing, and resource optimization
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

**Admin Dashboard Features:**
- **AI Consultant Assistant**: Advanced AI-powered chat interface with:
  - 8 pre-defined quick actions (Cost Estimate, Security Review, Best Practices, etc.)
  - Context management with client name and meeting type
  - Fullscreen mode for focused conversations
  - Real-time connection status monitoring
  - Debounced input for optimal performance
  - Session persistence across page reloads
- Real-time chat with AI assistant (polling-based)
- Inquiry and report management with AI-generated reports
- System metrics and performance monitoring
- Email delivery tracking with AWS SES integration
- Meeting preparation tools with competitive analysis
- Quality assurance dashboard with peer review system
- Integration management for third-party services
- Performance optimization tools with intelligent caching

## API Endpoints

### Public Endpoints
- `POST /api/v1/inquiries` - Create new inquiry
- `GET /api/v1/inquiries/{id}` - Get inquiry details
- `GET /api/v1/config/services` - Get available service types
- `GET /health` - Health check

### Admin Endpoints (Protected)
- `POST /api/v1/auth/login` - Admin login with JWT token generation
- `GET /api/v1/admin/inquiries` - List all inquiries with filtering and pagination
- `GET /api/v1/admin/metrics` - System metrics and performance data
- `GET /api/v1/admin/email-status/{id}` - Email delivery status tracking
- `GET /api/v1/admin/reports/{id}/download/{format}` - Download reports (PDF/HTML)

### Chat Endpoints (Protected)
- `POST /api/v1/admin/chat/sessions` - Create new chat session
- `GET /api/v1/admin/chat/sessions` - List chat sessions with metadata
- `GET /api/v1/admin/chat/sessions/{id}/history` - Get chat history with pagination
- `GET /api/v1/admin/chat/metrics` - Chat system performance metrics
- `POST /api/v1/admin/chat/polling` - HTTP polling for chat messages
- `POST /api/v1/admin/chat/send` - Send message with context and quick actions

### AI Consultant Endpoints (Protected)
- `POST /api/v1/admin/simple-chat/messages` - Send message to AI assistant with context
- `GET /api/v1/admin/simple-chat/messages` - Retrieve chat messages by session ID

### Advanced Admin Features (Protected)
- `GET /api/v1/admin/meeting-prep/*` - AI-powered meeting preparation tools
- `GET /api/v1/admin/quality-assurance/*` - Quality assurance and peer review system
- `GET /api/v1/admin/integrations` - Third-party integration management
- `GET /api/v1/admin/cost-analysis` - Cost analysis and optimization recommendations
- `GET /api/v1/admin/performance/*` - Performance optimization and monitoring

## Development

### Build Scripts

The project includes convenient build scripts for different scenarios:

#### Application-Only Build
```bash
# Build just the main application (no tests)
./backend/scripts/build_app_only.sh
```

#### Test Categories
```bash
# Run specific test categories
./backend/scripts/run_test_categories.sh [CATEGORY] [OPTIONS]

# Examples:
./backend/scripts/run_test_categories.sh unit              # Unit tests only
./backend/scripts/run_test_categories.sh integration -v   # Integration tests (verbose)
./backend/scripts/run_test_categories.sh email --coverage # Email tests with coverage
./backend/scripts/run_test_categories.sh all              # All test categories
./backend/scripts/run_test_categories.sh list             # List available tests
```

Available test categories:
- **unit**: Unit tests (`*_test.go` files in `internal/`)
- **integration**: Integration tests (`testing/integration/`)
- **email**: Email system tests (`testing/email/`)
- **performance**: Performance tests (`testing/performance/`)
- **all**: Run all test categories

### File Organization

The project follows a clean separation between application code and test utilities:

```
backend/
├── cmd/server/              # Main application entry point
├── internal/               # Application code
│   ├── config/            # Configuration management
│   ├── domain/            # Domain models and business logic
│   ├── handlers/          # HTTP request handlers
│   ├── interfaces/        # Service interfaces and contracts
│   ├── repositories/      # Data access layer
│   ├── services/          # Business logic services
│   ├── server/            # Server setup and middleware
│   └── storage/           # Storage implementations
├── testing/               # Organized test utilities (NEW)
│   ├── integration/       # Integration test executables
│   ├── email/            # Email system test executables
│   ├── performance/      # Performance test executables
│   └── shared/           # Shared test utilities and mocks
├── scripts/              # Build and deployment scripts
└── docs/                 # API and system documentation
```

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

# Authentication
ADMIN_USERNAME=admin
ADMIN_PASSWORD=cloudadmin
JWT_EXPIRATION=24h

# AWS Bedrock (for AI report generan)
AWS_BEARER_TOKEN_BEDROCK=your-bedrock-token
AWS-key
et-key
AWS_REGION=us-e

# AWS SES (for email notifications)
SES_SENDER_EMAIL=noreply@yourdomain.com
SES_REGION=us-east-1

# Database Configuration (for produ
DATABASE_URL=postgreng
REDIS_URL=redis://localhost:6379

guration
REACT_APP_API_UR1
61

# Performance and Monitoring
ENABLE_METRICS=true
ENABLE_CHAT_LOGGING=true
L=3600
```CACHE_TTlhost:80locas://=wT_APP_WS_URLREAC806calhost:L=http://lod Confironten# Fud_consulti2/cloalhost:543@locer:passwordsql://uson)ctiast-1ecrS_KEY=your-sT_ACCES_SECREAWScess=your-acIDCCESS_KEY__A

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: React with TypeScript and Redux Toolkit
- **Real-time Communication**: HTTP polling for reliable chat functionality
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

### Quick Fixes

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

### Build Issues

For comprehensive build and compilation troubleshooting, see: **[TROUBLESHOOTING_BUILD_ISSUES.md](./TROUBLESHOOTING_BUILD_ISSUES.md)**

Common build issues include:
- Multiple `package main` redeclaration errors
- Go module dependency conflicts
- TypeScript compilation errors
- Test execution problems
- Docker build context issues

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

### Build Verification

```bash
# Test application build
./backend/scripts/build_app_only.sh

# Test specific test categories
./backend/scripts/run_test_categories.sh list
./backend/scripts/run_test_categories.sh unit

# Test frontend build
cd frontend && npm run build
```

### Testing

For comprehensive testing instructions and guidelines, see: **[TESTING_GUIDE.md](./TESTING_GUIDE.md)**

Quick test commands:
```bash
# Run all tests
./backend/scripts/run_test_categories.sh all

# Run specific test categories
./backend/scripts/run_test_categories.sh unit
./backend/scripts/run_test_categories.sh integration
./backend/scripts/run_test_categories.sh email
./backend/scripts/run_test_categories.sh performance

# List available tests
./backend/scripts/run_test_categories.sh list
```

## Stopping the Application

```bash
./stop-local.sh
# or
docker-compose down
```

## Recent Updates

### Code Quality Improvements
- **Authentication System**: Cleaned up debug logging in AuthContext for production readiness
- **Error Handling**: Improved error logging while removing verbose debug output
- **Console Output**: Cleaner browser console experience for end users

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally with `./start-local.sh`
5. Submit a pull request

## License

This project is licensed under the MIT License.