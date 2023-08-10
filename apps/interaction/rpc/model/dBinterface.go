package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DBaction struct {
	favorite FavoriteModel
	comment  CommentModel
}

// 初始化数据库信息
func NewDBaction(conn sqlx.SqlConn, c cache.ClusterConf) *DBaction {
	ret := &DBaction{
		favorite: NewFavoriteModel(conn, c), //创建点赞表的接口
		comment:  NewCommentModel(conn, c),  //创建评论表的接口
	}
	return ret
}

// 调用favorite对数据库查询函数,查询是否点赞
// ErrNotFound错误是数据库未查询到所返回，这里需要捕获（这是正常操作，同下）
// 只有出现未知错误才会返回并打印日志，否则正常操作error = nil(同下)
func (d *DBaction) IsFavorite(ctx context.Context, userId, videoId int64) (bool, error) {
	f, err := d.favorite.FindOneById(ctx, userId, videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, userId, videoId, err)
		return false, err
	}
	return f.Behavior == "1", nil
}

// 调用favorite对数据库查询，用户 点赞数量
func (d *DBaction) FavoriteCountByUserId(ctx context.Context, userId int64) (int64, error) {
	count, err := d.favorite.userORvideoCount(ctx, userId, true)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, userId)
		return 0, err
	}
	return count, nil
}

// 调用favorite对数据库查询，视频 点赞数量
func (d *DBaction) FavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	count, err := d.favorite.userORvideoCount(ctx, videoId, false)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, videoId)
		return 0, err
	}
	return count, nil
}

// 调用comment对数据库查询，视频 评论数量
func (d *DBaction) CommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	count, err := d.comment.Count(ctx, videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, videoId)
		return 0, err
	}
	return count, nil
}

// 调用favorite对数据库查询用户是否点赞过
// 并对数据库actionType操作
// 如果取消操作，则更新Behavior 后续等待特定时间再删除记录
// 只有出现未知错误是false否则都是true
func (d *DBaction) FavoriteAction(ctx context.Context, userId, videoId int64, actionType int32) (bool, error) {
	f, err := d.favorite.FindOneById(ctx, userId, videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, userId, videoId, actionType)
		return false, err
	}
	//如果点赞过，并且用户操作也是点赞 直接返回
	if f.Behavior == "1" && actionType == 1 {
		return true, nil
	}
	//如果没点赞，并且用户操作也是取消点赞直接返回
	if f.Behavior == "2" && actionType == 2 {
		return true, nil
	}
	//剩下的都是相反操作，if 点赞 ->取消  else 点赞
	if f.Behavior == "1" {
		f.Behavior = "2"
		err = d.favorite.IndirectUpdate(ctx, f)
	}else if f.Behavior == "2"{//有可能不在表中
		f.Behavior = "1"
		err = d.favorite.IndirectUpdate(ctx, f)
	}else {//不在表中
		_, err = d.favorite.IndirectInsert(ctx, &Favorite{
			UserId:  userId,
			VideoId: videoId,
			Behavior: "1",
		})
	}
	if err != nil && err != ErrNotFound{
		logc.Error(ctx, err)
	}
	return true ,  nil
}

// 调用favorite对数据库查询用户点赞视频列表
func (d *DBaction) FavoriteList(ctx context.Context, userId int64) ([]int64, error) {
	ret, err := d.favorite.FindVideos(ctx, userId)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, userId)
		return nil, err
	}
	return ret, nil
}

// 执行评论/取消操作
// 成功返回comment结构体，（评论成功 查询结果，取消成功 初始值）
// 如果用户可选参数没有赋值，将会返回地址错误
func (d *DBaction) CommentAction(ctx context.Context, userId, videoId int64, actionType int32, commentText *string, commentId *int64) (*Comment, error) {
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

	if err != nil && err != ErrNotFound {
		logc.Error(ctx, commentId, actionType)
		return nil, err
	}

	return ret, nil
}

func (d *DBaction) CommentList(ctx context.Context, userId, videoId int64) ([]*Comment, error) {
	ret, err := d.comment.CommentList(ctx, videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(ctx, err)
		return nil, err
	}
	return ret, nil
}
//定时清楚favorite中无用数据
func (d *DBaction) CleanUnusedFavorite(ctx context.Context )(error){
	err := d.favorite.FlushAndClean(ctx)
	if err != nil{
		logc.Error(ctx)
	}
	return err;
}