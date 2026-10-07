package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID  `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	Email        string     `gorm:"column:email;type:varchar(254);not null;unique" json:"email"`
	PasswordHash string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	FullName     string     `gorm:"column:full_name;type:varchar(100);not null" json:"full_name"`
	DateOfBirth  *time.Time `gorm:"column:date_of_birth;type:date" json:"date_of_birth"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()" json:"created_at"`

	// Associations
	Posts []Post `gorm:"foreignKey:UserID;references:UserID" json:"posts"`
}

func (User) TableName() string {
	return "users"
}
