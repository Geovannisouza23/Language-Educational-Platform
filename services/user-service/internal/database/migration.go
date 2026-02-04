package database

import (
	"github.com/language-platform/user-service/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserProfile{},
	)
}
