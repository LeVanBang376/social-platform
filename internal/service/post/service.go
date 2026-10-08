package post

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	postRepository "social-platform/internal/repository/post"
	"social-platform/internal/response"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	redisClient *redis.Client
	repository  *postRepository.Repository
}

func NewService(
	db *gorm.DB,
	redisClient *redis.Client,
	repository *postRepository.Repository,
) *Service {
	return &Service{
		db:          db,
		redisClient: redisClient,
		repository:  repository,
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
	redisKey := fmt.Sprintf("post:%d", postID)
	data, err := s.redisClient.Get(ctx, redisKey).Bytes()

	if err == nil {
		var res dto.PostResponse
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	// Cache miss
	post, err := s.repository.FindByID(
		ctx,
		s.db,
		postID,
	)
	if err != nil {
		return nil, err
	}

	res := dto.PostResponse{
		PostID:       post.PostID,
		UserID:       post.UserID,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
	}

	redisData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	_ = s.redisClient.Set(ctx, redisKey, redisData, 10*time.Minute).Err()

	return &res, nil
}

func (s *Service) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	pagination *response.Pagination,
) ([]*dto.PostResponse, error) {
	posts, err := s.repository.FindByUserID(
		ctx,
		s.db,
		userID,
		pagination,
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
			&dto.PostResponse{
				PostID:       post.PostID,
				UserID:       post.UserID,
				Content:      post.Content,
				CreatedAt:    post.CreatedAt,
				LikeCount:    post.LikeCount,
				CommentCount: post.CommentCount,
			},
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
		&post.Post,
	); err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("post:%d", postID)
	if err := s.redisClient.Del(ctx, redisKey).Err(); err != nil {
		log.Printf("failed to delete post cache: %v", err)
	}

	return &dto.PostResponse{
		PostID:       post.PostID,
		UserID:       post.UserID,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
	}, nil
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
