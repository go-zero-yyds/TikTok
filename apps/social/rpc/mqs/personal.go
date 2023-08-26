package mqs

import (
	"TikTok/apps/social/rpc/internal/svc"
	"context"
	"encoding/json"
	"strconv"
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

func (l *PersonalSuccess) Consume(key, val string) error {
	var message map[string][]string
	err := json.Unmarshal([]byte(val) , &message)
	if err != nil {
		return err;
	}
	for k, v := range message {
		userId, err := strconv.ParseInt(k, 10, 64)
		if err != nil && len(v) != 2 {
			continue
		} 
		actionType := v[0]
		value := v[1]
		isSuccess , err := strconv.ParseBool(value)
		if err != nil{
			continue
		}
		//special action auto follow bots
		if actionType == "register"{
			for k := range l.svcCtx.Bot.Robots{
				l.svcCtx.DBAction.FollowAction(l.ctx, userId, int64(k), "1")
				l.svcCtx.DBAction.SendMessage(l.ctx , k , userId , l.svcCtx.Bot.Robots[k].DisplayPrologue())
			}
			continue
		}
		if actionType != "avatar" && actionType != "backgroundImage" && actionType != "signature" {
			continue
		}
		if isSuccess{
			actionType += " 设置成功"
		}else{
			actionType += " 设置失败"
		}
		//这里有问题，无法动态变化
		l.svcCtx.DBAction.SendMessage(l.ctx, 0 , userId , actionType)
	}
	return nil
}
