#!/bin/bash

# ===========================================
# Production Deployment Script
# ===========================================

set -e

echo "🚀 Starting production deployment..."

# Configuration
ENVIRONMENT=${1:-production}
NAMESPACE="language-platform"
DOCKER_REGISTRY=${DOCKER_REGISTRY:-"your-docker-registry"}
VERSION=${VERSION:-$(git rev-parse --short HEAD)}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    commands=("docker" "kubectl" "git")
    for cmd in "${commands[@]}"; do
        if ! command -v $cmd &> /dev/null; then
            log_error "$cmd is not installed"
            exit 1
        fi
    done
    
    log_info "✅ All prerequisites met"
}

load_environment() {
    log_info "Loading environment variables..."
    
    if [ -f ".env.$ENVIRONMENT" ]; then
        source ".env.$ENVIRONMENT"
        log_info "✅ Environment loaded"
    else
        log_error ".env.$ENVIRONMENT file not found"
        exit 1
    fi
}

build_services() {
    log_info "Building Docker images..."
    
    services=(
        "auth-service"
        "user-service"
        "course-service"
        "task-service"
        "progress-service"
        "notification-service"
        "file-service"
        "video-service"
    )
    
    for service in "${services[@]}"; do
        log_info "Building $service..."
        docker build -t ${DOCKER_REGISTRY}/${service}:${VERSION} \
                     -t ${DOCKER_REGISTRY}/${service}:latest \
                     ./services/${service}
        
        if [ $? -eq 0 ]; then
            log_info "✅ $service built successfully"
        else
            log_error "Failed to build $service"
            exit 1
        fi
    done
}

push_images() {
    log_info "Pushing images to registry..."
    
    docker login ${DOCKER_REGISTRY}
    
    services=(
        "auth-service"
        "user-service"
        "course-service"
        "task-service"
        "progress-service"
        "notification-service"
        "file-service"
        "video-service"
    )
    
    for service in "${services[@]}"; do
        log_info "Pushing $service..."
        docker push ${DOCKER_REGISTRY}/${service}:${VERSION}
        docker push ${DOCKER_REGISTRY}/${service}:latest
    done
    
    log_info "✅ All images pushed"
}

setup_kubernetes() {
    log_info "Setting up Kubernetes resources..."
    
    # Create namespace if it doesn't exist
    kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
    
    # Setup secrets
    log_info "Setting up secrets..."
    ./scripts/setup-k8s-secrets.sh
    
    log_info "✅ Kubernetes setup complete"
}

deploy_to_kubernetes() {
    log_info "Deploying to Kubernetes..."
    
    # Apply base configurations
    kubectl apply -k infra/kubernetes/base
    
    # Apply environment-specific overlays
    kubectl apply -k infra/kubernetes/overlays/${ENVIRONMENT}
    
    # Update image tags
    services=(
        "auth-service"
        "user-service"
        "course-service"
        "task-service"
        "progress-service"
        "notification-service"
        "file-service"
        "video-service"
    )
    
    for service in "${services[@]}"; do
        kubectl set image deployment/${service} \
            ${service}=${DOCKER_REGISTRY}/${service}:${VERSION} \
            -n ${NAMESPACE}
    done
    
    log_info "✅ Deployment initiated"
}

wait_for_rollout() {
    log_info "Waiting for rollout to complete..."
    
    services=(
        "auth-service"
        "user-service"
        "course-service"
        "task-service"
        "progress-service"
        "notification-service"
        "file-service"
        "video-service"
    )
    
    for service in "${services[@]}"; do
        log_info "Checking $service..."
        kubectl rollout status deployment/${service} -n ${NAMESPACE} --timeout=5m
    done
    
    log_info "✅ All services rolled out successfully"
}

run_migrations() {
    log_info "Running database migrations..."
    
    # Auth Service migrations
    kubectl run migration-auth --rm -i --restart=Never \
        --image=${DOCKER_REGISTRY}/auth-service:${VERSION} \
        --command -- dotnet ef database update \
        -n ${NAMESPACE}
    
    log_info "✅ Migrations complete"
}

smoke_test() {
    log_info "Running smoke tests..."
    
    # Get service endpoints
    API_URL=$(kubectl get ingress -n ${NAMESPACE} -o jsonpath='{.items[0].spec.rules[0].host}')
    
    # Test health endpoints
    services=("auth" "users" "courses" "tasks" "progress" "notifications" "files" "sessions")
    
    for service in "${services[@]}"; do
        log_info "Testing $service..."
        response=$(curl -s -o /dev/null -w "%{http_code}" https://${API_URL}/api/health/${service})
        
        if [ "$response" == "200" ]; then
            log_info "✅ $service is healthy"
        else
            log_warn "⚠️  $service returned $response"
        fi
    done
}

rollback() {
    log_error "Deployment failed. Rolling back..."
    
    services=(
        "auth-service"
        "user-service"
        "course-service"
        "task-service"
        "progress-service"
        "notification-service"
        "file-service"
        "video-service"
    )
    
    for service in "${services[@]}"; do
        kubectl rollout undo deployment/${service} -n ${NAMESPACE}
    done
    
    log_info "✅ Rollback complete"
}

# Main execution
main() {
    log_info "==================================="
    log_info "Production Deployment Script"
    log_info "Environment: $ENVIRONMENT"
    log_info "Version: $VERSION"
    log_info "==================================="
    echo ""
    
    check_prerequisites
    load_environment
    
    # Build and push
    build_services
    push_images
    
    # Deploy
    setup_kubernetes
    deploy_to_kubernetes
    
    # Wait and verify
    wait_for_rollout
    
    # Post-deployment
    run_migrations
    smoke_test
    
    echo ""
    log_info "==================================="
    log_info "✅ Deployment completed successfully!"
    log_info "==================================="
    echo ""
    log_info "📊 View status:"
    log_info "   kubectl get pods -n ${NAMESPACE}"
    echo ""
    log_info "📋 View logs:"
    log_info "   kubectl logs -f deployment/auth-service -n ${NAMESPACE}"
    echo ""
    log_info "🔄 Rollback if needed:"
    log_info "   ./scripts/rollback.sh"
    echo ""
}

# Trap errors and rollback
trap 'rollback' ERR

# Run main
main
