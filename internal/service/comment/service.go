package comment

import (
	"context"
	"errors"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	commentRepository "social-platform/internal/repository/comment"
	postRepository "social-platform/internal/repository/post"
	"social-platform/internal/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound           = errors.New("post not found")
	ErrParentCommentNotFound  = errors.New("parent comment not found")
	ErrParentCommentWrongPost = errors.New("parent comment does not belong to this post")
	ErrNestedReply            = errors.New("nested replies are not allowed")
	ErrCommentNotFound        = errors.New("comment not found")
	ErrNotCommentOwner        = errors.New("you are not allowed to modify this comment")
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}

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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrParentCommentNotFound
			}

			return nil, err
		}

		// Parent comment must belong to the same post.
		if parentComment.PostID != postID {
			return nil, ErrParentCommentWrongPost
		}

		// Only allow one level of replies.
		if parentComment.ParentCommentID != nil {
			return nil, ErrNestedReply
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

func (s *Service) FindByPostID(
	ctx context.Context,
	postID int64,
	pagination *response.Pagination,
) ([]*dto.CommentResponse, error) {
	comments, err := s.repository.FindByPostID(
		ctx,
		s.db,
		postID,
		pagination,
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
			&dto.CommentResponse{
				CommentID:       comment.CommentID,
				PostID:          comment.PostID,
				UserID:          comment.UserID,
				ParentCommentID: comment.ParentCommentID,
				Content:         comment.Content,
				CreatedAt:       comment.CreatedAt,
				ReplyCount:      comment.ReplyCount,
				LikeCount:       comment.LikeCount,
			},
		)
	}

	return responses, nil
}

func (s *Service) FindReplies(
	ctx context.Context,
	postID int64,
	commentID int64,
	pagination *response.Pagination,
) ([]*dto.CommentResponse, error) {
	comments, err := s.repository.FindReplies(
		ctx,
		s.db,
		postID,
		commentID,
		pagination,
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
			&dto.CommentResponse{
				CommentID:       comment.CommentID,
				PostID:          comment.PostID,
				UserID:          comment.UserID,
				ParentCommentID: comment.ParentCommentID,
				Content:         comment.Content,
				CreatedAt:       comment.CreatedAt,
				LikeCount:       comment.LikeCount,
			},
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}

		return nil, err
	}

	if comment.UserID != userID {
		return nil, ErrNotCommentOwner
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}

		return err
	}

	if comment.UserID != userID {
		return ErrNotCommentOwner
	}

	return s.repository.Delete(
		ctx,
		s.db,
		commentID,
	)
}
