## Quick Start Guide

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- .NET 8.0 SDK (for local development)
- Node.js 20+ (for frontend)
- kubectl (for Kubernetes deployment)

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/your-org/language-platform.git
cd language-platform
```

2. **Start all services**
```bash
chmod +x scripts/local-dev.sh
./scripts/local-dev.sh
```

3. **Run migrations**
```bash
chmod +x scripts/migrate.sh
./scripts/migrate.sh
```

4. **Seed database** (optional)
```bash
chmod +x scripts/seed-db.sh
./scripts/seed-db.sh
```

5. **Access the services**
- API Gateway: http://localhost
- Auth Service: http://localhost:5001
- User Service: http://localhost:8001
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090

### Testing the API

#### Register a new user
```bash
curl -X POST http://localhost/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "Test123!",
    "confirmPassword": "Test123!",
    "role": "Student"
  }'
```

#### Login
```bash
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "Test123!"
  }'
```

Save the `accessToken` from the response.

#### Get user profile
```bash
curl -X GET http://localhost/api/v1/users/{USER_ID}/profile \
  -H "Authorization: Bearer {ACCESS_TOKEN}"
```

### Development Workflow

#### Working on Auth Service (C#)
```bash
cd services/auth-service
dotnet restore
dotnet run
```

#### Working on User Service (Go)
```bash
cd services/user-service
go mod download
go run cmd/main.go
```

#### Running Tests

**Go services:**
```bash
cd services/user-service
go test ./...
```

**C# services:**
```bash
cd services/auth-service
dotnet test
```

### Troubleshooting

#### Services not starting
```bash
# Check Docker logs
docker-compose -f infra/docker/docker-compose.yml logs -f

# Restart a specific service
docker-compose -f infra/docker/docker-compose.yml restart auth-service
```

#### Database connection issues
```bash
# Check PostgreSQL is running
docker-compose -f infra/docker/docker-compose.yml ps postgres

# Access PostgreSQL
docker-compose -f infra/docker/docker-compose.yml exec postgres psql -U postgres
```

#### Clear all data and restart
```bash
docker-compose -f infra/docker/docker-compose.yml down -v
./scripts/local-dev.sh
```

### Production Deployment

#### Build and push Docker images
```bash
# Auth Service
cd services/auth-service
docker build -t ghcr.io/your-org/auth-service:v1.0.0 .
docker push ghcr.io/your-org/auth-service:v1.0.0

# User Service
cd services/user-service
docker build -t ghcr.io/your-org/user-service:v1.0.0 .
docker push ghcr.io/your-org/user-service:v1.0.0
```

#### Deploy to Kubernetes
```bash
# Create secrets file
cp infra/kubernetes/overlays/production/secrets.env.example \
   infra/kubernetes/overlays/production/secrets.env

# Edit secrets.env with your production values
nano infra/kubernetes/overlays/production/secrets.env

# Deploy
kubectl apply -k infra/kubernetes/overlays/production/

# Check status
kubectl get pods -n language-platform
kubectl get services -n language-platform
```

### Monitoring

#### View Prometheus metrics
1. Open http://localhost:9090
2. Try queries like:
   - `up` - Check service health
   - `http_requests_total` - Total HTTP requests
   - `http_request_duration_seconds` - Request duration

#### View Grafana dashboards
1. Open http://localhost:3000
2. Login with admin/admin
3. Import dashboards from `infra/grafana/dashboards/`

### Next Steps

- [ ] Configure email service for notifications
- [ ] Set up video conferencing integration
- [ ] Implement course-service
- [ ] Implement task-service
- [ ] Implement progress-service
- [ ] Build frontend with Next.js
- [ ] Build mobile app
- [ ] Set up production monitoring
- [ ] Configure CI/CD pipeline
