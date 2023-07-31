package logic

import (
	"context"

	"rpc/apps/social/rpc/internal/svc"
	"rpc/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendRelationActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRelationActionLogic {
	return &SendRelationActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendRelationActionLogic) SendRelationAction(in *social.RelationActionReq) (*social.RelationActionResp, error) {
	// todo: add your logic here and delete this line

	return &social.RelationActionResp{}, nil
}
