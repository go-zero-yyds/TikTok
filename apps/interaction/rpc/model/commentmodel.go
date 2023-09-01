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
	_                               CommentModel = (*customCommentModel)(nil)
	cacheCommentCommentIdListPrefix              = "cache:comment:commentId:list"
)

type (
	// CommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCommentModel.
	CommentModel interface {
		commentModel
		CommentList(ctx context.Context, videoId int64) ([]*Comment, error)
		TranInsert(ctx context.Context, s sqlx.Session, data *Comment, keys *[]string) (sql.Result, error)
		TranUpdateDel(ctx context.Context, s sqlx.Session, data *Comment, keys *[]string) error
	}

	customCommentModel struct {
		*defaultCommentModel
	}
)

// NewCommentModel returns a model for the database table.
func NewCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CommentModel {
	return &customCommentModel{
		defaultCommentModel: newCommentModel(conn, c, opts...),
	}
}

// CommentList 查看视频评论列表, 因为没分页，限制1000条。
func (m *defaultCommentModel) CommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	resp := make([]*Comment, 0)
	key := fmt.Sprintf("%s%v", cacheCommentCommentIdListPrefix, videoId)
	if err := m.CachedConn.GetCache(key, &resp); err == nil {
		return resp, nil
	}

	query := fmt.Sprintf("select * from %s where `video_id` = ? and is_deleted = '0' order by create_time desc limit 1000", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, videoId)
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

func (m *defaultCommentModel) TranInsert(ctx context.Context, s sqlx.Session, data *Comment, keys *[]string) (sql.Result, error) {
	commentCommentIdKey := fmt.Sprintf("%s%v", cacheCommentCommentIdPrefix, data.CommentId)
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, commentRowsExpectAutoSet)
	ret, err := s.ExecCtx(ctx, query, data.CommentId, data.UserId, data.VideoId, data.Content, data.IpAddress, data.Location, "0")
	if err != nil {
		return nil, err
	}
	*keys = append(*keys, commentCommentIdKey)
	m.OnChangeDeleteCache(ctx, data.VideoId)
	return ret, err
}

func (m *defaultCommentModel) TranUpdateDel(ctx context.Context, s sqlx.Session, data *Comment, keys *[]string) error {

	commentCommentIdKey := fmt.Sprintf("%s%v", cacheCommentCommentIdPrefix, data.CommentId)
	query := fmt.Sprintf("update %s set is_deleted = '1' where `comment_id` = ?", m.table)
	_, err := s.ExecCtx(ctx, query, data.CommentId)
	*keys = append(*keys, commentCommentIdKey)
	m.OnChangeDeleteCache(ctx, data.VideoId)
	return err
}

func (m *defaultCommentModel) OnChangeDeleteCache(ctx context.Context, videoId int64) {
	deleteKeys := []string{
		fmt.Sprintf("%s%v", cacheCommentCommentIdListPrefix, videoId),
	}
	m.CachedConn.DelCacheCtx(ctx, deleteKeys...)
}
