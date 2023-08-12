package logic

import (
	"context"
	"fmt"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCommentList 查看视频的所有评论，按发布时间倒序
func (l *GetCommentListLogic) GetCommentList(in *interaction.CommentListReq) (*interaction.CommentListResp, error) {
	commentList, err := l.svcCtx.DBAction.CommentList(l.ctx, in.VideoId)
	if err != nil {
		return nil, err
	}
	ret := make([]*interaction.Comment, 0, len(commentList))
	for _, v := range commentList {
		ret = append(ret, &interaction.Comment{
			Id:         v.CommentId,
			UserId:     v.UserId,
			Content:    v.Content,
			CreateDate: fmt.Sprintf("%v", v.CreateDate.Unix()),
		})
	}
	return &interaction.CommentListResp{
		CommentList: ret,
	}, nil
}
