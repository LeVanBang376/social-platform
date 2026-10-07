package handler

import (
	"errors"
	"fmt"
	"net/http"

	"social-platform/internal/middleware"
	"social-platform/internal/response"
	followService "social-platform/internal/service/follow"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FollowHandler struct {
	service *followService.Service
}

func NewFollowHandler(
	service *followService.Service,
) *FollowHandler {
	return &FollowHandler{
		service: service,
	}
}

// Follow godoc
// @Summary      Follow a user
// @Description  Follow another user
// @Tags         Follow
// @Produce      json
// @Security     CookieAuth
// @Param        user_id  path  string  true  "User ID"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}  "Invalid user ID"
// @Failure      404  {object}  map[string]interface{}  "User not found"
// @Failure      409  {object}  map[string]interface{}  "Already following user"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /users/{user_id}/follow [post]
func (h *FollowHandler) Follow(c *gin.Context) {
	followerID := middleware.GetUserID(c)

	if followerID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	followingID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid user ID",
		)
		return
	}

	follow, err := h.service.Follow(
		c.Request.Context(),
		followerID,
		followingID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusBadRequest,
				"Cannot follow yourself",
			)
			return
		}

		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			fmt.Sprintf("Internal server error: %s", err.Error()),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusCreated,
		"Follow user successfully",
		follow,
	)
}

// Unfollow godoc
// @Summary      Unfollow a user
// @Description  Unfollow another user
// @Tags         Follow
// @Produce      json
// @Security     CookieAuth
// @Param        user_id  path  string  true  "User ID"
// @Success      204  "Successfully unfollowed"
// @Failure      400  {object}  map[string]interface{}  "Invalid user ID"
// @Failure      404  {object}  map[string]interface{}  "Follow relationship not found"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /users/{user_id}/follow [delete]
func (h *FollowHandler) Unfollow(c *gin.Context) {
	followerID := middleware.GetUserID(c)

	if followerID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	followingID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid user ID",
		)
		return
	}

	if err := h.service.Unfollow(
		c.Request.Context(),
		followerID,
		followingID,
	); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusInternalServerError,
			fmt.Sprintf("Internal server error: %s", err.Error()),
		)
		return
	}

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Unfollow user successfully",
	)
}
