package logic

import (
	"context"
	"strconv"

	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/interaction/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFavoriteActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFavoriteActionLogic {
	return &SendFavoriteActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendFavoriteAction 调用dB接口中函数，进行点赞/取消操作
// 只有在底层数据库出现未知错误会返回err
func (l *SendFavoriteActionLogic) SendFavoriteAction(in *interaction.FavoriteActionReq) (*interaction.FavoriteActionResp, error) {
	success, err := l.svcCtx.DBAction.FavoriteAction(l.ctx, in.UserId, in.VideoId, strconv.Itoa(int(in.ActionType)))
	if err != nil {
		return nil, err
	}
	return &interaction.FavoriteActionResp{
		IsSucceed: success,
	}, nil
}
