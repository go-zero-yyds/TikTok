package test

import (
	"TikTok/apps/app/api/apivars"
	apiVideo "TikTok/apps/app/api/internal/logic/video"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/test/mock"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/app/api/utils/auth"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/user/rpc/user"
	"TikTok/apps/video/rpc/video"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVideoFeedEmpty(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	videoMock := mock.NewMockVideo(ctl)
	//videoMock.EXPECT().GetWorkCountByUserId(gomock.Any(), gomock.Any()).Return(&video.WorkCountByUserIdResp{WorkCount: 0}, nil)
	videoMock.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(&video.FeedResp{
		VideoList: make([]*video.BasicVideoInfo, 0),
		NextTime:  nil,
	}, nil)

	l := apiVideo.NewFeedLogic(context.TODO(), &svc.ServiceContext{
		VideoRPC: videoMock,
	})

	feed, err := l.Feed(&types.FeedRequest{
		LatestTime: 0,
		Token:      "",
	})
	assert.NoError(t, err)
	assert.Equal(t, &types.FeedResponse{
		RespStatus: types.RespStatus(apivars.Success),
		VideoList:  make([]types.Video, 0),
		NextTime:   0,
	}, feed)
}

func TestVideoFeed(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	u := user.BasicUserInfo{
		Id:              114,
		Name:            "i_ku_yo",
		Avatar:          nil,
		BackgroundImage: nil,
		Signature:       "1919",
	}
	userMock := mock.NewMockUser(ctl)
	userMock.EXPECT().Detail(gomock.Any(), gomock.Any()).Return(&user.BasicUserInfoResp{User: &u}, nil)

	v := video.BasicVideoInfo{
		Id:       123,
		UserId:   u.Id,
		PlayUrl:  "/key",
		CoverUrl: "/key",
		Title:    "suki",
	}
	videoMock := mock.NewMockVideo(ctl)
	//videoMock.EXPECT().GetWorkCountByUserId(gomock.Any(), gomock.Any()).Return(&video.WorkCountByUserIdResp{WorkCount: 0}, nil)
	videoMock.EXPECT().GetPublishList(gomock.Any(), gomock.Any()).Return(&video.PublishListResp{VideoList: []*video.BasicVideoInfo{&v}}, nil)
	videoMock.EXPECT().GetFeed(gomock.Any(), gomock.Any()).Return(&video.FeedResp{
		VideoList: []*video.BasicVideoInfo{&v},
		NextTime:  nil,
	}, nil)

	interactionMock := mock.NewMockInteraction(ctl)
	interactionMock.EXPECT().GetFavoriteCountByUserId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByUserIdResp{FavoriteCount: 514}, nil)
	interactionMock.EXPECT().GetFavoriteCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByVideoIdResp{FavoriteCount: 5}, nil)
	interactionMock.EXPECT().GetFavoriteCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByVideoIdResp{FavoriteCount: 5}, nil)
	//interactionMock.EXPECT().IsFavorite(gomock.Any(), gomock.Any()).Return(&interaction.IsFavoriteResp{IsFavorite: true}, nil)
	interactionMock.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.CommentCountByVideoIdResp{CommentCount: 6}, nil)

	socialMock := mock.NewMockSocial(ctl)
	//socialMock.EXPECT().IsFollow(gomock.Any(), gomock.Any()).Return(&social.IsFollowResp{IsFollow: true}, nil)
	socialMock.EXPECT().GetFollowCount(gomock.Any(), gomock.Any()).Return(&social.FollowCountResp{FollowCount: 19}, nil)
	socialMock.EXPECT().GetFollowerCount(gomock.Any(), gomock.Any()).Return(&social.FollowerCountResp{FollowerCount: 91}, nil)

	fsMock := mock.NewMockFileSystem(ctl)
	fsMock.EXPECT().GetDownloadLink(gomock.Any()).Return("ok", nil)
	fsMock.EXPECT().GetDownloadLink(gomock.Any()).Return("ok", nil)
	l := apiVideo.NewFeedLogic(context.TODO(), &svc.ServiceContext{
		UserRPC:        userMock,
		VideoRPC:       videoMock,
		InteractionRPC: interactionMock,
		SocialRPC:      socialMock,
		FS:             fsMock,
	})

	feed, err := l.Feed(&types.FeedRequest{
		LatestTime: 0,
		Token:      "",
	})
	assert.NoError(t, err)
	assert.Equal(t, &types.FeedResponse{
		RespStatus: types.RespStatus(apivars.Success),
		VideoList: []types.Video{
			{
				ID: v.Id,
				Author: types.User{
					ID:              u.Id,
					Name:            u.Name,
					FollowCount:     19,
					FollowerCount:   91,
					IsFollow:        false,
					Avatar:          u.GetAvatar(),
					BackgroundImage: u.GetBackgroundImage(),
					Signature:       u.Signature,
					TotalFavorited:  5,
					WorkCount:       1,
					FavoriteCount:   514,
				},
				PlayURL:       "ok",
				CoverURL:      "ok",
				FavoriteCount: 5,
				CommentCount:  6,
				IsFavorite:    false,
				Title:         "suki",
			},
		},
		NextTime: 0,
	}, feed)
}

