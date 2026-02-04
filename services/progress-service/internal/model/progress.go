package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Progress struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	StudentID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_id"`
	CourseID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
	CompletedLessons  int            `gorm:"default:0" json:"completed_lessons"`
	TotalLessons      int            `gorm:"default:0" json:"total_lessons"`
	CompletedTasks    int            `gorm:"default:0" json:"completed_tasks"`
	TotalTasks        int            `gorm:"default:0" json:"total_tasks"`
	AverageScore      float64        `gorm:"type:decimal(5,2);default:0" json:"average_score"`
	StudyTimeMinutes  int            `gorm:"default:0" json:"study_time_minutes"`
	Streak            int            `gorm:"default:0" json:"streak"`
	LastStudyDate     *time.Time     `json:"last_study_date,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type LessonProgress struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_id"`
	LessonID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"lesson_id"`
	Completed   bool           `gorm:"default:false" json:"completed"`
	TimeSpent   int            `gorm:"default:0;comment:Minutes" json:"time_spent"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Achievement struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	StudentID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_id"`
	Title       string         `gorm:"size:100;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	IconURL     string         `gorm:"size:500" json:"icon_url"`
	UnlockedAt  time.Time      `json:"unlocked_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Progress) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (l *LessonProgress) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (a *Achievement) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
