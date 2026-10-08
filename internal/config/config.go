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

	return cfg
}
