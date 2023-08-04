package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteListLogic {
	return &GetFavoriteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//调用dB接口中函数，获取用户点赞视频列表
//只有在底层数据库出现未知错误会返回err
func (l *GetFavoriteListLogic) GetFavoriteList(in *interaction.FavoriteListReq) (*interaction.FavoriteListResp, error) {
	videolist , err := l.svcCtx.DBact.FavoriteList(in.UserId)
	if err != nil{
		return nil ,err
	}
	return &interaction.FavoriteListResp{
		VideoList: videolist,
	}, nil
}
