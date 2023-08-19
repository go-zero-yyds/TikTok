package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FollowModel = (*customFollowModel)(nil)
var (
	cacheFollowCountIdPrefix   = "cache:follow:count:id:"
	cacheFollowerCountIdPrefix = "cache:follower:count:id:"
)

type (
	// FollowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowModel.
	FollowModel interface {
		followModel
		InsertIgnore(ctx context.Context, data *Follow) (sql.Result, error)
		InsertOrUpdate(ctx context.Context, data *Follow) (sql.Result, error)
		FindByFollowCount(ctx context.Context, id int64) (int64, error)
		FindByFollowerCount(ctx context.Context, id int64) (int64, error)
		FindFollowList(ctx context.Context, userId int64) ([]int64, error)
		FindFollowerList(ctx context.Context, userId int64) ([]int64, error)
		FindFriendList(ctx context.Context, userId int64) ([]int64, error)
	}

	customFollowModel struct {
		*defaultFollowModel
	}
)

// NewFollowModel returns a model for the database table.
func NewFollowModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FollowModel {
	return &customFollowModel{
		defaultFollowModel: newFollowModel(conn, c, opts...),
	}
}

func (m *defaultFollowModel) InsertIgnore(ctx context.Context, data *Follow) (sql.Result, error) {
	followIdKey := fmt.Sprintf("%s%v", cacheFollowIdPrefix, data.Id)
	followUserIdToUserIdKey := fmt.Sprintf("%s%v:%v", cacheFollowUserIdToUserIdPrefix, data.UserId, data.ToUserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert ignore into %s (%s) values (?, ?, ?)", m.table, followRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.ToUserId, data.Behavior)
	}, followIdKey, followUserIdToUserIdKey)
	return ret, err
}

func (m *defaultFollowModel) InsertOrUpdate(ctx context.Context, data *Follow) (sql.Result, error) {
	followIdKey := fmt.Sprintf("%s%v", cacheFollowIdPrefix, data.Id)
	followCountIdKey := fmt.Sprintf("%s%v", cacheFollowCountIdPrefix, data.UserId)
	followerCountIdKey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, data.ToUserId)
	followUserIdToUserIdKey := fmt.Sprintf("%s%v:%v", cacheFollowUserIdToUserIdPrefix, data.UserId, data.ToUserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf(`
			INSERT INTO %s (%s)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE behavior = ?;
		`, m.table, followRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.ToUserId, data.Behavior, data.Behavior)
	}, followIdKey, followUserIdToUserIdKey, followCountIdKey, followerCountIdKey)
	return ret, err
}

func (m *defaultFollowModel) FindByFollowCount(ctx context.Context, id int64) (int64, error) {
	followCountIdKey := fmt.Sprintf("%s%v", cacheFollowCountIdPrefix, id)
	var resp int64
	err := m.QueryRowCtx(ctx, &resp, followCountIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `behavior` = '1'", m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultFollowModel) FindByFollowerCount(ctx context.Context, id int64) (int64, error) {
	followCountIdKey := fmt.Sprintf("%s%v", cacheFollowerCountIdPrefix, id)
	var resp int64
	err := m.QueryRowCtx(ctx, &resp, followCountIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select count(*) from %s where `to_user_id` = ? and `behavior` = '1'", m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

// FindFollowList 查看用户关注id列表
func (m *defaultFollowModel) FindFollowList(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("select to_user_id from %s where `user_id` = ? and behavior = '1'", m.table)
	var resp []int64
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindFollowerList 查看用户粉丝id列表
func (m *defaultFollowModel) FindFollowerList(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("select user_id from %s where `to_user_id` = ? and behavior = '1'", m.table)
	var resp []int64
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindFriendList 查看用户粉丝id列表
func (m *defaultFollowModel) FindFriendList(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT f1.to_user_id
		FROM %s f1
		JOIN %s f2 ON f1.to_user_id = f2.user_id AND f1.user_id = f2.to_user_id
		WHERE f1.user_id = ? AND f1.behavior = '1' AND f2.behavior = '1';
	`, m.table, m.table)
	var resp []int64
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
