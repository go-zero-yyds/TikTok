package svc

import (
	"TikTok/apps/interaction/rpc/internal/config"
	"TikTok/apps/interaction/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	Snowflake *snowflake.Node
	DBAction  *model.DBAction
	Rds       *redis.Redis
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	snowflake.Epoch = c.SnowflakeStartTime
	node, err := snowflake.NewNode(int64(c.SnowflakeNode))
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:    c,
		Snowflake: node,
		DBAction:  model.NewDBAction(sqlx.NewMysql(c.DBSource), c.Cache),
		Rds:       redis.MustNewRedis(c.Redis.RedisConf),
	}, nil
}
