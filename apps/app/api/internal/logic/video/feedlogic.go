package video

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/video/rpc/video"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

type FeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FeedLogic {
	return &FeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FeedLogic) Feed(req *types.FeedRequest) (resp *types.FeedResponse, err error) {
	respStatus := apivars.Success
	var feedReq video.FeedReq
	if req.LatestTime != 0 {
		feedReq.LatestTime = &req.LatestTime
	}
	if req.Token != "" {
		userID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
		if err == nil {
			feedReq.UserId = &userID
		}
	}
	feedBasicList, err := l.svcCtx.VideoRPC.GetFeed(l.ctx, &feedReq)
	if err != nil {
		return nil, err
	}

	feedList, err := GetVideoInfoList(feedBasicList.VideoList, feedReq.UserId, l.svcCtx, l.ctx)

	if err == apivars.ErrSomeData {
		return &types.FeedResponse{
			RespStatus: types.RespStatus(apivars.ErrSomeData),
			VideoList:  feedList,
			NextTime:   feedBasicList.GetNextTime(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &types.FeedResponse{
		RespStatus: types.RespStatus(respStatus),
		VideoList:  feedList,
		NextTime:   feedBasicList.GetNextTime(),
	}, nil
}

// GetVideoInfoList 批量补充 video.BasicVideoInfo 切片的信息，转换为 types.Video 切片。
func GetVideoInfoList(feedBasicList []*video.BasicVideoInfo,
	tokenID *int64, svcCtx *svc.ServiceContext, ctx context.Context) ([]types.Video, error) {

	if feedBasicList == nil {
		return make([]types.Video, 0), nil
	}
	var e *apivars.RespVar
	size := len(feedBasicList)
	feedList, err := mr.MapReduce(func(source chan<- IdxVideo) {
		for i, bv := range feedBasicList {
			source <- IdxVideo{
				BasicVideoInfo: bv,
				idx:            i,
			}
		}
	}, func(item IdxVideo, writer mr.Writer[IdxApiVideo], cancel func(error)) {
		videoInfo, err := TryGetVideoInfo(tokenID, item.BasicVideoInfo, svcCtx, ctx)
		if err != nil {
			e = &apivars.ErrSomeData
			if err != apivars.ErrSomeData {
				return
			}
		}
		writer.Write(IdxApiVideo{
			Video: videoInfo,
			idx:   item.idx,
		})
	}, func(pipe <-chan IdxApiVideo, writer mr.Writer[[]types.Video], cancel func(error)) {
		vs := make([]types.Video, size)
		for item := range pipe {
			v := item
			vs[v.idx] = *v.Video
		}
		writer.Write(vs)
	})

	if err != nil {
		logc.Errorf(ctx, "转换视频列表失败: %v", err)
		return nil, err
	}
	if e != nil {
		return feedList, *e
	}
	return feedList, nil
}

// TryGetVideoInfo 补充 video.BasicVideoInfo 的信息，转换为 types.Video
func TryGetVideoInfo(tokenID *int64, basicVideo *video.BasicVideoInfo, svcCtx *svc.ServiceContext, ctx context.Context) (resp *types.Video, err error) {
	if basicVideo == nil {
		return nil, apivars.ErrInternal
	}

	res := types.Video{
		ID:            basicVideo.Id,
		Author:        types.User{},
		PlayURL:       basicVideo.PlayUrl,
		CoverURL:      basicVideo.CoverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         basicVideo.Title,
	}

	var e *apivars.RespVar

	// 启动goroutines并发调用四个函数
	var wg sync.WaitGroup
	wg.Add(4)

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		ID := int64(-1)
		if tokenID != nil {
			ID = *tokenID
		}
		author, err := user.TryGetUserInfo(ID, basicVideo.UserId, svcCtx, ctx)
		res.Author = *author
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		favoriteCount, err := GetFavoriteCount(svcCtx, ctx, basicVideo.Id)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.FavoriteCount = favoriteCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		commentCount, err := GetCommentCount(svcCtx, ctx, basicVideo.Id)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.CommentCount = commentCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		if tokenID == nil {
			return
		}
		isFavorite, err := GetIsFavorite(svcCtx, ctx, *tokenID, basicVideo.Id)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.IsFavorite = isFavorite
	})

	playLink, err := svcCtx.FS.GetDownloadLink(basicVideo.PlayUrl)
	if err != nil {
		return nil, err
	}
	res.PlayURL = playLink

	coverLink, err := svcCtx.FS.GetDownloadLink(basicVideo.CoverUrl)
	if err != nil {
		return nil, err
	}
	res.CoverURL = coverLink
	wg.Wait()
	if e != nil {
		return &res, *e
	}
	return &res, nil

}

// GetFavoriteCount 点赞数量
func GetFavoriteCount(svcCtx *svc.ServiceContext, ctx context.Context, videoID int64) (int64, error) {
	favoriteCount, err := svcCtx.InteractionRPC.GetFavoriteCountByVideoId(ctx,
		&interaction.FavoriteCountByVideoIdReq{VideoId: videoID})
	if err != nil {
		logc.Errorf(ctx, "获取点赞数失败: %v", err)
		return 0, err
	}
	return favoriteCount.FavoriteCount, nil
}

// GetCommentCount 评论数量
func GetCommentCount(svcCtx *svc.ServiceContext, ctx context.Context, videoID int64) (int64, error) {
	commentCount, err := svcCtx.InteractionRPC.GetCommentCountByVideoId(ctx,
		&interaction.CommentCountByVideoIdReq{VideoId: videoID})
	if err != nil {
		logc.Errorf(ctx, "获取评论数失败: %v", err)
		return 0, err
	}
	return commentCount.CommentCount, nil
}

// GetIsFavorite 是否点赞
func GetIsFavorite(svcCtx *svc.ServiceContext, ctx context.Context, tokenID int64, toVideoId int64) (bool, error) {

	isFavorite, err := svcCtx.InteractionRPC.IsFavorite(ctx, &interaction.IsFavoriteReq{
		UserId:  tokenID,
		VideoId: toVideoId,
	})
	if err != nil {
		logc.Errorf(ctx, "获取是否点赞失败: %v", err)
		return false, err
	}
	return isFavorite.IsFavorite, nil
}

type IdxVideo struct {
	*video.BasicVideoInfo
	idx int
}
type IdxApiVideo struct {
	*types.Video
	idx int
}
