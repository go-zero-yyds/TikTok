package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowCountLogic {
	return &GetFollowCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowCountLogic) GetFollowCount(in *social.FollowCountReq) (*social.FollowCountResp, error) {

	count, err := l.svcCtx.DBAction.FollowCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &social.FollowCountResp{
		FollowCount: count,
	}, nil
}
