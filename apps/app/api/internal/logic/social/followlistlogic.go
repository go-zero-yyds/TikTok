package social

import (
	"context"

	"rpc/apps/app/api/internal/svc"
	"rpc/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowListLogic {
	return &FollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowListLogic) FollowList(req *types.RelationFollowListRequest) (resp *types.RelationFollowListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
