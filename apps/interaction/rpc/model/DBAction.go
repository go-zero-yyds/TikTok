package model

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DBAction struct {
	favorite   FavoriteModel
	comment    CommentModel
	userLikes  UserLikesModel
	videoStats VideoStatsModel
	conn       sqlc.CachedConn
}

// NewDBAction 初始化数据库信息
func NewDBAction(conn sqlx.SqlConn, c cache.ClusterConf) *DBAction {
	ret := &DBAction{
		favorite:   NewFavoriteModel(conn, c), //创建点赞表的接口
		comment:    NewCommentModel(conn, c),  //创建评论表的接口
		userLikes:  NewUserLikesModel(conn, c),
		videoStats: NewVideoStatsModel(conn, c),
		conn:       sqlc.NewConn(conn, c),
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
	return f != nil && f.Behavior == FavoriteTypeFollowing, nil
}

// FavoriteCountByUserId 查询，用户 点赞数量
func (d *DBAction) FavoriteCountByUserId(ctx context.Context, userId int64) (int64, error) {
	count, err := d.userLikes.FindOne(ctx, userId)
	if errors.Is(err, ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		logc.Error(ctx, userId)
		return 0, err
	}
	return count.LikeCount, nil
}

// FavoriteCountByVideoId 调用favorite对数据库查询，视频 点赞数量
func (d *DBAction) FavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	count, err := d.videoStats.FindOne(ctx, videoId)
	if errors.Is(err, ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		logc.Error(ctx, videoId)
		return 0, err
	}
	return count.LikeCount, nil
}

// CommentCountByVideoId 调用comment对数据库查询，视频 评论数量
func (d *DBAction) CommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	count, err := d.videoStats.FindOne(ctx, videoId)
	if errors.Is(err, ErrNotFound) {
		return 0, nil
	}
	if err != nil {
		logc.Error(ctx, videoId)
		return 0, err
	}
	return count.CommentCount, nil
}

// FavoriteAction 调用favorite对数据库查询用户是否点赞过
// 并对数据库actionType操作
// 如果取消操作，则更新Behavior 后续等待特定时间再删除记录
// 只有出现未知错误是false否则都是true
func (d *DBAction) FavoriteAction(ctx context.Context, userId, videoId int64, actionType string) (bool, error) {
	data := &Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	var res bool
	keys := make([]string, 0)
	err := d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var (
			result sql.Result
			err    error
		)
		if actionType == FavoriteTypeFollowing {
			data.Behavior = FavoriteTypeFollowing
			result, err = d.favorite.TranInsertOrUpdate(ctx, session, data, &keys)
		} else {
			data.Behavior = FavoriteTypeNotFollowing
			result, err = d.favorite.TranEmptyOrUpdate(ctx, session, data, &keys)
		}
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		res = affected != int64(0)
		if !res {
			return nil
		}
		if actionType == FavoriteTypeFollowing {
			err := d.videoStats.TranIncrLikeCount(ctx, session, videoId, &keys)
			if err != nil {
				return err
			}
			err = d.userLikes.TranIncrCount(ctx, session, userId, &keys)
			if err != nil {
				return err
			}
		} else {
			err := d.videoStats.TranDecrLikeCount(ctx, session, videoId, &keys)
			if err != nil {
				return err
			}
			err = d.userLikes.TranDecrCount(ctx, session, userId, &keys)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	if res {
		_ = d.conn.DelCacheCtx(ctx, keys...)
	}
	return res, nil
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
func (d *DBAction) CommentAction(ctx context.Context, c *Comment, actionType int32) error {
	var err error
	keys := make([]string, 0)
	if actionType == 1 {
		err = d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
			_, err := d.comment.TranInsert(ctx, session, c, &keys)
			if err != nil {
				return err
			}
			err = d.videoStats.TranIncrCommentCount(ctx, session, c.VideoId, &keys)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	} else if actionType == 2 {

		v, err := d.comment.FindOne(ctx, c.CommentId)

		if err != nil || v.UserId != c.UserId {
			return ErrIllegalArgument
		}
		if err != nil {
			return err
		}
		err = d.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
			err := d.comment.TranUpdateDel(ctx, session, c, &keys)
			if err != nil {
				return err
			}
			err = d.videoStats.TranDecrCommentCount(ctx, session, c.VideoId, &keys)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	if err != nil {
		logc.Error(ctx, c.CommentId, actionType)
		return err
	}
	_ = d.conn.DelCacheCtx(ctx, keys...)
	return nil
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
