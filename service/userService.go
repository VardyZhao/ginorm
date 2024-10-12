package service

import (
	"ginorm/constant"
	"ginorm/model"
	"ginorm/repository"
	"ginorm/request"
	"ginorm/serializer"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserService 管理用户注册服务
type UserService struct{}

// Register 用户注册
func (service *UserService) Register(registerRequest request.RegisterRequest) serializer.Response {

	repo := repository.NewUserRepository()

	// 校验nickname唯一性
	count, err := repo.CountUserByNickname(registerRequest.Nickname)
	if err != nil {
		return serializer.ParamErr("校验昵称失败", err)
	}
	if count > 0 {
		return serializer.Response{
			Code: 40001,
			Msg:  "昵称被占用",
		}
	}

	// 校验username唯一性
	count, err = repo.CountUserByNickname(registerRequest.Username)
	if err != nil {
		return serializer.ParamErr("校验用户名失败", err)
	}
	if count > 0 {
		return serializer.Response{
			Code: 40001,
			Msg:  "用户名已经注册",
		}
	}

	user := model.User{
		Nickname: registerRequest.Nickname,
		Username: registerRequest.Username,
		Status:   constant.UserStatusActive,
	}
	// 加密密码
	if err := service.setPassword(user, registerRequest.Password); err != nil {
		return serializer.Err(
			serializer.CodeEncryptError,
			"密码加密失败",
			err,
		)
	}

	// 创建用户
	if err := repo.CreateUser(&user); err != nil {
		return serializer.ParamErr("注册失败", err)
	}

	return serializer.BuildUserResponse(user)
}

// Login 用户登录函数
func (service *UserService) Login(loginRequest request.LoginRequest, c *gin.Context) serializer.Response {
	repo := repository.NewUserRepository()
	user, err := repo.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	if service.checkPassword(user, loginRequest.Password) == false {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	// 设置session
	service.setSession(c, user)

	return serializer.BuildUserResponse(user)
}

func (service *UserService) Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
}

func (service *UserService) GetUser(ID interface{}) (model.User, error) {
	repo := repository.NewUserRepository()
	user, err := repo.GetUser(ID)
	return user, err
}

func (service *UserService) setSession(c *gin.Context, user model.User) {
	s := sessions.Default(c)
	s.Clear()
	s.Set("user_id", user.ID)
	s.Save()
}

// SetPassword 设置密码
func (service *UserService) setPassword(user model.User, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constant.PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (service *UserService) checkPassword(user model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
