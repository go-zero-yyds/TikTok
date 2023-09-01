package logic

import (
	"TikTok/apps/social/rmq/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"encoding/json"
	"strconv"
)

type RobotsResponse struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRobotsResponse(ctx context.Context, svcCtx *svc.ServiceContext) *RobotsResponse {
	return &RobotsResponse{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// userid  touserid content
func (l *RobotsResponse) Consume(key, val string) error {
	var message map[string][]string
	err := json.Unmarshal([]byte(val), &message)
	if err != nil {
		return err
	}
	for k, v := range message {
		if len(v) != 2 {
			continue
		}
		userId, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		touserId, err := strconv.ParseInt(v[0], 10, 64)
		if err != nil {
			continue
		}
		go func(userId  , touserId int64 , v []string) {
			action, data, err := l.svcCtx.Bot.ProcessIfMessageForRobot(l.ctx, userId, touserId, v[1], l.svcCtx.KqPusherClient, l.svcCtx.FS)
			if err == nil && action {
				if data != "" { // 机器人回发消息
					l.svcCtx.SocialRPC.SendMessageAction(l.ctx, &social.MessageActionReq{
						UserId:     touserId,
						ToUserId:   userId,
						Content:    data,
						ActionType: 1,
					})
				}
			}
		}(userId , touserId , v)
	}
	return nil
}
