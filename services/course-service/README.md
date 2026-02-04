# Course Service

Microservice for managing courses, lessons, and enrollments in the Language Learning Platform.

## Features

- **Course Management**: Create, read, update, delete courses
- **Lesson Management**: Add lessons to courses with content and videos
- **Enrollment System**: Students can enroll/unenroll from courses
- **Access Control**: Role-based permissions (Teachers, Students, Admins)
- **Course Publishing**: Draft and publish workflow
- **Caching**: Redis for performance optimization
- **Multi-language Support**: Courses in different languages
- **Level System**: Beginner, Intermediate, Advanced

## API Endpoints

### Courses

#### Public Endpoints
- `GET /api/v1/courses` - List all published courses
- `GET /api/v1/courses/:id` - Get course details
- `GET /api/v1/courses/:id/lessons` - Get course lessons

#### Teacher Endpoints (Authentication Required)
- `POST /api/v1/courses` - Create new course
- `PUT /api/v1/courses/:id` - Update course
- `DELETE /api/v1/courses/:id` - Delete course
- `PUT /api/v1/courses/:id/publish` - Publish course
- `POST /api/v1/courses/:id/lessons` - Add lesson to course
- `PUT /api/v1/courses/:id/lessons/:lessonId` - Update lesson
- `DELETE /api/v1/courses/:id/lessons/:lessonId` - Delete lesson
- `GET /api/v1/courses/:id/enrollments` - View course enrollments

#### Student Endpoints (Authentication Required)
- `POST /api/v1/courses/:id/enroll` - Enroll in course
- `DELETE /api/v1/courses/:id/unenroll` - Unenroll from course
- `GET /api/v1/courses/my-courses` - Get my enrolled courses

## Data Models

### Course
```go
{
  "id": "uuid",
  "teacher_id": "uuid",
  "title": "Spanish for Beginners",
  "description": "Learn Spanish from scratch",
  "language": "Spanish",
  "level": "Beginner",
  "status": "published",
  "price": 49.99,
  "duration": 30,
  "max_students": 20,
  "thumbnail_url": "https://...",
  "enrolled_count": 15,
  "rating": 4.5,
  "rating_count": 10
}
```

### Lesson
```go
{
  "id": "uuid",
  "course_id": "uuid",
  "title": "Lesson 1: Introduction",
  "description": "Basic introduction",
  "content": "Lesson content here",
  "video_url": "https://...",
  "duration": 15,
  "order": 1,
  "is_published": true
}
```

### Enrollment
```go
{
  "id": "uuid",
  "course_id": "uuid",
  "student_id": "uuid",
  "status": "active",
  "progress": 45.5,
  "enrolled_at": "2024-02-04T10:00:00Z",
  "last_access_at": "2024-02-04T12:00:00Z"
}
```

## Environment Variables

- `COURSE_SERVICE_ENVIRONMENT` - Environment (development/production)
- `COURSE_SERVICE_PORT` - Server port (default: 8080)
- `COURSE_SERVICE_DATABASE_URL` - PostgreSQL connection string
- `COURSE_SERVICE_REDIS_URL` - Redis connection string
- `COURSE_SERVICE_JWT_SECRET` - JWT secret for token validation
- `COURSE_SERVICE_LOG_LEVEL` - Log level (debug/info/warn/error)

## Development

```bash
# Run locally
go run cmd/main.go

# Run tests
go test ./...

# Build
go build -o course-service cmd/main.go

# Run with Docker
docker build -t course-service .
docker run -p 8080:8080 course-service
```

## Database Migrations

Migrations run automatically on service startup. Tables created:
- `courses` - Course information
- `lessons` - Course lessons
- `enrollments` - Student enrollments

## Caching Strategy

- Course details cached for 5 minutes
- Cache invalidated on course updates
- Cache key format: `course:{id}`

## Dependencies

- Gin Web Framework
- GORM ORM
- PostgreSQL
- Redis
- JWT for authentication
- Logrus for logging
