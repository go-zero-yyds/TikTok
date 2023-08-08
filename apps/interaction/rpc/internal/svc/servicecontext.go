package svc

import (
	"TikTok/apps/interaction/rpc/internal/config"
	"TikTok/apps/interaction/rpc/model"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	Snowflake *snowflake.Node
	DBact     *model.DBaction
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	return &ServiceContext{
		Config:    c,
		Snowflake: node,
		DBact:     model.NewDBaction(sqlx.NewMysql(c.DBsource), c.Cache),
	}, nil
}
