package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/pkg/jwt"
	"template-vue3-gin-fullstack/backend/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// tokenBlacklist 本地缓存，减少 Redis 查询频率
type tokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

func newTokenBlacklist() *tokenBlacklist {
	tb := &tokenBlacklist{
		tokens: make(map[string]time.Time),
	}
	// 定期清理过期条目
	go tb.cleanup()
	return tb
}

func (tb *tokenBlacklist) Add(token string, exp time.Duration) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens[token] = time.Now().Add(exp)
}

func (tb *tokenBlacklist) IsBlacklisted(token string) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	exp, exists := tb.tokens[token]
	if !exists {
		return false
	}
	return time.Now().Before(exp)
}

func (tb *tokenBlacklist) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		tb.mu.Lock()
		now := time.Now()
		for token, exp := range tb.tokens {
			if now.After(exp) {
				delete(tb.tokens, token)
			}
		}
		tb.mu.Unlock()
	}
}

var globalBlacklist = newTokenBlacklist()

func JWT(secret string, rdb *redis.Client) gin.HandlerFunc {
	jwtMgr := jwt.NewJWT(secret, config.GetAccessTokenDuration(), config.GetRefreshTokenDuration())

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未授权：缺少Token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "未授权：Token格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 先检查本地缓存，减少 Redis 查询
		if globalBlacklist.IsBlacklisted(tokenString) {
			response.Error(c, http.StatusUnauthorized, "未授权：Token已失效")
			c.Abort()
			return
		}

		// 本地缓存未命中时，检查 Redis
		if rdb != nil {
			ctx := context.Background()
			blacklistKey := fmt.Sprintf("jwt_blacklist:%s", tokenString)
			exists, _ := rdb.Exists(ctx, blacklistKey).Result()
			if exists > 0 {
				response.Error(c, http.StatusUnauthorized, "未授权：Token已失效")
				c.Abort()
				return
			}
		}

		claims, err := jwtMgr.ParseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "未授权：Token无效或已过期")
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "未授权：Token解析失败")
			c.Abort()
			return
		}

		c.Set("user_id", uint(userID))
		c.Set("token", tokenString)
		c.Next()
	}
}

func BlacklistToken(rdb *redis.Client, token string, exp time.Duration) error {
	// 同时更新本地缓存和 Redis
	globalBlacklist.Add(token, exp)

	if rdb != nil {
		ctx := context.Background()
		key := fmt.Sprintf("jwt_blacklist:%s", token)
		return rdb.Set(ctx, key, "1", exp).Err()
	}
	return nil
}