package user_session

import (
	"context"

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
	session *model.UserSession,
) error {
	return db.WithContext(ctx).
		Create(session).
		Error
}

func (r *Repository) FindByRefreshTokenHash(
	ctx context.Context,
	db *gorm.DB,
	refreshTokenHash string,
) (*model.UserSession, error) {
	var session model.UserSession

	err := db.WithContext(ctx).
		Where(
			"refresh_token_hash = ?",
			refreshTokenHash,
		).
		First(&session).
		Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *Repository) Revoke(
	ctx context.Context,
	db *gorm.DB,
	sessionID uuid.UUID,
) error {
	return db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("id = ?", sessionID).
		Update(
			"revoked_at",
			gorm.Expr("NOW()"),
		).
		Error
}

func (r *Repository) UpdateLastUsedAt(
	ctx context.Context,
	db *gorm.DB,
	sessionID uuid.UUID,
) error {
	return db.WithContext(ctx).
		Model(&model.UserSession{}).
		Where("id = ?", sessionID).
		Update(
			"last_used_at",
			gorm.Expr("NOW()"),
		).
		Error
}
