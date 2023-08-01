package logic

import (
	"context"

	"TikTok/apps/video/rpc/internal/svc"
	"TikTok/apps/video/rpc/video"

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
	videoId := in.VideoId
	videoInfo, err := l.svcCtx.Model.FindOne(l.ctx, videoId)
	if err != nil {
		return nil, err
	}

	return &video.BasicVideoInfoResp{
		Video: &video.BasicVideoInfo{
			Id:       videoInfo.VideoId,
			UserId:   videoInfo.UserId,
			PlayUrl:  videoInfo.PlayUrl,
			CoverUrl: videoInfo.CoverUrl,
			Title:    videoInfo.Title,
		},
	}, nil
}
