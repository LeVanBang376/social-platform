package dto

import (
	"time"

	"github.com/google/uuid"

	"social-platform/internal/model"
)

type FindUsersQuery struct {
	Email    *string `form:"email"`
	FullName *string `form:"full_name"`
}

type CreateUserRequest struct {
	Email       string     `json:"email" binding:"required,email,max=254"`
	Password    string     `json:"password" binding:"required,min=8,max=100"`
	FullName    string     `json:"full_name" binding:"required,max=100"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type UpdateUserRequest struct {
	FullName    *string    `json:"full_name" binding:"omitempty,max=100"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type UserResponse struct {
	UserID      uuid.UUID  `json:"user_id"`
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at"`
}

func FromUserModelToResponse(user *model.User) *UserResponse {
	return &UserResponse{
		UserID:      user.UserID,
		Email:       user.Email,
		FullName:    user.FullName,
		DateOfBirth: user.DateOfBirth,
		CreatedAt:   user.CreatedAt,
	}
}
