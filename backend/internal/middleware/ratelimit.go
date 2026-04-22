package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	defaultLimit = 60
	authLimit    = 5
)

func RateLimiter(rdb *redis.Client, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path

		var limit int
		var key string

		if isAuthEndpoint(path) {
			limit = authLimit
			key = fmt.Sprintf("rate_limit:auth:%s", ip)
		} else {
			limit = defaultLimit
			key = fmt.Sprintf("rate_limit:%s", ip)
		}

		ctx := context.Background()

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if int(count) > limit {
			if isAuthEndpoint(path) {
				response.Error(c, http.StatusTooManyRequests, "登录请求过于频繁，请稍后再试")
			} else {
				response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

func isAuthEndpoint(path string) bool {
	authPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}
	for _, authPath := range authPaths {
		if path == authPath {
			return true
		}
	}
	return false
}