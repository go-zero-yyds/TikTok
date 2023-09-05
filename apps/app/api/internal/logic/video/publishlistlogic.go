package video

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/video/rpc/video"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
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

	tokenID := int64(-1)

	if req.Token != "" {
		tokenID, err = l.svcCtx.JwtAuth.ParseToken(req.Token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = nil
			} else {
				return nil, err
			}

		}
	}
	publishList, err := l.svcCtx.VideoRPC.GetPublishList(l.ctx, &video.PublishListReq{UserId: req.UserID})

	videoInfoList, err := GetVideoInfoList(publishList.VideoList, &tokenID, l.svcCtx, l.ctx)

	if err == apivars.ErrSomeData {
		return &types.PublishListResponse{
			RespStatus: types.RespStatus(apivars.ErrSomeData),
			VideoList:  videoInfoList,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &types.PublishListResponse{
		RespStatus: types.RespStatus(apivars.Success),
		VideoList:  videoInfoList,
	}, nil
}
