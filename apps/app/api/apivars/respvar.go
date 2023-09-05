package apivars

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
		StatusCode: 0,
		StatusMsg:  "部分数据异常",
	}
	UserNotFound = RespErr{
		StatusCode: 100,
		StatusMsg:  "用户不存在",
	}
	UserValidation = RespErr{
		StatusCode: 101,
		StatusMsg:  "密码错误",
	}
	DuplicateUsername = RespErr{
		StatusCode: 102,
		StatusMsg:  "用户已存在",
	}
	NotLogged = RespErr{
		StatusCode: 104,
		StatusMsg:  "请登录",
	}
	TokenSignatureInvalid = RespErr{
		StatusCode: 105,
		StatusMsg:  "token无效，尝试重新登录",
	}
	UsernameRuleError = RespErr{
		StatusCode: 106,
		StatusMsg:  "用户名格式不符合规范",
	}
	PasswordRuleError = RespErr{
		StatusCode: 107,
		StatusMsg:  "密码格式不符合规范",
	}
	VideoNotFound = RespErr{
		StatusCode: 200,
		StatusMsg:  "视频不存在",
	}
	DataNotVideo = RespErr{
		StatusCode: 201,
		StatusMsg:  "不是视频",
	}
	AlreadyLiked = RespErr{
		StatusCode: 300,
		StatusMsg:  "已经点赞过了",
	}
	AlreadyUnLiked = RespErr{
		StatusCode: 301,
		StatusMsg:  "已经取消过点赞了",
	}
	ErrNotFriend = RespErr{
		StatusCode: 400,
		StatusMsg:  "不是好友",
	}
	IllegalArgument = RespErr{
		StatusCode: 401,
		StatusMsg:  "非法的请求",
	}
	NoFollowMyself = RespErr{
		StatusCode: 402,
		StatusMsg:  "不能关注自己",
	}
	TextRuleError = RespErr{
		StatusCode: 2400,
		StatusMsg:  "内容格式不符合规范",
	}
	// InternalError 内部异常造成的错误 兜底用
	InternalError = RespErr{
		StatusCode: 500,
		StatusMsg:  "内部错误",
	}
)
