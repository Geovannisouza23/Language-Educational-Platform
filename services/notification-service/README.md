# Language Learning Platform - Notification Service

Microservice for managing notifications (email, push, in-app).

## Features

- 📧 Email notifications via SendGrid
- 📲 Push notifications (FCM ready)
- 💬 In-app notifications
- 📋 Notification templates
- 🔔 User notification preferences
- 📊 Delivery tracking
- 🌐 Multi-language support

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Email**: SendGrid API
- **Push**: Firebase Cloud Messaging (FCM)

## Data Models

### Notification
```go
- ID (UUID)
- UserID (UUID)
- Type (email|push|in_app)
- Status (pending|sent|failed|read)
- Title
- Message
- Data (JSON metadata)
- SentAt
- ReadAt
```

### EmailTemplate
```go
- ID (UUID)
- Name (unique)
- Subject
- Body (HTML template)
```

## API Endpoints

### Notifications
- `POST /api/v1/notifications/send` - Send notification
- `GET /api/v1/notifications` - List user notifications
- `POST /api/v1/notifications/:id/read` - Mark as read
- `DELETE /api/v1/notifications/:id` - Delete notification

### Templates
- `POST /api/v1/templates` - Create email template
- `GET /api/v1/templates` - List templates
- `GET /api/v1/templates/:name` - Get template

## Notification Types

### Email Notifications
- Welcome email
- Email verification
- Password reset
- Course enrollment confirmation
- Assignment due reminder
- Grade notification
- New course announcement

### Push Notifications
- New message
- Live class starting
- Assignment graded
- Achievement unlocked

### In-App Notifications
- Teacher replied to question
- New lesson available
- Course update

## Environment Variables

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/notifications_db
SENDGRID_API_KEY=your-api-key
EMAIL_FROM=noreply@platform.com
FCM_SERVER_KEY=your-fcm-key
```

## Development

```bash
go mod download
go run cmd/main.go
```

## Email Templates

Templates use Go template syntax:

```html
<!DOCTYPE html>
<html>
<body>
  <h1>Welcome {{.Name}}!</h1>
  <p>Thanks for joining our platform.</p>
</body>
</html>
```
