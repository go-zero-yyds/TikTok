package model

import (
	"TikTok/apps/social/rpc/social"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type DBAction struct {
	follow  FollowModel
	message MessageModel
}

// NewDBAction 初始化数据库信息
func NewDBAction(conn sqlx.SqlConn, c cache.ClusterConf) *DBAction {
	ret := &DBAction{
		follow:  NewFollowModel(conn, c, cache.WithExpiry(time.Minute)),
		message: NewMessageModel(conn, c, cache.WithExpiry(time.Minute)),
	}
	return ret
}

// FollowAction 调用 follow 对数据库查询用户是否点赞过
// 并对数据库actionType操作
// 出现未知错误和没有修改是false否则都是true
func (d *DBAction) FollowAction(ctx context.Context, userId, toUserId int64, actionType string) (bool, error) {
	result, err := d.follow.InsertOrUpdate(ctx, &Follow{
		UserId:   userId,
		ToUserId: toUserId,
		Behavior: actionType,
	})
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected != int64(0), nil
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
	count, err := d.follow.FindByFollowCount(ctx, userId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId, err)
		return 0, err
	}
	return count, nil
}

// FollowerCount 对数据库查询关注数量
func (d *DBAction) FollowerCount(ctx context.Context, userId int64) (int64, error) {
	count, err := d.follow.FindByFollowerCount(ctx, userId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId, err)
		return 0, err
	}
	return count, nil
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
	_, err := d.message.Insert(ctx, &Message{
		FromUserId: userId,
		ToUserId:   toUserId,
		Content:    content,
	})
	if err != nil {
		return err
	}
	return nil
}
