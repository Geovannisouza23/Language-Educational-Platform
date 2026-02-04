package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/language-platform/course-service/internal/model"
	"github.com/language-platform/course-service/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseService interface {
	CreateCourse(ctx context.Context, teacherID uuid.UUID, req *model.CreateCourseRequest) (*model.Course, error)
	GetCourse(ctx context.Context, id uuid.UUID) (*model.CourseResponse, error)
	ListCourses(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*model.CourseResponse, int64, error)
	UpdateCourse(ctx context.Context, id, teacherID uuid.UUID, req *model.UpdateCourseRequest) (*model.Course, error)
	DeleteCourse(ctx context.Context, id, teacherID uuid.UUID) error
	PublishCourse(ctx context.Context, id, teacherID uuid.UUID) error

	CreateLesson(ctx context.Context, courseID, teacherID uuid.UUID, req *model.CreateLessonRequest) (*model.Lesson, error)
	GetCourseLessons(ctx context.Context, courseID uuid.UUID) ([]*model.Lesson, error)
	UpdateLesson(ctx context.Context, lessonID, courseID, teacherID uuid.UUID, req *model.UpdateLessonRequest) (*model.Lesson, error)
	DeleteLesson(ctx context.Context, lessonID, courseID, teacherID uuid.UUID) error

	EnrollInCourse(ctx context.Context, courseID, studentID uuid.UUID) (*model.Enrollment, error)
	UnenrollFromCourse(ctx context.Context, courseID, studentID uuid.UUID) error
	GetMyCourses(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*model.Enrollment, error)
	GetCourseEnrollments(ctx context.Context, courseID, teacherID uuid.UUID, limit, offset int) ([]*model.Enrollment, int64, error)
}

type courseService struct {
	courseRepo     repository.CourseRepository
	lessonRepo     repository.LessonRepository
	enrollmentRepo repository.EnrollmentRepository
	redis          *redis.Client
}

func NewCourseService(
	courseRepo repository.CourseRepository,
	lessonRepo repository.LessonRepository,
	enrollmentRepo repository.EnrollmentRepository,
	redis *redis.Client,
) CourseService {
	return &courseService{
		courseRepo:     courseRepo,
		lessonRepo:     lessonRepo,
		enrollmentRepo: enrollmentRepo,
		redis:          redis,
	}
}

func (s *courseService) CreateCourse(ctx context.Context, teacherID uuid.UUID, req *model.CreateCourseRequest) (*model.Course, error) {
	course := &model.Course{
		TeacherID:    teacherID,
		Title:        req.Title,
		Description:  req.Description,
		Language:     req.Language,
		Level:        req.Level,
		Status:       model.CourseDraft,
		Price:        req.Price,
		Duration:     req.Duration,
		MaxStudents:  req.MaxStudents,
		ThumbnailURL: req.ThumbnailURL,
	}

	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, fmt.Errorf("failed to create course: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"course_id":  course.ID,
		"teacher_id": teacherID,
		"title":      course.Title,
	}).Info("Course created")

	return course, nil
}

func (s *courseService) GetCourse(ctx context.Context, id uuid.UUID) (*model.CourseResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("course:%s", id.String())
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var response model.CourseResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			logrus.Debugf("Course %s loaded from cache", id)
			return &response, nil
		}
	}

	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("course not found")
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	response := &model.CourseResponse{
		ID:            course.ID,
		TeacherID:     course.TeacherID,
		Title:         course.Title,
		Description:   course.Description,
		Language:      course.Language,
		Level:         course.Level,
		Status:        course.Status,
		Price:         course.Price,
		Duration:      course.Duration,
		MaxStudents:   course.MaxStudents,
		ThumbnailURL:  course.ThumbnailURL,
		EnrolledCount: course.EnrolledCount,
		Rating:        course.Rating,
		RatingCount:   course.RatingCount,
		LessonCount:   len(course.Lessons),
		CreatedAt:     course.CreatedAt,
		UpdatedAt:     course.UpdatedAt,
	}

	// Cache the result
	data, _ := json.Marshal(response)
	s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	return response, nil
}

