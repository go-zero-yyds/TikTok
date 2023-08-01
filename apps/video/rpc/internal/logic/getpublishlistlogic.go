package logic

import (
	"context"
	"fmt"

	"TikTok/apps/video/rpc/internal/svc"
	"TikTok/apps/video/rpc/model"
	"TikTok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPublishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPublishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPublishListLogic {
	return &GetPublishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPublishListLogic) GetPublishList(in *video.PublishListReq) (*video.PublishListResp, error) {
	var publishListResp *video.PublishListResp
	mvideoList, err := l.svcCtx.Model.VideoListByUserId(context.Background(), in.UserId)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	for i, v := range mvideoList {
		publishListResp.VideoList[i] = convertToBasic(v)
	}
	return publishListResp, nil
}

func convertToBasic(mvideo *model.Video) *video.BasicVideoInfo {

	basic := &video.BasicVideoInfo{
		Id:       mvideo.VideoId,
		UserId:   mvideo.UserId,
		PlayUrl:  mvideo.PlayUrl,
		CoverUrl: mvideo.CoverUrl,
		Title:    mvideo.Title,
	}

	return basic
}
