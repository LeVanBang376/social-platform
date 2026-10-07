package post_like

import (
	"context"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	postLikeRepository "social-platform/internal/repository/post_like"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *postLikeRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *postLikeRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Like(
	ctx context.Context,
	postID int64,
	userID uuid.UUID,
) (*dto.PostLikeResponse, error) {
	like := &model.PostLike{
		PostID: postID,
		UserID: userID,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		like,
	); err != nil {
		return nil, err
	}

	return dto.FromPostLikeModelToResponse(like), nil
}

func (s *Service) Unlike(
	ctx context.Context,
	postID int64,
	userID uuid.UUID,
) error {
	like, err := s.repository.Find(
		ctx,
		s.db,
		postID,
		userID,
	)
	if err != nil {
		return err
	}

	return s.repository.Delete(
		ctx,
		s.db,
		like.PostID,
		like.UserID,
	)
}
