package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseStatus string
type CourseLevel string

const (
	CourseDraft     CourseStatus = "draft"
	CoursePublished CourseStatus = "published"
	CourseArchived  CourseStatus = "archived"

	LevelBeginner     CourseLevel = "Beginner"
	LevelIntermediate CourseLevel = "Intermediate"
	LevelAdvanced     CourseLevel = "Advanced"
)

type Course struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	TeacherID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Language      string         `gorm:"size:50;not null;index" json:"language"`
	Level         CourseLevel    `gorm:"size:20;not null;index" json:"level"`
	Status        CourseStatus   `gorm:"size:20;default:'draft';index" json:"status"`
	Price         float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	Duration      int            `gorm:"comment:Duration in days" json:"duration"`
	MaxStudents   int            `gorm:"default:0;comment:0 = unlimited" json:"max_students"`
	ThumbnailURL  string         `gorm:"size:500" json:"thumbnail_url"`
	EnrolledCount int            `gorm:"default:0" json:"enrolled_count"`
	Rating        float64        `gorm:"type:decimal(3,2);default:0" json:"rating"`
	RatingCount   int            `gorm:"default:0" json:"rating_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Lessons     []Lesson     `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID" json:"-"`
}

type Lesson struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CourseID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Content     string         `gorm:"type:text" json:"content"`
	VideoURL    string         `gorm:"size:500" json:"video_url"`
	Duration    int            `gorm:"comment:Duration in minutes" json:"duration"`
	Order       int            `gorm:"not null;default:0" json:"order"`
	IsPublished bool           `gorm:"default:false" json:"is_published"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Course Course `gorm:"foreignKey:CourseID" json:"-"`
}

type EnrollmentStatus string

const (
	EnrollmentActive    EnrollmentStatus = "active"
	EnrollmentCompleted EnrollmentStatus = "completed"
	EnrollmentDropped   EnrollmentStatus = "dropped"
)

type Enrollment struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	CourseID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"course_id"`
	StudentID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"student_id"`
	Status       EnrollmentStatus `gorm:"size:20;default:'active'" json:"status"`
	Progress     float64          `gorm:"type:decimal(5,2);default:0" json:"progress"`
	EnrolledAt   time.Time        `json:"enrolled_at"`
	CompletedAt  *time.Time       `json:"completed_at,omitempty"`
	LastAccessAt time.Time        `json:"last_access_at"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`

	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.EnrolledAt.IsZero() {
		e.EnrolledAt = time.Now()
	}
	if e.LastAccessAt.IsZero() {
		e.LastAccessAt = time.Now()
	}
	return nil
}

// DTOs
type CreateCourseRequest struct {
	Title        string      `json:"title" binding:"required,min=3,max=255"`
	Description  string      `json:"description" binding:"required"`
	Language     string      `json:"language" binding:"required"`
	Level        CourseLevel `json:"level" binding:"required,oneof=Beginner Intermediate Advanced"`
	Price        float64     `json:"price" binding:"min=0"`
	Duration     int         `json:"duration" binding:"min=1"`
	MaxStudents  int         `json:"max_students" binding:"min=0"`
	ThumbnailURL string      `json:"thumbnail_url"`
}

type UpdateCourseRequest struct {
	Title        *string      `json:"title,omitempty"`
	Description  *string      `json:"description,omitempty"`
	Language     *string      `json:"language,omitempty"`
	Level        *CourseLevel `json:"level,omitempty"`
	Price        *float64     `json:"price,omitempty"`
	Duration     *int         `json:"duration,omitempty"`
	MaxStudents  *int         `json:"max_students,omitempty"`
	ThumbnailURL *string      `json:"thumbnail_url,omitempty"`
}

type CreateLessonRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description"`
	Content     string `json:"content"`
	VideoURL    string `json:"video_url"`
	Duration    int    `json:"duration" binding:"min=0"`
	Order       int    `json:"order" binding:"min=0"`
}

type UpdateLessonRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Content     *string `json:"content,omitempty"`
	VideoURL    *string `json:"video_url,omitempty"`
	Duration    *int    `json:"duration,omitempty"`
	Order       *int    `json:"order,omitempty"`
	IsPublished *bool   `json:"is_published,omitempty"`
}

type CourseResponse struct {
	ID            uuid.UUID    `json:"id"`
	TeacherID     uuid.UUID    `json:"teacher_id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Language      string       `json:"language"`
	Level         CourseLevel  `json:"level"`
	Status        CourseStatus `json:"status"`
	Price         float64      `json:"price"`
	Duration      int          `json:"duration"`
	MaxStudents   int          `json:"max_students"`
	ThumbnailURL  string       `json:"thumbnail_url"`
	EnrolledCount int          `json:"enrolled_count"`
	Rating        float64      `json:"rating"`
	RatingCount   int          `json:"rating_count"`
	LessonCount   int          `json:"lesson_count"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
