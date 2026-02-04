# Language Learning Platform - Video Service

Microservice for video conferencing and live class management.

## Features

- 🎥 Schedule live video sessions
- 👥 Participant management
- 📹 Recording support
- 📊 Session analytics
- ⏱️ Duration tracking
- 🔗 Meeting URL generation
- 📅 Calendar integration ready

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Video**: Zoom/Agora/Jitsi SDK integration ready

## Data Models

### VideoSession
```go
- ID (UUID)
- CourseID (UUID)
- TeacherID (UUID)
- Title
- Description
- ScheduledAt
- Duration (minutes)
- Status (scheduled|live|ended|cancelled)
- MeetingURL
- RecordingURL
- StartedAt
- EndedAt
```

### Participant
```go
- ID (UUID)
- SessionID (UUID)
- UserID (UUID)
- JoinedAt
- LeftAt
- Duration (minutes)
```

## API Endpoints

### Sessions
- `POST /api/v1/sessions` - Create session (Teacher)
- `GET /api/v1/sessions/:id` - Get session details
- `GET /api/v1/sessions` - List sessions (filter by course)
- `PUT /api/v1/sessions/:id` - Update session
- `DELETE /api/v1/sessions/:id` - Cancel session
- `POST /api/v1/sessions/:id/start` - Start session
- `POST /api/v1/sessions/:id/end` - End session

### Participants
- `POST /api/v1/sessions/:id/join` - Join session
- `POST /api/v1/sessions/:id/leave` - Leave session
- `GET /api/v1/sessions/:id/participants` - List participants

## Environment Variables

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/videos_db
JWT_SECRET=your-secret-key
ZOOM_API_KEY=your-zoom-key
ZOOM_API_SECRET=your-zoom-secret
AGORA_APP_ID=your-agora-app-id
JITSI_DOMAIN=meet.jit.si
```

## Video Provider Integration

### Zoom
- OAuth authentication
- Meeting creation via REST API
- Webhook for events

### Agora
- RTC token generation
- Channel management

### Jitsi
- Self-hosted or meet.jit.si
- JWT authentication

## Session Lifecycle

1. **Scheduled**: Session created, waiting for start time
2. **Live**: Session started, participants can join
3. **Ended**: Session finished, recording available
4. **Cancelled**: Session cancelled by teacher

## Development

```bash
go mod download
go run cmd/main.go
```

## Create Session Example

```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": "uuid",
    "title": "Spanish Conversation Practice",
    "description": "Live conversation class",
    "scheduled_at": "2024-12-20T15:00:00Z",
    "duration": 60
  }'
```

## Features Roadmap

- [ ] Breakout rooms
- [ ] Screen sharing tracking
- [ ] Chat integration
- [ ] Polls and quizzes during session
- [ ] Whiteboard integration
- [ ] Recording auto-upload
