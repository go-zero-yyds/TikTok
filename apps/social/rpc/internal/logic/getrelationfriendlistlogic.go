package logic

import (
	"context"

	"rpc/apps/social/rpc/internal/svc"
	"rpc/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFriendListLogic {
	return &GetRelationFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationFriendListLogic) GetRelationFriendList(in *social.RelationFriendListReq) (*social.RelationFriendListResp, error) {
	// todo: add your logic here and delete this line

	return &social.RelationFriendListResp{}, nil
}
