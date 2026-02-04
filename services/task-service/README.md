# Language Learning Platform - Task Service

Microservice for managing assignments, quizzes, and student submissions.

## Features

- ✍️ Create assignments, quizzes, essays, speaking tasks
- 📤 Student submission management
- 📊 Automatic and manual grading
- 📅 Due date tracking
- 📎 File attachments support
- 🔒 Teacher-only task creation
- 🎯 Task type categorization

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **ORM**: GORM

## Data Models

### Task
```go
- ID (UUID)
- CourseID (UUID)
- TeacherID (UUID)
- Title
- Description
- Type (assignment|quiz|essay|speaking)
- Status (draft|published)
- DueDate
- MaxScore
- Content (JSON for quiz questions, etc.)
```

### Submission
```go
- ID (UUID)
- TaskID (UUID)
- StudentID (UUID)
- Content (submission text)
- FileURL (optional attachment)
- Status (pending|submitted|graded|revision)
- Score
- Feedback
- SubmittedAt
- GradedAt
```

## API Endpoints

### Tasks
- `POST /api/v1/tasks` - Create task (Teacher only)
- `GET /api/v1/tasks/:id` - Get task details
- `GET /api/v1/tasks` - List tasks (filter by course)
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task
- `POST /api/v1/tasks/:id/publish` - Publish task

### Submissions
- `POST /api/v1/tasks/:id/submit` - Submit assignment (Student)
- `GET /api/v1/submissions/:id` - Get submission
- `GET /api/v1/tasks/:id/submissions` - List all submissions for task
- `GET /api/v1/my-submissions` - Get student's submissions
- `POST /api/v1/submissions/:id/grade` - Grade submission (Teacher)

## Environment Variables

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/tasks_db
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key
```

## Development

```bash
# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go

# Start service
go run cmd/main.go

# Run tests
go test ./...
```

## Docker

```bash
docker build -t task-service .
docker run -p 8080:8080 task-service
```
