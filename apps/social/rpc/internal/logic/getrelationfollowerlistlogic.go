package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"context"
	"github.com/zeromicro/go-zero/core/logc"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRelationFollowerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRelationFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRelationFollowerListLogic {
	return &GetRelationFollowerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRelationFollowerListLogic) GetRelationFollowerList(in *social.RelationFollowerListReq) (*social.RelationFollowerListResp, error) {
	//查询 social 表中是否有该 user_id
	exist, err := l.svcCtx.CustomDB.QueryUserIdExistsInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回空
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFollowerListResp{UserList: nil}, nil
	}

	userList, err := l.svcCtx.CustomDB.QueryUsersOfFollowerListByUserId(l.ctx, in.UserId)
	//如果未找到粉丝/没有粉丝则直接返回空
	if userList == nil || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFollowerListResp{UserList: nil}, nil
	}

	return &social.RelationFollowerListResp{UserList: userList}, nil
}
