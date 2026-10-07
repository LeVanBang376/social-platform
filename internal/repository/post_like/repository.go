package post_like

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
	like *model.PostLike,
) error {
	return db.
		WithContext(ctx).
		Create(like).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	postID int64,
	userID uuid.UUID,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.PostLike{},
			"post_id = ? AND user_id = ?",
			postID,
			userID,
		).
		Error
}

func (r *Repository) Find(
	ctx context.Context,
	db *gorm.DB,
	postID int64,
	userID uuid.UUID,
) (*model.PostLike, error) {
	var like model.PostLike

	err := db.
		WithContext(ctx).
		First(
			&like,
			"post_id = ? AND user_id = ?",
			postID,
			userID,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &like, nil
}
