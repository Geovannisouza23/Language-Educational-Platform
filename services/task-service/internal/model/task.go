package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskType string
type TaskStatus string
type SubmissionStatus string

const (
	TaskTypeAssignment TaskType = "assignment"
	TaskTypeQuiz       TaskType = "quiz"
	TaskTypeEssay      TaskType = "essay"
	TaskTypeSpeaking   TaskType = "speaking"

	TaskDraft     TaskStatus = "draft"
	TaskPublished TaskStatus = "published"

	SubmissionPending   SubmissionStatus = "pending"
	SubmissionSubmitted SubmissionStatus = "submitted"
	SubmissionGraded    SubmissionStatus = "graded"
	SubmissionRevision  SubmissionStatus = "revision"
)

type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CourseID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
	TeacherID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Type        TaskType       `gorm:"size:20;not null" json:"type"`
	Status      TaskStatus     `gorm:"size:20;default:'draft'" json:"status"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	MaxScore    int            `gorm:"default:100" json:"max_score"`
	Content     string         `gorm:"type:text" json:"content"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Submissions []Submission `gorm:"foreignKey:TaskID" json:"submissions,omitempty"`
}

type Submission struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	TaskID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"task_id"`
	StudentID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"student_id"`
	Content      string           `gorm:"type:text" json:"content"`
	FileURL      string           `gorm:"size:500" json:"file_url,omitempty"`
	Status       SubmissionStatus `gorm:"size:20;default:'pending'" json:"status"`
	Score        *int             `json:"score,omitempty"`
	Feedback     string           `gorm:"type:text" json:"feedback,omitempty"`
	SubmittedAt  *time.Time       `json:"submitted_at,omitempty"`
	GradedAt     *time.Time       `json:"graded_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`

	Task Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (s *Submission) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type CreateTaskRequest struct {
	CourseID    uuid.UUID  `json:"course_id" binding:"required"`
	Title       string     `json:"title" binding:"required,min=3"`
	Description string     `json:"description"`
	Type        TaskType   `json:"type" binding:"required"`
	DueDate     *time.Time `json:"due_date"`
	MaxScore    int        `json:"max_score" binding:"min=1"`
	Content     string     `json:"content"`
}

type SubmitTaskRequest struct {
	Content string `json:"content"`
	FileURL string `json:"file_url"`
}

type GradeSubmissionRequest struct {
	Score    int    `json:"score" binding:"required,min=0"`
	Feedback string `json:"feedback"`
}
