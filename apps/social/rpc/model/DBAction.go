package model

import (
	"TikTok/apps/social/rpc/social"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type DBAction struct {
	follow     FollowModel
	message    MessageModel
	userStatus UserStatsModel
	conn       sqlc.CachedConn
}

// NewDBAction 初始化数据库信息
func NewDBAction(conn sqlx.SqlConn, c cache.ClusterConf) *DBAction {
	ret := &DBAction{
		follow:     NewFollowModel(conn, c, cache.WithExpiry(24*time.Hour), cache.WithNotFoundExpiry(24*time.Hour)),
		message:    NewMessageModel(conn, c, cache.WithExpiry(24*time.Hour), cache.WithNotFoundExpiry(24*time.Hour)),
		userStatus: NewUserStatsModel(conn, c, cache.WithExpiry(24*time.Hour), cache.WithNotFoundExpiry(24*time.Hour)),
		conn:       sqlc.NewConn(conn, c),
	}
	return ret
}

// FollowAction 调用 follow 对数据库查询用户是否点赞过
// 并对数据库actionType操作
// 出现未知错误和没有修改是false否则都是true
func (d *DBAction) FollowAction(ctx context.Context, userId, toUserId int64, actionType string) (bool, error) {
	if actionType == FollowTypeFollowing {
		return d.DoFollow(ctx, userId, toUserId)
	} else {
		return d.UnFollow(ctx, userId, toUserId)
	}
}

// DoFollow 关注对方
func (d *DBAction) DoFollow(ctx context.Context, userId, toUserId int64) (bool, error) {

	var res bool
	keys := make([]string, 0)

	err := d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		user := &Follow{
			UserId:   userId,
			ToUserId: toUserId,
			Behavior: FollowTypeFollowing,
		}
		toUser := &Follow{
			UserId:   toUserId,
			ToUserId: userId,
			Behavior: FollowTypeNotFollowing,
		}
		var follows []Follow
		follows, err := d.follow.TranLockByUserIdToUserId(ctx, session, userId, toUserId)
		if err != nil {
			return err
		}
		oldUserStatus := StatusStranger
		for _, v := range follows {
			if v.UserId == userId {
				if v.Behavior == FollowTypeFollowing {
					return nil
				}
				user.Id = v.Id
				oldUserStatus = v.Attribute
			} else {
				toUser.Id = v.Id
				toUser.Behavior = v.Behavior
			}
		}
		user.Attribute, toUser.Attribute = d.follow.StateMachine(oldUserStatus, FollowTypeFollowing)
		if (user.Attribute == StatusFriend || oldUserStatus == StatusFriend) && user.Attribute != oldUserStatus {
			d.follow.AddFriendKey(userId, toUserId, &keys)
		}
		// 创建或修改记录
		if len(follows) == 0 {
			err = nil
			_, err := d.follow.TranInsert(ctx, session, user, &keys)
			if err != nil {
				return err
			}
			_, err = d.follow.TranInsert(ctx, session, toUser, &keys)
			if err != nil {
				return err
			}
		} else {
			_, err := d.follow.TranUpdate(ctx, session, user, &keys)
			if err != nil {
				return err
			}
			_, err = d.follow.TranUpdate(ctx, session, toUser, &keys)
			if err != nil {
				return err
			}
		}
		// 更新记数
		err = d.userStatus.TranIncrCount(ctx, session, userId, toUserId, &keys)
		if err != nil {
			return err
		}
		res = true
		_ = d.conn.DelCacheCtx(ctx, keys...)
		return nil
	})
	if err != nil {
		return false, err
	}

	return res, nil
}

// UnFollow 取消关注对方
func (d *DBAction) UnFollow(ctx context.Context, userId, toUserId int64) (bool, error) {

	var res bool
	keys := make([]string, 0)

	err := d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		user := &Follow{
			UserId:   userId,
			ToUserId: toUserId,
			Behavior: FollowTypeNotFollowing,
		}
		toUser := &Follow{
			UserId:   toUserId,
			ToUserId: userId,
			Behavior: FollowTypeNotFollowing,
		}
		var follows []Follow
		follows, err := d.follow.TranLockByUserIdToUserId(ctx, session, userId, toUserId)
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		oldUserStatus := StatusStranger
		for _, v := range follows {
			if v.UserId == userId {
				if v.Behavior == FollowTypeNotFollowing {
					return nil
				}
				user.Id = v.Id
				oldUserStatus = v.Attribute
			} else {
				toUser.Id = v.Id
				toUser.Behavior = v.Behavior
			}
		}
		user.Attribute, toUser.Attribute = d.follow.StateMachine(oldUserStatus, FollowTypeNotFollowing)
		if (user.Attribute == StatusFriend || oldUserStatus == StatusFriend) && user.Attribute != oldUserStatus {
			d.follow.AddFriendKey(userId, toUserId, &keys)
		}
		// 修改记录

		_, err = d.follow.TranUpdate(ctx, session, user, &keys)
		if err != nil {
			return err
		}
		_, err = d.follow.TranUpdate(ctx, session, toUser, &keys)
		if err != nil {
			return err
		}

		// 更新记数
		err = d.userStatus.TranDecrCount(ctx, session, userId, toUserId, &keys)
		if err != nil {
			return err
		}
		res = true
		_ = d.conn.DelCacheCtx(ctx, keys...)
		return nil
	})
	if err != nil {
		return false, err
	}

	return res, nil
}

