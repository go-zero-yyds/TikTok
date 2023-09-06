package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrNotFriend = status.Error(100, "not friend")
