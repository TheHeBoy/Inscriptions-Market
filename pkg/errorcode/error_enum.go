package errorcode

// 通用错误码 1000
var (
	BAD_REQUEST            = ErrorCode{1000, "请求错误"}
	JSON_FORMAT_ERROR      = ErrorCode{1001, "Json 格式错误"}
	BAD_REQUEST_LIMIT      = ErrorCode{1002, "请求限制"}
	SERVICE_INTERNAL_ERROR = ErrorCode{1003, "服务器内存错误"}
)

// 权限模块 2000
var (
	AUTH_ACCOUNT_OR_PASSWD_NOT_EXIST = ErrorCode{2001, "账号不存在或密码错误"}
	AUTH_TOKEN_REFRESH_FAIL          = ErrorCode{2002, "令牌刷新失败"}
	AUTH_JWT_UNAUTHORIZED            = ErrorCode{2003, "授权错误"}
)

// 用户模块 3000
var (
	USER_NO_EXIST    = ErrorCode{3001, "用户不存在"}
	USER_CREATE_FAIL = ErrorCode{3002, "用户创建失败"}
	USER_UPDATE_FAIL = ErrorCode{3002, "用户更新失败"}
)

// 短信模块 4000
var (
	SMS_SEND_VERIFYCODE_FAIL = ErrorCode{4001, "发送短信验证码失败"}
)

// Email 模块 5000
var (
	Email_SEND_VERIFYCODE_FAIL = ErrorCode{5001, "发送 Email 验证码失败"}
)
