package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	CommentID       int64     `gorm:"column:comment_id;type:bigint;primaryKey;autoIncrement" json:"comment_id"`
	PostID          int64     `gorm:"column:post_id;type:bigint;not null" json:"post_id"`
	UserID          uuid.UUID `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	ParentCommentID *int64    `gorm:"column:parent_comment_id;type:bigint" json:"parent_comment_id"`
	Content         string    `gorm:"column:content;type:text;not null" json:"content"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`

	// Associations
	Post          Post          `gorm:"foreignKey:PostID;references:PostID" json:"post"`
	User          User          `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	ParentComment *Comment      `gorm:"foreignKey:ParentCommentID;references:CommentID" json:"parent_comment"`
	Replies       []Comment     `gorm:"foreignKey:ParentCommentID;references:CommentID" json:"replies"`
	Likes         []CommentLike `gorm:"foreignKey:CommentID;references:CommentID" json:"likes"`
}

func (Comment) TableName() string {
	return "comments"
}
