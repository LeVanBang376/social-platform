package routes

import (
	handler "social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(
	router *gin.Engine,
	handler *handler.UserHandler,
	authMiddleware gin.HandlerFunc,
) {
	users := router.Group("/users")
	{
		users.POST("", handler.Create)
		users.GET("/:user_id", handler.FindByID)
	}

	protected := users.Group("")
	protected.Use(authMiddleware)
	{
		protected.PUT("/:user_id", handler.Update)
	}
}
