package svc

import (
	"rpc/apps/video/model"
	"rpc/apps/video/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Model  model.VideoModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Model:  model.NewVideoModel(sqlx.NewMysql(c.DataSourse), c.Cache),
	}
}
