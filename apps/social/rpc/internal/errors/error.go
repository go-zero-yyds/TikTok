package errors

import "errors"

var (
	RecordNotFound   = errors.New("record not found")
	SQLOperateFailed = errors.New("sql operate failed")
	ParamsError      = errors.New("params error")
	Timeout          = errors.New("time out")
)
