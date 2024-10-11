package request

import (
	"ginorm/serializer"
)

type RegisterRequest struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=2,max=30"`
	Username        string `form:"username" json:"username" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

// Valid 验证表单
func (r *RegisterRequest) Valid() *serializer.Response {
	if r.PasswordConfirm != r.Password {
		return &serializer.Response{
			Code: 40001,
			Msg:  "两次输入的密码不相同",
		}
	}

	return nil
}
