package dto

import (
	"social-platform/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	Content         string `json:"content" binding:"required,max=5000"`
	ParentCommentID *int64 `json:"parent_comment_id"`
}

type UpdateCommentRequest struct {
	Content *string `json:"content" binding:"omitempty,max=5000"`
}

type CommentResponse struct {
	CommentID       int64     `json:"comment_id"`
	PostID          int64     `json:"post_id"`
	UserID          uuid.UUID `json:"user_id"`
	ParentCommentID *int64    `json:"parent_comment_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	ReplyCount      int64     `json:"reply_count"`
	LikeCount       int64     `json:"like_count"`
}

func FromCommentModelToResponse(
	comment *model.Comment,
) *CommentResponse {
	return &CommentResponse{
		CommentID:       comment.CommentID,
		PostID:          comment.PostID,
		UserID:          comment.UserID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		CreatedAt:       comment.CreatedAt,
	}
}
