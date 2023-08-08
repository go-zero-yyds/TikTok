package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	commentFieldNames          = builder.RawFieldNames(&Comment{})
	commentRows                = strings.Join(commentFieldNames, ",")
	commentRowsExpectAutoSet   = strings.Join(stringx.Remove(commentFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	commentRowsWithPlaceHolder = strings.Join(stringx.Remove(commentFieldNames, "`commentId`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheTiktokCommentCommentIdPrefix = "cache:tiktok:comment:commentId:"
)

type (
	commentModel interface {
		Insert(ctx context.Context, data *Comment) (sql.Result, error)
		FindOne(ctx context.Context, commentId int64) (*Comment, error)
		Delete(ctx context.Context, commentId int64) error
		Count(ctx context.Context, videoId int64) (int64, error)
		CommentList(ctx context.Context, videoId int64) ([]*Comment, error)
	}

	defaultCommentModel struct {
		sqlc.CachedConn
		table string
	}

	Comment struct {
		CommentId  int64     `db:"commentId"` //雪花算法生成id
		UserId     int64     `db:"userId"`    //用户id
		VideoId    int64     `db:"videoId"`	  //视频id
		CreateDate time.Time `db:"createDate"`//数据库获取的当前创建时间
		Content    string    `db:"content"`   //评论内容
	}
)

func newCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultCommentModel {
	return &defaultCommentModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`comment`",
	}
}

func (m *defaultCommentModel) withSession(session sqlx.Session) *defaultCommentModel {
	return &defaultCommentModel{
		CachedConn: m.CachedConn.WithSession(session),
		table:      "`comment`",
	}
}

func (m *defaultCommentModel) Delete(ctx context.Context, commentId int64) error {
	tiktokCommentCommentIdKey := fmt.Sprintf("%s%v", cacheTiktokCommentCommentIdPrefix, commentId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `commentId` = ?", m.table)
		return conn.ExecCtx(ctx, query, commentId)
	}, tiktokCommentCommentIdKey)
	return err
}

// 查找成功: 查询结果，nli
// 没查到  : 初始值,ErrNoFound
// 其他    : nil , err
// 【缓存】
func (m *defaultCommentModel) FindOne(ctx context.Context, commentId int64) (*Comment, error) {
	tiktokCommentCommentIdKey := fmt.Sprintf("%s%v", cacheTiktokCommentCommentIdPrefix, commentId)
	var resp Comment
	err := m.QueryRowCtx(ctx, &resp, tiktokCommentCommentIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `commentId` = ? limit 1", commentRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, commentId)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return &resp, ErrNotFound
	default:
		return nil, err
	}
}

// 数据库自动获取时间，插入不需要指定时间
func (m *defaultCommentModel) Insert(ctx context.Context, data *Comment) (sql.Result, error) {
	tiktokCommentCommentIdKey := fmt.Sprintf("%s%v", cacheTiktokCommentCommentIdPrefix, data.CommentId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, now() , ?)", m.table, commentRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.CommentId, data.UserId, data.VideoId, data.Content)
	}, tiktokCommentCommentIdKey)
	return ret, err
}

func (m *defaultCommentModel) Count(ctx context.Context, videoId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `videoId` = ? ", m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, videoId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultCommentModel) CommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	query := fmt.Sprintf("select * from %s where `videoId` = ? order by createDate desc ", m.table)
	resp := make([]*Comment, 0)
	err := m.QueryRowsNoCacheCtx(context.Background(), &resp, query, videoId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultCommentModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheTiktokCommentCommentIdPrefix, primary)
}

func (m *defaultCommentModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `commentId` = ? limit 1", commentRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultCommentModel) tableName() string {
	return m.table
}
