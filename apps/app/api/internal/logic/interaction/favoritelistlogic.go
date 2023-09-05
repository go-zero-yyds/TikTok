package interaction

import (
	"TikTok/apps/app/api/apivars"
	videoApi "TikTok/apps/app/api/internal/logic/video"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/interaction/rpc/interaction"
	"TikTok/apps/video/rpc/video"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type FavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteListLogic {
	return &FavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteListLogic) FavoriteList(req *types.FavoriteListRequest) (resp *types.FavoriteListResponse, err error) {

	tokenID := int64(-1)
	if req.Token != "" {
		tokenID, err = l.svcCtx.JwtAuth.ParseToken(req.Token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				err = nil
			} else {
				return nil, err
			}
		}
	}

	list, err := l.svcCtx.InteractionRPC.GetFavoriteList(l.ctx, &interaction.FavoriteListReq{UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	if list.VideoList == nil {
		list.VideoList = make([]int64, 0)
	}
	e := apivars.Success
	videoList, err := mr.MapReduce(func(source chan<- int64) {
		for _, bv := range list.VideoList {
			source <- bv
		}
	}, func(item int64, writer mr.Writer[*types.Video], cancel func(error)) {

		detail, err := l.svcCtx.VideoRPC.Detail(l.ctx, &video.BasicVideoInfoReq{VideoId: item})
		if err != nil {
			return
		}

		videoInfo, err := videoApi.TryGetVideoInfo(&tokenID, detail.Video, l.svcCtx, l.ctx)
		if err != nil {
			e = apivars.SomeDataErr
			if err != apivars.SomeDataErr {
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

	return &types.FavoriteListResponse{
		RespStatus: types.RespStatus(e),
		VideoList:  videoList,
	}, nil
}
