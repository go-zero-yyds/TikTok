package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSourse string
	Cache      cache.CacheConf
	Snow       SnowConf
}

type SnowConf struct {
	StartTime string
	Node      int64
}
