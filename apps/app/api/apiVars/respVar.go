package apiVars

type RespErr struct {
	StatusCode int32
	StatusMsg  string
}

func (r RespErr) Error() string {
	return r.StatusMsg
}

var (
	// Success 成功时返回
	Success = RespErr{
		StatusCode: 0,
		StatusMsg:  "成功",
	}
	// SomeDataErr 部分数据获取失败
	SomeDataErr = RespErr{
		StatusCode: 120,
		StatusMsg:  "部分数据异常",
	}
	UserNotFound = RespErr{
		StatusCode: 100,
		StatusMsg:  "用户不存在",
	}
	VideoNotFound = RespErr{
		StatusCode: 120,
		StatusMsg:  "视频不存在",
	}
	DataNotVideo = RespErr{
		StatusCode: 130,
		StatusMsg:  "不是视频",
	}
	UserValidation = RespErr{
		StatusCode: 200,
		StatusMsg:  "密码错误",
	}
	DuplicateUsername = RespErr{
		StatusCode: 300,
		StatusMsg:  "用户已存在",
	}
	NotLogged = RespErr{
		StatusCode: 400,
		StatusMsg:  "用户未登录",
	}
	// InternalError 内部异常造成的错误
	InternalError = RespErr{
		StatusCode: 500,
		StatusMsg:  "内部错误",
	}
	UsernameRuleError = RespErr{
		StatusCode: 401,
		StatusMsg:  "用户名格式不符合规范",
	}
	PasswordRuleError = RespErr{
		StatusCode: 402,
		StatusMsg:  "密码格式不符合规范",
	}
	TextRuleError = RespErr{
		StatusCode: 410,
		StatusMsg:  "内容格式不符合规范",
	}
)
