package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	_ FavoriteModel = (*customFavoriteModel)(nil)

	cacheTiktokFavoriteUserIdPrefix  = "cache:tiktok:favorite:userId:" //这两条记录拼接作为一个key
	cacheTiktokFavoriteVideoIdsuffix = ":videoId:"

	cacheTiktokFavoriteUserId  = "cache:tiktok:favorite:userId:" //单独为一个key 用来记录各自出现在数据库用的数量
	cacheTiktokFavoriteVideoId = "cache:tiktok:favorite:videoId:"
)

type (
	// FavoriteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFavoriteModel.
	FavoriteModel interface {
		favoriteModel
		FindOneById(ctx context.Context, userId, videoId int64) (*Favorite, error)        //检查用户是否给视频点赞
		userORvideoCount(ctx context.Context, Id int64, userORvideo bool) (int64, error) //用户/视频点赞数量
		FindVideos(ctx context.Context, userId int64) ([]int64, error)                   //用户点赞视频id列表
		IndirectUpdate(ctx context.Context, data *Favorite) (error)						 //管理增加缓存删除
		IndirectInsert(ctx context.Context, data *Favorite) (sql.Result, error)			 //管理更新缓存删除
	}

	customFavoriteModel struct {
		*defaultFavoriteModel
	}
)

// NewFavoriteModel returns a model for the database table.
func NewFavoriteModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FavoriteModel {
	return &customFavoriteModel{
		defaultFavoriteModel: newFavoriteModel(conn, c, opts...),
	}
}

// 查看用户点赞视频id列表
func (m *defaultFavoriteModel) FindVideos(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("select videoId from %s where `userId` = ? and behavior = '1'  ", m.table)
	var resp []int64
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// 检查是否存在user和video之间的联系
// 该函数总会返回对象
// 【缓存】
func (m *defaultFavoriteModel) FindOneById(ctx context.Context, userId, videoId int64) (*Favorite, error) {
	tiktokFavoriteUserIdVIdeoIdKey := fmt.Sprintf("%s%v%s%v", cacheTiktokFavoriteUserIdPrefix, userId, cacheTiktokFavoriteVideoIdsuffix, videoId)
	var resp Favorite
	err := m.QueryRowCtx(ctx, &resp, tiktokFavoriteUserIdVIdeoIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select * from  %s where userId = ? and videoId = ?  ", m.table)
		return conn.QueryRowCtx(ctx, &resp, query, userId, videoId)
	})
	return &resp, err
}

// 查看各自点赞数量 true : user  false : video
// 【缓存】【这部分性能不太好】
func (m *defaultFavoriteModel) userORvideoCount(ctx context.Context, Id int64, userORvideo bool) (int64, error) {
	var obj string
	var Key string
	if userORvideo {
		obj = "userId"
		Key = fmt.Sprintf("%s%v", cacheTiktokFavoriteUserId, Id)
	} else {
		obj = "videoId"
		Key = fmt.Sprintf("%s%v", cacheTiktokFavoriteVideoId, Id)
	}
	var resp int64
	err := m.QueryRowCtx(ctx, &resp, Key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select count(*) from %s where `%s` = ? and behavior = '1' ", m.table, obj)
		return conn.QueryRowCtx(ctx, &resp, query, Id)
	})
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

// 管理缓存删除
func (m *defaultFavoriteModel) IndirectInsert(ctx context.Context, data *Favorite) (sql.Result, error) {
	tiktokFavoriteUserIdVIdeoIdKey := fmt.Sprintf("%s%v%s%v", cacheTiktokFavoriteUserIdPrefix, data.UserId, cacheTiktokFavoriteVideoIdsuffix, data.VideoId)
	tiktokFavoriteUserIdKey := fmt.Sprintf("%s%v", cacheTiktokFavoriteUserId, data.UserId)
	tiktokFavoriteVideoIdKey := fmt.Sprintf("%s%v", cacheTiktokFavoriteVideoId, data.VideoId)
	//删除缓存
	m.DelCacheCtx(ctx ,tiktokFavoriteUserIdVIdeoIdKey ,  tiktokFavoriteUserIdKey , tiktokFavoriteVideoIdKey)
	return m.Insert(ctx , data)
}

func (m *defaultFavoriteModel) IndirectUpdate(ctx context.Context, data *Favorite) (error) {
	tiktokFavoriteUserIdVIdeoIdKey := fmt.Sprintf("%s%v%s%v", cacheTiktokFavoriteUserIdPrefix, data.UserId, cacheTiktokFavoriteVideoIdsuffix, data.VideoId)
	tiktokFavoriteUserIdKey := fmt.Sprintf("%s%v", cacheTiktokFavoriteUserId, data.UserId)
	tiktokFavoriteVideoIdKey := fmt.Sprintf("%s%v", cacheTiktokFavoriteVideoId, data.VideoId)
	//删除缓存
	m.DelCacheCtx(ctx ,tiktokFavoriteUserIdVIdeoIdKey ,  tiktokFavoriteUserIdKey , tiktokFavoriteVideoIdKey)
	return m.Update(ctx , data)
}
