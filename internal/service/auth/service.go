package auth

import (
	"context"
	"errors"
	"time"

	"social-platform/infrastructure/jwt"
	"social-platform/infrastructure/token"
	"social-platform/internal/dto"
	"social-platform/internal/model"
	userRepository "social-platform/internal/repository/user"
	userSessionRepository "social-platform/internal/repository/user_session"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const refreshTokenDuration = 7 * 24 * time.Hour

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *dto.UserResponse
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	db                *gorm.DB
	userRepository    *userRepository.Repository
	sessionRepository *userSessionRepository.Repository
	jwtService        *jwt.JWTService
}

func NewService(
	db *gorm.DB,
	userRepository *userRepository.Repository,
	sessionRepository *userSessionRepository.Repository,
	jwtService *jwt.JWTService,
) *Service {
	return &Service{
		db:                db,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		jwtService:        jwtService,
	}
}

func (s *Service) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(
		ctx,
		s.db,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return dto.FromUserModelToResponse(user), nil
}

func (s *Service) Login(
	ctx context.Context,
	req *dto.LoginRequest,
) (*LoginResult, error) {
	user, err := s.userRepository.FindByEmail(
		ctx,
		s.db,
		req.Email,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.jwtService.GenerateToken(
		user.UserID,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	session := &model.UserSession{
		ID:               uuid.New(),
		UserID:           user.UserID,
		RefreshTokenHash: token.HashRefreshToken(refreshToken),
		ExpiresAt:        time.Now().Add(refreshTokenDuration),
	}

	if err := s.sessionRepository.Create(
		ctx,
		s.db,
		session,
	); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.FromUserModelToResponse(user),
	}, nil
}

func (s *Service) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	refreshTokenHash := token.HashRefreshToken(refreshToken)

	session, err := s.sessionRepository.FindByRefreshTokenHash(
		ctx,
		s.db,
		refreshTokenHash,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Already logged out / invalid token.
			return nil
		}

		return err
	}

	if session.RevokedAt != nil {
		return nil
	}

	return s.sessionRepository.Revoke(
		ctx,
		s.db,
		session.ID,
	)
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (*RefreshResult, error) {
	refreshTokenHash := token.HashRefreshToken(refreshToken)

	session, err := s.sessionRepository.FindByRefreshTokenHash(
		ctx,
		s.db,
		refreshTokenHash,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}

		return nil, err
	}

	if session.RevokedAt != nil {
		return nil, errors.New("refresh token has been revoked")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	user, err := s.userRepository.FindByID(
		ctx,
		s.db,
		session.UserID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	accessToken, err := s.jwtService.GenerateToken(
		user.UserID,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newSession := &model.UserSession{
		ID:               uuid.New(),
		UserID:           user.UserID,
		RefreshTokenHash: token.HashRefreshToken(newRefreshToken),
		ExpiresAt:        time.Now().Add(refreshTokenDuration),
	}

	// Rotate refresh token atomically.
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.sessionRepository.Revoke(
			ctx,
			tx,
			session.ID,
		); err != nil {
			return err
		}

		if err := s.sessionRepository.Create(
			ctx,
			tx,
			newSession,
		); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
