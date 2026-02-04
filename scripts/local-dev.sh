#!/bin/bash

set -e

echo "🚀 Starting Language Platform in Development Mode..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Navigate to docker directory
cd "$(dirname "$0")/../infra/docker"

# Stop any running containers
echo "🛑 Stopping existing containers..."
docker-compose down

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up -d --build

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check service health
echo "🏥 Checking service health..."
services=("postgres:5432" "redis:6379" "auth-service:80" "user-service:8080")

for service in "${services[@]}"; do
    IFS=':' read -r name port <<< "$service"
    if docker-compose ps | grep -q "$name.*Up"; then
        echo "✅ $name is running"
    else
        echo "❌ $name is not running"
    fi
done

echo ""
echo "🎉 Language Platform is running!"
echo ""
echo "📍 Services:"
echo "   - API Gateway:    http://localhost"
echo "   - Auth Service:   http://localhost:5001"
echo "   - User Service:   http://localhost:8001"
echo "   - Course Service: http://localhost:8002"
echo "   - PostgreSQL:     localhost:5432"
echo "   - Redis:          localhost:6379"
echo "   - MinIO:          http://localhost:9001"
echo "   - Prometheus:     http://localhost:9090"
echo "   - Grafana:        http://localhost:3000 (admin/admin)"
echo ""
echo "📝 Logs: docker-compose -f infra/docker/docker-compose.yml logs -f [service-name]"
echo "🛑 Stop: docker-compose -f infra/docker/docker-compose.yml down"
echo ""
