package logic

import (
	"context"

	"rpc/apps/interaction/rpc/interaction"
	"rpc/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByUserIdLogic {
	return &GetFavoriteCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountByUserIdLogic) GetFavoriteCountByUserId(in *interaction.FavoriteCountByUserIdReq) (*interaction.FavoriteCountByUserIdResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.FavoriteCountByUserIdResp{}, nil
}
