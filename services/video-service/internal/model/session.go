package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionStatus string

const (
	SessionScheduled SessionStatus = "scheduled"
	SessionLive      SessionStatus = "live"
	SessionEnded     SessionStatus = "ended"
	SessionCancelled SessionStatus = "cancelled"
)

type VideoSession struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CourseID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
	TeacherID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	Title        string         `gorm:"size:255;not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	ScheduledAt  time.Time      `gorm:"not null" json:"scheduled_at"`
	Duration     int            `gorm:"comment:Minutes" json:"duration"`
	Status       SessionStatus  `gorm:"size:20;default:'scheduled'" json:"status"`
	MeetingURL   string         `gorm:"size:500" json:"meeting_url"`
	RecordingURL string         `gorm:"size:500" json:"recording_url,omitempty"`
	StartedAt    *time.Time     `json:"started_at,omitempty"`
	EndedAt      *time.Time     `json:"ended_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Participants []Participant `gorm:"foreignKey:SessionID" json:"participants,omitempty"`
}

type Participant struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	SessionID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"session_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	JoinedAt   *time.Time     `json:"joined_at,omitempty"`
	LeftAt     *time.Time     `json:"left_at,omitempty"`
	Duration   int            `gorm:"default:0;comment:Minutes" json:"duration"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (v *VideoSession) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}

func (p *Participant) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type CreateSessionRequest struct {
	CourseID    uuid.UUID `json:"course_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	Duration    int       `json:"duration" binding:"required,min=15"`
}
