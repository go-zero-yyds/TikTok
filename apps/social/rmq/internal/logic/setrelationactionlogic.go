package logic

import (
	"TikTok/apps/social/rmq/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"encoding/json"
	"strconv"

	"github.com/zeromicro/go-zero/core/threading"
)

type PersonalSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPersonalSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *PersonalSuccess {
	return &PersonalSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PersonalSuccess) Consume(_, val string) error {
	var message map[string][]string
	err := json.Unmarshal([]byte(val), &message)
	if err != nil {
		return err
	}
	for k, v := range message {
		userId, err := strconv.ParseInt(k, 10, 64)
		if err != nil && len(v) != 2 {
			continue
		}
		actionType := v[0]
		if actionType != "register" {
			continue
		}
		isSuccess, err := strconv.ParseBool(v[1])
		if err != nil || !isSuccess {
			continue
		}
		threading.GoSafeCtx(l.ctx, func() {
			// auto follow bots
			for botId := range l.svcCtx.Bot.Robots {
				_, _ = l.svcCtx.SocialRPC.SendRelationAction(l.ctx, &social.RelationActionReq{
					UserId:     userId,
					ToUserId:   botId,
					ActionType: 1,
				})
				_, _ = l.svcCtx.SocialRPC.SendRelationAction(l.ctx, &social.RelationActionReq{
					UserId:     botId,
					ToUserId:   userId,
					ActionType: 1,
				})
				_, _ = l.svcCtx.SocialRPC.SendMessageAction(l.ctx, &social.MessageActionReq{
					UserId:   botId,
					ToUserId: userId,
					Content:  l.svcCtx.Bot.Robots[botId].DisplayPrologue(),
				})
			}
		})
	}
	return nil
}
