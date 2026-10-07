package dto

import (
	"time"

	"github.com/google/uuid"

	"social-platform/internal/model"
)

type CreateFollowRequest struct {
	FollowingID uuid.UUID `json:"following_id" binding:"required"`
}

type FollowResponse struct {
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func FromFollowModelToResponse(follow *model.Follow) *FollowResponse {
	return &FollowResponse{
		FollowerID:  follow.FollowerID,
		FollowingID: follow.FollowingID,
		CreatedAt:   follow.CreatedAt,
	}
}
