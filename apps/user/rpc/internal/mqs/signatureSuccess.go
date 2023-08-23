package mqs

import (
	"TikTok/apps/user/rpc/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type SignatureSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignatureSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *SignatureSuccess {
	return &SignatureSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignatureSuccess) Consume(key, val string) error {
	logx.Infof("SignatureSuccess key :%s , val :%s", key, val)
	return nil
}
