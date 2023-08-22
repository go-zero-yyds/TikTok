package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoStatsModel = (*customVideoStatsModel)(nil)

type (
	// VideoStatsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoStatsModel.
	VideoStatsModel interface {
		videoStatsModel
		TranIncrLikeCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error
		TranDecrLikeCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error
		TranIncrCommentCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error
		TranDecrCommentCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error
	}

	customVideoStatsModel struct {
		*defaultVideoStatsModel
	}
)

// NewVideoStatsModel returns a model for the database table.
func NewVideoStatsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) VideoStatsModel {
	return &customVideoStatsModel{
		defaultVideoStatsModel: newVideoStatsModel(conn, c, opts...),
	}
}

func (m *customVideoStatsModel) TranIncrLikeCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error {
	videoStatsVideoIdKey := fmt.Sprintf("%s%v", cacheVideoStatsVideoIdPrefix, videoId)
	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 1, 0)
		ON DUPLICATE KEY UPDATE like_count = like_count + 1;`, m.table, videoStatsRowsExpectAutoSet)
	_, err := s.ExecCtx(ctx, query, videoId)
	*keys = append(*keys, videoStatsVideoIdKey)
	return err
}

func (m *customVideoStatsModel) TranDecrLikeCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error {
	videoStatsVideoIdKey := fmt.Sprintf("%s%v", cacheVideoStatsVideoIdPrefix, videoId)
	query := fmt.Sprintf(`
		UPDATE %s
		SET like_count = like_count - 1
		WHERE video_id = ?;`, m.table)
	_, err := s.ExecCtx(ctx, query, videoId)
	*keys = append(*keys, videoStatsVideoIdKey)
	return err
}

func (m *customVideoStatsModel) TranIncrCommentCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error {
	videoStatsVideoIdKey := fmt.Sprintf("%s%v", cacheVideoStatsVideoIdPrefix, videoId)
	query := fmt.Sprintf(`
		INSERT INTO %s (%s)
		VALUES (?, 0, 1)
		ON DUPLICATE KEY UPDATE comment_count = comment_count + 1;`, m.table, videoStatsRowsExpectAutoSet)
	_, err := s.ExecCtx(ctx, query, videoId)
	*keys = append(*keys, videoStatsVideoIdKey)
	return err
}

func (m *customVideoStatsModel) TranDecrCommentCount(ctx context.Context, s sqlx.Session, videoId int64, keys *[]string) error {
	videoStatsVideoIdKey := fmt.Sprintf("%s%v", cacheVideoStatsVideoIdPrefix, videoId)
	query := fmt.Sprintf(`
		UPDATE %s
		SET comment_count = comment_count - 1
		WHERE video_id = ?;`, m.table)
	_, err := s.ExecCtx(ctx, query, videoId)
	*keys = append(*keys, videoStatsVideoIdKey)
	return err
}
