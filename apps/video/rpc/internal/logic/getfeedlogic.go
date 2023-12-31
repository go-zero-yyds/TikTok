package logic

import (
	"context"
	"time"

	"TikTok/apps/video/rpc/internal/svc"
	"TikTok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

const minYear = 1
const maxYear = 9999

type GetFeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedLogic {
	return &GetFeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeedLogic) GetFeed(in *video.FeedReq) (*video.FeedResp, error) {
	// todo: add your logic here and delete this line
	var nowTime int64
	if in.LatestTime != nil {
		nowTime = *in.LatestTime
	} else {
		nowTime = time.Now().UnixMilli()
	}

	lastYear := time.UnixMilli(nowTime).Year()
	if lastYear < minYear || lastYear > maxYear {
		nowTime = time.Now().UnixMilli()
	}

	videoList, err := l.svcCtx.Model.VideoFeedByLastTime(l.ctx, nowTime)
	if err != nil {
		return &video.FeedResp{}, err
	}

	publishListResp := make([]*video.BasicVideoInfo, 0, 30)
	for _, v := range videoList {
		publishListResp = append(publishListResp, convertToBasic(v))
	}

	lastIndex := -1
	nextTime := int64(0)
	if len(publishListResp) > 0 {
		lastIndex = len(publishListResp) - 1
		nextTime = videoList[lastIndex].CreateTime.UnixMilli()
	}

	return &video.FeedResp{
		VideoList: publishListResp,
		NextTime:  &nextTime,
	}, nil
}