func (s *courseService) ListCourses(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*model.CourseResponse, int64, error) {
	courses, total, err := s.courseRepo.List(ctx, filters, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list courses: %w", err)
	}

	responses := make([]*model.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = &model.CourseResponse{
			ID:            course.ID,
			TeacherID:     course.TeacherID,
			Title:         course.Title,
			Description:   course.Description,
			Language:      course.Language,
			Level:         course.Level,
			Status:        course.Status,
			Price:         course.Price,
			Duration:      course.Duration,
			MaxStudents:   course.MaxStudents,
			ThumbnailURL:  course.ThumbnailURL,
			EnrolledCount: course.EnrolledCount,
			Rating:        course.Rating,
			RatingCount:   course.RatingCount,
			CreatedAt:     course.CreatedAt,
			UpdatedAt:     course.UpdatedAt,
		}
	}

	return responses, total, nil
}

func (s *courseService) UpdateCourse(ctx context.Context, id, teacherID uuid.UUID, req *model.UpdateCourseRequest) (*model.Course, error) {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return nil, fmt.Errorf("unauthorized: not course owner")
	}

	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = *req.Description
	}
	if req.Language != nil {
		course.Language = *req.Language
	}
	if req.Level != nil {
		course.Level = *req.Level
	}
	if req.Price != nil {
		course.Price = *req.Price
	}
	if req.Duration != nil {
		course.Duration = *req.Duration
	}
	if req.MaxStudents != nil {
		course.MaxStudents = *req.MaxStudents
	}
	if req.ThumbnailURL != nil {
		course.ThumbnailURL = *req.ThumbnailURL
	}

	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, fmt.Errorf("failed to update course: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("course:%s", id.String())
	s.redis.Del(ctx, cacheKey)

	return course, nil
}

func (s *courseService) DeleteCourse(ctx context.Context, id, teacherID uuid.UUID) error {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return fmt.Errorf("unauthorized: not course owner")
	}

	if err := s.courseRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("course:%s", id.String())
	s.redis.Del(ctx, cacheKey)

	return nil
}

func (s *courseService) PublishCourse(ctx context.Context, id, teacherID uuid.UUID) error {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return fmt.Errorf("unauthorized: not course owner")
	}

	// Check if course has at least one lesson
	if len(course.Lessons) == 0 {
		return fmt.Errorf("cannot publish course without lessons")
	}

	course.Status = model.CoursePublished
	if err := s.courseRepo.Update(ctx, course); err != nil {
		return fmt.Errorf("failed to publish course: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("course:%s", id.String())
	s.redis.Del(ctx, cacheKey)

	return nil
}

func (s *courseService) CreateLesson(ctx context.Context, courseID, teacherID uuid.UUID, req *model.CreateLessonRequest) (*model.Lesson, error) {
	// Verify course ownership
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return nil, fmt.Errorf("unauthorized: not course owner")
	}

	lesson := &model.Lesson{
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		VideoURL:    req.VideoURL,
		Duration:    req.Duration,
		Order:       req.Order,
		IsPublished: false,
	}

	if err := s.lessonRepo.Create(ctx, lesson); err != nil {
		return nil, fmt.Errorf("failed to create lesson: %w", err)
	}

	// Invalidate course cache
	cacheKey := fmt.Sprintf("course:%s", courseID.String())
	s.redis.Del(ctx, cacheKey)

	return lesson, nil
}

func (s *courseService) GetCourseLessons(ctx context.Context, courseID uuid.UUID) ([]*model.Lesson, error) {
	lessons, err := s.lessonRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}
	return lessons, nil
}

