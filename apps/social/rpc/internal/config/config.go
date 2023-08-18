package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MySQL       MySQLConfig
	RedisConfig redis.RedisConf
}

type MySQLConfig struct {
	Database string `json:",default=tiktok_social"`
	Account  string `json:",default=root"`
	Password string `json:",default=123456"`
	Host     string `json:",default=127.0.0.1"`
	Port     int16  `json:",default=3306"`
	Options  string `json:",charset=utf8mb4&parseTime=True&loc=Local"`
}
