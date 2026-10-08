package post

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"social-platform/internal/model"
	"social-platform/internal/response"
)

type Repository struct{}

type PostWithCounts struct {
	model.Post
	LikeCount    int64 `gorm:"column:like_count"`
	CommentCount int64 `gorm:"column:comment_count"`
}

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
) (*PostWithCounts, error) {
	var post PostWithCounts

	err := db.
		WithContext(ctx).
		Table("posts AS p").
		Select(`
			p.*,
			(
				SELECT COUNT(*)
				FROM post_likes AS pl
				WHERE pl.post_id = p.post_id
			) AS like_count,
			(
				SELECT COUNT(*)
				FROM comments AS c
				WHERE c.post_id = p.post_id
			) AS comment_count
		`).
		Where("p.post_id = ?", postID).
		Scan(&post).
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
	pagination *response.Pagination,
) ([]*PostWithCounts, error) {
	var posts []*PostWithCounts

	var total int64

	if err := db.
		WithContext(ctx).
		Model(&model.Post{}).
		Where("user_id = ?", userID).
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	pagination.SetTotal(total)

	if err := db.
		WithContext(ctx).
		Table("posts AS p").
		Select(`
			p.*,
			(
				SELECT COUNT(*)
				FROM post_likes AS pl
				WHERE pl.post_id = p.post_id
			) AS like_count,
			(
				SELECT COUNT(*)
				FROM comments AS c
				WHERE c.post_id = p.post_id
			) AS comment_count
		`).
		Where("p.user_id = ?", userID).
		Order("p.created_at DESC").
		Limit(pagination.PerPage).
		Offset(pagination.Offset()).
		Scan(&posts).
		Error; err != nil {
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