func TestVideoPublishList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	u := user.BasicUserInfo{
		Id:              114,
		Name:            "i_ku_yo",
		Avatar:          nil,
		BackgroundImage: nil,
		Signature:       "1919",
	}
	userMock := mock.NewMockUser(ctl)
	userMock.EXPECT().Detail(gomock.Any(), gomock.Any()).Return(&user.BasicUserInfoResp{User: &u}, nil)

	v := video.BasicVideoInfo{
		Id:       123,
		UserId:   u.Id,
		PlayUrl:  "/key",
		CoverUrl: "/key",
		Title:    "suki",
	}
	videoMock := mock.NewMockVideo(ctl)
	//videoMock.EXPECT().GetWorkCountByUserId(gomock.Any(), gomock.Any()).Return(&video.WorkCountByUserIdResp{WorkCount: 0}, nil)
	videoMock.EXPECT().GetPublishList(gomock.Any(), gomock.Any()).Return(&video.PublishListResp{VideoList: []*video.BasicVideoInfo{&v}}, nil)
	videoMock.EXPECT().GetPublishList(gomock.Any(), gomock.Any()).Return(&video.PublishListResp{
		VideoList: []*video.BasicVideoInfo{&v},
	}, nil)

	interactionMock := mock.NewMockInteraction(ctl)
	interactionMock.EXPECT().GetFavoriteCountByUserId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByUserIdResp{FavoriteCount: 514}, nil)
	interactionMock.EXPECT().GetFavoriteCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByVideoIdResp{FavoriteCount: 5}, nil)
	interactionMock.EXPECT().GetFavoriteCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByVideoIdResp{FavoriteCount: 5}, nil)
	interactionMock.EXPECT().IsFavorite(gomock.Any(), gomock.Any()).Return(&interaction.IsFavoriteResp{IsFavorite: true}, nil)
	interactionMock.EXPECT().GetCommentCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.CommentCountByVideoIdResp{CommentCount: 6}, nil)

	socialMock := mock.NewMockSocial(ctl)
	socialMock.EXPECT().IsFollow(gomock.Any(), gomock.Any()).Return(&social.IsFollowResp{IsFollow: true}, nil)
	socialMock.EXPECT().GetFollowCount(gomock.Any(), gomock.Any()).Return(&social.FollowCountResp{FollowCount: 19}, nil)
	socialMock.EXPECT().GetFollowerCount(gomock.Any(), gomock.Any()).Return(&social.FollowerCountResp{FollowerCount: 91}, nil)

	fsMock := mock.NewMockFileSystem(ctl)
	fsMock.EXPECT().GetDownloadLink(gomock.Any()).Return("ok", nil)
	fsMock.EXPECT().GetDownloadLink(gomock.Any()).Return("ok", nil)

	authJwt := auth.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600,
	}
	token, err := authJwt.CreateToken(810)
	assert.NoError(t, err, "Error create token")

	l := apiVideo.NewPublishListLogic(context.TODO(), &svc.ServiceContext{
		UserRPC:        userMock,
		VideoRPC:       videoMock,
		InteractionRPC: interactionMock,
		SocialRPC:      socialMock,
		FS:             fsMock,
		JwtAuth:        authJwt,
	})

	publishList, err := l.PublishList(&types.PublishListRequest{
		UserID: u.Id,
		Token:  token,
	})
	assert.NoError(t, err)
	assert.Equal(t, &types.PublishListResponse{
		RespStatus: types.RespStatus(apivars.Success),
		VideoList: []types.Video{
			{
				ID: v.Id,
				Author: types.User{
					ID:              u.Id,
					Name:            u.Name,
					FollowCount:     19,
					FollowerCount:   91,
					IsFollow:        true,
					Avatar:          u.GetAvatar(),
					BackgroundImage: u.GetBackgroundImage(),
					Signature:       u.Signature,
					TotalFavorited:  5,
					WorkCount:       1,
					FavoriteCount:   514,
				},
				PlayURL:       "ok",
				CoverURL:      "ok",
				FavoriteCount: 5,
				CommentCount:  6,
				IsFavorite:    true,
				Title:         "suki",
			},
		},
	}, publishList)
}

func TestVideoPublishListEmpty(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	videoMock := mock.NewMockVideo(ctl)
	//videoMock.EXPECT().GetWorkCountByUserId(gomock.Any(), gomock.Any()).Return(&video.WorkCountByUserIdResp{WorkCount: 0}, nil)
	videoMock.EXPECT().GetPublishList(gomock.Any(), gomock.Any()).Return(&video.PublishListResp{
		VideoList: make([]*video.BasicVideoInfo, 0),
	}, nil)

	authJwt := auth.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600,
	}
	token, err := authJwt.CreateToken(810)
	assert.NoError(t, err, "Error create token")

	l := apiVideo.NewPublishListLogic(context.TODO(), &svc.ServiceContext{
		VideoRPC: videoMock,
		JwtAuth:  authJwt,
	})

	publishList, err := l.PublishList(&types.PublishListRequest{
		Token: token,
	})
	assert.NoError(t, err)
	assert.Equal(t, &types.PublishListResponse{
		RespStatus: types.RespStatus(apivars.Success),
		VideoList:  make([]types.Video, 0),
	}, publishList)
}
