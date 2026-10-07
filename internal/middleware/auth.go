package middleware

import (
	"net/http"

	"social-platform/infrastructure/jwt"
	"social-platform/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextUserKey = "user"

func Auth(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")

		if err != nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Unauthorized",
			)
			c.Abort()
			return
		}

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			response.NonDataJSON(
				c.Writer,
				http.StatusUnauthorized,
				"Invalid or expired token",
			)
			c.Abort()
			return
		}

		c.Set(ContextUserKey, claims)

		c.Next()
	}
}

func GetClaims(c *gin.Context) *jwt.Claims {
	claims, exists := c.Get(ContextUserKey)

	if !exists {
		return nil
	}

	result, ok := claims.(*jwt.Claims)
	if !ok {
		return nil
	}

	return result
}

func GetUserID(c *gin.Context) uuid.UUID {
	claims := GetClaims(c)

	if claims == nil {
		return uuid.Nil
	}

	return claims.UserID
}
