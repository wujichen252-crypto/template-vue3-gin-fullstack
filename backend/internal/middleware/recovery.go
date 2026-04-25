package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"template-vue3-gin-fullstack/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger, stackStack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))

				c.Abort()
			}
		}()

		c.Next()
	}
}
