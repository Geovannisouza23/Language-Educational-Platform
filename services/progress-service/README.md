# Language Learning Platform - Progress Service

Microservice for tracking student learning progress and achievements.

## Features

- 📈 Track course completion percentage
- 🎯 Lesson and task completion tracking
- 📊 Performance analytics
- 🏆 Achievement and badge system
- 🔥 Learning streak tracking
- ⏱️ Study time monitoring
- 📉 Average score calculation

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis

## Data Models

### Progress
```go
- ID (UUID)
- StudentID (UUID)
- CourseID (UUID)
- CompletedLessons (int)
- TotalLessons (int)
- CompletedTasks (int)
- TotalTasks (int)
- AverageScore (float)
- StudyTimeMinutes (int)
- Streak (int)
- LastStudyDate
```

### LessonProgress
```go
- ID (UUID)
- StudentID (UUID)
- LessonID (UUID)
- Completed (bool)
- TimeSpent (minutes)
- CompletedAt
```

### Achievement
```go
- ID (UUID)
- StudentID (UUID)
- Title
- Description
- IconURL
- UnlockedAt
```

## API Endpoints

### Progress
- `GET /api/v1/progress/:courseId` - Get course progress
- `GET /api/v1/progress` - Get all student progress
- `POST /api/v1/progress/lesson/:lessonId/complete` - Mark lesson complete
- `GET /api/v1/achievements` - Get student achievements
- `GET /api/v1/stats` - Get overall statistics

## Environment Variables

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/progress_db
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key
```

## Achievement Triggers

- **First Steps**: Complete first lesson
- **Committed**: 7-day streak
- **Dedicated**: 30-day streak
- **Perfectionist**: 100% on 5 tasks
- **Course Master**: Complete entire course
- **Fast Learner**: Complete course in <30 days

## Development

```bash
go mod download
go run cmd/main.go
```
