package svc

import (
	"TikTok/apps/social/rpc/internal/config"
	"TikTok/apps/social/rpc/model"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel model.MessageModel
	FollowModel  model.FollowModel
	SocialModel  model.SocialModel
	CustomDB     model.CustomDB
	Redis        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	//MySQL配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.MySQL.Account, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.Database, c.MySQL.Options)

	//Redis配置
	redisConf := redis.RedisConf{
		Host:     c.RedisConfig.Host,
		Type:     c.RedisConfig.Type,
		Pass:     c.RedisConfig.Pass,
		Tls:      c.RedisConfig.Tls,
		NonBlock: c.RedisConfig.NonBlock,
	}
	rds := redis.MustNewRedis(redisConf)

	return &ServiceContext{
		Config:       c,
		MessageModel: model.NewMessageModel(sqlx.NewMysql(dsn)),
		FollowModel:  model.NewFollowModel(sqlx.NewMysql(dsn)),
		SocialModel:  model.NewSocialModel(sqlx.NewMysql(dsn)),
		CustomDB:     *model.NewCustomDB(sqlx.NewMysql(dsn)),
		Redis:        rds,
	}
}
