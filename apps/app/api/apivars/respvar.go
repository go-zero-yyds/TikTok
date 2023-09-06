package apivars

type RespVar struct {
	StatusCode int32
	StatusMsg  string
}

func (r RespVar) Error() string {
	return r.StatusMsg
}

var (
	// Success 成功时返回
	Success = RespVar{
		StatusCode: 0,
		StatusMsg:  "成功",
	}
	// ErrSomeData 部分数据获取失败
	ErrSomeData = RespVar{
		StatusCode: 0,
		StatusMsg:  "部分数据异常",
	}
	ErrUserNotFound = RespVar{
		StatusCode: 100,
		StatusMsg:  "用户不存在",
	}
	ErrUserValidation = RespVar{
		StatusCode: 101,
		StatusMsg:  "密码错误",
	}
	ErrDuplicateUsername = RespVar{
		StatusCode: 102,
		StatusMsg:  "用户已存在",
	}
	ErrNotLogged = RespVar{
		StatusCode: 104,
		StatusMsg:  "请登录",
	}
	ErrTokenSignatureInvalid = RespVar{
		StatusCode: 105,
		StatusMsg:  "token无效，尝试重新登录",
	}
	ErrUsernameRule = RespVar{
		StatusCode: 106,
		StatusMsg:  "用户名格式不符合规范",
	}
	ErrPasswordRule = RespVar{
		StatusCode: 107,
		StatusMsg:  "密码格式不符合规范",
	}
	ErrVideoNotFound = RespVar{
		StatusCode: 200,
		StatusMsg:  "视频不存在",
	}
	ErrDataNotVideo = RespVar{
		StatusCode: 201,
		StatusMsg:  "不是视频",
	}
	ErrAlreadyLiked = RespVar{
		StatusCode: 300,
		StatusMsg:  "已经点赞过了",
	}
	ErrAlreadyUnLiked = RespVar{
		StatusCode: 301,
		StatusMsg:  "已经取消过点赞了",
	}
	ErrNotFriend = RespVar{
		StatusCode: 400,
		StatusMsg:  "不是好友",
	}
	ErrIllegalArgument = RespVar{
		StatusCode: 401,
		StatusMsg:  "非法的请求",
	}
	ErrNoFollowMyself = RespVar{
		StatusCode: 402,
		StatusMsg:  "不能关注自己",
	}
	// ErrInternal 内部异常造成的错误 兜底用
	ErrInternal = RespVar{
		StatusCode: 500,
		StatusMsg:  "内部错误",
	}
	ErrTextRuleError = RespVar{
		StatusCode: 6234,
		StatusMsg:  "内容格式不符合规范",
	}
)
