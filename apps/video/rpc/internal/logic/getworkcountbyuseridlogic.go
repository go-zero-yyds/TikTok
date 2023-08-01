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
	videoCout, err := l.svcCtx.Model.CountByUserId(context.Background(), in.UserId)
	if err != nil {
		return nil, err
	}
	return &video.WorkCountByUserIdResp{
		WorkCount: videoCout,
	}, nil
}
