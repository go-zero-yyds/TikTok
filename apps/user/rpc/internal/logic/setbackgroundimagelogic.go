package logic

import (
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetBackgroundImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetBackgroundImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetBackgroundImageLogic {
	return &SetBackgroundImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetBackgroundImageLogic) SetBackgroundImage(in *user.SetBackgroundImageReq) (*user.SetBackgroundImageResp, error) {
	isSucceed := true
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId:          in.UserId,
		BackgroundImage: in.Url,
	}, "backgroundImage")

	if err != nil {
		isSucceed = false
	}

	return &user.SetBackgroundImageResp{
		IsSucceed: isSucceed,
	}, err

}
