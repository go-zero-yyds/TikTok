package svc

import (
	"TikTok/apps/app/api/internal/config"
	"TikTok/apps/app/api/utils/auth"
	"TikTok/apps/interaction/rpc/interactionclient"
	"TikTok/apps/social/rpc/socialclient"
	"TikTok/apps/user/rpc/userclient"
	"TikTok/apps/video/rpc/videoclient"
	"TikTok/pkg/FileSystem"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	UserRPC        userclient.User
	VideoRPC       videoclient.Video
	InteractionRPC interactionclient.Interaction
	SocialRPC      socialclient.Social
	JwtAuth        auth.JwtAuth
	FS             FileSystem.FileSystem
}

func NewServiceContext(c config.Config) *ServiceContext {
	var fs FileSystem.FileSystem
	if c.FS.AwsS3.Endpoint != "" {
		fs = FileSystem.NewS3(c.FS.AwsS3.Endpoint, c.FS.AwsS3.Bucket, c.FS.Prefix, c.FS.AwsS3.AccessKeyID, c.FS.AwsS3.AccessKeySecret)
	}
	return &ServiceContext{
		Config:         c,
		UserRPC:        userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		VideoRPC:       videoclient.NewVideo(zrpc.MustNewClient(c.VideoRPC)),
		InteractionRPC: interactionclient.NewInteraction(zrpc.MustNewClient(c.InteractionRPC)),
		SocialRPC:      socialclient.NewSocial(zrpc.MustNewClient(c.SocialRPC)),
		JwtAuth: auth.JwtAuth{
			AccessSecret: []byte(c.JwtAuth.AccessSecret),
			AccessExpire: c.JwtAuth.AccessExpire,
		},
		FS: fs,
	}
}
