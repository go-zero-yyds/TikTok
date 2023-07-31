package logic

import (
	"context"

	"rpc/apps/user/rpc/internal/svc"
	"rpc/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAvatarLogic {
	return &SetAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetAvatarLogic) SetAvatar(in *user.SetAvatarReq) (*user.SetAvatarResp, error) {
	// todo: add your logic here and delete this line

	return &user.SetAvatarResp{}, nil
}
