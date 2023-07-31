package logic

import (
	"context"

	"TikTok/apps/video/rpc/internal/svc"
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
	// todo: add your logic here and delete this line

	return &video.PublishActionResp{}, nil
}
