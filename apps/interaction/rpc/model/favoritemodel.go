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

var (
	_ FavoriteModel = (*customFavoriteModel)(nil)

	cacheFavoriteCountUserIdPrefix  = "cache:favorite:count:userId:"
	cacheFavoriteCountVideoIdPrefix = "cache:favorite:count:videoId:"
)

type (
	// FavoriteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFavoriteModel.
	FavoriteModel interface {
		favoriteModel
		UserOrVideoCount(ctx context.Context, Id int64, userORvideo bool) (int64, error) //用户/视频点赞数量
		FindVideos(ctx context.Context, userId int64) ([]int64, error)                   //用户点赞视频id列表
		FlushAndClean(ctx context.Context) error                                         //特定时间重置数据，删除没用数据
		InsertOrUpdate(ctx context.Context, data *Favorite) (sql.Result, error)
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

// FindVideos 查看用户点赞视频id列表
func (m *defaultFavoriteModel) FindVideos(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("select videoId from %s where `userId` = ? and behavior = '1'  ", m.table)
	var resp []int64
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// UserOrVideoCount 查看各自点赞数量 true : user  false : video
// 【缓存】【这部分性能不太好】
func (m *defaultFavoriteModel) UserOrVideoCount(ctx context.Context, Id int64, userOrVideo bool) (int64, error) {
	var obj string
	var Key string
	if userOrVideo {
		obj = "userId"
		Key = fmt.Sprintf("%s%v", cacheFavoriteCountUserIdPrefix, Id)
	} else {
		obj = "videoId"
		Key = fmt.Sprintf("%s%v", cacheFavoriteCountVideoIdPrefix, Id)
	}
	var resp int64
	err := m.QueryRowCtx(ctx, &resp, Key, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select count(*) from %s where `%s` = ? and behavior = '1' ", m.table, obj)
		return conn.QueryRowCtx(ctx, &resp, query, Id)
	})
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

// FlushAndClean 删除数据库中所有behavior为2的值，减少冗余
func (m *defaultFavoriteModel) FlushAndClean(ctx context.Context) error {
	//这里不删除缓存中数据
	query := fmt.Sprintf("delete from %s where behavior = '2' ", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query)
	return err
}

func (m *defaultFavoriteModel) InsertOrUpdate(ctx context.Context, data *Favorite) (sql.Result, error) {
	favoriteCountUserIdKey := fmt.Sprintf("%s%v", cacheFavoriteCountUserIdPrefix, data.UserId)
	favoriteCountVideoIdKey := fmt.Sprintf("%s%v", cacheFavoriteCountVideoIdPrefix, data.VideoId)
	favoriteFavoriteIdKey := fmt.Sprintf("%s%v", cacheFavoriteFavoriteIdPrefix, data.FavoriteId)
	favoriteUserIdVideoIdKey := fmt.Sprintf("%s%v:%v", cacheFavoriteUserIdVideoIdPrefix, data.UserId, data.VideoId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf(`
			INSERT INTO %s (%s)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE behavior = ?;
		`, m.table, favoriteRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.VideoId, data.Behavior, data.Behavior)
	}, favoriteFavoriteIdKey, favoriteUserIdVideoIdKey, favoriteCountUserIdKey, favoriteCountVideoIdKey)
	return ret, err
}
