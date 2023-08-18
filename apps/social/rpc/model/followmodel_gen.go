// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	followFieldNames          = builder.RawFieldNames(&Follow{})
	followRows                = strings.Join(followFieldNames, ",")
	followRowsExpectAutoSet   = strings.Join(stringx.Remove(followFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	followRowsWithPlaceHolder = strings.Join(stringx.Remove(followFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	followModel interface {
		Insert(ctx context.Context, data *Follow) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Follow, error)
		Update(ctx context.Context, data *Follow) error
		Delete(ctx context.Context, id int64) error
	}

	defaultFollowModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Follow struct {
		Id       int64 `db:"id"`         // 字段ID
		UserId   int64 `db:"user_id"`    // 用户ID
		ToUserId int64 `db:"to_user_id"` // 关注者ID
		Status   []uint8  `db:"status"`     // 关注状态 0=>没关注 1=>关注
	}
)

func newFollowModel(conn sqlx.SqlConn) *defaultFollowModel {
	return &defaultFollowModel{
		conn:  conn,
		table: "`follow`",
	}
}

func (m *defaultFollowModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultFollowModel) FindOne(ctx context.Context, id int64) (*Follow, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", followRows, m.table)
	var resp Follow
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFollowModel) Insert(ctx context.Context, data *Follow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, followRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ToUserId, data.Status)
	return ret, err
}

func (m *defaultFollowModel) Update(ctx context.Context, data *Follow) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, followRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ToUserId, data.Status, data.Id)
	return err
}

func (m *defaultFollowModel) tableName() string {
	return m.table
}