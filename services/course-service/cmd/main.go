package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/language-platform/course-service/internal/config"
	"github.com/language-platform/course-service/internal/database"
	"github.com/language-platform/course-service/internal/handler"
	"github.com/language-platform/course-service/internal/middleware"
	"github.com/language-platform/course-service/internal/repository"
	"github.com/language-platform/course-service/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	setupLogger(cfg.LogLevel)

	logrus.Info("Starting Course Service...")

	// Initialize database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	courseRepo := repository.NewCourseRepository(db)
	lessonRepo := repository.NewLessonRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)

	// Initialize services
	courseService := service.NewCourseService(courseRepo, lessonRepo, enrollmentRepo, redisClient)

	// Initialize handlers
	courseHandler := handler.NewCourseHandler(courseService)

	// Setup router
	router := setupRouter(cfg, courseHandler)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logrus.Infof("Server listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

func setupRouter(cfg *config.Config, courseHandler *handler.CourseHandler) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "course-service",
			"time":    time.Now().Unix(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		courses := api.Group("/courses")
		{
			// Public endpoints
			courses.GET("", courseHandler.ListCourses)
			courses.GET("/:id", courseHandler.GetCourse)
			courses.GET("/:id/lessons", courseHandler.GetCourseLessons)

			// Protected endpoints (require authentication)
			protected := courses.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				// Teacher endpoints
				teacher := protected.Group("")
				teacher.Use(middleware.RoleMiddleware("Teacher", "Admin"))
				{
					teacher.POST("", courseHandler.CreateCourse)
					teacher.PUT("/:id", courseHandler.UpdateCourse)
					teacher.DELETE("/:id", courseHandler.DeleteCourse)
					teacher.POST("/:id/lessons", courseHandler.CreateLesson)
					teacher.PUT("/:id/lessons/:lessonId", courseHandler.UpdateLesson)
					teacher.DELETE("/:id/lessons/:lessonId", courseHandler.DeleteLesson)
					teacher.PUT("/:id/publish", courseHandler.PublishCourse)
				}

				// Student endpoints
				protected.POST("/:id/enroll", courseHandler.EnrollInCourse)
				protected.DELETE("/:id/unenroll", courseHandler.UnenrollFromCourse)
				protected.GET("/:id/enrollments", courseHandler.GetCourseEnrollments)
				protected.GET("/my-courses", courseHandler.GetMyCourses)
			}
		}
	}

	return router
}

func setupLogger(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
