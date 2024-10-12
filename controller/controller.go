package controller

import (
	"ginorm/entity/response"
	bussErrors "ginorm/errors"
	"github.com/gin-gonic/gin"
)

func FailWithErr(ctx *gin.Context, err error) {
	ctx.Error(err)
}

func Fail(ctx *gin.Context, args ...interface{}) {
	if len(args) < 2 {
		panic("Controller Fail function's params error")
	}
	code := args[0].(int)
	msg := args[1].(string)
	data := args[2]

	ctx.Error(bussErrors.NewBusinessError(code, msg, data))
}

func Success(ctx *gin.Context, data interface{}, msg ...string) {
	success := "ok"
	if len(msg) > 0 {
		success = msg[0]
	}
	ctx.JSON(200, &response.Response{
		Code: 0,
		Msg:  success,
		Data: data,
	})
}
