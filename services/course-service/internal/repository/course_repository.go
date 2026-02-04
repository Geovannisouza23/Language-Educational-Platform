package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/language-platform/course-service/internal/model"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error)
	GetByTeacherID(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]*model.Course, error)
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*model.Course, int64, error)
	Update(ctx context.Context, course *model.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementEnrollmentCount(ctx context.Context, courseID uuid.UUID) error
	DecrementEnrollmentCount(ctx context.Context, courseID uuid.UUID) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *courseRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).
		Preload("Lessons", "is_published = ?", true).
		First(&course, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) GetByTeacherID(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]*model.Course, error) {
	var courses []*model.Course
	err := r.db.WithContext(ctx).
		Preload("Lessons").
		Where("teacher_id = ?", teacherID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&courses).Error
	return courses, err
}

func (r *courseRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*model.Course, int64, error) {
	var courses []*model.Course
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Course{})

	// Apply filters
	if language, ok := filters["language"].(string); ok && language != "" {
		query = query.Where("language = ?", language)
	}
	if level, ok := filters["level"].(string); ok && level != "" {
		query = query.Where("level = ?", level)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	} else {
		// Default: only show published courses
		query = query.Where("status = ?", model.CoursePublished)
	}
	if teacherID, ok := filters["teacher_id"].(string); ok && teacherID != "" {
		query = query.Where("teacher_id = ?", teacherID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&courses).Error

	return courses, total, err
}

func (r *courseRepository) Update(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *courseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Course{}, "id = ?", id).Error
}

func (r *courseRepository) IncrementEnrollmentCount(ctx context.Context, courseID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Course{}).
		Where("id = ?", courseID).
		UpdateColumn("enrolled_count", gorm.Expr("enrolled_count + ?", 1)).
		Error
}

func (r *courseRepository) DecrementEnrollmentCount(ctx context.Context, courseID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Course{}).
		Where("id = ?", courseID).
		Where("enrolled_count > ?", 0).
		UpdateColumn("enrolled_count", gorm.Expr("enrolled_count - ?", 1)).
		Error
}
