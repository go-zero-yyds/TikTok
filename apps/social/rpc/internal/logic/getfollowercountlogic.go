package logic

import (
	"context"

	"rpc/apps/social/rpc/internal/svc"
	"rpc/apps/social/rpc/social"

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
	// todo: add your logic here and delete this line

	return &social.FollowerCountResp{}, nil
}
