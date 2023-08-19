package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageActionLogic {
	return &SendMessageActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageActionLogic) SendMessageAction(in *social.MessageActionReq) (*social.MessageActionResp, error) {
	err := l.svcCtx.DBAction.SendMessage(l.ctx, in.UserId, in.ToUserId, in.Content)
	if err != nil {
		return nil, err
	}
	return &social.MessageActionResp{
		IsSucceed: true,
	}, nil
}
