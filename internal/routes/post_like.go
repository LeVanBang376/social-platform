package routes

import (
	handler "social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPostLikeRoutes(
	router *gin.Engine,
	handler *handler.PostLikeHandler,
	authMiddleware gin.HandlerFunc,
) {
	posts := router.Group("/posts")
	posts.Use(authMiddleware)
	{
		posts.POST("/:post_id/like", handler.Like)
		posts.DELETE("/:post_id/like", handler.Unlike)
	}
}
