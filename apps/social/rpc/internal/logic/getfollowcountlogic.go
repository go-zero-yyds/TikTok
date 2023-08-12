package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowCountLogic {
	return &GetFollowCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFollowCount 获取用户关注数量
func (l *GetFollowCountLogic) GetFollowCount(in *social.FollowCountReq) (*social.FollowCountResp, error) {
	//查询 social 表中是否有该 user_id
	exist, err := l.svcCtx.CustomDB.QueryUserIdIsExistInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回失败
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.FollowCountResp{FollowCount: -1}, nil
	}

	//查询 social 表中用户的 follow_count
	socialStruct, err := l.svcCtx.CustomDB.QueryFieldByUserIdInSocial(l.ctx, in.UserId, "follow_count")
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.FollowCountResp{FollowCount: -1}, nil
	}

	return &social.FollowCountResp{FollowCount: socialStruct.FollowCount}, nil
}
