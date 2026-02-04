package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/language-platform/course-service/internal/model"
	"github.com/language-platform/course-service/internal/service"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(service service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.service.CreateCourse(c.Request.Context(), teacherID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	course, err := h.service.GetCourse(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	// Build filters
	filters := make(map[string]interface{})
	if language := c.Query("language"); language != "" {
		filters["language"] = language
	}
	if level := c.Query("level"); level != "" {
		filters["level"] = level
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if teacherID := c.Query("teacher_id"); teacherID != "" {
		filters["teacher_id"] = teacherID
	}

	courses, total, err := h.service.ListCourses(c.Request.Context(), filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.service.UpdateCourse(c.Request.Context(), id, teacherID, &req)
	if err != nil {
		if err.Error() == "unauthorized: not course owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.DeleteCourse(c.Request.Context(), id, teacherID); err != nil {
		if err.Error() == "unauthorized: not course owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func (h *CourseHandler) PublishCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.PublishCourse(c.Request.Context(), id, teacherID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course published successfully"})
}

func (h *CourseHandler) CreateLesson(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson, err := h.service.CreateLesson(c.Request.Context(), courseID, teacherID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lesson)
}

func (h *CourseHandler) GetCourseLessons(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	lessons, err := h.service.GetCourseLessons(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}

func (h *CourseHandler) UpdateLesson(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson, err := h.service.UpdateLesson(c.Request.Context(), lessonID, courseID, teacherID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lesson)
}

func (h *CourseHandler) DeleteLesson(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	lessonID, err := uuid.Parse(c.Param("lessonId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.DeleteLesson(c.Request.Context(), lessonID, courseID, teacherID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lesson deleted successfully"})
}

func (h *CourseHandler) EnrollInCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	studentID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	enrollment, err := h.service.EnrollInCourse(c.Request.Context(), courseID, studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Successfully enrolled",
		"enrollment": enrollment,
	})
}

func (h *CourseHandler) UnenrollFromCourse(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	studentID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.UnenrollFromCourse(c.Request.Context(), courseID, studentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unenrolled"})
}

func (h *CourseHandler) GetMyCourses(c *gin.Context) {
	studentID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	enrollments, err := h.service.GetMyCourses(c.Request.Context(), studentID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
}

func (h *CourseHandler) GetCourseEnrollments(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	teacherID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	enrollments, total, err := h.service.GetCourseEnrollments(c.Request.Context(), courseID, teacherID, limit, offset)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments": enrollments,
		"total":       total,
		"limit":       limit,
		"offset":      offset,
	})
}

// Helper function to get user ID from context
func getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, nil
	}

	// Try to parse as UUID directly
	if userID, ok := userIDStr.(uuid.UUID); ok {
		return userID, nil
	}

	// Try to parse as string
	if userIDStr, ok := userIDStr.(string); ok {
		return uuid.Parse(userIDStr)
	}

	return uuid.Nil, nil
}
