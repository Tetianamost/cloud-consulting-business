#!/bin/bash

# Build script similar to quota-service approach
set -e

PROJECT_NAME="cloud-consulting-business"
COMMIT_HASH=$(git rev-parse --short=8 HEAD)
DATE=$(date -u +%Y%m%d)
VERSION="0.0.1-${DATE}-${COMMIT_HASH}"

echo "Building Go binary for amd64..."

# Create build directory
mkdir -p .build

# Build Go binary using Docker container with amd64 platform
docker run --rm \
    --platform linux/amd64 \
    -v $(pwd):/workspace \
    -w /workspace/backend \
    -e CGO_ENABLED=0 \
    -e GOOS=linux \
    -e GOARCH=amd64 \
    public.ecr.aws/docker/library/golang:1.24-alpine \
    go build -ldflags "-X main.Version=${VERSION}" -a -o ../.build/server cmd/server/main.go

echo "Building frontend..."

# Build frontend using Docker container
docker run --rm \
    --platform linux/amd64 \
    -v $(pwd):/workspace \
    -w /workspace/frontend \
    public.ecr.aws/docker/library/node:18-alpine \
    sh -c "npm install && npm run build"

echo "Build completed successfully!"
