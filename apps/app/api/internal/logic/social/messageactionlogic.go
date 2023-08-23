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

type MessageActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageActionLogic) MessageAction(req *types.MessageActionRequest) (resp *types.MessageActionResponse, err error) {

	// 参数检查
	matched, err := regexp.MatchString("^\\d{19}$", strconv.FormatInt(req.ToUserID, 10)) //是否为19位纯数字
	if strconv.FormatInt(req.ToUserID, 10) == "" || matched == false {
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apiVars.UserIdRuleError),
		}, nil
	} else if req.ActionType != 1 { //是否有除1的数字
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apiVars.ActionTypeRuleError),
		}, nil
	} else if req.Content == "" { //内容是否为空
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apiVars.TextIsNull),
		}, nil
	}

	if req.Token == "" {
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apiVars.NotLogged),
		}, nil
	}

	tokenID, err := l.svcCtx.JwtAuth.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.SocialRPC.SendMessageAction(l.ctx, &social.MessageActionReq{
		UserId:     tokenID,
		ToUserId:   req.ToUserID,
		ActionType: req.ActionType,
		Content:    req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &types.MessageActionResponse{RespStatus: types.RespStatus(apiVars.Success)}, nil
}
