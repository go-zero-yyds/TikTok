package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *interaction.CommentCountByVideoIdReq) (*interaction.CommentCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.CommentCountByVideoIdResp{}, nil
}
