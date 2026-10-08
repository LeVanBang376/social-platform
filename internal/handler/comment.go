package handler

import (
	"errors"
	"net/http"
	"strconv"

	"social-platform/internal/dto"
	"social-platform/internal/middleware"
	"social-platform/internal/response"
	commentService "social-platform/internal/service/comment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	service *commentService.Service
}

func NewCommentHandler(
	service *commentService.Service,
) *CommentHandler {
	return &CommentHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create comment
// @Description  Create a new comment on a post
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        post_id path int true "Post ID"
// @Param        request body dto.CreateCommentRequest true "Create comment request"
// @Success      201 {object} dto.CommentResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "Post or parent comment not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id}/comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
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

	var req dto.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	comment, err := h.service.Create(
		c.Request.Context(),
		postID,
		userID,
		&req,
	)
	if err != nil {
		if errors.Is(err, commentService.ErrPostNotFound) ||
			errors.Is(err, commentService.ErrParentCommentNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Post or parent comment not found",
			)
			return
		}

		if errors.Is(err, commentService.ErrParentCommentWrongPost) {
			response.NonDataJSON(
				c.Writer,
				http.StatusBadRequest,
				"Parent comment does not belong to this post",
			)
			return
		}

		if errors.Is(err, commentService.ErrNestedReply) {
			response.NonDataJSON(
				c.Writer,
				http.StatusBadRequest,
				"Nested replies are not allowed",
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
		"Create comment successful",
		comment,
	)
}

// FindByPostID godoc
// @Summary      Get comments by post
// @Description  Get paginated comments belonging to a specific post
// @Tags         Comments
// @Produce      json
// @Param        post_id path int true "Post ID"
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of comments per page" default(10) maximum(200)
// @Success      200 {array} dto.CommentResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      404 {object} map[string]interface{} "Post not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id}/comments [get]
func (h *CommentHandler) FindByPostID(c *gin.Context) {
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

	pagination := response.NewPagination(c.Request)

	comments, err := h.service.FindByPostID(
		c.Request.Context(),
		postID,
		pagination,
	)
	if err != nil {
		if errors.Is(err, commentService.ErrPostNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Post not found",
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

	response.PaginatedJSON(
		c.Writer,
		http.StatusOK,
		"Get post comments successful",
		comments,
		pagination,
	)
}

// FindReplies godoc
// @Summary      Get comment replies
// @Description  Get paginated replies of a comment
// @Tags         Comments
// @Produce      json
// @Param        post_id path int true "Post ID"
// @Param        comment_id path int true "Comment ID"
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Number of replies per page" default(10) maximum(200)
// @Success      200 {array} dto.CommentResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      404 {object} map[string]interface{} "Comment not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id}/comments/{comment_id}/replies [get]
func (h *CommentHandler) FindReplies(c *gin.Context) {
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

	pagination := response.NewPagination(c.Request)

	comments, err := h.service.FindReplies(
		c.Request.Context(),
		postID,
		commentID,
		pagination,
	)
	if err != nil {
		if errors.Is(err, commentService.ErrCommentNotFound) {
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
			"Internal server error",
		)
		return
	}

	response.PaginatedJSON(
		c.Writer,
		http.StatusOK,
		"Get comment replies successful",
		comments,
		pagination,
	)
}

// Update godoc
// @Summary      Update comment
// @Description  Update a comment. Only the comment owner can update it.
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        comment_id path int true "Comment ID"
// @Param        request body dto.UpdateCommentRequest true "Update comment request"
// @Success      200 {object} dto.CommentResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Comment not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /comments/{comment_id} [put]
func (h *CommentHandler) Update(c *gin.Context) {
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

	var req dto.UpdateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	comment, err := h.service.Update(
		c.Request.Context(),
		commentID,
		userID,
		&req,
	)
	if err != nil {
		if errors.Is(err, commentService.ErrCommentNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Comment not found",
			)
			return
		}

		if errors.Is(err, commentService.ErrNotCommentOwner) {
			response.NonDataJSON(
				c.Writer,
				http.StatusForbidden,
				"You are not the owner of this comment",
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
		"Update comment successful",
		comment,
	)
}

// Delete godoc
// @Summary      Delete comment
// @Description  Delete a comment. Only the comment owner can delete it.
// @Tags         Comments
// @Security     CookieAuth
// @Param        comment_id path int true "Comment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Comment not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /comments/{comment_id} [delete]
func (h *CommentHandler) Delete(c *gin.Context) {
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

	if err := h.service.Delete(
		c.Request.Context(),
		commentID,
		userID,
	); err != nil {
		if errors.Is(err, commentService.ErrCommentNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Comment not found",
			)
			return
		}

		if errors.Is(err, commentService.ErrNotCommentOwner) {
			response.NonDataJSON(
				c.Writer,
				http.StatusForbidden,
				"You are not the owner of this comment",
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
		"Delete comment successful",
	)
}
