package consumer

import (
	"TikTok/apps/social/rmq/internal/config"
	"TikTok/apps/social/rmq/internal/logic"
	"TikTok/apps/social/rmq/internal/svc"
	"context"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {
    return []service.Service{
		kq.MustNewQueue(c.MessageForRobots , logic.NewRobotsResponse(ctx , svcContext)),
		kq.MustNewQueue(c.PersonalCallback , logic.NewPersonalSuccess(ctx, svcContext)),
    }

}