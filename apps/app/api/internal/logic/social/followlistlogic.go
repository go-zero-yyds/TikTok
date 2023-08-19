package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/social/rpc/social"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

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
	if req.Token == "" {
		return &types.RelationFollowListResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
			UserList:   make([]types.User, 0),
		}, nil
	}
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.SocialRPC.GetRelationFollowList(l.ctx, &social.RelationFollowListReq{UserId: tokenID})
	if err != nil {
		return nil, err
	}
	infoList, err := GetUserInfoList(list.UserList, tokenID, l.svcCtx, l.ctx)
	if err != nil {
		return nil, err
	}
	return &types.RelationFollowListResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		UserList:   infoList,
	}, nil
}
