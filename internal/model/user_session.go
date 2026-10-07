package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;index"`
	RefreshTokenHash string    `gorm:"type:varchar(64);not null;uniqueIndex"`
	ExpiresAt        time.Time `gorm:"not null"`
	RevokedAt        *time.Time
	LastUsedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
