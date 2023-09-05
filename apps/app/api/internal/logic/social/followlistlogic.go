package social

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/social/rpc/social"
	"context"

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

	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)
	list, err := l.svcCtx.SocialRPC.GetRelationFollowList(l.ctx, &social.RelationFollowListReq{UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	infoList, err := GetUserInfoList(list.UserList, tokenID, l.svcCtx, l.ctx)
	if err != nil {
		return nil, err
	}
	return &types.RelationFollowListResponse{
		RespStatus: types.RespStatus(apivars.Success),
		UserList:   infoList,
	}, nil
}
