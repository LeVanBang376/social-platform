package routes

import (
	handler "social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(
	router *gin.Engine,
	handler *handler.CommentHandler,
	authMiddleware gin.HandlerFunc,
) {
	// Public post comment routes
	posts := router.Group("/posts")
	{
		posts.GET("/:post_id/comments", handler.FindByPostID)
		posts.GET(
			"/:post_id/comments/:comment_id/replies",
			handler.FindReplies,
		)
	}

	// Protected post comment routes
	protectedPosts := posts.Group("")
	protectedPosts.Use(authMiddleware)
	{
		protectedPosts.POST("/:post_id/comments", handler.Create)
	}

	// Protected comment routes
	comments := router.Group("/comments")
	protectedComments := comments.Group("")
	protectedComments.Use(authMiddleware)
	{
		protectedComments.PUT("/:comment_id", handler.Update)
		protectedComments.DELETE("/:comment_id", handler.Delete)
	}
}
