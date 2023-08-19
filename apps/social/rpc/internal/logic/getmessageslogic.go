package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

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

	list, err := l.svcCtx.DBAction.MessageList(l.ctx, in.UserId, in.ToUserId, in.PreMsgTime)
	if err != nil {
		return nil, err
	}
	return &social.MessageChatResp{
		MessageList: list,
	}, nil
}
