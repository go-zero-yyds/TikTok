package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFollowerListLogic {
	return &GetRelationFollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationFollowerListLogic) GetRelationFollowerList(in *social.RelationFollowerListReq) (*social.RelationFollowerListResp, error) {
	// todo: add your logic here and delete this line

	return &social.RelationFollowerListResp{}, nil
}
