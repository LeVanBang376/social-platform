package comment_like

import (
	"context"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	commentLikeRepository "social-platform/internal/repository/comment_like"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		return err
	}

	return s.repository.Delete(
		ctx,
		s.db,
		like.CommentID,
		like.UserID,
	)
}
