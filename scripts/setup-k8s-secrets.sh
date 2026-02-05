#!/bin/bash

# ===========================================
# Kubernetes Secrets Setup Script
# ===========================================

set -e

echo "🔐 Setting up Kubernetes secrets..."

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl first."
    exit 1
fi

# Check if namespace exists
if ! kubectl get namespace language-platform &> /dev/null; then
    echo "📦 Creating namespace..."
    kubectl create namespace language-platform
fi

# Load environment variables
if [ -f config/env/.env.production ]; then
  source config/env/.env.production
else
  echo "❌ config/env/.env.production not found"
    exit 1
fi

# Validate required variables
required_vars=(
    "DB_PASSWORD"
    "JWT_SECRET"
    "SENDGRID_API_KEY"
    "AWS_ACCESS_KEY_ID"
    "AWS_SECRET_ACCESS_KEY"
)

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "❌ Required variable $var is not set"
        exit 1
    fi
done

echo "✅ All required variables are set"

# Create database secret
echo "📝 Creating database secret..."
kubectl create secret generic db-secret \
  --from-literal=auth-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/auth_db?sslmode=require" \
  --from-literal=users-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/users_db?sslmode=require" \
  --from-literal=courses-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/courses_db?sslmode=require" \
  --from-literal=tasks-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/tasks_db?sslmode=require" \
  --from-literal=progress-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/progress_db?sslmode=require" \
  --from-literal=notifications-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/notifications_db?sslmode=require" \
  --from-literal=files-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/files_db?sslmode=require" \
  --from-literal=videos-url="postgres://langplatform:${DB_PASSWORD}@${DATABASE_HOST}:5432/videos_db?sslmode=require" \
  -n language-platform \
  --dry-run=client -o yaml | kubectl apply -f -

# Create JWT secret
echo "📝 Creating JWT secret..."
kubectl create secret generic jwt-secret \
  --from-literal=secret="${JWT_SECRET}" \
  -n language-platform \
  --dry-run=client -o yaml | kubectl apply -f -

# Create SendGrid secret
echo "📝 Creating SendGrid secret..."
kubectl create secret generic sendgrid-secret \
  --from-literal=api-key="${SENDGRID_API_KEY}" \
  -n language-platform \
  --dry-run=client -o yaml | kubectl apply -f -

# Create AWS secret
echo "📝 Creating AWS secret..."
kubectl create secret generic aws-secret \
  --from-literal=access-key-id="${AWS_ACCESS_KEY_ID}" \
  --from-literal=secret-access-key="${AWS_SECRET_ACCESS_KEY}" \
  -n language-platform \
  --dry-run=client -o yaml | kubectl apply -f -

# Create Zoom secret (if configured)
if [ ! -z "${ZOOM_API_KEY}" ]; then
    echo "📝 Creating Zoom secret..."
    kubectl create secret generic zoom-secret \
      --from-literal=api-key="${ZOOM_API_KEY}" \
      --from-literal=api-secret="${ZOOM_API_SECRET}" \
      -n language-platform \
      --dry-run=client -o yaml | kubectl apply -f -
fi

# Create Redis secret (if password is set)
if [ ! -z "${REDIS_PASSWORD}" ]; then
    echo "📝 Creating Redis secret..."
    kubectl create secret generic redis-secret \
      --from-literal=password="${REDIS_PASSWORD}" \
      -n language-platform \
      --dry-run=client -o yaml | kubectl apply -f -
fi

# Create Sentry secret (if configured)
if [ ! -z "${SENTRY_DSN}" ]; then
    echo "📝 Creating Sentry secret..."
    kubectl create secret generic sentry-secret \
      --from-literal=dsn="${SENTRY_DSN}" \
      -n language-platform \
      --dry-run=client -o yaml | kubectl apply -f -
fi

echo ""
echo "✅ All secrets created successfully!"
echo ""
echo "📋 View secrets:"
echo "   kubectl get secrets -n language-platform"
echo ""
echo "🔍 Describe a secret:"
echo "   kubectl describe secret db-secret -n language-platform"
echo ""
