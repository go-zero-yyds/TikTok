package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DBSource           string
	Cache              cache.CacheConf
	SnowflakeNode      int
	SnowflakeStartTime int64
	CleanTime          string
}
