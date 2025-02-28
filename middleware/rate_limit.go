package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// 创建一个IP到限速器的映射
var (
	limiterMap = make(map[string]*rate.Limiter)
	mu         sync.Mutex
)

// 获取限速器
func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiterMap[ip]
	if !exists {
		// 创建新的限速器，r为每秒允许的请求数，b为桶容量
		r := rate.Limit(viper.GetFloat64("ratelimit.requests_per_second"))
		b := viper.GetInt("ratelimit.burst")
		limiter = rate.NewLimiter(r, b)
		limiterMap[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware 限速中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
