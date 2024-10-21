package service

import (
	"fmt"
	"ginorm/constant"
	"ginorm/entity/dto"
	"ginorm/entity/request"
	"ginorm/errors"
	"ginorm/model"
	"ginorm/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserService 管理用户注册服务
type UserService struct{}

// Register 用户注册
func (service *UserService) Register(registerRequest request.RegisterRequest) (dto.UserDTO, *errors.BusinessError) {

	repo := repository.NewUserRepository()
	fmt.Printf("type: %t\n", repo)
	userDTO := dto.UserDTO{}

	// 校验nickname唯一性
	count, err := repo.CountUserByNickname(registerRequest.Nickname)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return userDTO, errors.NewBusinessError(errors.CodeUserDuplicatedNickname, errors.MsgUserDuplicatedNickname)
	}

	// 校验username唯一性
	count, err = repo.CountUserByNickname(registerRequest.Username)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return userDTO, errors.NewBusinessError(errors.CodeUserDuplicatedUsername, errors.MsgUserDuplicatedUsername)
	}

	user := model.User{
		Nickname: registerRequest.Nickname,
		Username: registerRequest.Username,
		Status:   constant.UserStatusActive,
	}
	// 加密密码
	if err := service.setPassword(&user, registerRequest.Password); err != nil {
		return userDTO, errors.NewBusinessError(errors.CodeEncryptError, errors.MsgEncryptError)
	}

	// 创建用户
	if err := repo.CreateUser(&user); err != nil {
		panic(err)
	}

	return dto.BuildUserDTO(&user), nil
}

// Login 用户登录函数
func (service *UserService) Login(loginRequest request.LoginRequest, c *gin.Context) (dto.UserDTO, *errors.BusinessError) {
	repo := repository.NewUserRepository()
	userDTO := dto.UserDTO{}

	user, err := repo.GetUserByUsername(loginRequest.Username)
	if err != nil {
		return userDTO, errors.NewBusinessError(errors.CodeUserLoginFail, errors.MsgUserLoginFail)
	}

	if service.checkPassword(user, loginRequest.Password) == false {
		return userDTO, errors.NewBusinessError(errors.CodeUserLoginFail, errors.MsgUserLoginFail)
	}

	// 设置session
	service.setSession(c, user)

	return dto.BuildUserDTO(&user), nil
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
func (service *UserService) setPassword(user *model.User, password string) error {
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
