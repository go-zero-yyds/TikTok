package logic

import (
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetSignatureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSignatureLogic {
	return &SetSignatureLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetSignatureLogic) SetSignature(in *user.SetSignatureReq) (*user.SetSignatureResp, error) {
	isSucceed := true
	err := l.svcCtx.UserModel.UpdateByUserId(l.ctx, &model.User{
		UserId:    in.UserId,
		Signature: in.Content,
	}, "signature")

	if err != nil {
		isSucceed = false
	}

	return &user.SetSignatureResp{
		IsSucceed: isSucceed,
	}, err

}
