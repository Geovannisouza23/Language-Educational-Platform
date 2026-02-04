# API Contracts

## Authentication Service

### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123",
  "confirmPassword": "securePassword123",
  "role": "Student"
}

Response 200:
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "base64-encoded-token",
  "expiresAt": "2024-02-04T12:00:00Z",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "Student",
    "emailVerified": false
  }
}
```

### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}

Response 200:
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "base64-encoded-token",
  "expiresAt": "2024-02-04T12:00:00Z",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "Student",
    "emailVerified": true
  }
}
```

### Refresh Token
```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refreshToken": "base64-encoded-token"
}

Response 200:
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "new-base64-encoded-token",
  "expiresAt": "2024-02-04T12:00:00Z",
  "user": { ... }
}
```

## User Service

### Get User Profile
```http
GET /api/v1/users/{id}/profile
Authorization: Bearer {accessToken}

Response 200:
{
  "id": "uuid",
  "userId": "uuid",
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "+1234567890",
  "dateOfBirth": "1990-01-01",
  "country": "USA",
  "city": "New York",
  "bio": "Language enthusiast",
  "avatarUrl": "https://...",
  "language": "en",
  "timezone": "America/New_York",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-02-01T00:00:00Z"
}
```

### Update User Profile
```http
PUT /api/v1/users/{id}/profile
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "+1234567890",
  "country": "USA",
  "city": "New York",
  "bio": "Updated bio",
  "language": "en"
}

Response 200:
{
  "id": "uuid",
  "userId": "uuid",
  ...
}
```

### List Users
```http
GET /api/v1/users?limit=10&offset=0
Authorization: Bearer {accessToken}

Response 200:
{
  "users": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "role": "Student",
      "isActive": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "profile": { ... }
    }
  ],
  "limit": 10,
  "offset": 0
}
```

## Course Service

### Create Course
```http
POST /api/v1/courses
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "title": "Spanish for Beginners",
  "description": "Learn Spanish from scratch",
  "language": "Spanish",
  "level": "Beginner",
  "price": 49.99,
  "duration": 30,
  "maxStudents": 20,
  "thumbnailUrl": "https://..."
}

Response 201:
{
  "id": "uuid",
  "title": "Spanish for Beginners",
  "teacherId": "uuid",
  "status": "draft",
  ...
}
```

### List Courses
```http
GET /api/v1/courses?language=Spanish&level=Beginner&limit=10
Authorization: Bearer {accessToken}

Response 200:
{
  "courses": [
    {
      "id": "uuid",
      "title": "Spanish for Beginners",
      "teacher": {
        "id": "uuid",
        "name": "John Doe"
      },
      "enrolledCount": 15,
      "rating": 4.5,
      ...
    }
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

### Enroll in Course
```http
POST /api/v1/courses/{id}/enroll
Authorization: Bearer {accessToken}

Response 200:
{
  "message": "Successfully enrolled",
  "enrollment": {
    "id": "uuid",
    "courseId": "uuid",
    "studentId": "uuid",
    "enrolledAt": "2024-02-04T00:00:00Z",
    "status": "active"
  }
}
```

## Task Service

### Create Task
```http
POST /api/v1/tasks
Authorization: Bearer {accessToken}
Content-Type: application/json

{
  "courseId": "uuid",
  "title": "Homework 1",
  "description": "Complete exercises 1-10",
  "type": "assignment",
  "dueDate": "2024-02-10T23:59:59Z",
  "maxScore": 100
}

Response 201:
{
  "id": "uuid",
  "courseId": "uuid",
  "title": "Homework 1",
  "status": "published",
  ...
}
```

### Submit Task
```http
POST /api/v1/tasks/{id}/submit
Authorization: Bearer {accessToken}
Content-Type: multipart/form-data

{
  "file": [binary],
  "notes": "My submission notes"
}

Response 200:
{
  "submissionId": "uuid",
  "status": "submitted",
  "submittedAt": "2024-02-05T10:00:00Z"
}
```

## Progress Service

### Get Student Progress
```http
GET /api/v1/progress/student/{id}/course/{courseId}
Authorization: Bearer {accessToken}

Response 200:
{
  "studentId": "uuid",
  "courseId": "uuid",
  "completedLessons": 5,
  "totalLessons": 10,
  "completionPercentage": 50,
  "averageScore": 85.5,
  "streak": 7,
  "lastActivityAt": "2024-02-04T10:00:00Z"
}
```

## Common Response Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no body
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service down

## Error Response Format

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "Specific error details"
  }
}
```

## Authentication

All endpoints (except auth endpoints) require JWT token:

```
Authorization: Bearer {accessToken}
```

Token should be included in the `Authorization` header.