// IsFollow 调用favorite对数据库查询函数,查询是否点赞
// ErrNotFound错误是数据库未查询到所返回，这里需要捕获（这是正常操作，同下）
// 只有出现未知错误才会返回并打印日志，否则正常操作error = nil(同下)
func (d *DBAction) IsFollow(ctx context.Context, userId, toUserId int64) (bool, error) {
	f, err := d.follow.FindOneByUserIdToUserId(ctx, userId, toUserId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId, toUserId, err)
		return false, err
	}
	return f != nil && f.Behavior == "1", nil
}

// FollowCount 对数据库查询关注数量
func (d *DBAction) FollowCount(ctx context.Context, userId int64) (int64, error) {
	count, err := d.userStatus.FindOne(ctx, userId)
	if errors.Is(err, ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		logc.Error(ctx, userId)
		return 0, err
	}
	return count.FollowCount, nil
}

// FollowerCount 对数据库查询关注数量
func (d *DBAction) FollowerCount(ctx context.Context, userId int64) (int64, error) {
	count, err := d.userStatus.FindOne(ctx, userId)
	if errors.Is(err, ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		logc.Error(ctx, userId)
		return 0, err
	}
	return count.FollowerCount, nil
}

// FollowList 调用favorite对数据库查询用户点赞视频列表
func (d *DBAction) FollowList(ctx context.Context, userId int64) ([]int64, error) {
	res, err := d.follow.FindFollowList(ctx, userId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return nil, err
	}
	return res, nil
}

// FollowerList 调用favorite对数据库查询用户点赞视频列表
func (d *DBAction) FollowerList(ctx context.Context, userId int64) ([]int64, error) {
	res, err := d.follow.FindFollowerList(ctx, userId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return nil, err
	}
	return res, nil
}

// FriendList 调用favorite对数据库查询用户点赞视频列表
func (d *DBAction) FriendList(ctx context.Context, userId int64) ([]int64, error) {
	res, err := d.follow.FindFriendList(ctx, userId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return nil, err
	}
	return res, nil
}

func (d *DBAction) NowMessage(ctx context.Context, userId, toUserId int64) (*social.FriendUser, error) {
	v, err := d.message.FindNowMessage(ctx, userId, toUserId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return nil, err
	}
	if errors.Is(err, ErrNotFound) {
		return &social.FriendUser{
			UserId:  toUserId,
			MsgType: 0,
		}, nil
	}
	res := &social.FriendUser{
		UserId:  toUserId,
		Message: &v.Content,
	}
	if v.FromUserId != userId {
		res.MsgType = 0
	} else {
		res.MsgType = 1
	}
	return res, nil
}

func (d *DBAction) MessageList(ctx context.Context, userId, toUserId, preTime int64) ([]*social.Message, error) {
	msgList, err := d.message.FindMessageList(ctx, userId, toUserId, preTime)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return nil, err
	}
	var res []*social.Message
	if errors.Is(err, ErrNotFound) {
		return res, nil
	}
	for _, v := range msgList {
		unixMilli := v.CreateTime.UnixMilli()
		res = append(res, &social.Message{
			Id:         v.Id,
			ToUserId:   v.ToUserId,
			FromUserId: v.FromUserId,
			Content:    v.Content,
			CreateTime: unixMilli,
		})
	}
	return res, nil
}

func (d *DBAction) SendMessage(ctx context.Context, userId, toUserId int64, content string) error {
	follow, err := d.follow.FindOneByUserIdToUserId(ctx, userId, toUserId)
	if err != nil {
		return err
	}
	if follow.Attribute != StatusFriend {
		return ErrNotFriend
	}
	_, err = d.message.Insert(ctx, &Message{
		FromUserId: userId,
		ToUserId:   toUserId,
		Content:    content,
	})
	if err != nil {
		return err
	}
	_ = d.conn.DelCacheCtx(ctx, d.message.GetNowMessageCacheKey(userId, toUserId),
		d.message.GetNowMessageListCacheKey(userId, toUserId))
	return nil
}
