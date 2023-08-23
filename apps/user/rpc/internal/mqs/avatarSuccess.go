package mqs

import (
	"TikTok/apps/user/rpc/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type AvatarSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAvatarSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *AvatarSuccess {
	return &AvatarSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AvatarSuccess) Consume(key, val string) error {
	logx.Infof("AvatarSuccess key :%s , val :%s", key, val)
	return nil
}
