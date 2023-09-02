package main

import (
	"TikTok/apps/social/rmq/internal/config"
	"TikTok/apps/social/rmq/internal/consumer"
	"TikTok/apps/social/rmq/internal/svc"
	"context"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/robot.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)


	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range consumer.Consumers(c, context.Background(), ctx) {
		serviceGroup.Add(mq)
		fmt.Println(mq)
	}
	fmt.Printf("Starting bots server ...\n")
	serviceGroup.Start()

}