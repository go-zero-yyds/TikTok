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

// 定义用户状态常量
const (
	StatusStranger = "0"
	StatusFollow   = "1"
	StatusFan      = "2"
	StatusFriend   = "3"
)

// 定义关注类型常量
const (
	FollowTypeNotFollowing = "0"
	FollowTypeFollowing    = "1"
)

var (
	_                               FollowModel = (*customFollowModel)(nil)
	cacheFollowUserIdFollowPrefix               = "cache:follow:userId:follow:"
	cacheFollowUserIdFollowerPrefix             = "cache:follow:userId:follower:"
	cacheFollowUserIdFriendPrefix               = "cache:follow:userId:friend:"
)

type (
	// FollowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFollowModel.
	FollowModel interface {
		followModel
		FindFollowList(ctx context.Context, userId int64) ([]int64, error)
		FindFollowerList(ctx context.Context, userId int64) ([]int64, error)
		FindFriendList(ctx context.Context, userId int64) ([]int64, error)
		TranLockByUserIdToUserId(ctx context.Context, s sqlx.Session, userId int64, toUserId int64) ([]Follow, error)
		TranInsert(ctx context.Context, s sqlx.Session, data *Follow, keys *[]string) (sql.Result, error)
		TranUpdate(ctx context.Context, s sqlx.Session, newData *Follow, keys *[]string) (sql.Result, error)
		StateMachine(userStatus, userFollowType string) (newUserStatus string, newToUserStatus string)
		AddFriendKey(userId, toUserId int64, keys *[]string)
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

func (m *defaultFollowModel) TranLockByUserIdToUserId(ctx context.Context, s sqlx.Session, userId int64, toUserId int64) ([]Follow, error) {
	//followUserIdToUserIdKey := fmt.Sprintf("%s%v:%v", cacheFollowUserIdToUserIdPrefix, userId, toUserId)
	var resp []Follow
	query := fmt.Sprintf(`select %s from %s where
 			(user_id = ? and to_user_id = ?) or (to_user_id = ? and user_id = ?)
            for update`, followRows, m.table)
	err := s.QueryRowsCtx(ctx, &resp, query, userId, toUserId, userId, toUserId)
	if err != nil {
		return nil, err
	}

	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFollowModel) TranInsert(ctx context.Context, s sqlx.Session, data *Follow, keys *[]string) (sql.Result, error) {
	followIdKey := fmt.Sprintf("%s%v", cacheFollowIdPrefix, data.Id)
	followUserIdToUserIdKey := fmt.Sprintf("%s%v:%v", cacheFollowUserIdToUserIdPrefix, data.UserId, data.ToUserId)
	followListKey := fmt.Sprintf("%s%v", cacheFollowUserIdFollowPrefix, data.UserId)
	followerListKey := fmt.Sprintf("%s%v", cacheFollowUserIdFollowPrefix, data.ToUserId)
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, followRowsExpectAutoSet)
	ret, err := s.ExecCtx(ctx, query, data.UserId, data.ToUserId, data.Behavior, data.Attribute)
	if err != nil {
		return nil, err
	}
	*keys = append(*keys, followIdKey, followUserIdToUserIdKey, followListKey, followerListKey)
	return ret, err
}
func (m *defaultFollowModel) AddFriendKey(userId, toUserId int64, keys *[]string) {
	*keys = append(*keys, fmt.Sprintf("%s%v", cacheFollowUserIdFriendPrefix, userId),
		fmt.Sprintf("%s%v", cacheFollowUserIdFriendPrefix, toUserId))
}

// TranUpdate 更新关注记录没有则不操作
func (m *defaultFollowModel) TranUpdate(ctx context.Context, s sqlx.Session, newData *Follow, keys *[]string) (result sql.Result, err error) {
	followIdKey := fmt.Sprintf("%s%v", cacheFollowIdPrefix, newData.Id)
	//followUserIdToUserIdKey := fmt.Sprintf("%s%v:%v", cacheFollowUserIdToUserIdPrefix, newData.UserId, newData.ToUserId)
	followListKey := fmt.Sprintf("%s%v", cacheFollowUserIdFollowPrefix, newData.UserId)
	followerListKey := fmt.Sprintf("%s%v", cacheFollowUserIdFollowPrefix, newData.ToUserId)
	query := fmt.Sprintf("update %s set %s where `user_id` = ? and to_user_id = ?", m.table, followRowsWithPlaceHolder)
	ret, err := s.ExecCtx(ctx, query, newData.UserId, newData.ToUserId, newData.Behavior, newData.Attribute, newData.UserId, newData.ToUserId)
	*keys = append(*keys, followIdKey, followListKey, followerListKey)
	return ret, err
}

// FindFollowList 查看用户关注id列表
func (m *defaultFollowModel) FindFollowList(ctx context.Context, userId int64) ([]int64, error) {
	var resp []int64
	key := fmt.Sprintf("%s%v", cacheFollowUserIdFollowPrefix, userId)
	if err := m.CachedConn.GetCache(key, &resp); err == nil {
		return resp, nil
	}

	query := fmt.Sprintf("select to_user_id from %s where `user_id` = ? and (attribute = '%s' or attribute = '%s') LIMIT 1000",
		m.table, StatusFollow, StatusFriend)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		_ = m.CachedConn.SetCacheCtx(ctx, key, &resp)
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindFollowerList 查看用户粉丝id列表
func (m *defaultFollowModel) FindFollowerList(ctx context.Context, userId int64) ([]int64, error) {
	var resp []int64
	key := fmt.Sprintf("%s%v", cacheFollowUserIdFollowerPrefix, userId)
	if err := m.CachedConn.GetCache(key, &resp); err == nil {
		return resp, nil
	}

	query := fmt.Sprintf("select to_user_id from %s where `user_id` = ? and (attribute = '%s' or attribute = '%s') LIMIT 1000",
		m.table, StatusFan, StatusFriend)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		_ = m.CachedConn.SetCacheCtx(ctx, key, &resp)
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindFriendList 查看用户好友id列表
func (m *defaultFollowModel) FindFriendList(ctx context.Context, userId int64) ([]int64, error) {
	var resp []int64
	key := fmt.Sprintf("%s%v", cacheFollowUserIdFriendPrefix, userId)
	if err := m.CachedConn.GetCache(key, &resp); err == nil {
		return resp, nil
	}

	query := fmt.Sprintf("select to_user_id from %s where `user_id` = ? and attribute = '%s' LIMIT 1000", m.table, StatusFriend)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		_ = m.CachedConn.SetCacheCtx(ctx, key, &resp)
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// StateMachine 根据用户状态和关注类型返回新的用户状态
func (m *defaultFollowModel) StateMachine(userStatus, userFollowType string) (newUserStatus string, newToUserStatus string) {
	switch userStatus {
	case StatusStranger:
		if userFollowType == FollowTypeFollowing {
			newUserStatus = StatusFollow
			newToUserStatus = StatusFan
		} else {
			newUserStatus = StatusStranger
			newToUserStatus = StatusStranger
		}
	case StatusFollow:
		if userFollowType == FollowTypeFollowing {
			newUserStatus = StatusFollow
			newToUserStatus = StatusFan
		} else {
			newUserStatus = StatusStranger
			newToUserStatus = StatusStranger
		}
	case StatusFan:
		if userFollowType == FollowTypeFollowing {
			newUserStatus = StatusFriend
			newToUserStatus = StatusFriend
		} else {
			newUserStatus = StatusFan
			newToUserStatus = StatusFollow
		}
	case StatusFriend:
		if userFollowType == FollowTypeFollowing {
			newUserStatus = StatusFriend
			newToUserStatus = StatusFriend
		} else {
			newUserStatus = StatusFan
			newToUserStatus = StatusFollow
		}
	default:
		newUserStatus = StatusStranger
		newToUserStatus = StatusStranger
	}

	return newUserStatus, newToUserStatus
}
