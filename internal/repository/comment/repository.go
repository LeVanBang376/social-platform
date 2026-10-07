package comment

import (
	"context"

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
	comment *model.Comment,
) error {
	return db.
		WithContext(ctx).
		Create(comment).
		Error
}

func (r *Repository) FindByID(
	ctx context.Context,
	db *gorm.DB,
	commentID int64,
) (*model.Comment, error) {
	var comment model.Comment

	err := db.
		WithContext(ctx).
		First(&comment, "comment_id = ?", commentID).
		Error

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *Repository) FindByPostID(
	ctx context.Context,
	db *gorm.DB,
	postID int64,
	limit int,
	offset int,
) ([]*model.Comment, error) {
	var comments []*model.Comment

	err := db.
		WithContext(ctx).
		Where(
			"post_id = ? AND parent_comment_id IS NULL",
			postID,
		).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).
		Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) FindReplies(
	ctx context.Context,
	db *gorm.DB,
	parentCommentID int64,
) ([]*model.Comment, error) {
	var comments []*model.Comment

	err := db.
		WithContext(ctx).
		Where("parent_comment_id = ?", parentCommentID).
		Order("created_at ASC").
		Find(&comments).
		Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) Update(
	ctx context.Context,
	db *gorm.DB,
	comment *model.Comment,
) error {
	return db.
		WithContext(ctx).
		Save(comment).
		Error
}

func (r *Repository) Delete(
	ctx context.Context,
	db *gorm.DB,
	commentID int64,
) error {
	return db.
		WithContext(ctx).
		Delete(
			&model.Comment{},
			"comment_id = ?",
			commentID,
		).
		Error
}
