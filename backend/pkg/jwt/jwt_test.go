package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewJWT(t *testing.T) {
	secret := "test-secret"
	accessDuration := time.Hour
	refreshDuration := time.Hour * 24

	j := NewJWT(secret, accessDuration, refreshDuration)

	assert.NotNil(t, j)
	assert.Equal(t, []byte(secret), j.secret)
	assert.Equal(t, accessDuration, j.accessDuration)
	assert.Equal(t, refreshDuration, j.refreshDuration)
}

func TestJWT_GenerateToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "正常生成Token",
			userID:  123,
			wantErr: false,
		},
		{
			name:    "用户ID为0",
			userID:  0,
			wantErr: false,
		},
		{
			name:    "大用户ID",
			userID:  999999999,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtMgr.GenerateToken(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				// 验证 token 格式 (xxx.yyy.zzz)
				parts := strings.Split(token, ".")
				assert.Equal(t, 3, len(parts))
			}
		})
	}
}

func TestJWT_ParseToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	t.Run("解析有效Token", func(t *testing.T) {
		userID := uint(123)
		token, err := jwtMgr.GenerateToken(userID)
		assert.NoError(t, err)

		claims, err := jwtMgr.ParseToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, "123", claims.Subject)
	})

	t.Run("解析无效Token", func(t *testing.T) {
		claims, err := jwtMgr.ParseToken("invalid-token")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("解析空Token", func(t *testing.T) {
		claims, err := jwtMgr.ParseToken("")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("解析过期Token", func(t *testing.T) {
		// 创建一个已过期 token
		expiredMgr := NewJWT(secret, -time.Hour, time.Hour*24)
		token, err := expiredMgr.GenerateToken(123)
		assert.NoError(t, err)

		claims, err := jwtMgr.ParseToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("使用错误密钥解析", func(t *testing.T) {
		jwtMgr1 := NewJWT("secret-1", time.Hour, time.Hour*24)
		jwtMgr2 := NewJWT("secret-2", time.Hour, time.Hour*24)

		token, err := jwtMgr1.GenerateToken(123)
		assert.NoError(t, err)

		claims, err := jwtMgr2.ParseToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("篡改Token", func(t *testing.T) {
		token, err := jwtMgr.GenerateToken(123)
		assert.NoError(t, err)

		// 篡改 token 内容
		tamperedToken := token[:len(token)-5] + "xxxxx"

		claims, err := jwtMgr.ParseToken(tamperedToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestJWT_GenerateRefreshToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "正常生成RefreshToken",
			userID:  456,
			wantErr: false,
		},
		{
			name:    "用户ID为0",
			userID:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwtMgr.GenerateRefreshToken(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// 验证可以解析
				claims, err := jwtMgr.ParseToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
			}
		})
	}
}

func TestJWT_IsRefreshToken(t *testing.T) {
	secret := "test-secret-key-for-testing"
	accessDuration := time.Hour
	refreshDuration := time.Hour * 24 * 7 // 7天
	jwtMgr := NewJWT(secret, accessDuration, refreshDuration)

	t.Run("有效的RefreshToken", func(t *testing.T) {
		token, err := jwtMgr.GenerateRefreshToken(123)
		assert.NoError(t, err)

		isRefresh := jwtMgr.IsRefreshToken(token)
		assert.True(t, isRefresh)
	})

	t.Run("AccessToken不是RefreshToken", func(t *testing.T) {
		token, err := jwtMgr.GenerateToken(123)
		assert.NoError(t, err)

		isRefresh := jwtMgr.IsRefreshToken(token)
		assert.False(t, isRefresh)
	})

	t.Run("无效Token", func(t *testing.T) {
		isRefresh := jwtMgr.IsRefreshToken("invalid-token")
		assert.False(t, isRefresh)
	})

	t.Run("空Token", func(t *testing.T) {
		isRefresh := jwtMgr.IsRefreshToken("")
		assert.False(t, isRefresh)
	})

	t.Run("过期Token", func(t *testing.T) {
		expiredMgr := NewJWT(secret, time.Hour, -time.Hour)
		token, err := expiredMgr.GenerateRefreshToken(123)
		assert.NoError(t, err)

		isRefresh := jwtMgr.IsRefreshToken(token)
		assert.False(t, isRefresh)
	})
}

func TestClaims_Valid(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	t.Run("Token包含正确的声明", func(t *testing.T) {
		userID := uint(789)
		token, err := jwtMgr.GenerateToken(userID)
		assert.NoError(t, err)

		claims, err := jwtMgr.ParseToken(token)
		assert.NoError(t, err)

		// 验证时间声明
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.NotBefore)

		// 验证过期时间在未来
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))

		// 验证签发时间在现在或过去
		assert.True(t, claims.IssuedAt.Time.Before(time.Now().Add(time.Second)))
	})
}

func TestJWT_DifferentSigningMethods(t *testing.T) {
	secret := "test-secret-key-for-testing"
	jwtMgr := NewJWT(secret, time.Hour, time.Hour*24)

	t.Run("使用HS256签名", func(t *testing.T) {
		token, err := jwtMgr.GenerateToken(123)
		assert.NoError(t, err)

		// 解析但不验证，检查签名方法
		parser := jwt.NewParser()
		tokenObj, _, err := parser.ParseUnverified(token, &Claims{})
		assert.NoError(t, err)

		assert.Equal(t, jwt.SigningMethodHS256, tokenObj.Method)
	})
}

func TestJWT_EdgeCases(t *testing.T) {
	t.Run("超长密钥", func(t *testing.T) {
		longSecret := strings.Repeat("a", 1000)
		jwtMgr := NewJWT(longSecret, time.Hour, time.Hour*24)

		token, err := jwtMgr.GenerateToken(123)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := jwtMgr.ParseToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(123), claims.UserID)
	})

	t.Run("极短有效期", func(t *testing.T) {
		jwtMgr := NewJWT("secret", time.Nanosecond, time.Hour)
		token, err := jwtMgr.GenerateToken(123)
		assert.NoError(t, err)

		// 稍微等待确保过期
		time.Sleep(time.Millisecond * 10)

		claims, err := jwtMgr.ParseToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("零有效期", func(t *testing.T) {
		jwtMgr := NewJWT("secret", 0, time.Hour)
		token, err := jwtMgr.GenerateToken(123)
		assert.NoError(t, err)

		// 立即过期
		claims, err := jwtMgr.ParseToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
