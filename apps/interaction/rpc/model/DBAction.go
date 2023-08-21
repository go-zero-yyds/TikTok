package model

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DBAction struct {
	favorite FavoriteModel
	comment  CommentModel
}

// NewDBAction 初始化数据库信息
func NewDBAction(conn sqlx.SqlConn, c cache.ClusterConf) *DBAction {
	ret := &DBAction{
		favorite: NewFavoriteModel(conn, c), //创建点赞表的接口
		comment:  NewCommentModel(conn, c),  //创建评论表的接口
	}
	return ret
}

// IsFavorite 调用favorite对数据库查询函数,查询是否点赞
// ErrNotFound错误是数据库未查询到所返回，这里需要捕获（这是正常操作，同下）
// 只有出现未知错误才会返回并打印日志，否则正常操作error = nil(同下)
func (d *DBAction) IsFavorite(ctx context.Context, userId, videoId int64) (bool, error) {
	f, err := d.favorite.FindOneByUserIdVideoId(ctx, userId, videoId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId, videoId, err)
		return false, err
	}
	return f != nil && f.Behavior == "1", nil
}

// FavoriteCountByUserId 调用favorite对数据库查询，用户 点赞数量
func (d *DBAction) FavoriteCountByUserId(ctx context.Context, userId int64) (int64, error) {
	count, err := d.favorite.UserOrVideoCount(ctx, userId, true)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return 0, err
	}
	return count, nil
}

// FavoriteCountByVideoId 调用favorite对数据库查询，视频 点赞数量
func (d *DBAction) FavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	count, err := d.favorite.UserOrVideoCount(ctx, videoId, false)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, videoId)
		return 0, err
	}
	return count, nil
}

// CommentCountByVideoId 调用comment对数据库查询，视频 评论数量
func (d *DBAction) CommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	count, err := d.comment.Count(ctx, videoId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, videoId)
		return 0, err
	}
	return count, nil
}

// FavoriteAction 调用favorite对数据库查询用户是否点赞过
// 并对数据库actionType操作
// 如果取消操作，则更新Behavior 后续等待特定时间再删除记录
// 只有出现未知错误是false否则都是true
func (d *DBAction) FavoriteAction(ctx context.Context, userId, videoId int64, actionType string) (bool, error) {
	result, err := d.favorite.InsertOrUpdate(ctx, &Favorite{
		UserId:   userId,
		VideoId:  videoId,
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

// FavoriteList 调用favorite对数据库查询用户点赞视频列表
func (d *DBAction) FavoriteList(ctx context.Context, userId int64) ([]int64, error) {
	ret, err := d.favorite.FindVideos(ctx, userId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, userId)
		return nil, err
	}
	return ret, nil
}

// CommentAction 执行评论/取消操作
// 成功返回comment结构体，（评论成功 查询结果，取消成功 初始值）
// 如果用户可选参数没有赋值，将会返回地址错误
func (d *DBAction) CommentAction(ctx context.Context, userId, videoId int64, actionType int32, commentText *string, commentId *int64) (*Comment, error) {
	var ret *Comment
	var err error
	if actionType == 1 {
		ret = &Comment{
			CommentId:  *commentId,
			UserId:     userId,
			VideoId:    videoId,
			Content:    *commentText,
			CreateDate: time.Now(),
		}
		_, err = d.comment.Insert(ctx, ret)
	} else if actionType == 2 {
		err = d.comment.Delete(ctx, *commentId)
	}

	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, commentId, actionType)
		return nil, err
	}

	return ret, nil
}

func (d *DBAction) CommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	ret, err := d.comment.CommentList(ctx, videoId)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logc.Error(ctx, err)
		return nil, err
	}
	return ret, nil
}

// CleanUnusedFavorite 定时清楚favorite中无用数据
func (d *DBAction) CleanUnusedFavorite(ctx context.Context) error {
	err := d.favorite.FlushAndClean(ctx)
	if err != nil {
		logc.Error(ctx)
	}
	return err
}
