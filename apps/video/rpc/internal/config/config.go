package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DBSource  string
	Cache     cache.CacheConf
	Snowflake SnowflakeConf

	Redisconf    redis.RedisConf
	HotVideoConf struct {
		Windowsize   int
		HotScore     int
		ExpireMinute int
		JobSpec      string
	}
	Kafka struct {
		Brokers []string
		Topic   string
	}
}

type SnowflakeConf struct {
	StartTime int64
	Node      int64
}
