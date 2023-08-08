package model

import (
	"context"
	"fmt"

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
		favorite: newFavoriteModel(conn, c), //创建点赞表的接口
		comment:  newCommentModel(conn, c),  //创建评论表的接口
	}
	return ret
}

// 调用favorite对数据库查询函数,查询是否点赞
// ErrNotFound错误是数据库未查询到所返回，这里需要捕获（这是正常操作，同下）
// 只有出现未知错误才会返回并打印日志，否则正常操作error = nil(同下)
func (d *DBaction) IsFavorite(userId, videoId int64) (bool, error) {
	exist, err := d.favorite.IsExist(context.Background(), userId, videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), userId, videoId, err)
		return false, err
	}
	return exist, nil
}

// 调用favorite对数据库查询，用户 点赞数量
func (d *DBaction) FavoriteCountByUserId(userId int64) (int64, error) {
	count, err := d.favorite.userORvideoCount(context.Background(), userId, true)
	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), userId)
		return 0, err
	}
	return count, nil
}

// 调用favorite对数据库查询，视频 点赞数量
func (d *DBaction) FavoriteCountByVideoId(videoId int64) (int64, error) {
	count, err := d.favorite.userORvideoCount(context.Background(), videoId, false)
	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), videoId)
		return 0, err
	}
	return count, nil
}

// 调用comment对数据库查询，视频 评论数量
func (d *DBaction) CommentCountByVideoId(videoId int64) (int64, error) {
	count, err := d.comment.Count(context.Background(), videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), videoId)
		return 0, err
	}
	return count, nil
}

// 调用favorite对数据库查询用户是否点赞过
// 并对数据库actionType操作
// 只有出现未知错误是false否则都是true
func (d *DBaction) FavoriteAction(userId, videoId int64, actionType int32) (bool, error) {
	exist, err := d.favorite.IsExist(context.Background(), userId, videoId)
	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), userId, videoId, actionType)
		return false, err
	}
	//如果点赞过，并且用户操作也是点赞 直接返回
	if exist && actionType == 1 {
		return true, nil
	}
	//如果没点赞，并且用户操作也是取消点赞直接返回
	if !exist && actionType == 2 {
		return true, nil
	}
	//剩下的都是相反操作，if 点赞 ->取消  else 点赞
	if exist {
		err = d.favorite.Delete(context.Background(), userId, videoId)
	} else {
		_, err = d.favorite.Insert(context.Background(), &Favorite{
			UserId:  userId,
			VideoId: videoId,
		})
	}
	if err != nil {
		logc.Error(context.Background(), err)
	}
	return err == nil, nil
}

// 调用favorite对数据库查询用户点赞视频列表
func (d *DBaction) FavoriteList(userId int64) ([]int64, error) {
	ret, err := d.favorite.FindVideos(context.Background(), userId)
	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), userId)
		return nil, err
	}
	return ret, nil
}

// 执行评论/取消操作
// 成功返回comment结构体，（评论成功 查询结果，取消成功 初始值）
// 如果用户可选参数没有赋值，将会返回地址错误
func (d *DBaction) CommentAction(userId, videoId int64, actionType int32, commentText *string, commentId *int64) (*Comment, error) {
	var ret *Comment
	var err error
	if actionType == 1 {
		_, err = d.comment.Insert(context.Background(), &Comment{
			CommentId: *commentId,
			UserId:    userId,
			VideoId:   videoId,
			Content:   *commentText,
		})
	} else if actionType == 2 {
		err = d.comment.Delete(context.Background(), *commentId)
	}

	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), commentId, actionType)
		return nil, err
	}
	ret, err = d.comment.FindOne(context.Background(), *commentId)

	if err != nil && err != ErrNotFound {
		logc.Error(context.Background(), commentId, actionType)
		return nil, err
	}
	return ret, nil
}

func (d *DBaction) CommentList(userId, videoId int64) ([]*Comment, error) {
	ret, err := d.comment.CommentList(context.Background(), videoId)
	if err != nil && err != ErrNotFound {
		fmt.Println("error ", err)
		return nil, err
	}
	return ret, nil
}
