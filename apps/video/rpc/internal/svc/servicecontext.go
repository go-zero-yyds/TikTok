package svc

import (
	"TikTok/apps/video/rpc/internal/config"
	"TikTok/apps/video/rpc/model"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	Model     model.VideoModel
	Snowflake *snowflake.Node
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	return &ServiceContext{
		Config:    c,
		Model:     model.NewVideoModel(sqlx.NewMysql(c.DataSourse), c.Cache),
		Snowflake: node,
	}, nil
}
