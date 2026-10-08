package main

import (
	"context"
	"log"
	"os/signal"
	"strconv"
	"syscall"

	"social-platform/infrastructure/redis"
	"social-platform/internal/config"
	"social-platform/internal/service/email"
	"social-platform/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()

	redisClient := redis.NewClient(
		cfg.RedisAddr,
	)
	defer redisClient.Close()

	smtpPort, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		log.Fatalf("invalid SMTP_PORT: %v", err)
	}

	emailService := email.NewEmailService(
		cfg.SMTPHost,
		smtpPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)

	emailWorker := worker.NewEmailWorker(
		redisClient,
		emailService,
	)

	// Start workers. Email 3 workers
	emailWorkers := 3

	for i := 1; i <= emailWorkers; i++ {
		consumerName := "email-worker-" + strconv.Itoa(i)

		go emailWorker.Start(
			ctx,
			consumerName,
		)
	}

	log.Printf(
		"email workers started: %d",
		emailWorkers,
	)

	// ============================================
	// Wait for Shutdown Signal
	// ============================================

	<-ctx.Done()

	log.Println("shutting down email workers...")
}
