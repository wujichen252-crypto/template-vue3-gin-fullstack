package handler

import (
	"errors"
	"net/http"
	"strconv"
	"template-vue3-gin-fullstack/backend/config"
	"template-vue3-gin-fullstack/backend/internal/service"
	"template-vue3-gin-fullstack/backend/pkg/jwt"
	"template-vue3-gin-fullstack/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	svc    service.UserService
	rdb    *redis.Client
	cfg    *config.Config
	jwtMgr *jwt.JWT
}

func NewUserHandler(svc service.UserService, rdb *redis.Client, cfg *config.Config) *UserHandler {
	return &UserHandler{
		svc:    svc,
		rdb:    rdb,
		cfg:    cfg,
		jwtMgr: jwt.NewJWT(cfg.JWT.Secret, config.GetAccessTokenDuration(), config.GetRefreshTokenDuration()),
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	user, _, err := h.svc.Register(req.Username, req.Password, req.Email)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "用户名或邮箱已被注册" {
			response.Error(c, http.StatusConflict, errMsg)
			return
		}
		response.Error(c, http.StatusInternalServerError, "注册失败")
		return
	}

	token, err := h.jwtMgr.GenerateToken(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Token生成失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user_info": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"status":     user.Status,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	user, _, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "登录失败")
		return
	}

	accessToken, err := h.jwtMgr.GenerateToken(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Token生成失败")
		return
	}

	refreshToken, err := h.jwtMgr.GenerateRefreshToken(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "RefreshToken生成失败")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_info": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"status":     user.Status,
		},
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	user, err := h.svc.GetUserInfo(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"avatar_url": user.AvatarURL,
		"status":     user.Status,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if !h.jwtMgr.IsRefreshToken(req.RefreshToken) {
		response.Error(c, http.StatusUnauthorized, "无效的RefreshToken")
		return
	}

	claims, err := h.jwtMgr.ParseToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "RefreshToken解析失败")
		return
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token解析失败")
		return
	}

	if err := h.svc.RefreshToken(uint(userID)); err != nil {
		response.Error(c, http.StatusUnauthorized, "用户状态异常")
		return
	}

	accessToken, err := h.jwtMgr.GenerateToken(uint(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Token生成失败")
		return
	}

	refreshToken, err := h.jwtMgr.GenerateRefreshToken(uint(userID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "RefreshToken生成失败")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		response.Error(c, http.StatusBadRequest, "无效的Token")
		return
	}

	if err := h.svc.Logout(token.(string), config.GetAccessTokenDuration()); err != nil {
		response.Error(c, http.StatusInternalServerError, "登出失败")
		return
	}

	response.Success(c, gin.H{"message": "登出成功"})
}

func (h *UserHandler) GetUserIDFromToken(c *gin.Context) (uint, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("Authorization header is required")
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		tokenString = authHeader
	}

	claims, err := h.jwtMgr.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}