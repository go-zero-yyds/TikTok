package logic

import (
	"context"

	"TikTok/apps/video/rpc/internal/svc"
	"TikTok/apps/video/rpc/model"
	"TikTok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPublishActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPublishActionLogic {
	return &SendPublishActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendPublishActionLogic) SendPublishAction(in *video.PublishActionReq) (*video.PublishActionResp, error) {

	videoId := l.svcCtx.Snowflake.Generate().Int64()

	videoData := &model.Video{
		VideoId:  videoId,
		UserId:   in.UserId,
		PlayUrl:  in.PlayUrl,
		CoverUrl: in.CoverUrl,
		Title:    in.Title,
	}

	_, err := l.svcCtx.Model.Insert(l.ctx, videoData)
	if err != nil {
		return &video.PublishActionResp{
			IsSucceed: false,
		}, err
	}

	return &video.PublishActionResp{
		IsSucceed: true,
	}, nil
}
