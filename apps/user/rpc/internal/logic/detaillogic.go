package logic

import (
	"context"

	"rpc/apps/user/rpc/internal/svc"
	"rpc/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *user.BasicUserInfoReq) (*user.BasicUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &user.BasicUserInfoResp{}, nil
}
