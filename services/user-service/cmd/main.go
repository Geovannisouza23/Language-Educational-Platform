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
	"github.com/language-platform/user-service/internal/config"
	"github.com/language-platform/user-service/internal/database"
	"github.com/language-platform/user-service/internal/handler"
	"github.com/language-platform/user-service/internal/middleware"
	"github.com/language-platform/user-service/internal/repository"
	"github.com/language-platform/user-service/internal/service"
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

	logrus.Info("Starting User Service...")

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
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, redisClient)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := setupRouter(cfg, userHandler)

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

func setupRouter(cfg *config.Config, userHandler *handler.UserHandler) *gin.Engine {
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
			"service": "user-service",
			"time":    time.Now().Unix(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/:id/profile", userHandler.GetUserProfile)
			users.PUT("/:id/profile", userHandler.UpdateUserProfile)
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
