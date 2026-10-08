package model

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetOTP struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	OTPHash   string     `gorm:"type:varchar(64);not null"`
	ExpiresAt time.Time  `gorm:"not null;index"`
	UsedAt    *time.Time `gorm:""`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
}
