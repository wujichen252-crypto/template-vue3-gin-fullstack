package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, gin.H{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("Status code: got %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != http.StatusOK {
		t.Errorf("Response code: got %d, want %d", resp.Code, http.StatusOK)
	}

	if resp.Msg != "ok" {
		t.Errorf("Response msg: got %s, want ok", resp.Msg)
	}

	if resp.RequestID == "" {
		t.Error("RequestID should not be empty")
	}
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, http.StatusBadRequest, "参数错误")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code: got %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Response code: got %d, want %d", resp.Code, http.StatusBadRequest)
	}

	if resp.Msg != "参数错误" {
		t.Errorf("Response msg: got %s, want 参数错误", resp.Msg)
	}

	if resp.Data != nil {
		t.Error("Data should be nil for error response")
	}

	if resp.RequestID == "" {
		t.Error("RequestID should not be empty")
	}
}