package test

import (
	"TikTok/apps/app/api/apiVars"
	apiUser "TikTok/apps/app/api/internal/logic/user"
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

func TestUserRegister(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	userMock := mock.NewMockUser(ctl)
	userMock.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&user.RegisterResp{UserId: int64(114514)}, nil)
	authJwt := auth.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600,
	}
	l := apiUser.NewRegisterLogic(context.TODO(), &svc.ServiceContext{
		UserRPC: userMock,
		JwtAuth: authJwt,
	})

	res, err := l.Register(&types.UserRegisterRequest{
		Username: "1145141919810",
		Password: "1145141919810",
	})
	assert.NoError(t, err, "Error register")
	assert.Equal(t, res.UserID, int64(114514))
	tokenID, err := authJwt.ParseToken(res.Token)
	if err != nil {
		return
	}
	assert.NoError(t, err, "Error register")
	assert.Equal(t, tokenID, int64(114514))
}

func TestUserLogin(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	userMock := mock.NewMockUser(ctl)
	userMock.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&user.LoginResp{UserId: int64(114514)}, nil)
	authJwt := auth.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600,
	}
	l := apiUser.NewLoginLogic(context.TODO(), &svc.ServiceContext{
		UserRPC: userMock,
		JwtAuth: authJwt,
	})

	res, err := l.Login(&types.UserLoginRequest{
		Username: "1145141919810",
		Password: "1145141919810",
	})
	assert.NoError(t, err, "Error Login")
	assert.Equal(t, res.UserID, int64(114514))
	tokenID, err := authJwt.ParseToken(res.Token)
	if err != nil {
		return
	}
	assert.NoError(t, err, "Error Login")
	assert.Equal(t, tokenID, int64(114514))
}

func TestUserDetailNotVideoList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	userMock := mock.NewMockUser(ctl)
	userMock.EXPECT().Detail(gomock.Any(), gomock.Any()).Return(&user.BasicUserInfoResp{User: &user.BasicUserInfo{
		Id:              114,
		Name:            "i_ku_yo",
		Avatar:          nil,
		BackgroundImage: nil,
		Signature:       "1919",
	}}, nil)

	videoMock := mock.NewMockVideo(ctl)
	//videoMock.EXPECT().GetWorkCountByUserId(gomock.Any(), gomock.Any()).Return(&video.WorkCountByUserIdResp{WorkCount: 0}, nil)
	videoMock.EXPECT().GetPublishList(gomock.Any(), gomock.Any()).Return(&video.PublishListResp{VideoList: make([]*video.BasicVideoInfo, 0)}, nil)

	interactionMock := mock.NewMockInteraction(ctl)
	interactionMock.EXPECT().GetFavoriteCountByUserId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByUserIdResp{FavoriteCount: 514}, nil)

	socialMock := mock.NewMockSocial(ctl)
	socialMock.EXPECT().IsFollow(gomock.Any(), gomock.Any()).Return(&social.IsFollowResp{IsFollow: true}, nil)
	socialMock.EXPECT().GetFollowCount(gomock.Any(), gomock.Any()).Return(&social.FollowCountResp{FollowCount: 19}, nil)
	socialMock.EXPECT().GetFollowerCount(gomock.Any(), gomock.Any()).Return(&social.FollowerCountResp{FollowerCount: 91}, nil)

	authJwt := auth.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600,
	}
	l := apiUser.NewDetailLogic(context.TODO(), &svc.ServiceContext{
		UserRPC:        userMock,
		VideoRPC:       videoMock,
		InteractionRPC: interactionMock,
		SocialRPC:      socialMock,
		JwtAuth:        authJwt,
	})

	token, err := authJwt.CreateToken(810)
	assert.NoError(t, err, "Error create token")
	res, err := l.Detail(&types.UserRequest{
		UserID: 114,
		Token:  token,
	})
	assert.NoError(t, err, "Error detail")
	assert.Equal(t, &types.UserResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		User: types.User{
			ID:              114,
			Name:            "i_ku_yo",
			FollowCount:     19,
			FollowerCount:   91,
			IsFollow:        true,
			Avatar:          "",
			BackgroundImage: "",
			Signature:       "1919",
			TotalFavorited:  0,
			WorkCount:       0,
			FavoriteCount:   514,
		},
	}, res)
}

func TestUserDetail(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	userMock := mock.NewMockUser(ctl)
	userMock.EXPECT().Detail(gomock.Any(), gomock.Any()).Return(&user.BasicUserInfoResp{User: &user.BasicUserInfo{
		Id:              114,
		Name:            "i_ku_yo",
		Avatar:          nil,
		BackgroundImage: nil,
		Signature:       "1919",
	}}, nil)

	videoMock := mock.NewMockVideo(ctl)
	//videoMock.EXPECT().GetWorkCountByUserId(gomock.Any(), gomock.Any()).Return(&video.WorkCountByUserIdResp{WorkCount: 0}, nil)
	videoMock.EXPECT().GetPublishList(gomock.Any(), gomock.Any()).Return(&video.PublishListResp{VideoList: []*video.BasicVideoInfo{
		{
			Id:       1,
			UserId:   114,
			PlayUrl:  "",
			CoverUrl: "",
			Title:    "1",
		},
		{
			Id:       2,
			UserId:   114,
			PlayUrl:  "",
			CoverUrl: "",
			Title:    "2",
		},
	}}, nil)

	interactionMock := mock.NewMockInteraction(ctl)
	interactionMock.EXPECT().GetFavoriteCountByUserId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByUserIdResp{FavoriteCount: 514}, nil)
	interactionMock.EXPECT().GetFavoriteCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByVideoIdResp{FavoriteCount: 5}, nil)
	interactionMock.EXPECT().GetFavoriteCountByVideoId(gomock.Any(), gomock.Any()).Return(&interaction.FavoriteCountByVideoIdResp{FavoriteCount: 5}, nil)

	socialMock := mock.NewMockSocial(ctl)
	socialMock.EXPECT().IsFollow(gomock.Any(), gomock.Any()).Return(&social.IsFollowResp{IsFollow: true}, nil)
	socialMock.EXPECT().GetFollowCount(gomock.Any(), gomock.Any()).Return(&social.FollowCountResp{FollowCount: 19}, nil)
	socialMock.EXPECT().GetFollowerCount(gomock.Any(), gomock.Any()).Return(&social.FollowerCountResp{FollowerCount: 91}, nil)

	authJwt := auth.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600,
	}
	l := apiUser.NewDetailLogic(context.TODO(), &svc.ServiceContext{
		UserRPC:        userMock,
		VideoRPC:       videoMock,
		InteractionRPC: interactionMock,
		SocialRPC:      socialMock,
		JwtAuth:        authJwt,
	})

	token, err := authJwt.CreateToken(810)
	assert.NoError(t, err, "Error create token")
	res, err := l.Detail(&types.UserRequest{
		UserID: 114,
		Token:  token,
	})
	assert.NoError(t, err, "Error detail")
	assert.Equal(t, &types.UserResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		User: types.User{
			ID:              114,
			Name:            "i_ku_yo",
			FollowCount:     19,
			FollowerCount:   91,
			IsFollow:        true,
			Avatar:          "",
			BackgroundImage: "",
			Signature:       "1919",
			TotalFavorited:  10,
			WorkCount:       2,
			FavoriteCount:   514,
		},
	}, res)
}
