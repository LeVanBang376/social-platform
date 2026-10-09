package post_like

import (
	"context"
	"errors"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	redisinvalidate "social-platform/internal/redis_invalidate"
	postLikeRepository "social-platform/internal/repository/post_like"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrAlreadyLiked = errors.New("post already liked")
	ErrNotLiked     = errors.New("post like not found")
)

type Service struct {
	db          *gorm.DB
	redisClient *redis.Client
	repository  *postLikeRepository.Repository
}

func NewService(
	db *gorm.DB,
	redisClient *redis.Client,
	repository *postLikeRepository.Repository,
) *Service {
	return &Service{
		db:          db,
		redisClient: redisClient,
		repository:  repository,
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrAlreadyLiked
		}

		return nil, err
	}

	// Invalidate post detail cache after liking.
	redisinvalidate.InvalidatePostDetailCache(
		ctx,
		s.redisClient,
		postID,
	)

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotLiked
		}

		return err
	}

	if err := s.repository.Delete(
		ctx,
		s.db,
		like.PostID,
		like.UserID,
	); err != nil {
		return err
	}

	// Invalidate post detail cache after unliking.
	redisinvalidate.InvalidatePostDetailCache(
		ctx,
		s.redisClient,
		postID,
	)

	return nil
}
