package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"social-platform/internal/dto"
	"social-platform/internal/middleware"
	"social-platform/internal/response"
	postLikeService "social-platform/internal/service/post_like"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = dto.PostLikeResponse{}

type PostLikeHandler struct {
	service *postLikeService.Service
}

func NewPostLikeHandler(
	service *postLikeService.Service,
) *PostLikeHandler {
	return &PostLikeHandler{
		service: service,
	}
}

// Like godoc
// @Summary      Like a post
// @Description  Like a post
// @Tags         Post Likes
// @Produce      json
// @Security     CookieAuth
// @Param        post_id path int true "Post ID"
// @Success      201 {object} dto.PostLikeResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id}/like [post]
func (h *PostLikeHandler) Like(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	postID, err := strconv.ParseInt(
		c.Param("post_id"),
		10,
		64,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid post ID",
		)
		return
	}

	like, err := h.service.Like(
		c.Request.Context(),
		postID,
		userID,
	)
	if err != nil {
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
		"Like post successful",
		like,
	)
}

// Unlike godoc
// @Summary      Unlike a post
// @Description  Remove user's like from a post
// @Tags         Post Likes
// @Produce      json
// @Security     CookieAuth
// @Param        post_id path int true "Post ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Like not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id}/like [delete]
func (h *PostLikeHandler) Unlike(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	postID, err := strconv.ParseInt(
		c.Param("post_id"),
		10,
		64,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid post ID",
		)
		return
	}

	err = h.service.Unlike(
		c.Request.Context(),
		postID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Post like not found",
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

	response.NonDataJSON(
		c.Writer,
		http.StatusOK,
		"Unlike post successful",
	)
}
