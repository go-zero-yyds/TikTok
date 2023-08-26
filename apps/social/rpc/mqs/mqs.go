package mqs

import (
	"TikTok/apps/social/rpc/internal/config"
	"TikTok/apps/social/rpc/internal/svc"
	"context"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {
    return []service.Service{
		kq.MustNewQueue(c.KqCPersonalCallbackConf , NewPersonalSuccess(ctx , svcContext)),
    }

}