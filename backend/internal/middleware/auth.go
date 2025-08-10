package middleware

import (
	"github.com/gin-gonic/gin"
	"igropoisk_backend/internal/auth"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": " Authorization header is must"})
			c.Abort()
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}
		token := parts[1]
		claims, err := auth.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)

		c.Next()
	}

}
