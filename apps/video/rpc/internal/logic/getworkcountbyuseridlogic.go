package logic

import (
	"context"

	"TikTok/apps/video/rpc/internal/svc"
	"TikTok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkCountByUserIdLogic {
	return &GetWorkCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkCountByUserIdLogic) GetWorkCountByUserId(in *video.WorkCountByUserIdReq) (*video.WorkCountByUserIdResp, error) {
	// todo: add your logic here and delete this line

	return &video.WorkCountByUserIdResp{}, nil
}
