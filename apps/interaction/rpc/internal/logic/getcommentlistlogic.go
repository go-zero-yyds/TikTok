package logic

import (
	"context"

	"rpc/apps/interaction/rpc/interaction"
	"rpc/apps/interaction/rpc/internal/svc"

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

func (l *GetCommentListLogic) GetCommentList(in *interaction.CommentListReq) (*interaction.CommentListResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.CommentListResp{}, nil
}
