package errors

const (
	// 客户端错误
	CodeLoginInvalid = 401
	CodeNoRightErr   = 403

	// 公用错误码
	CodeParamsError  = 1001
	CodeJsonError    = 1002
	CodeEncryptError = 1003

	// 用户相关
	CodeUserDuplicatedNickname = 10001
	CodeUserDuplicatedUsername = 10002
	CodeUserLoginFail          = 10003
)

const (
	MsgLoginInvalid = "登录状态已失效，请重新登录"
	MsgNoRightErr   = "没有相关权限"

	MsgParamsError  = "参数错误"
	MsgJsonError    = "json类型不正确"
	MsgEncryptError = "加密错误"

	MsgUserDuplicatedNickname = "昵称被占用"
	MsgUserDuplicatedUsername = "用户名已经注册"
	MsgUserLoginFail          = "账号或密码错误"
)
