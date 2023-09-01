package logic

import (
	"context"
	"encoding/json"
	"strconv"

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
	
	if in.ToUserId < l.svcCtx.Config.RobotMaxId { 
		message := make(map[string][]string)
		message[strconv.FormatInt(in.UserId, 10)] = []string{strconv.FormatInt(in.ToUserId, 10), in.Content}
		data, err := json.Marshal(message)
		if err == nil {
			//推送消息
			err = l.svcCtx.KqPusherClient.Push(string(data))
		}
		if err != nil {
			return &social.MessageActionResp{
				IsSucceed: false,
			}, nil
		}
	}
	return &social.MessageActionResp{
		IsSucceed: true,
	}, nil
}
