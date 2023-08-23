package mqs

import (
	"TikTok/apps/user/rpc/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type BackgroundImageSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBackgroundImageSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *BackgroundImageSuccess {
	return &BackgroundImageSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BackgroundImageSuccess) Consume(key, val string) error {
	logx.Infof("BackgroundImageSuccess key :%s , val :%s", key, val)
	return nil
}
