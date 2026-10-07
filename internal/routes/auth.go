package routes

import (
	"social-platform/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	handler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
		auth.POST("/logout", handler.Logout)

		protected := auth.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("/me", handler.Me)
		}
	}
}
