package logic

import (
	"TikTok/apps/social/rmq/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"encoding/json"
	"strconv"

	"github.com/zeromicro/go-zero/core/threading"
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

// Consume userId  toUserId content
func (l *RobotsResponse) Consume(_, val string) error {
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
		threading.GoSafeCtx(l.ctx, func() {
			action, data, err := l.svcCtx.Bot.ProcessIfMessageForRobot(l.ctx, userId, touserId, v[1], l.svcCtx.KqPusherClient, l.svcCtx.FS)
			if err == nil && action {
				if data != "" { // 机器人回发消息
					_, _ = l.svcCtx.SocialRPC.SendMessageAction(l.ctx, &social.MessageActionReq{
						UserId:     touserId,
						ToUserId:   userId,
						Content:    data,
						ActionType: 1,
					})
				}
			}
		})
	}
	return nil
}
