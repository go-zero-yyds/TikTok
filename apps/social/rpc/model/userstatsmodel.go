package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserStatsModel = (*customUserStatsModel)(nil)

type (
	// UserStatsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserStatsModel.
	UserStatsModel interface {
		userStatsModel
		TranIncrCount(ctx context.Context, s sqlx.Session, userId, toUserId int64, keys *[]string) error
		TranDecrCount(ctx context.Context, s sqlx.Session, userId, toUserId int64, keys *[]string) error
	}

	customUserStatsModel struct {
		*defaultUserStatsModel
	}
)

// NewUserStatsModel returns a model for the database table.
func NewUserStatsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserStatsModel {
	return &customUserStatsModel{
		defaultUserStatsModel: newUserStatsModel(conn, c, opts...),
	}
}

func (m *defaultUserStatsModel) TranIncrCount(ctx context.Context, s sqlx.Session, userId, toUserId int64, keys *[]string) error {
	userStatsUserIdKey := fmt.Sprintf("%s%v", cacheUserStatsUserIdPrefix, userId)
	toUserStatsUserIdKey := fmt.Sprintf("%s%v", cacheUserStatsUserIdPrefix, toUserId)
	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 1, 0)
		ON DUPLICATE KEY UPDATE follow_count = follow_count + 1;`, m.table, userStatsRowsExpectAutoSet)
	_, err := s.ExecCtx(ctx, query, userId)
	query = fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 0, 1)
		ON DUPLICATE KEY UPDATE follower_count = follower_count + 1;`, m.table, userStatsRowsExpectAutoSet)
	_, err = s.ExecCtx(ctx, query, toUserId)
	*keys = append(*keys, userStatsUserIdKey, toUserStatsUserIdKey)
	return err
}

func (m *defaultUserStatsModel) TranDecrCount(ctx context.Context, s sqlx.Session, userId, toUserId int64, keys *[]string) error {
	userStatsUserIdKey := fmt.Sprintf("%s%v", cacheUserStatsUserIdPrefix, userId)
	toUserStatsUserIdKey := fmt.Sprintf("%s%v", cacheUserStatsUserIdPrefix, toUserId)
	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 1, 0)
		ON DUPLICATE KEY UPDATE follow_count = follow_count - 1;`, m.table, userStatsRowsExpectAutoSet)
	_, err := s.ExecCtx(ctx, query, userId)
	query = fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 0, 1)
		ON DUPLICATE KEY UPDATE follower_count = follower_count - 1;`, m.table, userStatsRowsExpectAutoSet)
	_, err = s.ExecCtx(ctx, query, toUserId)
	*keys = append(*keys, userStatsUserIdKey, toUserStatsUserIdKey)
	return err
}
