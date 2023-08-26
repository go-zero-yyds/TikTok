package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

var ErrNotFound = sqlx.ErrNotFound

var (
	UserNotFound      = status.Error(100, "用户不存在")
	UserValidation    = status.Error(200, "密码错误")
	DuplicateUsername = status.Error(300, "用户已存在")
	UnmarshalError    = status.Error(400, "解析错误")
)
