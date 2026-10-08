package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"social-platform/infrastructure/jwt"
	"social-platform/infrastructure/token"
	"social-platform/internal/dto"
	"social-platform/internal/model"
	passwordResetOTPRepository "social-platform/internal/repository/password_reset_otp"
	userRepository "social-platform/internal/repository/user"
	userSessionRepository "social-platform/internal/repository/user_session"
	"social-platform/internal/service/email"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	db                         *gorm.DB
	userRepository             *userRepository.Repository
	sessionRepository          *userSessionRepository.Repository
	passwordResetOTPRepository *passwordResetOTPRepository.Repository
	redisClient                *redis.Client
	jwtService                 *jwt.JWTService
}

func NewService(
	db *gorm.DB,
	userRepository *userRepository.Repository,
	sessionRepository *userSessionRepository.Repository,
	passwordResetOTPRepository *passwordResetOTPRepository.Repository,
	redisClient *redis.Client,
	jwtService *jwt.JWTService,
) *Service {
	return &Service{
		db:                         db,
		userRepository:             userRepository,
		sessionRepository:          sessionRepository,
		passwordResetOTPRepository: passwordResetOTPRepository,
		redisClient:                redisClient,
		jwtService:                 jwtService,
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

func (s *Service) ForgotPassword(
	ctx context.Context,
	req *dto.ForgotPasswordRequest,
) error {
	user, err := s.userRepository.FindByEmail(
		ctx,
		s.db,
		req.Email,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Không reveal user có tồn tại hay không.
			return nil
		}

		return err
	}

	// Generate 6-digit OTP
	n, err := rand.Int(
		rand.Reader,
		big.NewInt(1000000),
	)
	if err != nil {
		return err
	}

	otp := fmt.Sprintf("%06d", n)

	// Hash OTP
	hash := sha256.Sum256([]byte(otp))
	otpHash := hex.EncodeToString(hash[:])

	// Save OTP
	resetOTP := &model.PasswordResetOTP{
		ID:        uuid.New(),
		UserID:    user.UserID,
		OTPHash:   otpHash,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.passwordResetOTPRepository.Create(
		ctx,
		s.db,
		resetOTP,
	); err != nil {
		return err
	}

	// Create email message
	message := email.EmailMessage{
		To:      user.Email,
		Subject: "Password Reset OTP",
		Body: fmt.Sprintf(
			"Your password reset OTP is: %s\n\n"+
				"This OTP will expire in 5 minutes.",
			otp,
		),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish email job to Redis Stream
	if err := s.redisClient.XAdd(
		ctx,
		&redis.XAddArgs{
			Stream: "email_jobs",
			Values: map[string]interface{}{
				"data": string(data),
			},
		},
	).Err(); err != nil {
		return err
	}

	return nil
}

func (s *Service) ResetPassword(
	ctx context.Context,
	req *dto.ResetPasswordRequest,
) error {
	// Find user
	user, err := s.userRepository.FindByEmail(
		ctx,
		s.db,
		req.Email,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid email or OTP")
		}

		return err
	}

	// Hash OTP provided by user
	hash := sha256.Sum256([]byte(req.OTP))
	otpHash := hex.EncodeToString(hash[:])

	// Find valid OTP
	resetOTP, err := s.passwordResetOTPRepository.FindValid(
		ctx,
		s.db,
		user.UserID,
		time.Now(),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired OTP")
		}

		return err
	}

	// Compare OTP hash
	if otpHash != resetOTP.OTPHash {
		return errors.New("invalid or expired OTP")
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// Update password and mark OTP as used atomically
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.UpdatePasswordHash(
			ctx,
			tx,
			user.UserID,
			string(passwordHash),
		); err != nil {
			return err
		}

		if err := s.passwordResetOTPRepository.MarkAsUsed(
			ctx,
			tx,
			resetOTP.ID,
		); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
