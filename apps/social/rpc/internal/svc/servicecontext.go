package svc

import (
	"TikTok/apps/social/rpc/internal/config"
	"TikTok/apps/social/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel model.MessageModel
	FollowModel  model.FollowModel
	SocialModel  model.SocialModel
	CustomDB     model.CustomDB
	Conn         sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		MessageModel: model.NewMessageModel(sqlx.NewMysql(c.DataSource)),
		FollowModel:  model.NewFollowModel(sqlx.NewMysql(c.DataSource)),
		SocialModel:  model.NewSocialModel(sqlx.NewMysql(c.DataSource)),
		CustomDB:     *model.NewCustomDB(sqlx.NewMysql(c.DataSource)),
		Conn:         sqlx.NewMysql(c.DataSource),
	}
}
