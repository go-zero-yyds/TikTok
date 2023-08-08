package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	favoriteFieldNames          = builder.RawFieldNames(&Favorite{})
	favoriteRows                = strings.Join(favoriteFieldNames, ",")
	favoriteRowsExpectAutoSet   = strings.Join(stringx.Remove(favoriteFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	favoriteRowsWithPlaceHolder = strings.Join(stringx.Remove(favoriteFieldNames, "`userId`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTiktokFavoriteUserIdPrefix  = "cache:tiktok:favorite:userId:" //这两条记录拼接作为一个key
	cacheTiktokFavoriteVideoIdsuffix = ":videoId:"

	cacheTiktokFavoriteUserId  =  "cache:tiktok:favorite:userId:" //单独为一个key 用来记录各自出现在数据库用的数量
	cacheTiktokFavoriteVideoId =  "cache:tiktok:favorite:videoId:"
)

type (
	favoriteModel interface {
		Insert(ctx context.Context, data *Favorite) (sql.Result, error)
		Delete(ctx context.Context, userId, videoId int64) error
		IsExist(ctx context.Context, userId, videoId int64) (bool, error)                //检查用户是否给视频点赞
		userORvideoCount(ctx context.Context, Id int64, userORvideo bool) (int64, error) //用户/视频点赞数量
		FindVideos(ctx context.Context, userId int64) ([]int64, error)                   //用户点赞视频id列表
	}

	defaultFavoriteModel struct {
		sqlc.CachedConn
		table string
	}

	Favorite struct { //联合主码
		UserId  int64 `db:"userId"` //用户id
		VideoId int64 `db:"videoId"`//视频id
	}
)

func newFavoriteModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultFavoriteModel {
	return &defaultFavoriteModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`favorite`",
	}
}

func (m *defaultFavoriteModel) withSession(session sqlx.Session) *defaultFavoriteModel {
	return &defaultFavoriteModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`favorite`",
	}
}

// 删除记录并删除缓存
func (m *defaultFavoriteModel) Delete(ctx context.Context, userId, videoId int64) error {
	tiktokFavoriteUserIdVIdeoIdKey := fmt.Sprintf("%s%v%s%v", cacheTiktokFavoriteUserIdPrefix, userId, cacheTiktokFavoriteVideoIdsuffix, videoId)
	tiktokFavoriteUserIdKey := fmt.Sprintf("%s%v" , cacheTiktokFavoriteUserId,userId)
	tiktokFavoriteVideoIdKey:= fmt.Sprintf("%s%v" , cacheTiktokFavoriteVideoId,videoId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `userId` = ? and videoId = ?", m.table)
		return conn.ExecCtx(ctx, query, userId, videoId)
	}, tiktokFavoriteUserIdVIdeoIdKey, tiktokFavoriteUserIdKey,tiktokFavoriteVideoIdKey)
	if err != nil {
		return err
	}
	return err
}
//插入记录，并删除缓存
func (m *defaultFavoriteModel) Insert(ctx context.Context, data *Favorite) (sql.Result, error) {
	tiktokFavoriteUserIdVIdeoIdKey := fmt.Sprintf("%s%v%s%v", cacheTiktokFavoriteUserIdPrefix, data.UserId, cacheTiktokFavoriteVideoIdsuffix, data.VideoId)
	tiktokFavoriteUserIdKey := fmt.Sprintf("%s%v" , cacheTiktokFavoriteUserId,data.UserId)
	tiktokFavoriteVideoIdKey:= fmt.Sprintf("%s%v" , cacheTiktokFavoriteVideoId,data.VideoId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, favoriteRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.VideoId)
	}, tiktokFavoriteUserIdVIdeoIdKey,tiktokFavoriteUserIdKey , tiktokFavoriteVideoIdKey)
	return ret, err
}

// 查看用户点赞视频id列表
func (m *defaultFavoriteModel) FindVideos(ctx context.Context, userId int64) ([]int64, error) {
	query := fmt.Sprintf("select videoId from %s where `userId` = ? ", m.table)
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
// 【缓存】
func (m *defaultFavoriteModel) IsExist(ctx context.Context, userId, videoId int64) (bool, error) {
	tiktokFavoriteUserIdVIdeoIdKey := fmt.Sprintf("%s%v%s%v", cacheTiktokFavoriteUserIdPrefix, userId, cacheTiktokFavoriteVideoIdsuffix, videoId)
	var resp Favorite
	err := m.QueryRowCtx(ctx, &resp, tiktokFavoriteUserIdVIdeoIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select * from  %s where userId = ? and videoId = ? ", m.table)
		return conn.QueryRowCtx(ctx, &resp, query, userId, videoId)
	})
	switch err {
	case nil:
		return true, nil
	case sqlc.ErrNotFound:
		return false, ErrNotFound
	default:
		return false, err
	}
}

// 查看各自点赞数量 true : user  false : video
// 【缓存】【这部分性能不太好】
func (m *defaultFavoriteModel) userORvideoCount(ctx context.Context, Id int64, userORvideo bool) (int64, error) {
	var obj string
	var Key string
	if userORvideo {
		obj = "userId"
		Key = fmt.Sprintf("%s%v" , cacheTiktokFavoriteUserId, Id)
	} else {
		obj = "videoId"
		Key = fmt.Sprintf("%s%v" , cacheTiktokFavoriteVideoId, Id)
	}
	var resp int64
	err := m.QueryRowCtx(ctx, &resp , Key , func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select count(*) from %s where `%s` = ?  ", m.table, obj)
		return conn.QueryRowCtx(ctx , &resp , query , Id)
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

func (m *defaultFavoriteModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTiktokFavoriteUserIdPrefix, primary)
}

func (m *defaultFavoriteModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `userId` = ? limit 1", favoriteRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultFavoriteModel) tableName() string {
	return m.table
}
