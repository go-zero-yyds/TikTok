package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ MessageModel = (*customMessageModel)(nil)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		FindNowMessage(ctx context.Context, userId int64, toUserId int64) (*Message, error)
		FindMessageList(ctx context.Context, userId int64, toUserId, preTime int64) ([]Message, error)
	}

	customMessageModel struct {
		*defaultMessageModel
	}
)

// NewMessageModel returns a model for the database table.
func NewMessageModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MessageModel {
	return &customMessageModel{
		defaultMessageModel: newMessageModel(conn, c, opts...),
	}
}

func (m *defaultMessageModel) FindNowMessage(ctx context.Context, userId int64, toUserId int64) (*Message, error) {

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s m
		WHERE (m.from_user_id = ? AND m.to_user_id = ?)
		   OR (m.from_user_id = ? AND m.to_user_id = ?)
		ORDER BY m.create_time DESC
		LIMIT 1;
	`, messageRows, m.table)
	var resp Message
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, toUserId, toUserId, userId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return &resp, ErrNotFound
	default:
		return nil, err
	}

}

func (m *defaultMessageModel) FindMessageList(ctx context.Context, userId int64, toUserId, preTime int64) ([]Message, error) {
	t := time.UnixMilli(preTime)
	query := fmt.Sprintf(`
		SELECT %s
		FROM %s m
		WHERE (m.from_user_id = ? AND m.to_user_id = ? AND m.create_time > ?)
		   OR (m.from_user_id = ? AND m.to_user_id = ? AND m.create_time > ?)
		ORDER BY m.create_time ASC LIMIT 1000;
	`, messageRows, m.table)
	var resp []Message
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, toUserId, t, toUserId, userId, t)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return resp, ErrNotFound
	default:
		return nil, err
	}

}
