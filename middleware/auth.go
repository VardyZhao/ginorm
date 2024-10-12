package middleware

import (
	"ginorm/entity/response"
	"ginorm/errors"
	"ginorm/model"
	"ginorm/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid := session.Get("user_id")
		if uid != nil {
			us := service.UserService{}
			user, err := us.GetUser(uid)
			if err == nil {
				c.Set("user", &user)
			}
		}
		c.Next()
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get("user"); user != nil {
			if _, ok := user.(*model.User); ok {
				c.Next()
				return
			}
		}

		c.JSON(200, &response.Response{
			Code: errors.CodeLoginInvalid,
			Msg:  errors.MsgLoginInvalid,
		})
		c.Abort()
	}
}
