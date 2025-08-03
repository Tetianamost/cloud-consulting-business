#!/bin/bash

# AI Consultant Live Chat System Deployment Script
# This script deploys the chat system using Docker Compose or Kubernetes

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENVIRONMENT="${ENVIRONMENT:-development}"
DEPLOYMENT_TYPE="${DEPLOYMENT_TYPE:-docker}"
NAMESPACE="${NAMESPACE:-ai-consultant-chat}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    cat << EOF
AI Consultant Live Chat System Deployment Script

Usage: $0 [OPTIONS]

Options:
    -e, --environment    Environment to deploy (development, staging, production)
    -t, --type          Deployment type (docker, kubernetes)
    -n, --namespace     Kubernetes namespace (default: ai-consultant-chat)
    -h, --help          Show this help message

Environment Variables:
    ENVIRONMENT         Deployment environment
    DEPLOYMENT_TYPE     Type of deployment (docker or kubernetes)
    NAMESPACE          Kubernetes namespace

Examples:
    # Deploy to development using Docker Compose
    $0 -e development -t docker

    # Deploy to production using Kubernetes
    $0 -e production -t kubernetes -n production

    # Deploy using environment variables
    ENVIRONMENT=staging DEPLOYMENT_TYPE=kubernetes $0
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -t|--type)
            DEPLOYMENT_TYPE="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(development|staging|production)$ ]]; then
    log_error "Invalid environment: $ENVIRONMENT. Must be development, staging, or production."
    exit 1
fi

# Validate deployment type
if [[ ! "$DEPLOYMENT_TYPE" =~ ^(docker|kubernetes)$ ]]; then
    log_error "Invalid deployment type: $DEPLOYMENT_TYPE. Must be docker or kubernetes."
    exit 1
fi

log_info "Starting deployment of AI Consultant Live Chat System"
log_info "Environment: $ENVIRONMENT"
log_info "Deployment Type: $DEPLOYMENT_TYPE"
log_info "Project Root: $PROJECT_ROOT"

# Change to project root
cd "$PROJECT_ROOT"

# Load environment-specific configuration
ENV_FILE="config/environments/${ENVIRONMENT}.env"
if [[ -f "$ENV_FILE" ]]; then
    log_info "Loading environment configuration from $ENV_FILE"
    set -a
    source "$ENV_FILE"
    set +a
else
    log_warning "Environment file $ENV_FILE not found, using defaults"
fi

# Pre-deployment checks
log_info "Running pre-deployment checks..."

# Check if required files exist
required_files=(
    "docker-compose.chat.yml"
    "backend/Dockerfile"
    "frontend/Dockerfile"
    "redis.conf"
)

for file in "${required_files[@]}"; do
    if [[ ! -f "$file" ]]; then
        log_error "Required file not found: $file"
        exit 1
    fi
done

# Docker deployment
deploy_docker() {
    log_info "Deploying using Docker Compose..."
    
    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    
    # Check if docker-compose is available
    if ! command -v docker-compose >/dev/null 2>&1; then
        log_error "docker-compose is not installed. Please install docker-compose and try again."
        exit 1
    fi
    
    # Build and start services
    log_info "Building and starting services..."
    docker-compose -f docker-compose.chat.yml --env-file "$ENV_FILE" up -d --build
    
    # Wait for services to be healthy
    log_info "Waiting for services to be healthy..."
    sleep 30
    
    # Check service health
    services=("backend-chat" "frontend-chat" "db-chat" "redis-chat")
    for service in "${services[@]}"; do
        if docker-compose -f docker-compose.chat.yml ps "$service" | grep -q "Up (healthy)"; then
            log_success "$service is healthy"
        else
            log_warning "$service may not be healthy, checking logs..."
            docker-compose -f docker-compose.chat.yml logs --tail=20 "$service"
        fi
    done
    
    log_success "Docker deployment completed!"
    log_info "Services are available at:"
    log_info "  - Frontend: http://localhost:3006"
    log_info "  - Backend API: http://localhost:8061"
    log_info "  - WebSocket: ws://localhost:8061/api/v1/admin/chat/ws"
    log_info "  - Prometheus: http://localhost:9090"
    log_info "  - Grafana: http://localhost:3000"
}

