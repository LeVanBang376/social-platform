package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"social-platform/infrastructure/db"
	"social-platform/internal/config"

	"github.com/gin-gonic/gin"
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
