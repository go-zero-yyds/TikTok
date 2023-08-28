package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	DBSource  string
	Cache     cache.CacheConf
	Redisconf redis.RedisConf
	Kafka     kq.KqConf

	VideoRPC zrpc.RpcClientConf
}
