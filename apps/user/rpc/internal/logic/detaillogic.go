package logic

import (
	"context"

	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
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
	res, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "用户不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	return &user.BasicUserInfoResp{
		User: res,
	}, nil
}
