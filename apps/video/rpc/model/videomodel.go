package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoModel = (*customVideoModel)(nil)

type (
	// VideoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoModel.
	VideoModel interface {
		videoModel
		CountByUserId(ctx context.Context, userId int64) (int64, error)
		VideoListByUserId(ctx context.Context, userId int64) ([]*Video, error)
	}

	customVideoModel struct {
		*defaultVideoModel
	}
)

// NewVideoModel returns a model for the database table.
func NewVideoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) VideoModel {
	return &customVideoModel{
		defaultVideoModel: newVideoModel(conn, c, opts...),
	}
}

func (m *customVideoModel) CountByUserId(ctx context.Context, userId int64) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE user_id = ?", m.table)
	err := m.QueryRowNoCache(&count, query, userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *customVideoModel) VideoListByUserId(ctx context.Context, userId int64) ([]*Video, error) {
	var videoList []*Video
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = ?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &videoList, query, userId)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}
