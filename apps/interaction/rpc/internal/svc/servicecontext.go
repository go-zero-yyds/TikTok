package svc

import (
	"TikTok/apps/interaction/rpc/internal/config"
	"TikTok/apps/interaction/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	DBAction *model.DBAction
	Rds      *redis.Redis
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	r := redis.MustNewRedis(c.Redis.RedisConf)
	return &ServiceContext{
		Config:   c,
		DBAction: model.NewDBAction(r, sqlx.NewMysql(c.DBSource), c.Cache),
		Rds:      r,
	}, nil
}
