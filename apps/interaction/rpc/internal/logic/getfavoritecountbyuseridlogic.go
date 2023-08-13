package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFavoriteCountByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFavoriteCountByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoriteCountByUserIdLogic {
	return &GetFavoriteCountByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFavoriteCountByUserId 调用dB接口中函数，获取用户点赞数量
// 只有在底层数据库出现未知错误会返回err
func (l *GetFavoriteCountByUserIdLogic) GetFavoriteCountByUserId(in *interaction.FavoriteCountByUserIdReq) (*interaction.FavoriteCountByUserIdResp, error) {
	count, err := l.svcCtx.DBAction.FavoriteCountByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &interaction.FavoriteCountByUserIdResp{
		FavoriteCount: count,
	}, nil
}
