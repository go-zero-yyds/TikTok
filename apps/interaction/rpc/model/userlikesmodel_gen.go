// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userLikesFieldNames          = builder.RawFieldNames(&UserLikes{})
	userLikesRows                = strings.Join(userLikesFieldNames, ",")
	userLikesRowsExpectAutoSet   = strings.Join(stringx.Remove(userLikesFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	userLikesRowsWithPlaceHolder = strings.Join(stringx.Remove(userLikesFieldNames, "`user_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheUserLikesUserIdPrefix = "cache:userLikes:userId:"
)

type (
	userLikesModel interface {
		Insert(ctx context.Context, data *UserLikes) (sql.Result, error)
		FindOne(ctx context.Context, userId int64) (*UserLikes, error)
		Update(ctx context.Context, data *UserLikes) error
		Delete(ctx context.Context, userId int64) error
	}

	defaultUserLikesModel struct {
		sqlc.CachedConn
		table string
	}

	UserLikes struct {
		UserId    int64 `db:"user_id"`    // 用户id
		LikeCount int64 `db:"like_count"` // 用户点赞数
	}
)

func newUserLikesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultUserLikesModel {
	return &defaultUserLikesModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`user_likes`",
	}
}

func (m *defaultUserLikesModel) withSession(session sqlx.Session) *defaultUserLikesModel {
	return &defaultUserLikesModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`user_likes`",
	}
}

func (m *defaultUserLikesModel) Delete(ctx context.Context, userId int64) error {
	userLikesUserIdKey := fmt.Sprintf("%s%v", cacheUserLikesUserIdPrefix, userId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, userId)
	}, userLikesUserIdKey)
	return err
}

func (m *defaultUserLikesModel) FindOne(ctx context.Context, userId int64) (*UserLikes, error) {
	userLikesUserIdKey := fmt.Sprintf("%s%v", cacheUserLikesUserIdPrefix, userId)
	var resp UserLikes
	err := m.QueryRowCtx(ctx, &resp, userLikesUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userLikesRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, userId)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserLikesModel) Insert(ctx context.Context, data *UserLikes) (sql.Result, error) {
	userLikesUserIdKey := fmt.Sprintf("%s%v", cacheUserLikesUserIdPrefix, data.UserId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, userLikesRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.LikeCount)
	}, userLikesUserIdKey)
	return ret, err
}

func (m *defaultUserLikesModel) Update(ctx context.Context, data *UserLikes) error {
	userLikesUserIdKey := fmt.Sprintf("%s%v", cacheUserLikesUserIdPrefix, data.UserId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `user_id` = ?", m.table, userLikesRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.LikeCount, data.UserId)
	}, userLikesUserIdKey)
	return err
}

func (m *defaultUserLikesModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheUserLikesUserIdPrefix, primary)
}

func (m *defaultUserLikesModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", userLikesRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserLikesModel) tableName() string {
	return m.table
}
