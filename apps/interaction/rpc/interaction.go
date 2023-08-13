package main

import (
	"context"
	"flag"
	"fmt"

	"TikTok/apps/interaction/rpc/cron"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/config"
	"TikTok/apps/interaction/rpc/internal/server"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/interaction.yaml", "the config file")

func main() {
	var cfg logx.LogConf
	_ = conf.FillDefault(&cfg)
	cfg.Mode = "file"
	logc.MustSetup(cfg)
	defer logc.Close()
	
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx , err := svc.NewServiceContext(c)
	if err != nil{
		logc.Error(context.Background() , "loading svc error")
		return 
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		interaction.RegisterInteractionServer(grpcServer, server.NewInteractionServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	//启动定时任务
	cron.InitCron(context.Background() , ctx , c.CleanTime)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
