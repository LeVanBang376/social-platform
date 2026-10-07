package post

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
	post *model.Post,
) error {
	return db.
		WithContext(ctx).
		Create(post).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	postID int64,
) (*model.Post, error) {
	var post model.Post

	err := db.
		WithContext(ctx).
		First(&post, "post_id = ?", postID).
		Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *Repository) FindByUserID(
	ctx context.Context,
	db *gorm.DB,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]*model.Post, error) {
	var posts []*model.Post

	err := db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).
		Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	post *model.Post,
) error {
	return db.
		WithContext(ctx).
		Save(post).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	postID int64,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Post{},
			"post_id = ?",
			postID,
		).
		Error
}
