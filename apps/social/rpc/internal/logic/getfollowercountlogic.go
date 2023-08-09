package logic

import (
	"TikTok/apps/social/rpc/internal/errors"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"strconv"

	"TikTok/apps/social/rpc/internal/svc"
	"TikTok/apps/social/rpc/social"

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

func (l *GetFollowerCountLogic) GetFollowerCount(in *social.FollowerCountReq) (*social.FollowerCountResp, error) {
	//查询 social 表中是否有该 user_id
	exist, err := l.svcCtx.CustomDB.QueryUserIdExistsInSocial(l.ctx, in.UserId)

	//如果不存在则直接返回失败
	if exist == false || err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.FollowerCountResp{FollowerCount: -1}, nil
	}

	//查询 social 表中用户的 follower_count
	socialStruct, err := l.svcCtx.CustomDB.QueryFieldByUserIdInSocial(l.ctx, in.UserId, "follower_count")
	if err != nil {
		logc.Error(l.ctx, errors.RecordNotFound, in.UserId)
		return &social.FollowerCountResp{FollowerCount: -1}, nil
	}

	//string转int64
	followerCount, _ := strconv.ParseInt(socialStruct.FollowerCount, 10, 64)

	return &social.FollowerCountResp{FollowerCount: followerCount}, nil
}
