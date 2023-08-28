package interaction

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/video/rpc/video"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteActionLogic {
	return &FavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteActionLogic) FavoriteAction(req *types.FavoriteActionRequest) (resp *types.FavoriteActionResponse, err error) {
	// todo: add your logic here and delete this line
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.VideoRPC.Detail(l.ctx, &video.BasicVideoInfoReq{VideoId: req.VideoID})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.InteractionRPC.SendFavoriteAction(l.ctx, &interaction.FavoriteActionReq{
		UserId:     tokenID,
		VideoId:    req.VideoID,
		ActionType: req.ActionType,
	})
	if err != nil {
		return nil, err
	}

	// 通知热门视频服务
	_, err = l.svcCtx.VideoRPC.NotifyHotVideo(
		l.ctx,
		&video.NotifyHotVideoReq{
			UserId:     tokenID,
			VideoId:    req.VideoID,
			ActionType: req.ActionType,
		},
	)

	if err != nil {
		l.Logger.Errorf("notify hot video failed: %v", err)
	}

	return &types.FavoriteActionResponse{
		RespStatus: types.RespStatus(apiVars.Success),
	}, nil
}
