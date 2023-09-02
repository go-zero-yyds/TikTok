package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowerListLogic {
	return &FollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowerListLogic) FollowerList(req *types.RelationFollowerListRequest) (resp *types.RelationFollowerListResponse, err error) {

	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)
	list, err := l.svcCtx.SocialRPC.GetRelationFollowerList(l.ctx, &social.RelationFollowerListReq{UserId: req.UserID})

	if err != nil {
		return nil, err
	}
	infoList, err := GetUserInfoList(list.UserList, tokenID, l.svcCtx, l.ctx)
	if err != nil {
		return nil, err
	}
	return &types.RelationFollowerListResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		UserList:   infoList,
	}, nil
}

// GetUserInfoList 根据 userID 切片，转换为 types.User 切片。
func GetUserInfoList(userList []int64,
	userID int64, svcCtx *svc.ServiceContext, ctx context.Context) ([]types.User, error) {

	if userList == nil {
		return make([]types.User, 0), nil
	}
	var e *apiVars.RespErr

	userInfoList, err := mr.MapReduce(func(source chan<- int64) {
		for _, bv := range userList {
			source <- bv
		}
	}, func(item int64, writer mr.Writer[*types.User], cancel func(error)) {
		userInfo, err := user.TryGetUserInfo(userID, item, svcCtx, ctx)
		if err != nil {
			e = &apiVars.SomeDataErr
			if err != apiVars.SomeDataErr {
				return
			}
		}
		writer.Write(userInfo)
	}, func(pipe <-chan *types.User, writer mr.Writer[[]types.User], cancel func(error)) {
		var vs []types.User
		for item := range pipe {
			v := item
			vs = append(vs, *v)
		}
		writer.Write(vs)
	})

	if err != nil {
		logc.Errorf(ctx, "转换用户列表失败: %v", err)
		return nil, err
	}
	if e == nil {
		return userInfoList, nil
	}
	return userInfoList, *e
}
