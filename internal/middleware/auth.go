package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renatopnasc/made2share-api/internal/config"
)

func VerifyAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("_HttpSID")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session cookie not found"})
			return
		}

		token, err := config.GetRedisDB().Get(config.Ctx, sessionID).Result()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			return
		}

		c.Set("token", token)

		c.Next()
	}
}
