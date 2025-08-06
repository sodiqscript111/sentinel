package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sentinel/redis"
)

func BlockIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	err := redisclient.Rdb.SAdd(redisclient.Ctx, "blocked_ips", ip).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not block IP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "IP blocked successfully"})
}
