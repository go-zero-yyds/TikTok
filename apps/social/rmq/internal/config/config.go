package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	MessageForRobots kq.KqConf
	PersonalCallback kq.KqConf
	SocialRPC        zrpc.RpcClientConf
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
