package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/internal/model"
	"template-vue3-gin-fullstack/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(username, password, email string) (*model.User, error) {
	args := m.Called(username, password, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(username, password string) (*model.User, error) {
	args := m.Called(username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUserInfo(ctx context.Context, id uint) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) ClearUserCache(ctx context.Context, userID uint) {
	m.Called(ctx, userID)
}

func (m *MockUserService) RefreshToken(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserService) Logout(ctx context.Context, token string, exp time.Duration) error {
	args := m.Called(ctx, token, exp)
	return args.Error(0)
}

// setupTest 创建测试所需的 Handler 和 Mock
func setupTest() (*UserHandler, *MockUserService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockUserService)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key-for-unit-testing-only",
			AccessExpire:  2,
			RefreshExpire: 168,
		},
	}
	handler := NewUserHandler(mockSvc, nil, cfg)
	router := gin.New()
	return handler, mockSvc, router
}

func TestUserHandler_Register(t *testing.T) {
	handler, mockSvc, router := setupTest()
	router.POST("/register", handler.Register)

	tests := []struct {
		name         string
		reqBody      map[string]interface{}
		mockSetup    func()
		wantStatus   int
		wantErrMsg   string
		wantToken    bool
		wantUserInfo bool
	}{
		{
			name: "注册成功",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
				"email":    "test@example.com",
			},
			mockSetup: func() {
			mockSvc.On("Register", "testuser", "password123", "test@example.com").
				Return(&model.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					AvatarURL: "",
					Status:    1,
				}, nil).Once()
		},
			wantStatus:   http.StatusOK,
			wantToken:    true,
			wantUserInfo: true,
		},
		{
			name: "参数错误-缺少用户名",
			reqBody: map[string]interface{}{
				"password": "password123",
				"email":    "test@example.com",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "参数错误-用户名太短",
			reqBody: map[string]interface{}{
				"username": "ab",
				"password": "password123",
				"email":    "test@example.com",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "参数错误-密码太短",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "123",
				"email":    "test@example.com",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "参数错误-邮箱格式错误",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
				"email":    "invalid-email",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "用户名已存在",
			reqBody: map[string]interface{}{
				"username": "existinguser",
				"password": "password123",
				"email":    "existing@example.com",
			},
			mockSetup: func() {
			mockSvc.On("Register", "existinguser", "password123", "existing@example.com").
				Return(nil, errors.New("用户名或邮箱已被注册")).Once()
		},
			wantStatus: http.StatusConflict,
			wantErrMsg: "用户名或邮箱已被注册",
		},
		{
			name: "注册失败-其他错误",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
				"email":    "test@example.com",
			},
			mockSetup: func() {
			mockSvc.On("Register", "testuser", "password123", "test@example.com").
				Return(nil, errors.New("数据库错误")).Once()
		},
			wantStatus: http.StatusInternalServerError,
			wantErrMsg: "注册失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.wantErrMsg != "" {
				assert.Equal(t, tt.wantErrMsg, resp["msg"])
			}

			if tt.wantToken {
				data, ok := resp["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["token"])
			}

			if tt.wantUserInfo {
				data, ok := resp["data"].(map[string]interface{})
				assert.True(t, ok)
				userInfo, ok := data["user_info"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotNil(t, userInfo["id"])
				assert.NotNil(t, userInfo["username"])
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	handler, mockSvc, router := setupTest()
	router.POST("/login", handler.Login)

	tests := []struct {
		name           string
		reqBody        map[string]interface{}
		mockSetup      func()
		wantStatus     int
		wantErrMsg     string
		wantAccessToken bool
		wantRefreshToken bool
	}{
		{
			name: "登录成功",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			mockSetup: func() {
			mockSvc.On("Login", "testuser", "password123").
				Return(&model.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					AvatarURL: "",
					Status:    1,
				}, nil).Once()
		},
			wantStatus:       http.StatusOK,
			wantAccessToken:  true,
			wantRefreshToken: true,
		},
		{
			name: "参数错误-缺少用户名",
			reqBody: map[string]interface{}{
				"password": "password123",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "参数错误-缺少密码",
			reqBody: map[string]interface{}{
				"username": "testuser",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "登录失败-用户不存在或密码错误",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpassword",
			},
			mockSetup: func() {
			mockSvc.On("Login", "testuser", "wrongpassword").
				Return(nil, errors.New("用户不存在")).Once()
		},
			wantStatus: http.StatusUnauthorized,
			wantErrMsg: "登录失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.wantErrMsg != "" {
				assert.Equal(t, tt.wantErrMsg, resp["msg"])
			}

			if tt.wantAccessToken || tt.wantRefreshToken {
				data, ok := resp["data"].(map[string]interface{})
				assert.True(t, ok)
				if tt.wantAccessToken {
					assert.NotEmpty(t, data["access_token"])
				}
				if tt.wantRefreshToken {
					assert.NotEmpty(t, data["refresh_token"])
				}
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUserInfo(t *testing.T) {
	handler, mockSvc, router := setupTest()
	router.GET("/userinfo", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		handler.GetUserInfo(c)
	})

	tests := []struct {
		name         string
		mockSetup    func()
		wantStatus   int
		wantErrMsg   string
		wantUserInfo bool
	}{
		{
			name: "获取成功",
			mockSetup: func() {
				mockSvc.On("GetUserInfo", mock.Anything, uint(1)).Return(&model.User{
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					AvatarURL: "",
					Status:    1,
				}, nil).Once()
			},
			wantStatus:   http.StatusOK,
			wantUserInfo: true,
		},
		{
			name: "用户不存在",
			mockSetup: func() {
				mockSvc.On("GetUserInfo", mock.Anything, uint(1)).
					Return(nil, errors.New("用户不存在")).Once()
			},
			wantStatus: http.StatusNotFound,
			wantErrMsg: "用户不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.wantErrMsg != "" {
				assert.Equal(t, tt.wantErrMsg, resp["msg"])
			}

			if tt.wantUserInfo {
				data, ok := resp["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotNil(t, data["id"])
				assert.NotNil(t, data["username"])
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUserInfo_Unauthorized(t *testing.T) {
	handler, _, router := setupTest()
	router.GET("/userinfo", func(c *gin.Context) {
		// 不设置 user_id，模拟未授权
		handler.GetUserInfo(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "未授权", resp["msg"])
}

func TestUserHandler_RefreshToken(t *testing.T) {
	handler, mockSvc, router := setupTest()
	router.POST("/refresh", handler.RefreshToken)

	// 先生成一个有效的 refresh token
	jwtMgr := jwt.NewJWT("test-secret-key-for-unit-testing-only",
		time.Duration(2)*time.Hour,
		time.Duration(168)*time.Hour)
	validRefreshToken, _ := jwtMgr.GenerateRefreshToken(1)

	tests := []struct {
		name             string
		reqBody          map[string]interface{}
		mockSetup        func()
		wantStatus       int
		wantErrMsg       string
		wantAccessToken  bool
		wantRefreshToken bool
	}{
		{
			name: "刷新成功",
			reqBody: map[string]interface{}{
				"refresh_token": validRefreshToken,
			},
			mockSetup: func() {
				mockSvc.On("RefreshToken", uint(1)).Return(nil).Once()
			},
			wantStatus:       http.StatusOK,
			wantAccessToken:  true,
			wantRefreshToken: true,
		},
		{
			name: "参数错误-缺少refresh_token",
			reqBody: map[string]interface{}{
				"refresh_token": "",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
			wantErrMsg: "参数错误",
		},
		{
			name: "无效的RefreshToken",
			reqBody: map[string]interface{}{
				"refresh_token": "invalid.token.here",
			},
			mockSetup:  func() {},
			wantStatus: http.StatusUnauthorized,
			wantErrMsg: "无效的RefreshToken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.wantErrMsg != "" {
				assert.Equal(t, tt.wantErrMsg, resp["msg"])
			}

			if tt.wantAccessToken || tt.wantRefreshToken {
				data, ok := resp["data"].(map[string]interface{})
				assert.True(t, ok)
				if tt.wantAccessToken {
					assert.NotEmpty(t, data["access_token"])
				}
				if tt.wantRefreshToken {
					assert.NotEmpty(t, data["refresh_token"])
				}
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Logout(t *testing.T) {
	handler, mockSvc, router := setupTest()
	router.POST("/logout", func(c *gin.Context) {
		c.Set("token", "test-token-123")
		handler.Logout(c)
	})

	tests := []struct {
		name       string
		mockSetup  func()
		wantStatus int
		wantErrMsg string
	}{
		{
			name: "登出成功",
			mockSetup: func() {
			mockSvc.On("Logout", mock.AnythingOfType("context.backgroundCtx"), "test-token-123", mock.AnythingOfType("time.Duration")).
				Return(nil).Once()
		},
			wantStatus: http.StatusOK,
		},
		{
			name: "登出失败",
			mockSetup: func() {
			mockSvc.On("Logout", mock.AnythingOfType("context.backgroundCtx"), "test-token-123", mock.AnythingOfType("time.Duration")).
				Return(errors.New("Redis错误")).Once()
		},
			wantStatus: http.StatusInternalServerError,
			wantErrMsg: "登出失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.wantErrMsg != "" {
				assert.Equal(t, tt.wantErrMsg, resp["msg"])
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Logout_NoToken(t *testing.T) {
	handler, _, router := setupTest()
	router.POST("/logout", func(c *gin.Context) {
		// 不设置 token
		handler.Logout(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "无效的Token", resp["msg"])
}
