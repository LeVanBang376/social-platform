package model

import (
	"time"

	"github.com/google/uuid"
)

type CommentLike struct {
	CommentID int64     `gorm:"column:comment_id;type:bigint;primaryKey" json:"comment_id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`

	// Associations
	Comment Comment `gorm:"foreignKey:CommentID;references:CommentID" json:"comment"`
	User    User    `gorm:"foreignKey:UserID;references:UserID" json:"user"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}
