#!/bin/bash

# Application-Only Build Script
# This script builds only the main application without running any tests
# Requirement: 4.4 - Script for building application only

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Change to backend directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$BACKEND_DIR"

print_status $BLUE "=== Application-Only Build ==="
print_status $BLUE "Building Cloud Consulting Backend (Application Only)..."
echo

# Check Go version
GO_VERSION=$(go version 2>/dev/null || echo "Go not found")
print_status $BLUE "Go Version: $GO_VERSION"

# Clean up any existing binaries
print_status $BLUE "Cleaning up existing binaries..."
rm -f server main backend-server bin/server

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the main application only
print_status $BLUE "Compiling main application..."
echo "Command: go build -o bin/server ./cmd/server"

if go build -o bin/server ./cmd/server; then
    print_status $GREEN "✅ Application build successful!"
    
    # Show binary info
    if [ -f "bin/server" ]; then
        size=$(ls -lh bin/server | awk '{print $5}')
        print_status $BLUE "Binary size: $size"
        print_status $BLUE "Binary location: ./bin/server"
        
        # Test if binary is executable
        if [ -x "bin/server" ]; then
            print_status $GREEN "Binary is executable"
        else
            print_status $YELLOW "Warning: Binary may not be executable"
        fi
    fi
    
    echo
    print_status $GREEN "🚀 Application ready to run!"
    print_status $BLUE "Usage:"
    print_status $BLUE "  ./bin/server                    # Run with default config"
    print_status $BLUE "  ./bin/server --help             # Show help"
    print_status $BLUE "  ENV_VAR=value ./bin/server      # Run with environment variables"
    
else
    print_status $RED "❌ Application build failed!"
    print_status $RED "Please check the error messages above and fix any compilation issues."
    exit 1
fi

echo
print_status $BLUE "Build completed in: $BACKEND_DIR"
print_status $BLUE "Binary location: $BACKEND_DIR/bin/server"