#!/bin/bash

# Docker Compose Development Script
# This script helps manage the development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists
if [ ! -f .env ]; then
    print_warning ".env file not found. Creating from .env.example..."
    if [ -f .env.example ]; then
        cp .env.example .env
        print_warning "Please edit .env file with your AWS credentials before running the services"
        exit 1
    else
        print_error ".env.example file not found. Please create .env file manually."
        exit 1
    fi
fi

# Function to show help
show_help() {
    echo "Docker Compose Development Helper"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up          Start all services (backend + frontend)"
    echo "  up-backend  Start only backend service"
    echo "  up-full     Start all services including database and cache"
    echo "  down        Stop all services"
    echo "  logs        Show logs for all services"
    echo "  logs-backend Show logs for backend only"
    echo "  logs-frontend Show logs for frontend only"
    echo "  build       Build all services"
    echo "  clean       Stop services and remove volumes"
    echo "  test        Run health checks"
    echo "  help        Show this help message"
}

# Function to test services
test_services() {
    print_status "Testing backend health..."
    if curl -f http://localhost:8061/health > /dev/null 2>&1; then
        print_status "✓ Backend is healthy"
    else
        print_error "✗ Backend health check failed"
    fi

    print_status "Testing frontend..."
    if curl -f http://localhost:3000 > /dev/null 2>&1; then
        print_status "✓ Frontend is accessible"
    else
        print_error "✗ Frontend is not accessible"
    fi

    print_status "Testing API endpoint..."
    if curl -f http://localhost:8061/api/v1/config/services > /dev/null 2>&1; then
        print_status "✓ API endpoints are working"
    else
        print_error "✗ API endpoints are not working"
    fi
}

# Main command handling
case "${1:-help}" in
    up)
        print_status "Starting backend and frontend services..."
        docker-compose up -d backend frontend
        print_status "Services started. Backend: http://localhost:8061, Frontend: http://localhost:3000"
        ;;
    up-backend)
        print_status "Starting backend service only..."
        docker-compose up -d backend
        print_status "Backend started at http://localhost:8061"
        ;;
    up-full)
        print_status "Starting all services including database and cache..."
        docker-compose --profile database --profile cache up -d
        print_status "All services started"
        ;;
    down)
        print_status "Stopping all services..."
        docker-compose down
        print_status "Services stopped"
        ;;
    logs)
        docker-compose logs -f
        ;;
    logs-backend)
        docker-compose logs -f backend
        ;;
    logs-frontend)
        docker-compose logs -f frontend
        ;;
    build)
        print_status "Building all services..."
        docker-compose build
        print_status "Build complete"
        ;;
    clean)
        print_status "Stopping services and cleaning up..."
        docker-compose down -v --remove-orphans
        docker system prune -f
        print_status "Cleanup complete"
        ;;
    test)
        test_services
        ;;
    help)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac