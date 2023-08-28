package main

import (
	"flag"

	"TikTok/apps/video/rmq/internal/config"
	"TikTok/apps/video/rmq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/video.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	q := kq.MustNewQueue(c.Kafka, kq.WithHandle(svc.NewService(c).Consume))
	defer q.Stop()
	q.Start()
}
