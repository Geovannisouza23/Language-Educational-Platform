package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/language-platform/course-service/internal/model"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *model.Lesson) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Lesson, error)
	GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*model.Lesson, error)
	Update(ctx context.Context, lesson *model.Lesson) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) Create(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *lessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Lesson, error) {
	var lesson model.Lesson
	err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *lessonRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]*model.Lesson, error) {
	var lessons []*model.Lesson
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Where("is_published = ?", true).
		Order("`order` ASC").
		Find(&lessons).Error
	return lessons, err
}

func (r *lessonRepository) Update(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Save(lesson).Error
}

func (r *lessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Lesson{}, "id = ?", id).Error
}
