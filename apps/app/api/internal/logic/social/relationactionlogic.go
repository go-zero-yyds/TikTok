package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/social/rpc/social"
	"TikTok/apps/user/rpc/model"
	"TikTok/apps/user/rpc/user"
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

	// 参数检查
	if req.Token == "" {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserRPC.Detail(l.ctx, &user.BasicUserInfoReq{UserId: req.ToUserID})
	if errors.Is(err, model.UserNotFound) {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apiVars.UserNotFound),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SocialRPC.SendRelationAction(l.ctx, &social.RelationActionReq{
		UserId:     tokenID,
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
	})

	if err != nil {
		return nil, err
	}
	return &types.RelationActionResponse{
		RespStatus: types.RespStatus(apiVars.Success),
	}, nil
}
