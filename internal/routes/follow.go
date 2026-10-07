package routes

import (
	handler "social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterFollowRoutes(
	router *gin.Engine,
	handler *handler.FollowHandler,
	authMiddleware gin.HandlerFunc,
) {
	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		users.POST("/:user_id/follow", handler.Follow)
		users.DELETE("/:user_id/follow", handler.Unfollow)
	}
}
