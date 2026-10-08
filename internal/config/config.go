package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	RedisAddr      string
	JWTSecret      string
	AllowedOrigins []string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	cfg := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.RedisAddr == "" {
		log.Fatal("REDIS_ADDR is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	if cfg.SMTPHost == "" {
		log.Fatal("SMTP_HOST is required")
	}

	if cfg.SMTPPort == "" {
		log.Fatal("SMTP_PORT is required")
	}

	if cfg.SMTPUsername == "" {
		log.Fatal("SMTP_USERNAME is required")
	}

	if cfg.SMTPPassword == "" {
		log.Fatal("SMTP_PASSWORD is required")
	}

	return cfg
}
