package comment_like

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
	like *model.CommentLike,
) error {
	return db.
		WithContext(ctx).
		Create(like).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	commentID int64,
	userID uuid.UUID,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.CommentLike{},
			"comment_id = ? AND user_id = ?",
			commentID,
			userID,
		).
		Error
}

func (r *Repository) Find(
	ctx context.Context,
	db *gorm.DB,
	commentID int64,
	userID uuid.UUID,
) (*model.CommentLike, error) {
	var like model.CommentLike

	err := db.
		WithContext(ctx).
		First(
			&like,
			"comment_id = ? AND user_id = ?",
			commentID,
			userID,
		).
		Error

	if err != nil {
		return nil, err
	}

	return &like, nil
}
