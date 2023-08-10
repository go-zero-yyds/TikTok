package cron

import (
	"TikTok/apps/interaction/rpc/internal/svc"
	"context"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logc"
)

func InitCron(ctx context.Context , l *svc.ServiceContext, spec string) {
	c := cron.New()

	// 添加一个每天凌晨3点执行的任务
	_, err := c.AddFunc(spec , func() {
		// 这里执行你的删除操作
		err := l.DBact.CleanUnusedFavorite(ctx)
		logc.Info(ctx , "clean favorite data ")
		if err != nil{
			logc.Error(ctx, err)
		}
	})
	if err != nil {
		logc.Errorf(context.TODO(), "Add task error ,  message: %v", err)
	}

	c.Start()
}