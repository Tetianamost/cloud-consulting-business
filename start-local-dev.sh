#!/bin/bash

# Cloud Consulting Backend - Local Development (No Docker)

echo "🚀 Starting Cloud Consulting Backend and Frontend (Local Development)..."
echo "=================================================================="
# // kill this frontend port as well - lsof -ti:3007 | xargs kill -9 || true 


# Function to kill processes on specific ports
kill_port() {
    local port=$1
    echo "🔍 Checking for processes on port $port..."
    local pids=$(lsof -ti:$port 2>/dev/null)
    if [ ! -z "$pids" ]; then
        echo "🔪 Killing processes on port $port: $pids"
        echo "$pids" | xargs kill -9 2>/dev/null || true
        sleep 2
        # Double check if processes are still running
        local remaining=$(lsof -ti:$port 2>/dev/null)
        if [ ! -z "$remaining" ]; then
            echo "⚠️  Some processes still running on port $port: $remaining"
            echo "$remaining" | xargs kill -9 2>/dev/null || true
            sleep 1
        fi
    else
        echo "✅ Port $port is free"
    fi
}

# Kill any existing processes on our ports
echo "🧹 Cleaning up existing processes..."
kill_port 8061
kill_port 3006
kill_port 3007

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ and try again."
    echo "   Download from: https://golang.org/dl/"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18+ and try again."
    echo "   Download from: https://nodejs.org/"
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    echo "📋 Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "⚠️  No .env file found. Using default configuration."
fi

# Set default environment variables if not set
export PORT=${PORT:-8061}
export REACT_APP_API_URL=${REACT_APP_API_URL:-http://localhost:8061}
export JWT_SECRET=${JWT_SECRET:-cloud-consulting-demo-secret}
export GIN_MODE=${GIN_MODE:-debug}
export LOG_LEVEL=${LOG_LEVEL:-4}
export FRONTEND_PORT=3007
export CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS:-"http://localhost:3007, http://localhost:3006, http://localhost:3000"}

echo "🔧 Configuration:"
echo "  Backend Port: $PORT"
echo "  Frontend URL: http://localhost:$FRONTEND_PORT"
echo "  API URL: $REACT_APP_API_URL"

# Function to kill background processes on exit
cleanup() {
    echo ""
    echo "🛑 Stopping services..."
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
    fi
    
    # Also kill any processes still running on our ports
    kill_port 8061
    kill_port 3006
    kill_port 3007
    
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Start Backend
echo "🔨 Starting Backend..."

# Check if backend directory exists
if [ ! -d "backend" ]; then
    echo "❌ Backend directory not found. Make sure you're running this script from the project root."
    exit 1
fi

cd backend

# Download Go dependencies
echo "📦 Installing Go dependencies..."
if ! go mod download; then
    echo "❌ Failed to download Go dependencies"
    cd ..
    exit 1
fi

# Start the backend in the background
echo "🚀 Starting Go server on port $PORT..."
go run ./cmd/server/main.go &
BACKEND_PID=$!

# Wait a moment and check if the process is still running
sleep 2
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "❌ Backend process died immediately. Check for compilation errors."
    cd ..
    exit 1
fi

cd ..

# Wait a moment for backend to start (increased timeout for full initialization)
sleep 6

# Check if backend is running
if ! curl -f http://localhost:$PORT/health > /dev/null 2>&1; then
    echo "❌ Backend failed to start. Check the logs above."
    cleanup
fi

echo "✅ Backend is running on http://localhost:$PORT"

# Start Frontend
echo "🔨 Starting Frontend..."

# Check if frontend directory exists
if [ ! -d "frontend" ]; then
    echo "❌ Frontend directory not found. Make sure you're running this script from the project root."
    cleanup
fi

cd frontend

# Install npm dependencies
echo "� Inastalling npm dependencies..."
if ! npm install; then
    echo "❌ Failed to install npm dependencies"
    cd ..
    cleanup
fi

# Set the port for the frontend
export PORT=$FRONTEND_PORT

# Start the frontend in the background
echo "🚀 Starting React development server on port $FRONTEND_PORT..."
npm start &
FRONTEND_PID=$!

cd ..

# Wait for frontend to start
echo "⏳ Waiting for frontend to start..."
sleep 10

# Check if frontend is running
for i in {1..30}; do
    if curl -f http://localhost:$FRONTEND_PORT > /dev/null 2>&1; then
        echo "✅ Frontend is running!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Frontend failed to start after 30 attempts"
        cleanup
    fi
    sleep 2
done

echo ""
echo "🎉 Services are running successfully!"
echo "=================================================="
echo "📱 Frontend: http://localhost:$FRONTEND_PORT"
echo "🔧 Backend API: http://localhost:$PORT"
echo "🏥 Health Check: http://localhost:$PORT/health"
echo "🔐 Admin Login: http://localhost:$FRONTEND_PORT/admin/login"
echo ""
echo "Admin Credentials:"
echo "  Username: admin"
echo "  Password: cloudadmin"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for user to stop
wait