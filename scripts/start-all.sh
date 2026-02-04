#!/bin/bash

echo "🚀 Starting Language Learning Platform..."

# Navigate to docker directory
cd "$(dirname "$0")/../infra/docker"

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cat > .env << EOF
DATABASE_PASSWORD=postgres
JWT_SECRET=your-secret-key-change-in-production-$(openssl rand -hex 32)
SENDGRID_API_KEY=your-sendgrid-api-key
EOF
fi

# Start services
echo "🐳 Starting Docker containers..."
docker-compose up -d

# Wait for services
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check service health
echo "🏥 Checking service health..."
services=(
    "http://localhost:5001/health"
    "http://localhost:8001/health"
    "http://localhost:8002/health"
    "http://localhost:8003/health"
    "http://localhost:8004/health"
    "http://localhost:8005/health"
    "http://localhost:8006/health"
    "http://localhost:8007/health"
)

for service in "${services[@]}"; do
    if curl -f -s "$service" > /dev/null; then
        echo "✅ $service is healthy"
    else
        echo "⚠️  $service is not responding"
    fi
done

echo ""
echo "✨ Platform is running!"
echo ""
echo "📊 Service URLs:"
echo "  Auth Service:         http://localhost:5001"
echo "  User Service:         http://localhost:8001"
echo "  Course Service:       http://localhost:8002"
echo "  Task Service:         http://localhost:8003"
echo "  Progress Service:     http://localhost:8004"
echo "  Notification Service: http://localhost:8005"
echo "  File Service:         http://localhost:8006"
echo "  Video Service:        http://localhost:8007"
echo ""
echo "🌐 Infrastructure:"
echo "  API Gateway:          http://localhost"
echo "  Prometheus:           http://localhost:9090"
echo "  Grafana:              http://localhost:3000 (admin/admin)"
echo "  MinIO Console:        http://localhost:9001 (minioadmin/minioadmin)"
echo ""
echo "📱 To start the frontend:"
echo "  cd frontend/web && npm install && npm run dev"
echo ""
echo "🛑 To stop all services:"
echo "  docker-compose down"
echo ""
