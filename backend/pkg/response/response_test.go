package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		data     interface{}
		wantCode int
		wantMsg  string
	}{
		{
			name:     "返回对象数据",
			data:     gin.H{"key": "value"},
			wantCode: http.StatusOK,
			wantMsg:  "ok",
		},
		{
			name:     "返回数组数据",
			data:     []string{"item1", "item2"},
			wantCode: http.StatusOK,
			wantMsg:  "ok",
		},
		{
			name:     "返回 nil",
			data:     nil,
			wantCode: http.StatusOK,
			wantMsg:  "ok",
		},
		{
			name:     "返回嵌套对象",
			data:     gin.H{"user": gin.H{"id": 1, "name": "test"}},
			wantCode: http.StatusOK,
			wantMsg:  "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Success(c, tt.data)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantMsg, resp.Msg)
			assert.NotEmpty(t, resp.RequestID)
			assert.Len(t, resp.RequestID, 36) // UUID 长度

			// 验证 Content-Type
			contentType := w.Header().Get("Content-Type")
			assert.Contains(t, contentType, "application/json")
		})
	}
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		code     int
		msg      string
		wantCode int
		wantMsg  string
	}{
		{
			name:     "400 参数错误",
			code:     http.StatusBadRequest,
			msg:      "参数错误",
			wantCode: http.StatusBadRequest,
			wantMsg:  "参数错误",
		},
		{
			name:     "401 未授权",
			code:     http.StatusUnauthorized,
			msg:      "请先登录",
			wantCode: http.StatusUnauthorized,
			wantMsg:  "请先登录",
		},
		{
			name:     "403 禁止访问",
			code:     http.StatusForbidden,
			msg:      "权限不足",
			wantCode: http.StatusForbidden,
			wantMsg:  "权限不足",
		},
		{
			name:     "404 资源不存在",
			code:     http.StatusNotFound,
			msg:      "用户不存在",
			wantCode: http.StatusNotFound,
			wantMsg:  "用户不存在",
		},
		{
			name:     "500 服务器错误",
			code:     http.StatusInternalServerError,
			msg:      "服务器内部错误",
			wantCode: http.StatusInternalServerError,
			wantMsg:  "服务器内部错误",
		},
		{
			name:     "409 资源冲突",
			code:     http.StatusConflict,
			msg:      "用户名已存在",
			wantCode: http.StatusConflict,
			wantMsg:  "用户名已存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tt.code, tt.msg)

			assert.Equal(t, tt.code, w.Code)

			var resp Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantMsg, resp.Msg)
			assert.Nil(t, resp.Data)
			assert.NotEmpty(t, resp.RequestID)
		})
	}
}

func TestResponseStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, gin.H{"test": "data"})

	// 验证 JSON 结构
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 验证必需字段存在
	assert.Contains(t, resp, "code")
	assert.Contains(t, resp, "data")
	assert.Contains(t, resp, "msg")
	assert.Contains(t, resp, "request_id")

	// 验证字段类型
	assert.IsType(t, float64(0), resp["code"])
	assert.IsType(t, "", resp["msg"])
	assert.IsType(t, "", resp["request_id"])
}

func TestGenerateRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 生成多个 RequestID 验证唯一性
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		Success(c, nil)

		var resp Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// 验证 UUID 格式 (简单检查长度和唯一性)
		assert.NotEmpty(t, resp.RequestID)
		assert.Len(t, resp.RequestID, 36)

		// 验证唯一性
		assert.False(t, ids[resp.RequestID], "RequestID 应该唯一")
		ids[resp.RequestID] = true
	}
}

func TestSuccessWithComplexData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	complexData := gin.H{
		"users": []gin.H{
			{"id": 1, "name": "user1"},
			{"id": 2, "name": "user2"},
		},
		"pagination": gin.H{
			"page":     1,
			"pageSize": 10,
			"total":    100,
		},
	}

	Success(c, complexData)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, data, "users")
	assert.Contains(t, data, "pagination")
}

func TestErrorWithEmptyMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, "")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "", resp.Msg)
	assert.NotEmpty(t, resp.RequestID)
}

func TestResponseHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success 响应头", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		Success(c, nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("Error 响应头", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		Error(c, http.StatusNotFound, "not found")

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}
