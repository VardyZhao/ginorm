package middleware

import (
	"bytes"
	"ginorm/constant"
	"ginorm/logger"
	"ginorm/util"
	"github.com/gin-gonic/gin"
	"io"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取或生成 traceID
		traceID := c.GetHeader(constant.HeaderTraceId)
		if traceID == "" {
			traceID = util.GenerateTraceID()
		}
		// traceID 写入上下文
		c.Set(constant.TraceId, traceID)

		// 读取请求体
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // 还原请求体

		// 创建日志上下文
		ctx := map[string]interface{}{
			"request_url":    c.Request.URL.String(),
			"request_method": c.Request.Method,
			"request_header": c.Request.Header,
			"request_body":   string(bodyBytes), // 将请求体保存为字符串
			"client_ip":      c.ClientIP(),
		}

		logger.Log(logger.LevelInfo, "access", ctx, c)

		// 执行后续处理器
		c.Next()
	}
}
