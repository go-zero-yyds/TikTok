package user

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/user/rpc/user"
	"TikTok/apps/video/rpc/video"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"sync"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserRequest) (resp *types.UserResponse, err error) {
	// 解析token
	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	return GetUserInfo(tokenID, req.UserID, l.svcCtx, l.ctx)
}

func GetUserInfo(tokenID, toUserId int64, svcCtx *svc.ServiceContext, ctx context.Context) (resp *types.UserResponse, err error) {
	res, err := GetBasicUserInfo(svcCtx, ctx, toUserId)
	if err != nil {
		return nil, err
	}

	// 启动goroutines并发调用五个函数
	var wg sync.WaitGroup
	wg.Add(5)

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		// 错误降级, 不影响获取user的基本信息。
		isFollow, _ := GetIsFollow(svcCtx, ctx, tokenID, toUserId)
		res.User.IsFollow = isFollow
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		// 错误降级, 不影响获取user的基本信息。
		followCount, _ := GetFollowCount(svcCtx, ctx, toUserId)
		res.User.FollowCount = followCount

	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		// 错误降级, 不影响获取user的基本信息。
		followerCount, _ := GetFollowerCount(svcCtx, ctx, toUserId)
		res.User.FollowerCount = followerCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		// 错误降级, 不影响获取user的基本信息。
		favoriteCount, _ := GetFavoriteCount(svcCtx, ctx, toUserId)
		res.User.FavoriteCount = favoriteCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		videoList, err := GetPublishList(svcCtx, ctx, toUserId)
		// 错误降级, 不影响获取user的基本信息。
		if err != nil {
			return
		}
		res.User.WorkCount = int64(len(videoList))
		totalFavorited, err := mr.MapReduce(func(source chan<- int64) {
			for _, info := range videoList {
				source <- info.Id
			}
		}, func(item int64, writer mr.Writer[int64], cancel func(error)) {
			videoFavoriteCount, err := svcCtx.InteractionRPC.GetFavoriteCountByVideoId(
				ctx, &interaction.FavoriteCountByVideoIdReq{VideoId: item})
			if err != nil {
				logc.Errorf(ctx, "获取视频点赞数失败: %v", err)
				cancel(err)
				return
			}
			writer.Write(videoFavoriteCount.FavoriteCount)
		}, func(pipe <-chan int64, writer mr.Writer[int64], cancel func(error)) {
			sum := int64(0)
			for item := range pipe {
				sum = sum + item
			}
			writer.Write(sum)
		})
		res.User.TotalFavorited = totalFavorited
	})
	wg.Wait()

	return res, nil

}

// GetPublishList 获取视频列表，根据userId
func GetPublishList(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) ([]*video.BasicVideoInfo, error) {

	videoList, err := svcCtx.VideoRPC.GetPublishList(ctx, &video.PublishListReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取被获赞数失败: %v", err)
		return nil, err
	}
	return videoList.GetVideoList(), nil
}

// GetFavoriteCount 点赞数量
func GetFavoriteCount(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (int64, error) {
	favoriteCount, err := svcCtx.InteractionRPC.GetFavoriteCountByUserId(ctx,
		&interaction.FavoriteCountByUserIdReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取点赞数失败: %v", err)
		return 0, err
	}
	return favoriteCount.GetFavoriteCount(), nil
}

func GetFollowerCount(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (int64, error) {
	followerCount, err := svcCtx.SocialRPC.GetFollowerCount(ctx, &social.FollowerCountReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取被关注数失败: %v", err)
		return 0, err
	}
	return followerCount.GetFollowerCount(), nil
}

// GetFollowCount 获取关注数
func GetFollowCount(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (int64, error) {
	followCount, err := svcCtx.SocialRPC.GetFollowCount(ctx, &social.FollowCountReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取关注数失败: %v", err)
		return 0, err
	}
	return followCount.GetFollowCount(), nil
}

// GetIsFollow 是否关注
func GetIsFollow(svcCtx *svc.ServiceContext, ctx context.Context, tokenID int64, toUserId int64) (bool, error) {

	isFollow, err := svcCtx.SocialRPC.IsFollow(ctx, &social.IsFollowReq{
		UserId:   tokenID,
		ToUserId: toUserId,
	})
	if err != nil {
		logc.Errorf(ctx, "获取是否关注失败: %v", err)
		return false, err
	}
	return isFollow.GetIsFollow(), nil
}

// GetBasicUserInfo 获取用户基本信息
func GetBasicUserInfo(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (*types.UserResponse, error) {
	basicUserInfo, err := svcCtx.UserRPC.Detail(ctx, &user.BasicUserInfoReq{UserId: toUserId})
	if err != nil {
		return nil, err
	}

	res := &types.UserResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		User: types.User{
			ID:              basicUserInfo.GetUser().GetId(),
			Name:            basicUserInfo.GetUser().Name,
			Avatar:          basicUserInfo.GetUser().GetAvatar(),
			BackgroundImage: basicUserInfo.GetUser().GetBackgroundImage(),
			Signature:       basicUserInfo.GetUser().GetSignature(),
		},
	}
	return res, nil
}
