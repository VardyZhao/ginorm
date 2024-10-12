package api

import (
	"ginorm/controller"
	"github.com/gin-gonic/gin"
)

// Ping 状态检查页面
func Ping(ctx *gin.Context) {
	controller.Success(ctx, nil, "Pong")
}
