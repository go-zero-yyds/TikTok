package logic

import (
	"context"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFavoriteLogic {
	return &IsFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// IsFavorite 调用dB接口中函数，获取用户是否给某个视频点赞
// 只有在底层数据库出现未知错误会返回err
func (l *IsFavoriteLogic) IsFavorite(in *interaction.IsFavoriteReq) (*interaction.IsFavoriteResp, error) {
	exist, err := l.svcCtx.DBAction.IsFavorite(l.ctx, in.UserId, in.VideoId)
	if err != nil {
		return nil, err
	}
	return &interaction.IsFavoriteResp{
		IsFavorite: exist,
	}, nil
}
