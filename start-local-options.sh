#!/bin/bash

# Cloud Consulting Backend - Local Development Options

echo "🚀 Cloud Consulting Backend - Local Development"
echo "=============================================="
echo ""
echo "Choose your development setup:"
echo ""
echo "1) 🐳 Docker with local builds (no registry pulls)"
echo "2) 💻 Native development (Go + Node.js directly)"
echo "3) 🐳 Docker with registry images (original - requires internet)"
echo ""
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo ""
        echo "🐳 Starting with Docker (local builds)..."
        echo "This will build everything locally without pulling from registries."
        echo ""
        
        # Check if Docker is running
        if ! docker info > /dev/null 2>&1; then
            echo "❌ Docker is not running. Please start Docker and try again."
            exit 1
        fi
        
        # Load environment variables
        if [ -f .env ]; then
            echo "📋 Loading environment variables from .env file..."
            export $(cat .env | grep -v '^#' | xargs)
        fi
        
        # Clean up any existing containers
        echo "🧹 Cleaning up existing containers..."
        docker-compose -f docker-compose.local.yml down --remove-orphans
        
        # Build and start the services
        echo "🔨 Building and starting services (this may take a while for first build)..."
        docker-compose -f docker-compose.local.yml up --build
        ;;
        
    2)
        echo ""
        echo "💻 Starting native development..."
        echo "This requires Go and Node.js to be installed locally."
        echo ""
        ./start-local-dev.sh
        ;;
        
    3)
        echo ""
        echo "🐳 Starting with Docker (registry images)..."
        echo "This requires internet connection to pull images."
        echo ""
        ./start-local.sh
        ;;
        
    *)
        echo "❌ Invalid choice. Please run the script again and choose 1, 2, or 3."
        exit 1
        ;;
esac