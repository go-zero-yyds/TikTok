package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		UpdateByUserId(ctx context.Context, userId, field, value string) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *defaultUserModel) UpdateByUserId(ctx context.Context, userId, field, value string) error {
	var result int64
	update := fmt.Sprintf("UPDATE %s SET %s = ? WHERE user_id = ?", m.table, field)
	err := m.QueryRowsNoCacheCtx(ctx, &result, update, value, userId)
	if err != nil {
		return err
	}

	return nil
}
