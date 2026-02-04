package database

import (
	"github.com/language-platform/course-service/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Course{},
		&model.Lesson{},
		&model.Enrollment{},
	)
}
