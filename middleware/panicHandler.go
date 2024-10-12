package middleware

import (
	"fmt"
	"ginorm/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func PanicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				util.Log().Error("Panic: %v\n", err)
				debug.PrintStack()
				// todo 生产环境隐藏底层报错
				detail := ""
				if err != nil && gin.Mode() != gin.ReleaseMode {
					detail = fmt.Sprintf("%v", err)
				}
				// 返回统一的 JSON 错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":   http.StatusInternalServerError,
					"msg":    "Internal Server Error",
					"detail": detail,
				})

				// 防止继续处理请求
				c.Abort()
			}
		}()
		// 继续处理请求
		c.Next()
	}
}
