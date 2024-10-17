package constant

const (
	// TraceId trace id字段名
	TraceId = "trace_id"

	// HeaderTraceId 传入的header头的trace id字段名
	HeaderTraceId = "X-Trace-Id"

	// PasswordCost 密码加密难度
	PasswordCost = 12

	// UserStatusActive 激活用户
	UserStatusActive string = "active"

	// UserStatusInactive 未激活用户
	UserStatusInactive string = "inactive"

	// UserStatusSuspend 被封禁用户
	UserStatusSuspend string = "suspend"
)
