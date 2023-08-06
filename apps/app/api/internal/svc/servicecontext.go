package svc

import (
	"TikTok/apps/app/api/internal/config"
	"TikTok/apps/app/api/utils/auth"
	"TikTok/apps/app/api/utils/oss"
	"TikTok/apps/interaction/rpc/interactionclient"
	"TikTok/apps/social/rpc/socialclient"
	"TikTok/apps/user/rpc/userclient"
	"TikTok/apps/video/rpc/videoclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	UserRPC        userclient.User
	VideoRPC       videoclient.Video
	InteractionRPC interactionclient.Interaction
	SocialRPC      socialclient.Social
	JwtAuth        auth.JwtAuth
	OSS            *oss.S3
}

func NewServiceContext(c config.Config) *ServiceContext {
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
		OSS: oss.NewS3(c.OSS.Endpoint, c.OSS.Bucket, c.OSS.AccessKeyID, c.OSS.AccessKeySecret),
	}
}
