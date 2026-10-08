package handler

import (
	"errors"
	"net/http"

	"social-platform/internal/middleware"
	"social-platform/internal/response"
	followService "social-platform/internal/service/follow"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Failure      400  {object}  map[string]interface{}  "Invalid user ID or cannot follow yourself"
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
		if errors.Is(err, followService.ErrCannotFollowSelf) {
			response.NonDataJSON(
				c.Writer,
				http.StatusBadRequest,
				"Cannot follow yourself",
			)
			return
		}

		if errors.Is(err, followService.ErrAlreadyFollowing) {
			response.NonDataJSON(
				c.Writer,
				http.StatusConflict,
				"Already following user",
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
// @Success      200  {object} map[string]interface{} "Successfully unfollowed"
// @Failure      400  {object} map[string]interface{}  "Invalid user ID"
// @Failure      404  {object} map[string]interface{}  "Follow relationship not found"
// @Failure      500  {object} map[string]interface{}  "Internal server error"
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
		if errors.Is(err, followService.ErrNotFollowing) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Follow relationship not found",
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
		"Unfollow user successfully",
	)
}
