package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowLogic {
	return &IsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowLogic) IsFollow(in *social.IsFollowReq) (*social.IsFollowResp, error) {
	isFollow, err := l.svcCtx.DBAction.IsFollow(l.ctx, in.UserId, in.ToUserId)
	if err != nil {
		return nil, err
	}
	return &social.IsFollowResp{
		IsFollow: isFollow,
	}, nil
}
