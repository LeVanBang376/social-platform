package comment_like

import (
	"context"
	"errors"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	commentLikeRepository "social-platform/internal/repository/comment_like"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAlreadyLiked = errors.New("comment already liked")
	ErrNotLiked     = errors.New("comment like not found")
)

type Service struct {
	db         *gorm.DB
	repository *commentLikeRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *commentLikeRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Like(
	ctx context.Context,
	commentID int64,
	userID uuid.UUID,
) (*dto.CommentLikeResponse, error) {
	like := &model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		like,
	); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrAlreadyLiked
		}

		return nil, err
	}

	return dto.FromCommentLikeModelToResponse(like), nil
}

func (s *Service) Unlike(
	ctx context.Context,
	commentID int64,
	userID uuid.UUID,
) error {
	like, err := s.repository.Find(
		ctx,
		s.db,
		commentID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotLiked
		}

		return err
	}

	return s.repository.Delete(
		ctx,
		s.db,
		like.CommentID,
		like.UserID,
	)
}
