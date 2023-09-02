package model

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrIllegalArgument = errors.New("非法操作") // 尝试删除不是自己发的评论
