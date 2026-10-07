package routes

import (
	handler "social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPostRoutes(
	router *gin.Engine,
	handler *handler.PostHandler,
	authMiddleware gin.HandlerFunc,
) {
	// Public post routes
	posts := router.Group("/posts")
	{
		posts.GET("/:post_id", handler.FindByID)
	}

	// Public user post routes
	users := router.Group("/users")
	{
		users.GET("/:user_id/posts", handler.FindByUserID)
	}

	// Protected post routes
	protectedPosts := posts.Group("")
	protectedPosts.Use(authMiddleware)
	{
		protectedPosts.POST("", handler.Create)
		protectedPosts.PUT("/:post_id", handler.Update)
		protectedPosts.DELETE("/:post_id", handler.Delete)
	}
}
