package logic

import (
	"context"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFriendListLogic {
	return &GetRelationFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationFriendListLogic) GetRelationFriendList(in *social.RelationFriendListReq) (*social.RelationFriendListResp, error) {
	list, err := l.svcCtx.DBAction.FollowList(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	var userList []*social.FriendUser
	for _, v := range list {
		message, err := l.svcCtx.DBAction.NowMessage(l.ctx, in.UserId, v)
		if err != nil {
			return nil, err
		}
		userList = append(userList, message)
	}
	return &social.RelationFriendListResp{
		UserList: userList,
	}, nil
}
