package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "document"
	FileTypeAudio    FileType = "audio"
)

type File struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Filename    string         `gorm:"size:255;not null" json:"filename"`
	OriginalName string        `gorm:"size:255;not null" json:"original_name"`
	FileType    FileType       `gorm:"size:20;not null" json:"file_type"`
	MimeType    string         `gorm:"size:100" json:"mime_type"`
	Size        int64          `gorm:"not null" json:"size"`
	URL         string         `gorm:"size:500;not null" json:"url"`
	BucketName  string         `gorm:"size:100" json:"bucket_name"`
	ObjectKey   string         `gorm:"size:500" json:"object_key"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type UploadResponse struct {
	ID       uuid.UUID `json:"id"`
	Filename string    `json:"filename"`
	URL      string    `json:"url"`
	Size     int64     `json:"size"`
	FileType FileType  `json:"file_type"`
}
