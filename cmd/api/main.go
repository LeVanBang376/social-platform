// @title           Social Platform APIs
// @version         1.0
// @description     Social Platform Backend APIs
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey CookieAuth
// @in              cookie
// @name            access_token

package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"social-platform/infrastructure/db"
	"social-platform/infrastructure/jwt"
	"social-platform/internal/config"
	handler "social-platform/internal/handler"
	"social-platform/internal/middleware"
	"social-platform/internal/repository/follow"
	"social-platform/internal/repository/user"
	"social-platform/internal/repository/user_session"
	"social-platform/internal/routes"
	authService "social-platform/internal/service/auth"
	followService "social-platform/internal/service/follow"
	userService "social-platform/internal/service/user"

	_ "social-platform/docs"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	shutdownTimeout = 15 * time.Second
)

func main() {
	// ============================================
	// Config
	// ============================================

	cfg := config.Load()

	// ============================================
	// Context
	// ============================================

	rootCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// ============================================
	// Database
	// ============================================

	database, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL!")

	// ============================================
	// Infrastructure
	// ============================================

	jwtService := jwt.NewJWTService(cfg.JWTSecret)

	// ============================================
	// Repositories
	// ============================================

	userRepo := user.NewRepository()
	userSessionRepo := user_session.NewRepository()
	followRepo := follow.NewRepository()

	// ============================================
	// Services
	// ============================================

	userSvc := userService.NewService(
		database,
		userRepo,
	)

	authSvc := authService.NewService(
		database,
		userRepo,
		userSessionRepo,
		jwtService,
	)

	followSvc := followService.NewService(database, followRepo)

	// ============================================
	// Handlers
	// ============================================

	userHdl := handler.NewUserHandler(
		userSvc,
	)

	authHdl := handler.NewAuthHandler(
		authSvc,
	)

	followHdl := handler.NewFollowHandler(followSvc)

	// ============================================
	// Middleware
	// ============================================

	authMiddleware := middleware.Auth(
		jwtService,
	)

	// ============================================
	// Gin
	// ============================================

	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	// ============================================
	// Health check
	// ============================================

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// ============================================
	// Swagger
	// ============================================

	router.GET(
		"/swagger/*any",
		gin.WrapH(httpSwagger.WrapHandler),
	)

	// ============================================
	// Routes
	// ============================================

	routes.RegisterAuthRoutes(
		router,
		authHdl,
		authMiddleware,
	)

	// User routes
	routes.RegisterUserRoutes(
		router,
		userHdl,
		authMiddleware,
	)

	// Follow routes
	routes.RegisterFollowRoutes(
		router,
		followHdl,
		authMiddleware,
	)

	// ============================================
	// HTTP Server
	// ============================================

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server starting on :8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ============================================
	// Graceful Shutdown
	// ============================================

	<-rootCtx.Done()

	log.Println("Received shutdown signal, shutting down...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	log.Println("Server shut down gracefully.")
}
