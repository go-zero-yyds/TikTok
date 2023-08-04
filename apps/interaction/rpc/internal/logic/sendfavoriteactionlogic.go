package logic

import (
	"context"

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

//调用dB接口中函数，进行点赞/取消操作
//只有在底层数据库出现未知错误会返回err
func (l *SendFavoriteActionLogic) SendFavoriteAction(in *interaction.FavoriteActionReq) (*interaction.FavoriteActionResp, error) {
	success , err := l.svcCtx.DBact.FavoriteAction(in.UserId, in.VideoId, in.ActionType)
	if err != nil{
		return nil , err
	}
	return &interaction.FavoriteActionResp{
		IsSucceed: success,
	}, nil
}
