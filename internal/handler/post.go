package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"social-platform/internal/dto"
	"social-platform/internal/middleware"
	"social-platform/internal/response"
	postService "social-platform/internal/service/post"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostHandler struct {
	service *postService.Service
}

func NewPostHandler(
	service *postService.Service,
) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create post
// @Description  Create a new post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        request body dto.CreatePostRequest true "Create post request"
// @Success      201 {object} dto.PostResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == uuid.Nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	var req dto.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	post, err := h.service.Create(
		c.Request.Context(),
		&req,
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
		"Create post successful",
		post,
	)
}

// FindByID godoc
// @Summary      Get post by ID
// @Description  Get a post by its ID
// @Tags         Posts
// @Produce      json
// @Param        post_id path int true "Post ID"
// @Success      200 {object} dto.PostResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      404 {object} map[string]interface{} "Post not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id} [get]
func (h *PostHandler) FindByID(c *gin.Context) {
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

	post, err := h.service.FindByID(
		c.Request.Context(),
		postID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
			fmt.Sprintf("Internal server error: %s", err.Error()),
		)
		return
	}

	response.JSON(
		c.Writer,
		http.StatusOK,
		"Get post successful",
		post,
	)
}

// FindByUserID godoc
// @Summary      Get posts by user
// @Description  Get posts created by a specific user
// @Tags         Posts
// @Produce      json
// @Param        user_id path string true "User ID"
// @Param        limit query int false "Number of posts"
// @Param        offset query int false "Offset"
// @Success      200 {array} dto.PostResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{user_id}/posts [get]
func (h *PostHandler) FindByUserID(c *gin.Context) {
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

	var query dto.FindPostsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid query parameters",
		)
		return
	}

	posts, err := h.service.FindByUserID(
		c.Request.Context(),
		userID,
		&query,
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
		http.StatusOK,
		"Get user posts successful",
		posts,
	)
}

// Update godoc
// @Summary      Update post
// @Description  Update a post. Only the post owner can update it.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        post_id path int true "Post ID"
// @Param        request body dto.UpdatePostRequest true "Update post request"
// @Success      200 {object} dto.PostResponse
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Post not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id} [put]
func (h *PostHandler) Update(c *gin.Context) {
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

	var req dto.UpdatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.NonDataJSON(
			c.Writer,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	post, err := h.service.Update(
		c.Request.Context(),
		postID,
		userID,
		&req,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Post not found",
			)
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusForbidden,
				"You are not the owner of this post",
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
		http.StatusOK,
		"Update post successful",
		post,
	)
}

// Delete godoc
// @Summary      Delete post
// @Description  Delete a post. Only the post owner can delete it.
// @Tags         Posts
// @Security     CookieAuth
// @Param        post_id path int true "Post ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "Bad request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      403 {object} map[string]interface{} "Forbidden"
// @Failure      404 {object} map[string]interface{} "Post not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /posts/{post_id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
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

	if err := h.service.Delete(
		c.Request.Context(),
		postID,
		userID,
	); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NonDataJSON(
				c.Writer,
				http.StatusNotFound,
				"Post not found",
			)
			return
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			response.NonDataJSON(
				c.Writer,
				http.StatusForbidden,
				"You are not the owner of this post",
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
		"Delete post successful",
	)
}
