package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DBSource     string
	Cache        cache.CacheConf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	KqCPersonalCallbackConf kq.KqConf
	FS             struct {
		AwsS3 struct {
			Endpoint        string
			AccessKeyID     string
			AccessKeySecret string
			Bucket          string
		} `json:",optional"`
		Webdav struct {
			URL                string
			User               string
			Password           string
			DownloadLinkPrefix string
		} `json:",optional"`
		Prefix string
	}
}
