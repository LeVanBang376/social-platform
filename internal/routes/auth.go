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
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST(
			"/reset-password",
			handler.ResetPassword,
		)

		protected := auth.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("/me", handler.Me)
		}
	}
}
