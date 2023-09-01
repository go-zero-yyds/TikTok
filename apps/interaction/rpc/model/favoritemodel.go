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
	_                               FavoriteModel = (*customFavoriteModel)(nil)
	cacheFavoriteUserIdVideosPrefix               = "cache:userLikes:userId:videos:"
)

// 定义关注类型常量
const (
	FavoriteTypeNotFollowing = "0"
	FavoriteTypeFollowing    = "1"
)

type (
	// FavoriteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFavoriteModel.
	FavoriteModel interface {
		favoriteModel
		FindVideos(ctx context.Context, userId int64) ([]int64, error)
		FlushAndClean(ctx context.Context) error
		TranInsertOrUpdate(ctx context.Context, s sqlx.Session, data *Favorite, keys *[]string) (sql.Result, error)
		TranEmptyOrUpdate(ctx context.Context, s sqlx.Session, newData *Favorite, keys *[]string) (result sql.Result, err error)
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

// FindVideos 查看用户点赞视频id列表, 因为没分页，限制1000条。
func (m *defaultFavoriteModel) FindVideos(ctx context.Context, userId int64) ([]int64, error) {
	var resp []int64
	key := fmt.Sprintf("%s%v", cacheFavoriteUserIdVideosPrefix, userId)
	if err := m.CachedConn.GetCache(key, &resp); err == nil {
		return resp, nil
	}

	query := fmt.Sprintf("select video_id from %s where `user_id` = ? and behavior = '%s' limit 1000", m.table, FavoriteTypeFollowing)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		m.CachedConn.SetCacheCtx(ctx, key, &resp)
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FlushAndClean 删除数据库中所有behavior为 FavoriteTypeNotFollowing 的值，减少冗余
func (m *defaultFavoriteModel) FlushAndClean(ctx context.Context) error {
	//这里不删除缓存中数据
	query := fmt.Sprintf("delete from %s where behavior = '%s' LIMIT 5000", m.table, FavoriteTypeNotFollowing)
	_, err := m.ExecNoCacheCtx(ctx, query)
	return err
}

// TranInsertOrUpdate 插入一条关注记录或者更新关注记录
func (m *customFavoriteModel) TranInsertOrUpdate(ctx context.Context, s sqlx.Session, data *Favorite, keys *[]string) (sql.Result, error) {
	favoriteFavoriteIdKey := fmt.Sprintf("%s%v", cacheFavoriteFavoriteIdPrefix, data.FavoriteId)
	favoriteUserIdVideoIdKey := fmt.Sprintf("%s%v:%v", cacheFavoriteUserIdVideoIdPrefix, data.UserId, data.VideoId)
	query := fmt.Sprintf(`
			INSERT INTO %s (%s)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE behavior = ?;
		`, m.table, favoriteRowsExpectAutoSet)
	ret, err := s.ExecCtx(ctx, query, data.UserId, data.VideoId, data.Behavior, data.Behavior)
	*keys = append(*keys, favoriteFavoriteIdKey, favoriteUserIdVideoIdKey)
	m.OnChangeDeleteCache(ctx, data.UserId)
	return ret, err
}

// TranEmptyOrUpdate 更新关注记录没有则不操作
func (m *defaultFavoriteModel) TranEmptyOrUpdate(ctx context.Context, s sqlx.Session, newData *Favorite, keys *[]string) (result sql.Result, err error) {
	favoriteFavoriteIdKey := fmt.Sprintf("%s%v", cacheFavoriteFavoriteIdPrefix, newData.FavoriteId)
	favoriteUserIdVideoIdKey := fmt.Sprintf("%s%v:%v", cacheFavoriteUserIdVideoIdPrefix, newData.UserId, newData.VideoId)
	query := fmt.Sprintf("update %s set %s where `user_id` = ? and video_id = ?", m.table, favoriteRowsWithPlaceHolder)
	ret, err := s.ExecCtx(ctx, query, newData.UserId, newData.VideoId, newData.Behavior, newData.UserId, newData.VideoId)
	*keys = append(*keys, favoriteFavoriteIdKey, favoriteUserIdVideoIdKey)
	m.OnChangeDeleteCache(ctx, newData.UserId)
	return ret, err
}

func (m *defaultFavoriteModel) OnChangeDeleteCache(ctx context.Context, userId int64) {
	deleteKeys := []string{
		fmt.Sprintf("%s%v", cacheFavoriteUserIdVideosPrefix, userId),
	}
	m.CachedConn.DelCacheCtx(ctx, deleteKeys...)
}
