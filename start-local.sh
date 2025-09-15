#!/bin/bash

# Cloud Consulting Backend - Local Development Startup Script

echo "🚀 Starting Cloud Consulting Backend and Frontend..."
echo "=================================================="

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    echo "📋 Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  No .env file found. Using default configuration."
fi

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose down --remove-orphans

# Build and start the services
echo "🔨 Building and starting services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check service health
echo "🔍 Checking service health..."

# Check backend health
echo "Checking backend health..."
for i in {1..30}; do
    if curl -f http://localhost:8061/health > /dev/null 2>&1; then
        echo "✅ Backend is healthy!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Backend health check failed after 30 attempts"
        echo "Backend logs:"
        docker-compose logs backend
        exit 1
    fi
    sleep 2
done

# Check frontend health
echo "Checking frontend health..."
for i in {1..30}; do
    if curl -f http://localhost:3006 > /dev/null 2>&1; then
        echo "✅ Frontend is healthy!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Frontend health check failed after 30 attempts"
        echo "Frontend logs:"
        docker-compose logs frontend
        exit 1
    fi
    sleep 2
done

echo ""
echo "🎉 Services are running successfully!"
echo "=================================================="
echo "📱 Frontend: http://localhost:3006"
echo "🔧 Backend API: http://localhost:8061"
echo "🏥 Health Check: http://localhost:8061/health"
echo "🔐 Admin Login: http://localhost:3006/admin/login"
echo ""
echo "Admin Credentials:"
echo "  Username: admin"
echo "  Password: cloudadmin"
echo ""
echo "📊 To view logs:"
echo "  docker-compose logs -f backend"
echo "  docker-compose logs -f frontend"
echo ""
echo "🛑 To stop services:"
echo "  docker-compose down"
echo ""
echo "Happy coding! 🚀"