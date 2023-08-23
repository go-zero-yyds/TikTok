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

var _ CommentModel = (*customCommentModel)(nil)

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

func (m *defaultCommentModel) CommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	query := fmt.Sprintf("select * from %s where `video_id` = ? and is_deleted = '0' order by create_date desc ", m.table)
	resp := make([]*Comment, 0)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, videoId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCommentModel) TranInsert(ctx context.Context, s sqlx.Session, data *Comment, keys *[]string) (sql.Result, error) {
	commentCommentIdKey := fmt.Sprintf("%s%v", cacheCommentCommentIdPrefix, data.CommentId)
	commentCommentIdUserIdKey := fmt.Sprintf("%s%v:%v", cacheCommentCommentIdUserIdPrefix, data.CommentId, data.UserId)
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, commentRowsExpectAutoSet)
	ret, err := s.ExecCtx(ctx, query, data.UserId, data.VideoId, data.CreateDate, data.Content, "0", data.IpAddress, data.Location)
	if err != nil {
		return nil, err
	}
	*keys = append(*keys, commentCommentIdKey, commentCommentIdUserIdKey)
	return ret, err
}

func (m *defaultCommentModel) TranUpdateDel(ctx context.Context, s sqlx.Session, data *Comment, keys *[]string) error {

	commentCommentIdKey := fmt.Sprintf("%s%v", cacheCommentCommentIdPrefix, data.CommentId)
	commentCommentIdUserIdKey := fmt.Sprintf("%s%v:%v", cacheCommentCommentIdUserIdPrefix, data.CommentId, data.UserId)
	query := fmt.Sprintf("update %s set is_deleted = '1' where `comment_id` = ?", m.table)
	_, err := s.ExecCtx(ctx, query, data.CommentId)
	*keys = append(*keys, commentCommentIdKey, commentCommentIdUserIdKey)
	return err
}
