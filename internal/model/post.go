package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostID    int64     `gorm:"column:post_id;type:bigint;primaryKey;autoIncrement" json:"post_id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	Content   string    `gorm:"column:content;type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`

	// Associations
	User     User       `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Comments []Comment  `gorm:"foreignKey:PostID;references:PostID" json:"comments"`
	Likes    []PostLike `gorm:"foreignKey:PostID;references:PostID" json:"likes"`
}

func (Post) TableName() string {
	return "posts"
}
