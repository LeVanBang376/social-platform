package routes

import (
	handler "social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCommentLikeRoutes(
	router *gin.Engine,
	handler *handler.CommentLikeHandler,
	authMiddleware gin.HandlerFunc,
) {
	comments := router.Group("/comments")
	comments.Use(authMiddleware)
	{
		comments.POST("/:comment_id/likes", handler.Like)
		comments.DELETE("/:comment_id/likes", handler.Unlike)
	}
}
