package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"social-platform/internal/dto"
	"social-platform/internal/middleware"
	"social-platform/internal/response"
	commentLikeService "social-platform/internal/service/comment_like"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ = dto.CommentLikeResponse{}

type CommentLikeHandler struct {
	service *commentLikeService.Service
}

func NewCommentLikeHandler(
	service *commentLikeService.Service,
) *CommentLikeHandler {
	return &CommentLikeHandler{
		service: service,
	}
}

// Like godoc
// @Summary      Like comment
// @Description  Like a comment
// @Tags         Comment Likes
// @Security     CookieAuth
// @Produce      json
// @Param        comment_id path int true "Comment ID"
// @Success      201 {object} dto.CommentLikeResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /comments/{comment_id}/likes [post]
func (h *CommentLikeHandler) Like(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	commentID, err := strconv.ParseInt(
		c.Param("comment_id"),
		10,
		64,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid comment ID",
		)
		return
	}

	like, err := h.service.Like(
		c.Request.Context(),
		commentID,
		userID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Comment not found",
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
		"Like comment successful",
		like,
	)
}

// Unlike godoc
// @Summary      Unlike comment
// @Description  Remove the current user's like from a comment
// @Tags         Comment Likes
// @Security     CookieAuth
// @Param        comment_id path int true "Comment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Like not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /comments/{comment_id}/likes [delete]
func (h *CommentLikeHandler) Unlike(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	commentID, err := strconv.ParseInt(
		c.Param("comment_id"),
		10,
		64,
	)
	if err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid comment ID",
		)
		return
	}

	if err := h.service.Unlike(
		c.Request.Context(),
		commentID,
		userID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Comment like not found",
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
		"Unlike comment successful",
	)
}
