package dto

import (
	"time"

	"github.com/google/uuid"

	"social-platform/internal/model"
)

type FindCommentsQuery struct {
	ParentCommentID *int64 `form:"parent_comment_id"`
	Limit           int    `form:"limit,default=20" binding:"omitempty,min=1,max=100"`
	Offset          int    `form:"offset,default=0" binding:"omitempty,min=0"`
}

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

	Replies []*CommentResponse `json:"replies,omitempty"`
}

func FromCommentModelToResponse(comment *model.Comment) *CommentResponse {
	replies := make(
		[]*CommentResponse,
		0,
		len(comment.Replies),
	)

	for _, reply := range comment.Replies {
		replies = append(
			replies,
			FromCommentModelToResponse(&reply),
		)
	}

	return &CommentResponse{
		CommentID:       comment.CommentID,
		PostID:          comment.PostID,
		UserID:          comment.UserID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		CreatedAt:       comment.CreatedAt,
		Replies:         replies,
	}
}
