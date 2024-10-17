package middleware

import (
	"fmt"
	"ginorm/entity/response"
	"ginorm/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				logger.Writer.Error(fmt.Sprintf("Panic: %v\n", err))
				// 生产环境隐藏底层报错
				detail := ""
				if err != nil && gin.Mode() != gin.ReleaseMode {
					detail = fmt.Sprintf("%v", err)
				}
				// 返回统一的 JSON 错误响应
				c.JSON(http.StatusInternalServerError, &response.Response{
					Code: http.StatusInternalServerError,
					Data: detail,
					Msg:  "Internal Server Error",
				})

				// 防止继续处理请求
				c.Abort()
			}
		}()
		// 继续处理请求
		c.Next()
	}
}
