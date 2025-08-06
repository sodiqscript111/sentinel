package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	redisclient "sentinel/redis"
)

func BlocklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		isBlocked, err := redisclient.Rdb.SIsMember(redisclient.Ctx, "blocked_ips", ip).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if isBlocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your IP has been blocked"})
			c.Abort()
			return
		}

		c.Next()
	}
}
