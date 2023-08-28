package svc

import (
	"TikTok/apps/video/rpc/internal/config"
	"TikTok/apps/video/rpc/model"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	Model       model.VideoModel
	Snowflake   *snowflake.Node
	RedisClient *redis.Client
	KafkaPusher *kq.Pusher
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	snowflake.Epoch = c.Snowflake.StartTime
	node, err := snowflake.NewNode(c.Snowflake.Node)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:    c,
		Model:     model.NewVideoModel(sqlx.NewMysql(c.DBSource), c.Cache),
		Snowflake: node,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     c.Redisconf.Host,
			Password: c.Redisconf.Pass,
		}),
		KafkaPusher: kq.NewPusher(c.Kafka.Brokers, c.Kafka.Topic),
	}, nil
}
