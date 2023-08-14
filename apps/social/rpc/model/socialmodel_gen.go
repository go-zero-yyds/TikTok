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
	socialFieldNames          = builder.RawFieldNames(&Social{})
	socialRows                = strings.Join(socialFieldNames, ",")
	socialRowsExpectAutoSet   = strings.Join(stringx.Remove(socialFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	socialRowsWithPlaceHolder = strings.Join(stringx.Remove(socialFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	socialModel interface {
		Insert(ctx context.Context, data *Social) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Social, error)
		Update(ctx context.Context, data *Social) error
		Delete(ctx context.Context, id int64) error
	}

	defaultSocialModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Social struct {
		Id            int64 `db:"id"`             // 字段ID
		UserId        int64 `db:"user_id"`        // 用户ID，由雪花算法生成
		FollowCount   int64 `db:"follow_count"`   // 关注总数
		FollowerCount int64 `db:"follower_count"` // 粉丝总数
	}
)

func newSocialModel(conn sqlx.SqlConn) *defaultSocialModel {
	return &defaultSocialModel{
		conn:  conn,
		table: "`social`",
	}
}

func (m *defaultSocialModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultSocialModel) FindOne(ctx context.Context, id int64) (*Social, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", socialRows, m.table)
	var resp Social
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

func (m *defaultSocialModel) Insert(ctx context.Context, data *Social) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, socialRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.FollowCount, data.FollowerCount)
	return ret, err
}

func (m *defaultSocialModel) Update(ctx context.Context, data *Social) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, socialRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.FollowCount, data.FollowerCount, data.Id)
	return err
}

func (m *defaultSocialModel) tableName() string {
	return m.table
}
