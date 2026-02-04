package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string
type NotificationStatus string

const (
	NotificationEmail    NotificationType = "email"
	NotificationPush     NotificationType = "push"
	NotificationInApp    NotificationType = "in_app"

	NotificationPending  NotificationStatus = "pending"
	NotificationSent     NotificationStatus = "sent"
	NotificationFailed   NotificationStatus = "failed"
	NotificationRead     NotificationStatus = "read"
)

type Notification struct {
	ID        uuid.UUID          `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      NotificationType   `gorm:"size:20;not null" json:"type"`
	Status    NotificationStatus `gorm:"size:20;default:'pending'" json:"status"`
	Title     string             `gorm:"size:255;not null" json:"title"`
	Message   string             `gorm:"type:text;not null" json:"message"`
	Data      string             `gorm:"type:jsonb" json:"data,omitempty"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
	ReadAt    *time.Time         `json:"read_at,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type EmailTemplate struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Subject   string    `gorm:"size:255;not null" json:"subject"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

type SendNotificationRequest struct {
	UserID  uuid.UUID        `json:"user_id" binding:"required"`
	Type    NotificationType `json:"type" binding:"required"`
	Title   string           `json:"title" binding:"required"`
	Message string           `json:"message" binding:"required"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
