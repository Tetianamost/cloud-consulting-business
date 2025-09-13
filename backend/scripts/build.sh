#!/bin/bash

# Simple Build Script for Cloud Consulting Backend
# This script builds the main application with basic verification

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Change to backend directory
cd "$(dirname "$0")/.."

print_status $BLUE "Building Cloud Consulting Backend..."

# Clean up any existing binaries
rm -f server main backend-server

# Build the main application
print_status $BLUE "Compiling main application..."
if go build -o server ./cmd/server; then
    print_status $GREEN "✅ Build successful!"
    
    # Show binary info
    size=$(ls -lh server | awk '{print $5}')
    print_status $BLUE "Binary size: $size"
    print_status $BLUE "Binary location: ./server"
    
    echo
    print_status $GREEN "🚀 Ready to run with: ./server"
else
    print_status $RED "❌ Build failed!"
    exit 1
fi