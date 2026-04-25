package middleware

import (
	"context"
	"fmt"
	"net/http"
	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	defaultLimit = 60
	authLimit    = 5
	windowSize   = 60
)

var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local windowStart = now - window

redis.call('ZREMRANGEBYSCORE', key, 0, windowStart)

local count = redis.call('ZCARD', key)

if count < limit then
    redis.call('ZADD', key, now, now .. ':' .. math.random(1000000))
    redis.call('EXPIRE', key, window)
    return 1
end

return 0
`)

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
		now := time.Now().UnixMilli()

		allowed, err := slidingWindowScript.Run(ctx, rdb, []string{key}, limit, windowSize*1000, now).Int()
		if err != nil {
			c.Next()
			return
		}

		if allowed == 0 {
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
