package follow

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"social-platform/internal/model"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(
	ctx context.Context,
	db *gorm.DB,
	follow *model.Follow,
) error {
	return db.
		WithContext(ctx).
		Create(follow).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	followerID uuid.UUID,
	followingID uuid.UUID,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Follow{},
			"follower_id = ? AND following_id = ?",
			followerID,
			followingID,
		).
		Error
}

func (r *Repository) Find(
	ctx context.Context,
	db *gorm.DB,
	followerID uuid.UUID,
	followingID uuid.UUID,
) (*model.Follow, error) {
	var follow model.Follow

	err := db.
		WithContext(ctx).
		First(
			&follow,
			"follower_id = ? AND following_id = ?",
			followerID,
			followingID,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &follow, nil
}
