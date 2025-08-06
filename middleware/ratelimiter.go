package middleware

import (
	"fmt"
	"net/http"
	"sentinel/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(limit int, abuseThreshold int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rateKey := fmt.Sprintf("ratelimit:%s", ip)
		abuseKey := fmt.Sprintf("abusecount:%s", ip)

		abuseCount, _ := redisclient.Rdb.Incr(redisclient.Ctx, abuseKey).Result()
		redisclient.Rdb.Expire(redisclient.Ctx, abuseKey, window)

		if abuseCount >= int64(abuseThreshold) {
			_ = redisclient.Rdb.SAdd(redisclient.Ctx, "blocked_ips", ip).Err()
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Your IP has been blocked for abuse",
			})
			c.Abort()
			return
		}

		count, _ := redisclient.Rdb.Get(redisclient.Ctx, rateKey).Int()
		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			c.Abort()
			return
		}

		pipe := redisclient.Rdb.TxPipeline()
		pipe.Incr(redisclient.Ctx, rateKey)
		pipe.Expire(redisclient.Ctx, rateKey, window)
		_, err := pipe.Exec(redisclient.Ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Redis  error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
