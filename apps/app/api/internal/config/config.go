package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// RPC 配置
	VideoRPC       zrpc.RpcClientConf
	UserRPC        zrpc.RpcClientConf
	InteractionRPC zrpc.RpcClientConf
	SocialRPC      zrpc.RpcClientConf
	// JWT 配置
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	FS struct {
		Type            string
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
		Prefix          string
		Bucket          string
	}
}
