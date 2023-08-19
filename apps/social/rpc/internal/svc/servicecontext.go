package svc

import (
	"TikTok/apps/social/rpc/internal/config"
	"TikTok/apps/social/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	DBAction *model.DBAction
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		DBAction: model.NewDBAction(sqlx.NewMysql(c.DBSource), c.Cache),
	}
}
