package comment

import (
	"context"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	commentRepository "social-platform/internal/repository/comment"
	postRepository "social-platform/internal/repository/post"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db             *gorm.DB
	repository     *commentRepository.Repository
	postRepository *postRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *commentRepository.Repository,
	postRepository *postRepository.Repository,
) *Service {
	return &Service{
		db:             db,
		repository:     repository,
		postRepository: postRepository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	postID int64,
	userID uuid.UUID,
	req *dto.CreateCommentRequest,
) (*dto.CommentResponse, error) {
	// Make sure the post exists.
	if _, err := s.postRepository.FindByID(
		ctx,
		s.db,
		postID,
	); err != nil {
		return nil, err
	}

	// If this is a reply, validate the parent comment.
	if req.ParentCommentID != nil {
		parentComment, err := s.repository.FindByID(
			ctx,
			s.db,
			*req.ParentCommentID,
		)
		if err != nil {
			return nil, err
		}

		// Parent comment must belong to the same post.
		if parentComment.PostID != postID {
			return nil, gorm.ErrInvalidData
		}

		// Only allow one level of replies.
		if parentComment.ParentCommentID != nil {
			return nil, gorm.ErrInvalidData
		}
	}

	comment := &model.Comment{
		PostID:          postID,
		UserID:          userID,
		Content:         req.Content,
		ParentCommentID: req.ParentCommentID,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		comment,
	); err != nil {
		return nil, err
	}

	return dto.FromCommentModelToResponse(comment), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	commentID int64,
) (*dto.CommentResponse, error) {
	comment, err := s.repository.FindByID(
		ctx,
		s.db,
		commentID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromCommentModelToResponse(comment), nil
}

func (s *Service) FindByPostID(
	ctx context.Context,
	postID int64,
	query *dto.FindCommentsQuery,
) ([]*dto.CommentResponse, error) {
	comments, err := s.repository.FindByPostID(
		ctx,
		s.db,
		postID,
		query.Limit,
		query.Offset,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.CommentResponse,
		0,
		len(comments),
	)

	for _, comment := range comments {
		responses = append(
			responses,
			dto.FromCommentModelToResponse(comment),
		)
	}

	return responses, nil
}

func (s *Service) FindReplies(
	ctx context.Context,
	commentID int64,
) ([]*dto.CommentResponse, error) {
	comments, err := s.repository.FindReplies(
		ctx,
		s.db,
		commentID,
	)
	if err != nil {
		return nil, err
	}

	responses := make(
		[]*dto.CommentResponse,
		0,
		len(comments),
	)

	for _, comment := range comments {
		responses = append(
			responses,
			dto.FromCommentModelToResponse(comment),
		)
	}

	return responses, nil
}

func (s *Service) Update(
	ctx context.Context,
	commentID int64,
	userID uuid.UUID,
	req *dto.UpdateCommentRequest,
) (*dto.CommentResponse, error) {
	comment, err := s.repository.FindByID(
		ctx,
		s.db,
		commentID,
	)
	if err != nil {
		return nil, err
	}

	if comment.UserID != userID {
		return nil, gorm.ErrInvalidData
	}

	if req.Content != nil {
		comment.Content = *req.Content
	}

	if err := s.repository.Update(
		ctx,
		s.db,
		comment,
	); err != nil {
		return nil, err
	}

	return dto.FromCommentModelToResponse(comment), nil
}

func (s *Service) Delete(
	ctx context.Context,
	commentID int64,
	userID uuid.UUID,
) error {
	comment, err := s.repository.FindByID(
		ctx,
		s.db,
		commentID,
	)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return gorm.ErrInvalidData
	}

	return s.repository.Delete(
		ctx,
		s.db,
		commentID,
	)
}
