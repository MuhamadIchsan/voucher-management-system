package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SIMPLE MIDDLEWARE TO CHECK BEARER TOKEN HEADER
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Check if format Bearer "token"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Dummy validation, just check token length
		if len(token) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// CONTINUE
		c.Next()
	}
}
