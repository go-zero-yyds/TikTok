package logic

import (
	"context"
	"strconv"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendRelationActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRelationActionLogic {
	return &SendRelationActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendRelationActionLogic) SendRelationAction(in *social.RelationActionReq) (*social.RelationActionResp, error) {
	res, err := l.svcCtx.DBAction.FollowAction(l.ctx, in.UserId, in.ToUserId, strconv.Itoa(int(in.ActionType)))
	if err != nil {
		return nil, err
	}
	return &social.RelationActionResp{
		IsSucceed: res,
	}, nil
}
