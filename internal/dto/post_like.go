package dto

import (
	"time"

	"github.com/google/uuid"

	"social-platform/internal/model"
)

type CreatePostLikeRequest struct {
	PostID int64 `json:"post_id" binding:"required"`
}

type PostLikeResponse struct {
	PostID    int64     `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func FromPostLikeModelToResponse(
	like *model.PostLike,
) *PostLikeResponse {
	return &PostLikeResponse{
		PostID:    like.PostID,
		UserID:    like.UserID,
		CreatedAt: like.CreatedAt,
	}
}
