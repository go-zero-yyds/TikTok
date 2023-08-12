package svc

import (
	"TikTok/apps/user/rpc/internal/config"
	"TikTok/apps/user/rpc/model"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
	Snowflake *snowflake.Node
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	node, err := snowflake.NewNode(c.Node)
	if err != nil {
		return nil, err
	}
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.Cache),
		Snowflake: node,
	}, nil
}
