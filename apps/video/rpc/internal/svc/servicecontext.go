package svc

import (
	"TikTok/apps/video/rpc/internal/config"
	"TikTok/apps/video/rpc/model"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	Model     model.VideoModel
	Snowflake *snowflake.Node
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	var startTime time.Time
	startTime, err := time.Parse("2006-01-02", c.Snow.StartTime)
	if err != nil {
		return nil, err
	}

	snowflake.Epoch = startTime.UnixNano() / 1000000
	node, err := snowflake.NewNode(c.Snow.Node)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:    c,
		Model:     model.NewVideoModel(sqlx.NewMysql(c.DataSourse), c.Cache),
		Snowflake: node,
	}, nil
}
