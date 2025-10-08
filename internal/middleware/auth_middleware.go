package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SIMPLE MIDDLEWARE TO CHECK BEARER TOKEN HEADER
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow preflight requests to pass through
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// Check if format Bearer [token]
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			return
		}

		token := parts[1]

		// Dummy validation, just check token length
		if len(strings.TrimSpace(token)) < 10 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// CONTINUE
		c.Next()
	}
}
