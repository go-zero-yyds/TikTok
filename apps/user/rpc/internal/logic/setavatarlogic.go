package logic

import (
	"TikTok/apps/user/rpc/model"
	"context"

	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/user"

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
	errUpdate := l.svcCtx.UserModel.Update(l.ctx, &model.User{
		UserId: in.UserId,
		//Avatar: in.Avatar,
	})

	if errUpdate != nil {
		return nil, errUpdate
	}

	return &user.SetAvatarResp{
		IsSucceed: true,
	}, nil
}
