package model

import (
	"time"

	"github.com/google/uuid"
)

type PostLike struct {
	PostID    int64     `gorm:"column:post_id;type:bigint;primaryKey" json:"post_id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`

	// Associations
	Post Post `gorm:"foreignKey:PostID;references:PostID" json:"post"`
	User User `gorm:"foreignKey:UserID;references:UserID" json:"user"`
}

func (PostLike) TableName() string {
	return "post_likes"
}
