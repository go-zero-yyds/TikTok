package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByVideoIdLogic {
	return &GetFavoriteCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFavoriteCountByVideoIdLogic) GetFavoriteCountByVideoId(in *interaction.FavoriteCountByVideoIdReq) (*interaction.FavoriteCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line

	return &interaction.FavoriteCountByVideoIdResp{}, nil
}
