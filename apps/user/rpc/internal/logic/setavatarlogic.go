package logic

import (
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"context"

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
	isSucceed := true
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId: in.UserId,
		Avatar: in.Url,
	}, "avatar")

	if err != nil {
		isSucceed = false
	}

	return &user.SetAvatarResp{
		IsSucceed: isSucceed,
	}, err
}
