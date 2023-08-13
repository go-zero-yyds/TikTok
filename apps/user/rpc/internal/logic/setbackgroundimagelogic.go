package logic

import (
	"context"

	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/user"

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

func (l *SetBackgroundImageLogic) SetBackgroundImage(in *user.BackgroundImageReq) (*user.BackgroundImageResp, error) {
	// todo: add your logic here and delete this line

	return &user.BackgroundImageResp{}, nil
}
