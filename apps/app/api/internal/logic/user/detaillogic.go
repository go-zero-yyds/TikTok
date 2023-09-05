package user

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
	"TikTok/apps/video/rpc/video"
	"context"
	"errors"
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

	// 参数检查
	if req.Token == "" {
		return &types.UserResponse{
			RespStatus: types.RespStatus(apivars.ErrNotLogged),
		}, nil
	}

	// 解析token
	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)
	userInfo, err := TryGetUserInfo(tokenID, req.UserID, l.svcCtx, l.ctx)
	if errors.Is(err, model.UserNotFound) {
		return &types.UserResponse{
			RespStatus: types.RespStatus(apivars.ErrUserNotFound),
		}, nil
	}
	if err == apivars.ErrSomeData {
		return &types.UserResponse{
			RespStatus: types.RespStatus(apivars.ErrSomeData),
			User:       *userInfo,
		}, nil
	}

	if err != nil {
		return nil, err
	}
	return &types.UserResponse{
		RespStatus: types.RespStatus(apivars.Success),
		User:       *userInfo,
	}, nil
}

// TryGetUserInfo 尝试获取全部用户信息
// 部分非必要信息未获取到时，返回 apivars.ErrSomeData
func TryGetUserInfo(tokenID, toUserId int64, svcCtx *svc.ServiceContext, ctx context.Context) (resp *types.User, err error) {
	res, err := GetBasicUserInfo(svcCtx, ctx, toUserId)
	if err != nil {
		return nil, err
	}

	var e *apivars.RespVar

	// 启动goroutines并发调用五个函数
	var wg sync.WaitGroup
	wg.Add(5)

	if tokenID != -1 {
		threading.GoSafeCtx(ctx, func() {
			defer wg.Done()
			// 错误降级, 不影响获取user的基本信息。
			isFollow, err := GetIsFollow(svcCtx, ctx, tokenID, toUserId)
			if err != nil {
				e = &apivars.ErrSomeData
				return
			}
			res.IsFollow = isFollow
		})
	} else {
		wg.Done()
	}

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		followCount, err := GetFollowCount(svcCtx, ctx, toUserId)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.FollowCount = followCount

	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		followerCount, err := GetFollowerCount(svcCtx, ctx, toUserId)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.FollowerCount = followerCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		favoriteCount, err := GetFavoriteCount(svcCtx, ctx, toUserId)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.FavoriteCount = favoriteCount
	})

	threading.GoSafeCtx(ctx, func() {
		defer wg.Done()
		videoList, err := GetPublishList(svcCtx, ctx, toUserId)
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.WorkCount = int64(len(videoList))
		totalFavorited, err := mr.MapReduce(func(source chan<- int64) {
			for _, info := range videoList {
				source <- info.Id
			}
		}, func(item int64, writer mr.Writer[int64], cancel func(error)) {
			videoFavoriteCount, err := svcCtx.InteractionRPC.GetFavoriteCountByVideoId(
				ctx, &interaction.FavoriteCountByVideoIdReq{VideoId: item})
			if err != nil {
				e = &apivars.ErrSomeData
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
		if err != nil {
			e = &apivars.ErrSomeData
			return
		}
		res.TotalFavorited = totalFavorited
	})
	wg.Wait()
	if e != nil {
		return res, *e
	}
	return res, nil

}

// GetPublishList 获取视频列表，根据userId
func GetPublishList(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) ([]*video.BasicVideoInfo, error) {

	videoList, err := svcCtx.VideoRPC.GetPublishList(ctx, &video.PublishListReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取被获赞数失败: %v", err)
		return nil, err
	}
	return videoList.VideoList, nil
}

// GetFavoriteCount 点赞数量
func GetFavoriteCount(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (int64, error) {
	favoriteCount, err := svcCtx.InteractionRPC.GetFavoriteCountByUserId(ctx,
		&interaction.FavoriteCountByUserIdReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取点赞数失败: %v", err)
		return 0, err
	}
	return favoriteCount.FavoriteCount, nil
}

func GetFollowerCount(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (int64, error) {
	followerCount, err := svcCtx.SocialRPC.GetFollowerCount(ctx, &social.FollowerCountReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取被关注数失败: %v", err)
		return 0, err
	}
	return followerCount.FollowerCount, nil
}

// GetFollowCount 获取关注数
func GetFollowCount(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (int64, error) {
	followCount, err := svcCtx.SocialRPC.GetFollowCount(ctx, &social.FollowCountReq{UserId: toUserId})
	if err != nil {
		logc.Errorf(ctx, "获取关注数失败: %v", err)
		return 0, err
	}
	return followCount.FollowCount, nil
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
	return isFollow.IsFollow, nil
}

// GetBasicUserInfo 获取用户基本信息
func GetBasicUserInfo(svcCtx *svc.ServiceContext, ctx context.Context, toUserId int64) (*types.User, error) {
	basicUserInfo, err := svcCtx.UserRPC.Detail(ctx, &user.BasicUserInfoReq{UserId: toUserId})
	if err != nil {
		return nil, err
	}

	res := &types.User{
		ID:              basicUserInfo.User.Id,
		Name:            basicUserInfo.User.Name,
		Avatar:          basicUserInfo.User.GetAvatar(),
		BackgroundImage: basicUserInfo.User.GetBackgroundImage(),
		Signature:       basicUserInfo.User.Signature,
	}
	res.Avatar, err = svcCtx.FS.GetDownloadLink(res.Avatar)
	if err != nil {
		return nil, err
	}
	res.BackgroundImage, err = svcCtx.FS.GetDownloadLink(res.BackgroundImage)
	if err != nil {
		return nil, err
	}
	return res, nil
}
