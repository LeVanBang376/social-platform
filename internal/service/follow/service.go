package follow

import (
	"context"
	"errors"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	followRepository "social-platform/internal/repository/follow"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrAlreadyFollowing = errors.New("Already following user")
var ErrNotFollowing = errors.New("Not following user")

type Service struct {
	db         *gorm.DB
	repository *followRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *followRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Follow(
	ctx context.Context,
	followerID uuid.UUID,
	followingID uuid.UUID,
) (*dto.FollowResponse, error) {
	if followerID == followingID {
		return nil, gorm.ErrInvalidData
	}

	follow := &model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		follow,
	); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrAlreadyFollowing
		}

		return nil, err
	}

	return dto.FromFollowModelToResponse(follow), nil
}

func (s *Service) Unfollow(
	ctx context.Context,
	followerID uuid.UUID,
	followingID uuid.UUID,
) error {
	follow, err := s.repository.Find(
		ctx,
		s.db,
		followerID,
		followingID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFollowing
		}

		return err
	}

	return s.repository.Delete(
		ctx,
		s.db,
		follow.FollowerID,
		follow.FollowingID,
	)
}
