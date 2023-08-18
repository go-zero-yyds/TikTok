package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
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

// GetRelationFollowList 获取关注的人的ID
func (l *GetRelationFollowListLogic) GetRelationFollowList(in *social.RelationFollowListReq) (*social.RelationFollowListResp, error) {

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId)
	}

	userList, err := l.svcCtx.CustomDB.QueryFollowListOfUserByUserId(l.ctx, in.UserId)
	//如果未找到关注/没有关注则直接返回空
	if userList == nil || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFollowListResp{UserList: nil}, nil
	}
	return &social.RelationFollowListResp{UserList: userList}, nil
}
