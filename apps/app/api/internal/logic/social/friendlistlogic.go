package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/logic/user"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/mr"
	"regexp"
	"strconv"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.RelationFriendListRequest) (resp *types.RelationFriendListResponse, err error) {

	// 参数检查
	matched, err := regexp.MatchString("^\\d+$", strconv.FormatInt(req.UserID, 10)) //是否为纯数字
	if strconv.FormatInt(req.UserID, 10) == "" || matched == false {
		return &types.RelationFriendListResponse{
			RespStatus: types.RespStatus(apiVars.UserIdRuleError),
		}, nil
	}

	if req.Token == "" {
		return &types.RelationFriendListResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.SocialRPC.GetRelationFriendList(l.ctx, &social.RelationFriendListReq{UserId: tokenID})
	if err != nil {
		return nil, err
	}
	infoList, err := GetFriendInfoList(list.UserList, tokenID, l.svcCtx, l.ctx)
	if err != nil {
		return nil, err
	}
	return &types.RelationFriendListResponse{
		RespStatus: types.RespStatus(apiVars.Success),
		UserList:   infoList,
	}, nil
}

// GetFriendInfoList 根据 userID 切片，转换为 types.User 切片。
func GetFriendInfoList(userList []*social.FriendUser,
	userID int64, svcCtx *svc.ServiceContext, ctx context.Context) ([]types.FriendUser, error) {

	if userList == nil {
		return nil, apiVars.InternalError
	}
	var e *apiVars.RespErr

	userInfoList, err := mr.MapReduce(func(source chan<- *social.FriendUser) {
		for _, bv := range userList {
			source <- bv
		}
	}, func(item *social.FriendUser, writer mr.Writer[*types.FriendUser], cancel func(error)) {
		userInfo, err := user.TryGetUserInfo(userID, item.UserId, svcCtx, ctx)
		if err != nil {
			e = &apiVars.SomeDataErr
			if err != apiVars.SomeDataErr {
				return
			}
		}
		writer.Write(&types.FriendUser{
			User:    *userInfo,
			Message: item.GetMessage(),
			MsgType: item.MsgType,
		})
	}, func(pipe <-chan *types.FriendUser, writer mr.Writer[[]types.FriendUser], cancel func(error)) {
		var vs []types.FriendUser
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
