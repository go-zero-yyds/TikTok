package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrVideoNotFound = errors.New("video not found")
	ErrNotFound      = sqlx.ErrNotFound
)
