#!/bin/bash

# Build script for completely offline Docker deployment

echo "🔨 Building Cloud Consulting Backend for Offline Docker Deployment"
echo "================================================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ to build the backend."
    echo "   Download from: https://golang.org/dl/"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ to build the frontend."
    echo "   Download from: https://nodejs.org/"
    exit 1
fi

echo "📦 Building Backend Binary..."
cd backend

# Build the Go binary for Linux (for Docker)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server/main.go

if [ $? -ne 0 ]; then
    echo "❌ Backend build failed"
    exit 1
fi

echo "✅ Backend binary built successfully"
cd ..

echo "📦 Building Frontend..."
cd frontend

# Install dependencies and build
npm install
npm run build

if [ $? -ne 0 ]; then
    echo "❌ Frontend build failed"
    exit 1
fi

echo "✅ Frontend built successfully"
cd ..

echo "🐳 Creating minimal Docker images..."

# Create a simple docker-compose for the built artifacts
cat > docker-compose.offline.yml << 'EOF'
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.offline
    ports:
      - "8061:8061"
    environment:
      - PORT=8061
      - LOG_LEVEL=4
      - GIN_MODE=debug
      - CORS_ALLOWED_ORIGINS=http://localhost:3007
      - JWT_SECRET=cloud-consulting-demo-secret
    networks:
      - consulting-network
    restart: unless-stopped

  frontend:
    image: nginx:alpine
    ports:
      - "3006:80"
    volumes:
      - ./frontend/build:/usr/share/nginx/html:ro
      - ./frontend/nginx.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      - backend
    networks:
      - consulting-network
    restart: unless-stopped

networks:
  consulting-network:
    driver: bridge
EOF

echo "✅ Docker configuration created"
echo ""
echo "🚀 To start the application:"
echo "   docker-compose -f docker-compose.offline.yml up"
echo ""
echo "📱 Frontend: http://localhost:3006"
echo "🔧 Backend: http://localhost:8061"
echo "🔐 Admin: http://localhost:3006/admin/login (admin/cloudadmin)"