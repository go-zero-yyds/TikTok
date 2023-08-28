package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"TikTok/apps/video/rpc/internal/config"
	"TikTok/apps/video/rpc/internal/logic"
	"TikTok/apps/video/rpc/internal/svc"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
)

var configFile = flag.String("f", "apps/video/rpc/etc/video.yaml", "the config file")

// 定时任务
func main() {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx, cancel := context.WithCancel(context.Background())
	svcCtx, err := svc.NewServiceContext(c)
	if err != nil {
		logc.Errorf(context.Background(), "init service context error ,  message: %v", err)
		return
	}
	hotVideo := logic.NewHotVideoLogic(ctx, svcCtx)
	tasks := cron.New() // 创建了cron实例，用于处理定时任务
	// 添加了定时任务函数，根据配置中的JobSpec规定的时间间隔执行。
	// 在定时任务函数中，获取当前时间ctime，然后调用hotVideo.ScoreCalculation(ctime)进行热门视频计算。
	_, err = tasks.AddFunc(svcCtx.Config.HotVideoConf.JobSpec, func() {
		// 热门视频计算
		// 获取当前时间
		ctime := time.Now()
		hotVideo.ScoreCalculation(ctime)
	})
	if err != nil {
		logc.Errorf(context.TODO(), "Add task error ,  message: %v", err)
	}
	tasks.Start()
	logc.Infof(context.Background(), "cron start success")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signo := <-sig
	logc.Infof(context.Background(), "got signal %#v, exit", signo)

	tasks.Stop()
	cancel()
	// wait for all tasks to finish or just kill it?
	logc.Infof(context.Background(), "shutting down")
}
