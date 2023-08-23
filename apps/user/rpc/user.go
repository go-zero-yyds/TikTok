package main

import (
	"TikTok/apps/user/rpc/internal/mqs"
	"context"
	"flag"
	"fmt"

	"TikTok/apps/user/rpc/internal/config"
	"TikTok/apps/user/rpc/internal/server"
	"TikTok/apps/user/rpc/internal/svc"
	"TikTok/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "apps/user/rpc/etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx, err := svc.NewServiceContext(c)
	ctx := context.Background()

	if err != nil {
		logc.Error(context.Background(), "loading svc error")
		return
	}

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(svcCtx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range mqs.Consumers(c, ctx, svcCtx) {
		serviceGroup.Add(mq)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
	serviceGroup.Start()
}
