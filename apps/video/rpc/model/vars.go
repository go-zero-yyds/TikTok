package model

import (
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrVideoNotFound = status.Error(100, "video not found")
	ErrNotFound      = sqlx.ErrNotFound
)
