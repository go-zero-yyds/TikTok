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
		StatusCode: 0,
		StatusMsg:  "部分数据异常",
	}
	// InternalError 内部异常造成的错误
	InternalError = RespErr{
		StatusCode: 500,
		StatusMsg:  "内部错误",
	}
)
