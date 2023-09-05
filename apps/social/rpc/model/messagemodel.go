package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	_                                    MessageModel = (*customMessageModel)(nil)
	cacheMessageUserIdToUserIdPrefix                  = "cache:message:id:userId:toUserId:"
	cacheMessageUserIdToUserIdListPrefix              = "cache:message:id:userId:toUserId:list"
)

type (
	// MessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMessageModel.
	MessageModel interface {
		messageModel
		FindNowMessage(ctx context.Context, userId int64, toUserId int64) (*Message, error)
		FindMessageList(ctx context.Context, userId int64, toUserId, preTime int64) ([]Message, error)
		GetNowMessageCacheKey(userId int64, toUserId int64) string
		GetNowMessageListCacheKey(userId int64, toUserId int64) string
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

// func (m *defaultMessageModel) FindNowMessage(ctx context.Context, userId int64, toUserId int64) (*Message, error) {
//
//	query := fmt.Sprintf(`
//		SELECT %s
//		FROM %s m
//		WHERE (m.from_user_id = ? AND m.to_user_id = ?)
//		   OR (m.from_user_id = ? AND m.to_user_id = ?)
//		ORDER BY m.create_time DESC
//		LIMIT 1;
//	`, messageRows, m.table)
//	var resp Message
//	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, toUserId, toUserId, userId)
//	switch {
//	case err == nil:
//		return &resp, nil
//	case errors.Is(err, sqlc.ErrNotFound):
//		return &resp, ErrNotFound
//	default:
//		return nil, err
//	}
//
// }

func (m *defaultMessageModel) GetNowMessageCacheKey(userId int64, toUserId int64) string {
	if userId > toUserId {
		userId = userId ^ toUserId
		toUserId = userId ^ toUserId
		userId = userId ^ toUserId
	}
	return fmt.Sprintf("%s%v:%v", cacheMessageUserIdToUserIdPrefix, userId, toUserId)
}
func (m *defaultMessageModel) GetNowMessageListCacheKey(userId int64, toUserId int64) string {
	if userId > toUserId {
		userId = userId ^ toUserId
		toUserId = userId ^ toUserId
		userId = userId ^ toUserId
	}
	return fmt.Sprintf("%s%v:%v", cacheMessageUserIdToUserIdListPrefix, userId, toUserId)
}
func (m *defaultMessageModel) FindNowMessage(ctx context.Context, userId int64, toUserId int64) (*Message, error) {
	messageUserIdToUserIdKey := m.GetNowMessageCacheKey(userId, toUserId)
	var resp Message
	err := m.QueryRowIndexCtx(ctx, &resp, messageUserIdToUserIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf(`
			SELECT %s
			FROM %s m
			WHERE (m.from_user_id = ? AND m.to_user_id = ?)
			   OR (m.from_user_id = ? AND m.to_user_id = ?)
			ORDER BY m.create_time DESC
			LIMIT 1;
		`, messageRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, userId, toUserId, toUserId, userId); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultMessageModel) FindMessageList(ctx context.Context, userId int64, toUserId, preTime int64) ([]Message, error) {
	t := time.UnixMilli(preTime)
	var resp []Message
	key := m.GetNowMessageListCacheKey(userId, toUserId)
	if preTime == 0 {
		if err := m.CachedConn.GetCache(key, &resp); err == nil {
			return resp, nil
		}
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM (
				SELECT %s
				FROM %s m
				WHERE (m.from_user_id = ? AND m.to_user_id = ? AND m.create_time > ?)
				   OR (m.from_user_id = ? AND m.to_user_id = ? AND m.create_time > ?)
				ORDER BY m.create_time DESC LIMIT 1000
		) AS latest_messages
		ORDER BY create_time ASC;
	`, messageRows, messageRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, toUserId, t, toUserId, userId, t)
	switch {
	case err == nil:
		if preTime == 0 {
			_ = m.CachedConn.SetCacheCtx(ctx, key, &resp)
		}
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return resp, ErrNotFound
	default:
		return nil, err
	}

}
