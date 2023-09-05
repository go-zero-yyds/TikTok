package social

import (
	"TikTok/apps/app/api/apiVars"
	"TikTok/apps/app/api/internal/middleware"
	"TikTok/apps/app/api/internal/svc"
	"TikTok/apps/app/api/internal/types"
	"TikTok/apps/social/rpc/social"
	"context"

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
	if len(req.Content) > 15000 { //内容是否为符合规范
		return &types.MessageActionResponse{
			RespStatus: types.RespStatus(apiVars.TextRuleError),
		}, nil
	}

	tokenID := l.ctx.Value(middleware.TokenIDKey).(int64)

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
