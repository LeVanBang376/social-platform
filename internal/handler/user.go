package handler

import (
	"errors"
	"net/http"

	"social-platform/internal/dto"
	"social-platform/internal/response"
	userService "social-platform/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	service *userService.Service
}

func NewUserHandler(service *userService.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create user
// @Description  Create a new user account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUserRequest  true  "User information"
// @Success      201      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]interface{}  "Invalid request"
// @Failure      409      {object}  map[string]interface{}  "Email already exists"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	user, err := h.service.Create(
		c.Request.Context(),
		&req,
	)
	if err != nil {
		if err.Error() == "email already exists" {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to create user",
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusCreated,
		"User created successfully",
		user,
	)
}

// FindByID godoc
// @Summary      Get user by ID
// @Description  Get a user by their ID
// @Tags         Users
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]interface{}  "Invalid user ID"
// @Failure      404      {object}  map[string]interface{}  "User not found"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /users/{user_id} [get]
func (h *UserHandler) FindByID(c *gin.Context) {
	userID, err := uuid.Parse(
		c.Param("user_id"),
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid user ID",
		)
		return
	}

	user, err := h.service.FindByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"User not found",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to get user",
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get user successful",
		user,
	)
}

// Update godoc
// @Summary      Update user
// @Description  Update user information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user_id  path      string                 true  "User ID"
// @Param        request  body      dto.UpdateUserRequest  true  "User information"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]interface{}  "Invalid request"
// @Failure      404      {object}  map[string]interface{}  "User not found"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /users/{user_id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	userID, err := uuid.Parse(
		c.Param("user_id"),
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid user ID",
		)
		return
	}

	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	user, err := h.service.Update(
		c.Request.Context(),
		userID,
		&req,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"User not found",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			"Failed to update user",
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"User updated successfully",
		user,
	)
}
