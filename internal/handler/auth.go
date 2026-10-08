package handler

import (
	"errors"
	"net/http"

	"social-platform/internal/dto"
	"social-platform/internal/middleware"
	"social-platform/internal/response"
	authService "social-platform/internal/service/auth"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"

	accessTokenMaxAge  = 15 * 60
	refreshTokenMaxAge = 7 * 24 * 60 * 60
)

type AuthHandler struct {
	service *authService.Service
}

func NewAuthHandler(service *authService.Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and set access and refresh tokens in HttpOnly cookies
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Login credentials"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]interface{}  "Invalid request"
// @Failure      401      {object}  map[string]interface{}  "Invalid email or password"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	res, err := h.service.Login(
		c.Request.Context(),
		&req,
	)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCredentials) {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid email or password",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	h.setAuthCookies(
		c,
		res.AccessToken,
		res.RefreshToken,
	)

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Login successful",
		res.User,
	)
}

// Me godoc
// @Summary      Get current user
// @Description  Get information of the currently authenticated user
// @Tags         Auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.GetClaims(c)

	if claims == nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	user, err := h.service.GetUserByID(
		c.Request.Context(),
		claims.UserID,
	)
	if err != nil {
		if errors.Is(err, authService.ErrUserNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Unauthorized",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get current user successful",
		user,
	)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Rotate the refresh token and issue a new access token
// @Tags         Auth
// @Produce      json
// @Success      200
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Refresh token not found",
		)
		return
	}

	res, err := h.service.Refresh(
		c.Request.Context(),
		refreshToken,
	)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidRefreshToken) ||
			errors.Is(err, authService.ErrRefreshTokenRevoked) ||
			errors.Is(err, authService.ErrRefreshTokenExpired) ||
			errors.Is(err, authService.ErrUserNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Unauthorized",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	h.setAuthCookies(
		c,
		res.AccessToken,
		res.RefreshToken,
	)

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Token refreshed successfully",
	)
}

// Logout godoc
// @Summary      User logout
// @Description  Revoke the current refresh token and clear authentication cookies
// @Tags         Auth
// @Produce      json
// @Success      200      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)

	if err == nil {
		if err := h.service.Logout(
			c.Request.Context(),
			refreshToken,
		); err != nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusInternalServerError,
				"Logout failed",
			)
			return
		}
	}

	h.clearAuthCookies(c)

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Logout successful",
	)
}

// ForgotPassword godoc
// @Summary      Forgot password
// @Description  Send a password reset OTP to the user's email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ForgotPasswordRequest  true  "Forgot password request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}  "Invalid request"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	if err := h.service.ForgotPassword(
		c.Request.Context(),
		&req,
	); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to process forgot password request",
		)
		return
	}

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"If the email exists, a password reset OTP has been sent",
	)
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using a 6-digit OTP
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ResetPasswordRequest  true  "Reset password request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}  "Invalid request"
// @Failure      401      {object}  map[string]interface{}  "Invalid or expired OTP"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	if err := h.service.ResetPassword(
		c.Request.Context(),
		&req,
	); err != nil {
		if errors.Is(err, authService.ErrInvalidEmailOrOTP) {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid email or OTP",
			)
			return
		}

		if errors.Is(err, authService.ErrInvalidOrExpiredOTP) {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid or expired OTP",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Password reset successful",
	)
}

func (h *AuthHandler) setAuthCookies(
	c *gin.Context,
	accessToken string,
	refreshToken string,
) {
	c.SetCookie(
		accessTokenCookie,
		accessToken,
		accessTokenMaxAge,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		refreshTokenCookie,
		refreshToken,
		refreshTokenMaxAge,
		"/",
		"",
		false,
		true,
	)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie(
		accessTokenCookie,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		refreshTokenCookie,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}
