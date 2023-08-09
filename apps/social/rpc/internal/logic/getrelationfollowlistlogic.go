package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"context"
	"github.com/zeromicro/go-zero/core/logc"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFollowListLogic {
	return &GetRelationFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationFollowListLogic) GetRelationFollowList(in *social.RelationFollowListReq) (*social.RelationFollowListResp, error) {
	//查询 social 表中是否有该 user_id
	exist, err := l.svcCtx.CustomDB.QueryUserIdExistsInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回空
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFollowListResp{UserList: nil}, nil
	}

	userList, err := l.svcCtx.CustomDB.QueryUsersOfFollowListByUserId(l.ctx, in.UserId)
	//如果未找到关注/没有关注则直接返回空
	if userList == nil || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFollowListResp{UserList: nil}, nil
	}
	return &social.RelationFollowListResp{UserList: userList}, nil
}
