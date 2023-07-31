package logic

import (
	"context"

	"rpc/apps/social/rpc/internal/svc"
	"rpc/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *social.MessageChatReq) (*social.MessageChatResp, error) {
	// todo: add your logic here and delete this line

	return &social.MessageChatResp{}, nil
}
