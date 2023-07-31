package logic

import (
	"context"

	"rpc/apps/interaction/rpc/interaction"
	"rpc/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFavoriteActionLogic {
	return &SendFavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendFavoriteActionLogic) SendFavoriteAction(in *interaction.FavoriteActionReq) (*interaction.FavoriteActionResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.FavoriteActionResp{}, nil
}
