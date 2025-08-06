package main

import (
	"github.com/gin-gonic/gin"
	"sentinel/controllers"
	"sentinel/middleware"
	redisclient "sentinel/redis"
	"time"
)

func main() {
	redisclient.InitRedis()

	r := gin.Default()
	r.Use(
		middleware.BlocklistMiddleware(),
		middleware.RateLimiter(10, 100, 1*time.Minute),
	)

	r.GET("/ping", controllers.Ping)
	r.POST("/admin/block", controllers.BlockIP)

	r.Run(":8080")
}
