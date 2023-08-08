package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteLogic {
	return &IsFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFavoriteLogic) IsFavorite(in *interaction.IsFavoriteReq) (*interaction.IsFavoriteResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.IsFavoriteResp{}, nil
}
