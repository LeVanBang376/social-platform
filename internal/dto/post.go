package dto

import (
	"time"

	"github.com/google/uuid"

	"social-platform/internal/model"
)

type FindPostsQuery struct {
	UserID *uuid.UUID `form:"user_id"`
	Limit  int        `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
	Offset int        `form:"offset,default=0" binding:"omitempty,min=0"`
}

type CreatePostRequest struct {
	Content string `json:"content" binding:"required,max=5000"`
}

type UpdatePostRequest struct {
	Content *string `json:"content" binding:"omitempty,max=5000"`
}

type PostResponse struct {
	PostID    int64     `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func FromPostModelToResponse(post *model.Post) *PostResponse {
	return &PostResponse{
		PostID:    post.PostID,
		UserID:    post.UserID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
	}
}
