package middleware

import (
	"ginorm/constant"
	"ginorm/util"
	"github.com/gin-gonic/gin"
)

// Logger 记录access日志，增加trace id
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(constant.HeaderTraceId)
		if traceID == "" {
			traceID = util.GenerateTraceID()
		}
		// trace id塞入上下文
		c.Set(constant.TraceId, traceID)
		// 记录access日志

		c.Next()
	}
}
