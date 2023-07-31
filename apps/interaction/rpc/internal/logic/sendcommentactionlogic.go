package logic

import (
	"context"

	"rpc/apps/interaction/rpc/interaction"
	"rpc/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCommentActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCommentActionLogic {
	return &SendCommentActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCommentActionLogic) SendCommentAction(in *interaction.CommentActionReq) (*interaction.CommentActionResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.CommentActionResp{}, nil
}
