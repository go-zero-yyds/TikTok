package logic

import (
	"context"

	"rpc/apps/social/rpc/internal/svc"
	"rpc/apps/social/rpc/social"

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
	// todo: add your logic here and delete this line

	return &social.IsFollowResp{}, nil
}
