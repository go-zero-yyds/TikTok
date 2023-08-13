package logic

import (
	"context"
	"errors"

	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"

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
	res, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, model.UserNotFound
		}
		return nil, err
	}

	return &user.BasicUserInfoResp{
		User: &user.BasicUserInfo{
			Id:              res.UserId,
			Name:            res.Username,
			Avatar:          &res.Avatar,
			BackgroundImage: &res.BackgroundImage,
			Signature:       res.Signature,
		},
	}, nil
}
