package user

import (
	"context"
	"errors"

	"social-platform/internal/dto"
	"social-platform/internal/model"
	userRepository "social-platform/internal/repository/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Service struct {
	db         *gorm.DB
	repository *userRepository.Repository
}

func NewService(
	db *gorm.DB,
	repository *userRepository.Repository,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	req *dto.CreateUserRequest,
) (*dto.UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UserID:       uuid.New(),
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FullName:     req.FullName,
		DateOfBirth:  req.DateOfBirth,
	}

	if err := s.repository.Create(
		ctx,
		s.db,
		user,
	); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailAlreadyExists
		}

		return nil, err
	}

	return dto.FromUserModelToResponse(user), nil
}

func (s *Service) FindByID(
	ctx context.Context,
	userID uuid.UUID,
) (*dto.UserResponse, error) {
	user, err := s.repository.FindByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return dto.FromUserModelToResponse(user), nil
}

func (s *Service) Update(
	ctx context.Context,
	userID uuid.UUID,
	req *dto.UpdateUserRequest,
) (*dto.UserResponse, error) {
	user, err := s.repository.FindByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.DateOfBirth != nil {
		user.DateOfBirth = req.DateOfBirth
	}

	if err := s.repository.Update(
		ctx,
		s.db,
		user,
	); err != nil {
		return nil, err
	}

	return dto.FromUserModelToResponse(user), nil
}