# Kubernetes deployment
deploy_kubernetes() {
    log_info "Deploying using Kubernetes..."
    
    # Check if kubectl is available
    if ! command -v kubectl >/dev/null 2>&1; then
        log_error "kubectl is not installed. Please install kubectl and try again."
        exit 1
    fi
    
    # Check if helm is available
    if ! command -v helm >/dev/null 2>&1; then
        log_error "helm is not installed. Please install helm and try again."
        exit 1
    fi
    
    # Check cluster connectivity
    if ! kubectl cluster-info >/dev/null 2>&1; then
        log_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    # Create namespace if it doesn't exist
    if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
        log_info "Creating namespace: $NAMESPACE"
        kubectl create namespace "$NAMESPACE"
    fi
    
    # Add required Helm repositories
    log_info "Adding Helm repositories..."
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    
    # Deploy using Helm
    log_info "Deploying Helm chart..."
    helm upgrade --install ai-consultant-chat ./k8s/chat-system \
        --namespace "$NAMESPACE" \
        --values "./k8s/chat-system/values-${ENVIRONMENT}.yaml" \
        --wait \
        --timeout=10m
    
    # Wait for rollout to complete
    log_info "Waiting for deployment to complete..."
    kubectl rollout status deployment/ai-consultant-chat-backend -n "$NAMESPACE" --timeout=300s
    kubectl rollout status deployment/ai-consultant-chat-frontend -n "$NAMESPACE" --timeout=300s
    
    # Check pod status
    log_info "Checking pod status..."
    kubectl get pods -n "$NAMESPACE"
    
    # Get service information
    log_info "Getting service information..."
    kubectl get services -n "$NAMESPACE"
    
    # Get ingress information
    if kubectl get ingress -n "$NAMESPACE" >/dev/null 2>&1; then
        log_info "Ingress information:"
        kubectl get ingress -n "$NAMESPACE"
    fi
    
    log_success "Kubernetes deployment completed!"
    log_info "Use 'kubectl get all -n $NAMESPACE' to check all resources"
}

# Run database migrations
run_migrations() {
    log_info "Running database migrations..."
    
    if [[ "$DEPLOYMENT_TYPE" == "docker" ]]; then
        # Run migrations using Docker
        docker-compose -f docker-compose.chat.yml exec -T db-chat psql -U postgres -d consulting -f /docker-entrypoint-initdb.d/02-chat-migration.sql
        docker-compose -f docker-compose.chat.yml exec -T db-chat psql -U postgres -d consulting -f /docker-entrypoint-initdb.d/03-enhanced-chat-migration.sql
    else
        # Migrations are handled by the migration job in Kubernetes
        kubectl wait --for=condition=complete job/ai-consultant-chat-migration -n "$NAMESPACE" --timeout=300s
    fi
    
    log_success "Database migrations completed!"
}

# Health check
health_check() {
    log_info "Running health checks..."
    
    if [[ "$DEPLOYMENT_TYPE" == "docker" ]]; then
        # Check Docker services
        backend_health=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8061/health || echo "000")
        frontend_health=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3006/ || echo "000")
        
        if [[ "$backend_health" == "200" ]]; then
            log_success "Backend health check passed"
        else
            log_error "Backend health check failed (HTTP $backend_health)"
        fi
        
        if [[ "$frontend_health" == "200" ]]; then
            log_success "Frontend health check passed"
        else
            log_error "Frontend health check failed (HTTP $frontend_health)"
        fi
    else
        # Check Kubernetes services
        kubectl get pods -n "$NAMESPACE" -o wide
        
        # Port forward for health check
        kubectl port-forward -n "$NAMESPACE" service/ai-consultant-chat-backend 8061:8061 &
        PF_PID=$!
        sleep 5
        
        backend_health=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8061/health || echo "000")
        kill $PF_PID 2>/dev/null || true
        
        if [[ "$backend_health" == "200" ]]; then
            log_success "Backend health check passed"
        else
            log_error "Backend health check failed (HTTP $backend_health)"
        fi
    fi
}

# Main deployment logic
main() {
    case "$DEPLOYMENT_TYPE" in
        docker)
            deploy_docker
            ;;
        kubernetes)
            deploy_kubernetes
            ;;
    esac
    
    # Run migrations
    run_migrations
    
    # Health check
    health_check
    
    log_success "Deployment completed successfully!"
    
    # Show next steps
    log_info "Next steps:"
    log_info "1. Verify all services are running correctly"
    log_info "2. Check logs for any errors"
    log_info "3. Test the chat functionality"
    log_info "4. Monitor system metrics"
}

# Trap to cleanup on exit
cleanup() {
    log_info "Cleaning up..."
    # Kill any background processes
    jobs -p | xargs -r kill 2>/dev/null || true
}

trap cleanup EXIT

# Run main function
main