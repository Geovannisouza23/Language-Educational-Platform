package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Role      string         `gorm:"not null" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

type UserProfile struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	FirstName   string         `gorm:"size:100" json:"first_name"`
	LastName    string         `gorm:"size:100" json:"last_name"`
	PhoneNumber string         `gorm:"size:20" json:"phone_number"`
	DateOfBirth *time.Time     `json:"date_of_birth"`
	Country     string         `gorm:"size:50" json:"country"`
	City        string         `gorm:"size:100" json:"city"`
	Bio         string         `gorm:"type:text" json:"bio"`
	AvatarURL   string         `json:"avatar_url"`
	Language    string         `gorm:"size:10;default:'en'" json:"language"`
	Timezone    string         `gorm:"size:50" json:"timezone"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (p *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Profile   *UserProfile `json:"profile,omitempty"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Role     *string `json:"role,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UpdateProfileRequest struct {
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	PhoneNumber *string    `json:"phone_number,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Country     *string    `json:"country,omitempty"`
	City        *string    `json:"city,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	Language    *string    `json:"language,omitempty"`
	Timezone    *string    `json:"timezone,omitempty"`
}
