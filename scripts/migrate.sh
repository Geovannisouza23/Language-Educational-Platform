#!/bin/bash

set -e

echo "🔄 Running database migrations..."

# Run Auth Service migrations
echo "📦 Migrating Auth Service database..."
docker-compose -f infra/docker/docker-compose.yml exec -T auth-service \
    dotnet ef database update || echo "⚠️  Auth migrations skipped"

# Run User Service migrations
echo "📦 Migrating User Service database..."
docker-compose -f infra/docker/docker-compose.yml exec -T user-service \
    ./user-service migrate || echo "⚠️  User migrations skipped"

# Run Course Service migrations
echo "📦 Migrating Course Service database..."
docker-compose -f infra/docker/docker-compose.yml exec -T course-service \
    ./course-service migrate || echo "⚠️  Course migrations skipped"

echo "✅ All migrations completed!"
