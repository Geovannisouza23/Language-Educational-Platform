# Language Learning Platform - File Service

Microservice for file upload, storage, and retrieval.

## Features

- 📁 File upload (images, videos, documents, audio)
- 🖼️ Image processing and resizing
- 🎥 Video transcoding support
- 📄 Document preview generation
- 🔒 Secure file access
- 💾 S3/MinIO storage
- 🗑️ File deletion and lifecycle
- 📊 Storage quota management

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Storage**: MinIO / AWS S3
- **Image Processing**: Go imaging libraries

## Data Models

### File
```go
- ID (UUID)
- UserID (UUID)
- Filename (generated)
- OriginalName
- FileType (image|video|document|audio)
- MimeType
- Size (bytes)
- URL (CDN/direct URL)
- BucketName
- ObjectKey
```

## API Endpoints

### Files
- `POST /api/v1/upload` - Upload file
- `GET /api/v1/files/:id` - Get file metadata
- `GET /api/v1/files` - List user files
- `DELETE /api/v1/files/:id` - Delete file
- `GET /api/v1/download/:id` - Download file

## Supported File Types

### Images
- JPEG, PNG, GIF, WebP
- Max size: 10MB
- Auto-resize to multiple resolutions

### Videos
- MP4, WebM, MOV
- Max size: 500MB
- Generate thumbnails

### Documents
- PDF, DOC, DOCX, PPT, PPTX
- Max size: 50MB

### Audio
- MP3, WAV, OGG
- Max size: 100MB

## Environment Variables

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/files_db
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=uploads
MINIO_USE_SSL=false
CDN_URL=https://cdn.example.com
MAX_UPLOAD_SIZE=524288000
```

## Storage Structure

```
uploads/
├── images/
│   ├── original/
│   ├── large/      # 1920x1080
│   ├── medium/     # 1280x720
│   └── thumbnail/  # 320x240
├── videos/
├── documents/
└── audio/
```

## Development

```bash
go mod download
go run cmd/main.go
```

## Upload Example

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer TOKEN" \
  -F "file=@image.jpg"
```

## Security

- File type validation
- Virus scanning (ClamAV integration ready)
- Size limits enforcement
- Signed URLs for private files
- Rate limiting per user
