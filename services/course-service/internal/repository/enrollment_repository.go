package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/language-platform/course-service/internal/model"
	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *model.Enrollment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Enrollment, error)
	GetByCourseAndStudent(ctx context.Context, courseID, studentID uuid.UUID) (*model.Enrollment, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*model.Enrollment, error)
	GetByCourseID(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]*model.Enrollment, int64, error)
	Update(ctx context.Context, enrollment *model.Enrollment) error
	Delete(ctx context.Context, id uuid.UUID) error
	IsEnrolled(ctx context.Context, courseID, studentID uuid.UUID) (bool, error)
}

type enrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *model.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *enrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	err := r.db.WithContext(ctx).
		Preload("Course").
		First(&enrollment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetByCourseAndStudent(ctx context.Context, courseID, studentID uuid.UUID) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	err := r.db.WithContext(ctx).
		Where("course_id = ? AND student_id = ?", courseID, studentID).
		First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) GetByStudentID(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*model.Enrollment, error) {
	var enrollments []*model.Enrollment
	err := r.db.WithContext(ctx).
		Preload("Course").
		Where("student_id = ?", studentID).
		Where("status = ?", model.EnrollmentActive).
		Limit(limit).
		Offset(offset).
		Order("enrolled_at DESC").
		Find(&enrollments).Error
	return enrollments, err
}

func (r *enrollmentRepository) GetByCourseID(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]*model.Enrollment, int64, error) {
	var enrollments []*model.Enrollment
	var total int64

	query := r.db.WithContext(ctx).
		Model(&model.Enrollment{}).
		Where("course_id = ?", courseID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Limit(limit).
		Offset(offset).
		Order("enrolled_at DESC").
		Find(&enrollments).Error

	return enrollments, total, err
}

func (r *enrollmentRepository) Update(ctx context.Context, enrollment *model.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}

func (r *enrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Enrollment{}, "id = ?", id).Error
}

func (r *enrollmentRepository) IsEnrolled(ctx context.Context, courseID, studentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Enrollment{}).
		Where("course_id = ? AND student_id = ? AND status = ?", courseID, studentID, model.EnrollmentActive).
		Count(&count).Error
	return count > 0, err
}
