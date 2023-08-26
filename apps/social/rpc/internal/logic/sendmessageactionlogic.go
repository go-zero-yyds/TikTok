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
	//返回true代表有bot执行了操作，否则未执行
	action , data, err := l.svcCtx.Bot.ProcessIfMessageForRobot(l.ctx, in.UserId, in.ToUserId, in.Content,l.svcCtx.KqPusherClient , l.svcCtx.FS)
	if err != nil {
		return nil ,err
	}
	if action{
		if data != "" { // 机器人回发消息
			l.svcCtx.DBAction.SendMessage(l.ctx , in.ToUserId , in.UserId , data)
		}
	}
	return &social.MessageActionResp{
		IsSucceed: true,
	}, nil
}
