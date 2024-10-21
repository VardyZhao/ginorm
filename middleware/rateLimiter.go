package middleware

import (
	"ginorm/entity/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// 为每个 IP 地址创建限流器
var visitors = make(map[string]*rate.Limiter)
var apiLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

// RateLimit 限流中间件：基于 IP 地址限流
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitorLimiter(ip)
		apiLimiter := getApiLimiter(c.Request.URL.Path)

		if !limiter.Allow() || !apiLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, &response.Response{
				Code: http.StatusTooManyRequests,
				Msg:  "Too many requests from this IP, try again later.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 获取或创建一个限流器
func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 3) // 每秒1个请求，最多积累3个令牌
		visitors[ip] = limiter

		// 设置过期时间，清理旧的限流器
		go func() {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			delete(visitors, ip)
			mu.Unlock()
		}()
	}
	return limiter
}

func getApiLimiter(path string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exists := apiLimiters[path]
	if !exists {
		limiter = rate.NewLimiter(1, 3)
		apiLimiters[path] = limiter
		go func() {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			delete(apiLimiters, path)
			mu.Unlock()
		}()
	}
	return limiter
}
