package svc

import (
	"TikTok/apps/social/rpc/internal/config"
	"TikTok/apps/social/rpc/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	DBAction *model.DBAction
	KqPusherClient *kq.Pusher

}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		DBAction: model.NewDBAction(sqlx.NewMysql(c.DBSource), c.Cache),
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
