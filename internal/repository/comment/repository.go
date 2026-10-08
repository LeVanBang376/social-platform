package comment

import (
	"context"

	"gorm.io/gorm"

	"social-platform/internal/model"
	"social-platform/internal/response"
)

type Repository struct{}

type CommentWithCounts struct {
	model.Comment
	LikeCount  int64 `gorm:"column:like_count"`
	ReplyCount int64 `gorm:"column:reply_count"`
}

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
	pagination *response.Pagination,
) ([]*CommentWithCounts, error) {
	var comments []*CommentWithCounts

	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Comment{}).
		Where(
			"post_id = ? AND parent_comment_id IS NULL",
			postID,
		).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	pagination.SetTotal(total)

	if err := db.
		WithContext(ctx).
		Table("comments AS c").
		Select(`
			c.*,
			(
				SELECT COUNT(*)
				FROM comment_likes AS cl
				WHERE cl.comment_id = c.comment_id
			) AS like_count,
			(
				SELECT COUNT(*)
				FROM comments AS r
				WHERE r.parent_comment_id = c.comment_id
			) AS reply_count
		`).
		Where(
			"c.post_id = ? AND c.parent_comment_id IS NULL",
			postID,
		).
		Order("c.created_at DESC").
		Limit(pagination.PerPage).
		Offset(pagination.Offset()).
		Scan(&comments).
		Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) FindReplies(
	ctx context.Context,
	db *gorm.DB,
	postID int64,
	parentCommentID int64,
	pagination *response.Pagination,
) ([]*CommentWithCounts, error) {
	var comments []*CommentWithCounts

	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Comment{}).
		Where(
			"post_id = ? AND parent_comment_id = ?",
			postID,
			parentCommentID,
		).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	pagination.SetTotal(total)

	if err := db.
		WithContext(ctx).
		Table("comments AS c").
		Select(`
			c.*,
			(
				SELECT COUNT(*)
				FROM comment_likes AS cl
				WHERE cl.comment_id = c.comment_id
			) AS like_count
		`).
		Where(
			"c.post_id = ? AND c.parent_comment_id = ?",
			postID,
			parentCommentID,
		).
		Order("c.created_at ASC").
		Limit(pagination.PerPage).
		Offset(pagination.Offset()).
		Scan(&comments).
		Error; err != nil {
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
