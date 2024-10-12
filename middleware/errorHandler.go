package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"ginorm/config"
	"ginorm/entity/response"
	bussErrors "ginorm/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler 统一拦截业务错误并返回 JSON 响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 业务错误处理
			var be *bussErrors.BusinessError
			if errors.As(err, &be) {
				c.JSON(200, &response.Response{
					Code: be.Code,
					Data: be.Data,
					Msg:  be.Msg,
				})
			}

			// 参数校验错误处理
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				// 表单校验错误
				formError := ""
				for _, e := range ve {
					field := config.T(fmt.Sprintf("Field.%s", e.Field()))
					tag := config.T(fmt.Sprintf("Tag.Valid.%s", e.Tag()))
					formError += fmt.Sprintf("%s%s", field, tag) + "; "
				}
				c.JSON(200, &response.Response{
					Code: bussErrors.CodeParamsError,
					Data: formError,
					Msg:  bussErrors.MsgParamsError,
				})
			}

			// json格式校验处理
			var je *json.UnmarshalTypeError
			if errors.As(err, &je) {
				c.JSON(200, &response.Response{
					Code: bussErrors.CodeJsonError,
					Data: be.Data,
					Msg:  bussErrors.MsgJsonError,
				})
			}
		}
	}
}
