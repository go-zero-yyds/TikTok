package svc

import (
	"TikTok/apps/social/rpc/internal/config"
	"TikTok/apps/social/rpc/model"
	"TikTok/apps/social/rpc/bot"
	"TikTok/pkg/FileSystem"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	DBAction *model.DBAction
	KqPusherClient *kq.Pusher
	Bot *Robot.BossRobot
	FS  FileSystem.FileSystem
}

func NewServiceContext(c config.Config) *ServiceContext {
	var fs FileSystem.FileSystem
	if c.FS.AwsS3.Endpoint != "" {
		fs = FileSystem.NewS3(c.FS.AwsS3.Endpoint, c.FS.AwsS3.Bucket, c.FS.Prefix, c.FS.AwsS3.AccessKeyID, c.FS.AwsS3.AccessKeySecret)
	} else {
		fs = FileSystem.New(c.FS.Webdav.URL, c.FS.Webdav.User, c.FS.Webdav.Password, c.FS.Prefix, c.FS.Webdav.DownloadLinkPrefix)
	}
	return &ServiceContext{
		Config:   c,
		DBAction: model.NewDBAction(sqlx.NewMysql(c.DBSource), c.Cache),
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		Bot: Robot.NewBossRobot(),
		FS: fs,
	}
}
