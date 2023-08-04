package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByVideoIdLogic {
	return &GetFavoriteCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//调用dB接口中函数，获取视频点赞数量
//只有在底层数据库出现未知错误会返回err
func (l *GetFavoriteCountByVideoIdLogic) GetFavoriteCountByVideoId(in *interaction.FavoriteCountByVideoIdReq) (*interaction.FavoriteCountByVideoIdResp, error) {
	count , err := l.svcCtx.DBact.FavoriteCountByVideoId(in.VideoId)
	if err != nil{
		return nil , err;
	}
	return &interaction.FavoriteCountByVideoIdResp{
		FavoriteCount: count,
	}, nil
}
