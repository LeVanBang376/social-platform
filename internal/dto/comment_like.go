package dto

import (
	"time"

	"github.com/google/uuid"

	"social-platform/internal/model"
)

type CreateCommentLikeRequest struct {
	CommentID int64 `json:"comment_id" binding:"required"`
}

type CommentLikeResponse struct {
	CommentID int64     `json:"comment_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func FromCommentLikeModelToResponse(
	like *model.CommentLike,
) *CommentLikeResponse {
	return &CommentLikeResponse{
		CommentID: like.CommentID,
		UserID:    like.UserID,
		CreatedAt: like.CreatedAt,
	}
}
