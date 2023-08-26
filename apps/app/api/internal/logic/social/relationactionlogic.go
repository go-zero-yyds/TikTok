package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/social/rpc/social"
	"context"
	"regexp"
	"strconv"

	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"

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
	matched, err := regexp.MatchString("^\\d+$", strconv.FormatInt(req.ToUserID, 10)) //是否为纯数字
	if strconv.FormatInt(req.ToUserID, 10) == "" || matched == false {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apiVars.UserIdRuleError),
		}, nil
	} else if (req.ActionType != 1) && (req.ActionType != 2) { //是否有除1或2的数字
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apiVars.ActionTypeRuleError),
		}, nil
	}

	if req.Token == "" {
		return &types.RelationActionResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
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
