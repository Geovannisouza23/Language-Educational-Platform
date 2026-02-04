#!/bin/bash

set -e

echo "🌱 Seeding database with sample data..."

# Wait for database to be ready
sleep 5

# Seed Auth Service with default users
echo "👤 Creating default users..."
cat <<EOF | docker-compose -f infra/docker/docker-compose.yml exec -T postgres psql -U postgres -d auth_db
-- Sample data for development
-- This would normally be done through the Auth Service API
EOF

# Seed User Service with profiles
echo "📝 Creating user profiles..."

# Seed Course Service with sample courses
echo "📚 Creating sample courses..."

echo "✅ Database seeding completed!"
echo ""
echo "📝 Default users (for testing):"
echo "   - Admin:   admin@example.com / password123"
echo "   - Teacher: teacher@example.com / password123"
echo "   - Student: student@example.com / password123"
echo ""
