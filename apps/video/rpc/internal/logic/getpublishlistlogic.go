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
	publishListResp := &video.PublishListResp{}
	videoList, err := l.svcCtx.Model.VideoListByUserId(context.Background(), in.UserId)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	for _, v := range videoList {
		publishListResp.VideoList = append(publishListResp.VideoList, convertToBasic(v))
	}
	return publishListResp, nil
}

func convertToBasic(v *model.Video) *video.BasicVideoInfo {
	basic := &video.BasicVideoInfo{
		Id:       v.VideoId,
		UserId:   v.UserId,
		PlayUrl:  v.PlayUrl,
		CoverUrl: v.CoverUrl,
		Title:    v.Title,
	}

	return basic
}
