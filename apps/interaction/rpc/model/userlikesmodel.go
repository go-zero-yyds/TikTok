package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserLikesModel = (*customUserLikesModel)(nil)

type (
	// UserLikesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLikesModel.
	UserLikesModel interface {
		userLikesModel
		TranIncrCount(ctx context.Context, s sqlx.Session, userId int64) error
		TranDecrCount(ctx context.Context, s sqlx.Session, userId int64) error
	}

	customUserLikesModel struct {
		*defaultUserLikesModel
	}
)

// NewUserLikesModel returns a model for the database table.
func NewUserLikesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserLikesModel {
	return &customUserLikesModel{
		defaultUserLikesModel: newUserLikesModel(conn, c, opts...),
	}
}

func (m *defaultUserLikesModel) TranIncrCount(ctx context.Context, s sqlx.Session, userId int64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 1)
		ON DUPLICATE KEY UPDATE like_count = like_count + 1;`, m.table, userLikesRowsExpectAutoSet)
	_, err := s.ExecCtx(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *defaultUserLikesModel) TranDecrCount(ctx context.Context, s sqlx.Session, userId int64) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET like_count = like_count - 1
		WHERE user_id = ?;`, m.table)
	_, err := s.ExecCtx(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}
