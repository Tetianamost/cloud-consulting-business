#!/bin/bash

# Cloud Consulting Backend - Stop Local Development

echo "🛑 Stopping Cloud Consulting Backend and Frontend..."
echo "=================================================="

# Stop and remove containers
docker-compose down --remove-orphans

# Optional: Remove volumes (uncomment if you want to clean up data)
# docker-compose down --volumes

echo "✅ Services stopped successfully!"
echo ""
echo "To start again, run: ./start-local.sh"