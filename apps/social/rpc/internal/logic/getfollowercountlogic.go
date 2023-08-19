package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerCountLogic {
	return &GetFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowerCountLogic) GetFollowerCount(in *social.FollowerCountReq) (*social.FollowerCountResp, error) {

	count, err := l.svcCtx.DBAction.FollowerCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &social.FollowerCountResp{
		FollowerCount: count,
	}, nil
}
