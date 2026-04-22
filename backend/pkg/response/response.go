package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      http.StatusOK,
		Data:      data,
		Msg:       "ok",
		RequestID: generateRequestID(),
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code:      code,
		Data:      nil,
		Msg:       msg,
		RequestID: generateRequestID(),
	})
}

func generateRequestID() string {
	return uuid.New().String()
}
