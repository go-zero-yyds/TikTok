package svc

import (
	"TikTok/apps/social/rpc/dao"
	"TikTok/apps/social/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	MessageModel dao.MessageModel
	FollowModel  dao.FollowModel
	SocialModel  dao.SocialModel
	CustomDB     dao.CustomDB
	Conn         sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		MessageModel: dao.NewMessageModel(sqlx.NewMysql(c.DataSource)),
		FollowModel:  dao.NewFollowModel(sqlx.NewMysql(c.DataSource)),
		SocialModel:  dao.NewSocialModel(sqlx.NewMysql(c.DataSource)),
		CustomDB:     *dao.NewCustomDB(sqlx.NewMysql(c.DataSource)),
		Conn:         sqlx.NewMysql(c.DataSource),
	}
}
