package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/logic/common"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowerCountLogic {
	return &GetFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFollowerCount 获取用户粉丝数量
func (l *GetFollowerCountLogic) GetFollowerCount(in *social.FollowerCountReq) (*social.FollowerCountResp, error) {

	//验证用户存在性并注册
	check := common.NewValidateAndRegisterStruct(l.ctx, l.svcCtx)
	ok := check.ValidateAndRegister(in.UserId)
	if ok != true {
		logc.Error(l.ctx, errors.SQLOperateFailed, in.UserId)
	}

	//查询 social 表中用户的 follower_count
	socialStruct, err := l.svcCtx.CustomDB.QueryFieldByUserIdInSocial(l.ctx, in.UserId, "follower_count")
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.FollowerCountResp{FollowerCount: -1}, nil
	}

	return &social.FollowerCountResp{FollowerCount: socialStruct.FollowerCount}, nil
}
