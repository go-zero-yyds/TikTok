package svc

import (
	"TikTok/apps/social/rmq/internal/config"
	"TikTok/apps/social/rmq/robot"
	"TikTok/apps/social/rpc/socialclient"
	"TikTok/pkg/FileSystem"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	KqPusherClient *kq.Pusher
	Bot 		   *robot.BossRobot
	FS  		    FileSystem.FileSystem
	SocialRPC       socialclient.Social	

}

func NewServiceContext(c config.Config) *ServiceContext {
	var fs FileSystem.FileSystem
	if c.FS.AwsS3.Endpoint != "" {
		fs = FileSystem.NewS3(c.FS.AwsS3.Endpoint, c.FS.AwsS3.Bucket, c.FS.Prefix, c.FS.AwsS3.AccessKeyID, c.FS.AwsS3.AccessKeySecret)
	} else {
		fs = FileSystem.New(c.FS.Webdav.URL, c.FS.Webdav.User, c.FS.Webdav.Password, c.FS.Prefix, c.FS.Webdav.DownloadLinkPrefix)
	}
	pusher :=  kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic)
	return &ServiceContext{
		KqPusherClient:pusher,
		Bot: robot.NewBossRobot(pusher),
		FS: fs,
		SocialRPC:socialclient.NewSocial(zrpc.MustNewClient(c.SocialRPC)),
	}
}
