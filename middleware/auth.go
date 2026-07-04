package middleware

import (
	"net/http"
	"strings"

	"ticket-system/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts requests to authenticate them using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		userID, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the userID context variable
		c.Set("userID", userID)
		c.Next()
	}
}
