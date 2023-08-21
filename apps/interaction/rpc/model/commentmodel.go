package model

import (
	"context"
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
		Count(ctx context.Context, videoId int64) (int64, error)
		CommentList(ctx context.Context, videoId int64) ([]*Comment, error)
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

func (m *defaultCommentModel) Count(ctx context.Context, videoId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `videoId` = ? ", m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, videoId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultCommentModel) CommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	query := fmt.Sprintf("select * from %s where `videoId` = ? order by createDate desc ", m.table)
	resp := make([]*Comment, 0)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, videoId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