func (s *courseService) UpdateLesson(ctx context.Context, lessonID, courseID, teacherID uuid.UUID, req *model.UpdateLessonRequest) (*model.Lesson, error) {
	// Verify course ownership
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return nil, fmt.Errorf("unauthorized: not course owner")
	}

	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	if lesson.CourseID != courseID {
		return nil, fmt.Errorf("lesson does not belong to this course")
	}

	if req.Title != nil {
		lesson.Title = *req.Title
	}
	if req.Description != nil {
		lesson.Description = *req.Description
	}
	if req.Content != nil {
		lesson.Content = *req.Content
	}
	if req.VideoURL != nil {
		lesson.VideoURL = *req.VideoURL
	}
	if req.Duration != nil {
		lesson.Duration = *req.Duration
	}
	if req.Order != nil {
		lesson.Order = *req.Order
	}
	if req.IsPublished != nil {
		lesson.IsPublished = *req.IsPublished
	}

	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, fmt.Errorf("failed to update lesson: %w", err)
	}

	// Invalidate course cache
	cacheKey := fmt.Sprintf("course:%s", courseID.String())
	s.redis.Del(ctx, cacheKey)

	return lesson, nil
}

func (s *courseService) DeleteLesson(ctx context.Context, lessonID, courseID, teacherID uuid.UUID) error {
	// Verify course ownership
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return fmt.Errorf("unauthorized: not course owner")
	}

	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return fmt.Errorf("failed to get lesson: %w", err)
	}

	if lesson.CourseID != courseID {
		return fmt.Errorf("lesson does not belong to this course")
	}

	if err := s.lessonRepo.Delete(ctx, lessonID); err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}

	// Invalidate course cache
	cacheKey := fmt.Sprintf("course:%s", courseID.String())
	s.redis.Del(ctx, cacheKey)

	return nil
}

func (s *courseService) EnrollInCourse(ctx context.Context, courseID, studentID uuid.UUID) (*model.Enrollment, error) {
	// Check if course exists and is published
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found")
	}

	if course.Status != model.CoursePublished {
		return nil, fmt.Errorf("course is not available for enrollment")
	}

	// Check if already enrolled
	enrolled, err := s.enrollmentRepo.IsEnrolled(ctx, courseID, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check enrollment: %w", err)
	}

	if enrolled {
		return nil, fmt.Errorf("already enrolled in this course")
	}

	// Check max students limit
	if course.MaxStudents > 0 && course.EnrolledCount >= course.MaxStudents {
		return nil, fmt.Errorf("course is full")
	}

	enrollment := &model.Enrollment{
		CourseID:  courseID,
		StudentID: studentID,
		Status:    model.EnrollmentActive,
		Progress:  0,
	}

	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	// Increment enrollment count
	if err := s.courseRepo.IncrementEnrollmentCount(ctx, courseID); err != nil {
		logrus.Errorf("Failed to increment enrollment count: %v", err)
	}

	// Invalidate course cache
	cacheKey := fmt.Sprintf("course:%s", courseID.String())
	s.redis.Del(ctx, cacheKey)

	return enrollment, nil
}

func (s *courseService) UnenrollFromCourse(ctx context.Context, courseID, studentID uuid.UUID) error {
	enrollment, err := s.enrollmentRepo.GetByCourseAndStudent(ctx, courseID, studentID)
	if err != nil {
		return fmt.Errorf("enrollment not found")
	}

	if err := s.enrollmentRepo.Delete(ctx, enrollment.ID); err != nil {
		return fmt.Errorf("failed to unenroll: %w", err)
	}

	// Decrement enrollment count
	if err := s.courseRepo.DecrementEnrollmentCount(ctx, courseID); err != nil {
		logrus.Errorf("Failed to decrement enrollment count: %v", err)
	}

	// Invalidate course cache
	cacheKey := fmt.Sprintf("course:%s", courseID.String())
	s.redis.Del(ctx, cacheKey)

	return nil
}

func (s *courseService) GetMyCourses(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]*model.Enrollment, error) {
	enrollments, err := s.enrollmentRepo.GetByStudentID(ctx, studentID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollments: %w", err)
	}
	return enrollments, nil
}

func (s *courseService) GetCourseEnrollments(ctx context.Context, courseID, teacherID uuid.UUID, limit, offset int) ([]*model.Enrollment, int64, error) {
	// Verify course ownership
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get course: %w", err)
	}

	if course.TeacherID != teacherID {
		return nil, 0, fmt.Errorf("unauthorized: not course owner")
	}

	enrollments, total, err := s.enrollmentRepo.GetByCourseID(ctx, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get enrollments: %w", err)
	}

	return enrollments, total, nil
}
