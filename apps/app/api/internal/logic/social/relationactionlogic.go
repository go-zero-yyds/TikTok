package social

import (
	"TikTok/apps/app/api/apivars"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	sm "TikTok/apps/social/rpc/model"
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/user/rpc/model"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionRequest) (resp *types.RelationActionResponse, err error) {
	if errors.Is(err, model.UserNotFound) {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apivars.UserNotFound),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)
	if tokenID == req.ToUserID {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apivars.NoFollowMyself),
		}, nil
	}
	_, err = l.svcCtx.SocialRPC.SendRelationAction(l.ctx, &social.RelationActionReq{
		UserId:     tokenID,
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
	})
	if errors.Is(err, sm.ErrNotFriend) {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apivars.ErrNotFriend),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &types.RelationActionResponse{
		RespStatus: types.RespStatus(apivars.Success),
	}, nil
}
