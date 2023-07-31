package logic

import (
	"context"

	"rpc/apps/video/rpc/internal/svc"
	"rpc/apps/video/rpc/video"

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

func (l *DetailLogic) Detail(in *video.BasicVideoInfoReq) (*video.BasicVideoInfoResp, error) {
	// todo: add your logic here and delete this line

	return &video.BasicVideoInfoResp{}, nil
}
