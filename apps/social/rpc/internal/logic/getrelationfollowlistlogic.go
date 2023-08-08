package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFollowListLogic {
	return &GetRelationFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationFollowListLogic) GetRelationFollowList(in *social.RelationFollowListReq) (*social.RelationFollowListResp, error) {
	// todo: add your logic here and delete this line

	return &social.RelationFollowListResp{}, nil
}
