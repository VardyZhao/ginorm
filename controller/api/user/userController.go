package user

import (
	"ginorm/controller"
	"ginorm/entity/dto"
	"ginorm/entity/request"
	"ginorm/model"
	"ginorm/service"

	"github.com/gin-gonic/gin"
)

// Register 用户注册接口
func Register(c *gin.Context) {
	var registerRequest request.RegisterRequest
	if err := c.ShouldBind(&registerRequest); err != nil {
		controller.FailWithErr(c, err)
		return
	}
	if err := registerRequest.Valid(); err != nil {
		controller.FailWithErr(c, err)
		return
	}

	userService := service.UserService{}
	res, err := userService.Register(registerRequest)
	if err != nil {
		controller.FailWithErr(c, err)
		return
	}
	controller.Success(c, res)
}

// Login 用户登录接口
func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		controller.FailWithErr(c, err)
		return
	}

	userService := service.UserService{}
	res, loginErr := userService.Login(loginRequest, c)
	if loginErr != nil {
		controller.FailWithErr(c, loginErr)
		return
	}
	controller.Success(c, res)
}

// Profile 用户详情
func Profile(c *gin.Context) {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			res := dto.BuildUserDTO(u)
			controller.Success(c, res)
		}
	}
}

// Logout 用户登出
func Logout(c *gin.Context) {
	userService := service.UserService{}
	userService.Logout(c)
	controller.Success(c, nil, "登出成功")
}
