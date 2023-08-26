package video

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/video/rpc/video"
	"context"
	"regexp"
	"strconv"

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

	// 参数检查
	matched, err := regexp.MatchString("^\\d+$", strconv.FormatInt(req.UserID, 10)) //是否为纯数字
	if strconv.FormatInt(req.UserID, 10) == "" || matched == false {
		return &types.PublishListResponse{
			RespStatus: types.RespStatus(apiVars.UserIdRuleError),
		}, nil
	}

	if req.Token == "" {
		return &types.PublishListResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
			VideoList:  make([]types.Video, 0),
		}, nil
	}
	// 解析token
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	publishList, err := l.svcCtx.VideoRPC.GetPublishList(l.ctx, &video.PublishListReq{UserId: req.UserID})

	videoInfoList, err := GetVideoInfoList(publishList.VideoList, &tokenID, l.svcCtx, l.ctx)

	if err == apiVars.SomeDataErr {
		return &types.PublishListResponse{
			RespStatus: types.RespStatus(apiVars.SomeDataErr),
			VideoList:  videoInfoList,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &types.PublishListResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		VideoList:  videoInfoList,
	}, nil
}
