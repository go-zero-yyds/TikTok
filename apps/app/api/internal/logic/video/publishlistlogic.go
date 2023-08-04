package video

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/video/rpc/video"
	"context"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishListLogic {
	return &PublishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishListLogic) PublishList(req *types.PublishListRequest) (resp *types.PublishListResponse, err error) {
	// todo: add your logic here and delete this line
	// 解析token
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	publishList, err := l.svcCtx.VideoRPC.GetPublishList(l.ctx, &video.PublishListReq{UserId: req.UserID})

	videoInfoList, err := GetVideoInfoList(&publishList.VideoList, &tokenID, l.svcCtx, l.ctx)

	if err == apiVars.SomeDataErr {
		return &types.PublishListResponse{
			RespStatus: types.RespStatus(apiVars.SomeDataErr),
			VideoList:  *videoInfoList,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &types.PublishListResponse{
		RespStatus: types.RespStatus(apiVars.SomeDataErr),
		VideoList:  *videoInfoList,
	}, nil
}
