package video

import (
	"TikTok/apps/app/api/apiVars"
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
	// todo: add your logic here and delete this line
	respStatus := apiVars.Success
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

	feedList, err := mr.MapReduce(func(source chan<- *video.BasicVideoInfo) {
		for _, bv := range feedBasicList.VideoList {
			source <- bv
		}
	}, func(item *video.BasicVideoInfo, writer mr.Writer[*types.Video], cancel func(error)) {
		videoInfo, err := TryGetVideoInfo(feedReq.UserId, item, l.svcCtx, l.ctx)
		if err != nil {
			respStatus = apiVars.SomeDataErr
			if err != apiVars.SomeDataErr {
				return
			}
		}
		writer.Write(videoInfo)
	}, func(pipe <-chan *types.Video, writer mr.Writer[[]types.Video], cancel func(error)) {
		var vs []types.Video
		for item := range pipe {
			v := item
			vs = append(vs, *v)
		}
		writer.Write(vs)
	})
	return &types.FeedResponse{
		RespStatus: types.RespStatus(respStatus),
		VideoList:  feedList,
		NextTime:   0,
	}, nil
}

func TryGetVideoInfo(tokenID *int64, basicVideo *video.BasicVideoInfo, svcCtx *svc.ServiceContext, ctx context.Context) (resp *types.Video, err error) {
	if basicVideo == nil {
		return nil, apiVars.InternalError
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

	var e *apiVars.RespErr

	// 启动goroutines并发调用四个函数
	var wg sync.WaitGroup
	wg.Add(3)

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		favoriteCount, err := GetFavoriteCount(svcCtx, ctx, basicVideo.Id)
		if err != nil {
			e = &apiVars.SomeDataErr
			return
		}
		res.FavoriteCount = favoriteCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		commentCount, err := GetCommentCount(svcCtx, ctx, basicVideo.Id)
		if err != nil {
			e = &apiVars.SomeDataErr
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
			e = &apiVars.SomeDataErr
			return
		}
		res.IsFavorite = isFavorite
	})
	wg.Wait()
	return &res, *e

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
		logc.Errorf(ctx, "获取是否关注失败: %v", err)
		return false, err
	}
	return isFavorite.IsFavorite, nil
}
