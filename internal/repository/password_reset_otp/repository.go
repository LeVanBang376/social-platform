package passwordresetotp

import (
	"context"
	"time"

	"social-platform/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	otp *model.PasswordResetOTP,
) error {
	return db.WithContext(ctx).
		Create(otp).
		Error
}

func (r *Repository) FindValid(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
	now time.Time,
) (*model.PasswordResetOTP, error) {
	var otp model.PasswordResetOTP

	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		Order("created_at DESC").
		First(&otp).
		Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *Repository) MarkAsUsed(
	ctx context.Context,
	db *gorm.DB,
	id uuid.UUID,
) error {
	return db.WithContext(ctx).
		Model(&model.PasswordResetOTP{}).
		Where("id = ?", id).
		Update("used_at", time.Now()).
		Error
}
