package post

import (
	"context"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	postRepository "social-platform/internal/repository/post"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db         *gorm.DB
	repository *postRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *postRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreatePostRequest,
	userID uuid.UUID,
) (*dto.PostResponse, error) {
	post := &model.Post{
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		post,
	); err != nil {
		return nil, err
	}

	return dto.FromPostModelToResponse(post), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	postID int64,
) (*dto.PostResponse, error) {
	post, err := s.repository.FindByID(
		ctx,
		s.db,
		postID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromPostModelToResponse(post), nil
}

func (s *Service) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	query *dto.FindPostsQuery,
) ([]*dto.PostResponse, error) {
	posts, err := s.repository.FindByUserID(
		ctx,
		s.db,
		userID,
		query.Limit,
		query.Offset,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.PostResponse,
		0,
		len(posts),
	)

	for _, post := range posts {
		responses = append(
			responses,
			dto.FromPostModelToResponse(post),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	postID int64,
	userID uuid.UUID,
	req *dto.UpdatePostRequest,
) (*dto.PostResponse, error) {
	post, err := s.repository.FindByID(
		ctx,
		s.db,
		postID,
	)
	if err != nil {
		return nil, err
	}

	// Only the post owner can update it.
	if post.UserID != userID {
		return nil, gorm.ErrInvalidData
	}

	if req.Content != nil {
		post.Content = *req.Content
	}

	if err := s.repository.Update(
		ctx,
		s.db,
		post,
	); err != nil {
		return nil, err
	}

	return dto.FromPostModelToResponse(post), nil
}

func (s *Service) Delete(
	ctx context.Context,
	postID int64,
	userID uuid.UUID,
) error {
	post, err := s.repository.FindByID(
		ctx,
		s.db,
		postID,
	)
	if err != nil {
		return err
	}

	// Only the post owner can delete it.
	if post.UserID != userID {
		return gorm.ErrInvalidData
	}

	return s.repository.Delete(
		ctx,
		s.db,
		post.PostID,
	)
}
