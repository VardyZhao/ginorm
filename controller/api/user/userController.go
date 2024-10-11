package user

import (
	"ginorm/controller/api"
	"ginorm/request"
	"ginorm/serializer"
	"ginorm/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Register 用户注册接口
func Register(c *gin.Context) {
	var registerRequest request.RegisterRequest
	if err := c.ShouldBind(&registerRequest); err == nil {
		if err := registerRequest.Valid(); err != nil {
			c.JSON(200, err)
			return
		}
	}

	userService := service.UserService{}
	res := userService.Register(registerRequest)
	c.JSON(200, res)
}

// Login 用户登录接口
func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.JSON(200, api.ErrorResponse(err))
		return
	}

	userService := service.UserService{}
	res := userService.Login(loginRequest, c)
	c.JSON(200, res)
}

// Profile 用户详情
func Profile(c *gin.Context) {
	user := api.CurrentUser(c)
	res := serializer.BuildUserResponse(*user)
	c.JSON(200, res)
}

// Logout 用户登出
func Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "登出成功",
	})
}
