package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/pkg/jwt"
	"template-vue3-gin-fullstack/backend/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

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

		ctx := context.Background()
		blacklistKey := fmt.Sprintf("jwt_blacklist:%s", tokenString)
		exists, _ := rdb.Exists(ctx, blacklistKey).Result()
		if exists > 0 {
			response.Error(c, http.StatusUnauthorized, "未授权：Token已失效")
			c.Abort()
			return
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
	ctx := context.Background()
	key := fmt.Sprintf("jwt_blacklist:%s", token)
	return rdb.Set(ctx, key, "1", exp).Err()
}