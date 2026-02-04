# Architecture Documentation

## Overview

Language Learning Platform is a microservices-based SaaS platform for language education with support for multiple teachers, students, and languages.

## System Architecture

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │  Web (Next)  │  │    Mobile    │  │   Admin Panel   │   │
│  └──────────────┘  └──────────────┘  └─────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS
┌────────────────────────▼────────────────────────────────────┐
│                    API Gateway (NGINX)                       │
│  - Authentication / Authorization                            │
│  - Rate Limiting                                             │
│  - SSL Termination                                           │
│  - Load Balancing                                            │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
┌───────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
│Auth Service  │ │User Service │ │Course Svc   │
│   (C#)       │ │    (Go)     │ │    (Go)     │
└───────┬──────┘ └──────┬──────┘ └──────┬──────┘
        │                │                │
        └────────────────┼────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
┌───────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
│  PostgreSQL  │ │    Redis    │ │   MinIO/S3  │
└──────────────┘ └─────────────┘ └─────────────┘
```

## Microservices

### 1. Auth Service (C# / ASP.NET Core)

**Responsibilities:**
- User authentication and authorization
- JWT token generation and validation
- Refresh token management
- Role-based access control (RBAC)
- Email verification

**Technology Stack:**
- ASP.NET Core 8.0
- Entity Framework Core
- PostgreSQL
- Redis (session storage)
- BCrypt (password hashing)

**Endpoints:**
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - User logout
- `POST /api/auth/verify-email` - Email verification
- `GET /api/auth/me` - Get current user info

### 2. User Service (Go)

**Responsibilities:**
- User profile management
- User CRUD operations
- Profile information (name, avatar, preferences)
- User search and filtering

**Technology Stack:**
- Go 1.21+
- Gin Web Framework
- GORM
- PostgreSQL
- Redis (caching)

**Endpoints:**
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user details
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users/:id/profile` - Get user profile
- `PUT /api/v1/users/:id/profile` - Update user profile

### 3. Course Service (Go)

**Responsibilities:**
- Course management (CRUD)
- Lesson organization
- Course enrollment
- Teacher-course association
- Multi-language support

**Technology Stack:**
- Go 1.21+
- Gin Web Framework
- GORM
- PostgreSQL

**Endpoints:**
- `POST /api/v1/courses` - Create course
- `GET /api/v1/courses` - List courses
- `GET /api/v1/courses/:id` - Get course details
- `PUT /api/v1/courses/:id` - Update course
- `DELETE /api/v1/courses/:id` - Delete course
- `POST /api/v1/courses/:id/enroll` - Enroll student
- `GET /api/v1/courses/:id/students` - List enrolled students

### 4. Task Service (Go)

**Responsibilities:**
- Assignment creation and management
- Task submission
- Deadline management
- Task types (quiz, essay, speaking, etc.)

### 5. Progress Service (Go)

**Responsibilities:**
- Student progress tracking
- Completion status
- Performance analytics
- Learning streak
- Achievement system

### 6. Notification Service (Go)

**Responsibilities:**
- Email notifications
- Push notifications
- In-app notifications
- Notification preferences

### 7. File Service (Go)

**Responsibilities:**
- File upload/download
- Image processing
- Document storage
- Video storage

### 8. Video Service (Go)

**Responsibilities:**
- Video conferencing integration
- Live class scheduling
- Recording management
- Chat during sessions

## Data Models

### Auth Service

```sql
Users
- id (UUID, PK)
- email (VARCHAR, UNIQUE)
- password_hash (VARCHAR)
- role_id (INT, FK)
- is_active (BOOLEAN)
- email_verified (BOOLEAN)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

Roles
- id (INT, PK)
- name (VARCHAR)
- description (TEXT)

RefreshTokens
- id (UUID, PK)
- user_id (UUID, FK)
- token (VARCHAR, UNIQUE)
- expires_at (TIMESTAMP)
- is_revoked (BOOLEAN)
```

### User Service

```sql
UserProfiles
- id (UUID, PK)
- user_id (UUID, UNIQUE, FK)
- first_name (VARCHAR)
- last_name (VARCHAR)
- phone_number (VARCHAR)
- date_of_birth (DATE)
- country (VARCHAR)
- city (VARCHAR)
- bio (TEXT)
- avatar_url (VARCHAR)
- language (VARCHAR)
- timezone (VARCHAR)
```

## Security

### Authentication Flow

1. User sends credentials to Auth Service
2. Auth Service validates and generates JWT + Refresh Token
3. Client stores tokens securely
4. Client includes JWT in Authorization header for API calls
5. Gateway/Services validate JWT
6. Refresh token used to get new JWT when expired

### Authorization

- Role-Based Access Control (RBAC)
- Roles: Admin, Teacher, Student
- JWT contains user ID, email, role
- Each service validates permissions based on role

## Scalability

### Horizontal Scaling
- All services are stateless
- Can be scaled independently
- Load balancing via NGINX/Kubernetes

### Caching Strategy
- Redis for session data
- User profile caching (5 min TTL)
- Course data caching
- Invalidation on updates

### Database Strategy
- Separate database per service
- Connection pooling
- Read replicas for heavy read operations

## Observability

### Logging
- Structured logging (JSON format)
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Centralized log aggregation

### Monitoring
- Prometheus for metrics collection
- Grafana for visualization
- Custom dashboards per service

### Tracing
- Jaeger for distributed tracing
- Request ID propagation
- Performance bottleneck identification

## Deployment

### Container Orchestration
- Kubernetes for production
- Docker Compose for local development

### CI/CD Pipeline
- GitHub Actions
- Automated testing
- Docker image building
- Automated deployment to staging/production

### Environments
- Development (local Docker Compose)
- Staging (Kubernetes cluster)
- Production (Kubernetes cluster with HA)

## Communication Patterns

### Synchronous
- REST APIs for client-service communication
- Service-to-service HTTP calls (with circuit breakers)

### Asynchronous (Future)
- Message queue (RabbitMQ/Kafka) for events
- Event-driven architecture for notifications
- Background job processing

## Disaster Recovery

### Backup Strategy
- Daily database backups
- Point-in-time recovery
- S3/MinIO for file backups

### High Availability
- Multi-replica deployments
- Database replication
- Redis Sentinel for cache HA
- Load balancing across instances
