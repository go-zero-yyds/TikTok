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

// GetRelationFollowerList 获取粉丝们的用户ID
func (l *GetRelationFollowerListLogic) GetRelationFollowerList(in *social.RelationFollowerListReq) (*social.RelationFollowerListResp, error) {

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId)
	}

	userList, err := l.svcCtx.CustomDB.QueryFollowerListOfUserByUserId(l.ctx, in.UserId)
	//如果未找到粉丝/没有粉丝则直接返回空
	if userList == nil || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.RelationFollowerListResp{UserList: nil}, nil
	}

	return &social.RelationFollowerListResp{UserList: userList}, nil
}
